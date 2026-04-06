package action

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	centymo "github.com/erniealice/centymo-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	revenuepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/revenue"
	revenuelineitempb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/revenue_line_item"
)

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func testLabels() centymo.RevenueLabels {
	return centymo.RevenueLabels{
		Errors: centymo.RevenueErrorLabels{
			PermissionDenied:        "Permission denied",
			InvalidFormData:         "Invalid form data",
			NotFound:                "Not found",
			IDRequired:              "ID is required",
			NoIDsProvided:           "No IDs provided",
			InvalidStatus:           "Invalid status",
			InvalidTargetStatus:     "Invalid target status",
			NoItemsCannotComplete:   "Cannot complete: no items",
			HasPaymentsCannotCancel: "Cannot cancel: has payments",
			BulkHasPayments:         "%d of %d have payments",
			BulkNoItems:             "%d of %d have no items",
		},
	}
}

func testRoutes() centymo.RevenueRoutes {
	return centymo.RevenueRoutes{
		AddURL:    "/sales/add",
		EditURL:   "/sales/{id}/edit",
		DeleteURL: "/sales/delete",
		DetailURL: "/sales/{id}",
	}
}

// mockDB is a minimal DataSource mock for tests.
type mockDB struct {
	listSimple func(ctx context.Context, collection string) ([]map[string]any, error)
}

func (m *mockDB) ListSimple(ctx context.Context, collection string) ([]map[string]any, error) {
	if m.listSimple != nil {
		return m.listSimple(ctx, collection)
	}
	return nil, nil
}
func (m *mockDB) Create(_ context.Context, _ string, _ map[string]any) (map[string]any, error) {
	return nil, nil
}
func (m *mockDB) Read(_ context.Context, _ string, _ string) (map[string]any, error) {
	return nil, nil
}
func (m *mockDB) Update(_ context.Context, _ string, _ string, _ map[string]any) (map[string]any, error) {
	return nil, nil
}
func (m *mockDB) Delete(_ context.Context, _ string, _ string) error { return nil }
func (m *mockDB) HardDelete(_ context.Context, _ string, _ string) error {
	return nil
}

// ctxWithPerms returns a context with the given permission codes.
func ctxWithPerms(codes ...string) context.Context {
	perms := types.NewUserPermissions(codes)
	return view.WithUserPermissions(context.Background(), perms)
}

// ctxNoPerms returns a context with an empty permissions set (denies everything).
func ctxNoPerms() context.Context {
	perms := types.NewUserPermissions(nil)
	return view.WithUserPermissions(context.Background(), perms)
}

// postForm creates a POST request with URL-encoded form fields.
func postForm(target string, values url.Values) *http.Request {
	body := values.Encode()
	req := httptest.NewRequest(http.MethodPost, target, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req
}

// ---------------------------------------------------------------------------
// NewAddAction — permission denied
// ---------------------------------------------------------------------------

func TestNewAddAction_PermissionDenied(t *testing.T) {
	t.Parallel()

	deps := &Deps{
		Routes: testRoutes(),
		Labels: testLabels(),
		DB:     &mockDB{},
	}

	v := NewAddAction(deps)
	ctx := ctxNoPerms()
	vc := &view.ViewContext{
		Request: httptest.NewRequest(http.MethodGet, "/sales/add", nil),
	}

	result := v.Handle(ctx, vc)
	if result.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("StatusCode = %d, want %d", result.StatusCode, http.StatusUnprocessableEntity)
	}
	if result.Headers["HX-Error-Message"] != "Permission denied" {
		t.Errorf("HX-Error-Message = %q, want %q", result.Headers["HX-Error-Message"], "Permission denied")
	}
}

// ---------------------------------------------------------------------------
// NewAddAction — POST with missing fields
// ---------------------------------------------------------------------------

