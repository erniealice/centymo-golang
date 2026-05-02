package revenue

import (
	"context"
	"net/http"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	"github.com/erniealice/centymo-golang"
	revenueaction "github.com/erniealice/centymo-golang/views/revenue/action"
	revenuedashboard "github.com/erniealice/centymo-golang/views/revenue/dashboard"
	revenuedetail "github.com/erniealice/centymo-golang/views/revenue/detail"
	revenuelist "github.com/erniealice/centymo-golang/views/revenue/list"
	revenuepayment "github.com/erniealice/centymo-golang/views/revenue/payment"
	revenuesearch "github.com/erniealice/centymo-golang/views/revenue/search"
	revenuesettings "github.com/erniealice/centymo-golang/views/revenue/settings"
	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	documenttemplatepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/template"
	clientpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client"
	inventoryitempb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/inventory_item"
	inventoryserialpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/inventory_serial"
	serialhistorypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/serial_history"
	jobactivitypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/operation/job_activity"
	pricelistpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/price_list"
	priceproductpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/price_product"
	productpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product"
	revenuepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/revenue"
	revenuelineitempb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/revenue_line_item"
	priceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_plan"
	productpriceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/product_price_plan"
	subscriptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription"
	"github.com/erniealice/hybra-golang/views/attachment"
	"github.com/erniealice/hybra-golang/views/auditlog"
)

// PaymentTermOption is re-exported from action for use by callers wiring ModuleDeps.
type PaymentTermOption = revenueaction.PaymentTermOption

