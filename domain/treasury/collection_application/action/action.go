// Package action holds the receive-and-apply handlers: the drawer (GET/POST),
// the read-only application preview, and the application reversal.
//
// Every handler re-checks its permission first (layer 2, fail closed); the use
// cases remain the authoritative gate (strict verbs).
package action

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	collectionapplicationpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/collection_application"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	ra "github.com/erniealice/centymo-golang/domain/treasury/collection_application"
)

// CollectionsTableID is the collection list table refreshed when no detail
// URL is wired.
const CollectionsTableID = "collections-table"

// Deps holds the handler dependencies.
type Deps struct {
	Routes       ra.Routes
	Labels       ra.Labels
	CommonLabels pyeza.CommonLabels
	UseCases     *ra.UseCases

	// SearchClientURL is the async client search route (revenue.search_client);
	// empty falls back to a select of every client (paged through).
	SearchClientURL string

	// CollectionDetailURL is the collection detail route template ({id}); the
	// success redirect target. Empty refreshes the collection list instead.
	CollectionDetailURL string
	// DefaultCurrency preselects the currency select.
	DefaultCurrency string

	// Now is a test seam; nil uses the wall clock.
	Now func() time.Time
}

func (d *Deps) now() time.Time {
	if d.Now != nil {
		return d.Now()
	}
	return time.Now()
}

// ReceiveApplyData is the template data of the receive-and-apply drawer.
type ReceiveApplyData struct {
	FormAction  string
	WorkspaceID string // injected by the ViewAdapter (action_workspace_guard)
	PreviewURL  string

	Labels        ra.Labels
	SearchURL     string // async client search; empty renders the ClientOptions select
	ClientID      string // preselected client (search mode)
	ClientLabel   string
	ClientOptions []types.SelectOption
	MethodOptions []types.SelectOption
	Amount        string
	Currency      string
	PaymentDate   string
	Reference     string

	CommonLabels any
}

// PlanRow is one line of the application plan.
type PlanRow struct {
	Number     string
	TargetKind string
	DueDate    string
	Balance    string
	Apply      string
	Rank       int32
}

// PreviewData is the template data of the application preview fragment.
type PreviewData struct {
	Labels      ra.Labels
	Rows        []PlanRow
	Applied     string
	Unapplied   string
	NoOpenItems bool
	Error       string
}

func denied(deps *Deps) view.ViewResult {
	return view.HTMXError(deps.CommonLabels.Errors.PermissionDenied)
}

func unavailable(deps *Deps) view.ViewResult {
	return view.ViewResult{
		StatusCode: http.StatusServiceUnavailable,
		Headers:    map[string]string{"HX-Error-Message": deps.CommonLabels.Errors.General},
	}
}

func message(deps *Deps, err error) string {
	kind := ra.ErrorKind(err)
	if kind == ra.ErrKindPermissionDenied {
		return deps.CommonLabels.Errors.PermissionDenied
	}
	if kind == ra.ErrKindUnknown {
		log.Printf("collection_application: %v", err)
	}
	if m := deps.Labels.ErrorMessage(kind); m != "" {
		return m
	}
	return deps.CommonLabels.Errors.General
}

// parseAmount parses a decimal major-unit form value into centavos through the
// shared pyeza parser (rejects malformed and negative input); ok is false unless
// the amount is > 0.
func parseAmount(s string) (int64, bool) {
	c, err := types.ParseCentavos(strings.ReplaceAll(s, ",", ""))
	if err != nil || c <= 0 {
		return 0, false
	}
	return c, true
}

type receiveFields struct {
	ClientID, Currency, Date, Method, Reference string
	Amount                                      int64
	AmountOK                                    bool
}

func readFields(r *http.Request) receiveFields {
	f := receiveFields{
		ClientID:  strings.TrimSpace(r.FormValue("client_id")),
		Currency:  strings.TrimSpace(r.FormValue("currency")),
		Date:      strings.TrimSpace(r.FormValue("payment_date")),
		Method:    strings.TrimSpace(r.FormValue("collection_method_id")),
		Reference: strings.TrimSpace(r.FormValue("reference_number")),
	}
	f.Amount, f.AmountOK = parseAmount(r.FormValue("amount"))
	return f
}

