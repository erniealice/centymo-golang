package action

import (
	"errors"
	"testing"
)

// tierT returns a translate func backed by a flat message map — mirroring
// viewCtx.T: a hit returns the value, a miss echoes the key back.
func tierT(msgs map[string]string) func(string) string {
	return func(key string) string {
		if v, ok := msgs[key]; ok {
			return v
		}
		return key
	}
}

// generalMessages mirrors packages/lyngua/translations/en/general/
// subscription_group.json (validation block) — the vertical-neutral tier.
func generalMessages() map[string]string {
	return map[string]string{
		"subscription_group.validation.delete_blocked_references":     "Cannot delete this section while it still has active references: {dependents}. Remove them first.",
		"subscription_group.validation.dependent_member_one":          "member",
		"subscription_group.validation.dependent_member_many":         "members",
		"subscription_group.validation.dependent_teaching_staff_one":  "staff assignment",
		"subscription_group.validation.dependent_teaching_staff_many": "staff assignments",
		"subscription_group.validation.dependent_access_grant_one":    "access grant",
		"subscription_group.validation.dependent_access_grant_many":   "access grants",
	}
}

// educationMessages mirrors the general → education cascade: education overrides
// only the member + teaching-staff nouns; the template and access-grant nouns
// are inherited from general (education/subscription_group.json sets neither).
func educationMessages() map[string]string {
	m := generalMessages()
	m["subscription_group.validation.dependent_member_one"] = "enrolled student"
	m["subscription_group.validation.dependent_member_many"] = "enrolled students"
	m["subscription_group.validation.dependent_teaching_staff_one"] = "teacher"
	m["subscription_group.validation.dependent_teaching_staff_many"] = "teachers"
	return m
}

// espynaGeneric reproduces the runtime espyna DeleteSubscriptionGroup message:
// the NoOp Translator returns the Go-default template with {dependents}
// substituted from the Go-default (vertical-neutral) nouns.
func espynaGeneric(dependents string) string {
	return "Cannot delete this group while it still has active dependents: " + dependents + ". Remove them first."
}

func TestParseBlockedDependents(t *testing.T) {
	cases := []struct {
		name   string
		msg    string
		want   dependentCounts
		wantOK bool
	}{
		{"all three", espynaGeneric("3 members, 2 staff assignments, 1 access grant"), dependentCounts{3, 2, 1}, true},
		{"one member", espynaGeneric("1 member"), dependentCounts{1, 0, 0}, true},
		{"one staff assignment", espynaGeneric("1 staff assignment"), dependentCounts{0, 1, 0}, true},
		{"two access grants", espynaGeneric("2 access grants"), dependentCounts{0, 0, 2}, true},
		{"members + grants only", espynaGeneric("5 members, 4 access grants"), dependentCounts{5, 0, 4}, true},
		{"not the guard error", "some other failure", dependentCounts{}, false},
		{"marker but no counts", "blocked by active dependents somewhere", dependentCounts{}, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := parseBlockedDependents(tc.msg)
			if ok != tc.wantOK {
				t.Fatalf("ok = %v, want %v", ok, tc.wantOK)
			}
			if ok && got != tc.want {
				t.Fatalf("counts = %+v, want %+v", got, tc.want)
			}
		})
	}
}

func TestFormatBlockedDependents_GeneralTier(t *testing.T) {
	t2 := tierT(generalMessages())

	got, ok := formatBlockedDependents(dependentCounts{3, 2, 1}, t2)
	if !ok {
		t.Fatal("ok = false, want true")
	}
	want := "Cannot delete this section while it still has active references: 3 members, 2 staff assignments, 1 access grant. Remove them first."
	if got != want {
		t.Fatalf("general tier =\n  %q\nwant\n  %q", got, want)
	}

	// Singular member clause.
	got, _ = formatBlockedDependents(dependentCounts{1, 0, 0}, t2)
	want = "Cannot delete this section while it still has active references: 1 member. Remove them first."
	if got != want {
		t.Fatalf("general singular =\n  %q\nwant\n  %q", got, want)
	}
}

func TestFormatBlockedDependents_EducationTier(t *testing.T) {
	t2 := tierT(educationMessages())

	// Plural, all three dimensions. member/staff use the education overrides;
	// access grant is inherited from general; the template is inherited too.
	got, ok := formatBlockedDependents(dependentCounts{3, 2, 1}, t2)
	if !ok {
		t.Fatal("ok = false, want true")
	}
	want := "Cannot delete this section while it still has active references: 3 enrolled students, 2 teachers, 1 access grant. Remove them first."
	if got != want {
		t.Fatalf("education tier =\n  %q\nwant\n  %q", got, want)
	}

	// Singular of each dimension.
	got, _ = formatBlockedDependents(dependentCounts{1, 1, 1}, t2)
	want = "Cannot delete this section while it still has active references: 1 enrolled student, 1 teacher, 1 access grant. Remove them first."
	if got != want {
		t.Fatalf("education singular =\n  %q\nwant\n  %q", got, want)
	}
}

func TestFormatBlockedDependents_MessagesNotWired(t *testing.T) {
	// nil-backed translator echoes keys — the template guard must reject so the
	// caller can fall back to the espyna text instead of surfacing a raw key.
	if _, ok := formatBlockedDependents(dependentCounts{1, 0, 0}, tierT(nil)); ok {
		t.Fatal("ok = true, want false when the template key does not resolve")
	}
}

func TestMapDeleteError(t *testing.T) {
	eduT := tierT(educationMessages())

	// Referential-block error → re-rendered with education nouns.
	blockErr := errors.New(espynaGeneric("3 members, 2 staff assignments, 1 access grant"))
	got := mapDeleteError(blockErr, eduT)
	want := "Cannot delete this section while it still has active references: 3 enrolled students, 2 teachers, 1 access grant. Remove them first."
	if got != want {
		t.Fatalf("mapped =\n  %q\nwant\n  %q", got, want)
	}

	// Non-guard error passes through verbatim.
	if got := mapDeleteError(errors.New("boom"), eduT); got != "boom" {
		t.Fatalf("passthrough = %q, want %q", got, "boom")
	}

	// Guard error but flat messages unavailable → keep the espyna text.
	if got := mapDeleteError(blockErr, tierT(nil)); got != blockErr.Error() {
		t.Fatalf("fallback = %q, want the espyna text %q", got, blockErr.Error())
	}
}
