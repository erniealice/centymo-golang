# centymo-golang

Commerce domain package for Ryta OS applications. Provides reusable views, templates, route configurations, labels, and adapters for commerce-related modules (products, inventory, sales, price lists, plans, subscriptions, collections, disbursements, expenditures).

**Module path:** `github.com/erniealice/centymo-golang`

> **Centavo/Cent-Mode Amounts:** All finance-related amounts (revenue, collection, disbursement, expenditure) are stored in the database in centavo/cent mode. For example, a value of `5000000` in the DB represents `PHP 50,000.00` when displayed. Views use `centymo.FormatCentavoAmount()` to divide by 100 and format with comma separators before rendering.

**Dependencies:**
- `github.com/erniealice/pyeza-golang` -- UI framework (view system, template engine, types)
- `github.com/erniealice/esqyma` -- Proto schemas (product, inventory, revenue, treasury, subscription, expenditure)
- `google.golang.org/protobuf`

## Package Structure

```
centymo-golang/
  go.mod
  pkgdir.go                  # runtime.Caller(0)-based package directory resolution
  datasource.go              # DataSource interface (legacy raw DB access)
  routes.go                  # Route constant definitions (URL patterns)
  routes_config.go           # Route config structs + Default*Routes() constructors
  labels.go                  # All label structs + Default*Labels() + MapTableLabels/MapBulkConfig helpers
  htmx.go                    # HTMX response helpers (HTMXSuccess, HTMXError)
  assets.go                  # CopyStyles(), CopyStaticAssets() for CSS/JS asset distribution
  product_adapters.go        # Raw-DB-to-proto adapters for Product CRUD
  pricelist_adapters.go      # Raw-DB-to-proto adapters for PriceList/PriceProduct CRUD
  assets/
    css/
      sales-detail.css
      centymo-sales-dashboard.css
      inventory-detail.css
      centymo-inventory-dashboard.css
      pricelist-detail.css
      product-detail.css
      variant-detail.css
  views/
    product/
      embed.go               # //go:embed templates/*.html
      templates/              # 14 HTML templates
      list/page.go            # Product list view
      action/action.go        # Product CRUD actions (add, edit, delete, bulk delete, status)
      detail/
        page.go               # Product detail view (info, options, variants tabs)
        options.go            # Options tab table builder
        option_page.go        # Option detail page (option values management)
        option_action.go      # Option CRUD actions
        option_value_action.go # Option value CRUD actions
        attributes.go         # Attribute assign/remove views
        helpers.go            # Shared helpers (option count, etc.)
        variant/
          deps.go             # Shared variant deps struct
          page.go             # Variant detail page (info, pricing, stock, images tabs)
          action.go           # Variant CRUD actions (assign, edit, remove)
          images.go           # Variant image upload/delete
          item/
            page.go           # Variant stock detail page (inventory item in variant context)
            serial/
              page.go         # Serial detail page
    inventory/
      embed.go
      templates/              # 10 HTML templates
      list/page.go            # Inventory list view (per-location)
      action/
        action.go             # Inventory item CRUD (add, edit, delete, bulk delete)
        serial.go             # Serial CRUD (assign, edit, remove, table)
        transaction.go        # Transaction actions (assign, table)
        depreciation.go       # Depreciation actions (assign, edit)
        status.go             # Status activate/deactivate (single + bulk)
      detail/
        page.go               # Inventory detail view (info, attributes, serials, transactions, depreciation, audit tabs)
        product_detail.go     # Product-context inventory detail (viewed from product variant)
      dashboard/page.go       # Inventory dashboard (stats, chart, movements, alerts as HTMX partials)
      movements/
        page.go               # Global movement history view
        filter.go             # Filtered movement table (HTMX partial)
        export.go             # CSV export handler (raw http.HandlerFunc, bypasses view layer)
    sales/
      embed.go
      templates/              # 8 HTML templates
      list/page.go            # Sales list view (ongoing, complete, cancelled statuses)
      action/
        action.go             # Sales CRUD + status actions
        payment.go            # Payment CRUD actions (add, edit, remove, table)
      detail/
        page.go               # Sales detail view (info, line items, payments tabs)
        line_items.go         # Line item CRUD views (add, edit, remove, discount, table)
      dashboard/page.go       # Sales dashboard
    pricelist/
      embed.go
      templates/              # 4 HTML templates
      list/page.go            # Price list list view
      action/
        action.go             # Price list CRUD (add, edit, delete, bulk delete)
        price_product.go      # Price product add/delete
      detail/page.go          # Price list detail view (info, products tabs)
    plan/
      embed.go
      templates/list.html
      list/page.go            # Plan list view
    subscription/
      embed.go
      templates/list.html
      list/page.go            # Subscription list view
    collection/
      embed.go
      templates/              # 3 HTML templates (list, detail, drawer-form)
      module.go               # Full module (ModuleDeps, Module, NewModule, RegisterRoutes)
      list/page.go
      detail/page.go
      action/action.go
    disbursement/
      embed.go
      templates/              # 3 HTML templates (list, detail, drawer-form)
      module.go               # Full module (ModuleDeps, Module, NewModule, RegisterRoutes)
      list/page.go
      detail/page.go
      action/action.go
    expenditure/
      embed.go
      templates/list.html
      module.go               # Module with purchase + expense list views
      list/page.go
  services/
    checkout/
      types.go                # CheckoutDeps, CheckoutItem, CheckoutRequest, CheckoutResult, WebhookResult, OrderData
      service.go              # Service: PlaceOrder, HandlePaymentWebhook, GetOrder, reserveStock
      serial.go               # reserveSerials, reserveSerialsForItem (best-effort serial reservation)
```