func TestNewAddAction_POST_MissingFields(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		form   url.Values
		wantOK bool
	}{
		{
			name:   "completely empty form",
			form:   url.Values{},
			wantOK: true, // no client-side required field validation; CreateRevenue called
		},
		{
			name:   "missing name",
			form:   url.Values{"currency": {"PHP"}, "status": {"draft"}},
			wantOK: true,
		},
		{
			name:   "missing currency",
			form:   url.Values{"name": {"Test Sale"}, "status": {"draft"}},
			wantOK: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			createCalled := false
			deps := &Deps{
				Routes: testRoutes(),
				Labels: testLabels(),
				DB:     &mockDB{},
				CreateRevenue: func(_ context.Context, _ *revenuepb.CreateRevenueRequest) (*revenuepb.CreateRevenueResponse, error) {
					createCalled = true
					return &revenuepb.CreateRevenueResponse{
						Success: true,
						Data:    []*revenuepb.Revenue{{Id: "rev-new"}},
					}, nil
				},
			}

			v := NewAddAction(deps)
			ctx := ctxWithPerms("invoice:create")
			vc := &view.ViewContext{
				Request: postForm("/sales/add", tt.form),
			}

			result := v.Handle(ctx, vc)
			if tt.wantOK {
				// Successful creation redirects
				if result.StatusCode != http.StatusOK {
					t.Errorf("StatusCode = %d, want %d", result.StatusCode, http.StatusOK)
				}
				if !createCalled {
					t.Error("CreateRevenue should have been called")
				}
			}
		})
	}
}

// ---------------------------------------------------------------------------
// NewAddAction — POST create failure propagates error
// ---------------------------------------------------------------------------

func TestNewAddAction_POST_CreateFailure(t *testing.T) {
	t.Parallel()

	deps := &Deps{
		Routes: testRoutes(),
		Labels: testLabels(),
		DB:     &mockDB{},
		CreateRevenue: func(_ context.Context, _ *revenuepb.CreateRevenueRequest) (*revenuepb.CreateRevenueResponse, error) {
			return nil, fmt.Errorf("database unavailable")
		},
	}

	v := NewAddAction(deps)
	ctx := ctxWithPerms("invoice:create")
	vc := &view.ViewContext{
		Request: postForm("/sales/add", url.Values{"name": {"Test"}}),
	}

	result := v.Handle(ctx, vc)
	if result.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("StatusCode = %d, want %d", result.StatusCode, http.StatusUnprocessableEntity)
	}
}

// ---------------------------------------------------------------------------
// NewEditAction — permission denied
// ---------------------------------------------------------------------------

func TestNewEditAction_PermissionDenied(t *testing.T) {
	t.Parallel()

	deps := &Deps{
		Routes: testRoutes(),
		Labels: testLabels(),
		DB:     &mockDB{},
	}

	v := NewEditAction(deps)
	ctx := ctxNoPerms()
	vc := &view.ViewContext{
		Request: httptest.NewRequest(http.MethodGet, "/sales/rev-001/edit", nil),
	}

	result := v.Handle(ctx, vc)
	if result.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("StatusCode = %d, want %d", result.StatusCode, http.StatusUnprocessableEntity)
	}
}

// ---------------------------------------------------------------------------
// NewDeleteAction — permission denied
// ---------------------------------------------------------------------------

func TestNewDeleteAction_PermissionDenied(t *testing.T) {
	t.Parallel()

	deps := &Deps{
		Routes: testRoutes(),
		Labels: testLabels(),
		DB:     &mockDB{},
	}

	v := NewDeleteAction(deps)
	ctx := ctxNoPerms()
	vc := &view.ViewContext{
		Request: httptest.NewRequest(http.MethodPost, "/sales/delete?id=rev-001", nil),
	}

	result := v.Handle(ctx, vc)
	if result.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("StatusCode = %d, want %d", result.StatusCode, http.StatusUnprocessableEntity)
	}
}

// ---------------------------------------------------------------------------
// NewDeleteAction — missing ID
// ---------------------------------------------------------------------------

func TestNewDeleteAction_MissingID(t *testing.T) {
	t.Parallel()

	deps := &Deps{
		Routes: testRoutes(),
		Labels: testLabels(),
		DB:     &mockDB{},
	}

	v := NewDeleteAction(deps)
	ctx := ctxWithPerms("invoice:delete")
	vc := &view.ViewContext{
		Request: postForm("/sales/delete", url.Values{}),
	}

	result := v.Handle(ctx, vc)
	if result.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("StatusCode = %d, want %d", result.StatusCode, http.StatusUnprocessableEntity)
	}
	if result.Headers["HX-Error-Message"] != "ID is required" {
		t.Errorf("HX-Error-Message = %q, want %q", result.Headers["HX-Error-Message"], "ID is required")
	}
}

