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

	suppliercontractlinepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/supplier_contract_line"
	scpslpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/supplier_contract_price_schedule_line"
)

// LineFormData is the template data for the SCPSL drawer form.
type LineFormData struct {
	FormAction                      string
	IsEdit                          bool
	ID                              string
	SupplierContractPriceScheduleID string

	// Fields
	SupplierContractLineID string
	Currency               string
	UnitPrice              string
	MinimumAmount          string
	Quantity               string
	CycleValueOverride     string
	CycleUnitOverride      string
	Notes                  string

	// Options
	ContractLines []types.SelectOption

	Labels       centymo.SupplierContractPriceScheduleLineFormLabels
	NounLabels   centymo.SupplierContractPriceScheduleLinesLabels
	CommonLabels pyeza.CommonLabels
}

// Deps holds dependencies for SCPSL action handlers.
type Deps struct {
	Routes       centymo.SupplierContractPriceScheduleRoutes
	Labels       centymo.SupplierContractPriceScheduleLabels
	CommonLabels pyeza.CommonLabels

	CreateSupplierContractPriceScheduleLine func(ctx context.Context, req *scpslpb.CreateSupplierContractPriceScheduleLineRequest) (*scpslpb.CreateSupplierContractPriceScheduleLineResponse, error)
	ReadSupplierContractPriceScheduleLine   func(ctx context.Context, req *scpslpb.ReadSupplierContractPriceScheduleLineRequest) (*scpslpb.ReadSupplierContractPriceScheduleLineResponse, error)
	UpdateSupplierContractPriceScheduleLine func(ctx context.Context, req *scpslpb.UpdateSupplierContractPriceScheduleLineRequest) (*scpslpb.UpdateSupplierContractPriceScheduleLineResponse, error)
	DeleteSupplierContractPriceScheduleLine func(ctx context.Context, req *scpslpb.DeleteSupplierContractPriceScheduleLineRequest) (*scpslpb.DeleteSupplierContractPriceScheduleLineResponse, error)

	// Optional — for the contract-line picker in the line drawer form.
	ListSupplierContractLines func(ctx context.Context, req *suppliercontractlinepb.ListSupplierContractLinesRequest) (*suppliercontractlinepb.ListSupplierContractLinesResponse, error)
}

// NewAddAction handles GET+POST for adding a SCPSL.
// URL pattern: /action/supplier-contract-price-schedule/{id}/lines/add  ({id}=schedule id)
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		scheduleID := viewCtx.Request.PathValue("id")
		if scheduleID == "" {
			return view.Error(fmt.Errorf("missing schedule id"))
		}

		l := deps.Labels
		if viewCtx.Request.Method == http.MethodGet {
			fd := buildEmptyLineFormData(ctx, deps, l)
			fd.FormAction = viewCtx.Request.URL.Path
			fd.SupplierContractPriceScheduleID = scheduleID
			return view.OK("supplier-contract-price-schedule-line-drawer-form", fd)
		}

		// POST
		r := viewCtx.Request
		if err := r.ParseForm(); err != nil {
			return view.Error(fmt.Errorf("parse form: %w", err))
		}

		unitPrice := parseCentavos(r.FormValue("unit_price"))
		var minPtr *int64
		if v := r.FormValue("minimum_amount"); v != "" {
			m := parseCentavos(v)
			minPtr = &m
		}
		var qtyPtr *float64
		if v := r.FormValue("quantity"); v != "" {
			if q, err := strconv.ParseFloat(v, 64); err == nil {
				qtyPtr = &q
			}
		}
		var cycValPtr *int32
		if v := r.FormValue("cycle_value_override"); v != "" {
			if cv, err := strconv.ParseInt(v, 10, 32); err == nil {
				ci := int32(cv)
				cycValPtr = &ci
			}
		}
		var cycUnitPtr *string
		if v := r.FormValue("cycle_unit_override"); v != "" {
			cu := v
			cycUnitPtr = &cu
		}
		var notesPtr *string
		if v := r.FormValue("notes"); v != "" {
			n := v
			notesPtr = &n
		}

		req := &scpslpb.CreateSupplierContractPriceScheduleLineRequest{
			Data: &scpslpb.SupplierContractPriceScheduleLine{
				SupplierContractPriceScheduleId: scheduleID,
				SupplierContractLineId:          r.FormValue("supplier_contract_line_id"),
				Currency:                        r.FormValue("currency"),
				UnitPrice:                       unitPrice,
				MinimumAmount:                   minPtr,
				Quantity:                        qtyPtr,
				CycleValueOverride:              cycValPtr,
				CycleUnitOverride:               cycUnitPtr,
				Notes:                           notesPtr,
				Active:                          true,
			},
		}

		_, err := deps.CreateSupplierContractPriceScheduleLine(ctx, req)
		if err != nil {
			log.Printf("CreateSupplierContractPriceScheduleLine: %v", err)
			return centymo.HTMXError(err.Error())
		}
		return centymo.HTMXSuccess("schedule-lines-table")
	})
}

