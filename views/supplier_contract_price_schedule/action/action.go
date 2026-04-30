package action

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	centymo "github.com/erniealice/centymo-golang"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	suppliercontractpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/supplier_contract"
	scpspb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/supplier_contract_price_schedule"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// StatusOption is a select option for the schedule status enum.
// Description is required by the shared form-group select template (reads
// .Description on every option to set data-description) — keep the field
// even when unused so template execution doesn't panic.
type StatusOption struct {
	Value       string
	Label       string
	Selected    bool
	Description string
}

// FormData is the template data for the SCPS drawer form.
type FormData struct {
	FormAction string
	IsEdit     bool
	ID         string

	// Identity
	Name        string
	Description string
	InternalID  string

	// Scoping
	SupplierContractID    string
	SupplierContractLabel string
	SupplierContracts     []types.SelectOption

	// Validity (date-only — half-open [start, end), end blank = open-ended)
	DateStart   string
	DateEnd     string
	OpenEnded   bool

	// Money / location
	Currency           string
	LocationID         string
	Locations          []types.SelectOption

	// Lifecycle
	Status         string
	StatusOptions  []StatusOption
	SequenceNumber string

	// Notes
	Notes string

	Labels       centymo.SupplierContractPriceScheduleFormLabels
	StatusLabels centymo.SupplierContractPriceScheduleStatusLabels
	CommonLabels pyeza.CommonLabels
}

// Deps holds all dependencies for the SCPS action handlers.
type Deps struct {
	Routes       centymo.SupplierContractPriceScheduleRoutes
	Labels       centymo.SupplierContractPriceScheduleLabels
	CommonLabels pyeza.CommonLabels

	// CRUD
	CreateSupplierContractPriceSchedule func(ctx context.Context, req *scpspb.CreateSupplierContractPriceScheduleRequest) (*scpspb.CreateSupplierContractPriceScheduleResponse, error)
	ReadSupplierContractPriceSchedule   func(ctx context.Context, req *scpspb.ReadSupplierContractPriceScheduleRequest) (*scpspb.ReadSupplierContractPriceScheduleResponse, error)
	UpdateSupplierContractPriceSchedule func(ctx context.Context, req *scpspb.UpdateSupplierContractPriceScheduleRequest) (*scpspb.UpdateSupplierContractPriceScheduleResponse, error)
	DeleteSupplierContractPriceSchedule func(ctx context.Context, req *scpspb.DeleteSupplierContractPriceScheduleRequest) (*scpspb.DeleteSupplierContractPriceScheduleResponse, error)

	// Workflow
	ActivateSupplierContractPriceSchedule  func(ctx context.Context, id string) error
	SupersedeSupplierContractPriceSchedule func(ctx context.Context, id, reason string) error

	// Status setter (drawer-driven)
	SetSupplierContractPriceScheduleStatus func(ctx context.Context, id, status string) error

	// Lookups for drawer dropdowns
	ListSupplierContracts func(ctx context.Context, req *suppliercontractpb.ListSupplierContractsRequest) (*suppliercontractpb.ListSupplierContractsResponse, error)
}

