package subscription_group_product_plan

// ---------------------------------------------------------------------------
// Subscription Group Product Plan (the CLASS entity) labels
//
// docs/plan/20260724-section-assignment-merged/plan.md + centymo.md. Every
// field here is a Go zero-fallback default (DefaultLabels) — the lyngua keys
// (general + education tiers) land separately per the house convention; tier
// overlays replace these strings without touching the struct shape. No
// vertical noun (Class/Subject/Semester/Teacher) appears in any Go identifier
// — only inside the DEFAULT string VALUES below, which lyngua fully replaces
// per tier.
// ---------------------------------------------------------------------------

// Labels holds all labels for the subscription_group_product_plan module.
type Labels struct {
	Page    PageLabels    `json:"page"`
	Buttons ButtonLabels  `json:"buttons"`
	Columns ColumnLabels  `json:"columns"`
	Empty   EmptyLabels   `json:"empty"`
	Form    FormLabels    `json:"form"`
	Bulk    BulkLabels    `json:"bulk"`
	Confirm ConfirmLabels `json:"confirm"`
	Tabs    TabLabels     `json:"tabs"`
	Detail  DetailLabels  `json:"detail"`
	Picker  PickerLabels  `json:"picker"`
	Assign  AssignLabels  `json:"assign"`
	Roster  RosterLabels  `json:"roster"`
	Errors  ErrorLabels   `json:"errors"`
}

// PageLabels — the S7 admin list page chrome.
type PageLabels struct {
	Title         string `json:"title"`
	Subtitle      string `json:"subtitle"`
	ActiveTitle   string `json:"active_title"`
	InactiveTitle string `json:"inactive_title"`
}

// ButtonLabels — row/toolbar actions shared across S7 and S4.
type ButtonLabels struct {
	View       string `json:"view"`
	Add        string `json:"add"`
	Edit       string `json:"edit"`
	Delete     string `json:"delete"`
	BulkDelete string `json:"bulk_delete"`
	Exclude    string `json:"exclude"`
	Restore    string `json:"restore"`
}

// ColumnLabels — the S7 admin table columns.
type ColumnLabels struct {
	SubscriptionGroupId string `json:"subscription_group_id"` // "Section"
	ProductPlanId       string `json:"product_plan_id"`       // "Offering"
	JobTemplateId       string `json:"job_template_id"`       // "Curriculum"
	Status              string `json:"status"`
	AssignmentsCount    string `json:"assignments_count"` // "Teachers"
	DateCreated         string `json:"date_created"`
	Actions             string `json:"actions"`
}

type EmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

// FormLabels — the S7 admin add/edit drawer.
type FormLabels struct {
	SectionIdentity     string `json:"section_identity"`
	SubscriptionGroupId string `json:"subscription_group_id"`
	SubscriptionGroupPH string `json:"subscription_group_placeholder"`
	SubscriptionGroupSr string `json:"subscription_group_search"`
	SubscriptionGroupIn string `json:"subscription_group_info"`
	ProductPlanId       string `json:"product_plan_id"`
	ProductPlanPH       string `json:"product_plan_placeholder"`
	ProductPlanSr       string `json:"product_plan_search"`
	ProductPlanIn       string `json:"product_plan_info"`
	JobTemplateId       string `json:"job_template_id"`
	JobTemplatePH       string `json:"job_template_placeholder"`
	JobTemplateSr       string `json:"job_template_search"`
	JobTemplateIn       string `json:"job_template_info"`
	Status              string `json:"status"`
	StatusInfo          string `json:"status_info"`
	StatusActive        string `json:"status_active"`
	StatusExcluded      string `json:"status_excluded"`
	Active              string `json:"active"`
	ActiveInfo          string `json:"active_info"`
}

type BulkLabels struct {
	DeleteTitle   string `json:"delete_title"`
	DeleteMessage string `json:"delete_message"`
}

// ConfirmLabels — S6 confirm sheets (exclude/restore/remove class, remove
// assignment). Message templates carry {{section}}/{{offering}} placeholders,
// resolved by the view layer (never string-built in lyngua).
type ConfirmLabels struct {
	DeleteTitle         string `json:"delete_title"`
	DeleteMessage       string `json:"delete_message"`
	ExcludeTitle        string `json:"exclude_title"`
	ExcludeMessage      string `json:"exclude_message"`
	ExcludeWarning      string `json:"exclude_warning"` // live-course warning (S6)
	RestoreTitle        string `json:"restore_title"`
	RestoreMessage      string `json:"restore_message"`
	RemoveBlockedInUse  string `json:"remove_blocked_in_use"`
	RemoveAssignTitle   string `json:"remove_assign_title"`
	RemoveAssignMessage string `json:"remove_assign_message"`
}

// TabLabels — the S4 class detail page tab strip. Canonical keys stay
// "info"/"staff" (no URL churn); the button label is tier-overridden
// (education renders "Teachers").
type TabLabels struct {
	Info  string `json:"info"`
	Staff string `json:"staff"`
}