// NewEditAction handles GET+POST for editing a SCPSL.
// URL pattern: /action/supplier-contract-price-schedule/{id}/lines/edit/{lid}
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		lineID := viewCtx.Request.PathValue("lid")
		scheduleID := viewCtx.Request.PathValue("id")
		if lineID == "" || scheduleID == "" {
			return view.Error(fmt.Errorf("missing id or lid"))
		}

		l := deps.Labels
		if viewCtx.Request.Method == http.MethodGet {
			resp, err := deps.ReadSupplierContractPriceScheduleLine(ctx, &scpslpb.ReadSupplierContractPriceScheduleLineRequest{
				Data: &scpslpb.SupplierContractPriceScheduleLine{Id: lineID},
			})
			if err != nil {
				return view.Error(fmt.Errorf("failed to read schedule line: %w", err))
			}
			data := resp.GetData()
			if len(data) == 0 {
				return view.Error(fmt.Errorf("schedule line not found"))
			}
			line := data[0]

			fd := buildEmptyLineFormData(ctx, deps, l)
			fd.FormAction = viewCtx.Request.URL.Path
			fd.IsEdit = true
			fd.ID = lineID
			fd.SupplierContractPriceScheduleID = scheduleID
			fd.SupplierContractLineID = line.GetSupplierContractLineId()
			fd.Currency = line.GetCurrency()
			fd.UnitPrice = formatCentavos(line.GetUnitPrice())
			if line.MinimumAmount != nil {
				fd.MinimumAmount = formatCentavos(line.GetMinimumAmount())
			}
			if line.Quantity != nil {
				fd.Quantity = strconv.FormatFloat(line.GetQuantity(), 'f', 2, 64)
			}
			if line.CycleValueOverride != nil {
				fd.CycleValueOverride = strconv.FormatInt(int64(line.GetCycleValueOverride()), 10)
			}
			if line.CycleUnitOverride != nil {
				fd.CycleUnitOverride = line.GetCycleUnitOverride()
			}
			fd.Notes = line.GetNotes()
			return view.OK("supplier-contract-price-schedule-line-drawer-form", fd)
		}

		// POST
		r := viewCtx.Request
		if err := r.ParseForm(); err != nil {
			return view.Error(fmt.Errorf("parse form: %w", err))
		}

		unitPrice := parseCentavos(r.FormValue("unit_price"))
		var minPtr *int64
		if v := r.FormValue("minimum_amount"); v != "" {
			m := parseCentavos(v)
			minPtr = &m
		}
		var qtyPtr *float64
		if v := r.FormValue("quantity"); v != "" {
			if q, err := strconv.ParseFloat(v, 64); err == nil {
				qtyPtr = &q
			}
		}
		var cycValPtr *int32
		if v := r.FormValue("cycle_value_override"); v != "" {
			if cv, err := strconv.ParseInt(v, 10, 32); err == nil {
				ci := int32(cv)
				cycValPtr = &ci
			}
		}
		var cycUnitPtr *string
		if v := r.FormValue("cycle_unit_override"); v != "" {
			cu := v
			cycUnitPtr = &cu
		}
		var notesPtr *string
		if v := r.FormValue("notes"); v != "" {
			n := v
			notesPtr = &n
		}

		req := &scpslpb.UpdateSupplierContractPriceScheduleLineRequest{
			Data: &scpslpb.SupplierContractPriceScheduleLine{
				Id:                              lineID,
				SupplierContractPriceScheduleId: scheduleID,
				SupplierContractLineId:          r.FormValue("supplier_contract_line_id"),
				Currency:                        r.FormValue("currency"),
				UnitPrice:                       unitPrice,
				MinimumAmount:                   minPtr,
				Quantity:                        qtyPtr,
				CycleValueOverride:              cycValPtr,
				CycleUnitOverride:               cycUnitPtr,
				Notes:                           notesPtr,
			},
		}

		_, err := deps.UpdateSupplierContractPriceScheduleLine(ctx, req)
		if err != nil {
			log.Printf("UpdateSupplierContractPriceScheduleLine %s: %v", lineID, err)
			return centymo.HTMXError(err.Error())
		}
		return centymo.HTMXSuccess("schedule-lines-table")
	})
}

// NewDeleteAction handles POST /action/supplier-contract-price-schedule/{id}/lines/delete.
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		if viewCtx.Request.Method != http.MethodPost {
			return view.Error(fmt.Errorf("method not allowed"))
		}
		lineID := viewCtx.Request.FormValue("id")
		if lineID == "" {
			return view.Error(fmt.Errorf("missing id"))
		}
		_, err := deps.DeleteSupplierContractPriceScheduleLine(ctx, &scpslpb.DeleteSupplierContractPriceScheduleLineRequest{
			Data: &scpslpb.SupplierContractPriceScheduleLine{Id: lineID},
		})
		if err != nil {
			log.Printf("DeleteSupplierContractPriceScheduleLine %s: %v", lineID, err)
			return centymo.HTMXError(err.Error())
		}
		return centymo.HTMXSuccess("schedule-lines-table")
	})
}

// --- helpers -----------------------------------------------------------------

func buildEmptyLineFormData(ctx context.Context, deps *Deps, l centymo.SupplierContractPriceScheduleLabels) *LineFormData {
	fd := &LineFormData{
		Labels:       l.Lines.LineForm,
		NounLabels:   l.Lines,
		CommonLabels: deps.CommonLabels,
	}

	if deps.ListSupplierContractLines != nil {
		resp, err := deps.ListSupplierContractLines(ctx, &suppliercontractlinepb.ListSupplierContractLinesRequest{})
		if err == nil {
			for _, line := range resp.GetData() {
				fd.ContractLines = append(fd.ContractLines, types.SelectOption{
					Value: line.GetId(),
					Label: line.GetDescription(),
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
