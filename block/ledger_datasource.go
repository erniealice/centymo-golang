package block

import (
	"context"
	"time"

	consumer "github.com/erniealice/espyna-golang/consumer"
	expreportpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/reporting/expenditure_report"
	reportpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/reporting/gross_profit"
	agingpb    "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/reporting/receivables_aging"
	payagingpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/reporting/payables_aging"
	revreportpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/reporting/revenue_report"
	collsumpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/reporting/collection_summary"
	disbreportpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/reporting/disbursement_report"
	suppstmtpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/reporting/supplier_statement"
	fycha "github.com/erniealice/fycha-golang"
)

// espynaLedgerDataSource adapts consumer.LedgerReportingService to the
// fycha.DataSource interface. centymo.Block() uses this when the
// AppContext supplies a LedgerReportingSvc so that fycha report views
// (gross profit, revenue report, etc.) can be registered alongside the
// centymo commerce domain routes.
type espynaLedgerDataSource struct {
	svc consumer.LedgerReportingService
}

// Ensure espynaLedgerDataSource satisfies fycha.DataSource at compile time.
var _ fycha.DataSource = (*espynaLedgerDataSource)(nil)

// GetGrossProfitReport delegates to the espyna LedgerReportingAdapter.
func (d *espynaLedgerDataSource) GetGrossProfitReport(
	ctx context.Context,
	req *reportpb.GrossProfitReportRequest,
) (*reportpb.GrossProfitReportResponse, error) {
	return d.svc.GetGrossProfitReport(ctx, req)
}

// GetRevenueReport delegates to the espyna LedgerReportingAdapter.
func (d *espynaLedgerDataSource) GetRevenueReport(
	ctx context.Context,
	req *revreportpb.RevenueReportRequest,
) (*revreportpb.RevenueReportResponse, error) {
	return d.svc.GetRevenueReport(ctx, req)
}

// GetExpenditureReport delegates to the espyna LedgerReportingAdapter.
func (d *espynaLedgerDataSource) GetExpenditureReport(
	ctx context.Context,
	req *expreportpb.ExpenditureReportRequest,
) (*expreportpb.ExpenditureReportResponse, error) {
	return d.svc.GetExpenditureReport(ctx, req)
}

// GetDisbursementReport delegates to the espyna LedgerReportingAdapter.
func (d *espynaLedgerDataSource) GetDisbursementReport(
	ctx context.Context,
	req *disbreportpb.DisbursementReportRequest,
) (*disbreportpb.DisbursementReportResponse, error) {
	return d.svc.GetDisbursementReport(ctx, req)
}

// GetReceivablesAgingReport delegates to the espyna LedgerReportingAdapter.
func (d *espynaLedgerDataSource) GetReceivablesAgingReport(
	ctx context.Context,
	req *agingpb.ReceivablesAgingRequest,
) (*agingpb.ReceivablesAgingResponse, error) {
	return d.svc.GetReceivablesAgingReport(ctx, req)
}

// GetPayablesAgingReport delegates to the espyna LedgerReportingAdapter.
func (d *espynaLedgerDataSource) GetPayablesAgingReport(
	ctx context.Context,
	req *payagingpb.PayablesAgingRequest,
) (*payagingpb.PayablesAgingResponse, error) {
	return d.svc.GetPayablesAgingReport(ctx, req)
}

// GetCollectionSummaryReport delegates to the espyna LedgerReportingAdapter.
func (d *espynaLedgerDataSource) GetCollectionSummaryReport(
	ctx context.Context,
	req *collsumpb.CollectionSummaryRequest,
) (*collsumpb.CollectionSummaryResponse, error) {
	return d.svc.GetCollectionSummaryReport(ctx, req)
}

// GetSupplierStatement delegates to the espyna LedgerReportingAdapter.
func (d *espynaLedgerDataSource) GetSupplierStatement(
	ctx context.Context,
	req *suppstmtpb.SupplierStatementRequest,
) (*suppstmtpb.SupplierStatementResponse, error) {
	return d.svc.GetSupplierStatement(ctx, req)
}

// GetSupplierBalances delegates to the espyna LedgerReportingAdapter.
func (d *espynaLedgerDataSource) GetSupplierBalances(ctx context.Context) (map[string]int64, error) {
	return d.svc.GetSupplierBalances(ctx)
}

// ListRevenue delegates to the espyna LedgerReportingAdapter.
func (d *espynaLedgerDataSource) ListRevenue(
	ctx context.Context,
	start, end *time.Time,
) ([]map[string]any, error) {
	return d.svc.ListRevenue(ctx, start, end)
}

// ListExpenses delegates to the espyna LedgerReportingAdapter.
func (d *espynaLedgerDataSource) ListExpenses(
	ctx context.Context,
	start, end *time.Time,
) ([]map[string]any, error) {
	return d.svc.ListExpenses(ctx, start, end)
}

// newLedgerDataSource wraps a consumer.LedgerReportingService as a fycha.DataSource.
// Returns nil if svc is nil (report views will be skipped gracefully).
func newLedgerDataSource(svc consumer.LedgerReportingService) fycha.DataSource {
	if svc == nil {
		return nil
	}
	return &espynaLedgerDataSource{svc: svc}
}
