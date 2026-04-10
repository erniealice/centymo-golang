package action

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"

	"github.com/erniealice/centymo-golang"

	clientpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client"
	inventoryitempb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/inventory_item"
	inventoryserialpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/inventory_serial"
	serialhistorypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/serial_history"
	pricelistpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/price_list"
	priceproductpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/price_product"
	productpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product"
	revenuepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/revenue"
	revenuelineitempb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/revenue_line_item"
	priceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_plan"
	productpriceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/product_price_plan"
	subscriptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription"
)

// PaymentTermOption is a minimal struct for rendering payment term options in the form.
type PaymentTermOption struct {
	Id      string
	Name    string
	NetDays int32
}

// FormInner holds nested form labels accessed via .Labels.Form.* in templates.
type FormInner struct {
	SectionInfo                string
	CurrencyPlaceholder        string
	StatusDraft                string
	StatusComplete             string
	StatusCancelled            string
	CustomerNamePlaceholder    string
	CustomerSearchPlaceholder  string
	CustomerNoResults          string
	LocationPlaceholder        string
	LocationSearchPlaceholder  string
	LocationNoResults          string
}

// FormLabels holds i18n labels for the drawer form template.
type FormLabels struct {
	Customer                  string
	Date                      string
	Currency                  string
	Reference                 string
	ReferencePlaceholder      string
	Status                    string
	Notes                     string
	NotesPlaceholder          string
	Location                  string
	PaymentTerms              string
	SelectPaymentTerm         string
	DueDate                   string
	Subscription              string
	SubscriptionNoResults     string
	RevenueType               string
	RevenueTypeOneTime        string
	RevenueTypeFromEngagement string
	Form                      FormInner
}

// FormData is the template data for the sales drawer form.
type FormData struct {
	FormAction            string
	IsEdit                bool
	ID                    string
	Name                  string
	ClientID              string
	ClientLabel           string
	SearchClientURL       string
	SubscriptionID        string
	SubscriptionLabel     string
	SearchSubscriptionURL string
	ReferenceNumber       string
	Date                  string
	Currency              string
	Status                string
	Notes                 string
	LocationID            string
	LocationLabel         string
	SearchLocationURL     string
	PaymentTerms          []*PaymentTermOption
	SelectedPaymentTermID string
	DueDateString         string
	RevenueType           string
	Labels                FormLabels
	CommonLabels          any
}

// Deps holds dependencies for sales action handlers.
type Deps struct {
	Routes centymo.RevenueRoutes
	Labels centymo.RevenueLabels
	DB     centymo.DataSource // KEEP — used for location, revenue_payment, and collection_method operations

	// Payment terms dropdown (optional — gracefully degrades when nil)
	ListPaymentTerms func(ctx context.Context) ([]*PaymentTermOption, error)

	// Client search for autocomplete (optional — gracefully degrades when nil)
	ListClients         func(ctx context.Context, req *clientpb.ListClientsRequest) (*clientpb.ListClientsResponse, error)
	SearchClientsByName func(ctx context.Context, req *clientpb.SearchClientsByNameRequest) (*clientpb.SearchClientsByNameResponse, error)

	// Subscription search for revenue form autocomplete (optional — gracefully degrades when nil)
	ListSubscriptions func(ctx context.Context, req *subscriptionpb.ListSubscriptionsRequest) (*subscriptionpb.ListSubscriptionsResponse, error)

	// Subscription auto-populate (optional — gracefully degrades when nil)
	ReadSubscription     func(ctx context.Context, req *subscriptionpb.ReadSubscriptionRequest) (*subscriptionpb.ReadSubscriptionResponse, error)
	ReadPricePlan        func(ctx context.Context, req *priceplanpb.ReadPricePlanRequest) (*priceplanpb.ReadPricePlanResponse, error)
	ListProductPricePlans func(ctx context.Context, req *productpriceplanpb.ListProductPricePlansRequest) (*productpriceplanpb.ListProductPricePlansResponse, error)
	ReadProduct          func(ctx context.Context, req *productpb.ReadProductRequest) (*productpb.ReadProductResponse, error)
	ListProducts         func(ctx context.Context, req *productpb.ListProductsRequest) (*productpb.ListProductsResponse, error)

	// Typed revenue operations
	CreateRevenue func(ctx context.Context, req *revenuepb.CreateRevenueRequest) (*revenuepb.CreateRevenueResponse, error)
	ReadRevenue   func(ctx context.Context, req *revenuepb.ReadRevenueRequest) (*revenuepb.ReadRevenueResponse, error)
	UpdateRevenue func(ctx context.Context, req *revenuepb.UpdateRevenueRequest) (*revenuepb.UpdateRevenueResponse, error)
	DeleteRevenue func(ctx context.Context, req *revenuepb.DeleteRevenueRequest) (*revenuepb.DeleteRevenueResponse, error)

	// Typed line item operations
	CreateRevenueLineItem func(ctx context.Context, req *revenuelineitempb.CreateRevenueLineItemRequest) (*revenuelineitempb.CreateRevenueLineItemResponse, error)
	ListRevenueLineItems  func(ctx context.Context, req *revenuelineitempb.ListRevenueLineItemsRequest) (*revenuelineitempb.ListRevenueLineItemsResponse, error)

	// Typed inventory operations
	ReadInventoryItem            func(ctx context.Context, req *inventoryitempb.ReadInventoryItemRequest) (*inventoryitempb.ReadInventoryItemResponse, error)
	UpdateInventoryItem          func(ctx context.Context, req *inventoryitempb.UpdateInventoryItemRequest) (*inventoryitempb.UpdateInventoryItemResponse, error)
	ListInventoryItems           func(ctx context.Context, req *inventoryitempb.ListInventoryItemsRequest) (*inventoryitempb.ListInventoryItemsResponse, error)
	UpdateInventorySerial        func(ctx context.Context, req *inventoryserialpb.UpdateInventorySerialRequest) (*inventoryserialpb.UpdateInventorySerialResponse, error)
	CreateInventorySerialHistory func(ctx context.Context, req *serialhistorypb.CreateInventorySerialHistoryRequest) (*serialhistorypb.CreateInventorySerialHistoryResponse, error)

	// Price lookup for line item (optional — gracefully degrades when nil)
	FindApplicablePriceList func(ctx context.Context, req *pricelistpb.FindApplicablePriceListRequest) (*pricelistpb.FindApplicablePriceListResponse, error)
	ListPriceProducts       func(ctx context.Context, req *priceproductpb.ListPriceProductsRequest) (*priceproductpb.ListPriceProductsResponse, error)
}

