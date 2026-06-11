package centymo

// labels.go — RESIDUAL after centymo W7.
//
// W1-W7 dissolved this god-file per domain. What remains:
//   - LocationMap / LocationDisplayName (entydad-bound; WL deferral, not yet landed)
//
// Procurement-domain label sections (SupplierSubscription, CostSchedule,
// SupplierPlan, CostPlan, SupplierProductPlan, SupplierProductCostPlan)
// moved to domain/procurement/<entity>_labels.go in W7.
// Expenditure-domain labels moved to domain/expenditure/<entity>_labels.go in W6.
// ProcurementLabels (composition app) moved to domain/procurement/procurement.go in W7.

var LocationMap = map[string]string{
	"ayala-central-bloc": "Ayala Central Bloc",
	"sm-city-cebu":       "SM City Cebu",
	"ayala-center-cebu":  "Ayala Center Cebu",
	"robinsons-galleria": "Robinsons Galleria",
}

// LocationDisplayName returns the display name for a location slug.
func LocationDisplayName(slug string) string {
	if name, ok := LocationMap[slug]; ok {
		return name
	}
	return slug
}
