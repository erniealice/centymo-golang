package action

import (
	"context"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"

	centymo "github.com/erniealice/centymo-golang"

	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	locationpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/location"
	productplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_plan"
	priceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_plan"
	productpriceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/product_price_plan"
)

// LocationOption is a minimal struct for rendering location options in the price plan form.
type LocationOption struct {
	Id   string
	Name string
}

// ProductPlanPriceRow represents a product-plan entry with an optional price for the product pricing section.
type ProductPlanPriceRow struct {
	ProductPlanID string
	ProductName   string
	Price         string // display value (centavos ÷100), empty on add
}

// PricePlanFormLabels holds i18n labels for the price plan drawer form template.
type PricePlanFormLabels struct {
	Name                string
	NamePlaceholder     string
	Description         string
	DescPlaceholder     string
	Amount              string
	AmountPlaceholder   string
	Currency            string
	CurrencyPlaceholder string
	DurationValue       string
	DurationUnit        string
	Location            string
	LocationPlaceholder string
	SelectLocation      string
	Active              string
}

// PricePlanFormData is the template data for the price plan drawer form.
type PricePlanFormData struct {
	FormAction           string
	IsEdit               bool
	ID                   string
	PlanID               string
	Name                 string
	Description          string
	Amount               string
	Currency             string
	DurationValue        string
	DurationUnit         string
	Active               bool
	Locations            []*LocationOption
	SelectedLocationID   string
	SelectedLocationLabel string
	LocationOptions      []map[string]any
	ProductPlans         []ProductPlanPriceRow
	Labels               PricePlanFormLabels
	CommonLabels         any
}

// PricePlanDeps holds dependencies for price plan action handlers.
type PricePlanDeps struct {
	Routes                centymo.PlanRoutes
	Labels                centymo.PlanLabels
	CreatePricePlan       func(ctx context.Context, req *priceplanpb.CreatePricePlanRequest) (*priceplanpb.CreatePricePlanResponse, error)
	ReadPricePlan         func(ctx context.Context, req *priceplanpb.ReadPricePlanRequest) (*priceplanpb.ReadPricePlanResponse, error)
	UpdatePricePlan       func(ctx context.Context, req *priceplanpb.UpdatePricePlanRequest) (*priceplanpb.UpdatePricePlanResponse, error)
	DeletePricePlan       func(ctx context.Context, req *priceplanpb.DeletePricePlanRequest) (*priceplanpb.DeletePricePlanResponse, error)
	ListLocations         func(ctx context.Context, req *locationpb.ListLocationsRequest) (*locationpb.ListLocationsResponse, error)
	ListProductPlans      func(ctx context.Context, req *productplanpb.ListProductPlansRequest) (*productplanpb.ListProductPlansResponse, error)
	CreateProductPricePlan func(ctx context.Context, req *productpriceplanpb.CreateProductPricePlanRequest) (*productpriceplanpb.CreateProductPricePlanResponse, error)
	ListProductPricePlans  func(ctx context.Context, req *productpriceplanpb.ListProductPricePlansRequest) (*productpriceplanpb.ListProductPricePlansResponse, error)
}

// pricePlanFormLabels converts centymo.PricePlanFormLabels into the local type.
func pricePlanFormLabels(l centymo.PricePlanFormLabels) PricePlanFormLabels {
	return PricePlanFormLabels{
		Name:                l.Name,
		NamePlaceholder:     l.NamePlaceholder,
		Description:         l.Description,
		DescPlaceholder:     l.DescPlaceholder,
		Amount:              l.Amount,
		AmountPlaceholder:   l.AmountPlaceholder,
		Currency:            l.Currency,
		CurrencyPlaceholder: l.CurrencyPlaceholder,
		DurationValue:       l.DurationValue,
		DurationUnit:        l.DurationUnit,
		Location:            l.Location,
		LocationPlaceholder: l.LocationPlaceholder,
		SelectLocation:      l.SelectLocation,
		Active:              l.Active,
	}
}

// loadLocationOptions fetches the location list and converts to options.
// Returns nil slice on error (graceful degradation).
func loadLocationOptions(ctx context.Context, deps *PricePlanDeps) []*LocationOption {
	if deps.ListLocations == nil {
		return nil
	}
	resp, err := deps.ListLocations(ctx, &locationpb.ListLocationsRequest{})
	if err != nil {
		log.Printf("Failed to load locations for price plan form: %v", err)
		return nil
	}
	var options []*LocationOption
	for _, loc := range resp.GetData() {
		options = append(options, &LocationOption{
			Id:   loc.GetId(),
			Name: loc.GetName(),
		})
	}
	return options
}

