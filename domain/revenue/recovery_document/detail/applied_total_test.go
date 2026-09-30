package detail

import (
	"testing"

	collectionapplicationpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/collection_application"
)

// Only active APPLIED cash applications reduce the balance; an inactive row is ignored.
func TestAppliedTotalCountsOnlyActiveAppliedCash(t *testing.T) {
	applied := collectionapplicationpb.ApplicationStatus_APPLICATION_STATUS_APPLIED
	cash := collectionapplicationpb.ApplicationKind_APPLICATION_KIND_CASH
	got := appliedTotal([]*collectionapplicationpb.CollectionApplication{
		{Amount: 100, Active: true, Status: applied, ApplicationKind: cash},
		{Amount: 40, Active: false, Status: applied, ApplicationKind: cash},
		{Amount: 7, Active: true, Status: collectionapplicationpb.ApplicationStatus_APPLICATION_STATUS_REVERSED, ApplicationKind: cash},
	})
	if got != 100 {
		t.Fatalf("appliedTotal = %d, want 100", got)
	}
}
