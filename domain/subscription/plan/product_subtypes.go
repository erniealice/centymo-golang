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
	Product            string                  `json:"product"`
	ProductPlaceholder string                  `json:"productPlaceholder"`
	SelectProduct      string                  `json:"selectProduct"`
	Active             string                  `json:"active"`
	ProductKindLabel   string                  `json:"productKindLabel"`
	ProductKind        ProductKindOptionLabels `json:"productKind"`

	// Model D — variant picker on the ProductPlan drawer form
	VariantSelectLabel       string `json:"variantSelectLabel"`
	VariantSelectPlaceholder string `json:"variantSelectPlaceholder"`
	VariantSelectInfo        string `json:"variantSelectInfo"`
}

// ProductKindOptionLabels provides translated labels for each product_kind
// enum value, used to build the kind selector on the add/edit drawer AND
// to map product_kind values to display labels in table cells.
type ProductKindOptionLabels struct {
	Service        string `json:"service"`
	StockedGood    string `json:"stockedGood"`
	NonStockedGood string `json:"nonStockedGood"`
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
