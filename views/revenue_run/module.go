// Package revenuerun wires the Revenue Run (invoice-run) view module:
// Surface B — workspace queue page (Phase 7).
// Surface D — run history list + detail pages (Phase 4).
package revenuerun

import (
	"context"

	centymo "github.com/erniealice/centymo-golang"
	rrshared "github.com/erniealice/centymo-golang/views/revenue_run/shared"
	revenuerundetail "github.com/erniealice/centymo-golang/views/revenue_run/detail"
	revenuerunlist "github.com/erniealice/centymo-golang/views/revenue_run/list"
	revenuerunqueue "github.com/erniealice/centymo-golang/views/revenue_run/queue"
	revenuerunqueueaction "github.com/erniealice/centymo-golang/views/revenue_run/queue/action"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
)

// ---------------------------------------------------------------------------
// Re-export shared view-typed data shapes so block.go callers can reference
// them via the top-level revenuerun package (e.g. revenuerunmod.RevenueRunRow).
// ---------------------------------------------------------------------------

// RevenueRunRow is the view-layer representation of a single revenue run.
type RevenueRunRow = rrshared.RevenueRunRow

// RevenueRunWithAttempts bundles a run and its attempt list for the detail page.
type RevenueRunWithAttempts = rrshared.RevenueRunWithAttempts

// RevenueRunAttemptRow is the view-layer representation of a single run attempt.
type RevenueRunAttemptRow = rrshared.RevenueRunAttemptRow

// RevenueRow is a minimal invoice row for the Invoices tab.
type RevenueRow = rrshared.RevenueRow

// ListRevenueRunsScope carries filter parameters for the list page.
type ListRevenueRunsScope = rrshared.ListRevenueRunsScope

// ---------------------------------------------------------------------------
// Re-export queue-local types so block.go can reference them without
// importing the queue sub-package directly.
// ---------------------------------------------------------------------------

// QueueClientRecord is a minimal client row for the queue fan-out.
type QueueClientRecord = revenuerunqueue.ClientRecord

// QueueCandidateInput is the per-period input shape for the queue fan-out.
type QueueCandidateInput = revenuerunqueue.RevenueRunCandidateInput

// BatchRunInput is the per-client input for GenerateRevenueRun in the batch handler.
type BatchRunInput = revenuerunqueueaction.GenerateRevenueRunInput

// BatchRunOutput is the per-client output from GenerateRevenueRun.
type BatchRunOutput = revenuerunqueueaction.GenerateRevenueRunOutput

// ---------------------------------------------------------------------------
// ModuleDeps — typed callbacks; no espyna/proto types cross this boundary.
// ---------------------------------------------------------------------------

// ModuleDeps holds all dependencies for the revenue-run view module.
type ModuleDeps struct {
	Routes       centymo.RevenueRunRoutes
	Labels       centymo.RevenueRunLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// Surface D callbacks — run history list + detail.

	// ListRevenueRuns returns a page of run rows matching the given scope.
	ListRevenueRuns func(ctx context.Context, scope ListRevenueRunsScope) ([]RevenueRunRow, string, error)

	// ReadRevenueRun fetches a single run plus all its attempts by run ID.
	ReadRevenueRun func(ctx context.Context, id string) (*RevenueRunWithAttempts, error)

	// ListRevenueByRunID fetches revenue records whose run_id matches the given ID.
	// Used to populate the Invoices tab on the detail page.
	ListRevenueByRunID func(ctx context.Context, runID string) ([]RevenueRow, error)

	// Surface B callbacks — workspace queue page.

	// ClientDetailURLTemplate is the path template for the client detail page
	// (e.g. "/app/clients/detail/{id}"). Optional — rows are not linked when empty.
	ClientDetailURLTemplate string

	// ClientDrawerURLTemplate is the path template for the Surface-A per-client
	// revenue-run drawer (e.g. "/action/client/revenue-run/{id}").
	// Populated by WithClientRevenueRunDrawerURL BlockOption.
	ClientDrawerURLTemplate string

	// ListClients returns all clients visible to the current workspace user.
	ListClients func(ctx context.Context, cursor string) ([]QueueClientRecord, string, error)

	// ListRevenueRunCandidates returns pending billing periods for one client.
	// Called per-client in a bounded fan-out goroutine on the queue page.
	ListRevenueRunCandidates func(ctx context.Context, clientID, asOfDate string) ([]QueueCandidateInput, error)

	// GenerateRevenueRun executes the revenue run for a single client.
	// Called by the batch-run POST handler.
	GenerateRevenueRun func(ctx context.Context, in BatchRunInput) (*BatchRunOutput, error)

	// Attachment operations.
	UploadFile       func(ctx context.Context, bucket, key string, content []byte, contentType string) error
	ListAttachments  func(ctx context.Context, moduleKey, foreignKey string) (*attachmentpb.ListAttachmentsResponse, error)
	CreateAttachment func(ctx context.Context, req *attachmentpb.CreateAttachmentRequest) (*attachmentpb.CreateAttachmentResponse, error)
	DeleteAttachment func(ctx context.Context, req *attachmentpb.DeleteAttachmentRequest) (*attachmentpb.DeleteAttachmentResponse, error)
	NewAttachmentID  func() string
}

