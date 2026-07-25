package detail

// staff.go — the section's "Teaching Staff"/"Subjects" tab.
//
// M4 ROW-SOURCE FLIP (plan.md §2, centymo.md §3): the LIVE tab lists active
// subscription_group_product_plan (class) rows — SectionSGPPRow /
// SectionSGPPTabData / buildSGPPTabData / countSectionSGPPs, below the
// "M4 row-source flip" banner. Two columns (offering + Teachers, ALL
// assignments comma-separated with phase labels); row actions View / Assign /
// Exclude / Restore / Remove wired to the subscription_group_product_plan
// module's own routes (cross-domain closures on DetailViewDeps). A section
// with ZERO class rows (pre-M3-backfill) falls back to today's DERIVED
// offering rows with every action disabled + a "not yet set up" hint.
//
// LEGACY (§6.1–§6.3, kept byte-identical, no longer surfaced): the original
// per-offering table that assigns an eligible servicer to every product_plan
// of the section's plan — SectionAssignmentRow / SectionStaffTabData /
// buildStaffTabData / countSectionOfferings / renderAssignDrawer /
// NewAssignAction. NewAssignAction stays REGISTERED (the year-grain upsert
// route still resolves) but no row action opens it anymore; the functions
// below are retained for their existing unit coverage (staff_test.go) and as
// the fallback row source's data provider (listPlanOfferings, offeringLabel).
//
// All vertical vocabulary (subject / teacher / class) lives in lyngua only —
// code identifiers below use "sgpp" (the entity's own short form, matching
// the sibling subscription_group_product_plan package's own convention) for
// consistency across the feature, per that package's own labels.go precedent.

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/erniealice/centymo-golang/domain/subscription/subscription_group"
	"github.com/erniealice/espyna-golang/consumer"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/render"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	jobtemplatephasepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/operation/job_template_phase"
	productplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_plan"
	productplanstaffpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_plan_staff"
	productvariantpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_variant"
	subscriptiongrouppb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription_group"
	sgpppb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription_group_product_plan"
	sgppspb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription_group_product_plan_staff"
)

// sgppEntity is the permission entity for the class row itself
// (subscription_group_product_plan) — distinct from sgppsEntity, which
// governs the assignment edge (plan.md §5 D-5).
const sgppEntity = "subscription_group_product_plan"

// sgppTableID is the M4 table-card id (the refresh-target token).
const sgppTableID = "subscription-group-sgpp-table"

// sgppsEntity is the permission entity for the class-edge (subscription_group_
// product_plan_staff). The table reads with :list and the drawer writes with
// :{create,update,delete}, fail-closed.
const sgppsEntity = "subscription_group_product_plan_staff"

// rolePrimary is the default role a fresh (or auto-defaulted) assignment carries.
// Role tokens are exactly {"primary","access"}: primary = teacher-of-record,
// access = visibility-only. Labels come from lyngua.
const (
	rolePrimary = "primary"
	roleAccess  = "access"
)

// staffAssignTab is the reserved tab token that serves the per-offering assign
// drawer through the ALREADY-REGISTERED tab-action GET route (no new route): the
// row action GETs /action/subscription-group/{id}/tab/staff-assign?pp={plan}
// into #sheetContent. NewTabAction dispatches this token to renderAssignDrawer.
const staffAssignTab = "staff-assign"

// staffTableID is the table-card id (and the refresh-target token). refreshTable
// resolves the enclosing tab panel by id (staffTabID); the table-card renders as
// "<id>-card".
const (
	staffTableID = "subscription-group-staff-table"
	staffTabID   = "subscription-group-staff-tab"
)

// emDash is the placeholder rendered in an unassigned row's teacher/role cells.
const emDash = "—"

// SectionAssignmentRow is one offering (product_plan) row: the subject label, the
// current servicer (empty when unassigned) with its resolved name, the eligible
// pool as drawer autocomplete options, and the existing edge id for edit/clear.
type SectionAssignmentRow struct {
	ProductPlanID    string
	SubjectLabel     string // product_plan.name (data, not a lyngua label)
	CurrentStaffID   string // "" when unassigned
	CurrentStaffName string // resolved person name (falls back to id); "" when unassigned
	CurrentRole      string // "" when unassigned
	SgppsID          string // the active class-edge id, for edit/clear
	Assigned         bool
	Eligible         []types.SelectOption // {Value: staff_id, Label: person name}; empty ⇒ gate message
	EligibilityURL   string               // where to set eligibility when the pool is empty
}

// SectionStaffTabData is the whole tab payload (§6.2). The Table is the shared
// table-card config; RefreshURL re-renders this whole tab panel after a drawer
// POST (reused tab-action GET). AssignActionURL is the signed assign POST path
// the drawer form targets.
type SectionStaffTabData struct {
	SubscriptionGroupID string
	AssignActionURL     string // AssignURL resolved with the section id
	RefreshURL          string // tab-action GET for tab=staff, refreshes this panel
	Rows                []SectionAssignmentRow
	Table               *types.TableConfig
	CanEdit             bool
	Authorized          bool // false ⇒ render the read-only "no permission" note
	Labels              subscription_group.Labels
}

// staffAssignDrawer is the per-offering assign drawer payload. WorkspaceID/Nonce
// are injected reflectively by the render pipeline (top-level fields) so the
// form's {{actionForm}} can sign the assign path.
type staffAssignDrawer struct {
	WorkspaceID       string
	Nonce             string
	AssignActionURL   string
	ProductPlanID     string
	OfferingLabel     string
	CurrentStaffID    string
	CurrentStaffLabel string
	Assigned          bool
	TeacherOptions    []types.SelectOption
	RoleOptions       []types.SelectOption
	Labels            subscription_group.Labels
	CommonLabels      pyeza.CommonLabels
}

// eligibleEntry is an internal eligible-pool row: the staff id + role (role
// drives the sole-primary auto-default) with the resolved person name.
type eligibleEntry struct {
	staffID string
	role    string
	name    string
}

// buildStaffTabData assembles the tab for one section by composing existing List
// use cases (§6.2): ListProductPlans(plan_id) = the offering rows,
// ListSubscriptionGroupProductPlanStaffs(section, active) = current assignments,
// ListProductPlanStaffs(product_plan_id, active) = the eligible pool per offering.
// Fail-closed on sgppsEntity:list INSIDE the tab body (a note, never a full-page
// Forbidden — which would be wrong for an HTMX tab swap).
func buildStaffTabData(ctx context.Context, deps *DetailViewDeps, sg *subscriptiongrouppb.SubscriptionGroup, l subscription_group.Labels) *SectionStaffTabData {
	groupID := sg.GetId()
	data := &SectionStaffTabData{
		SubscriptionGroupID: groupID,
		AssignActionURL:     route.ResolveURL(deps.Routes.AssignURL, "id", groupID),
		RefreshURL:          route.ResolveURL(deps.Routes.TabActionURL, "id", groupID, "tab", "staff"),
		Labels:              l,
	}

	perms := view.GetUserPermissions(ctx)
	if perms == nil || !perms.Can(sgppsEntity, "list") {
		// Fail-closed read gate: no table, no drawers, a read-only note.
		return data
	}
	data.Authorized = true
	data.CanEdit = perms.CanAny(sgppsEntity+":create", sgppsEntity+":update")

	planID := sg.GetPlanId()
	if planID == "" || deps.ListProductPlans == nil {
		return data
	}

	plans := listPlanOfferings(ctx, deps, planID)
	current := listActiveAssignments(ctx, deps, groupID)                // product_plan_id → active edge
	names := staffNames(ctx, deps)                                      // staff_id → person name
	eligible := listEligiblePools(ctx, deps, offeringIDs(plans), names) // product_plan_id → []eligibleEntry

	rows := make([]SectionAssignmentRow, 0, len(plans))
	for _, pp := range plans {
		ppID := pp.GetId()
		row := SectionAssignmentRow{
			ProductPlanID:  ppID,
			SubjectLabel:   offeringLabel(pp),
			EligibilityURL: eligibilityURL(deps.ProductPlanStaffListURL, ppID),
		}
		if edge := current[ppID]; edge != nil {
			row.CurrentStaffID = edge.GetStaffId()
			row.CurrentRole = edge.GetRole()
			row.SgppsID = edge.GetId()
			row.Assigned = row.CurrentStaffID != ""
			if row.Assigned {
				n := names[row.CurrentStaffID]
				if n == "" {
					n = row.CurrentStaffID
				}
				row.CurrentStaffName = n
			}
		}
		row.Eligible = eligibleOptions(eligible[ppID], row.CurrentStaffID, row.Assigned)
		rows = append(rows, row)
	}

	data.Rows = rows
	data.Table = buildStaffTable(deps, data, l)
	return data
}

