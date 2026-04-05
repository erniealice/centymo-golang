package action

import (
	"context"
	"fmt"
	"log"
	"net/http"

	centymo "github.com/erniealice/centymo-golang"

	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"

	expenditurepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/expenditure"
	disbursementpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/disbursement"
)

// PayFormData is the template data for the pay drawer form.
type PayFormData struct {
	FormAction       string
	ExpenditureID    string
	Name             string
	Amount           string
	Currency         string
	DisbursementType string
	Labels           centymo.DisbursementFormLabels
	CommonLabels     any
}

// PayDeps holds dependencies for the expenditure pay action.
type PayDeps struct {
	ExpenditureRoutes  centymo.ExpenditureRoutes
	DisbursementRoutes centymo.DisbursementRoutes
	DisbursementLabels centymo.DisbursementLabels
	ReadExpenditure    func(ctx context.Context, req *expenditurepb.ReadExpenditureRequest) (*expenditurepb.ReadExpenditureResponse, error)
	CreateDisbursement func(ctx context.Context, req *disbursementpb.CreateDisbursementRequest) (*disbursementpb.CreateDisbursementResponse, error)
}

// expenditureTypeToDisbursementType maps an expenditure type to a sensible disbursement type default.
func expenditureTypeToDisbursementType(expenditureType string) string {
	switch expenditureType {
	case "purchase":
		return "supplier"
	case "payroll":
		return "payroll"
	case "rent":
		return "rent"
	case "utilities":
		return "utilities"
	default:
		return "other"
	}
}

// NewPayAction creates a disbursement pre-linked to an expenditure.
//
// GET: returns a pre-filled drawer form.
// POST: creates the disbursement and redirects to disbursement detail.
func NewPayAction(deps *PayDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("disbursement", "create") {
			return centymo.HTMXError(deps.DisbursementLabels.Errors.PermissionDenied)
		}

		id := viewCtx.Request.PathValue("id")

		// Read the expenditure to pre-fill the form.
		resp, err := deps.ReadExpenditure(ctx, &expenditurepb.ReadExpenditureRequest{
			Data: &expenditurepb.Expenditure{Id: id},
		})
		if err != nil {
			log.Printf("PayAction: failed to read expenditure %s: %v", id, err)
			return centymo.HTMXError("Expense not found")
		}
		data := resp.GetData()
		if len(data) == 0 {
			return centymo.HTMXError("Expense not found")
		}
		exp := data[0]

		if viewCtx.Request.Method == http.MethodGet {
			amountStr := fmt.Sprintf("%.2f", float64(exp.GetTotalAmount())/100.0)
			disbursementType := expenditureTypeToDisbursementType(exp.GetExpenditureType())
			payeeName := exp.GetName()
			if payeeName == "" {
				payeeName = exp.GetReferenceNumber()
			}

			return view.OK("expense-pay-drawer-form", &PayFormData{
				FormAction:       route.ResolveURL(deps.ExpenditureRoutes.PayURL, "id", id),
				ExpenditureID:    id,
				Name:             "Payment - " + payeeName,
				Amount:           amountStr,
				Currency:         exp.GetCurrency(),
				DisbursementType: disbursementType,
				Labels:           deps.DisbursementLabels.Form,
				CommonLabels:     nil, // injected by ViewAdapter
			})
		}

		// POST — create the disbursement pre-linked to this expenditure.
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.DisbursementLabels.Errors.InvalidFormData)
		}

		r := viewCtx.Request

		createResp, err := deps.CreateDisbursement(ctx, &disbursementpb.CreateDisbursementRequest{
			Data: &disbursementpb.Disbursement{
				ReferenceNumber:      r.FormValue("reference_number"),
				Name:                 r.FormValue("payee"),
				Amount:               parseAmount(r.FormValue("amount")),
				Currency:             r.FormValue("currency"),
				DisbursementMethodId: r.FormValue("disbursement_method"),
				ApprovedBy:           r.FormValue("approved_by"),
				DisbursementType:     r.FormValue("disbursement_type"),
				ExpenditureId:        id,
				Status:               "draft",
			},
		})
		if err != nil {
			log.Printf("PayAction: failed to create disbursement for expenditure %s: %v", id, err)
			return centymo.HTMXError(err.Error())
		}

		newID := ""
		if respData := createResp.GetData(); len(respData) > 0 {
			newID = respData[0].GetId()
		}
		if newID != "" {
			return view.ViewResult{
				StatusCode: http.StatusOK,
				Headers: map[string]string{
					"HX-Trigger":  `{"formSuccess":true}`,
					"HX-Redirect": route.ResolveURL(deps.DisbursementRoutes.DetailURL, "id", newID),
				},
			}
		}

		return centymo.HTMXSuccess("")
	})
}
