# centymo-golang

Commerce domain package for Ichizen OS. Owns **seven esqyma proto domains** -- product, subscription, revenue, expenditure, inventory, procurement, and treasury -- making it the largest package in the monorepo (~477 Go files, 27 compose engine units).

**Module path:** `github.com/erniealice/centymo-golang`

> **Centavo convention:** All monetary amounts (revenue, collection, disbursement, expenditure, prices, accruals) are stored as integer centavos. A DB value of `5000000` is `PHP 50,000.00` on screen. Views divide by 100 at render time (`formatCentavos`/`formatCentavoDisplay` helpers). `product.price` is centavos too (confirmed in proto). No exceptions.

## Domain ownership

centymo maps to seven esqyma proto domains (`proto/v1/domain/<d>/`):

| Domain | What it covers |
|--------|----------------|
| `product` | Product catalog, product lines, price lists, resources |
| `subscription` | Plans, price plans, price schedules, product price plans, subscriptions |
| `revenue` | Invoices (revenue), revenue line items, revenue runs, revenue payments |
| `expenditure` | Purchases, expenses, supplier contracts, procurement requests, purchase orders, expense recognition, accrued expenses, supplier billing events |
| `inventory` | Inventory items, serials, transactions, depreciation |
| `procurement` | Supplier plans, cost plans, cost schedules, supplier product plans, supplier subscriptions |
| `treasury` | Collections (inflows), disbursements (outflows), advances dashboard |

No other proto domain lives here. Entity management (clients, suppliers, locations, workspaces) lives in entydad; scheduling (events) in cyta; operations (jobs) in entydad. centymo reads cross-domain entities via typed `UseCases.Entity.*` closures, never by importing those packages.

## Package structure (Option B)

Under Option B the ENTITY is the contract package. Each `domain/<d>/<e>/` directory is one esqyma entity. The domain facade (`domain/<d>/<d>.go`) re-exports entity-local types as Go type aliases so consumers never change their import paths.