func formLabels(t func(string) string) FormLabels {
	return FormLabels{
		Customer:                  t("revenue.form.customer"),
		Date:                      t("revenue.form.date"),
		Currency:                  t("revenue.form.currency"),
		Reference:                 t("revenue.form.reference"),
		ReferencePlaceholder:      t("revenue.form.referencePlaceholder"),
		Status:                    t("revenue.form.status"),
		Notes:                     t("revenue.form.notes"),
		NotesPlaceholder:          t("revenue.form.notesPlaceholder"),
		Location:                  t("revenue.form.location"),
		PaymentTerms:              t("revenue.form.paymentTerms"),
		SelectPaymentTerm:         t("revenue.form.selectPaymentTerm"),
		DueDate:                   t("revenue.form.dueDate"),
		Subscription:              t("revenue.form.subscription"),
		SubscriptionNoResults:     t("revenue.form.subscriptionNoResults"),
		RevenueType:               t("revenue.form.revenueType"),
		RevenueTypeOneTime:        t("revenue.form.revenueTypeOneTime"),
		RevenueTypeFromEngagement: t("revenue.form.revenueTypeFromEngagement"),
		Form: FormInner{
			SectionInfo:               t("revenue.form.sectionInfo"),
			CurrencyPlaceholder:       t("revenue.form.currencyPlaceholder"),
			StatusDraft:               t("revenue.form.statusDraft"),
			StatusComplete:            t("revenue.form.statusComplete"),
			StatusCancelled:           t("revenue.form.statusCancelled"),
			CustomerNamePlaceholder:   t("revenue.form.customerNamePlaceholder"),
			CustomerSearchPlaceholder: t("revenue.form.customerSearchPlaceholder"),
			CustomerNoResults:         t("revenue.form.customerNoResults"),
			LocationPlaceholder:       t("revenue.form.locationPlaceholder"),
			LocationSearchPlaceholder: t("revenue.form.locationSearchPlaceholder"),
			LocationNoResults:         t("revenue.form.locationNoResults"),
		},
	}
}

// loadPaymentTerms fetches payment term options. Returns nil on error (graceful degradation).
func loadPaymentTerms(ctx context.Context, deps *Deps) []*PaymentTermOption {
	if deps.ListPaymentTerms == nil {
		return nil
	}
	terms, err := deps.ListPaymentTerms(ctx)
	if err != nil {
		log.Printf("Failed to load payment terms: %v", err)
		return nil
	}
	return terms
}

