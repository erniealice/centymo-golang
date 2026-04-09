package category

import (
	"context"
	"fmt"
	"log"

	centymo "github.com/erniealice/centymo-golang"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	expenditurecategorypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/expenditure_category"
)

// ListViewDeps holds view dependencies for the expenditure category list.
type ListViewDeps struct {
	Routes                    centymo.ExpenditureRoutes
	ListExpenditureCategories func(ctx context.Context, req *expenditurecategorypb.ListExpenditureCategoriesRequest) (*expenditurecategorypb.ListExpenditureCategoriesResponse, error)
	Labels                    centymo.ExpenditureLabels
	CommonLabels              pyeza.CommonLabels
	TableLabels               types.TableLabels
}

// PageData holds the data for the expenditure category list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

// NewView creates the expenditure category list view.
func NewView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)

		resp, err := deps.ListExpenditureCategories(ctx, &expenditurecategorypb.ListExpenditureCategoriesRequest{})
		if err != nil {
			log.Printf("Failed to list expenditure categories: %v", err)
			return view.Error(fmt.Errorf("failed to load expenditure categories: %w", err))
		}

		l := deps.Labels.Category
		columns := categoryColumns(l)
		rows := buildTableRows(resp.GetData(), l, deps.Routes, perms)
		types.ApplyColumnStyles(columns, rows)

		tableConfig := &types.TableConfig{
			ID:                   "expenditure-categories-table",
			RefreshURL:           deps.Routes.ExpenseCategoryTableURL,
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
			DefaultSortColumn:    "code",
			DefaultSortDirection: "asc",
			Labels:               deps.TableLabels,
			EmptyState: types.TableEmptyState{
				Title:   l.Empty.Title,
				Message: l.Empty.Message,
			},
			PrimaryAction: &types.PrimaryAction{
				Label:           l.Buttons.AddCategory,
				ActionURL:       deps.Routes.ExpenseCategoryAddURL,
				Icon:            "icon-plus",
				Disabled:        !perms.Can("expenditure_category", "create"),
				DisabledTooltip: l.Errors.PermissionDenied,
			},
		}
		types.ApplyTableSettings(tableConfig)

		heading := l.Page.Heading
		caption := l.Page.Caption
		if heading == "" {
			heading = "Expense Categories"
		}
		if caption == "" {
			caption = "Manage expense categories"
		}

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          heading,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "expense",
				ActiveSubNav:   "expense-categories",
				HeaderTitle:    heading,
				HeaderSubtitle: caption,
				HeaderIcon:     "icon-tag",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "expenditure-category-list-content",
			Table:           tableConfig,
		}

		return view.OK("expenditure-category-list", pageData)
	})
}

func categoryColumns(l centymo.ExpenditureCategoryLabels) []types.TableColumn {
	code := l.Columns.Code
	name := l.Columns.Name
	description := l.Columns.Description
	status := l.Columns.Status
	if code == "" {
		code = "Code"
	}
	if name == "" {
		name = "Name"
	}
	if description == "" {
		description = "Description"
	}
	if status == "" {
		status = "Status"
	}
	return []types.TableColumn{
		{Key: "code", Label: code, Sortable: true, WidthClass: "col-2xl"},
		{Key: "name", Label: name, Sortable: true},
		{Key: "description", Label: description, Sortable: false},
		{Key: "status", Label: status, Sortable: true, WidthClass: "col-lg"},
	}
}

func buildTableRows(
	categories []*expenditurecategorypb.ExpenditureCategory,
	l centymo.ExpenditureCategoryLabels,
	routes centymo.ExpenditureRoutes,
	perms *types.UserPermissions,
) []types.TableRow {
	editLabel := l.Actions.Edit
	deleteLabel := l.Actions.Delete
	permDenied := l.Errors.PermissionDenied
	deleteTitle := l.Confirm.DeleteTitle
	deleteMsg := l.Confirm.DeleteMessage
	if editLabel == "" {
		editLabel = "Edit"
	}
	if deleteLabel == "" {
		deleteLabel = "Delete"
	}
	if deleteTitle == "" {
		deleteTitle = "Delete Category"
	}

	rows := []types.TableRow{}
	for _, cat := range categories {
		id := cat.GetId()
		code := cat.GetCode()
		name := cat.GetName()
		description := cat.GetDescription()
		statusVal := "active"
		statusVariant := "success"
		if !cat.GetActive() {
			statusVal = "inactive"
			statusVariant = "default"
		}

		actions := []types.TableAction{
			{
				Type:            "edit",
				Label:           editLabel,
				Action:          "edit",
				URL:             routes.ExpenseCategoryEditURL,
				DrawerTitle:     editLabel,
				Disabled:        !perms.Can("expenditure_category", "update"),
				DisabledTooltip: permDenied,
			},
			{
				Type:            "delete",
				Label:           deleteLabel,
				Action:          "delete",
				URL:             routes.ExpenseCategoryDeleteURL,
				ItemName:        name,
				ConfirmTitle:    deleteTitle,
				ConfirmMessage:  fmt.Sprintf(deleteMsg, name),
				Disabled:        !perms.Can("expenditure_category", "delete"),
				DisabledTooltip: permDenied,
			},
		}

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: code},
				{Type: "text", Value: name},
				{Type: "text", Value: description},
				{Type: "badge", Value: statusVal, Variant: statusVariant},
			},
			DataAttrs: map[string]string{
				"code":        code,
				"name":        name,
				"description": description,
				"status":      statusVal,
			},
			Actions: actions,
		})
	}
	return rows
}