// buildLocationAutoCompleteOptions converts []*LocationOption to the auto-complete format.
func buildLocationAutoCompleteOptions(locations []*LocationOption, selectedID string) []map[string]any {
	opts := make([]map[string]any, 0, len(locations))
	for _, loc := range locations {
		opts = append(opts, map[string]any{
			"Value":    loc.Id,
			"Label":    loc.Name,
			"Selected": loc.Id == selectedID,
		})
	}
	return opts
}

// findLocationLabel returns the name of the location with the given ID, or empty string.
func findLocationLabel(locations []*LocationOption, id string) string {
	for _, loc := range locations {
		if loc.Id == id {
			return loc.Name
		}
	}
	return ""
}

// loadProductPlansForPlan loads all product_plan records for a given plan and returns price rows.
// If pricePlanID is non-empty, also loads existing product_price_plan records to pre-fill prices.
func loadProductPlansForPlan(ctx context.Context, deps *PricePlanDeps, planID string, pricePlanID string) []ProductPlanPriceRow {
	if deps.ListProductPlans == nil {
		return nil
	}
	resp, err := deps.ListProductPlans(ctx, &productplanpb.ListProductPlansRequest{
		Filters: &commonpb.FilterRequest{
			Logic: commonpb.FilterLogic_AND,
			Filters: []*commonpb.TypedFilter{
				{
					Field: "plan_id",
					FilterType: &commonpb.TypedFilter_StringFilter{
						StringFilter: &commonpb.StringFilter{
							Value:    planID,
							Operator: commonpb.StringOperator_STRING_EQUALS,
						},
					},
				},
			},
		},
	})
	if err != nil {
		log.Printf("Failed to load product plans for plan %s: %v", planID, err)
		return nil
	}

	// Build a map of product_plan_id → existing price (centavos) when editing
	existingPrices := map[string]int64{}
	if pricePlanID != "" && deps.ListProductPricePlans != nil {
		ppResp, err := deps.ListProductPricePlans(ctx, &productpriceplanpb.ListProductPricePlansRequest{
			Filters: &commonpb.FilterRequest{
				Logic: commonpb.FilterLogic_AND,
				Filters: []*commonpb.TypedFilter{
					{
						Field: "price_plan_id",
						FilterType: &commonpb.TypedFilter_StringFilter{
							StringFilter: &commonpb.StringFilter{
								Value:    pricePlanID,
								Operator: commonpb.StringOperator_STRING_EQUALS,
							},
						},
					},
				},
			},
		})
		if err == nil {
			for _, ppp := range ppResp.GetData() {
				// Key by product_id — map to price
				existingPrices[ppp.GetProductId()] = ppp.GetPrice()
			}
		}
	}

	rows := make([]ProductPlanPriceRow, 0, len(resp.GetData()))
	for _, pp := range resp.GetData() {
		priceStr := ""
		if price, ok := existingPrices[pp.GetProductId()]; ok && price > 0 {
			priceStr = strconv.FormatFloat(float64(price)/100.0, 'f', 2, 64)
		}
		rows = append(rows, ProductPlanPriceRow{
			ProductPlanID: pp.GetId(),
			ProductName:   pp.GetName(),
			Price:         priceStr,
		})
	}
	return rows
}

