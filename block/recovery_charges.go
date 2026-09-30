package block

import (
	"context"
	"fmt"

	"github.com/erniealice/espyna-golang/consumer/compose"
	clientpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client"
	subscriptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription"

	revenuedomain "github.com/erniealice/centymo-golang/domain/revenue"
	recoverydocumentpkg "github.com/erniealice/centymo-golang/domain/revenue/recovery_document"
	revenuepkg "github.com/erniealice/centymo-golang/domain/revenue/revenue"
	subscriptiondom "github.com/erniealice/centymo-golang/domain/subscription"
	billablechargepkg "github.com/erniealice/centymo-golang/domain/subscription/billable_charge"
	treasurydomain "github.com/erniealice/centymo-golang/domain/treasury"
	collectionpkg "github.com/erniealice/centymo-golang/domain/treasury/collection"
	collectionapplicationpkg "github.com/erniealice/centymo-golang/domain/treasury/collection_application"
)

// Usage & pass-through S1 — billable charges, recovery documents, receive and
// apply. The three units are opt-in (WithRecoveryCharges); an app that never
// passes the option mounts none of them and its route map is unchanged.

// WithRecoveryCharges mounts the billable charge list, the recovery document
// list/detail and the receive-and-apply drawer (which also becomes the
// collection list's primary action) when enabled is true. Never implied by
// another option; false is a no-op, so an app can pass its capability flag
// verbatim (fycha WithChargePolicies precedent).
func WithRecoveryCharges(enabled bool) EngineOption {
	return func(c *engineConfig) { c.recoveryCharges = enabled }
}

func requireBillableCharge(uc *UseCases) error {
	if uc == nil {
		return fmt.Errorf("centymo billable_charge: UseCases not supplied")
	}
	return requireUnit("billable_charge", map[string]bool{
		"UseCases.BillableCharge.ListBillableCharges":      uc.BillableCharge.ListBillableCharges != nil,
		"UseCases.BillableCharge.AdjustBillableCharge":     uc.BillableCharge.AdjustBillableCharge != nil,
		"UseCases.RecoveryDocument.IssueRecoveryDocuments": uc.RecoveryDocument.IssueRecoveryDocuments != nil,
		"UseCases.DocumentSeries.ListDocumentSeries":       uc.DocumentSeries.ListDocumentSeries != nil,
	})
}

func requireRecoveryDocument(uc *UseCases) error {
	if uc == nil {
		return fmt.Errorf("centymo recovery_document: UseCases not supplied")
	}
	return requireUnit("recovery_document", map[string]bool{
		"UseCases.RecoveryDocument.ListRecoveryDocuments":           uc.RecoveryDocument.ListRecoveryDocuments != nil,
		"UseCases.RecoveryDocument.ReadRecoveryDocument":            uc.RecoveryDocument.ReadRecoveryDocument != nil,
		"UseCases.RecoveryDocument.VoidRecoveryDocument":            uc.RecoveryDocument.VoidRecoveryDocument != nil,
		"UseCases.CollectionApplication.ListCollectionApplications": uc.CollectionApplication.ListCollectionApplications != nil,
	})
}

func requireReceiveApply(uc *UseCases) error {
	if uc == nil {
		return fmt.Errorf("centymo collection_application: UseCases not supplied")
	}
	return requireUnit("collection_application", map[string]bool{
		"UseCases.CollectionApplication.ReceiveAndApplyCollection":    uc.CollectionApplication.ReceiveAndApplyCollection != nil,
		"UseCases.CollectionApplication.PreviewCollectionApplication": uc.CollectionApplication.PreviewCollectionApplication != nil,
		"UseCases.CollectionApplication.ReverseCollectionApplication": uc.CollectionApplication.ReverseCollectionApplication != nil,
		// R5 m10: without the client lister the picker renders empty on a 200 page.
		"UseCases.Entity.Client.ListClients": uc.Entity.Client.ListClients != nil,
	})
}

// clientLister / subscriptionLister are the optional display-name resolvers.
func clientLister(uc *UseCases) func(context.Context, *clientpb.ListClientsRequest) (*clientpb.ListClientsResponse, error) {
	return uc.Entity.Client.ListClients
}

func subscriptionLister(uc *UseCases) func(context.Context, *subscriptionpb.ListSubscriptionsRequest) (*subscriptionpb.ListSubscriptionsResponse, error) {
	return uc.Subscription.ListSubscriptions
}

// RecoveryChargeUnits returns the three opt-in units (billable charge,
// recovery document, receive and apply).
func RecoveryChargeUnits(uc *UseCases, infra *Infra) []compose.Unit {
	return []compose.Unit{BillableChargeUnit(uc, infra), RecoveryDocumentUnit(uc, infra), CollectionApplicationUnit(uc, infra)}
}

