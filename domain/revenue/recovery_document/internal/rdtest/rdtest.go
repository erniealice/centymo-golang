// Package rdtest holds the in-memory fakes and request helpers shared by the
// view tests of this module (test-only).
package rdtest

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"testing/fstest"

	clientpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client"
	recoverydocumentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/recovery_document"
	recoverydocumentlinepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/recovery_document_line"
	collectionapplicationpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/collection_application"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	revenue "github.com/erniealice/centymo-golang/domain/revenue"
	rd "github.com/erniealice/centymo-golang/domain/revenue/recovery_document"
	rdaction "github.com/erniealice/centymo-golang/domain/revenue/recovery_document/action"
	rddetail "github.com/erniealice/centymo-golang/domain/revenue/recovery_document/detail"
	rdlist "github.com/erniealice/centymo-golang/domain/revenue/recovery_document/list"
	ra "github.com/erniealice/centymo-golang/domain/treasury/collection_application"
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

type Fake struct {
	// Clients feeds ListClients (paged by the request) and ReadClient.
	Clients         []*clientpb.Client
	ReadClients     []string // ids read through ReadClient
	ListClientPages int
	Docs            []*recoverydocumentpb.RecoveryDocument
	Lines           map[string][]*recoverydocumentlinepb.RecoveryDocumentLine
	Apps            []*collectionapplicationpb.CollectionApplication
	VoidErr         error
	VoidedReqs      []*recoverydocumentpb.VoidRecoveryDocumentRequest
}

func (f *Fake) UseCases() *rd.UseCases {
	return &rd.UseCases{
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
		ReadClient: func(_ context.Context, r *clientpb.ReadClientRequest) (*clientpb.ReadClientResponse, error) {
			f.ReadClients = append(f.ReadClients, r.GetData().GetId())
			for _, c := range f.Clients {
				if c.GetId() == r.GetData().GetId() {
					return &clientpb.ReadClientResponse{Data: []*clientpb.Client{c}, Success: true}, nil
				}
			}
			return &clientpb.ReadClientResponse{Success: true}, nil
		},
		ListRecoveryDocuments: func(_ context.Context, r *recoverydocumentpb.ListRecoveryDocumentsRequest) (*recoverydocumentpb.ListRecoveryDocumentsResponse, error) {
			var out []*recoverydocumentpb.RecoveryDocument
			for _, d := range f.Docs {
				if r.Status == nil || d.GetStatus() == r.GetStatus() {
					out = append(out, d)
				}
			}
			return &recoverydocumentpb.ListRecoveryDocumentsResponse{Data: out, Success: true}, nil
		},
		ReadRecoveryDocument: func(_ context.Context, r *recoverydocumentpb.ReadRecoveryDocumentRequest) (*recoverydocumentpb.ReadRecoveryDocumentResponse, error) {
			for _, d := range f.Docs {
				if d.GetId() == r.GetData().GetId() {
					var notes []*recoverydocumentpb.RecoveryDocument
					for _, n := range f.Docs {
						if n.GetCorrectsDocumentId() == d.GetId() {
							notes = append(notes, n)
						}
					}
					return &recoverydocumentpb.ReadRecoveryDocumentResponse{Data: []*recoverydocumentpb.RecoveryDocument{d}, Lines: f.Lines[d.GetId()], CreditNotes: notes, Success: true}, nil
				}
			}
			return nil, CodedErr("not_found")
		},
		VoidRecoveryDocument: func(_ context.Context, r *recoverydocumentpb.VoidRecoveryDocumentRequest) (*recoverydocumentpb.VoidRecoveryDocumentResponse, error) {
			f.VoidedReqs = append(f.VoidedReqs, r)
			return &recoverydocumentpb.VoidRecoveryDocumentResponse{Success: f.VoidErr == nil}, f.VoidErr
		},
		ListCollectionApplications: func(context.Context, *collectionapplicationpb.ListCollectionApplicationsRequest) (*collectionapplicationpb.ListCollectionApplicationsResponse, error) {
			return &collectionapplicationpb.ListCollectionApplicationsResponse{Data: f.Apps, Success: true}, nil
		},
	}
}

const (
	StIssued = recoverydocumentpb.RecoveryDocumentStatus_RECOVERY_DOCUMENT_STATUS_ISSUED
	StVoid   = recoverydocumentpb.RecoveryDocumentStatus_RECOVERY_DOCUMENT_STATUS_VOID
)

func Doc(id, number string, st recoverydocumentpb.RecoveryDocumentStatus, typ recoverydocumentpb.RecoveryDocumentType, total int64, corrects string) *recoverydocumentpb.RecoveryDocument {
	d := &recoverydocumentpb.RecoveryDocument{Id: id, DocumentNumber: number, Status: st, DocumentType: typ, TotalAmount: total, Currency: "PHP", ClientId: "c1", Active: true}
	if corrects != "" {
		d.CorrectsDocumentId = &corrects
	}
	return d
}

