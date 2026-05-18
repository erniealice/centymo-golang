// Package expense_recognition_run wires the buying-side Expense Recognition Run
// view module — the Plan A mirror of the shipped Revenue Run module.
//
// Phase 1 (labels-first) only: this file exposes a thin alias for the
// canonical label struct defined in package centymo. Subsequent phases
// (2-7) extend this package with the four UI surfaces:
//
//   - Surface A: per-supplier drawer (entydad supplier statement)
//   - Surface B: workspace queue page (this package's queue/ sub-package)
//   - Surface C: per-supplier-subscription drawer (centymo supplier_subscription)
//   - Surface D: run history list + detail (this package's list/ + detail/ sub-packages)
//
// The Plan A reference is docs/plan/20260517-expense-run/plan.md.
package expense_recognition_run

import (
	centymo "github.com/erniealice/centymo-golang"
)

// Labels is a package-local alias for centymo.ExpenseRecognitionRunLabels.
// View sub-packages embed this alias on their *ViewDeps structs so that
// `deps.Labels.Detail.Title` reads naturally without a long type chain.
type Labels = centymo.ExpenseRecognitionRunLabels

// DefaultLabels is a package-local alias for
// centymo.DefaultExpenseRecognitionRunLabels. Composition wiring in
// apps/service-admin calls this to seed defaults before lyngua overlay.
var DefaultLabels = centymo.DefaultExpenseRecognitionRunLabels
