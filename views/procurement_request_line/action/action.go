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
	procurementrequestlinepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/procurement_request_line"
)

// LineFormData is the template data for the procurement request line drawer form.
type LineFormData struct {
	FormAction           string
	IsEdit               bool
	ID                   string
	ProcurementRequestID string

	// Core fields
	Description           string
	LineType              string
	ProductID             string
	Quantity              string
	EstimatedUnitPrice    string
	EstimatedTotalPrice   string
	ExpenditureCategoryID string
	LocationID            string
	LineNumber            string

	// SPS Wave 3 — F1 fulfillment_mode picker + RECURRING fields
	// FulfillmentMode is the short canonical token: "outright" | "stockable" | "recurring" | "petty"
	FulfillmentMode      string
	FulfillmentModeOptions []types.RadioOption

	// RECURRING-only fields (rendered conditionally when FulfillmentMode == "recurring")
	RecurringCycleValue string
	RecurringCycleUnit  string
	RecurringTermValue  string
	RecurringTermUnit   string
	RecurringUnitOptions []types.SelectOption

	// Options
	Products []types.SelectOption

	Labels       centymo.ProcurementRequestLabels
	CommonLabels pyeza.CommonLabels
}

// Deps holds dependencies for procurement request line action handlers.
type Deps struct {
	Routes                      centymo.ProcurementRequestRoutes
	Labels                      centymo.ProcurementRequestLabels
	CommonLabels                pyeza.CommonLabels
	CreateProcurementRequestLine func(ctx context.Context, req *procurementrequestlinepb.CreateProcurementRequestLineRequest) (*procurementrequestlinepb.CreateProcurementRequestLineResponse, error)
	ReadProcurementRequestLine   func(ctx context.Context, req *procurementrequestlinepb.ReadProcurementRequestLineRequest) (*procurementrequestlinepb.ReadProcurementRequestLineResponse, error)
	UpdateProcurementRequestLine func(ctx context.Context, req *procurementrequestlinepb.UpdateProcurementRequestLineRequest) (*procurementrequestlinepb.UpdateProcurementRequestLineResponse, error)
	DeleteProcurementRequestLine func(ctx context.Context, req *procurementrequestlinepb.DeleteProcurementRequestLineRequest) (*procurementrequestlinepb.DeleteProcurementRequestLineResponse, error)
	ListProducts                 func(ctx context.Context, req *productpb.ListProductsRequest) (*productpb.ListProductsResponse, error)
}

// NewAddAction handles GET+POST for adding a procurement request line.
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		requestID := viewCtx.Request.PathValue("id")
		if requestID == "" {
			return view.Error(fmt.Errorf("missing request id"))
		}

		l := deps.Labels
		if viewCtx.Request.Method == http.MethodGet {
			fd := buildEmptyLineFormData(ctx, deps, l)
			fd.FormAction = viewCtx.Request.URL.Path
			fd.ProcurementRequestID = requestID
			return view.OK("procurement-request-line-drawer-form", fd)
		}

		// POST
		r := viewCtx.Request
		if err := r.ParseForm(); err != nil {
			return view.Error(fmt.Errorf("parse form: %w", err))
		}

		qty, _ := strconv.ParseFloat(r.FormValue("quantity"), 64)
		lineNum, _ := strconv.ParseInt(r.FormValue("line_number"), 10, 32)
		unitPrice := parseCentavos(r.FormValue("estimated_unit_price"))
		totalPrice := parseCentavos(r.FormValue("estimated_total_price"))

		mode := r.FormValue("fulfillment_mode")
		if mode == "" {
			return centymo.HTMXError(deps.Labels.Lines.FormFulfillmentMode + " is required")
		}
		modeEnumValue := parseFulfillmentMode(mode)
		modeEnum := &modeEnumValue

		req := &procurementrequestlinepb.CreateProcurementRequestLineRequest{
			Data: &procurementrequestlinepb.ProcurementRequestLine{
				ProcurementRequestId:  requestID,
				Description:           r.FormValue("description"),
				LineType:              r.FormValue("line_type"),
				ProductId:             optionalString(r.FormValue("product_id")),
				Quantity:              qty,
				EstimatedUnitPrice:    unitPrice,
				EstimatedTotalPrice:   totalPrice,
				ExpenditureCategoryId: optionalString(r.FormValue("expenditure_category_id")),
				LocationId:            optionalString(r.FormValue("location_id")),
				LineNumber:            int32(lineNum),
				FulfillmentMode:       modeEnum,
			},
		}

		// RECURRING — required cycle/term fields
		if mode == "recurring" {
			cycleValue, termValue, err := parseRecurringFields(r)
			if err != nil {
				return centymo.HTMXError(err.Error())
			}
			cycleUnit := r.FormValue("recurring_cycle_unit")
			termUnit := r.FormValue("recurring_term_unit")
			req.Data.RecurringCycleValue = &cycleValue
			req.Data.RecurringCycleUnit = optionalString(cycleUnit)
			req.Data.RecurringTermValue = &termValue
			req.Data.RecurringTermUnit = optionalString(termUnit)
		}

		_, err := deps.CreateProcurementRequestLine(ctx, req)
		if err != nil {
			log.Printf("CreateProcurementRequestLine: %v", err)
			return view.Error(fmt.Errorf("failed to create request line: %w", err))
		}

		return centymo.HTMXSuccess("pr-lines-table")
	})
}

