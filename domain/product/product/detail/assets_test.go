package detail

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	product "github.com/erniealice/centymo-golang/domain/product/product"
	assetpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/asset/asset"
	productpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

func assetTestDeps() *DetailViewDeps {
	return &DetailViewDeps{
		Routes: product.Routes{TabActionURL: "/products/{id}/tab/{tab}"},
		Labels: product.Labels{
			AssetAssignment: product.AssetAssignmentLabels{
				Assets: "Assets", Add: "Assign Existing Asset", SelectAsset: "Select asset",
				Required: "Asset is required", Unavailable: "Asset assignment unavailable",
				Failed: "Asset operation failed",
				Name:   "Name", Number: "Number", Location: "Location",
			},
			Errors: product.ErrorLabels{PermissionDenied: "Permission denied", InvalidFormData: "Invalid form data"},
		},
	}
}

func assetCtx(perms ...string) context.Context {
	return view.WithUserPermissions(context.Background(), types.NewUserPermissions(perms))
}

func assetRequest(method string, form url.Values) *view.ViewContext {
	var body *strings.Reader
	if form == nil {
		body = strings.NewReader("")
	} else {
		body = strings.NewReader(form.Encode())
	}
	req := httptest.NewRequest(method, "/products/p1/tab/assets?mode=add", body)
	if form != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	req.SetPathValue("id", "p1")
	req.SetPathValue("tab", "assets")
	return &view.ViewContext{Request: req}
}

func TestBuildTabItems_HidesAssetsWhenCallbacksMissing(t *testing.T) {
	items := buildTabItems("p1", product.Labels{}, 0, 0, 0, 0, product.Routes{DetailURL: "/products/{id}", TabActionURL: "/products/{id}/tab/{tab}"}, false, false)
	for _, item := range items {
		if item.Key == "assets" {
			t.Fatal("assets tab should be hidden when callbacks are missing")
		}
	}
}

func TestBuildTabItems_AssetsImmediatelyAfterLines(t *testing.T) {
	items := buildTabItems("p1", product.Labels{AssetAssignment: product.AssetAssignmentLabels{Assets: "Assets"}}, 0, 0, 0, 2, product.Routes{DetailURL: "/products/{id}", TabActionURL: "/products/{id}/tab/{tab}"}, false, true)
	var lines, assets int
	for i, item := range items {
		if item.Key == "lines" {
			lines = i
		}
		if item.Key == "assets" {
			assets = i
		}
	}
	if assets != lines+1 {
		t.Fatalf("assets tab index = %d, lines index = %d; want immediately after", assets, lines)
	}
}

func TestNewTabAction_AssetReadDenied(t *testing.T) {
	deps := assetTestDeps()
	deps.ListProductAssets = func(context.Context, string) ([]*assetpb.Asset, error) { return nil, nil }
	deps.ListAssignableAssets = deps.ListProductAssets
	deps.AssignProductAsset = func(context.Context, string, string) error { return nil }
	result := NewTabAction(deps).Handle(assetCtx("product:read"), assetRequest(http.MethodGet, nil))
	if result.StatusCode != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", result.StatusCode, http.StatusForbidden)
	}
}

func TestHandleProductAssetAdd_AssetUpdateDeniedOnGETAndPOST(t *testing.T) {
	deps := assetTestDeps()
	deps.ListProductAssets = func(context.Context, string) ([]*assetpb.Asset, error) { return nil, nil }
	deps.ListAssignableAssets = func(context.Context, string) ([]*assetpb.Asset, error) { return []*assetpb.Asset{{Id: "a1"}}, nil }
	deps.AssignProductAsset = func(context.Context, string, string) error { return nil }
	ctx := assetCtx("product:read", "asset:read")
	for _, method := range []string{http.MethodGet, http.MethodPost} {
		result := handleProductAssetAdd(ctx, deps, assetRequest(method, url.Values{"asset_id": {"a1"}}), "p1")
		if result.StatusCode != http.StatusUnprocessableEntity {
			t.Errorf("%s status = %d, want %d", method, result.StatusCode, http.StatusUnprocessableEntity)
		}
	}
}

