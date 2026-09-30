// Package action holds the cost allocation drawers and writes: the allocate
// drawer (participants + weights), the server-side preview, save draft, publish
// and the read-only view of a published allocation.
//
// Every handler re-checks its permission first (fail closed, layer 2); the use
// cases stay the authoritative gate (CheckStrict, term-boundary and claim rules).
// Preview is the persisted DRAFT: the split amounts come from the use case's own
// largest-remainder response (the split function lives in espyna internals and is
// not reachable from centymo), so the preview is exactly what publish will use.
package action

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"

	ab "github.com/erniealice/centymo-golang/domain/expenditure/allocation_batch"
	"github.com/erniealice/centymo-golang/domain/expenditure/allocation_batch/form"
	csc "github.com/erniealice/centymo-golang/domain/expenditure/cost_source_component"
	"github.com/erniealice/centymo-golang/domain/expenditure/cost_source_component/table"
	allocationbatchpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/allocation_batch"
	allocationsharepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/allocation_share"
	costsourcecomponentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/cost_source_component"
	subscriptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription"
)

// Deps holds the action handler dependencies.
type Deps struct {
	Routes          ab.Routes
	Labels          ab.Labels
	ComponentLabels csc.Labels
	CommonLabels    pyeza.CommonLabels
	UseCases        *ab.UseCases
}

func denied(deps *Deps) view.ViewResult {
	return view.HTMXError(deps.CommonLabels.Errors.PermissionDenied)
}

func unavailable(deps *Deps) view.ViewResult {
	return view.ViewResult{
		StatusCode: http.StatusServiceUnavailable,
		Headers:    map[string]string{"HX-Error-Message": deps.Labels.Errors.Unavailable},
	}
}

func refuse(deps *Deps, err error) view.ViewResult {
	kind := ab.ErrorKind(err)
	if kind == ab.ErrUnknown {
		log.Printf("allocation_batch action: %v", err)
	}
	if kind == ab.ErrPermissionDenied {
		return denied(deps)
	}
	return view.HTMXError(deps.Labels.ErrorMessage(kind))
}

func claimedRefusal(deps *Deps, c *costsourcecomponentpb.CostSourceComponent) (view.ViewResult, bool) {
	switch c.GetClaimKind() {
	case costsourcecomponentpb.SourceClaimKind_SOURCE_CLAIM_KIND_ALLOCATION:
		return view.HTMXError(deps.Labels.Errors.SourceClaimedByAllocation), true
	case costsourcecomponentpb.SourceClaimKind_SOURCE_CLAIM_KIND_RECOGNITION:
		return view.HTMXError(deps.Labels.Errors.SourceClaimedByRecognition), true
	}
	return view.ViewResult{}, false
}

func componentLabel(deps *Deps, c *costsourcecomponentpb.CostSourceComponent) string {
	kind := table.KindLabel(deps.ComponentLabels, c.GetComponentKind())
	if d := c.GetDescription(); d != "" {
		return kind + " — " + d
	}
	return kind
}

func shareKindLabel(l ab.Labels, k allocationsharepb.AllocationShareKind) string {
	switch k {
	case allocationsharepb.AllocationShareKind_ALLOCATION_SHARE_KIND_RECOVERABLE:
		return l.Enums.ShareKindRecoverable
	case allocationsharepb.AllocationShareKind_ALLOCATION_SHARE_KIND_OWN_USE:
		return l.Enums.ShareKindOwnUse
	case allocationsharepb.AllocationShareKind_ALLOCATION_SHARE_KIND_VACANCY:
		return l.Enums.ShareKindVacancy
	case allocationsharepb.AllocationShareKind_ALLOCATION_SHARE_KIND_COMMON_LOSS:
		return l.Enums.ShareKindCommonLoss
	}
	return ""
}

// percent renders numerator/denominator as a two-decimal percentage with integer math.
func percent(num, den int64) string {
	if den <= 0 {
		return "—"
	}
	p := num * 10000 / den
	return fmt.Sprintf("%d.%02d%%", p/100, p%100)
}