// NewReceiveApplyAction handles GET (drawer) and POST (ReceiveAndApplyCollection).
func NewReceiveApplyAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("collection_application", "create") {
			return denied(deps)
		}
		uc := deps.UseCases
		if uc == nil || uc.ReceiveAndApplyCollection == nil || uc.PreviewCollectionApplication == nil {
			return unavailable(deps)
		}
		if viewCtx.Request.Method == http.MethodGet {
			clientID := strings.TrimSpace(viewCtx.Request.URL.Query().Get("client_id"))
			cur := deps.DefaultCurrency
			if cur == "" {
				cur = "PHP"
			}
			data := &ReceiveApplyData{
				FormAction:    deps.Routes.ReceiveApplyURL,
				PreviewURL:    deps.Routes.PreviewURL,
				Labels:        deps.Labels,
				MethodOptions: uc.MethodOptions(ctx),
				Currency:      cur,
				PaymentDate:   deps.now().Format("2006-01-02"),
			}
			if deps.SearchClientURL != "" {
				data.SearchURL = deps.SearchClientURL
				if data.ClientLabel = uc.SelectedClient(ctx, clientID); data.ClientLabel != "" {
					data.ClientID = clientID
				}
			} else {
				data.ClientOptions = uc.ClientOptions(ctx, clientID)
			}
			return view.OK("receive-apply-drawer", data)
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.HTMXError(deps.CommonLabels.Errors.InvalidFormData)
		}
		f := readFields(viewCtx.Request)
		if f.ClientID == "" {
			return view.HTMXError(deps.Labels.Errors.ClientRequired)
		}
		if !f.AmountOK {
			return view.HTMXError(deps.Labels.Errors.AmountInvalid)
		}
		req := &collectionapplicationpb.ReceiveAndApplyCollectionRequest{
			ClientId: f.ClientID, Amount: f.Amount, Currency: f.Currency,
		}
		if f.Date != "" {
			req.PaymentDate = &f.Date
		}
		if f.Method != "" {
			req.CollectionMethodId = &f.Method
		}
		if f.Reference != "" {
			req.ReferenceNumber = &f.Reference
		}
		resp, err := uc.ReceiveAndApplyCollection(ctx, req)
		if err != nil {
			return view.HTMXError(message(deps, err))
		}
		if id := resp.GetCollection().GetId(); id != "" && deps.CollectionDetailURL != "" {
			return view.ViewResult{
				StatusCode: http.StatusOK,
				Headers: map[string]string{
					"HX-Trigger":  `{"formSuccess":true}`,
					"HX-Redirect": route.ResolveURL(deps.CollectionDetailURL, "id", id),
				},
			}
		}
		return view.HTMXSuccess(CollectionsTableID)
	})
}

// NewPreviewAction renders the N9 application plan for the drawer's current
// values (GET; read-only PreviewCollectionApplication).
func NewPreviewAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("collection_application", "create") {
			return denied(deps)
		}
		uc := deps.UseCases
		if uc == nil || uc.PreviewCollectionApplication == nil {
			return unavailable(deps)
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.HTMXError(deps.CommonLabels.Errors.InvalidFormData)
		}
		f := readFields(viewCtx.Request)
		data := &PreviewData{Labels: deps.Labels}
		switch {
		case f.ClientID == "":
			data.Error = deps.Labels.Errors.ClientRequired
			return view.OK("receive-apply-preview", data)
		case !f.AmountOK:
			data.Error = deps.Labels.Errors.AmountInvalid
			return view.OK("receive-apply-preview", data)
		}
		req := &collectionapplicationpb.PreviewCollectionApplicationRequest{
			ClientId: f.ClientID, Amount: f.Amount, Currency: f.Currency,
		}
		resp, err := uc.PreviewCollectionApplication(ctx, req)
		if err != nil {
			data.Error = message(deps, err)
			return view.OK("receive-apply-preview", data)
		}
		plan := resp.GetPlan()
		cur := plan.GetCurrency()
		if cur == "" {
			cur = f.Currency
		}
		for _, a := range plan.GetAllocations() {
			kind := deps.Labels.Enums.TargetKindRevenue
			if a.GetTargetKind() == collectionapplicationpb.ApplicationTargetKind_APPLICATION_TARGET_KIND_RECOVERY_DOCUMENT {
				kind = deps.Labels.Enums.TargetKindRecoveryDocument
			}
			data.Rows = append(data.Rows, PlanRow{
				Number: a.GetNumber(), TargetKind: kind, DueDate: a.GetDueDate(),
				Balance: types.FormatMoney(a.GetBalance(), cur), Apply: types.FormatMoney(a.GetApply(), cur), Rank: a.GetRank(),
			})
		}
		data.NoOpenItems = len(data.Rows) == 0
		data.Applied = types.FormatMoney(plan.GetApplied(), cur)
		data.Unapplied = types.FormatMoney(plan.GetUnapplied(), cur)
		return view.OK("receive-apply-preview", data)
	})
}

// NewReverseAction handles POST (ReverseCollectionApplication) from a row
// action; success refreshes the applications table.
func NewReverseAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("collection_application", "reverse") {
			return denied(deps)
		}
		uc := deps.UseCases
		if uc == nil || uc.ReverseCollectionApplication == nil {
			return unavailable(deps)
		}
		id := viewCtx.Request.PathValue("id")
		if id == "" {
			return view.HTMXError(deps.CommonLabels.Errors.IDRequired)
		}
		if _, err := uc.ReverseCollectionApplication(ctx, &collectionapplicationpb.ReverseCollectionApplicationRequest{
			CollectionApplicationId: id,
		}); err != nil {
			return view.HTMXError(message(deps, err))
		}
		return view.HTMXSuccess(ra.ApplicationsTableID)
	})
}