// resolveClientLabel finds the display name for a client by ID.
func resolveClientLabel(ctx context.Context, clientID string, listClients func(ctx context.Context, req *clientpb.ListClientsRequest) (*clientpb.ListClientsResponse, error)) string {
	if clientID == "" || listClients == nil {
		return ""
	}
	resp, err := listClients(ctx, &clientpb.ListClientsRequest{})
	if err != nil {
		return clientID
	}
	for _, c := range resp.GetData() {
		if c.GetId() == clientID {
			if cn := c.GetName(); cn != "" {
				return cn
			}
			if u := c.GetUser(); u != nil {
				first := u.GetFirstName()
				last := u.GetLastName()
				if first != "" || last != "" {
					return strings.TrimSpace(first + " " + last)
				}
			}
			return clientID
		}
	}
	return clientID
}

// resolveSubscriptionLabel finds the display name for a subscription by ID.
func resolveSubscriptionLabel(ctx context.Context, subscriptionID string, listSubscriptions func(ctx context.Context, req *subscriptionpb.ListSubscriptionsRequest) (*subscriptionpb.ListSubscriptionsResponse, error)) string {
	if subscriptionID == "" || listSubscriptions == nil {
		return ""
	}
	resp, err := listSubscriptions(ctx, &subscriptionpb.ListSubscriptionsRequest{})
	if err != nil {
		return subscriptionID
	}
	for _, s := range resp.GetData() {
		if s.GetId() == subscriptionID {
			if name := s.GetName(); name != "" {
				return name
			}
			return subscriptionID
		}
	}
	return subscriptionID
}

// resolveLocationLabel finds the display name for a location by ID using the DB.
func resolveLocationLabel(ctx context.Context, locationID string, db centymo.DataSource) string {
	if locationID == "" || db == nil {
		return ""
	}
	records, err := db.ListSimple(ctx, "location")
	if err != nil {
		return locationID
	}
	for _, r := range records {
		id, _ := r["id"].(string)
		if id == locationID {
			name, _ := r["name"].(string)
			if name != "" {
				return name
			}
			return locationID
		}
	}
	return locationID
}

// NewAddAction creates the sales add action (GET = form, POST = create).
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("invoice", "create") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		if viewCtx.Request.Method == http.MethodGet {
			paymentTerms := loadPaymentTerms(ctx, deps)
			return view.OK("revenue-drawer-form", &FormData{
				FormAction:            deps.Routes.AddURL,
				Currency:              "PHP",
				Status:                "draft",
				SearchLocationURL:     deps.Routes.SearchLocationURL,
				PaymentTerms:          paymentTerms,
				SearchClientURL:       deps.Routes.SearchClientURL,
				SearchSubscriptionURL: deps.Routes.SearchSubscriptionURL,
				Labels:                formLabels(viewCtx.T),
				CommonLabels:          nil, // injected by ViewAdapter
			})
		}

		// POST — create sale
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
		}

		r := viewCtx.Request

		// Resolve client: prefer autocomplete client_id, fall back to freetext name
		customerName := r.FormValue("name")
		clientID := r.FormValue("client_id")
		if clientID != "" {
			if label := resolveClientLabel(ctx, clientID, deps.ListClients); label != "" {
				customerName = label
			}
		}

		// Read subscription_id from form (optional)
		subscriptionID := r.FormValue("subscription_id")

		revenueData := &revenuepb.Revenue{
			Name:            customerName,
			ClientId:        clientID,
			ReferenceNumber: strPtr(r.FormValue("reference_number")),
			RevenueDate:     strPtr(r.FormValue("revenue_date_string")),
			Currency:        r.FormValue("currency"),
			Status:          r.FormValue("status"),
			Notes:           strPtr(r.FormValue("notes")),
			LocationId: r.FormValue("location_id"),
			PaymentTermId: func() *string {
				if v := r.FormValue("payment_term_id"); v != "" {
					return &v
				}
				return nil
			}(),
		}
		if subscriptionID != "" {
			revenueData.SubscriptionId = &subscriptionID
		}

		resp, err := deps.CreateRevenue(ctx, &revenuepb.CreateRevenueRequest{
			Data: revenueData,
		})
		if err != nil {
			log.Printf("Failed to create sale: %v", err)
			return centymo.HTMXError(err.Error())
		}

		// Redirect to new sale detail with Items tab
		newID := ""
		if respData := resp.GetData(); len(respData) > 0 {
			newID = respData[0].GetId()
		}

		// Auto-populate line items from subscription's price plan (optional — gracefully degrades)
		if subscriptionID != "" && newID != "" {
			autoPopulateLineItems(r.Context(), deps, newID, subscriptionID)
		}

		if newID != "" {
			return view.ViewResult{
				StatusCode: http.StatusOK,
				Headers: map[string]string{
					"HX-Trigger":  `{"formSuccess":true}`,
					"HX-Redirect": route.ResolveURL(deps.Routes.DetailURL, "id", newID) + "?tab=items",
				},
			}
		}

		return centymo.HTMXSuccess("revenue-table")
	})
}

