// Package block — typed wiring contract for centymo.Block.
//
// This file is the answer to "what does centymo need from outside?"
// Service-admin's composition layer constructs a *UseCases value from
// espyna's consumer container; centymo's Block() consumes only this
// typed shape.
//
// IMPORTANT: shape this struct by what CENTYMO needs, NOT by mirroring
// espyna's *consumer.UseCases. Service-admin's adapter is the only
// place that knows both vocabularies. If espyna restructures its
// container, only that adapter changes.
package block

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	commonv1pb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	clientpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client"
	locationpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/location"
	paymenttermpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/payment_term"
	supplierpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/supplier"
	workspacepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/workspace"
	accruedexpensepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/accrued_expense"
	expenditurepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/expenditure"
	expenditurecategorypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/expenditure_category"
	expenditurelineitempb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/expenditure_line_item"
	expenserecognitionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/expense_recognition"
	expenserecognitionlinepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/expense_recognition_line"
	expenserecognitionrunpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/expense_recognition_run"
	procurementrequestpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/procurement_request"
	procurementrequestlinepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/procurement_request_line"
	purchaseorderpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/purchase_order"
	supplierbillingeventpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/supplier_billing_event"
	suppliercontractpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/supplier_contract"
	suppliercontractlinepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/supplier_contract_line"
	suppliercontractpriceschedulepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/supplier_contract_price_schedule"
	suppliercontractpriceschedulelinepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/supplier_contract_price_schedule_line"
	inventorydepreciationpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/inventory_depreciation"
	inventoryitempb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/inventory_item"
	inventoryserialpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/inventory_serial"
	inventorytransactionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/inventory_transaction"
	inventoryserialhistorypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/serial_history"
	jobpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/operation/job"
	jobactivitypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/operation/job_activity"
	jobphasepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/operation/job_phase"
	jobtemplatepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/operation/job_template"
	jobtemplatedepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/operation/job_template_phase"
	jobtemplaterelationpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/operation/job_template_relation"
	jobttemplatetaskpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/operation/job_template_task"
	costplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/procurement/cost_plan"
	costschedulepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/procurement/cost_schedule"
	supplierplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/procurement/supplier_plan"
	supplierproductcostplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/procurement/supplier_product_cost_plan"
	supplierproductplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/procurement/supplier_product_plan"
	suppliersubscriptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/procurement/supplier_subscription"
	linepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/line"
	lineworkspaceuserpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/line_workspace_user"
	plangrouppb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/plan_group"
	plangroupplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/plan_group_plan"
	pricelistpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/price_list"
	priceproductpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/price_product"
	productpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product"
	productattributepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_attribute"
	productlinepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_line"
	productoptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_option"
	productoptionvaluepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_option_value"
	productplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_plan"
	productplanstaffpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_plan_staff"
	productvariantpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_variant"
	productvariantimagepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_variant_image"
	productvariantoptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_variant_option"
	resourcepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/resource"
	revenuepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/revenue"
	revenuelineitempb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/revenue_line_item"
	revenuepaymentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/revenue_payment"
	revenuerunpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/revenue_run"
	revenuetaxlinepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/revenue_tax_line"
	billingeventpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/billing_event"
	planpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/plan"
	priceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_plan"
	priceschedulepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_schedule"
	pricescheduleworkspaceuserpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_schedule_workspace_user"
	productpriceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/product_price_plan"
	subscriptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription"
	subscriptiongrouppb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription_group"
	subscriptiongroupmemberpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription_group_member"
	subscriptiongroupproductplanstaffpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription_group_product_plan_staff"
	subscriptiongroupworkspaceuserpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription_group_workspace_user"
	collectionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/collection"
	collectionmethodpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/collection_method"
	disbursementpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/disbursement"

	expenseboard "github.com/erniealice/centymo-golang/domain/expenditure/expenditure/expense_dashboard"
	purchaseboard "github.com/erniealice/centymo-golang/domain/expenditure/expenditure/purchase_dashboard"
	productdashboard "github.com/erniealice/centymo-golang/domain/product/product/dashboard"
	treasurydomain "github.com/erniealice/centymo-golang/domain/treasury"
	collectiondashboard "github.com/erniealice/centymo-golang/domain/treasury/collection/dashboard"
)

// UseCases declares everything centymo's Block() needs from outside.
// Construction is service-admin's job; centymo only declares the shape.
//
// Naming conventions (locked-in 2026-05-10, matching plan.md §2):
//
//  1. Field names are SINGULAR matching the proto folder name.
//  2. Group struct types use the `<Entity>UseCases` suffix.
//  3. Closure signatures use proto request/response types — no block-local
//     transport types. Exception: dashboard closures use centymo view-layer
//     Request/Response types (espyna internals are unreachable from centymo).
//  4. Grouped by what CENTYMO does, NOT by espyna's container shape.
type UseCases struct {
	// ExtractUserID extracts the authenticated user ID from a request context.
	// Wired by service-admin's adapter to identity.Must(ctx).UserID.
	// Used by workflow action closures that need to record the acting user
	// (e.g. ApprovedBy, ActivatedBy) without importing the espyna consumer package.
	ExtractUserID func(context.Context) string

	// SetActive sets ONLY the `active` boolean column on the named collection. It
	// is a deliberate, auditable, capability-narrow replacement for the deleted
	// DataSource duck's generic Update — needed because proto3 omits false bools,
	// so the typed proto Update cannot clear `active`. The collection argument is
	// the table name (e.g. "subscription", "plan", "inventory_item",
	// "cost_schedule"). Bound by service-admin in W7; nil-safe — when unbound the
	// per-entity activate/deactivate toggles log + no-op (the row's active flag is
	// left unchanged rather than panicking).
	// 20260612-datasource-typed-path W6.
	SetActive func(ctx context.Context, collection string, id string, active bool) error
	// Domain CRUD + use-case groups (singular field, `XxxUseCases` type)
	// Ordered alphabetically for easy scanning.
	Collection        CollectionUseCases
	CollectionMethod  CollectionMethodUseCases
	Common            CommonUseCases
	Disbursement      DisbursementUseCases
	Entity            EntityUseCases
	Expenditure       ExpenditureUseCases
	Inventory         InventoryUseCases
	LineWorkspaceUser LineWorkspaceUserUseCases
	Operation         OperationUseCases
	Plan              PlanUseCases
	PlanGroup         PlanGroupUseCases
	PlanGroupPlan     PlanGroupPlanUseCases
	PricePlan         PricePlanUseCases
	PriceSchedule     PriceScheduleUseCases

	PriceScheduleWorkspaceUser PriceScheduleWorkspaceUserUseCases

	PriceList         PriceListUseCases
	Procurement       ProcurementUseCases
	Product           ProductUseCases
	ProductPlanStaff  ProductPlanStaffUseCases
	Revenue           RevenueUseCases
	RevenueRun        RevenueRunUseCases
	Subscription      SubscriptionUseCases
	SubscriptionGroup SubscriptionGroupUseCases

	SubscriptionGroupMember           SubscriptionGroupMemberUseCases
	SubscriptionGroupProductPlanStaff SubscriptionGroupProductPlanStaffUseCases
	SubscriptionGroupWorkspaceUser    SubscriptionGroupWorkspaceUserUseCases

	SupplierContract SupplierContractUseCases

	// 20260517-advance-cash-events Plan B Phase 3 — workspace advances dashboard.
	TreasuryAdvances TreasuryAdvancesUseCases

	// 20260517-expense-run Plan A Phase 4 — buying-side Expense Recognition Run.
	ExpenseRecognitionRun ExpenseRecognitionRunUseCases
}

// setActiveClosure adapts the capability-narrow UseCases.SetActive into the
// per-entity `func(ctx, id, active) error` shape the view-layer action Deps
// expect (SetSubscriptionActive, SetPlanActive, SetItemActive, …), binding the
// collection name. Nil-safe: when SetActive is unbound (service-admin wires it
// in W7) the returned closure logs and no-ops rather than panicking — the
// activate/deactivate toggle silently leaves the row unchanged, matching the
// fail-soft contract every other unwired typed closure in this package uses.
// 20260612-datasource-typed-path W6 — replaces the deleted DataSource duck's
// `Update(collection, id, {"active": active})` call at each toggle site.
func setActiveClosure(useCases *UseCases, collection string) func(context.Context, string, bool) error {
	return func(ctx context.Context, id string, active bool) error {
		if useCases == nil || useCases.SetActive == nil {
			log.Printf("centymo.Block: SetActive(%q) is not wired — active toggle no-op for id %s (bind UseCases.SetActive in service-admin)", collection, id)
			return nil
		}
		return useCases.SetActive(ctx, collection, id, active)
	}
}

// ExpenseRecognitionRunUseCases groups everything the centymo
// expense_recognition_run view module needs from outside.
//
// Mirror of RevenueRunUseCases. List/Read/ListAttempts are repo-direct
// pass-through calls; ListExpenseRunCandidates + GenerateExpenseRun are the
// application-layer use cases.
// Plan A 20260517-expense-run Phase 4.
type ExpenseRecognitionRunUseCases struct {
	ListExpenseRecognitionRuns        func(context.Context, *expenserecognitionrunpb.ListExpenseRecognitionRunsRequest) (*expenserecognitionrunpb.ListExpenseRecognitionRunsResponse, error)
	ReadExpenseRecognitionRun         func(context.Context, *expenserecognitionrunpb.ReadExpenseRecognitionRunRequest) (*expenserecognitionrunpb.ReadExpenseRecognitionRunResponse, error)
	ListExpenseRecognitionRunAttempts func(context.Context, *expenserecognitionrunpb.ListExpenseRecognitionRunAttemptsRequest) (*expenserecognitionrunpb.ListExpenseRecognitionRunAttemptsResponse, error)
	ListExpenseRunCandidates          func(context.Context, *expenserecognitionrunpb.ListExpenseRunCandidatesRequest) (*expenserecognitionrunpb.ListExpenseRunCandidatesResponse, error)
	GenerateExpenseRun                func(context.Context, *expenserecognitionrunpb.GenerateExpenseRunRequest) (*expenserecognitionrunpb.GenerateExpenseRunResponse, error)
}

// -- Common ------------------------------------------------------------------

type CommonUseCases struct {
	ListCategories func(context.Context, *commonv1pb.ListCategoriesRequest) (*commonv1pb.ListCategoriesResponse, error)
	ListAttributes func(context.Context, *commonv1pb.ListAttributesRequest) (*commonv1pb.ListAttributesResponse, error)
	ReadAttribute  func(context.Context, *commonv1pb.ReadAttributeRequest) (*commonv1pb.ReadAttributeResponse, error)
}

// -- Entity ------------------------------------------------------------------

type EntityUseCases struct {
	Client      ClientUseCases
	Location    LocationUseCases
	PaymentTerm PaymentTermUseCases
	Supplier    SupplierUseCases
	Workspace   WorkspaceUseCases
}

// PaymentTermUseCases groups the typed payment_term reads the revenue
// (client/both scope) + expenditure (supplier/both scope) payment-terms
// dropdowns need. Replaces the duck-typed DataSource.ListSimple("payment_term")
// path. The block scopes/filters the returned rows itself (by entity_scope), so
// a single unfiltered ListPaymentTerms closure serves both callers.
// 20260612-datasource-typed-path W6. Nil-safe — the dropdowns render empty
// (no payment-term options) when unwired.
type PaymentTermUseCases struct {
	ListPaymentTerms func(context.Context, *paymenttermpb.ListPaymentTermsRequest) (*paymenttermpb.ListPaymentTermsResponse, error)
}

type ClientUseCases struct {
	ListClients         func(context.Context, *clientpb.ListClientsRequest) (*clientpb.ListClientsResponse, error)
	ReadClient          func(context.Context, *clientpb.ReadClientRequest) (*clientpb.ReadClientResponse, error)
	SearchClientsByName func(context.Context, *clientpb.SearchClientsByNameRequest) (*clientpb.SearchClientsByNameResponse, error)
}

type LocationUseCases struct {
	ListLocations func(context.Context, *locationpb.ListLocationsRequest) (*locationpb.ListLocationsResponse, error)
}

type SupplierUseCases struct {
	ListSuppliers func(context.Context, *supplierpb.ListSuppliersRequest) (*supplierpb.ListSuppliersResponse, error)
}

type WorkspaceUseCases struct {
	ReadWorkspace func(context.Context, *workspacepb.ReadWorkspaceRequest) (*workspacepb.ReadWorkspaceResponse, error)
}

// -- Inventory ---------------------------------------------------------------