// countSectionOfferings returns the number of offering rows (active product_plan
// rows of the section's program) for the Teaching-Staff tab count badge. It
// mirrors the tab's row source (listPlanOfferings) and read gate (sgppsEntity:list)
// exactly, so the badge count matches the rendered row count and cannot drift.
// CHEAP by design — a single ListProductPlans call (the enrollments-count
// precedent), never a full buildStaffTabData run (which also fans out to
// assignments + names + eligible pools). Returns 0 (no badge — the tabs component
// renders nothing for a zero Count) when the table is not permitted, the program
// is unset, or the dep is unwired.
func countSectionOfferings(ctx context.Context, deps *DetailViewDeps, sg *subscriptiongrouppb.SubscriptionGroup) int {
	perms := view.GetUserPermissions(ctx)
	if perms == nil || !perms.Can(sgppsEntity, "list") {
		return 0
	}
	planID := sg.GetPlanId()
	if planID == "" || deps.ListProductPlans == nil {
		return 0
	}
	return len(listPlanOfferings(ctx, deps, planID))
}

// buildStaffTable maps the assignment rows into the shared table-card config:
// one row per offering (subject / teacher / role / state) + a per-row assign
// action that opens the drawer. It mirrors the STANDARD detail-tab table exactly
// — the Enrollments/Subscriptions tab's config (buildSubscriptionsTable): the same
// client-side toolbar/footer chrome (ShowSearch/ShowColumns/ShowDensity/ShowEntries)
// so the staff table is visually and behaviorally identical to every other standard
// table in the app. ShowActions is added on top (the assign row action; the
// Subscriptions roster has none). The columns keep NoSort/NoFilter (lookup-style,
// not indexable list fields — same rationale as the roster). The tab-scoped empty
// copy (l.Staff.Empty*) is retained rather than the generic l.Empty.*, and the
// Status column's Saved/Unassigned chips carry the "gaps at a glance" summary
// (W2g removed the redundant in-panel coverage badge).
func buildStaffTable(deps *DetailViewDeps, data *SectionStaffTabData, l subscription_group.Labels) *types.TableConfig {
	columns := []types.TableColumn{
		{Key: "subject", Label: l.Staff.ColumnSubject, NoSort: true, NoFilter: true, Width: "36%"},
		{Key: "teacher", Label: l.Staff.ColumnServicer, NoSort: true, NoFilter: true},
		{Key: "role", Label: l.Staff.ColumnRole, NoSort: true, NoFilter: true, WidthClass: "col-md"},
		{Key: "state", Label: l.Staff.ColumnState, NoSort: true, NoFilter: true, WidthClass: "col-md"},
	}
	rows := make([]types.TableRow, 0, len(data.Rows))
	for _, r := range data.Rows {
		rows = append(rows, staffTableRow(deps, data, r, l))
	}
	cfg := &types.TableConfig{
		ID:          staffTableID,
		Columns:     columns,
		Rows:        rows,
		Labels:      deps.TableLabels,
		EmptyState:  types.TableEmptyState{Title: l.Staff.EmptyTitle, Message: l.Staff.EmptyMessage},
		ShowSearch:  true,
		ShowColumns: true,
		ShowDensity: true,
		ShowEntries: true,
		ShowActions: true,
	}
	types.ApplyColumnStyles(columns, rows)
	types.ApplyTableSettings(cfg)
	return cfg
}

// staffTableRow builds one table row for an offering. The teacher cell renders the
// current servicer name, a muted em-dash when unassigned, or the empty-pool gate
// link when there is no eligible staff. The assign action is fail-closed: disabled
// when the operator lacks create/update, or when the pool is empty (nothing to
// assign — the gate link is the next step instead).
func staffTableRow(deps *DetailViewDeps, data *SectionStaffTabData, r SectionAssignmentRow, l subscription_group.Labels) types.TableRow {
	emptyPool := len(r.Eligible) == 0

	subjectCell := types.TableCell{Type: "text", Value: r.SubjectLabel}

	var teacherCell types.TableCell
	switch {
	case emptyPool:
		teacherCell = gateCell(r, l)
	case r.Assigned:
		teacherCell = types.TableCell{Type: "text", Value: r.CurrentStaffName}
	default:
		teacherCell = types.TableCell{Type: "text", Value: emDash}
	}

	roleValue := emDash
	if r.Assigned {
		roleValue = roleLabel(l, r.CurrentRole)
	}
	roleCell := types.TableCell{Type: "text", Value: roleValue}

	stateValue, stateVariant := l.Staff.Unassigned, "warning"
	if r.Assigned {
		stateValue, stateVariant = l.Staff.Saved, "success"
	}
	stateCell := types.TableCell{Type: "badge", Value: stateValue, Variant: stateVariant}

	// Standard row-action convention (WCAG 2.1.1): a native <button
	// data-action="edit" data-edit-url=…> — the shared table component's default
	// action branch (table.html) — which table-actions.js opens into #sheetContent
	// via htmx.ajax + lf.Sheet.open(DrawerTitle). A native button is in the tab
	// order and activates on Enter/Space, unlike the HxGet branch's <a> without
	// href (keyboard-unreachable, role="generic"). assignDrawerURL already carries
	// ?pp=<offering>; the handler's appended &id=<row> is inert (renderAssignDrawer
	// reads pp).
	action := types.TableAction{
		Type:        "edit",
		Action:      "edit",
		Label:       l.Staff.AssignAction,
		TestID:      "sg-staff-assign-" + r.ProductPlanID,
		URL:         assignDrawerURL(deps, data.SubscriptionGroupID, r.ProductPlanID),
		DrawerTitle: l.Staff.AssignAction,
	}
	switch {
	case !data.CanEdit:
		action.Disabled = true
		action.DisabledTooltip = l.Staff.Unauthorized
	case emptyPool:
		action.Disabled = true
		action.DisabledTooltip = l.Staff.EmptyPool
	}

	return types.TableRow{
		ID: r.ProductPlanID,
		// "testid" renders as data-testid on the <tr> (the table component maps
		// every DataAttr key to data-<key>), giving the spec a stable per-row
		// selector; "assigned" drives the gaps-at-a-glance state.
		DataAttrs: map[string]string{
			"assigned": strconv.FormatBool(r.Assigned),
			"testid":   "sg-staff-row-" + r.ProductPlanID,
		},
		Cells:   []types.TableCell{subjectCell, teacherCell, roleCell, stateCell},
		Actions: []types.TableAction{action},
	}
}