## Domain Services (`services/`)

Domain services orchestrate multi-step business workflows across multiple aggregates. Unlike views (which render HTML) or use cases (which do single-entity CRUD), services coordinate multiple use cases into a complete business operation.

### Checkout Service (`services/checkout`)

**Import:** `github.com/erniealice/centymo-golang/services/checkout`

Orchestrates the full e-commerce checkout flow: create revenue record → create line items → reserve inventory stock → reserve serials → create payment session.

#### CheckoutDeps

Wire espyna use case `.Execute` functions into this struct:

```go
type CheckoutDeps struct {
    // Revenue
    CreateRevenue func(ctx, *revenuepb.CreateRevenueRequest) (*revenuepb.CreateRevenueResponse, error)
    UpdateRevenue func(ctx, *revenuepb.UpdateRevenueRequest) (*revenuepb.UpdateRevenueResponse, error)
    ReadRevenue   func(ctx, *revenuepb.ReadRevenueRequest) (*revenuepb.ReadRevenueResponse, error)
    ListRevenues  func(ctx, *revenuepb.ListRevenuesRequest) (*revenuepb.ListRevenuesResponse, error)

    // Revenue Line Items
    CreateLineItem func(ctx, *lineItempb.CreateRevenueLineItemRequest) (*lineItempb.CreateRevenueLineItemResponse, error)
    ListLineItems  func(ctx, *lineItempb.ListRevenueLineItemsRequest) (*lineItempb.ListRevenueLineItemsResponse, error)

    // Inventory (stock reservation)
    UpdateInventoryItem func(ctx, *inventoryItempb.UpdateInventoryItemRequest) (*inventoryItempb.UpdateInventoryItemResponse, error)
    ListInventoryItems  func(ctx, *inventoryItempb.ListInventoryItemsRequest) (*inventoryItempb.ListInventoryItemsResponse, error)

    // Inventory Serial (serial assignment -- optional, best-effort)
    ListSerials         func(ctx, *serialpb.ListInventorySerialsRequest) (*serialpb.ListInventorySerialsResponse, error)
    UpdateSerial        func(ctx, *serialpb.UpdateInventorySerialRequest) (*serialpb.UpdateInventorySerialResponse, error)
    CreateSerialHistory func(ctx, *serialHistorypb.CreateInventorySerialHistoryRequest) (*serialHistorypb.CreateInventorySerialHistoryResponse, error)

    // Payment (Maya integration -- optional)
    CreateCheckoutSession func(ctx, *paymentpb.CreateCheckoutSessionRequest) (*paymentpb.CreateCheckoutSessionResponse, error)
    ProcessWebhook        func(ctx, *paymentpb.ProcessWebhookRequest) (*paymentpb.ProcessWebhookResponse, error)
}
```

#### Service Methods

| Method | Description |
|--------|-------------|
| `NewService(deps CheckoutDeps) *Service` | Constructor |
| `PlaceOrder(ctx, CheckoutRequest) (*CheckoutResult, error)` | Full checkout: revenue → line items → stock reservation → serial reservation → payment session |
| `HandlePaymentWebhook(ctx, *paymentpb.ProcessWebhookRequest) (*WebhookResult, error)` | Process payment webhook, update revenue status (paid/cancelled) |
| `GetOrder(ctx, referenceNumber string) (*OrderData, error)` | Retrieve order by reference number with line items |

#### PlaceOrder Flow

1. Generate reference number (`ORD-XXXX-XXXX`)
2. Create Revenue record (status: "pending")
3. Create RevenueLineItem for each cart item
4. Reserve stock (decrement `quantity_available`, increment `quantity_reserved`)
5. Reserve serials (best-effort — update serial status to "reserved", create history entry)
6. Create payment session if `PaymentProvider` is set (Maya redirect URL)
7. Update revenue with checkout session ID

#### Consumer App Wiring