// NewEditAction handles GET+POST for editing a procurement request line.
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		lineID := viewCtx.Request.PathValue("lid")
		requestID := viewCtx.Request.PathValue("id")
		if lineID == "" || requestID == "" {
			return view.Error(fmt.Errorf("missing id or lid"))
		}

		l := deps.Labels
		if viewCtx.Request.Method == http.MethodGet {
			resp, err := deps.ReadProcurementRequestLine(ctx, &procurementrequestlinepb.ReadProcurementRequestLineRequest{
				Data: &procurementrequestlinepb.ProcurementRequestLine{Id: lineID},
			})
			if err != nil {
				return view.Error(fmt.Errorf("failed to read request line: %w", err))
			}
			data := resp.GetData()
			if len(data) == 0 {
				return view.Error(fmt.Errorf("procurement request line not found"))
			}
			line := data[0]

			fd := buildEmptyLineFormData(ctx, deps, l)
			fd.FormAction = viewCtx.Request.URL.Path
			fd.IsEdit = true
			fd.ID = lineID
			fd.ProcurementRequestID = requestID
			fd.Description = line.GetDescription()
			fd.LineType = line.GetLineType()
			fd.ProductID = line.GetProductId()
			fd.Quantity = strconv.FormatFloat(line.GetQuantity(), 'f', 2, 64)
			fd.EstimatedUnitPrice = formatCentavos(line.GetEstimatedUnitPrice())
			fd.EstimatedTotalPrice = formatCentavos(line.GetEstimatedTotalPrice())
			fd.ExpenditureCategoryID = line.GetExpenditureCategoryId()
			fd.LocationID = line.GetLocationId()
			fd.LineNumber = strconv.Itoa(int(line.GetLineNumber()))
			fd.FulfillmentMode = fulfillmentModeToToken(line.GetFulfillmentMode())
			if v := line.GetRecurringCycleValue(); v != 0 {
				fd.RecurringCycleValue = strconv.Itoa(int(v))
			}
			fd.RecurringCycleUnit = line.GetRecurringCycleUnit()
			if v := line.GetRecurringTermValue(); v != 0 {
				fd.RecurringTermValue = strconv.Itoa(int(v))
			}
			fd.RecurringTermUnit = line.GetRecurringTermUnit()
			return view.OK("procurement-request-line-drawer-form", fd)
		}

		// POST
		r := viewCtx.Request
		if err := r.ParseForm(); err != nil {
			return view.Error(fmt.Errorf("parse form: %w", err))
		}

		qty, _ := strconv.ParseFloat(r.FormValue("quantity"), 64)
		lineNum, _ := strconv.ParseInt(r.FormValue("line_number"), 10, 32)
		unitPrice := parseCentavos(r.FormValue("estimated_unit_price"))
		totalPrice := parseCentavos(r.FormValue("estimated_total_price"))

		mode := r.FormValue("fulfillment_mode")
		if mode == "" {
			return centymo.HTMXError(deps.Labels.Lines.FormFulfillmentMode + " is required")
		}
		modeEnumValue := parseFulfillmentMode(mode)
		modeEnum := &modeEnumValue

		req := &procurementrequestlinepb.UpdateProcurementRequestLineRequest{
			Data: &procurementrequestlinepb.ProcurementRequestLine{
				Id:                    lineID,
				ProcurementRequestId:  requestID,
				Description:           r.FormValue("description"),
				LineType:              r.FormValue("line_type"),
				ProductId:             optionalString(r.FormValue("product_id")),
				Quantity:              qty,
				EstimatedUnitPrice:    unitPrice,
				EstimatedTotalPrice:   totalPrice,
				ExpenditureCategoryId: optionalString(r.FormValue("expenditure_category_id")),
				LocationId:            optionalString(r.FormValue("location_id")),
				LineNumber:            int32(lineNum),
				FulfillmentMode:       modeEnum,
			},
		}

		// RECURRING — required cycle/term fields
		if mode == "recurring" {
			cycleValue, termValue, err := parseRecurringFields(r)
			if err != nil {
				return centymo.HTMXError(err.Error())
			}
			cycleUnit := r.FormValue("recurring_cycle_unit")
			termUnit := r.FormValue("recurring_term_unit")
			req.Data.RecurringCycleValue = &cycleValue
			req.Data.RecurringCycleUnit = optionalString(cycleUnit)
			req.Data.RecurringTermValue = &termValue
			req.Data.RecurringTermUnit = optionalString(termUnit)
		}

		_, err := deps.UpdateProcurementRequestLine(ctx, req)
		if err != nil {
			log.Printf("UpdateProcurementRequestLine %s: %v", lineID, err)
			return view.Error(fmt.Errorf("failed to update request line: %w", err))
		}

		return centymo.HTMXSuccess("pr-lines-table")
	})
}

