package block

import (
	"context"
	"time"

	consumer "github.com/erniealice/espyna-golang/consumer"
	reportpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/reporting/gross_profit"
	revreportpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/reporting/revenue_report"
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