```go
// In container.go
ucs := espynaContainer.GetUseCases()

checkoutDeps := centymoCheckout.CheckoutDeps{
    CreateRevenue:  ucs.Revenue.Revenue.CreateRevenue.Execute,
    UpdateRevenue:  ucs.Revenue.Revenue.UpdateRevenue.Execute,
    ReadRevenue:    ucs.Revenue.Revenue.ReadRevenue.Execute,
    ListRevenues:   ucs.Revenue.Revenue.ListRevenues.Execute,
    CreateLineItem: ucs.Revenue.RevenueLineItem.CreateRevenueLineItem.Execute,
    ListLineItems:  ucs.Revenue.RevenueLineItem.ListRevenueLineItems.Execute,
    // ... inventory + payment deps
}

checkoutSvc := centymoCheckout.NewService(checkoutDeps)
orderService := domain.NewCheckoutOrderService(checkoutSvc, cartService, checkoutDeps.UpdateRevenue)
```

**Consumer apps:** `retail-client`, `service-client` — both wrap the checkout service in an app-level `OrderService` adapter at `internal/domain/order_checkout.go`.

## DataSource Interface (Legacy)

Provides technology-agnostic data access for views. Consumer apps satisfy this by wrapping their database adapter. espyna's `DatabaseAdapter` matches this signature directly.

```go
type DataSource interface {
    ListSimple(ctx context.Context, collection string) ([]map[string]any, error)
    Create(ctx context.Context, collection string, data map[string]any) (map[string]any, error)
    Read(ctx context.Context, collection string, id string) (map[string]any, error)
    Update(ctx context.Context, collection string, id string, data map[string]any) (map[string]any, error)
    Delete(ctx context.Context, collection string, id string) error
    HardDelete(ctx context.Context, collection string, id string) error
}
```

Most views have migrated to typed proto functions. DataSource is retained for backward compatibility where needed (e.g., sales payments, collection methods).

## Route System

### Three-Level Route Configuration

1. **Generic defaults** from Go constants (`routes.go`) -- sensible URL patterns that work out of the box
2. **Industry-specific overrides** via JSON (loaded by consumer apps via lyngua)
3. **App-specific overrides** via Go field assignment after loading defaults

### Route Constants (`routes.go`)

All URL pattern constants are defined as package-level `const`. Key groups:

| Domain | List URL | Detail URL | Action Prefix |
|---|---|---|---|
| Product | `/app/products/list/{status}` | `/app/products/detail/{id}` | `/action/products/` |
| Inventory | `/app/inventory/list/{location}` | `/app/inventory/detail/{id}` | `/action/inventory/` |
| Sales | `/app/sales/list/{status}` | `/app/sales/detail/{id}` | `/action/sales/` |
| Price List | `/app/price-lists/list/{status}` | `/app/price-lists/{id}` | `/action/price-lists/` |
| Plan | `/app/plans/list/{status}` | `/app/plans/{id}` | `/action/plans/` |
| Subscription | `/app/subscriptions/list/{status}` | `/app/subscriptions/{id}` | `/action/subscriptions/` |
| Collection | `/app/collections/list/{status}` | `/app/collections/detail/{id}` | `/action/collections/` |
| Disbursement | `/app/disbursements/list/{status}` | `/app/disbursements/detail/{id}` | `/action/disbursements/` |
| Expenditure | `/app/purchases/list/{status}`, `/app/expenses/list/{status}` | -- | -- |

Product and inventory have extensive nested sub-routes for variants, options, option values, serials, transactions, depreciation, images, and stock detail. Product alone defines 40+ route constants.

### Route Config Structs (`routes_config.go`)

Each domain has a typed route struct with JSON tags for unmarshalling overrides, a `Default*Routes()` constructor, and a `RouteMap()` method returning `map[string]string` for template URL resolution.

| Struct | Field Count | Constructor | RouteMap Keys |
|---|---|---|---|
| `ProductRoutes` | 28 | `DefaultProductRoutes()` | `product.*`, `product.variant.*`, `product.attribute.*`, `product.option.*`, `product.option_value.*` |
| `InventoryRoutes` | 26 | `DefaultInventoryRoutes()` | `inventory.*`, `inventory.serial.*`, `inventory.transaction.*`, `inventory.depreciation.*`, `inventory.dashboard.*` |
| `SalesRoutes` | 17 | `DefaultSalesRoutes()` | `sales.*`, `sales.line_item.*`, `sales.payment.*` |
| `PriceListRoutes` | 9 | `DefaultPriceListRoutes()` | `price_list.*`, `price_list.price_product.*` |
| `ExpenditureRoutes` | 6 | `DefaultExpenditureRoutes()` | `expenditure.purchase.*`, `expenditure.expense.*` |
| `PlanRoutes` | 5 | `DefaultPlanRoutes()` | `plan.*` |
| `SubscriptionRoutes` | 5 | `DefaultSubscriptionRoutes()` | `subscription.*` |
| `CollectionRoutes` | 10 | `DefaultCollectionRoutes()` | `collection.*` |
| `DisbursementRoutes` | 10 | `DefaultDisbursementRoutes()` | `disbursement.*` |

**ProductRoutes** also includes `ActiveNav` and `ActiveSubNav` fields for sidebar navigation context.

Usage example:
```go
routes := centymo.DefaultSalesRoutes()
// Override specific fields for service-admin
routes.ListURL = "/app/bookings/list/{status}"
// Get all routes as a map for templates
routeMap := routes.RouteMap() // {"sales.list": "/app/bookings/list/{status}", ...}
```