// gateCell renders the empty-pool gate as an html cell so the stable
// sg-staff-gate-{plan} test id survives the move into the table (the typed cell
// types carry no per-cell test id). Every interpolated value is trusted (a
// composed /action URL, a UUID, and a lyngua label) and HTML-escaped; there is no
// user input, hence no XSS surface. Value carries the label for a sane CSV export.
func gateCell(r SectionAssignmentRow, l subscription_group.Labels) types.TableCell {
	label := l.Staff.EmptyPool
	testID := "sg-staff-gate-" + r.ProductPlanID
	var b strings.Builder
	if r.EligibilityURL != "" {
		b.WriteString(`<a class="centymo-sgstaff-gate" href="`)
		b.WriteString(template.HTMLEscapeString(r.EligibilityURL))
		b.WriteString(`" data-testid="`)
		b.WriteString(template.HTMLEscapeString(testID))
		b.WriteString(`">`)
		b.WriteString(template.HTMLEscapeString(label))
		b.WriteString(`</a>`)
	} else {
		b.WriteString(`<span class="centymo-sgstaff-gate" data-testid="`)
		b.WriteString(template.HTMLEscapeString(testID))
		b.WriteString(`">`)
		b.WriteString(template.HTMLEscapeString(label))
		b.WriteString(`</span>`)
	}
	return types.TableCell{Type: "html", HTML: template.HTML(b.String()), Value: label}
}

// assignDrawerURL builds the row action's drawer GET URL: the reserved
// staff-assign tab token on the registered tab-action route, plus ?pp=<offering>.
func assignDrawerURL(deps *DetailViewDeps, groupID, productPlanID string) string {
	base := route.ResolveURL(deps.Routes.TabActionURL, "id", groupID, "tab", staffAssignTab)
	return base + "?pp=" + productPlanID
}

// roleLabel maps a role token to its lyngua label; unknown tokens pass through.
func roleLabel(l subscription_group.Labels, role string) string {
	switch role {
	case rolePrimary:
		return l.Staff.RolePrimary
	case roleAccess:
		return l.Staff.RoleAccess
	default:
		return role
	}
}

// buildRoleOptions returns the bounded role <select> options {primary, access},
// defaulting to primary when the current role is empty/unknown.
func buildRoleOptions(l subscription_group.Labels, selected string) []types.SelectOption {
	if selected != rolePrimary && selected != roleAccess {
		selected = rolePrimary
	}
	defs := []struct{ value, label string }{
		{rolePrimary, l.Staff.RolePrimary},
		{roleAccess, l.Staff.RoleAccess},
	}
	opts := make([]types.SelectOption, 0, len(defs))
	for _, d := range defs {
		opts = append(opts, types.SelectOption{Value: d.value, Label: d.label, Selected: d.value == selected})
	}
	return opts
}

// listPlanOfferings returns the section program's offerings (product_plan rows),
// sorted by display label for a stable table order.
func listPlanOfferings(ctx context.Context, deps *DetailViewDeps, planID string) []*productplanpb.ProductPlan {
	if deps.ListProductPlans == nil {
		return nil
	}
	resp, err := deps.ListProductPlans(ctx, &productplanpb.ListProductPlansRequest{
		Filters: stringEqFilter("plan_id", planID),
	})
	if err != nil {
		log.Printf("staff tab: list product plans for plan %s: %v", planID, err)
		return nil
	}
	plans := make([]*productplanpb.ProductPlan, 0, len(resp.GetData()))
	for _, pp := range resp.GetData() {
		if pp == nil || !pp.GetActive() {
			continue
		}
		plans = append(plans, pp)
	}
	sort.SliceStable(plans, func(i, j int) bool {
		a, b := strings.ToLower(offeringLabel(plans[i])), strings.ToLower(offeringLabel(plans[j]))
		if a == b {
			return plans[i].GetId() < plans[j].GetId()
		}
		return a < b
	})
	return plans
}

// listActiveAssignments returns the section's current class edges keyed by
// product_plan_id (active only; the first active edge per offering wins).
func listActiveAssignments(ctx context.Context, deps *DetailViewDeps, groupID string) map[string]*sgppspb.SubscriptionGroupProductPlanStaff {
	out := map[string]*sgppspb.SubscriptionGroupProductPlanStaff{}
	if deps.ListSubscriptionGroupProductPlanStaffs == nil {
		return out
	}
	resp, err := deps.ListSubscriptionGroupProductPlanStaffs(ctx, &sgppspb.ListSubscriptionGroupProductPlanStaffsRequest{
		Filters: stringEqFilter("subscription_group_id", groupID),
	})
	if err != nil {
		log.Printf("staff tab: list class edges for section %s: %v", groupID, err)
		return out
	}
	for _, e := range resp.GetData() {
		if e == nil || !e.GetActive() || e.GetSubscriptionGroupId() != groupID {
			continue
		}
		if _, seen := out[e.GetProductPlanId()]; !seen {
			out[e.GetProductPlanId()] = e
		}
	}
	return out
}

// listEligiblePools returns the eligible servicer pool per offering (active
// product_plan_staff rows), keyed by product_plan_id, resolving person names from
// the caller-supplied staff-name batch. One LIST_IN call over all offering ids.
func listEligiblePools(ctx context.Context, deps *DetailViewDeps, offeringIDs []string, names map[string]string) map[string][]eligibleEntry {
	out := map[string][]eligibleEntry{}
	if len(offeringIDs) == 0 || deps.ListProductPlanStaffs == nil {
		return out
	}
	resp, err := deps.ListProductPlanStaffs(ctx, &productplanstaffpb.ListProductPlanStaffsRequest{
		Filters: stringInFilter("product_plan_id", offeringIDs),
	})
	if err != nil {
		log.Printf("staff tab: list eligibility pools: %v", err)
		return out
	}
	for _, pps := range resp.GetData() {
		if pps == nil || !pps.GetActive() {
			continue
		}
		ppID := pps.GetProductPlanId()
		staffID := pps.GetStaffId()
		if ppID == "" || staffID == "" {
			continue
		}
		name := names[staffID]
		if name == "" {
			name = staffID
		}
		out[ppID] = append(out[ppID], eligibleEntry{staffID: staffID, role: pps.GetRole(), name: name})
	}
	// Stable option order per offering: by person name, then id.
	for ppID := range out {
		entries := out[ppID]
		sort.SliceStable(entries, func(i, j int) bool {
			a, b := strings.ToLower(entries[i].name), strings.ToLower(entries[j].name)
			if a == b {
				return entries[i].staffID < entries[j].staffID
			}
			return a < b
		})
		out[ppID] = entries
	}
	return out
}

// eligibleOptions maps the eligible pool to drawer autocomplete options, marking
// the current servicer selected. Auto-default (§6.1): when the offering is
// unassigned and exactly one eligible role='primary' exists, that sole primary is
// pre-selected in the drawer (the operator still saves to persist it — no write
// happens on render).
func eligibleOptions(pool []eligibleEntry, currentStaffID string, assigned bool) []types.SelectOption {
	if len(pool) == 0 {
		return nil
	}
	autoDefault := ""
	if !assigned {
		primaries := 0
		soleID := ""
		for _, e := range pool {
			if e.role == rolePrimary {
				primaries++
				soleID = e.staffID
			}
		}
		if primaries == 1 {
			autoDefault = soleID
		}
	}
	opts := make([]types.SelectOption, 0, len(pool))
	for _, e := range pool {
		selected := (assigned && e.staffID == currentStaffID) || (!assigned && e.staffID == autoDefault)
		opts = append(opts, types.SelectOption{Value: e.staffID, Label: e.name, Selected: selected})
	}
	return opts
}

// staffNames returns the workspace staff id → person-name map (nil-safe).
func staffNames(ctx context.Context, deps *DetailViewDeps) map[string]string {
	if deps.ListStaffNames == nil {
		return nil
	}
	return deps.ListStaffNames(ctx)
}

// offeringIDs collects the product_plan ids for the LIST_IN eligibility query.
func offeringIDs(plans []*productplanpb.ProductPlan) []string {
	ids := make([]string, 0, len(plans))
	for _, pp := range plans {
		if id := pp.GetId(); id != "" {
			ids = append(ids, id)
		}
	}
	return ids
}

