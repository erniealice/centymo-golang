package action

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"

	centymo "github.com/erniealice/centymo-golang"
	"github.com/erniealice/centymo-golang/views/collection_method/form"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	cmpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/collection_method"
)

// Deps holds all dependencies for the collection method action handlers.
// All use-case closures are nil-safe: the espyna collection_method use cases
// are NOT yet implemented (W1/W2 added proto + DB columns only). Until they
// land, the drawer renders + the HTMX fragment swap works (compile/structural
// correctness), but POST submit + lifecycle actions are no-ops that report a
// "not yet wired" error rather than panicking.
type Deps struct {
	Routes       centymo.CollectionMethodRoutes
	Labels       centymo.CollectionMethodLabels
	CommonLabels pyeza.CommonLabels

	CreateCollectionMethod func(ctx context.Context, req *cmpb.CreateCollectionMethodRequest) (*cmpb.CreateCollectionMethodResponse, error)
	ReadCollectionMethod   func(ctx context.Context, req *cmpb.ReadCollectionMethodRequest) (*cmpb.ReadCollectionMethodResponse, error)
	UpdateCollectionMethod func(ctx context.Context, req *cmpb.UpdateCollectionMethodRequest) (*cmpb.UpdateCollectionMethodResponse, error)
	DeleteCollectionMethod func(ctx context.Context, req *cmpb.DeleteCollectionMethodRequest) (*cmpb.DeleteCollectionMethodResponse, error)
}

const tableID = "collection-methods-table"

// NewAddAction handles GET+POST /action/collection-method/add.
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("collection_method", "create") {
			return centymo.HTMXError(fmt.Sprintf(deps.CommonLabels.Errors.MissingPermission, "collection_method:create"))
		}
		if viewCtx.Request.Method == http.MethodGet {
			fd := buildEmptyFormData(deps)
			fd.FormAction = deps.Routes.AddURL
			fd.Lifecycle = "COLLECTION_METHOD_LIFECYCLE_DRAFT"
			fd.Category = "COLLECTION_METHOD_CATEGORY_STANDARD"
			fd.Fragment = buildFragmentData(deps, fd.Category)
			return view.OK("collection-method-drawer-form", fd)
		}

		req, err := parseCollectionMethod(viewCtx.Request, "")
		if err != nil {
			return view.Error(err)
		}
		if deps.CreateCollectionMethod == nil {
			return centymo.HTMXError("Collection method create is not wired yet (espyna use cases pending).")
		}
		if _, err := deps.CreateCollectionMethod(ctx, &cmpb.CreateCollectionMethodRequest{Data: req}); err != nil {
			log.Printf("CreateCollectionMethod: %v", err)
			return view.Error(fmt.Errorf("failed to create collection method: %w", err))
		}
		return centymo.HTMXSuccess(tableID)
	})
}

// NewEditAction handles GET+POST /action/collection-method/edit/{id}.
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("collection_method", "update") {
			return centymo.HTMXError(fmt.Sprintf(deps.CommonLabels.Errors.MissingPermission, "collection_method:update"))
		}
		id := viewCtx.Request.PathValue("id")
		if id == "" {
			return view.Error(fmt.Errorf("missing id"))
		}

		if viewCtx.Request.Method == http.MethodGet {
			if deps.ReadCollectionMethod == nil {
				return centymo.HTMXError("Collection method read is not wired yet (espyna use cases pending).")
			}
			resp, err := deps.ReadCollectionMethod(ctx, &cmpb.ReadCollectionMethodRequest{
				Data: &cmpb.CollectionMethod{Id: id},
			})
			if err != nil {
				return view.Error(fmt.Errorf("failed to read collection method: %w", err))
			}
			data := resp.GetData()
			if len(data) == 0 {
				return view.Error(fmt.Errorf("collection method not found"))
			}
			fd := buildFormDataFromMethod(deps, data[0])
			fd.FormAction = route.ResolveURL(deps.Routes.EditURL, "id", id)
			fd.IsEdit = true
			fd.ID = id
			return view.OK("collection-method-drawer-form", fd)
		}

		req, err := parseCollectionMethod(viewCtx.Request, id)
		if err != nil {
			return view.Error(err)
		}
		if deps.UpdateCollectionMethod == nil {
			return centymo.HTMXError("Collection method update is not wired yet (espyna use cases pending).")
		}
		if _, err := deps.UpdateCollectionMethod(ctx, &cmpb.UpdateCollectionMethodRequest{Data: req}); err != nil {
			log.Printf("UpdateCollectionMethod %s: %v", id, err)
			return view.Error(fmt.Errorf("failed to update collection method: %w", err))
		}
		return centymo.HTMXSuccess(tableID)
	})
}

