package block

// wiring.go wires dashboard use cases from the espyna UseCases aggregate
// into centymo module ModuleDeps callbacks.
//
// Since the dashboard use-case request/response types live in espyna's
// internal packages (unreachable from centymo), we use reflection to:
//  1. Dereference the use-case pointer field (e.g. useCases.Treasury.CashDashboard)
//  2. Build the Execute request via reflect.New + field-name assignment
//  3. Call Execute reflectively
//  4. Copy matching fields from the response to the view-layer Response type
//
// All helpers are nil-safe: if the Dashboard field is nil the callback is
// left unset and the dashboard view renders empty state (its existing behaviour).

import (
	"context"
	"reflect"
	"time"

	consumer "github.com/erniealice/espyna-golang/consumer"

	collectiondashboard "github.com/erniealice/centymo-golang/views/collection/dashboard"
	expenseboard "github.com/erniealice/centymo-golang/views/expenditure/expense_dashboard"
	purchaseboard "github.com/erniealice/centymo-golang/views/expenditure/purchase_dashboard"
	productdashboard "github.com/erniealice/centymo-golang/views/product/dashboard"

	collectionmod "github.com/erniealice/centymo-golang/views/collection"
	expendituremod "github.com/erniealice/centymo-golang/views/expenditure"
	productmod "github.com/erniealice/centymo-golang/views/product"

	collectionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/collection"
	expenditurepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/expenditure"
	productpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product"
)

// callDashboardExecute calls a use-case's Execute method via reflection.
// The useCase value must be a non-nil pointer to a use-case struct.
// workspaceID and now are set on the request struct by field name.
// An optional extraFields map allows setting additional request fields (e.g. Kind).
// Returns the dereferenced response (as reflect.Value) and an error.
func callDashboardExecute(
	useCase reflect.Value,
	ctx context.Context,
	workspaceID string,
	now time.Time,
	extraFields map[string]any,
) (reflect.Value, error) {
	m := useCase.MethodByName("Execute")
	if !m.IsValid() {
		return reflect.Value{}, nil
	}
	reqType := m.Type().In(1).Elem() // *Request → Request
	reqPtr := reflect.New(reqType)
	if f := reqPtr.Elem().FieldByName("WorkspaceID"); f.IsValid() && f.CanSet() {
		f.SetString(workspaceID)
	}
	if f := reqPtr.Elem().FieldByName("Now"); f.IsValid() && f.CanSet() {
		f.Set(reflect.ValueOf(now))
	}
	for field, val := range extraFields {
		if f := reqPtr.Elem().FieldByName(field); f.IsValid() && f.CanSet() {
			f.Set(reflect.ValueOf(val))
		}
	}
	results := m.Call([]reflect.Value{reflect.ValueOf(ctx), reqPtr})
	if len(results) < 2 {
		return reflect.Value{}, nil
	}
	if !results[1].IsNil() {
		return reflect.Value{}, results[1].Interface().(error)
	}
	resp := results[0]
	if resp.Kind() == reflect.Ptr && !resp.IsNil() {
		return resp.Elem(), nil
	}
	return resp, nil
}

// int64FieldC reads an int64 field by name from a reflect.Value (struct).
func int64FieldC(v reflect.Value, name string) int64 {
	if !v.IsValid() {
		return 0
	}
	f := v.FieldByName(name)
	if !f.IsValid() {
		return 0
	}
	return f.Int()
}

// stringFieldC reads a string field by name from a reflect.Value (struct).
func stringFieldC(v reflect.Value, name string) string {
	if !v.IsValid() {
		return ""
	}
	f := v.FieldByName(name)
	if !f.IsValid() {
		return ""
	}
	return f.String()
}

// float64SliceFieldC reads a []float64 field by name.
func float64SliceFieldC(v reflect.Value, name string) []float64 {
	if !v.IsValid() {
		return nil
	}
	f := v.FieldByName(name)
	if !f.IsValid() || f.IsNil() {
		return nil
	}
	if s, ok := f.Interface().([]float64); ok {
		return s
	}
	return nil
}

// stringSliceFieldC reads a []string field by name.
func stringSliceFieldC(v reflect.Value, name string) []string {
	if !v.IsValid() {
		return nil
	}
	f := v.FieldByName(name)
	if !f.IsValid() || f.IsNil() {
		return nil
	}
	if s, ok := f.Interface().([]string); ok {
		return s
	}
	return nil
}

// ---------------------------------------------------------------------------
// Cash (collection) dashboard wiring
// ---------------------------------------------------------------------------

