package action

import (
	"context"
	"log"
	"net/http"

	"github.com/erniealice/pyeza-golang/view"

	centymo "github.com/erniealice/centymo-golang"
)

// DepreciationFormLabels holds i18n labels for the depreciation drawer form.
type DepreciationFormLabels struct {
	Method       string
	CostBasis    string
	SalvageValue string
	UsefulLife   string
	StartDate    string
}

// DepreciationFormData is the template data for the depreciation drawer form.
type DepreciationFormData struct {
	FormAction    string
	IsEdit        bool
	ID            string
	Method        string
	CostBasis     string
	SalvageValue  string
	UsefulLife    string
	StartDate     string
	Labels        DepreciationFormLabels
	MethodOptions []SelectOption
	CommonLabels  any
}

func depreciationFormLabels(t func(string) string) DepreciationFormLabels {
	return DepreciationFormLabels{
		Method:       t("inventory.depreciation.method"),
		CostBasis:    t("inventory.depreciation.costBasis"),
		SalvageValue: t("inventory.depreciation.salvageValue"),
		UsefulLife:   t("inventory.depreciation.usefulLife"),
		StartDate:    t("inventory.depreciation.startDate"),
	}
}

func depreciationMethodOptions(t func(string) string) []SelectOption {
	return []SelectOption{
		{Value: "straight_line", Label: t("inventory.depreciation.methodStraightLine")},
		{Value: "declining_balance", Label: t("inventory.depreciation.methodDecliningBalance")},
		{Value: "sum_of_years", Label: t("inventory.depreciation.methodSumOfYears")},
	}
}

// NewDepreciationAssignAction creates the depreciation configure action (GET = form, POST = create).
func NewDepreciationAssignAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		inventoryItemID := viewCtx.Request.PathValue("id")

		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("depreciation-drawer-form", &DepreciationFormData{
				FormAction:    "/action/inventory/detail/" + inventoryItemID + "/depreciation/assign",
				Method:        "straight_line",
				Labels:        depreciationFormLabels(viewCtx.T),
				MethodOptions: depreciationMethodOptions(viewCtx.T),
				CommonLabels:  nil,
			})
		}

		// POST - create depreciation record
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError("Invalid form data")
		}

		r := viewCtx.Request
		data := map[string]any{
			"inventory_item_id": inventoryItemID,
			"method":            r.FormValue("method"),
			"cost_basis":        r.FormValue("cost_basis"),
			"salvage_value":     r.FormValue("salvage_value"),
			"useful_life_months": r.FormValue("useful_life_months"),
			"start_date":        r.FormValue("start_date"),
		}

		_, err := deps.DB.Create(ctx, "inventory_depreciation", data)
		if err != nil {
			log.Printf("Failed to create depreciation: %v", err)
			return centymo.HTMXError("Failed to configure depreciation")
		}

		// Redirect back to depreciation tab to show the new config
		return view.ViewResult{
			StatusCode: http.StatusOK,
			Headers: map[string]string{
				"HX-Trigger":  `{"formSuccess":true}`,
				"HX-Redirect": "/app/inventory/detail/" + inventoryItemID + "?tab=depreciation",
			},
		}
	})
}

// NewDepreciationEditAction creates the depreciation edit action (GET = form, POST = update).
func NewDepreciationEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		inventoryItemID := viewCtx.Request.PathValue("id")
		depreciationID := viewCtx.Request.PathValue("did")

		if viewCtx.Request.Method == http.MethodGet {
			record, err := deps.DB.Read(ctx, "inventory_depreciation", depreciationID)
			if err != nil {
				log.Printf("Failed to read depreciation %s: %v", depreciationID, err)
				return centymo.HTMXError("Depreciation record not found")
			}

			method, _ := record["method"].(string)
			costBasis := anyToString(record["cost_basis"])
			salvageValue := anyToString(record["salvage_value"])
			usefulLife := anyToString(record["useful_life_months"])
			startDate, _ := record["start_date"].(string)

			return view.OK("depreciation-drawer-form", &DepreciationFormData{
				FormAction:    "/action/inventory/detail/" + inventoryItemID + "/depreciation/edit/" + depreciationID,
				IsEdit:        true,
				ID:            depreciationID,
				Method:        method,
				CostBasis:     costBasis,
				SalvageValue:  salvageValue,
				UsefulLife:    usefulLife,
				StartDate:     startDate,
				Labels:        depreciationFormLabels(viewCtx.T),
				MethodOptions: depreciationMethodOptions(viewCtx.T),
				CommonLabels:  nil,
			})
		}

		// POST - update depreciation record
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError("Invalid form data")
		}

		r := viewCtx.Request
		data := map[string]any{
			"method":            r.FormValue("method"),
			"cost_basis":        r.FormValue("cost_basis"),
			"salvage_value":     r.FormValue("salvage_value"),
			"useful_life_months": r.FormValue("useful_life_months"),
			"start_date":        r.FormValue("start_date"),
		}

		_, err := deps.DB.Update(ctx, "inventory_depreciation", depreciationID, data)
		if err != nil {
			log.Printf("Failed to update depreciation %s: %v", depreciationID, err)
			return centymo.HTMXError("Failed to update depreciation")
		}

		return view.ViewResult{
			StatusCode: http.StatusOK,
			Headers: map[string]string{
				"HX-Trigger":  `{"formSuccess":true}`,
				"HX-Redirect": "/app/inventory/detail/" + inventoryItemID + "?tab=depreciation",
			},
		}
	})
}
