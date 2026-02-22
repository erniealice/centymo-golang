package movements

import (
	"context"

	"github.com/erniealice/pyeza-golang/view"
)

// NewFilterView creates the HTMX partial for filtered movements table.
// Returns only the table-card HTML (no app-shell wrapper).
func NewFilterView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		r := viewCtx.Request

		dateFrom := r.URL.Query().Get("date_from")
		dateTo := r.URL.Query().Get("date_to")
		location := r.URL.Query().Get("location")
		txType := r.URL.Query().Get("type")
		search := r.URL.Query().Get("search")

		tableConfig := buildFilteredTable(ctx, deps, dateFrom, dateTo, location, txType, search)

		// Return just the table partial — the HTMX swap replaces #movements-table-wrapper
		data := struct {
			Table interface{}
		}{
			Table: tableConfig,
		}

		return view.OK("inventory-movements-table", data)
	})
}