// ModuleDeps holds all dependencies for the revenue module.
type ModuleDeps struct {
	Routes          centymo.RevenueRoutes
	DB              centymo.DataSource // KEEP — used for revenue_payment, collection_method, location
	GetListPageData func(ctx context.Context, req *revenuepb.GetRevenueListPageDataRequest) (*revenuepb.GetRevenueListPageDataResponse, error)
	Labels          centymo.RevenueLabels
	CommonLabels    pyeza.CommonLabels
	TableLabels     types.TableLabels

	// Payment terms dropdown (optional — gracefully degrades when nil)
	ListPaymentTerms func(ctx context.Context) ([]*PaymentTermOption, error)

	// Client search for autocomplete (optional — gracefully degrades when nil)
	ListClients         func(ctx context.Context, req *clientpb.ListClientsRequest) (*clientpb.ListClientsResponse, error)
	SearchClientsByName func(ctx context.Context, req *clientpb.SearchClientsByNameRequest) (*clientpb.SearchClientsByNameResponse, error)

	// Subscription search for revenue form autocomplete (optional — gracefully degrades when nil)
	ListSubscriptions func(ctx context.Context, req *subscriptionpb.ListSubscriptionsRequest) (*subscriptionpb.ListSubscriptionsResponse, error)

	// Subscription auto-populate (optional — gracefully degrades when nil)
	ReadSubscription      func(ctx context.Context, req *subscriptionpb.ReadSubscriptionRequest) (*subscriptionpb.ReadSubscriptionResponse, error)
	ReadPricePlan         func(ctx context.Context, req *priceplanpb.ReadPricePlanRequest) (*priceplanpb.ReadPricePlanResponse, error)
	ListProductPricePlans func(ctx context.Context, req *productpriceplanpb.ListProductPricePlansRequest) (*productpriceplanpb.ListProductPricePlansResponse, error)
	ReadProduct           func(ctx context.Context, req *productpb.ReadProductRequest) (*productpb.ReadProductResponse, error)
	ListProducts          func(ctx context.Context, req *productpb.ListProductsRequest) (*productpb.ListProductsResponse, error)

	// Typed revenue operations (for detail + action views)
	CreateRevenue func(ctx context.Context, req *revenuepb.CreateRevenueRequest) (*revenuepb.CreateRevenueResponse, error)
	ReadRevenue   func(ctx context.Context, req *revenuepb.ReadRevenueRequest) (*revenuepb.ReadRevenueResponse, error)
	UpdateRevenue func(ctx context.Context, req *revenuepb.UpdateRevenueRequest) (*revenuepb.UpdateRevenueResponse, error)
	DeleteRevenue func(ctx context.Context, req *revenuepb.DeleteRevenueRequest) (*revenuepb.DeleteRevenueResponse, error)

	// Typed line item operations
	CreateRevenueLineItem func(ctx context.Context, req *revenuelineitempb.CreateRevenueLineItemRequest) (*revenuelineitempb.CreateRevenueLineItemResponse, error)
	ReadRevenueLineItem   func(ctx context.Context, req *revenuelineitempb.ReadRevenueLineItemRequest) (*revenuelineitempb.ReadRevenueLineItemResponse, error)
	UpdateRevenueLineItem func(ctx context.Context, req *revenuelineitempb.UpdateRevenueLineItemRequest) (*revenuelineitempb.UpdateRevenueLineItemResponse, error)
	DeleteRevenueLineItem func(ctx context.Context, req *revenuelineitempb.DeleteRevenueLineItemRequest) (*revenuelineitempb.DeleteRevenueLineItemResponse, error)
	ListRevenueLineItems  func(ctx context.Context, req *revenuelineitempb.ListRevenueLineItemsRequest) (*revenuelineitempb.ListRevenueLineItemsResponse, error)

	// Typed inventory operations (for action views — stock deduction on status change)
	ReadInventoryItem            func(ctx context.Context, req *inventoryitempb.ReadInventoryItemRequest) (*inventoryitempb.ReadInventoryItemResponse, error)
	UpdateInventoryItem          func(ctx context.Context, req *inventoryitempb.UpdateInventoryItemRequest) (*inventoryitempb.UpdateInventoryItemResponse, error)
	ListInventoryItems           func(ctx context.Context, req *inventoryitempb.ListInventoryItemsRequest) (*inventoryitempb.ListInventoryItemsResponse, error)
	UpdateInventorySerial        func(ctx context.Context, req *inventoryserialpb.UpdateInventorySerialRequest) (*inventoryserialpb.UpdateInventorySerialResponse, error)
	CreateInventorySerialHistory func(ctx context.Context, req *serialhistorypb.CreateInventorySerialHistoryRequest) (*serialhistorypb.CreateInventorySerialHistoryResponse, error)

	// Document generation (wraps fycha.DocumentService.ProcessBytes)
	GenerateDoc func(templateData []byte, data map[string]any) ([]byte, error)

	// Optional: load custom default template from storage
	LoadDefaultTemplate func(ctx context.Context, purpose string) ([]byte, error)

	// Document template CRUD operations
	ListDocumentTemplates  func(ctx context.Context, req *documenttemplatepb.ListDocumentTemplatesRequest) (*documenttemplatepb.ListDocumentTemplatesResponse, error)
	CreateDocumentTemplate func(ctx context.Context, req *documenttemplatepb.CreateDocumentTemplateRequest) (*documenttemplatepb.CreateDocumentTemplateResponse, error)
	UpdateDocumentTemplate func(ctx context.Context, req *documenttemplatepb.UpdateDocumentTemplateRequest) (*documenttemplatepb.UpdateDocumentTemplateResponse, error)
	DeleteDocumentTemplate func(ctx context.Context, req *documenttemplatepb.DeleteDocumentTemplateRequest) (*documenttemplatepb.DeleteDocumentTemplateResponse, error)

	// Storage operations for template file upload
	UploadTemplate func(ctx context.Context, bucketName, objectKey string, content []byte, contentType string) error

	// Email sending for invoice delivery
	SendEmail func(ctx context.Context, to []string, subject, htmlBody, textBody string, attachmentName string, attachmentData []byte) error

	// Attachment operations
	UploadFile       func(ctx context.Context, bucket, key string, content []byte, contentType string) error
	ListAttachments  func(ctx context.Context, moduleKey, foreignKey string) (*attachmentpb.ListAttachmentsResponse, error)
	CreateAttachment func(ctx context.Context, req *attachmentpb.CreateAttachmentRequest) (*attachmentpb.CreateAttachmentResponse, error)
	DeleteAttachment func(ctx context.Context, req *attachmentpb.DeleteAttachmentRequest) (*attachmentpb.DeleteAttachmentResponse, error)
	NewID            func() string

	// Audit history
	ListAuditHistory func(ctx context.Context, req *auditlog.ListAuditRequest) (*auditlog.ListAuditResponse, error)

	// Price lookup for line item (optional — gracefully degrades when nil)
	FindApplicablePriceList func(ctx context.Context, req *pricelistpb.FindApplicablePriceListRequest) (*pricelistpb.FindApplicablePriceListResponse, error)
	ListPriceProducts       func(ctx context.Context, req *priceproductpb.ListPriceProductsRequest) (*priceproductpb.ListPriceProductsResponse, error)

	// Job activity lookup for "from_activities" revenue type (optional — gracefully degrades when nil)
	ReadJobActivity func(ctx context.Context, req *jobactivitypb.ReadJobActivityRequest) (*jobactivitypb.ReadJobActivityResponse, error)

	// RecognizeRevenueFromSubscription delegates auto-population of line items
	// from a subscription's price plan to the espyna use case
	// (skip_header=true mode). When wired, the manual revenue-add flow's
	// autoPopulateLineItems path goes through the use case so the recognize
	// drawer + manual flow share one source of truth.
	RecognizeRevenueFromSubscription func(ctx context.Context, req *revenuepb.CreateRevenueWithLineItemsRequest) (*revenuepb.CreateRevenueWithLineItemsResponse, error)
}

