// Package expenditure is the expenditure-domain consumer facade (centymo restructure).
//
// PURE RE-EXPORT — zero behaviour. The expenditure domain's data/route types,
// Default* constructors, and URL consts moved into per-entity packages under
// domain/expenditure/<entity>/ with entity-local names (the <Entity> prefix stripped).
// This facade re-adds the original prefixed names so existing consumers
// (block/, service-admin) keep resolving expenditure.<Entity>Labels /
// expenditure.Default<Entity>Routes() / expenditure.<Entity>ListURL unchanged.
//
// An entity package MUST NEVER import this facade (that would be an import
// cycle expenditure -> <entity> -> expenditure); cross-entity references go DIRECT to the
// sibling package.
package expenditure

import (
	accruedexpensepkg "github.com/erniealice/centymo-golang/domain/expenditure/accrued_expense"
	expenditurepkg "github.com/erniealice/centymo-golang/domain/expenditure/expenditure"
	expenserecognitionpkg "github.com/erniealice/centymo-golang/domain/expenditure/expense_recognition"
	expenserecognitionrunpkg "github.com/erniealice/centymo-golang/domain/expenditure/expense_recognition_run"
	procurementrequestpkg "github.com/erniealice/centymo-golang/domain/expenditure/procurement_request"
	purchaseorderpkg "github.com/erniealice/centymo-golang/domain/expenditure/purchase_order"
	supplierbillingeventpkg "github.com/erniealice/centymo-golang/domain/expenditure/supplier_billing_event"
	suppliercontractpkg "github.com/erniealice/centymo-golang/domain/expenditure/supplier_contract"
	suppliercontractpriceschedulepkg "github.com/erniealice/centymo-golang/domain/expenditure/supplier_contract_price_schedule"
)

