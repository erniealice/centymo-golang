// Package revenue is the revenue-domain consumer facade (centymo restructure).
//
// PURE RE-EXPORT — zero behaviour. The revenue domain's data/route types,
// Default* constructors, and URL consts moved into per-entity packages under
// domain/revenue/<entity>/ with entity-local names (the <Entity> prefix stripped).
// This facade re-adds the original prefixed names so existing consumers
// (block/, service-admin) keep resolving revenue.<Entity>Labels /
// revenue.Default<Entity>Routes() / revenue.<Entity>ListURL unchanged.
//
// An entity package MUST NEVER import this facade (that would be an import
// cycle revenue -> <entity> -> revenue); cross-entity references go DIRECT to the
// sibling package.
package revenue

import (
	revenuepkg "github.com/erniealice/centymo-golang/domain/revenue/revenue"
	revenuerunpkg "github.com/erniealice/centymo-golang/domain/revenue/revenue_run"
)

// Re-exported data/route types (type aliases — identity-preserving).
type (
	RevenueActionLabels            = revenuepkg.ActionLabels
	RevenueBulkLabels              = revenuepkg.BulkLabels
	RevenueButtonLabels            = revenuepkg.ButtonLabels
	RevenueColumnLabels            = revenuepkg.ColumnLabels
	RevenueConfirmLabels           = revenuepkg.ConfirmLabels
	RevenueDashboardLabels         = revenuepkg.DashboardLabels
	RevenueDetailLabels            = revenuepkg.DetailLabels
	RevenueEmptyLabels             = revenuepkg.EmptyLabels
	RevenueErrorLabels             = revenuepkg.ErrorLabels
	RevenueFormLabels              = revenuepkg.FormLabels
	RevenueLabels                  = revenuepkg.Labels
	RevenuePageLabels              = revenuepkg.PageLabels
	RevenueRoutes                  = revenuepkg.Routes
	RevenueRunActionLabels         = revenuerunpkg.ActionLabels
	RevenueRunDetailLabels         = revenuerunpkg.DetailLabels
	RevenueRunDetailTabLabels      = revenuerunpkg.DetailTabLabels
	RevenueRunErrorLabels          = revenuerunpkg.ErrorLabels
	RevenueRunInvoicesTabLabels    = revenuerunpkg.InvoicesTabLabels
	RevenueRunLabels               = revenuerunpkg.Labels
	RevenueRunListColumnLabels     = revenuerunpkg.ListColumnLabels
	RevenueRunListEmptyLabels      = revenuerunpkg.ListEmptyLabels
	RevenueRunListEmptyStateLabels = revenuerunpkg.ListEmptyStateLabels
	RevenueRunListFilterLabels     = revenuerunpkg.ListFilterLabels
	RevenueRunListLabels           = revenuerunpkg.ListLabels
	RevenueRunOutcomeLabels        = revenuerunpkg.OutcomeLabels
	RevenueRunQueueBulkLabels      = revenuerunpkg.QueueBulkLabels
	RevenueRunQueueColumnLabels    = revenuerunpkg.QueueColumnLabels
	RevenueRunQueueEmptyLabels     = revenuerunpkg.QueueEmptyLabels
	RevenueRunQueueLabels          = revenuerunpkg.QueueLabels
	RevenueRunResultsTabLabels     = revenuerunpkg.ResultsTabLabels
	RevenueRunRoutes               = revenuerunpkg.Routes
	RevenueRunScopeKindLabels      = revenuerunpkg.ScopeKindLabels
	RevenueRunSelectionsTabLabels  = revenuerunpkg.SelectionsTabLabels
	RevenueRunStatusBadgeLabels    = revenuerunpkg.StatusBadgeLabels
	RevenueRunSummaryLabels        = revenuerunpkg.SummaryLabels
	RevenueSettingsLabels          = revenuepkg.SettingsLabels
)

