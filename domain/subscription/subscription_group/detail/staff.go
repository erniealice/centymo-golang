package detail

// staff.go — the "Teaching Staff" tab (§6.1–§6.3): a section-centric table that
// assigns an ELIGIBLE servicer to every offering (product_plan) of the section's
// program (plan). The tab renders the shared pyeza table-card (one row per
// offering: subject / current teacher / role / state) with NO in-panel header —
// the tab-strip chip names the tab + count and the Status column's Saved/Unassigned
// chips carry the coverage gap signal (W2g); a per-row action opens an assign
// DRAWER (offering read-only + eligible-only teacher autocomplete + primary/access
// role). The drawer POSTs the
// existing assign upsert (§6.3) and the response closes the drawer + refreshes the
// tab. All page-data is assembled from existing List use cases (§6.2, no new
// proto/DB). All vertical vocabulary (subject / teacher) lives in lyngua only.

import (
	"context"
	"html/template"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/erniealice/centymo-golang/domain/subscription/subscription_group"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	productplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_plan"
	productplanstaffpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_plan_staff"
	subscriptiongrouppb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription_group"
	sgppspb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription_group_product_plan_staff"
)

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
