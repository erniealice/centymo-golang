// Package action holds the recoverable cost line (cost_source_component) write
// handlers: add, edit and delete drawers/actions of the expenditure detail
// "Recoverable costs" tab.
//
// Every handler re-checks its permission first (fail closed, layer 2). The
// use cases stay the authoritative gate; claimed lines are refused there
// (`claimed`) and the table disables their edit/delete affordances.
package action

import (
	"context"
	"log"
	"net/http"
	"strings"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"

	csc "github.com/erniealice/centymo-golang/domain/expenditure/cost_source_component"
	"github.com/erniealice/centymo-golang/domain/expenditure/cost_source_component/form"
	"github.com/erniealice/centymo-golang/domain/expenditure/cost_source_component/table"
	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	costsourcecomponentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/cost_source_component"
	expenditurelineitempb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/expenditure_line_item"
)

// Deps holds the action handler dependencies.
type Deps struct {
	Routes       csc.Routes
	Labels       csc.Labels
	CommonLabels pyeza.CommonLabels
	UseCases     *csc.UseCases

	// ReadParent reads the URL expenditure once through a workspace-scoped read and returns its
	// currency ("" when unset) and whether it exists in the actor's workspace. Required: the Add
	// drawer GET fails closed with 404 when it is nil or reports !exists (the drawer lists the
	// parent's bill lines, and expenditure_line_item is a column-less table (tenant-boundary
	// census), so the URL parent must be proven in the workspace first (R5 B3; expenditure detail
	// precedent)).
	ReadParent func(ctx context.Context, expenditureID string) (currency string, exists bool)
}

func notFound(deps *Deps) view.ViewResult {
	return view.ViewResult{
		StatusCode: http.StatusNotFound,
		Headers:    map[string]string{"HX-Error-Message": deps.Labels.Errors.NotFound},
	}
}

func denied(deps *Deps) view.ViewResult {
	return view.HTMXError(deps.CommonLabels.Errors.PermissionDenied)
}

func unavailable(deps *Deps) view.ViewResult {
	return view.ViewResult{
		StatusCode: http.StatusServiceUnavailable,
		Headers:    map[string]string{"HX-Error-Message": deps.Labels.Errors.Unavailable},
	}
}

func refuse(deps *Deps, err error) view.ViewResult {
	kind := csc.ErrorKind(err)
	if kind == csc.ErrUnknown {
		log.Printf("cost_source_component action: %v", err)
	}
	if kind == csc.ErrPermissionDenied {
		return denied(deps)
	}
	return view.HTMXError(deps.Labels.ErrorMessage(kind))
}

func options(l csc.Labels, selectedKind, selectedFact string) (kinds, facts []form.Option) {
	for _, k := range []costsourcecomponentpb.CostSourceComponentKind{
		costsourcecomponentpb.CostSourceComponentKind_COST_SOURCE_COMPONENT_KIND_ENERGY,
		costsourcecomponentpb.CostSourceComponentKind_COST_SOURCE_COMPONENT_KIND_WATER,
		costsourcecomponentpb.CostSourceComponentKind_COST_SOURCE_COMPONENT_KIND_SERVICE_FEE,
		costsourcecomponentpb.CostSourceComponentKind_COST_SOURCE_COMPONENT_KIND_OTHER,
	} {
		kinds = append(kinds, form.Option{Value: k.String(), Label: table.KindLabel(l, k), Selected: k.String() == selectedKind})
	}
	facts = append(facts, form.Option{Value: "", Label: ""})
	for _, f := range []costsourcecomponentpb.CostTaxFact{
		costsourcecomponentpb.CostTaxFact_COST_TAX_FACT_VATABLE,
		costsourcecomponentpb.CostTaxFact_COST_TAX_FACT_EXEMPT,
		costsourcecomponentpb.CostTaxFact_COST_TAX_FACT_ZERO_RATED,
		costsourcecomponentpb.CostTaxFact_COST_TAX_FACT_NOT_APPLICABLE,
	} {
		facts = append(facts, form.Option{Value: f.String(), Label: table.TaxFactLabel(l, f), Selected: f.String() == selectedFact})
	}
	return kinds, facts
}

