package detail

import (
	"context"
	"fmt"
	"log"

	"github.com/erniealice/hybra-golang/views/attachment"
	"github.com/erniealice/hybra-golang/views/auditlog"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	pricelistpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/price_list"
	priceproductpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/price_product"

	"github.com/erniealice/centymo-golang"
	lynguaV1 "github.com/erniealice/lyngua/golang/v1"
)

// DetailViewDeps holds view dependencies.
type DetailViewDeps struct {
	Routes            centymo.PriceListRoutes
	ReadPriceList     func(ctx context.Context, req *pricelistpb.ReadPriceListRequest) (*pricelistpb.ReadPriceListResponse, error)
	ListPriceProducts func(ctx context.Context, req *priceproductpb.ListPriceProductsRequest) (*priceproductpb.ListPriceProductsResponse, error)
	Labels            centymo.PriceListLabels
	CommonLabels      pyeza.CommonLabels
	TableLabels       types.TableLabels

	attachment.AttachmentOps
	auditlog.AuditOps
}

// PageData holds the data for the price list detail page.
type PageData struct {
	types.PageData
	ContentTemplate     string
	PriceList           *pricelistpb.PriceList
	ActiveTab           string
	TabItems            []pyeza.TabItem
	ID                  string
	PricesTable         *types.TableConfig
	AttachmentTable     *types.TableConfig
	AttachmentUploadURL string
	Labels              centymo.PriceListLabels
	// Audit history tab
	AuditEntries    []auditlog.AuditEntryView
	AuditHasNext    bool
	AuditNextCursor string
	AuditHistoryURL string
}