```
centymo-golang/
  placement_test.go              # B-STRICT placement gate -- the ONLY test at root
  go.mod / go.sum
  assets.go                      # AssetsFS //go:embed stub (root residual -> pyeza Wave P)
  routes.go  routes_config.go    # root compatibility shim for external consumer (entydad)
  domain/
    product/                     # package product -- facade for the product domain
      product.go                 # facade: type ProductLabels = product.Labels, etc.
      product_module.go          # NewProductModule() assembler
      resource_module.go         # NewResourceModule() assembler
      price_list_module.go       # NewPriceListModule() assembler
      product/                   # entity: product/product (3 sidebar mounts)
        labels.go  routes.go  descriptor.go  embed.go
        list/  detail/  form/  action/  dashboard/  templates/
      line/                      # entity: product/line (2 sidebar mounts)
        labels.go  embed.go  module.go
        list/  detail/  form/  action/  templates/
      price_list/                # entity: product/price_list
        labels.go  routes.go  descriptor.go  embed.go
        list/  detail/  form/  action/  templates/
      resource/                  # entity: product/resource
        labels.go  routes.go  descriptor.go  embed.go
        list/  form/  action/  templates/
    subscription/                # package subscription -- facade
      subscription.go            # facade: aliases for all subscription entities
      price_plan_module.go       # NewPricePlanModule() assembler
      price_schedule_module.go   # NewPriceScheduleModule() assembler
      plan/                      # entity: subscription/plan (2 sidebar mounts)
        labels.go  routes.go  descriptor.go  embed.go
        list/  detail/  form/  action/  templates/
      price_plan/                # entity: subscription/price_plan
        labels.go  routes.go  descriptor.go  embed.go
        list/  detail/  form/  action/  templates/
      price_schedule/            # entity: subscription/price_schedule (2 sidebar mounts)
        labels.go  routes.go  descriptor.go  embed.go
        list/  detail/  form/  action/  templates/
      subscription/              # entity: subscription/subscription
        labels.go  routes.go  descriptor.go  embed.go
        list/  detail/  form/  action/  customize/  recognize/  search/
        spawn_jobs/  spawn_cycle_jobs/  revenue_run/  ad_hoc_actions/  templates/
      product_price_plan/        # entity: subscription/product_price_plan
        labels.go  embed.go  action.go  parent.go  table.go  templates/
      client_packages/           # view projection (legacyAllow -- no esqyma entity)
        labels.go
    revenue/                     # package revenue -- facade
      revenue.go                 # facade: aliases for revenue entities
      revenue_module.go          # NewRevenueModule() assembler
      revenue_run_module.go      # NewRevenueRunModule() assembler
      revenue/                   # entity: revenue/revenue
        labels.go  routes.go  descriptor.go  embed.go
        list/  detail/  form/  action/  dashboard/  payment/  search/  settings/  templates/
      revenue_run/               # entity: revenue/revenue_run
        labels.go  routes.go  descriptor.go  embed.go
        list/  detail/  queue/  shared/  templates/
    expenditure/                 # package expenditure -- facade (14 entity dirs)
      expenditure.go             # facade: aliases for all expenditure entities
      expenditure_module.go  accrued_expense_module.go  expense_recognition_module.go
      expense_recognition_run_module.go  supplier_billing_event_module.go
      supplier_contract_module.go  supplier_contract_price_schedule_module.go
      procurement_request_module.go  purchase_order_module.go
      expenditure/               # entity: expenditure/expenditure
        list/  detail/  form/  action/  category/  pay/  settings/
        purchase_dashboard/  expense_dashboard/  templates/
      accrued_expense/           # entity: expenditure/accrued_expense
      accrued_expense_settlement/  # view sub-flow (legacyAllow)
      expense_recognition/       # entity: expenditure/expense_recognition
      expense_recognition_line/  # entity: expenditure/expense_recognition_line
      expense_recognition_run/   # entity: expenditure/expense_recognition_run
      procurement_request/       # entity: expenditure/procurement_request
      procurement_request_line/  # entity: expenditure/procurement_request_line
      purchase_order/            # entity: expenditure/purchase_order
      supplier_contract/         # entity: expenditure/supplier_contract
      supplier_contract_line/    # entity: expenditure/supplier_contract_line
      supplier_contract_price_schedule/      # entity
      supplier_contract_price_schedule_line/ # entity
      supplier_billing_event/    # entity: expenditure/supplier_billing_event
    inventory/                   # package inventory -- facade
      inventory.go               # facade: aliases for inventory entities
      inventory_module.go        # NewInventoryModule() assembler
      inventory/                 # aggregate view (legacyAllow -- maps to inventory_item/serial/transaction/depreciation)
        list/  detail/  form/  action/  dashboard/  serial/  transaction/
        movements/  depreciation/  templates/
    procurement/                 # package procurement -- facade
      procurement.go             # facade: aliases for procurement entities
      cost_plan_module.go  cost_schedule_module.go  supplier_plan_module.go
      supplier_product_plan_module.go  supplier_subscription_module.go
      procurementdashboard_module.go
      cost_plan/  cost_schedule/  supplier_plan/  supplier_product_plan/
      supplier_product_cost_plan/  supplier_subscription/
      procurementdashboard/      # domain-view: procurement dashboard
        dashboard/  renewals/  utilization/  variance/  recurrence_drafts/  templates/
    treasury/                    # package treasury -- facade
      treasury.go                # facade: aliases for treasury entities
      labels.go  routes.go       # (legacyAllow -- root label/route stubs pending repoint)
      collection_module.go       # NewCollectionModule() assembler
      disbursement_module.go     # NewDisbursementModule() assembler
      collection/                # entity: treasury/collection
        list/  detail/  form/  action/  dashboard/  templates/
      disbursement/              # entity: treasury/disbursement
        list/  detail/  form/  action/  templates/
      treasuryadvancesdashboard/ # domain-view: advances dashboard
        descriptor.go  templates/
      shared/                    # cross-entity shared types for treasury
    shared/                      # charter'd cross-domain leaf (DataSource + Location*)
  block/
    block.go                     # Block() entry point, inline modules
    options.go                   # BlockOption, WithX() funcs, blockConfig
    usecases.go                  # *UseCases typed wiring contract + RequireFor + MustValidate
    wiring.go                    # dashboard wire helpers + location resolver
    catalog.go                   # compose-v2 catalog (27 XxxUnit funcs)
    mustvalidate_test.go         # fail-closed wiring tests
    infra.go                     # Infra shared deps for compose units
    product.go                   # wireProductModules (Product 3-mount + ProductLine 2-mount)
    plan.go                      # wirePlanModules (PricePlan, PriceSchedule, PriceList, Plan)
    subscription.go              # wireSubscriptionModule
    revenue_run.go               # wireRevenueRunModule
    supplier_commitment.go       # wireSupplierCommitmentModules (PO + receipt + returns)
    supplier_contract_price_schedule.go  # wireSupplierContractPriceScheduleModules
    expense_recognition.go       # wireExpenseRecognitionModules
    accrued_expense.go           # wireAccruedExpenseModules
    supplier_subscription.go     # wireSupplierSubscriptionModules (procurement graph)
    supplier_billing_event.go    # wireSupplierBillingEventModule
    advances_dashboard.go        # wireAdvancesDashboardModule
    route_loading_test.go        # route-loading sanity checks
  services/
    checkout/
      serial.go  service.go  types.go  # deferred espyna checkout surface
  tests/                         # Playwright E2E test infrastructure
```