// billLines lists the bill lines of one expenditure for the optional picker. The caller has already
// proven the expenditure in the workspace. The list is filtered server-side by expenditure_id (the
// adapter honors Filters; it does not honor Pagination, so one call returns at most PageSize rows:
// a bill with more lines than that gets a truncated picker, logged, never a foreign row).
func billLines(ctx context.Context, deps *Deps, expenditureID, selected string) []form.Option {
	if deps.UseCases == nil || deps.UseCases.ListExpenditureLineItems == nil {
		return nil
	}
	resp, err := deps.UseCases.ListExpenditureLineItems(ctx, &expenditurelineitempb.ListExpenditureLineItemsRequest{
		ExpenditureId: &expenditureID,
		Filters: &commonpb.FilterRequest{Filters: []*commonpb.TypedFilter{{
			Field: "expenditure_id",
			FilterType: &commonpb.TypedFilter_StringFilter{StringFilter: &commonpb.StringFilter{
				Value: expenditureID, Operator: commonpb.StringOperator_STRING_EQUALS, CaseSensitive: true,
			}},
		}}},
	})
	if err != nil {
		log.Printf("cost_source_component action: list bill lines for %s: %v", expenditureID, err)
		return nil
	}
	if len(resp.GetData()) >= csc.PageSize {
		log.Printf("cost_source_component action: bill lines of %s reach the %d-row read cap; picker truncated", expenditureID, csc.PageSize)
	}
	out := []form.Option{{Value: "", Label: deps.Labels.Form.ExpenditureLinePlaceholder}}
	for _, li := range resp.GetData() {
		if li.GetExpenditureId() != expenditureID {
			continue
		}
		out = append(out, form.Option{Value: li.GetId(), Label: li.GetDescription(), Selected: li.GetId() == selected})
	}
	if len(out) == 1 {
		return nil // no bill lines: hide the picker
	}
	return out
}

func amountLabel(l csc.Labels, currency string) string {
	if currency == "" {
		return l.Form.AmountLabel
	}
	return l.Form.AmountLabel + " (" + currency + ")"
}

func currencyOf(ctx context.Context, deps *Deps, expenditureID string) string {
	if deps.ReadParent == nil {
		return ""
	}
	currency, _ := deps.ReadParent(ctx, expenditureID)
	return currency
}

func blank(s string) *string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	return &s
}

// parse turns the drawer POST into a component (Id / ExpenditureId set by the caller).
func parse(deps *Deps, r *http.Request) (*costsourcecomponentpb.CostSourceComponent, bool) {
	if err := r.ParseForm(); err != nil {
		return nil, false
	}
	kindVal, ok := costsourcecomponentpb.CostSourceComponentKind_value[r.FormValue("component_kind")]
	if !ok || kindVal == 0 {
		return nil, false
	}
	amount, ok := csc.ParseAmountCentavos(r.FormValue("amount"))
	if !ok {
		return nil, false
	}
	from, to := strings.TrimSpace(r.FormValue("service_from")), strings.TrimSpace(r.FormValue("service_to"))
	if from == "" || to == "" {
		return nil, false
	}
	c := &costsourcecomponentpb.CostSourceComponent{
		ComponentKind: costsourcecomponentpb.CostSourceComponentKind(kindVal),
		Description:   blank(r.FormValue("description")),
		BasisUnit:     blank(r.FormValue("basis_unit")),
		Amount:        amount,
		Currency:      strings.TrimSpace(r.FormValue("currency")),
		ServiceFrom:   &from,
		ServiceTo:     &to,
	}
	scaled, scale, present, valid := csc.ParseQuantity(r.FormValue("basis_quantity"))
	if !valid {
		return nil, false
	}
	if present {
		c.BasisQuantityScaled, c.BasisScale = &scaled, &scale
	}
	if v, ok := costsourcecomponentpb.CostTaxFact_value[r.FormValue("tax_fact")]; ok && v != 0 {
		f := costsourcecomponentpb.CostTaxFact(v)
		c.TaxFact = &f
	}
	c.ExpenditureLineItemId = blank(r.FormValue("expenditure_line_item_id"))
	return c, true
}

// NewAddAction handles GET (drawer) / POST (create) for one expenditure ({id}).
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("expenditure", "update") {
			return denied(deps)
		}
		expenditureID := viewCtx.Request.PathValue("id")
		if viewCtx.Request.Method == http.MethodGet {
			if deps.ReadParent == nil {
				return notFound(deps)
			}
			currency, exists := deps.ReadParent(ctx, expenditureID)
			if !exists {
				return notFound(deps)
			}
			kinds, facts := options(deps.Labels, "", "")
			return view.OK("cost-source-component-drawer-form", &form.Data{
				FormAction: route.ResolveURL(deps.Routes.AddURL, "id", expenditureID), ExpenditureID: expenditureID,
				Currency: currency, Kinds: kinds, TaxFacts: facts,
				BillLines: billLines(ctx, deps, expenditureID, ""), Labels: deps.Labels, CommonLabels: deps.CommonLabels,
				AmountLabel: amountLabel(deps.Labels, currency), SubmitLabel: deps.Labels.Buttons.Add,
			})
		}
		if deps.UseCases == nil || deps.UseCases.CreateCostSourceComponent == nil {
			return unavailable(deps)
		}
		c, ok := parse(deps, viewCtx.Request)
		if !ok {
			return view.HTMXError(deps.Labels.Errors.FormInvalid)
		}
		c.ExpenditureId = expenditureID
		if c.Currency == "" {
			c.Currency = currencyOf(ctx, deps, expenditureID)
		}
		if _, err := deps.UseCases.CreateCostSourceComponent(ctx, &costsourcecomponentpb.CreateCostSourceComponentRequest{Data: c}); err != nil {
			return refuse(deps, err)
		}
		return view.HTMXSuccess(table.TableID)
	})
}