// NewView creates the price list detail view.
func NewView(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")

		tab := viewCtx.Request.URL.Query().Get("tab")
		if tab == "" {
			tab = "basic"
		}

		resp, err := deps.ReadPriceList(ctx, &pricelistpb.ReadPriceListRequest{
			Data: &pricelistpb.PriceList{Id: id},
		})
		if err != nil {
			log.Printf("Failed to read price list %s: %v", id, err)
			return view.Error(fmt.Errorf("failed to load price list: %w", err))
		}

		data := resp.GetData()
		if len(data) == 0 {
			return view.Error(fmt.Errorf("price list not found"))
		}
		priceList := data[0]

		name := priceList.GetName()
		description := priceList.GetDescription()

		tabItems := buildTabItems(id, deps.Labels, deps.Routes)

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          name,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "sale",
				HeaderTitle:    name,
				HeaderSubtitle: description,
				HeaderIcon:     "icon-tag",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "pricelist-detail-content",
			PriceList:       priceList,
			ActiveTab:       tab,
			TabItems:        tabItems,
			ID:              id,
			Labels:          deps.Labels,
		}

		// Populate tab-specific data
		switch tab {
		case "prices":
			perms := view.GetUserPermissions(ctx)
			pricesTable, err := buildPricesTable(ctx, deps, id, deps.Routes, perms)
			if err != nil {
				log.Printf("Failed to load price products for price list %s: %v", id, err)
			}
			pageData.PricesTable = pricesTable

		case "attachments":
			if deps.ListAttachments != nil {
				cfg := attachmentConfig(deps)
				resp, err := deps.ListAttachments(ctx, cfg.EntityType, id)
				if err != nil {
					log.Printf("Failed to list attachments for %s %s: %v", cfg.EntityType, id, err)
				}
				var items []*attachmentpb.Attachment
				if resp != nil {
					items = resp.GetData()
				}
				pageData.AttachmentTable = attachment.BuildTable(items, cfg, id)
			}
			pageData.AttachmentUploadURL = route.ResolveURL(deps.Routes.AttachmentUploadURL, "id", id)

		case "audit-history":
			if deps.ListAuditHistory != nil {
				cursor := viewCtx.Request.URL.Query().Get("cursor")
				auditResp, err := deps.ListAuditHistory(ctx, &auditlog.ListAuditRequest{
					EntityType:  "price_list",
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

		// KB help content
		if viewCtx.Translations != nil {
			if provider, ok := viewCtx.Translations.(*lynguaV1.TranslationProvider); ok {
				if kb, _ := provider.LoadKBIfExists(viewCtx.Lang, viewCtx.BusinessType, "pricelist-detail"); kb != nil {
					pageData.HasHelp = true
					pageData.HelpContent = kb.Body
				}
			}
		}

		return view.OK("pricelist-detail", pageData)
	})
}

// NewTabAction creates the tab action view (partial — returns only the tab content).
// Handles GET /action/price-lists/{id}/tab/{tab}
func NewTabAction(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")
		tab := viewCtx.Request.PathValue("tab")
		if tab == "" {
			tab = "basic"
		}

		resp, err := deps.ReadPriceList(ctx, &pricelistpb.ReadPriceListRequest{
			Data: &pricelistpb.PriceList{Id: id},
		})
		if err != nil {
			log.Printf("Failed to read price list %s: %v", id, err)
			return view.Error(fmt.Errorf("failed to load price list: %w", err))
		}

		data := resp.GetData()
		if len(data) == 0 {
			return view.Error(fmt.Errorf("price list not found"))
		}
		priceList := data[0]

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion: viewCtx.CacheVersion,
				CommonLabels: deps.CommonLabels,
			},
			PriceList: priceList,
			ActiveTab: tab,
			TabItems:  buildTabItems(id, deps.Labels, deps.Routes),
			ID:        id,
			Labels:    deps.Labels,
		}

		switch tab {
		case "prices":
			perms := view.GetUserPermissions(ctx)
			pricesTable, err := buildPricesTable(ctx, deps, id, deps.Routes, perms)
			if err != nil {
				log.Printf("Failed to load price products for price list %s: %v", id, err)
			}
			pageData.PricesTable = pricesTable

		case "attachments":
			if deps.ListAttachments != nil {
				cfg := attachmentConfig(deps)
				resp, err := deps.ListAttachments(ctx, cfg.EntityType, id)
				if err != nil {
					log.Printf("Failed to list attachments for %s %s: %v", cfg.EntityType, id, err)
				}
				var items []*attachmentpb.Attachment
				if resp != nil {
					items = resp.GetData()
				}
				pageData.AttachmentTable = attachment.BuildTable(items, cfg, id)
			}
			pageData.AttachmentUploadURL = route.ResolveURL(deps.Routes.AttachmentUploadURL, "id", id)

		case "audit-history":
			if deps.ListAuditHistory != nil {
				cursor := viewCtx.Request.URL.Query().Get("cursor")
				auditResp, err := deps.ListAuditHistory(ctx, &auditlog.ListAuditRequest{
					EntityType:  "price_list",
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

		templateName := "pricelist-detail-" + tab
		if tab == "attachments" {
			templateName = "attachment-tab"
		}
		if tab == "audit-history" {
			templateName = "audit-history-tab"
		}
		return view.OK(templateName, pageData)
	})
}

func buildTabItems(id string, labels centymo.PriceListLabels, routes centymo.PriceListRoutes) []pyeza.TabItem {
	base := route.ResolveURL(routes.DetailURL, "id", id)
	action := route.ResolveURL(routes.TabActionURL, "id", id, "tab", "")
	return []pyeza.TabItem{
		{Key: "basic", Label: labels.Detail.BasicInfo, Href: base + "?tab=basic", HxGet: action + "basic"},
		{Key: "prices", Label: labels.Detail.Prices, Href: base + "?tab=prices", HxGet: action + "prices"},
		{Key: "attachments", Label: labels.Detail.TabAttachments, Href: base + "?tab=attachments", HxGet: action + "attachments", Icon: "icon-paperclip"},
		{Key: "audit-history", Label: "History", Href: base + "?tab=audit-history", HxGet: action + "audit-history", Icon: "icon-clock"},
	}
}

func buildPricesTable(ctx context.Context, deps *DetailViewDeps, priceListID string, routes centymo.PriceListRoutes, perms *types.UserPermissions) (*types.TableConfig, error) {
	resp, err := deps.ListPriceProducts(ctx, &priceproductpb.ListPriceProductsRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to list price products: %w", err)
	}

	l := deps.Labels
	columns := []types.TableColumn{
		{Key: "product_name", Label: l.Detail.ProductName, Sortable: true},
		{Key: "amount", Label: l.Detail.Amount, Sortable: true, WidthClass: "col-4xl"},
		{Key: "currency", Label: l.Detail.Currency, Sortable: true, WidthClass: "col-2xl"},
	}

	rows := []types.TableRow{}
	deleteURL := route.ResolveURL(routes.PriceProductDeleteURL, "id", priceListID)
	for _, pp := range resp.GetData() {
		// Filter price products belonging to this price list
		if pp.GetPriceListId() != priceListID {
			continue
		}
		id := pp.GetId()
		productName := pp.GetName()
		amount := fmt.Sprintf("%d", pp.GetAmount())
		currency := pp.GetCurrency()

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: productName},
				{Type: "text", Value: amount},
				{Type: "text", Value: currency},
			},
			Actions: []types.TableAction{
				{Type: "delete", Label: l.Detail.RemoveLabel, Action: "delete", URL: deleteURL, ItemName: productName, Disabled: !perms.Can("price_list", "delete"), DisabledTooltip: l.Errors.PermissionDenied},
			},
		})
	}
	types.ApplyColumnStyles(columns, rows)

	addURL := route.ResolveURL(routes.PriceProductAddURL, "id", priceListID)
	tableConfig := &types.TableConfig{
		ID:                   "price-products-table",
		Columns:              columns,
		Rows:                 rows,
		ShowSearch:           true,
		ShowActions:          true,
		DefaultSortColumn:    "product_name",
		DefaultSortDirection: "asc",
		Labels:               deps.TableLabels,
		EmptyState: types.TableEmptyState{
			Title:   l.Detail.EmptyTitle,
			Message: l.Detail.EmptyMessage,
		},
		PrimaryAction: &types.PrimaryAction{
			Label:           l.Detail.AddPrice,
			ActionURL:       addURL,
			Icon:            "icon-plus",
			Disabled:        !perms.Can("price_list", "create"),
			DisabledTooltip: l.Errors.PermissionDenied,
		},
	}
	types.ApplyTableSettings(tableConfig)

	return tableConfig, nil
}
