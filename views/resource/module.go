package resource

import (
	"context"

	centymo "github.com/erniealice/centymo-golang"
	resourceaction "github.com/erniealice/centymo-golang/views/resource/action"
	resourcelist "github.com/erniealice/centymo-golang/views/resource/list"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	view "github.com/erniealice/pyeza-golang/view"

	resourcepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/resource"
)

// ModuleDeps holds all dependencies for the resource module.
type ModuleDeps struct {
	Routes       centymo.ResourceRoutes
	Labels       centymo.ResourceLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	ListResources  func(ctx context.Context, req *resourcepb.ListResourcesRequest) (*resourcepb.ListResourcesResponse, error)
	ReadResource   func(ctx context.Context, req *resourcepb.ReadResourceRequest) (*resourcepb.ReadResourceResponse, error)
	CreateResource func(ctx context.Context, req *resourcepb.CreateResourceRequest) (*resourcepb.CreateResourceResponse, error)
	UpdateResource func(ctx context.Context, req *resourcepb.UpdateResourceRequest) (*resourcepb.UpdateResourceResponse, error)
	DeleteResource func(ctx context.Context, req *resourcepb.DeleteResourceRequest) (*resourcepb.DeleteResourceResponse, error)
}

// Module holds all constructed resource views.
type Module struct {
	routes        centymo.ResourceRoutes
	List          view.View
	Table         view.View
	Add           view.View
	Edit          view.View
	Delete        view.View
	BulkDelete    view.View
	SetStatus     view.View
	BulkSetStatus view.View
}

// NewModule creates the resource module with all views wired.
func NewModule(deps *ModuleDeps) *Module {
	actionDeps := &resourceaction.Deps{
		Routes:         deps.Routes,
		Labels:         deps.Labels,
		CreateResource: deps.CreateResource,
		ReadResource:   deps.ReadResource,
		UpdateResource: deps.UpdateResource,
		DeleteResource: deps.DeleteResource,
	}

	listDeps := &resourcelist.ListViewDeps{
		Routes:        deps.Routes,
		ListResources: deps.ListResources,
		Labels:        deps.Labels,
		CommonLabels:  deps.CommonLabels,
		TableLabels:   deps.TableLabels,
	}
	listView := resourcelist.NewView(listDeps)
	tableView := resourcelist.NewTableView(listDeps)

	return &Module{
		routes:        deps.Routes,
		List:          listView,
		Table:         tableView,
		Add:           resourceaction.NewAddAction(actionDeps),
		Edit:          resourceaction.NewEditAction(actionDeps),
		Delete:        resourceaction.NewDeleteAction(actionDeps),
		BulkDelete:    resourceaction.NewBulkDeleteAction(actionDeps),
		SetStatus:     resourceaction.NewSetStatusAction(actionDeps),
		BulkSetStatus: resourceaction.NewBulkSetStatusAction(actionDeps),
	}
}

// RegisterRoutes registers all resource routes.
func (m *Module) RegisterRoutes(r view.RouteRegistrar) {
	r.GET(m.routes.ListURL, m.List)
	r.GET(m.routes.TableURL, m.Table)
	r.GET(m.routes.AddURL, m.Add)
	r.POST(m.routes.AddURL, m.Add)
	r.GET(m.routes.EditURL, m.Edit)
	r.POST(m.routes.EditURL, m.Edit)
	r.POST(m.routes.DeleteURL, m.Delete)
	r.POST(m.routes.BulkDeleteURL, m.BulkDelete)
	r.POST(m.routes.SetStatusURL, m.SetStatus)
	r.POST(m.routes.BulkSetStatusURL, m.BulkSetStatus)
}
