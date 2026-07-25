package detail

import (
	"context"
	"strings"
	"testing"

	"github.com/erniealice/centymo-golang/domain/subscription/subscription_group"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	clientpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client"
	clientattributepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client_attribute"
	subscriptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription"
	subscriptiongroupmemberpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription_group_member"
)

func strp(s string) *string { return &s }

func testLabels() subscription_group.Labels {
	return subscription_group.Labels{
		Columns: subscription_group.ColumnLabels{Client: "Student", Subscription: "Enrollment", Status: "Status"},
		Empty:   subscription_group.EmptyLabels{Title: "No enrollments", Message: "None"},
	}
}

// rosterFixtureDeps wires stub closures over a fixed section: three members in a
// deliberately unsorted order (Zulu/male, Mango/female, Aardvark/male). Clients
// carry no Name column so the flat path renders "First Last" while the optioned
// path renders "Last, First".
func rosterFixtureDeps(opts subscription_group.Options) *DetailViewDeps {
	members := []*subscriptiongroupmemberpb.SubscriptionGroupMember{
		{Id: "mem1", SubscriptionGroupId: "sec1", ClientId: "cm2", SubscriptionId: "sub1", Active: true},
		{Id: "mem2", SubscriptionGroupId: "sec1", ClientId: "cf1", SubscriptionId: "sub2", Active: true},
		{Id: "mem3", SubscriptionGroupId: "sec1", ClientId: "cm1", SubscriptionId: "sub3", Active: true},
	}
	clients := []*clientpb.Client{
		{Id: "cm1", FirstName: strp("Al"), LastName: strp("Aardvark"), Active: true},
		{Id: "cm2", FirstName: strp("Bob"), LastName: strp("Zulu"), Active: true},
		{Id: "cf1", FirstName: strp("Ann"), LastName: strp("Mango"), Active: true},
	}
	subs := []*subscriptionpb.Subscription{
		{Id: "sub1", Code: strp("AAA111"), Active: true},
		{Id: "sub2", Code: strp("BBB222"), Active: true},
		{Id: "sub3", Code: strp("CCC333"), Active: true},
	}
	attrs := []*clientattributepb.ClientAttribute{
		{Id: "a1", ClientId: "cm1", AttributeId: "ag", Value: "male"},
		{Id: "a2", ClientId: "cm2", AttributeId: "ag", Value: "male"},
		{Id: "a3", ClientId: "cf1", AttributeId: "ag", Value: "female"},
	}
	return &DetailViewDeps{
		Labels:       testLabels(),
		CommonLabels: pyeza.CommonLabels{Status: pyeza.StatusLabels{Active: "Active", Inactive: "Inactive"}},
		TableLabels:  types.TableLabels{},
		Options:      opts,
		ListSubscriptionGroupMembers: func(ctx context.Context, req *subscriptiongroupmemberpb.ListSubscriptionGroupMembersRequest) (*subscriptiongroupmemberpb.ListSubscriptionGroupMembersResponse, error) {
			return &subscriptiongroupmemberpb.ListSubscriptionGroupMembersResponse{Data: members}, nil
		},
		ListClients: func(ctx context.Context, req *clientpb.ListClientsRequest) (*clientpb.ListClientsResponse, error) {
			return &clientpb.ListClientsResponse{Data: clients}, nil
		},
		ListSubscriptions: func(ctx context.Context, req *subscriptionpb.ListSubscriptionsRequest) (*subscriptionpb.ListSubscriptionsResponse, error) {
			return &subscriptionpb.ListSubscriptionsResponse{Data: subs}, nil
		},
		ListClientAttributes: func(ctx context.Context, req *clientattributepb.ListClientAttributesRequest) (*clientattributepb.ListClientAttributesResponse, error) {
			return &clientattributepb.ListClientAttributesResponse{Data: attrs}, nil
		},
		ResolveAttributeIDByCode: func(ctx context.Context, code string) (string, error) {
			if code == "gender" {
				return "ag", nil
			}
			return "", nil
		},
	}
}