type InventoryUseCases struct {
	ListInventoryItems                func(context.Context, *inventoryitempb.ListInventoryItemsRequest) (*inventoryitempb.ListInventoryItemsResponse, error)
	CreateInventoryItem               func(context.Context, *inventoryitempb.CreateInventoryItemRequest) (*inventoryitempb.CreateInventoryItemResponse, error)
	ReadInventoryItem                 func(context.Context, *inventoryitempb.ReadInventoryItemRequest) (*inventoryitempb.ReadInventoryItemResponse, error)
	UpdateInventoryItem               func(context.Context, *inventoryitempb.UpdateInventoryItemRequest) (*inventoryitempb.UpdateInventoryItemResponse, error)
	DeleteInventoryItem               func(context.Context, *inventoryitempb.DeleteInventoryItemRequest) (*inventoryitempb.DeleteInventoryItemResponse, error)
	ListInventorySerials              func(context.Context, *inventoryserialpb.ListInventorySerialsRequest) (*inventoryserialpb.ListInventorySerialsResponse, error)
	CreateInventorySerial             func(context.Context, *inventoryserialpb.CreateInventorySerialRequest) (*inventoryserialpb.CreateInventorySerialResponse, error)
	ReadInventorySerial               func(context.Context, *inventoryserialpb.ReadInventorySerialRequest) (*inventoryserialpb.ReadInventorySerialResponse, error)
	UpdateInventorySerial             func(context.Context, *inventoryserialpb.UpdateInventorySerialRequest) (*inventoryserialpb.UpdateInventorySerialResponse, error)
	DeleteInventorySerial             func(context.Context, *inventoryserialpb.DeleteInventorySerialRequest) (*inventoryserialpb.DeleteInventorySerialResponse, error)
	ListInventoryTransactions         func(context.Context, *inventorytransactionpb.ListInventoryTransactionsRequest) (*inventorytransactionpb.ListInventoryTransactionsResponse, error)
	CreateInventoryTransaction        func(context.Context, *inventorytransactionpb.CreateInventoryTransactionRequest) (*inventorytransactionpb.CreateInventoryTransactionResponse, error)
	GetInventoryMovementsListPageData func(context.Context, *inventorytransactionpb.GetInventoryMovementsListPageDataRequest) (*inventorytransactionpb.GetInventoryMovementsListPageDataResponse, error)
	ListInventoryDepreciations        func(context.Context, *inventorydepreciationpb.ListInventoryDepreciationsRequest) (*inventorydepreciationpb.ListInventoryDepreciationsResponse, error)
	CreateInventoryDepreciation       func(context.Context, *inventorydepreciationpb.CreateInventoryDepreciationRequest) (*inventorydepreciationpb.CreateInventoryDepreciationResponse, error)
	ReadInventoryDepreciation         func(context.Context, *inventorydepreciationpb.ReadInventoryDepreciationRequest) (*inventorydepreciationpb.ReadInventoryDepreciationResponse, error)
	UpdateInventoryDepreciation       func(context.Context, *inventorydepreciationpb.UpdateInventoryDepreciationRequest) (*inventorydepreciationpb.UpdateInventoryDepreciationResponse, error)
	CreateInventorySerialHistory      func(context.Context, *inventoryserialhistorypb.CreateInventorySerialHistoryRequest) (*inventoryserialhistorypb.CreateInventorySerialHistoryResponse, error)
}

// -- Revenue -----------------------------------------------------------------

type RevenueUseCases struct {
	GetListPageData                  func(context.Context, *revenuepb.GetRevenueListPageDataRequest) (*revenuepb.GetRevenueListPageDataResponse, error)
	CreateRevenue                    func(context.Context, *revenuepb.CreateRevenueRequest) (*revenuepb.CreateRevenueResponse, error)
	ReadRevenue                      func(context.Context, *revenuepb.ReadRevenueRequest) (*revenuepb.ReadRevenueResponse, error)
	UpdateRevenue                    func(context.Context, *revenuepb.UpdateRevenueRequest) (*revenuepb.UpdateRevenueResponse, error)
	DeleteRevenue                    func(context.Context, *revenuepb.DeleteRevenueRequest) (*revenuepb.DeleteRevenueResponse, error)
	RecognizeRevenueFromSubscription func(context.Context, *revenuepb.CreateRevenueWithLineItemsRequest) (*revenuepb.CreateRevenueWithLineItemsResponse, error)
	CreateRevenueLineItem            func(context.Context, *revenuelineitempb.CreateRevenueLineItemRequest) (*revenuelineitempb.CreateRevenueLineItemResponse, error)
	ReadRevenueLineItem              func(context.Context, *revenuelineitempb.ReadRevenueLineItemRequest) (*revenuelineitempb.ReadRevenueLineItemResponse, error)
	UpdateRevenueLineItem            func(context.Context, *revenuelineitempb.UpdateRevenueLineItemRequest) (*revenuelineitempb.UpdateRevenueLineItemResponse, error)
	DeleteRevenueLineItem            func(context.Context, *revenuelineitempb.DeleteRevenueLineItemRequest) (*revenuelineitempb.DeleteRevenueLineItemResponse, error)
	ListRevenueLineItems             func(context.Context, *revenuelineitempb.ListRevenueLineItemsRequest) (*revenuelineitempb.ListRevenueLineItemsResponse, error)
	ListRevenueTaxLines              func(context.Context, *revenuetaxlinepb.ListRevenueTaxLinesRequest) (*revenuetaxlinepb.ListRevenueTaxLinesResponse, error)
	// Ex-helpers promoted to proto-defined use cases in Phase 0:
	ListRevenueRunCandidates func(context.Context, *revenuerunpb.ListRevenueRunCandidatesRequest) (*revenuerunpb.ListRevenueRunCandidatesResponse, error)
	GenerateRevenueRun       func(context.Context, *revenuerunpb.GenerateRevenueRunRequest) (*revenuerunpb.GenerateRevenueRunResponse, error)
	// 20260612-datasource-typed-path W5 — revenue_payment CRUD, typed path
	// replacing the centymo DataSource duck's ListSimple/Create/Read/Update/
	// Delete on the "revenue_payment" collection. Wired by service-admin in W7.
	// Nil-safe — the payment drawer + detail tab degrade to an empty state when
	// unwired (mock builds, half-wired composition root).
	RevenuePayment RevenuePaymentUseCases
}

// RevenuePaymentUseCases groups the typed revenue_payment CRUD closures the
// revenue payment drawer + detail-tab views need. Replaces the duck-typed
// DataSource.{Create,Read,Update,Delete,ListSimple}("revenue_payment") path.
// 20260612-datasource-typed-path W5.
type RevenuePaymentUseCases struct {
	CreateRevenuePayment func(context.Context, *revenuepaymentpb.CreateRevenuePaymentRequest) (*revenuepaymentpb.CreateRevenuePaymentResponse, error)
	ReadRevenuePayment   func(context.Context, *revenuepaymentpb.ReadRevenuePaymentRequest) (*revenuepaymentpb.ReadRevenuePaymentResponse, error)
	UpdateRevenuePayment func(context.Context, *revenuepaymentpb.UpdateRevenuePaymentRequest) (*revenuepaymentpb.UpdateRevenuePaymentResponse, error)
	DeleteRevenuePayment func(context.Context, *revenuepaymentpb.DeleteRevenuePaymentRequest) (*revenuepaymentpb.DeleteRevenuePaymentResponse, error)
	ListRevenuePayments  func(context.Context, *revenuepaymentpb.ListRevenuePaymentsRequest) (*revenuepaymentpb.ListRevenuePaymentsResponse, error)
}

// RevenueRunUseCases — repo-direct operations on the RevenueRun entity.
// ListRevenueRuns, ReadRevenueRun, and ListRevenueRunAttempts are pass-through
// calls on the RevenueRun domain service (not application-layer use cases).
// In service-admin's adapter, these are wired via
// uc.Revenue.Revenue.GenerateRevenueRun.RevenueRunRepo().
type RevenueRunUseCases struct {
	ListRevenueRuns        func(context.Context, *revenuerunpb.ListRevenueRunsRequest) (*revenuerunpb.ListRevenueRunsResponse, error)
	ReadRevenueRun         func(context.Context, *revenuerunpb.ReadRevenueRunRequest) (*revenuerunpb.ReadRevenueRunResponse, error)
	ListRevenueRunAttempts func(context.Context, *revenuerunpb.ListRevenueRunAttemptsRequest) (*revenuerunpb.ListRevenueRunAttemptsResponse, error)
}

// -- Product -----------------------------------------------------------------

