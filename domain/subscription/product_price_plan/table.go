package product_price_plan

import (
	"context"
	"fmt"
	"log"

	sib_subscription_price_plan "github.com/erniealice/centymo-golang/domain/subscription/price_plan"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	productpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product"
	productplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_plan"
	productpriceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/product_price_plan"
)

// TableDeps holds the dependencies required to build the product-prices tab table.
// The parent caller constructs this from its own DetailViewDeps.
type TableDeps struct {
	// Route URLs for edit action and tab refresh.
	// PlanProductPriceEditURL is the URL template for the per-row edit action
	// (e.g. PriceScheduleRoutes.PlanProductPriceEditURL or PricePlanRoutes.ProductPriceEditURL).
	PlanProductPriceEditURL string
	// PlanTabActionURL is the URL template used to build the table refresh URL
	// (the tab action that re-renders this tab body on HTMX refresh).
	PlanTabActionURL string
	// ProductPricesTabSlug is the canonical slug for the product-prices tab
	// (pre-resolved by the parent from its labels, e.g. "product-prices").
	ProductPricesTabSlug string

	// Table ID — callers pass their own so different contexts get distinct IDs
	// (e.g. "price-schedule-plan-product-prices-table" vs "price-plan-product-prices-table").
	TableID string

	// Column labels sourced from the parent's schedule/plan label structs.
	ColumnProduct   string
	ColumnPrice     string
	ColumnCurrency  string
	ColumnTreatment string
	ColumnEffective string

	// Edit action label and empty-state strings.
	EditLabel    string
	EmptyTitle   string
	EmptyMsg     string
	Unauthorized string

	// Route path params — the parent passes the already-resolved IDs so the
	// table builder never needs to read URL path values itself.
	// sid = schedule/parent ID (first param in URL template, named "id").
	// ppid = price_plan ID (second param, named "ppid").
	SID  string
	PPID string

	// Use-case functions.
	ListProducts          func(ctx context.Context, req *productpb.ListProductsRequest) (*productpb.ListProductsResponse, error)
	ListProductPlans      func(ctx context.Context, req *productplanpb.ListProductPlansRequest) (*productplanpb.ListProductPlansResponse, error)
	ListProductPricePlans func(ctx context.Context, req *productpriceplanpb.ListProductPricePlansRequest) (*productpriceplanpb.ListProductPricePlansResponse, error)

	// Shared labels.
	PlanLabels             sib_subscription_price_plan.Labels
	ProductPricePlanLabels Labels
	TableLabels            types.TableLabels
}