## Dual-mount entities

Several entities appear under multiple sidebar accordion groups, each with namespace-shifted routes so both mounts coexist on the same ServeMux without duplicate registrations. The sidebar mount is controlled by lyngua route overrides + a per-mount `ActiveNav` field.

| Entity | Mounts | Route namespace examples |
|--------|--------|--------------------------|
| Product | **service** (default), **inventory** (`product_inventory`), **supplies** (`product_supplies`) | `/app/products/*`, `/app/inventory/products/*`, `/app/inventory/supplies/*` |
| ProductLine | **service** (default), **inventory** (`product_line_inventory`) | `/app/product-lines/*`, `/app/inventory/product-lines/*` |
| Plan | **service** (default), **bundles** (`plan_bundle`) | `/app/plans/*`, `/app/inventory/bundles/*` |
| PriceSchedule | **service** (default), **inventory** (`price_schedule_inventory`) | `/app/price-schedules/*`, `/app/inventory/price-schedules/*` |

Each mount gets its own `DefaultXxxRoutes()` factory and optional lyngua label overlay (e.g. `product_inventory.json` sparse-overlays on top of `product.json`).

## Compose engine units (27)

Each entity has a `descriptor.go` that returns a `compose.Unit` (routes pointer, labels pointer, lyngua JSON bindings, template FS, navigation contribution). `block/catalog.go` exposes 27 `XxxUnit(uc, infra)` functions that wrap the descriptor with a `Mount` closure mirroring the block.go wire helpers. The compose engine calls these at startup.

| Domain | Units |
|--------|-------|
| product | `product.product`, `product.resource`, `product.price_list` |
| subscription | `subscription.plan`, `subscription.price_plan`, `subscription.price_schedule`, `subscription.subscription` |
| revenue | `revenue.revenue`, `revenue.revenue_run` |
| expenditure | `expenditure.expenditure`, `expenditure.expense_recognition`, `expenditure.expense_recognition_run`, `expenditure.accrued_expense`, `expenditure.procurement_request`, `expenditure.supplier_contract`, `expenditure.supplier_contract_price_schedule`, `expenditure.supplier_billing_event` |
| inventory | `inventory.inventory` |
| procurement | `procurement.cost_plan`, `procurement.cost_schedule`, `procurement.supplier_plan`, `procurement.supplier_product_plan`, `procurement.supplier_subscription`, `procurement.procurementdashboard` |
| treasury | `treasury.collection`, `treasury.disbursement`, `treasury.treasuryadvancesdashboard` |

## Placement gate (`placement_test.go`)

centymo carries a **B-STRICT** placement gate (v2, Option B). `legacyAllow` holds dated residuals pending resolution; the target state is empty (STRICT).

| Rule | What it checks |
|------|----------------|
| **R1** Empty root | No package `.go` files at module root -- only `_test.go` permitted |
| **R2** Canonical dirs | Every first-level dir is an allowed infra surface; every `domain/<d>` is an esqyma proto domain |
| **R2'** Entity dirs | Every `domain/<d>/<child>/` DIR is an esqyma entity of domain `<d>`, `shared`, or a domain-view (name starts with `<d>`) |
| **R3'** Entity contract | No real `*Labels`/`*Routes` type declaration at the domain root -- only alias re-exports (`type X = pkg.Y`) are allowed |
| **R4** No god-files | No `.go` file (excl. `_test.go`) may exceed 1,200 lines |
| **R5** Facade exists | A facade `domain/<d>/<d>.go` must exist for every domain dir with >=1 entity subdir |
| **R6** No cycles | Enforced by `lint-no-domain-cycles.sh` (external, go-list based) |

`crossCutting = false` -- the domain variant applies. esqyma's `proto/v1/domain/` is located at test time so the rules never drift from the live proto tree.

Current `legacyAllow` residuals (all EXPIRES 2026-07-15):