## View Sub-Packages

Each view sub-package follows the same pattern:
- **Deps struct** -- holds all dependencies (routes, labels, proto function closures, common labels, table labels)
- **PageData struct** -- embeds `types.PageData`, adds domain-specific fields (Table, ContentTemplate, etc.)
- **NewView()** constructor -- returns a `view.View` (via `view.ViewFunc`)
- Template output uses `view.OK("template-name", pageData)`
- Action views use `centymo.HTMXSuccess(tableID)` and `centymo.HTMXError(message)` for HTMX responses

### product/list

```go
type Deps struct {
    Routes       centymo.ProductRoutes
    ListProducts func(ctx context.Context, req *productpb.ListProductsRequest) (*productpb.ListProductsResponse, error)
    GetInUseIDs  func(ctx context.Context, ids []string) (map[string]bool, error)
    RefreshURL   string
    Labels       centymo.ProductLabels
    CommonLabels pyeza.CommonLabels
    TableLabels  types.TableLabels
}
```
- **Constructor:** `NewView(deps *Deps) view.View`
- **Templates:** `product-list` (full page), `product-list-content` (HTMX partial)
- **Table ID:** `products-table`
- Supports status filtering (`active`/`inactive`), in-use deletion protection, bulk actions, RBAC permission checks

### product/action

```go
type Deps struct {
    Routes           centymo.ProductRoutes
    Labels           centymo.ProductLabels
    CreateProduct    func(ctx context.Context, req *productpb.CreateProductRequest) (*productpb.CreateProductResponse, error)
    ReadProduct      func(ctx context.Context, req *productpb.ReadProductRequest) (*productpb.ReadProductResponse, error)
    UpdateProduct    func(ctx context.Context, req *productpb.UpdateProductRequest) (*productpb.UpdateProductResponse, error)
    DeleteProduct    func(ctx context.Context, req *productpb.DeleteProductRequest) (*productpb.DeleteProductResponse, error)
    SetProductActive func(ctx context.Context, id string, active bool) error
}
```
- **Constructors:** `NewAddAction`, `NewEditAction`, `NewDeleteAction`, `NewBulkDeleteAction`, `NewSetStatusAction`, `NewBulkSetStatusAction`
- **Template:** `product-drawer-form` (GET returns form, POST processes action)
- Add/Edit handle both GET (render form) and POST (process submission)
- `SetProductActive` uses raw map update instead of protobuf because proto3's protojson omits bool fields with value `false`

### product/detail

```go
type Deps struct {
    Routes       centymo.ProductRoutes
    ReadProduct  func(ctx context.Context, req *productpb.ReadProductRequest) (*productpb.ReadProductResponse, error)
    Labels       centymo.ProductLabels
    CommonLabels pyeza.CommonLabels
    TableLabels  types.TableLabels
    DB           centymo.DataSource
    ListProductVariants       func(...)
    ListProductOptions        func(...)
    ListProductOptionValues   func(...)
    ListProductVariantOptions func(...)
}
```
- **Constructors:** `NewView` (full page), `NewTabAction` (HTMX partial for tab switching)
- **Templates:** `product-detail` (full page), `product-tab-info`, `product-tab-options`, `product-tab-variants`
- Tabs: Info (product fields), Options (option CRUD table), Variants (variant CRUD table with option labels)
- **Sub-views** with separate deps: `OptionsDeps` (option CRUD), `AttributeDeps` (attribute assign/remove)

### product/detail/variant

Deeply nested product variant management with its own `Deps` struct shared across:
- `NewPageView` -- Variant detail page with tabs (info, pricing, stock, images)
- `NewTabAction` -- HTMX tab switching
- `NewAssignView`, `NewEditView`, `NewRemoveView` -- Variant CRUD
- `NewTableView` -- Variants table (HTMX partial)
- `NewImageUploadAction`, `NewImageDeleteAction` -- Image management
- `item/NewPageView`, `item/NewTabAction` -- Stock detail (inventory item in variant context)
- `item/serial/NewPageView` -- Serial detail

### inventory/list

```go
type Deps struct {
    Routes             centymo.InventoryRoutes
    ListInventoryItems func(ctx context.Context, req *inventoryitempb.ListInventoryItemsRequest) (*inventoryitempb.ListInventoryItemsResponse, error)
    RefreshURL         string
    Labels             centymo.InventoryLabels
    CommonLabels       pyeza.CommonLabels
    TableLabels        types.TableLabels
}
```
- **Constructor:** `NewView(deps *Deps) view.View`
- **Templates:** `inventory-list` / `inventory-list-content`
- **Table ID:** `inventory-table`
- List is per-location (path param `{location}`, default `ayala-central-bloc`)
- Columns: product name, SKU, item type (serialized/non-serialized/consumable), on hand, available, reorder level, status
- Low stock alert indicator when available quantity is at or below reorder level

### inventory/action

