package subscription_group_product_plan

import (
	"sort"

	jobtemplatephasepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/operation/job_template_phase"
)

// CoveredPhase is one phase anchor a class can assign staff against (plan.md
// §2.5 the coverage rule). ID/Name/Order come straight from the
// job_template_phase row; VariantID is its output_product_variant_id (used to
// filter a strand class's anchors down to its own semester).
type CoveredPhase struct {
	ID    string
	Name  string
	Order int32
}

// ResolveCoveredPhases filters+sorts a template's active phases down to a
// class's own covered grain (plan.md §2.5): a strand class (offeringVariantID
// set) sees only phases whose output_product_variant_id matches its own
// variant; an umbrella class (offeringVariantID empty) sees every phase — the
// whole duration. Pure function, no I/O — callers batch-fetch `phases` once
// (espyna.md §1b call #4) and pass the class's own product_plan.product_variant_id.
func ResolveCoveredPhases(phases []*jobtemplatephasepb.JobTemplatePhase, offeringVariantID string) []CoveredPhase {
	out := make([]CoveredPhase, 0, len(phases))
	for _, p := range phases {
		if p == nil || !p.GetActive() {
			continue
		}
		if offeringVariantID != "" && p.GetOutputProductVariantId() != "" && p.GetOutputProductVariantId() != offeringVariantID {
			continue
		}
		out = append(out, CoveredPhase{ID: p.GetId(), Name: p.GetName(), Order: p.GetPhaseOrder()})
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Order != out[j].Order {
			return out[i].Order < out[j].Order
		}
		return out[i].ID < out[j].ID
	})
	return out
}
