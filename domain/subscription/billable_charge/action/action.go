// Package action holds the billable charge write handlers: issue recovery
// documents (drawer + POST) and adjust an issued charge (drawer + POST).
//
// Every handler re-checks its permission first (layer 2, fail closed); the use
// cases remain the authoritative gate (strict verbs).
package action

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	recoverydocumentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/recovery_document"
	billablechargepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/billable_charge"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	rd "github.com/erniealice/centymo-golang/domain/revenue/recovery_document"
	bc "github.com/erniealice/centymo-golang/domain/subscription/billable_charge"
	"github.com/erniealice/centymo-golang/domain/subscription/billable_charge/form"
)

// TableID is the list table refreshed after a successful write.
const TableID = "billable-charges-table"

// Deps holds the action handler dependencies.
type Deps struct {
	Routes bc.Routes
	Labels bc.Labels
	// RecoveryLabels maps the issue use case's refusal codes
	// (`recovery_document.errors.<code>`).
	RecoveryLabels rd.Labels
	CommonLabels   pyeza.CommonLabels
	UseCases       *bc.UseCases

	// Now and NewKey are seams for tests; nil uses the wall clock / crypto rand.
	Now    func() time.Time
	NewKey func() string
}

func (d *Deps) now() time.Time {
	if d.Now != nil {
		return d.Now()
	}
	return time.Now()
}

func (d *Deps) newKey() string {
	if d.NewKey != nil {
		return d.NewKey()
	}
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		// A key that cannot be minted must not silently degrade to a constant.
		return fmt.Sprintf("k%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(b[:])
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

func invalid(deps *Deps) view.ViewResult {
	return view.HTMXError(deps.CommonLabels.Errors.InvalidFormData)
}

func refuseIssue(deps *Deps, err error) view.ViewResult {
	kind := rd.ErrorKind(err)
	if kind == rd.ErrKindPermissionDenied {
		return denied(deps)
	}
	if kind == rd.ErrKindUnknown {
		log.Printf("billable_charge issue: %v", err)
	}
	if m := deps.RecoveryLabels.ErrorMessage(kind); m != "" {
		return view.HTMXError(m)
	}
	return view.HTMXError(deps.CommonLabels.Errors.General)
}

func refuseAdjust(deps *Deps, err error) view.ViewResult {
	kind := bc.ErrorKind(err)
	if kind == bc.ErrKindPermissionDenied {
		return denied(deps)
	}
	if kind == bc.ErrKindUnknown {
		log.Printf("billable_charge adjust: %v", err)
	}
	if m := deps.Labels.ErrorMessage(kind); m != "" {
		return view.HTMXError(m)
	}
	return view.HTMXError(deps.CommonLabels.Errors.General)
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

func decimal(centavos int64) string {
	return fmt.Sprintf("%d.%02d", centavos/100, centavos%100)
}

// ---------------------------------------------------------------------------
// Issue recovery documents
// ---------------------------------------------------------------------------

// NewIssueAction handles GET (drawer listing the open charges) and POST
// (IssueRecoveryDocuments with the drawer's minted issuance key).
func NewIssueAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("recovery_document", "issue") {
			return denied(deps)
		}
		uc := deps.UseCases
		if uc == nil || uc.ListBillableCharges == nil || uc.ListDocumentSeries == nil || uc.IssueRecoveryDocuments == nil {
			return unavailable(deps)
		}
		if viewCtx.Request.Method == http.MethodGet {
			return issueDrawer(ctx, deps, viewCtx)
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return invalid(deps)
		}
		r := viewCtx.Request
		ids := dedupe(r.Form["charge_id"])
		seriesID := strings.TrimSpace(r.FormValue("document_series_id"))
		key := strings.TrimSpace(r.FormValue("issuance_key"))
		if len(ids) == 0 {
			return view.HTMXError(deps.RecoveryLabels.Errors.NothingSelected)
		}
		if seriesID == "" || key == "" {
			return invalid(deps)
		}
		req := &recoverydocumentpb.IssueRecoveryDocumentsRequest{
			BillableChargeIds: ids,
			DocumentSeriesId:  seriesID,
			IssuanceKey:       key,
		}
		if v := strings.TrimSpace(r.FormValue("issue_date")); v != "" {
			req.IssueDate = &v
		}
		if v := strings.TrimSpace(r.FormValue("due_date")); v != "" {
			req.DueDate = &v
		}
		if _, err := uc.IssueRecoveryDocuments(ctx, req); err != nil {
			return refuseIssue(deps, err)
		}
		return view.HTMXSuccess(TableID)
	})
}