// offeringLabel is the offering's display name (product_plan.name), falling back
// to its id so a row is never blank.
func offeringLabel(pp *productplanpb.ProductPlan) string {
	if n := strings.TrimSpace(pp.GetName()); n != "" {
		return n
	}
	return pp.GetId()
}

// eligibilityURL builds the empty-pool gate link to the product_plan_staff
// management surface for one offering. Returns "" when no base is wired (the
// template then renders the gate text without a link).
func eligibilityURL(base, productPlanID string) string {
	if base == "" {
		return ""
	}
	resolved := route.ResolveURL(base, "status", "active")
	if productPlanID != "" {
		resolved += "?product_plan_id=" + productPlanID
	}
	return resolved
}

// stringEqFilter builds a single field == value FilterRequest.
func stringEqFilter(field, value string) *commonpb.FilterRequest {
	return &commonpb.FilterRequest{Filters: []*commonpb.TypedFilter{{
		Field: field,
		FilterType: &commonpb.TypedFilter_StringFilter{
			StringFilter: &commonpb.StringFilter{Value: value, Operator: commonpb.StringOperator_STRING_EQUALS},
		},
	}}}
}

// stringInFilter builds a single field IN (values) FilterRequest.
func stringInFilter(field string, values []string) *commonpb.FilterRequest {
	return &commonpb.FilterRequest{Filters: []*commonpb.TypedFilter{{
		Field: field,
		FilterType: &commonpb.TypedFilter_ListFilter{
			ListFilter: &commonpb.ListFilter{Values: values, Operator: commonpb.ListOperator_LIST_IN},
		},
	}}}
}

// renderAssignDrawer serves the per-offering assign drawer (GET, via the reserved
// staff-assign tab token — see staffAssignTab). It re-derives the section tab data
// (so the eligible pool + current assignment come straight from the page-data),
// finds the addressed offering, and returns the drawer form. Fail-closed: the
// operator must hold create OR update on the class-edge.
func renderAssignDrawer(ctx context.Context, deps *DetailViewDeps, viewCtx *view.ViewContext, groupID string) view.ViewResult {
	perms := view.GetUserPermissions(ctx)
	if perms == nil || !perms.CanAny(sgppsEntity+":create", sgppsEntity+":update") {
		return view.HTMXError(deps.Labels.Errors.Unauthorized)
	}
	productPlanID := strings.TrimSpace(viewCtx.Request.URL.Query().Get("pp"))
	if productPlanID == "" {
		return view.HTMXError(deps.Labels.Errors.NotFound)
	}
	resp, err := deps.ReadSubscriptionGroup(ctx, &subscriptiongrouppb.ReadSubscriptionGroupRequest{
		Data: &subscriptiongrouppb.SubscriptionGroup{Id: groupID},
	})
	if err != nil || len(resp.GetData()) == 0 {
		return view.HTMXError(deps.Labels.Errors.NotFound)
	}
	tab := buildStaffTabData(ctx, deps, resp.GetData()[0], deps.Labels)
	row, ok := findRow(tab.Rows, productPlanID)
	if !ok {
		return view.HTMXError(deps.Labels.Errors.NotFound)
	}
	// The autocomplete renders its trigger + hidden input from Value/SelectedLabel
	// (NOT the options' Selected flag), so surface the pre-selected teacher the
	// page-data already resolved: the current servicer when assigned, or the
	// sole-primary auto-default when unassigned (eligibleOptions marks it).
	teacherValue, teacherLabel := row.CurrentStaffID, row.CurrentStaffName
	for _, o := range row.Eligible {
		if o.Selected {
			teacherValue, teacherLabel = o.Value, o.Label
			break
		}
	}
	return view.OK("subscription-group-staff-assign-drawer", &staffAssignDrawer{
		AssignActionURL:   tab.AssignActionURL,
		ProductPlanID:     productPlanID,
		OfferingLabel:     row.SubjectLabel,
		CurrentStaffID:    teacherValue,
		CurrentStaffLabel: teacherLabel,
		Assigned:          row.Assigned,
		TeacherOptions:    row.Eligible,
		RoleOptions:       buildRoleOptions(deps.Labels, row.CurrentRole),
		Labels:            deps.Labels,
		CommonLabels:      deps.CommonLabels,
	})
}

// NewAssignAction handles POST /action/subscription-group/assign/{id} — the
// class-edge upsert (§6.3). The section id is authoritative from the signed path;
// product_plan_id/staff_id/role come from the drawer body. Permission-gated
// fail-closed; on success it returns a header-only response that closes the drawer
// and refreshes the whole staff tab panel (refreshTable → staffTabID).
func NewAssignAction(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		groupID := viewCtx.Request.PathValue("id")
		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.HTMXError(deps.Labels.Errors.UpdateFailed)
		}
		productPlanID := strings.TrimSpace(viewCtx.Request.FormValue("product_plan_id"))
		staffID := strings.TrimSpace(viewCtx.Request.FormValue("staff_id"))
		role := strings.TrimSpace(viewCtx.Request.FormValue("role"))
		if role == "" {
			role = rolePrimary
		}
		// The drawer's explicit Clear control (shown only when an edge exists)
		// submits clear=1 — a request to soft-delete the active edge (§6.1
		// "Clear"). It overrides whatever staff_id the hidden autocomplete input
		// still carries, routing to the empty-staff (delete) branch below.
		if viewCtx.Request.FormValue("clear") == "1" {
			staffID = ""
		}

		// Fail-closed permission gate (Layer 2/3): a clear (empty staff) needs
		// :delete; a set needs :create OR :update (the espyna use case action-
		// gates the precise create-vs-update branch as the Layer-4 backstop).
		if staffID == "" {
			if !perms.Can(sgppsEntity, "delete") {
				return view.HTMXError(deps.Labels.Errors.Unauthorized)
			}
		} else if !perms.CanAny(sgppsEntity+":create", sgppsEntity+":update") {
			return view.HTMXError(deps.Labels.Errors.Unauthorized)
		}
		if productPlanID == "" {
			return view.HTMXError(deps.Labels.Errors.UpdateFailed)
		}
		if deps.AssignGroupServicer == nil {
			return view.HTMXError(deps.Labels.Errors.UpdateFailed)
		}

		if _, err := deps.AssignGroupServicer(ctx, groupID, productPlanID, staffID, role); err != nil {
			log.Printf("staff tab: assign servicer (group=%s plan=%s): %v", groupID, productPlanID, err)
			return view.HTMXError(err.Error())
		}

		// Close the drawer + refresh the whole staff tab panel (the table, whose
		// Status chips reflect the new coverage). refreshTable resolves the panel by
		// id (staffTabID) and re-GETs its data-refresh-url (the tab-action for tab=staff).
		return view.HTMXSuccess(staffTabID)
	})
}

// findRow returns the tab row for one offering.
func findRow(rows []SectionAssignmentRow, productPlanID string) (SectionAssignmentRow, bool) {
	for _, r := range rows {
		if r.ProductPlanID == productPlanID {
			return r, true
		}
	}
	return SectionAssignmentRow{}, false
}

// =============================================================================
// M4 row-source flip (plan.md §2, centymo.md §3) — everything below reads
// active subscription_group_product_plan (class) rows instead of deriving
// offering rows from the plan. See the package doc comment at the top of this
// file for the legacy/live split.
// =============================================================================

