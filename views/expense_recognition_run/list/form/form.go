// Package form holds the data types shared between the
// expense-recognition-run list page view and its table-card template.
//
// Mirror of packages/centymo-golang/views/revenue_run/list/form/form.go.
// Plan A 20260517-expense-run Phase 4 / Surface D.
package form

import (
	"github.com/erniealice/pyeza-golang/types"
)

// PageData is the full data context passed to the
// expense-recognition-run-list template.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}
