package dashboard

import (
	"context"

	centymo "github.com/erniealice/centymo-golang"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// Deps holds view dependencies.
type Deps struct {
	Labels       centymo.RevenueLabels
	CommonLabels pyeza.CommonLabels
}

// PageData holds the data for the sales dashboard page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Labels          centymo.RevenueDashboardLabels
}

// NewView creates the sales dashboard view.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion: viewCtx.CacheVersion,
				Title:        deps.Labels.Dashboard.Title,
				CurrentPath:  viewCtx.CurrentPath,
				ActiveNav:    "sales",
				ActiveSubNav: "dashboard",
				HeaderTitle:  deps.Labels.Dashboard.Title,
				HeaderIcon:   "icon-shopping-bag",
				CommonLabels: deps.CommonLabels,
			},
			ContentTemplate: "sales-dashboard-content",
			Labels:          deps.Labels.Dashboard,
		}

		return view.OK("sales-dashboard", pageData)
	})
}