// DetailLabels — the S4 Info tab.
type DetailLabels struct {
	Title            string `json:"title"`
	Header           string `json:"header"` // "{{offering}} — {{section}}" template
	Section          string `json:"section"`
	Offering         string `json:"offering"`
	Curriculum       string `json:"curriculum"`
	Period           string `json:"period"`
	Plan             string `json:"plan"`
	Status           string `json:"status"`
	AssignmentsCount string `json:"assignments_count"`
	DateCreated      string `json:"date_created"`
	DateModified     string `json:"date_modified"`
	NoGroup          string `json:"no_group"`
	NoPlan           string `json:"no_plan"`
	NoTemplate       string `json:"no_template"`
	GradeSheetLink   string `json:"grade_sheet_link"`
	VariantsPrefix   string `json:"variants_prefix"` // "umbrella · variants: "
}

// PickerLabels — S2 the add-offerings picker drawer.
type PickerLabels struct {
	Title            string `json:"title"` // "Add subjects — {{section}}"
	SectionContext   string `json:"section_context"`
	PlanContext      string `json:"plan_context"`
	PeriodContext    string `json:"period_context"`
	SelectAll        string `json:"select_all"` // "Select all ({{count}})"
	CurriculumPrefix string `json:"curriculum_prefix"`
	Ambiguous        string `json:"ambiguous"` // hint shown beside a constrained template select
	EmptyTitle       string `json:"empty_title"`
	EmptyMessage     string `json:"empty_message"`
	SubmitOne        string `json:"submit_one"`  // "Add 1 subject"
	SubmitMany       string `json:"submit_many"` // "Add {{count}} subjects"
}

// AssignLabels — the S3/S5 shared assign drawer (D-2 toggle form).
type AssignLabels struct {
	Title             string `json:"title"` // "Assign teacher — {{offering}}"
	SectionContext    string `json:"section_context"`
	SubjectContext    string `json:"subject_context"`
	CurriculumContext string `json:"curriculum_context"`
	CurrentHeading    string `json:"current_heading"`
	EmptyCurrent      string `json:"empty_current"`
	ModeAll           string `json:"mode_all"` // "All semesters"
	ModePer           string `json:"mode_per"` // "Per semester"
	AllPhasesChip     string `json:"all_phases_chip"`
	Staff             string `json:"staff"`
	StaffPH           string `json:"staff_placeholder"`
	StaffSearch       string `json:"staff_search"`
	Role              string `json:"role"`
	RolePrimary       string `json:"role_primary"`   // "Teacher of record"
	RoleSecondary     string `json:"role_secondary"` // "Co-teacher"
	ClearAction       string `json:"clear_action"`
	ReadOnlySlot      string `json:"read_only_slot"` // "Manage on the class page →"
	Unauthorized      string `json:"unauthorized"`
}

// RosterLabels — the S5 Teachers tab (the class page's multi-staff surface).
type RosterLabels struct {
	ColumnStaff  string `json:"column_staff"`
	ColumnRole   string `json:"column_role"`
	ColumnPhase  string `json:"column_phase"`
	EmptyTitle   string `json:"empty_title"`
	EmptyMessage string `json:"empty_message"`
	EmptyPool    string `json:"empty_pool"`
	AddAction    string `json:"add_action"`
	EditAction   string `json:"edit_action"`
	RemoveAction string `json:"remove_action"`
	Unauthorized string `json:"unauthorized"`
}

type ErrorLabels struct {
	NotFound         string `json:"not_found"`
	LoadFailed       string `json:"load_failed"`
	Unauthorized     string `json:"unauthorized"`
	CreateFailed     string `json:"create_failed"`
	UpdateFailed     string `json:"update_failed"`
	DeleteFailed     string `json:"delete_failed"`
	InUse            string `json:"in_use"`
	Duplicate        string `json:"duplicate"`
	UncoveredPhase   string `json:"uncovered_phase"`
	IneligibleStaff  string `json:"ineligible_staff"`
	ValidationFailed string `json:"validation_failed"`
}