// NewAddAction handles GET+POST /action/supplier-contract-price-schedule/add.
// Optional ?supplier_contract_id= query parameter pre-selects the parent
// contract when the drawer is opened from the supplier_contract detail tab.
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		l := deps.Labels
		if viewCtx.Request.Method == http.MethodGet {
			fd := buildEmptyFormData(ctx, deps, l)
			fd.FormAction = deps.Routes.AddURL
			fd.Status = "SUPPLIER_CONTRACT_PRICE_SCHEDULE_STATUS_SCHEDULED"
			fd.Currency = "PHP"
			fd.SequenceNumber = "1"
			// Pre-select parent contract if provided via query string.
			if scID := viewCtx.Request.URL.Query().Get("supplier_contract_id"); scID != "" {
				fd.SupplierContractID = scID
			}
			return view.OK("supplier-contract-price-schedule-drawer-form", fd)
		}

		// POST
		r := viewCtx.Request
		if err := r.ParseForm(); err != nil {
			return view.Error(fmt.Errorf("parse form: %w", err))
		}

		seqNum, _ := strconv.ParseInt(r.FormValue("sequence_number"), 10, 32)
		if seqNum <= 0 {
			seqNum = 1
		}
		openEnded := r.FormValue("open_ended") == "true"

		schedule := &scpspb.SupplierContractPriceSchedule{
			SupplierContractId: r.FormValue("supplier_contract_id"),
			Name:               r.FormValue("name"),
			Description:        optionalString(r.FormValue("description")),
			DateTimeStart:      parseDateUTC(r.FormValue("date_start"), false),
			DateTimeEnd:        parseEndDate(r.FormValue("date_end"), openEnded),
			LocationId:         optionalString(r.FormValue("location_id")),
			Currency:           r.FormValue("currency"),
			Status:             parseStatus(r.FormValue("status")),
			SequenceNumber:     int32(seqNum),
			Notes:              optionalString(r.FormValue("notes")),
			Active:             true,
		}
		if internal := r.FormValue("internal_id"); internal != "" {
			schedule.InternalId = internal
		}

		_, err := deps.CreateSupplierContractPriceSchedule(ctx, &scpspb.CreateSupplierContractPriceScheduleRequest{
			Data: schedule,
		})
		if err != nil {
			log.Printf("CreateSupplierContractPriceSchedule: %v", err)
			return centymo.HTMXError(err.Error())
		}
		return centymo.HTMXSuccess("supplier-contract-price-schedules-table")
	})
}

// NewEditAction handles GET+POST /action/supplier-contract-price-schedule/edit/{id}.
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		l := deps.Labels
		id := viewCtx.Request.PathValue("id")
		if id == "" {
			return view.Error(fmt.Errorf("missing id"))
		}

		if viewCtx.Request.Method == http.MethodGet {
			resp, err := deps.ReadSupplierContractPriceSchedule(ctx, &scpspb.ReadSupplierContractPriceScheduleRequest{
				Data: &scpspb.SupplierContractPriceSchedule{Id: id},
			})
			if err != nil {
				return view.Error(fmt.Errorf("failed to read price schedule: %w", err))
			}
			data := resp.GetData()
			if len(data) == 0 {
				return view.Error(fmt.Errorf("price schedule not found"))
			}
			s := data[0]

			fd := buildEmptyFormData(ctx, deps, l)
			fd.FormAction = route.ResolveURL(deps.Routes.EditURL, "id", id)
			fd.IsEdit = true
			fd.ID = id
			fd.Name = s.GetName()
			fd.Description = s.GetDescription()
			fd.InternalID = s.GetInternalId()
			fd.SupplierContractID = s.GetSupplierContractId()
			fd.DateStart = formatDateUTC(s.GetDateTimeStart())
			if end := s.GetDateTimeEnd(); end == nil {
				fd.OpenEnded = true
				fd.DateEnd = ""
			} else {
				fd.DateEnd = formatDateUTC(end)
			}
			fd.LocationID = s.GetLocationId()
			fd.Currency = s.GetCurrency()
			fd.Status = s.GetStatus().String()
			fd.SequenceNumber = strconv.FormatInt(int64(s.GetSequenceNumber()), 10)
			fd.Notes = s.GetNotes()
			return view.OK("supplier-contract-price-schedule-drawer-form", fd)
		}

		// POST
		r := viewCtx.Request
		if err := r.ParseForm(); err != nil {
			return view.Error(fmt.Errorf("parse form: %w", err))
		}

		seqNum, _ := strconv.ParseInt(r.FormValue("sequence_number"), 10, 32)
		if seqNum <= 0 {
			seqNum = 1
		}
		openEnded := r.FormValue("open_ended") == "true"

		schedule := &scpspb.SupplierContractPriceSchedule{
			Id:                 id,
			SupplierContractId: r.FormValue("supplier_contract_id"),
			Name:               r.FormValue("name"),
			Description:        optionalString(r.FormValue("description")),
			DateTimeStart:      parseDateUTC(r.FormValue("date_start"), false),
			DateTimeEnd:        parseEndDate(r.FormValue("date_end"), openEnded),
			LocationId:         optionalString(r.FormValue("location_id")),
			Currency:           r.FormValue("currency"),
			Status:             parseStatus(r.FormValue("status")),
			SequenceNumber:     int32(seqNum),
			Notes:              optionalString(r.FormValue("notes")),
		}

		_, err := deps.UpdateSupplierContractPriceSchedule(ctx, &scpspb.UpdateSupplierContractPriceScheduleRequest{
			Data: schedule,
		})
		if err != nil {
			log.Printf("UpdateSupplierContractPriceSchedule %s: %v", id, err)
			return centymo.HTMXError(err.Error())
		}
		return centymo.HTMXSuccess("supplier-contract-price-schedules-table")
	})
}

