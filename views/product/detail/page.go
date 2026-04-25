package detail

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/erniealice/hybra-golang/views/attachment"
	"github.com/erniealice/hybra-golang/views/auditlog"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	centymo "github.com/erniealice/centymo-golang"
	lynguaV1 "github.com/erniealice/lyngua/golang/v1"

	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	linepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/line"
	productpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product"
	productlinepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_line"
	productoptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_option"
	productoptionvaluepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_option_value"
	productvariantpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_variant"
	productvariantoptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_variant_option"
)

// DetailViewDeps holds view dependencies.
type DetailViewDeps struct {
	Routes       centymo.ProductRoutes
	ReadProduct  func(ctx context.Context, req *productpb.ReadProductRequest) (*productpb.ReadProductResponse, error)
	Labels       centymo.ProductLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels
	// DataSource for backward compat
	DB centymo.DataSource

	// Typed proto funcs
	ListProductVariants       func(ctx context.Context, req *productvariantpb.ListProductVariantsRequest) (*productvariantpb.ListProductVariantsResponse, error)
	ListProductOptions        func(ctx context.Context, req *productoptionpb.ListProductOptionsRequest) (*productoptionpb.ListProductOptionsResponse, error)
	ListProductOptionValues   func(ctx context.Context, req *productoptionvaluepb.ListProductOptionValuesRequest) (*productoptionvaluepb.ListProductOptionValuesResponse, error)
	ListLines                 func(ctx context.Context, req *linepb.ListLinesRequest) (*linepb.ListLinesResponse, error)
	ListProductLines          func(ctx context.Context, req *productlinepb.ListProductLinesRequest) (*productlinepb.ListProductLinesResponse, error)
	CreateProductLine         func(ctx context.Context, req *productlinepb.CreateProductLineRequest) (*productlinepb.CreateProductLineResponse, error)
	UpdateProductLine         func(ctx context.Context, req *productlinepb.UpdateProductLineRequest) (*productlinepb.UpdateProductLineResponse, error)
	DeleteProductLine         func(ctx context.Context, req *productlinepb.DeleteProductLineRequest) (*productlinepb.DeleteProductLineResponse, error)
	ListProductVariantOptions func(ctx context.Context, req *productvariantoptionpb.ListProductVariantOptionsRequest) (*productvariantoptionpb.ListProductVariantOptionsResponse, error)

	// PermissionEntity is the first argument to perms.Can(entity, action) for
	// the detail-page action buttons (edit, delete, variant assign). Defaults
	// to "product". See centymo-golang/views/product/module.go ModuleDeps.
	PermissionEntity string

	attachment.AttachmentOps
	auditlog.AuditOps
}

// permEntity returns the configured PermissionEntity with a safe default.
func (d *DetailViewDeps) permEntity() string {
	if d == nil || d.PermissionEntity == "" {
		return "product"
	}
	return d.PermissionEntity
}

// PageData holds the data for the product detail page.
type PageData struct {
	types.PageData
	ContentTemplate     string
	Product             *productpb.Product
	Labels              centymo.ProductLabels
	ActiveTab           string
	TabItems            []pyeza.TabItem
	ID                  string
	ProductName         string
	ProductDesc         string
	ProductPrice        string
	ProductCurrency     string
	ProductStatus       string
	StatusVariant       string
	LineName            string // resolved name of the product's primary line (from product.line_id)
	// Model D — display-formatted unit of measure and variant mode for the Info tab.
	ProductUnit        string
	ProductVariantMode string
	// Label fallbacks surfaced to the template when lyngua doesn't supply
	// Detail.Unit / Detail.VariantMode. Kept on PageData so the template can
	// render a consistent label without sprinkling fallback strings in HTML.
	UnitRowLabel        string
	VariantModeRowLabel string
	VariantsTable       *types.TableConfig
	OptionsTable        *types.TableConfig
	LinesTable          *types.TableConfig
	AttachmentTable     *types.TableConfig
	AttachmentUploadURL string
	// Audit history tab
	AuditEntries    []auditlog.AuditEntryView
	AuditHasNext    bool
	AuditNextCursor string
	AuditHistoryURL string
}