Multiple action files with separate concerns:
- `action.go` -- Inventory item add/edit/delete/bulk-delete (full CRUD)
- `serial.go` -- Serial assign/edit/remove/table
- `transaction.go` -- Transaction assign/table
- `depreciation.go` -- Depreciation assign/edit
- `status.go` -- Single + bulk activate/deactivate

### inventory/dashboard

Dashboard with HTMX lazy-loaded partials:
- `NewView` -- Full dashboard page
- `NewDashboardStatsAction` -- Stats cards (total items, low stock, etc.)
- `NewDashboardChartAction` -- Chart data
- `NewDashboardMovementsAction` -- Recent movements table
- `NewDashboardAlertsAction` -- Stock alerts

### inventory/movements

Global transaction history with SQL-based filtering:
- `NewView` -- Full movements page
- `NewFilterView` -- Filtered table (HTMX partial)
- `NewExportHandler` -- CSV export (`http.HandlerFunc`, bypasses view/template layer)

### sales/list

```go
type Deps struct {
    Routes       centymo.SalesRoutes
    ListRevenues func(ctx context.Context, req *revenuepb.ListRevenuesRequest) (*revenuepb.ListRevenuesResponse, error)
    RefreshURL   string
    Labels       centymo.SalesLabels
    CommonLabels pyeza.CommonLabels
    TableLabels  types.TableLabels
}
```
- **Constructor:** `NewView`
- **Templates:** `sales-list` / `sales-list-content`
- **Table ID:** `sales-table`
- Three statuses: `ongoing`, `complete`, `cancelled`
- Columns: reference, customer, date, amount, status

### sales/detail

Sales detail page with line items and payments tables, each with full CRUD sub-views:
- `NewView` (full page), `NewTabAction` (tab switching)
- `NewLineItemTableView`, `NewLineItemAddView`, `NewLineItemEditView`, `NewLineItemRemoveView`, `NewLineItemDiscountView`

### sales/action

- `Deps` struct includes inventory use cases for stock deduction on status change (complete/cancel)
- `PaymentDeps` -- Separate deps for payment CRUD (uses `centymo.DataSource` for `revenue_payment` collection)

### pricelist/list

```go
type Deps struct {
    Routes         centymo.PriceListRoutes
    ListPriceLists func(ctx context.Context, req *pricelistpb.ListPriceListsRequest) (*pricelistpb.ListPriceListsResponse, error)
    GetInUseIDs    func(ctx context.Context, ids []string) (map[string]bool, error)
    RefreshURL     string
    Labels         centymo.PriceListLabels
    CommonLabels   pyeza.CommonLabels
    TableLabels    types.TableLabels
}
```
- **Table ID:** `price-lists-table`
- Columns: name, date start, date end, status

### plan/list, subscription/list

Simple list views for plan and subscription entities. No detail/action views -- they use inline table actions only.
- **Plan** -- Columns: name, interval (fulfillment type), price (description), status. Table ID: `plans-table`
- **Subscription** -- Columns: customer, plan, start date, status. Table ID: `subscriptions-table`

### collection, disbursement (package-level modules)

Full modules (`module.go`) with list, detail, action views. Treasury domain entities for money-in (collection) and money-out (disbursement). Each provides `ModuleDeps`, `Module`, `NewModule()`, and `RegisterRoutes()`.

### expenditure (package-level module)

Module with purchase and expense list views. Distinguishes type via `ExpenditureType` field (`"purchase"` or `"expense"`) in the shared list view deps. Dashboard views currently reuse the list view.

## Template System

Templates are embedded via `//go:embed templates/*.html` in each view sub-package's `embed.go`. The pyeza template engine discovers them from the embedded FS.

### All Templates

| View Package | Template File | Purpose |
|---|---|---|
| product | `list.html` | Product list (full page + HTMX content partial) |
| product | `detail.html` | Product detail page with tab container |
| product | `product-drawer-form.html` | Product add/edit drawer form |
| product | `variant-drawer-form.html` | Variant assign/edit drawer form |
| product | `variant-detail.html` | Variant detail page with tabs |
| product | `variant-stock-detail.html` | Variant stock (inventory item) detail |
| product | `serial-detail.html` | Serial detail within variant stock |
| product | `option-detail.html` | Option detail page (option values table) |
| product | `option-drawer-form.html` | Option add/edit drawer form |
| product | `option-value-drawer-form.html` | Option value add/edit form |
| product | `attribute-drawer-form.html` | Attribute assign form |
| inventory | `list.html` | Inventory list (full page + content partial) |
| inventory | `detail.html` | Inventory detail page with tabs |
| inventory | `dashboard.html` | Inventory dashboard with lazy-loaded partials |
| inventory | `movements.html` | Movements history page |
| inventory | `inventory-drawer-form.html` | Inventory item add/edit form |
| inventory | `serial-drawer-form.html` | Serial assign/edit form |
| inventory | `transaction-drawer-form.html` | Transaction form |
| inventory | `depreciation-drawer-form.html` | Depreciation form |
| sales | `list.html` | Sales list (full page + content partial) |
| sales | `detail.html` | Sales detail page with tabs |
| sales | `dashboard.html` | Sales dashboard |
| sales | `sales-drawer-form.html` | Sales add/edit form |
| sales | `line-item-drawer-form.html` | Line item add/edit form |
| sales | `line-item-discount-form.html` | Line item discount form |
| sales | `payment-drawer-form.html` | Payment add/edit form |
| pricelist | `list.html` | Price list list (full page + content partial) |
| pricelist | `detail.html` | Price list detail with products tab |
| pricelist | `pricelist-drawer-form.html` | Price list add/edit form |
| pricelist | `price-product-drawer-form.html` | Price product add form |
| plan | `list.html` | Plan list |
| subscription | `list.html` | Subscription list |
| expenditure | `list.html` | Expenditure list (shared for purchase + expense) |
| collection | `list.html` | Collection list |
| collection | `detail.html` | Collection detail |
| collection | `drawer-form.html` | Collection add/edit form |
| disbursement | `list.html` | Disbursement list |
| disbursement | `detail.html` | Disbursement detail |
| disbursement | `drawer-form.html` | Disbursement add/edit form |

