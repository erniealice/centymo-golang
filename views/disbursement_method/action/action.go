package action

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	centymo "github.com/erniealice/centymo-golang"
	"github.com/erniealice/centymo-golang/views/disbursement_method/form"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	dmpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/disbursement_method"
)

// Deps holds all dependencies for the disbursement method action handlers.
// CRUD closures are nil-safe (espyna disbursement_method use cases pending).
type Deps struct {
	Routes       centymo.DisbursementMethodRoutes
	Labels       centymo.DisbursementMethodLabels
	CommonLabels pyeza.CommonLabels

	CreateDisbursementMethod func(ctx context.Context, req *dmpb.CreateDisbursementMethodRequest) (*dmpb.CreateDisbursementMethodResponse, error)
	ReadDisbursementMethod   func(ctx context.Context, req *dmpb.ReadDisbursementMethodRequest) (*dmpb.ReadDisbursementMethodResponse, error)
	UpdateDisbursementMethod func(ctx context.Context, req *dmpb.UpdateDisbursementMethodRequest) (*dmpb.UpdateDisbursementMethodResponse, error)
	DeleteDisbursementMethod func(ctx context.Context, req *dmpb.DeleteDisbursementMethodRequest) (*dmpb.DeleteDisbursementMethodResponse, error)
}

const tableID = "disbursement-methods-table"

// NewAddAction handles GET+POST /action/disbursement-method/add.
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("disbursement_method", "create") {
			return centymo.HTMXError(fmt.Sprintf(deps.CommonLabels.Errors.MissingPermission, "disbursement_method:create"))
		}
		if viewCtx.Request.Method == http.MethodGet {
			fd := buildEmptyFormData(deps)
			fd.FormAction = deps.Routes.AddURL
			fd.Lifecycle = "DISBURSEMENT_METHOD_LIFECYCLE_DRAFT"
			fd.Category = "DISBURSEMENT_METHOD_CATEGORY_STANDARD"
			fd.Fragment = buildFragmentData(deps, fd.Category)
			return view.OK("disbursement-method-drawer-form", fd)
		}

		req, err := parseDisbursementMethod(viewCtx.Request, "")
		if err != nil {
			return view.Error(err)
		}
		if deps.CreateDisbursementMethod == nil {
			return centymo.HTMXError("Payment method create is not wired yet (espyna use cases pending).")
		}
		if _, err := deps.CreateDisbursementMethod(ctx, &dmpb.CreateDisbursementMethodRequest{Data: req}); err != nil {
			log.Printf("CreateDisbursementMethod: %v", err)
			return view.Error(fmt.Errorf("failed to create disbursement method: %w", err))
		}
		return centymo.HTMXSuccess(tableID)
	})
}

// NewEditAction handles GET+POST /action/disbursement-method/edit/{id}.
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("disbursement_method", "update") {
			return centymo.HTMXError(fmt.Sprintf(deps.CommonLabels.Errors.MissingPermission, "disbursement_method:update"))
		}
		id := viewCtx.Request.PathValue("id")
		if id == "" {
			return view.Error(fmt.Errorf("missing id"))
		}

		if viewCtx.Request.Method == http.MethodGet {
			if deps.ReadDisbursementMethod == nil {
				return centymo.HTMXError("Payment method read is not wired yet (espyna use cases pending).")
			}
			resp, err := deps.ReadDisbursementMethod(ctx, &dmpb.ReadDisbursementMethodRequest{
				Data: &dmpb.DisbursementMethod{Id: id},
			})
			if err != nil {
				return view.Error(fmt.Errorf("failed to read disbursement method: %w", err))
			}
			data := resp.GetData()
			if len(data) == 0 {
				return view.Error(fmt.Errorf("disbursement method not found"))
			}
			fd := buildFormDataFromMethod(deps, data[0])
			fd.FormAction = route.ResolveURL(deps.Routes.EditURL, "id", id)
			fd.IsEdit = true
			fd.ID = id
			return view.OK("disbursement-method-drawer-form", fd)
		}

		req, err := parseDisbursementMethod(viewCtx.Request, id)
		if err != nil {
			return view.Error(err)
		}
		if deps.UpdateDisbursementMethod == nil {
			return centymo.HTMXError("Payment method update is not wired yet (espyna use cases pending).")
		}
		if _, err := deps.UpdateDisbursementMethod(ctx, &dmpb.UpdateDisbursementMethodRequest{Data: req}); err != nil {
			log.Printf("UpdateDisbursementMethod %s: %v", id, err)
			return view.Error(fmt.Errorf("failed to update disbursement method: %w", err))
		}
		return centymo.HTMXSuccess(tableID)
	})
}

// NewFragmentAction handles GET /action/disbursement-method/fragment?cat=X — the
// §A template-slot HTMX swap (buying side).
func NewFragmentAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("disbursement_method", "create") && !perms.Can("disbursement_method", "update") {
			return centymo.HTMXError(fmt.Sprintf(deps.CommonLabels.Errors.MissingPermission, "disbursement_method:create"))
		}
		cat := viewCtx.Request.URL.Query().Get("cat")
		if cat == "" {
			cat = viewCtx.Request.URL.Query().Get("category")
		}
		return view.OK("disbursement-method-kind-fragment", buildFragmentData(deps, cat))
	})
}