// SectionSGPPRow is one row of the M4 section tab: an offering (a class, once
// materialized) with its variant badge, exception-only exclude tag, and ALL of
// its active assignments rendered comma-separated with phase labels
// (centymo.md §3 — "nothing is hidden, first-wins retired"). Materialized is
// false only for the transitional-fallback rows a zero-class section renders
// (derived from the plan universe, no sgpp row exists yet); every row action
// is disabled on those rows.
type SectionSGPPRow struct {
	ID              string // subscription_group_product_plan id; "" on a fallback row
	ProductPlanID   string
	OfferingLabel   string // product_plan.name (data, not a lyngua label)
	VariantLabel    string // "" when the offering carries no variant (umbrella)
	Excluded        bool   // sgpp.status == EXCLUDED
	Materialized    bool   // false ⇒ transitional fallback row
	AssignmentsText string // comma-separated "Name (Phase)"; "" when unstaffed
	InUse           bool   // gates the guarded Remove action

	// -- Umbrella grouping (the owner's "one row for Arts, one for Design"):
	// product identity is what lets two strand offerings of the SAME product
	// collapse into one rendered row with a chip per variant. ProductID is the
	// grouping key (product_plan.product_id); ProductLabel is the umbrella's
	// own name IF the adapter hydrated product_plan.product (it does not
	// today), otherwise "" and the label is derived from the members' offering
	// names — see productDisplayLabel.
	ProductID    string
	ProductLabel string

	// assignments is the row's resolved assignment entries, kept so a grouped
	// row can merge its members' teachers through the SAME joinAssignments
	// composer (one joiner, one lyngua separator) instead of re-splitting the
	// already-joined AssignmentsText.
	assignments []sgppAssignmentDisplay
}

// SectionSGPPTabData is the whole M4 tab payload.
type SectionSGPPTabData struct {
	SubscriptionGroupID string
	PickerURL           string // primary action + empty-state CTA (the S2 add-offerings picker)
	RefreshURL          string // tab-action GET for tab=staff, refreshes this panel
	Rows                []SectionSGPPRow
	Table               *types.TableConfig
	Authorized          bool // false ⇒ render the read-only "no permission" note
	Labels              subscription_group.Labels
}

// buildSGPPTabData assembles the M4 tab for one section: active class rows
// (or, transitionally, the derived offering universe when none exist yet).
// Fail-closed on sgppEntity:list INSIDE the tab body (a note, never a
// full-page Forbidden — which would be wrong for an HTMX tab swap).
func buildSGPPTabData(ctx context.Context, deps *DetailViewDeps, sg *subscriptiongrouppb.SubscriptionGroup, l subscription_group.Labels) *SectionSGPPTabData {
	groupID := sg.GetId()
	data := &SectionSGPPTabData{
		SubscriptionGroupID: groupID,
		RefreshURL:          route.ResolveURL(deps.Routes.TabActionURL, "id", groupID, "tab", "staff"),
		Labels:              l,
	}
	if deps.SGPPPickerURL != nil {
		data.PickerURL = deps.SGPPPickerURL(groupID)
	}

	perms := view.GetUserPermissions(ctx)
	if perms == nil || !perms.Can(sgppEntity, "list") {
		return data
	}
	data.Authorized = true

	data.Rows = listSGPPRows(ctx, deps, sg, l)
	data.Table = buildSGPPTable(deps, data, perms, l)
	// Stamp WorkspaceID/Nonce directly: the generic pipeline injector
	// (pyeza render.Pipeline.InjectPageData -> injectTableConfigContext) only
	// recurses into ANONYMOUS embedded struct fields when hunting for a
	// nested *types.TableConfig — it does NOT recurse through a NAMED field
	// like PageData.Staff (*SectionSGPPTabData). Table therefore never
	// received its WorkspaceID/Nonce via that generic pass (verified live:
	// the rendered tab's table-card carried neither data-ws-id nor
	// data-action-tokens), which left {{rowActionTokens}} signing an empty
	// map and every row action needing a POST (Exclude/Restore/Remove)
	// failing the action_workspace_guard with 409 "Workspace context missing
	// from form; please reload" — a real, user-facing break, not a test
	// artifact. Every OTHER working table-card in this app has its
	// *types.TableConfig as a TOP-LEVEL field of the struct passed to
	// view.OK, where the generic injector's direct type-match already
	// applies; this class table is the first one nested two levels deep
	// (PageData.Staff.Table), so it needs this explicit stamp.
	if data.Table != nil {
		if data.Table.WorkspaceID == "" {
			data.Table.WorkspaceID = consumer.GetWorkspaceIDFromContext(ctx)
		}
		if data.Table.Nonce == "" {
			data.Table.Nonce = render.NonceFromContext(ctx)
		}
	}
	return data
}

// countSectionSGPPs returns the section tab's M4 count badge: non-excluded
// class rows, or (transitionally) the derived offering count when the section
// has no class rows yet. Mirrors buildSGPPTabData's row source + read gate
// exactly so the badge can never drift from the rendered rows. CHEAP by
// design — never a full buildSGPPTabData run.
func countSectionSGPPs(ctx context.Context, deps *DetailViewDeps, sg *subscriptiongrouppb.SubscriptionGroup) int {
	perms := view.GetUserPermissions(ctx)
	if perms == nil || !perms.Can(sgppEntity, "list") {
		return 0
	}
	sgpps := listActiveSGPPs(ctx, deps, sg.GetId())
	if len(sgpps) > 0 {
		// Count RENDERED rows, not class rows: two strands of one umbrella
		// product collapse into a single row, so the badge must count groups.
		// Shares sgppDisplayRows + groupSGPPRows with the table so the badge
		// cannot drift from what is on screen; still cheap (two batched
		// LIST_IN calls, no assignment/phase/name/in-use fan-out).
		count := 0
		for _, g := range groupSGPPRows(sgppDisplayRows(ctx, deps, sgpps)) {
			for _, m := range g.Members {
				if !m.Excluded {
					count++
					break
				}
			}
		}
		return count
	}
	planID := sg.GetPlanId()
	if planID == "" || deps.ListProductPlans == nil {
		return 0
	}
	return len(listPlanOfferings(ctx, deps, planID))
}

// listActiveSGPPs returns the section's active class rows (both ACTIVE and
// EXCLUDED status — EXCLUDED rows stay visible so they can be restored,
// centymo.md §3).
func listActiveSGPPs(ctx context.Context, deps *DetailViewDeps, groupID string) []*sgpppb.SubscriptionGroupProductPlan {
	if deps.ListSubscriptionGroupProductPlans == nil || groupID == "" {
		return nil
	}
	resp, err := deps.ListSubscriptionGroupProductPlans(ctx, &sgpppb.ListSubscriptionGroupProductPlansRequest{
		Filters: andFilters(stringEqFilter("subscription_group_id", groupID), boolEqFilter("active", true)),
	})
	if err != nil {
		log.Printf("class tab: list sgpps for section %s: %v", groupID, err)
		return nil
	}
	out := make([]*sgpppb.SubscriptionGroupProductPlan, 0, len(resp.GetData()))
	for _, c := range resp.GetData() {
		if c != nil && c.GetActive() && c.GetSubscriptionGroupId() == groupID {
			out = append(out, c)
		}
	}
	return out
}

