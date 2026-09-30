package block

import (
	"fmt"
	"sort"
	"strings"
)

// requireUnit is the one centymo Mount/RequireFor completeness helper (R4 rule
// 6b, R5 M7): it returns an error naming every missing closure in sorted
// order, or nil. Each unit builds one require<Unit>(uc) error on top of it.
// checks maps a closure name to "is wired".
func requireUnit(unit string, checks map[string]bool) error {
	var missing []string
	for name, ok := range checks {
		if !ok {
			missing = append(missing, name)
		}
	}
	if len(missing) == 0 {
		return nil
	}
	sort.Strings(missing)
	return fmt.Errorf("centymo %s: missing use cases: %s", unit, strings.Join(missing, ", "))
}
