package expenditure

import "github.com/erniealice/espyna-golang/consumer/compose"

// Describe returns the composition-v2 descriptor for the expenditure entity.
func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "expenditure.expenditure",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "expenditure"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "expenditure.json", Key: "expenditure"},
		LabelName: "ExpenditureLabels",
		Templates: TemplatesFS,
		Nav: compose.NavContrib{
			Permission: "purchase:list",
			AppEntry: &compose.AppEntry{
				Key: "purchase", Route: "expenditure.purchase.dashboard",
				Label: "Purchases", Icon: "icon-shopping-cart",
				Permission: "purchase:list",
			},
			Items: []compose.NavItem{
				// purchase app — dashboard
				{Key: "dashboard", Route: "expenditure.purchase.dashboard",
					Label: "Dashboard", Icon: "icon-layout-dashboard", Permission: "purchase:list", LabelKey: "dashboard_label", IconKey: "dashboard_icon"},
				// purchase app — Purchase Orders by status
				{Key: "po-draft", Route: "expenditure.purchase_order.list", Params: map[string]string{"status": "draft"},
					Label: "Draft", Icon: "icon-file-text", Permission: "purchase:list"},
				{Key: "po-approved", Route: "expenditure.purchase_order.list", Params: map[string]string{"status": "approved"},
					Label: "Approved", Icon: "icon-check-circle", Permission: "purchase:list", LabelKey: "approved_label"},
				{Key: "po-received", Route: "expenditure.purchase_order.list", Params: map[string]string{"status": "fully_received"},
					Label: "Received", Icon: "icon-package", Permission: "purchase:list"},
				{Key: "po-closed", Route: "expenditure.purchase_order.list", Params: map[string]string{"status": "closed"},
					Label: "Closed", Icon: "icon-archive", Permission: "purchase:list"},
				// purchase app — Purchases (expenditures) by status
				{Key: "purchases-all", Route: "expenditure.purchase.list", Params: map[string]string{"status": "all"},
					Label: "All", Icon: "icon-list", Permission: "purchase:list", LabelKey: "all_label", IconKey: "all_icon"},
				{Key: "purchases-pending", Route: "expenditure.purchase.list", Params: map[string]string{"status": "pending"},
					Label: "Pending", Icon: "icon-clock", Permission: "purchase:list", LabelKey: "pending_label", IconKey: "purchases_pending_icon"},
				{Key: "purchases-approved", Route: "expenditure.purchase.list", Params: map[string]string{"status": "approved"},
					Label: "Approved", Icon: "icon-check-circle", Permission: "purchase:list", LabelKey: "approved_label", IconKey: "purchases_approved_icon"},
				{Key: "purchases-paid", Route: "expenditure.purchase.list", Params: map[string]string{"status": "paid"},
					Label: "Paid", Icon: "icon-dollar-sign", Permission: "purchase:list", LabelKey: "paid_label", IconKey: "purchases_paid_icon"},
				// purchase settings
				{Key: "purchase-templates", Route: "purchases.settings.templates",
					Label: "Purchase Templates", Icon: "icon-file", Permission: "purchase:list", LabelKey: "purchase_templates_label", IconKey: "purchase_templates_icon"},
				// expense app — dashboard + items
				{Key: "expense-dashboard", Route: "expenditure.expense.dashboard",
					Label: "Dashboard", Icon: "icon-layout-dashboard", Permission: "expense:list", LabelKey: "dashboard_label", IconKey: "dashboard_icon"},
				{Key: "expenses-all", Route: "expenditure.expense.list", Params: map[string]string{"status": "all"},
					Label: "All", Icon: "icon-list", Permission: "expense:list", LabelKey: "all_label", IconKey: "all_icon"},
				{Key: "expenses-pending", Route: "expenditure.expense.list", Params: map[string]string{"status": "pending"},
					Label: "Pending", Icon: "icon-clock", Permission: "expense:list", LabelKey: "pending_label", IconKey: "expenses_pending_icon"},
				{Key: "expenses-approved", Route: "expenditure.expense.list", Params: map[string]string{"status": "approved"},
					Label: "Approved", Icon: "icon-check-circle", Permission: "expense:list", LabelKey: "approved_label", IconKey: "expenses_approved_icon"},
				{Key: "expenses-paid", Route: "expenditure.expense.list", Params: map[string]string{"status": "paid"},
					Label: "Paid", Icon: "icon-dollar-sign", Permission: "expense:list", LabelKey: "paid_label", IconKey: "expenses_paid_icon"},
				// expense settings
				{Key: "expense-categories", Route: "expenditure.expense_category.list",
					Label: "Expense Categories", Icon: "icon-tag", Permission: "expenditure_category:list", LabelKey: "expense_categories_label", IconKey: "expense_categories_icon"},
			},
		},
	}
}