// listSGPPRows builds the tab's row set: active class rows resolved to
// display form, or the transitional fallback when none exist yet. The base
// display projection (sgppDisplayRows) is shared with the count badge; only
// the table path pays for the assignment/phase/name/in-use fan-out below.
func listSGPPRows(ctx context.Context, deps *DetailViewDeps, sg *subscriptiongrouppb.SubscriptionGroup, l subscription_group.Labels) []SectionSGPPRow {
	groupID := sg.GetId()
	sgpps := listActiveSGPPs(ctx, deps, groupID)
	if len(sgpps) == 0 {
		return fallbackSGPPRows(ctx, deps, sg)
	}

	rows := sgppDisplayRows(ctx, deps, sgpps)

	sgppIDs := make([]string, 0, len(sgpps))
	templateIDs := make([]string, 0, len(sgpps))
	for _, c := range sgpps {
		sgppIDs = append(sgppIDs, c.GetId())
		if tid := c.GetJobTemplateId(); tid != "" {
			templateIDs = append(templateIDs, tid)
		}
	}

	phasesByID := resolveSGPPPhases(ctx, deps, templateIDs)
	names := staffNames(ctx, deps)
	assignments := listSGPPAssignments(ctx, deps, sgppIDs, phasesByID, names)

	var inUseIDs map[string]bool
	if deps.GetSubscriptionGroupProductPlanInUseIDs != nil {
		inUseIDs, _ = deps.GetSubscriptionGroupProductPlanInUseIDs(ctx, sgppIDs)
	}

	for i := range rows {
		r := &rows[i]
		r.assignments = assignments[r.ID]
		r.AssignmentsText = joinAssignments(r.assignments, l)
		r.InUse = inUseIDs[r.ID]
	}
	sort.SliceStable(rows, func(i, j int) bool {
		a, b := strings.ToLower(rows[i].OfferingLabel), strings.ToLower(rows[j].OfferingLabel)
		if a == b {
			return rows[i].ID < rows[j].ID
		}
		return a < b
	})
	return rows
}

// sgppDisplayRows is the CHEAP base projection of a section's class rows —
// exactly two batched LIST_IN calls (offerings, variants), no assignment/
// phase/name/in-use fan-out — shared by the table builder (listSGPPRows
// enriches it) and the count badge (countSectionSGPPs groups it as-is), so
// the two can never drift.
func sgppDisplayRows(ctx context.Context, deps *DetailViewDeps, sgpps []*sgpppb.SubscriptionGroupProductPlan) []SectionSGPPRow {
	planIDs := make([]string, 0, len(sgpps))
	for _, c := range sgpps {
		planIDs = append(planIDs, c.GetProductPlanId())
	}
	plans := resolveProductPlansByID(ctx, deps, planIDs)
	variantIDs := make([]string, 0, len(plans))
	for _, pp := range plans {
		if v := pp.GetProductVariantId(); v != "" {
			variantIDs = append(variantIDs, v)
		}
	}
	variantNames := resolveVariantNames(ctx, deps, variantIDs)

	rows := make([]SectionSGPPRow, 0, len(sgpps))
	for _, c := range sgpps {
		pp := plans[c.GetProductPlanId()]
		rows = append(rows, SectionSGPPRow{
			ID:            c.GetId(),
			ProductPlanID: c.GetProductPlanId(),
			OfferingLabel: offeringLabelFromPlan(pp, c.GetProductPlanId()),
			VariantLabel:  variantNames[pp.GetProductVariantId()],
			Excluded:      c.GetStatus() == sgpppb.SubscriptionGroupProductPlanStatus_SUBSCRIPTION_GROUP_PRODUCT_PLAN_STATUS_EXCLUDED,
			Materialized:  true,
			ProductID:     pp.GetProductId(),
			ProductLabel:  strings.TrimSpace(pp.GetProduct().GetName()),
		})
	}
	return rows
}

// fallbackSGPPRows renders the TRANSITIONAL fallback (plan.md §2 row-source
// fallback, pre-M3-backfill): today's derived plan-offering universe, with
// every action disabled + a "not yet set up" hint. Reuses the legacy
// listPlanOfferings/offeringLabel helpers verbatim — zero duplicate logic.
func fallbackSGPPRows(ctx context.Context, deps *DetailViewDeps, sg *subscriptiongrouppb.SubscriptionGroup) []SectionSGPPRow {
	planID := sg.GetPlanId()
	if planID == "" || deps.ListProductPlans == nil {
		return nil
	}
	plans := listPlanOfferings(ctx, deps, planID)
	rows := make([]SectionSGPPRow, 0, len(plans))
	for _, pp := range plans {
		rows = append(rows, SectionSGPPRow{
			ProductPlanID: pp.GetId(),
			OfferingLabel: offeringLabel(pp),
			Materialized:  false,
			ProductID:     pp.GetProductId(),
			ProductLabel:  strings.TrimSpace(pp.GetProduct().GetName()),
		})
	}
	return rows
}

// resolveProductPlansByID resolves exactly the referenced offerings (LIST_IN
// by id), keyed by id.
func resolveProductPlansByID(ctx context.Context, deps *DetailViewDeps, ids []string) map[string]*productplanpb.ProductPlan {
	out := map[string]*productplanpb.ProductPlan{}
	if len(ids) == 0 || deps.ListProductPlans == nil {
		return out
	}
	resp, err := deps.ListProductPlans(ctx, &productplanpb.ListProductPlansRequest{Filters: stringInFilter("id", ids)})
	if err != nil {
		log.Printf("class tab: list offerings: %v", err)
		return out
	}
	for _, pp := range resp.GetData() {
		if pp != nil {
			out[pp.GetId()] = pp
		}
	}
	return out
}

// resolveVariantNames resolves exactly the referenced variants (LIST_IN by
// id), keyed by id. Label = SKU, falling back to the id.
func resolveVariantNames(ctx context.Context, deps *DetailViewDeps, variantIDs []string) map[string]string {
	out := map[string]string{}
	if len(variantIDs) == 0 || deps.ListProductVariants == nil {
		return out
	}
	resp, err := deps.ListProductVariants(ctx, &productvariantpb.ListProductVariantsRequest{Filters: stringInFilter("id", variantIDs)})
	if err != nil {
		log.Printf("class tab: list variants: %v", err)
		return out
	}
	for _, v := range resp.GetData() {
		if v == nil {
			continue
		}
		label := v.GetSku()
		if label == "" {
			label = v.GetId()
		}
		out[v.GetId()] = label
	}
	return out
}

// sgppPhase is a job_template_phase's display projection (id -> name/order).
type sgppPhase struct {
	Name  string
	Order int32
}

// resolveSGPPPhases resolves exactly the referenced templates' active phases
// (LIST_IN by job_template_id), keyed by phase id — the same source the sgpp
// module's own coverage.go reads (job-sampling ban, plan.md §2.5: never a
// member job).
func resolveSGPPPhases(ctx context.Context, deps *DetailViewDeps, templateIDs []string) map[string]sgppPhase {
	out := map[string]sgppPhase{}
	if len(templateIDs) == 0 || deps.ListJobTemplatePhases == nil {
		return out
	}
	resp, err := deps.ListJobTemplatePhases(ctx, &jobtemplatephasepb.ListJobTemplatePhasesRequest{
		Filters: andFilters(stringInFilter("job_template_id", templateIDs), boolEqFilter("active", true)),
	})
	if err != nil {
		log.Printf("class tab: list phases: %v", err)
		return out
	}
	for _, p := range resp.GetData() {
		if p != nil {
			out[p.GetId()] = sgppPhase{Name: p.GetName(), Order: p.GetPhaseOrder()}
		}
	}
	return out
}

// sgppAssignmentDisplay is one resolved assignment entry ready to join into
// a row's comma-separated Teachers text.
type sgppAssignmentDisplay struct {
	staffName  string
	phaseLabel string // "" for a phase-NULL (whole-coverage) assignment — no parenthetical rendered
	order      int32
}