// customRosterDeps wires stub closures over caller-supplied fixtures. ListClients
// returns exactly the passed slice (a workspace-scoped adapter would already have
// dropped out-of-workspace clients), while ListClientAttributes returns its slice
// verbatim (the client_attribute EAV table is NOT workspace-scoped) — the shape
// the foreign-client leak test needs.
func customRosterDeps(opts subscription_group.Options,
	members []*subscriptiongroupmemberpb.SubscriptionGroupMember,
	clients []*clientpb.Client,
	subs []*subscriptionpb.Subscription,
	attrs []*clientattributepb.ClientAttribute,
) *DetailViewDeps {
	return &DetailViewDeps{
		Labels:       testLabels(),
		CommonLabels: pyeza.CommonLabels{Status: pyeza.StatusLabels{Active: "Active", Inactive: "Inactive"}},
		TableLabels:  types.TableLabels{},
		Options:      opts,
		ListSubscriptionGroupMembers: func(ctx context.Context, req *subscriptiongroupmemberpb.ListSubscriptionGroupMembersRequest) (*subscriptiongroupmemberpb.ListSubscriptionGroupMembersResponse, error) {
			return &subscriptiongroupmemberpb.ListSubscriptionGroupMembersResponse{Data: members}, nil
		},
		ListClients: func(ctx context.Context, req *clientpb.ListClientsRequest) (*clientpb.ListClientsResponse, error) {
			return &clientpb.ListClientsResponse{Data: clients}, nil
		},
		ListSubscriptions: func(ctx context.Context, req *subscriptionpb.ListSubscriptionsRequest) (*subscriptionpb.ListSubscriptionsResponse, error) {
			return &subscriptionpb.ListSubscriptionsResponse{Data: subs}, nil
		},
		ListClientAttributes: func(ctx context.Context, req *clientattributepb.ListClientAttributesRequest) (*clientattributepb.ListClientAttributesResponse, error) {
			return &clientattributepb.ListClientAttributesResponse{Data: attrs}, nil
		},
		ResolveAttributeIDByCode: func(ctx context.Context, code string) (string, error) {
			if code == "gender" {
				return "ag", nil
			}
			return "", nil
		},
	}
}

func rosterCtx() context.Context {
	perms := types.NewUserPermissions([]string{"subscription_group_member:list"})
	return view.WithUserPermissions(context.Background(), perms)
}

// TestBuildSubscriptionsTable_Banded pins the optioned path: Male band leads
// (per GroupValueOrder, not value-alpha), rows sort last_name asc within a band,
// the client column renders "Last, First", and band ids are sg-band-<slug>.
func TestBuildSubscriptionsTable_Banded(t *testing.T) {
	deps := rosterFixtureDeps(subscription_group.Options{Roster: subscription_group.RowOptions{
		GroupByField:    "client_attributes.gender",
		GroupValueOrder: []string{"male", "female"},
		SortField:       "last_name",
	}})
	cfg := buildSubscriptionsTable(rosterCtx(), deps, "sec1", testLabels())

	if len(cfg.Rows) != 0 {
		t.Errorf("banded table should use Groups, got %d flat rows", len(cfg.Rows))
	}
	if len(cfg.Groups) != 2 {
		t.Fatalf("want 2 bands, got %d", len(cfg.Groups))
	}
	if cfg.Groups[0].Title != "male" || cfg.Groups[0].ID != "sg-band-male" {
		t.Errorf("band 0 = %q/%q, want male/sg-band-male (leads per GroupValueOrder)", cfg.Groups[0].Title, cfg.Groups[0].ID)
	}
	if cfg.Groups[0].DataAttrs["testid"] != "sg-band-male" {
		t.Errorf("band 0 testid = %q, want sg-band-male", cfg.Groups[0].DataAttrs["testid"])
	}
	if cfg.Groups[1].Title != "female" || cfg.Groups[1].ID != "sg-band-female" {
		t.Errorf("band 1 = %q/%q, want female/sg-band-female", cfg.Groups[1].Title, cfg.Groups[1].ID)
	}
	male := cfg.Groups[0].Rows
	if len(male) != 2 {
		t.Fatalf("male band: want 2 rows, got %d", len(male))
	}
	if male[0].Cells[0].Value != "Aardvark, Al" {
		t.Errorf("male[0] client = %q, want %q (last_name asc + Last, First)", male[0].Cells[0].Value, "Aardvark, Al")
	}
	if male[1].Cells[0].Value != "Zulu, Bob" {
		t.Errorf("male[1] client = %q, want %q", male[1].Cells[0].Value, "Zulu, Bob")
	}
	female := cfg.Groups[1].Rows
	if len(female) != 1 || female[0].Cells[0].Value != "Mango, Ann" {
		t.Errorf("female band = %+v, want single [Mango, Ann]", female)
	}
}

