package views_test

// Verification for the Plan 6 CSP onclick= -> lf.on() delegation sweep.
//
// Asserts that:
//  1. The full centymo template set (converted templates + the shared
//     "centymo-lf-delegation" define that lives in views/revenue/templates/
//     lf-delegation.html) PARSES with no html/template error — this catches any
//     malformed attribute / unbalanced action introduced by the sweep.
//  2. The "centymo-lf-delegation" template is DEFINED in the shared set, so every
//     {{template "centymo-lf-delegation" .}} include the sweep added resolves at
//     execute time (html/template does not flag an undefined {{template}} at parse
//     time, so this is checked explicitly via Lookup).
//
// This is a throwaway parse/wiring guard for the sweep; it constructs the
// renderer exactly like apps/service-admin container.go does (pyeza.SharedFS +
// the per-view embed.FS list), but only for the centymo views the sweep touched.
//
// NOTE: this test is self-contained to the centymo module and does NOT touch the
// service-admin renderer registration.

import (
	"io/fs"
	"testing"

	pyeza "github.com/erniealice/pyeza-golang"

	accrued_expense "github.com/erniealice/centymo-golang/views/accrued_expense"
	accrued_expense_settlement "github.com/erniealice/centymo-golang/views/accrued_expense_settlement"
	expenditure "github.com/erniealice/centymo-golang/views/expenditure"
	expense_recognition "github.com/erniealice/centymo-golang/views/expense_recognition"
	expense_recognition_line "github.com/erniealice/centymo-golang/views/expense_recognition_line"
	inventory "github.com/erniealice/centymo-golang/domain/inventory/views/inventory"
	plan "github.com/erniealice/centymo-golang/domain/subscription/views/plan"
	price_plan "github.com/erniealice/centymo-golang/domain/subscription/views/price_plan"
	price_schedule "github.com/erniealice/centymo-golang/domain/subscription/views/price_schedule"
	procurement_request "github.com/erniealice/centymo-golang/views/procurement_request"
	procurement_request_line "github.com/erniealice/centymo-golang/views/procurement_request_line"
	product "github.com/erniealice/centymo-golang/domain/product/views/product"
	purchase_order "github.com/erniealice/centymo-golang/views/purchase_order"
	revenue "github.com/erniealice/centymo-golang/domain/revenue/views/revenue"
	subscription "github.com/erniealice/centymo-golang/domain/subscription/views/subscription"
	supplier_contract "github.com/erniealice/centymo-golang/views/supplier_contract"
	supplier_contract_line "github.com/erniealice/centymo-golang/views/supplier_contract_line"
	supplier_contract_price_schedule "github.com/erniealice/centymo-golang/views/supplier_contract_price_schedule"
	supplier_contract_price_schedule_line "github.com/erniealice/centymo-golang/views/supplier_contract_price_schedule_line"
)

func TestLfDelegationSweepParsesAndDefineResolves(t *testing.T) {
	filesystems := []fs.FS{
		pyeza.SharedFS, // base components: app-shell, tabs, table-card, sheet-*, icons, etc.

		// Every centymo view the onclick sweep touched (those whose FS is in the
		// service-admin renderer list). revenue carries the shared define.
		revenue.TemplatesFS,
		inventory.TemplatesFS,
		product.TemplatesFS,
		plan.TemplatesFS,
		price_plan.TemplateFS,
		price_schedule.TemplateFS,
		subscription.TemplatesFS,
		expenditure.TemplatesFS,
		purchase_order.TemplatesFS,
		supplier_contract.TemplatesFS,
		supplier_contract_line.TemplatesFS,
		procurement_request.TemplatesFS,
		procurement_request_line.TemplatesFS,
		supplier_contract_price_schedule.TemplatesFS,
		supplier_contract_price_schedule_line.TemplatesFS,
		expense_recognition.TemplatesFS,
		expense_recognition_line.TemplatesFS,
		accrued_expense.TemplatesFS,
		accrued_expense_settlement.TemplatesFS,
	}

	r := pyeza.NewHTMLRendererFromFS(filesystems...)
	if err := r.Init(); err != nil {
		t.Fatalf("renderer.Init() failed to parse the swept centymo template set: %v", err)
	}

	tmpls := r.GetTemplates()
	if tmpls == nil {
		t.Fatal("renderer.GetTemplates() returned nil")
	}

	// The shared delegation define must exist so every {{template
	// "centymo-lf-delegation" .}} the sweep added resolves at execute time.
	if tmpls.Lookup("centymo-lf-delegation") == nil {
		t.Fatal(`template "centymo-lf-delegation" is not defined in the parsed set ` +
			`(it lives in views/revenue/templates/lf-delegation.html)`)
	}

	// Spot-check a representative converted top-level template per hook category
	// is present (Sheet open/close, nav, download, movements, stop-propagation).
	for _, name := range []string{
		"accrued-expense-detail-content",          // Sheet.open by title
		"accrued-expense-drawer-form",             // Sheet.close
		"inventory-detail-content",                // Sheet.open(this) trigger form
		"inventory-movements-content",             // movements clear + export hooks
		"revenue-detail-content",                  // downloadInvoice delegated trigger
		"product-detail-tab-body",                 // guarded-open -> data hook
		"supplier-contract-detail-content",        // multiple Sheet.open
		"variant-drawer-form",                     // stop-propagation
	} {
		if tmpls.Lookup(name) == nil {
			t.Errorf("expected converted template %q to be defined in the parsed set", name)
		}
	}
}
