package list

import (
	"context"
	"fmt"
	"strconv"

	sgpp "github.com/erniealice/centymo-golang/domain/subscription/subscription_group_product_plan"
	espynahttp "github.com/erniealice/espyna-golang/contrib/http"
	"github.com/erniealice/espyna-golang/shared/tableparams"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	sgpppb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription_group_product_plan"
)

// sgppEntity mirrors action.sgppEntity (kept package-local to avoid an
// import cycle with action/).
const sgppEntity = "subscription_group_product_plan"

// ListViewDeps holds dependencies for the subscription_group_product_plan (the
// class) S7 admin list views. Name-resolution closures are batched LIST_IN
// maps (espyna.md §1b — never a per-row call); a nil closure degrades a
// column to the raw id, never an error.
type ListViewDeps struct {
	Routes                            sgpp.Routes
	ListSubscriptionGroupProductPlans func(ctx context.Context, req *sgpppb.ListSubscriptionGroupProductPlansRequest) (*sgpppb.ListSubscriptionGroupProductPlansResponse, error)
	Labels                            sgpp.Labels
	CommonLabels                      pyeza.CommonLabels
	TableLabels                       types.TableLabels

	// GetSubscriptionGroupProductPlanInUseIDs gates the bulk-delete
	// RequiresDataAttr idiom (plan.md §2 in-use guard).
	GetSubscriptionGroupProductPlanInUseIDs func(ctx context.Context, ids []string) (map[string]bool, error)
	// ListSubscriptionGroupNames / ListProductPlanNames / ListJobTemplateNames
	// resolve id -> display name for the Section/Offering/Curriculum columns
	// (one batch call each, not per-row).
	ListSubscriptionGroupNames func(ctx context.Context) map[string]string
	ListProductPlanNames       func(ctx context.Context) map[string]string
	ListJobTemplateNames       func(ctx context.Context) map[string]string
	// CountAssignmentsBySGPPID resolves sgppID -> active-assignment count for
	// the Teachers column (one batch call, filtered on the v2 FK f12 —
	// espyna.md §1b call #2).
	CountAssignmentsBySGPPID func(ctx context.Context, sgppIDs []string) (map[string]int, error)
}

// PageData holds the data for the subscription_group_product_plan list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

var sgppSearchFields = []string{"subscription_group_id", "product_plan_id", "job_template_id"}

// NewView creates the full-page subscription_group_product_plan list view.
func NewView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can(sgppEntity, "list") {
			return view.Forbidden(sgppEntity + ":list")
		}
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}
		columns := sgppColumns(deps.Labels)
		p, err := espynahttp.ParseTableParamsWithFilters(viewCtx.Request, types.SortableKeys(columns), types.FilterableKeys(columns), "date_created", "desc")
		if err != nil {
			return view.Error(err)
		}
		tableConfig, err := buildTableConfig(ctx, deps, status, columns, p)
		if err != nil {
			return view.Error(err)
		}

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          statusPageTitle(deps.Labels, status),
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      deps.Routes.ActiveNav,
				ActiveSubNav:   statusSubNav(deps.Routes.ActiveSubNav, status),
				HeaderTitle:    statusPageTitle(deps.Labels, status),
				HeaderSubtitle: deps.Labels.Page.Subtitle,
				HeaderIcon:     "icon-book-open",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "sgpp-list-content",
			Table:           tableConfig,
		}

		return view.OK("sgpp-list", pageData)
	})
}

// NewTableView creates the HTMX table-only subscription_group_product_plan
// list view.
func NewTableView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}
		columns := sgppColumns(deps.Labels)
		p, err := espynahttp.ParseTableParamsWithFilters(viewCtx.Request, types.SortableKeys(columns), types.FilterableKeys(columns), "date_created", "desc")
		if err != nil {
			return view.Error(err)
		}
		tableConfig, err := buildTableConfig(ctx, deps, status, columns, p)
		if err != nil {
			return view.Error(err)
		}
		return view.OK("table-card", tableConfig)
	})
}

