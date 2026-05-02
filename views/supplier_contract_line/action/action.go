package action

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"

	centymo "github.com/erniealice/centymo-golang"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	productpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product"
	suppliercontractlinepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/supplier_contract_line"

	"github.com/erniealice/centymo-golang/views/supplier_contract_line/form"
)

// Deps holds dependencies for line item action handlers.
type Deps struct {
	Routes                     centymo.SupplierContractRoutes
	Labels                     centymo.SupplierContractLabels
	CommonLabels               pyeza.CommonLabels
	CreateSupplierContractLine func(ctx context.Context, req *suppliercontractlinepb.CreateSupplierContractLineRequest) (*suppliercontractlinepb.CreateSupplierContractLineResponse, error)
	ReadSupplierContractLine   func(ctx context.Context, req *suppliercontractlinepb.ReadSupplierContractLineRequest) (*suppliercontractlinepb.ReadSupplierContractLineResponse, error)
	UpdateSupplierContractLine func(ctx context.Context, req *suppliercontractlinepb.UpdateSupplierContractLineRequest) (*suppliercontractlinepb.UpdateSupplierContractLineResponse, error)
	DeleteSupplierContractLine func(ctx context.Context, req *suppliercontractlinepb.DeleteSupplierContractLineRequest) (*suppliercontractlinepb.DeleteSupplierContractLineResponse, error)
	ListProducts               func(ctx context.Context, req *productpb.ListProductsRequest) (*productpb.ListProductsResponse, error)
}

// NewAddAction handles GET+POST for adding a supplier contract line.
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		contractID := viewCtx.Request.PathValue("id")
		if contractID == "" {
			return view.Error(fmt.Errorf("missing contract id"))
		}

		l := deps.Labels
		if viewCtx.Request.Method == http.MethodGet {
			fd := buildEmptyLineFormData(ctx, deps, l)
			fd.FormAction = viewCtx.Request.URL.Path
			fd.SupplierContractID = contractID
			return view.OK("supplier-contract-line-drawer-form", fd)
		}

		// POST
		r := viewCtx.Request
		if err := r.ParseForm(); err != nil {
			return view.Error(fmt.Errorf("parse form: %w", err))
		}

		qty, _ := strconv.ParseFloat(r.FormValue("quantity"), 64)
		lineNum, _ := strconv.ParseInt(r.FormValue("line_number"), 10, 32)
		unitPrice := parseCentavos(r.FormValue("unit_price"))

		req := &suppliercontractlinepb.CreateSupplierContractLineRequest{
			Data: &suppliercontractlinepb.SupplierContractLine{
				SupplierContractId: contractID,
				Description:        r.FormValue("description"),
				LineType:           r.FormValue("line_type"),
				ProductId:          optionalString(r.FormValue("product_id")),
				Quantity:           qty,
				UnitPrice:          unitPrice,
				StartDate:          optionalString(r.FormValue("start_date")),
				EndDate:            optionalString(r.FormValue("end_date")),
				ExpenseAccountId:   optionalString(r.FormValue("expense_account_id")),
				LineNumber:         int32(lineNum),
				Active:             true,
			},
		}

		_, err := deps.CreateSupplierContractLine(ctx, req)
		if err != nil {
			log.Printf("CreateSupplierContractLine: %v", err)
			return view.Error(fmt.Errorf("failed to create contract line: %w", err))
		}

		return centymo.HTMXSuccess("supplier-contracts-table")
	})
}

