package pricelist

import (
	"context"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	"github.com/erniealice/centymo-golang"
	pricelistaction "github.com/erniealice/centymo-golang/views/pricelist/action"
	pricelistdetail "github.com/erniealice/centymo-golang/views/pricelist/detail"
	pricelistlist "github.com/erniealice/centymo-golang/views/pricelist/list"
	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	pricelistpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/price_list"
	priceproductpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/price_product"
	productpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product"
	"github.com/erniealice/hybra-golang/views/attachment"
)

// ModuleDeps holds all dependencies for the price list module.
type ModuleDeps struct {
	Routes       centymo.PriceListRoutes
	Labels       centymo.PriceListLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels
	// Deletable state
	GetInUseIDs func(ctx context.Context, ids []string) (map[string]bool, error)
	// Price List CRUD
	ListPriceLists  func(ctx context.Context, req *pricelistpb.ListPriceListsRequest) (*pricelistpb.ListPriceListsResponse, error)
	ReadPriceList   func(ctx context.Context, req *pricelistpb.ReadPriceListRequest) (*pricelistpb.ReadPriceListResponse, error)
	CreatePriceList func(ctx context.Context, req *pricelistpb.CreatePriceListRequest) (*pricelistpb.CreatePriceListResponse, error)
	UpdatePriceList func(ctx context.Context, req *pricelistpb.UpdatePriceListRequest) (*pricelistpb.UpdatePriceListResponse, error)
	DeletePriceList func(ctx context.Context, req *pricelistpb.DeletePriceListRequest) (*pricelistpb.DeletePriceListResponse, error)
	// Price Product
	ListPriceProducts  func(ctx context.Context, req *priceproductpb.ListPriceProductsRequest) (*priceproductpb.ListPriceProductsResponse, error)
	CreatePriceProduct func(ctx context.Context, req *priceproductpb.CreatePriceProductRequest) (*priceproductpb.CreatePriceProductResponse, error)
	DeletePriceProduct func(ctx context.Context, req *priceproductpb.DeletePriceProductRequest) (*priceproductpb.DeletePriceProductResponse, error)
	ListProducts       func(ctx context.Context, req *productpb.ListProductsRequest) (*productpb.ListProductsResponse, error)

	// Attachment operations
	UploadFile       func(ctx context.Context, bucket, key string, content []byte, contentType string) error
	ListAttachments  func(ctx context.Context, moduleKey, foreignKey string) (*attachmentpb.ListAttachmentsResponse, error)
	CreateAttachment func(ctx context.Context, req *attachmentpb.CreateAttachmentRequest) (*attachmentpb.CreateAttachmentResponse, error)
	DeleteAttachment func(ctx context.Context, req *attachmentpb.DeleteAttachmentRequest) (*attachmentpb.DeleteAttachmentResponse, error)
	NewID            func() string
}

// Module holds all constructed price list views.
type Module struct {
	routes           centymo.PriceListRoutes
	List             view.View
	TableView        view.View
	Detail           view.View
	TabAction        view.View
	Add              view.View
	Edit             view.View
	Delete           view.View
	BulkDelete       view.View
	PriceProductAdd  view.View
	PriceProductDel  view.View
	AttachmentUpload view.View
	AttachmentDelete view.View
}

func NewModule(deps *ModuleDeps) *Module {
	actionDeps := &pricelistaction.Deps{
		Routes:          deps.Routes,
		Labels:          deps.Labels,
		CreatePriceList: deps.CreatePriceList,
		ReadPriceList:   deps.ReadPriceList,
		UpdatePriceList: deps.UpdatePriceList,
		DeletePriceList: deps.DeletePriceList,
	}
	ppDeps := &pricelistaction.PriceProductDeps{
		Routes:             deps.Routes,
		Labels:             deps.Labels,
		CreatePriceProduct: deps.CreatePriceProduct,
		DeletePriceProduct: deps.DeletePriceProduct,
		ListProducts:       deps.ListProducts,
	}

	detailDeps := &pricelistdetail.DetailViewDeps{
		ReadPriceList:     deps.ReadPriceList,
		ListPriceProducts: deps.ListPriceProducts,
		Labels:            deps.Labels,
		CommonLabels:      deps.CommonLabels,
		TableLabels:       deps.TableLabels,
		AttachmentOps: attachment.AttachmentOps{
			UploadFile:       deps.UploadFile,
			ListAttachments:  deps.ListAttachments,
			CreateAttachment: deps.CreateAttachment,
			DeleteAttachment: deps.DeleteAttachment,
			NewAttachmentID:  deps.NewID,
		},
	}

	listDeps := &pricelistlist.ListViewDeps{
		Routes:         deps.Routes,
		ListPriceLists: deps.ListPriceLists,
		GetInUseIDs:    deps.GetInUseIDs,
		Labels:         deps.Labels,
		CommonLabels:   deps.CommonLabels,
		TableLabels:    deps.TableLabels,
	}

	return &Module{
		routes:           deps.Routes,
		List:             pricelistlist.NewView(listDeps),
		TableView:        pricelistlist.NewTableView(listDeps),
		Detail:           pricelistdetail.NewView(detailDeps),
		TabAction:        pricelistdetail.NewTabAction(detailDeps),
		Add:              pricelistaction.NewAddAction(actionDeps),
		Edit:             pricelistaction.NewEditAction(actionDeps),
		Delete:           pricelistaction.NewDeleteAction(actionDeps),
		BulkDelete:       pricelistaction.NewBulkDeleteAction(actionDeps),
		PriceProductAdd:  pricelistaction.NewPriceProductAddAction(ppDeps),
		PriceProductDel:  pricelistaction.NewPriceProductDeleteAction(ppDeps),
		AttachmentUpload: pricelistdetail.NewAttachmentUploadAction(detailDeps),
		AttachmentDelete: pricelistdetail.NewAttachmentDeleteAction(detailDeps),
	}
}

func (m *Module) RegisterRoutes(r view.RouteRegistrar) {
	r.GET(m.routes.ListURL, m.List)
	r.POST(m.routes.TableURL, m.TableView)
	r.GET(m.routes.DetailURL, m.Detail)
	r.GET(m.routes.TabActionURL, m.TabAction)
	r.GET(m.routes.AddURL, m.Add)
	r.POST(m.routes.AddURL, m.Add)
	r.GET(m.routes.EditURL, m.Edit)
	r.POST(m.routes.EditURL, m.Edit)
	r.POST(m.routes.DeleteURL, m.Delete)
	r.POST(m.routes.BulkDeleteURL, m.BulkDelete)
	r.GET(m.routes.PriceProductAddURL, m.PriceProductAdd)
	r.POST(m.routes.PriceProductAddURL, m.PriceProductAdd)
	r.POST(m.routes.PriceProductDeleteURL, m.PriceProductDel)
	// Attachments
	if m.AttachmentUpload != nil {
		r.GET(m.routes.AttachmentUploadURL, m.AttachmentUpload)
		r.POST(m.routes.AttachmentUploadURL, m.AttachmentUpload)
		r.POST(m.routes.AttachmentDeleteURL, m.AttachmentDelete)
	}
}
