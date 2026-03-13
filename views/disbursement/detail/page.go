package detail

import (
	"context"
	"fmt"
	"log"

	centymo "github.com/erniealice/centymo-golang"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/fycha-golang/views/attachment"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	disbursementpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/disbursement"
)

// Deps holds view dependencies.
type Deps struct {
	Routes           centymo.DisbursementRoutes
	ReadDisbursement func(ctx context.Context, req *disbursementpb.ReadDisbursementRequest) (*disbursementpb.ReadDisbursementResponse, error)
	Labels           centymo.DisbursementLabels
	CommonLabels     pyeza.CommonLabels
	TableLabels      types.TableLabels

	// Attachment deps
	UploadFile       func(ctx context.Context, bucket, key string, content []byte, contentType string) error
	ListAttachments  func(ctx context.Context, moduleKey, foreignKey string) (*attachmentpb.ListAttachmentsResponse, error)
	CreateAttachment func(ctx context.Context, req *attachmentpb.CreateAttachmentRequest) (*attachmentpb.CreateAttachmentResponse, error)
	DeleteAttachment func(ctx context.Context, req *attachmentpb.DeleteAttachmentRequest) (*attachmentpb.DeleteAttachmentResponse, error)
	NewID            func() string
}

// PageData holds the data for the disbursement detail page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Disbursement    map[string]any
	Labels          centymo.DisbursementLabels
	ActiveTab       string
	TabItems        []pyeza.TabItem

	// Convenience fields for template rendering
	Reference     string
	StatusLabel   string
	StatusVariant string
	Amount        string
	Currency      string

	AuditTable          *types.TableConfig
	AttachmentTable     *types.TableConfig
	AttachmentUploadURL string
}

// disbursementToMap converts a Disbursement protobuf to a map[string]any for template use.
func disbursementToMap(d *disbursementpb.Disbursement) map[string]any {
	return map[string]any{
		"id":                      d.GetId(),
		"name":                    d.GetName(),
		"reference_number":        d.GetReferenceNumber(),
		"amount":                  fmt.Sprintf("%.2f", d.GetAmount()),
		"currency":                d.GetCurrency(),
		"status":                  d.GetStatus(),
		"disbursement_method_id":  d.GetDisbursementMethodId(),
		"disbursement_type":       d.GetDisbursementType(),
		"expenditure_id":          d.GetExpenditureId(),
		"approved_by":             d.GetApprovedBy(),
		"active":                  d.GetActive(),
		"date_created_string":     d.GetDateCreatedString(),
		"date_modified_string":    d.GetDateModifiedString(),
	}
}

// NewView creates the disbursement detail view.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")

		resp, err := deps.ReadDisbursement(ctx, &disbursementpb.ReadDisbursementRequest{
			Data: &disbursementpb.Disbursement{Id: id},
		})
		if err != nil {
			log.Printf("Failed to read disbursement %s: %v", id, err)
			return view.Error(fmt.Errorf("failed to load disbursement: %w", err))
		}
		data := resp.GetData()
		if len(data) == 0 {
			log.Printf("Disbursement %s not found", id)
			return view.Error(fmt.Errorf("disbursement not found"))
		}
		record := data[0]
		disbursement := disbursementToMap(record)

		refNumber := record.GetReferenceNumber()
		status := record.GetStatus()
		currency := record.GetCurrency()
		amount := fmt.Sprintf("%.2f", record.GetAmount())

		l := deps.Labels
		headerTitle := l.Detail.TitlePrefix + refNumber

		activeTab := viewCtx.QueryParams["tab"]
		if activeTab == "" {
			activeTab = "info"
		}
		tabItems := buildTabItems(l, id, deps.Routes)

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          headerTitle,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "cash",
				HeaderTitle:    headerTitle,
				HeaderSubtitle: l.Detail.PageTitle,
				HeaderIcon:     "icon-arrow-up-right",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "disbursement-detail-content",
			Disbursement:    disbursement,
			Labels:          l,
			ActiveTab:       activeTab,
			TabItems:        tabItems,
			Reference:       refNumber,
			StatusLabel:     status,
			StatusVariant:   statusVariant(status),
			Amount:          amount,
			Currency:        currency,
		}

		// Load tab-specific data
		switch activeTab {
		case "info":
			// Disbursement map has everything
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
		}

		return view.OK("disbursement-detail", pageData)
	})
}