// NewPricePlanAddAction creates the price plan add action (GET = form, POST = create).
// URL: /action/plans/{id}/pricelists/add
func NewPricePlanAddAction(deps *PricePlanDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("price_plan", "create") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		planID := viewCtx.Request.PathValue("id")

		if viewCtx.Request.Method == http.MethodGet {
			locations := loadLocationOptions(ctx, deps)
			return view.OK("price-plan-drawer-form", &PricePlanFormData{
				FormAction:      route.ResolveURL(deps.Routes.PricePlanAddURL, "id", planID),
				PlanID:          planID,
				Active:          true,
				Currency:        "PHP",
				DurationUnit:    "months",
				Locations:       locations,
				LocationOptions: buildLocationAutoCompleteOptions(locations, ""),
				ProductPlans:    loadProductPlansForPlan(ctx, deps, planID, ""),
				Labels:          pricePlanFormLabels(deps.Labels.PricePlanForm),
				CommonLabels:    nil, // injected by ViewAdapter
			})
		}

		// POST — create price plan
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
		}

		r := viewCtx.Request
		active := r.FormValue("active") == "true"

		amount := int64(0)
		if v, err := strconv.ParseFloat(r.FormValue("amount"), 64); err == nil {
			amount = int64(math.Round(v * 100))
		}

		durationValue := int32(0)
		if v, err := strconv.ParseInt(r.FormValue("duration_value"), 10, 32); err == nil {
			durationValue = int32(v)
		}

		currency := r.FormValue("currency")

		pp := &priceplanpb.PricePlan{
			PlanId:        planID,
			Name:          r.FormValue("name"),
			Description:   r.FormValue("description"),
			Amount:        amount,
			Currency:      currency,
			DurationValue: durationValue,
			DurationUnit:  r.FormValue("duration_unit"),
			Active:        active,
		}
		if locID := r.FormValue("location_id"); locID != "" {
			pp.LocationId = &locID
		}

		createResp, err := deps.CreatePricePlan(ctx, &priceplanpb.CreatePricePlanRequest{
			Data: pp,
		})
		if err != nil {
			log.Printf("Failed to create price plan for plan %s: %v", planID, err)
			return centymo.HTMXError(err.Error())
		}

		// Get the new price_plan_id from the response
		newPricePlanID := ""
		if createResp != nil && len(createResp.GetData()) > 0 {
			newPricePlanID = createResp.GetData()[0].GetId()
		}

		// Create product_price_plan records for each product_prices[xxx] form value
		if newPricePlanID != "" && deps.CreateProductPricePlan != nil {
			productPlans := loadProductPlansForPlan(ctx, deps, planID, "")
			// Build a map of product_plan_id → product_id for lookup
			ppIDToProductID := map[string]string{}
			for _, row := range productPlans {
				// We need the product_id, not just the name — re-load product plans to get it
				_ = row
			}
			// Re-load to get product_id per product_plan_id
			if deps.ListProductPlans != nil {
				ppResp, err := deps.ListProductPlans(ctx, &productplanpb.ListProductPlansRequest{
					Filters: &commonpb.FilterRequest{
						Logic: commonpb.FilterLogic_AND,
						Filters: []*commonpb.TypedFilter{
							{
								Field: "plan_id",
								FilterType: &commonpb.TypedFilter_StringFilter{
									StringFilter: &commonpb.StringFilter{
										Value:    planID,
										Operator: commonpb.StringOperator_STRING_EQUALS,
									},
								},
							},
						},
					},
				})
				if err == nil {
					for _, pp := range ppResp.GetData() {
						ppIDToProductID[pp.GetId()] = pp.GetProductId()
					}
				}
			}

			for key, values := range r.Form {
				if !strings.HasPrefix(key, "product_prices[") {
					continue
				}
				if len(values) == 0 || values[0] == "" {
					continue
				}
				// Extract product_plan_id from key: product_prices[{id}]
				trimmed := strings.TrimPrefix(key, "product_prices[")
				trimmed = strings.TrimSuffix(trimmed, "]")
				productPlanID := trimmed

				productID, ok := ppIDToProductID[productPlanID]
				if !ok || productID == "" {
					continue
				}

				priceVal := int64(0)
				if v, err := strconv.ParseFloat(values[0], 64); err == nil {
					priceVal = int64(math.Round(v * 100))
				}

				_, err := deps.CreateProductPricePlan(ctx, &productpriceplanpb.CreateProductPricePlanRequest{
					Data: &productpriceplanpb.ProductPricePlan{
						PricePlanId: newPricePlanID,
						ProductId:   productID,
						Price:       priceVal,
						Currency:    currency,
					},
				})
				if err != nil {
					log.Printf("Failed to create product_price_plan for product_plan %s: %v", productPlanID, err)
				}
			}
		}

		return centymo.HTMXSuccess("plan-pricelists-table")
	})
}

