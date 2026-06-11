# centymo-golang

Commerce domain package for Ichizen OS applications. Ships the views,
templates, labels, route config, and per-entity wiring modules for the commerce
verticals: **product**, **inventory**, **revenue**, **subscription**,
**procurement**, **treasury**, and **expenditure**.

**Module path:** `github.com/erniealice/centymo-golang`

> **Centavo amounts:** All finance amounts (revenue, collection, disbursement,
> expenditure, prices) are stored as integer centavos. A DB value of `5000000`
> is `PHP 50,000.00` on screen. Views divide by 100 at render time (per-view
> `formatCentavos`/`formatCentavoDisplay` helpers); `product.price` is centavos
> too (confirmed in proto). No exceptions.

**Dependencies**
- `github.com/erniealice/pyeza-golang` — UI framework (view system, template engine, types, `AppOption`/`AppContext`)
- `github.com/erniealice/esqyma` — proto schemas (the source of truth for domains + entities)
- `github.com/erniealice/lyngua-golang` — translations
- `google.golang.org/protobuf`

---

## Layout — Option B (domain → entity vertical slices)

The package is organized **by domain, then by entity** (a vertical slice per
entity), mirroring esqyma's `proto/v1/domain/<domain>/<entity>/` taxonomy. The
**entity** — not the layer — is the unit of code. This shape is enforced at test
time by `placement_test.go` (the canonical package-cleanup gate) plus
`lint-no-domain-cycles.sh`.

```
centymo-golang/
  go.mod
  placement_test.go            # Option-B structural gate (R1–R6), proto-derived
  assets.go                    # AssetsFS //go:embed  (root residual → pyeza Wave P)
  datasource.go                # DataSource alias = shared.DataSource  (root residual)
  labels.go                    # LocationMap / LocationDisplayName alias = shared.*  (root residual)
  routes.go  routes_config.go  # root compatibility shim for external consumers (entydad)

  domain/<domain>/
    <domain>.go                # FACADE: package <domain>, alias re-exports
                               #   `type ProductLabels = product.Labels` etc.
                               #   (build-enforced completeness; a missing alias
                               #    is a compile error)
    <entity>/                  # one vertical slice per esqyma entity of <domain>
      labels.go  routes.go     # the entity's contract types (Labels, Routes)
      embed.go                 # //go:embed templates/*.html
      <entity>_module.go       # hoisted assembler (package <entity>), OR module/
      list/    page.go         # list view
      detail/  page.go         # detail view (tabs)
      form/    form.go         # create/edit form
      action/  action.go       # CRUD + workflow action handlers
      templates/*.html         # the entity's HTML
    shared/                    # cross-entity shared types for this domain (opt-in)

  block/                       # ASSEMBLER (see below)
  services/checkout/           # service surface (deferred espyna checkout)
  domain/shared/               # charter'd cross-domain leaf (DataSource + Location*)
```

### Domain → entity map (realized tree)

| Domain | Entities / domain-views |
|---|---|
| `product` | `product`, `line`, `price_list`, `resource` |
| `inventory` | `inventory` (aggregate view over esqyma `inventory_item`/`_serial`/`_transaction`/`_depreciation`) |
| `revenue` | `revenue`, `revenue_run` |
| `subscription` | `subscription`, `plan`, `price_plan`, `price_schedule`, `product_price_plan`, `client_packages` |
| `procurement` | `cost_plan`, `cost_schedule`, `supplier_plan`, `supplier_product_plan`, `supplier_product_cost_plan`, `supplier_subscription`, `procurementdashboard` |
| `treasury` | `collection`, `disbursement`, `treasuryadvancesdashboard`, `shared` |
| `expenditure` | `accrued_expense`, `accrued_expense_settlement`, `expenditure`, `expense_recognition`, `expense_recognition_line`, `expense_recognition_run`, `procurement_request`, `procurement_request_line`, `purchase_order`, `supplier_contract`, `supplier_contract_line`, `supplier_contract_price_schedule`, `supplier_contract_price_schedule_line`, `supplier_billing_event` |

Entity directory names match esqyma proto entity names. Two non-entity dir
kinds are allowed under a domain: `shared/` (cross-entity types) and a
**disambiguated domain-level view** whose name starts with the domain name
(e.g. `procurementdashboard`, `treasuryadvancesdashboard`) — the prefix avoids a
bare `package dashboard` collision across domains.

### The domain facade (`domain/<d>/<d>.go`)

Each domain dir with ≥1 entity ships a hand-written `<d>.go` (package `<d>`)
that re-exports the entity contract types via Go **aliases**:

```go
package product

type ProductLabels = product.Labels   // alias re-export — the facade
type LineRoutes    = line.Routes
```

Aliases (not fresh `struct` decls) keep the contract type's identity in its
entity package while giving consumers one import. Completeness is build-enforced:
a missing alias is a compile error, not a silent gap.

---

## `block/` — the assembler / composition entry point

`block.Block()` is the Lego composition entry point. A consumer app (e.g.
`service-admin`) calls `centymo.Block(opts...)` to register the commerce
modules' routes/views onto a pyeza `AppContext`. The wiring is split across
companion files (all `package block`): `block.go` (entry + inline modules),
`options.go` (`BlockOption` / `WithX()` / `blockConfig`), and one
`wireXxxModules` file per module group (`product.go`, `revenue_run.go`,
`subscription.go`, `supplier_commitment.go`, `expense_recognition.go`, …).

### `usecases.go` — typed wiring contract (fail-CLOSED)

`block.UseCases` is the typed struct of proto-shaped function-field ports that
`Block()` needs from outside; `service-admin`'s `buildCentymoUseCases` constructs
it. Two guards make a missing port loud instead of silently rendering empty:

- **`RequireFor(cfg)`** — the required-vs-optional **policy**: for each enabled
  module (`if cfg.wantInventory() { … }`) it asserts each REQUIRED closure is
  non-nil, accumulating a named-field list. A closure is OPTIONAL iff it is never
  asserted (those are documented `// OPTIONAL — NOT in RequireFor`).
- **`MustValidate(cfg)`** — the fail-CLOSED **posture** around `RequireFor`,
  mirroring service-admin's `AUTHZ_ENFORCE` boot-guard. A missing REQUIRED
  closure **PANICS** in dev/test (`testing.Testing()`, or `CENTYMO_BLOCK_STRICT`
  truthy) and, in prod, logs a screaming `FATAL:` line **and** returns the error
  so boot halts — never silently registers an empty feature. `Block()` calls
  `MustValidate` at entry. OPTIONAL nils are never flagged.

---

## Verifying the structure

From the package root:

```bash
go build ./...                          # clean
go vet ./...                            # clean
go test -run Placement ./...            # Option-B placement gate (R1–R5) — PASS
go test ./block/ -run 'MustValidate|RequireFor'   # fail-closed wiring guard — PASS
bash ../../docs/orchestrate/20260610-package-cleanup/lint-no-domain-cycles.sh .   # R6 no intra-domain cycles
```

`placement_test.go` derives the esqyma domain + per-domain entity sets **live**
from `packages/esqyma/proto/v1/domain/` at test time, so the rules can never
drift from proto. Its `legacyAllow` map is the shrinking migration ledger: every
entry carries a dated `EXPIRES` stamp and the capstone target is **empty**
(STRICT). Remaining residuals are root infra stubs awaiting an import repoint,
the non-Go `docs/` dir, the deferred `services/` surface, two pre-existing
oversized detail view handlers (awaiting a view-split wave), and a few
view-only sub-entities with no 1:1 esqyma entity.
