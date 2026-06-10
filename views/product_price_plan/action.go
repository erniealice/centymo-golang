package product_price_plan

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"

	centymo "github.com/erniealice/centymo-golang"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	productpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product"
	productplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_plan"
	planpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/plan"
	priceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_plan"
	priceschedulepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_schedule"
	productpriceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/product_price_plan"
)

// Deps holds all dependencies for the ProductPricePlan CRUD action handlers.
// Parent modules (price_schedule, price_plan) construct this from their own
// DetailViewDeps and pass it to NewAddAction / NewEditAction / NewDeleteAction.
type Deps struct {
	// URL templates for form actions and HTMX refresh targets.
	// AddURL / EditURL are templates with {id}, {ppid}, and (for edit) {pppid} params.
	AddURL  string
	EditURL string
	// RefreshTableID is the HTMX target ID to refresh after a successful CUD.
	// e.g. "price-schedule-plan-product-prices-table" or "price-plan-product-prices-table".
	RefreshTableID string

	// DrawerTemplateName is the template name for the add/edit drawer form.
	// e.g. "price-schedule-plan-product-price-drawer" or "price-plan-product-price-drawer".
	DrawerTemplateName string

	// URL path params — the parent passes the pre-resolved IDs. "id" is the
	// schedule/parent ID and "ppid" is the price_plan ID. Both are read from
	// the incoming request path values by each handler (not from Deps); these
	// fields are unused by the handlers themselves (reserved for future use).

	// Labels.
	PlanLabels             centymo.PricePlanLabels
	ProductPricePlanLabels centymo.ProductPricePlanLabels
	ScheduleDetailLabels   centymo.PriceScheduleDetailLabels
	CommonLabels           pyeza.CommonLabels
	TableLabels            types.TableLabels

	// Use-case functions.
	ReadPricePlan          func(ctx context.Context, req *priceplanpb.ReadPricePlanRequest) (*priceplanpb.ReadPricePlanResponse, error)
	ReadPriceSchedule      func(ctx context.Context, req *priceschedulepb.ReadPriceScheduleRequest) (*priceschedulepb.ReadPriceScheduleResponse, error)
	ListPlans              func(ctx context.Context, req *planpb.ListPlansRequest) (*planpb.ListPlansResponse, error)
	ListProducts           func(ctx context.Context, req *productpb.ListProductsRequest) (*productpb.ListProductsResponse, error)
	ListProductPlans       func(ctx context.Context, req *productplanpb.ListProductPlansRequest) (*productplanpb.ListProductPlansResponse, error)
	ListProductPricePlans  func(ctx context.Context, req *productpriceplanpb.ListProductPricePlansRequest) (*productpriceplanpb.ListProductPricePlansResponse, error)
	CreateProductPricePlan func(ctx context.Context, req *productpriceplanpb.CreateProductPricePlanRequest) (*productpriceplanpb.CreateProductPricePlanResponse, error)
	UpdateProductPricePlan func(ctx context.Context, req *productpriceplanpb.UpdateProductPricePlanRequest) (*productpriceplanpb.UpdateProductPricePlanResponse, error)
	DeleteProductPricePlan func(ctx context.Context, req *productpriceplanpb.DeleteProductPricePlanRequest) (*productpriceplanpb.DeleteProductPricePlanResponse, error)

	// GetPricePlanInUseIDs checks whether a PricePlan is referenced by active
	// subscriptions. When true, per-item price editing is rejected on POST.
	GetPricePlanInUseIDs func(ctx context.Context, ids []string) (map[string]bool, error)
}