func buildTableConfig(ctx context.Context, deps *ListViewDeps, status string, columns []types.TableColumn, p tableparams.TableQueryParams) (*types.TableConfig, error) {
	perms := view.GetUserPermissions(ctx)
	listParams := espynahttp.ToListParams(p, sgppSearchFields)

	activeValue := status != "inactive"
	if listParams.Filters == nil {
		listParams.Filters = &commonpb.FilterRequest{}
	}
	listParams.Filters.Filters = append(listParams.Filters.Filters, &commonpb.TypedFilter{
		Field: "active",
		FilterType: &commonpb.TypedFilter_BooleanFilter{
			BooleanFilter: &commonpb.BooleanFilter{Value: activeValue},
		},
	})

	resp, err := deps.ListSubscriptionGroupProductPlans(ctx, &sgpppb.ListSubscriptionGroupProductPlansRequest{
		Search:     listParams.Search,
		Filters:    listParams.Filters,
		Sort:       listParams.Sort,
		Pagination: listParams.Pagination,
	})
	if err != nil {
		return nil, err
	}

	items := resp.GetData()

	var inUseIDs map[string]bool
	var assignmentCounts map[string]int
	itemIDs := make([]string, 0, len(items))
	for _, item := range items {
		itemIDs = append(itemIDs, item.GetId())
	}
	if deps.GetSubscriptionGroupProductPlanInUseIDs != nil {
		inUseIDs, _ = deps.GetSubscriptionGroupProductPlanInUseIDs(ctx, itemIDs)
	}
	if deps.CountAssignmentsBySGPPID != nil {
		assignmentCounts, _ = deps.CountAssignmentsBySGPPID(ctx, itemIDs)
	}
	var groupNames, planNames, templateNames map[string]string
	if deps.ListSubscriptionGroupNames != nil {
		groupNames = deps.ListSubscriptionGroupNames(ctx)
	}
	if deps.ListProductPlanNames != nil {
		planNames = deps.ListProductPlanNames(ctx)
	}
	if deps.ListJobTemplateNames != nil {
		templateNames = deps.ListJobTemplateNames(ctx)
	}

	l := deps.Labels
	rows := buildTableRows(items, status, l, deps.CommonLabels, deps.Routes, inUseIDs, assignmentCounts, groupNames, planNames, templateNames, perms)
	types.ApplyColumnStyles(columns, rows)

	bulkCfg := pyeza.MapBulkConfig(deps.CommonLabels)
	bulkCfg.Actions = buildBulkActions(l, deps.CommonLabels, deps.Routes)

	refreshURL := route.ResolveURL(deps.Routes.TableURL, "status", status)

	var primaryAction *types.PrimaryAction
	if status == "active" {
		primaryAction = &types.PrimaryAction{
			Label:           l.Buttons.Add,
			ActionURL:       deps.Routes.AddURL,
			Icon:            "icon-plus",
			Disabled:        !perms.Can(sgppEntity, "create"),
			DisabledTooltip: fmt.Sprintf(deps.CommonLabels.Errors.MissingPermission, sgppEntity+":create"),
		}
	}

	tableConfig := &types.TableConfig{
		ID:                   "subscription-group-product-plans-table",
		RefreshURL:           refreshURL,
		Columns:              columns,
		Rows:                 rows,
		ShowSearch:           true,
		ShowActions:          true,
		ShowFilters:          true,
		ShowSort:             true,
		ShowColumns:          true,
		ShowExport:           true,
		ShowDensity:          true,
		ShowEntries:          true,
		DefaultSortColumn:    "date_created",
		DefaultSortDirection: "desc",
		Labels:               deps.TableLabels,
		EmptyState: types.TableEmptyState{
			Title:   l.Empty.Title,
			Message: l.Empty.Message,
		},
		PrimaryAction: primaryAction,
		BulkActions:   &bulkCfg,
	}
	types.ApplyTableSettings(tableConfig)
	return tableConfig, nil
}

func sgppColumns(l sgpp.Labels) []types.TableColumn {
	return []types.TableColumn{
		{Key: "subscription_group_id", Label: l.Columns.SubscriptionGroupId},
		{Key: "product_plan_id", Label: l.Columns.ProductPlanId},
		{Key: "job_template_id", Label: l.Columns.JobTemplateId, NoSort: true, NoFilter: true},
		{Key: "status", Label: l.Columns.Status, NoSort: true, NoFilter: true},
		{Key: "assignments_count", Label: l.Columns.AssignmentsCount, NoSort: true, NoFilter: true},
	}
}

