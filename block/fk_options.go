package block

// fk_options.go — the FK-picker option producers for the assignment/servicing
// drawers (subscription_group_product_plan_staff, product_plan_staff,
// subscription_group_workspace_user).
//
// Each drawer renders `auto-complete` pickers instead of raw UUID inputs. The
// picker options are `[]form.Pair{ID, Label}` supplied by a lightweight
// `func(ctx) []form.Pair` closure the catalog.go wiring binds onto the module
// deps. This file holds the shared machinery those closures reuse:
//
//   - one generic pair-builder (buildOptionPairs) — a single map loop reused by
//     every producer, no per-entity duplication;
//   - two generic adapters (pairsInto, pairsIntoDisabled) — turn the
//     block-neutral pairs into each module's own form.Pair type at the wiring
//     site, so the block layer never imports a drawer module's form package;
//     the second variant also carries the not-selectable marker;
//   - the four per-entity producers (subscription_group, product_plan, staff,
//     workspace_user), each wrapping a List use case reachable on the centymo
//     UseCases and mapping every row to an id + human label.
//
// Kept out of engineblock.go / usecases.go for the god-file budget, mirroring
// subscription_group_options.go.

import (
	"context"
	"strings"

	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	staffpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/staff"
	workspaceuserpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/workspace_user"
	jobtemplatepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/operation/job_template"
	productplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_plan"
	subscriptiongrouppb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription_group"
)

// staffOptionPickerLimit lifts GetStaffListPageData off its default 50-row page
// so the FK picker shows the whole staff roster, not a page. The page-data list
// variant is used only because its CTE hydrates the nested user for a
// human-readable label. 100 is the hard ceiling the use case enforces
// (GetStaffListPageData rejects Limit > 100); a workspace with more than 100
// staff would need the picker to accumulate pages (future enhancement).
const staffOptionPickerLimit = 100

// optionPair is the block-neutral id/label shape the FK-picker producers emit.
// Each drawer module maps it to that module's own form.Pair via pairsInto at the
// catalog.go wiring site, so the block layer stays free of any drawer module's
// form package.
//
// Disabled marks an option that must render but must not be selectable (the
// auto-complete component greys it out and skips it in keyboard nav). Producers
// that have no lifecycle axis leave it false, so every existing picker keeps all
// of its options selectable.
type optionPair struct {
	ID       string
	Label    string
	Disabled bool
}

// buildOptionPairs is the single reusable pair-builder: one map loop over a list
// use case's rows, extracting an id + human label from each. Rows with an empty
// id are dropped; an empty label falls back to the id so a picker option is
// never blank. Generic over the row type — every producer reuses this loop
// instead of hand-rolling its own.
func buildOptionPairs[T any](rows []T, id func(T) string, label func(T) string) []optionPair {
	pairs := make([]optionPair, 0, len(rows))
	for _, row := range rows {
		rid := id(row)
		if rid == "" {
			continue
		}
		rlabel := label(row)
		if rlabel == "" {
			rlabel = rid
		}
		pairs = append(pairs, optionPair{ID: rid, Label: rlabel})
	}
	return pairs
}

// pairsInto adapts block-neutral option pairs into a module-specific pair slice
// via the module's own constructor. Keeps each producer wiring in catalog.go to
// a single expression with no hand-rolled conversion loop. Drops the Disabled
// flag — for pickers whose options are all selectable.
func pairsInto[P any](src []optionPair, mk func(id, label string) P) []P {
	out := make([]P, 0, len(src))
	for _, p := range src {
		out = append(out, mk(p.ID, p.Label))
	}
	return out
}

// pairsIntoDisabled is the pairsInto sibling that carries the Disabled flag
// through to the module pair. Used only by the pickers whose producer marks
// non-selectable rows (subscription group).
func pairsIntoDisabled[P any](src []optionPair, mk func(id, label string, disabled bool) P) []P {
	out := make([]P, 0, len(src))
	for _, p := range src {
		out = append(out, mk(p.ID, p.Label, p.Disabled))
	}
	return out
}

// subscriptionGroupStatusCurrent is the one subscription_group lifecycle status
// whose rows are selectable in a picker. "completed" / "draft" rows still
// render, greyed out.
const subscriptionGroupStatusCurrent = "current"

// subscriptionGroupOptionPairs lists the workspace's subscription groups as
// id/name pairs. Label = group name. Nil-safe.
//
// The lifecycle status is MARKED, never filtered: a group that is not "current"
// comes back as a Disabled pair so the picker still shows it (an existing row
// pointing at a completed cohort keeps a readable label) while refusing new
// selections. Consumers that only resolve names (pairsToNames) ignore the flag
// and stay status-agnostic.
func subscriptionGroupOptionPairs(
	ctx context.Context,
	list func(context.Context, *subscriptiongrouppb.ListSubscriptionGroupsRequest) (*subscriptiongrouppb.ListSubscriptionGroupsResponse, error),
) []optionPair {
	if list == nil {
		return nil
	}
	resp, err := list(ctx, &subscriptiongrouppb.ListSubscriptionGroupsRequest{})
	if err != nil || resp == nil {
		return nil
	}
	groups := resp.GetData()
	pairs := make([]optionPair, 0, len(groups))
	for _, g := range groups {
		id := g.GetId()
		if id == "" {
			continue
		}
		label := g.GetName()
		if label == "" {
			label = id
		}
		pairs = append(pairs, optionPair{
			ID:       id,
			Label:    label,
			Disabled: g.GetStatus() != subscriptionGroupStatusCurrent,
		})
	}
	return pairs
}

