package detail

import (
	"context"
	"sort"
	"strings"

	sgpp "github.com/erniealice/centymo-golang/domain/subscription/subscription_group_product_plan"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	jobtemplatephasepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/operation/job_template_phase"
	productplanstaffpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_plan_staff"
	sgpppb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription_group_product_plan"
	sgppspb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription_group_product_plan_staff"
)

// assignmentsTableID is the S5 Teachers-tab table-card id.
const assignmentsTableID = "sgpp-assignments-table"

// assignmentEntry is one resolved assignment row, kept internal for
// sorting (by phase order then staff name — centymo.md §7) before mapping to
// table rows.
type assignmentEntry struct {
	id, staffName, roleLabel, phaseLabel string
	order                                int32
}

// roleDisplayLabel maps a role token to its lyngua label (mirrors
// action.roleLabel — a separate copy since detail/ and action/ are distinct
// packages).
func roleDisplayLabel(l sgpp.Labels, role string) string {
	if role == "secondary" {
		return l.Assign.RoleSecondary
	}
	return l.Assign.RolePrimary
}

// buildAssignmentsTable renders the S5 Teachers tab: one row per assignment
// (staff · role · phase chip), Edit/Remove row actions. Staff names resolve
// via f13 (product_plan_staff_id) -> its staff_id -> name, NEVER legacy f10
// (espyna.md §1b) — the view must survive the M7 legacy retirement.
func buildAssignmentsTable(ctx context.Context, deps *DetailViewDeps, sgppID string, record *sgpppb.SubscriptionGroupProductPlan, pageData *PageData) *types.TableConfig {
	l := deps.Labels
	columns := []types.TableColumn{
		{Key: "staff", Label: l.Roster.ColumnStaff, NoSort: true, NoFilter: true},
		{Key: "role", Label: l.Roster.ColumnRole, NoSort: true, NoFilter: true},
		{Key: "phase", Label: l.Roster.ColumnPhase, NoSort: true, NoFilter: true},
	}

	perms := view.GetUserPermissions(ctx)
	if perms == nil || !perms.Can(sgppsEntity, "list") {
		return &types.TableConfig{
			ID:         assignmentsTableID,
			Columns:    columns,
			Rows:       []types.TableRow{},
			Labels:     deps.TableLabels,
			EmptyState: types.TableEmptyState{Title: l.Roster.EmptyTitle, Message: l.Roster.Unauthorized},
		}
	}

	var names map[string]string
	if deps.ListStaffNames != nil {
		names = deps.ListStaffNames(ctx)
	}

	// ppsStaffByID resolves a product_plan_staff row id -> its staff id — the
	// only sanctioned path to a display name (never a.GetStaffId(), the
	// dual-write legacy mirror).
	ppsStaffByID := map[string]string{}
	if deps.ListProductPlanStaffs != nil && record.GetProductPlanId() != "" {
		if resp, err := deps.ListProductPlanStaffs(ctx, &productplanstaffpb.ListProductPlanStaffsRequest{
			Filters: sgpp.StringEqFilter("product_plan_id", record.GetProductPlanId()),
		}); err == nil && resp != nil {
			for _, pps := range resp.GetData() {
				if pps != nil {
					ppsStaffByID[pps.GetId()] = pps.GetStaffId()
				}
			}
		}
	}

	var phases []sgpp.CoveredPhase
	if deps.ListJobTemplatePhases != nil && record.GetJobTemplateId() != "" {
		if resp, err := deps.ListJobTemplatePhases(ctx, &jobtemplatephasepb.ListJobTemplatePhasesRequest{
			Filters: sgpp.AndFilters(sgpp.StringEqFilter("job_template_id", record.GetJobTemplateId()), sgpp.BoolEqFilter("active", true)),
		}); err == nil && resp != nil {
			phases = sgpp.ResolveCoveredPhases(resp.GetData(), pageData.OfferingVariantID)
		}
	}

	var entries []assignmentEntry
	if deps.ListSubscriptionGroupProductPlanStaffs != nil {
		if resp, err := deps.ListSubscriptionGroupProductPlanStaffs(ctx, &sgppspb.ListSubscriptionGroupProductPlanStaffsRequest{
			Filters: sgpp.AndFilters(sgpp.StringEqFilter("subscription_group_product_plan_id", sgppID), sgpp.BoolEqFilter("active", true)),
		}); err == nil && resp != nil {
			for _, a := range resp.GetData() {
				if a == nil || a.GetSubscriptionGroupProductPlanId() != sgppID || !a.GetActive() {
					continue
				}
				staffID := ppsStaffByID[a.GetProductPlanStaffId()]
				staffName := names[staffID]
				if staffName == "" {
					staffName = staffID
				}
				if staffName == "" {
					staffName = a.GetProductPlanStaffId()
				}
				order := int32(-1)
				phLabel := l.Assign.AllPhasesChip
				if pid := a.GetJobTemplatePhaseId(); pid != "" {
					phLabel = pid
					for _, p := range phases {
						if p.ID == pid {
							phLabel = p.Name
							order = p.Order
							break
						}
					}
				}
				entries = append(entries, assignmentEntry{
					id:         a.GetId(),
					staffName:  staffName,
					roleLabel:  roleDisplayLabel(l, a.GetRole()),
					phaseLabel: phLabel,
					order:      order,
				})
			}
		}
	}
	sort.SliceStable(entries, func(i, j int) bool {
		if entries[i].order != entries[j].order {
			return entries[i].order < entries[j].order
		}
		return strings.ToLower(entries[i].staffName) < strings.ToLower(entries[j].staffName)
	})

	canEdit := perms.Can(sgppsEntity, "update")
	canDelete := perms.Can(sgppsEntity, "delete")
	assignBase := route.ResolveURL(deps.Routes.AssignURL, "sgppId", sgppID)

	rows := make([]types.TableRow, 0, len(entries))
	for _, e := range entries {
		editAction := types.TableAction{
			Type: "edit", Action: "edit", Label: l.Roster.EditAction,
			URL:             assignBase + "?edit=" + e.id,
			DrawerTitle:     l.Roster.EditAction,
			TestID:          "sgpp-assignment-edit-" + e.id,
			Disabled:        !canEdit,
			DisabledTooltip: l.Roster.Unauthorized,
		}
		removeAction := types.TableAction{
			Type: "delete", Action: "delete", Label: l.Roster.RemoveAction,
			URL:             assignBase + "?clear=1",
			ItemName:        e.staffName,
			ConfirmTitle:    l.Confirm.RemoveAssignTitle,
			ConfirmMessage:  strings.ReplaceAll(l.Confirm.RemoveAssignMessage, "{{offering}}", pageData.OfferingName),
			TestID:          "sgpp-assignment-remove-" + e.id,
			Disabled:        !canDelete,
			DisabledTooltip: l.Roster.Unauthorized,
		}
		rows = append(rows, types.TableRow{
			ID: e.id,
			Cells: []types.TableCell{
				{Type: "text", Value: e.staffName},
				{Type: "text", Value: e.roleLabel},
				{Type: "text", Value: e.phaseLabel},
			},
			Actions: []types.TableAction{editAction, removeAction},
		})
	}

	cfg := &types.TableConfig{
		ID:          assignmentsTableID,
		Columns:     columns,
		Rows:        rows,
		Labels:      deps.TableLabels,
		EmptyState:  types.TableEmptyState{Title: l.Roster.EmptyTitle, Message: l.Roster.EmptyMessage},
		ShowActions: true,
	}
	types.ApplyColumnStyles(columns, rows)
	types.ApplyTableSettings(cfg)
	return cfg
}