// Re-exported data/route types (type aliases — identity-preserving).
type (
	AccruedExpenseActionLabels                    = accruedexpensepkg.ActionLabels
	AccruedExpenseBalanceLabels                   = accruedexpensepkg.BalanceLabels
	AccruedExpenseButtonLabels                    = accruedexpensepkg.ButtonLabels
	AccruedExpenseColumnLabels                    = accruedexpensepkg.ColumnLabels
	AccruedExpenseConfirmLabels                   = accruedexpensepkg.ConfirmLabels
	AccruedExpenseDetailLabels                    = accruedexpensepkg.DetailLabels
	AccruedExpenseEmptyLabels                     = accruedexpensepkg.EmptyLabels
	AccruedExpenseErrorLabels                     = accruedexpensepkg.ErrorLabels
	AccruedExpenseFormLabels                      = accruedexpensepkg.FormLabels
	AccruedExpenseLabels                          = accruedexpensepkg.Labels
	AccruedExpensePageLabels                      = accruedexpensepkg.PageLabels
	AccruedExpenseRoutes                          = accruedexpensepkg.Routes
	AccruedExpenseSettlementLabels                = accruedexpensepkg.SettlementLabels
	AccruedExpenseStatusLabels                    = accruedexpensepkg.StatusLabels
	AccruedExpenseTabLabels                       = accruedexpensepkg.TabLabels
	ExpenditureActionLabels                       = expenditurepkg.ActionLabels
	ExpenditureBulkLabels                         = expenditurepkg.BulkLabels
	ExpenditureButtonLabels                       = expenditurepkg.ButtonLabels
	ExpenditureCategoryActionLabels               = expenditurepkg.CategoryActionLabels
	ExpenditureCategoryButtonLabels               = expenditurepkg.CategoryButtonLabels
	ExpenditureCategoryColumnLabels               = expenditurepkg.CategoryColumnLabels
	ExpenditureCategoryConfirmLabels              = expenditurepkg.CategoryConfirmLabels
	ExpenditureCategoryEmptyLabels                = expenditurepkg.CategoryEmptyLabels
	ExpenditureCategoryErrorLabels                = expenditurepkg.CategoryErrorLabels
	ExpenditureCategoryFormLabels                 = expenditurepkg.CategoryFormLabels
	ExpenditureCategoryLabels                     = expenditurepkg.CategoryLabels
	ExpenditureCategoryPageLabels                 = expenditurepkg.CategoryPageLabels
	ExpenditureColumnLabels                       = expenditurepkg.ColumnLabels
	ExpenditureDetailLabels                       = expenditurepkg.DetailLabels
	ExpenditureDisbursementCategoryLabels         = expenditurepkg.DisbursementCategoryLabels
	ExpenditureDisbursementFormLabels             = expenditurepkg.DisbursementFormLabels
	ExpenditureEmptyLabels                        = expenditurepkg.EmptyLabels
	ExpenditureErrorLabels                        = expenditurepkg.ErrorLabels
	ExpenditureFormLabels                         = expenditurepkg.FormLabels
	ExpenditureLabelNames                         = expenditurepkg.LabelNames
	ExpenditureLabels                             = expenditurepkg.Labels
	ExpenditureLineItemFormLabels                 = expenditurepkg.LineItemFormLabels
	ExpenditurePageLabels                         = expenditurepkg.PageLabels
	ExpenditurePaymentMethodLabels                = expenditurepkg.PaymentMethodLabels
	ExpenditureRoutes                             = expenditurepkg.Routes
	ExpenditureScheduleLabels                     = expenditurepkg.ScheduleLabels
	ExpenditureStatusLabels                       = expenditurepkg.StatusLabels
	ExpenditureTypeLabels                         = expenditurepkg.TypeLabels
	ExpenseDashboardLabels                        = expenditurepkg.ExpenseDashboardLabels
	ExpenseRecognitionActionLabels                = expenserecognitionpkg.ActionLabels
	ExpenseRecognitionButtonLabels                = expenserecognitionpkg.ButtonLabels
	ExpenseRecognitionColumnLabels                = expenserecognitionpkg.ColumnLabels
	ExpenseRecognitionConfirmLabels               = expenserecognitionpkg.ConfirmLabels
	ExpenseRecognitionDetailLabels                = expenserecognitionpkg.DetailLabels
	ExpenseRecognitionEmptyLabels                 = expenserecognitionpkg.EmptyLabels
	ExpenseRecognitionErrorLabels                 = expenserecognitionpkg.ErrorLabels
	ExpenseRecognitionLabels                      = expenserecognitionpkg.Labels
	ExpenseRecognitionLineLabels                  = expenserecognitionpkg.LineLabels
	ExpenseRecognitionPageLabels                  = expenserecognitionpkg.PageLabels
	ExpenseRecognitionRoutes                      = expenserecognitionpkg.Routes
	ExpenseRecognitionRunActionLabels             = expenserecognitionrunpkg.ActionLabels
	ExpenseRecognitionRunBillsTabLabels           = expenserecognitionrunpkg.BillsTabLabels
	ExpenseRecognitionRunButtonLabels             = expenserecognitionrunpkg.ButtonLabels
	ExpenseRecognitionRunColumnLabels             = expenserecognitionrunpkg.ColumnLabels
	ExpenseRecognitionRunConfirmationLabels       = expenserecognitionrunpkg.ConfirmationLabels
	ExpenseRecognitionRunDetailLabels             = expenserecognitionrunpkg.DetailLabels
	ExpenseRecognitionRunDetailTabHintLabels      = expenserecognitionrunpkg.DetailTabHintLabels
	ExpenseRecognitionRunDetailTabLabels          = expenserecognitionrunpkg.DetailTabLabels
	ExpenseRecognitionRunDrawerLabels             = expenserecognitionrunpkg.DrawerLabels
	ExpenseRecognitionRunEmptyLabels              = expenserecognitionrunpkg.EmptyLabels
	ExpenseRecognitionRunEntityLabels             = expenserecognitionrunpkg.EntityLabels
	ExpenseRecognitionRunErrorLabels              = expenserecognitionrunpkg.ErrorLabels
	ExpenseRecognitionRunFilterLabels             = expenserecognitionrunpkg.FilterLabels
	ExpenseRecognitionRunLabels                   = expenserecognitionrunpkg.Labels
	ExpenseRecognitionRunListColumnLabels         = expenserecognitionrunpkg.ListColumnLabels
	ExpenseRecognitionRunListEmptyLabels          = expenserecognitionrunpkg.ListEmptyLabels
	ExpenseRecognitionRunListEmptyStateLabels     = expenserecognitionrunpkg.ListEmptyStateLabels
	ExpenseRecognitionRunListFilterLabels         = expenserecognitionrunpkg.ListFilterLabels
	ExpenseRecognitionRunListLabels               = expenserecognitionrunpkg.ListLabels
	ExpenseRecognitionRunOutcomeLabels            = expenserecognitionrunpkg.OutcomeLabels
	ExpenseRecognitionRunPageLabels               = expenserecognitionrunpkg.PageLabels
	ExpenseRecognitionRunQueueBulkLabels          = expenserecognitionrunpkg.QueueBulkLabels
	ExpenseRecognitionRunQueueColumnLabels        = expenserecognitionrunpkg.QueueColumnLabels
	ExpenseRecognitionRunQueueEmptyLabels         = expenserecognitionrunpkg.QueueEmptyLabels
	ExpenseRecognitionRunQueueLabels              = expenserecognitionrunpkg.QueueLabels
	ExpenseRecognitionRunRecognitionsTabLabels    = expenserecognitionrunpkg.RecognitionsTabLabels
	ExpenseRecognitionRunResultsTabLabels         = expenserecognitionrunpkg.ResultsTabLabels
	ExpenseRecognitionRunRoutes                   = expenserecognitionrunpkg.Routes
	ExpenseRecognitionRunScopeKindLabels          = expenserecognitionrunpkg.ScopeKindLabels
	ExpenseRecognitionRunSearchLabels             = expenserecognitionrunpkg.SearchLabels
	ExpenseRecognitionRunSelectionsTabLabels      = expenserecognitionrunpkg.SelectionsTabLabels
	ExpenseRecognitionRunSourceKindLabels         = expenserecognitionrunpkg.SourceKindLabels
	ExpenseRecognitionRunStatusBadgeLabels        = expenserecognitionrunpkg.StatusBadgeLabels
	ExpenseRecognitionRunSubscriptionDrawerLabels = expenserecognitionrunpkg.SubscriptionDrawerLabels
	ExpenseRecognitionRunSummaryLabels            = expenserecognitionrunpkg.SummaryLabels
	ExpenseRecognitionRunSupplierDrawerLabels     = expenserecognitionrunpkg.SupplierDrawerLabels
	ExpenseRecognitionRunSuppressionLabels        = expenserecognitionrunpkg.SuppressionLabels
	ExpenseRecognitionRunToastLabels              = expenserecognitionrunpkg.ToastLabels
	ExpenseRecognitionSourceLabels                = expenserecognitionpkg.SourceLabels
	ExpenseRecognitionStatusLabels                = expenserecognitionpkg.StatusLabels
	ExpenseRecognitionTabLabels                   = expenserecognitionpkg.TabLabels
	ProcurementRequestColumnLabels                = procurementrequestpkg.ColumnLabels
	ProcurementRequestDetailLabels                = procurementrequestpkg.DetailLabels
	ProcurementRequestEmptyLabels                 = procurementrequestpkg.EmptyLabels
	ProcurementRequestFilterLabels                = procurementrequestpkg.FilterLabels
	ProcurementRequestFormLabels                  = procurementrequestpkg.FormLabels
	ProcurementRequestFulfillmentModeHintLabels   = procurementrequestpkg.FulfillmentModeHintLabels
	ProcurementRequestFulfillmentModeLabels       = procurementrequestpkg.FulfillmentModeLabels
	ProcurementRequestFulfillmentStrategyLabels   = procurementrequestpkg.FulfillmentStrategyLabels
	ProcurementRequestLabels                      = procurementrequestpkg.Labels
	ProcurementRequestLineLabels                  = procurementrequestpkg.LineLabels
	ProcurementRequestPageLabels                  = procurementrequestpkg.PageLabels
	ProcurementRequestPolicyDecisionLabels        = procurementrequestpkg.PolicyDecisionLabels
	ProcurementRequestRoutes                      = procurementrequestpkg.Routes
	ProcurementRequestSpawnLabels                 = procurementrequestpkg.SpawnLabels
	ProcurementRequestSpawnedPOLabels             = procurementrequestpkg.SpawnedPOLabels
	ProcurementRequestTabLabels                   = procurementrequestpkg.TabLabels
	PurchaseDashboardLabels                       = expenditurepkg.PurchaseDashboardLabels
	PurchaseOrderActionLabels                     = purchaseorderpkg.ActionLabels
	PurchaseOrderBulkLabels                       = purchaseorderpkg.BulkLabels
	PurchaseOrderButtonLabels                     = purchaseorderpkg.ButtonLabels
	PurchaseOrderColumnLabels                     = purchaseorderpkg.ColumnLabels
	PurchaseOrderDetailLabels                     = purchaseorderpkg.DetailLabels
	PurchaseOrderEmptyLabels                      = purchaseorderpkg.EmptyLabels
	PurchaseOrderErrorLabels                      = purchaseorderpkg.ErrorLabels
	PurchaseOrderFormLabels                       = purchaseorderpkg.FormLabels
	PurchaseOrderLabelNames                       = purchaseorderpkg.LabelNames
	PurchaseOrderLabels                           = purchaseorderpkg.Labels
	PurchaseOrderLineItemLabels                   = purchaseorderpkg.LineItemLabels
	PurchaseOrderLineTypeLabels                   = purchaseorderpkg.LineTypeLabels
	PurchaseOrderPOTypeLabels                     = purchaseorderpkg.POTypeLabels
	PurchaseOrderPageLabels                       = purchaseorderpkg.PageLabels
	PurchaseOrderReceiptLabels                    = purchaseorderpkg.ReceiptLabels
	PurchaseOrderStatusLabels                     = purchaseorderpkg.StatusLabels
	SupplierBillingEventActionLabels              = supplierbillingeventpkg.ActionLabels
	SupplierBillingEventColumnLabels              = supplierbillingeventpkg.ColumnLabels
	SupplierBillingEventDetailLabels              = supplierbillingeventpkg.DetailLabels
	SupplierBillingEventEmptyLabels               = supplierbillingeventpkg.EmptyLabels
	SupplierBillingEventErrorLabels               = supplierbillingeventpkg.ErrorLabels
	SupplierBillingEventLabels                    = supplierbillingeventpkg.Labels
	SupplierBillingEventPageLabels                = supplierbillingeventpkg.PageLabels
	SupplierBillingEventStatusLabels              = supplierbillingeventpkg.StatusLabels
	SupplierBillingEventTriggerLabels             = supplierbillingeventpkg.TriggerLabels
	SupplierContractColumnLabels                  = suppliercontractpkg.ColumnLabels
	SupplierContractDetailLabels                  = suppliercontractpkg.DetailLabels
	SupplierContractEmptyLabels                   = suppliercontractpkg.EmptyLabels
	SupplierContractFormLabels                    = suppliercontractpkg.FormLabels
	SupplierContractLabels                        = suppliercontractpkg.Labels
	SupplierContractLineLabels                    = suppliercontractpkg.LineLabels
	SupplierContractLinkedExpenditureLabels       = suppliercontractpkg.LinkedExpenditureLabels
	SupplierContractLinkedPOLabels                = suppliercontractpkg.LinkedPOLabels
	SupplierContractPageLabels                    = suppliercontractpkg.PageLabels
	SupplierContractPriceScheduleButtonLabels     = suppliercontractpriceschedulepkg.ButtonLabels
	SupplierContractPriceScheduleColumnLabels     = suppliercontractpriceschedulepkg.ColumnLabels
	SupplierContractPriceScheduleDetailLabels     = suppliercontractpriceschedulepkg.DetailLabels
	SupplierContractPriceScheduleEmptyLabels      = suppliercontractpriceschedulepkg.EmptyLabels
	SupplierContractPriceScheduleErrorLabels      = suppliercontractpriceschedulepkg.ErrorLabels
	SupplierContractPriceScheduleFilterLabels     = suppliercontractpriceschedulepkg.FilterLabels
	SupplierContractPriceScheduleFormLabels       = suppliercontractpriceschedulepkg.FormLabels
	SupplierContractPriceScheduleLabels           = suppliercontractpriceschedulepkg.Labels
	SupplierContractPriceScheduleLineFormLabels   = suppliercontractpriceschedulepkg.LineFormLabels
	SupplierContractPriceScheduleLinesLabels      = suppliercontractpriceschedulepkg.LinesLabels
	SupplierContractPriceScheduleNounLabels       = suppliercontractpriceschedulepkg.NounLabels
	SupplierContractPriceSchedulePageLabels       = suppliercontractpriceschedulepkg.PageLabels
	SupplierContractPriceScheduleRoutes           = suppliercontractpriceschedulepkg.Routes
	SupplierContractPriceScheduleStatusLabels     = suppliercontractpriceschedulepkg.StatusLabels
	SupplierContractPriceScheduleTabLabels        = suppliercontractpriceschedulepkg.TabLabels
	SupplierContractRoutes                        = suppliercontractpkg.Routes
	SupplierContractTabLabels                     = suppliercontractpkg.TabLabels
)