// fixedRows are the rows every allocation offers besides the participants.
func fixedRows(l ab.Labels) []form.Row {
	return []form.Row{
		{Kind: allocationsharepb.AllocationShareKind_ALLOCATION_SHARE_KIND_OWN_USE.String(), KindLabel: l.Enums.ShareKindOwnUse, Name: l.Form.OwnUseRowLabel, Slug: "own-use", Numerator: "0"},
		{Kind: allocationsharepb.AllocationShareKind_ALLOCATION_SHARE_KIND_VACANCY.String(), KindLabel: l.Enums.ShareKindVacancy, Name: l.Form.VacancyRowLabel, Slug: "vacancy", Numerator: "0"},
		{Kind: allocationsharepb.AllocationShareKind_ALLOCATION_SHARE_KIND_COMMON_LOSS.String(), KindLabel: l.Enums.ShareKindCommonLoss, Name: l.Form.CommonLossRowLabel, Slug: "common-loss", Numerator: "0"},
	}
}

func rowKey(kind, subscriptionID string) string { return kind + "|" + subscriptionID }

// NewAllocateAction handles GET (drawer) / POST (save draft) for one component ({id}).
func NewAllocateAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("allocation_batch", "create") && !perms.Can("allocation_batch", "update") {
			return denied(deps)
		}
		if deps.UseCases == nil {
			return unavailable(deps)
		}
		componentID := viewCtx.Request.PathValue("id")
		if viewCtx.Request.Method == http.MethodGet {
			return drawer(ctx, deps, perms, componentID)
		}
		if _, _, _, err := saveDraft(ctx, deps, perms, viewCtx.Request, componentID); err != nil {
			return refuseSave(deps, err)
		}
		return view.HTMXSuccess(table.TableID)
	})
}

// NewPreviewAction handles POST: persist the draft and render the split preview panel. The request
// comes from a plain button (not the form), so every outcome — including a refusal — is a 200
// fragment for the panel.
func NewPreviewAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		panel := func(preview *form.Preview, msg string) view.ViewResult {
			return view.OK("allocation-batch-preview", &form.PreviewPage{Preview: preview, Labels: deps.Labels, Error: msg})
		}
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("allocation_batch", "create") && !perms.Can("allocation_batch", "update") {
			return panel(nil, deps.CommonLabels.Errors.PermissionDenied)
		}
		if deps.UseCases == nil {
			return panel(nil, deps.Labels.Errors.Unavailable)
		}
		componentID := viewCtx.Request.PathValue("id")
		comp, _, shares, err := saveDraft(ctx, deps, perms, viewCtx.Request, componentID)
		if err != nil {
			return panel(nil, saveMessage(deps, err))
		}
		return panel(buildPreview(deps, comp, shares, labelsFromForm(viewCtx.Request)), "")
	})
}

// NewPublishAction handles POST: persist the draft, then publish it.
func NewPublishAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("allocation_batch", "publish") || (!perms.Can("allocation_batch", "create") && !perms.Can("allocation_batch", "update")) {
			return denied(deps)
		}
		if deps.UseCases == nil || deps.UseCases.PublishAllocationBatch == nil {
			return unavailable(deps)
		}
		componentID := viewCtx.Request.PathValue("id")
		_, batch, _, err := saveDraft(ctx, deps, perms, viewCtx.Request, componentID)
		if err != nil {
			return refuseSave(deps, err)
		}
		if _, err := deps.UseCases.PublishAllocationBatch(ctx, &allocationbatchpb.PublishAllocationBatchRequest{AllocationBatchId: batch.GetId()}); err != nil {
			return refuse(deps, err)
		}
		return view.HTMXSuccess(table.TableID)
	})
}