type ProductUseCases struct {
	// Product CRUD
	ListProducts  func(context.Context, *productpb.ListProductsRequest) (*productpb.ListProductsResponse, error)
	ReadProduct   func(context.Context, *productpb.ReadProductRequest) (*productpb.ReadProductResponse, error)
	CreateProduct func(context.Context, *productpb.CreateProductRequest) (*productpb.CreateProductResponse, error)
	UpdateProduct func(context.Context, *productpb.UpdateProductRequest) (*productpb.UpdateProductResponse, error)
	DeleteProduct func(context.Context, *productpb.DeleteProductRequest) (*productpb.DeleteProductResponse, error)
	// ProductVariant CRUD
	ListProductVariants  func(context.Context, *productvariantpb.ListProductVariantsRequest) (*productvariantpb.ListProductVariantsResponse, error)
	ReadProductVariant   func(context.Context, *productvariantpb.ReadProductVariantRequest) (*productvariantpb.ReadProductVariantResponse, error)
	CreateProductVariant func(context.Context, *productvariantpb.CreateProductVariantRequest) (*productvariantpb.CreateProductVariantResponse, error)
	UpdateProductVariant func(context.Context, *productvariantpb.UpdateProductVariantRequest) (*productvariantpb.UpdateProductVariantResponse, error)
	DeleteProductVariant func(context.Context, *productvariantpb.DeleteProductVariantRequest) (*productvariantpb.DeleteProductVariantResponse, error)
	// ProductVariantOption
	ListProductVariantOptions  func(context.Context, *productvariantoptionpb.ListProductVariantOptionsRequest) (*productvariantoptionpb.ListProductVariantOptionsResponse, error)
	CreateProductVariantOption func(context.Context, *productvariantoptionpb.CreateProductVariantOptionRequest) (*productvariantoptionpb.CreateProductVariantOptionResponse, error)
	// DeleteProductVariantOption replaces the duck-typed
	// DataSource.HardDelete("product_variant_option", id) path used by the
	// variant detail's deleteVariantOptions cleanup. The product_variant_option
	// proto Delete is a hard delete on this junction table. Bound by
	// service-admin in W7; nil-safe — when unbound the cleanup logs + skips.
	// 20260612-datasource-typed-path W6.
	DeleteProductVariantOption func(context.Context, *productvariantoptionpb.DeleteProductVariantOptionRequest) (*productvariantoptionpb.DeleteProductVariantOptionResponse, error)
	// ProductOption CRUD
	ListProductOptions  func(context.Context, *productoptionpb.ListProductOptionsRequest) (*productoptionpb.ListProductOptionsResponse, error)
	ReadProductOption   func(context.Context, *productoptionpb.ReadProductOptionRequest) (*productoptionpb.ReadProductOptionResponse, error)
	CreateProductOption func(context.Context, *productoptionpb.CreateProductOptionRequest) (*productoptionpb.CreateProductOptionResponse, error)
	UpdateProductOption func(context.Context, *productoptionpb.UpdateProductOptionRequest) (*productoptionpb.UpdateProductOptionResponse, error)
	DeleteProductOption func(context.Context, *productoptionpb.DeleteProductOptionRequest) (*productoptionpb.DeleteProductOptionResponse, error)
	// ProductOptionValue CRUD
	ListProductOptionValues  func(context.Context, *productoptionvaluepb.ListProductOptionValuesRequest) (*productoptionvaluepb.ListProductOptionValuesResponse, error)
	ReadProductOptionValue   func(context.Context, *productoptionvaluepb.ReadProductOptionValueRequest) (*productoptionvaluepb.ReadProductOptionValueResponse, error)
	CreateProductOptionValue func(context.Context, *productoptionvaluepb.CreateProductOptionValueRequest) (*productoptionvaluepb.CreateProductOptionValueResponse, error)
	UpdateProductOptionValue func(context.Context, *productoptionvaluepb.UpdateProductOptionValueRequest) (*productoptionvaluepb.UpdateProductOptionValueResponse, error)
	DeleteProductOptionValue func(context.Context, *productoptionvaluepb.DeleteProductOptionValueRequest) (*productoptionvaluepb.DeleteProductOptionValueResponse, error)
	// ProductAttribute
	ListProductAttributes  func(context.Context, *productattributepb.ListProductAttributesRequest) (*productattributepb.ListProductAttributesResponse, error)
	CreateProductAttribute func(context.Context, *productattributepb.CreateProductAttributeRequest) (*productattributepb.CreateProductAttributeResponse, error)
	DeleteProductAttribute func(context.Context, *productattributepb.DeleteProductAttributeRequest) (*productattributepb.DeleteProductAttributeResponse, error)
	// Line CRUD (product "lines" — top-level category grouping)
	ListLines  func(context.Context, *linepb.ListLinesRequest) (*linepb.ListLinesResponse, error)
	ReadLine   func(context.Context, *linepb.ReadLineRequest) (*linepb.ReadLineResponse, error)
	CreateLine func(context.Context, *linepb.CreateLineRequest) (*linepb.CreateLineResponse, error)
	UpdateLine func(context.Context, *linepb.UpdateLineRequest) (*linepb.UpdateLineResponse, error)
	DeleteLine func(context.Context, *linepb.DeleteLineRequest) (*linepb.DeleteLineResponse, error)
	// ProductLine CRUD (product ↔ line membership)
	ListProductLines  func(context.Context, *productlinepb.ListProductLinesRequest) (*productlinepb.ListProductLinesResponse, error)
	ReadProductLine   func(context.Context, *productlinepb.ReadProductLineRequest) (*productlinepb.ReadProductLineResponse, error)
	CreateProductLine func(context.Context, *productlinepb.CreateProductLineRequest) (*productlinepb.CreateProductLineResponse, error)
	UpdateProductLine func(context.Context, *productlinepb.UpdateProductLineRequest) (*productlinepb.UpdateProductLineResponse, error)
	DeleteProductLine func(context.Context, *productlinepb.DeleteProductLineRequest) (*productlinepb.DeleteProductLineResponse, error)
	// ProductVariantImage
	ListProductVariantImages  func(context.Context, *productvariantimagepb.ListProductVariantImagesRequest) (*productvariantimagepb.ListProductVariantImagesResponse, error)
	CreateProductVariantImage func(context.Context, *productvariantimagepb.CreateProductVariantImageRequest) (*productvariantimagepb.CreateProductVariantImageResponse, error)
	DeleteProductVariantImage func(context.Context, *productvariantimagepb.DeleteProductVariantImageRequest) (*productvariantimagepb.DeleteProductVariantImageResponse, error)
	// ProductPlan CRUD
	ListProductPlans  func(context.Context, *productplanpb.ListProductPlansRequest) (*productplanpb.ListProductPlansResponse, error)
	ReadProductPlan   func(context.Context, *productplanpb.ReadProductPlanRequest) (*productplanpb.ReadProductPlanResponse, error)
	CreateProductPlan func(context.Context, *productplanpb.CreateProductPlanRequest) (*productplanpb.CreateProductPlanResponse, error)
	UpdateProductPlan func(context.Context, *productplanpb.UpdateProductPlanRequest) (*productplanpb.UpdateProductPlanResponse, error)
	DeleteProductPlan func(context.Context, *productplanpb.DeleteProductPlanRequest) (*productplanpb.DeleteProductPlanResponse, error)
	// PriceList CRUD
	FindApplicablePriceList func(context.Context, *pricelistpb.FindApplicablePriceListRequest) (*pricelistpb.FindApplicablePriceListResponse, error)
	ListPriceLists          func(context.Context, *pricelistpb.ListPriceListsRequest) (*pricelistpb.ListPriceListsResponse, error)
	ReadPriceList           func(context.Context, *pricelistpb.ReadPriceListRequest) (*pricelistpb.ReadPriceListResponse, error)
	CreatePriceList         func(context.Context, *pricelistpb.CreatePriceListRequest) (*pricelistpb.CreatePriceListResponse, error)
	UpdatePriceList         func(context.Context, *pricelistpb.UpdatePriceListRequest) (*pricelistpb.UpdatePriceListResponse, error)
	DeletePriceList         func(context.Context, *pricelistpb.DeletePriceListRequest) (*pricelistpb.DeletePriceListResponse, error)
	// PriceProduct
	ListPriceProducts  func(context.Context, *priceproductpb.ListPriceProductsRequest) (*priceproductpb.ListPriceProductsResponse, error)
	CreatePriceProduct func(context.Context, *priceproductpb.CreatePriceProductRequest) (*priceproductpb.CreatePriceProductResponse, error)
	DeletePriceProduct func(context.Context, *priceproductpb.DeletePriceProductRequest) (*priceproductpb.DeletePriceProductResponse, error)
	// Resource CRUD
	ListResources  func(context.Context, *resourcepb.ListResourcesRequest) (*resourcepb.ListResourcesResponse, error)
	ReadResource   func(context.Context, *resourcepb.ReadResourceRequest) (*resourcepb.ReadResourceResponse, error)
	CreateResource func(context.Context, *resourcepb.CreateResourceRequest) (*resourcepb.CreateResourceResponse, error)
	UpdateResource func(context.Context, *resourcepb.UpdateResourceRequest) (*resourcepb.UpdateResourceResponse, error)
	DeleteResource func(context.Context, *resourcepb.DeleteResourceRequest) (*resourcepb.DeleteResourceResponse, error)
	// Dashboard — centymo view-layer types (espyna internals are unreachable).
	// Nil-safe: service dashboard renders empty state when unset.
	GetServiceDashboard func(context.Context, *productdashboard.Request) (*productdashboard.Response, error)
}

// -- Plan --------------------------------------------------------------------

type PlanUseCases struct {
	ListPlans         func(context.Context, *planpb.ListPlansRequest) (*planpb.ListPlansResponse, error)
	ReadPlan          func(context.Context, *planpb.ReadPlanRequest) (*planpb.ReadPlanResponse, error)
	CreatePlan        func(context.Context, *planpb.CreatePlanRequest) (*planpb.CreatePlanResponse, error)
	UpdatePlan        func(context.Context, *planpb.UpdatePlanRequest) (*planpb.UpdatePlanResponse, error)
	DeletePlan        func(context.Context, *planpb.DeletePlanRequest) (*planpb.DeletePlanResponse, error)
	SearchPlansByName func(context.Context, *planpb.SearchPlansByNameRequest) (*planpb.SearchPlansByNameResponse, error)
	// Ex-helper promoted to proto-defined use case in Phase 0:
	CustomizePlanForClient func(context.Context, *planpb.CustomizePlanForClientRequest) (*planpb.CustomizePlanForClientResponse, error)
}

// -- PricePlan ---------------------------------------------------------------

type PricePlanUseCases struct {
	ListPricePlans         func(context.Context, *priceplanpb.ListPricePlansRequest) (*priceplanpb.ListPricePlansResponse, error)
	ReadPricePlan          func(context.Context, *priceplanpb.ReadPricePlanRequest) (*priceplanpb.ReadPricePlanResponse, error)
	CreatePricePlan        func(context.Context, *priceplanpb.CreatePricePlanRequest) (*priceplanpb.CreatePricePlanResponse, error)
	UpdatePricePlan        func(context.Context, *priceplanpb.UpdatePricePlanRequest) (*priceplanpb.UpdatePricePlanResponse, error)
	DeletePricePlan        func(context.Context, *priceplanpb.DeletePricePlanRequest) (*priceplanpb.DeletePricePlanResponse, error)
	ListProductPricePlans  func(context.Context, *productpriceplanpb.ListProductPricePlansRequest) (*productpriceplanpb.ListProductPricePlansResponse, error)
	CreateProductPricePlan func(context.Context, *productpriceplanpb.CreateProductPricePlanRequest) (*productpriceplanpb.CreateProductPricePlanResponse, error)
	UpdateProductPricePlan func(context.Context, *productpriceplanpb.UpdateProductPricePlanRequest) (*productpriceplanpb.UpdateProductPricePlanResponse, error)
	DeleteProductPricePlan func(context.Context, *productpriceplanpb.DeleteProductPricePlanRequest) (*productpriceplanpb.DeleteProductPricePlanResponse, error)
}

// -- PriceSchedule -----------------------------------------------------------

type PriceScheduleUseCases struct {
	ListPriceSchedules           func(context.Context, *priceschedulepb.ListPriceSchedulesRequest) (*priceschedulepb.ListPriceSchedulesResponse, error)
	ReadPriceSchedule            func(context.Context, *priceschedulepb.ReadPriceScheduleRequest) (*priceschedulepb.ReadPriceScheduleResponse, error)
	CreatePriceSchedule          func(context.Context, *priceschedulepb.CreatePriceScheduleRequest) (*priceschedulepb.CreatePriceScheduleResponse, error)
	UpdatePriceSchedule          func(context.Context, *priceschedulepb.UpdatePriceScheduleRequest) (*priceschedulepb.UpdatePriceScheduleResponse, error)
	DeletePriceSchedule          func(context.Context, *priceschedulepb.DeletePriceScheduleRequest) (*priceschedulepb.DeletePriceScheduleResponse, error)
	ListSubscriptionsByPricePlan func(context.Context, *subscriptionpb.ListSubscriptionsByPricePlanRequest) (*subscriptionpb.ListSubscriptionsByPricePlanResponse, error)
}

// -- SubscriptionGroup -------------------------------------------------------
//
// The education "section / cohort" cohort. Simple single-aggregate CRUD +
// List; no cross-entity orchestration. Closures use the subscription_group
// proto req/resp types and bind to espyna's
// uc.Subscription.SubscriptionGroup.* use cases (see engineblock.go).

type SubscriptionGroupUseCases struct {
	ListSubscriptionGroups  func(context.Context, *subscriptiongrouppb.ListSubscriptionGroupsRequest) (*subscriptiongrouppb.ListSubscriptionGroupsResponse, error)
	ReadSubscriptionGroup   func(context.Context, *subscriptiongrouppb.ReadSubscriptionGroupRequest) (*subscriptiongrouppb.ReadSubscriptionGroupResponse, error)
	CreateSubscriptionGroup func(context.Context, *subscriptiongrouppb.CreateSubscriptionGroupRequest) (*subscriptiongrouppb.CreateSubscriptionGroupResponse, error)
	UpdateSubscriptionGroup func(context.Context, *subscriptiongrouppb.UpdateSubscriptionGroupRequest) (*subscriptiongrouppb.UpdateSubscriptionGroupResponse, error)
	DeleteSubscriptionGroup func(context.Context, *subscriptiongrouppb.DeleteSubscriptionGroupRequest) (*subscriptiongrouppb.DeleteSubscriptionGroupResponse, error)
}

// -- SubscriptionGroupMember -------------------------------------------------
//
// Roster of members within a subscription_group (education cohort/section).
// Simple single-aggregate CRUD + List. Closures use the
// subscription_group_member proto req/resp types and bind to espyna's
// uc.Subscription.SubscriptionGroupMember.* use cases (see engineblock.go).

type SubscriptionGroupMemberUseCases struct {
	ListSubscriptionGroupMembers  func(context.Context, *subscriptiongroupmemberpb.ListSubscriptionGroupMembersRequest) (*subscriptiongroupmemberpb.ListSubscriptionGroupMembersResponse, error)
	ReadSubscriptionGroupMember   func(context.Context, *subscriptiongroupmemberpb.ReadSubscriptionGroupMemberRequest) (*subscriptiongroupmemberpb.ReadSubscriptionGroupMemberResponse, error)
	CreateSubscriptionGroupMember func(context.Context, *subscriptiongroupmemberpb.CreateSubscriptionGroupMemberRequest) (*subscriptiongroupmemberpb.CreateSubscriptionGroupMemberResponse, error)
	UpdateSubscriptionGroupMember func(context.Context, *subscriptiongroupmemberpb.UpdateSubscriptionGroupMemberRequest) (*subscriptiongroupmemberpb.UpdateSubscriptionGroupMemberResponse, error)
	DeleteSubscriptionGroupMember func(context.Context, *subscriptiongroupmemberpb.DeleteSubscriptionGroupMemberRequest) (*subscriptiongroupmemberpb.DeleteSubscriptionGroupMemberResponse, error)
}

// -- SubscriptionGroupWorkspaceUser ------------------------------------------
//
// Workspace-user (staff) access grants on a subscription_group. Single-aggregate
// CRUD + List; binds to uc.Subscription.SubscriptionGroupWorkspaceUser.*.