// listSGPPAssignments resolves ALL active assignments for the given sgpps
// (filtered on the v2 FK f12, espyna.md §1b), keyed by class id. Staff names
// resolve via f13 (product_plan_staff_id) -> pps.staff_id -> name — NEVER
// legacy f10, so the view survives M7.
func listSGPPAssignments(ctx context.Context, deps *DetailViewDeps, sgppIDs []string, phasesByID map[string]sgppPhase, names map[string]string) map[string][]sgppAssignmentDisplay {
	out := map[string][]sgppAssignmentDisplay{}
	if len(sgppIDs) == 0 || deps.ListSubscriptionGroupProductPlanStaffs == nil {
		return out
	}
	resp, err := deps.ListSubscriptionGroupProductPlanStaffs(ctx, &sgppspb.ListSubscriptionGroupProductPlanStaffsRequest{
		Filters: andFilters(stringInFilter("subscription_group_product_plan_id", sgppIDs), boolEqFilter("active", true)),
	})
	if err != nil {
		log.Printf("class tab: list assignments: %v", err)
		return out
	}
	data := resp.GetData()

	ppsIDs := make([]string, 0, len(data))
	seen := map[string]bool{}
	for _, a := range data {
		if a == nil {
			continue
		}
		if id := a.GetProductPlanStaffId(); id != "" && !seen[id] {
			seen[id] = true
			ppsIDs = append(ppsIDs, id)
		}
	}
	ppsStaffByID := resolvePPSStaffIDs(ctx, deps, ppsIDs)

	for _, a := range data {
		if a == nil || !a.GetActive() {
			continue
		}
		sgppID := a.GetSubscriptionGroupProductPlanId()
		if sgppID == "" {
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

		phaseLabel, order := "", int32(-1)
		if pid := a.GetJobTemplatePhaseId(); pid != "" {
			if p, ok := phasesByID[pid]; ok {
				phaseLabel, order = p.Name, p.Order
			} else {
				phaseLabel = pid
			}
		}
		out[sgppID] = append(out[sgppID], sgppAssignmentDisplay{staffName: staffName, phaseLabel: phaseLabel, order: order})
	}
	return out
}

// resolvePPSStaffIDs resolves exactly the referenced eligibility rows
// (LIST_IN by id — espyna.md §1b call #6), keyed by pps id -> staff id.
func resolvePPSStaffIDs(ctx context.Context, deps *DetailViewDeps, ppsIDs []string) map[string]string {
	out := map[string]string{}
	if len(ppsIDs) == 0 || deps.ListProductPlanStaffs == nil {
		return out
	}
	resp, err := deps.ListProductPlanStaffs(ctx, &productplanstaffpb.ListProductPlanStaffsRequest{Filters: stringInFilter("id", ppsIDs)})
	if err != nil {
		log.Printf("class tab: resolve eligibility rows: %v", err)
		return out
	}
	for _, pps := range resp.GetData() {
		if pps != nil {
			out[pps.GetId()] = pps.GetStaffId()
		}
	}
	return out
}

// joinAssignments renders a class's assignments as the comma-separated
// Teachers text (centymo.md §3 — phase-scoped entries carry a parenthetical,
// phase-NULL entries render bare). Sorted by phase order (phase-NULL first),
// then staff name.
func joinAssignments(entries []sgppAssignmentDisplay, l subscription_group.Labels) string {
	if len(entries) == 0 {
		return ""
	}
	sort.SliceStable(entries, func(i, j int) bool {
		if entries[i].order != entries[j].order {
			return entries[i].order < entries[j].order
		}
		return strings.ToLower(entries[i].staffName) < strings.ToLower(entries[j].staffName)
	})
	sep := l.Staff.ListSeparator
	if sep == "" {
		sep = ", "
	}
	parts := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.phaseLabel != "" {
			parts = append(parts, e.staffName+" ("+e.phaseLabel+")")
		} else {
			parts = append(parts, e.staffName)
		}
	}
	return strings.Join(parts, sep)
}

// offeringLabelFromPlan mirrors offeringLabel but tolerates a nil plan (an
// unresolved product_plan id), falling back to the raw id.
func offeringLabelFromPlan(pp *productplanpb.ProductPlan, fallbackID string) string {
	if pp == nil {
		return fallbackID
	}
	return offeringLabel(pp)
}

// buildSGPPTable maps the M4 rows into the shared table-card config: two
// columns (Subject, Teachers), the S2 picker as the toolbar primary action,
// and per-row View/Assign/Exclude-Restore/Remove actions.
//
// Umbrella grouping (the owner's "one row per product"): class rows are first
// grouped by product identity (groupSGPPRows). A group of ONE renders through
// sgppTableRow untouched — byte-identical to the ungrouped tab, so every
// single-offering product (and all of a plan whose offerings are un-varianted)
// is a strict no-op. Only a multi-offering product collapses into the grouped
// row (sgppGroupTableRow): one chip per variant, merged Teachers text, ONE
// View (the lead member's class page), and a whole-product Remove.
func buildSGPPTable(deps *DetailViewDeps, data *SectionSGPPTabData, perms *types.UserPermissions, l subscription_group.Labels) *types.TableConfig {
	columns := []types.TableColumn{
		{Key: "subject", Label: l.Staff.ColumnSubject, NoSort: true, NoFilter: true, Width: "40%"},
		{Key: "assignments", Label: l.Staff.ColumnServicer, NoSort: true, NoFilter: true},
	}
	groups := groupSGPPRows(data.Rows)
	rows := make([]types.TableRow, 0, len(groups))
	for _, g := range groups {
		if len(g.Members) == 1 {
			rows = append(rows, sgppTableRow(deps, data.SubscriptionGroupID, g.Members[0], perms, l))
			continue
		}
		rows = append(rows, sgppGroupTableRow(deps, data.SubscriptionGroupID, g, perms, l))
	}
	cfg := &types.TableConfig{
		ID:          sgppTableID,
		Columns:     columns,
		Rows:        rows,
		Labels:      deps.TableLabels,
		EmptyState:  types.TableEmptyState{Title: l.Staff.EmptyTitle, Message: l.Staff.EmptyMessage},
		ShowSearch:  true,
		ShowColumns: true,
		ShowDensity: true,
		ShowEntries: true,
		ShowActions: true,
	}
	if data.PickerURL != "" {
		cfg.PrimaryAction = &types.PrimaryAction{
			Label:           l.Staff.AddAction,
			ActionURL:       data.PickerURL,
			Icon:            "icon-plus",
			Disabled:        !perms.Can(sgppEntity, "create"),
			DisabledTooltip: l.Staff.Unauthorized,
		}
	}
	types.ApplyColumnStyles(columns, rows)
	types.ApplyTableSettings(cfg)
	return cfg
}

// sgppSubjectCell renders the Subject column: offering name + variant badge
// (strand only) + the exception-only signal — either the EXCLUDED tag or the
// transitional-fallback "not yet set up" hint (never both; centymo.md §3 —
// default-state badges are noise). Every interpolated value is a trusted
// composed label (offering/variant names, lyngua labels) and HTML-escaped.
func sgppSubjectCell(r SectionSGPPRow, l subscription_group.Labels) types.TableCell {
	var b strings.Builder
	b.WriteString(`<span class="centymo-sgpp-subject-name">`)
	b.WriteString(template.HTMLEscapeString(r.OfferingLabel))
	b.WriteString(`</span>`)
	if r.VariantLabel != "" {
		b.WriteString(` <span class="badge badge--muted">`)
		b.WriteString(template.HTMLEscapeString(r.VariantLabel))
		b.WriteString(`</span>`)
	}
	switch {
	case !r.Materialized:
		b.WriteString(` <span class="text-muted centymo-sgpp-unmaterialized-hint" data-testid="sgpp-unmaterialized-` + template.HTMLEscapeString(r.ProductPlanID) + `" title="`)
		b.WriteString(template.HTMLEscapeString(l.Staff.NotMaterializedHint))
		b.WriteString(`">`)
		b.WriteString(template.HTMLEscapeString(l.Staff.NotMaterializedHint))
		b.WriteString(`</span>`)
	case r.Excluded:
		b.WriteString(` <span class="badge badge--warning" data-testid="sgpp-excluded-tag-` + template.HTMLEscapeString(r.ID) + `">`)
		b.WriteString(template.HTMLEscapeString(l.Staff.ExcludedTag))
		b.WriteString(`</span>`)
	}
	return types.TableCell{Type: "html", HTML: template.HTML(b.String()), Value: r.OfferingLabel}
}