func buildTableRows(items []*sgpppb.SubscriptionGroupProductPlan, status string, l sgpp.Labels, cl pyeza.CommonLabels, routes sgpp.Routes, inUseIDs map[string]bool, assignmentCounts map[string]int, groupNames, planNames, templateNames map[string]string, perms *types.UserPermissions) []types.TableRow {
	rows := []types.TableRow{}
	for _, c := range items {
		recordStatus := "active"
		if !c.GetActive() {
			recordStatus = "inactive"
		}

		id := c.GetId()
		groupID := c.GetSubscriptionGroupId()
		planID := c.GetProductPlanId()
		templateID := c.GetJobTemplateId()

		groupLabel := resolveName(groupNames, groupID)
		planLabel := resolveName(planNames, planID)
		templateLabel := resolveName(templateNames, templateID)

		isInUse := inUseIDs[id]

		statusValue, statusVariant := l.Form.StatusActive, "success"
		if c.GetStatus() == sgpppb.SubscriptionGroupProductPlanStatus_SUBSCRIPTION_GROUP_PRODUCT_PLAN_STATUS_EXCLUDED {
			statusValue, statusVariant = l.Form.StatusExcluded, "warning"
		}

		cells := []types.TableCell{
			{Type: "text", Value: groupLabel},
			{Type: "text", Value: planLabel},
			{Type: "text", Value: templateLabel},
			{Type: "badge", Value: statusValue, Variant: statusVariant},
			{Type: "text", Value: strconv.Itoa(assignmentCounts[id])},
		}

		rows = append(rows, types.TableRow{
			ID:    id,
			Cells: cells,
			DataAttrs: map[string]string{
				"status":    recordStatus,
				"deletable": strconv.FormatBool(!isInUse),
			},
			Actions: buildRowActions(id, groupID, planLabel, isInUse, l, cl, routes, perms),
		})
	}
	return rows
}

func resolveName(names map[string]string, id string) string {
	if id == "" {
		return ""
	}
	if names != nil {
		if n, ok := names[id]; ok && n != "" {
			return n
		}
	}
	return id
}

func buildRowActions(id, groupID, name string, isInUse bool, l sgpp.Labels, cl pyeza.CommonLabels, routes sgpp.Routes, perms *types.UserPermissions) []types.TableAction {
	actions := []types.TableAction{
		{Type: "view", Label: l.Buttons.View, Action: "view", Href: route.ResolveURL(routes.DetailURL, "id", groupID, "sgppId", id)},
		{Type: "edit", Label: l.Buttons.Edit, Action: "edit", URL: route.ResolveURL(routes.EditURL, "id", id), DrawerTitle: l.Buttons.Edit,
			Disabled: !perms.Can(sgppEntity, "update"), DisabledTooltip: fmt.Sprintf(cl.Errors.MissingPermission, sgppEntity+":update")},
	}

	deleteAction := types.TableAction{
		Type:     "delete",
		Label:    l.Buttons.Delete,
		Action:   "delete",
		URL:      routes.DeleteURL,
		ItemName: name,
	}
	if isInUse {
		deleteAction.Disabled = true
		deleteAction.DisabledTooltip = l.Errors.InUse
	} else if !perms.Can(sgppEntity, "delete") {
		deleteAction.Disabled = true
		deleteAction.DisabledTooltip = fmt.Sprintf(cl.Errors.MissingPermission, sgppEntity+":delete")
	}
	actions = append(actions, deleteAction)
	return actions
}

func buildBulkActions(l sgpp.Labels, cl pyeza.CommonLabels, routes sgpp.Routes) []types.BulkAction {
	return []types.BulkAction{
		{
			Key:              "delete",
			Label:            cl.Bulk.Delete,
			Icon:             "icon-trash-2",
			Variant:          "danger",
			Endpoint:         routes.BulkDeleteURL,
			ConfirmTitle:     l.Bulk.DeleteTitle,
			ConfirmMessage:   l.Bulk.DeleteMessage,
			RequiresDataAttr: "deletable",
		},
	}
}

func statusPageTitle(l sgpp.Labels, status string) string {
	switch status {
	case "active":
		return l.Page.ActiveTitle
	case "inactive":
		return l.Page.InactiveTitle
	default:
		return l.Page.Title
	}
}

func statusSubNav(base, status string) string {
	if base == "" {
		return status
	}
	return base + "-" + status
}
