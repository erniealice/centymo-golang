package dashboard

import (
	"context"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// Deps holds view dependencies.
type Deps struct {
	CommonLabels pyeza.CommonLabels
}

// PageData holds the data for the inventory dashboard page.
type PageData struct {
	types.PageData
	ContentTemplate string
}

// NewView creates the inventory dashboard view.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion: viewCtx.CacheVersion,
				Title:        "Inventory Dashboard",
				CurrentPath:  viewCtx.CurrentPath,
				ActiveNav:    "inventory",
				ActiveSubNav: "dashboard",
				HeaderTitle:  "Inventory Dashboard",
				HeaderIcon:   "icon-briefcase",
				CommonLabels: deps.CommonLabels,
			},
			ContentTemplate: "inventory-dashboard-content",
		}

		return view.OK("inventory-dashboard", pageData)
	})
}
