// Package treasury is the treasury-domain consumer facade (centymo restructure).
//
// PURE RE-EXPORT — zero behaviour. The treasury domain's data/route types,
// Default* constructors, and URL consts moved into per-entity packages under
// domain/treasury/<entity>/ with entity-local names (the <Entity> prefix stripped).
// This facade re-adds the original prefixed names so existing consumers
// (block/, service-admin) keep resolving treasury.<Entity>Labels /
// treasury.Default<Entity>Routes() / treasury.<Entity>ListURL unchanged.
//
// An entity package MUST NEVER import this facade (that would be an import
// cycle treasury -> <entity> -> treasury); cross-entity references go DIRECT to the
// sibling package.
package treasury

import (
	advancesdashboardpkg "github.com/erniealice/centymo-golang/domain/treasury/advancesdashboard"
	collectionpkg "github.com/erniealice/centymo-golang/domain/treasury/collection"
	disbursementpkg "github.com/erniealice/centymo-golang/domain/treasury/disbursement"
	sharedpkg "github.com/erniealice/centymo-golang/domain/treasury/shared"
)

// Re-exported shared advance contract (lives in domain/treasury/shared — used by
// both collection + disbursement and cross-domain by subscription/expenditure).
type (
	AdvanceCancelViewInput          = sharedpkg.AdvanceCancelViewInput
	AdvanceCancelViewOutput         = sharedpkg.AdvanceCancelViewOutput
	AdvanceEnumLabels               = sharedpkg.AdvanceEnumLabels
	AdvanceKindLabels               = sharedpkg.AdvanceKindLabels
	AdvanceKindRootLabels           = sharedpkg.AdvanceKindRootLabels
	AdvanceProrationPolicyLabels    = sharedpkg.AdvanceProrationPolicyLabels
	AdvanceRecognizeMilestoneInput  = sharedpkg.AdvanceRecognizeMilestoneInput
	AdvanceRecognizeMilestoneOutput = sharedpkg.AdvanceRecognizeMilestoneOutput
	AdvanceRefundViewInput          = sharedpkg.AdvanceRefundViewInput
	AdvanceRefundViewOutput         = sharedpkg.AdvanceRefundViewOutput
	AdvanceSettleViewInput          = sharedpkg.AdvanceSettleViewInput
	AdvanceSettleViewOutput         = sharedpkg.AdvanceSettleViewOutput
	AdvanceStatusLabels             = sharedpkg.AdvanceStatusLabels
	TreasuryAdvanceActionLabels     = sharedpkg.TreasuryAdvanceActionLabels
	TreasuryAdvanceLabels           = sharedpkg.TreasuryAdvanceLabels
)

// Re-exported shared advance Default* constructors.
var (
	DefaultAdvanceEnumLabels                 = sharedpkg.DefaultAdvanceEnumLabels
	DefaultAdvanceKindLabels                 = sharedpkg.DefaultAdvanceKindLabels
	DefaultAdvanceKindRootLabels             = sharedpkg.DefaultAdvanceKindRootLabels
	DefaultAdvanceProrationPolicyLabels      = sharedpkg.DefaultAdvanceProrationPolicyLabels
	DefaultAdvanceStatusLabels               = sharedpkg.DefaultAdvanceStatusLabels
	DefaultTreasuryCollectionAdvanceLabels   = sharedpkg.DefaultTreasuryCollectionAdvanceLabels
	DefaultTreasuryDisbursementAdvanceLabels = sharedpkg.DefaultTreasuryDisbursementAdvanceLabels
)

// Re-exported data/route types (type aliases — identity-preserving).
type (
	AdvancesDashboardLabels        = advancesdashboardpkg.Labels
	AdvancesDashboardSectionLabels = advancesdashboardpkg.SectionLabels
	AdvancesDashboardTableLabels   = advancesdashboardpkg.TableLabels
	CashDashboardLabels            = collectionpkg.CashDashboardLabels
	CollectionActionLabels         = collectionpkg.ActionLabels
	CollectionBulkLabels           = collectionpkg.BulkLabels
	CollectionButtonLabels         = collectionpkg.ButtonLabels
	CollectionColumnLabels         = collectionpkg.ColumnLabels
	CollectionConfirmLabels        = collectionpkg.ConfirmLabels
	CollectionDetailLabels         = collectionpkg.DetailLabels
	CollectionEmptyLabels          = collectionpkg.EmptyLabels
	CollectionErrorLabels          = collectionpkg.ErrorLabels
	CollectionFormLabels           = collectionpkg.FormLabels
	CollectionLabels               = collectionpkg.Labels
	CollectionPageLabels           = collectionpkg.PageLabels
	CollectionRoutes               = collectionpkg.Routes
	CollectionStatusLabels         = collectionpkg.StatusLabels
	DisbursementActionLabels       = disbursementpkg.ActionLabels
	DisbursementBulkLabels         = disbursementpkg.BulkLabels
	DisbursementButtonLabels       = disbursementpkg.ButtonLabels
	DisbursementColumnLabels       = disbursementpkg.ColumnLabels
	DisbursementConfirmLabels      = disbursementpkg.ConfirmLabels
	DisbursementDetailLabels       = disbursementpkg.DetailLabels
	DisbursementEmptyLabels        = disbursementpkg.EmptyLabels
	DisbursementErrorLabels        = disbursementpkg.ErrorLabels
	DisbursementFormLabels         = disbursementpkg.FormLabels
	DisbursementLabels             = disbursementpkg.Labels
	DisbursementPageLabels         = disbursementpkg.PageLabels
	DisbursementRoutes             = disbursementpkg.Routes
	DisbursementStatusLabels       = disbursementpkg.StatusLabels
	TreasuryAdvancesRoutes         = advancesdashboardpkg.Routes
)

