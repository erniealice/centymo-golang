package action

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	centymo "github.com/erniealice/centymo-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	subscriptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription"
)

func TestFormatDateForInput(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		dateString string
		dateMillis int64
		want       string
	}{
		{
			name:       "date only string",
			dateString: "2026-04-03",
			want:       "2026-04-03",
		},
		{
			name:       "timestamp string trimmed for input",
			dateString: "2026-04-03T09:30:00Z",
			want:       "2026-04-03",
		},
		{
			name:       "falls back to millis",
			dateMillis: time.Date(2026, time.April, 3, 0, 0, 0, 0, time.UTC).UnixMilli(),
			want:       "2026-04-03",
		},
		{
			name: "empty when no date is available",
			want: "",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := formatDateForInput(tt.dateString, tt.dateMillis)
			if got != tt.want {
				t.Fatalf("formatDateForInput(%q, %d) = %q, want %q", tt.dateString, tt.dateMillis, got, tt.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func testLabels() centymo.SubscriptionLabels {
	return centymo.SubscriptionLabels{
		Errors: centymo.SubscriptionErrorLabels{
			PermissionDenied: "Permission denied",
			InvalidFormData:  "Invalid form data",
			NotFound:         "Not found",
			IDRequired:       "ID is required",
		},
	}
}

func testRoutes() centymo.SubscriptionRoutes {
	return centymo.SubscriptionRoutes{
		AddURL:    "/subscriptions/add",
		EditURL:   "/subscriptions/{id}/edit",
		DeleteURL: "/subscriptions/delete",
	}
}

// ctxWithPerms creates a context with the given permission codes.
func ctxWithPerms(codes ...string) context.Context {
	perms := types.NewUserPermissions(codes)
	return view.WithUserPermissions(context.Background(), perms)
}

// ctxNoPerms creates a context with an empty permissions set (denies everything).
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
	}

	v := NewAddAction(deps)
	ctx := ctxNoPerms()
	vc := &view.ViewContext{
		Request: httptest.NewRequest(http.MethodGet, "/subscriptions/add", nil),
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
		wantOK bool // true if the action should succeed (no server-side required field validation)
	}{
		{
			name:   "completely empty form",
			form:   url.Values{},
			wantOK: true, // CreateSubscription is called; validation is server-side
		},
		{
			name:   "missing client_id",
			form:   url.Values{"price_plan_id": {"plan-001"}, "date_start_string": {"2027-01-01"}},
			wantOK: true,
		},
		{
			name:   "missing price_plan_id",
			form:   url.Values{"client_id": {"client-001"}, "date_start_string": {"2027-01-01"}},
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
				CreateSubscription: func(_ context.Context, req *subscriptionpb.CreateSubscriptionRequest) (*subscriptionpb.CreateSubscriptionResponse, error) {
					createCalled = true
					return &subscriptionpb.CreateSubscriptionResponse{Success: true}, nil
				},
			}

			v := NewAddAction(deps)
			ctx := ctxWithPerms("subscription:create")
			vc := &view.ViewContext{
				Request: postForm("/subscriptions/add", tt.form),
			}

			result := v.Handle(ctx, vc)
			if tt.wantOK {
				if result.StatusCode != http.StatusOK {
					t.Errorf("StatusCode = %d, want %d", result.StatusCode, http.StatusOK)
				}
				if !createCalled {
					t.Error("CreateSubscription should have been called")
				}
			}
		})
	}
}

// ---------------------------------------------------------------------------
// NewAddAction — POST with past start date
// ---------------------------------------------------------------------------

func TestNewAddAction_POST_PastStartDate(t *testing.T) {
	t.Parallel()

	deps := &Deps{
		Routes: testRoutes(),
		Labels: testLabels(),
		CreateSubscription: func(_ context.Context, _ *subscriptionpb.CreateSubscriptionRequest) (*subscriptionpb.CreateSubscriptionResponse, error) {
			t.Error("CreateSubscription should not be called for past start date")
			return &subscriptionpb.CreateSubscriptionResponse{Success: true}, nil
		},
	}

	v := NewAddAction(deps)
	ctx := ctxWithPerms("subscription:create")
	form := url.Values{
		"client_id":         {"client-001"},
		"price_plan_id":     {"plan-001"},
		"date_start_string": {"2020-01-01"}, // past date
	}
	vc := &view.ViewContext{
		Request: postForm("/subscriptions/add", form),
	}

	result := v.Handle(ctx, vc)
	if result.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("StatusCode = %d, want %d for past start date", result.StatusCode, http.StatusUnprocessableEntity)
	}
	if result.Headers["HX-Error-Message"] != "Start date cannot be in the past" {
		t.Errorf("HX-Error-Message = %q, want past-date error", result.Headers["HX-Error-Message"])
	}
}