## Label Structs (`labels.go`)

All labels are loaded from lyngua JSON translations and unmarshalled into typed structs. The file defines **120+ label structs** organized by domain.

### Top-Level Label Types

| Type | Sub-Struct Count | Default Constructor |
|---|---|---|
| `InventoryLabels` | 19 (Page, Buttons, Columns, Empty, Form, Actions, Bulk, Detail, Tabs, ItemType, Status, Serial, Transaction, Depreciation, Dashboard, Movements, Confirm, Errors, Breadcrumb) | loaded from lyngua JSON |
| `SalesLabels` | 11 (Page, Buttons, Columns, Empty, Form, Actions, Bulk, Detail, Confirm, Errors, Dashboard) | loaded from lyngua JSON |
| `ProductLabels` | 17 (Page, Buttons, Columns, Empty, Form, Actions, Bulk, Tabs, Detail, Status, Variant, Attribute, Confirm, Errors, Breadcrumb, Option, OptionValue) | loaded from lyngua JSON |
| `PriceListLabels` | 10 (Page, Buttons, Columns, Empty, Form, Actions, Bulk, Detail, Confirm, Errors) | loaded from lyngua JSON |
| `ExpenditureLabels` | 10 (Page, Buttons, Columns, Empty, Form, Status, Type, Actions, Bulk, Detail) | loaded from lyngua JSON |
| `CollectionLabels` | 12 (Page, Buttons, Columns, Empty, Form, Actions, Bulk, Detail, Status, Confirm, Errors) | `DefaultCollectionLabels()` |
| `DisbursementLabels` | 12 (Page, Buttons, Columns, Empty, Form, Actions, Bulk, Detail, Status, Confirm, Errors) | `DefaultDisbursementLabels()` |
| `PlanLabels` | 6 (Page, Buttons, Columns, Empty, Actions, Errors) | `DefaultPlanLabels()` |
| `SubscriptionLabels` | 6 (Page, Buttons, Columns, Empty, Actions, Errors) | `DefaultSubscriptionLabels()` |

### Helper Functions

```go
// MapTableLabels maps pyeza CommonLabels to centymo's types.TableLabels.
// Includes search, filters, sort, columns, export, density, pagination, etc.
func MapTableLabels(common pyeza.CommonLabels) types.TableLabels

// MapBulkConfig returns a BulkActionsConfig with labels from common bulk labels.
func MapBulkConfig(common pyeza.CommonLabels) types.BulkActionsConfig

// LocationDisplayName converts a location slug to display name.
// Known slugs: ayala-central-bloc, sm-city-cebu, ayala-center-cebu, robinsons-galleria
func LocationDisplayName(slug string) string
```

## HTMX Helpers (`htmx.go`)

```go
// HTMXSuccess returns a header-only response that signals the sheet to close
// and the table to refresh. Sends HX-Trigger: {"formSuccess":true,"refreshTable":"<tableID>"}
// with HTTP 200 status.
func HTMXSuccess(tableID string) view.ViewResult

// HTMXError returns a header-only response that signals a form error.
// Sends HX-Error-Message header with HTTP 422 status.
func HTMXError(message string) view.ViewResult
```

## Asset Distribution (`assets.go`)

```go
// CopyStyles copies centymo's CSS to {targetDir}/centymo/
func CopyStyles(targetDir string) error

// CopyStaticAssets copies centymo's JS to {targetDir}/centymo/
func CopyStaticAssets(targetDir string) error
```

Uses `runtime.Caller(0)` via `pkgdir.go` to locate the package source directory at runtime, same approach as pyeza-golang. Files are namespaced under a `centymo/` subdirectory in the target.

**CSS assets:** `sales-detail.css`, `centymo-sales-dashboard.css`, `inventory-detail.css`, `centymo-inventory-dashboard.css`, `pricelist-detail.css`, `product-detail.css`, `variant-detail.css`