// NewPricePlanEditAction creates the price plan edit action (GET = form, POST = update).
// URL: /action/plans/{id}/pricelists/edit/{ppid}
func NewPricePlanEditAction(deps *PricePlanDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("price_plan", "update") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		planID := viewCtx.Request.PathValue("id")
		ppID := viewCtx.Request.PathValue("ppid")

		if viewCtx.Request.Method == http.MethodGet {
			resp, err := deps.ReadPricePlan(ctx, &priceplanpb.ReadPricePlanRequest{
				Data: &priceplanpb.PricePlan{Id: ppID},
			})
			if err != nil {
				log.Printf("Failed to read price plan %s: %v", ppID, err)
				return centymo.HTMXError(deps.Labels.Errors.NotFound)
			}
			data := resp.GetData()
			if len(data) == 0 {
				return centymo.HTMXError(deps.Labels.Errors.NotFound)
			}
			pp := data[0]

			amountStr := strconv.FormatFloat(float64(pp.GetAmount())/100.0, 'f', 2, 64)
			durationStr := strconv.FormatInt(int64(pp.GetDurationValue()), 10)
			selectedLocationID := pp.GetLocationId()
			locations := loadLocationOptions(ctx, deps)

			return view.OK("price-plan-drawer-form", &PricePlanFormData{
				FormAction:            route.ResolveURL(deps.Routes.PricePlanEditURL, "id", planID, "ppid", ppID),
				IsEdit:                true,
				ID:                    ppID,
				PlanID:                planID,
				Name:                  pp.GetName(),
				Description:           pp.GetDescription(),
				Amount:                amountStr,
				Currency:              pp.GetCurrency(),
				DurationValue:         durationStr,
				DurationUnit:          pp.GetDurationUnit(),
				Active:                pp.GetActive(),
				Locations:             locations,
				SelectedLocationID:    selectedLocationID,
				SelectedLocationLabel: findLocationLabel(locations, selectedLocationID),
				LocationOptions:       buildLocationAutoCompleteOptions(locations, selectedLocationID),
				ProductPlans:          loadProductPlansForPlan(ctx, deps, planID, ppID),
				Labels:                pricePlanFormLabels(deps.Labels.PricePlanForm),
				CommonLabels:          nil, // injected by ViewAdapter
			})
		}

		// POST — update price plan
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
		}

		r := viewCtx.Request
		active := r.FormValue("active") == "true"

		amount := int64(0)
		if v, err := strconv.ParseFloat(r.FormValue("amount"), 64); err == nil {
			amount = int64(math.Round(v * 100))
		}

		durationValue := int32(0)
		if v, err := strconv.ParseInt(r.FormValue("duration_value"), 10, 32); err == nil {
			durationValue = int32(v)
		}

		currency := r.FormValue("currency")

		pp := &priceplanpb.PricePlan{
			Id:            ppID,
			PlanId:        planID,
			Name:          r.FormValue("name"),
			Description:   r.FormValue("description"),
			Amount:        amount,
			Currency:      currency,
			DurationValue: durationValue,
			DurationUnit:  r.FormValue("duration_unit"),
			Active:        active,
		}
		if locID := r.FormValue("location_id"); locID != "" {
			pp.LocationId = &locID
		}

		_, err := deps.UpdatePricePlan(ctx, &priceplanpb.UpdatePricePlanRequest{
			Data: pp,
		})
		if err != nil {
			log.Printf("Failed to update price plan %s: %v", ppID, err)
			return centymo.HTMXError(err.Error())
		}

		// Update product_price_plan records: delete+recreate pattern
		if deps.CreateProductPricePlan != nil {
			// Build product_plan_id → product_id lookup
			ppIDToProductID := map[string]string{}
			if deps.ListProductPlans != nil {
				ppResp, err := deps.ListProductPlans(ctx, &productplanpb.ListProductPlansRequest{
					Filters: &commonpb.FilterRequest{
						Logic: commonpb.FilterLogic_AND,
						Filters: []*commonpb.TypedFilter{
							{
								Field: "plan_id",
								FilterType: &commonpb.TypedFilter_StringFilter{
									StringFilter: &commonpb.StringFilter{
										Value:    planID,
										Operator: commonpb.StringOperator_STRING_EQUALS,
									},
								},
							},
						},
					},
				})
				if err == nil {
					for _, pp := range ppResp.GetData() {
						ppIDToProductID[pp.GetId()] = pp.GetProductId()
					}
				}
			}

			for key, values := range r.Form {
				if !strings.HasPrefix(key, "product_prices[") {
					continue
				}
				if len(values) == 0 || values[0] == "" {
					continue
				}
				trimmed := strings.TrimPrefix(key, "product_prices[")
				trimmed = strings.TrimSuffix(trimmed, "]")
				productPlanID := trimmed

				productID, ok := ppIDToProductID[productPlanID]
				if !ok || productID == "" {
					continue
				}

				priceVal := int64(0)
				if v, err := strconv.ParseFloat(values[0], 64); err == nil {
					priceVal = int64(math.Round(v * 100))
				}

				_, err := deps.CreateProductPricePlan(ctx, &productpriceplanpb.CreateProductPricePlanRequest{
					Data: &productpriceplanpb.ProductPricePlan{
						PricePlanId: ppID,
						ProductId:   productID,
						Price:       priceVal,
						Currency:    currency,
					},
				})
				if err != nil {
					log.Printf("Failed to upsert product_price_plan for product_plan %s: %v", productPlanID, err)
				}
			}
		}

		return centymo.HTMXSuccess("plan-pricelists-table")
	})
}

// NewPricePlanDeleteAction creates the price plan delete action (POST only).
// URL: /action/plans/{id}/pricelists/delete  (id=price_plan_id via query param)
func NewPricePlanDeleteAction(deps *PricePlanDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("price_plan", "delete") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		ppID := viewCtx.Request.URL.Query().Get("id")
		if ppID == "" {
			_ = viewCtx.Request.ParseForm()
			ppID = viewCtx.Request.FormValue("id")
		}
		if ppID == "" {
			return centymo.HTMXError(deps.Labels.Errors.IDRequired)
		}

		_, err := deps.DeletePricePlan(ctx, &priceplanpb.DeletePricePlanRequest{
			Data: &priceplanpb.PricePlan{Id: ppID},
		})
		if err != nil {
			log.Printf("Failed to delete price plan %s: %v", ppID, err)
			return centymo.HTMXError(err.Error())
		}

		return centymo.HTMXSuccess("plan-pricelists-table")
	})
}