// NewDeleteAction handles POST /action/disbursement-method/delete.
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("disbursement_method", "update") {
			return centymo.HTMXError(fmt.Sprintf(deps.CommonLabels.Errors.MissingPermission, "disbursement_method:update"))
		}
		if viewCtx.Request.Method != http.MethodPost {
			return view.Error(fmt.Errorf("method not allowed"))
		}
		id := viewCtx.Request.FormValue("id")
		if id == "" {
			return view.Error(fmt.Errorf("missing id"))
		}
		if deps.DeleteDisbursementMethod == nil {
			return centymo.HTMXError("Payment method delete is not wired yet (espyna use cases pending).")
		}
		if _, err := deps.DeleteDisbursementMethod(ctx, &dmpb.DeleteDisbursementMethodRequest{
			Data: &dmpb.DisbursementMethod{Id: id},
		}); err != nil {
			log.Printf("DeleteDisbursementMethod %s: %v", id, err)
			return view.Error(fmt.Errorf("failed to delete disbursement method: %w", err))
		}
		return centymo.HTMXSuccess(tableID)
	})
}

// --- form helpers ------------------------------------------------------------

func buildEmptyFormData(deps *Deps) *form.Data {
	l := deps.Labels
	fd := &form.Data{
		Labels:       l.Form,
		FragLabels:   l.Fragment,
		CommonLabels: deps.CommonLabels,
		FragmentURL:  deps.Routes.FragmentURL,
	}
	// Buying-side categories: bank-account (STANDARD) / check / advance.
	// CHECK has no proto template_details variant yet — see action note.
	fd.CategoryOptions = []types.SelectOption{
		{Value: "DISBURSEMENT_METHOD_CATEGORY_STANDARD", Label: l.Form.CategoryStandard},
		{Value: "DISBURSEMENT_METHOD_CATEGORY_ADVANCE", Label: l.Form.CategoryAdvance},
	}
	fd.PostingKindOptions = []types.SelectOption{
		{Value: "DISBURSEMENT_METHOD_POSTING_KIND_CASH", Label: "Cash"},
		{Value: "DISBURSEMENT_METHOD_POSTING_KIND_ADVANCE_DRAWDOWN", Label: "Advance Drawdown"},
		{Value: "DISBURSEMENT_METHOD_POSTING_KIND_CLAIM_AR", Label: "Claim / AP"},
		{Value: "DISBURSEMENT_METHOD_POSTING_KIND_DEFERRED_RECEIVABLE", Label: "Deferred Payable"},
	}
	fd.TaxEffectOptions = []types.SelectOption{
		{Value: "DISBURSEMENT_METHOD_TAX_EFFECT_KIND_NONE", Label: "None"},
		{Value: "DISBURSEMENT_METHOD_TAX_EFFECT_KIND_INCLUSIVE", Label: "Tax Inclusive"},
		{Value: "DISBURSEMENT_METHOD_TAX_EFFECT_KIND_EXCLUSIVE", Label: "Tax Exclusive"},
	}
	fd.LifecycleOptions = []types.SelectOption{
		{Value: "DISBURSEMENT_METHOD_LIFECYCLE_DRAFT", Label: "Draft"},
		{Value: "DISBURSEMENT_METHOD_LIFECYCLE_ACTIVE", Label: "Active"},
		{Value: "DISBURSEMENT_METHOD_LIFECYCLE_CLOSED", Label: "Closed"},
		{Value: "DISBURSEMENT_METHOD_LIFECYCLE_ARCHIVED", Label: "Archived"},
	}
	fd.SourceOptions = []types.SelectOption{
		{Value: "DISBURSEMENT_METHOD_SOURCE_WORKSPACE", Label: "Workspace"},
		{Value: "DISBURSEMENT_METHOD_SOURCE_SYSTEM", Label: "System"},
		{Value: "DISBURSEMENT_METHOD_SOURCE_VENDOR_TEMPLATE", Label: "Vendor Template"},
	}
	fd.VersionStatusOptions = []types.SelectOption{
		{Value: "DISBURSEMENT_METHOD_VERSION_STATUS_DRAFT", Label: "Draft"},
		{Value: "DISBURSEMENT_METHOD_VERSION_STATUS_PUBLISHED", Label: "Published"},
		{Value: "DISBURSEMENT_METHOD_VERSION_STATUS_SUPERSEDED", Label: "Superseded"},
	}
	return fd
}