## DB Adapters

Raw-DB-to-proto adapter functions for use when espyna use cases are unavailable (e.g., `noop` build tag). These convert between `DataSource` `map[string]any` records and proto types.

### Product Adapters (`product_adapters.go`)

| Function | Returns |
|---|---|
| `ProductDBListAdapter(db DataSource)` | `func(ctx, *ListProductsRequest) (*ListProductsResponse, error)` |
| `ProductDBReadAdapter(db DataSource)` | `func(ctx, *ReadProductRequest) (*ReadProductResponse, error)` |
| `ProductDBCreateAdapter(db DataSource)` | `func(ctx, *CreateProductRequest) (*CreateProductResponse, error)` |
| `ProductDBUpdateAdapter(db DataSource)` | `func(ctx, *UpdateProductRequest) (*UpdateProductResponse, error)` |
| `ProductDBDeleteAdapter(db DataSource)` | `func(ctx, *DeleteProductRequest) (*DeleteProductResponse, error)` |

### Price List Adapters (`pricelist_adapters.go`)

| Function | Returns |
|---|---|
| `PriceListDBListAdapter(db DataSource)` | `func(ctx, *ListPriceListsRequest) (*ListPriceListsResponse, error)` |
| `PriceListDBReadAdapter(db DataSource)` | `func(ctx, *ReadPriceListRequest) (*ReadPriceListResponse, error)` |
| `PriceListDBCreateAdapter(db DataSource)` | `func(ctx, *CreatePriceListRequest) (*CreatePriceListResponse, error)` |
| `PriceListDBUpdateAdapter(db DataSource)` | `func(ctx, *UpdatePriceListRequest) (*UpdatePriceListResponse, error)` |
| `PriceListDBDeleteAdapter(db DataSource)` | `func(ctx, *DeletePriceListRequest) (*DeletePriceListResponse, error)` |
| `PriceProductDBListAdapter(db DataSource)` | `func(ctx, *ListPriceProductsRequest) (*ListPriceProductsResponse, error)` |

Internal helpers: `mapToProduct`, `productToMap`, `mapToPriceList`, `priceListToMap`, `mapToPriceProduct`, `mapStr`, `mapFloat64`, `mapInt64`, `mapBool`.

## Consumer App Wiring Guide

### Step 1: Define app-level module

Each consumer app creates a presentation module (e.g., `apps/retail-admin/internal/presentation/product/module.go`) that wraps centymo views:

```go
package product

import (
    "github.com/erniealice/centymo-golang"
    productaction "github.com/erniealice/centymo-golang/views/product/action"
    productdetail "github.com/erniealice/centymo-golang/views/product/detail"
    productlist "github.com/erniealice/centymo-golang/views/product/list"
)

type ModuleDeps struct {
    Routes       centymo.ProductRoutes
    DB           centymo.DataSource
    Labels       centymo.ProductLabels
    CommonLabels pyeza.CommonLabels
    TableLabels  types.TableLabels
    // Proto function closures for all CRUD operations
    ListProducts     func(ctx context.Context, req *productpb.ListProductsRequest) (*productpb.ListProductsResponse, error)
    ReadProduct      func(ctx context.Context, req *productpb.ReadProductRequest) (*productpb.ReadProductResponse, error)
    CreateProduct    func(ctx context.Context, req *productpb.CreateProductRequest) (*productpb.CreateProductResponse, error)
    UpdateProduct    func(ctx context.Context, req *productpb.UpdateProductRequest) (*productpb.UpdateProductResponse, error)
    DeleteProduct    func(ctx context.Context, req *productpb.DeleteProductRequest) (*productpb.DeleteProductResponse, error)
    SetProductActive func(ctx context.Context, id string, active bool) error
    // ... variant, option, attribute, image use cases
}

type Module struct {
    routes centymo.ProductRoutes
    List   view.View
    Detail view.View
    Add    view.View
    Edit   view.View
    // ... all view fields
}

func NewModule(deps *ModuleDeps) *Module {
    actionDeps := &productaction.Deps{
        Routes:           deps.Routes,
        CreateProduct:    deps.CreateProduct,
        ReadProduct:      deps.ReadProduct,
        UpdateProduct:    deps.UpdateProduct,
        DeleteProduct:    deps.DeleteProduct,
        SetProductActive: deps.SetProductActive,
    }

    return &Module{
        routes: deps.Routes,
        List: productlist.NewView(&productlist.Deps{
            Routes:       deps.Routes,
            ListProducts: deps.ListProducts,
            RefreshURL:   deps.Routes.ListURL,
            Labels:       deps.Labels,
            CommonLabels: deps.CommonLabels,
            TableLabels:  deps.TableLabels,
        }),
        Detail: productdetail.NewView(&productdetail.Deps{...}),
        Add:    productaction.NewAddAction(actionDeps),
        Edit:   productaction.NewEditAction(actionDeps),
        // ...
    }
}

func (m *Module) RegisterRoutes(r view.RouteRegistrar) {
    r.GET(m.routes.ListURL, m.List)
    r.GET(m.routes.DetailURL, m.Detail)
    r.GET(m.routes.AddURL, m.Add)
    r.POST(m.routes.AddURL, m.Add)
    r.GET(m.routes.EditURL, m.Edit)
    r.POST(m.routes.EditURL, m.Edit)
    r.POST(m.routes.DeleteURL, m.Delete)
    // ...
}
```