func readOne(ctx context.Context, deps *Deps, id string) (*costsourcecomponentpb.CostSourceComponent, string) {
	if deps.UseCases == nil || deps.UseCases.ReadCostSourceComponent == nil {
		return nil, csc.ErrUnknown
	}
	resp, err := deps.UseCases.ReadCostSourceComponent(ctx, &costsourcecomponentpb.ReadCostSourceComponentRequest{Data: &costsourcecomponentpb.CostSourceComponent{Id: id}})
	if err != nil {
		return nil, csc.ErrorKind(err)
	}
	if len(resp.GetData()) == 0 || resp.GetData()[0] == nil {
		return nil, csc.ErrNotFound
	}
	return resp.GetData()[0], ""
}

// NewEditAction handles GET (drawer) / POST (update) for one component ({id}).
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("expenditure", "update") {
			return denied(deps)
		}
		id := viewCtx.Request.PathValue("id")
		if viewCtx.Request.Method == http.MethodGet {
			cur, kind := readOne(ctx, deps, id)
			if kind != "" {
				if kind == csc.ErrUnknown {
					return unavailable(deps)
				}
				if kind == csc.ErrPermissionDenied {
					return denied(deps)
				}
				return view.HTMXError(deps.Labels.ErrorMessage(kind))
			}
			if table.Claimed(cur) {
				return view.HTMXError(deps.Labels.ErrorMessage(csc.ErrClaimed))
			}
			factName := ""
			if cur.TaxFact != nil {
				factName = cur.GetTaxFact().String()
			}
			kinds, facts := options(deps.Labels, cur.GetComponentKind().String(), factName)
			qty := ""
			if cur.BasisQuantityScaled != nil {
				qty = csc.FormatQuantity(cur.GetBasisQuantityScaled(), cur.GetBasisScale())
			}
			return view.OK("cost-source-component-drawer-form", &form.Data{
				FormAction: route.ResolveURL(deps.Routes.EditURL, "id", id), IsEdit: true, ID: id,
				ExpenditureID: cur.GetExpenditureId(), Currency: cur.GetCurrency(),
				Description: cur.GetDescription(), BasisUnit: cur.GetBasisUnit(), Quantity: qty,
				Amount: csc.FormatAmount(cur.GetAmount()), ServiceFrom: cur.GetServiceFrom(), ServiceTo: cur.GetServiceTo(),
				Kinds: kinds, TaxFacts: facts,
				BillLines:   billLines(ctx, deps, cur.GetExpenditureId(), cur.GetExpenditureLineItemId()),
				AmountLabel: amountLabel(deps.Labels, cur.GetCurrency()), SubmitLabel: deps.Labels.Buttons.Edit,
				Labels: deps.Labels, CommonLabels: deps.CommonLabels,
			})
		}
		if deps.UseCases == nil || deps.UseCases.UpdateCostSourceComponent == nil {
			return unavailable(deps)
		}
		c, ok := parse(deps, viewCtx.Request)
		if !ok {
			return view.HTMXError(deps.Labels.Errors.FormInvalid)
		}
		c.Id = id
		if _, err := deps.UseCases.UpdateCostSourceComponent(ctx, &costsourcecomponentpb.UpdateCostSourceComponentRequest{Data: c}); err != nil {
			return refuse(deps, err)
		}
		return view.HTMXSuccess(table.TableID)
	})
}

// NewDeleteAction handles POST (delete) for one component ({id}).
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("expenditure", "update") {
			return denied(deps)
		}
		if deps.UseCases == nil || deps.UseCases.DeleteCostSourceComponent == nil {
			return unavailable(deps)
		}
		id := viewCtx.Request.PathValue("id")
		if id == "" {
			return view.HTMXError(deps.Labels.Errors.FormInvalid)
		}
		if _, err := deps.UseCases.DeleteCostSourceComponent(ctx, &costsourcecomponentpb.DeleteCostSourceComponentRequest{Data: &costsourcecomponentpb.CostSourceComponent{Id: id}}); err != nil {
			return refuse(deps, err)
		}
		return view.HTMXSuccess(table.TableID)
	})
}