type SubscriptionGroupWorkspaceUserUseCases struct {
	ListSubscriptionGroupWorkspaceUsers  func(context.Context, *subscriptiongroupworkspaceuserpb.ListSubscriptionGroupWorkspaceUsersRequest) (*subscriptiongroupworkspaceuserpb.ListSubscriptionGroupWorkspaceUsersResponse, error)
	ReadSubscriptionGroupWorkspaceUser   func(context.Context, *subscriptiongroupworkspaceuserpb.ReadSubscriptionGroupWorkspaceUserRequest) (*subscriptiongroupworkspaceuserpb.ReadSubscriptionGroupWorkspaceUserResponse, error)
	CreateSubscriptionGroupWorkspaceUser func(context.Context, *subscriptiongroupworkspaceuserpb.CreateSubscriptionGroupWorkspaceUserRequest) (*subscriptiongroupworkspaceuserpb.CreateSubscriptionGroupWorkspaceUserResponse, error)
	UpdateSubscriptionGroupWorkspaceUser func(context.Context, *subscriptiongroupworkspaceuserpb.UpdateSubscriptionGroupWorkspaceUserRequest) (*subscriptiongroupworkspaceuserpb.UpdateSubscriptionGroupWorkspaceUserResponse, error)
	DeleteSubscriptionGroupWorkspaceUser func(context.Context, *subscriptiongroupworkspaceuserpb.DeleteSubscriptionGroupWorkspaceUserRequest) (*subscriptiongroupworkspaceuserpb.DeleteSubscriptionGroupWorkspaceUserResponse, error)
}

// -- SubscriptionGroupProductPlanStaff ---------------------------------------
//
// Staff-to-product-plan assignments scoped to a subscription_group. Single-
// aggregate CRUD + List; binds to uc.Subscription.SubscriptionGroupProductPlanStaff.*.

type SubscriptionGroupProductPlanStaffUseCases struct {
	ListSubscriptionGroupProductPlanStaffs  func(context.Context, *subscriptiongroupproductplanstaffpb.ListSubscriptionGroupProductPlanStaffsRequest) (*subscriptiongroupproductplanstaffpb.ListSubscriptionGroupProductPlanStaffsResponse, error)
	ReadSubscriptionGroupProductPlanStaff   func(context.Context, *subscriptiongroupproductplanstaffpb.ReadSubscriptionGroupProductPlanStaffRequest) (*subscriptiongroupproductplanstaffpb.ReadSubscriptionGroupProductPlanStaffResponse, error)
	CreateSubscriptionGroupProductPlanStaff func(context.Context, *subscriptiongroupproductplanstaffpb.CreateSubscriptionGroupProductPlanStaffRequest) (*subscriptiongroupproductplanstaffpb.CreateSubscriptionGroupProductPlanStaffResponse, error)
	UpdateSubscriptionGroupProductPlanStaff func(context.Context, *subscriptiongroupproductplanstaffpb.UpdateSubscriptionGroupProductPlanStaffRequest) (*subscriptiongroupproductplanstaffpb.UpdateSubscriptionGroupProductPlanStaffResponse, error)
	DeleteSubscriptionGroupProductPlanStaff func(context.Context, *subscriptiongroupproductplanstaffpb.DeleteSubscriptionGroupProductPlanStaffRequest) (*subscriptiongroupproductplanstaffpb.DeleteSubscriptionGroupProductPlanStaffResponse, error)
}

// -- PriceScheduleWorkspaceUser ----------------------------------------------
//
// Workspace-user (staff) access grants on a price_schedule. Single-aggregate
// CRUD + List; binds to uc.Subscription.PriceScheduleWorkspaceUser.*.

type PriceScheduleWorkspaceUserUseCases struct {
	ListPriceScheduleWorkspaceUsers  func(context.Context, *pricescheduleworkspaceuserpb.ListPriceScheduleWorkspaceUsersRequest) (*pricescheduleworkspaceuserpb.ListPriceScheduleWorkspaceUsersResponse, error)
	ReadPriceScheduleWorkspaceUser   func(context.Context, *pricescheduleworkspaceuserpb.ReadPriceScheduleWorkspaceUserRequest) (*pricescheduleworkspaceuserpb.ReadPriceScheduleWorkspaceUserResponse, error)
	CreatePriceScheduleWorkspaceUser func(context.Context, *pricescheduleworkspaceuserpb.CreatePriceScheduleWorkspaceUserRequest) (*pricescheduleworkspaceuserpb.CreatePriceScheduleWorkspaceUserResponse, error)
	UpdatePriceScheduleWorkspaceUser func(context.Context, *pricescheduleworkspaceuserpb.UpdatePriceScheduleWorkspaceUserRequest) (*pricescheduleworkspaceuserpb.UpdatePriceScheduleWorkspaceUserResponse, error)
	DeletePriceScheduleWorkspaceUser func(context.Context, *pricescheduleworkspaceuserpb.DeletePriceScheduleWorkspaceUserRequest) (*pricescheduleworkspaceuserpb.DeletePriceScheduleWorkspaceUserResponse, error)
}

// -- PlanGroup ---------------------------------------------------------------
//
// A named grouping of plans (product catalog bundling). Single-aggregate CRUD +
// List; binds to uc.Product.PlanGroup.*.

type PlanGroupUseCases struct {
	ListPlanGroups  func(context.Context, *plangrouppb.ListPlanGroupsRequest) (*plangrouppb.ListPlanGroupsResponse, error)
	ReadPlanGroup   func(context.Context, *plangrouppb.ReadPlanGroupRequest) (*plangrouppb.ReadPlanGroupResponse, error)
	CreatePlanGroup func(context.Context, *plangrouppb.CreatePlanGroupRequest) (*plangrouppb.CreatePlanGroupResponse, error)
	UpdatePlanGroup func(context.Context, *plangrouppb.UpdatePlanGroupRequest) (*plangrouppb.UpdatePlanGroupResponse, error)
	DeletePlanGroup func(context.Context, *plangrouppb.DeletePlanGroupRequest) (*plangrouppb.DeletePlanGroupResponse, error)
}

// -- PlanGroupPlan -----------------------------------------------------------
//
// Membership rows linking a plan_group to its plans. Single-aggregate CRUD +
// List; binds to uc.Product.PlanGroupPlan.*.

type PlanGroupPlanUseCases struct {
	ListPlanGroupPlans  func(context.Context, *plangroupplanpb.ListPlanGroupPlansRequest) (*plangroupplanpb.ListPlanGroupPlansResponse, error)
	ReadPlanGroupPlan   func(context.Context, *plangroupplanpb.ReadPlanGroupPlanRequest) (*plangroupplanpb.ReadPlanGroupPlanResponse, error)
	CreatePlanGroupPlan func(context.Context, *plangroupplanpb.CreatePlanGroupPlanRequest) (*plangroupplanpb.CreatePlanGroupPlanResponse, error)
	UpdatePlanGroupPlan func(context.Context, *plangroupplanpb.UpdatePlanGroupPlanRequest) (*plangroupplanpb.UpdatePlanGroupPlanResponse, error)
	DeletePlanGroupPlan func(context.Context, *plangroupplanpb.DeletePlanGroupPlanRequest) (*plangroupplanpb.DeletePlanGroupPlanResponse, error)
}

// -- ProductPlanStaff --------------------------------------------------------
//
// Staff assignments on a product_plan. Single-aggregate CRUD + List; binds to
// uc.Product.ProductPlanStaff.*.

type ProductPlanStaffUseCases struct {
	ListProductPlanStaffs  func(context.Context, *productplanstaffpb.ListProductPlanStaffsRequest) (*productplanstaffpb.ListProductPlanStaffsResponse, error)
	ReadProductPlanStaff   func(context.Context, *productplanstaffpb.ReadProductPlanStaffRequest) (*productplanstaffpb.ReadProductPlanStaffResponse, error)
	CreateProductPlanStaff func(context.Context, *productplanstaffpb.CreateProductPlanStaffRequest) (*productplanstaffpb.CreateProductPlanStaffResponse, error)
	UpdateProductPlanStaff func(context.Context, *productplanstaffpb.UpdateProductPlanStaffRequest) (*productplanstaffpb.UpdateProductPlanStaffResponse, error)
	DeleteProductPlanStaff func(context.Context, *productplanstaffpb.DeleteProductPlanStaffRequest) (*productplanstaffpb.DeleteProductPlanStaffResponse, error)
}

// -- LineWorkspaceUser -------------------------------------------------------
//
// Workspace-user (staff) access grants on a product line. Single-aggregate CRUD
// + List; binds to uc.Product.LineWorkspaceUser.*.

type LineWorkspaceUserUseCases struct {
	ListLineWorkspaceUsers  func(context.Context, *lineworkspaceuserpb.ListLineWorkspaceUsersRequest) (*lineworkspaceuserpb.ListLineWorkspaceUsersResponse, error)
	ReadLineWorkspaceUser   func(context.Context, *lineworkspaceuserpb.ReadLineWorkspaceUserRequest) (*lineworkspaceuserpb.ReadLineWorkspaceUserResponse, error)
	CreateLineWorkspaceUser func(context.Context, *lineworkspaceuserpb.CreateLineWorkspaceUserRequest) (*lineworkspaceuserpb.CreateLineWorkspaceUserResponse, error)
	UpdateLineWorkspaceUser func(context.Context, *lineworkspaceuserpb.UpdateLineWorkspaceUserRequest) (*lineworkspaceuserpb.UpdateLineWorkspaceUserResponse, error)
	DeleteLineWorkspaceUser func(context.Context, *lineworkspaceuserpb.DeleteLineWorkspaceUserRequest) (*lineworkspaceuserpb.DeleteLineWorkspaceUserResponse, error)
}

// -- PriceList ---------------------------------------------------------------

type PriceListUseCases struct {
	// PriceList CRUD fields live in ProductUseCases.
	// This stub is kept for the standalone wantPriceList() module path.
}

// -- Subscription ------------------------------------------------------------

type SubscriptionUseCases struct {
	GetSubscriptionListPageData func(context.Context, *subscriptionpb.GetSubscriptionListPageDataRequest) (*subscriptionpb.GetSubscriptionListPageDataResponse, error)
	GetSubscriptionItemPageData func(context.Context, *subscriptionpb.GetSubscriptionItemPageDataRequest) (*subscriptionpb.GetSubscriptionItemPageDataResponse, error)
	CreateSubscription          func(context.Context, *subscriptionpb.CreateSubscriptionRequest) (*subscriptionpb.CreateSubscriptionResponse, error)
	ReadSubscription            func(context.Context, *subscriptionpb.ReadSubscriptionRequest) (*subscriptionpb.ReadSubscriptionResponse, error)
	UpdateSubscription          func(context.Context, *subscriptionpb.UpdateSubscriptionRequest) (*subscriptionpb.UpdateSubscriptionResponse, error)
	DeleteSubscription          func(context.Context, *subscriptionpb.DeleteSubscriptionRequest) (*subscriptionpb.DeleteSubscriptionResponse, error)
	ListSubscriptions           func(context.Context, *subscriptionpb.ListSubscriptionsRequest) (*subscriptionpb.ListSubscriptionsResponse, error)
	// BillingEvent server methods (milestone billing).
	ListBillingEventsBySubscription func(context.Context, *billingeventpb.ListBillingEventsBySubscriptionRequest) (*billingeventpb.ListBillingEventsBySubscriptionResponse, error)
	SetBillingEventStatus           func(context.Context, *billingeventpb.SetBillingEventStatusRequest) (*billingeventpb.SetBillingEventStatusResponse, error)
	// Ex-helpers promoted to proto-defined use cases in Phase 0:
	MaterializeJobsForSubscription         func(context.Context, *subscriptionpb.MaterializeJobsForSubscriptionRequest) (*subscriptionpb.MaterializeJobsForSubscriptionResponse, error)
	MaterializeInstanceJobsForSubscription func(context.Context, *subscriptionpb.MaterializeInstanceJobsForSubscriptionRequest) (*subscriptionpb.MaterializeInstanceJobsForSubscriptionResponse, error)
}

// -- Collection (treasury) ---------------------------------------------------