// ---------------------------------------------------------------------------
// Module — holds constructed view instances.
// ---------------------------------------------------------------------------

// Module holds all constructed revenue-run views.
type Module struct {
	routes    centymo.RevenueRunRoutes
	// Surface D.
	List      view.View
	Table     view.View
	Detail    view.View
	TabAction view.View
	// Surface B.
	Queue     view.View
	QueueTable view.View
	BatchRun  view.View
	// Attachments.
	AttachmentUpload view.View
	AttachmentDelete view.View
}

// NewModule constructs the revenue-run module from the given deps.
func NewModule(deps *ModuleDeps) *Module {
	listDeps := &revenuerunlist.ListViewDeps{
		Routes:          deps.Routes,
		Labels:          deps.Labels,
		CommonLabels:    deps.CommonLabels,
		TableLabels:     deps.TableLabels,
		ListRevenueRuns: deps.ListRevenueRuns,
	}
	detailDeps := &revenuerundetail.DetailViewDeps{
		Routes:             deps.Routes,
		Labels:             deps.Labels,
		CommonLabels:       deps.CommonLabels,
		TableLabels:        deps.TableLabels,
		ReadRevenueRun:     deps.ReadRevenueRun,
		ListRevenueByRunID: deps.ListRevenueByRunID,
	}
	detailDeps.UploadFile = deps.UploadFile
	detailDeps.ListAttachments = deps.ListAttachments
	detailDeps.CreateAttachment = deps.CreateAttachment
	detailDeps.DeleteAttachment = deps.DeleteAttachment
	detailDeps.NewAttachmentID = deps.NewAttachmentID
	queueDeps := &revenuerunqueue.QueueViewDeps{
		Routes:                   deps.Routes,
		Labels:                   deps.Labels,
		CommonLabels:             deps.CommonLabels,
		TableLabels:              deps.TableLabels,
		ClientDetailURLTemplate:  deps.ClientDetailURLTemplate,
		ClientDrawerURLTemplate:  deps.ClientDrawerURLTemplate,
		ListClients:              deps.ListClients,
		ListRevenueRunCandidates: deps.ListRevenueRunCandidates,
	}
	batchRunDeps := &revenuerunqueueaction.BatchRunDeps{
		Routes:             deps.Routes,
		Labels:             deps.Labels,
		GenerateRevenueRun: deps.GenerateRevenueRun,
	}
	m := &Module{
		routes:     deps.Routes,
		List:       revenuerunlist.NewView(listDeps),
		Table:      revenuerunlist.NewTableView(listDeps),
		Detail:     revenuerundetail.NewView(detailDeps),
		TabAction:  revenuerundetail.NewTabAction(detailDeps),
		Queue:      revenuerunqueue.NewView(queueDeps),
		QueueTable: revenuerunqueue.NewTableView(queueDeps),
		BatchRun:   revenuerunqueueaction.NewBatchRunAction(batchRunDeps),
	}
	if deps.UploadFile != nil {
		m.AttachmentUpload = revenuerundetail.NewAttachmentUploadAction(detailDeps)
		m.AttachmentDelete = revenuerundetail.NewAttachmentDeleteAction(detailDeps)
	}
	return m
}

// RegisterRoutes registers all revenue-run routes on the given registrar.
func (m *Module) RegisterRoutes(r view.RouteRegistrar) {
	// Surface D — run history list + detail.
	r.GET(m.routes.ListURL, m.List)
	r.GET(m.routes.ListTableURL, m.Table)
	r.POST(m.routes.ListTableURL, m.Table)
	r.GET(m.routes.DetailURL, m.Detail)
	r.GET(m.routes.DetailTabActionURL, m.TabAction)
	// Surface B — workspace queue page.
	r.GET(m.routes.QueueURL, m.Queue)
	r.GET(m.routes.QueueTableURL, m.QueueTable)
	r.POST(m.routes.QueueTableURL, m.QueueTable)
	r.POST(m.routes.SubmitBatchURL, m.BatchRun)
	if m.AttachmentUpload != nil {
		r.GET(m.routes.AttachmentUploadURL, m.AttachmentUpload)
		r.POST(m.routes.AttachmentUploadURL, m.AttachmentUpload)
		r.POST(m.routes.AttachmentDeleteURL, m.AttachmentDelete)
	}
}
