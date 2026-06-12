package block

// wiring.go wires dashboard use cases from centymo's *UseCases into
// module ModuleDeps callbacks.
//
// Option 4.4.A (2026-05-10): the old reflection bridge is replaced with
// explicit typed closures. Each dashboard function reads the centymo view-layer
// callback from *UseCases and assigns it directly — no reflect, no interface{}.
//
// All helpers are nil-safe: if the dashboard callback is nil the ModuleDeps
// field is left unset and the dashboard view renders its empty state.

import (
	"context"
	"time"

	expendituredomain "github.com/erniealice/centymo-golang/domain/expenditure"
	expenseboard "github.com/erniealice/centymo-golang/domain/expenditure/expenditure/expense_dashboard"
	purchaseboard "github.com/erniealice/centymo-golang/domain/expenditure/expenditure/purchase_dashboard"
	productdom "github.com/erniealice/centymo-golang/domain/product"
	productdashboard "github.com/erniealice/centymo-golang/domain/product/product/dashboard"
	shared "github.com/erniealice/centymo-golang/domain/shared"
	treasurydomain "github.com/erniealice/centymo-golang/domain/treasury"
	collectiondashboard "github.com/erniealice/centymo-golang/domain/treasury/collection/dashboard"
	locationpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/location"
)

// ---------------------------------------------------------------------------
// Location name resolver (DB-backed display names for inventory/product views)
// ---------------------------------------------------------------------------

// buildLocationResolver returns a shared.LocationResolver backed by the typed
// espyna location use-case (useCases.Entity.Location.ListLocations). It maps a
// location id to its human name; unknown ids (e.g. the inventory list's human
// slug "ayala-central-bloc", which is not a location row) pass through unchanged
// — preserving the LocationDisplayName stub's behaviour.
//
// Nil-safe: when ListLocations is unwired it returns nil, so every consumer
// falls back to shared.LocationDisplayName via shared.ResolveLocationName.
//
// Each call lists locations and builds a fresh id→name map (matching the
// existing per-call resolveLocationLabel / loadLocations precedents). Batch
// denormalisation onto inventory_item is a deferred follow-up (plan Q5).
func buildLocationResolver(useCases *UseCases) shared.LocationResolver {
	if useCases == nil || useCases.Entity.Location.ListLocations == nil {
		return nil
	}
	list := useCases.Entity.Location.ListLocations
	return func(ctx context.Context, id string) string {
		if id == "" {
			return id
		}
		resp, err := list(ctx, &locationpb.ListLocationsRequest{})
		if err != nil || resp == nil {
			return id
		}
		for _, loc := range resp.GetData() {
			if loc.GetId() == id {
				if name := loc.GetName(); name != "" {
					return name
				}
				return id
			}
		}
		return id
	}
}

// ---------------------------------------------------------------------------
// Cash (collection) dashboard wiring
// ---------------------------------------------------------------------------

// wireCashDashboard sets collectionDeps.GetCashDashboardPageData from
// useCases.Collection.GetCashDashboard if non-nil.
func wireCashDashboard(deps *treasurydomain.CollectionModuleDeps, useCases *UseCases) {
	if useCases == nil || useCases.Collection.GetCashDashboard == nil {
		return
	}
	cb := useCases.Collection.GetCashDashboard
	deps.GetCashDashboardPageData = func(ctx context.Context, req *collectiondashboard.Request) (*collectiondashboard.Response, error) {
		if req == nil {
			req = &collectiondashboard.Request{Now: time.Now()}
		}
		return cb(ctx, req)
	}
}

// ---------------------------------------------------------------------------
// Service (product kind=service) dashboard wiring
// ---------------------------------------------------------------------------

// wireServiceDashboard sets productDeps.GetServiceDashboardPageData from
// useCases.Product.GetServiceDashboard if non-nil.
func wireServiceDashboard(deps *productdom.ProductModuleDeps, useCases *UseCases) {
	if useCases == nil || useCases.Product.GetServiceDashboard == nil {
		return
	}
	cb := useCases.Product.GetServiceDashboard
	deps.GetServiceDashboardPageData = func(ctx context.Context, req *productdashboard.Request) (*productdashboard.Response, error) {
		if req == nil {
			req = &productdashboard.Request{Now: time.Now()}
		}
		return cb(ctx, req)
	}
}

// ---------------------------------------------------------------------------
// Purchase dashboard wiring (kind="purchase")
// ---------------------------------------------------------------------------

// wirePurchaseDashboard sets expDeps.GetPurchaseDashboardPageData from
// useCases.Expenditure.GetPurchaseDashboard if non-nil.
func wirePurchaseDashboard(deps *expendituredomain.ExpenditureModuleDeps, useCases *UseCases) {
	if useCases == nil || useCases.Expenditure.GetPurchaseDashboard == nil {
		return
	}
	cb := useCases.Expenditure.GetPurchaseDashboard
	deps.GetPurchaseDashboardPageData = func(ctx context.Context, req *purchaseboard.Request) (*purchaseboard.Response, error) {
		if req == nil {
			req = &purchaseboard.Request{Now: time.Now()}
		}
		return cb(ctx, req)
	}
}

// ---------------------------------------------------------------------------
// Expense dashboard wiring (kind="expense")
// ---------------------------------------------------------------------------

// wireExpenseDashboard sets expDeps.GetExpenseDashboardPageData from
// useCases.Expenditure.GetExpenseDashboard if non-nil.
func wireExpenseDashboard(deps *expendituredomain.ExpenditureModuleDeps, useCases *UseCases) {
	if useCases == nil || useCases.Expenditure.GetExpenseDashboard == nil {
		return
	}
	cb := useCases.Expenditure.GetExpenseDashboard
	deps.GetExpenseDashboardPageData = func(ctx context.Context, req *expenseboard.Request) (*expenseboard.Response, error) {
		if req == nil {
			req = &expenseboard.Request{Now: time.Now()}
		}
		return cb(ctx, req)
	}
}