func dedupe(in []string) []string {
	seen := map[string]bool{}
	out := make([]string, 0, len(in))
	for _, s := range in {
		s = strings.TrimSpace(s)
		if s != "" && !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	return out
}

func issueDrawer(ctx context.Context, deps *Deps, viewCtx *view.ViewContext) view.ViewResult {
	uc := deps.UseCases
	charges, err := uc.ListAllByStatus(ctx, billablechargepb.BillableChargeStatus_BILLABLE_CHARGE_STATUS_OPEN)
	if err != nil {
		log.Printf("billable_charge issue drawer: list open: %v", err)
		return view.HTMXError(deps.CommonLabels.Errors.General)
	}
	series, err := uc.ActiveRecoverySeries(ctx)
	if err != nil {
		log.Printf("billable_charge issue drawer: list series: %v", err)
		return view.HTMXError(deps.CommonLabels.Errors.General)
	}

	// Preselect the ?id= charges (row action); none given = every open charge.
	preselect := map[string]bool{}
	for _, id := range viewCtx.Request.URL.Query()["id"] {
		preselect[id] = true
	}
	clients := uc.ClientNames(ctx)
	subs := uc.SubscriptionNames(ctx)

	rows := make([]form.ChargeRow, 0, len(charges))
	clientSet := map[string]bool{}
	checked := 0
	for _, c := range charges {
		on := len(preselect) == 0 || preselect[c.GetId()]
		if on {
			checked++
			clientSet[c.GetClientId()] = true
		}
		rows = append(rows, form.ChargeRow{
			ID:           c.GetId(),
			Client:       name(clients, c.GetClientId()),
			Subscription: name(subs, c.GetSubscriptionId()),
			Period:       period(c),
			Amount:       types.FormatMoney(c.GetAmount(), c.GetCurrency()),
			Checked:      on,
		})
	}
	opts := make([]types.SelectOption, 0, len(series))
	for i, s := range series {
		label := s.GetCode()
		if n := s.GetName(); n != "" {
			label = s.GetCode() + " — " + n
		}
		opts = append(opts, types.SelectOption{Value: s.GetId(), Label: label, Selected: i == 0})
	}
	l := deps.Labels
	return view.OK("billable-charge-issue-drawer", &form.IssueData{
		FormAction:   deps.Routes.IssueURL,
		Labels:       l,
		Series:       opts,
		Charges:      rows,
		IssueDate:    deps.now().Format("2006-01-02"),
		IssuanceKey:  deps.newKey(),
		Summary:      bc.Format(l.Form.IssueSummary, strconv.Itoa(checked), strconv.Itoa(len(clientSet))),
		NoSeries:     len(opts) == 0,
		NoCharges:    len(rows) == 0,
		CommonLabels: nil, // injected by ViewAdapter
	})
}

func name(names map[string]string, id string) string {
	if n := names[id]; n != "" {
		return n
	}
	return id
}

func period(c *billablechargepb.BillableCharge) string {
	switch {
	case c.GetServiceFrom() != "" && c.GetServiceTo() != "":
		return c.GetServiceFrom() + " – " + c.GetServiceTo()
	default:
		return c.GetServiceFrom()
	}
}

// ---------------------------------------------------------------------------
// Adjust an issued original charge
// ---------------------------------------------------------------------------

// NewAdjustAction handles GET (drawer) and POST (AdjustBillableCharge).
func NewAdjustAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("billable_charge", "adjust") {
			return denied(deps)
		}
		uc := deps.UseCases
		if uc == nil || uc.ListBillableCharges == nil || uc.AdjustBillableCharge == nil {
			return unavailable(deps)
		}
		id := viewCtx.Request.PathValue("id")
		if id == "" {
			return invalid(deps)
		}
		if viewCtx.Request.Method == http.MethodGet {
			return adjustDrawer(ctx, deps, id)
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return invalid(deps)
		}
		amount, ok := parseAmount(viewCtx.Request.FormValue("new_amount"))
		reason := strings.TrimSpace(viewCtx.Request.FormValue("reason"))
		if !ok || reason == "" {
			return invalid(deps)
		}
		if _, err := uc.AdjustBillableCharge(ctx, &billablechargepb.AdjustBillableChargeRequest{
			BillableChargeId: id, NewAmount: amount, Reason: reason,
		}); err != nil {
			return refuseAdjust(deps, err)
		}
		return view.HTMXSuccess(TableID)
	})
}

func adjustDrawer(ctx context.Context, deps *Deps, id string) view.ViewResult {
	uc := deps.UseCases
	found, err := uc.ReadByID(ctx, id)
	if err != nil {
		log.Printf("billable_charge adjust drawer: read %s: %v", id, err)
		return view.HTMXError(deps.CommonLabels.Errors.General)
	}
	if found == nil {
		return view.HTMXError(deps.Labels.Errors.NotFound)
	}
	if found.GetStatus() != billablechargepb.BillableChargeStatus_BILLABLE_CHARGE_STATUS_ISSUED ||
		found.GetChargeKind() != billablechargepb.BillableChargeKind_BILLABLE_CHARGE_KIND_ORIGINAL {
		return view.HTMXError(deps.Labels.Errors.NotIssued)
	}
	clients := uc.ClientNames(ctx)
	return view.OK("billable-charge-adjust-drawer", &form.AdjustData{
		FormAction:    route.ResolveURL(deps.Routes.AdjustURL, "id", id),
		Labels:        deps.Labels,
		ID:            id,
		Client:        name(clients, found.GetClientId()),
		CurrentAmount: types.FormatMoney(found.GetAmount(), found.GetCurrency()),
		NewAmount:     decimal(found.GetAmount()),
		Currency:      found.GetCurrency(),
	})
}
