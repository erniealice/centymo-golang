package detail

import (
	"context"
	"log"

	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	csc "github.com/erniealice/centymo-golang/domain/expenditure/cost_source_component"
	csctable "github.com/erniealice/centymo-golang/domain/expenditure/cost_source_component/table"
	"github.com/erniealice/centymo-golang/domain/expenditure/expenditure"
	expenditurepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/expenditure"
)

// RecoverableCostsData is the template data of "expense-tab-recoverable-costs".
type RecoverableCostsData struct {
	Labels        csc.Labels
	Table         *types.TableConfig
	AddURL        string
	CanAdd        bool
	AddTooltip    string
	Failed        bool
	ComponentsSum string
	BillTotal     string
	Difference    string
	Currency      string
}

// recoverableCostsEnabled reports whether the opt-in tab is wired for this app.
func recoverableCostsEnabled(deps *DetailViewDeps) bool {
	return deps.CostSourceComponents != nil && deps.CostSourceComponents.ListCostSourceComponents != nil
}

// recoverableCostsTabLabel is the tab label; its English default lives in DefaultLabels (general tier key
// expenditure.tabs.recoverable_costs).
func recoverableCostsTabLabel(l expenditure.Labels) string {
	return l.Tabs.RecoverableCosts
}

// populateRecoverableCosts loads the components of the expenditure and builds the tab.
// When the tab is not wired it leaves pageData.RecoverableCosts nil (the tab template
// then renders nothing; the tab item itself is not offered).
func populateRecoverableCosts(ctx context.Context, deps *DetailViewDeps, exp *expenditurepb.Expenditure, pageData *PageData) {
	if !recoverableCostsEnabled(deps) {
		return
	}
	perms := view.GetUserPermissions(ctx)
	l := deps.CostSourceComponentLabels
	id := exp.GetId()
	data := &RecoverableCostsData{Labels: l, Currency: exp.GetCurrency()}
	data.AddURL = route.ResolveURL(deps.CostSourceComponentRoutes.AddURL, "id", id)
	data.CanAdd = perms.Can("expenditure", "update")
	if !data.CanAdd {
		data.AddTooltip = deps.CommonLabels.Errors.PermissionDenied
	}

	rows, err := deps.CostSourceComponents.ListForExpenditure(ctx, id)
	if err != nil {
		log.Printf("expenditure detail: recoverable costs of %s: %v", id, err)
		data.Failed = true
	}
	var sum int64
	for _, c := range rows {
		sum += c.GetAmount()
	}
	data.ComponentsSum = csc.FormatAmount(sum)
	data.BillTotal = csc.FormatAmount(exp.GetTotalAmount())
	data.Difference = csc.FormatAmount(exp.GetTotalAmount() - sum)

	data.Table = csctable.Build(rows, csctable.Input{
		ExpenditureID: id, Currency: exp.GetCurrency(), Labels: l, Routes: deps.CostSourceComponentRoutes,
		TableLabels: deps.TableLabels, AllocateURL: deps.AllocationBatchRoutes.AllocateURL,
		ViewAllocationURL: deps.AllocationBatchRoutes.ViewURL, Perms: perms, NoPermission: deps.CommonLabels.Errors.PermissionDenied,
		MissingPermission: deps.CommonLabels.Errors.MissingPermission,
	})
	pageData.RecoverableCosts = data
}