func buildFormDataFromMethod(deps *Deps, m *dmpb.DisbursementMethod) *form.Data {
	fd := buildEmptyFormData(deps)
	fd.Name = m.GetName()
	fd.Category = m.GetCategory().String()
	fd.PostingKind = m.GetPostingKind().String()
	fd.TaxEffectKind = m.GetTaxEffectKind().String()
	fd.BalanceAccount = m.GetBalanceAccountId()
	fd.TargetAccount = m.GetTargetAccountId()
	fd.Lifecycle = m.GetLifecycle().String()
	fd.Source = m.GetSource().String()
	fd.TemplateCode = m.GetTemplateCode()
	fd.Revision = strconv.FormatInt(int64(m.GetRevision()), 10)
	fd.VersionStatus = m.GetVersionStatus().String()
	fd.Supersedes = m.GetSupersedesDisbursementMethodId()
	fd.Fragment = buildFragmentDataFromMethod(deps, m)
	return fd
}

func buildFragmentData(deps *Deps, category string) form.FragmentData {
	return form.FragmentData{
		Category:   category,
		FragLabels: deps.Labels.Fragment,
	}
}

func buildFragmentDataFromMethod(deps *Deps, m *dmpb.DisbursementMethod) form.FragmentData {
	fd := buildFragmentData(deps, m.GetCategory().String())
	if ap := m.GetAdvanceProgram(); ap != nil {
		fd.AdvanceKind = ap.GetAdvanceKind().String()
		fd.DefaultBalanceAcct = ap.GetDefaultBalanceAccountId()
		fd.DefaultTargetAcct = ap.GetDefaultTargetAccountId()
		if ap.GetDefaultPeriodCount() > 0 {
			fd.DefaultPeriodCount = strconv.FormatInt(int64(ap.GetDefaultPeriodCount()), 10)
		}
		fd.DefaultPeriodUnit = ap.GetDefaultPeriodUnit()
	} else if ba := m.GetBankAccount(); ba != nil {
		fd.BankName = ba.GetBankName()
	}
	return fd
}

// parseDisbursementMethod builds a proto from the submitted form.
func parseDisbursementMethod(r *http.Request, id string) (*dmpb.DisbursementMethod, error) {
	if err := r.ParseForm(); err != nil {
		return nil, fmt.Errorf("parse form: %w", err)
	}
	m := &dmpb.DisbursementMethod{
		Id:            id,
		Name:          r.FormValue("name"),
		PostingKind:   dmpb.DisbursementMethodPostingKind(dmpb.DisbursementMethodPostingKind_value[r.FormValue("posting_kind")]),
		Category:      dmpb.DisbursementMethodCategory(dmpb.DisbursementMethodCategory_value[r.FormValue("category")]),
		TaxEffectKind: dmpb.DisbursementMethodTaxEffectKind(dmpb.DisbursementMethodTaxEffectKind_value[r.FormValue("tax_effect_kind")]),
		Lifecycle:     dmpb.DisbursementMethodLifecycle(dmpb.DisbursementMethodLifecycle_value[r.FormValue("lifecycle")]),
		Source:        dmpb.DisbursementMethodSource(dmpb.DisbursementMethodSource_value[r.FormValue("source")]),
		TemplateCode:  r.FormValue("template_code"),
		VersionStatus: dmpb.DisbursementMethodVersionStatus(dmpb.DisbursementMethodVersionStatus_value[r.FormValue("version_status")]),
	}
	if rev, err := strconv.ParseInt(r.FormValue("revision"), 10, 32); err == nil {
		m.Revision = int32(rev)
	}
	if v := r.FormValue("balance_account_id"); v != "" {
		m.BalanceAccountId = &v
	}
	if v := r.FormValue("target_account_id"); v != "" {
		m.TargetAccountId = &v
	}
	if v := r.FormValue("supersedes_disbursement_method_id"); v != "" {
		m.SupersedesDisbursementMethodId = &v
	}

	// Kind-specific template_details oneof (buying side: bank_account / advance_program).
	switch r.FormValue("category") {
	case "DISBURSEMENT_METHOD_CATEGORY_ADVANCE":
		ap := &dmpb.DisbursementMethodAdvanceProgramDetails{}
		if v := r.FormValue("default_balance_account_id"); v != "" {
			ap.DefaultBalanceAccountId = &v
		}
		if v := r.FormValue("default_target_account_id"); v != "" {
			ap.DefaultTargetAccountId = &v
		}
		if pc, err := strconv.ParseInt(r.FormValue("default_period_count"), 10, 32); err == nil && pc > 0 {
			pcc := int32(pc)
			ap.DefaultPeriodCount = &pcc
		}
		if v := r.FormValue("default_period_unit"); v != "" {
			ap.DefaultPeriodUnit = &v
		}
		m.TemplateDetails = &dmpb.DisbursementMethod_AdvanceProgram{AdvanceProgram: ap}
	case "DISBURSEMENT_METHOD_CATEGORY_STANDARD":
		// bank_account lives in the legacy `method_details` oneof (field 9),
		// NOT `template_details` — it predates the Stage-1 template extension.
		if bn := r.FormValue("bank_name"); bn != "" {
			m.MethodDetails = &dmpb.DisbursementMethod_BankAccount{
				BankAccount: &dmpb.DisbursementBankAccountDetails{BankName: bn},
			}
		}
	}
	return m, nil
}
