package list

import (
	"context"
	"fmt"
	"log"

	"github.com/erniealice/centymo-golang"

	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// Deps holds view dependencies.
type Deps struct {
	DB     centymo.DataSource
	Routes centymo.PlanRoutes
	Labels centymo.PlanLabels
}

// PageData holds the data for the plan list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

// NewView creates the plan list view.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)

		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}

		records, err := deps.DB.ListSimple(ctx, "plan")
		if err != nil {
			log.Printf("Failed to list plans: %v", err)
			return view.Error(fmt.Errorf("failed to load plans: %w", err))
		}

		l := deps.Labels
		columns := planColumns(l)
		rows := buildTableRows(records, status, l, deps.Routes, perms)
		types.ApplyColumnStyles(columns, rows)

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          statusTitle(l, status),
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "services",
				ActiveSubNav:   "plans-" + status,
				HeaderTitle:    statusTitle(l, status),
				HeaderSubtitle: statusSubtitle(l, status),
				HeaderIcon:     "icon-file-text",
			},
			ContentTemplate: "plan-list-content",
			Table: &types.TableConfig{
				ID:          "plans-table",
				Columns:     columns,
				Rows:        rows,
				ShowSearch:  true,
				ShowActions: true,
				EmptyState: types.TableEmptyState{
					Title:   l.Empty.Title,
					Message: l.Empty.Message,
				},
				PrimaryAction: &types.PrimaryAction{
					Label:           l.Buttons.AddPlan,
					ActionURL:       deps.Routes.AddURL,
					Icon:            "icon-plus",
					Disabled:        !perms.Can("plan", "create"),
					DisabledTooltip: l.Errors.NoPermission,
				},
			},
		}

		return view.OK("plan-list", pageData)
	})
}

func planColumns(l centymo.PlanLabels) []types.TableColumn {
	return []types.TableColumn{
		{Key: "name", Label: l.Columns.Name, Sortable: true},
		{Key: "interval", Label: l.Columns.Interval, Sortable: true, Width: "150px"},
		{Key: "price", Label: l.Columns.Price, Sortable: true, Width: "120px"},
		{Key: "status", Label: l.Columns.Status, Sortable: true, Width: "120px"},
	}
}

func buildTableRows(records []map[string]any, status string, l centymo.PlanLabels, routes centymo.PlanRoutes, perms *types.UserPermissions) []types.TableRow {
	rows := []types.TableRow{}
	for _, record := range records {
		recordStatus, _ := record["status"].(string)
		if recordStatus != "" && recordStatus != status {
			continue
		}

		id, _ := record["id"].(string)
		name, _ := record["name"].(string)
		interval, _ := record["interval"].(string)
		price, _ := record["price"].(string)
		if recordStatus == "" {
			recordStatus = status
		}
		if interval == "" {
			interval = "monthly"
		}

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: name},
				{Type: "badge", Value: interval, Variant: intervalVariant(interval)},
				{Type: "text", Value: price},
				{Type: "badge", Value: recordStatus, Variant: statusVariant(recordStatus)},
			},
			DataAttrs: map[string]string{
				"name":     name,
				"interval": interval,
				"price":    price,
				"status":   recordStatus,
			},
			Actions: []types.TableAction{
				{Type: "view", Label: l.Actions.View, Action: "view", Href: route.ResolveURL(routes.DetailURL, "id", id)},
				{Type: "edit", Label: l.Actions.Edit, Action: "edit", URL: route.ResolveURL(routes.EditURL, "id", id), DrawerTitle: l.Actions.Edit, Disabled: !perms.Can("plan", "update"), DisabledTooltip: l.Errors.NoPermission},
				{Type: "delete", Label: l.Actions.Delete, Action: "delete", URL: routes.DeleteURL, ItemName: name, Disabled: !perms.Can("plan", "delete"), DisabledTooltip: l.Errors.NoPermission},
			},
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