// autoPopulateLineItems creates revenue line items from a subscription's price plan.
// It fetches the subscription to get the price_plan_id, then looks for associated
// ProductPricePlan records. If itemized plans exist, each becomes its own line item;
// otherwise a single bundle line item is created from the PricePlan amount.
// All errors are logged and silently ignored — this is best-effort enrichment.
func autoPopulateLineItems(ctx context.Context, deps *Deps, revenueID, subscriptionID string) {
	if deps.ReadSubscription == nil || deps.ListRevenueLineItems == nil {
		return
	}

	// 1. Read subscription to get price_plan_id
	subResp, err := deps.ReadSubscription(ctx, &subscriptionpb.ReadSubscriptionRequest{
		Data: &subscriptionpb.Subscription{Id: subscriptionID},
	})
	if err != nil || subResp.GetData() == nil || len(subResp.GetData()) == 0 {
		log.Printf("autoPopulateLineItems: failed to read subscription %s: %v", subscriptionID, err)
		return
	}
	pricePlanID := subResp.GetData()[0].GetPricePlanId()
	if pricePlanID == "" {
		return
	}

	// 2. List ProductPricePlans and filter by price_plan_id in Go
	var items []*productpriceplanpb.ProductPricePlan
	if deps.ListProductPricePlans != nil {
		pppResp, err := deps.ListProductPricePlans(ctx, &productpriceplanpb.ListProductPricePlansRequest{})
		if err == nil {
			for _, ppp := range pppResp.GetData() {
				if ppp.GetPricePlanId() == pricePlanID {
					items = append(items, ppp)
				}
			}
		} else {
			log.Printf("autoPopulateLineItems: failed to list product price plans: %v", err)
		}
	}

	if len(items) > 0 && deps.CreateRevenueLineItem != nil {
		// Itemized mode — one line item per ProductPricePlan
		for _, ppp := range items {
			desc := "Subscription item"
			if deps.ReadProduct != nil {
				prodResp, err := deps.ReadProduct(ctx, &productpb.ReadProductRequest{
					Data: &productpb.Product{Id: ppp.GetProductId()},
				})
				if err == nil && prodResp.GetData() != nil && len(prodResp.GetData()) > 0 {
					if name := prodResp.GetData()[0].GetName(); name != "" {
						desc = name
					}
				}
			}

			price := ppp.GetPrice()
			productID := ppp.GetProductId()
			pppID := ppp.GetId()

			_, err := deps.CreateRevenueLineItem(ctx, &revenuelineitempb.CreateRevenueLineItemRequest{
				Data: &revenuelineitempb.RevenueLineItem{
					RevenueId:          revenueID,
					ProductId:          &productID,
					Description:        desc,
					Quantity:           1,
					UnitPrice:          price,
					TotalPrice:         price,
					LineItemType:       "item",
					ProductPricePlanId: &pppID,
				},
			})
			if err != nil {
				log.Printf("autoPopulateLineItems: failed to create line item for product %s: %v", productID, err)
			}
		}
	} else if deps.ReadPricePlan != nil && deps.CreateRevenueLineItem != nil {
		// Bundle mode — single line item from PricePlan.amount
		ppResp, err := deps.ReadPricePlan(ctx, &priceplanpb.ReadPricePlanRequest{
			Data: &priceplanpb.PricePlan{Id: pricePlanID},
		})
		if err != nil || ppResp.GetData() == nil || len(ppResp.GetData()) == 0 {
			log.Printf("autoPopulateLineItems: failed to read price plan %s: %v", pricePlanID, err)
			return
		}
		pp := ppResp.GetData()[0]
		amount := pp.GetAmount()
		name := pp.GetName()
		if name == "" {
			name = "Subscription"
		}

		_, err = deps.CreateRevenueLineItem(ctx, &revenuelineitempb.CreateRevenueLineItemRequest{
			Data: &revenuelineitempb.RevenueLineItem{
				RevenueId:    revenueID,
				Description:  name,
				Quantity:     1,
				UnitPrice:    amount,
				TotalPrice:   amount,
				LineItemType: "item",
			},
		})
		if err != nil {
			log.Printf("autoPopulateLineItems: failed to create bundle line item: %v", err)
		}
	}
}