// DefaultLabels returns Labels with sensible English defaults (house
// zero-fallback convention). Tiers override every field via lyngua.
func DefaultLabels() Labels {
	return Labels{
		Page: PageLabels{
			Title:         "Classes",
			Subtitle:      "Manage the classes delivered on each section",
			ActiveTitle:   "Active Classes",
			InactiveTitle: "Inactive Classes",
		},
		Buttons: ButtonLabels{
			View:       "View",
			Add:        "Add Class",
			Edit:       "Edit Class",
			Delete:     "Delete Class",
			BulkDelete: "Delete Classes",
			Exclude:    "Exclude",
			Restore:    "Restore",
		},
		Columns: ColumnLabels{
			SubscriptionGroupId: "Section",
			ProductPlanId:       "Offering",
			JobTemplateId:       "Curriculum",
			Status:              "Status",
			AssignmentsCount:    "Teachers",
			DateCreated:         "Date Created",
			Actions:             "Actions",
		},
		Empty: EmptyLabels{
			Title:   "No Classes",
			Message: "No classes to display.",
		},
		Form: FormLabels{
			SectionIdentity:     "Class details",
			SubscriptionGroupId: "Section",
			SubscriptionGroupPH: "Select a section...",
			SubscriptionGroupSr: "Filter...",
			SubscriptionGroupIn: "The section this class is delivered on.",
			ProductPlanId:       "Offering",
			ProductPlanPH:       "Select an offering...",
			ProductPlanSr:       "Filter...",
			ProductPlanIn:       "The catalog offering this class delivers.",
			JobTemplateId:       "Curriculum",
			JobTemplatePH:       "Select a curriculum...",
			JobTemplateSr:       "Filter...",
			JobTemplateIn:       "The curriculum template naming this class's phases.",
			Status:              "Status",
			StatusInfo:          "Excluded classes leave the roster surface but keep their history.",
			StatusActive:        "Offered",
			StatusExcluded:      "Not offered",
			Active:              "Active",
			ActiveInfo:          "Inactive classes are hidden everywhere.",
		},
		Bulk: BulkLabels{
			DeleteTitle:   "Delete Classes",
			DeleteMessage: "Permanently delete the selected classes? This cannot be undone.",
		},
		Confirm: ConfirmLabels{
			DeleteTitle:         "Delete Class",
			DeleteMessage:       "Permanently delete this class? This cannot be undone.",
			ExcludeTitle:        "Exclude Class",
			ExcludeMessage:      "Exclude {{offering}} from {{section}}?",
			ExcludeWarning:      "This class has a live grade sheet — excluding hides it from the roster but does not delete grades.",
			RestoreTitle:        "Restore Class",
			RestoreMessage:      "Restore {{offering}} on {{section}}?",
			RemoveBlockedInUse:  "This class has active teachers or a live grade sheet and cannot be removed.",
			RemoveAssignTitle:   "Remove Assignment",
			RemoveAssignMessage: "Remove this teacher from {{offering}}?",
		},
		Tabs: TabLabels{
			Info:  "Info",
			Staff: "Staff",
		},
		Detail: DetailLabels{
			Title:            "Class",
			Header:           "{{offering}} — {{section}}",
			Section:          "Section",
			Offering:         "Offering",
			Curriculum:       "Curriculum",
			Period:           "Period",
			Plan:             "Plan",
			Status:           "Status",
			AssignmentsCount: "Teachers",
			DateCreated:      "Created",
			DateModified:     "Modified",
			NoGroup:          "No section",
			NoPlan:           "No offering",
			NoTemplate:       "No curriculum",
			GradeSheetLink:   "View grade sheet",
			VariantsPrefix:   "variants: ",
		},
		Picker: PickerLabels{
			Title:            "Add classes — {{section}}",
			SectionContext:   "Section",
			PlanContext:      "Plan",
			PeriodContext:    "Period",
			SelectAll:        "Select all ({{count}})",
			CurriculumPrefix: "curriculum: ",
			Ambiguous:        "Multiple curricula match — choose one",
			EmptyTitle:       "Nothing to add",
			EmptyMessage:     "Every offering on the plan is already on this section.",
			SubmitOne:        "Add 1 class",
			SubmitMany:       "Add {{count}} classes",
		},
		Assign: AssignLabels{
			Title:             "Assign staff — {{offering}}",
			SectionContext:    "Section",
			SubjectContext:    "Offering",
			CurriculumContext: "Curriculum",
			CurrentHeading:    "Current",
			EmptyCurrent:      "No staff assigned yet.",
			ModeAll:           "All semesters",
			ModePer:           "Per semester",
			AllPhasesChip:     "All semesters",
			Staff:             "Staff",
			StaffPH:           "Select staff...",
			StaffSearch:       "Filter...",
			Role:              "Role",
			RolePrimary:       "Teacher of record",
			RoleSecondary:     "Co-teacher",
			ClearAction:       "Clear",
			ReadOnlySlot:      "Manage on the class page →",
			Unauthorized:      "You are not authorized to perform this action",
		},
		Roster: RosterLabels{
			ColumnStaff:  "Staff",
			ColumnRole:   "Role",
			ColumnPhase:  "Semester",
			EmptyTitle:   "No Staff",
			EmptyMessage: "No staff assigned.",
			EmptyPool:    "No eligible staff yet",
			AddAction:    "Add staff",
			EditAction:   "Edit",
			RemoveAction: "Remove",
			Unauthorized: "You are not authorized to perform this action",
		},
		Errors: ErrorLabels{
			NotFound:         "Class not found",
			LoadFailed:       "Failed to load class",
			Unauthorized:     "You are not authorized to perform this action",
			CreateFailed:     "Failed to create class",
			UpdateFailed:     "Failed to update class",
			DeleteFailed:     "Failed to delete class",
			InUse:            "This class is in use and cannot be deleted.",
			Duplicate:        "This offering is already on this section.",
			UncoveredPhase:   "That semester does not belong to this class's curriculum.",
			IneligibleStaff:  "That staff member is not eligible for this offering.",
			ValidationFailed: "Could not save — check the highlighted fields.",
		},
	}
}
