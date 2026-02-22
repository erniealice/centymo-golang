package centymo

import (
	"context"
	"fmt"
	"strconv"

	productpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product"
)

// Raw-DB adapter functions for product CRUD.
// Used when Product use cases are unavailable (e.g. noop build tag).
// These convert between DataSource map[string]any records and proto types.

// ProductDBListAdapter returns a proto-typed list function backed by raw DB.
func ProductDBListAdapter(db DataSource) func(context.Context, *productpb.ListProductsRequest) (*productpb.ListProductsResponse, error) {
	return func(ctx context.Context, req *productpb.ListProductsRequest) (*productpb.ListProductsResponse, error) {
		records, err := db.ListSimple(ctx, "product")
		if err != nil {
			return nil, err
		}
		var items []*productpb.Product
		for _, r := range records {
			items = append(items, mapToProduct(r))
		}
		return &productpb.ListProductsResponse{Data: items, Success: true}, nil
	}
}

// ProductDBReadAdapter returns a proto-typed read function backed by raw DB.
func ProductDBReadAdapter(db DataSource) func(context.Context, *productpb.ReadProductRequest) (*productpb.ReadProductResponse, error) {
	return func(ctx context.Context, req *productpb.ReadProductRequest) (*productpb.ReadProductResponse, error) {
		id := ""
		if req.Data != nil {
			id = req.Data.Id
		}
		record, err := db.Read(ctx, "product", id)
		if err != nil {
			return nil, err
		}
		return &productpb.ReadProductResponse{
			Data:    []*productpb.Product{mapToProduct(record)},
			Success: true,
		}, nil
	}
}

// ProductDBCreateAdapter returns a proto-typed create function backed by raw DB.
func ProductDBCreateAdapter(db DataSource) func(context.Context, *productpb.CreateProductRequest) (*productpb.CreateProductResponse, error) {
	return func(ctx context.Context, req *productpb.CreateProductRequest) (*productpb.CreateProductResponse, error) {
		if req.Data == nil {
			return nil, fmt.Errorf("missing product data")
		}
		data := productToMap(req.Data)
		record, err := db.Create(ctx, "product", data)
		if err != nil {
			return nil, err
		}
		return &productpb.CreateProductResponse{
			Data:    []*productpb.Product{mapToProduct(record)},
			Success: true,
		}, nil
	}
}

// ProductDBUpdateAdapter returns a proto-typed update function backed by raw DB.
func ProductDBUpdateAdapter(db DataSource) func(context.Context, *productpb.UpdateProductRequest) (*productpb.UpdateProductResponse, error) {
	return func(ctx context.Context, req *productpb.UpdateProductRequest) (*productpb.UpdateProductResponse, error) {
		if req.Data == nil {
			return nil, fmt.Errorf("missing product data")
		}
		data := productToMap(req.Data)
		record, err := db.Update(ctx, "product", req.Data.Id, data)
		if err != nil {
			return nil, err
		}
		return &productpb.UpdateProductResponse{
			Data:    []*productpb.Product{mapToProduct(record)},
			Success: true,
		}, nil
	}
}

// ProductDBDeleteAdapter returns a proto-typed delete function backed by raw DB.
func ProductDBDeleteAdapter(db DataSource) func(context.Context, *productpb.DeleteProductRequest) (*productpb.DeleteProductResponse, error) {
	return func(ctx context.Context, req *productpb.DeleteProductRequest) (*productpb.DeleteProductResponse, error) {
		id := ""
		if req.Data != nil {
			id = req.Data.Id
		}
		err := db.Delete(ctx, "product", id)
		if err != nil {
			return nil, err
		}
		return &productpb.DeleteProductResponse{Success: true}, nil
	}
}

// --- internal helpers ---

func mapToProduct(m map[string]any) *productpb.Product {
	p := &productpb.Product{
		Id:       mapStr(m, "id"),
		Name:     mapStr(m, "name"),
		Price:    mapFloat64(m, "price"),
		Currency: mapStr(m, "currency"),
	}
	if v := mapStr(m, "description"); v != "" {
		p.Description = &v
	}
	if v, ok := m["active"]; ok {
		p.Active = mapBool(v)
	}
	return p
}

func productToMap(p *productpb.Product) map[string]any {
	m := map[string]any{
		"name":     p.Name,
		"price":    p.Price,
		"currency": p.Currency,
		"active":   p.Active,
	}
	if p.Description != nil {
		m["description"] = *p.Description
	}
	return m
}

func mapFloat64(m map[string]any, key string) float64 {
	v, ok := m[key]
	if !ok || v == nil {
		return 0
	}
	switch n := v.(type) {
	case float64:
		return n
	case float32:
		return float64(n)
	case int64:
		return float64(n)
	case int:
		return float64(n)
	case string:
		f, _ := strconv.ParseFloat(n, 64)
		return f
	default:
		return 0
	}
}
