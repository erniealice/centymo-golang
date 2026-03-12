package action

// This file provides JSON search handlers for the auto-complete component.
// They accept ?q=searchterm and return JSON: [{"value":"id","label":"Name"}, ...]

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	clientpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client"
	planpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/plan"
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
		query := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("q")))

		if deps.ListClients == nil {
			writeJSON(w, []searchOption{})
			return
		}

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

			// Build label: prefer company_name, fallback to user name, fallback to ID
			label := c.GetId()
			companyName := c.GetCompanyName()
			if companyName != "" {
				label = companyName
			} else if u := c.GetUser(); u != nil {
				first := u.GetFirstName()
				last := u.GetLastName()
				if first != "" || last != "" {
					label = strings.TrimSpace(first + " " + last)
				}
			}

			// Filter by query (if provided)
			if query != "" {
				labelLower := strings.ToLower(label)
				match := strings.Contains(labelLower, query)
				// Also check individual name parts for broader matching
				if !match {
					if cn := c.GetCompanyName(); cn != "" {
						match = strings.Contains(strings.ToLower(cn), query)
					}
				}
				if !match {
					if u := c.GetUser(); u != nil {
						match = strings.Contains(strings.ToLower(u.GetFirstName()), query) ||
							strings.Contains(strings.ToLower(u.GetLastName()), query)
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

// NewSearchPlansAction returns an http.HandlerFunc that searches plans
// by name and returns JSON results for the auto-complete component.
func NewSearchPlansAction(deps *Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		query := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("q")))

		if deps.ListPlans == nil {
			writeJSON(w, []searchOption{})
			return
		}

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

			// Filter by query (if provided)
			if query != "" && !strings.Contains(strings.ToLower(name), query) {
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
}

// writeJSON marshals data as JSON and writes it to the response writer.
func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("search: failed to encode JSON response: %v", err)
	}
}
