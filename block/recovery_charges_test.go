package block

import (
	"context"
	clientpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client"
	"net/http"
	"strings"
	"testing"

	"github.com/erniealice/espyna-golang/consumer/compose"
	documentseriespb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/document_series"
	recoverydocumentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/recovery_document"
	billablechargepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/billable_charge"
	collectionapplicationpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/collection_application"
	"github.com/erniealice/pyeza-golang/view"
)

func recoveryUnitKeys(units []compose.Unit) []string {
	var keys []string
	for _, u := range units {
		switch u.Key {
		case "subscription.billable_charge", "revenue.recovery_document", "treasury.collection_application":
			keys = append(keys, u.Key)
		}
	}
	return keys
}

// The option gates the three units: absent/false mounts none of them (other
// apps' route maps are unchanged), true mounts all three.
func TestWithRecoveryChargesGatesTheUnits(t *testing.T) {
	if got := recoveryUnitKeys(AllUnits(&UseCases{}, nil)); len(got) != 0 {
		t.Fatalf("default AllUnits mounts %v", got)
	}
	if got := recoveryUnitKeys(AllUnits(&UseCases{}, nil, WithRecoveryCharges(false))); len(got) != 0 {
		t.Fatalf("WithRecoveryCharges(false) mounts %v", got)
	}
	if got := recoveryUnitKeys(AllUnits(&UseCases{}, nil, WithRecoveryCharges(true))); len(got) != 3 {
		t.Fatalf("WithRecoveryCharges(true) mounts %v", got)
	}
}

// R4: a compose-v2 Unit fails closed in Mount — a missing closure is a boot
// error naming the field, never a silently dead page.
func TestRecoveryUnitsFailClosedInMountWhenClosuresAreMissing(t *testing.T) {
	want := map[string]string{
		"subscription.billable_charge":    "UseCases.BillableCharge.ListBillableCharges",
		"revenue.recovery_document":       "UseCases.RecoveryDocument.ReadRecoveryDocument",
		"treasury.collection_application": "UseCases.CollectionApplication.ReceiveAndApplyCollection",
	}
	seen := 0
	for _, u := range AllUnits(&UseCases{}, nil, WithRecoveryCharges(true)) {
		field, ok := want[u.Key]
		if !ok {
			continue
		}
		seen++
		err := u.Mount(&compose.MountContext{Routes: &recordingRegistrar{}})
		if err == nil || !strings.Contains(err.Error(), field) {
			t.Errorf("%s Mount(empty UseCases) error = %v, want it to name %s", u.Key, err, field)
		}
	}
	if seen != 3 {
		t.Fatalf("saw %d recovery units", seen)
	}
	// A nil aggregate fails closed too (never a nil dereference).
	for _, u := range RecoveryChargeUnits(nil, nil) {
		if err := u.Mount(&compose.MountContext{Routes: &recordingRegistrar{}}); err == nil {
			t.Errorf("%s Mount(nil UseCases) must fail", u.Key)
		}
	}
}

// R5 m10: without ListClients the receive-and-apply picker would render empty on
// a 200 page, so Mount fails closed naming the closure.
func TestReceiveApplyRequiresTheClientLister(t *testing.T) {
	uc := stubbedRecoveryUseCases()
	uc.Entity.Client.ListClients = nil
	for _, u := range RecoveryChargeUnits(uc, nil) {
		err := u.Mount(&compose.MountContext{Routes: &recordingRegistrar{}})
		if u.Key == "treasury.collection_application" {
			if err == nil || !strings.Contains(err.Error(), "UseCases.Entity.Client.ListClients") {
				t.Errorf("receive-apply Mount without ListClients = %v, want it to name ListClients", err)
			}
		}
	}
}

type recordingRegistrar struct{ gets, posts []string }

func (r *recordingRegistrar) GET(p string, _ view.View, _ ...string)  { r.gets = append(r.gets, p) }
func (r *recordingRegistrar) POST(p string, _ view.View, _ ...string) { r.posts = append(r.posts, p) }

type muxRegistrar struct{ mux *http.ServeMux }

func (m muxRegistrar) handle(method, path string, v view.View) {
	m.mux.HandleFunc(method+" "+path, func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) })
}
func (m muxRegistrar) GET(path string, v view.View, _ ...string)  { m.handle(http.MethodGet, path, v) }
func (m muxRegistrar) POST(path string, v view.View, _ ...string) { m.handle(http.MethodPost, path, v) }