// Re-exported URL route consts (const-identity preserved).
const (
	AdvanceCollectionListURL                  = advancesdashboardpkg.AdvanceCollectionListURL
	AdvanceDisbursementListURL                = advancesdashboardpkg.AdvanceDisbursementListURL
	AdvancesDashboardURL                      = advancesdashboardpkg.AdvancesDashboardURL
	CollectionAddURL                          = collectionpkg.AddURL
	CollectionAttachmentDeleteURL             = collectionpkg.AttachmentDeleteURL
	CollectionAttachmentUploadURL             = collectionpkg.AttachmentUploadURL
	CollectionBulkDeleteURL                   = collectionpkg.BulkDeleteURL
	CollectionBulkSetStatusURL                = collectionpkg.BulkSetStatusURL
	CollectionDashboardURL                    = collectionpkg.DashboardURL
	CollectionDeleteURL                       = collectionpkg.DeleteURL
	CollectionDetailURL                       = collectionpkg.DetailURL
	CollectionEditURL                         = collectionpkg.EditURL
	CollectionListURL                         = collectionpkg.ListURL
	CollectionSetStatusURL                    = collectionpkg.SetStatusURL
	CollectionTabActionURL                    = collectionpkg.TabActionURL
	DisbursementAddURL                        = disbursementpkg.AddURL
	DisbursementAttachmentDeleteURL           = disbursementpkg.AttachmentDeleteURL
	DisbursementAttachmentUploadURL           = disbursementpkg.AttachmentUploadURL
	DisbursementBulkDeleteURL                 = disbursementpkg.BulkDeleteURL
	DisbursementBulkSetStatusURL              = disbursementpkg.BulkSetStatusURL
	DisbursementDashboardURL                  = disbursementpkg.DashboardURL
	DisbursementDeleteURL                     = disbursementpkg.DeleteURL
	DisbursementDetailURL                     = disbursementpkg.DetailURL
	DisbursementEditURL                       = disbursementpkg.EditURL
	DisbursementListURL                       = disbursementpkg.ListURL
	DisbursementSetStatusURL                  = disbursementpkg.SetStatusURL
	DisbursementTabActionURL                  = disbursementpkg.TabActionURL
	TreasuryCollectionAdvanceScheduleTabURL   = collectionpkg.TreasuryCollectionAdvanceScheduleTabURL
	TreasuryCollectionCancelURL               = collectionpkg.TreasuryCollectionCancelURL
	TreasuryCollectionRefundURL               = collectionpkg.TreasuryCollectionRefundURL
	TreasuryCollectionSettleURL               = collectionpkg.TreasuryCollectionSettleURL
	TreasuryDisbursementAdvanceScheduleTabURL = disbursementpkg.TreasuryDisbursementAdvanceScheduleTabURL
	TreasuryDisbursementCancelURL             = disbursementpkg.TreasuryDisbursementCancelURL
	TreasuryDisbursementRefundURL             = disbursementpkg.TreasuryDisbursementRefundURL
	TreasuryDisbursementSettleURL             = disbursementpkg.TreasuryDisbursementSettleURL
)

// Re-exported Default* constructors (function values).
var (
	DefaultAdvancesDashboardLabels = advancesdashboardpkg.DefaultLabels
	DefaultCollectionLabels        = collectionpkg.DefaultLabels
	DefaultCollectionRoutes        = collectionpkg.DefaultRoutes
	DefaultDisbursementLabels      = disbursementpkg.DefaultLabels
	DefaultDisbursementRoutes      = disbursementpkg.DefaultRoutes
	DefaultTreasuryAdvancesRoutes  = advancesdashboardpkg.DefaultRoutes
)