// NewEditAction creates the sales edit action (GET = form, POST = update).
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("invoice", "update") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		id := viewCtx.Request.PathValue("id")

		if viewCtx.Request.Method == http.MethodGet {
			readResp, err := deps.ReadRevenue(ctx, &revenuepb.ReadRevenueRequest{
				Data: &revenuepb.Revenue{Id: id},
			})
			if err != nil {
				log.Printf("Failed to read sale %s: %v", id, err)
				return centymo.HTMXError(deps.Labels.Errors.NotFound)
			}
			readData := readResp.GetData()
			if len(readData) == 0 {
				return centymo.HTMXError(deps.Labels.Errors.NotFound)
			}
			record := readData[0]

			paymentTerms := loadPaymentTerms(ctx, deps)
			selectedPaymentTermID := record.GetPaymentTermId()
			existingClientID := record.GetClientId()
			clientLabel := resolveClientLabel(ctx, existingClientID, deps.ListClients)
			existingSubscriptionID := record.GetSubscriptionId()
			subscriptionLabel := resolveSubscriptionLabel(ctx, existingSubscriptionID, deps.ListSubscriptions)
			existingLocationID := record.GetLocationId()
			locationLabel := resolveLocationLabel(ctx, existingLocationID, deps.DB)
			revenueType := "one_time"
			if existingSubscriptionID != "" {
				revenueType = "from_engagement"
			}
			return view.OK("revenue-drawer-form", &FormData{
				FormAction:            route.ResolveURL(deps.Routes.EditURL, "id", id),
				IsEdit:                true,
				ID:                    id,
				Name:                  record.GetName(),
				ClientID:              existingClientID,
				ClientLabel:           clientLabel,
				SearchClientURL:       deps.Routes.SearchClientURL,
				SubscriptionID:        existingSubscriptionID,
				SubscriptionLabel:     subscriptionLabel,
				SearchSubscriptionURL: deps.Routes.SearchSubscriptionURL,
				ReferenceNumber:       record.GetReferenceNumber(),
				Date:                  record.GetRevenueDate(),
				Currency:              record.GetCurrency(),
				Status:                record.GetStatus(),
				Notes:                 record.GetNotes(),
				LocationID:            existingLocationID,
				LocationLabel:         locationLabel,
				SearchLocationURL:     deps.Routes.SearchLocationURL,
				PaymentTerms:          paymentTerms,
				SelectedPaymentTermID: selectedPaymentTermID,
				DueDateString:         record.GetDueDate(),
				RevenueType:           revenueType,
				Labels:                formLabels(viewCtx.T),
				CommonLabels:          nil, // injected by ViewAdapter
			})
		}

		// POST — update sale
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
		}

		r := viewCtx.Request

		// Resolve client: prefer autocomplete client_id, fall back to freetext name
		customerName := r.FormValue("name")
		clientID := r.FormValue("client_id")
		if clientID != "" {
			if label := resolveClientLabel(ctx, clientID, deps.ListClients); label != "" {
				customerName = label
			}
		}

		_, err := deps.UpdateRevenue(ctx, &revenuepb.UpdateRevenueRequest{
			Data: &revenuepb.Revenue{
				Id:              id,
				Name:            customerName,
				ClientId:        clientID,
				ReferenceNumber: strPtr(r.FormValue("reference_number")),
				RevenueDate:     strPtr(r.FormValue("revenue_date_string")),
				Currency:        r.FormValue("currency"),
				Status:          r.FormValue("status"),
				Notes:      strPtr(r.FormValue("notes")),
				LocationId: r.FormValue("location_id"),
				PaymentTermId: func() *string {
					if v := r.FormValue("payment_term_id"); v != "" {
						return &v
					}
					return nil
				}(),
			},
		})
		if err != nil {
			log.Printf("Failed to update sale %s: %v", id, err)
			return centymo.HTMXError(err.Error())
		}

		// Redirect to detail page (preserves current tab)
		return view.ViewResult{
			StatusCode: http.StatusOK,
			Headers: map[string]string{
				"HX-Trigger":  `{"formSuccess":true}`,
				"HX-Redirect": route.ResolveURL(deps.Routes.DetailURL, "id", id),
			},
		}
	})
}

// NewDeleteAction creates the sales delete action (POST only).
// The row ID comes via query param (?id=xxx) appended by table-actions.js.
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("invoice", "delete") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		id := viewCtx.Request.URL.Query().Get("id")
		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
		}
		if id == "" {
			return centymo.HTMXError(deps.Labels.Errors.IDRequired)
		}

		_, err := deps.DeleteRevenue(ctx, &revenuepb.DeleteRevenueRequest{
			Data: &revenuepb.Revenue{Id: id},
		})
		if err != nil {
			log.Printf("Failed to delete sale %s: %v", id, err)
			return centymo.HTMXError(err.Error())
		}

		return centymo.HTMXSuccess("revenue-table")
	})
}

