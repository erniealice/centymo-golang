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
	expendituremodmodule "github.com/erniealice/centymo-golang/domain/expenditure/expenditure/module"
	productmodmodule "github.com/erniealice/centymo-golang/domain/product/product/module"
	collectionmodmodule "github.com/erniealice/centymo-golang/domain/treasury/collection/module"
	"time"

	expenseboard "github.com/erniealice/centymo-golang/domain/expenditure/expenditure/expense_dashboard"
	purchaseboard "github.com/erniealice/centymo-golang/domain/expenditure/expenditure/purchase_dashboard"
	productdashboard "github.com/erniealice/centymo-golang/domain/product/product/dashboard"
	collectiondashboard "github.com/erniealice/centymo-golang/domain/treasury/collection/dashboard"
)

// ---------------------------------------------------------------------------
// Cash (collection) dashboard wiring
// ---------------------------------------------------------------------------

// wireCashDashboard sets collectionDeps.GetCashDashboardPageData from
// useCases.Collection.GetCashDashboard if non-nil.
func wireCashDashboard(deps *collectionmodmodule.ModuleDeps, useCases *UseCases) {
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
func wireServiceDashboard(deps *productmodmodule.ModuleDeps, useCases *UseCases) {
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
func wirePurchaseDashboard(deps *expendituremodmodule.ModuleDeps, useCases *UseCases) {
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
func wireExpenseDashboard(deps *expendituremodmodule.ModuleDeps, useCases *UseCases) {
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
