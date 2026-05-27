package centymo

// routes_treasury_method.go — Stage-1 (Wave 4) route configs for the treasury
// domain-rebuild Method management views: collection_method (selling side) +
// disbursement_method (buying side).
//
// These views live under the NEW `/treasury_collection/methods/*` and
// `/treasury_disbursement/methods/*` URL prefixes (D-2.10 2-app split / D-4.22
// routing SSOT). The service-admin container injects the concrete URLs from
// apps/.../composition/routes_treasury_collection.go +
// routes_treasury_disbursement.go; the Default*() helpers below mirror those
// canonical paths so the package is usable standalone (tests, other consumers).
//
// Stage 1 scope = template CRUD + the §A drawer-form template-slot pattern.
// Eligibility / Grants / Instances / Approvals tabs are later stages.

// ---------------------------------------------------------------------------
// CollectionMethodRoutes (selling side, pages.md §B-5)
// ---------------------------------------------------------------------------

// CollectionMethodRoutes holds all route paths for the collection_method views.
type CollectionMethodRoutes struct {
	ActiveNav    string `json:"active_nav"`
	ActiveSubNav string `json:"active_sub_nav"`

	ListURL   string `json:"list_url"`
	DetailURL string `json:"detail_url"`
	AddURL    string `json:"add_url"`
	EditURL   string `json:"edit_url"`
	DeleteURL string `json:"delete_url"`
	// FragmentURL serves the kind-specific drawer fragment (§A template-slot).
	// Driven by the category select's hx-get; ?cat=<CATEGORY> selects which
	// fragment HTML to swap into #kind-specific-slot.
	FragmentURL string `json:"fragment_url"`
	// TabActionURL serves HTMX detail-tab swaps (info / versions / activity in
	// Stage 1; later tabs added in Stages 2/3/4/6).
	TabActionURL string `json:"tab_action_url"`

	// Lifecycle workflow actions (header buttons; pages.md §B-5 tab 1).
	// Wired nil-safe until the espyna Publish/Close/Archive/Revise use cases land.
	PublishURL string `json:"publish_url"`
	CloseURL   string `json:"close_url"`
	ArchiveURL string `json:"archive_url"`
	ReviseURL  string `json:"revise_url"`
}

// DefaultCollectionMethodRoutes returns the canonical `/treasury_collection/methods/*`
// route set per pages.md §B-5 + routes_treasury_collection.go SSOT.
func DefaultCollectionMethodRoutes() CollectionMethodRoutes {
	return CollectionMethodRoutes{
		ActiveNav:    "treasury_collection",
		ActiveSubNav: "active",
		ListURL:      "/treasury_collection/methods/{status}",
		DetailURL:    "/treasury_collection/methods/detail/{id}",
		AddURL:       "/action/collection-method/add",
		EditURL:      "/action/collection-method/edit/{id}",
		DeleteURL:    "/action/collection-method/delete",
		FragmentURL:  "/action/collection-method/fragment",
		TabActionURL: "/action/collection-method/detail/{id}/tab/{tab}",
		PublishURL:   "/action/collection-method/publish/{id}",
		CloseURL:     "/action/collection-method/close/{id}",
		ArchiveURL:   "/action/collection-method/archive/{id}",
		ReviseURL:    "/action/collection-method/revise/{id}",
	}
}

// RouteMap returns dot-notation keys → paths for route resolution.
func (r CollectionMethodRoutes) RouteMap() map[string]string {
	return map[string]string{
		"collection_method.list":     r.ListURL,
		"collection_method.detail":   r.DetailURL,
		"collection_method.add":      r.AddURL,
		"collection_method.edit":     r.EditURL,
		"collection_method.delete":   r.DeleteURL,
		"collection_method.fragment": r.FragmentURL,
		"collection_method.tab":      r.TabActionURL,
		"collection_method.publish":  r.PublishURL,
		"collection_method.close":    r.CloseURL,
		"collection_method.archive":  r.ArchiveURL,
		"collection_method.revise":   r.ReviseURL,
	}
}

// ---------------------------------------------------------------------------
// DisbursementMethodRoutes (buying side, pages.md §C-5)
// ---------------------------------------------------------------------------

// DisbursementMethodRoutes holds all route paths for the disbursement_method views.
// Buying-side asymmetry (D-4.9): no audience_mode / eligibility / grants.
type DisbursementMethodRoutes struct {
	ActiveNav    string `json:"active_nav"`
	ActiveSubNav string `json:"active_sub_nav"`

	ListURL      string `json:"list_url"`
	DetailURL    string `json:"detail_url"`
	AddURL       string `json:"add_url"`
	EditURL      string `json:"edit_url"`
	DeleteURL    string `json:"delete_url"`
	FragmentURL  string `json:"fragment_url"`
	TabActionURL string `json:"tab_action_url"`

	PublishURL string `json:"publish_url"`
	CloseURL   string `json:"close_url"`
	ArchiveURL string `json:"archive_url"`
	ReviseURL  string `json:"revise_url"`
}

// DefaultDisbursementMethodRoutes returns the canonical `/treasury_disbursement/methods/*`
// route set per pages.md §C-5 + routes_treasury_disbursement.go SSOT.
func DefaultDisbursementMethodRoutes() DisbursementMethodRoutes {
	return DisbursementMethodRoutes{
		ActiveNav:    "treasury_disbursement",
		ActiveSubNav: "active",
		ListURL:      "/treasury_disbursement/methods/{status}",
		DetailURL:    "/treasury_disbursement/methods/detail/{id}",
		AddURL:       "/action/disbursement-method/add",
		EditURL:      "/action/disbursement-method/edit/{id}",
		DeleteURL:    "/action/disbursement-method/delete",
		FragmentURL:  "/action/disbursement-method/fragment",
		TabActionURL: "/action/disbursement-method/detail/{id}/tab/{tab}",
		PublishURL:   "/action/disbursement-method/publish/{id}",
		CloseURL:     "/action/disbursement-method/close/{id}",
		ArchiveURL:   "/action/disbursement-method/archive/{id}",
		ReviseURL:    "/action/disbursement-method/revise/{id}",
	}
}

// RouteMap returns dot-notation keys → paths for route resolution.
func (r DisbursementMethodRoutes) RouteMap() map[string]string {
	return map[string]string{
		"disbursement_method.list":     r.ListURL,
		"disbursement_method.detail":   r.DetailURL,
		"disbursement_method.add":      r.AddURL,
		"disbursement_method.edit":     r.EditURL,
		"disbursement_method.delete":   r.DeleteURL,
		"disbursement_method.fragment": r.FragmentURL,
		"disbursement_method.tab":      r.TabActionURL,
		"disbursement_method.publish":  r.PublishURL,
		"disbursement_method.close":    r.CloseURL,
		"disbursement_method.archive":  r.ArchiveURL,
		"disbursement_method.revise":   r.ReviseURL,
	}
}