// ---------------------------------------------------------------------------
// NewSetStatusAction — permission denied
// ---------------------------------------------------------------------------

func TestNewSetStatusAction_PermissionDenied(t *testing.T) {
	t.Parallel()

	deps := &Deps{
		Routes: testRoutes(),
		Labels: testLabels(),
		DB:     &mockDB{},
	}

	v := NewSetStatusAction(deps)
	ctx := ctxNoPerms()
	vc := &view.ViewContext{
		Request: httptest.NewRequest(http.MethodPost, "/sales/set-status?id=rev-001&status=complete", nil),
	}

	result := v.Handle(ctx, vc)
	if result.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("StatusCode = %d, want %d", result.StatusCode, http.StatusUnprocessableEntity)
	}
}

// ---------------------------------------------------------------------------
// NewSetStatusAction — invalid status values
// ---------------------------------------------------------------------------

func TestNewSetStatusAction_InvalidStatus(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		status string
	}{
		{"empty status", ""},
		{"invalid status", "bogus"},
		{"uppercase", "Complete"},
		{"pending is not valid", "pending"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			deps := &Deps{
				Routes: testRoutes(),
				Labels: testLabels(),
				DB:     &mockDB{},
			}

			v := NewSetStatusAction(deps)
			ctx := ctxWithPerms("invoice:update")
			vc := &view.ViewContext{
				Request: httptest.NewRequest(
					http.MethodPost,
					"/sales/set-status?id=rev-001&status="+tt.status,
					nil,
				),
			}

			result := v.Handle(ctx, vc)
			if result.StatusCode != http.StatusUnprocessableEntity {
				t.Errorf("StatusCode = %d, want %d for status %q", result.StatusCode, http.StatusUnprocessableEntity, tt.status)
			}
			if result.Headers["HX-Error-Message"] != "Invalid status" {
				t.Errorf("HX-Error-Message = %q, want %q", result.Headers["HX-Error-Message"], "Invalid status")
			}
		})
	}
}

// ---------------------------------------------------------------------------
// NewSetStatusAction — missing ID
// ---------------------------------------------------------------------------

func TestNewSetStatusAction_MissingID(t *testing.T) {
	t.Parallel()

	deps := &Deps{
		Routes: testRoutes(),
		Labels: testLabels(),
		DB:     &mockDB{},
	}

	v := NewSetStatusAction(deps)
	ctx := ctxWithPerms("invoice:update")
	vc := &view.ViewContext{
		Request: postForm("/sales/set-status", url.Values{"status": {"complete"}}),
	}

	result := v.Handle(ctx, vc)
	if result.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("StatusCode = %d, want %d", result.StatusCode, http.StatusUnprocessableEntity)
	}
	if result.Headers["HX-Error-Message"] != "ID is required" {
		t.Errorf("HX-Error-Message = %q, want %q", result.Headers["HX-Error-Message"], "ID is required")
	}
}

// ---------------------------------------------------------------------------
// NewSetStatusAction — complete with zero line items
// ---------------------------------------------------------------------------

func TestNewSetStatusAction_CompleteWithNoItems(t *testing.T) {
	t.Parallel()

	deps := &Deps{
		Routes: testRoutes(),
		Labels: testLabels(),
		DB:     &mockDB{},
		ListRevenueLineItems: func(_ context.Context, _ *revenuelineitempb.ListRevenueLineItemsRequest) (*revenuelineitempb.ListRevenueLineItemsResponse, error) {
			return &revenuelineitempb.ListRevenueLineItemsResponse{
				Data: nil, // no line items
			}, nil
		},
	}

	v := NewSetStatusAction(deps)
	ctx := ctxWithPerms("invoice:update")
	vc := &view.ViewContext{
		Request: httptest.NewRequest(http.MethodPost, "/sales/set-status?id=rev-001&status=complete", nil),
	}

	result := v.Handle(ctx, vc)
	if result.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("StatusCode = %d, want %d", result.StatusCode, http.StatusUnprocessableEntity)
	}
	if result.Headers["HX-Error-Message"] != "Cannot complete: no items" {
		t.Errorf("HX-Error-Message = %q, want %q", result.Headers["HX-Error-Message"], "Cannot complete: no items")
	}
}