type CollectionUseCases struct {
	ListCollections  func(context.Context, *collectionpb.ListCollectionsRequest) (*collectionpb.ListCollectionsResponse, error)
	ReadCollection   func(context.Context, *collectionpb.ReadCollectionRequest) (*collectionpb.ReadCollectionResponse, error)
	CreateCollection func(context.Context, *collectionpb.CreateCollectionRequest) (*collectionpb.CreateCollectionResponse, error)
	UpdateCollection func(context.Context, *collectionpb.UpdateCollectionRequest) (*collectionpb.UpdateCollectionResponse, error)
	DeleteCollection func(context.Context, *collectionpb.DeleteCollectionRequest) (*collectionpb.DeleteCollectionResponse, error)
	// Dashboard — centymo view-layer types (espyna internals are unreachable).
	// Nil-safe: cash dashboard renders empty state when unset.
	GetCashDashboard func(context.Context, *collectiondashboard.Request) (*collectiondashboard.Response, error)

	// 20260517-advance-cash-events Plan B Phase 3 — UNSCHEDULED workflow
	// closures. Service-admin wires these from espyna's
	// SettleUnscheduledAdvance / RefundUnscheduledAdvance / CancelAdvance use
	// cases on the treasury_collection side. Nil-safe — the views surface
	// disabled buttons + helpful tooltips when unwired.
	SettleUnscheduledAdvance func(ctx context.Context, in AdvanceSettleInput) (*AdvanceSettleOutput, error)
	RefundUnscheduledAdvance func(ctx context.Context, in AdvanceRefundInput) (*AdvanceRefundOutput, error)
	CancelAdvance            func(ctx context.Context, in AdvanceCancelInput) (*AdvanceCancelOutput, error)
}

// -- CollectionMethod (treasury) ---------------------------------------------

// CollectionMethodUseCases groups the typed collection_method reads the revenue
// payment drawer needs (Read for the payment_method name lookup, List for the
// drawer's method-select options). Replaces the duck-typed
// DataSource.{Read,ListSimple}("collection_method") path.
// 20260612-datasource-typed-path W5. Nil-safe — the drawer renders an empty
// method list + falls back to the raw method id when unwired.
type CollectionMethodUseCases struct {
	ReadCollectionMethod  func(context.Context, *collectionmethodpb.ReadCollectionMethodRequest) (*collectionmethodpb.ReadCollectionMethodResponse, error)
	ListCollectionMethods func(context.Context, *collectionmethodpb.ListCollectionMethodsRequest) (*collectionmethodpb.ListCollectionMethodsResponse, error)
}

// -- Disbursement (treasury) -------------------------------------------------

type DisbursementUseCases struct {
	ListDisbursements  func(context.Context, *disbursementpb.ListDisbursementsRequest) (*disbursementpb.ListDisbursementsResponse, error)
	ReadDisbursement   func(context.Context, *disbursementpb.ReadDisbursementRequest) (*disbursementpb.ReadDisbursementResponse, error)
	CreateDisbursement func(context.Context, *disbursementpb.CreateDisbursementRequest) (*disbursementpb.CreateDisbursementResponse, error)
	UpdateDisbursement func(context.Context, *disbursementpb.UpdateDisbursementRequest) (*disbursementpb.UpdateDisbursementResponse, error)
	DeleteDisbursement func(context.Context, *disbursementpb.DeleteDisbursementRequest) (*disbursementpb.DeleteDisbursementResponse, error)

	// 20260517-advance-cash-events Plan B Phase 3 — UNSCHEDULED workflow
	// closures, buying-side mirror of CollectionUseCases. Nil-safe.
	SettleUnscheduledAdvance func(ctx context.Context, in AdvanceSettleInput) (*AdvanceSettleOutput, error)
	RefundUnscheduledAdvance func(ctx context.Context, in AdvanceRefundInput) (*AdvanceRefundOutput, error)
	CancelAdvance            func(ctx context.Context, in AdvanceCancelInput) (*AdvanceCancelOutput, error)
}

// ---------------------------------------------------------------------------
// 20260517-advance-cash-events Plan B Phase 3 — view-typed input/output shapes
// for the UNSCHEDULED workflow drawers. Mirrored between treasury_collection
// and treasury_disbursement: the only side-specific value is the underlying
// use case the service-admin adapter binds the closure to.
// ---------------------------------------------------------------------------

// AdvanceSettleInput captures the operator-supplied fields the Settle drawer
// posts to the workflow closure.
type AdvanceSettleInput struct {
	AdvanceID       string
	Amount          int64  // centavos
	TargetAccountID string // optional fund / GL account
	Reason          string
}

// AdvanceSettleOutput is the response shape the view renders into a toast.
type AdvanceSettleOutput struct {
	NewRemainingAmount  int64
	NewRecognizedAmount int64
	NewStatus           string // ACTIVE | PARTIALLY_SETTLED | SETTLED
}

// AdvanceRefundInput captures the Refund drawer fields.
type AdvanceRefundInput struct {
	AdvanceID          string
	Amount             int64
	RefundMethod       string
	DestinationAccount string
	Reason             string
}

// AdvanceRefundOutput is the response shape the view renders into a toast.
type AdvanceRefundOutput struct {
	NewRemainingAmount int64
	NewStatus          string // ACTIVE | PARTIALLY_SETTLED | REFUNDED
}

// AdvanceCancelInput captures the Cancel drawer fields.
type AdvanceCancelInput struct {
	AdvanceID string
	Reason    string
}

// AdvanceCancelOutput is the response shape the view renders into a toast.
type AdvanceCancelOutput struct {
	NewStatus string // CANCELLED
}

// -- Treasury Advances (workspace dashboard) ---------------------------------

// TreasuryAdvancesUseCases groups the workspace-level dashboard callback.
// Lives in a struct (not flat on UseCases) so the AdvancesDashboard module
// can be selectively enabled via WithTreasuryAdvances() without coupling to
// the existing Collection/Disbursement modules.
type TreasuryAdvancesUseCases struct {
	// GetAdvancesDashboard returns the workspace-level summary. Nil-safe —
	// the dashboard view renders empty state when unset.
	GetAdvancesDashboard func(ctx context.Context, asOfDate string) (*AdvancesDashboardData, error)

	// 20260517-advance-cash-events Plan B Phase 7 — MILESTONE recognize
	// closures. Service-admin wires these from espyna's
	// RecognizeMilestoneAdvanceCollection / RecognizeMilestoneAdvanceDisbursement
	// use cases. Nil-safe — the view-layer button surfaces a disabled state
	// + helpful tooltip when unwired. The view-typed input/output shapes
	// live at the centymo package root (see advance_actions.go) so the
	// per-package view modules can import them without circling through
	// block/.
	RecognizeMilestoneAdvanceCollection   func(ctx context.Context, in treasurydomain.AdvanceRecognizeMilestoneInput) (*treasurydomain.AdvanceRecognizeMilestoneOutput, error)
	RecognizeMilestoneAdvanceDisbursement func(ctx context.Context, in treasurydomain.AdvanceRecognizeMilestoneInput) (*treasurydomain.AdvanceRecognizeMilestoneOutput, error)
}

// AdvancesDashboardData is the view-typed return shape for the workspace
// summary callback. Service-admin's block-level shim converts proto
// TreasuryCollection / TreasuryDisbursement rows into these view-friendly
// rows.
type AdvancesDashboardData struct {
	Outflows               []AdvancesDashboardRow
	Inflows                []AdvancesDashboardRow
	OutflowTotalRemaining  int64
	InflowTotalRemaining   int64
	OutflowActiveCount     int
	InflowActiveCount      int
	OutflowFullyRecognized int
	InflowFullyRecognized  int
	Currency               string
}

// AdvancesDashboardRow is the per-row shape the dashboard renders.
type AdvancesDashboardRow struct {
	ID               string
	ReferenceNumber  string
	CounterpartyName string
	Kind             string // raw enum string (e.g. "TIME_BASED")
	Status           string // raw enum string (e.g. "ACTIVE")
	Currency         string
	TotalAmount      int64
	RemainingAmount  int64
	RecognizedAmount int64
}

// -- Expenditure -------------------------------------------------------------

type ExpenditureUseCases struct {
	ListExpenditures          func(context.Context, *expenditurepb.ListExpendituresRequest) (*expenditurepb.ListExpendituresResponse, error)
	CreateExpenditure         func(context.Context, *expenditurepb.CreateExpenditureRequest) (*expenditurepb.CreateExpenditureResponse, error)
	ReadExpenditure           func(context.Context, *expenditurepb.ReadExpenditureRequest) (*expenditurepb.ReadExpenditureResponse, error)
	UpdateExpenditure         func(context.Context, *expenditurepb.UpdateExpenditureRequest) (*expenditurepb.UpdateExpenditureResponse, error)
	DeleteExpenditure         func(context.Context, *expenditurepb.DeleteExpenditureRequest) (*expenditurepb.DeleteExpenditureResponse, error)
	ListExpenditureCategories func(context.Context, *expenditurecategorypb.ListExpenditureCategoriesRequest) (*expenditurecategorypb.ListExpenditureCategoriesResponse, error)
	CreateExpenditureCategory func(context.Context, *expenditurecategorypb.CreateExpenditureCategoryRequest) (*expenditurecategorypb.CreateExpenditureCategoryResponse, error)
	ReadExpenditureCategory   func(context.Context, *expenditurecategorypb.ReadExpenditureCategoryRequest) (*expenditurecategorypb.ReadExpenditureCategoryResponse, error)
	UpdateExpenditureCategory func(context.Context, *expenditurecategorypb.UpdateExpenditureCategoryRequest) (*expenditurecategorypb.UpdateExpenditureCategoryResponse, error)
	DeleteExpenditureCategory func(context.Context, *expenditurecategorypb.DeleteExpenditureCategoryRequest) (*expenditurecategorypb.DeleteExpenditureCategoryResponse, error)
	CreateExpenditureLineItem func(context.Context, *expenditurelineitempb.CreateExpenditureLineItemRequest) (*expenditurelineitempb.CreateExpenditureLineItemResponse, error)
	ReadExpenditureLineItem   func(context.Context, *expenditurelineitempb.ReadExpenditureLineItemRequest) (*expenditurelineitempb.ReadExpenditureLineItemResponse, error)
	UpdateExpenditureLineItem func(context.Context, *expenditurelineitempb.UpdateExpenditureLineItemRequest) (*expenditurelineitempb.UpdateExpenditureLineItemResponse, error)
	DeleteExpenditureLineItem func(context.Context, *expenditurelineitempb.DeleteExpenditureLineItemRequest) (*expenditurelineitempb.DeleteExpenditureLineItemResponse, error)
	ListExpenditureLineItems  func(context.Context, *expenditurelineitempb.ListExpenditureLineItemsRequest) (*expenditurelineitempb.ListExpenditureLineItemsResponse, error)
	// Purchase orders (linked POs on supplier contract / procurement request detail).
	ListPurchaseOrders func(context.Context, *purchaseorderpb.ListPurchaseOrdersRequest) (*purchaseorderpb.ListPurchaseOrdersResponse, error)
	// Expense recognition CRUD + workflow ops (accessed via wantExpenseRecognition).
	ListExpenseRecognitions   func(context.Context, *expenserecognitionpb.ListExpenseRecognitionsRequest) (*expenserecognitionpb.ListExpenseRecognitionsResponse, error)
	ReadExpenseRecognition    func(context.Context, *expenserecognitionpb.ReadExpenseRecognitionRequest) (*expenserecognitionpb.ReadExpenseRecognitionResponse, error)
	DeleteExpenseRecognition  func(context.Context, *expenserecognitionpb.DeleteExpenseRecognitionRequest) (*expenserecognitionpb.DeleteExpenseRecognitionResponse, error)
	ReverseExpenseRecognition func(context.Context, *expenserecognitionpb.ReverseExpenseRecognitionRequest) (*expenserecognitionpb.ReverseExpenseRecognitionResponse, error)
	RecognizeFromExpenditure  func(context.Context, *expenserecognitionpb.RecognizeFromExpenditureRequest) (*expenserecognitionpb.RecognizeFromExpenditureResponse, error)
	RecognizeFromContract     func(context.Context, *expenserecognitionpb.RecognizeFromContractRequest) (*expenserecognitionpb.RecognizeFromContractResponse, error)
	// Expense recognition line CRUD.
	ListExpenseRecognitionLines  func(context.Context, *expenserecognitionlinepb.ListExpenseRecognitionLinesRequest) (*expenserecognitionlinepb.ListExpenseRecognitionLinesResponse, error)
	ReadExpenseRecognitionLine   func(context.Context, *expenserecognitionlinepb.ReadExpenseRecognitionLineRequest) (*expenserecognitionlinepb.ReadExpenseRecognitionLineResponse, error)
	CreateExpenseRecognitionLine func(context.Context, *expenserecognitionlinepb.CreateExpenseRecognitionLineRequest) (*expenserecognitionlinepb.CreateExpenseRecognitionLineResponse, error)
	UpdateExpenseRecognitionLine func(context.Context, *expenserecognitionlinepb.UpdateExpenseRecognitionLineRequest) (*expenserecognitionlinepb.UpdateExpenseRecognitionLineResponse, error)
	DeleteExpenseRecognitionLine func(context.Context, *expenserecognitionlinepb.DeleteExpenseRecognitionLineRequest) (*expenserecognitionlinepb.DeleteExpenseRecognitionLineResponse, error)
	// Accrued expense CRUD + workflow ops (accessed via wantAccruedExpense).
	ListAccruedExpenses  func(context.Context, *accruedexpensepb.ListAccruedExpensesRequest) (*accruedexpensepb.ListAccruedExpensesResponse, error)
	ReadAccruedExpense   func(context.Context, *accruedexpensepb.ReadAccruedExpenseRequest) (*accruedexpensepb.ReadAccruedExpenseResponse, error)
	CreateAccruedExpense func(context.Context, *accruedexpensepb.CreateAccruedExpenseRequest) (*accruedexpensepb.CreateAccruedExpenseResponse, error)
	UpdateAccruedExpense func(context.Context, *accruedexpensepb.UpdateAccruedExpenseRequest) (*accruedexpensepb.UpdateAccruedExpenseResponse, error)
	DeleteAccruedExpense func(context.Context, *accruedexpensepb.DeleteAccruedExpenseRequest) (*accruedexpensepb.DeleteAccruedExpenseResponse, error)
	SettleAccrual        func(context.Context, *accruedexpensepb.SettleAccrualRequest) (*accruedexpensepb.SettleAccrualResponse, error)
	ReverseAccrual       func(context.Context, *accruedexpensepb.ReverseAccrualRequest) (*accruedexpensepb.ReverseAccrualResponse, error)
	AccrueFromContract   func(context.Context, *accruedexpensepb.AccrueFromContractRequest) (*accruedexpensepb.AccrueFromContractResponse, error)
	// Accrued expense settlement (child of accrued expense; same proto package).
	ListAccruedExpenseSettlements  func(context.Context, *accruedexpensepb.ListAccruedExpenseSettlementsRequest) (*accruedexpensepb.ListAccruedExpenseSettlementsResponse, error)
	CreateAccruedExpenseSettlement func(context.Context, *accruedexpensepb.CreateAccruedExpenseSettlementRequest) (*accruedexpensepb.CreateAccruedExpenseSettlementResponse, error)
	ReadAccruedExpenseSettlement   func(context.Context, *accruedexpensepb.ReadAccruedExpenseSettlementRequest) (*accruedexpensepb.ReadAccruedExpenseSettlementResponse, error)
	UpdateAccruedExpenseSettlement func(context.Context, *accruedexpensepb.UpdateAccruedExpenseSettlementRequest) (*accruedexpensepb.UpdateAccruedExpenseSettlementResponse, error)
	DeleteAccruedExpenseSettlement func(context.Context, *accruedexpensepb.DeleteAccruedExpenseSettlementRequest) (*accruedexpensepb.DeleteAccruedExpenseSettlementResponse, error)
	// Dashboards — centymo view-layer types.
	GetPurchaseDashboard func(context.Context, *purchaseboard.Request) (*purchaseboard.Response, error)
	GetExpenseDashboard  func(context.Context, *expenseboard.Request) (*expenseboard.Response, error)

	// 20260517-advance-cash-events Plan B Phase 7 — SupplierBillingEvent
	// list + detail reads (buying-side MILESTONE anchor). Nil-safe — the
	// view module degrades to empty state when unwired.
	ListSupplierBillingEvents func(context.Context, *supplierbillingeventpb.ListSupplierBillingEventsRequest) (*supplierbillingeventpb.ListSupplierBillingEventsResponse, error)
	ReadSupplierBillingEvent  func(context.Context, *supplierbillingeventpb.ReadSupplierBillingEventRequest) (*supplierbillingeventpb.ReadSupplierBillingEventResponse, error)
}

