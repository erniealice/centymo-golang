package plan

import (
	"context"

	productpriceplan "github.com/erniealice/centymo-golang/domain/subscription/product_price_plan"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	priceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_plan"
	priceschedulepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_schedule"
)

// NewProductPriceAddAction handles add under the schedule namespace.
// Delegates to productpriceplan.NewAddAction.
func NewProductPriceAddAction(deps *DetailViewDeps) view.View {
	return productpriceplan.NewAddAction(pppDepsFromDetailDeps(deps))
}

// pppDepsFromDetailDeps translates DetailViewDeps into the subset required by
// the product_price_plan action package.
func pppDepsFromDetailDeps(deps *DetailViewDeps) *productpriceplan.Deps {
	return &productpriceplan.Deps{
		AddURL:                 deps.Routes.PlanProductPriceAddURL,
		EditURL:                deps.Routes.PlanProductPriceEditURL,
		RefreshTableID:         "price-schedule-plan-product-prices-table",
		DrawerTemplateName:     "price-schedule-plan-product-price-drawer",
		PlanLabels:             deps.PlanLabels,
		ProductPricePlanLabels: deps.ProductPricePlanLabels,
		ScheduleDetailLabels:   deps.ScheduleLabels.Detail,
		CommonLabels:           deps.CommonLabels,
		TableLabels:            deps.TableLabels,
		ReadPricePlan:          deps.ReadPricePlan,
		ReadPriceSchedule:      deps.ReadPriceSchedule,
		ListPlans:              deps.ListPlans,
		ListProducts:           deps.ListProducts,
		ListProductPlans:       deps.ListProductPlans,
		ListProductPricePlans:  deps.ListProductPricePlans,
		CreateProductPricePlan: deps.CreateProductPricePlan,
		UpdateProductPricePlan: deps.UpdateProductPricePlan,
		DeleteProductPricePlan: deps.DeleteProductPricePlan,
		GetPricePlanInUseIDs:   deps.GetPricePlanInUseIDs,
	}
}

// NewProductPriceEditAction handles edit under the schedule namespace.
// Delegates to productpriceplan.NewEditAction.
func NewProductPriceEditAction(deps *DetailViewDeps) view.View {
	return productpriceplan.NewEditAction(pppDepsFromDetailDeps(deps))
}

// NewProductPriceDeleteAction handles delete under the schedule namespace.
// Delegates to productpriceplan.NewDeleteAction.
func NewProductPriceDeleteAction(deps *DetailViewDeps) view.View {
	return productpriceplan.NewDeleteAction(pppDepsFromDetailDeps(deps))
}

// buildProductPricesTable assembles the product-prices tab table.
// Delegates to productpriceplan.BuildTable by resolving the parent context
// and constructing a TableDeps from this page's DetailViewDeps.
func buildProductPricesTable(ctx context.Context, deps *DetailViewDeps, sid, ppid string) *types.TableConfig {
	parent, _ := loadParentContext(ctx, deps, ppid)
	parentCtx := productpriceplan.ParentContext{
		Currency:              parent.Currency,
		BillingKind:           parent.BillingKind,
		AmountBasis:           parent.AmountBasis,
		BillingKindDisplay:    parent.BillingKindDisplay,
		AmountBasisDisplay:    parent.AmountBasisDisplay,
		BillingCycleDisplay:   parent.BillingCycleDisplay,
		TermDisplay:           parent.TermDisplay,
		ParentCurrencyDisplay: parent.ParentCurrencyDisplay,
		RateCardName:          parent.RateCardName,
	}
	tableDeps := &productpriceplan.TableDeps{
		PlanProductPriceEditURL: deps.Routes.PlanProductPriceEditURL,
		PlanTabActionURL:        deps.Routes.PlanTabActionURL,
		ProductPricesTabSlug:    deps.ScheduleLabels.Tabs.ResolveTabSlug("product-prices"),
		TableID:                 "price-schedule-plan-product-prices-table",
		ColumnProduct:           deps.ScheduleLabels.Detail.ProductPriceColumnProduct,
		ColumnPrice:             deps.ScheduleLabels.Detail.ProductPriceColumnPrice,
		ColumnCurrency:          deps.ScheduleLabels.Detail.ProductPriceColumnCurrency,
		ColumnTreatment:         deps.ScheduleLabels.Detail.ProductPriceColumnTreatment,
		ColumnEffective:         deps.ScheduleLabels.Detail.ProductPriceColumnEffective,
		EditLabel:               deps.ScheduleLabels.Detail.ProductPriceEdit,
		EmptyTitle:              deps.ScheduleLabels.Detail.ProductPriceEmptyTitle,
		EmptyMsg:                deps.ScheduleLabels.Detail.ProductPriceEmptyMsg,
		Unauthorized:            deps.PlanLabels.Errors.Unauthorized,
		SID:                     sid,
		PPID:                    ppid,
		ListProducts:            deps.ListProducts,
		ListProductPlans:        deps.ListProductPlans,
		ListProductPricePlans:   deps.ListProductPricePlans,
		PlanLabels:              deps.PlanLabels,
		ProductPricePlanLabels:  deps.ProductPricePlanLabels,
		TableLabels:             deps.TableLabels,
	}
	return productpriceplan.BuildTable(ctx, tableDeps, parentCtx)
}

// parentPricePlanContext captures the parent PricePlan fields the PPP drawer
// needs to know about: currency (locks the per-line currency), billing_kind
// (decides whether billing_treatment renders), amount_basis (drives the
// banner explaining what the line prices mean), and pre-formatted display
// strings for the read-only context block above the editable fields.
type parentPricePlanContext struct {
	Currency    string
	BillingKind string
	AmountBasis string

	// Display strings — empty when the corresponding source data is missing.
	BillingKindDisplay    string
	AmountBasisDisplay    string
	BillingCycleDisplay   string
	TermDisplay           string
	ParentCurrencyDisplay string
	RateCardName          string
}

func loadParentContext(ctx context.Context, deps *DetailViewDeps, pricePlanID string) (parentPricePlanContext, bool) {
	if deps.ReadPricePlan == nil || pricePlanID == "" {
		return parentPricePlanContext{}, false
	}
	resp, err := deps.ReadPricePlan(ctx, &priceplanpb.ReadPricePlanRequest{
		Data: &priceplanpb.PricePlan{Id: pricePlanID},
	})
	if err != nil || len(resp.GetData()) == 0 {
		return parentPricePlanContext{}, false
	}
	pp := resp.GetData()[0]
	pc := parentPricePlanContext{
		Currency:              pp.GetBillingCurrency(),
		BillingKind:           pp.GetBillingKind().String(),
		AmountBasis:           pp.GetAmountBasis().String(),
		ParentCurrencyDisplay: pp.GetBillingCurrency(),
		BillingKindDisplay:    productpriceplan.FormatBillingKindLabel(pp.GetBillingKind().String(), deps.PlanLabels.Form),
		AmountBasisDisplay:    productpriceplan.FormatAmountBasisLabel(pp.GetAmountBasis().String(), deps.PlanLabels.Form),
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