// productPlanOptionPairs lists the workspace's product plans as id/name pairs.
// Label = plan name. Nil-safe.
func productPlanOptionPairs(
	ctx context.Context,
	list func(context.Context, *productplanpb.ListProductPlansRequest) (*productplanpb.ListProductPlansResponse, error),
) []optionPair {
	if list == nil {
		return nil
	}
	resp, err := list(ctx, &productplanpb.ListProductPlansRequest{})
	if err != nil || resp == nil {
		return nil
	}
	return buildOptionPairs(resp.GetData(),
		(*productplanpb.ProductPlan).GetId,
		(*productplanpb.ProductPlan).GetName)
}

// staffOptionPairs lists the workspace's staff as id/person-name pairs. The
// page-data list variant is used because its CTE hydrates the nested user;
// label = "First Last" (falling back to the staff id when the user is not
// hydrated, via buildOptionPairs). Nil-safe.
func staffOptionPairs(
	ctx context.Context,
	list func(context.Context, *staffpb.GetStaffListPageDataRequest) (*staffpb.GetStaffListPageDataResponse, error),
) []optionPair {
	if list == nil {
		return nil
	}
	resp, err := list(ctx, &staffpb.GetStaffListPageDataRequest{
		Pagination: &commonpb.PaginationRequest{Limit: staffOptionPickerLimit},
	})
	if err != nil || resp == nil {
		return nil
	}
	return buildOptionPairs(resp.GetStaffList(),
		(*staffpb.Staff).GetId,
		staffPersonName)
}

// staffPersonName renders "First Last" from the hydrated user, or "" when the
// user is absent (buildOptionPairs then falls back to the staff id).
func staffPersonName(s *staffpb.Staff) string {
	u := s.GetUser()
	return strings.TrimSpace(u.GetFirstName() + " " + u.GetLastName())
}

// workspaceUserOptionPairs lists the workspace's users as id/person-name pairs.
// ListWorkspaceUsers hydrates the nested user; label = "First Last", falling
// back to the email address, then (via buildOptionPairs) the row id. Nil-safe.
func workspaceUserOptionPairs(
	ctx context.Context,
	list func(context.Context, *workspaceuserpb.ListWorkspaceUsersRequest) (*workspaceuserpb.ListWorkspaceUsersResponse, error),
) []optionPair {
	if list == nil {
		return nil
	}
	resp, err := list(ctx, &workspaceuserpb.ListWorkspaceUsersRequest{})
	if err != nil || resp == nil {
		return nil
	}
	return buildOptionPairs(resp.GetData(),
		(*workspaceuserpb.WorkspaceUser).GetId,
		workspaceUserPersonName)
}

// workspaceUserPersonName renders "First Last" from the hydrated user, or the
// email address when the name is blank.
func workspaceUserPersonName(w *workspaceuserpb.WorkspaceUser) string {
	u := w.GetUser()
	name := strings.TrimSpace(u.GetFirstName() + " " + u.GetLastName())
	if name == "" {
		name = strings.TrimSpace(u.GetEmailAddress())
	}
	return name
}

// jobTemplateOptionPairs lists the workspace's job templates (curricula) as
// id/name pairs — the subscription_group_product_plan S7 admin form's
// curriculum picker. Label = template name. Nil-safe.
func jobTemplateOptionPairs(
	ctx context.Context,
	list func(context.Context, *jobtemplatepb.ListJobTemplatesRequest) (*jobtemplatepb.ListJobTemplatesResponse, error),
) []optionPair {
	if list == nil {
		return nil
	}
	resp, err := list(ctx, &jobtemplatepb.ListJobTemplatesRequest{})
	if err != nil || resp == nil {
		return nil
	}
	return buildOptionPairs(resp.GetData(),
		(*jobtemplatepb.JobTemplate).GetId,
		(*jobtemplatepb.JobTemplate).GetName)
}

// pairsToNames adapts block-neutral option pairs into an id -> label map —
// the shape the subscription_group_product_plan S7 list's name-resolution
// batches (ListSubscriptionGroupNames / ListProductPlanNames /
// ListJobTemplateNames) need.
func pairsToNames(pairs []optionPair) map[string]string {
	out := make(map[string]string, len(pairs))
	for _, p := range pairs {
		out[p.ID] = p.Label
	}
	return out
}