// NewAddAction handles GET (render drawer) and POST (submit) for adding a
// ProductPricePlan under the schedule namespace.
// Path values read: "id" (schedule/parent id), "ppid" (price_plan id).
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("product_price_plan", "create") {
			return view.HTMXError(deps.PlanLabels.Errors.Unauthorized)
		}
		if deps.CreateProductPricePlan == nil {
			return view.HTMXError(deps.PlanLabels.Messages.CreateNotAvailable)
		}
		sid := viewCtx.Request.PathValue("id")
		ppid := viewCtx.Request.PathValue("ppid")

		pplLabels := deps.ProductPricePlanLabels.Form
		if viewCtx.Request.Method == http.MethodGet {
			planName, planDesc := lookupPackageNameDesc(ctx, deps, ppid)
			parent, _ := loadParentContext(ctx, deps, ppid)
			currency := parent.Currency
			if currency == "" {
				currency = "PHP"
			}
			showTreatment := parent.BillingKind != "BILLING_KIND_ONE_TIME"
			return view.OK(deps.DrawerTemplateName, &productPriceFormData{
				FormAction:             route.ResolveURL(deps.AddURL, "id", sid, "ppid", ppid),
				ScheduleID:             sid,
				PricePlanID:            ppid,
				Currency:               currency,
				CommonLabels:           deps.CommonLabels,
				PlanName:               planName,
				PlanDescription:        planDesc,
				ParentBillingKind:      parent.BillingKind,
				ParentAmountBasis:      parent.AmountBasis,
				ShowTreatment:          showTreatment,
				BasisBannerMessage:     basisBannerMessage(parent.AmountBasis, deps.ScheduleDetailLabels),
				BillingKindDisplay:     parent.BillingKindDisplay,
				AmountBasisDisplay:     parent.AmountBasisDisplay,
				BillingCycleDisplay:    parent.BillingCycleDisplay,
				TermDisplay:            parent.TermDisplay,
				ParentCurrencyDisplay:  parent.ParentCurrencyDisplay,
				RateCardName:           parent.RateCardName,
				ProductPricePlanLabels: pplLabels,
				PriceScheduleLabels:    deps.ScheduleDetailLabels,
			})
		}

		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.HTMXError(deps.PlanLabels.Errors.Unauthorized)
		}
		productPlanID := viewCtx.Request.FormValue("product_plan_id")
		if productPlanID == "" {
			// Backward-compatible: old form posts may still send product_id;
			// resolve to a product_plan_id on this plan when possible.
			if legacyProductID := viewCtx.Request.FormValue("product_id"); legacyProductID != "" {
				productPlanID = resolveProductPlanIDForProduct(ctx, deps, ppid, legacyProductID)
			}
		}
		if productPlanID == "" {
			return view.HTMXError(deps.PlanLabels.Messages.ProductRequired)
		}
		priceCentavos, ok := parsePriceCentavos(viewCtx.Request.FormValue("price"))
		if !ok {
			return view.HTMXError(deps.PlanLabels.Messages.InvalidPrice)
		}
		currency := viewCtx.Request.FormValue("currency")
		if currency == "" {
			currency = "PHP"
		}
		dateStart := viewCtx.Request.FormValue("date_start")
		dateEnd := viewCtx.Request.FormValue("date_end")
		billingTreatment := viewCtx.Request.FormValue("billing_treatment")
		parent, _ := loadParentContext(ctx, deps, ppid)
		// Currency must match parent PricePlan.billing_currency (proto invariant).
		if parent.Currency != "" && currency != parent.Currency {
			return view.HTMXError(deps.PlanLabels.Messages.CurrencyMismatch)
		}
		// billing_treatment is meaningless when parent has no cycles. Drop the
		// posted value so we never persist a stale treatment on a ONE_TIME plan.
		if parent.BillingKind == "BILLING_KIND_ONE_TIME" {
			billingTreatment = ""
		}
		record := &productpriceplanpb.ProductPricePlan{
			PricePlanId:     ppid,
			ProductPlanId:   productPlanID,
			BillingAmount:   priceCentavos,
			BillingCurrency: currency,
			Active:          true,
		}
		if billingTreatment != "" {
			if bt, ok := productpriceplanpb.BillingTreatment_value[billingTreatment]; ok {
				record.BillingTreatment = productpriceplanpb.BillingTreatment(bt)
			}
		}
		if dateStart != "" {
			record.DateStart = &dateStart
		}
		if dateEnd != "" {
			record.DateEnd = &dateEnd
		}
		if _, err := deps.CreateProductPricePlan(ctx, &productpriceplanpb.CreateProductPricePlanRequest{Data: record}); err != nil {
			log.Printf("Failed to create product price plan for plan %s (parent %s): %v", ppid, sid, err)
			return view.HTMXError(err.Error())
		}
		return view.HTMXSuccess(deps.RefreshTableID)
	})
}