// NewViewAction handles GET: the read-only published allocation of one component ({id}).
func NewViewAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("allocation_batch", "read") {
			return denied(deps)
		}
		if deps.UseCases == nil {
			return unavailable(deps)
		}
		componentID := viewCtx.Request.PathValue("id")
		comp, err := deps.UseCases.ReadComponent(ctx, componentID)
		if err != nil {
			log.Printf("allocation_batch view: read component %s: %v", componentID, err)
			return view.HTMXError(deps.Labels.Errors.NotFound)
		}
		batches, err := deps.UseCases.BatchesOf(ctx, componentID)
		if err != nil {
			return refuse(deps, err)
		}
		var published *allocationbatchpb.AllocationBatch
		for _, b := range batches {
			if b.GetStatus() == allocationbatchpb.AllocationBatchStatus_ALLOCATION_BATCH_STATUS_PUBLISHED {
				published = b
				break
			}
		}
		if published == nil {
			return view.HTMXError(deps.Labels.Errors.NotFound)
		}
		shares, err := deps.UseCases.SharesOf(ctx, published.GetId())
		if err != nil {
			return refuse(deps, err)
		}
		names := map[string]string{}
		for _, s := range shares {
			id := s.GetSubscriptionId()
			if id == "" || names[id] != "" || deps.UseCases.GetSubscriptionItemPageData == nil {
				continue
			}
			if resp, err := deps.UseCases.GetSubscriptionItemPageData(ctx, &subscriptionpb.GetSubscriptionItemPageDataRequest{SubscriptionId: id}); err == nil {
				names[id] = resp.GetSubscription().GetName()
			}
		}
		return view.OK("allocation-batch-view", &form.ViewData{
			ComponentLabel: componentLabel(deps, comp),
			ServicePeriod:  table.ServicePeriod(deps.ComponentLabels, comp.GetServiceFrom(), comp.GetServiceTo()),
			Revision:       strings.Replace(deps.Labels.Detail.RevisionLabel, "{0}", strconv.Itoa(int(published.GetRevision())), 1),
			Status:         deps.Labels.Enums.StatusPublished,
			Preview:        buildPreview(deps, comp, shares, names),
			Labels:         deps.Labels, CommonLabels: deps.CommonLabels,
		})
	})
}

func drawer(ctx context.Context, deps *Deps, perms interface{ Can(string, string) bool }, componentID string) view.ViewResult {
	comp, err := deps.UseCases.ReadComponent(ctx, componentID)
	if err != nil {
		log.Printf("allocation_batch drawer: read component %s: %v", componentID, err)
		return view.HTMXError(deps.Labels.Errors.NotFound)
	}
	if res, claimed := claimedRefusal(deps, comp); claimed {
		return res
	}
	participants, err := deps.UseCases.Participants(ctx, comp.GetServiceFrom(), comp.GetServiceTo())
	if err != nil {
		return refuse(deps, err)
	}
	batches, err := deps.UseCases.BatchesOf(ctx, componentID)
	if err != nil {
		return refuse(deps, err)
	}
	var draft *allocationbatchpb.AllocationBatch
	for _, b := range batches { // newest revision first
		if b.GetStatus() == allocationbatchpb.AllocationBatchStatus_ALLOCATION_BATCH_STATUS_DRAFT {
			draft = b
			break
		}
	}
	saved := map[string]*allocationsharepb.AllocationShare{}
	var draftShares []*allocationsharepb.AllocationShare
	denominator := ""
	if draft != nil {
		draftShares, err = deps.UseCases.SharesOf(ctx, draft.GetId())
		if err != nil {
			return refuse(deps, err)
		}
		for _, s := range draftShares {
			saved[rowKey(s.GetShareKind().String(), s.GetSubscriptionId())] = s
			if denominator == "" && s.GetBasisDenominator() > 0 {
				denominator = strconv.FormatInt(s.GetBasisDenominator(), 10)
			}
		}
	}

	var rows []form.Row
	for _, p := range participants {
		kind := allocationsharepb.AllocationShareKind_ALLOCATION_SHARE_KIND_RECOVERABLE
		name := p.SubscriptionName
		if p.ClientName != "" {
			name += " — " + p.ClientName
		}
		row := form.Row{Kind: kind.String(), KindLabel: deps.Labels.Enums.ShareKindRecoverable, SubscriptionID: p.SubscriptionID, ClientID: p.ClientID, Name: name, Slug: p.SubscriptionID, Numerator: "0"}
		if s := saved[rowKey(kind.String(), p.SubscriptionID)]; s != nil {
			row.Numerator = strconv.FormatInt(s.GetBasisNumerator(), 10)
		}
		rows = append(rows, row)
	}
	for _, f := range fixedRows(deps.Labels) {
		if s := saved[rowKey(f.Kind, "")]; s != nil {
			f.Numerator = strconv.FormatInt(s.GetBasisNumerator(), 10)
		}
		rows = append(rows, f)
	}
	if len(participants) == 0 {
		rows = nil // nobody can share this cost: show the empty state, no inputs
	}

	data := &form.Data{
		FormAction:  route.ResolveURL(deps.Routes.AllocateURL, "id", componentID),
		PreviewURL:  route.ResolveURL(deps.Routes.PreviewURL, "id", componentID),
		PublishURL:  route.ResolveURL(deps.Routes.PublishURL, "id", componentID),
		ComponentID: componentID, ComponentLabel: componentLabel(deps, comp),
		ServicePeriod: table.ServicePeriod(deps.ComponentLabels, comp.GetServiceFrom(), comp.GetServiceTo()),
		SourceAmount:  csc.FormatAmount(comp.GetAmount()), Currency: comp.GetCurrency(),
		Denominator: denominator, Rows: rows,
		CanPublish: perms.Can("allocation_batch", "publish"), PublishTooltip: fmt.Sprintf(deps.CommonLabels.Errors.MissingPermission, "allocation_batch:publish"),
		Labels: deps.Labels, CommonLabels: deps.CommonLabels,
	}
	if draft != nil {
		data.HasDraft = true
		data.Revision = strings.Replace(deps.Labels.Detail.RevisionLabel, "{0}", strconv.Itoa(int(draft.GetRevision())), 1)
		names := map[string]string{}
		for _, r := range rows {
			names[rowKey(r.Kind, r.SubscriptionID)] = r.Name
		}
		data.Preview = buildPreview(deps, comp, draftShares, names)
	}
	return view.OK("allocation-batch-drawer-form", data)
}

