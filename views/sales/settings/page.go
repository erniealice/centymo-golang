package settings

import (
	"context"
	"fmt"
	"log"

	centymo "github.com/erniealice/centymo-golang"

	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	documenttemplatepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/document_template"
)

// TemplateData holds display-friendly fields for a single document template.
type TemplateData struct {
	ID              string
	Name            string
	TemplateType    string
	DocumentPurpose string
	OriginalFile    string
	FileSizeBytes   int64
	IsDefault       bool
	SetDefaultURL   string // pre-resolved URL for set-default action
}

// PageData holds the data for the sales settings templates page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

// NewView creates the sales settings templates list view.
// Route: GET /app/sales/settings/templates
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		templates := loadTemplateList(ctx, deps)
		l := deps.Labels.Settings

		columns := settingsColumns(l)
		rows := buildSettingsRows(templates, l, deps.Routes)
		types.ApplyColumnStyles(columns, rows)

		tableConfig := &types.TableConfig{
			ID:          "templates-table",
			RefreshURL:  deps.Routes.SettingsTemplatesURL,
			Columns:     columns,
			Rows:        rows,
			ShowSearch:  true,
			ShowActions: true,
			ShowEntries: true,
			Labels:      deps.TableLabels,
			EmptyState: types.TableEmptyState{
				Title:   l.EmptyTitle,
				Message: l.EmptyMessage,
			},
			PrimaryAction: &types.PrimaryAction{
				Label:     l.UploadTemplate,
				ActionURL: deps.Routes.SettingsTemplateUploadURL,
				Icon:      "icon-upload",
			},
		}
		types.ApplyTableSettings(tableConfig)

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          l.PageTitle,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "sales",
				HeaderTitle:    l.PageTitle,
				HeaderSubtitle: l.Caption,
				HeaderIcon:     "icon-file-text",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "sales-settings-templates-content",
			Table:           tableConfig,
		}

		return view.OK("sales-settings-templates", pageData)
	})
}

func settingsColumns(l centymo.SalesSettingsLabels) []types.TableColumn {
	return []types.TableColumn{
		{Key: "name", Label: l.TemplateName, Sortable: true},
		{Key: "type", Label: l.TemplateType, Sortable: true, Width: "120px"},
		{Key: "purpose", Label: l.Purpose, Sortable: true, Width: "120px"},
		{Key: "status", Label: l.DefaultBadge, Sortable: true, Width: "120px"},
	}
}

func buildSettingsRows(templates []TemplateData, l centymo.SalesSettingsLabels, routes centymo.SalesRoutes) []types.TableRow {
	rows := []types.TableRow{}
	for _, t := range templates {
		actions := []types.TableAction{}
		if !t.IsDefault {
			actions = append(actions, types.TableAction{
				Type:           "activate",
				Label:          l.SetDefault,
				Action:         "set-default",
				URL:            t.SetDefaultURL,
				ConfirmTitle:   l.SetDefault,
				ConfirmMessage: fmt.Sprintf("Set \"%s\" as the default template?", t.Name),
			})
		}
		actions = append(actions, types.TableAction{
			Type:           "delete",
			Label:          l.Delete,
			Action:         "delete",
			URL:            routes.SettingsTemplateDeleteURL,
			ItemName:       t.Name,
			ConfirmTitle:   l.Delete,
			ConfirmMessage: l.DeleteConfirm,
		})

		statusValue := ""
		statusVariant := ""
		if t.IsDefault {
			statusValue = l.DefaultBadge
			statusVariant = "info"
		}

		rows = append(rows, types.TableRow{
			ID: t.ID,
			Cells: []types.TableCell{
				{Type: "text", Value: t.Name},
				{Type: "text", Value: t.TemplateType},
				{Type: "text", Value: t.DocumentPurpose},
				{Type: "badge", Value: statusValue, Variant: statusVariant},
			},
			DataAttrs: map[string]string{
				"name":    t.Name,
				"type":    t.TemplateType,
				"purpose": t.DocumentPurpose,
				"status":  statusValue,
			},
			Actions: actions,
		})
	}
	return rows
}

// loadTemplateList loads all active invoice document templates.
func loadTemplateList(ctx context.Context, deps *Deps) []TemplateData {
	if deps.ListDocumentTemplates == nil {
		return nil
	}

	resp, err := deps.ListDocumentTemplates(ctx, &documenttemplatepb.ListDocumentTemplatesRequest{})
	if err != nil {
		log.Printf("Failed to list document templates: %v", err)
		return nil
	}

	var templates []TemplateData
	for _, t := range resp.GetData() {
		if !t.GetActive() {
			continue
		}
		if t.GetDocumentPurpose() != "invoice" {
			continue
		}
		templates = append(templates, TemplateData{
			ID:              t.GetId(),
			Name:            t.GetName(),
			TemplateType:    t.GetTemplateType(),
			DocumentPurpose: t.GetDocumentPurpose(),
			OriginalFile:    t.GetOriginalFilename(),
			FileSizeBytes:   t.GetFileSizeBytes(),
			IsDefault:       t.GetIsDefault(),
			SetDefaultURL:   route.ResolveURL(deps.Routes.SettingsTemplateDefaultURL, "id", t.GetId()),
		})
	}
	return templates
}