// TestBuildSubscriptionsTable_ZeroValue pins the fail-safe: no options → flat
// table, member order preserved, plain "First Last" display (no bands, no sort).
func TestBuildSubscriptionsTable_ZeroValue(t *testing.T) {
	deps := rosterFixtureDeps(subscription_group.Options{})
	cfg := buildSubscriptionsTable(rosterCtx(), deps, "sec1", testLabels())

	if len(cfg.Groups) != 0 {
		t.Errorf("zero-value table should be flat, got %d bands", len(cfg.Groups))
	}
	if len(cfg.Rows) != 3 {
		t.Fatalf("want 3 flat rows, got %d", len(cfg.Rows))
	}
	// Member order preserved: mem1(cm2 Bob Zulu), mem2(cf1 Ann Mango), mem3(cm1 Al Aardvark).
	want := []string{"Bob Zulu", "Ann Mango", "Al Aardvark"}
	for i, w := range want {
		if got := cfg.Rows[i].Cells[0].Value; got != w {
			t.Errorf("row[%d] client = %q, want %q (flat display, member order)", i, got, w)
		}
	}
	// Subscription + status cells (assert every cell, not just client).
	wantSub := []string{"AAA111", "BBB222", "CCC333"}
	for i, w := range wantSub {
		if got := cfg.Rows[i].Cells[1].Value; got != w {
			t.Errorf("row[%d] subscription = %q, want %q", i, got, w)
		}
		if got := cfg.Rows[i].Cells[2].Value; got != "Active" {
			t.Errorf("row[%d] status = %q, want Active", i, got)
		}
	}
	// Table display settings applied.
	if !cfg.ShowSearch || !cfg.ShowColumns || !cfg.ShowDensity || !cfg.ShowEntries {
		t.Errorf("table settings not applied: search=%v columns=%v density=%v entries=%v",
			cfg.ShowSearch, cfg.ShowColumns, cfg.ShowDensity, cfg.ShowEntries)
	}
}

// TestBuildSubscriptionsTable_ForeignClientNoLeak proves the security fix: a
// membership whose client_id is NOT returned by the workspace-scoped ListClients
// (a foreign/other-workspace client) never has its client_attribute value banded,
// even though the unscoped EAV stub holds that value.
func TestBuildSubscriptionsTable_ForeignClientNoLeak(t *testing.T) {
	members := []*subscriptiongroupmemberpb.SubscriptionGroupMember{
		{Id: "mem_local", SubscriptionGroupId: "sec1", ClientId: "cm1", SubscriptionId: "sub1", Active: true},
		{Id: "mem_foreign", SubscriptionGroupId: "sec1", ClientId: "cforeign", SubscriptionId: "sub1", Active: true},
	}
	clients := []*clientpb.Client{ // workspace-scoped: cforeign is absent
		{Id: "cm1", FirstName: strp("Al"), LastName: strp("Alpha"), Active: true},
	}
	subs := []*subscriptionpb.Subscription{{Id: "sub1", Code: strp("AAA111"), Active: true}}
	attrs := []*clientattributepb.ClientAttribute{ // unscoped EAV: holds the foreign value
		{Id: "a1", ClientId: "cm1", AttributeId: "ag", Value: "male"},
		{Id: "a2", ClientId: "cforeign", AttributeId: "ag", Value: "female"},
	}
	deps := customRosterDeps(subscription_group.Options{Roster: subscription_group.RowOptions{
		GroupByField:    "client_attributes.gender",
		GroupValueOrder: []string{"male", "female"},
		SortField:       "last_name",
	}}, members, clients, subs, attrs)

	cfg := buildSubscriptionsTable(rosterCtx(), deps, "sec1", testLabels())

	var maleRows, noValRows []types.TableRow
	for _, g := range cfg.Groups {
		if strings.EqualFold(g.Title, "female") {
			t.Fatalf("foreign client's gender leaked as a band: %q", g.Title)
		}
		switch g.Title {
		case "male":
			maleRows = g.Rows
		case "—":
			noValRows = g.Rows
		}
	}
	if len(maleRows) != 1 || maleRows[0].Cells[0].Value != "Alpha, Al" {
		t.Errorf("male band = %+v, want single [Alpha, Al] (local client only)", maleRows)
	}
	if len(noValRows) != 1 || noValRows[0].Cells[0].Value != "cforeign" {
		t.Errorf("no-value band = %+v, want single foreign member rendered as raw id", noValRows)
	}
}