// NewDeleteAction handles POST /action/supplier-contract-price-schedule/delete.
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		if viewCtx.Request.Method != http.MethodPost {
			return view.Error(fmt.Errorf("method not allowed"))
		}
		id := viewCtx.Request.FormValue("id")
		if id == "" {
			return view.Error(fmt.Errorf("missing id"))
		}
		_, err := deps.DeleteSupplierContractPriceSchedule(ctx, &scpspb.DeleteSupplierContractPriceScheduleRequest{
			Data: &scpspb.SupplierContractPriceSchedule{Id: id},
		})
		if err != nil {
			log.Printf("DeleteSupplierContractPriceSchedule %s: %v", id, err)
			return centymo.HTMXError(err.Error())
		}
		return centymo.HTMXSuccess("supplier-contract-price-schedules-table")
	})
}

// NewSetStatusAction handles POST /action/supplier-contract-price-schedule/set-status.
func NewSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		if viewCtx.Request.Method != http.MethodPost {
			return view.Error(fmt.Errorf("method not allowed"))
		}
		id := viewCtx.Request.URL.Query().Get("id")
		status := viewCtx.Request.URL.Query().Get("status")
		if id == "" || status == "" {
			return view.Error(fmt.Errorf("missing id or status"))
		}
		if deps.SetSupplierContractPriceScheduleStatus != nil {
			if err := deps.SetSupplierContractPriceScheduleStatus(ctx, id, status); err != nil {
				return centymo.HTMXError(err.Error())
			}
		}
		return centymo.HTMXSuccess("supplier-contract-price-schedules-table")
	})
}

// NewBulkSetStatusAction handles POST /action/supplier-contract-price-schedule/bulk-set-status.
func NewBulkSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		if viewCtx.Request.Method != http.MethodPost {
			return view.Error(fmt.Errorf("method not allowed"))
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.Error(fmt.Errorf("parse form: %w", err))
		}
		ids := viewCtx.Request.Form["id"]
		targetStatus := viewCtx.Request.FormValue("target_status")
		if len(ids) == 0 || targetStatus == "" {
			return view.Error(fmt.Errorf("missing ids or target_status"))
		}
		if deps.SetSupplierContractPriceScheduleStatus != nil {
			for _, id := range ids {
				if err := deps.SetSupplierContractPriceScheduleStatus(ctx, id, targetStatus); err != nil {
					log.Printf("BulkSetStatus %s → %s: %v", id, targetStatus, err)
				}
			}
		}
		return centymo.HTMXSuccess("supplier-contract-price-schedules-table")
	})
}

// --- helpers -----------------------------------------------------------------

func buildEmptyFormData(ctx context.Context, deps *Deps, l centymo.SupplierContractPriceScheduleLabels) *FormData {
	fd := &FormData{
		Labels:       l.Form,
		StatusLabels: l.Status,
		CommonLabels: deps.CommonLabels,
	}

	fd.StatusOptions = []StatusOption{
		{Value: "SUPPLIER_CONTRACT_PRICE_SCHEDULE_STATUS_SCHEDULED", Label: l.Status.Scheduled},
		{Value: "SUPPLIER_CONTRACT_PRICE_SCHEDULE_STATUS_ACTIVE", Label: l.Status.Active},
		{Value: "SUPPLIER_CONTRACT_PRICE_SCHEDULE_STATUS_SUPERSEDED", Label: l.Status.Superseded},
		{Value: "SUPPLIER_CONTRACT_PRICE_SCHEDULE_STATUS_CANCELLED", Label: l.Status.Cancelled},
	}

	// Load supplier contracts for the dropdown
	if deps.ListSupplierContracts != nil {
		resp, err := deps.ListSupplierContracts(ctx, &suppliercontractpb.ListSupplierContractsRequest{})
		if err == nil {
			for _, c := range resp.GetData() {
				fd.SupplierContracts = append(fd.SupplierContracts, types.SelectOption{
					Value: c.GetId(),
					Label: c.GetName(),
				})
			}
		}
	}

	return fd
}

