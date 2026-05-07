package action

import (
	"context"
	"math/rand"
	"strings"
	"time"

	pyeza "github.com/erniealice/pyeza-golang"
	pyezatypes "github.com/erniealice/pyeza-golang/types"

	centymo "github.com/erniealice/centymo-golang"
	"github.com/erniealice/centymo-golang/views/subscription/form"

	clientpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client"
	jobtemplatepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/operation/job_template"
	jobtemplatephasepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/operation/job_template_phase"
	jobtemplaterelationpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/operation/job_template_relation"
	jobtemplatetaskpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/operation/job_template_task"
	revenuepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/revenue"
	billingeventpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/billing_event"
	planpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/plan"
	priceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_plan"
	priceschedulepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_schedule"
	subscriptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Deps holds dependencies for subscription action handlers.
type Deps struct {
	Routes       centymo.SubscriptionRoutes
	Labels       centymo.SubscriptionLabels
	CommonLabels pyeza.CommonLabels

	CreateSubscription  func(ctx context.Context, req *subscriptionpb.CreateSubscriptionRequest) (*subscriptionpb.CreateSubscriptionResponse, error)
	ReadSubscription    func(ctx context.Context, req *subscriptionpb.ReadSubscriptionRequest) (*subscriptionpb.ReadSubscriptionResponse, error)
	// GetSubscriptionItemPageData returns the subscription with its joined
	// Client (+ User) and PricePlan (+ Plan) populated. Edit drawer uses it
	// to render the customer name (not the bare client_id) without depending
	// on a separate ListClients-and-iterate fallback.
	GetSubscriptionItemPageData func(ctx context.Context, req *subscriptionpb.GetSubscriptionItemPageDataRequest) (*subscriptionpb.GetSubscriptionItemPageDataResponse, error)
	UpdateSubscription  func(ctx context.Context, req *subscriptionpb.UpdateSubscriptionRequest) (*subscriptionpb.UpdateSubscriptionResponse, error)
	DeleteSubscription  func(ctx context.Context, req *subscriptionpb.DeleteSubscriptionRequest) (*subscriptionpb.DeleteSubscriptionResponse, error)
	ListClients         func(ctx context.Context, req *clientpb.ListClientsRequest) (*clientpb.ListClientsResponse, error)
	ListPlans           func(ctx context.Context, req *planpb.ListPlansRequest) (*planpb.ListPlansResponse, error)
	ReadPlan            func(ctx context.Context, req *planpb.ReadPlanRequest) (*planpb.ReadPlanResponse, error)
	SearchClientsByName func(ctx context.Context, req *clientpb.SearchClientsByNameRequest) (*clientpb.SearchClientsByNameResponse, error)
	SearchPlansByName   func(ctx context.Context, req *planpb.SearchPlansByNameRequest) (*planpb.SearchPlansByNameResponse, error)
	ListPricePlans      func(ctx context.Context, req *priceplanpb.ListPricePlansRequest) (*priceplanpb.ListPricePlansResponse, error)
	ReadPricePlan       func(ctx context.Context, req *priceplanpb.ReadPricePlanRequest) (*priceplanpb.ReadPricePlanResponse, error)
	ListPriceSchedules  func(ctx context.Context, req *priceschedulepb.ListPriceSchedulesRequest) (*priceschedulepb.ListPriceSchedulesResponse, error)

	// SetSubscriptionActive performs a raw DB update of the active field.
	// Required for set-status and bulk-set-status handlers.
	// Uses raw update (not proto) because proto3 omits bool=false on serialization.
	SetSubscriptionActive func(ctx context.Context, id string, active bool) error

	// GetInUseIDs checks whether subscription IDs are referenced by dependent records.
	// Used by the bulk-delete handler to skip in-use rows.
	GetInUseIDs func(ctx context.Context, ids []string) (map[string]bool, error)

	// RecognizeRevenueFromSubscription invokes the espyna use case that
	// materializes a Revenue + N RevenueLineItems for a billing period.
	// Set when the centymo block wiring threads the use case through. Used
	// by NewRecognizeAction (drawer GET dry-run + POST commit) and by the
	// existing manual revenue-add flow's auto-populate path (skip_header=true).
	RecognizeRevenueFromSubscription func(ctx context.Context, req *revenuepb.CreateRevenueWithLineItemsRequest) (*revenuepb.CreateRevenueWithLineItemsResponse, error)

	// CustomClientPriceScheduleLabelSuffix carries the lyngua-resolved suffix
	// appended to a client's name when constructing the default custom
	// PriceSchedule name (e.g. "Price Schedule" / "Rate Cards"). Read by the
	// customize handler; sourced from PriceScheduleLabels.Form by block.go.
	CustomClientPriceScheduleLabelSuffix string

	// CustomizePlanForClient invokes the espyna use case that clones the
	// source Plan + PricePlan into a client-scoped copy and (optionally)
	// repoints the subscription onto the new PricePlan. See plan §4.
	// Wired by the centymo block when the use case is available.
	CustomizePlanForClient func(ctx context.Context, req *CustomizePlanForClientRequest) (*CustomizePlanForClientResponse, error)

	// 2026-04-29 milestone-billing plan §5 / Phase D — BillingEvent operations
	// for the subscription Package tab Milestones section + recognize drawer
	// milestone select. nil-safe: when unset (no adapter registered), the
	// drawer falls back to the legacy non-milestone branches and the Package
	// tab milestone section is skipped.
	ListBillingEventsBySubscription func(ctx context.Context, req *billingeventpb.ListBillingEventsBySubscriptionRequest) (*billingeventpb.ListBillingEventsBySubscriptionResponse, error)
	SetBillingEventStatus           func(ctx context.Context, req *billingeventpb.SetBillingEventStatusRequest) (*billingeventpb.SetBillingEventStatusResponse, error)

	// 2026-04-29 auto-spawn-jobs-from-subscription plan §5 / Phase D —
	// JobTemplate read deps used by:
	//   1. The Spawn Jobs section detection on the create form (resolves
	//      Plan.job_template_id → JobTemplate + JobTemplateRelation rows
	//      → phase + task counts).
	//   2. The retroactive spawn drawer (lists detected templates).
	// All are nil-safe — when unwired, the section/drawer hide.
	ReadJobTemplate          func(ctx context.Context, req *jobtemplatepb.ReadJobTemplateRequest) (*jobtemplatepb.ReadJobTemplateResponse, error)
	ListJobTemplatePhases    func(ctx context.Context, req *jobtemplatephasepb.ListByJobTemplateRequest) (*jobtemplatephasepb.ListByJobTemplateResponse, error)
	ListJobTemplateTasks     func(ctx context.Context, req *jobtemplatetaskpb.ListJobTemplateTasksByPhaseRequest) (*jobtemplatetaskpb.ListJobTemplateTasksByPhaseResponse, error)
	ListJobTemplateRelations func(ctx context.Context, req *jobtemplaterelationpb.ListJobTemplateRelationsByParentRequest) (*jobtemplaterelationpb.ListJobTemplateRelationsByParentResponse, error)

	// MaterializeJobsForSubscription is the espyna use case wired through
	// the subscription block. Used by the retroactive spawn handler. nil-safe.
	MaterializeJobsForSubscription func(ctx context.Context, subscriptionID string, spawnJobs bool) (jobCount int, skippedReason string, err error)

	// ListRevenueRunCandidates enumerates un-invoiced billing periods for the
	// subscription. Wired by block.go when the use case is available. nil-safe —
	// the revenue-run drawer falls back to an empty candidate list.
	ListRevenueRunCandidates func(ctx context.Context, scope RevenueRunScopeAction) ([]RevenueRunCandidateAction, string, error)

	// GenerateRevenueRun executes the batch revenue generation run for this
	// subscription. Wired by block.go. nil-safe — POST returns an error when unset.
	GenerateRevenueRun func(ctx context.Context, scope RevenueRunScopeAction, sels RevenueRunSelectionsAction) (*RevenueRunResultAction, error)

	// 2026-04-30 cyclic-subscription-jobs plan §5.3 / Phase D — adapter for
	// espyna's MaterializeInstanceJobsForSubscription consumer. Wired by
	// centymo block.go after both this Phase D and the parallel block.go-
	// owning agent finish (main thread coordinates). nil-safe — the cycle-
	// spawn / backfill action handlers gate on it explicitly.
	MaterializeInstanceJobsForSubscription MaterializeInstanceJobsForSubscriptionAdapter
}

// MaterializeInstanceJobsRequest mirrors the espyna consumer-surface request.
// Keeping a centymo-local struct so this package does not import espyna directly.
// 2026-04-30 cyclic-subscription-jobs plan §5.3 / Phase D.
type MaterializeInstanceJobsRequest struct {
	SubscriptionID   string
	CyclePeriodStart string
	Backfill         bool
	// 2026-05-01 ad-hoc-subscription-billing plan §3.2 — operator-supplied
	// usage request date for AD_HOC plans. Empty defaults to today UTC.
	UsageRequestDate string
}

// MaterializeInstanceJobsResponse mirrors the espyna consumer-surface response.
type MaterializeInstanceJobsResponse struct {
	SpawnedCycleCount         int
	SpawnedJobCount           int
	OnceAtStartJobCount       int
	EngagementWasNewlyCreated bool
	SkippedReason             string
	BackfillCappedAt          int32
}

// MaterializeInstanceJobsForSubscriptionAdapter is the function-pointer type
// the centymo block.go wires once the espyna consumer is available.
type MaterializeInstanceJobsForSubscriptionAdapter func(
	ctx context.Context, req *MaterializeInstanceJobsRequest,
) (*MaterializeInstanceJobsResponse, error)

// CustomizePlanForClientRequest mirrors the espyna use-case request shape
// (plan §4.1). Centymo handlers build this and pass it through Deps.
// The `derivedName` carries the per-tier "{Client.name} - {suffix}" label
// resolved on the centymo side from typed labels (plan §4.4.1 step 2-3).
type CustomizePlanForClientRequest struct {
	SourcePlanID      string
	SourcePricePlanID string
	ClientID          string
	SubscriptionID    string
	NewScheduleName   string
}

// CustomizePlanForClientResponse mirrors the espyna use-case response shape
// (plan §4.1). Only the fields centymo's POST handler needs are surfaced;
// extend if a future caller needs the cloned proto records.
type CustomizePlanForClientResponse struct {
	NewPlanID      string
	NewPricePlanID string
	NewScheduleID  string
	Reused         bool
}

// ---------------------------------------------------------------------------
// Revenue-run view-local types
//
// Defined here (in the action package) so that subscriptionaction.Deps can
// carry the callback signatures. The revenue_run sub-package imports these
// via its own identically-shaped local types; block.go translates consumer.*
// shapes into these types when wiring the callbacks. Plan rule D12: types are
// duplicated per package rather than shared to keep view packages decoupled
// from each other.
// ---------------------------------------------------------------------------

// RevenueRunScopeAction is the view-layer scope for revenue-run callbacks.
// (Defined in subscriptionaction to avoid an import cycle with the revenue_run
// sub-package; the sub-package has its own mirror type.)
type RevenueRunScopeAction struct {
	WorkspaceID    string
	ClientID       string
	SubscriptionID string
	AsOfDate       string
	Cursor         string
	Limit          int32
}

// RevenueRunCandidateAction is the view-layer representation of one pending period.
type RevenueRunCandidateAction struct {
	SubscriptionID    string
	SubscriptionName  string
	ClientID          string
	ClientName        string
	PlanName          string
	BillingCycleLabel string
	Currency          string
	PeriodStart       string
	PeriodEnd         string
	PeriodLabel       string
	PeriodMarker      string
	Amount            int64
	AmountDisplay     string
	LineItemCount     int
	Eligible          bool
	BlockerReason     string
}

// SelectedRevenueRunCandidateAction is one confirmed selection.
type SelectedRevenueRunCandidateAction struct {
	SubscriptionID string
	PeriodStart    string
	PeriodEnd      string
	PeriodMarker   string
}

// RevenueRunSelectionsAction carries operator selections for GenerateRevenueRun.
type RevenueRunSelectionsAction struct {
	ExplicitList []SelectedRevenueRunCandidateAction
	FilterToken  string
}

// RevenueRunResultAction is the output of a GenerateRevenueRun call.
type RevenueRunResultAction struct {
	RunID   string
	Status  string
	Created int32
	Skipped int32
	Errored int32
}

// buildFormLabels builds a form.Labels from centymo.SubscriptionLabels.
// This is the dep-bearing helper that resolves typed labels into the
// form package's flat Labels shape. Lives here (not in form/) because
// centymo.SubscriptionLabels is a dependency-bearing type.
func buildFormLabels(l centymo.SubscriptionLabels) form.Labels {
	return form.Labels{
		Customer:                  l.Form.Customer,
		CustomerPlaceholder:       l.Form.CustomerPlaceholder,
		Plan:                      l.Form.Plan,
		PlanPlaceholder:           l.Form.PlanPlaceholder,
		StartDate:                 l.Form.StartDate,
		EndDate:                   l.Form.EndDate,
		StartTime:                 l.Form.StartTime,
		EndTime:                   l.Form.EndTime,
		TimePlaceholder:           l.Form.TimePlaceholder,
		Timezone:                  l.Form.Timezone,
		Notes:                     l.Form.Notes,
		NotesPlaceholder:          l.Form.NotesPlaceholder,
		CustomerSearchPlaceholder: l.Form.CustomerSearchPlaceholder,
		PlanSearchPlaceholder:     l.Form.PlanSearchPlaceholder,
		CustomerNoResults:         l.Form.CustomerNoResults,
		PlanNoResults:             l.Form.PlanNoResults,
		Code:                      l.Form.Code,
		CodePlaceholder:           l.Form.CodePlaceholder,
		CustomerInfo:              l.Form.CustomerInfo,
		PlanInfo:                  l.Form.PlanInfo,
		CodeInfo:                  l.Form.CodeInfo,
		StartDateInfo:             l.Form.StartDateInfo,
		EndDateInfo:               l.Form.EndDateInfo,
		StartTimeInfo:             l.Form.StartTimeInfo,
		EndTimeInfo:               l.Form.EndTimeInfo,
		NotesInfo:                 l.Form.NotesInfo,
		StartDateRowHelp:          l.Form.StartDateRowHelp,
		EndDateRowHelp:            l.Form.EndDateRowHelp,
		PlanGroupForClient:        l.Form.PlanGroupForClient,
		PlanGroupGeneral:          l.Form.PlanGroupGeneral,
		PlanClientScopeNotice:     l.Form.PlanClientScopeNotice,
		EditLockedReason:          l.Form.EditLockedReason,
		SpawnJobsSectionTitle:     l.Form.SpawnJobsSectionTitle,
		SpawnJobsToggle:           l.Form.SpawnJobsToggle,
		SpawnJobsHelpText:         l.Form.SpawnJobsHelpText,
		SpawnJobsSummary:          l.Form.SpawnJobsSummary,
		SpawnJobsNone:             l.Form.SpawnJobsNone,
	}
}

// generateCode returns a random 7-character uppercase alphanumeric code,
// using chars that are visually unambiguous (no O, I, 0, 1).
// resolvePlanClientScopeNotice substitutes the {{.Currency}} placeholder in
// the picker scope notice with the client's billing currency code (e.g.
// "PHP"). When the client has no billing currency, the placeholder + its
// surrounding parenthesis decoration are stripped so the sentence still reads
// cleanly: "Plans below match this client's billing currency."
func resolvePlanClientScopeNotice(template, currency string) string {
	if template == "" {
		return template
	}
	if currency == "" {
		// Strip "(<placeholder>)" or " (<placeholder>)" if present;
		// otherwise leave the template alone.
		template = strings.ReplaceAll(template, " ({{.Currency}})", "")
		return strings.ReplaceAll(template, "({{.Currency}})", "")
	}
	return strings.ReplaceAll(template, "{{.Currency}}", currency)
}

func generateCode() string {
	const chars = "23456789ABCDEFGHJKLMNPQRSTUVWXYZ"
	b := make([]byte, 7)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

// resolveClientBillingCurrency finds the billing_currency for a client by ID.
// Returns empty string when the client has no billing_currency set (caller should
// treat empty as "no currency filter" in the drawer's plan search).
func resolveClientBillingCurrency(ctx context.Context, clientID string, listClients func(ctx context.Context, req *clientpb.ListClientsRequest) (*clientpb.ListClientsResponse, error)) string {
	if clientID == "" || listClients == nil {
		return ""
	}
	resp, err := listClients(ctx, &clientpb.ListClientsRequest{})
	if err != nil {
		return ""
	}
	for _, c := range resp.GetData() {
		if c.GetId() == clientID {
			return c.GetBillingCurrency()
		}
	}
	return ""
}

// resolveClientLabel finds the display name for a client by ID.
func resolveClientLabel(ctx context.Context, clientID string, listClients func(ctx context.Context, req *clientpb.ListClientsRequest) (*clientpb.ListClientsResponse, error)) string {
	if clientID == "" || listClients == nil {
		return ""
	}
	resp, err := listClients(ctx, &clientpb.ListClientsRequest{})
	if err != nil {
		return clientID
	}
	for _, c := range resp.GetData() {
		if c.GetId() == clientID {
			if cn := c.GetName(); cn != "" {
				return cn
			}
			if u := c.GetUser(); u != nil {
				first := u.GetFirstName()
				last := u.GetLastName()
				if first != "" || last != "" {
					return strings.TrimSpace(first + " " + last)
				}
			}
			return clientID
		}
	}
	return clientID
}

// resolvePlanLabel finds the display name for a plan by ID.
func resolvePlanLabel(ctx context.Context, planID string, listPlans func(ctx context.Context, req *planpb.ListPlansRequest) (*planpb.ListPlansResponse, error)) string {
	if planID == "" || listPlans == nil {
		return ""
	}
	resp, err := listPlans(ctx, &planpb.ListPlansRequest{})
	if err != nil {
		return planID
	}
	for _, p := range resp.GetData() {
		if p.GetId() == planID {
			return p.GetName()
		}
	}
	return planID
}

// resolvePricePlanName looks up a PricePlan by ID and returns its display name.
// Prefers a single-row ReadPricePlan lookup over a full list scan.
// Falls back to the legacy ListPlans scan only when ReadPricePlan is nil or errors.
func resolvePricePlanName(ctx context.Context, pricePlanID string, deps *Deps) string {
	if pricePlanID == "" {
		return ""
	}
	// Single-row lookup is preferred over ListPricePlans for one-shot resolution.
	if deps.ReadPricePlan != nil {
		if resp, err := deps.ReadPricePlan(ctx, &priceplanpb.ReadPricePlanRequest{
			Data: &priceplanpb.PricePlan{Id: pricePlanID},
		}); err == nil && resp != nil && len(resp.GetData()) > 0 {
			pp := resp.GetData()[0]
			if name := pp.GetName(); name != "" {
				return name
			}
			if pl := pp.GetPlan(); pl != nil && pl.GetName() != "" {
				return pl.GetName()
			}
			if deps.ReadPlan != nil && pp.GetPlanId() != "" {
				planID := pp.GetPlanId()
				if rr, err := deps.ReadPlan(ctx, &planpb.ReadPlanRequest{Data: &planpb.Plan{Id: &planID}}); err == nil && len(rr.GetData()) > 0 {
					if n := rr.GetData()[0].GetName(); n != "" {
						return n
					}
				}
			}
			return pricePlanID
		}
	}
	// Last-resort fallback: legacy Plan list (handles cases where the
	// submitted ID is actually a plan_id rather than a price_plan_id).
	if deps.ListPlans != nil {
		if name := resolvePlanLabel(ctx, pricePlanID, deps.ListPlans); name != "" && name != pricePlanID {
			return name
		}
	}
	return pricePlanID
}

// splitTimestampForInputs renders ts in tz as a (date, time, RFC3339) triple
// suitable for the drawer's two-input grid + hidden ISO field. Nil ts → empties.
func splitTimestampForInputs(ts *timestamppb.Timestamp, tz *time.Location) (date, t, iso string) {
	if ts == nil || !ts.IsValid() {
		return "", "", ""
	}
	moment := ts.AsTime().In(tz)
	return moment.Format(pyezatypes.DateInputLayout), moment.Format(pyezatypes.TimeInputLayout), moment.Format(time.RFC3339)
}

// parseFormDateTime combines a date input ("2026-04-17"), a time input ("09:00"),
// and an explicit RFC3339 ISO string (set by JS with the chosen TZ offset) into
// a *timestamppb.Timestamp. The hidden ISO wins when present so the operator's
// chosen offset is preserved exactly. Falls back to date+time-in-tz when JS is
// disabled or the hidden field is empty. Empty all → nil.
//
// 2026-04-28 date+time field plan §4 — when no time is provided, isEnd
// switches the default between 00:00:00 (start) and 23:59:59 (end) so that an
// "end" date without a time still includes the full day.
func parseFormDateTime(date, t, iso string, tz *time.Location, isEnd bool) *timestamppb.Timestamp {
	if iso != "" {
		if parsed, err := time.Parse(time.RFC3339, iso); err == nil {
			return timestamppb.New(parsed.UTC())
		}
	}
	if date == "" {
		return nil
	}
	if t == "" {
		if isEnd {
			t = "23:59:59"
		} else {
			t = "00:00:00"
		}
	} else if len(t) == 5 {
		// Browser time inputs default to HH:MM precision; pad seconds.
		t = t + ":00"
	}
	parsed, err := time.ParseInLocation("2006-01-02 15:04:05", date+" "+t, tz)
	if err != nil {
		return nil
	}
	return timestamppb.New(parsed.UTC())
}

// formatDateForInput extracts a YYYY-MM-DD date string from either a date
// string (which may be a full timestamp or just a date) or a millisecond
// timestamp. If dateString is non-empty it is used directly — if it is longer
// than 10 chars the first 10 characters are returned (the date portion of an
// ISO 8601 timestamp). Falls back to dateMillis when dateString is empty;
// zero and negative millisecond values return "".
func formatDateForInput(dateString string, dateMillis int64) string {
	if dateString != "" {
		if len(dateString) > 10 {
			return dateString[:10]
		}
		return dateString
	}
	if dateMillis <= 0 {
		return ""
	}
	t := time.UnixMilli(dateMillis).UTC()
	return t.Format(pyezatypes.DateInputLayout)
}

// strPtr returns a pointer to a string.
func strPtr(s string) *string {
	return &s
}

// SpawnJobsDetection bundles the result of resolving a PricePlan → Plan →
// JobTemplate (+ JobTemplateRelation) chain for the Spawn Jobs section.
// 2026-04-29 auto-spawn-jobs-from-subscription plan §5.1.
type SpawnJobsDetection struct {
	Available     bool     // true when at least one JobTemplate resolves
	TemplateNames []string // root + child template names (display order)
	JobCount      int
	PhaseCount    int
	TaskCount     int
}

// detectSpawnJobs walks PricePlan → Plan → JobTemplate (+ relations) for the
// drawer's Spawn Jobs section. Returns Available=false when any link is
// missing or any read dep is unwired. Reads are best-effort — errors are
// swallowed and surface as Available=false (the section is hidden).
func detectSpawnJobs(ctx context.Context, deps *Deps, pricePlanID string) SpawnJobsDetection {
	out := SpawnJobsDetection{}
	if deps == nil || pricePlanID == "" || deps.ReadPricePlan == nil || deps.ReadPlan == nil || deps.ReadJobTemplate == nil {
		return out
	}
	ppResp, err := deps.ReadPricePlan(ctx, &priceplanpb.ReadPricePlanRequest{
		Data: &priceplanpb.PricePlan{Id: pricePlanID},
	})
	if err != nil || ppResp == nil || len(ppResp.GetData()) == 0 {
		return out
	}
	planID := ppResp.GetData()[0].GetPlanId()
	if planID == "" {
		return out
	}
	planResp, err := deps.ReadPlan(ctx, &planpb.ReadPlanRequest{Data: &planpb.Plan{Id: &planID}})
	if err != nil || planResp == nil || len(planResp.GetData()) == 0 {
		return out
	}
	rootTemplateID := planResp.GetData()[0].GetJobTemplateId()
	if rootTemplateID == "" {
		return out
	}

	// Collect root + active children via JobTemplateRelation.
	templateIDs := []string{rootTemplateID}
	if deps.ListJobTemplateRelations != nil {
		relResp, err := deps.ListJobTemplateRelations(ctx, &jobtemplaterelationpb.ListJobTemplateRelationsByParentRequest{
			ParentTemplateId: rootTemplateID,
		})
		if err == nil && relResp != nil {
			for _, rel := range relResp.GetJobTemplateRelations() {
				if !rel.GetActive() {
					continue
				}
				cid := rel.GetChildTemplateId()
				if cid != "" && cid != rootTemplateID {
					templateIDs = append(templateIDs, cid)
				}
			}
		}
	}

	for _, tid := range templateIDs {
		tplResp, err := deps.ReadJobTemplate(ctx, &jobtemplatepb.ReadJobTemplateRequest{
			Data: &jobtemplatepb.JobTemplate{Id: tid},
		})
		if err != nil || tplResp == nil || len(tplResp.GetData()) == 0 {
			continue
		}
		tpl := tplResp.GetData()[0]
		if !tpl.GetActive() {
			continue
		}
		out.TemplateNames = append(out.TemplateNames, tpl.GetName())
		out.JobCount++
		if deps.ListJobTemplatePhases == nil {
			continue
		}
		phaseResp, err := deps.ListJobTemplatePhases(ctx, &jobtemplatephasepb.ListByJobTemplateRequest{JobTemplateId: tid})
		if err != nil || phaseResp == nil {
			continue
		}
		phases := phaseResp.GetJobTemplatePhases()
		out.PhaseCount += len(phases)
		if deps.ListJobTemplateTasks == nil {
			continue
		}
		for _, ph := range phases {
			tasksResp, err := deps.ListJobTemplateTasks(ctx, &jobtemplatetaskpb.ListJobTemplateTasksByPhaseRequest{
				JobTemplatePhaseId: ph.GetId(),
			})
			if err != nil || tasksResp == nil {
				continue
			}
			out.TaskCount += len(tasksResp.GetJobTemplateTasks())
		}
	}

	out.Available = out.JobCount > 0
	return out
}

// resolveSpawnJobsSummary renders the SpawnJobsSummary template
// ("Spawning {{.JobCount}} Job(s) from {{.TemplateNames}} — includes
// {{.PhaseCount}} phases, {{.TaskCount}} tasks.") with the detected counts.
// Falls back to an empty string when the template is empty.
func resolveSpawnJobsSummary(template string, det SpawnJobsDetection) string {
	if template == "" || !det.Available {
		return ""
	}
	names := strings.Join(det.TemplateNames, ", ")
	r := strings.NewReplacer(
		"{{.JobCount}}", itoa(det.JobCount),
		"{{.TemplateNames}}", names,
		"{{.PhaseCount}}", itoa(det.PhaseCount),
		"{{.TaskCount}}", itoa(det.TaskCount),
	)
	return r.Replace(template)
}

func itoa(n int) string {
	// fmt.Sprintf is fine but avoids importing fmt at this scope.
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var b [20]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		b[i] = '-'
	}
	return string(b[i:])
}

// resolveSubscriptionLabel returns a "code · name" label for the drawer header.
// Shared by spawn_jobs and spawn_cycle_jobs handlers.
func resolveSubscriptionLabel(ctx context.Context, deps *Deps, subscriptionID string) string {
	if deps == nil || deps.ReadSubscription == nil {
		return subscriptionID
	}
	resp, err := deps.ReadSubscription(ctx, &subscriptionpb.ReadSubscriptionRequest{
		Data: &subscriptionpb.Subscription{Id: subscriptionID},
	})
	if err != nil || resp == nil || len(resp.GetData()) == 0 {
		return subscriptionID
	}
	s := resp.GetData()[0]
	if c := s.GetCode(); c != "" {
		if n := s.GetName(); n != "" {
			return c + " · " + n
		}
		return c
	}
	if n := s.GetName(); n != "" {
		return n
	}
	return subscriptionID
}