// NewDeleteAction handles POST for deleting a procurement request line.
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		if viewCtx.Request.Method != http.MethodPost {
			return view.Error(fmt.Errorf("method not allowed"))
		}
		lineID := viewCtx.Request.FormValue("id")
		if lineID == "" {
			return view.Error(fmt.Errorf("missing id"))
		}
		_, err := deps.DeleteProcurementRequestLine(ctx, &procurementrequestlinepb.DeleteProcurementRequestLineRequest{
			Data: &procurementrequestlinepb.ProcurementRequestLine{Id: lineID},
		})
		if err != nil {
			log.Printf("DeleteProcurementRequestLine %s: %v", lineID, err)
			return view.Error(fmt.Errorf("failed to delete request line: %w", err))
		}
		return centymo.HTMXSuccess("pr-lines-table")
	})
}

// NewRetrySpawnAction is the SPS Wave 3 CRIT-3 retry placeholder. The actual
// per-line spawn-retry use case ships in a later wave; for now this handler
// logs the intent and refreshes the lines table so the failed-status indicator
// is re-rendered. Wired so the operator-facing button works end-to-end UX-wise.
func NewRetrySpawnAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		if viewCtx.Request.Method != http.MethodPost {
			return view.Error(fmt.Errorf("method not allowed"))
		}
		lineID := viewCtx.Request.PathValue("lid")
		requestID := viewCtx.Request.PathValue("id")
		if lineID == "" || requestID == "" {
			return view.Error(fmt.Errorf("missing id or lid"))
		}
		log.Printf("RetrySpawn placeholder: pr=%s line=%s — actual retry use case is out of SPS Wave 3 scope", requestID, lineID)
		return centymo.HTMXSuccess("pr-lines-table")
	})
}

// --- helpers -----------------------------------------------------------------