// Re-exported URL route consts (const-identity preserved).
const (
	RevenueAddURL                     = revenuepkg.AddURL
	RevenueAttachmentDeleteURL        = revenuepkg.AttachmentDeleteURL
	RevenueAttachmentUploadURL        = revenuepkg.AttachmentUploadURL
	RevenueBulkDeleteURL              = revenuepkg.BulkDeleteURL
	RevenueBulkSetStatusURL           = revenuepkg.BulkSetStatusURL
	RevenueDashboardURL               = revenuepkg.DashboardURL
	RevenueDeleteURL                  = revenuepkg.DeleteURL
	RevenueDetailURL                  = revenuepkg.DetailURL
	RevenueEditURL                    = revenuepkg.EditURL
	RevenueEmailURL                   = revenuepkg.EmailURL
	RevenueInvoiceDownloadURL         = revenuepkg.InvoiceDownloadURL
	RevenueLineItemAddURL             = revenuepkg.LineItemAddURL
	RevenueLineItemDiscountURL        = revenuepkg.LineItemDiscountURL
	RevenueLineItemEditURL            = revenuepkg.LineItemEditURL
	RevenueLineItemRemoveURL          = revenuepkg.LineItemRemoveURL
	RevenueLineItemTableURL           = revenuepkg.LineItemTableURL
	RevenueListURL                    = revenuepkg.ListURL
	RevenuePaymentAddURL              = revenuepkg.PaymentAddURL
	RevenuePaymentEditURL             = revenuepkg.PaymentEditURL
	RevenuePaymentRemoveURL           = revenuepkg.PaymentRemoveURL
	RevenuePaymentTableURL            = revenuepkg.PaymentTableURL
	RevenuePriceLookupURL             = revenuepkg.PriceLookupURL
	RevenueRecomputeTaxesURL          = revenuepkg.RecomputeTaxesURL
	RevenueRunAttachmentDeleteURL     = revenuerunpkg.AttachmentDeleteURL
	RevenueRunAttachmentUploadURL     = revenuerunpkg.AttachmentUploadURL
	RevenueRunDetailTabActionURL      = revenuerunpkg.DetailTabActionURL
	RevenueRunDetailURL               = revenuerunpkg.DetailURL
	RevenueRunListTableURL            = revenuerunpkg.ListTableURL
	RevenueRunListURL                 = revenuerunpkg.ListURL
	RevenueRunQueueTableURL           = revenuerunpkg.QueueTableURL
	RevenueRunQueueURL                = revenuerunpkg.QueueURL
	RevenueRunSubmitBatchURL          = revenuerunpkg.SubmitBatchURL
	RevenueSearchClientURL            = revenuepkg.SearchClientURL
	RevenueSearchLocationURL          = revenuepkg.SearchLocationURL
	RevenueSearchProductURL           = revenuepkg.SearchProductURL
	RevenueSearchSubscriptionURL      = revenuepkg.SearchSubscriptionURL
	RevenueSetStatusURL               = revenuepkg.SetStatusURL
	RevenueSettingsTemplateDefaultURL = revenuepkg.SettingsTemplateDefaultURL
	RevenueSettingsTemplateDeleteURL  = revenuepkg.SettingsTemplateDeleteURL
	RevenueSettingsTemplateUploadURL  = revenuepkg.SettingsTemplateUploadURL
	RevenueSettingsTemplatesURL       = revenuepkg.SettingsTemplatesURL
	RevenueSummaryURL                 = revenuepkg.SummaryURL
	RevenueTabActionURL               = revenuepkg.TabActionURL
	RevenueTableURL                   = revenuepkg.TableURL
)

// Re-exported Default* constructors (function values).
var (
	DefaultRevenueRoutes    = revenuepkg.DefaultRoutes
	DefaultRevenueRunLabels = revenuerunpkg.DefaultLabels
	DefaultRevenueRunRoutes = revenuerunpkg.DefaultRoutes
)
