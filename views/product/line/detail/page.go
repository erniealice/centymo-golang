package detail

import (
	"context"
	"fmt"
	"log"

	centymo "github.com/erniealice/centymo-golang"
	"github.com/erniealice/hybra-golang/views/attachment"
	"github.com/erniealice/hybra-golang/views/auditlog"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	linepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/line"
)

// DetailViewDeps holds view dependencies.
type DetailViewDeps struct {
	Routes       centymo.ProductLineRoutes
	ReadLine     func(ctx context.Context, req *linepb.ReadLineRequest) (*linepb.ReadLineResponse, error)
	Labels       centymo.ProductLineLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	attachment.AttachmentOps
	auditlog.AuditOps
}

// PageData holds the data for the line detail page.
type PageData struct {
	types.PageData
	ContentTemplate     string
	Line                map[string]any
	Labels              centymo.ProductLineLabels
	ActiveTab           string
	TabItems            []pyeza.TabItem
	AuditTable          *types.TableConfig
	AttachmentTable     *types.TableConfig
	AttachmentUploadURL string
	AuditEntries        []auditlog.AuditEntryView
	AuditHasNext        bool
	AuditNextCursor     string
	AuditHistoryURL     string
}

func lineToMap(line *linepb.Line) map[string]any {
	status := "active"
	if !line.GetActive() {
		status = "inactive"
	}
	return map[string]any{
		"id":                   line.GetId(),
		"name":                 line.GetName(),
		"description":          line.GetDescription(),
		"status":               status,
		"date_created_string":  line.GetDateCreatedString(),
		"date_modified_string": line.GetDateModifiedString(),
	}
}

// NewView creates the line detail view.
func NewView(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")
		resp, err := deps.ReadLine(ctx, &linepb.ReadLineRequest{Data: &linepb.Line{Id: id}})
		if err != nil {
			log.Printf("Failed to read line %s: %v", id, err)
			return view.Error(fmt.Errorf("failed to load line: %w", err))
		}
		data := resp.GetData()
		if len(data) == 0 {
			return view.Error(fmt.Errorf("line not found"))
		}
		line := lineToMap(data[0])

		l := deps.Labels
		headerTitle, _ := line["name"].(string)
		if headerTitle == "" {
			headerTitle = l.Detail.PageTitle
		}

		activeTab := viewCtx.QueryParams["tab"]
		if activeTab == "" {
			activeTab = "info"
		}

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          headerTitle,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      deps.Routes.ActiveNav,
				HeaderTitle:    headerTitle,
				HeaderSubtitle: l.Detail.PageTitle,
				HeaderIcon:     "icon-layers",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "product-line-detail-content",
			Line:            line,
			Labels:          l,
			ActiveTab:       activeTab,
			TabItems:        buildTabItems(l, id, deps.Routes),
		}

		switch activeTab {
		case "attachments":
			if deps.ListAttachments != nil {
				cfg := attachmentConfig(deps)
				resp, err := deps.ListAttachments(ctx, cfg.EntityType, id)
				if err != nil {
					log.Printf("Failed to list attachments: %v", err)
				}
				var items []*attachmentpb.Attachment
				if resp != nil {
					items = resp.GetData()
				}
				pageData.AttachmentTable = attachment.BuildTable(items, cfg, id)
			}
			pageData.AttachmentUploadURL = route.ResolveURL(deps.Routes.AttachmentUploadURL, "id", id)
		case "audit":
			pageData.AuditTable = buildAuditTable(l, deps.TableLabels)
		case "audit-history":
			if deps.ListAuditHistory != nil {
				cursor := viewCtx.QueryParams["cursor"]
				auditResp, err := deps.ListAuditHistory(ctx, &auditlog.ListAuditRequest{
					EntityType:  "line",
					EntityID:    id,
					Limit:       20,
					CursorToken: cursor,
				})
				if err != nil {
					log.Printf("Failed to load audit history: %v", err)
				}
				if auditResp != nil {
					pageData.AuditEntries = auditResp.Entries
					pageData.AuditHasNext = auditResp.HasNext
					pageData.AuditNextCursor = auditResp.NextCursor
				}
			}
			pageData.AuditHistoryURL = route.ResolveURL(deps.Routes.TabActionURL, "id", id, "tab", "") + "audit-history"
		}

		return view.OK("product-line-detail", pageData)
	})
}

