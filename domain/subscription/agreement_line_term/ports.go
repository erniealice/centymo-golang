package agreement_line_term

import (
	"context"

	policypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/charge_policy"
	versionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/charge_policy_version"
	agreementlinetermpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/agreement_line_term"
	productpriceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/product_price_plan"
)

// UseCases is the set of espyna closures the subscription "Charge terms" tab
// consumes, exactly the proto request/response shape of the existing use cases.
// ListAgreementLineTerms is required (the block's Mount refuses a nil one); the
// other three only enrich names and degrade to ids when unbound.
type UseCases struct {
	ListAgreementLineTerms func(context.Context, *agreementlinetermpb.ListAgreementLineTermsRequest) (*agreementlinetermpb.ListAgreementLineTermsResponse, error)

	ReadChargePolicyVersion func(context.Context, *versionpb.ReadChargePolicyVersionRequest) (*versionpb.ReadChargePolicyVersionResponse, error)
	ReadChargePolicy        func(context.Context, *policypb.ReadChargePolicyRequest) (*policypb.ReadChargePolicyResponse, error)
	ListProductPricePlans   func(context.Context, *productpriceplanpb.ListProductPricePlansRequest) (*productpriceplanpb.ListProductPricePlansResponse, error)
}