// NewFragmentAction handles GET /action/collection-method/fragment?cat=X — the
// §A template-slot HTMX swap. The category select fires hx-get on change and
// this handler returns ONLY the kind-specific fragment HTML for #kind-specific-slot.
func NewFragmentAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("collection_method", "create") && !perms.Can("collection_method", "update") {
			return centymo.HTMXError(fmt.Sprintf(deps.CommonLabels.Errors.MissingPermission, "collection_method:create"))
		}
		cat := viewCtx.Request.URL.Query().Get("cat")
		if cat == "" {
			cat = viewCtx.Request.URL.Query().Get("category")
		}
		fragData := buildFragmentData(deps, cat)
		return view.OK("collection-method-kind-fragment", fragData)
	})
}

// NewDeleteAction handles POST /action/collection-method/delete.
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("collection_method", "update") {
			return centymo.HTMXError(fmt.Sprintf(deps.CommonLabels.Errors.MissingPermission, "collection_method:update"))
		}
		if viewCtx.Request.Method != http.MethodPost {
			return view.Error(fmt.Errorf("method not allowed"))
		}
		id := viewCtx.Request.FormValue("id")
		if id == "" {
			return view.Error(fmt.Errorf("missing id"))
		}
		if deps.DeleteCollectionMethod == nil {
			return centymo.HTMXError("Collection method delete is not wired yet (espyna use cases pending).")
		}
		if _, err := deps.DeleteCollectionMethod(ctx, &cmpb.DeleteCollectionMethodRequest{
			Data: &cmpb.CollectionMethod{Id: id},
		}); err != nil {
			log.Printf("DeleteCollectionMethod %s: %v", id, err)
			return view.Error(fmt.Errorf("failed to delete collection method: %w", err))
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
	fd.CategoryOptions = []types.SelectOption{
		{Value: "COLLECTION_METHOD_CATEGORY_STANDARD", Label: l.Form.CategoryStandard},
		{Value: "COLLECTION_METHOD_CATEGORY_VOUCHER", Label: l.Form.CategoryVoucher},
		{Value: "COLLECTION_METHOD_CATEGORY_ADVANCE", Label: l.Form.CategoryAdvance},
		{Value: "COLLECTION_METHOD_CATEGORY_CARD", Label: l.Form.CategoryCard},
	}
	fd.PostingKindOptions = []types.SelectOption{
		{Value: "COLLECTION_METHOD_POSTING_KIND_CASH", Label: "Cash"},
		{Value: "COLLECTION_METHOD_POSTING_KIND_ADVANCE_DRAWDOWN", Label: "Advance Drawdown"},
		{Value: "COLLECTION_METHOD_POSTING_KIND_CLAIM_AR", Label: "Claim / AR"},
		{Value: "COLLECTION_METHOD_POSTING_KIND_DEFERRED_RECEIVABLE", Label: "Deferred Receivable"},
	}
	fd.AudienceModeOptions = []types.SelectOption{
		{Value: "COLLECTION_METHOD_AUDIENCE_MODE_OPEN", Label: "Open"},
		{Value: "COLLECTION_METHOD_AUDIENCE_MODE_RESTRICTED", Label: "Restricted"},
		{Value: "COLLECTION_METHOD_AUDIENCE_MODE_SINGLE_CLIENT", Label: "Single Client"},
	}
	fd.TaxEffectOptions = []types.SelectOption{
		{Value: "COLLECTION_METHOD_TAX_EFFECT_KIND_NONE", Label: "None"},
		{Value: "COLLECTION_METHOD_TAX_EFFECT_KIND_INCLUSIVE", Label: "Tax Inclusive"},
		{Value: "COLLECTION_METHOD_TAX_EFFECT_KIND_EXCLUSIVE", Label: "Tax Exclusive"},
	}
	fd.LifecycleOptions = []types.SelectOption{
		{Value: "COLLECTION_METHOD_LIFECYCLE_DRAFT", Label: "Draft"},
		{Value: "COLLECTION_METHOD_LIFECYCLE_ACTIVE", Label: "Active"},
		{Value: "COLLECTION_METHOD_LIFECYCLE_CLOSED", Label: "Closed"},
		{Value: "COLLECTION_METHOD_LIFECYCLE_ARCHIVED", Label: "Archived"},
	}
	fd.SourceOptions = []types.SelectOption{
		{Value: "COLLECTION_METHOD_SOURCE_WORKSPACE", Label: "Workspace"},
		{Value: "COLLECTION_METHOD_SOURCE_SYSTEM", Label: "System"},
		{Value: "COLLECTION_METHOD_SOURCE_VENDOR_TEMPLATE", Label: "Vendor Template"},
	}
	fd.VersionStatusOptions = []types.SelectOption{
		{Value: "COLLECTION_METHOD_VERSION_STATUS_DRAFT", Label: "Draft"},
		{Value: "COLLECTION_METHOD_VERSION_STATUS_PUBLISHED", Label: "Published"},
		{Value: "COLLECTION_METHOD_VERSION_STATUS_SUPERSEDED", Label: "Superseded"},
	}
	return fd
}

