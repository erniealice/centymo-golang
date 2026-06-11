// Package product_price_plan is a templates-only module.
//
// ProductPricePlan is the canonical rate-card record (4-way join across
// Plan/PriceSchedule/Product trees, FK target of revenue_line_item, carries
// billing_treatment per line). Per the buying/selling parity audit
// (docs/plan/20260509-buying-selling-parity-audit/plan.md §5.4 θ), the
// drawer-form template lives in its own module so the buying-side mirror
// (views/supplier_product_cost_plan/) can match shape.
//
// All action handlers, page builders, and form-construction logic remain in
// views/price_plan/ — this module hosts only HTML templates.
package product_price_plan

import "embed"

//go:embed templates/*
var TemplatesFS embed.FS