// ---------------------------------------------------------------------------
// NewAddAction — POST with invalid date format (not rejected, bad parse is skipped)
// ---------------------------------------------------------------------------

func TestNewAddAction_POST_InvalidDateFormat(t *testing.T) {
	t.Parallel()

	createCalled := false
	deps := &Deps{
		Routes: testRoutes(),
		Labels: testLabels(),
		CreateSubscription: func(_ context.Context, _ *subscriptionpb.CreateSubscriptionRequest) (*subscriptionpb.CreateSubscriptionResponse, error) {
			createCalled = true
			return &subscriptionpb.CreateSubscriptionResponse{Success: true}, nil
		},
	}

	v := NewAddAction(deps)
	ctx := ctxWithPerms("subscription:create")
	form := url.Values{
		"client_id":         {"client-001"},
		"price_plan_id":     {"plan-001"},
		"date_start_string": {"not-a-date"}, // invalid date format
	}
	vc := &view.ViewContext{
		Request: postForm("/subscriptions/add", form),
	}

	result := v.Handle(ctx, vc)
	// Invalid date format causes time.Parse to fail, so validation is skipped
	// and CreateSubscription proceeds.
	if result.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d, want %d (invalid date skips validation)", result.StatusCode, http.StatusOK)
	}
	if !createCalled {
		t.Error("CreateSubscription should be called when date parse fails (validation skipped)")
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
	}

	v := NewEditAction(deps)
	ctx := ctxNoPerms()
	vc := &view.ViewContext{
		Request: httptest.NewRequest(http.MethodGet, "/subscriptions/sub-001/edit", nil),
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
	}

	v := NewDeleteAction(deps)
	ctx := ctxNoPerms()
	vc := &view.ViewContext{
		Request: httptest.NewRequest(http.MethodPost, "/subscriptions/delete?id=sub-001", nil),
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
	}

	v := NewDeleteAction(deps)
	ctx := ctxWithPerms("subscription:delete")
	// POST with no id param
	vc := &view.ViewContext{
		Request: postForm("/subscriptions/delete", url.Values{}),
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
// formatDateForInput — additional edge cases
// ---------------------------------------------------------------------------

func TestFormatDateForInput_EdgeCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		dateString string
		dateMillis int64
		want       string
	}{
		{
			name:       "short date string returns as-is",
			dateString: "2026-04",
			want:       "2026-04",
		},
		{
			name:       "negative millis returns empty",
			dateMillis: -1,
			want:       "",
		},
		{
			name:       "zero millis returns empty",
			dateMillis: 0,
			want:       "",
		},
		{
			name:       "date string takes precedence over millis",
			dateString: "2025-01-15",
			dateMillis: time.Date(2026, time.June, 1, 0, 0, 0, 0, time.UTC).UnixMilli(),
			want:       "2025-01-15",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := formatDateForInput(tt.dateString, tt.dateMillis)
			if got != tt.want {
				t.Fatalf("formatDateForInput(%q, %d) = %q, want %q", tt.dateString, tt.dateMillis, got, tt.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// generateCode — sanity check
// ---------------------------------------------------------------------------

func TestGenerateCode(t *testing.T) {
	t.Parallel()

	for i := 0; i < 50; i++ {
		code := generateCode()
		if len(code) != 7 {
			t.Errorf("generateCode() = %q, length %d, want 7", code, len(code))
		}
		// Check no ambiguous chars
		for _, c := range code {
			if c == 'O' || c == 'I' || c == '0' || c == '1' {
				t.Errorf("generateCode() = %q, contains ambiguous char %q", code, string(c))
			}
		}
	}
}
