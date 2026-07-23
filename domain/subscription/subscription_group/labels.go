package subscription_group

import "strings"

// ---------------------------------------------------------------------------
// Subscription Group (section / cohort) labels
// ---------------------------------------------------------------------------

// Labels holds all labels for the subscription_group module.
type Labels struct {
	Page    PageLabels     `json:"page"`
	Buttons ButtonLabels   `json:"buttons"`
	Columns ColumnLabels   `json:"columns"`
	Empty   EmptyLabels    `json:"empty"`
	Form    FormLabels     `json:"form"`
	Bulk    BulkLabels     `json:"bulk"`
	Confirm ConfirmLabels  `json:"confirm"`
	Tabs    TabLabels      `json:"tabs"`
	Staff   StaffTabLabels `json:"staff"`
	Detail  DetailLabels   `json:"detail"`
	Errors  ErrorLabels    `json:"errors"`
}

type PageLabels struct {
	Title         string `json:"title"`
	Subtitle      string `json:"subtitle"`
	ActiveTitle   string `json:"active_title"`
	InactiveTitle string `json:"inactive_title"`
}

type ButtonLabels struct {
	View       string `json:"view"`
	Add        string `json:"add"`
	Edit       string `json:"edit"`
	Delete     string `json:"delete"`
	BulkDelete string `json:"bulk_delete"`
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`
}

type ColumnLabels struct {
	Name          string `json:"name"`
	Kind          string `json:"kind"`
	Capacity      string `json:"capacity"`
	PriceSchedule string `json:"price_schedule"` // education value: "Academic Year"
	Client        string `json:"client"`         // subscriptions tab; education value: "Student"
	Subscription  string `json:"subscription"`   // subscriptions tab; education value: "Enrollment"
	Status        string `json:"status"`
	DateCreated   string `json:"date_created"`
	Actions       string `json:"actions"`
}

type EmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

// FormLabels holds the drawer-form field labels. Section/cohort vocabulary —
// a subscription_group is a class roster / patient panel / project team
// anchored to a program (plan) and a period (price_schedule).
type FormLabels struct {
	SectionIdentity   string `json:"section_identity"`
	SectionAnchors    string `json:"section_anchors"`
	SectionCapacity   string `json:"section_capacity"`
	Name              string `json:"name"`
	NamePlaceholder   string `json:"name_placeholder"`
	NameInfo          string `json:"name_info"`
	Kind              string `json:"kind"`
	KindPlaceholder   string `json:"kind_placeholder"`
	KindInfo          string `json:"kind_info"`
	KindCohort        string `json:"kind_cohort"`
	KindRoster        string `json:"kind_roster"`
	KindPanel         string `json:"kind_panel"`
	KindProjectTeam   string `json:"kind_project_team"`
	Plan              string `json:"plan"`
	PlanPlaceholder   string `json:"plan_placeholder"`
	PlanSearch        string `json:"plan_search"`
	PlanInfo          string `json:"plan_info"`
	PriceSchedule     string `json:"price_schedule"`
	PriceSchedulePH   string `json:"price_schedule_placeholder"`
	PriceScheduleSrch string `json:"price_schedule_search"`
	PriceScheduleInfo string `json:"price_schedule_info"`
	CapacityMode      string `json:"capacity_mode"`
	CapacityModeInfo  string `json:"capacity_mode_info"`
	CapUnlimited      string `json:"cap_unlimited"`
	CapClosed         string `json:"cap_closed"`
	CapCapped         string `json:"cap_capped"`
	MaxCapacity       string `json:"max_capacity"`
	MaxCapacityPH     string `json:"max_capacity_placeholder"`
	MaxCapacityInfo   string `json:"max_capacity_info"`
	Active            string `json:"active"`
	ActiveInfo        string `json:"active_info"`
}

type BulkLabels struct {
	DeleteTitle       string `json:"delete_title"`
	DeleteMessage     string `json:"delete_message"`
	ActivateTitle     string `json:"activate_title"`
	ActivateMessage   string `json:"activate_message"`
	DeactivateTitle   string `json:"deactivate_title"`
	DeactivateMessage string `json:"deactivate_message"`
}

type ConfirmLabels struct {
	DeleteTitle       string `json:"delete_title"`
	DeleteMessage     string `json:"delete_message"`
	ActivateTitle     string `json:"activate_title"`
	ActivateMessage   string `json:"activate_message"`
	DeactivateTitle   string `json:"deactivate_title"`
	DeactivateMessage string `json:"deactivate_message"`
}

type TabLabels struct {
	Info              string `json:"info"`
	Subscriptions     string `json:"subscriptions"`      // label — education: "Enrollments"; general: "Subscriptions"
	SubscriptionsSlug string `json:"subscriptions_slug"` // URL slug — education: "enrollments"; general (empty) → "subscriptions"
	Staff             string `json:"staff"`              // tab button — general: "Staff"; education: "Teaching Staff"
	Audit             string `json:"audit"`
	Attachments       string `json:"attachments"`
}

// StaffTabLabels holds the "Teaching Staff" assignment-grid labels (§6.1). The
// grid assigns an eligible servicer to each offering (product_plan) of a cohort
// (subscription_group). Generic Go field names; vertical vocabulary (subject/
// teacher) enters only via the lyngua overrides.
type StaffTabLabels struct {
	ColumnSubject      string `json:"column_subject"`      // offering column — general "Offering"; education "Subject"; also the drawer offering field label
	ColumnServicer     string `json:"column_servicer"`     // servicer column — general "Staff"; education "Teacher"; also the drawer teacher field label
	ColumnRole         string `json:"column_role"`         // role column; also the drawer role field label
	ColumnState        string `json:"column_state"`        // assignment-state column header (Saved / Unassigned)
	EmptyPool          string `json:"empty_pool"`          // empty-pool gate link text (one row's eligible pool is empty)
	EmptyTitle         string `json:"empty_title"`         // tab empty-state title — the section has zero offerings to staff
	EmptyMessage       string `json:"empty_message"`       // tab empty-state message — the section has zero offerings to staff
	Saved              string `json:"saved"`               // assigned-row state
	Unassigned         string `json:"unassigned"`          // unassigned-row state
	AssignAction       string `json:"assign_action"`       // per-row action tooltip + assign-drawer title
	RolePrimary        string `json:"role_primary"`        // role enum label — teacher-of-record
	RoleAccess         string `json:"role_access"`         // role enum label — visibility-only access
	TeacherPlaceholder string `json:"teacher_placeholder"` // drawer teacher autocomplete placeholder
	TeacherSearch      string `json:"teacher_search"`      // drawer teacher autocomplete filter placeholder
	ClearAction        string `json:"clear_action"`        // drawer Clear control — soft-deletes the active edge (§6.1 "Clear")
	Unauthorized       string `json:"unauthorized"`        // read-only / no-permission note
}

// ResolveTabSlug returns the URL slug for a canonical tab key. The "subscriptions"
// tab re-slugs per tier (education ships "enrollments"); other tabs pass through.
func (t TabLabels) ResolveTabSlug(canonical string) string {
	if canonical == "subscriptions" {
		if s := strings.TrimSpace(t.SubscriptionsSlug); s != "" {
			return s
		}
	}
	return canonical
}

// CanonicalizeTab maps an incoming URL tab slug back to its canonical key so
// template lookups + equality checks stay tier-agnostic.
func (t TabLabels) CanonicalizeTab(slug string) string {
	if slug == "" {
		return ""
	}
	if s := strings.TrimSpace(t.SubscriptionsSlug); s != "" && slug == s {
		return "subscriptions"
	}
	return slug
}

type DetailLabels struct {
	Title          string `json:"title"`
	DateCreated    string `json:"date_created"`
	DateModified   string `json:"date_modified"`
	NoPlan         string `json:"no_plan"`
	NoSchedule     string `json:"no_schedule"`
	NoKind         string `json:"no_kind"`
	CapacityValue  string `json:"capacity_value"`   // e.g. "%d seats" (CAPPED)
	CapacityModeNF string `json:"capacity_mode_nf"` // fallback when mode unspecified
	NoSubtitle     string `json:"no_subtitle"`
}

type ErrorLabels struct {
	NotFound     string `json:"not_found"`
	LoadFailed   string `json:"load_failed"`
	Unauthorized string `json:"unauthorized"`
	CreateFailed string `json:"create_failed"`
	UpdateFailed string `json:"update_failed"`
	DeleteFailed string `json:"delete_failed"`
	InUse        string `json:"in_use"`
}

// DefaultLabels returns Labels with sensible English defaults using the
// section/cohort vocabulary. Tiers override field names via lyngua.
func DefaultLabels() Labels {
	return Labels{
		Page: PageLabels{
			Title:         "Sections",
			Subtitle:      "Manage your sections and cohorts",
			ActiveTitle:   "Active Sections",
			InactiveTitle: "Inactive Sections",
		},
		Buttons: ButtonLabels{
			View:       "View",
			Add:        "Add Section",
			Edit:       "Edit Section",
			Delete:     "Delete Section",
			BulkDelete: "Delete Sections",
			Activate:   "Activate",
			Deactivate: "Deactivate",
		},
		Columns: ColumnLabels{
			Name:          "Name",
			Kind:          "Type",
			Capacity:      "Capacity",
			PriceSchedule: "Price Schedule",
			Client:        "Client",
			Subscription:  "Subscription",
			Status:        "Status",
			DateCreated:   "Date Created",
			Actions:       "Actions",
		},
		Empty: EmptyLabels{
			Title:   "No Sections",
			Message: "No sections to display.",
		},
		Form: FormLabels{
			SectionIdentity:   "Section details",
			SectionAnchors:    "Program & period",
			SectionCapacity:   "Capacity",
			Name:              "Name",
			NamePlaceholder:   "Enter section name",
			NameInfo:          "A short display name for this section or cohort.",
			Kind:              "Type",
			KindPlaceholder:   "Select a type...",
			KindInfo:          "The shape of this cohort: a class roster, a panel, or a project team.",
			KindCohort:        "Cohort",
			KindRoster:        "Roster",
			KindPanel:         "Panel",
			KindProjectTeam:   "Project team",
			Plan:              "Program",
			PlanPlaceholder:   "Select a program...",
			PlanSearch:        "Filter...",
			PlanInfo:          "The program (plan) this section realizes. Section identity is program × period.",
			PriceSchedule:     "Period",
			PriceSchedulePH:   "Select a period...",
			PriceScheduleSrch: "Filter...",
			PriceScheduleInfo: "The billing period this section is anchored to (e.g. the academic-year price schedule).",
			CapacityMode:      "Capacity mode",
			CapacityModeInfo:  "Unlimited admits everyone; Closed admits no one; Capped enforces the seat limit below.",
			CapUnlimited:      "Unlimited",
			CapClosed:         "Closed",
			CapCapped:         "Capped",
			MaxCapacity:       "Maximum seats",
			MaxCapacityPH:     "e.g. 30",
			MaxCapacityInfo:   "Read only when the capacity mode is Capped.",
			Active:            "Active",
			ActiveInfo:        "Inactive sections are hidden from new enrollments.",
		},
		Bulk: BulkLabels{
			DeleteTitle:       "Delete Sections",
			DeleteMessage:     "Permanently delete the selected sections? This cannot be undone.",
			ActivateTitle:     "Activate Sections",
			ActivateMessage:   "Activate the selected sections?",
			DeactivateTitle:   "Deactivate Sections",
			DeactivateMessage: "Deactivate the selected sections?",
		},
		Confirm: ConfirmLabels{
			DeleteTitle:       "Delete Section",
			DeleteMessage:     "Permanently delete this section? This cannot be undone.",
			ActivateTitle:     "Activate Section",
			ActivateMessage:   "Activate {{name}}?",
			DeactivateTitle:   "Deactivate Section",
			DeactivateMessage: "Deactivate {{name}}?",
		},
		Tabs: TabLabels{
			Info:          "Info",
			Subscriptions: "Subscriptions",
			Staff:         "Staff",
			Audit:         "Audit",
			Attachments:   "Attachments",
		},
		Staff: StaffTabLabels{
			ColumnSubject:      "Offering",
			ColumnServicer:     "Staff",
			ColumnRole:         "Role",
			ColumnState:        "State",
			EmptyPool:          "No eligible staff — set eligibility first",
			EmptyTitle:         "No offerings",
			EmptyMessage:       "This section has no offerings to staff yet.",
			Saved:              "Saved",
			Unassigned:         "Unassigned",
			AssignAction:       "Assign",
			RolePrimary:        "Primary",
			RoleAccess:         "Access",
			TeacherPlaceholder: "Select staff...",
			TeacherSearch:      "Filter...",
			ClearAction:        "Clear",
			Unauthorized:       "You do not have permission to view staff assignments",
		},
		Detail: DetailLabels{
			Title:          "Section",
			DateCreated:    "Date Created",
			DateModified:   "Date Modified",
			NoPlan:         "No program",
			NoSchedule:     "No period",
			NoKind:         "—",
			CapacityValue:  "%d seats",
			CapacityModeNF: "Unlimited",
			NoSubtitle:     "No description provided",
		},
		Errors: ErrorLabels{
			NotFound:     "Section not found",
			LoadFailed:   "Failed to load section",
			Unauthorized: "You are not authorized to perform this action",
			CreateFailed: "Failed to create section",
			UpdateFailed: "Failed to update section",
			DeleteFailed: "Failed to delete section",
			InUse:        "This section is in use and cannot be deleted.",
		},
	}
}