// Module holds all constructed revenue views.
type Module struct {
	routes             centymo.RevenueRoutes
	Dashboard          view.View
	List               view.View
	Table              view.View
	Detail             view.View
	TabAction          view.View
	Add                view.View
	Edit               view.View
	Delete             view.View
	BulkDelete         view.View
	SetStatus          view.View
	BulkSetStatus      view.View
	LineItemTable      view.View
	LineItemAdd        view.View
	LineItemEdit       view.View
	LineItemRemove     view.View
	LineItemDiscount   view.View
	PaymentTable       view.View
	PaymentAdd         view.View
	PaymentEdit        view.View
	PaymentRemove      view.View
	InvoiceDownload       http.HandlerFunc
	SendEmailHandler      http.HandlerFunc
	SearchClients         http.HandlerFunc
	SearchSubscriptions   http.HandlerFunc
	SearchLocations       http.HandlerFunc
	SearchProducts        http.HandlerFunc
	PriceLookup           http.HandlerFunc
	SettingsTemplates     view.View
	SettingsUpload     view.View
	SettingsDelete     view.View
	SettingsSetDefault view.View
	AttachmentUpload   view.View
	AttachmentDelete   view.View
}

func NewModule(deps *ModuleDeps) *Module {
	actionDeps := &revenueaction.Deps{
		Routes:                       deps.Routes,
		Labels:                       deps.Labels,
		DB:                           deps.DB,
		ListPaymentTerms:             deps.ListPaymentTerms,
		ListClients:                  deps.ListClients,
		SearchClientsByName:          deps.SearchClientsByName,
		ListSubscriptions:            deps.ListSubscriptions,
		ReadSubscription:             deps.ReadSubscription,
		ReadPricePlan:                deps.ReadPricePlan,
		ListProductPricePlans:        deps.ListProductPricePlans,
		ReadProduct:                  deps.ReadProduct,
		ListProducts:                 deps.ListProducts,
		CreateRevenue:                deps.CreateRevenue,
		ReadRevenue:                  deps.ReadRevenue,
		UpdateRevenue:                deps.UpdateRevenue,
		DeleteRevenue:                deps.DeleteRevenue,
		CreateRevenueLineItem:        deps.CreateRevenueLineItem,
		ListRevenueLineItems:         deps.ListRevenueLineItems,
		ReadInventoryItem:            deps.ReadInventoryItem,
		UpdateInventoryItem:          deps.UpdateInventoryItem,
		ListInventoryItems:           deps.ListInventoryItems,
		UpdateInventorySerial:        deps.UpdateInventorySerial,
		CreateInventorySerialHistory: deps.CreateInventorySerialHistory,
		FindApplicablePriceList:          deps.FindApplicablePriceList,
		ListPriceProducts:                deps.ListPriceProducts,
		ReadJobActivity:                  deps.ReadJobActivity,
		RecognizeRevenueFromSubscription: deps.RecognizeRevenueFromSubscription,
	}
	paymentDeps := &revenuepayment.Deps{Routes: deps.Routes, DB: deps.DB, Labels: deps.Labels}
	searchDeps := &revenuesearch.Deps{
		DB:                  deps.DB,
		ListClients:         deps.ListClients,
		SearchClientsByName: deps.SearchClientsByName,
		ListSubscriptions:   deps.ListSubscriptions,
		ListProducts:        deps.ListProducts,
	}
	detailDeps := &revenuedetail.DetailViewDeps{
		Routes:               deps.Routes,
		DB:                   deps.DB,
		Labels:               deps.Labels,
		CommonLabels:         deps.CommonLabels,
		TableLabels:          deps.TableLabels,
		ReadRevenue:          deps.ReadRevenue,
		ListRevenueLineItems: deps.ListRevenueLineItems,
		AttachmentOps: attachment.AttachmentOps{
			UploadFile:       deps.UploadFile,
			ListAttachments:  deps.ListAttachments,
			CreateAttachment: deps.CreateAttachment,
			DeleteAttachment: deps.DeleteAttachment,
			NewAttachmentID:  deps.NewID,
		},
		AuditOps: auditlog.AuditOps{
			ListAuditHistory: deps.ListAuditHistory,
		},
	}
	lineItemDeps := &revenuedetail.LineItemDeps{
		Routes:                deps.Routes,
		Labels:                deps.Labels,
		CommonLabels:          deps.CommonLabels,
		TableLabels:           deps.TableLabels,
		ListInventoryItems:    deps.ListInventoryItems,
		SearchProductURL:      deps.Routes.SearchProductURL,
		ReadRevenue:           deps.ReadRevenue,
		UpdateRevenue:         deps.UpdateRevenue,
		CreateRevenueLineItem: deps.CreateRevenueLineItem,
		ReadRevenueLineItem:   deps.ReadRevenueLineItem,
		UpdateRevenueLineItem: deps.UpdateRevenueLineItem,
		DeleteRevenueLineItem: deps.DeleteRevenueLineItem,
		ListRevenueLineItems:  deps.ListRevenueLineItems,
	}

	// Invoice download handler (nil-guarded)
	var invoiceDownload http.HandlerFunc
	if deps.GenerateDoc != nil {
		invoiceDownload = revenueaction.NewInvoiceDownloadHandler(&revenueaction.InvoiceDownloadDeps{
			Routes:               deps.Routes,
			Labels:               deps.Labels,
			ReadRevenue:          deps.ReadRevenue,
			ListRevenueLineItems: deps.ListRevenueLineItems,
			GenerateDoc:          deps.GenerateDoc,
			LoadDefaultTemplate:  deps.LoadDefaultTemplate,
		})
	}

	// Send email handler (nil-guarded)
	var sendEmailHandler http.HandlerFunc
	if deps.GenerateDoc != nil && deps.SendEmail != nil {
		sendEmailHandler = revenueaction.NewSendEmailHandler(&revenueaction.SendEmailDeps{
			Routes:               deps.Routes,
			Labels:               deps.Labels,
			ReadRevenue:          deps.ReadRevenue,
			ListRevenueLineItems: deps.ListRevenueLineItems,
			GenerateDoc:          deps.GenerateDoc,
			LoadDefaultTemplate:  deps.LoadDefaultTemplate,
			SendEmail:            deps.SendEmail,
		})
	}

	// Settings views (nil-guarded)
	var settingsTemplates, settingsUpload, settingsDelete, settingsSetDefault view.View
	if deps.ListDocumentTemplates != nil {
		settingsDeps := &revenuesettings.SettingsViewDeps{
			Routes:                 deps.Routes,
			Labels:                 deps.Labels,
			CommonLabels:           deps.CommonLabels,
			TableLabels:            deps.TableLabels,
			ListDocumentTemplates:  deps.ListDocumentTemplates,
			CreateDocumentTemplate: deps.CreateDocumentTemplate,
			UpdateDocumentTemplate: deps.UpdateDocumentTemplate,
			DeleteDocumentTemplate: deps.DeleteDocumentTemplate,
			UploadTemplate:         deps.UploadTemplate,
		}
		settingsTemplates = revenuesettings.NewView(settingsDeps)
		settingsUpload = revenuesettings.NewUploadAction(settingsDeps)
		settingsDelete = revenuesettings.NewDeleteAction(settingsDeps)
		settingsSetDefault = revenuesettings.NewSetDefaultAction(settingsDeps)
	}

	return &Module{
		routes:    deps.Routes,
		Dashboard: revenuedashboard.NewView(&revenuedashboard.Deps{Labels: deps.Labels, CommonLabels: deps.CommonLabels}),
		List: revenuelist.NewView(&revenuelist.ListViewDeps{
			Routes: deps.Routes, GetListPageData: deps.GetListPageData,
			Labels: deps.Labels, CommonLabels: deps.CommonLabels, TableLabels: deps.TableLabels,
		}),
		Table: revenuelist.NewTableView(&revenuelist.ListViewDeps{
			Routes: deps.Routes, GetListPageData: deps.GetListPageData,
			Labels: deps.Labels, CommonLabels: deps.CommonLabels, TableLabels: deps.TableLabels,
		}),
		Detail:             revenuedetail.NewView(detailDeps),
		TabAction:          revenuedetail.NewTabAction(detailDeps),
		Add:                revenueaction.NewAddAction(actionDeps),
		Edit:               revenueaction.NewEditAction(actionDeps),
		Delete:             revenueaction.NewDeleteAction(actionDeps),
		BulkDelete:         revenueaction.NewBulkDeleteAction(actionDeps),
		SetStatus:          revenueaction.NewSetStatusAction(actionDeps),
		BulkSetStatus:      revenueaction.NewBulkSetStatusAction(actionDeps),
		LineItemTable:      revenuedetail.NewLineItemTableView(lineItemDeps),
		LineItemAdd:        revenuedetail.NewLineItemAddView(lineItemDeps),
		LineItemEdit:       revenuedetail.NewLineItemEditView(lineItemDeps),
		LineItemRemove:     revenuedetail.NewLineItemRemoveView(lineItemDeps),
		LineItemDiscount:   revenuedetail.NewLineItemDiscountView(lineItemDeps),
		PaymentTable:       revenuepayment.NewTableAction(paymentDeps),
		PaymentAdd:         revenuepayment.NewAddAction(paymentDeps),
		PaymentEdit:        revenuepayment.NewEditAction(paymentDeps),
		PaymentRemove:      revenuepayment.NewRemoveAction(paymentDeps),
		InvoiceDownload:     invoiceDownload,
		SendEmailHandler:    sendEmailHandler,
		SearchClients:       revenuesearch.NewSearchClientsAction(searchDeps),
		SearchSubscriptions: revenuesearch.NewSearchSubscriptionsAction(searchDeps),
		SearchLocations:     revenuesearch.NewSearchLocationsAction(searchDeps),
		SearchProducts:      revenuesearch.NewSearchProductsAction(searchDeps),
		PriceLookup:         revenueaction.NewPriceLookupAction(actionDeps),
		SettingsTemplates:   settingsTemplates,
		SettingsUpload:     settingsUpload,
		SettingsDelete:     settingsDelete,
		SettingsSetDefault: settingsSetDefault,
		AttachmentUpload:   revenuedetail.NewAttachmentUploadAction(detailDeps),
		AttachmentDelete:   revenuedetail.NewAttachmentDeleteAction(detailDeps),
	}
}