func buildTabItems(l centymo.ProductLineLabels, id string, routes centymo.ProductLineRoutes) []pyeza.TabItem {
	base := route.ResolveURL(routes.DetailURL, "id", id)
	action := route.ResolveURL(routes.TabActionURL, "id", id, "tab", "")
	return []pyeza.TabItem{
		{Key: "info", Label: l.Detail.TabBasicInfo, Href: base + "?tab=info", HxGet: action + "info", Icon: "icon-info"},
		{Key: "attachments", Label: l.Detail.TabAttachments, Href: base + "?tab=attachments", HxGet: action + "attachments", Icon: "icon-paperclip"},
		{Key: "audit", Label: l.Detail.TabAuditTrail, Href: base + "?tab=audit", HxGet: action + "audit", Icon: "icon-clock"},
		{Key: "audit-history", Label: "History", Href: base + "?tab=audit-history", HxGet: action + "audit-history", Icon: "icon-clock"},
	}
}

// NewTabAction creates the tab action view.
func NewTabAction(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")
		tab := viewCtx.Request.PathValue("tab")
		if tab == "" {
			tab = "info"
		}

		resp, err := deps.ReadLine(ctx, &linepb.ReadLineRequest{Data: &linepb.Line{Id: id}})
		if err != nil {
			log.Printf("Failed to read line %s: %v", id, err)
			return view.Error(fmt.Errorf("failed to load line: %w", err))
		}
		data := resp.GetData()
		if len(data) == 0 {
			return view.Error(fmt.Errorf("line not found"))
		}
		line := lineToMap(data[0])

		l := deps.Labels
		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion: viewCtx.CacheVersion,
				CommonLabels: deps.CommonLabels,
			},
			Line:      line,
			Labels:    l,
			ActiveTab: tab,
			TabItems:  buildTabItems(l, id, deps.Routes),
		}

		switch tab {
		case "attachments":
			if deps.ListAttachments != nil {
				cfg := attachmentConfig(deps)
				resp, err := deps.ListAttachments(ctx, cfg.EntityType, id)
				if err != nil {
					log.Printf("Failed to list attachments: %v", err)
				}
				var items []*attachmentpb.Attachment
				if resp != nil {
					items = resp.GetData()
				}
				pageData.AttachmentTable = attachment.BuildTable(items, cfg, id)
			}
			pageData.AttachmentUploadURL = route.ResolveURL(deps.Routes.AttachmentUploadURL, "id", id)
		case "audit":
			pageData.AuditTable = buildAuditTable(l, deps.TableLabels)
		case "audit-history":
			if deps.ListAuditHistory != nil {
				cursor := viewCtx.QueryParams["cursor"]
				auditResp, err := deps.ListAuditHistory(ctx, &auditlog.ListAuditRequest{
					EntityType:  "line",
					EntityID:    id,
					Limit:       20,
					CursorToken: cursor,
				})
				if err != nil {
					log.Printf("Failed to load audit history: %v", err)
				}
				if auditResp != nil {
					pageData.AuditEntries = auditResp.Entries
					pageData.AuditHasNext = auditResp.HasNext
					pageData.AuditNextCursor = auditResp.NextCursor
				}
			}
			pageData.AuditHistoryURL = route.ResolveURL(deps.Routes.TabActionURL, "id", id, "tab", "") + "audit-history"
		}

		templateName := "product-line-tab-" + tab
		if tab == "attachments" {
			templateName = "attachment-tab"
		}
		if tab == "audit-history" {
			templateName = "audit-history-tab"
		}
		return view.OK(templateName, pageData)
	})
}

func buildAuditTable(l centymo.ProductLineLabels, tableLabels types.TableLabels) *types.TableConfig {
	return &types.TableConfig{
		ID: "product-line-audit-table",
		Columns: []types.TableColumn{
			{Key: "date", Label: "Date"},
			{Key: "action", Label: l.Detail.AuditAction, NoSort: true},
			{Key: "user", Label: l.Detail.AuditUser, NoSort: true},
		},
		Rows:       []types.TableRow{},
		ShowSearch: false, ShowActions: false, ShowFilters: false, ShowSort: false,
		ShowColumns: false, ShowExport: false, ShowDensity: false, ShowEntries: false,
		Labels:     tableLabels,
		EmptyState: types.TableEmptyState{Title: l.Detail.AuditEmptyTitle, Message: l.Detail.AuditEmptyMessage},
	}
}
