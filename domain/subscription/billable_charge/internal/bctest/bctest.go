// Package bctest holds the in-memory fakes and request helpers shared by the
// view tests of this module (test-only).
package bctest

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"testing/fstest"

	clientpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client"
	enumspb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/enums"
	documentseriespb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/document_series"
	recoverydocumentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/recovery_document"
	billablechargepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/billable_charge"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	rd "github.com/erniealice/centymo-golang/domain/revenue/recovery_document"
	subscription "github.com/erniealice/centymo-golang/domain/subscription"
	bc "github.com/erniealice/centymo-golang/domain/subscription/billable_charge"
	bcaction "github.com/erniealice/centymo-golang/domain/subscription/billable_charge/action"
	bclist "github.com/erniealice/centymo-golang/domain/subscription/billable_charge/list"
)

type CodedErr string

func (e CodedErr) Error() string { return "x: " + string(e) }

func (e CodedErr) ErrorCode() string { return string(e) }

type MuxRegistrar struct{ Mux *http.ServeMux }

func (m MuxRegistrar) handle(method, path string, v view.View) {
	m.Mux.HandleFunc(method+" "+path, func(w http.ResponseWriter, r *http.Request) {
		res := v.Handle(r.Context(), &view.ViewContext{Request: r})
		w.WriteHeader(res.StatusCode)
	})
}

func (m MuxRegistrar) GET(path string, v view.View, _ ...string) { m.handle(http.MethodGet, path, v) }

func (m MuxRegistrar) POST(path string, v view.View, _ ...string) { m.handle(http.MethodPost, path, v) }

// fake is an in-memory set of the use cases the views consume.
type Fake struct {
	// Clients feeds ListClients (paged by the request); ListClientPages counts the calls.
	Clients         []*clientpb.Client
	ListClientPages int

	Charges      []*billablechargepb.BillableCharge
	Series       []*documentseriespb.DocumentSeries
	IssueErr     error
	AdjustErr    error
	IssuedReqs   []*recoverydocumentpb.IssueRecoveryDocumentsRequest
	AdjustedReqs []*billablechargepb.AdjustBillableChargeRequest
}

func (f *Fake) UseCases() *bc.UseCases {
	return &bc.UseCases{
		ListClients: func(_ context.Context, r *clientpb.ListClientsRequest) (*clientpb.ListClientsResponse, error) {
			f.ListClientPages++
			limit, page := int(r.GetPagination().GetLimit()), int(r.GetPagination().GetOffset().GetPage())
			if limit <= 0 || page <= 0 {
				return &clientpb.ListClientsResponse{Data: f.Clients, Success: true}, nil
			}
			lo, hi := (page-1)*limit, page*limit
			if lo > len(f.Clients) {
				lo = len(f.Clients)
			}
			if hi > len(f.Clients) {
				hi = len(f.Clients)
			}
			return &clientpb.ListClientsResponse{Data: f.Clients[lo:hi], Success: true}, nil
		},
		ListBillableCharges: func(_ context.Context, r *billablechargepb.ListBillableChargesRequest) (*billablechargepb.ListBillableChargesResponse, error) {
			var out []*billablechargepb.BillableCharge
			for _, c := range f.Charges {
				if r.Status != nil && c.GetStatus() != r.GetStatus() {
					continue
				}
				if len(r.GetFilters().GetFilters()) > 0 && c.GetId() != r.GetFilters().GetFilters()[0].GetStringFilter().GetValue() {
					continue
				}
				out = append(out, c)
			}
			return &billablechargepb.ListBillableChargesResponse{Data: out, Success: true}, nil
		},
		AdjustBillableCharge: func(_ context.Context, r *billablechargepb.AdjustBillableChargeRequest) (*billablechargepb.AdjustBillableChargeResponse, error) {
			f.AdjustedReqs = append(f.AdjustedReqs, r)
			return &billablechargepb.AdjustBillableChargeResponse{Success: f.AdjustErr == nil}, f.AdjustErr
		},
		IssueRecoveryDocuments: func(_ context.Context, r *recoverydocumentpb.IssueRecoveryDocumentsRequest) (*recoverydocumentpb.IssueRecoveryDocumentsResponse, error) {
			f.IssuedReqs = append(f.IssuedReqs, r)
			return &recoverydocumentpb.IssueRecoveryDocumentsResponse{Success: f.IssueErr == nil}, f.IssueErr
		},
		ListDocumentSeries: func(context.Context, *documentseriespb.ListDocumentSeriesRequest) (*documentseriespb.ListDocumentSeriesResponse, error) {
			return &documentseriespb.ListDocumentSeriesResponse{Data: f.Series, Success: true}, nil
		},
	}
}

func Charge(id string, st billablechargepb.BillableChargeStatus, kind billablechargepb.BillableChargeKind, amount int64) *billablechargepb.BillableCharge {
	client, from, to := "c-"+id, "2026-01-01", "2026-01-31"
	return &billablechargepb.BillableCharge{Id: id, Status: st, ChargeKind: kind, Amount: amount, Currency: "PHP", ClientId: &client, ServiceFrom: &from, ServiceTo: &to}
}