// -- Supplier Contract -------------------------------------------------------

type SupplierContractUseCases struct {
	ListSupplierContracts                   func(context.Context, *suppliercontractpb.ListSupplierContractsRequest) (*suppliercontractpb.ListSupplierContractsResponse, error)
	ReadSupplierContract                    func(context.Context, *suppliercontractpb.ReadSupplierContractRequest) (*suppliercontractpb.ReadSupplierContractResponse, error)
	CreateSupplierContract                  func(context.Context, *suppliercontractpb.CreateSupplierContractRequest) (*suppliercontractpb.CreateSupplierContractResponse, error)
	UpdateSupplierContract                  func(context.Context, *suppliercontractpb.UpdateSupplierContractRequest) (*suppliercontractpb.UpdateSupplierContractResponse, error)
	DeleteSupplierContract                  func(context.Context, *suppliercontractpb.DeleteSupplierContractRequest) (*suppliercontractpb.DeleteSupplierContractResponse, error)
	ApproveSupplierContract                 func(context.Context, *suppliercontractpb.ApproveSupplierContractRequest) (*suppliercontractpb.ApproveSupplierContractResponse, error)
	TerminateSupplierContract               func(context.Context, *suppliercontractpb.TerminateSupplierContractRequest) (*suppliercontractpb.TerminateSupplierContractResponse, error)
	ListSupplierContractLines               func(context.Context, *suppliercontractlinepb.ListSupplierContractLinesRequest) (*suppliercontractlinepb.ListSupplierContractLinesResponse, error)
	ReadSupplierContractLine                func(context.Context, *suppliercontractlinepb.ReadSupplierContractLineRequest) (*suppliercontractlinepb.ReadSupplierContractLineResponse, error)
	CreateSupplierContractLine              func(context.Context, *suppliercontractlinepb.CreateSupplierContractLineRequest) (*suppliercontractlinepb.CreateSupplierContractLineResponse, error)
	UpdateSupplierContractLine              func(context.Context, *suppliercontractlinepb.UpdateSupplierContractLineRequest) (*suppliercontractlinepb.UpdateSupplierContractLineResponse, error)
	DeleteSupplierContractLine              func(context.Context, *suppliercontractlinepb.DeleteSupplierContractLineRequest) (*suppliercontractlinepb.DeleteSupplierContractLineResponse, error)
	ListSupplierContractPriceSchedules      func(context.Context, *suppliercontractpriceschedulepb.ListSupplierContractPriceSchedulesRequest) (*suppliercontractpriceschedulepb.ListSupplierContractPriceSchedulesResponse, error)
	ReadSupplierContractPriceSchedule       func(context.Context, *suppliercontractpriceschedulepb.ReadSupplierContractPriceScheduleRequest) (*suppliercontractpriceschedulepb.ReadSupplierContractPriceScheduleResponse, error)
	CreateSupplierContractPriceSchedule     func(context.Context, *suppliercontractpriceschedulepb.CreateSupplierContractPriceScheduleRequest) (*suppliercontractpriceschedulepb.CreateSupplierContractPriceScheduleResponse, error)
	UpdateSupplierContractPriceSchedule     func(context.Context, *suppliercontractpriceschedulepb.UpdateSupplierContractPriceScheduleRequest) (*suppliercontractpriceschedulepb.UpdateSupplierContractPriceScheduleResponse, error)
	DeleteSupplierContractPriceSchedule     func(context.Context, *suppliercontractpriceschedulepb.DeleteSupplierContractPriceScheduleRequest) (*suppliercontractpriceschedulepb.DeleteSupplierContractPriceScheduleResponse, error)
	ActivateSupplierContractPriceSchedule   func(context.Context, *suppliercontractpriceschedulepb.ActivateSupplierContractPriceScheduleRequest) (*suppliercontractpriceschedulepb.ActivateSupplierContractPriceScheduleResponse, error)
	SupersedeSupplierContractPriceSchedule  func(context.Context, *suppliercontractpriceschedulepb.SupersedeSupplierContractPriceScheduleRequest) (*suppliercontractpriceschedulepb.SupersedeSupplierContractPriceScheduleResponse, error)
	ListSupplierContractPriceScheduleLines  func(context.Context, *suppliercontractpriceschedulelinepb.ListSupplierContractPriceScheduleLinesRequest) (*suppliercontractpriceschedulelinepb.ListSupplierContractPriceScheduleLinesResponse, error)
	ReadSupplierContractPriceScheduleLine   func(context.Context, *suppliercontractpriceschedulelinepb.ReadSupplierContractPriceScheduleLineRequest) (*suppliercontractpriceschedulelinepb.ReadSupplierContractPriceScheduleLineResponse, error)
	CreateSupplierContractPriceScheduleLine func(context.Context, *suppliercontractpriceschedulelinepb.CreateSupplierContractPriceScheduleLineRequest) (*suppliercontractpriceschedulelinepb.CreateSupplierContractPriceScheduleLineResponse, error)
	UpdateSupplierContractPriceScheduleLine func(context.Context, *suppliercontractpriceschedulelinepb.UpdateSupplierContractPriceScheduleLineRequest) (*suppliercontractpriceschedulelinepb.UpdateSupplierContractPriceScheduleLineResponse, error)
	DeleteSupplierContractPriceScheduleLine func(context.Context, *suppliercontractpriceschedulelinepb.DeleteSupplierContractPriceScheduleLineRequest) (*suppliercontractpriceschedulelinepb.DeleteSupplierContractPriceScheduleLineResponse, error)
	// Expense recognition + lines (triggered from contract workflow).
	ListExpenseRecognitions      func(context.Context, *expenserecognitionpb.ListExpenseRecognitionsRequest) (*expenserecognitionpb.ListExpenseRecognitionsResponse, error)
	ReadExpenseRecognition       func(context.Context, *expenserecognitionpb.ReadExpenseRecognitionRequest) (*expenserecognitionpb.ReadExpenseRecognitionResponse, error)
	ListExpenseRecognitionLines  func(context.Context, *expenserecognitionlinepb.ListExpenseRecognitionLinesRequest) (*expenserecognitionlinepb.ListExpenseRecognitionLinesResponse, error)
	ReadExpenseRecognitionLine   func(context.Context, *expenserecognitionlinepb.ReadExpenseRecognitionLineRequest) (*expenserecognitionlinepb.ReadExpenseRecognitionLineResponse, error)
	ListProcurementRequests      func(context.Context, *procurementrequestpb.ListProcurementRequestsRequest) (*procurementrequestpb.ListProcurementRequestsResponse, error)
	ReadProcurementRequest       func(context.Context, *procurementrequestpb.ReadProcurementRequestRequest) (*procurementrequestpb.ReadProcurementRequestResponse, error)
	CreateProcurementRequest     func(context.Context, *procurementrequestpb.CreateProcurementRequestRequest) (*procurementrequestpb.CreateProcurementRequestResponse, error)
	UpdateProcurementRequest     func(context.Context, *procurementrequestpb.UpdateProcurementRequestRequest) (*procurementrequestpb.UpdateProcurementRequestResponse, error)
	DeleteProcurementRequest     func(context.Context, *procurementrequestpb.DeleteProcurementRequestRequest) (*procurementrequestpb.DeleteProcurementRequestResponse, error)
	SubmitProcurementRequest     func(context.Context, *procurementrequestpb.SubmitProcurementRequestRequest) (*procurementrequestpb.SubmitProcurementRequestResponse, error)
	ApproveProcurementRequest    func(context.Context, *procurementrequestpb.ApproveProcurementRequestRequest) (*procurementrequestpb.ApproveProcurementRequestResponse, error)
	RejectProcurementRequest     func(context.Context, *procurementrequestpb.RejectProcurementRequestRequest) (*procurementrequestpb.RejectProcurementRequestResponse, error)
	SpawnProcurementRequestPO    func(context.Context, *procurementrequestpb.SpawnPurchaseOrderRequest) (*procurementrequestpb.SpawnPurchaseOrderResponse, error)
	ListProcurementRequestLines  func(context.Context, *procurementrequestlinepb.ListProcurementRequestLinesRequest) (*procurementrequestlinepb.ListProcurementRequestLinesResponse, error)
	ReadProcurementRequestLine   func(context.Context, *procurementrequestlinepb.ReadProcurementRequestLineRequest) (*procurementrequestlinepb.ReadProcurementRequestLineResponse, error)
	CreateProcurementRequestLine func(context.Context, *procurementrequestlinepb.CreateProcurementRequestLineRequest) (*procurementrequestlinepb.CreateProcurementRequestLineResponse, error)
	UpdateProcurementRequestLine func(context.Context, *procurementrequestlinepb.UpdateProcurementRequestLineRequest) (*procurementrequestlinepb.UpdateProcurementRequestLineResponse, error)
	DeleteProcurementRequestLine func(context.Context, *procurementrequestlinepb.DeleteProcurementRequestLineRequest) (*procurementrequestlinepb.DeleteProcurementRequestLineResponse, error)
}