### Step 2: Wire in composition root (`views.go`)

```go
func RegisterAllRoutes(deps *ViewDeps, routes *RouteRegistry) {
    centymoTableLabels := centymo.MapTableLabels(deps.CommonLabels)

    productmod.NewModule(&productmod.ModuleDeps{
        Routes:       centymo.DefaultProductRoutes(),
        DB:           deps.CentymoDB,
        Labels:       deps.ProductLabels,
        CommonLabels: deps.CommonLabels,
        TableLabels:  centymoTableLabels,
        ListProducts: deps.ProductListProducts,
        ReadProduct:  deps.ProductReadProduct,
        // ... all use case functions
    }).RegisterRoutes(routes)
}
```

### Step 3: Load labels in container

Labels are loaded from lyngua JSON at startup:

```go
// In container.go
var productLabels centymo.ProductLabels
lyngua.LoadPath("centymo.product", &productLabels)
```

### Step 4: Copy static assets at startup

```go
centymo.CopyStyles(filepath.Join("assets", "css"))
centymo.CopyStaticAssets(filepath.Join("assets", "js"))
```

### Inline wiring (no app-level module)

For simpler views (plan, subscription, collection list), wire directly in views.go:

```go
planRoutes := centymo.DefaultPlanRoutes()
planLabels := centymo.DefaultPlanLabels()
routes.GET(planRoutes.ListURL, planlist.NewView(&planlist.Deps{
    Routes:       planRoutes,
    ListPlans:    deps.PlanListPlans,
    Labels:       planLabels,
    CommonLabels: deps.CommonLabels,
    TableLabels:  centymoTableLabels,
}))
```

### Using package-level modules (collection, disbursement, expenditure)

These domains provide their own `module.go`:

```go
import collectionmod "github.com/erniealice/centymo-golang/views/collection"

collectionmod.NewModule(&collectionmod.ModuleDeps{
    Routes:           centymo.DefaultCollectionRoutes(),
    Labels:           centymo.DefaultCollectionLabels(),
    CommonLabels:     deps.CommonLabels,
    TableLabels:      centymoTableLabels,
    CreateCollection: deps.CreateCollection,
    ReadCollection:   deps.ReadCollection,
    UpdateCollection: deps.UpdateCollection,
    DeleteCollection: deps.DeleteCollection,
    ListCollections:  deps.ListCollections,
}).RegisterRoutes(routes)
```

## Consumer Apps

### retail-admin

Uses **all** centymo modules via app-level modules in `apps/retail-admin/internal/presentation/`:
- `product/module.go` -- Full product management (list, detail, variants, options, option values, attributes, images, stock, serial)
- `inventory/module.go` -- Full inventory management (list, detail, dashboard, movements, serials, transactions, depreciation, CSV export)
- `sales/module.go` -- Full sales management (list, detail, dashboard, line items, payments, stock deduction)
- `pricelist/module.go` -- Price list management (list, detail, price products)
- Expenditure -- Package-level module (purchase + expense lists)
- Plan, Subscription, Collection -- Inline wired in views.go (list views only)

### service-admin

Uses centymo modules via `apps/service-admin/internal/presentation/`:
- `product/module.go` -- Product (services in service context)
- `inventory/module.go` -- Inventory management
- `sales/module.go` -- Sales (bookings in service context)
- `pricelist/module.go` -- Price lists

Service-admin applies lyngua route/label overrides for service industry terminology (e.g., `products` -> `services`, `Sales` -> `Bookings`).

## Domain Entities Summary

| Entity | Views | Module Level | Data Access | Templates |
|---|---|---|---|---|
| Product | List, Detail, CRUD, Variants, Options, Option Values, Attributes, Images, Stock, Serial | App-level | Proto typed functions | 14 |
| Inventory | List, Detail, CRUD, Dashboard, Movements, Serials, Transactions, Depreciation | App-level | Proto typed functions | 10 |
| Sales | List, Detail, CRUD, Dashboard, Line Items, Payments | App-level | Proto + DataSource (payments) | 8 |
| Price List | List, Detail, CRUD, Price Products | App-level | Proto typed functions | 4 |
| Collection | List, Detail, CRUD, Status | Package-level | Proto typed functions | 3 |
| Disbursement | List, Detail, CRUD, Status | Package-level | Proto typed functions | 3 |
| Expenditure | Purchase List, Expense List | Package-level | Proto typed functions | 1 |
| Plan | List only | Inline | Proto typed functions | 1 |
| Subscription | List only | Inline | Proto typed functions | 1 |
