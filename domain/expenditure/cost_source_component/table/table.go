// Package table builds the recoverable-costs table of the expenditure detail
// "Recoverable costs" tab and serves its HTMX refresh view. The builder lives
// here (not in expenditure/detail) so the tab and the refresh route render the
// same table.
package table

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	csc "github.com/erniealice/centymo-golang/domain/expenditure/cost_source_component"
	costsourcecomponentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/cost_source_component"
)

// TableID is the DOM id of the table card (refreshTable trigger target).
const TableID = "cost-source-components-table"

// Input carries everything Build needs.
type Input struct {
	ExpenditureID string
	Currency      string
	Labels        csc.Labels
	Routes        csc.Routes
	TableLabels   types.TableLabels
	// AllocateURL / ViewAllocationURL are the allocation_batch drawer routes with
	// an {id} = cost source component id token; empty hides the action.
	AllocateURL       string
	ViewAllocationURL string
	Perms             *types.UserPermissions
	NoPermission      string
	// MissingPermission is the printf template ("Missing permission: %s") of the per-code tooltip
	// used by the cross-entity allocation actions (empty falls back to NoPermission).
	MissingPermission string
}

// missingTip is the disabled tooltip naming the missing permission code.
func missingTip(in Input, code string) string {
	if in.MissingPermission == "" {
		return in.NoPermission
	}
	return fmt.Sprintf(in.MissingPermission, code)
}

// KindLabel returns the display label of a component kind.
func KindLabel(l csc.Labels, k costsourcecomponentpb.CostSourceComponentKind) string {
	switch k {
	case costsourcecomponentpb.CostSourceComponentKind_COST_SOURCE_COMPONENT_KIND_ENERGY:
		return l.Enums.ComponentKindEnergy
	case costsourcecomponentpb.CostSourceComponentKind_COST_SOURCE_COMPONENT_KIND_WATER:
		return l.Enums.ComponentKindWater
	case costsourcecomponentpb.CostSourceComponentKind_COST_SOURCE_COMPONENT_KIND_SERVICE_FEE:
		return l.Enums.ComponentKindServiceFee
	case costsourcecomponentpb.CostSourceComponentKind_COST_SOURCE_COMPONENT_KIND_OTHER:
		return l.Enums.ComponentKindOther
	}
	return ""
}

// TaxFactLabel returns the display label of a tax fact ("" when unset).
func TaxFactLabel(l csc.Labels, f costsourcecomponentpb.CostTaxFact) string {
	switch f {
	case costsourcecomponentpb.CostTaxFact_COST_TAX_FACT_VATABLE:
		return l.Enums.TaxFactVatable
	case costsourcecomponentpb.CostTaxFact_COST_TAX_FACT_EXEMPT:
		return l.Enums.TaxFactExempt
	case costsourcecomponentpb.CostTaxFact_COST_TAX_FACT_ZERO_RATED:
		return l.Enums.TaxFactZeroRated
	case costsourcecomponentpb.CostTaxFact_COST_TAX_FACT_NOT_APPLICABLE:
		return l.Enums.TaxFactNotApplicable
	}
	return ""
}

// ServicePeriod renders "{from} to {to}" through the detail.service_period label.
func ServicePeriod(l csc.Labels, from, to string) string {
	if from == "" && to == "" {
		return "—"
	}
	return strings.NewReplacer("{0}", from, "{1}", to).Replace(l.Detail.ServicePeriod)
}

// Claimed reports whether a component is claimed by an allocation or a recognition.
func Claimed(c *costsourcecomponentpb.CostSourceComponent) bool {
	return c.GetClaimKind() != costsourcecomponentpb.SourceClaimKind_SOURCE_CLAIM_KIND_UNSPECIFIED && c.ClaimKind != nil
}

func claimCell(l csc.Labels, c *costsourcecomponentpb.CostSourceComponent) types.TableCell {
	switch c.GetClaimKind() {
	case costsourcecomponentpb.SourceClaimKind_SOURCE_CLAIM_KIND_ALLOCATION:
		return types.TableCell{Type: "badge", Value: l.Enums.ClaimKindAllocation, Variant: "success"}
	case costsourcecomponentpb.SourceClaimKind_SOURCE_CLAIM_KIND_RECOGNITION:
		return types.TableCell{Type: "badge", Value: l.Enums.ClaimKindRecognition, Variant: "info"}
	}
	return types.TableCell{Type: "text", Value: l.Enums.ClaimNone}
}