// -- Procurement (P3 supplier subscriptions) ---------------------------------

// ProcurementUseCases groups the P3 supplier-subscription modules
// (CostSchedule, SupplierPlan, CostPlan, SupplierProductPlan,
// SupplierProductCostPlan, SupplierSubscription). Mirrors espyna's
// ProcurementUseCases shape so the container adapter is a direct field-to-field
// mapping.
type ProcurementUseCases struct {
	CostSchedule            CostScheduleUseCases
	SupplierPlan            SupplierPlanUseCases
	CostPlan                CostPlanUseCases
	SupplierProductPlan     SupplierProductPlanUseCases
	SupplierProductCostPlan SupplierProductCostPlanUseCases
	SupplierSubscription    SupplierSubscriptionUseCases
}

type CostScheduleUseCases struct {
	ListCostSchedules           func(context.Context, *costschedulepb.ListCostSchedulesRequest) (*costschedulepb.ListCostSchedulesResponse, error)
	ReadCostSchedule            func(context.Context, *costschedulepb.ReadCostScheduleRequest) (*costschedulepb.ReadCostScheduleResponse, error)
	CreateCostSchedule          func(context.Context, *costschedulepb.CreateCostScheduleRequest) (*costschedulepb.CreateCostScheduleResponse, error)
	UpdateCostSchedule          func(context.Context, *costschedulepb.UpdateCostScheduleRequest) (*costschedulepb.UpdateCostScheduleResponse, error)
	DeleteCostSchedule          func(context.Context, *costschedulepb.DeleteCostScheduleRequest) (*costschedulepb.DeleteCostScheduleResponse, error)
	GetCostScheduleListPageData func(context.Context, *costschedulepb.GetCostScheduleListPageDataRequest) (*costschedulepb.GetCostScheduleListPageDataResponse, error)
	GetCostScheduleItemPageData func(context.Context, *costschedulepb.GetCostScheduleItemPageDataRequest) (*costschedulepb.GetCostScheduleItemPageDataResponse, error)
}

type SupplierPlanUseCases struct {
	ListSupplierPlans           func(context.Context, *supplierplanpb.ListSupplierPlansRequest) (*supplierplanpb.ListSupplierPlansResponse, error)
	ReadSupplierPlan            func(context.Context, *supplierplanpb.ReadSupplierPlanRequest) (*supplierplanpb.ReadSupplierPlanResponse, error)
	CreateSupplierPlan          func(context.Context, *supplierplanpb.CreateSupplierPlanRequest) (*supplierplanpb.CreateSupplierPlanResponse, error)
	UpdateSupplierPlan          func(context.Context, *supplierplanpb.UpdateSupplierPlanRequest) (*supplierplanpb.UpdateSupplierPlanResponse, error)
	DeleteSupplierPlan          func(context.Context, *supplierplanpb.DeleteSupplierPlanRequest) (*supplierplanpb.DeleteSupplierPlanResponse, error)
	GetSupplierPlanListPageData func(context.Context, *supplierplanpb.GetSupplierPlanListPageDataRequest) (*supplierplanpb.GetSupplierPlanListPageDataResponse, error)
	GetSupplierPlanItemPageData func(context.Context, *supplierplanpb.GetSupplierPlanItemPageDataRequest) (*supplierplanpb.GetSupplierPlanItemPageDataResponse, error)
}

type CostPlanUseCases struct {
	ListCostPlans           func(context.Context, *costplanpb.ListCostPlansRequest) (*costplanpb.ListCostPlansResponse, error)
	ReadCostPlan            func(context.Context, *costplanpb.ReadCostPlanRequest) (*costplanpb.ReadCostPlanResponse, error)
	CreateCostPlan          func(context.Context, *costplanpb.CreateCostPlanRequest) (*costplanpb.CreateCostPlanResponse, error)
	UpdateCostPlan          func(context.Context, *costplanpb.UpdateCostPlanRequest) (*costplanpb.UpdateCostPlanResponse, error)
	DeleteCostPlan          func(context.Context, *costplanpb.DeleteCostPlanRequest) (*costplanpb.DeleteCostPlanResponse, error)
	GetCostPlanListPageData func(context.Context, *costplanpb.GetCostPlanListPageDataRequest) (*costplanpb.GetCostPlanListPageDataResponse, error)
	GetCostPlanItemPageData func(context.Context, *costplanpb.GetCostPlanItemPageDataRequest) (*costplanpb.GetCostPlanItemPageDataResponse, error)
}

type SupplierProductPlanUseCases struct {
	ListSupplierProductPlans           func(context.Context, *supplierproductplanpb.ListSupplierProductPlansRequest) (*supplierproductplanpb.ListSupplierProductPlansResponse, error)
	ReadSupplierProductPlan            func(context.Context, *supplierproductplanpb.ReadSupplierProductPlanRequest) (*supplierproductplanpb.ReadSupplierProductPlanResponse, error)
	CreateSupplierProductPlan          func(context.Context, *supplierproductplanpb.CreateSupplierProductPlanRequest) (*supplierproductplanpb.CreateSupplierProductPlanResponse, error)
	UpdateSupplierProductPlan          func(context.Context, *supplierproductplanpb.UpdateSupplierProductPlanRequest) (*supplierproductplanpb.UpdateSupplierProductPlanResponse, error)
	DeleteSupplierProductPlan          func(context.Context, *supplierproductplanpb.DeleteSupplierProductPlanRequest) (*supplierproductplanpb.DeleteSupplierProductPlanResponse, error)
	GetSupplierProductPlanListPageData func(context.Context, *supplierproductplanpb.GetSupplierProductPlanListPageDataRequest) (*supplierproductplanpb.GetSupplierProductPlanListPageDataResponse, error)
	GetSupplierProductPlanItemPageData func(context.Context, *supplierproductplanpb.GetSupplierProductPlanItemPageDataRequest) (*supplierproductplanpb.GetSupplierProductPlanItemPageDataResponse, error)
}

type SupplierProductCostPlanUseCases struct {
	ListSupplierProductCostPlans           func(context.Context, *supplierproductcostplanpb.ListSupplierProductCostPlansRequest) (*supplierproductcostplanpb.ListSupplierProductCostPlansResponse, error)
	ReadSupplierProductCostPlan            func(context.Context, *supplierproductcostplanpb.ReadSupplierProductCostPlanRequest) (*supplierproductcostplanpb.ReadSupplierProductCostPlanResponse, error)
	CreateSupplierProductCostPlan          func(context.Context, *supplierproductcostplanpb.CreateSupplierProductCostPlanRequest) (*supplierproductcostplanpb.CreateSupplierProductCostPlanResponse, error)
	UpdateSupplierProductCostPlan          func(context.Context, *supplierproductcostplanpb.UpdateSupplierProductCostPlanRequest) (*supplierproductcostplanpb.UpdateSupplierProductCostPlanResponse, error)
	DeleteSupplierProductCostPlan          func(context.Context, *supplierproductcostplanpb.DeleteSupplierProductCostPlanRequest) (*supplierproductcostplanpb.DeleteSupplierProductCostPlanResponse, error)
	GetSupplierProductCostPlanItemPageData func(context.Context, *supplierproductcostplanpb.GetSupplierProductCostPlanItemPageDataRequest) (*supplierproductcostplanpb.GetSupplierProductCostPlanItemPageDataResponse, error)
}

type SupplierSubscriptionUseCases struct {
	ListSupplierSubscriptions           func(context.Context, *suppliersubscriptionpb.ListSupplierSubscriptionsRequest) (*suppliersubscriptionpb.ListSupplierSubscriptionsResponse, error)
	ReadSupplierSubscription            func(context.Context, *suppliersubscriptionpb.ReadSupplierSubscriptionRequest) (*suppliersubscriptionpb.ReadSupplierSubscriptionResponse, error)
	CreateSupplierSubscription          func(context.Context, *suppliersubscriptionpb.CreateSupplierSubscriptionRequest) (*suppliersubscriptionpb.CreateSupplierSubscriptionResponse, error)
	UpdateSupplierSubscription          func(context.Context, *suppliersubscriptionpb.UpdateSupplierSubscriptionRequest) (*suppliersubscriptionpb.UpdateSupplierSubscriptionResponse, error)
	DeleteSupplierSubscription          func(context.Context, *suppliersubscriptionpb.DeleteSupplierSubscriptionRequest) (*suppliersubscriptionpb.DeleteSupplierSubscriptionResponse, error)
	GetSupplierSubscriptionListPageData func(context.Context, *suppliersubscriptionpb.GetSupplierSubscriptionListPageDataRequest) (*suppliersubscriptionpb.GetSupplierSubscriptionListPageDataResponse, error)
	GetSupplierSubscriptionItemPageData func(context.Context, *suppliersubscriptionpb.GetSupplierSubscriptionItemPageDataRequest) (*suppliersubscriptionpb.GetSupplierSubscriptionItemPageDataResponse, error)
}

// -- Operation ---------------------------------------------------------------

type OperationUseCases struct {
	JobTemplate         JobTemplateUseCases
	JobTemplatePhase    JobTemplatePhaseUseCases
	JobTemplateTask     JobTemplateTaskUseCases
	JobTemplateRelation JobTemplateRelationUseCases
	Job                 JobUseCases
	JobPhase            JobPhaseUseCases
	JobActivity         JobActivityUseCases
}

type JobTemplateUseCases struct {
	ListJobTemplates func(context.Context, *jobtemplatepb.ListJobTemplatesRequest) (*jobtemplatepb.ListJobTemplatesResponse, error)
	ReadJobTemplate  func(context.Context, *jobtemplatepb.ReadJobTemplateRequest) (*jobtemplatepb.ReadJobTemplateResponse, error)
}

type JobTemplatePhaseUseCases struct {
	ListByJobTemplate func(context.Context, *jobtemplatedepb.ListByJobTemplateRequest) (*jobtemplatedepb.ListByJobTemplateResponse, error)
}

type JobTemplateTaskUseCases struct {
	ListByPhase func(context.Context, *jobttemplatetaskpb.ListJobTemplateTasksByPhaseRequest) (*jobttemplatetaskpb.ListJobTemplateTasksByPhaseResponse, error)
}

// JobTemplateRelationUseCases wraps the JobTemplateRelationDomainServiceServer
// ListByParent method as a typed closure. In service-admin's adapter, this is
// wired directly from uc.Operation.JobTemplateRelation.ListByParent (the
// domain server method has the exact same signature).
type JobTemplateRelationUseCases struct {
	ListByParent func(context.Context, *jobtemplaterelationpb.ListJobTemplateRelationsByParentRequest) (*jobtemplaterelationpb.ListJobTemplateRelationsByParentResponse, error)
}

type JobUseCases struct {
	GetJobsByOrigin func(context.Context, *jobpb.GetJobsByOriginRequest) (*jobpb.GetJobsByOriginResponse, error)
}

type JobPhaseUseCases struct {
	ListByJob func(context.Context, *jobphasepb.ListJobPhasesByJobRequest) (*jobphasepb.ListJobPhasesByJobResponse, error)
}

type JobActivityUseCases struct {
	ReadJobActivity func(context.Context, *jobactivitypb.ReadJobActivityRequest) (*jobactivitypb.ReadJobActivityResponse, error)
}

// ---------------------------------------------------------------------------
// RequireFor — startup-time completeness validator
// ---------------------------------------------------------------------------

