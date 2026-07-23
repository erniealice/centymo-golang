package subscription_group

import "strings"

// Options — app-configurable presentation for the subscription_group roster
// (the section enrollments tab), set by the consuming app through the block's
// EngineBlock option. Generic reference grammar with a fail-safe contract: an
// empty field disables that behavior and an unrecognized reference is ignored —
// the roster renders exactly as without it, never an error. The zero value is
// today's flat roster (apps that set no options are unaffected).
//
// Reference forms (all generic — no vertical nouns):
//   - attribute ref     "client_attributes.<code>"   (attribute.code)
//   - client-column ref "last_name"
//   - direction         "asc" | "desc"               (default "asc")
type Options struct {
	// Roster configures the enrollments-tab row presentation (group bands + sort).
	Roster RowOptions
}

// RowOptions — roster row presentation. GroupByField partitions rows into value
// bands ("client_attributes.<code>"); SortField orders rows within a band (a
// client-column ref, e.g. "last_name"); SortDirection is "asc"|"desc" (default
// asc).
type RowOptions struct {
	GroupByField string
	// GroupValueOrder pins the band order to these attribute values
	// (case-insensitive; listed values lead in list order). Values not listed
	// follow in ascending value order; the no-value band stays last. Empty =
	// ascending value order (the fail-safe default).
	GroupValueOrder []string
	SortField       string
	SortDirection   string
}

// ClientAttributeFieldPrefix is the reference-form prefix for the client-
// attribute family.
const ClientAttributeFieldPrefix = "client_attributes."

// ClientAttributeCode extracts <code> from a "client_attributes.<code>"
// reference. ok is false for an empty or foreign-form reference.
func ClientAttributeCode(field string) (code string, ok bool) {
	rest, found := strings.CutPrefix(field, ClientAttributeFieldPrefix)
	if !found || rest == "" {
		return "", false
	}
	return rest, true
}

// GroupByAttributeCode returns the group-by attribute code and whether banding
// is configured (GroupByField is a "client_attributes.<code>" reference). A
// false result means the flat roster (the fail-safe default).
func (r RowOptions) GroupByAttributeCode() (code string, ok bool) {
	return ClientAttributeCode(r.GroupByField)
}

// SortByLastName reports whether rows sort by the client's last name within a
// band (the single implemented SortField). Any other value keeps member order.
func (r RowOptions) SortByLastName() bool {
	return strings.TrimSpace(r.SortField) == "last_name"
}

// Direction returns the normalized row sort direction ("asc"|"desc"), defaulting
// to "asc" for an empty or unrecognized value.
func (r RowOptions) Direction() string {
	if strings.EqualFold(strings.TrimSpace(r.SortDirection), "desc") {
		return "desc"
	}
	return "asc"
}

// AttributeCodes returns the distinct client-attribute codes referenced by the
// Roster options (group_by, then sort), in first-use order. SortField is
// normally a plain client column ("last_name"), so this usually yields just the
// group_by code.
func (o Options) AttributeCodes() []string {
	var out []string
	seen := map[string]bool{}
	for _, f := range []string{o.Roster.GroupByField, o.Roster.SortField} {
		if code, ok := ClientAttributeCode(f); ok && !seen[code] {
			seen[code] = true
			out = append(out, code)
		}
	}
	return out
}
