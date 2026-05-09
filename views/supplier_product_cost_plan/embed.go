// Package supplier_product_cost_plan is a templates-only module.
//
// SupplierProductCostPlan is the buying-side canonical rate-card record
// (mirror of selling-side ProductPricePlan). Per the buying/selling parity
// audit (docs/plan/20260509-buying-selling-parity-audit/plan.md §5.4 θ),
// all action handlers, page builders, and form-construction logic live in
// views/cost_plan/ — this module hosts only HTML templates.
package supplier_product_cost_plan

import "embed"

//go:embed templates/*
var TemplatesFS embed.FS