// wireCashDashboard sets collectionDeps.GetCashDashboardPageData if
// useCases.Treasury.CashDashboard is non-nil.
func wireCashDashboard(deps *collectionmod.ModuleDeps, useCases *consumer.UseCases) {
	if useCases == nil || useCases.Treasury == nil || useCases.Treasury.CashDashboard == nil {
		return
	}
	uc := reflect.ValueOf(useCases.Treasury.CashDashboard)
	deps.GetCashDashboardPageData = func(ctx context.Context, req *collectiondashboard.Request) (*collectiondashboard.Response, error) {
		now := time.Now()
		if req != nil && !req.Now.IsZero() {
			now = req.Now
		}
		resp, err := callDashboardExecute(uc, ctx, "", now, nil)
		if err != nil || !resp.IsValid() {
			return nil, err
		}
		// Stats sub-struct
		var stats collectiondashboard.Stats
		if s := resp.FieldByName("Stats"); s.IsValid() {
			stats.Pending = s.FieldByName("Pending").Int()
			stats.Overdue = s.FieldByName("Overdue").Int()
			stats.CollectedToday = s.FieldByName("CollectedToday").Int()
			stats.CollectedThisWeek = s.FieldByName("CollectedThisWeek").Int()
		}
		// Recent: []proto — type is the same across packages
		var recent []*collectionpb.Collection
		if f := resp.FieldByName("Recent"); f.IsValid() && !f.IsNil() {
			if v, ok := f.Interface().([]*collectionpb.Collection); ok {
				recent = v
			}
		}
		return &collectiondashboard.Response{
			Stats:       stats,
			DailyLabels: stringSliceFieldC(resp, "DailyLabels"),
			DailyValues: float64SliceFieldC(resp, "DailyValues"),
			ModeLabels:  stringSliceFieldC(resp, "ModeLabels"),
			ModeValues:  float64SliceFieldC(resp, "ModeValues"),
			Recent:      recent,
		}, nil
	}
}

// ---------------------------------------------------------------------------
// Service (product kind=service) dashboard wiring
// ---------------------------------------------------------------------------

// wireServiceDashboard sets productDeps.GetServiceDashboardPageData if
// useCases.Product.Dashboard is non-nil.
func wireServiceDashboard(deps *productmod.ModuleDeps, useCases *consumer.UseCases) {
	if useCases == nil || useCases.Product == nil || useCases.Product.Dashboard == nil {
		return
	}
	uc := reflect.ValueOf(useCases.Product.Dashboard)
	deps.GetServiceDashboardPageData = func(ctx context.Context, req *productdashboard.Request) (*productdashboard.Response, error) {
		now := time.Now()
		if req != nil && !req.Now.IsZero() {
			now = req.Now
		}
		resp, err := callDashboardExecute(uc, ctx, "", now, nil)
		if err != nil || !resp.IsValid() {
			return nil, err
		}
		// Stats sub-struct
		var stats productdashboard.Stats
		if s := resp.FieldByName("Stats"); s.IsValid() {
			stats.TotalActive = s.FieldByName("TotalActive").Int()
			stats.TopRevenueName = s.FieldByName("TopRevenueName").String()
			stats.TopRevenueValue = s.FieldByName("TopRevenueValue").Int()
			stats.LineCount = s.FieldByName("LineCount").Int()
			stats.RecentlyAddedCnt = s.FieldByName("RecentlyAddedCnt").Int()
		}
		// TopRevenue: []TopRevenueRow (same field shape, different package type)
		var topRevenue []productdashboard.TopRevenueRow
		if f := resp.FieldByName("TopRevenue"); f.IsValid() && !f.IsNil() {
			for i := 0; i < f.Len(); i++ {
				s := f.Index(i)
				topRevenue = append(topRevenue, productdashboard.TopRevenueRow{
					ProductID:   s.FieldByName("ProductID").String(),
					ProductName: s.FieldByName("ProductName").String(),
					Total:       s.FieldByName("Total").Int(),
				})
			}
		}
		// Recent: []*productpb.Product — same proto type
		var recent []*productpb.Product
		if f := resp.FieldByName("Recent"); f.IsValid() && !f.IsNil() {
			if v, ok := f.Interface().([]*productpb.Product); ok {
				recent = v
			}
		}
		return &productdashboard.Response{
			Stats:      stats,
			LineLabels: stringSliceFieldC(resp, "LineLabels"),
			LineValues: float64SliceFieldC(resp, "LineValues"),
			TopRevenue: topRevenue,
			Recent:     recent,
		}, nil
	}
}

// ---------------------------------------------------------------------------
// Purchase dashboard wiring (kind="purchase")
// ---------------------------------------------------------------------------