// Re-exported URL route consts (const-identity preserved).
const (
	AccruedExpenseAccrueFromContractURL              = accruedexpensepkg.AccrueFromContractURL
	AccruedExpenseAddURL                             = accruedexpensepkg.AddURL
	AccruedExpenseAttachmentDeleteURL                = accruedexpensepkg.AttachmentDeleteURL
	AccruedExpenseAttachmentUploadURL                = accruedexpensepkg.AttachmentUploadURL
	AccruedExpenseBulkSetStatusURL                   = accruedexpensepkg.BulkSetStatusURL
	AccruedExpenseDeleteURL                          = accruedexpensepkg.DeleteURL
	AccruedExpenseDetailURL                          = accruedexpensepkg.DetailURL
	AccruedExpenseEditURL                            = accruedexpensepkg.EditURL
	AccruedExpenseListURL                            = accruedexpensepkg.ListURL
	AccruedExpenseReverseURL                         = accruedexpensepkg.ReverseURL
	AccruedExpenseSetStatusURL                       = accruedexpensepkg.SetStatusURL
	AccruedExpenseSettleURL                          = accruedexpensepkg.SettleURL
	AccruedExpenseSettlementAddURL                   = accruedexpensepkg.SettlementAddURL
	AccruedExpenseSettlementDeleteURL                = accruedexpensepkg.SettlementDeleteURL
	AccruedExpenseSettlementEditURL                  = accruedexpensepkg.SettlementEditURL
	AccruedExpenseTabActionURL                       = accruedexpensepkg.TabActionURL
	ExpenditureAttachmentDeleteURL                   = expenditurepkg.AttachmentDeleteURL
	ExpenditureAttachmentUploadURL                   = expenditurepkg.AttachmentUploadURL
	ExpenditureExpenseAddURL                         = expenditurepkg.ExpenseAddURL
	ExpenditureExpenseCategoryAddURL                 = expenditurepkg.ExpenseCategoryAddURL
	ExpenditureExpenseCategoryDeleteURL              = expenditurepkg.ExpenseCategoryDeleteURL
	ExpenditureExpenseCategoryEditURL                = expenditurepkg.ExpenseCategoryEditURL
	ExpenditureExpenseCategoryListURL                = expenditurepkg.ExpenseCategoryListURL
	ExpenditureExpenseCategoryTableURL               = expenditurepkg.ExpenseCategoryTableURL
	ExpenditureExpenseDashboardURL                   = expenditurepkg.ExpenseDashboardURL
	ExpenditureExpenseDeleteURL                      = expenditurepkg.ExpenseDeleteURL
	ExpenditureExpenseDetailURL                      = expenditurepkg.ExpenseDetailURL
	ExpenditureExpenseEditURL                        = expenditurepkg.ExpenseEditURL
	ExpenditureExpenseLineItemAddURL                 = expenditurepkg.ExpenseLineItemAddURL
	ExpenditureExpenseLineItemEditURL                = expenditurepkg.ExpenseLineItemEditURL
	ExpenditureExpenseLineItemRemoveURL              = expenditurepkg.ExpenseLineItemRemoveURL
	ExpenditureExpenseLineItemTableURL               = expenditurepkg.ExpenseLineItemTableURL
	ExpenditureExpenseListURL                        = expenditurepkg.ExpenseListURL
	ExpenditureExpensePayURL                         = expenditurepkg.ExpensePayURL
	ExpenditureExpenseSetStatusURL                   = expenditurepkg.ExpenseSetStatusURL
	ExpenditureExpenseTabActionURL                   = expenditurepkg.ExpenseTabActionURL
	ExpenditureExpenseTableURL                       = expenditurepkg.ExpenseTableURL
	ExpenditurePurchaseDashboardURL                  = expenditurepkg.PurchaseDashboardURL
	ExpenditurePurchaseListURL                       = expenditurepkg.PurchaseListURL
	ExpenditureSettingsTemplateDefaultURL            = expenditurepkg.SettingsTemplateDefaultURL
	ExpenditureSettingsTemplateDeleteURL             = expenditurepkg.SettingsTemplateDeleteURL
	ExpenditureSettingsTemplateUploadURL             = expenditurepkg.SettingsTemplateUploadURL
	ExpenditureSettingsTemplatesURL                  = expenditurepkg.SettingsTemplatesURL
	ExpenseRecognitionAttachmentDeleteURL            = expenserecognitionpkg.AttachmentDeleteURL
	ExpenseRecognitionAttachmentUploadURL            = expenserecognitionpkg.AttachmentUploadURL
	ExpenseRecognitionDeleteURL                      = expenserecognitionpkg.DeleteURL
	ExpenseRecognitionDetailURL                      = expenserecognitionpkg.DetailURL
	ExpenseRecognitionLineAddURL                     = expenserecognitionpkg.LineAddURL
	ExpenseRecognitionLineDeleteURL                  = expenserecognitionpkg.LineDeleteURL
	ExpenseRecognitionLineEditURL                    = expenserecognitionpkg.LineEditURL
	ExpenseRecognitionListURL                        = expenserecognitionpkg.ListURL
	ExpenseRecognitionRecognizeFromContractURL       = expenserecognitionpkg.RecognizeFromContractURL
	ExpenseRecognitionRecognizeFromExpenditureURL    = expenserecognitionpkg.RecognizeFromExpenditureURL
	ExpenseRecognitionReverseURL                     = expenserecognitionpkg.ReverseURL
	ExpenseRecognitionRunDetailTabActionURL          = expenserecognitionrunpkg.DetailTabActionURL
	ExpenseRecognitionRunDetailURL                   = expenserecognitionrunpkg.DetailURL
	ExpenseRecognitionRunGenerateURL                 = expenserecognitionrunpkg.GenerateURL
	ExpenseRecognitionRunListTableURL                = expenserecognitionrunpkg.ListTableURL
	ExpenseRecognitionRunListURL                     = expenserecognitionrunpkg.ListURL
	ExpenseRecognitionRunNewURL                      = expenserecognitionrunpkg.NewURL
	ExpenseRecognitionRunPerSubscriptionDrawerURL    = expenserecognitionrunpkg.PerSubscriptionDrawerURL
	ExpenseRecognitionRunPerSupplierDrawerURL        = expenserecognitionrunpkg.PerSupplierDrawerURL
	ExpenseRecognitionRunQueueTableURL               = expenserecognitionrunpkg.QueueTableURL
	ExpenseRecognitionRunQueueURL                    = expenserecognitionrunpkg.QueueURL
	ExpenseRecognitionRunSubmitBatchURL              = expenserecognitionrunpkg.SubmitBatchURL
	ExpenseRecognitionTabActionURL                   = expenserecognitionpkg.TabActionURL
	ExpensesSummaryURL                               = expenditurepkg.ExpensesSummaryURL
	ProcurementRequestAddURL                         = procurementrequestpkg.AddURL
	ProcurementRequestApproveURL                     = procurementrequestpkg.ApproveURL
	ProcurementRequestAttachmentDeleteURL            = procurementrequestpkg.AttachmentDeleteURL
	ProcurementRequestAttachmentUploadURL            = procurementrequestpkg.AttachmentUploadURL
	ProcurementRequestBulkSetStatusURL               = procurementrequestpkg.BulkSetStatusURL
	ProcurementRequestDeleteURL                      = procurementrequestpkg.DeleteURL
	ProcurementRequestDetailURL                      = procurementrequestpkg.DetailURL
	ProcurementRequestEditURL                        = procurementrequestpkg.EditURL
	ProcurementRequestLineAddURL                     = procurementrequestpkg.LineAddURL
	ProcurementRequestLineDeleteURL                  = procurementrequestpkg.LineDeleteURL
	ProcurementRequestLineEditURL                    = procurementrequestpkg.LineEditURL
	ProcurementRequestLineRetrySpawnURL              = procurementrequestpkg.LineRetrySpawnURL
	ProcurementRequestListURL                        = procurementrequestpkg.ListURL
	ProcurementRequestRejectURL                      = procurementrequestpkg.RejectURL
	ProcurementRequestSetStatusURL                   = procurementrequestpkg.SetStatusURL
	ProcurementRequestSpawnPOURL                     = procurementrequestpkg.SpawnPOURL
	ProcurementRequestSubmitURL                      = procurementrequestpkg.SubmitURL
	ProcurementRequestTabActionURL                   = procurementrequestpkg.TabActionURL
	PurchaseOrderAddURL                              = expenditurepkg.PurchaseOrderAddURL
	PurchaseOrderAttachmentDeleteURL                 = expenditurepkg.PurchaseOrderAttachmentDeleteURL
	PurchaseOrderAttachmentUploadURL                 = expenditurepkg.PurchaseOrderAttachmentUploadURL
	PurchaseOrderConfirmReceiptURL                   = expenditurepkg.PurchaseOrderConfirmReceiptURL
	PurchaseOrderDeleteURL                           = expenditurepkg.PurchaseOrderDeleteURL
	PurchaseOrderDetailURL                           = expenditurepkg.PurchaseOrderDetailURL
	PurchaseOrderEditURL                             = expenditurepkg.PurchaseOrderEditURL
	PurchaseOrderLineItemAddURL                      = expenditurepkg.PurchaseOrderLineItemAddURL
	PurchaseOrderLineItemEditURL                     = expenditurepkg.PurchaseOrderLineItemEditURL
	PurchaseOrderLineItemRemoveURL                   = expenditurepkg.PurchaseOrderLineItemRemoveURL
	PurchaseOrderLineItemTableURL                    = expenditurepkg.PurchaseOrderLineItemTableURL
	PurchaseOrderListURL                             = expenditurepkg.PurchaseOrderListURL
	PurchaseOrderSetStatusURL                        = expenditurepkg.PurchaseOrderSetStatusURL
	PurchaseOrderTabActionURL                        = expenditurepkg.PurchaseOrderTabActionURL
	PurchaseOrderTableURL                            = expenditurepkg.PurchaseOrderTableURL
	PurchasesSummaryURL                              = expenditurepkg.PurchasesSummaryURL
	SupplierBillingEventDetailURL                    = supplierbillingeventpkg.DetailURL
	SupplierBillingEventListURL                      = supplierbillingeventpkg.ListURL
	SupplierBillingEventRecognizeURL                 = supplierbillingeventpkg.RecognizeURL
	SupplierContractAddURL                           = suppliercontractpkg.AddURL
	SupplierContractApproveURL                       = suppliercontractpkg.ApproveURL
	SupplierContractAttachmentDeleteURL              = suppliercontractpkg.AttachmentDeleteURL
	SupplierContractAttachmentUploadURL              = suppliercontractpkg.AttachmentUploadURL
	SupplierContractBulkSetStatusURL                 = suppliercontractpkg.BulkSetStatusURL
	SupplierContractDeleteURL                        = suppliercontractpkg.DeleteURL
	SupplierContractDetailURL                        = suppliercontractpkg.DetailURL
	SupplierContractEditURL                          = suppliercontractpkg.EditURL
	SupplierContractLineAddURL                       = suppliercontractpkg.LineAddURL
	SupplierContractLineDeleteURL                    = suppliercontractpkg.LineDeleteURL
	SupplierContractLineEditURL                      = suppliercontractpkg.LineEditURL
	SupplierContractListURL                          = suppliercontractpkg.ListURL
	SupplierContractPriceScheduleActivateURL         = suppliercontractpriceschedulepkg.ActivateURL
	SupplierContractPriceScheduleAddURL              = suppliercontractpriceschedulepkg.AddURL
	SupplierContractPriceScheduleAttachmentDeleteURL = suppliercontractpriceschedulepkg.AttachmentDeleteURL
	SupplierContractPriceScheduleAttachmentUploadURL = suppliercontractpriceschedulepkg.AttachmentUploadURL
	SupplierContractPriceScheduleBulkSetStatusURL    = suppliercontractpriceschedulepkg.BulkSetStatusURL
	SupplierContractPriceScheduleDeleteURL           = suppliercontractpriceschedulepkg.DeleteURL
	SupplierContractPriceScheduleDetailURL           = suppliercontractpriceschedulepkg.DetailURL
	SupplierContractPriceScheduleEditURL             = suppliercontractpriceschedulepkg.EditURL
	SupplierContractPriceScheduleLineAddURL          = suppliercontractpriceschedulepkg.LineAddURL
	SupplierContractPriceScheduleLineDeleteURL       = suppliercontractpriceschedulepkg.LineDeleteURL
	SupplierContractPriceScheduleLineEditURL         = suppliercontractpriceschedulepkg.LineEditURL
	SupplierContractPriceScheduleListURL             = suppliercontractpriceschedulepkg.ListURL
	SupplierContractPriceScheduleSetStatusURL        = suppliercontractpriceschedulepkg.SetStatusURL
	SupplierContractPriceScheduleSupersedeURL        = suppliercontractpriceschedulepkg.SupersedeURL
	SupplierContractPriceScheduleTabActionURL        = suppliercontractpriceschedulepkg.TabActionURL
	SupplierContractSetStatusURL                     = suppliercontractpkg.SetStatusURL
	SupplierContractTabActionURL                     = suppliercontractpkg.TabActionURL
	SupplierContractTerminateURL                     = suppliercontractpkg.TerminateURL
)