func buildFormDataFromMethod(deps *Deps, m *cmpb.CollectionMethod) *form.Data {
	fd := buildEmptyFormData(deps)
	fd.Name = m.GetName()
	fd.Category = m.GetCategory().String()
	fd.PostingKind = m.GetPostingKind().String()
	fd.AudienceMode = m.GetAudienceMode().String()
	fd.TaxEffectKind = m.GetTaxEffectKind().String()
	fd.EligibilityRule = m.GetDefaultEligibilityRuleId()
	fd.BalanceAccount = m.GetBalanceAccountId()
	fd.TargetAccount = m.GetTargetAccountId()
	fd.Lifecycle = m.GetLifecycle().String()
	fd.Source = m.GetSource().String()
	fd.TemplateCode = m.GetTemplateCode()
	fd.Revision = strconv.FormatInt(int64(m.GetRevision()), 10)
	fd.VersionStatus = m.GetVersionStatus().String()
	fd.Supersedes = m.GetSupersedesCollectionMethodId()
	fd.Fragment = buildFragmentDataFromMethod(deps, m)
	return fd
}

func buildFragmentData(deps *Deps, category string) form.FragmentData {
	return form.FragmentData{
		Category:   category,
		FragLabels: deps.Labels.Fragment,
	}
}

func buildFragmentDataFromMethod(deps *Deps, m *cmpb.CollectionMethod) form.FragmentData {
	fd := buildFragmentData(deps, m.GetCategory().String())
	if vp := m.GetVoucherProgram(); vp != nil {
		if vp.GetDefaultFaceValueCentavos() > 0 {
			fd.DefaultFaceValue = formatCentavos(vp.GetDefaultFaceValueCentavos())
		}
		if vp.GetDefaultExpiryDays() > 0 {
			fd.DefaultExpiryDays = strconv.FormatInt(int64(vp.GetDefaultExpiryDays()), 10)
		}
	} else if ap := m.GetAdvanceProgram(); ap != nil {
		fd.AdvanceKind = ap.GetAdvanceKind().String()
		fd.DefaultBalanceAcct = ap.GetDefaultBalanceAccountId()
		fd.DefaultTargetAcct = ap.GetDefaultTargetAccountId()
		if ap.GetDefaultPeriodCount() > 0 {
			fd.DefaultPeriodCount = strconv.FormatInt(int64(ap.GetDefaultPeriodCount()), 10)
		}
		fd.DefaultPeriodUnit = ap.GetDefaultPeriodUnit()
	}
	return fd
}

