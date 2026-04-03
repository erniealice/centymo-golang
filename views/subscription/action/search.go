package action

// This file provides JSON search handlers for the auto-complete component.
// They accept ?q=searchterm and return JSON: [{"value":"id","label":"Name"}, ...]

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	clientpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client"
	planpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/plan"
	priceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_plan"
)

// searchOption is the JSON shape returned by the search handlers.
type searchOption struct {
	Value string `json:"value"`
	Label string `json:"label"`
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
// and returns JSON results for the auto-complete component.
// Returns price_plan.id as value with label showing name + formatted amount + currency.
func NewSearchPlansAction(deps *Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		query := strings.TrimSpace(r.URL.Query().Get("q"))

		if deps.ListPricePlans == nil {
			// Fallback to plan search if ListPricePlans not wired
			if deps.ListPlans != nil {
				searchPlansLegacy(ctx, w, query, deps)
			} else {
				writeJSON(w, []searchOption{})
			}
			return
		}

		resp, err := deps.ListPricePlans(ctx, &priceplanpb.ListPricePlansRequest{})
		if err != nil {
			log.Printf("search price plans: failed to list price plans: %v", err)
			writeJSON(w, []searchOption{})
			return
		}

		queryLower := strings.ToLower(query)
		var results []searchOption
		for _, pp := range resp.GetData() {
			if !pp.GetActive() {
				continue
			}

			name := pp.GetName()
			if queryLower != "" && !strings.Contains(strings.ToLower(name), queryLower) {
				continue
			}

			// Format: "Plan Name — 15,000.00 PHP"
			amount := pp.GetAmount() / 100
			currency := pp.GetCurrency()
			label := fmt.Sprintf("%s — %s %s", name, formatAmount(amount), currency)

			results = append(results, searchOption{
				Value: pp.GetId(),
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