// RequireFor validates that the fields needed for cfg's enabled modules are
// non-nil. Called at Block() entry, BEFORE any module wiring runs.
// Missing field → startup error naming the exact field path.
func (u *UseCases) RequireFor(cfg *blockConfig) error {
	if u == nil {
		return fmt.Errorf("centymo.Block: WithUseCases(...) was not supplied")
	}

	var missing []string
	check := func(ok bool, name string) {
		if !ok {
			missing = append(missing, name)
		}
	}

	// ExtractUserID is required when supplier workflow modules are enabled.
	// Those modules call it inside action closures (ApproveSupplierContract,
	// ActivateSupplierContractPriceSchedule, ApproveProcurementRequest) without
	// nil-guards; a nil here would cause a runtime panic — fail at startup instead.
	if cfg.wantSupplierContract() || cfg.wantSupplierContractPriceSchedule() || cfg.wantProcurementRequest() {
		check(u.ExtractUserID != nil, "UseCases.ExtractUserID")
	}

	if cfg.wantInventory() {
		check(u.Inventory.ListInventoryItems != nil, "UseCases.Inventory.ListInventoryItems")
		check(u.Inventory.CreateInventoryItem != nil, "UseCases.Inventory.CreateInventoryItem")
		check(u.Inventory.ReadInventoryItem != nil, "UseCases.Inventory.ReadInventoryItem")
		check(u.Inventory.UpdateInventoryItem != nil, "UseCases.Inventory.UpdateInventoryItem")
		check(u.Inventory.DeleteInventoryItem != nil, "UseCases.Inventory.DeleteInventoryItem")
	}

	if cfg.wantRevenue() {
		check(u.Revenue.GetListPageData != nil, "UseCases.Revenue.GetListPageData")
		check(u.Revenue.CreateRevenue != nil, "UseCases.Revenue.CreateRevenue")
		check(u.Revenue.ReadRevenue != nil, "UseCases.Revenue.ReadRevenue")
		check(u.Revenue.UpdateRevenue != nil, "UseCases.Revenue.UpdateRevenue")
		check(u.Revenue.DeleteRevenue != nil, "UseCases.Revenue.DeleteRevenue")
	}

	if cfg.wantPricePlan() {
		check(u.PricePlan.ListPricePlans != nil, "UseCases.PricePlan.ListPricePlans")
		check(u.PricePlan.ReadPricePlan != nil, "UseCases.PricePlan.ReadPricePlan")
		check(u.PricePlan.CreatePricePlan != nil, "UseCases.PricePlan.CreatePricePlan")
		check(u.PricePlan.UpdatePricePlan != nil, "UseCases.PricePlan.UpdatePricePlan")
		check(u.PricePlan.DeletePricePlan != nil, "UseCases.PricePlan.DeletePricePlan")
	}

	if cfg.wantPriceSchedule() {
		check(u.PriceSchedule.ListPriceSchedules != nil, "UseCases.PriceSchedule.ListPriceSchedules")
		check(u.PriceSchedule.ReadPriceSchedule != nil, "UseCases.PriceSchedule.ReadPriceSchedule")
		check(u.PriceSchedule.CreatePriceSchedule != nil, "UseCases.PriceSchedule.CreatePriceSchedule")
		check(u.PriceSchedule.UpdatePriceSchedule != nil, "UseCases.PriceSchedule.UpdatePriceSchedule")
		check(u.PriceSchedule.DeletePriceSchedule != nil, "UseCases.PriceSchedule.DeletePriceSchedule")
	}

	if cfg.wantPriceList() {
		check(u.Product.ListPriceLists != nil, "UseCases.Product.ListPriceLists")
		check(u.Product.ReadPriceList != nil, "UseCases.Product.ReadPriceList")
		check(u.Product.CreatePriceList != nil, "UseCases.Product.CreatePriceList")
		check(u.Product.UpdatePriceList != nil, "UseCases.Product.UpdatePriceList")
		check(u.Product.DeletePriceList != nil, "UseCases.Product.DeletePriceList")
		check(u.Product.ListPriceProducts != nil, "UseCases.Product.ListPriceProducts")
		check(u.Product.CreatePriceProduct != nil, "UseCases.Product.CreatePriceProduct")
		check(u.Product.DeletePriceProduct != nil, "UseCases.Product.DeletePriceProduct")
	}

	if cfg.wantPlan() {
		check(u.Plan.ListPlans != nil, "UseCases.Plan.ListPlans")
		check(u.Plan.ReadPlan != nil, "UseCases.Plan.ReadPlan")
		check(u.Plan.CreatePlan != nil, "UseCases.Plan.CreatePlan")
		check(u.Plan.UpdatePlan != nil, "UseCases.Plan.UpdatePlan")
		check(u.Plan.DeletePlan != nil, "UseCases.Plan.DeletePlan")
	}

	if cfg.wantSubscription() {
		check(u.Subscription.GetSubscriptionListPageData != nil, "UseCases.Subscription.GetSubscriptionListPageData")
		check(u.Subscription.CreateSubscription != nil, "UseCases.Subscription.CreateSubscription")
		check(u.Subscription.ReadSubscription != nil, "UseCases.Subscription.ReadSubscription")
		check(u.Subscription.UpdateSubscription != nil, "UseCases.Subscription.UpdateSubscription")
		check(u.Subscription.DeleteSubscription != nil, "UseCases.Subscription.DeleteSubscription")
	}

	if cfg.wantCollection() {
		check(u.Collection.ListCollections != nil, "UseCases.Collection.ListCollections")
		check(u.Collection.CreateCollection != nil, "UseCases.Collection.CreateCollection")
		check(u.Collection.ReadCollection != nil, "UseCases.Collection.ReadCollection")
		check(u.Collection.UpdateCollection != nil, "UseCases.Collection.UpdateCollection")
		check(u.Collection.DeleteCollection != nil, "UseCases.Collection.DeleteCollection")
	}

	if cfg.wantDisbursement() {
		check(u.Disbursement.ListDisbursements != nil, "UseCases.Disbursement.ListDisbursements")
		check(u.Disbursement.CreateDisbursement != nil, "UseCases.Disbursement.CreateDisbursement")
		check(u.Disbursement.ReadDisbursement != nil, "UseCases.Disbursement.ReadDisbursement")
		check(u.Disbursement.UpdateDisbursement != nil, "UseCases.Disbursement.UpdateDisbursement")
		check(u.Disbursement.DeleteDisbursement != nil, "UseCases.Disbursement.DeleteDisbursement")
	}

	if cfg.wantExpenditure() {
		check(u.Expenditure.ListExpenditures != nil, "UseCases.Expenditure.ListExpenditures")
		check(u.Expenditure.CreateExpenditure != nil, "UseCases.Expenditure.CreateExpenditure")
		check(u.Expenditure.ReadExpenditure != nil, "UseCases.Expenditure.ReadExpenditure")
		check(u.Expenditure.UpdateExpenditure != nil, "UseCases.Expenditure.UpdateExpenditure")
		check(u.Expenditure.DeleteExpenditure != nil, "UseCases.Expenditure.DeleteExpenditure")
	}

	if cfg.wantRevenueRun() {
		check(u.RevenueRun.ListRevenueRuns != nil, "UseCases.RevenueRun.ListRevenueRuns")
		check(u.RevenueRun.ReadRevenueRun != nil, "UseCases.RevenueRun.ReadRevenueRun")
		check(u.RevenueRun.ListRevenueRunAttempts != nil, "UseCases.RevenueRun.ListRevenueRunAttempts")
		check(u.Revenue.ListRevenueRunCandidates != nil, "UseCases.Revenue.ListRevenueRunCandidates")
		check(u.Revenue.GenerateRevenueRun != nil, "UseCases.Revenue.GenerateRevenueRun")
	}

	if cfg.wantSupplierContract() {
		check(u.SupplierContract.ListSupplierContracts != nil, "UseCases.SupplierContract.ListSupplierContracts")
		check(u.SupplierContract.CreateSupplierContract != nil, "UseCases.SupplierContract.CreateSupplierContract")
		check(u.SupplierContract.ReadSupplierContract != nil, "UseCases.SupplierContract.ReadSupplierContract")
		check(u.SupplierContract.UpdateSupplierContract != nil, "UseCases.SupplierContract.UpdateSupplierContract")
		check(u.SupplierContract.DeleteSupplierContract != nil, "UseCases.SupplierContract.DeleteSupplierContract")
	}

	if cfg.wantSupplierContractPriceSchedule() {
		check(u.SupplierContract.ListSupplierContractPriceSchedules != nil, "UseCases.SupplierContract.ListSupplierContractPriceSchedules")
		check(u.SupplierContract.CreateSupplierContractPriceSchedule != nil, "UseCases.SupplierContract.CreateSupplierContractPriceSchedule")
	}

	if cfg.wantExpenseRecognition() {
		check(u.Expenditure.ListExpenseRecognitions != nil, "UseCases.Expenditure.ListExpenseRecognitions")
	}

	if cfg.wantAccruedExpense() {
		check(u.Expenditure.ListAccruedExpenses != nil, "UseCases.Expenditure.ListAccruedExpenses")
	}

	if cfg.wantCostSchedule() {
		check(u.Procurement.CostSchedule.ListCostSchedules != nil, "UseCases.Procurement.CostSchedule.ListCostSchedules")
		check(u.Procurement.CostSchedule.CreateCostSchedule != nil, "UseCases.Procurement.CostSchedule.CreateCostSchedule")
	}

	if cfg.wantSupplierPlan() {
		check(u.Procurement.SupplierPlan.ListSupplierPlans != nil, "UseCases.Procurement.SupplierPlan.ListSupplierPlans")
		check(u.Procurement.SupplierPlan.CreateSupplierPlan != nil, "UseCases.Procurement.SupplierPlan.CreateSupplierPlan")
	}

	if cfg.wantCostPlan() {
		check(u.Procurement.CostPlan.ListCostPlans != nil, "UseCases.Procurement.CostPlan.ListCostPlans")
		check(u.Procurement.CostPlan.CreateCostPlan != nil, "UseCases.Procurement.CostPlan.CreateCostPlan")
	}

	if cfg.wantSupplierSubscription() {
		check(u.Procurement.SupplierSubscription.ListSupplierSubscriptions != nil, "UseCases.Procurement.SupplierSubscription.ListSupplierSubscriptions")
		check(u.Procurement.SupplierSubscription.CreateSupplierSubscription != nil, "UseCases.Procurement.SupplierSubscription.CreateSupplierSubscription")
	}

	if len(missing) > 0 {
		return fmt.Errorf("centymo.Block: incomplete UseCases — missing %v", missing)
	}
	return nil
}

// MustValidate is the FAIL-CLOSED enforcement wrapper around RequireFor. It is
// the seam-level guard that makes a missing REQUIRED closure impossible to
// ignore — mirroring the AUTHZ_ENFORCE boot-guard in service-admin's
// container.go (a missing security precondition is a boot REFUSAL, never a
// silent degrade).
//
// Why a wrapper and not just `return RequireFor(...)`: a bare returned error is
// fail-OPEN by convention. A caller can drop it (`_ =`, an ignored value, a
// future app that doesn't check) and the block silently registers an empty
// feature — the exact nil-closure trap the architecture roast (burn #1) named.
// MustValidate removes that escape hatch:
//
//   - In dev/test (running under `go test`, OR CENTYMO_BLOCK_STRICT truthy) a
//     missing REQUIRED closure PANICS with the full field list. A panic cannot
//     be silently dropped, prints a stack trace at the offending wiring site,
//     and fails the test/CI loudly. This is where a developer wiring a new
//     entity discovers a gap — at their desk, not in prod.
//   - In prod a missing REQUIRED closure logs a screaming FATAL line at the
//     seam (so even a caller that drops the returned error leaves an
//     unmissable log record) AND returns the error so Block() propagates it and
//     NewServiceAdmin halts boot with a clear "domain block failed" message.
//
// OPTIONAL ports are NEVER flagged — that required-vs-optional discrimination
// lives entirely in RequireFor, which only asserts a field when its enabling
// cfg.wantXxx() module is on. MustValidate adds posture, not policy: it changes
// HOW a gap fails, not WHICH fields gate.
func (u *UseCases) MustValidate(cfg *blockConfig) error {
	err := u.RequireFor(cfg)
	if err == nil {
		return nil
	}
	if blockStrictMode() {
		// Dev/test: loud, uncatchable-by-accident, stack-traced.
		panic("FATAL: " + err.Error() + " — REQUIRED block wiring is nil. " +
			"Fix the closure assignment in service-admin's buildCentymoUseCases " +
			"(adapters.go) before this reaches prod.")
	}
	// Prod: scream at the seam, then return so boot halts. The log line is the
	// belt to the returned-error's suspenders (a dropped error still screams).
	log.Printf("FATAL: %v — refusing to register centymo modules with a nil "+
		"REQUIRED closure (fail-closed wiring).", err)
	return err
}

// blockStrictMode reports whether the fail-closed wiring guard should PANIC
// (dev/test) rather than return-and-log (prod) on a missing REQUIRED closure.
//
// True when running under `go test` (testing.Testing(), Go 1.21+ — zero env
// coupling, auto-on in every test + CI run) OR when CENTYMO_BLOCK_STRICT is set
// to an explicit truthy value (the dev escape hatch for `go run` smoke tests).
// The env matching mirrors container.go's authzEnforceEnabled — anything else
// (unset, "", "0", "false") is prod posture.
func blockStrictMode() bool {
	if testing.Testing() {
		return true
	}
	switch os.Getenv("CENTYMO_BLOCK_STRICT") {
	case "1", "true", "TRUE", "True", "yes", "on":
		return true
	default:
		return false
	}
}