type ProductLineFormLabels struct {
	Title           string
	Line            string
	LinePlaceholder string
	SortOrder       string
	Active          string
}

type ProductLineFormData struct {
	FormAction   string
	IsEdit       bool
	ID           string
	ProductID    string
	LineID       string
	SortOrder    string
	Active       bool
	LineOptions  []types.SelectOption
	Labels       ProductLineFormLabels
	CommonLabels pyeza.CommonLabels
}

// NewView creates the product detail view (full page).
func NewView(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")

		activeTab := viewCtx.Request.URL.Query().Get("tab")
		if activeTab == "" || activeTab == "pricing" {
			activeTab = "info"
		}

		pageData, err := buildPageData(ctx, deps, id, activeTab, viewCtx)
		if err != nil {
			return view.Error(err)
		}

		// KB help content
		if viewCtx.Translations != nil {
			if provider, ok := viewCtx.Translations.(*lynguaV1.TranslationProvider); ok {
				if kb, _ := provider.LoadKBIfExists(viewCtx.Lang, viewCtx.BusinessType, "product-detail"); kb != nil {
					pageData.HasHelp = true
					pageData.HelpContent = kb.Body
				}
			}
		}

		return view.OK("product-detail", pageData)
	})
}

// NewTabAction creates the tab action view (partial — returns only the tab content).
// Handles GET /action/products/detail/{id}/tab/{tab}
func NewTabAction(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")
		tab := viewCtx.Request.PathValue("tab")
		if tab == "" {
			tab = "info"
		}

		if tab == "lines" {
			switch viewCtx.Request.URL.Query().Get("mode") {
			case "add":
				return handleProductLineAssociationAdd(ctx, deps, viewCtx, id)
			case "edit":
				return handleProductLineAssociationEdit(ctx, deps, viewCtx, id)
			case "delete":
				return handleProductLineAssociationDelete(ctx, deps, viewCtx, id)
			}
		}

		pageData, err := buildPageData(ctx, deps, id, tab, viewCtx)
		if err != nil {
			return view.Error(err)
		}

		// The response partial swaps the clicked tab's body into #tabContent
		// and OOB-swaps the tab bar so its .active class updates to the
		// newly-selected tab. Without the OOB swap, the bar stays frozen on
		// whichever tab was active at full-page render.
		return view.OK("product-detail-tab-response", pageData)
	})
}

