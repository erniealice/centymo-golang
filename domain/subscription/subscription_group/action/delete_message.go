package action

// delete_message.go — view-layer mapping of the espyna DeleteSubscriptionGroup
// referential-guard error onto a tier-label-backed message.
//
// The espyna use case fails closed while a subscription_group still has ACTIVE
// dependents and returns a VERTICAL-NEUTRAL enumerating message. Its use-case
// Translator is the port NoOp, so at runtime that message carries the Go-default
// (general-tier) nouns: "group" / "members" / "staff assignments" /
// "access grants". The per-vertical nouns already exist as lyngua keys
// (subscription_group.validation.*, incl. the education overrides
// "enrolled students" / "teachers"). Per the owner decision the relabel happens
// HERE — the view layer, where the business-type-resolved flat messages
// (viewCtx.T) live — NOT by rewiring the espyna use case or changing its
// message contract (which stays the vertical-neutral fallback).
//
// The code stays generic: the vocabulary is sourced ONLY from lyngua keys
// (resolved through the translate func passed in), never hard-coded vertical
// nouns. If the error is not the referential block, or the flat messages are
// unavailable, the original espyna text passes through unchanged.

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// lyngua flat keys (business-type resolved via viewCtx.T). The dependent-noun
// keys take a "_one"/"_many" suffix chosen by count.
const (
	blockedTemplateKey = "subscription_group.validation.delete_blocked_references"
	memberNounKey      = "subscription_group.validation.dependent_member"
	teachingNounKey    = "subscription_group.validation.dependent_teaching_staff"
	accessGrantNounKey = "subscription_group.validation.dependent_access_grant"
)

// espynaBlockedMarker is a stable fragment of the espyna Go-DEFAULT template
// ("...while it still has active dependents..."). Its presence identifies the
// referential-block error. If the espyna use-case Translator is ever re-wired
// (so the message arrives already tier-flavored, e.g. "...active references..."),
// this marker no longer matches and the message passes through untouched — the
// re-wired use case would already carry the correct tier nouns.
const espynaBlockedMarker = "active dependents"

// Generic-noun count extractors. These mirror the espyna guard's Go-DEFAULT
// nouns (the vertical-neutral message contract) — the ONLY place the generic
// English vocabulary is referenced, and only to READ the counts back out; the
// rendered nouns come from lyngua.
var (
	reMembers      = regexp.MustCompile(`(\d+)\s+members?\b`)
	reTeaching     = regexp.MustCompile(`(\d+)\s+staff assignments?\b`)
	reAccessGrants = regexp.MustCompile(`(\d+)\s+access grants?\b`)
)

// dependentCounts is the per-dependent-type tally parsed from the guard message.
// Vertical-neutral field names; the tier nouns come from lyngua.
type dependentCounts struct {
	members       int
	teachingStaff int
	accessGrants  int
}

func (c dependentCounts) total() int { return c.members + c.teachingStaff + c.accessGrants }

// parseBlockedDependents recognises the espyna referential-guard error and
// extracts the per-type counts. ok=false when the error is not that guard (or
// carries no counts), signalling the caller to keep the original message.
func parseBlockedDependents(msg string) (dependentCounts, bool) {
	if !strings.Contains(msg, espynaBlockedMarker) {
		return dependentCounts{}, false
	}
	var c dependentCounts
	if m := reMembers.FindStringSubmatch(msg); m != nil {
		c.members, _ = strconv.Atoi(m[1])
	}
	if m := reTeaching.FindStringSubmatch(msg); m != nil {
		c.teachingStaff, _ = strconv.Atoi(m[1])
	}
	if m := reAccessGrants.FindStringSubmatch(msg); m != nil {
		c.accessGrants, _ = strconv.Atoi(m[1])
	}
	if c.total() == 0 {
		return dependentCounts{}, false
	}
	return c, true
}

// formatBlockedDependents renders the enumerating message from the tier flat
// messages: each non-zero dimension becomes a "<n> <noun>" clause (singular/
// plural per count), joined and substituted into the delete-blocked template.
// ok=false when the template key does not resolve (flat messages not wired), so
// the caller falls back to the espyna text rather than surfacing a raw key.
func formatBlockedDependents(c dependentCounts, t func(string) string) (string, bool) {
	tmpl := t(blockedTemplateKey)
	if tmpl == blockedTemplateKey || !strings.Contains(tmpl, "{dependents}") {
		return "", false
	}
	var clauses []string
	if c.members > 0 {
		clauses = append(clauses, fmt.Sprintf("%d %s", c.members, t(nounKey(memberNounKey, c.members))))
	}
	if c.teachingStaff > 0 {
		clauses = append(clauses, fmt.Sprintf("%d %s", c.teachingStaff, t(nounKey(teachingNounKey, c.teachingStaff))))
	}
	if c.accessGrants > 0 {
		clauses = append(clauses, fmt.Sprintf("%d %s", c.accessGrants, t(nounKey(accessGrantNounKey, c.accessGrants))))
	}
	return strings.ReplaceAll(tmpl, "{dependents}", strings.Join(clauses, ", ")), true
}

// nounKey appends the singular ("_one") or plural ("_many") suffix by count.
func nounKey(base string, n int) string {
	if n == 1 {
		return base + "_one"
	}
	return base + "_many"
}

// mapDeleteError converts a DeleteSubscriptionGroup error into a tier-label-
// backed message when it is the referential block; otherwise it returns the
// error text verbatim. t is the business-type-resolved flat translator
// (viewCtx.T) — education renders "enrolled students" / "teachers".
func mapDeleteError(err error, t func(string) string) string {
	msg := err.Error()
	c, ok := parseBlockedDependents(msg)
	if !ok {
		return msg
	}
	formatted, ok := formatBlockedDependents(c, t)
	if !ok {
		return msg
	}
	return formatted
}
