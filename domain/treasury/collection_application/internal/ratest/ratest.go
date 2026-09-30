// Package ratest holds the in-memory fakes and request helpers shared by the
// view tests of this module (test-only).
package ratest

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"testing/fstest"

	clientpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client"
	collectionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/collection"
	collectionapplicationpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/collection_application"
	collectionmethodpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/collection_method"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	treasury "github.com/erniealice/centymo-golang/domain/treasury"
	ra "github.com/erniealice/centymo-golang/domain/treasury/collection_application"
	raaction "github.com/erniealice/centymo-golang/domain/treasury/collection_application/action"
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
	// Clients overrides the two default clients; ListClients pages it by the
	// request, ReadClient reads from it. ListClientPages counts list calls.
	Clients         []*clientpb.Client
	ReadClients     []string
	ListClientPages int

	ReceiveErr    error
	PreviewErr    error
	ReverseErr    error
	ReceivedReqs  []*collectionapplicationpb.ReceiveAndApplyCollectionRequest
	PreviewedReqs []*collectionapplicationpb.PreviewCollectionApplicationRequest
	ReversedReqs  []*collectionapplicationpb.ReverseCollectionApplicationRequest
}

func (f *Fake) clients() []*clientpb.Client {
	if f.Clients != nil {
		return f.Clients
	}
	return []*clientpb.Client{{Id: "c1", Name: Strp("Tenant One")}, {Id: "c2", Name: Strp("Tenant Two")}}
}

func (f *Fake) UseCases() *ra.UseCases {
	return &ra.UseCases{
		ReceiveAndApplyCollection: func(_ context.Context, r *collectionapplicationpb.ReceiveAndApplyCollectionRequest) (*collectionapplicationpb.ReceiveAndApplyCollectionResponse, error) {
			f.ReceivedReqs = append(f.ReceivedReqs, r)
			if f.ReceiveErr != nil {
				return nil, f.ReceiveErr
			}
			return &collectionapplicationpb.ReceiveAndApplyCollectionResponse{Collection: &collectionpb.Collection{Id: "col-1"}, Success: true}, nil
		},
		PreviewCollectionApplication: func(_ context.Context, r *collectionapplicationpb.PreviewCollectionApplicationRequest) (*collectionapplicationpb.PreviewCollectionApplicationResponse, error) {
			f.PreviewedReqs = append(f.PreviewedReqs, r)
			if f.PreviewErr != nil {
				return nil, f.PreviewErr
			}
			return &collectionapplicationpb.PreviewCollectionApplicationResponse{Plan: &collectionapplicationpb.CollectionApplicationPlan{
				Currency: "PHP", Applied: 150000, Unapplied: 50000,
				Allocations: []*collectionapplicationpb.CollectionApplicationPlanAllocation{
					{TargetKind: collectionapplicationpb.ApplicationTargetKind_APPLICATION_TARGET_KIND_REVENUE, TargetId: "r1", Number: "INV-1", DueDate: "2026-01-05", Balance: 100000, Apply: 100000, Rank: 1},
					{TargetKind: collectionapplicationpb.ApplicationTargetKind_APPLICATION_TARGET_KIND_RECOVERY_DOCUMENT, TargetId: "d1", Number: "REC-000001", DueDate: "2026-01-05", Balance: 80000, Apply: 50000, Rank: 2},
				},
			}, Success: true}, nil
		},
		ReverseCollectionApplication: func(_ context.Context, r *collectionapplicationpb.ReverseCollectionApplicationRequest) (*collectionapplicationpb.ReverseCollectionApplicationResponse, error) {
			f.ReversedReqs = append(f.ReversedReqs, r)
			return &collectionapplicationpb.ReverseCollectionApplicationResponse{Success: f.ReverseErr == nil}, f.ReverseErr
		},
		ListClients: func(_ context.Context, r *clientpb.ListClientsRequest) (*clientpb.ListClientsResponse, error) {
			f.ListClientPages++
			all := f.clients()
			limit, page := int(r.GetPagination().GetLimit()), int(r.GetPagination().GetOffset().GetPage())
			if limit <= 0 || page <= 0 {
				return &clientpb.ListClientsResponse{Data: all, Success: true}, nil
			}
			lo, hi := (page-1)*limit, page*limit
			if lo > len(all) {
				lo = len(all)
			}
			if hi > len(all) {
				hi = len(all)
			}
			return &clientpb.ListClientsResponse{Data: all[lo:hi], Success: true}, nil
		},
		ReadClient: func(_ context.Context, r *clientpb.ReadClientRequest) (*clientpb.ReadClientResponse, error) {
			f.ReadClients = append(f.ReadClients, r.GetData().GetId())
			for _, c := range f.clients() {
				if c.GetId() == r.GetData().GetId() {
					return &clientpb.ReadClientResponse{Data: []*clientpb.Client{c}, Success: true}, nil
				}
			}
			return &clientpb.ReadClientResponse{Success: true}, nil
		},
		ListCollectionMethods: func(context.Context, *collectionmethodpb.ListCollectionMethodsRequest) (*collectionmethodpb.ListCollectionMethodsResponse, error) {
			return &collectionmethodpb.ListCollectionMethodsResponse{Data: []*collectionmethodpb.CollectionMethod{{Id: "m1", Name: "Bank transfer"}}, Success: true}, nil
		},
	}
}

func Strp(s string) *string { return &s }

func CtxWith(codes ...string) context.Context {
	return view.WithUserPermissions(context.Background(), types.NewUserPermissions(codes))
}

func Common() pyeza.CommonLabels {
	var cl pyeza.CommonLabels
	cl.Errors.MissingPermission = "Missing permission: %s"
	cl.Errors.PermissionDenied = "Permission denied"
	cl.Errors.General = "Something went wrong"
	cl.Errors.InvalidFormData = "Invalid form"
	cl.Errors.IDRequired = "ID required"
	cl.Currency.Options = []types.SelectOption{{Value: "PHP", Label: "PHP"}, {Value: "USD", Label: "USD"}}
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

func ActionDeps(f *Fake) *raaction.Deps {
	return &raaction.Deps{
		Routes: ra.DefaultRoutes(), Labels: ra.DefaultLabels(), CommonLabels: Common(), UseCases: f.UseCases(),
		CollectionDetailURL: "/collections/detail/{id}",
	}
}

func NewModule(f *Fake) *treasury.CollectionApplicationModule {
	return treasury.NewCollectionApplicationModule(&treasury.CollectionApplicationModuleDeps{
		Routes: ra.DefaultRoutes(), Labels: ra.DefaultLabels(), CommonLabels: Common(), UseCases: f.UseCases(),
	})
}

func Renderer(t *testing.T) *pyeza.HTMLRenderer {
	t.Helper()
	shell := fstest.MapFS{"app-shell.html": {Data: []byte(`{{define "app-shell"}}[shell]{{end}}`)}}
	r := pyeza.NewHTMLRendererFromFS(pyeza.SharedFS, shell, ra.TemplatesFS)
	r.SetRouteMap(ra.DefaultRoutes().RouteMap())
	if err := r.Init(); err != nil {
		t.Fatalf("renderer init: %v", err)
	}
	return r
}

func RenderTo(t *testing.T, name string, data any) string {
	t.Helper()
	if v := reflect.ValueOf(data); v.Kind() == reflect.Ptr {
		if f := v.Elem().FieldByName("CommonLabels"); f.IsValid() && f.CanSet() {
			f.Set(reflect.ValueOf(Common()))
		}
	}
	w := httptest.NewRecorder()
	if err := Renderer(t).Render(w, name, data); err != nil {
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