// TestBuildSubscriptionsTable_BandedFirstNameTie pins the last→first contract:
// two members share a last name; the display Name is set to the INVERSE order to
// prove first name (not display name) breaks the tie.
func TestBuildSubscriptionsTable_BandedFirstNameTie(t *testing.T) {
	members := []*subscriptiongroupmemberpb.SubscriptionGroupMember{
		{Id: "mem_zoe", SubscriptionGroupId: "sec1", ClientId: "s_zoe", SubscriptionId: "sub1", Active: true},
		{Id: "mem_aaron", SubscriptionGroupId: "sec1", ClientId: "s_aaron", SubscriptionId: "sub1", Active: true},
	}
	clients := []*clientpb.Client{
		{Id: "s_aaron", Name: strp("Zulu Display"), FirstName: strp("Aaron"), LastName: strp("Smith"), Active: true},
		{Id: "s_zoe", Name: strp("Alpha Display"), FirstName: strp("Zoe"), LastName: strp("Smith"), Active: true},
	}
	subs := []*subscriptionpb.Subscription{{Id: "sub1", Code: strp("AAA111"), Active: true}}
	attrs := []*clientattributepb.ClientAttribute{
		{Id: "a1", ClientId: "s_aaron", AttributeId: "ag", Value: "male"},
		{Id: "a2", ClientId: "s_zoe", AttributeId: "ag", Value: "male"},
	}
	deps := customRosterDeps(subscription_group.Options{Roster: subscription_group.RowOptions{
		GroupByField:    "client_attributes.gender",
		GroupValueOrder: []string{"male"},
		SortField:       "last_name",
	}}, members, clients, subs, attrs)

	cfg := buildSubscriptionsTable(rosterCtx(), deps, "sec1", testLabels())
	if len(cfg.Groups) != 1 {
		t.Fatalf("want 1 band, got %d (%+v)", len(cfg.Groups), cfg.Groups)
	}
	rows := cfg.Groups[0].Rows
	if len(rows) != 2 {
		t.Fatalf("want 2 rows, got %d", len(rows))
	}
	if rows[0].Cells[0].Value != "Smith, Aaron" || rows[1].Cells[0].Value != "Smith, Zoe" {
		t.Errorf("rows = [%q, %q], want [Smith, Aaron ; Smith, Zoe] (first-name tie-break, not display)", rows[0].Cells[0].Value, rows[1].Cells[0].Value)
	}
}

// TestBuildSubscriptionsTable_MultiMembershipSameClient: one client with two
// memberships → two distinct rows, both carrying the shared client's band value.
func TestBuildSubscriptionsTable_MultiMembershipSameClient(t *testing.T) {
	members := []*subscriptiongroupmemberpb.SubscriptionGroupMember{
		{Id: "memA", SubscriptionGroupId: "sec1", ClientId: "cm1", SubscriptionId: "sub1", Active: true},
		{Id: "memB", SubscriptionGroupId: "sec1", ClientId: "cm1", SubscriptionId: "sub2", Active: true},
	}
	clients := []*clientpb.Client{{Id: "cm1", FirstName: strp("Al"), LastName: strp("Alpha"), Active: true}}
	subs := []*subscriptionpb.Subscription{
		{Id: "sub1", Code: strp("AAA111"), Active: true},
		{Id: "sub2", Code: strp("BBB222"), Active: true},
	}
	attrs := []*clientattributepb.ClientAttribute{{Id: "a1", ClientId: "cm1", AttributeId: "ag", Value: "male"}}
	deps := customRosterDeps(subscription_group.Options{Roster: subscription_group.RowOptions{
		GroupByField:    "client_attributes.gender",
		GroupValueOrder: []string{"male"},
		SortField:       "last_name",
	}}, members, clients, subs, attrs)

	cfg := buildSubscriptionsTable(rosterCtx(), deps, "sec1", testLabels())
	if len(cfg.Groups) != 1 || cfg.Groups[0].Title != "male" {
		t.Fatalf("want single male band, got %+v", cfg.Groups)
	}
	rows := cfg.Groups[0].Rows
	if len(rows) != 2 {
		t.Fatalf("want 2 rows (one per membership), got %d", len(rows))
	}
	for _, r := range rows {
		if r.Cells[0].Value != "Alpha, Al" {
			t.Errorf("row client = %q, want Alpha, Al (both memberships map to the shared client)", r.Cells[0].Value)
		}
	}
	if rows[0].ID == rows[1].ID {
		t.Errorf("rows share id %q, want distinct member ids", rows[0].ID)
	}
}
