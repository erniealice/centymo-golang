package list

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	centymo "github.com/erniealice/centymo-golang"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	dmpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/disbursement_method"
)

// ListViewDeps holds dependencies for the disbursement method list view.
type ListViewDeps struct {
	Routes centymo.DisbursementMethodRoutes
	// ListDisbursementMethods is nil until the espyna disbursement_method use
	// cases land. The view degrades to an empty list when nil.
	ListDisbursementMethods func(ctx context.Context, req *dmpb.ListDisbursementMethodsRequest) (*dmpb.ListDisbursementMethodsResponse, error)
	Labels                  centymo.DisbursementMethodLabels
	CommonLabels            pyeza.CommonLabels
	TableLabels             types.TableLabels
}

// PageData holds the data for the disbursement method list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

// NewView creates the disbursement method list view (pages.md §C-5 list).
// Buying-side asymmetry (D-4.9): no audience_mode column.
func NewView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("disbursement_method", "list") {
			return view.Forbidden("disbursement_method:list")
		}

		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}

		var methods []*dmpb.DisbursementMethod
		if deps.ListDisbursementMethods != nil {
			resp, err := deps.ListDisbursementMethods(ctx, &dmpb.ListDisbursementMethodsRequest{})
			if err != nil {
				log.Printf("Failed to list disbursement methods: %v", err)
				return view.Error(fmt.Errorf("failed to load disbursement methods: %w", err))
			}
			methods = resp.GetData()
		}

		methods = filterByLifecycle(methods, status)

		l := deps.Labels
		columns := disbursementMethodColumns(l)
		rows := buildTableRows(methods, deps.Routes, l)
		types.ApplyColumnStyles(columns, rows)

		var primaryAction *types.PrimaryAction
		if deps.Routes.AddURL != "" {
			primaryAction = &types.PrimaryAction{
				Label:     l.Page.AddButton,
				ActionURL: deps.Routes.AddURL,
			}
		}

		tableConfig := &types.TableConfig{
			ID:                   "disbursement-methods-table",
			RefreshURL:           deps.Routes.ListURL,
			Columns:              columns,
			Rows:                 rows,
			PrimaryAction:        primaryAction,
			ShowSearch:           true,
			ShowActions:          true,
			ShowFilters:          true,
			ShowSort:             true,
			ShowColumns:          true,
			ShowEntries:          true,
			DefaultSortColumn:    "name",
			DefaultSortDirection: "asc",
			Labels:               deps.TableLabels,
			EmptyState: types.TableEmptyState{
				Title:   l.Empty.Title,
				Message: l.Empty.Message,
			},
		}
		types.ApplyTableSettings(tableConfig)

		heading := statusPageTitle(l, status)
		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          heading,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      deps.Routes.ActiveNav,
				ActiveSubNav:   status,
				HeaderTitle:    heading,
				HeaderSubtitle: l.Page.Caption,
				HeaderIcon:     "icon-credit-card",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "disbursement-method-list-content",
			Table:           tableConfig,
		}

		return view.OK("disbursement-method-list", pageData)
	})
}

func filterByLifecycle(methods []*dmpb.DisbursementMethod, status string) []*dmpb.DisbursementMethod {
	if status == "all" || status == "" {
		return methods
	}
	want := ""
	switch status {
	case "active":
		want = "DISBURSEMENT_METHOD_LIFECYCLE_ACTIVE"
	case "draft":
		want = "DISBURSEMENT_METHOD_LIFECYCLE_DRAFT"
	case "archived":
		want = "DISBURSEMENT_METHOD_LIFECYCLE_ARCHIVED"
	case "closed":
		want = "DISBURSEMENT_METHOD_LIFECYCLE_CLOSED"
	default:
		return methods
	}
	var out []*dmpb.DisbursementMethod
	for _, m := range methods {
		if m.GetLifecycle().String() == want {
			out = append(out, m)
		}
	}
	return out
}

func disbursementMethodColumns(l centymo.DisbursementMethodLabels) []types.TableColumn {
	return []types.TableColumn{
		{Key: "template_code", Label: l.Columns.TemplateCode, WidthClass: "col-2xl"},
		{Key: "name", Label: l.Columns.Name},
		{Key: "category", Label: l.Columns.Category, WidthClass: "col-2xl"},
		{Key: "posting_kind", Label: l.Columns.PostingKind, WidthClass: "col-3xl"},
		{Key: "lifecycle", Label: l.Columns.Lifecycle, WidthClass: "col-2xl"},
		{Key: "source", Label: l.Columns.Source, WidthClass: "col-2xl"},
		{Key: "revision", Label: l.Columns.Revision, Align: "right", WidthClass: "col-xs"},
	}
}

func buildTableRows(methods []*dmpb.DisbursementMethod, routes centymo.DisbursementMethodRoutes, l centymo.DisbursementMethodLabels) []types.TableRow {
	rows := []types.TableRow{}
	for _, m := range methods {
		id := m.GetId()
		name := m.GetName()
		category := enumShort(m.GetCategory().String())
		lifecycle := enumShort(m.GetLifecycle().String())

		nameCell := types.TableCell{Type: "text", Value: name}
		if routes.DetailURL != "" {
			href := strings.Replace(routes.DetailURL, "{id}", id, 1)
			nameCell = types.TableCell{Type: "link", Value: name, Href: href}
		}

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: m.GetTemplateCode()},
				nameCell,
				{Type: "badge", Value: category, Variant: "default"},
				{Type: "text", Value: enumShort(m.GetPostingKind().String())},
				{Type: "badge", Value: lifecycle, Variant: lifecycleVariant(m.GetLifecycle().String())},
				{Type: "text", Value: enumShort(m.GetSource().String())},
				{Type: "number", Value: strconv.FormatInt(int64(m.GetRevision()), 10)},
			},
			DataAttrs: map[string]string{
				"method-id": id,
				"name":      name,
				"category":  category,
				"lifecycle": lifecycle,
			},
		})
	}
	return rows
}

func enumShort(s string) string {
	if s == "" {
		return ""
	}
	markers := []string{"_CATEGORY_", "_POSTING_KIND_", "_LIFECYCLE_", "_SOURCE_", "_TAX_EFFECT_KIND_", "_VERSION_STATUS_"}
	for _, mk := range markers {
		if i := strings.Index(s, mk); i >= 0 {
			return s[i+len(mk):]
		}
	}
	return s
}

func lifecycleVariant(lifecycle string) string {
	switch lifecycle {
	case "DISBURSEMENT_METHOD_LIFECYCLE_ACTIVE":
		return "success"
	case "DISBURSEMENT_METHOD_LIFECYCLE_DRAFT":
		return "default"
	case "DISBURSEMENT_METHOD_LIFECYCLE_CLOSED":
		return "warning"
	case "DISBURSEMENT_METHOD_LIFECYCLE_ARCHIVED":
		return "danger"
	default:
		return "default"
	}
}

func statusPageTitle(l centymo.DisbursementMethodLabels, status string) string {
	switch status {
	case "active":
		return l.Page.HeadingActive
	case "draft":
		return l.Page.HeadingDraft
	case "archived":
		return l.Page.HeadingArchived
	default:
		return l.Page.Heading
	}
}