// errNoBatch marks a create that returned no batch (an adapter contract breach; logged, generic message).
var errNoBatch = errors.New("no batch returned")

// saveErr keeps a coded use-case error distinct from a malformed form.
type saveErr struct {
	form        bool
	denied      bool
	unavailable bool
	err         error
}

func (e *saveErr) Error() string {
	if e.err != nil {
		return e.err.Error()
	}
	return "invalid form"
}

func refuseSave(deps *Deps, err error) view.ViewResult {
	return view.HTMXError(saveMessage(deps, err))
}

// saveMessage maps a save failure to its user-facing label.
func saveMessage(deps *Deps, err error) string {
	se, ok := err.(*saveErr)
	switch {
	case ok && se.form:
		return deps.Labels.Errors.FormInvalid
	case ok && se.denied:
		return deps.CommonLabels.Errors.PermissionDenied
	case ok && se.unavailable:
		return deps.Labels.Errors.Unavailable
	case ok:
		err = se.err
	}
	kind := ab.ErrorKind(err)
	if kind == ab.ErrUnknown {
		log.Printf("allocation_batch action: %v", err)
	}
	if kind == ab.ErrPermissionDenied {
		return deps.CommonLabels.Errors.PermissionDenied
	}
	return deps.Labels.ErrorMessage(kind)
}

// labelsFromForm maps "kind|subscription" -> the row label the drawer posted.
func labelsFromForm(r *http.Request) map[string]string {
	out := map[string]string{}
	kinds, subs, labels := r.Form["share_kind"], r.Form["subscription_id"], r.Form["row_label"]
	for i := range kinds {
		if i < len(subs) && i < len(labels) {
			out[rowKey(kinds[i], subs[i])] = labels[i]
		}
	}
	return out
}

// sharesFromForm parses the posted weight rows. A blank total weight defaults to
// the sum of the weights (the use case requires the weights to add up to it).
func sharesFromForm(r *http.Request) ([]*allocationsharepb.AllocationShare, bool) {
	if err := r.ParseForm(); err != nil {
		return nil, false
	}
	kinds, subs, clients, nums := r.Form["share_kind"], r.Form["subscription_id"], r.Form["client_id"], r.Form["numerator"]
	if len(kinds) == 0 || len(kinds) != len(subs) || len(kinds) != len(clients) || len(kinds) != len(nums) {
		return nil, false
	}
	var sum int64
	shares := make([]*allocationsharepb.AllocationShare, 0, len(kinds))
	for i := range kinds {
		kv, ok := allocationsharepb.AllocationShareKind_value[kinds[i]]
		if !ok || kv == 0 {
			return nil, false
		}
		n := int64(0)
		if t := strings.TrimSpace(nums[i]); t != "" {
			v, err := strconv.ParseInt(t, 10, 64)
			if err != nil || v < 0 {
				return nil, false
			}
			n = v
		}
		sum += n
		s := &allocationsharepb.AllocationShare{ShareKind: allocationsharepb.AllocationShareKind(kv), BasisNumerator: n}
		if v := strings.TrimSpace(subs[i]); v != "" {
			s.SubscriptionId = &v
		}
		if v := strings.TrimSpace(clients[i]); v != "" {
			s.ClientId = &v
		}
		shares = append(shares, s)
	}
	den := sum
	if t := strings.TrimSpace(r.FormValue("denominator")); t != "" {
		v, err := strconv.ParseInt(t, 10, 64)
		if err != nil {
			return nil, false
		}
		den = v
	}
	for _, s := range shares {
		s.BasisDenominator = den
	}
	return shares, true
}

