// Package search provides HTMX search-as-you-type endpoints for subscriptions.
// Extracted per S1 pattern — 476 LOC exceeds the ~150 LOC trigger threshold.
package search

// This file provides JSON search handlers for the auto-complete component.
// They accept ?q=searchterm and return JSON: [{"value":"id","label":"Name"}, ...]

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	pyezatypes "github.com/erniealice/pyeza-golang/types"

	clientpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client"
	planpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/plan"
	priceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_plan"
	priceschedulepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_schedule"
)

// Deps is the dependency subset needed by the search handlers.
type Deps struct {
	SearchClientsByName func(ctx context.Context, req *clientpb.SearchClientsByNameRequest) (*clientpb.SearchClientsByNameResponse, error)
	ListClients         func(ctx context.Context, req *clientpb.ListClientsRequest) (*clientpb.ListClientsResponse, error)
	ListPricePlans      func(ctx context.Context, req *priceplanpb.ListPricePlansRequest) (*priceplanpb.ListPricePlansResponse, error)
	ListPriceSchedules  func(ctx context.Context, req *priceschedulepb.ListPriceSchedulesRequest) (*priceschedulepb.ListPriceSchedulesResponse, error)
	ListPlans           func(ctx context.Context, req *planpb.ListPlansRequest) (*planpb.ListPlansResponse, error)
	ReadPlan            func(ctx context.Context, req *planpb.ReadPlanRequest) (*planpb.ReadPlanResponse, error)

	// PlanGroupForClient is the lyngua-resolved label template for the
	// client-scoped group header, e.g. "For {{.ClientName}}".
	// PlanGroupGeneral is the fixed label for the general group, e.g.
	// "General packages". Both are sourced from SubscriptionLabels.Form.
	PlanGroupForClient string
	PlanGroupGeneral   string
}

// option is the JSON shape returned by the search handlers.
type option struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

// groupedResult is the grouped JSON shape for the price plan auto-complete.
type groupedResult struct {
	Group   string   `json:"group"`
	Options []option `json:"options"`
}

const searchResultLimit = 20

// NewSearchClientsAction returns an http.HandlerFunc that searches clients
// by company_name, user first_name, or last_name and returns JSON results
// for the auto-complete component.
func NewSearchClientsAction(deps *Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		query := strings.TrimSpace(r.URL.Query().Get("q"))

		if deps.SearchClientsByName != nil {
			resp, err := deps.SearchClientsByName(ctx, &clientpb.SearchClientsByNameRequest{
				Query: query,
			})
			if err != nil {
				log.Printf("search clients: failed to search clients by name: %v", err)
				writeJSON(w, []option{})
				return
			}
			var results []option
			for _, r := range resp.GetResults() {
				results = append(results, option{
					Value: r.GetId(),
					Label: r.GetLabel(),
				})
			}
			if results == nil {
				results = []option{}
			}
			writeJSON(w, results)
			return
		}

		if deps.ListClients == nil {
			writeJSON(w, []option{})
			return
		}

		queryLower := strings.ToLower(query)
		resp, err := deps.ListClients(ctx, &clientpb.ListClientsRequest{})
		if err != nil {
			log.Printf("search clients: failed to list clients: %v", err)
			writeJSON(w, []option{})
			return
		}

		var results []option
		for _, c := range resp.GetData() {
			if !c.GetActive() {
				continue
			}

			label := c.GetId()
			companyName := c.GetName()
			if companyName != "" {
				label = companyName
			} else if u := c.GetUser(); u != nil {
				first := u.GetFirstName()
				last := u.GetLastName()
				if first != "" || last != "" {
					label = strings.TrimSpace(first + " " + last)
				}
			}

			if queryLower != "" {
				labelLower := strings.ToLower(label)
				match := strings.Contains(labelLower, queryLower)
				if !match {
					if cn := c.GetName(); cn != "" {
						match = strings.Contains(strings.ToLower(cn), queryLower)
					}
				}
				if !match {
					if u := c.GetUser(); u != nil {
						match = strings.Contains(strings.ToLower(u.GetFirstName()), queryLower) ||
							strings.Contains(strings.ToLower(u.GetLastName()), queryLower)
					}
				}
				if !match {
					continue
				}
			}

			results = append(results, option{
				Value: c.GetId(),
				Label: label,
			})

			if len(results) >= searchResultLimit {
				break
			}
		}

		if results == nil {
			results = []option{}
		}
		writeJSON(w, results)
	}
}

