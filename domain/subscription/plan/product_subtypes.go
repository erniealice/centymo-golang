package plan

// plan_product_subtypes.go — Plan-aggregate sub-types whose NAMES collide with
// the product domain under the placement test's mechanical longest-match
// (ProductPlanForm -> product_plan, ProductKindOption -> product), but which are
// functionally owned by the subscription/plan aggregate (nested in PlanLabels,
// built by DefaultPlanLabels). They stay in package subscription for cohesion
// (moving them to domain/product would force a subscription->product sibling
// import). Excused in placement_test legacyAllow by basename pending a W9
// naming-resolution pass. centymo W4. Pure structural relocation, no rename.

// ProductPlanFormLabels holds translatable labels for the ProductPlan add/edit form within a plan.
type ProductPlanFormLabels struct {
	// Name is the offering-name field label (D-R2-1). Generic in code; tiers
	// override the vocabulary via lyngua (e.g. "Name"). Empty user input still
	// falls back to the product name in the action handler.
	Name               string                  `json:"name"`
	NamePlaceholder    string                  `json:"name_placeholder"`
	Product            string                  `json:"product"`
	ProductPlaceholder string                  `json:"product_placeholder"`
	SelectProduct      string                  `json:"select_product"`
	Active             string                  `json:"active"`
	ProductKindLabel   string                  `json:"product_kind_label"`
	ProductKind        ProductKindOptionLabels `json:"product_kind"`

	// Model D — variant picker on the ProductPlan drawer form
	VariantSelectLabel       string `json:"variant_select_label"`
	VariantSelectPlaceholder string `json:"variant_select_placeholder"`
	VariantSelectInfo        string `json:"variant_select_info"`
}

// ProductKindOptionLabels provides translated labels for each product_kind
// enum value, used to build the kind selector on the add/edit drawer AND
// to map product_kind values to display labels in table cells.
type ProductKindOptionLabels struct {
	Service        string `json:"service"`
	StockedGood    string `json:"stocked_good"`
	NonStockedGood string `json:"non_stocked_good"`
	Consumable     string `json:"consumable"`
}

// Label returns the translated label for a product_kind value
// ("service" | "stocked_good" | "non_stocked_good" | "consumable").
// Unknown values round-trip through as-is so callers always get a string.
func (k ProductKindOptionLabels) Label(kind string) string {
	switch kind {
	case "service":
		return k.Service
	case "stocked_good":
		return k.StockedGood
	case "non_stocked_good":
		return k.NonStockedGood
	case "consumable":
		return k.Consumable
	}
	return kind
}
