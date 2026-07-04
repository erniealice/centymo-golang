package procurement_request

import "github.com/erniealice/espyna-golang/consumer/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "expenditure.procurement_request",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "procurement_request"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "procurement_request.json", Key: "procurementRequest"},
		LabelName: "ProcurementRequestLabels",
		Templates: TemplatesFS,
		Nav: compose.NavContrib{
			Permission: "procurement_request:list",
			Items: []compose.NavItem{
				// procurement app — Requests inbox by lifecycle state
				{Key: "pr-pending", Route: "procurement_request.list", Params: map[string]string{"status": "pending_approval"},
					Label: "Pending Approval", Icon: "icon-clock", Permission: "procurement_request:list", LabelKey: "pending_approval_label", IconKey: "procurement_requests_pending_icon"},
				{Key: "pr-draft", Route: "procurement_request.list", Params: map[string]string{"status": "draft"},
					Label: "Draft", Icon: "icon-file-text", Permission: "procurement_request:list", LabelKey: "draft_label", IconKey: "procurement_requests_draft_icon"},
				{Key: "pr-approved", Route: "procurement_request.list", Params: map[string]string{"status": "approved"},
					Label: "Approved", Icon: "icon-check-circle", Permission: "procurement_request:list", LabelKey: "approved_label", IconKey: "procurement_requests_approved_icon"},
				{Key: "pr-rejected", Route: "procurement_request.list", Params: map[string]string{"status": "rejected"},
					Label: "Rejected", Icon: "icon-x-circle", Permission: "procurement_request:list", LabelKey: "rejected_label", IconKey: "procurement_requests_rejected_icon"},
			},
		},
	}
}
