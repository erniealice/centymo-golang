package block

import (
	"context"
	"sort"
	"strings"

	"github.com/erniealice/espyna-golang/consumer"
	assetpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/asset/asset"
	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
)

func bindProductAssets(dst *UseCases, uc *consumer.UseCases) {
	if uc == nil || uc.Asset == nil || uc.Asset.Asset == nil {
		return
	}
	ops := uc.Asset.Asset
	if ops.GetAssetListPageData == nil || ops.AssignProduct == nil {
		return
	}
	list := productAssetLister(ops.GetAssetListPageData.Execute)
	dst.Product.ListProductAssets = func(ctx context.Context, productID string) ([]*assetpb.Asset, error) {
		return list(ctx, productID, false)
	}
	dst.Product.ListAssignableAssets = func(ctx context.Context, productID string) ([]*assetpb.Asset, error) {
		return list(ctx, productID, true)
	}
	dst.Product.AssignProductAsset = func(ctx context.Context, productID, assetID string) error {
		return ops.AssignProduct.Execute(ctx, assetID, productID)
	}
}

// Read every page in both activation states; assignment is determined by the
// product link, not the asset's activation flag. Scoped use cases retain ctx.
func productAssetLister(load func(context.Context, *assetpb.GetAssetListPageDataRequest) (*assetpb.GetAssetListPageDataResponse, error)) func(context.Context, string, bool) ([]*assetpb.Asset, error) {
	return func(ctx context.Context, productID string, unassigned bool) ([]*assetpb.Asset, error) {
		var result []*assetpb.Asset
		seen := map[string]bool{}
		for _, active := range []bool{true, false} {
			for page := int32(1); ; page++ {
				resp, err := load(ctx, &assetpb.GetAssetListPageDataRequest{
					Filters:    &commonpb.FilterRequest{Filters: []*commonpb.TypedFilter{{Field: "active", FilterType: &commonpb.TypedFilter_BooleanFilter{BooleanFilter: &commonpb.BooleanFilter{Value: active}}}}},
					Pagination: &commonpb.PaginationRequest{Limit: 100, Method: &commonpb.PaginationRequest_Offset{Offset: &commonpb.OffsetPagination{Page: page}}},
				})
				if err != nil {
					return nil, err
				}
				for _, a := range resp.GetAssetList() {
					if a == nil || seen[a.GetId()] {
						continue
					}
					seen[a.GetId()] = true
					if (unassigned && a.GetProductId() == "") || (!unassigned && a.GetProductId() == productID) {
						result = append(result, a)
					}
				}
				if !resp.GetPagination().GetHasNext() {
					break
				}
			}
		}
		sort.SliceStable(result, func(i, j int) bool {
			return strings.ToLower(result[i].GetName()) < strings.ToLower(result[j].GetName())
		})
		return result, nil
	}
}
