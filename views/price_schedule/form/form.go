package form

import (
	centymo "github.com/erniealice/centymo-golang"
	"strings"
)

// LocationOption is a location entry for the drawer's location picker.
type LocationOption struct {
	Id   string
	Name string
}

// Data is the template data for the price_schedule drawer form.
type Data struct {
	FormAction string
	IsEdit     bool
	ID         string
	Name       string
	Description string
	// Date + optional time inputs (2026-04-28 date+time field plan).
	// The drawer renders <input type="date"> + <input type="time"> side
	// by side; time is OPTIONAL. Empty time defaults to 00:00:00 for
	// DateStart and 23:59:59 for DateEnd so an end-only date covers the
	// full day. parseScheduleDateTime() applies the rule on POST.
	DateStartDate         string
	DateStartTime         string
	DateEndDate           string
	DateEndTime           string
	Active                bool
	Locations             []*LocationOption
	SelectedLocationID    string
	SelectedLocationLabel string
	LocationOptions       []map[string]any

	// 2026-04-27 plan-client-scope plan §6.7 / §4.4.1.
	ClientID        string
	ClientLabel     string
	ClientOptions   []map[string]any
	SearchClientURL string
	// SuggestNameURL is the GET endpoint that the Client picker hits via
	// HTMX to refresh the Name input with the per-tier derived name
	// "{ClientName} - {customClientPriceScheduleLabelSuffix}".
	SuggestNameURL string

	// Scope (2026-04-28) — "location" or "client". Drives the radio that
	// mutually excludes the Location and Client pickers. Default "location"
	// for new schedules; for edit, derived from record.client_id presence.
	Scope string

	// 2026-05-03 — When the drawer is opened from a specific client's detail
	// page (?client_id=...), ClientLocked is true. The template hides the
	// scope radio and renders the Client field as a read-only display row
	// with a hidden client_id input — the Client picker auto-complete is
	// suppressed entirely.
	ClientLocked bool

	Labels       centymo.PriceScheduleFormLabels
	CommonLabels any
}

// BuildLocationAutoCompleteOptions converts a slice of LocationOption into the
// map shape expected by the auto-complete component.
func BuildLocationAutoCompleteOptions(locations []*LocationOption, selectedID string) []map[string]any {
	opts := make([]map[string]any, 0, len(locations))
	for _, loc := range locations {
		opts = append(opts, map[string]any{
			"Value":    loc.Id,
			"Label":    loc.Name,
			"Selected": loc.Id == selectedID,
		})
	}
	return opts
}

// FindLocationLabel returns the display name of the location with the given ID,
// or empty string when not found.
func FindLocationLabel(locations []*LocationOption, id string) string {
	for _, loc := range locations {
		if loc.Id == id {
			return loc.Name
		}
	}
	return ""
}

// BuildDerivedScheduleName produces "{ClientName} - {suffix} - {timestamp}"
// per plan §4.4.1, with a 2026-05-03 trailing timestamp suffix to avoid
// duplicate names when an operator creates multiple schedules for the same
// client (or repeatedly adds master schedules from the standalone list).
// Any of the three components may be empty; the function joins only the
// non-empty parts with " - ". An entirely-empty input returns "".
func BuildDerivedScheduleName(clientName, suffix, timestamp string) string {
	parts := []string{}
	if v := strings.TrimSpace(clientName); v != "" {
		parts = append(parts, v)
	}
	if v := strings.TrimSpace(suffix); v != "" {
		parts = append(parts, v)
	}
	if v := strings.TrimSpace(timestamp); v != "" {
		parts = append(parts, v)
	}
	return strings.Join(parts, " - ")
}
