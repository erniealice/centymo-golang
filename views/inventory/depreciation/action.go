// Package depreciation handles the depreciation feature for inventory items.
// Drawer template: depreciation-drawer-form.html (stays flat at view root).
package depreciation

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/erniealice/pyeza-golang/route"
	pyeza "github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	centymo "github.com/erniealice/centymo-golang"
	depreciationform "github.com/erniealice/centymo-golang/views/inventory/depreciation/form"

	inventorydepreciationpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/inventory_depreciation"
)

// Deps is the dependency subset needed by the depreciation feature.
type Deps struct {
	Routes centymo.InventoryRoutes
	Labels centymo.InventoryLabels

	CreateInventoryDepreciation func(ctx context.Context, req *inventorydepreciationpb.CreateInventoryDepreciationRequest) (*inventorydepreciationpb.CreateInventoryDepreciationResponse, error)
	ReadInventoryDepreciation   func(ctx context.Context, req *inventorydepreciationpb.ReadInventoryDepreciationRequest) (*inventorydepreciationpb.ReadInventoryDepreciationResponse, error)
	UpdateInventoryDepreciation func(ctx context.Context, req *inventorydepreciationpb.UpdateInventoryDepreciationRequest) (*inventorydepreciationpb.UpdateInventoryDepreciationResponse, error)
}

func formLabels(t func(string) string, d centymo.InventoryDepreciationLabels) depreciationform.Labels {
	return depreciationform.Labels{
		Method:       t("inventory.depreciation.method"),
		CostBasis:    t("inventory.depreciation.costBasis"),
		SalvageValue: t("inventory.depreciation.salvageValue"),
		UsefulLife:   t("inventory.depreciation.usefulLife"),
		StartDate:    t("inventory.depreciation.startDate"),
		// Info fields sourced from centymo.InventoryDepreciationLabels (populated from lyngua JSON + defaults).
		MethodInfo:       d.MethodInfo,
		CostBasisInfo:    d.CostBasisInfo,
		SalvageValueInfo: d.SalvageValueInfo,
		UsefulLifeInfo:   d.UsefulLifeInfo,
		StartDateInfo:    d.StartDateInfo,
	}
}

func methodOptions(t func(string) string) []pyeza.SelectOption {
	return []pyeza.SelectOption{
		{Value: "straight_line", Label: t("inventory.depreciation.methodStraightLine")},
		{Value: "declining_balance", Label: t("inventory.depreciation.methodDecliningBalance")},
		{Value: "sum_of_years", Label: t("inventory.depreciation.methodSumOfYears")},
	}
}

// NewAssignAction creates the depreciation configure action (GET = form, POST = create).
func NewAssignAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("inventory_item", "create") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		inventoryItemID := viewCtx.Request.PathValue("id")

		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("depreciation-drawer-form", &depreciationform.Data{
				FormAction:    route.ResolveURL(deps.Routes.DepreciationAssignURL, "id", inventoryItemID),
				Method:        "straight_line",
				Labels:        formLabels(viewCtx.T, deps.Labels.Depreciation),
				MethodOptions: methodOptions(viewCtx.T),
				CommonLabels:  nil,
			})
		}

		// POST - create depreciation record
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
		}

		r := viewCtx.Request
		costBasisF, _ := strconv.ParseFloat(r.FormValue("cost_basis"), 64)
		salvageValueF, _ := strconv.ParseFloat(r.FormValue("salvage_value"), 64)
		usefulLife, _ := strconv.ParseInt(r.FormValue("useful_life_months"), 10, 32)

		data := &inventorydepreciationpb.InventoryDepreciation{
			InventoryItemId:  inventoryItemID,
			Method:           r.FormValue("method"),
			CostBasis:        int64(math.Round(costBasisF * 100)),
			SalvageValue:     int64(math.Round(salvageValueF * 100)),
			UsefulLifeMonths: int32(usefulLife),
			StartDate:        r.FormValue("start_date"),
		}

		_, err := deps.CreateInventoryDepreciation(ctx, &inventorydepreciationpb.CreateInventoryDepreciationRequest{Data: data})
		if err != nil {
			log.Printf("Failed to create depreciation: %v", err)
			return centymo.HTMXError(err.Error())
		}

		// Redirect back to depreciation tab to show the new config
		return view.ViewResult{
			StatusCode: http.StatusOK,
			Headers: map[string]string{
				"HX-Trigger":  `{"formSuccess":true}`,
				"HX-Redirect": route.ResolveURL(deps.Routes.DetailURL, "id", inventoryItemID) + "?tab=depreciation",
			},
		}
	})
}

// NewEditAction creates the depreciation edit action (GET = form, POST = update).
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("inventory_item", "update") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		inventoryItemID := viewCtx.Request.PathValue("id")
		depreciationID := viewCtx.Request.PathValue("did")

		if viewCtx.Request.Method == http.MethodGet {
			resp, err := deps.ReadInventoryDepreciation(ctx, &inventorydepreciationpb.ReadInventoryDepreciationRequest{
				Data: &inventorydepreciationpb.InventoryDepreciation{Id: depreciationID},
			})
			if err != nil {
				log.Printf("Failed to read depreciation %s: %v", depreciationID, err)
				return centymo.HTMXError(deps.Labels.Errors.NotFound)
			}
			records := resp.GetData()
			if len(records) == 0 {
				return centymo.HTMXError(deps.Labels.Errors.NotFound)
			}
			record := records[0]

			return view.OK("depreciation-drawer-form", &depreciationform.Data{
				FormAction:    route.ResolveURL(deps.Routes.DepreciationEditURL, "id", inventoryItemID, "did", depreciationID),
				IsEdit:        true,
				ID:            depreciationID,
				Method:        record.GetMethod(),
				CostBasis:     fmt.Sprintf("%.2f", float64(record.GetCostBasis())/100.0),
				SalvageValue:  fmt.Sprintf("%.2f", float64(record.GetSalvageValue())/100.0),
				UsefulLife:    fmt.Sprintf("%d", record.GetUsefulLifeMonths()),
				StartDate:     record.GetStartDate(),
				Labels:        formLabels(viewCtx.T, deps.Labels.Depreciation),
				MethodOptions: methodOptions(viewCtx.T),
				CommonLabels:  nil,
			})
		}

		// POST - update depreciation record
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
		}

		r := viewCtx.Request
		costBasisF, _ := strconv.ParseFloat(r.FormValue("cost_basis"), 64)
		salvageValueF, _ := strconv.ParseFloat(r.FormValue("salvage_value"), 64)
		usefulLife, _ := strconv.ParseInt(r.FormValue("useful_life_months"), 10, 32)

		data := &inventorydepreciationpb.InventoryDepreciation{
			Id:               depreciationID,
			Method:           r.FormValue("method"),
			CostBasis:        int64(math.Round(costBasisF * 100)),
			SalvageValue:     int64(math.Round(salvageValueF * 100)),
			UsefulLifeMonths: int32(usefulLife),
			StartDate:        r.FormValue("start_date"),
		}

		_, err := deps.UpdateInventoryDepreciation(ctx, &inventorydepreciationpb.UpdateInventoryDepreciationRequest{Data: data})
		if err != nil {
			log.Printf("Failed to update depreciation %s: %v", depreciationID, err)
			return centymo.HTMXError(err.Error())
		}

		return view.ViewResult{
			StatusCode: http.StatusOK,
			Headers: map[string]string{
				"HX-Trigger":  `{"formSuccess":true}`,
				"HX-Redirect": route.ResolveURL(deps.Routes.DetailURL, "id", inventoryItemID) + "?tab=depreciation",
			},
		}
	})
}