// buildPageData loads product data and builds the PageData for the given active tab.
func buildPageData(ctx context.Context, deps *DetailViewDeps, id, activeTab string, viewCtx *view.ViewContext) (*PageData, error) {
	resp, err := deps.ReadProduct(ctx, &productpb.ReadProductRequest{
		Data: &productpb.Product{Id: id},
	})
	if err != nil {
		log.Printf("Failed to read product %s: %v", id, err)
		return nil, fmt.Errorf("failed to load product: %w", err)
	}

	data := resp.GetData()
	if len(data) == 0 {
		return nil, fmt.Errorf("product not found")
	}
	product := data[0]

	name := product.GetName()
	description := product.GetDescription()
	currency := product.GetCurrency()
	if currency == "" {
		currency = "PHP"
	}
	// Model D — Product.price is optional. When variant_mode = "configurable",
	// per-SKU pricing lives on ProductVariant.price_override and the detail
	// card surfaces a localized "Varies" string instead of a number.
	// Phase D will add the VariantPriceVaries label; the key is referenced
	// here with a sensible English fallback for now.
	var priceFormatted string
	if product.GetVariantMode() == "configurable" {
		priceFormatted = deps.Labels.Form.VariantPriceVaries
		if priceFormatted == "" {
			priceFormatted = "Varies"
		}
	} else if product.Price != nil {
		priceFormatted = FormatPrice(currency, float64(product.GetPrice())/100.0)
	} else {
		priceFormatted = ""
	}

	productStatus := "active"
	if !product.GetActive() {
		productStatus = "inactive"
	}
	StatusVariant := "success"
	if productStatus == "inactive" {
		StatusVariant = "warning"
	}

	// Resolve the product's primary line name from product.line_id.
	productLineName := ""
	if lineID := product.GetLineId(); lineID != "" && deps.ListLines != nil {
		if lineResp, lerr := deps.ListLines(ctx, &linepb.ListLinesRequest{}); lerr == nil {
			for _, line := range lineResp.GetData() {
				if line != nil && line.GetId() == lineID {
					productLineName = line.GetName()
					break
				}
			}
		}
	}

	// Get counts for tab badges
	variantCount := 0
	if deps.ListProductVariants != nil {
		varResp, err := deps.ListProductVariants(ctx, &productvariantpb.ListProductVariantsRequest{})
		if err == nil {
			for _, v := range varResp.GetData() {
				if v.GetProductId() == id {
					variantCount++
				}
			}
		}
	}

	optionCount := getOptionCountTyped(ctx, deps, id)
	lineCount, linesTable := buildLinesTable(ctx, deps, id)

	// Model D — Info tab display strings for unit + variant mode.
	// Falls back to English when lyngua hasn't overlaid the keys yet.
	productUnit := product.GetUnit()
	unitRowLabel := deps.Labels.Detail.Unit
	if unitRowLabel == "" {
		// Form.UnitLabel is already loaded from lyngua (see Phase D) — reuse as
		// a working fallback until Detail.Unit translations land.
		unitRowLabel = deps.Labels.Form.UnitLabel
	}
	if unitRowLabel == "" {
		unitRowLabel = "Unit of measure"
	}
	variantModeRowLabel := deps.Labels.Detail.VariantMode
	if variantModeRowLabel == "" {
		variantModeRowLabel = deps.Labels.Form.VariantModeLabel
	}
	if variantModeRowLabel == "" {
		variantModeRowLabel = "Variant Mode"
	}
	var productVariantModeDisplay string
	switch product.GetVariantMode() {
	case "configurable":
		productVariantModeDisplay = deps.Labels.Form.VariantModeConfigurable
		if productVariantModeDisplay == "" {
			productVariantModeDisplay = "Configurable"
		}
	default:
		productVariantModeDisplay = deps.Labels.Form.VariantModeNone
		if productVariantModeDisplay == "" {
			productVariantModeDisplay = "Simple"
		}
	}

	l := deps.Labels
	// Model D — gate the Options and Variants tabs on Product.variant_mode.
	// Simple products omit them entirely so the detail page stays clean.
	//
	// Bug 1 safeguard: Phase A's migration backfilled every pre-existing row
	// with variant_mode="none" regardless of whether the product already had
	// product_option / product_variant rows. Legacy / mis-migrated products
	// would otherwise lose access to their existing configuration. Show the
	// tabs when the product is configurable OR has any existing options/variants.
	showVariantTabs := product.GetVariantMode() == "configurable" || optionCount > 0 || variantCount > 0
	tabItems := buildTabItems(id, l, variantCount, optionCount, lineCount, deps.Routes, showVariantTabs)

	// Header subtitle: use the product description, or fall back to the
	// "No description provided" lyngua label. Without this fallback the
	// shared header template uses CommonLabels.Header.WelcomeBack, which
	// reads as "Welcome back" on a detail page.
	headerSubtitle := description
	if headerSubtitle == "" {
		headerSubtitle = deps.Labels.Detail.NoDescriptionSubtitle
		if headerSubtitle == "" {
			headerSubtitle = "No description provided"
		}
	}

	pageData := &PageData{
		PageData: types.PageData{
			CacheVersion:   viewCtx.CacheVersion,
			Title:          name,
			CurrentPath:    viewCtx.CurrentPath,
			ActiveNav:      deps.Routes.ActiveNav,
			ActiveSubNav:   deps.Routes.ActiveSubNav,
			HeaderTitle:    name,
			HeaderSubtitle: headerSubtitle,
			HeaderIcon:     "icon-package",
			CommonLabels:   deps.CommonLabels,
		},
		ContentTemplate: "product-detail-content",
		Product:         product,
		Labels:          l,
		ActiveTab:       activeTab,
		TabItems:        tabItems,
		ID:              id,
		ProductName:     name,
		ProductDesc:     description,
		ProductPrice:    priceFormatted,
		ProductCurrency: currency,
		ProductStatus:       productStatus,
		StatusVariant:       StatusVariant,
		LineName:            productLineName,
		ProductUnit:         productUnit,
		ProductVariantMode:  productVariantModeDisplay,
		UnitRowLabel:        unitRowLabel,
		VariantModeRowLabel: variantModeRowLabel,
	}

	// Load tab-specific data
	perms := view.GetUserPermissions(ctx)
	switch activeTab {
	case "variants":
		tableConfig := BuildVariantsTable(ctx, deps, id, perms)
		pageData.VariantsTable = tableConfig
	case "options":
		tableConfig := buildOptionsTable(ctx, deps, id)
		pageData.OptionsTable = tableConfig
	case "lines":
		pageData.LinesTable = linesTable
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
				EntityType:  "product",
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

	return pageData, nil
}

func buildTabItems(id string, l centymo.ProductLabels, variantCount, optionCount, lineCount int, routes centymo.ProductRoutes, showVariantTabs bool) []pyeza.TabItem {
	base := route.ResolveURL(routes.DetailURL, "id", id)
	action := route.ResolveURL(routes.TabActionURL, "id", id, "tab", "")
	items := []pyeza.TabItem{
		{Key: "info", Label: l.Tabs.Info, Href: base + "?tab=info", HxGet: action + "info", Icon: "icon-info", Count: 0, Disabled: false},
	}
	// Options + Variants tabs render when the product is variant-configurable
	// OR has any existing option/variant rows (defensive for legacy rows whose
	// variant_mode was not inferred during the Phase A migration backfill).
	if showVariantTabs {
		items = append(items,
			pyeza.TabItem{Key: "options", Label: l.Tabs.Options, Href: base + "?tab=options", HxGet: action + "options", Icon: "icon-settings", Count: optionCount, Disabled: false},
			pyeza.TabItem{Key: "variants", Label: l.Tabs.Variants, Href: base + "?tab=variants", HxGet: action + "variants", Icon: "icon-layers", Count: variantCount, Disabled: false},
		)
	}
	items = append(items,
		pyeza.TabItem{Key: "lines", Label: "Lines", Href: base + "?tab=lines", HxGet: action + "lines", Icon: "icon-layers", Count: lineCount, Disabled: false},
		pyeza.TabItem{Key: "attachments", Label: l.Tabs.Attachments, Href: base + "?tab=attachments", HxGet: action + "attachments", Icon: "icon-paperclip", Count: 0, Disabled: false},
		pyeza.TabItem{Key: "audit-history", Label: "History", Href: base + "?tab=audit-history", HxGet: action + "audit-history", Icon: "icon-clock"},
	)
	return items
}

func buildLinesTable(ctx context.Context, deps *DetailViewDeps, productID string) (int, *types.TableConfig) {
	if deps.ListLines == nil || deps.ListProductLines == nil {
		return 0, nil
	}

	lineResp, err := deps.ListLines(ctx, &linepb.ListLinesRequest{})
	if err != nil {
		log.Printf("Failed to list product lines: %v", err)
		return 0, nil
	}
	lineNameByID := map[string]string{}
	for _, line := range lineResp.GetData() {
		if line != nil {
			lineNameByID[line.GetId()] = line.GetName()
		}
	}

	assocResp, err := deps.ListProductLines(ctx, &productlinepb.ListProductLinesRequest{})
	if err != nil {
		log.Printf("Failed to list product line associations: %v", err)
		return 0, nil
	}

	perms := view.GetUserPermissions(ctx)
	rows := []types.TableRow{}
	for _, assoc := range assocResp.GetData() {
		if assoc == nil || assoc.GetProductId() != productID {
			continue
		}
		lineID := assoc.GetLineId()
		lineName := lineNameByID[lineID]
		if lineName == "" {
			lineName = lineID
		}
		sortOrder := ""
		if assoc.GetSortOrder() != 0 {
			sortOrder = fmt.Sprintf("%d", assoc.GetSortOrder())
		}
		rows = append(rows, types.TableRow{
			ID: assoc.GetId(),
			Cells: []types.TableCell{
				{Type: "text", Value: lineName},
				{Type: "text", Value: lineID},
				{Type: "text", Value: sortOrder},
			},
			DataAttrs: map[string]string{
				"line_id": lineID,
			},
			Actions: []types.TableAction{
				{
					Type:   "view",
					Label:  deps.Labels.Actions.View,
					Action: "view",
					Href:   route.ResolveURL(deps.Routes.DetailURL, "id", lineID),
				},
				{
					Type:            "edit",
					Label:           deps.Labels.Actions.Edit,
					Action:          "edit",
					URL:             route.ResolveURL(deps.Routes.TabActionURL, "id", productID, "tab", "lines") + "?mode=edit&plid=" + assoc.GetId(),
					DrawerTitle:     "Edit Line",
					Disabled:        !perms.Can("product_line", "update"),
					DisabledTooltip: deps.Labels.Errors.PermissionDenied,
				},
				{
					Type:            "delete",
					Label:           deps.Labels.Actions.Delete,
					Action:          "delete",
					URL:             route.ResolveURL(deps.Routes.TabActionURL, "id", productID, "tab", "lines") + "?mode=delete&plid=" + assoc.GetId(),
					ItemName:        lineName,
					ConfirmTitle:    deps.Labels.Actions.Delete,
					ConfirmMessage:  fmt.Sprintf("Remove %s from this product?", lineName),
					Disabled:        !perms.Can("product_line", "delete"),
					DisabledTooltip: deps.Labels.Errors.PermissionDenied,
				},
			},
		})
	}

	columns := []types.TableColumn{
		{Key: "name", Label: "Line", Sortable: false},
		{Key: "line_id", Label: "Line ID", Sortable: false},
		{Key: "sort_order", Label: "Sort Order", Sortable: false, WidthClass: "col-2xl"},
	}
	types.ApplyColumnStyles(columns, rows)

	table := &types.TableConfig{
		ID:          "product-lines-table",
		Columns:     columns,
		Rows:        rows,
		ShowActions: true,
		ShowEntries: true,
		PrimaryAction: &types.PrimaryAction{
			Label:           "Add Line",
			ActionURL:       route.ResolveURL(deps.Routes.TabActionURL, "id", productID, "tab", "lines") + "?mode=add",
			Icon:            "icon-plus",
			Disabled:        !perms.Can("product_line", "create"),
			DisabledTooltip: deps.Labels.Errors.PermissionDenied,
		},
	}
	types.ApplyTableSettings(table)
	return len(rows), table
}

func handleProductLineAssociationAdd(ctx context.Context, deps *DetailViewDeps, viewCtx *view.ViewContext, productID string) view.ViewResult {
	perms := view.GetUserPermissions(ctx)
	if !perms.Can("product_line", "create") {
		return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
	}
	if deps.CreateProductLine == nil {
		return centymo.HTMXError("Product line create is not available")
	}

	if viewCtx.Request.Method == http.MethodGet {
		return view.OK("product-line-association-drawer-form", &ProductLineFormData{
			FormAction:   route.ResolveURL(deps.Routes.TabActionURL, "id", productID, "tab", "lines") + "?mode=add",
			ProductID:    productID,
			SortOrder:    "0",
			Active:       true,
			LineOptions:  loadLineOptions(ctx, deps, ""),
			Labels:       productLineFormLabels("Add Line"),
			CommonLabels: deps.CommonLabels,
		})
	}

	if err := viewCtx.Request.ParseForm(); err != nil {
		return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
	}

	lineID := viewCtx.Request.FormValue("line_id")
	if lineID == "" {
		return centymo.HTMXError("Line is required")
	}

	productLine := &productlinepb.ProductLine{
		ProductId: productID,
		LineId:    lineID,
		Active:    true,
	}
	if v := viewCtx.Request.FormValue("sort_order"); v != "" {
		if so, err := strconv.ParseInt(v, 10, 32); err == nil {
			productLine.SortOrder = int32(so)
		}
	}

	if _, err := deps.CreateProductLine(ctx, &productlinepb.CreateProductLineRequest{Data: productLine}); err != nil {
		log.Printf("Failed to create product line association for product %s: %v", productID, err)
		return centymo.HTMXError(err.Error())
	}

	return centymo.HTMXSuccess("product-lines-table")
}

func handleProductLineAssociationEdit(ctx context.Context, deps *DetailViewDeps, viewCtx *view.ViewContext, productID string) view.ViewResult {
	perms := view.GetUserPermissions(ctx)
	if !perms.Can("product_line", "update") {
		return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
	}
	if deps.UpdateProductLine == nil {
		return centymo.HTMXError("Product line update is not available")
	}

	assocID := viewCtx.Request.URL.Query().Get("plid")
	if assocID == "" {
		assocID = viewCtx.Request.FormValue("plid")
	}
	if assocID == "" {
		return centymo.HTMXError(deps.Labels.Errors.IDRequired)
	}

	assoc, err := findProductLineAssociation(ctx, deps, assocID)
	if err != nil {
		return centymo.HTMXError(err.Error())
	}

	if viewCtx.Request.Method == http.MethodGet {
		return view.OK("product-line-association-drawer-form", &ProductLineFormData{
			FormAction:   route.ResolveURL(deps.Routes.TabActionURL, "id", productID, "tab", "lines") + "?mode=edit&plid=" + assocID,
			IsEdit:       true,
			ID:           assocID,
			ProductID:    productID,
			LineID:       assoc.GetLineId(),
			SortOrder:    fmt.Sprintf("%d", assoc.GetSortOrder()),
			Active:       assoc.GetActive(),
			LineOptions:  loadLineOptions(ctx, deps, assoc.GetLineId()),
			Labels:       productLineFormLabels("Edit Line"),
			CommonLabels: deps.CommonLabels,
		})
	}

	if err := viewCtx.Request.ParseForm(); err != nil {
		return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
	}

	lineID := viewCtx.Request.FormValue("line_id")
	if lineID == "" {
		return centymo.HTMXError("Line is required")
	}

	updated := &productlinepb.ProductLine{
		Id:        assocID,
		ProductId: productID,
		LineId:    lineID,
		Active:    viewCtx.Request.FormValue("active") == "true",
	}
	if v := viewCtx.Request.FormValue("sort_order"); v != "" {
		if so, err := strconv.ParseInt(v, 10, 32); err == nil {
			updated.SortOrder = int32(so)
		}
	}

	if _, err := deps.UpdateProductLine(ctx, &productlinepb.UpdateProductLineRequest{Data: updated}); err != nil {
		log.Printf("Failed to update product line association %s: %v", assocID, err)
		return centymo.HTMXError(err.Error())
	}

	return centymo.HTMXSuccess("product-lines-table")
}

func handleProductLineAssociationDelete(ctx context.Context, deps *DetailViewDeps, viewCtx *view.ViewContext, productID string) view.ViewResult {
	perms := view.GetUserPermissions(ctx)
	if !perms.Can("product_line", "delete") {
		return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
	}
	if deps.DeleteProductLine == nil {
		return centymo.HTMXError("Product line delete is not available")
	}

	assocID := viewCtx.Request.URL.Query().Get("plid")
	if assocID == "" {
		_ = viewCtx.Request.ParseForm()
		assocID = viewCtx.Request.FormValue("plid")
	}
	if assocID == "" {
		return centymo.HTMXError(deps.Labels.Errors.IDRequired)
	}

	if _, err := deps.DeleteProductLine(ctx, &productlinepb.DeleteProductLineRequest{Data: &productlinepb.ProductLine{Id: assocID}}); err != nil {
		log.Printf("Failed to delete product line association %s for product %s: %v", assocID, productID, err)
		return centymo.HTMXError(err.Error())
	}

	return centymo.HTMXSuccess("product-lines-table")
}

func findProductLineAssociation(ctx context.Context, deps *DetailViewDeps, assocID string) (*productlinepb.ProductLine, error) {
	if deps.ListProductLines == nil {
		return nil, fmt.Errorf("product line associations not available")
	}

	resp, err := deps.ListProductLines(ctx, &productlinepb.ListProductLinesRequest{})
	if err != nil {
		log.Printf("Failed to list product line associations: %v", err)
		return nil, fmt.Errorf("failed to load product line associations")
	}

	for _, assoc := range resp.GetData() {
		if assoc != nil && assoc.GetId() == assocID {
			return assoc, nil
		}
	}

	return nil, fmt.Errorf("product line association not found")
}

func loadLineOptions(ctx context.Context, deps *DetailViewDeps, selectedID string) []types.SelectOption {
	if deps.ListLines == nil {
		return nil
	}

	resp, err := deps.ListLines(ctx, &linepb.ListLinesRequest{})
	if err != nil {
		log.Printf("Failed to load lines for product association form: %v", err)
		return nil
	}

	options := make([]types.SelectOption, 0, len(resp.GetData()))
	for _, line := range resp.GetData() {
		if line == nil {
			continue
		}
		if !line.GetActive() && line.GetId() != selectedID {
			continue
		}
		label := line.GetName()
		if label == "" {
			label = line.GetId()
		} else if line.GetId() != "" {
			label = fmt.Sprintf("%s (%s)", label, line.GetId())
		}
		options = append(options, types.SelectOption{
			Value:    line.GetId(),
			Label:    label,
			Selected: line.GetId() == selectedID,
		})
	}

	return options
}

func productLineFormLabels(title string) ProductLineFormLabels {
	return ProductLineFormLabels{
		Title:           title,
		Line:            "Line",
		LinePlaceholder: "Select a line",
		SortOrder:       "Sort Order",
		Active:          "Active",
	}
}

// ---------------------------------------------------------------------------
// Variants tab table
// ---------------------------------------------------------------------------

func BuildVariantsTable(ctx context.Context, deps *DetailViewDeps, productID string, perms *types.UserPermissions) *types.TableConfig {
	l := deps.Labels

	columns := []types.TableColumn{
		{Key: "sku", Label: l.Variant.SKU, Sortable: true},
		{Key: "priceOverride", Label: l.Variant.PriceOverride, Sortable: true, WidthClass: "col-4xl"},
		{Key: "options", Label: l.Detail.OptionsLabel, Sortable: false},
		{Key: "status", Label: l.Columns.Status, Sortable: true, WidthClass: "col-2xl"},
	}

	rows := []types.TableRow{}

	// Build variant → option value labels map
	variantOptionLabels := make(map[string][]string)
	if deps.ListProductVariantOptions != nil && deps.ListProductOptionValues != nil {
		voResp, _ := deps.ListProductVariantOptions(ctx, &productvariantoptionpb.ListProductVariantOptionsRequest{})
		ovResp, _ := deps.ListProductOptionValues(ctx, &productoptionvaluepb.ListProductOptionValuesRequest{})
		valLabelMap := make(map[string]string)
		if ovResp != nil {
			for _, ov := range ovResp.GetData() {
				vid := ov.GetId()
				label := ov.GetLabel()
				if vid != "" {
					valLabelMap[vid] = label
				}
			}
		}
		if voResp != nil {
			for _, vo := range voResp.GetData() {
				varID := vo.GetProductVariantId()
				valID := vo.GetProductOptionValueId()
				if label, ok := valLabelMap[valID]; ok && varID != "" {
					variantOptionLabels[varID] = append(variantOptionLabels[varID], label)
				}
			}
		}
	}

	if deps.ListProductVariants != nil {
		varResp, err := deps.ListProductVariants(ctx, &productvariantpb.ListProductVariantsRequest{})
		if err != nil {
			log.Printf("Failed to list product variants: %v", err)
		} else {
			for _, v := range varResp.GetData() {
				pid := v.GetProductId()
				if pid != productID {
					continue
				}

				vid := v.GetId()
				sku := v.GetSku()
				priceStr := ""
				if po := v.GetPriceOverride(); po != 0 {
					priceStr = fmt.Sprintf("%v", po)
				}
				active := v.GetActive()
				status := "active"
				if !active {
					status = "inactive"
				}

				optionsDisplay := strings.Join(variantOptionLabels[vid], centymo.OptionValueSeparator)

				actions := []types.TableAction{
					{
						Type: "view", Label: l.Actions.View,
						Href: route.ResolveURL(deps.Routes.VariantDetailURL, "id", productID, "vid", vid),
					},
					{
						Type: "edit", Label: l.Variant.Edit, Action: "edit",
						URL:             route.ResolveURL(deps.Routes.VariantEditURL, "id", productID, "vid", vid),
						DrawerTitle:     l.Variant.Edit,
						Disabled:        !perms.Can(deps.permEntity(), "update"),
						DisabledTooltip: l.Errors.PermissionDenied,
					},
					{
						Type: "delete", Label: l.Variant.Remove, Action: "delete",
						URL:             route.ResolveURL(deps.Routes.VariantRemoveURL, "id", productID),
						ItemName:        sku,
						ConfirmTitle:    l.Variant.Remove,
						ConfirmMessage:  fmt.Sprintf(l.Confirm.DeactivateMessage, sku),
						Disabled:        !perms.Can(deps.permEntity(), "delete"),
						DisabledTooltip: l.Errors.PermissionDenied,
					},
				}

				rows = append(rows, types.TableRow{
					ID: vid,
					Cells: []types.TableCell{
						{Type: "text", Value: sku},
						{Type: "text", Value: priceStr},
						{Type: "text", Value: optionsDisplay},
						{Type: "badge", Value: status, Variant: StatusVariant(status)},
					},
					DataAttrs: map[string]string{
						"sku":    sku,
						"status": status,
					},
					Actions: actions,
				})
			}
		}
	}

	types.ApplyColumnStyles(columns, rows)

	tableConfig := &types.TableConfig{
		ID:                   "product-variants-table",
		RefreshURL:           route.ResolveURL(deps.Routes.VariantTableURL, "id", productID),
		Columns:              columns,
		Rows:                 rows,
		ShowSearch:           true,
		ShowActions:          true,
		ShowFilters:          false,
		ShowSort:             true,
		ShowColumns:          true,
		ShowExport:           false,
		ShowDensity:          true,
		ShowEntries:          true,
		DefaultSortColumn:    "sku",
		DefaultSortDirection: "asc",
		Labels:               deps.TableLabels,
		EmptyState: types.TableEmptyState{
			Title:   l.Variant.Empty,
			Message: l.Detail.EmptyVariantsMessage,
		},
		PrimaryAction: &types.PrimaryAction{
			Label:           l.Variant.Assign,
			ActionURL:       route.ResolveURL(deps.Routes.VariantAssignURL, "id", productID),
			Icon:            "icon-plus",
			Disabled:        !perms.Can(deps.permEntity(), "create"),
			DisabledTooltip: l.Errors.PermissionDenied,
		},
	}
	types.ApplyTableSettings(tableConfig)

	return tableConfig
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func FormatPrice(currency string, price float64) string {
	if currency == "" {
		currency = "PHP"
	}
	raw := fmt.Sprintf("%.2f", price)
	parts := strings.SplitN(raw, ".", 2)
	intPart := parts[0]
	decPart := "00"
	if len(parts) > 1 {
		decPart = parts[1]
	}

	n := len(intPart)
	if n <= 3 {
		return currency + " " + intPart + "." + decPart
	}
	var result []byte
	for i, c := range intPart {
		if i > 0 && (n-i)%3 == 0 {
			result = append(result, ',')
		}
		result = append(result, byte(c))
	}
	return currency + " " + string(result) + "." + decPart
}

func StatusVariant(status string) string {
	switch status {
	case "active":
		return "success"
	case "inactive":
		return "warning"
	default:
		return "default"
	}
}
