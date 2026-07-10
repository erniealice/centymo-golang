package product_price_plan

import "github.com/erniealice/espyna-golang/consumer/compose"

// Describe returns the compose descriptor for ProductPricePlan — the canonical
// rate-card record (the 4-way join across Plan/PriceSchedule/Product, FK target
// of revenue_line_item, carrying billing_treatment per line).
//
// ProductPricePlan owns NO routes or handlers of its own: its action handlers
// and page builders live in views/price_plan/, and its drawer form is a
// templates-only module (embed.go), registered directly in each app's renderer.
// So this is a DATA-ONLY unit — it carries only the Labels pointer plus its
// lyngua binding, with no Routes, no RouteJSON, no Nav, no Templates, and no
// Mount (compose permits a nil Mount for "a data-only unit that contributes
// routes/labels/nav but registers no handlers itself").
//
// Its single job is to give the engine a mount whose Labels pointer receives
// the product_price_plan.json tier overlay, so sibling units (PricePlanUnit /
// PriceScheduleUnit) can resolve the POST-OVERLAY labels via
// compose.LabelsOf[*Labels](mc, "subscription.product_price_plan") instead of
// the raw DefaultLabels(), which skips the lyngua tier overlay.
func Describe() compose.Unit {
	l := DefaultLabels()
	return compose.Unit{
		Key:       "subscription.product_price_plan",
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "product_price_plan.json", Key: "product_price_plan"},
		LabelName: "ProductPricePlanLabels",
	}
}