// NewEditAction handles GET (render drawer) and POST (submit) for editing a
// ProductPricePlan under the schedule namespace.
// Path values read: "id" (schedule/parent id), "ppid" (price_plan id), "pppid" (product_price_plan id).
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("product_price_plan", "update") {
			return view.HTMXError(deps.PlanLabels.Errors.Unauthorized)
		}
		if deps.UpdateProductPricePlan == nil {
			return view.HTMXError(deps.PlanLabels.Messages.UpdateNotAvailable)
		}
		sid := viewCtx.Request.PathValue("id")
		ppid := viewCtx.Request.PathValue("ppid")
		pppid := viewCtx.Request.PathValue("pppid")

		existing, err := findProductPricePlan(ctx, deps, pppid)
		if err != nil {
			return view.HTMXError(err.Error())
		}

		pplLabels := deps.ProductPricePlanLabels.Form
		if viewCtx.Request.Method == http.MethodGet {
			parent, _ := loadParentContext(ctx, deps, ppid)
			currency := existing.GetBillingCurrency()
			if currency == "" {
				currency = parent.Currency
			}
			if currency == "" {
				currency = "PHP"
			}
			showTreatment := parent.BillingKind != "BILLING_KIND_ONE_TIME"
			planName, planDesc := lookupPackageNameDesc(ctx, deps, ppid)
			// Model D — resolve product + variant via the referenced ProductPlan row.
			existingProductPlanID := existing.GetProductPlanId()
			prodName, prodDesc, variantName := lookupProductPlanDisplay(ctx, deps, existingProductPlanID)

			pricingLocked := false
			pricingLockedReason := ""
			if deps.GetPricePlanInUseIDs != nil {
				if inUse, _ := deps.GetPricePlanInUseIDs(ctx, []string{ppid}); inUse[ppid] {
					pricingLocked = true
					pricingLockedReason = deps.PlanLabels.Messages.ItemPricingLockedReason
				}
			}

			return view.OK(deps.DrawerTemplateName, &productPriceFormData{
				FormAction:          route.ResolveURL(deps.EditURL, "id", sid, "ppid", ppid, "pppid", pppid),
				IsEdit:              true,
				ID:                  pppid,
				ScheduleID:          sid,
				PricePlanID:         ppid,
				ProductPlanID:       existingProductPlanID,
				Price:               fmt.Sprintf("%.2f", float64(existing.GetBillingAmount())/100.0),
				Currency:            currency,
				CommonLabels:        deps.CommonLabels,
				PlanName:            planName,
				PlanDescription:     planDesc,
				ProductName:         prodName,
				ProductDescription:  prodDesc,
				VariantName:         variantName,
				PricingLocked:       pricingLocked,
				PricingLockedReason: pricingLockedReason,
				// Wave 2: populate billing treatment and dates from existing record.
				BillingTreatment:       existing.GetBillingTreatment().String(),
				DateStart:              existing.GetDateStart(),
				DateEnd:                existing.GetDateEnd(),
				ParentBillingKind:      parent.BillingKind,
				ParentAmountBasis:      parent.AmountBasis,
				ShowTreatment:          showTreatment,
				BasisBannerMessage:     basisBannerMessage(parent.AmountBasis, deps.ScheduleDetailLabels),
				BillingKindDisplay:     parent.BillingKindDisplay,
				AmountBasisDisplay:     parent.AmountBasisDisplay,
				BillingCycleDisplay:    parent.BillingCycleDisplay,
				TermDisplay:            parent.TermDisplay,
				ParentCurrencyDisplay:  parent.ParentCurrencyDisplay,
				RateCardName:           parent.RateCardName,
				ProductPricePlanLabels: pplLabels,
				PriceScheduleLabels:    deps.ScheduleDetailLabels,
			})
		}

		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.HTMXError(deps.PlanLabels.Errors.Unauthorized)
		}
		// Server-side lock enforcement: if the parent PricePlan is in use by an
		// active subscription, reject price/currency changes (client may have
		// bypassed the disabled inputs).
		if deps.GetPricePlanInUseIDs != nil {
			if inUse, _ := deps.GetPricePlanInUseIDs(ctx, []string{ppid}); inUse[ppid] {
				return view.HTMXError(deps.PlanLabels.Messages.InUseCannotModify)
			}
		}
		// The catalog-line assignment is display-only in the drawer — preserve
		// the existing product_plan_id.
		priceCentavos, ok := parsePriceCentavos(viewCtx.Request.FormValue("price"))
		if !ok {
			return view.HTMXError(deps.PlanLabels.Messages.InvalidPrice)
		}
		currency := viewCtx.Request.FormValue("currency")
		if currency == "" {
			currency = "PHP"
		}
		dateStart := viewCtx.Request.FormValue("date_start")
		dateEnd := viewCtx.Request.FormValue("date_end")
		billingTreatment := viewCtx.Request.FormValue("billing_treatment")
		parent, _ := loadParentContext(ctx, deps, ppid)
		if parent.Currency != "" && currency != parent.Currency {
			return view.HTMXError(deps.PlanLabels.Messages.CurrencyMismatch)
		}
		if parent.BillingKind == "BILLING_KIND_ONE_TIME" {
			billingTreatment = ""
		}
		updated := &productpriceplanpb.ProductPricePlan{
			Id:              pppid,
			PricePlanId:     ppid,
			ProductPlanId:   existing.GetProductPlanId(),
			BillingAmount:   priceCentavos,
			BillingCurrency: currency,
			Active:          existing.GetActive(),
		}
		if billingTreatment != "" {
			if bt, ok := productpriceplanpb.BillingTreatment_value[billingTreatment]; ok {
				updated.BillingTreatment = productpriceplanpb.BillingTreatment(bt)
			}
		}
		if dateStart != "" {
			updated.DateStart = &dateStart
		}
		if dateEnd != "" {
			updated.DateEnd = &dateEnd
		}
		if _, err := deps.UpdateProductPricePlan(ctx, &productpriceplanpb.UpdateProductPricePlanRequest{Data: updated}); err != nil {
			log.Printf("Failed to update product price plan %s: %v", pppid, err)
			return view.HTMXError(err.Error())
		}
		return view.HTMXSuccess(deps.RefreshTableID)
	})
}