// NewBulkDeleteAction creates the sales bulk delete action (POST only).
// Selected IDs come as multiple "id" form fields from bulk-action.js.
func NewBulkDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("invoice", "delete") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		_ = viewCtx.Request.ParseMultipartForm(32 << 20)

		ids := viewCtx.Request.Form["id"]
		if len(ids) == 0 {
			return centymo.HTMXError(deps.Labels.Errors.NoIDsProvided)
		}

		for _, id := range ids {
			_, err := deps.DeleteRevenue(ctx, &revenuepb.DeleteRevenueRequest{
				Data: &revenuepb.Revenue{Id: id},
			})
			if err != nil {
				log.Printf("Failed to delete sale %s: %v", id, err)
			}
		}

		return centymo.HTMXSuccess("revenue-table")
	})
}

// NewSetStatusAction creates the sales status update action (POST only).
// Expects query params: ?id={saleId}&status={draft|complete|cancelled}
//
// Business rules:
//   - D20: Block completion with zero line items
//   - D21: Block cancellation if payments exist
//   - D5: Deduct stock on completion
//   - D6: Release serials on cancellation
func NewSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("invoice", "update") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		id := viewCtx.Request.URL.Query().Get("id")
		targetStatus := viewCtx.Request.URL.Query().Get("status")

		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
			targetStatus = viewCtx.Request.FormValue("status")
		}
		if id == "" {
			return centymo.HTMXError(deps.Labels.Errors.IDRequired)
		}
		if targetStatus != "draft" && targetStatus != "complete" && targetStatus != "cancelled" {
			return centymo.HTMXError(deps.Labels.Errors.InvalidStatus)
		}

		// D20: Block completion with zero items
		if targetStatus == "complete" {
			lineItems, err := getLineItemsForRevenueTyped(ctx, deps.ListRevenueLineItems, id)
			if err != nil {
				log.Printf("Failed to list line items for sale %s: %v", id, err)
				return centymo.HTMXError(err.Error())
			}
			if len(lineItems) == 0 {
				return centymo.HTMXError(deps.Labels.Errors.NoItemsCannotComplete)
			}

			// Update status
			if _, err := deps.UpdateRevenue(ctx, &revenuepb.UpdateRevenueRequest{
				Data: &revenuepb.Revenue{Id: id, Status: targetStatus},
			}); err != nil {
				log.Printf("Failed to update sale status %s: %v", id, err)
				return centymo.HTMXError(err.Error())
			}

			// D5: Deduct stock on completion
			deductStockForLineItems(ctx, deps, id, lineItems)

			return centymo.HTMXSuccess("revenue-table")
		}

		// D21: Block cancellation if payments exist
		if targetStatus == "cancelled" {
			payments, err := getPaymentsForRevenue(ctx, deps.DB, id)
			if err != nil {
				log.Printf("Failed to list payments for sale %s: %v", id, err)
				return centymo.HTMXError(err.Error())
			}
			if len(payments) > 0 {
				return centymo.HTMXError(deps.Labels.Errors.HasPaymentsCannotCancel)
			}

			// Update status
			if _, err := deps.UpdateRevenue(ctx, &revenuepb.UpdateRevenueRequest{
				Data: &revenuepb.Revenue{Id: id, Status: targetStatus},
			}); err != nil {
				log.Printf("Failed to update sale status %s: %v", id, err)
				return centymo.HTMXError(err.Error())
			}

			// D6: Release serials on cancellation
			lineItems, err := getLineItemsForRevenueTyped(ctx, deps.ListRevenueLineItems, id)
			if err != nil {
				log.Printf("Failed to list line items for serial release on sale %s: %v", id, err)
			} else {
				releaseSerialsForLineItems(ctx, deps, id, lineItems)
			}

			return centymo.HTMXSuccess("revenue-table")
		}

		// Default: draft — just update status
		if _, err := deps.UpdateRevenue(ctx, &revenuepb.UpdateRevenueRequest{
			Data: &revenuepb.Revenue{Id: id, Status: targetStatus},
		}); err != nil {
			log.Printf("Failed to update sale status %s: %v", id, err)
			return centymo.HTMXError(err.Error())
		}

		return centymo.HTMXSuccess("revenue-table")
	})
}