// NewSearchPlansAction returns an http.HandlerFunc that searches price plans
// and returns grouped JSON results for the auto-complete component.
// Response shape: [{"group":"Schedule Name","options":[{"value":"id","label":"..."}]}]
// Falls back to flat [{"value","label"}] via searchPlansLegacy when ListPricePlans is nil.
func NewSearchPlansAction(deps *Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		q := strings.TrimSpace(r.URL.Query().Get("q"))
		startISO := strings.TrimSpace(r.URL.Query().Get("date_time_start_iso"))
		endISO := strings.TrimSpace(r.URL.Query().Get("date_time_end_iso"))
		billingCurrency := strings.ToUpper(strings.TrimSpace(r.URL.Query().Get("billing_currency")))
		clientIDFilter := strings.TrimSpace(r.URL.Query().Get("client_id"))
		clientNameParam := strings.TrimSpace(r.URL.Query().Get("client_name"))

		tz := pyezatypes.LocationFromContext(ctx)
		reqStart, hasStart := parseRFC3339(startISO)
		reqEnd, hasEnd := parseRFC3339(endISO)
		hasDateFilter := hasStart || hasEnd
		if hasDateFilter {
			if !hasStart {
				reqStart = reqEnd
			}
			if !hasEnd {
				reqEnd = reqStart
			}
		}
		_ = tz

		if deps.ListPricePlans == nil {
			if deps.ListPlans != nil {
				searchPlansLegacy(ctx, w, q, deps)
			} else {
				writeJSON(w, []option{})
			}
			return
		}

		resp, err := deps.ListPricePlans(ctx, &priceplanpb.ListPricePlansRequest{})
		if err != nil {
			log.Printf("search price plans: failed to list price plans: %v", err)
			writeJSON(w, []groupedResult{})
			return
		}

		scheduleByID := map[string]*priceschedulepb.PriceSchedule{}
		if deps.ListPriceSchedules != nil {
			schedResp, schedErr := deps.ListPriceSchedules(ctx, &priceschedulepb.ListPriceSchedulesRequest{})
			if schedErr != nil {
				log.Printf("search price plans: failed to list price schedules: %v", schedErr)
			} else {
				for _, s := range schedResp.GetData() {
					scheduleByID[s.GetId()] = s
				}
			}
		}

		planNameByID := map[string]string{}
		resolvePlanName := func(planID string) string {
			if planID == "" || deps.ReadPlan == nil {
				return ""
			}
			if name, ok := planNameByID[planID]; ok {
				return name
			}
			id := planID
			readResp, readErr := deps.ReadPlan(ctx, &planpb.ReadPlanRequest{Data: &planpb.Plan{Id: &id}})
			if readErr != nil || len(readResp.GetData()) == 0 {
				planNameByID[planID] = ""
				return ""
			}
			name := readResp.GetData()[0].GetName()
			planNameByID[planID] = name
			return name
		}

		queryLower := strings.ToLower(q)
		_ = clientNameParam // reserved for future per-client group header (currently unused — grouping is per-schedule across both tiers).

		// 2026-05-03 (revised) — Per-schedule grouping for both standalone
		// (no client) and client-context callers. When client_id is present we
		// still apply a cross-client rejection filter so a *different* client's
		// scoped plans never appear, but client-scoped AND general-scope plans
		// for the chosen client both surface, grouped under their parent
		// PriceSchedule. See discussion 2026-05-03: a client engagement may
		// legitimately attach a general-scope plan, so we don't drop those.

		type groupEntry struct {
			schedID   string
			schedName string
			dateStart string
			options   []option
		}
		groupMap := map[string]*groupEntry{}
		var groupOrder []string

		totalOptions := 0

		for _, pp := range resp.GetData() {
			if !pp.GetActive() {
				continue
			}
			if totalOptions >= searchResultLimit {
				break
			}

			displayName := pp.GetName()
			if displayName == "" {
				if pl := pp.GetPlan(); pl != nil {
					displayName = pl.GetName()
				}
			}
			if displayName == "" {
				displayName = resolvePlanName(pp.GetPlanId())
			}
			if displayName == "" {
				displayName = "(Unnamed plan)"
			}

			if queryLower != "" && !strings.Contains(strings.ToLower(displayName), queryLower) {
				continue
			}

			if billingCurrency != "" && strings.ToUpper(pp.GetBillingCurrency()) != billingCurrency {
				continue
			}

			schedID := pp.GetPriceScheduleId()
			var sched *priceschedulepb.PriceSchedule
			if schedID != "" {
				sched = scheduleByID[schedID]
			}

			// 2026-05-03 — Apply cross-client rejection when clientIDFilter
			// is set. A plan scoped to a *different* client (directly or via
			// its parent schedule) is never selectable in this drawer.
			// General-scope plans pass through.
			if clientIDFilter != "" {
				ppClient := pp.GetClientId()
				schedClient := ""
				if sched != nil {
					schedClient = sched.GetClientId()
				}
				if (ppClient != "" && ppClient != clientIDFilter) || (schedClient != "" && schedClient != clientIDFilter) {
					continue
				}
			}

			// Apply date filter.
			if hasDateFilter {
				if schedID == "" {
					continue
				}
				if sched == nil {
					continue
				}
				schedStartTS := sched.GetDateTimeStart()
				schedEndTS := sched.GetDateTimeEnd()
				if schedStartTS == nil {
					continue
				}
				schedStart := schedStartTS.AsTime()
				if schedStart.After(reqStart) {
					continue
				}
				if schedEndTS != nil {
					if reqEnd.After(schedEndTS.AsTime()) {
						continue
					}
				}
			}

			groupKey := schedID
			if groupKey == "" {
				groupKey = "__unscheduled__"
			}

			if _, exists := groupMap[groupKey]; !exists {
				entry := &groupEntry{schedID: schedID}
				if sched != nil {
					schedStart := pyezatypes.FormatTimestampInTZ(sched.GetDateTimeStart(), tz, pyezatypes.DateInputLayout)
					schedEnd := pyezatypes.FormatTimestampInTZ(sched.GetDateTimeEnd(), tz, pyezatypes.DateInputLayout)
					entry.dateStart = schedStart
					name := sched.GetName()
					if name == "" {
						name = schedStart
						if schedEnd != "" {
							name += " → " + schedEnd
						}
					}
					entry.schedName = name
				} else {
					entry.schedName = "Unscheduled"
				}
				groupMap[groupKey] = entry
				groupOrder = append(groupOrder, groupKey)
			}

			amount := float64(pp.GetBillingAmount()) / 100.0
			currency := pp.GetBillingCurrency()
			var label string
			if currency != "" {
				label = fmt.Sprintf("%s · ₱%s %s", displayName, formatAmount(amount), currency)
			} else {
				label = fmt.Sprintf("%s · ₱%s", displayName, formatAmount(amount))
			}

			groupMap[groupKey].options = append(groupMap[groupKey].options, option{
				Value: pp.GetId(),
				Label: label,
			})
			totalOptions++
		}

		sort.SliceStable(groupOrder, func(i, j int) bool {
			a := groupMap[groupOrder[i]]
			b := groupMap[groupOrder[j]]
			if a.schedID == "" && b.schedID != "" {
				return false
			}
			if a.schedID != "" && b.schedID == "" {
				return true
			}
			if a.dateStart != b.dateStart {
				return a.dateStart < b.dateStart
			}
			return a.schedName < b.schedName
		})

		for _, entry := range groupMap {
			sort.SliceStable(entry.options, func(i, j int) bool {
				return entry.options[i].Label < entry.options[j].Label
			})
		}

		results := make([]groupedResult, 0, len(groupOrder))
		for _, key := range groupOrder {
			entry := groupMap[key]
			results = append(results, groupedResult{
				Group:   entry.schedName,
				Options: entry.options,
			})
		}

		writeJSON(w, results)
	}
}