// BillableChargeUnit mounts Revenue › Billable charges (open · issued) with the
// issue-recovery-documents and adjust drawers.
func BillableChargeUnit(uc *UseCases, _ *Infra) compose.Unit {
	u := billablechargepkg.Describe()
	u.Mount = func(mc *compose.MountContext) error {
		if err := requireBillableCharge(uc); err != nil {
			return err
		}
		r := u.Routes.(*billablechargepkg.Routes)
		l := u.Labels.(*billablechargepkg.Labels)
		recoveryLabels := recoverydocumentpkg.DefaultLabels()
		if rl, ok := compose.LabelsOf[*recoverydocumentpkg.Labels](mc, "revenue.recovery_document"); ok {
			recoveryLabels = *rl
		}
		subscriptiondom.NewBillableChargeModule(&subscriptiondom.BillableChargeModuleDeps{
			Routes: *r, Labels: *l, RecoveryLabels: recoveryLabels,
			CommonLabels: mc.Common, TableLabels: mc.Table,
			UseCases: &billablechargepkg.UseCases{
				ListBillableCharges:    uc.BillableCharge.ListBillableCharges,
				AdjustBillableCharge:   uc.BillableCharge.AdjustBillableCharge,
				IssueRecoveryDocuments: uc.RecoveryDocument.IssueRecoveryDocuments,
				ListDocumentSeries:     uc.DocumentSeries.ListDocumentSeries,
				ListClients:            clientLister(uc),
				ListSubscriptions:      subscriptionLister(uc),
			},
		}).RegisterRoutes(mc.Routes)
		return nil
	}
	return u
}

// RecoveryDocumentUnit mounts Revenue › Recovery documents (issued · void),
// the detail page and the void drawer.
func RecoveryDocumentUnit(uc *UseCases, _ *Infra) compose.Unit {
	u := recoverydocumentpkg.Describe()
	u.Mount = func(mc *compose.MountContext) error {
		if err := requireRecoveryDocument(uc); err != nil {
			return err
		}
		r := u.Routes.(*recoverydocumentpkg.Routes)
		l := u.Labels.(*recoverydocumentpkg.Labels)
		deps := &revenuedomain.RecoveryDocumentModuleDeps{
			Routes: *r, Labels: *l, CommonLabels: mc.Common, TableLabels: mc.Table,
			ApplicationLabels: collectionapplicationpkg.DefaultLabels(),
			UseCases: &recoverydocumentpkg.UseCases{
				ListRecoveryDocuments:      uc.RecoveryDocument.ListRecoveryDocuments,
				ReadRecoveryDocument:       uc.RecoveryDocument.ReadRecoveryDocument,
				VoidRecoveryDocument:       uc.RecoveryDocument.VoidRecoveryDocument,
				ListCollectionApplications: uc.CollectionApplication.ListCollectionApplications,
				ListClients:                clientLister(uc),
				ReadClient:                 uc.Entity.Client.ReadClient,
			},
		}
		if al, ok := compose.LabelsOf[*collectionapplicationpkg.Labels](mc, "treasury.collection_application"); ok {
			deps.ApplicationLabels = *al
		}
		if ar, ok := compose.RoutesOf[*collectionapplicationpkg.Routes](mc, "treasury.collection_application"); ok {
			deps.ReverseURL = ar.ReverseURL
			deps.ReceiveApplyURL = ar.ReceiveApplyURL
		}
		deps.ApplicationsTableID = collectionapplicationpkg.ApplicationsTableID
		revenuedomain.NewRecoveryDocumentModule(deps).RegisterRoutes(mc.Routes)
		return nil
	}
	return u
}

// CollectionApplicationUnit mounts the receive-and-apply drawer, its application
// preview and the application reversal.
func CollectionApplicationUnit(uc *UseCases, _ *Infra) compose.Unit {
	u := collectionapplicationpkg.Describe()
	u.Mount = func(mc *compose.MountContext) error {
		if err := requireReceiveApply(uc); err != nil {
			return err
		}
		r := u.Routes.(*collectionapplicationpkg.Routes)
		l := u.Labels.(*collectionapplicationpkg.Labels)
		deps := &treasurydomain.CollectionApplicationModuleDeps{
			Routes: *r, Labels: *l, CommonLabels: mc.Common,
			UseCases: &collectionapplicationpkg.UseCases{
				ReceiveAndApplyCollection:    uc.CollectionApplication.ReceiveAndApplyCollection,
				PreviewCollectionApplication: uc.CollectionApplication.PreviewCollectionApplication,
				ReverseCollectionApplication: uc.CollectionApplication.ReverseCollectionApplication,
				ListClients:                  clientLister(uc),
				ReadClient:                   uc.Entity.Client.ReadClient,
				ListCollectionMethods:        uc.CollectionMethod.ListCollectionMethods,
			},
		}
		if cr, ok := compose.RoutesOf[*collectionpkg.Routes](mc, "treasury.collection"); ok {
			deps.CollectionDetailURL = cr.DetailURL
		}
		// The drawer's client picker is the async client search revenue already
		// mounts (C11); absent (revenue not composed) it falls back to a paged select.
		if rr, ok := compose.RoutesOf[*revenuepkg.Routes](mc, "revenue.revenue"); ok {
			deps.SearchClientURL = rr.SearchClientURL
		}
		treasurydomain.NewCollectionApplicationModule(deps).RegisterRoutes(mc.Routes)
		return nil
	}
	return u
}
