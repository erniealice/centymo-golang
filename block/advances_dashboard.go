// Package block — advances-dashboard domain wiring (20260517-advance-cash-events Plan B Phase 3).
//
// Holds wireAdvancesDashboardModule (the lifted body of the `if
// cfg.wantTreasuryAdvances()` branch of Block()) plus the proto<->view
// translators that only the advances-dashboard wiring calls. Mirrors the
// revenue_run.go file shape exactly so a reader following one wiring file
// can follow the other.
package block

import (
	"context"
	"log"

	consumerapp "github.com/erniealice/espyna-golang/consumer/app"
	"github.com/erniealice/pyeza-golang/types"

	treasurydomain "github.com/erniealice/centymo-golang/domain/treasury"
	advancesdashboardmod "github.com/erniealice/centymo-golang/domain/treasury/treasuryadvancesdashboard"
)

// advancesDashboardWiring holds everything wireAdvancesDashboardModule needs
// from the surrounding Block() scope. Kept private; never re-exported.
type advancesDashboardWiring struct {
	routes             treasurydomain.TreasuryAdvancesRoutes
	labels             treasurydomain.AdvancesDashboardLabels
	enumLabels         treasurydomain.AdvanceEnumLabels
	collectionRoutes   treasurydomain.CollectionRoutes
	disbursementRoutes treasurydomain.DisbursementRoutes
	centymoTableLabels types.TableLabels
	functionalCurrency func(ctx context.Context) string
}

// wireAdvancesDashboardModule lifts the body of `if cfg.wantTreasuryAdvances()`
// from Block(). Behaviour-preserving: same construction order, same registration
// order, same callbacks. block.go calls this exactly once when
// cfg.wantTreasuryAdvances().
//
// The dashboard is fed by `useCases.TreasuryAdvances.GetAdvancesDashboard` — a
// view-typed closure the service-admin adapter wires from espyna's
// ListAdvanceCollectionsByWorkspace / ListAdvanceDisbursementsByWorkspace
// repositories. The closure can be nil; the view degrades to empty state.
func wireAdvancesDashboardModule(ctx *consumerapp.AppContext, cfg *blockConfig, useCases *UseCases, w advancesDashboardWiring) {
	deps := &advancesdashboardmod.ModuleDeps{
		Routes:                        w.routes,
		Labels:                        w.labels,
		EnumLabels:                    w.enumLabels,
		CommonLabels:                  ctx.Common,
		TableLabels:                   w.centymoTableLabels,
		CollectionDetailURLTemplate:   w.collectionRoutes.DetailURL,
		DisbursementDetailURLTemplate: w.disbursementRoutes.DetailURL,
		GetFunctionalCurrency:         w.functionalCurrency,
	}

	// Wire GetDashboard — translate the view-typed AdvancesDashboardData (from
	// the useCases struct) into the in-view advancesdashboardmod.DashboardResponse.
	// Nil-safe: the view renders empty state when the closure is unset.
	if useCases.TreasuryAdvances.GetAdvancesDashboard != nil {
		deps.GetDashboard = func(fctx context.Context, req advancesdashboardmod.DashboardRequest) (*advancesdashboardmod.DashboardResponse, error) {
			data, err := useCases.TreasuryAdvances.GetAdvancesDashboard(fctx, req.AsOfDate)
			if err != nil {
				log.Printf("advances_dashboard: GetAdvancesDashboard failed: %v", err)
				return nil, err
			}
			if data == nil {
				return &advancesdashboardmod.DashboardResponse{}, nil
			}
			return &advancesdashboardmod.DashboardResponse{
				Outflows: advancesdashboardmod.AdvancesSection{
					Rows: convertAdvancesRows(data.Outflows, w.disbursementRoutes.DetailURL),
				},
				Inflows: advancesdashboardmod.AdvancesSection{
					Rows: convertAdvancesRows(data.Inflows, w.collectionRoutes.DetailURL),
				},
				Position: advancesdashboardmod.AdvancesPosition{
					OutflowTotalRemaining:  data.OutflowTotalRemaining,
					InflowTotalRemaining:   data.InflowTotalRemaining,
					OutflowActiveCount:     data.OutflowActiveCount,
					InflowActiveCount:      data.InflowActiveCount,
					OutflowFullyRecognized: data.OutflowFullyRecognized,
					InflowFullyRecognized:  data.InflowFullyRecognized,
					Currency:               data.Currency,
				},
			}, nil
		}
	}

	module := advancesdashboardmod.NewModule(deps)
	module.RegisterRoutes(ctx.Routes)
}

// convertAdvancesRows translates the view-typed AdvancesDashboardRow values
// from UseCases into the in-view AdvanceRow values the dashboard page renders.
// Computes per-row utilization and deep-link DetailURL.
func convertAdvancesRows(in []AdvancesDashboardRow, detailURLPattern string) []advancesdashboardmod.AdvanceRow {
	out := make([]advancesdashboardmod.AdvanceRow, 0, len(in))
	for _, r := range in {
		pct := 0
		if r.TotalAmount > 0 {
			pct = int((r.RecognizedAmount * 100) / r.TotalAmount)
			if pct < 0 {
				pct = 0
			}
			if pct > 100 {
				pct = 100
			}
		}
		detailURL := ""
		if detailURLPattern != "" && r.ID != "" {
			detailURL = resolveDashboardDetailURL(detailURLPattern, r.ID)
		}
		out = append(out, advancesdashboardmod.AdvanceRow{
			ID:               r.ID,
			ReferenceNumber:  r.ReferenceNumber,
			CounterpartyName: r.CounterpartyName,
			Kind:             r.Kind, // raw → view re-maps to a label downstream
			KindRaw:          r.Kind,
			Status:           r.Status,
			StatusRaw:        r.Status,
			Currency:         r.Currency,
			TotalAmount:      r.TotalAmount,
			RemainingAmount:  r.RemainingAmount,
			RecognizedAmount: r.RecognizedAmount,
			UtilizationPct:   pct,
			DetailURL:        detailURL,
		})
	}
	return out
}

// resolveDashboardDetailURL is a tiny wrapper around the {id} substitution
// pattern used elsewhere in the block layer. Kept private so the dashboard
// wiring doesn't import pyeza/route directly.
func resolveDashboardDetailURL(pattern, id string) string {
	// Stable, dependency-free substitution: replace literal "{id}" once.
	// We avoid pyeza/route here because the URL is just one templated segment
	// and we don't want to import that subpackage in this leaf file.
	idx := stringIndexOf(pattern, "{id}")
	if idx < 0 {
		return pattern
	}
	return pattern[:idx] + id + pattern[idx+len("{id}"):]
}

// stringIndexOf is a tiny non-allocating substring search.
func stringIndexOf(haystack, needle string) int {
	hn := len(haystack)
	nn := len(needle)
	if nn == 0 || hn < nn {
		return -1
	}
	for i := 0; i+nn <= hn; i++ {
		if haystack[i:i+nn] == needle {
			return i
		}
	}
	return -1
}