// ---------------------------------------------------------------------------
// NewSetStatusAction — cancel with payments
// ---------------------------------------------------------------------------

func TestNewSetStatusAction_CancelWithPayments(t *testing.T) {
	t.Parallel()

	deps := &Deps{
		Routes: testRoutes(),
		Labels: testLabels(),
		DB: &mockDB{
			listSimple: func(_ context.Context, collection string) ([]map[string]any, error) {
				if collection == "revenue_payment" {
					return []map[string]any{
						{"id": "pay-001", "revenue_id": "rev-001", "amount": 100.0},
					}, nil
				}
				return nil, nil
			},
		},
	}

	v := NewSetStatusAction(deps)
	ctx := ctxWithPerms("invoice:update")
	vc := &view.ViewContext{
		Request: httptest.NewRequest(http.MethodPost, "/sales/set-status?id=rev-001&status=cancelled", nil),
	}

	result := v.Handle(ctx, vc)
	if result.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("StatusCode = %d, want %d", result.StatusCode, http.StatusUnprocessableEntity)
	}
	if result.Headers["HX-Error-Message"] != "Cannot cancel: has payments" {
		t.Errorf("HX-Error-Message = %q, want %q", result.Headers["HX-Error-Message"], "Cannot cancel: has payments")
	}
}

// ---------------------------------------------------------------------------
// NewBulkDeleteAction — permission denied
// ---------------------------------------------------------------------------

func TestNewBulkDeleteAction_PermissionDenied(t *testing.T) {
	t.Parallel()

	deps := &Deps{
		Routes: testRoutes(),
		Labels: testLabels(),
		DB:     &mockDB{},
	}

	v := NewBulkDeleteAction(deps)
	ctx := ctxNoPerms()
	vc := &view.ViewContext{
		Request: postForm("/sales/bulk-delete", url.Values{"id": {"rev-001", "rev-002"}}),
	}

	result := v.Handle(ctx, vc)
	if result.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("StatusCode = %d, want %d", result.StatusCode, http.StatusUnprocessableEntity)
	}
}

// ---------------------------------------------------------------------------
// NewBulkDeleteAction — no IDs provided
// ---------------------------------------------------------------------------

func TestNewBulkDeleteAction_NoIDs(t *testing.T) {
	t.Parallel()

	deps := &Deps{
		Routes: testRoutes(),
		Labels: testLabels(),
		DB:     &mockDB{},
	}

	v := NewBulkDeleteAction(deps)
	ctx := ctxWithPerms("invoice:delete")
	vc := &view.ViewContext{
		Request: postForm("/sales/bulk-delete", url.Values{}),
	}

	result := v.Handle(ctx, vc)
	if result.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("StatusCode = %d, want %d", result.StatusCode, http.StatusUnprocessableEntity)
	}
	if result.Headers["HX-Error-Message"] != "No IDs provided" {
		t.Errorf("HX-Error-Message = %q, want %q", result.Headers["HX-Error-Message"], "No IDs provided")
	}
}

// ---------------------------------------------------------------------------
// NewBulkSetStatusAction — invalid target status
// ---------------------------------------------------------------------------

func TestNewBulkSetStatusAction_InvalidStatus(t *testing.T) {
	t.Parallel()

	deps := &Deps{
		Routes: testRoutes(),
		Labels: testLabels(),
		DB:     &mockDB{},
	}

	v := NewBulkSetStatusAction(deps)
	ctx := ctxWithPerms("invoice:update")
	vc := &view.ViewContext{
		Request: postForm("/sales/bulk-set-status", url.Values{
			"id":            {"rev-001"},
			"target_status": {"invalid"},
		}),
	}

	result := v.Handle(ctx, vc)
	if result.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("StatusCode = %d, want %d", result.StatusCode, http.StatusUnprocessableEntity)
	}
	if result.Headers["HX-Error-Message"] != "Invalid target status" {
		t.Errorf("HX-Error-Message = %q, want %q", result.Headers["HX-Error-Message"], "Invalid target status")
	}
}