const (
	Open   = billablechargepb.BillableChargeStatus_BILLABLE_CHARGE_STATUS_OPEN
	Issued = billablechargepb.BillableChargeStatus_BILLABLE_CHARGE_STATUS_ISSUED
	Orig   = billablechargepb.BillableChargeKind_BILLABLE_CHARGE_KIND_ORIGINAL
	Corr   = billablechargepb.BillableChargeKind_BILLABLE_CHARGE_KIND_CORRECTION
)

func Seed() *Fake {
	return &Fake{
		Charges: []*billablechargepb.BillableCharge{
			Charge("o1", Open, Orig, 150000),
			Charge("o2", Open, Orig, 50000),
			Charge("i1", Issued, Orig, 100000),
			Charge("i2", Issued, Corr, -20000),
		},
		Series: []*documentseriespb.DocumentSeries{
			{Id: "s-active", Code: "REC", Status: documentseriespb.DocumentSeriesStatus_DOCUMENT_SERIES_STATUS_ACTIVE, DocumentKind: enumspb.ChargeDocumentKind_CHARGE_DOCUMENT_KIND_RECOVERY_DOCUMENT},
			{Id: "s-invoice", Code: "INV", Status: documentseriespb.DocumentSeriesStatus_DOCUMENT_SERIES_STATUS_ACTIVE, DocumentKind: enumspb.ChargeDocumentKind_CHARGE_DOCUMENT_KIND_INVOICE},
			{Id: "s-retired", Code: "OLD", Status: documentseriespb.DocumentSeriesStatus_DOCUMENT_SERIES_STATUS_RETIRED, DocumentKind: enumspb.ChargeDocumentKind_CHARGE_DOCUMENT_KIND_RECOVERY_DOCUMENT},
		},
	}
}

func CtxWith(codes ...string) context.Context {
	return view.WithUserPermissions(context.Background(), types.NewUserPermissions(codes))
}

func Common() pyeza.CommonLabels {
	var cl pyeza.CommonLabels
	cl.Errors.MissingPermission = "Missing permission: %s"
	cl.Errors.PermissionDenied = "Permission denied"
	cl.Errors.General = "Something went wrong"
	cl.Errors.InvalidFormData = "Invalid form"
	return cl
}

func Request(method, target, form string, pathValues ...string) *view.ViewContext {
	var r *http.Request
	if form != "" {
		r = httptest.NewRequest(method, target, strings.NewReader(form))
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		r = httptest.NewRequest(method, target, nil)
	}
	for i := 0; i+1 < len(pathValues); i += 2 {
		r.SetPathValue(pathValues[i], pathValues[i+1])
	}
	return &view.ViewContext{Request: r, CurrentPath: r.URL.Path}
}

func ActionDeps(f *Fake) *bcaction.Deps {
	return &bcaction.Deps{
		Routes: bc.DefaultRoutes(), Labels: bc.DefaultLabels(), RecoveryLabels: rd.DefaultLabels(),
		CommonLabels: Common(), UseCases: f.UseCases(),
		NewKey: func() string { return "key-fixed" },
	}
}

func ListDeps(f *Fake) *bclist.Deps {
	return &bclist.Deps{Routes: bc.DefaultRoutes(), Labels: bc.DefaultLabels(), CommonLabels: Common(), TableLabels: types.TableLabels{}, UseCases: f.UseCases()}
}

func NewModule(f *Fake) *subscription.BillableChargeModule {
	return subscription.NewBillableChargeModule(&subscription.BillableChargeModuleDeps{
		Routes: bc.DefaultRoutes(), Labels: bc.DefaultLabels(), RecoveryLabels: rd.DefaultLabels(),
		CommonLabels: Common(), TableLabels: types.TableLabels{}, UseCases: f.UseCases(),
	})
}

func Renderer(t *testing.T) *pyeza.HTMLRenderer {
	t.Helper()
	shell := fstest.MapFS{"app-shell.html": {Data: []byte(`{{define "app-shell"}}[shell]{{end}}`)}}
	r := pyeza.NewHTMLRendererFromFS(pyeza.SharedFS, shell, bc.TemplatesFS)
	r.SetRouteMap(bc.DefaultRoutes().RouteMap())
	if err := r.Init(); err != nil {
		t.Fatalf("renderer init: %v", err)
	}
	return r
}

func RenderTo(t *testing.T, r *pyeza.HTMLRenderer, name string, data any) string {
	t.Helper()
	w := httptest.NewRecorder()
	if err := r.Render(w, name, data); err != nil {
		t.Fatalf("render %s: %v", name, err)
	}
	return w.Body.String()
}

func MustContain(t *testing.T, html string, wants ...string) {
	t.Helper()
	for _, w := range wants {
		if !strings.Contains(html, w) {
			t.Errorf("rendered HTML lacks %q", w)
		}
	}
}

// injectCommon mimics the ViewAdapter, which injects CommonLabels into the view
// model by reflection in production.
func InjectCommon(data any) any {
	f := reflect.ValueOf(data).Elem().FieldByName("CommonLabels")
	if f.IsValid() && f.CanSet() {
		f.Set(reflect.ValueOf(Common()))
	}
	return data
}