| Key | Reason |
|-----|--------|
| `assets.go` | AssetsFS `//go:embed` stub at root -- move asset hosting to pyeza |
| `datasource.go` | DataSource legacy view-data port at root -- relocate to pyeza |
| `labels.go` | Residual root label aliases -- repoint service-admin import |
| `routes.go` | Residual root URL consts -- repoint service-admin import |
| `routes_config.go` | Root compatibility shim (type aliases) for external consumer entydad |
| `docs` | Planning markdown, not a Go concern |
| `services` | Deferred espyna checkout surface |
| `views` | Root views/ dir holding only a test file after the views/ collapse |
| `domain/shared` | Charter'd cross-domain leaf (DataSource + Location* helpers) |
| `domain/expenditure/accrued_expense_settlement` | Settlement sub-flow; esqyma has no `accrued_expense_settlement` entity |
| `domain/inventory/inventory` | Aggregate view; esqyma inventory entities are `inventory_item`/`_serial`/`_transaction`/`_depreciation` |
| `domain/subscription/client_packages` | Subscription-aggregate projection; no esqyma `client_packages` entity |
| `domain/subscription/price_plan/detail/page.go` | Pre-existing 1723-line detail handler (VIEW-SPLIT wave) |
| `domain/subscription/subscription/detail/page.go` | Pre-existing 1916-line detail handler (VIEW-SPLIT wave) |

## Fail-closed wiring (`block/usecases.go`)

`*UseCases` is the typed wiring contract between service-admin's composition layer and centymo's view modules. It declares ~200 proto-shaped function-field closures grouped by domain concern (Product, Revenue, Expenditure, Inventory, Procurement, Collection, Disbursement, Subscription, Plan, PricePlan, PriceSchedule, etc.).

Two guards make a missing port loud instead of silently rendering empty:

- **`RequireFor(cfg)`** -- lists every missing REQUIRED closure for the enabled modules. A closure is OPTIONAL iff it is never asserted (e.g. dashboards, cross-domain lookups) -- those degrade gracefully to empty-state.
- **`MustValidate(cfg)`** -- the fail-CLOSED posture wrapper:
  - **dev/test** (`testing.Testing()` true or `CENTYMO_BLOCK_STRICT` truthy): PANIC with the full field list -- uncatchable-by-accident, stack-traced, fails CI loudly.
  - **prod**: `log.Printf("FATAL: ...")` at the seam AND returns the error so `Block()` propagates and `NewServiceAdmin` halts boot.

OPTIONAL closures (dashboard callbacks, cross-domain entity reads, document-generation integrations) are never flagged -- they degrade gracefully to empty-state.

## Private services

`services/checkout` is a chartered private helper under `services/` (an allowed first-level directory). It holds stateless checkout serialization logic (`serial.go`, `service.go`, `types.go`). It is not exported as a separate module. Relocation to espyna is deferred.

## Dependencies

- `github.com/erniealice/pyeza-golang` -- UI framework (view system, template engine, compose engine, types)
- `github.com/erniealice/esqyma` -- proto schemas (product, subscription, revenue, expenditure, inventory, procurement, treasury domains)
- `github.com/erniealice/lyngua` -- translation/i18n (label + route overlays per business type)
- `github.com/erniealice/espyna-golang` -- typed use cases (via `reference.Checker` for referential integrity)
- `github.com/erniealice/hybra-golang` -- cross-cutting views (document template view)
- `github.com/erniealice/fycha-golang` -- accounting (indirect, via shared types)
- `google.golang.org/protobuf` -- proto runtime
- `golang.org/x/sync` -- concurrency utilities

## Role in the monorepo

centymo sits in the domain layer above pyeza and espyna. Consumer apps (e.g. `apps/service-admin`) call `block.Block()` to mount the commerce modules, supplying a `*UseCases` via `block.WithUseCases(...)`. The typed contract ensures any drift between espyna and centymo is a compile error, not a silent nil.

centymo is the heaviest domain package by entity count and file count. Its seven-domain span covers the full order-to-cash cycle (product catalog -> plan -> price schedule -> price plan -> subscription -> revenue recognition -> invoice -> collection) and the symmetric procure-to-pay cycle (supplier plan -> cost schedule -> cost plan -> supplier subscription -> expense recognition -> expenditure -> disbursement).

See `docs/wiki/articles/vertical-slices.md` for the full entity trace and `docs/wiki/articles/package-map.md` for the monorepo dependency graph.