// parseDateUTC parses YYYY-MM-DD as midnight UTC.
func parseDateUTC(date string, _ bool) *timestamppb.Timestamp {
	if date == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil
	}
	return timestamppb.New(t.UTC())
}

// parseEndDate handles the open-ended checkbox: when checked, end is nil.
// Otherwise parses as YYYY-MM-DD at end-of-day UTC (23:59:59).
func parseEndDate(date string, openEnded bool) *timestamppb.Timestamp {
	if openEnded || date == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil
	}
	// End is exclusive in the half-open window; render as end-of-day UTC.
	endOfDay := t.UTC().Add(24*time.Hour - time.Second)
	return timestamppb.New(endOfDay)
}

func formatDateUTC(ts *timestamppb.Timestamp) string {
	if ts == nil || !ts.IsValid() {
		return ""
	}
	return ts.AsTime().UTC().Format("2006-01-02")
}

func parseStatus(s string) scpspb.SupplierContractPriceScheduleStatus {
	if v, ok := scpspb.SupplierContractPriceScheduleStatus_value[s]; ok {
		return scpspb.SupplierContractPriceScheduleStatus(v)
	}
	return scpspb.SupplierContractPriceScheduleStatus_SUPPLIER_CONTRACT_PRICE_SCHEDULE_STATUS_SCHEDULED
}

func optionalString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// NewActivateAction handles POST /action/supplier-contract-price-schedule/activate/{id}.
func NewActivateAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		if viewCtx.Request.Method != http.MethodPost {
			return view.Error(fmt.Errorf("method not allowed"))
		}
		id := viewCtx.Request.PathValue("id")
		if id == "" {
			return view.Error(fmt.Errorf("missing id"))
		}
		if deps.ActivateSupplierContractPriceSchedule != nil {
			if err := deps.ActivateSupplierContractPriceSchedule(ctx, id); err != nil {
				log.Printf("ActivateSupplierContractPriceSchedule %s: %v", id, err)
				return centymo.HTMXError(err.Error())
			}
		}
		detailURL := route.ResolveURL(deps.Routes.DetailURL, "id", id)
		return view.ViewResult{
			StatusCode: http.StatusOK,
			Headers: map[string]string{
				"HX-Redirect": detailURL,
			},
		}
	})
}

// NewSupersedeAction handles POST /action/supplier-contract-price-schedule/supersede/{id}.
func NewSupersedeAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		if viewCtx.Request.Method != http.MethodPost {
			return view.Error(fmt.Errorf("method not allowed"))
		}
		id := viewCtx.Request.PathValue("id")
		if id == "" {
			return view.Error(fmt.Errorf("missing id"))
		}
		reason := ""
		if err := viewCtx.Request.ParseForm(); err == nil {
			reason = viewCtx.Request.FormValue("reason")
		}
		if deps.SupersedeSupplierContractPriceSchedule != nil {
			if err := deps.SupersedeSupplierContractPriceSchedule(ctx, id, reason); err != nil {
				log.Printf("SupersedeSupplierContractPriceSchedule %s: %v", id, err)
				return centymo.HTMXError(err.Error())
			}
		}
		detailURL := route.ResolveURL(deps.Routes.DetailURL, "id", id)
		return view.ViewResult{
			StatusCode: http.StatusOK,
			Headers: map[string]string{
				"HX-Redirect": detailURL,
			},
		}
	})
}