// Re-exported Default* constructors (function values).
var (
	DefaultAccruedExpenseLabels                = accruedexpensepkg.DefaultLabels
	DefaultAccruedExpenseRoutes                = accruedexpensepkg.DefaultRoutes
	DefaultExpenditureRoutes                   = expenditurepkg.DefaultRoutes
	DefaultExpenseRecognitionLabels            = expenserecognitionpkg.DefaultLabels
	DefaultExpenseRecognitionRoutes            = expenserecognitionpkg.DefaultRoutes
	DefaultExpenseRecognitionRunLabels         = expenserecognitionrunpkg.DefaultLabels
	DefaultExpenseRecognitionRunRoutes         = expenserecognitionrunpkg.DefaultRoutes
	DefaultProcurementRequestLabels            = procurementrequestpkg.DefaultLabels
	DefaultProcurementRequestRoutes            = procurementrequestpkg.DefaultRoutes
	DefaultSupplierBillingEventLabels          = supplierbillingeventpkg.DefaultLabels
	DefaultSupplierContractLabels              = suppliercontractpkg.DefaultLabels
	DefaultSupplierContractPriceScheduleLabels = suppliercontractpriceschedulepkg.DefaultLabels
	DefaultSupplierContractPriceScheduleRoutes = suppliercontractpriceschedulepkg.DefaultRoutes
	DefaultSupplierContractRoutes              = suppliercontractpkg.DefaultRoutes
)