// Build builds the recoverable-costs table config.
func Build(rows []*costsourcecomponentpb.CostSourceComponent, in Input) *types.TableConfig {
	l := in.Labels
	columns := []types.TableColumn{
		{Key: "component_kind", Label: l.Columns.ComponentKind, NoSort: true, WidthClass: "col-xl"},
		{Key: "description", Label: l.Columns.Description, NoSort: true},
		{Key: "basis", Label: l.Columns.Basis, NoSort: true, WidthClass: "col-xl"},
		{Key: "amount", Label: l.Columns.Amount, NoSort: true, WidthClass: "col-3xl"},
		{Key: "tax_fact", Label: l.Columns.TaxFact, NoSort: true, WidthClass: "col-xl"},
		{Key: "service_period", Label: l.Columns.ServicePeriod, NoSort: true, WidthClass: "col-3xl"},
		{Key: "claim", Label: l.Columns.Claim, NoSort: true, WidthClass: "col-xl"},
	}
	canUpdate := in.Perms.Can("expenditure", "update")
	// The allocation drawers belong to allocation_batch, not expenditure (cross-entity rule): the
	// handler requires create||update to allocate and read to view.
	canAllocate := in.Perms.Can("allocation_batch", "create") || in.Perms.Can("allocation_batch", "update")
	canViewAllocation := in.Perms.Can("allocation_batch", "read")
	tableRows := make([]types.TableRow, 0, len(rows))
	for _, c := range rows {
		id := c.GetId()
		claimed := Claimed(c)

		basis := "—"
		if c.BasisQuantityScaled != nil {
			basis = csc.FormatQuantity(c.GetBasisQuantityScaled(), c.GetBasisScale())
			if u := c.GetBasisUnit(); u != "" {
				basis += " " + u
			}
		}
		tax := TaxFactLabel(l, c.GetTaxFact())
		if tax == "" {
			tax = "—"
		}

		editTip, deleteTip := in.NoPermission, in.NoPermission
		if canUpdate && claimed {
			editTip, deleteTip = l.Detail.ClaimedNotice, l.Detail.ClaimedNotice
		}
		actions := []types.TableAction{
			{
				Type: "edit", Label: l.Buttons.Edit, Action: "edit", TestID: "cost-source-component-edit-" + id,
				URL: route.ResolveURL(in.Routes.EditURL, "id", id), DrawerTitle: l.Buttons.Edit,
				Disabled: !canUpdate || claimed, DisabledTooltip: editTip,
			},
			{
				Type: "delete", Label: l.Buttons.Delete, Action: "delete", TestID: "cost-source-component-delete-" + id,
				URL: route.ResolveURL(in.Routes.DeleteURL, "id", id), ItemName: c.GetDescription(),
				ConfirmTitle: l.Confirm.DeleteTitle, ConfirmMessage: l.Confirm.DeleteMsg,
				Disabled: !canUpdate || claimed, DisabledTooltip: deleteTip,
			},
		}
		switch {
		case claimed && in.ViewAllocationURL != "" && c.GetClaimKind() == costsourcecomponentpb.SourceClaimKind_SOURCE_CLAIM_KIND_ALLOCATION:
			actions = append(actions, types.TableAction{
				Type: "view", Label: l.Buttons.ViewAllocation, Action: "view-allocation", TestID: "cost-source-component-view-allocation-" + id,
				HxGet: route.ResolveURL(in.ViewAllocationURL, "id", id), HxTarget: "#sheetContent", HxSwap: "innerHTML",
				DrawerTitle: l.Buttons.ViewAllocation,
				Disabled:    !canViewAllocation, DisabledTooltip: missingTip(in, "allocation_batch:read"),
			})
		case !claimed && in.AllocateURL != "":
			actions = append(actions, types.TableAction{
				Type: "edit", Label: l.Buttons.Allocate, Action: "allocate", TestID: "cost-source-component-allocate-" + id,
				HxGet: route.ResolveURL(in.AllocateURL, "id", id), HxTarget: "#sheetContent", HxSwap: "innerHTML",
				DrawerTitle: l.Buttons.Allocate, Disabled: !canAllocate, DisabledTooltip: missingTip(in, "allocation_batch:create"),
			})
		}

		tableRows = append(tableRows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: KindLabel(l, c.GetComponentKind())},
				{Type: "text", Value: c.GetDescription()},
				{Type: "text", Value: basis},
				types.MoneyCell(float64(c.GetAmount()), firstNonEmpty(c.GetCurrency(), in.Currency), true),
				{Type: "text", Value: tax},
				{Type: "text", Value: ServicePeriod(l, c.GetServiceFrom(), c.GetServiceTo())},
				claimCell(l, c),
			},
			Actions: actions,
		})
	}
	types.ApplyColumnStyles(columns, tableRows)

	return &types.TableConfig{
		ID:         TableID,
		Columns:    columns,
		Rows:       tableRows,
		Labels:     in.TableLabels,
		RefreshURL: route.ResolveURL(in.Routes.TableURL, "id", in.ExpenditureID),
		EmptyState: types.TableEmptyState{Title: l.Empty.Title, Message: l.Empty.Message},
	}
}

func firstNonEmpty(a, b string) string {
	if a != "" {
		return a
	}
	return b
}

// Deps holds the dependencies of the refresh view.
type Deps struct {
	Routes            csc.Routes
	Labels            csc.Labels
	CommonLabels      pyeza.CommonLabels
	TableLabels       types.TableLabels
	UseCases          *csc.UseCases
	AllocateURL       string
	ViewAllocationURL string
	// ReadCurrency returns the expenditure currency for rows that carry none (optional).
	ReadCurrency func(ctx context.Context, expenditureID string) string
}

// NewView serves the table card fragment (HTMX refreshTable target).
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("expenditure", "read") {
			return view.Forbidden("expenditure:read")
		}
		if deps.UseCases == nil || deps.UseCases.ListCostSourceComponents == nil {
			return view.ViewResult{StatusCode: http.StatusServiceUnavailable, Headers: map[string]string{"HX-Error-Message": deps.Labels.Errors.Unavailable}}
		}
		expenditureID := viewCtx.Request.PathValue("id")
		rows, err := deps.UseCases.ListForExpenditure(ctx, expenditureID)
		if err != nil {
			log.Printf("cost_source_component table: list for expenditure %s: %v", expenditureID, err)
			return view.HTMXError(deps.Labels.Errors.Generic)
		}
		currency := ""
		if deps.ReadCurrency != nil {
			currency = deps.ReadCurrency(ctx, expenditureID)
		}
		return view.OK("table-card", Build(rows, Input{
			ExpenditureID: expenditureID, Currency: currency, Labels: deps.Labels, Routes: deps.Routes,
			TableLabels: deps.TableLabels, AllocateURL: deps.AllocateURL, ViewAllocationURL: deps.ViewAllocationURL,
			Perms: perms, NoPermission: deps.CommonLabels.Errors.PermissionDenied, MissingPermission: deps.CommonLabels.Errors.MissingPermission,
		}))
	})
}