func stubbedRecoveryUseCases() *UseCases {
	uc := &UseCases{}
	uc.BillableCharge.ListBillableCharges = func(context.Context, *billablechargepb.ListBillableChargesRequest) (*billablechargepb.ListBillableChargesResponse, error) {
		return &billablechargepb.ListBillableChargesResponse{}, nil
	}
	uc.BillableCharge.AdjustBillableCharge = func(context.Context, *billablechargepb.AdjustBillableChargeRequest) (*billablechargepb.AdjustBillableChargeResponse, error) {
		return &billablechargepb.AdjustBillableChargeResponse{}, nil
	}
	uc.RecoveryDocument.ListRecoveryDocuments = func(context.Context, *recoverydocumentpb.ListRecoveryDocumentsRequest) (*recoverydocumentpb.ListRecoveryDocumentsResponse, error) {
		return &recoverydocumentpb.ListRecoveryDocumentsResponse{}, nil
	}
	uc.RecoveryDocument.ReadRecoveryDocument = func(context.Context, *recoverydocumentpb.ReadRecoveryDocumentRequest) (*recoverydocumentpb.ReadRecoveryDocumentResponse, error) {
		return &recoverydocumentpb.ReadRecoveryDocumentResponse{}, nil
	}
	uc.RecoveryDocument.IssueRecoveryDocuments = func(context.Context, *recoverydocumentpb.IssueRecoveryDocumentsRequest) (*recoverydocumentpb.IssueRecoveryDocumentsResponse, error) {
		return &recoverydocumentpb.IssueRecoveryDocumentsResponse{}, nil
	}
	uc.RecoveryDocument.VoidRecoveryDocument = func(context.Context, *recoverydocumentpb.VoidRecoveryDocumentRequest) (*recoverydocumentpb.VoidRecoveryDocumentResponse, error) {
		return &recoverydocumentpb.VoidRecoveryDocumentResponse{}, nil
	}
	uc.DocumentSeries.ListDocumentSeries = func(context.Context, *documentseriespb.ListDocumentSeriesRequest) (*documentseriespb.ListDocumentSeriesResponse, error) {
		return &documentseriespb.ListDocumentSeriesResponse{}, nil
	}
	uc.CollectionApplication.ReceiveAndApplyCollection = func(context.Context, *collectionapplicationpb.ReceiveAndApplyCollectionRequest) (*collectionapplicationpb.ReceiveAndApplyCollectionResponse, error) {
		return &collectionapplicationpb.ReceiveAndApplyCollectionResponse{}, nil
	}
	uc.CollectionApplication.PreviewCollectionApplication = func(context.Context, *collectionapplicationpb.PreviewCollectionApplicationRequest) (*collectionapplicationpb.PreviewCollectionApplicationResponse, error) {
		return &collectionapplicationpb.PreviewCollectionApplicationResponse{}, nil
	}
	uc.CollectionApplication.ReverseCollectionApplication = func(context.Context, *collectionapplicationpb.ReverseCollectionApplicationRequest) (*collectionapplicationpb.ReverseCollectionApplicationResponse, error) {
		return &collectionapplicationpb.ReverseCollectionApplicationResponse{}, nil
	}
	uc.Entity.Client.ListClients = func(context.Context, *clientpb.ListClientsRequest) (*clientpb.ListClientsResponse, error) {
		return &clientpb.ListClientsResponse{}, nil
	}
	uc.CollectionApplication.ListCollectionApplications = func(context.Context, *collectionapplicationpb.ListCollectionApplicationsRequest) (*collectionapplicationpb.ListCollectionApplicationsResponse, error) {
		return &collectionapplicationpb.ListCollectionApplicationsResponse{}, nil
	}
	return uc
}

// Assembling the three units beside the collection unit through the real
// engine proves the route-key map has no collision, every nav route resolves,
// and the paths register on a live ServeMux (which panics on a conflict).
func TestRecoveryUnitsAssembleBesideCollectionWithoutConflict(t *testing.T) {
	uc := stubbedRecoveryUseCases()
	units := append(RecoveryChargeUnits(uc, nil), CollectionUnit(uc, &Infra{}))
	res, err := (&compose.Engine{}).Assemble(units, muxRegistrar{http.NewServeMux()})
	if err != nil {
		t.Fatalf("assemble: %v", err)
	}
	for _, key := range []string{
		"billable_charge.list", "billable_charge.issue", "billable_charge.adjust",
		"recovery_document.list", "recovery_document.detail", "recovery_document.void",
		"collection_application.receive_apply", "collection_application.preview", "collection_application.reverse",
		"collection.list",
	} {
		if res.RouteMap[key] == "" {
			t.Errorf("route %s missing from the assembled route map", key)
		}
	}
	// Sidebar entries are engine-resolved by unit + item key.
	for _, k := range [][2]string{{"subscription.billable_charge", "billable-charges"}, {"revenue.recovery_document", "recovery-documents"}} {
		if href, ok := res.ResolveNavItemHref(k[0], k[1]); !ok || href == "" {
			t.Errorf("nav %s/%s does not resolve", k[0], k[1])
		}
	}
}

// Without the option no recovery route exists (other apps' route maps unchanged).
func TestRouteMapWithoutTheOptionHasNoRecoveryRoutes(t *testing.T) {
	uc := stubbedRecoveryUseCases()
	res, err := (&compose.Engine{}).Assemble([]compose.Unit{CollectionUnit(uc, &Infra{})}, muxRegistrar{http.NewServeMux()})
	if err != nil {
		t.Fatalf("assemble: %v", err)
	}
	for k := range res.RouteMap {
		if strings.HasPrefix(k, "billable_charge.") || strings.HasPrefix(k, "recovery_document.") || strings.HasPrefix(k, "collection_application.") {
			t.Errorf("unexpected recovery route %s", k)
		}
	}
}