// NewDeleteAction handles POST for deleting a ProductPricePlan.
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("product_price_plan", "delete") {
			return view.HTMXError(deps.PlanLabels.Errors.Unauthorized)
		}
		if deps.DeleteProductPricePlan == nil {
			return view.HTMXError(deps.PlanLabels.Messages.DeleteNotAvailable)
		}
		_ = viewCtx.Request.ParseForm()
		pppid := viewCtx.Request.FormValue("id")
		if pppid == "" {
			pppid = viewCtx.Request.URL.Query().Get("id")
		}
		if pppid == "" {
			return view.HTMXError(deps.PlanLabels.Messages.IDRequired)
		}
		if _, err := deps.DeleteProductPricePlan(ctx, &productpriceplanpb.DeleteProductPricePlanRequest{
			Data: &productpriceplanpb.ProductPricePlan{Id: pppid},
		}); err != nil {
			log.Printf("Failed to delete product price plan %s: %v", pppid, err)
			return view.HTMXError(err.Error())
		}
		return view.HTMXSuccess(deps.RefreshTableID)
	})
}

// ---------------------------------------------------------------------------
// Internal helpers — mirror of the functions that remain in the parent package
// but operate on *Deps instead of *DetailViewDeps.
// ---------------------------------------------------------------------------

// productPriceFormData mirrors ProductPriceFormData in the parent package.
// Kept unexported because the parent packages already define this struct under
// their own package; the template only receives it via interface{} / any.
type productPriceFormData struct {
	FormAction    string
	WorkspaceID   string // injected by C1: populated by ViewAdapter.injectWorkspaceID for action_workspace_guard
	IsEdit        bool
	ID            string
	ScheduleID    string
	PricePlanID   string
	ProductPlanID string
	Price         string
	Currency      string
	CommonLabels  pyeza.CommonLabels

	// Display-only context (read-only).
	PlanName           string
	PlanDescription    string
	ProductName        string
	ProductDescription string
	VariantName        string

	// Wave 2: billing treatment + effective date fields.
	BillingTreatment string
	DateStart        string
	DateEnd          string

	// Parent PricePlan context.
	ParentBillingKind  string
	ParentAmountBasis  string
	ShowTreatment      bool
	BasisBannerMessage string

	// Read-only "package context" block.
	BillingKindDisplay    string
	AmountBasisDisplay    string
	BillingCycleDisplay   string
	TermDisplay           string
	ParentCurrencyDisplay string
	RateCardName          string

	// Labels.
	ProductPricePlanLabels centymo.ProductPricePlanFormLabels
	PriceScheduleLabels    centymo.PriceScheduleDetailLabels

	// Pricing lock.
	PricingLocked       bool
	PricingLockedReason string
}