// NewBulkSetStatusAction creates the sales bulk status update action (POST only).
// Selected IDs come as multiple "id" form fields; target status from "target_status" field.
//
// Business rules:
//   - D20: Block bulk completion if any sale has zero line items
//   - D21: Block bulk cancellation if any sale has payments
//   - D5: Deduct stock on completion for each sale
//   - D6: Release serials on cancellation for each sale
func NewBulkSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("invoice", "update") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		_ = viewCtx.Request.ParseMultipartForm(32 << 20)

		ids := viewCtx.Request.Form["id"]
		targetStatus := viewCtx.Request.FormValue("target_status")

		if len(ids) == 0 {
			return centymo.HTMXError(deps.Labels.Errors.NoIDsProvided)
		}
		if targetStatus != "draft" && targetStatus != "complete" && targetStatus != "cancelled" {
			return centymo.HTMXError(deps.Labels.Errors.InvalidTargetStatus)
		}

		// D21: Block bulk cancellation if any sale has payments
		if targetStatus == "cancelled" {
			withPayments := 0
			for _, id := range ids {
				payments, err := getPaymentsForRevenue(ctx, deps.DB, id)
				if err != nil {
					log.Printf("Failed to check payments for sale %s: %v", id, err)
					continue
				}
				if len(payments) > 0 {
					withPayments++
				}
			}
			if withPayments > 0 {
				return centymo.HTMXError(fmt.Sprintf(
					deps.Labels.Errors.BulkHasPayments,
					withPayments, len(ids),
				))
			}
		}

		// D20: Block bulk completion if any sale has zero line items
		if targetStatus == "complete" {
			emptyCount := 0
			for _, id := range ids {
				lineItems, err := getLineItemsForRevenueTyped(ctx, deps.ListRevenueLineItems, id)
				if err != nil {
					log.Printf("Failed to check line items for sale %s: %v", id, err)
					continue
				}
				if len(lineItems) == 0 {
					emptyCount++
				}
			}
			if emptyCount > 0 {
				return centymo.HTMXError(fmt.Sprintf(
					deps.Labels.Errors.BulkNoItems,
					emptyCount, len(ids),
				))
			}
		}

		// Update all statuses and apply side-effects
		for _, id := range ids {
			if _, err := deps.UpdateRevenue(ctx, &revenuepb.UpdateRevenueRequest{
				Data: &revenuepb.Revenue{Id: id, Status: targetStatus},
			}); err != nil {
				log.Printf("Failed to update sale status %s: %v", id, err)
				continue
			}

			// D5: Deduct stock on completion
			if targetStatus == "complete" {
				lineItems, err := getLineItemsForRevenueTyped(ctx, deps.ListRevenueLineItems, id)
				if err != nil {
					log.Printf("Failed to list line items for stock deduction on sale %s: %v", id, err)
					continue
				}
				deductStockForLineItems(ctx, deps, id, lineItems)
			}

			// D6: Release serials on cancellation
			if targetStatus == "cancelled" {
				lineItems, err := getLineItemsForRevenueTyped(ctx, deps.ListRevenueLineItems, id)
				if err != nil {
					log.Printf("Failed to list line items for serial release on sale %s: %v", id, err)
					continue
				}
				releaseSerialsForLineItems(ctx, deps, id, lineItems)
			}
		}

		return centymo.HTMXSuccess("revenue-table")
	})
}

// ---------------------------------------------------------------------------
// Helpers for status change business rules
// ---------------------------------------------------------------------------

// getLineItemsForRevenueTyped returns all revenue_line_item records for a given revenue ID using typed use case.
func getLineItemsForRevenueTyped(
	ctx context.Context,
	listFn func(ctx context.Context, req *revenuelineitempb.ListRevenueLineItemsRequest) (*revenuelineitempb.ListRevenueLineItemsResponse, error),
	revenueID string,
) ([]map[string]any, error) {
	resp, err := listFn(ctx, &revenuelineitempb.ListRevenueLineItemsRequest{
		RevenueId: &revenueID,
	})
	if err != nil {
		return nil, err
	}
	var items []map[string]any
	for _, item := range resp.GetData() {
		if item.GetRevenueId() == revenueID {
			items = append(items, map[string]any{
				"id":                  item.GetId(),
				"revenue_id":          item.GetRevenueId(),
				"description":         item.GetDescription(),
				"quantity":            fmt.Sprintf("%.0f", item.GetQuantity()),
				"unit_price":          fmt.Sprintf("%.2f", item.GetUnitPrice()),
				"cost_price":          fmt.Sprintf("%.2f", item.GetCostPrice()),
				"total":               fmt.Sprintf("%.2f", item.GetTotalPrice()),
				"line_item_type":      item.GetLineItemType(),
				"inventory_item_id":   item.GetInventoryItemId(),
				"inventory_serial_id": item.GetInventorySerialId(),
				"notes":               item.GetNotes(),
			})
		}
	}
	return items, nil
}