func Seed() *Fake {
	from, to, desc := "2026-01-01", "2026-01-31", "Electricity"
	return &Fake{
		Docs: []*recoverydocumentpb.RecoveryDocument{
			Doc("d1", "REC-000001", StIssued, recoverydocumentpb.RecoveryDocumentType_RECOVERY_DOCUMENT_TYPE_STATEMENT, 200000, ""),
			Doc("d2", "REC-000002", StVoid, recoverydocumentpb.RecoveryDocumentType_RECOVERY_DOCUMENT_TYPE_STATEMENT, 50000, ""),
			Doc("d3", "REC-000003", StIssued, recoverydocumentpb.RecoveryDocumentType_RECOVERY_DOCUMENT_TYPE_CREDIT_NOTE, -30000, "d1"),
		},
		Lines: map[string][]*recoverydocumentlinepb.RecoveryDocumentLine{
			"d1": {{Id: "l1", RecoveryDocumentId: "d1", Description: &desc, Amount: 200000, Currency: "PHP", ServiceFrom: &from, ServiceTo: &to}},
		},
		Apps: []*collectionapplicationpb.CollectionApplication{
			{Id: "a1", Active: true, Amount: 120000, Currency: "PHP", Status: collectionapplicationpb.ApplicationStatus_APPLICATION_STATUS_APPLIED, ApplicationKind: collectionapplicationpb.ApplicationKind_APPLICATION_KIND_CASH},
			{Id: "a2", Active: true, Amount: 40000, Currency: "PHP", Status: collectionapplicationpb.ApplicationStatus_APPLICATION_STATUS_REVERSED, ApplicationKind: collectionapplicationpb.ApplicationKind_APPLICATION_KIND_CASH},
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

func DetailDeps(f *Fake) *rddetail.Deps {
	return &rddetail.Deps{
		Routes: rd.DefaultRoutes(), Labels: rd.DefaultLabels(), CommonLabels: Common(), TableLabels: types.TableLabels{}, UseCases: f.UseCases(),
		ApplicationLabels: ra.DefaultLabels(), ReverseURL: ra.ReverseURL, ReceiveApplyURL: ra.ReceiveApplyURL, ApplicationsTableID: ra.ApplicationsTableID,
	}
}

func ListDeps(f *Fake) *rdlist.Deps {
	return &rdlist.Deps{Routes: rd.DefaultRoutes(), Labels: rd.DefaultLabels(), CommonLabels: Common(), TableLabels: types.TableLabels{}, UseCases: f.UseCases()}
}

func ActionDeps(f *Fake) *rdaction.Deps {
	return &rdaction.Deps{Routes: rd.DefaultRoutes(), Labels: rd.DefaultLabels(), CommonLabels: Common(), UseCases: f.UseCases()}
}

func NewModule(f *Fake) *revenue.RecoveryDocumentModule {
	return revenue.NewRecoveryDocumentModule(&revenue.RecoveryDocumentModuleDeps{
		Routes: rd.DefaultRoutes(), Labels: rd.DefaultLabels(), CommonLabels: Common(), TableLabels: types.TableLabels{},
		UseCases: f.UseCases(), ApplicationLabels: ra.DefaultLabels(), ReverseURL: ra.ReverseURL, ReceiveApplyURL: ra.ReceiveApplyURL, ApplicationsTableID: ra.ApplicationsTableID,
	})
}

func Renderer(t *testing.T) *pyeza.HTMLRenderer {
	t.Helper()
	shell := fstest.MapFS{"app-shell.html": {Data: []byte(`{{define "app-shell"}}[shell]{{end}}`)}}
	r := pyeza.NewHTMLRendererFromFS(pyeza.SharedFS, shell, rd.TemplatesFS)
	r.SetRouteMap(rd.DefaultRoutes().RouteMap())
	if err := r.Init(); err != nil {
		t.Fatalf("renderer init: %v", err)
	}
	return r
}

func RenderTo(t *testing.T, name string, data any) string {
	t.Helper()
	w := httptest.NewRecorder()
	if err := Renderer(t).Render(w, name, data); err != nil {
		t.Fatalf("render %s: %v", name, err)
	}
	return w.Body.String()
}

func InjectCommon(data any) any {
	f := reflect.ValueOf(data).Elem().FieldByName("CommonLabels")
	if f.IsValid() && f.CanSet() {
		f.Set(reflect.ValueOf(Common()))
	}
	return data
}

func MustContain(t *testing.T, html string, wants ...string) {
	t.Helper()
	for _, w := range wants {
		if !strings.Contains(html, w) {
			t.Errorf("rendered HTML lacks %q", w)
		}
	}
}
