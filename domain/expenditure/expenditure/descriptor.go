package expenditure

import "github.com/erniealice/pyeza-golang/compose"

// Describe returns the composition-v2 descriptor for the expenditure entity.
// Labels are not yet exposed via DefaultLabels() — the LabelJSON binding is
// left empty until a DefaultLabels factory is added.
func Describe() compose.Unit {
	r := DefaultRoutes()
	return compose.Unit{
		Key:       "expenditure.expenditure",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "expenditure"},
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
					Label: "Dashboard", Icon: "icon-layout-dashboard", Permission: "purchase:list"},
				// purchase app — Purchase Orders by status
				{Key: "po-draft", Route: "expenditure.purchase_order.list", Params: map[string]string{"status": "draft"},
					Label: "Draft", Icon: "icon-file-text", Permission: "purchase:list"},
				{Key: "po-approved", Route: "expenditure.purchase_order.list", Params: map[string]string{"status": "approved"},
					Label: "Approved", Icon: "icon-check-circle", Permission: "purchase:list"},
				{Key: "po-received", Route: "expenditure.purchase_order.list", Params: map[string]string{"status": "fully_received"},
					Label: "Received", Icon: "icon-package", Permission: "purchase:list"},
				{Key: "po-closed", Route: "expenditure.purchase_order.list", Params: map[string]string{"status": "closed"},
					Label: "Closed", Icon: "icon-archive", Permission: "purchase:list"},
				// purchase app — Purchases (expenditures) by status
				{Key: "purchases-all", Route: "expenditure.purchase.list", Params: map[string]string{"status": "all"},
					Label: "All", Icon: "icon-list", Permission: "purchase:list"},
				{Key: "purchases-pending", Route: "expenditure.purchase.list", Params: map[string]string{"status": "pending"},
					Label: "Pending", Icon: "icon-clock", Permission: "purchase:list"},
				{Key: "purchases-approved", Route: "expenditure.purchase.list", Params: map[string]string{"status": "approved"},
					Label: "Approved", Icon: "icon-check-circle", Permission: "purchase:list"},
				{Key: "purchases-paid", Route: "expenditure.purchase.list", Params: map[string]string{"status": "paid"},
					Label: "Paid", Icon: "icon-dollar-sign", Permission: "purchase:list"},
				// purchase settings
				{Key: "purchase-templates", Route: "purchases.settings.templates",
					Label: "Purchase Templates", Icon: "icon-file", Permission: "purchase:list"},
				// expense app — dashboard + items
				{Key: "expense-dashboard", Route: "expenditure.expense.dashboard",
					Label: "Dashboard", Icon: "icon-layout-dashboard", Permission: "expense:list"},
				{Key: "expenses-all", Route: "expenditure.expense.list", Params: map[string]string{"status": "all"},
					Label: "All", Icon: "icon-list", Permission: "expense:list"},
				{Key: "expenses-pending", Route: "expenditure.expense.list", Params: map[string]string{"status": "pending"},
					Label: "Pending", Icon: "icon-clock", Permission: "expense:list"},
				{Key: "expenses-approved", Route: "expenditure.expense.list", Params: map[string]string{"status": "approved"},
					Label: "Approved", Icon: "icon-check-circle", Permission: "expense:list"},
				{Key: "expenses-paid", Route: "expenditure.expense.list", Params: map[string]string{"status": "paid"},
					Label: "Paid", Icon: "icon-dollar-sign", Permission: "expense:list"},
				// expense settings
				{Key: "expense-categories", Route: "expenditure.expense_category.list",
					Label: "Expense Categories", Icon: "icon-tag", Permission: "expenditure_category:list"},
			},
		},
	}
}