func buildEmptyLineFormData(ctx context.Context, deps *Deps, l centymo.ProcurementRequestLabels) *LineFormData {
	fd := &LineFormData{
		Labels:       l,
		CommonLabels: deps.CommonLabels,
	}

	// SPS Wave 3 — F1 fulfillment_mode picker (4 radio options + per-mode hint).
	fd.FulfillmentModeOptions = []types.RadioOption{
		{Value: "outright", Label: l.FulfillmentMode.Outright, Description: l.FulfillmentModeHints.Outright},
		{Value: "stockable", Label: l.FulfillmentMode.Stockable, Description: l.FulfillmentModeHints.Stockable},
		{Value: "recurring", Label: l.FulfillmentMode.Recurring, Description: l.FulfillmentModeHints.Recurring},
		{Value: "petty", Label: l.FulfillmentMode.Petty, Description: l.FulfillmentModeHints.Petty},
	}

	// RECURRING — cycle/term unit dropdowns share the same vocabulary.
	fd.RecurringUnitOptions = []types.SelectOption{
		{Value: "day", Label: l.Lines.FormRecurringUnitDay},
		{Value: "week", Label: l.Lines.FormRecurringUnitWeek},
		{Value: "month", Label: l.Lines.FormRecurringUnitMonth},
		{Value: "year", Label: l.Lines.FormRecurringUnitYear},
	}

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

// parseFulfillmentMode maps the form's short token into the proto enum value.
// Falls back to UNSPECIFIED on unknown tokens; caller is expected to validate
// non-empty mode upstream.
func parseFulfillmentMode(token string) procurementrequestlinepb.ProcurementRequestLineFulfillmentMode {
	switch token {
	case "outright":
		return procurementrequestlinepb.ProcurementRequestLineFulfillmentMode_PROCUREMENT_REQUEST_LINE_FULFILLMENT_MODE_OUTRIGHT
	case "stockable":
		return procurementrequestlinepb.ProcurementRequestLineFulfillmentMode_PROCUREMENT_REQUEST_LINE_FULFILLMENT_MODE_STOCKABLE
	case "recurring":
		return procurementrequestlinepb.ProcurementRequestLineFulfillmentMode_PROCUREMENT_REQUEST_LINE_FULFILLMENT_MODE_RECURRING
	case "petty":
		return procurementrequestlinepb.ProcurementRequestLineFulfillmentMode_PROCUREMENT_REQUEST_LINE_FULFILLMENT_MODE_PETTY
	}
	return procurementrequestlinepb.ProcurementRequestLineFulfillmentMode_PROCUREMENT_REQUEST_LINE_FULFILLMENT_MODE_UNSPECIFIED
}

// fulfillmentModeToToken is the inverse of parseFulfillmentMode for edit-form
// pre-population — returns the short canonical token used in the radio group.
func fulfillmentModeToToken(mode procurementrequestlinepb.ProcurementRequestLineFulfillmentMode) string {
	switch mode {
	case procurementrequestlinepb.ProcurementRequestLineFulfillmentMode_PROCUREMENT_REQUEST_LINE_FULFILLMENT_MODE_OUTRIGHT:
		return "outright"
	case procurementrequestlinepb.ProcurementRequestLineFulfillmentMode_PROCUREMENT_REQUEST_LINE_FULFILLMENT_MODE_STOCKABLE:
		return "stockable"
	case procurementrequestlinepb.ProcurementRequestLineFulfillmentMode_PROCUREMENT_REQUEST_LINE_FULFILLMENT_MODE_RECURRING:
		return "recurring"
	case procurementrequestlinepb.ProcurementRequestLineFulfillmentMode_PROCUREMENT_REQUEST_LINE_FULFILLMENT_MODE_PETTY:
		return "petty"
	}
	return ""
}

// parseRecurringFields validates and extracts the 2 numeric inputs of the
// RECURRING-mode form. Returns a translated error message when either is
// missing or non-positive (matching plan validation requirement).
func parseRecurringFields(r *http.Request) (cycleValue int32, termValue int32, err error) {
	cv, errC := strconv.Atoi(r.FormValue("recurring_cycle_value"))
	if errC != nil || cv <= 0 {
		return 0, 0, fmt.Errorf("recurring_cycle_value must be a positive integer")
	}
	tv, errT := strconv.Atoi(r.FormValue("recurring_term_value"))
	if errT != nil || tv <= 0 {
		return 0, 0, fmt.Errorf("recurring_term_value must be a positive integer")
	}
	if r.FormValue("recurring_cycle_unit") == "" {
		return 0, 0, fmt.Errorf("recurring_cycle_unit is required")
	}
	if r.FormValue("recurring_term_unit") == "" {
		return 0, 0, fmt.Errorf("recurring_term_unit is required")
	}
	return int32(cv), int32(tv), nil
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
