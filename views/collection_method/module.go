package collection_method

import (
	"context"

	centymo "github.com/erniealice/centymo-golang"
	cmaction "github.com/erniealice/centymo-golang/views/collection_method/action"
	cmdetail "github.com/erniealice/centymo-golang/views/collection_method/detail"
	cmlist "github.com/erniealice/centymo-golang/views/collection_method/list"

	cmpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/collection_method"
	eligrulepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/collection_method_eligibility_rule"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// ModuleDeps holds all dependencies for the collection_method module
// (treasury-domain-rebuild Stage 1 + Stage 2, pages.md §B-5).
//
// All CRUD closures are OPTIONAL (nil-safe): the module mounts regardless,
// rendering empty/degraded states until the espyna use cases are wired.
type ModuleDeps struct {
	Routes       centymo.CollectionMethodRoutes
	Labels       centymo.CollectionMethodLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// Stage 1 — collection_method CRUD
	ListCollectionMethods  func(ctx context.Context, req *cmpb.ListCollectionMethodsRequest) (*cmpb.ListCollectionMethodsResponse, error)
	ReadCollectionMethod   func(ctx context.Context, req *cmpb.ReadCollectionMethodRequest) (*cmpb.ReadCollectionMethodResponse, error)
	CreateCollectionMethod func(ctx context.Context, req *cmpb.CreateCollectionMethodRequest) (*cmpb.CreateCollectionMethodResponse, error)
	UpdateCollectionMethod func(ctx context.Context, req *cmpb.UpdateCollectionMethodRequest) (*cmpb.UpdateCollectionMethodResponse, error)
	DeleteCollectionMethod func(ctx context.Context, req *cmpb.DeleteCollectionMethodRequest) (*cmpb.DeleteCollectionMethodResponse, error)

	// Stage 2 — collection_method_eligibility_rule CRUD (pages.md §B-5 tab 2)
	ListCollectionMethodEligibilityRules  func(ctx context.Context, req *eligrulepb.ListCollectionMethodEligibilityRulesRequest) (*eligrulepb.ListCollectionMethodEligibilityRulesResponse, error)
	ReadCollectionMethodEligibilityRule   func(ctx context.Context, req *eligrulepb.ReadCollectionMethodEligibilityRuleRequest) (*eligrulepb.ReadCollectionMethodEligibilityRuleResponse, error)
	CreateCollectionMethodEligibilityRule func(ctx context.Context, req *eligrulepb.CreateCollectionMethodEligibilityRuleRequest) (*eligrulepb.CreateCollectionMethodEligibilityRuleResponse, error)
	UpdateCollectionMethodEligibilityRule func(ctx context.Context, req *eligrulepb.UpdateCollectionMethodEligibilityRuleRequest) (*eligrulepb.UpdateCollectionMethodEligibilityRuleResponse, error)
	DeleteCollectionMethodEligibilityRule func(ctx context.Context, req *eligrulepb.DeleteCollectionMethodEligibilityRuleRequest) (*eligrulepb.DeleteCollectionMethodEligibilityRuleResponse, error)
}

// Module holds all constructed collection_method views.
type Module struct {
	routes    centymo.CollectionMethodRoutes
	List      view.View
	Detail    view.View
	TabAction view.View
	Add       view.View
	Edit      view.View
	Fragment  view.View
	Delete    view.View
	// Stage 2 — Eligibility Rules tab CRUD
	EligibilityRuleAdd    view.View
	EligibilityRuleEdit   view.View
	EligibilityRuleDelete view.View
}

// NewModule creates the collection_method module with all views wired.
func NewModule(deps *ModuleDeps) *Module {
	actionDeps := &cmaction.Deps{
		Routes:                 deps.Routes,
		Labels:                 deps.Labels,
		CommonLabels:           deps.CommonLabels,
		CreateCollectionMethod: deps.CreateCollectionMethod,
		ReadCollectionMethod:   deps.ReadCollectionMethod,
		UpdateCollectionMethod: deps.UpdateCollectionMethod,
		DeleteCollectionMethod: deps.DeleteCollectionMethod,
	}
	// Stage 2 — eligibility rule tab deps (nil closures degrade gracefully).
	eligRuleDeps := &cmdetail.EligibilityRuleTabDeps{
		Routes:                                deps.Routes,
		Labels:                                deps.Labels,
		CommonLabels:                          deps.CommonLabels,
		TableLabels:                           deps.TableLabels,
		ListCollectionMethodEligibilityRules:  deps.ListCollectionMethodEligibilityRules,
		ReadCollectionMethodEligibilityRule:   deps.ReadCollectionMethodEligibilityRule,
		CreateCollectionMethodEligibilityRule: deps.CreateCollectionMethodEligibilityRule,
		UpdateCollectionMethodEligibilityRule: deps.UpdateCollectionMethodEligibilityRule,
		DeleteCollectionMethodEligibilityRule: deps.DeleteCollectionMethodEligibilityRule,
	}
	detailDeps := &cmdetail.DetailViewDeps{
		Routes:               deps.Routes,
		Labels:               deps.Labels,
		CommonLabels:         deps.CommonLabels,
		TableLabels:          deps.TableLabels,
		ReadCollectionMethod: deps.ReadCollectionMethod,
		EligibilityRuleDeps:  eligRuleDeps,
	}
	listDeps := &cmlist.ListViewDeps{
		Routes:                deps.Routes,
		ListCollectionMethods: deps.ListCollectionMethods,
		Labels:                deps.Labels,
		CommonLabels:          deps.CommonLabels,
		TableLabels:           deps.TableLabels,
	}

	return &Module{
		routes:                deps.Routes,
		List:                  cmlist.NewView(listDeps),
		Detail:                cmdetail.NewView(detailDeps),
		TabAction:             cmdetail.NewTabAction(detailDeps),
		Add:                   cmaction.NewAddAction(actionDeps),
		Edit:                  cmaction.NewEditAction(actionDeps),
		Fragment:              cmaction.NewFragmentAction(actionDeps),
		Delete:                cmaction.NewDeleteAction(actionDeps),
		EligibilityRuleAdd:    cmdetail.NewEligibilityRuleAddAction(eligRuleDeps),
		EligibilityRuleEdit:   cmdetail.NewEligibilityRuleEditAction(eligRuleDeps),
		EligibilityRuleDelete: cmdetail.NewEligibilityRuleDeleteAction(eligRuleDeps),
	}
}

// RegisterRoutes registers all collection_method routes.
func (m *Module) RegisterRoutes(r view.RouteRegistrar) {
	r.GET(m.routes.ListURL, m.List)
	r.GET(m.routes.DetailURL, m.Detail)
	r.GET(m.routes.TabActionURL, m.TabAction)

	r.GET(m.routes.AddURL, m.Add)
	r.POST(m.routes.AddURL, m.Add)
	r.GET(m.routes.EditURL, m.Edit)
	r.POST(m.routes.EditURL, m.Edit)
	r.GET(m.routes.FragmentURL, m.Fragment)
	r.POST(m.routes.DeleteURL, m.Delete)

	// Stage 2 — Eligibility Rules tab CRUD routes (pages.md §B-5 tab 2).
	if m.routes.EligibilityRuleAddURL != "" {
		r.GET(m.routes.EligibilityRuleAddURL, m.EligibilityRuleAdd)
		r.POST(m.routes.EligibilityRuleAddURL, m.EligibilityRuleAdd)
	}
	if m.routes.EligibilityRuleEditURL != "" {
		r.GET(m.routes.EligibilityRuleEditURL, m.EligibilityRuleEdit)
		r.POST(m.routes.EligibilityRuleEditURL, m.EligibilityRuleEdit)
	}
	if m.routes.EligibilityRuleDeleteURL != "" {
		r.POST(m.routes.EligibilityRuleDeleteURL, m.EligibilityRuleDelete)
	}
}