// NewEditAction handles GET+POST for editing a supplier contract line.
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		lineID := viewCtx.Request.PathValue("lid")
		contractID := viewCtx.Request.PathValue("id")
		if lineID == "" || contractID == "" {
			return view.Error(fmt.Errorf("missing id or lid"))
		}

		l := deps.Labels
		if viewCtx.Request.Method == http.MethodGet {
			resp, err := deps.ReadSupplierContractLine(ctx, &suppliercontractlinepb.ReadSupplierContractLineRequest{
				Data: &suppliercontractlinepb.SupplierContractLine{Id: lineID},
			})
			if err != nil {
				return view.Error(fmt.Errorf("failed to read contract line: %w", err))
			}
			data := resp.GetData()
			if len(data) == 0 {
				return view.Error(fmt.Errorf("contract line not found"))
			}
			line := data[0]

			fd := buildEmptyLineFormData(ctx, deps, l)
			fd.FormAction = viewCtx.Request.URL.Path
			fd.IsEdit = true
			fd.ID = lineID
			fd.SupplierContractID = contractID
			fd.Description = line.GetDescription()
			fd.LineType = line.GetLineType()
			fd.ProductID = line.GetProductId()
			fd.Quantity = strconv.FormatFloat(line.GetQuantity(), 'f', 2, 64)
			fd.UnitPrice = formatCentavos(line.GetUnitPrice())
			fd.Treatment = line.GetTreatment().String()
			fd.StartDate = line.GetStartDate()
			fd.EndDate = line.GetEndDate()
			fd.ExpenseAccountID = line.GetExpenseAccountId()
			fd.LineNumber = strconv.Itoa(int(line.GetLineNumber()))
			return view.OK("supplier-contract-line-drawer-form", fd)
		}

		// POST
		r := viewCtx.Request
		if err := r.ParseForm(); err != nil {
			return view.Error(fmt.Errorf("parse form: %w", err))
		}

		qty, _ := strconv.ParseFloat(r.FormValue("quantity"), 64)
		lineNum, _ := strconv.ParseInt(r.FormValue("line_number"), 10, 32)
		unitPrice := parseCentavos(r.FormValue("unit_price"))

		req := &suppliercontractlinepb.UpdateSupplierContractLineRequest{
			Data: &suppliercontractlinepb.SupplierContractLine{
				Id:                 lineID,
				SupplierContractId: contractID,
				Description:        r.FormValue("description"),
				LineType:           r.FormValue("line_type"),
				ProductId:          optionalString(r.FormValue("product_id")),
				Quantity:           qty,
				UnitPrice:          unitPrice,
				StartDate:          optionalString(r.FormValue("start_date")),
				EndDate:            optionalString(r.FormValue("end_date")),
				ExpenseAccountId:   optionalString(r.FormValue("expense_account_id")),
				LineNumber:         int32(lineNum),
			},
		}

		_, err := deps.UpdateSupplierContractLine(ctx, req)
		if err != nil {
			log.Printf("UpdateSupplierContractLine %s: %v", lineID, err)
			return view.Error(fmt.Errorf("failed to update contract line: %w", err))
		}

		return centymo.HTMXSuccess("supplier-contracts-table")
	})
}

// NewDeleteAction handles POST /action/supplier-contract/{id}/lines/delete.
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		if viewCtx.Request.Method != http.MethodPost {
			return view.Error(fmt.Errorf("method not allowed"))
		}
		lineID := viewCtx.Request.FormValue("id")
		if lineID == "" {
			return view.Error(fmt.Errorf("missing id"))
		}
		_, err := deps.DeleteSupplierContractLine(ctx, &suppliercontractlinepb.DeleteSupplierContractLineRequest{
			Data: &suppliercontractlinepb.SupplierContractLine{Id: lineID},
		})
		if err != nil {
			log.Printf("DeleteSupplierContractLine %s: %v", lineID, err)
			return view.Error(fmt.Errorf("failed to delete contract line: %w", err))
		}
		return centymo.HTMXSuccess("supplier-contracts-table")
	})
}

// --- helpers -----------------------------------------------------------------

func buildEmptyLineFormData(ctx context.Context, deps *Deps, l centymo.SupplierContractLabels) *form.Data {
	fd := &form.Data{
		Labels:       l,
		CommonLabels: deps.CommonLabels,
		TreatmentOptions: []form.TreatmentOption{
			{Value: "SUPPLIER_CONTRACT_LINE_TREATMENT_RECURRING", Label: l.Lines.TreatmentRecurring},
			{Value: "SUPPLIER_CONTRACT_LINE_TREATMENT_ONE_TIME", Label: l.Lines.TreatmentOneTime},
			{Value: "SUPPLIER_CONTRACT_LINE_TREATMENT_USAGE_BASED", Label: l.Lines.TreatmentUsageBased},
			{Value: "SUPPLIER_CONTRACT_LINE_TREATMENT_MINIMUM_COMMITMENT", Label: l.Lines.TreatmentMinimumCommitment},
		},
	}

	// Load product options
	if deps.ListProducts != nil {
		resp, err := deps.ListProducts(ctx, &productpb.ListProductsRequest{})
		if err == nil {
			for _, p := range resp.GetData() {
				fd.Products = append(fd.Products, types.SelectOption{
					Value: p.GetId(),
					Label: p.GetName(),
				})
			}
		}
	}

	return fd
}

func parseCentavos(s string) int64 {
	if s == "" {
		return 0
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return int64(math.Round(f * 100))
}

func formatCentavos(v int64) string {
	if v == 0 {
		return ""
	}
	return strconv.FormatFloat(float64(v)/100.0, 'f', 2, 64)
}

func optionalString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func optionalStrVal(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}
