package action

// This file provides JSON search handlers for the auto-complete component.
// They accept ?q=searchterm and return JSON: [{"value":"id","label":"Name"}, ...]

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	clientpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client"
	subscriptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription"
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

// NewSearchSubscriptionsAction returns an http.HandlerFunc that searches active
// subscriptions filtered by client_id and returns JSON results for the auto-complete
// component. Requires ?client_id= query param; returns [] if missing.
func NewSearchSubscriptionsAction(deps *Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		clientID := r.URL.Query().Get("client_id")
		if clientID == "" {
			writeJSON(w, []searchOption{})
			return
		}

		if deps.ListSubscriptions == nil {
			writeJSON(w, []searchOption{})
			return
		}

		query := strings.TrimSpace(r.URL.Query().Get("q"))
		queryLower := strings.ToLower(query)

		resp, err := deps.ListSubscriptions(ctx, &subscriptionpb.ListSubscriptionsRequest{})
		if err != nil {
			log.Printf("search subscriptions: failed to list subscriptions: %v", err)
			writeJSON(w, []searchOption{})
			return
		}

		var results []searchOption
		for _, s := range resp.GetData() {
			if !s.GetActive() {
				continue
			}
			if s.GetClientId() != clientID {
				continue
			}
			name := s.GetName()
			if queryLower != "" && !strings.Contains(strings.ToLower(name), queryLower) {
				continue
			}
			results = append(results, searchOption{
				Value: s.GetId(),
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