// loadParentContext resolves parent PricePlan fields needed by the PPP drawer.
// Mirrors the same-named function in price_schedule/detail/plan/page.go but
// operates on *Deps so there is no import cycle.
func loadParentContext(ctx context.Context, deps *Deps, pricePlanID string) (ParentContext, bool) {
	if deps.ReadPricePlan == nil || pricePlanID == "" {
		return ParentContext{}, false
	}
	resp, err := deps.ReadPricePlan(ctx, &priceplanpb.ReadPricePlanRequest{
		Data: &priceplanpb.PricePlan{Id: pricePlanID},
	})
	if err != nil || len(resp.GetData()) == 0 {
		return ParentContext{}, false
	}
	pp := resp.GetData()[0]
	pc := ParentContext{
		Currency:              pp.GetBillingCurrency(),
		BillingKind:           pp.GetBillingKind().String(),
		AmountBasis:           pp.GetAmountBasis().String(),
		ParentCurrencyDisplay: pp.GetBillingCurrency(),
		BillingKindDisplay:    FormatBillingKindLabel(pp.GetBillingKind().String(), deps.PlanLabels.Form),
		AmountBasisDisplay:    FormatAmountBasisLabel(pp.GetAmountBasis().String(), deps.PlanLabels.Form),
	}
	if v := pp.GetBillingCycleValue(); v > 0 {
		pc.BillingCycleDisplay = pyeza.FormatDuration(v, pp.GetBillingCycleUnit(), deps.CommonLabels.DurationUnit)
	}
	if v := pp.GetDefaultTermValue(); v > 0 {
		pc.TermDisplay = pyeza.FormatDuration(v, pp.GetDefaultTermUnit(), deps.CommonLabels.DurationUnit)
	}
	if scheduleID := pp.GetPriceScheduleId(); scheduleID != "" && deps.ReadPriceSchedule != nil {
		if schedResp, err := deps.ReadPriceSchedule(ctx, &priceschedulepb.ReadPriceScheduleRequest{
			Data: &priceschedulepb.PriceSchedule{Id: scheduleID},
		}); err == nil && len(schedResp.GetData()) > 0 {
			pc.RateCardName = schedResp.GetData()[0].GetName()
		}
	}
	return pc, true
}

// findProductPricePlan looks up a single ProductPricePlan by its ID.
func findProductPricePlan(ctx context.Context, deps *Deps, pppid string) (*productpriceplanpb.ProductPricePlan, error) {
	if deps.ListProductPricePlans == nil {
		return nil, fmt.Errorf("product price plans not available")
	}
	resp, err := deps.ListProductPricePlans(ctx, &productpriceplanpb.ListProductPricePlansRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to load product price plans")
	}
	for _, item := range resp.GetData() {
		if item != nil && item.GetId() == pppid {
			return item, nil
		}
	}
	return nil, fmt.Errorf("product price plan not found")
}

// lookupPackageNameDesc resolves the display name + description for a price_plan,
// falling back to the linked Plan's values when the price_plan fields are blank.
func lookupPackageNameDesc(ctx context.Context, deps *Deps, pricePlanID string) (string, string) {
	if pricePlanID == "" || deps.ReadPricePlan == nil {
		return "", ""
	}
	resp, err := deps.ReadPricePlan(ctx, &priceplanpb.ReadPricePlanRequest{
		Data: &priceplanpb.PricePlan{Id: pricePlanID},
	})
	if err != nil || len(resp.GetData()) == 0 {
		return "", ""
	}
	pp := resp.GetData()[0]
	name := strings.TrimSpace(pp.GetName())
	desc := strings.TrimSpace(pp.GetDescription())
	if name != "" && desc != "" {
		return name, desc
	}
	planName, planDesc := lookupPlanNameDesc(ctx, deps, pp.GetPlanId())
	if name == "" {
		name = planName
	}
	if desc == "" {
		desc = planDesc
	}
	return name, desc
}