// sgppAssignmentsCell renders the Teachers column: the comma-separated
// assignment text with an aria-label carrying the full list (I-1 display
// truth), the unstaffed placeholder, or an em-dash on a fallback row.
func sgppAssignmentsCell(r SectionSGPPRow, l subscription_group.Labels) types.TableCell {
	if !r.Materialized {
		return types.TableCell{Type: "text", Value: emDash}
	}
	if r.AssignmentsText == "" {
		return types.TableCell{Type: "text", Value: l.Staff.Unstaffed}
	}
	aria := strings.ReplaceAll(l.Staff.ListAria, "{{names}}", r.AssignmentsText)
	var b strings.Builder
	b.WriteString(`<span data-testid="sgpp-row-assignments-`)
	b.WriteString(template.HTMLEscapeString(r.ID))
	b.WriteString(`" aria-label="`)
	b.WriteString(template.HTMLEscapeString(aria))
	b.WriteString(`">`)
	b.WriteString(template.HTMLEscapeString(r.AssignmentsText))
	b.WriteString(`</span>`)
	return types.TableCell{Type: "html", HTML: template.HTML(b.String()), Value: r.AssignmentsText}
}

// sgppRowKey returns the row's stable identity for test ids / the table row
// ID — the class id when materialized, else the offering id (fallback rows
// have no class id yet).
func sgppRowKey(r SectionSGPPRow) string {
	if r.ID != "" {
		return r.ID
	}
	return r.ProductPlanID
}

// sgppTableRow builds one M4 table row. Fallback (not-yet-materialized) rows
// render both actions disabled with the not-yet-materialized hint as the
// tooltip (centymo.md §3 — "actions disabled + a subtle not-yet-materialized
// hint"). Materialized rows wire View (navigate to S4) / Assign (opens the S3
// quick drawer) inline, and Exclude-or-Restore / Remove into the row's
// overflow (⋮) menu.
func sgppTableRow(deps *DetailViewDeps, groupID string, r SectionSGPPRow, perms *types.UserPermissions, l subscription_group.Labels) types.TableRow {
	key := sgppRowKey(r)
	row := types.TableRow{
		ID: key,
		DataAttrs: map[string]string{
			"testid": "sgpp-row-" + key,
			"status": sgppRowStatus(r),
		},
		Cells: []types.TableCell{sgppSubjectCell(r, l), sgppAssignmentsCell(r, l)},
	}

	if !r.Materialized {
		row.Actions = []types.TableAction{
			{Type: "view", Action: "view", Label: l.Staff.ViewAction, Disabled: true, DisabledTooltip: l.Staff.NotMaterializedHint, TestID: "sgpp-view-" + key},
			{Type: "edit", Action: "edit", Label: l.Staff.AssignAction, DrawerTitle: l.Staff.AssignAction, Disabled: true, DisabledTooltip: l.Staff.NotMaterializedHint, TestID: "sgpp-assign-" + key},
		}
		return row
	}

	canAssign := perms.CanAny(sgppsEntity+":create", sgppsEntity+":update")
	canManage := perms.Can(sgppEntity, "update")
	canDelete := perms.Can(sgppEntity, "delete")

	viewAction := types.TableAction{Type: "view", Action: "view", Label: l.Staff.ViewAction, TestID: "sgpp-view-" + r.ID}
	if deps.SGPPDetailURL != nil {
		viewAction.Href = deps.SGPPDetailURL(groupID, r.ID)
	} else {
		viewAction.Disabled = true
	}

	assignAction := types.TableAction{
		Type: "edit", Action: "edit", Label: l.Staff.AssignAction, DrawerTitle: l.Staff.AssignAction,
		TestID: "sgpp-assign-" + r.ID, Disabled: !canAssign, DisabledTooltip: l.Staff.Unauthorized,
	}
	if deps.SGPPAssignURL != nil {
		assignAction.URL = deps.SGPPAssignURL(r.ID)
	} else {
		assignAction.Disabled = true
	}

	row.Actions = append(row.Actions, viewAction, assignAction)

	if r.Excluded {
		restoreAction := types.TableAction{
			Type: "activate", Action: "activate", Label: l.Staff.RestoreAction, Overflow: true,
			ItemName: r.OfferingLabel, ConfirmTitle: l.Staff.RestoreConfirmTitle,
			ConfirmMessage:  fmt.Sprintf(l.Staff.RestoreConfirmMessage, r.OfferingLabel),
			TestID:          "sgpp-restore-" + r.ID,
			Disabled:        !canManage,
			DisabledTooltip: l.Staff.Unauthorized,
		}
		if deps.SGPPSetStatusURL != nil {
			restoreAction.URL = deps.SGPPSetStatusURL(r.ID, "active")
		} else {
			restoreAction.Disabled = true
		}
		row.Actions = append(row.Actions, restoreAction)
	} else {
		excludeAction := types.TableAction{
			Type: "deactivate", Action: "deactivate", Label: l.Staff.ExcludeAction, Overflow: true,
			ItemName: r.OfferingLabel, ConfirmTitle: l.Staff.ExcludeConfirmTitle,
			ConfirmMessage:  fmt.Sprintf(l.Staff.ExcludeConfirmMessage, r.OfferingLabel),
			TestID:          "sgpp-exclude-" + r.ID,
			Disabled:        !canManage,
			DisabledTooltip: l.Staff.Unauthorized,
		}
		if deps.SGPPSetStatusURL != nil {
			excludeAction.URL = deps.SGPPSetStatusURL(r.ID, "excluded")
		} else {
			excludeAction.Disabled = true
		}
		row.Actions = append(row.Actions, excludeAction)
	}

	removeAction := types.TableAction{
		Type: "delete", Action: "delete", Label: l.Staff.RemoveAction, Overflow: true,
		URL: deps.SGPPDeleteURL, ItemName: r.OfferingLabel,
		ConfirmTitle: l.Staff.RemoveConfirmTitle, ConfirmMessage: fmt.Sprintf(l.Staff.RemoveConfirmMessage, r.OfferingLabel),
		TestID: "sgpp-remove-" + r.ID,
	}
	switch {
	case r.InUse:
		removeAction.Disabled = true
		removeAction.DisabledTooltip = l.Staff.RemoveBlockedInUse
	case !canDelete:
		removeAction.Disabled = true
		removeAction.DisabledTooltip = l.Staff.Unauthorized
	case deps.SGPPDeleteURL == "":
		removeAction.Disabled = true
	}
	row.Actions = append(row.Actions, removeAction)

	return row
}

// sgppRowStatus renders the row's data-status attribute: "unmaterialized"
// for a fallback row, else the sgpp status token.
func sgppRowStatus(r SectionSGPPRow) string {
	if !r.Materialized {
		return "unmaterialized"
	}
	if r.Excluded {
		return "excluded"
	}
	return "active"
}

// boolEqFilter builds a single field == value boolean FilterRequest.
func boolEqFilter(field string, value bool) *commonpb.FilterRequest {
	return &commonpb.FilterRequest{Filters: []*commonpb.TypedFilter{{
		Field:      field,
		FilterType: &commonpb.TypedFilter_BooleanFilter{BooleanFilter: &commonpb.BooleanFilter{Value: value}},
	}}}
}

// andFilters merges filter requests into one (all conditions AND-ed — mirrors
// how the postgres adapter treats a Filters slice).
func andFilters(reqs ...*commonpb.FilterRequest) *commonpb.FilterRequest {
	out := &commonpb.FilterRequest{}
	for _, r := range reqs {
		if r == nil {
			continue
		}
		out.Filters = append(out.Filters, r.Filters...)
	}
	return out
}
