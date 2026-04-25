package action

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

// searchOption is the JSON shape returned by the search handlers.
type searchOption struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

// groupedSearchResult is the grouped JSON shape for the price plan auto-complete.
type groupedSearchResult struct {
	Group   string         `json:"group"`
	Options []searchOption `json:"options"`
}

const searchResultLimit = 20

// NewSearchClientsAction returns an http.HandlerFunc that searches clients
// by company_name, user first_name, or last_name and returns JSON results
// for the auto-complete component.
func NewSearchClientsAction(deps *Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		query := strings.TrimSpace(r.URL.Query().Get("q"))

		// Use proto search if available (SQL ILIKE, no full load)
		if deps.SearchClientsByName != nil {
			resp, err := deps.SearchClientsByName(ctx, &clientpb.SearchClientsByNameRequest{
				Query: query,
			})
			if err != nil {
				log.Printf("search clients: failed to search clients by name: %v", err)
				writeJSON(w, []searchOption{})
				return
			}
			var results []searchOption
			for _, r := range resp.GetResults() {
				results = append(results, searchOption{
					Value: r.GetId(),
					Label: r.GetLabel(),
				})
			}
			if results == nil {
				results = []searchOption{}
			}
			writeJSON(w, results)
			return
		}

		// Fallback: load all clients and filter in Go
		if deps.ListClients == nil {
			writeJSON(w, []searchOption{})
			return
		}

		queryLower := strings.ToLower(query)
		resp, err := deps.ListClients(ctx, &clientpb.ListClientsRequest{})
		if err != nil {
			log.Printf("search clients: failed to list clients: %v", err)
			writeJSON(w, []searchOption{})
			return
		}

		var results []searchOption
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

			results = append(results, searchOption{
				Value: c.GetId(),
				Label: label,
			})

			if len(results) >= searchResultLimit {
				break
			}
		}

		if results == nil {
			results = []searchOption{}
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
		// billing_currency filters results to PricePlans matching the client's
		// billing currency. Empty = no filter (show all currencies).
		billingCurrency := strings.ToUpper(strings.TrimSpace(r.URL.Query().Get("billing_currency")))

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
		_ = tz // tz only consumed below for group-label formatting

		if deps.ListPricePlans == nil {
			// Fallback to plan search if ListPricePlans not wired
			if deps.ListPlans != nil {
				searchPlansLegacy(ctx, w, q, deps)
			} else {
				writeJSON(w, []searchOption{})
			}
			return
		}

		resp, err := deps.ListPricePlans(ctx, &priceplanpb.ListPricePlansRequest{})
		if err != nil {
			log.Printf("search price plans: failed to list price plans: %v", err)
			writeJSON(w, []groupedSearchResult{})
			return
		}

		// Build schedule lookup map when ListPriceSchedules is wired.
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

		// Per-plan name cache for the embedded-plan fallback. ListPricePlans returns
		// flat rows without the joined plan, so we resolve names lazily via ReadPlan
		// and cache the result to avoid repeating the lookup for shared plan_ids.
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

		// group key → options (use insertion-stable order via a slice of keys)
		type groupEntry struct {
			schedID   string // empty = unscheduled
			schedName string // display label for the group header
			dateStart string // for sorting
			options   []searchOption
		}
		groupMap := map[string]*groupEntry{}
		var groupOrder []string // tracks insertion order of group keys

		totalOptions := 0

		for _, pp := range resp.GetData() {
			if !pp.GetActive() {
				continue
			}
			if totalOptions >= searchResultLimit {
				break
			}

			// Resolve display name: prefer pp.GetName(), fall back to embedded plan name,
			// then to the joined plan-name lookup (since ListPricePlans does not embed the plan),
			// then to a placeholder so the option is never rendered with an empty label.
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

			// Apply query filter.
			if queryLower != "" && !strings.Contains(strings.ToLower(displayName), queryLower) {
				continue
			}

			// Apply billing currency filter (when the drawer passes a client's billing_currency).
			if billingCurrency != "" && strings.ToUpper(pp.GetBillingCurrency()) != billingCurrency {
				continue
			}

			schedID := pp.GetPriceScheduleId()
			var sched *priceschedulepb.PriceSchedule
			if schedID != "" {
				sched = scheduleByID[schedID]
			}

			// Apply date filter (UTC timestamp comparison).
			if hasDateFilter {
				if schedID == "" {
					// Unscheduled plans are excluded when any date filter is set.
					continue
				}
				if sched == nil {
					// Unknown schedule referenced — skip.
					continue
				}
				schedStartTS := sched.GetDateTimeStart()
				schedEndTS := sched.GetDateTimeEnd()
				if schedStartTS == nil {
					continue
				}
				schedStart := schedStartTS.AsTime()
				// Schedule must cover the requested range:
				// schedStart <= reqStart AND (schedEnd == nil || reqEnd <= schedEnd)
				if schedStart.After(reqStart) {
					continue
				}
				if schedEndTS != nil {
					if reqEnd.After(schedEndTS.AsTime()) {
						continue
					}
				}
			}

			// Determine group key and label.
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

			// Build option label: "Name · ₱15,000.00 PHP"
			amount := float64(pp.GetBillingAmount()) / 100.0
			currency := pp.GetBillingCurrency()
			var label string
			if currency != "" {
				label = fmt.Sprintf("%s · ₱%s %s", displayName, formatAmount(amount), currency)
			} else {
				label = fmt.Sprintf("%s · ₱%s", displayName, formatAmount(amount))
			}

			groupMap[groupKey].options = append(groupMap[groupKey].options, searchOption{
				Value: pp.GetId(),
				Label: label,
			})
			totalOptions++
		}

		// Sort groups: by dateStart ascending then schedName; unscheduled last.
		sort.SliceStable(groupOrder, func(i, j int) bool {
			a := groupMap[groupOrder[i]]
			b := groupMap[groupOrder[j]]
			if a.schedID == "" && b.schedID != "" {
				return false // unscheduled always last
			}
			if a.schedID != "" && b.schedID == "" {
				return true
			}
			if a.dateStart != b.dateStart {
				return a.dateStart < b.dateStart
			}
			return a.schedName < b.schedName
		})

		// Sort options within each group by display name ascending.
		for _, entry := range groupMap {
			sort.SliceStable(entry.options, func(i, j int) bool {
				return entry.options[i].Label < entry.options[j].Label
			})
		}

		// Build result slice.
		results := make([]groupedSearchResult, 0, len(groupOrder))
		for _, key := range groupOrder {
			entry := groupMap[key]
			results = append(results, groupedSearchResult{
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
		writeJSON(w, []searchOption{})
		return
	}

	var results []searchOption
	for _, p := range resp.GetData() {
		if !p.GetActive() {
			continue
		}
		name := p.GetName()
		if queryLower != "" && !strings.Contains(strings.ToLower(name), queryLower) {
			continue
		}
		results = append(results, searchOption{
			Value: p.GetId(),
			Label: name,
		})
		if len(results) >= searchResultLimit {
			break
		}
	}
	if results == nil {
		results = []searchOption{}
	}
	writeJSON(w, results)
}

// formatAmount formats a float amount with thousands separators and 2 decimal places.
func formatAmount(amount float64) string {
	s := fmt.Sprintf("%.2f", amount)
	// Add thousands separators
	parts := strings.Split(s, ".")
	intPart := parts[0]
	decPart := parts[1]

	// Insert commas
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

// writeJSON marshals data as JSON and writes it to the response writer.
func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("search: failed to encode JSON response: %v", err)
	}
}

// parseRFC3339 returns the parsed UTC instant from an RFC3339 string or
// (zero, false) when empty/invalid.
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