// searchPlansLegacy is the old plan-based search fallback.
func searchPlansLegacy(ctx context.Context, w http.ResponseWriter, query string, deps *Deps) {
	queryLower := strings.ToLower(query)
	resp, err := deps.ListPlans(ctx, &planpb.ListPlansRequest{})
	if err != nil {
		log.Printf("search plans: failed to list plans: %v", err)
		writeJSON(w, []option{})
		return
	}

	var results []option
	for _, p := range resp.GetData() {
		if !p.GetActive() {
			continue
		}
		name := p.GetName()
		if queryLower != "" && !strings.Contains(strings.ToLower(name), queryLower) {
			continue
		}
		results = append(results, option{
			Value: p.GetId(),
			Label: name,
		})
		if len(results) >= searchResultLimit {
			break
		}
	}
	if results == nil {
		results = []option{}
	}
	writeJSON(w, results)
}

func formatAmount(amount float64) string {
	s := fmt.Sprintf("%.2f", amount)
	parts := strings.Split(s, ".")
	intPart := parts[0]
	decPart := parts[1]

	n := len(intPart)
	if n <= 3 {
		return intPart + "." + decPart
	}

	var result []byte
	for i, c := range intPart {
		if i > 0 && (n-i)%3 == 0 {
			result = append(result, ',')
		}
		result = append(result, byte(c))
	}
	return string(result) + "." + decPart
}

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("search: failed to encode JSON response: %v", err)
	}
}

func parseRFC3339(iso string) (time.Time, bool) {
	if iso == "" {
		return time.Time{}, false
	}
	t, err := time.Parse(time.RFC3339, iso)
	if err != nil {
		return time.Time{}, false
	}
	return t.UTC(), true
}