// wirePurchaseDashboard sets expDeps.GetPurchaseDashboardPageData if
// useCases.Expenditure.Dashboard is non-nil.
func wirePurchaseDashboard(deps *expendituremod.ModuleDeps, useCases *consumer.UseCases) {
	if useCases == nil || useCases.Expenditure == nil || useCases.Expenditure.Dashboard == nil {
		return
	}
	uc := reflect.ValueOf(useCases.Expenditure.Dashboard)
	deps.GetPurchaseDashboardPageData = func(ctx context.Context, req *purchaseboard.Request) (*purchaseboard.Response, error) {
		now := time.Now()
		if req != nil && !req.Now.IsZero() {
			now = req.Now
		}
		extra := map[string]any{"Kind": "purchase"}
		resp, err := callDashboardExecute(uc, ctx, "", now, extra)
		if err != nil || !resp.IsValid() {
			return nil, err
		}
		// Stats sub-struct
		var stats purchaseboard.Stats
		if s := resp.FieldByName("Stats"); s.IsValid() {
			stats.OpenCount = s.FieldByName("OpenCount").Int()
			stats.AwaitingCount = s.FieldByName("AwaitingCount").Int()
			stats.SpentMTD = s.FieldByName("TotalMTD").Int()
			stats.TopSupplierName = s.FieldByName("TopSupplierName").String()
			stats.TopSupplierTotal = s.FieldByName("TopSupplierTotal").Int()
		}
		// TopSuppliers: []TopSupplierRow (same field shape, different package type)
		var topSuppliers []purchaseboard.TopSupplierRow
		if f := resp.FieldByName("TopSuppliers"); f.IsValid() && !f.IsNil() {
			for i := 0; i < f.Len(); i++ {
				s := f.Index(i)
				topSuppliers = append(topSuppliers, purchaseboard.TopSupplierRow{
					SupplierID:   s.FieldByName("SupplierID").String(),
					SupplierName: s.FieldByName("SupplierName").String(),
					Total:        s.FieldByName("Total").Int(),
				})
			}
		}
		// Recent: []*expenditurepb.Expenditure — same proto type
		var recent []*expenditurepb.Expenditure
		if f := resp.FieldByName("Recent"); f.IsValid() && !f.IsNil() {
			if v, ok := f.Interface().([]*expenditurepb.Expenditure); ok {
				recent = v
			}
		}
		return &purchaseboard.Response{
			Stats:        stats,
			MonthLabels:  stringSliceFieldC(resp, "MonthLabels"),
			MonthValues:  float64SliceFieldC(resp, "MonthValues"),
			TopSuppliers: topSuppliers,
			Recent:       recent,
		}, nil
	}
}

// ---------------------------------------------------------------------------
// Expense dashboard wiring (kind="expense")
// ---------------------------------------------------------------------------

// wireExpenseDashboard sets expDeps.GetExpenseDashboardPageData if
// useCases.Expenditure.Dashboard is non-nil.
func wireExpenseDashboard(deps *expendituremod.ModuleDeps, useCases *consumer.UseCases) {
	if useCases == nil || useCases.Expenditure == nil || useCases.Expenditure.Dashboard == nil {
		return
	}
	uc := reflect.ValueOf(useCases.Expenditure.Dashboard)
	deps.GetExpenseDashboardPageData = func(ctx context.Context, req *expenseboard.Request) (*expenseboard.Response, error) {
		now := time.Now()
		if req != nil && !req.Now.IsZero() {
			now = req.Now
		}
		extra := map[string]any{"Kind": "expense"}
		resp, err := callDashboardExecute(uc, ctx, "", now, extra)
		if err != nil || !resp.IsValid() {
			return nil, err
		}
		// Stats sub-struct — expense uses different field mapping
		var stats expenseboard.Stats
		if s := resp.FieldByName("Stats"); s.IsValid() {
			stats.PendingApprovalCount = s.FieldByName("OpenCount").Int()
			stats.ApprovedMTD = s.FieldByName("TotalMTD").Int()
			stats.ReimbursableMTD = s.FieldByName("ReimbursableMTD").Int()
			stats.CategoriesUsed = s.FieldByName("CategoryCount").Int()
		}
		// Recent: []*expenditurepb.Expenditure — same proto type
		var recent []*expenditurepb.Expenditure
		if f := resp.FieldByName("Recent"); f.IsValid() && !f.IsNil() {
			if v, ok := f.Interface().([]*expenditurepb.Expenditure); ok {
				recent = v
			}
		}
		return &expenseboard.Response{
			Stats:          stats,
			CategoryLabels: stringSliceFieldC(resp, "CategoryLabels"),
			CategoryValues: float64SliceFieldC(resp, "CategoryValues"),
			Recent:         recent,
		}, nil
	}
}
