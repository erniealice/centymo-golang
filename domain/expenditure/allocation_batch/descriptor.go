package allocation_batch

import "github.com/erniealice/espyna-golang/consumer/compose"

// Describe returns the compose descriptor for the cost allocation drawers. The
// unit is opt-in (block.WithKnownCostRecovery) and carries no sidebar entry.
func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "expenditure.allocation_batch",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "allocation_batch"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "allocation_batch.json", Key: "allocation_batch"},
		LabelName: "AllocationBatchLabels",
		Templates: TemplatesFS,
	}
}