func (m *Module) RegisterRoutes(r view.RouteRegistrar) {
	r.GET(m.routes.DashboardURL, m.Dashboard)
	r.GET(m.routes.ListURL, m.List)
	r.GET(m.routes.TableURL, m.Table)
	r.POST(m.routes.TableURL, m.Table)
	r.GET(m.routes.DetailURL, m.Detail)
	r.GET(m.routes.TabActionURL, m.TabAction)
	r.GET(m.routes.AddURL, m.Add)
	r.POST(m.routes.AddURL, m.Add)
	r.GET(m.routes.EditURL, m.Edit)
	r.POST(m.routes.EditURL, m.Edit)
	r.POST(m.routes.DeleteURL, m.Delete)
	r.POST(m.routes.BulkDeleteURL, m.BulkDelete)
	r.POST(m.routes.SetStatusURL, m.SetStatus)
	r.POST(m.routes.BulkSetStatusURL, m.BulkSetStatus)
	// Line items
	r.GET(m.routes.LineItemTableURL, m.LineItemTable)
	r.GET(m.routes.LineItemAddURL, m.LineItemAdd)
	r.POST(m.routes.LineItemAddURL, m.LineItemAdd)
	r.GET(m.routes.LineItemEditURL, m.LineItemEdit)
	r.POST(m.routes.LineItemEditURL, m.LineItemEdit)
	r.POST(m.routes.LineItemRemoveURL, m.LineItemRemove)
	r.GET(m.routes.LineItemDiscountURL, m.LineItemDiscount)
	r.POST(m.routes.LineItemDiscountURL, m.LineItemDiscount)
	// Payments
	r.GET(m.routes.PaymentTableURL, m.PaymentTable)
	r.GET(m.routes.PaymentAddURL, m.PaymentAdd)
	r.POST(m.routes.PaymentAddURL, m.PaymentAdd)
	r.GET(m.routes.PaymentEditURL, m.PaymentEdit)
	r.POST(m.routes.PaymentEditURL, m.PaymentEdit)
	r.POST(m.routes.PaymentRemoveURL, m.PaymentRemove)
	// Settings (template management)
	if m.SettingsTemplates != nil {
		r.GET(m.routes.SettingsTemplatesURL, m.SettingsTemplates)
		r.GET(m.routes.SettingsTemplateUploadURL, m.SettingsUpload)
		r.POST(m.routes.SettingsTemplateUploadURL, m.SettingsUpload)
		r.POST(m.routes.SettingsTemplateDeleteURL, m.SettingsDelete)
		r.POST(m.routes.SettingsTemplateDefaultURL, m.SettingsSetDefault)
	}
	// Attachments
	if m.AttachmentUpload != nil {
		r.GET(m.routes.AttachmentUploadURL, m.AttachmentUpload)
		r.POST(m.routes.AttachmentUploadURL, m.AttachmentUpload)
		r.POST(m.routes.AttachmentDeleteURL, m.AttachmentDelete)
	}
	// Note: InvoiceDownload + SendEmailHandler are http.HandlerFunc — register via routes.HandleFunc() in views.go
}