func TestHandleProductAssetAdd_GETValidatesProductBeforeListing(t *testing.T) {
	deps := assetTestDeps()
	deps.ListProductAssets = func(context.Context, string) ([]*assetpb.Asset, error) { return nil, nil }
	deps.ReadProduct = func(context.Context, *productpb.ReadProductRequest) (*productpb.ReadProductResponse, error) {
		return nil, errors.New("product read failed")
	}
	deps.ListAssignableAssets = func(context.Context, string) ([]*assetpb.Asset, error) {
		t.Fatal("assignable assets should not load after product validation fails")
		return nil, nil
	}
	deps.AssignProductAsset = func(context.Context, string, string) error { return nil }
	result := handleProductAssetAdd(assetCtx("product:read", "asset:read", "asset:update"), deps, assetRequest(http.MethodGet, nil), "p1")
	if result.Headers["HX-Error-Message"] != "Asset operation failed" {
		t.Fatalf("HX error = %q, want translated failure", result.Headers["HX-Error-Message"])
	}
}

func TestHandleProductAssetAdd_GETIncludesValidatedProductAndPathOnlyAction(t *testing.T) {
	deps := assetTestDeps()
	deps.ReadProduct = func(context.Context, *productpb.ReadProductRequest) (*productpb.ReadProductResponse, error) {
		return &productpb.ReadProductResponse{Data: []*productpb.Product{{Id: "p1", Name: "Sample Product"}}}, nil
	}
	deps.ListProductAssets = func(context.Context, string) ([]*assetpb.Asset, error) { return nil, nil }
	deps.ListAssignableAssets = func(context.Context, string) ([]*assetpb.Asset, error) { return nil, nil }
	deps.AssignProductAsset = func(context.Context, string, string) error { return nil }
	result := handleProductAssetAdd(assetCtx("product:read", "asset:read", "asset:update"), deps, assetRequest(http.MethodGet, nil), "p1")
	if result.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", result.StatusCode, http.StatusOK)
	}
	data, ok := result.Data.(*ProductAssetFormData)
	if !ok {
		t.Fatalf("data type = %T, want *ProductAssetFormData", result.Data)
	}
	if data.ProductName != "Sample Product" {
		t.Fatalf("product name = %q, want %q", data.ProductName, "Sample Product")
	}
	if strings.Contains(data.FormAction, "?") {
		t.Fatalf("form action = %q, want path-only action", data.FormAction)
	}
}

func TestHandleProductAssetAdd_RejectsBlankSelection(t *testing.T) {
	deps := assetTestDeps()
	deps.ListAssignableAssets = func(context.Context, string) ([]*assetpb.Asset, error) { return nil, nil }
	deps.AssignProductAsset = func(context.Context, string, string) error { t.Fatal("assign should not be called"); return nil }
	result := handleProductAssetAdd(assetCtx("product:read", "asset:read", "asset:update"), deps, assetRequest(http.MethodPost, url.Values{}), "p1")
	if result.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want %d", result.StatusCode, http.StatusUnprocessableEntity)
	}
}

func TestHandleProductAssetAdd_POSTForwardsProductAndAsset(t *testing.T) {
	deps := assetTestDeps()
	deps.ListProductAssets = func(context.Context, string) ([]*assetpb.Asset, error) { return nil, nil }
	deps.ListAssignableAssets = func(context.Context, string) ([]*assetpb.Asset, error) { return nil, nil }
	var gotProduct, gotAsset string
	deps.AssignProductAsset = func(_ context.Context, productID, assetID string) error {
		gotProduct, gotAsset = productID, assetID
		return nil
	}
	result := handleProductAssetAdd(assetCtx("product:read", "asset:read", "asset:update"), deps, assetRequest(http.MethodPost, url.Values{"asset_id": {"a1"}}), "p1")
	if result.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", result.StatusCode, http.StatusOK)
	}
	if gotProduct != "p1" || gotAsset != "a1" {
		t.Fatalf("forwarded (%q, %q), want (%q, %q)", gotProduct, gotAsset, "p1", "a1")
	}
}

func TestBuildAssetsTable_PropagatesLoaderError(t *testing.T) {
	want := errors.New("asset loader failed")
	deps := assetTestDeps()
	deps.ListProductAssets = func(context.Context, string) ([]*assetpb.Asset, error) { return nil, want }
	_, _, err := buildAssetsTable(assetCtx("asset:read"), deps, "p1")
	if err == nil || err.Error() != "Asset operation failed" {
		t.Fatalf("error = %v, want translated loader failure", err)
	}
}