// BuildTable assembles the TableConfig for the product-prices tab. The parent
// resolves and passes parent (ParentContext) so this function never needs to
// reach into the parent package.
//
// sid is the schedule/parent ID (URL param "id") and ppid is the price_plan ID
// (URL param "ppid"). Both are provided by the caller since they are URL-template
// inputs owned by the parent.
func BuildTable(ctx context.Context, deps *TableDeps, parent ParentContext) *types.TableConfig {
	perms := view.GetUserPermissions(ctx)
	showTreatment := parent.BillingKind != "BILLING_KIND_ONE_TIME"

	columns := []types.TableColumn{
		{Key: "product", Label: deps.ColumnProduct},
		{Key: "price", Label: deps.ColumnPrice, WidthClass: "col-4xl", Align: "right"},
		{Key: "currency", Label: deps.ColumnCurrency, NoSort: true, WidthClass: "col-2xl"},
	}
	if showTreatment {
		columns = append(columns, types.TableColumn{Key: "treatment", Label: deps.ColumnTreatment, NoSort: true, WidthClass: "col-3xl"})
	}
	columns = append(columns, types.TableColumn{Key: "effective", Label: deps.ColumnEffective, NoSort: true, WidthClass: "col-4xl"})

	productNames := map[string]string{}
	if deps.ListProducts != nil {
		prodResp, err := deps.ListProducts(ctx, &productpb.ListProductsRequest{})
		if err == nil {
			for _, p := range prodResp.GetData() {
				if p != nil {
					productNames[p.GetId()] = p.GetName()
				}
			}
		}
	}

	// Model D — build product_plan_id → (product_id, variant_id) map so we
	// resolve row display via the catalog line's FK.
	type productPlanRef struct {
		productID string
		variantID string
	}
	productPlans := map[string]productPlanRef{}
	if deps.ListProductPlans != nil {
		ppResp, err := deps.ListProductPlans(ctx, &productplanpb.ListProductPlansRequest{})
		if err == nil {
			for _, pp := range ppResp.GetData() {
				if pp == nil {
					continue
				}
				productPlans[pp.GetId()] = productPlanRef{
					productID: pp.GetProductId(),
					variantID: pp.GetProductVariantId(),
				}
			}
		}
	}

	refreshURL := route.ResolveURL(deps.PlanTabActionURL, "id", deps.SID, "ppid", deps.PPID, "tab", deps.ProductPricesTabSlug)
	rows := []types.TableRow{}
	if deps.ListProductPricePlans != nil {
		pppResp, err := deps.ListProductPricePlans(ctx, &productpriceplanpb.ListProductPricePlansRequest{})
		if err != nil {
			log.Printf("Failed to list product price plans: %v", err)
		} else {
			pplLabels := deps.ProductPricePlanLabels.Form
			for _, item := range pppResp.GetData() {
				if item == nil || item.GetPricePlanId() != deps.PPID {
					continue
				}
				itemID := item.GetId()
				ref := productPlans[item.GetProductPlanId()]
				if embed := item.GetProductPlan(); embed != nil {
					if pid := embed.GetProductId(); pid != "" {
						ref.productID = pid
					}
					if vid := embed.GetProductVariantId(); vid != "" {
						ref.variantID = vid
					}
				}
				productName := productNames[ref.productID]
				if productName == "" {
					productName = ref.productID
				}
				if ref.variantID != "" {
					productName = fmt.Sprintf("%s (%s)", productName, ref.variantID)
				}
				itemCurrency := item.GetBillingCurrency()
				if itemCurrency == "" {
					itemCurrency = "PHP"
				}
				priceCell := types.MoneyCell(float64(item.GetBillingAmount()), itemCurrency, true)
				cells := []types.TableCell{
					{Type: "text", Value: productName},
					priceCell,
					{Type: "text", Value: itemCurrency},
				}
				if showTreatment {
					cells = append(cells, types.TableCell{Type: "text", Value: BillingTreatmentDisplay(item.GetBillingTreatment().String(), pplLabels)})
				}
				cells = append(cells, types.TableCell{Type: "text", Value: EffectiveRangeDisplay(item.GetDateStart(), item.GetDateEnd())})
				rows = append(rows, types.TableRow{
					ID:    itemID,
					Cells: cells,
					// No delete action: rows are auto-seeded from product_plan assignments,
					// so deletion here would desync the two tables. Use the plan's
					// Products tab to remove the product_plan link, which in turn
					// should remove its product_price_plan rows.
					Actions: []types.TableAction{
						{
							Type:            "edit",
							Label:           deps.EditLabel,
							Action:          "edit",
							URL:             route.ResolveURL(deps.PlanProductPriceEditURL, "id", deps.SID, "ppid", deps.PPID, "pppid", itemID),
							DrawerTitle:     deps.EditLabel,
							Disabled:        !perms.Can("product_price_plan", "update"),
							DisabledTooltip: deps.Unauthorized,
						},
					},
				})
			}
		}
	}

	types.ApplyColumnStyles(columns, rows)

	cfg := &types.TableConfig{
		ID:                   deps.TableID,
		RefreshURL:           refreshURL,
		Columns:              columns,
		Rows:                 rows,
		ShowSearch:           true,
		ShowActions:          true,
		ShowSort:             true,
		ShowColumns:          true,
		ShowEntries:          true,
		DefaultSortColumn:    "product",
		DefaultSortDirection: "asc",
		Labels:               deps.TableLabels,
		EmptyState: types.TableEmptyState{
			Title:   deps.EmptyTitle,
			Message: deps.EmptyMsg,
		},
		// No PrimaryAction: product_price_plan rows are auto-seeded from
		// product_plan assignments when the parent PricePlan is created, so
		// manual Add is disabled here — users Edit existing rows instead.
	}
	types.ApplyTableSettings(cfg)
	return cfg
}

// ---------------------------------------------------------------------------
// Display helpers — PPP-only formatters
// ---------------------------------------------------------------------------

// FormatBillingKindLabel maps a proto enum string to its translated display label.
func FormatBillingKindLabel(kind string, l sib_subscription_price_plan.FormLabels) string {
	switch kind {
	case "BILLING_KIND_ONE_TIME":
		return l.BillingKindOneTime
	case "BILLING_KIND_RECURRING":
		return l.BillingKindRecurring
	case "BILLING_KIND_CONTRACT":
		return l.BillingKindContract
	}
	return kind
}

// FormatAmountBasisLabel maps a proto enum string to its translated display label.
func FormatAmountBasisLabel(basis string, l sib_subscription_price_plan.FormLabels) string {
	switch basis {
	case "AMOUNT_BASIS_PER_CYCLE":
		return l.AmountBasisPerCycle
	case "AMOUNT_BASIS_TOTAL_PACKAGE":
		return l.AmountBasisTotalPackage
	case "AMOUNT_BASIS_DERIVED_FROM_LINES":
		return l.AmountBasisDerivedFromLines
	}
	return basis
}

// BillingTreatmentDisplay maps the proto enum string to its human label.
// Returns "—" for unspecified so the table cell stays visually quiet.
func BillingTreatmentDisplay(value string, l FormLabels) string {
	switch value {
	case "BILLING_TREATMENT_RECURRING":
		return l.BillingTreatmentRecurring
	case "BILLING_TREATMENT_ONE_TIME_INITIAL":
		return l.BillingTreatmentOneTimeInitial
	case "BILLING_TREATMENT_USAGE_BASED":
		return l.BillingTreatmentUsageBased
	}
	return "—"
}

// EffectiveRangeDisplay renders the per-line effective dates as "start → end",
// "from start", "until end", or "Always" when both ends are empty.
func EffectiveRangeDisplay(start, end string) string {
	switch {
	case start == "" && end == "":
		return "Always"
	case start != "" && end == "":
		return "from " + start
	case start == "" && end != "":
		return "until " + end
	default:
		return start + " → " + end
	}
}
