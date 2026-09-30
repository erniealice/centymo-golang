package block

import "testing"

// R5 M7: the one Mount/RequireFor helper names every missing closure, sorted, and is nil when complete.
func TestRequireUnitSortedAndNilWhenComplete(t *testing.T) {
	if err := requireUnit("x", map[string]bool{"A": true, "B": true}); err != nil {
		t.Fatalf("complete = %v", err)
	}
	for i := 0; i < 20; i++ { // map order is random: the output must not be
		err := requireUnit("cost_source_component", map[string]bool{"Zed": false, "Alpha": false, "Mid": true, "Beta": false})
		want := "centymo cost_source_component: missing use cases: Alpha, Beta, Zed"
		if err == nil || err.Error() != want {
			t.Fatalf("err = %v, want %q", err, want)
		}
	}
}