// lookupPlanNameDesc returns the linked Plan's name and description (trimmed).
func lookupPlanNameDesc(ctx context.Context, deps *Deps, planID string) (string, string) {
	if planID == "" || deps.ListPlans == nil {
		return "", ""
	}
	resp, err := deps.ListPlans(ctx, &planpb.ListPlansRequest{})
	if err != nil {
		return "", ""
	}
	for _, p := range resp.GetData() {
		if p == nil || p.GetId() != planID {
			continue
		}
		return strings.TrimSpace(p.GetName()), strings.TrimSpace(p.GetDescription())
	}
	return "", ""
}

// lookupProductPlanDisplay resolves product name, description, and (optional)
// variant SKU for a ProductPlan.id — used to render the read-only context
// rows on the product-price drawer under Model D.
func lookupProductPlanDisplay(ctx context.Context, deps *Deps, productPlanID string) (name, desc, variant string) {
	if productPlanID == "" || deps.ListProductPlans == nil {
		return "", "", ""
	}
	ppResp, err := deps.ListProductPlans(ctx, &productplanpb.ListProductPlansRequest{})
	if err != nil {
		return "", "", ""
	}
	var (
		productID string
		variantID string
	)
	for _, pp := range ppResp.GetData() {
		if pp != nil && pp.GetId() == productPlanID {
			productID = pp.GetProductId()
			variantID = pp.GetProductVariantId()
			break
		}
	}
	name, desc = lookupProductNameDesc(ctx, deps, productID)
	variant = variantID
	return name, desc, variant
}

// lookupProductNameDesc reads the Product and returns its trimmed name + description.
func lookupProductNameDesc(ctx context.Context, deps *Deps, productID string) (string, string) {
	if productID == "" || deps.ListProducts == nil {
		return "", ""
	}
	resp, err := deps.ListProducts(ctx, &productpb.ListProductsRequest{})
	if err != nil {
		return "", ""
	}
	for _, p := range resp.GetData() {
		if p == nil || p.GetId() != productID {
			continue
		}
		return strings.TrimSpace(p.GetName()), strings.TrimSpace(p.GetDescription())
	}
	return "", ""
}

// resolveProductPlanIDForProduct finds the ProductPlan row in the parent
// Plan of the given PricePlan that references the supplied product_id.
func resolveProductPlanIDForProduct(ctx context.Context, deps *Deps, pricePlanID, productID string) string {
	if productID == "" || deps.ListProductPlans == nil {
		return ""
	}
	planID := loadPricePlanPlanID(ctx, deps, pricePlanID)
	if planID == "" {
		return ""
	}
	ppResp, err := deps.ListProductPlans(ctx, &productplanpb.ListProductPlansRequest{})
	if err != nil {
		return ""
	}
	for _, pp := range ppResp.GetData() {
		if pp != nil && pp.GetPlanId() == planID && pp.GetProductId() == productID {
			return pp.GetId()
		}
	}
	return ""
}

// loadPricePlanPlanID reads the price plan to get its linked plan_id.
func loadPricePlanPlanID(ctx context.Context, deps *Deps, pricePlanID string) string {
	if deps.ReadPricePlan == nil {
		return ""
	}
	resp, err := deps.ReadPricePlan(ctx, &priceplanpb.ReadPricePlanRequest{
		Data: &priceplanpb.PricePlan{Id: pricePlanID},
	})
	if err != nil || len(resp.GetData()) == 0 {
		return ""
	}
	return resp.GetData()[0].GetPlanId()
}

// parsePriceCentavos parses a decimal price string and converts to centavos.
func parsePriceCentavos(s string) (int64, bool) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil || f < 0 {
		return 0, false
	}
	return int64(math.Round(f * 100)), true
}

// basisBannerMessage returns a one-line explanation about the parent's amount_basis.
func basisBannerMessage(amountBasis string, l centymo.PriceScheduleDetailLabels) string {
	switch amountBasis {
	case "AMOUNT_BASIS_DERIVED_FROM_LINES":
		return l.BasisBannerDerived
	case "AMOUNT_BASIS_TOTAL_PACKAGE":
		return l.BasisBannerTotalPackage
	case "AMOUNT_BASIS_PER_CYCLE":
		return l.BasisBannerPerCycle
	}
	return ""
}
