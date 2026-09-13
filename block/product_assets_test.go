package block

import (
	"context"
	"errors"
	assetpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/asset/asset"
	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	"testing"
)

func TestProductAssetListerUnassignedAndPagination(t *testing.T) {
	ctx := context.WithValue(context.Background(), struct{}{}, "scope")
	product := "product"
	other := "other"
	calls := 0
	list := productAssetLister(func(got context.Context, r *assetpb.GetAssetListPageDataRequest) (*assetpb.GetAssetListPageDataResponse, error) {
		if got != ctx {
			t.Fatal("lost scope")
		}
		calls++
		active := r.GetFilters().GetFilters()[0].GetBooleanFilter().GetValue()
		if !active {
			return &assetpb.GetAssetListPageDataResponse{AssetList: []*assetpb.Asset{{Id: "inactive"}}}, nil
		}
		if r.GetPagination().GetOffset().GetPage() == 1 {
			return &assetpb.GetAssetListPageDataResponse{AssetList: []*assetpb.Asset{{Id: "assigned", ProductId: &product}, {Id: "other", ProductId: &other}}, Pagination: &commonpb.PaginationResponse{HasNext: true}}, nil
		}
		return &assetpb.GetAssetListPageDataResponse{AssetList: []*assetpb.Asset{{Id: "free"}}}, nil
	})
	free, err := list(ctx, product, true)
	if err != nil || len(free) != 2 || calls != 3 {
		t.Fatalf("free=%v calls=%d err=%v", free, calls, err)
	}
	for _, a := range free {
		if a.GetProductId() != "" {
			t.Fatal("offered assigned asset")
		}
	}
	linked, err := list(ctx, product, false)
	if err != nil || len(linked) != 1 || linked[0].GetId() != "assigned" {
		t.Fatalf("linked=%v err=%v", linked, err)
	}
	failure := errors.New("lookup failed")
	_, err = productAssetLister(func(context.Context, *assetpb.GetAssetListPageDataRequest) (*assetpb.GetAssetListPageDataResponse, error) {
		return nil, failure
	})(ctx, product, true)
	if !errors.Is(err, failure) {
		t.Fatal("lookup failure hidden")
	}
}