// saveDraft persists the posted weights as the component's DRAFT allocation
// (Create when none exists, UpdateAllocationBatchShares otherwise) and returns
// the batch and its stored shares (with the use case's exact split amounts).
func saveDraft(ctx context.Context, deps *Deps, perms interface{ Can(string, string) bool }, r *http.Request, componentID string) (*costsourcecomponentpb.CostSourceComponent, *allocationbatchpb.AllocationBatch, []*allocationsharepb.AllocationShare, error) {
	uc := deps.UseCases
	if uc == nil || uc.CreateAllocationBatch == nil || uc.UpdateAllocationBatchShares == nil {
		return nil, nil, nil, &saveErr{unavailable: true}
	}
	shares, ok := sharesFromForm(r)
	if !ok {
		return nil, nil, nil, &saveErr{form: true}
	}
	comp, err := uc.ReadComponent(ctx, componentID)
	if err != nil {
		return nil, nil, nil, &saveErr{err: err}
	}
	batches, err := uc.BatchesOf(ctx, componentID)
	if err != nil {
		return nil, nil, nil, &saveErr{err: err}
	}
	var draft *allocationbatchpb.AllocationBatch
	for _, b := range batches {
		if b.GetStatus() == allocationbatchpb.AllocationBatchStatus_ALLOCATION_BATCH_STATUS_DRAFT {
			draft = b
			break
		}
	}
	if draft == nil {
		if !perms.Can("allocation_batch", "create") {
			return nil, nil, nil, &saveErr{denied: true}
		}
		resp, err := uc.CreateAllocationBatch(ctx, &allocationbatchpb.CreateAllocationBatchRequest{
			Data:   &allocationbatchpb.AllocationBatch{CostSourceComponentId: componentID},
			Shares: shares,
		})
		if err != nil {
			return nil, nil, nil, &saveErr{err: err}
		}
		if len(resp.GetData()) == 0 {
			return nil, nil, nil, &saveErr{err: errNoBatch}
		}
		return comp, resp.GetData()[0], resp.GetShares(), nil
	}
	if !perms.Can("allocation_batch", "update") {
		return nil, nil, nil, &saveErr{denied: true}
	}
	resp, err := uc.UpdateAllocationBatchShares(ctx, &allocationbatchpb.UpdateAllocationBatchSharesRequest{AllocationBatchId: draft.GetId(), Shares: shares})
	if err != nil {
		return nil, nil, nil, &saveErr{err: err}
	}
	batch := resp.GetData()
	if batch == nil {
		batch = draft
	}
	return comp, batch, resp.GetShares(), nil
}

// buildPreview renders stored shares (exact amounts) as the preview panel. names is
// keyed "kind|subscription" (or by subscription id for the read-only view).
func buildPreview(deps *Deps, comp *costsourcecomponentpb.CostSourceComponent, shares []*allocationsharepb.AllocationShare, names map[string]string) *form.Preview {
	p := &form.Preview{Currency: comp.GetCurrency(), Source: csc.FormatAmount(comp.GetAmount())}
	var total int64
	for _, s := range shares {
		name := names[rowKey(s.GetShareKind().String(), s.GetSubscriptionId())]
		if name == "" {
			name = names[s.GetSubscriptionId()]
		}
		if name == "" && s.GetSubscriptionId() == "" {
			name = shareKindLabel(deps.Labels, s.GetShareKind())
		}
		total += s.GetAmount()
		p.Lines = append(p.Lines, form.PreviewLine{
			KindLabel: shareKindLabel(deps.Labels, s.GetShareKind()),
			Name:      name,
			Weight:    strconv.FormatInt(s.GetBasisNumerator(), 10),
			Percent:   percent(s.GetBasisNumerator(), s.GetBasisDenominator()),
			Amount:    csc.FormatAmount(s.GetAmount()),
		})
	}
	p.Total = csc.FormatAmount(total)
	p.Balanced = total == comp.GetAmount()
	return p
}
