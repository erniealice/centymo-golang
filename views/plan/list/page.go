package list

import (
	"context"
	"fmt"
	"log"

	centymo "github.com/erniealice/centymo-golang"
	espynahttp "github.com/erniealice/espyna-golang/contrib/http"
	"github.com/erniealice/espyna-golang/tableparams"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	planpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/plan"
)

// ListViewDeps holds view dependencies.
type ListViewDeps struct {
	Routes       centymo.PlanRoutes
	ListPlans    func(ctx context.Context, req *planpb.ListPlansRequest) (*planpb.ListPlansResponse, error)
	GetInUseIDs  func(ctx context.Context, ids []string) (map[string]bool, error)
	Labels       centymo.PlanLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// Client name lookup for the optional Client column. Returns id → display
	// name. Optional — when nil the column shows the bare client_id.
	// Used by buildClientNameMap below.
	ListClientNames func(ctx context.Context) map[string]string

	// JobTemplate name lookup for the optional Job Template column. Returns
	// id → display name. Optional — when nil the column shows the bare
	// job_template_id (mirrors the ListClientNames pattern).
	// 2026-04-29 auto-spawn-jobs-from-subscription plan §5 — surfaces the
	// configured workflow on the Plan list so operators can see at a glance
	// which Plans will spawn jobs on subscription activation.
	ListJobTemplateNames func(ctx context.Context) map[string]string
}

// PageData holds the data for the plan list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

var planSearchFields = []string{"name", "description"}

// NewView creates the plan list view.
func NewView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}

		columns := planColumns(deps.Labels)
		p, err := espynahttp.ParseTableParams(viewCtx.Request, types.SortableKeys(columns), "name", "asc")
		if err != nil {
			return view.Error(err)
		}

		tableConfig, err := buildTableConfig(ctx, deps, columns, status, p)
		if err != nil {
			return view.Error(err)
		}

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          statusTitle(deps.Labels, status),
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      deps.Routes.ActiveNav,
				ActiveSubNav:   deps.Routes.ActiveSubNav + "-" + status,
				HeaderTitle:    statusTitle(deps.Labels, status),
				HeaderSubtitle: statusSubtitle(deps.Labels, status),
				HeaderIcon:     "icon-file-text",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "plan-list-content",
			Table:           tableConfig,
		}

		return view.OK("plan-list", pageData)
	})
}

// NewTableView creates a view that returns only the table-card HTML.
// Used as the refresh target after status/CRUD operations so that only the table
// is swapped (not the entire page content).
func NewTableView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}

		columns := planColumns(deps.Labels)
		p, err := espynahttp.ParseTableParams(viewCtx.Request, types.SortableKeys(columns), "name", "asc")
		if err != nil {
			return view.Error(err)
		}

		tableConfig, err := buildTableConfig(ctx, deps, columns, status, p)
		if err != nil {
			return view.Error(err)
		}

		return view.OK("table-card", tableConfig)
	})
}

