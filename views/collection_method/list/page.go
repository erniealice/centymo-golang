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

	cmpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/collection_method"
)

// ListViewDeps holds dependencies for the collection method list view.
type ListViewDeps struct {
	Routes centymo.CollectionMethodRoutes
	// ListCollectionMethods is nil until the espyna collection_method use cases
	// land (a separate wave). The view degrades to an empty list when nil.
	ListCollectionMethods func(ctx context.Context, req *cmpb.ListCollectionMethodsRequest) (*cmpb.ListCollectionMethodsResponse, error)
	Labels                centymo.CollectionMethodLabels
	CommonLabels          pyeza.CommonLabels
	TableLabels           types.TableLabels
}

// PageData holds the data for the collection method list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

// NewView creates the collection method list view (pages.md §B-5 list).
func NewView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("collection_method", "list") {
			return view.Forbidden("collection_method:list")
		}

		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}

		var methods []*cmpb.CollectionMethod
		if deps.ListCollectionMethods != nil {
			resp, err := deps.ListCollectionMethods(ctx, &cmpb.ListCollectionMethodsRequest{})
			if err != nil {
				log.Printf("Failed to list collection methods: %v", err)
				return view.Error(fmt.Errorf("failed to load collection methods: %w", err))
			}
			methods = resp.GetData()
		}

		methods = filterByLifecycle(methods, status)

		l := deps.Labels
		columns := collectionMethodColumns(l)
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
			ID:                   "collection-methods-table",
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
			ContentTemplate: "collection-method-list-content",
			Table:           tableConfig,
		}

		return view.OK("collection-method-list", pageData)
	})
}

// filterByLifecycle maps the URL {status} chip to lifecycle enum values.
func filterByLifecycle(methods []*cmpb.CollectionMethod, status string) []*cmpb.CollectionMethod {
	if status == "all" || status == "" {
		return methods
	}
	want := ""
	switch status {
	case "active":
		want = "COLLECTION_METHOD_LIFECYCLE_ACTIVE"
	case "draft":
		want = "COLLECTION_METHOD_LIFECYCLE_DRAFT"
	case "archived":
		want = "COLLECTION_METHOD_LIFECYCLE_ARCHIVED"
	case "closed":
		want = "COLLECTION_METHOD_LIFECYCLE_CLOSED"
	default:
		return methods
	}
	var out []*cmpb.CollectionMethod
	for _, m := range methods {
		if m.GetLifecycle().String() == want {
			out = append(out, m)
		}
	}
	return out
}

func collectionMethodColumns(l centymo.CollectionMethodLabels) []types.TableColumn {
	return []types.TableColumn{
		{Key: "template_code", Label: l.Columns.TemplateCode, WidthClass: "col-2xl"},
		{Key: "name", Label: l.Columns.Name},
		{Key: "category", Label: l.Columns.Category, WidthClass: "col-2xl"},
		{Key: "posting_kind", Label: l.Columns.PostingKind, WidthClass: "col-3xl"},
		{Key: "audience_mode", Label: l.Columns.AudienceMode, WidthClass: "col-3xl"},
		{Key: "lifecycle", Label: l.Columns.Lifecycle, WidthClass: "col-2xl"},
		{Key: "source", Label: l.Columns.Source, WidthClass: "col-2xl"},
		{Key: "revision", Label: l.Columns.Revision, Align: "right", WidthClass: "col-xs"},
	}
}

func buildTableRows(methods []*cmpb.CollectionMethod, routes centymo.CollectionMethodRoutes, l centymo.CollectionMethodLabels) []types.TableRow {
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
				{Type: "text", Value: enumShort(m.GetAudienceMode().String())},
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

// enumShort trims the long SCREAMING_CASE proto prefix to the trailing token(s)
// for display (e.g. COLLECTION_METHOD_CATEGORY_VOUCHER → VOUCHER).
func enumShort(s string) string {
	if s == "" {
		return ""
	}
	markers := []string{"_CATEGORY_", "_POSTING_KIND_", "_AUDIENCE_MODE_", "_LIFECYCLE_", "_SOURCE_", "_TAX_EFFECT_KIND_", "_VERSION_STATUS_"}
	for _, mk := range markers {
		if i := strings.Index(s, mk); i >= 0 {
			return s[i+len(mk):]
		}
	}
	return s
}

func lifecycleVariant(lifecycle string) string {
	switch lifecycle {
	case "COLLECTION_METHOD_LIFECYCLE_ACTIVE":
		return "success"
	case "COLLECTION_METHOD_LIFECYCLE_DRAFT":
		return "default"
	case "COLLECTION_METHOD_LIFECYCLE_CLOSED":
		return "warning"
	case "COLLECTION_METHOD_LIFECYCLE_ARCHIVED":
		return "danger"
	default:
		return "default"
	}
}

func statusPageTitle(l centymo.CollectionMethodLabels, status string) string {
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