func buildTabItems(l centymo.DisbursementLabels, id string, routes centymo.DisbursementRoutes) []pyeza.TabItem {
	base := route.ResolveURL(routes.DetailURL, "id", id)
	action := route.ResolveURL(routes.TabActionURL, "id", id, "tab", "")
	return []pyeza.TabItem{
		{Key: "info", Label: l.Detail.TabBasicInfo, Href: base + "?tab=info", HxGet: action + "info", Icon: "icon-info"},
		{Key: "attachments", Label: "Attachments", Href: base + "?tab=attachments", HxGet: action + "attachments", Icon: "icon-paperclip"},
		{Key: "audit", Label: l.Detail.TabAuditTrail, Href: base + "?tab=audit", HxGet: action + "audit", Icon: "icon-clock"},
	}
}

// NewTabAction creates the tab action view (partial — returns only the tab content).
func NewTabAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")
		tab := viewCtx.Request.PathValue("tab")
		if tab == "" {
			tab = "info"
		}

		resp, err := deps.ReadDisbursement(ctx, &disbursementpb.ReadDisbursementRequest{
			Data: &disbursementpb.Disbursement{Id: id},
		})
		if err != nil {
			log.Printf("Failed to read disbursement %s: %v", id, err)
			return view.Error(fmt.Errorf("failed to load disbursement: %w", err))
		}
		data := resp.GetData()
		if len(data) == 0 {
			log.Printf("Disbursement %s not found", id)
			return view.Error(fmt.Errorf("disbursement not found"))
		}
		record := data[0]
		disbursement := disbursementToMap(record)

		status := record.GetStatus()
		currency := record.GetCurrency()
		amount := fmt.Sprintf("%.2f", record.GetAmount())
		refNumber := record.GetReferenceNumber()

		l := deps.Labels
		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion: viewCtx.CacheVersion,
				CommonLabels: deps.CommonLabels,
			},
			Disbursement:  disbursement,
			Labels:        l,
			ActiveTab:     tab,
			TabItems:      buildTabItems(l, id, deps.Routes),
			Reference:     refNumber,
			StatusLabel:   status,
			StatusVariant: statusVariant(status),
			Amount:        amount,
			Currency:      currency,
		}

		switch tab {
		case "info":
			// disbursement map has everything
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
		}

		templateName := "disbursement-tab-" + tab
		if tab == "attachments" {
			templateName = "attachment-tab"
		}
		return view.OK(templateName, pageData)
	})
}

func buildAuditTable(l centymo.DisbursementLabels, tableLabels types.TableLabels) *types.TableConfig {
	columns := []types.TableColumn{
		{Key: "date", Label: l.Detail.Date, Sortable: true, Width: "160px"},
		{Key: "action", Label: l.Detail.AuditAction, Sortable: true},
		{Key: "user", Label: l.Detail.AuditUser, Sortable: true, Width: "180px"},
	}

	rows := []types.TableRow{}

	types.ApplyColumnStyles(columns, rows)

	cfg := &types.TableConfig{
		ID:                   "audit-trail-table",
		Columns:              columns,
		Rows:                 rows,
		ShowSearch:           true,
		ShowEntries:          true,
		DefaultSortColumn:    "date",
		DefaultSortDirection: "desc",
		Labels:               tableLabels,
		EmptyState: types.TableEmptyState{
			Title:   l.Detail.AuditEmptyTitle,
			Message: l.Detail.AuditEmptyMessage,
		},
	}
	types.ApplyTableSettings(cfg)

	return cfg
}

func statusVariant(status string) string {
	switch status {
	case "draft":
		return "default"
	case "pending":
		return "warning"
	case "approved":
		return "info"
	case "paid":
		return "success"
	case "cancelled":
		return "danger"
	case "overdue":
		return "danger"
	default:
		return "default"
	}
}
