package centymo

import (
	"context"
	"fmt"
	"strconv"

	pricelistpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/price_list"
	priceproductpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/price_product"
)

// Raw-DB adapter functions for price list CRUD.
// Used when Product use cases are unavailable (e.g. noop build tag).
// These convert between DataSource map[string]any records and proto types.

// PriceListDBListAdapter returns a proto-typed list function backed by raw DB.
func PriceListDBListAdapter(db DataSource) func(context.Context, *pricelistpb.ListPriceListsRequest) (*pricelistpb.ListPriceListsResponse, error) {
	return func(ctx context.Context, req *pricelistpb.ListPriceListsRequest) (*pricelistpb.ListPriceListsResponse, error) {
		records, err := db.ListSimple(ctx, "price_list")
		if err != nil {
			return nil, err
		}
		var items []*pricelistpb.PriceList
		for _, r := range records {
			items = append(items, mapToPriceList(r))
		}
		return &pricelistpb.ListPriceListsResponse{Data: items, Success: true}, nil
	}
}

// PriceListDBReadAdapter returns a proto-typed read function backed by raw DB.
func PriceListDBReadAdapter(db DataSource) func(context.Context, *pricelistpb.ReadPriceListRequest) (*pricelistpb.ReadPriceListResponse, error) {
	return func(ctx context.Context, req *pricelistpb.ReadPriceListRequest) (*pricelistpb.ReadPriceListResponse, error) {
		id := ""
		if req.Data != nil {
			id = req.Data.Id
		}
		record, err := db.Read(ctx, "price_list", id)
		if err != nil {
			return nil, err
		}
		return &pricelistpb.ReadPriceListResponse{
			Data:    []*pricelistpb.PriceList{mapToPriceList(record)},
			Success: true,
		}, nil
	}
}

// PriceListDBCreateAdapter returns a proto-typed create function backed by raw DB.
func PriceListDBCreateAdapter(db DataSource) func(context.Context, *pricelistpb.CreatePriceListRequest) (*pricelistpb.CreatePriceListResponse, error) {
	return func(ctx context.Context, req *pricelistpb.CreatePriceListRequest) (*pricelistpb.CreatePriceListResponse, error) {
		if req.Data == nil {
			return nil, fmt.Errorf("missing price list data")
		}
		data := priceListToMap(req.Data)
		record, err := db.Create(ctx, "price_list", data)
		if err != nil {
			return nil, err
		}
		return &pricelistpb.CreatePriceListResponse{
			Data:    []*pricelistpb.PriceList{mapToPriceList(record)},
			Success: true,
		}, nil
	}
}

// PriceListDBUpdateAdapter returns a proto-typed update function backed by raw DB.
func PriceListDBUpdateAdapter(db DataSource) func(context.Context, *pricelistpb.UpdatePriceListRequest) (*pricelistpb.UpdatePriceListResponse, error) {
	return func(ctx context.Context, req *pricelistpb.UpdatePriceListRequest) (*pricelistpb.UpdatePriceListResponse, error) {
		if req.Data == nil {
			return nil, fmt.Errorf("missing price list data")
		}
		data := priceListToMap(req.Data)
		record, err := db.Update(ctx, "price_list", req.Data.Id, data)
		if err != nil {
			return nil, err
		}
		return &pricelistpb.UpdatePriceListResponse{
			Data:    []*pricelistpb.PriceList{mapToPriceList(record)},
			Success: true,
		}, nil
	}
}

// PriceListDBDeleteAdapter returns a proto-typed delete function backed by raw DB.
func PriceListDBDeleteAdapter(db DataSource) func(context.Context, *pricelistpb.DeletePriceListRequest) (*pricelistpb.DeletePriceListResponse, error) {
	return func(ctx context.Context, req *pricelistpb.DeletePriceListRequest) (*pricelistpb.DeletePriceListResponse, error) {
		id := ""
		if req.Data != nil {
			id = req.Data.Id
		}
		err := db.Delete(ctx, "price_list", id)
		if err != nil {
			return nil, err
		}
		return &pricelistpb.DeletePriceListResponse{Success: true}, nil
	}
}

// PriceProductDBListAdapter returns a proto-typed list function for price products backed by raw DB.
func PriceProductDBListAdapter(db DataSource) func(context.Context, *priceproductpb.ListPriceProductsRequest) (*priceproductpb.ListPriceProductsResponse, error) {
	return func(ctx context.Context, req *priceproductpb.ListPriceProductsRequest) (*priceproductpb.ListPriceProductsResponse, error) {
		records, err := db.ListSimple(ctx, "price_product")
		if err != nil {
			// Table may not exist yet — return empty list
			return &priceproductpb.ListPriceProductsResponse{Data: nil, Success: true}, nil
		}
		var items []*priceproductpb.PriceProduct
		for _, r := range records {
			items = append(items, mapToPriceProduct(r))
		}
		return &priceproductpb.ListPriceProductsResponse{Data: items, Success: true}, nil
	}
}

// --- internal helpers ---

func mapToPriceList(m map[string]any) *pricelistpb.PriceList {
	pl := &pricelistpb.PriceList{
		Id:   mapStr(m, "id"),
		Name: mapStr(m, "name"),
	}
	if v := mapStr(m, "description"); v != "" {
		pl.Description = &v
	}
	if v := mapStr(m, "date_start_string"); v != "" {
		pl.DateStartString = v
	}
	if v := mapStr(m, "date_end_string"); v != "" {
		pl.DateEndString = &v
	}
	if v, ok := m["active"]; ok {
		pl.Active = mapBool(v)
	}
	return pl
}

func mapToPriceProduct(m map[string]any) *priceproductpb.PriceProduct {
	pp := &priceproductpb.PriceProduct{
		Id:        mapStr(m, "id"),
		ProductId: mapStr(m, "product_id"),
		Name:      mapStr(m, "name"),
		Currency:  mapStr(m, "currency"),
		Amount:    mapInt64(m, "amount"),
	}
	if v := mapStr(m, "description"); v != "" {
		pp.Description = &v
	}
	if v := mapStr(m, "price_list_id"); v != "" {
		pp.PriceListId = &v
	}
	if v, ok := m["active"]; ok {
		pp.Active = mapBool(v)
	}
	return pp
}

func priceListToMap(pl *pricelistpb.PriceList) map[string]any {
	m := map[string]any{
		"name": pl.Name,
	}
	if pl.Description != nil {
		m["description"] = *pl.Description
	}
	if pl.DateStartString != "" {
		m["date_start_string"] = pl.DateStartString
	}
	if pl.DateEndString != nil {
		m["date_end_string"] = *pl.DateEndString
	}
	m["active"] = pl.Active
	return m
}

func mapStr(m map[string]any, key string) string {
	if v, ok := m[key]; ok && v != nil {
		return fmt.Sprintf("%v", v)
	}
	return ""
}

func mapInt64(m map[string]any, key string) int64 {
	v, ok := m[key]
	if !ok || v == nil {
		return 0
	}
	switch n := v.(type) {
	case int64:
		return n
	case float64:
		return int64(n)
	case int:
		return int64(n)
	case string:
		i, _ := strconv.ParseInt(n, 10, 64)
		return i
	default:
		return 0
	}
}

func mapBool(v any) bool {
	switch b := v.(type) {
	case bool:
		return b
	case string:
		return b == "true" || b == "t"
	default:
		return false
	}
}