// parseCollectionMethod builds a proto from the submitted form. Enum fields are
// parsed via the proto enum value maps (string → enum number).
func parseCollectionMethod(r *http.Request, id string) (*cmpb.CollectionMethod, error) {
	if err := r.ParseForm(); err != nil {
		return nil, fmt.Errorf("parse form: %w", err)
	}
	m := &cmpb.CollectionMethod{
		Id:            id,
		Name:          r.FormValue("name"),
		PostingKind:   cmpb.CollectionMethodPostingKind(cmpb.CollectionMethodPostingKind_value[r.FormValue("posting_kind")]),
		Category:      cmpb.CollectionMethodCategory(cmpb.CollectionMethodCategory_value[r.FormValue("category")]),
		AudienceMode:  cmpb.CollectionMethodAudienceMode(cmpb.CollectionMethodAudienceMode_value[r.FormValue("audience_mode")]),
		TaxEffectKind: cmpb.CollectionMethodTaxEffectKind(cmpb.CollectionMethodTaxEffectKind_value[r.FormValue("tax_effect_kind")]),
		Lifecycle:     cmpb.CollectionMethodLifecycle(cmpb.CollectionMethodLifecycle_value[r.FormValue("lifecycle")]),
		Source:        cmpb.CollectionMethodSource(cmpb.CollectionMethodSource_value[r.FormValue("source")]),
		TemplateCode:  r.FormValue("template_code"),
		VersionStatus: cmpb.CollectionMethodVersionStatus(cmpb.CollectionMethodVersionStatus_value[r.FormValue("version_status")]),
	}
	if rev, err := strconv.ParseInt(r.FormValue("revision"), 10, 32); err == nil {
		m.Revision = int32(rev)
	}
	if v := r.FormValue("default_eligibility_rule_id"); v != "" {
		m.DefaultEligibilityRuleId = &v
	}
	if v := r.FormValue("balance_account_id"); v != "" {
		m.BalanceAccountId = &v
	}
	if v := r.FormValue("target_account_id"); v != "" {
		m.TargetAccountId = &v
	}
	if v := r.FormValue("supersedes_collection_method_id"); v != "" {
		m.SupersedesCollectionMethodId = &v
	}

	// Kind-specific template_details oneof.
	switch r.FormValue("category") {
	case "COLLECTION_METHOD_CATEGORY_VOUCHER":
		vp := &cmpb.CollectionMethodVoucherProgramDetails{}
		if fv := parseCentavos(r.FormValue("default_face_value")); fv > 0 {
			vp.DefaultFaceValueCentavos = &fv
		}
		if d, err := strconv.ParseInt(r.FormValue("default_expiry_days"), 10, 32); err == nil && d > 0 {
			dd := int32(d)
			vp.DefaultExpiryDays = &dd
		}
		m.TemplateDetails = &cmpb.CollectionMethod_VoucherProgram{VoucherProgram: vp}
	case "COLLECTION_METHOD_CATEGORY_ADVANCE":
		ap := &cmpb.CollectionMethodAdvanceProgramDetails{}
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
		m.TemplateDetails = &cmpb.CollectionMethod_AdvanceProgram{AdvanceProgram: ap}
	case "COLLECTION_METHOD_CATEGORY_CARD":
		// D-4.26: card template carries NO fields (brand is per-instance).
		m.TemplateDetails = &cmpb.CollectionMethod_CardType{CardType: &cmpb.CollectionMethodCardTypeDetails{}}
	}
	return m, nil
}

func parseCentavos(s string) int64 {
	if s == "" {
		return 0
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return int64(math.Round(f * 100))
}

func formatCentavos(v int64) string {
	if v == 0 {
		return ""
	}
	return strconv.FormatFloat(float64(v)/100.0, 'f', 2, 64)
}
