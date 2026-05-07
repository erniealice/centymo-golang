// Package form holds the data types shared between the revenue-run list page
// view and its table-card template.
package form

import (
	"github.com/erniealice/pyeza-golang/types"
)

// PageData is the full data context passed to the revenue-run-list template.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}