// buildTableConfig fetches plan data and builds the table configuration.
// All rows are returned regardless of client_id; the Client column is always present.
func buildTableConfig(ctx context.Context, deps *ListViewDeps, columns []types.TableColumn, status string, p tableparams.TableQueryParams) (*types.TableConfig, error) {
	perms := view.GetUserPermissions(ctx)

	listParams := espynahttp.ToListParams(p, planSearchFields)

	// Inject status filter for server-side pagination
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

	resp, err := deps.ListPlans(ctx, &planpb.ListPlansRequest{
		Search:     listParams.Search,
		Filters:    listParams.Filters,
		Sort:       listParams.Sort,
		Pagination: listParams.Pagination,
	})
	if err != nil {
		log.Printf("Failed to list plans: %v", err)
		return nil, fmt.Errorf("failed to load plans: %w", err)
	}

	items := resp.GetData()

	var inUseIDs map[string]bool
	if deps.GetInUseIDs != nil {
		var itemIDs []string
		for _, item := range items {
			itemIDs = append(itemIDs, item.GetId())
		}
		inUseIDs, _ = deps.GetInUseIDs(ctx, itemIDs)
	}

	clientNames := map[string]string{}
	if deps.ListClientNames != nil {
		clientNames = deps.ListClientNames(ctx)
	}

	templateNames := map[string]string{}
	if deps.ListJobTemplateNames != nil {
		templateNames = deps.ListJobTemplateNames(ctx)
	}

	l := deps.Labels
	rows := buildTableRows(items, status, l, deps.CommonLabels, deps.Routes, inUseIDs, perms, clientNames, templateNames)
	types.ApplyColumnStyles(columns, rows)

	bulkCfg := centymo.MapBulkConfig(deps.CommonLabels)
	bulkCfg.Actions = []types.BulkAction{
		{
			Key:              "activate",
			Label:            l.Status.Activate,
			Icon:             "icon-check-circle",
			Variant:          "success",
			Endpoint:         deps.Routes.BulkSetStatusURL,
			ExtraParamsJSON:  `{"target_status":"active"}`,
			ConfirmTitle:     l.Confirm.BulkActivate,
			ConfirmMessage:   l.Confirm.BulkActivateMessage,
			RequiresDataAttr: "activatable",
		},
		{
			Key:              "deactivate",
			Label:            l.Status.Deactivate,
			Icon:             "icon-x-circle",
			Variant:          "warning",
			Endpoint:         deps.Routes.BulkSetStatusURL,
			ExtraParamsJSON:  `{"target_status":"inactive"}`,
			ConfirmTitle:     l.Confirm.BulkDeactivate,
			ConfirmMessage:   l.Confirm.BulkDeactivateMessage,
			RequiresDataAttr: "deactivatable",
		},
		{
			Key:              "delete",
			Label:            l.Bulk.Delete,
			Icon:             "icon-trash-2",
			Variant:          "danger",
			Endpoint:         deps.Routes.BulkDeleteURL,
			ConfirmTitle:     l.Confirm.BulkDelete,
			ConfirmMessage:   l.Confirm.BulkDeleteMessage,
			RequiresDataAttr: "deletable",
		},
	}

	refreshURL := route.ResolveURL(deps.Routes.TableURL, "status", status)

	var primaryAction *types.PrimaryAction
	if status == "active" {
		primaryAction = &types.PrimaryAction{
			Label:           l.Buttons.AddPlan,
			ActionURL:       deps.Routes.AddURL,
			Icon:            "icon-plus",
			Disabled:        !perms.Can("plan", "create"),
			DisabledTooltip: l.Errors.NoPermission,
		}
	}

	tableConfig := &types.TableConfig{
		ID:                   "plans-table",
		RefreshURL:           refreshURL,
		Columns:              columns,
		Rows:                 rows,
		ShowSearch:           true,
		ShowActions:          true,
		ShowSort:             true,
		ShowColumns:          true,
		ShowDensity:          true,
		ShowEntries:          true,
		DefaultSortColumn:    "name",
		DefaultSortDirection: "asc",
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

func planColumns(l centymo.PlanLabels) []types.TableColumn {
	return []types.TableColumn{
		{Key: "name", Label: l.Columns.Name},
		{Key: "price", Label: l.Columns.Price, WidthClass: "col-9xl"},
		{Key: "client", Label: l.Form.ClientLabel, NoSort: true, WidthClass: "col-3xl"},
		{Key: "job_template", Label: l.Form.JobTemplate, NoSort: true, WidthClass: "col-3xl"},
	}
}

func buildTableRows(plans []*planpb.Plan, status string, l centymo.PlanLabels, cl pyeza.CommonLabels, routes centymo.PlanRoutes, inUseIDs map[string]bool, perms *types.UserPermissions, clientNames map[string]string, templateNames map[string]string) []types.TableRow {
	rows := []types.TableRow{}
	for _, p := range plans {
		active := p.GetActive()
		recordStatus := "active"
		if !active {
			recordStatus = "inactive"
		}

		id := p.GetId()
		name := p.GetName()
		description := p.GetDescription()
		clientID := p.GetClientId()
		clientLabel := ""
		if clientID != "" {
			if n, ok := clientNames[clientID]; ok {
				clientLabel = n
			} else {
				clientLabel = clientID
			}
		}

		tplLabel := ""
		if tid := p.GetJobTemplateId(); tid != "" {
			if n, ok := templateNames[tid]; ok {
				tplLabel = n
			} else {
				tplLabel = tid
			}
		}

		actions := []types.TableAction{
			{Type: "view", Label: l.Actions.View, Action: "view", Href: route.ResolveURL(routes.DetailURL, "id", id)},
			{Type: "edit", Label: l.Actions.Edit, Action: "edit", URL: route.ResolveURL(routes.EditURL, "id", id), DrawerTitle: l.Actions.Edit, Disabled: !perms.Can("plan", "update"), DisabledTooltip: l.Errors.NoPermission},
		}

		if recordStatus == "active" {
			actions = append(actions, types.TableAction{
				Type:            "clone",
				Label:           cl.Actions.Clone,
				Action:          "clone",
				URL:             route.ResolveURL(routes.EditURL, "id", id),
				DrawerTitle:     cl.Actions.Clone,
				Disabled:        !perms.Can("plan", "create"),
				DisabledTooltip: l.Errors.NoPermission,
			})
			actions = append(actions, types.TableAction{
				Type:            "deactivate",
				Label:           l.Actions.Deactivate,
				Action:          "deactivate",
				URL:             routes.SetStatusURL + "?status=inactive",
				ItemName:        name,
				ConfirmTitle:    l.Confirm.Deactivate,
				ConfirmMessage:  fmt.Sprintf(l.Confirm.DeactivateMessage, name),
				Disabled:        !perms.Can("plan", "update"),
				DisabledTooltip: l.Errors.NoPermission,
			})
		} else {
			actions = append(actions, types.TableAction{
				Type:            "activate",
				Label:           l.Actions.Activate,
				Action:          "activate",
				URL:             routes.SetStatusURL + "?status=active",
				ItemName:        name,
				ConfirmTitle:    l.Confirm.Activate,
				ConfirmMessage:  fmt.Sprintf(l.Confirm.ActivateMessage, name),
				Disabled:        !perms.Can("plan", "update"),
				DisabledTooltip: l.Errors.NoPermission,
			})
		}

		deleteAction := types.TableAction{
			Type:     "delete",
			Label:    l.Actions.Delete,
			Action:   "delete",
			URL:      routes.DeleteURL,
			ItemName: name,
		}
		if inUseIDs[id] {
			deleteAction.Disabled = true
			deleteAction.DisabledTooltip = l.Errors.CannotDelete
		}
		if !perms.Can("plan", "delete") {
			deleteAction.Disabled = true
			deleteAction.DisabledTooltip = l.Errors.NoPermission
		}
		actions = append(actions, deleteAction)

		cells := []types.TableCell{
			{Type: "text", Value: name},
			{Type: "text", Value: description},
		}
		if clientLabel != "" {
			cells = append(cells, types.TableCell{Type: "badge", Value: clientLabel, Variant: "info"})
		} else {
			cells = append(cells, types.TableCell{Type: "text", Value: ""})
		}
		if tplLabel != "" {
			cells = append(cells, types.TableCell{Type: "badge", Value: tplLabel, Variant: "info"})
		} else {
			cells = append(cells, types.TableCell{Type: "text", Value: "—"})
		}

		rows = append(rows, types.TableRow{
			ID:    id,
			Cells: cells,
			DataAttrs: map[string]string{
				"name":             name,
				"price":            description,
				"status":           recordStatus,
				"deletable":        boolAttr(!inUseIDs[id]),
				"activatable":      boolAttr(recordStatus == "inactive"),
				"deactivatable":    boolAttr(recordStatus == "active"),
				"client_id":        clientID,
				"job_template_id":  p.GetJobTemplateId(),
				"plan-job-template-name": tplLabel,
			},
			Actions: actions,
		})
	}
	return rows
}

func statusTitle(l centymo.PlanLabels, status string) string {
	switch status {
	case "active":
		return l.Page.HeadingActive
	case "inactive":
		return l.Page.HeadingInactive
	default:
		return l.Page.Heading
	}
}

func statusSubtitle(l centymo.PlanLabels, status string) string {
	switch status {
	case "active":
		return l.Page.CaptionActive
	case "inactive":
		return l.Page.CaptionInactive
	default:
		return l.Page.Caption
	}
}

func boolAttr(v bool) string {
	if v {
		return "true"
	}
	return "false"
}

func statusVariant(status string) string {
	switch status {
	case "active":
		return "success"
	case "inactive":
		return "warning"
	default:
		return "default"
	}
}

func intervalVariant(interval string) string {
	switch interval {
	case "monthly":
		return "info"
	case "annual":
		return "primary"
	default:
		return "default"
	}
}