// strPtr returns a pointer to a string.
func strPtr(s string) *string {
	return &s
}

// getPaymentsForRevenue returns all revenue_payment records for a given revenue ID.
func getPaymentsForRevenue(ctx context.Context, db centymo.DataSource, revenueID string) ([]map[string]any, error) {
	all, err := db.ListSimple(ctx, "revenue_payment")
	if err != nil {
		return nil, err
	}
	var payments []map[string]any
	for _, r := range all {
		rid, _ := r["revenue_id"].(string)
		if rid == revenueID {
			payments = append(payments, r)
		}
	}
	return payments, nil
}

// deductStockForLineItems decrements inventory quantities and marks serials as sold.
func deductStockForLineItems(ctx context.Context, deps *Deps, saleID string, lineItems []map[string]any) {
	for _, item := range lineItems {
		inventoryItemID, _ := item["inventory_item_id"].(string)
		serialID, _ := item["inventory_serial_id"].(string)

		// Deduct quantity from inventory item
		if inventoryItemID != "" {
			resp, err := deps.ReadInventoryItem(ctx, &inventoryitempb.ReadInventoryItemRequest{
				Data: &inventoryitempb.InventoryItem{Id: inventoryItemID},
			})
			if err != nil {
				log.Printf("Failed to read inventory item %s for stock deduction: %v", inventoryItemID, err)
				continue
			}
			data := resp.GetData()
			if len(data) == 0 {
				log.Printf("Inventory item %s not found for stock deduction", inventoryItemID)
				continue
			}
			invItem := data[0]

			lineQtyStr, _ := item["quantity"].(string)
			lineQty, _ := strconv.ParseFloat(lineQtyStr, 64)

			newQty := invItem.GetQuantityOnHand() - lineQty
			if _, err := deps.UpdateInventoryItem(ctx, &inventoryitempb.UpdateInventoryItemRequest{
				Data: &inventoryitempb.InventoryItem{
					Id:             inventoryItemID,
					QuantityOnHand: newQty,
				},
			}); err != nil {
				log.Printf("Failed to deduct stock for inventory item %s: %v", inventoryItemID, err)
			}
		}

		// Mark serial as sold and create history
		if serialID != "" {
			if _, err := deps.UpdateInventorySerial(ctx, &inventoryserialpb.UpdateInventorySerialRequest{
				Data: &inventoryserialpb.InventorySerial{
					Id:     serialID,
					Status: "sold",
				},
			}); err != nil {
				log.Printf("Failed to mark serial %s as sold: %v", serialID, err)
			}

			if _, err := deps.CreateInventorySerialHistory(ctx, &serialhistorypb.CreateInventorySerialHistoryRequest{
				Data: &serialhistorypb.InventorySerialHistory{
					InventorySerialId: serialID,
					InventoryItemId:   inventoryItemID,
					FromStatus:        "reserved",
					ToStatus:          "sold",
					ReferenceType:     "revenue",
					ReferenceId:       saleID,
					Notes:             "Auto: sale completed",
				},
			}); err != nil {
				log.Printf("Failed to create serial history for %s: %v", serialID, err)
			}
		}
	}
}

// releaseSerialsForLineItems marks serials as available and creates history records.
func releaseSerialsForLineItems(ctx context.Context, deps *Deps, saleID string, lineItems []map[string]any) {
	for _, item := range lineItems {
		serialID, _ := item["inventory_serial_id"].(string)
		if serialID == "" {
			continue
		}

		inventoryItemID, _ := item["inventory_item_id"].(string)

		if _, err := deps.UpdateInventorySerial(ctx, &inventoryserialpb.UpdateInventorySerialRequest{
			Data: &inventoryserialpb.InventorySerial{
				Id:     serialID,
				Status: "available",
			},
		}); err != nil {
			log.Printf("Failed to release serial %s: %v", serialID, err)
		}

		if _, err := deps.CreateInventorySerialHistory(ctx, &serialhistorypb.CreateInventorySerialHistoryRequest{
			Data: &serialhistorypb.InventorySerialHistory{
				InventorySerialId: serialID,
				InventoryItemId:   inventoryItemID,
				FromStatus:        "available",
				ToStatus:          "available",
				ReferenceType:     "revenue",
				ReferenceId:       saleID,
				Notes:             "Auto: sale cancelled",
			},
		}); err != nil {
			log.Printf("Failed to create serial history for %s: %v", serialID, err)
		}
	}
}
