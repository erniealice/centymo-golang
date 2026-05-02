// Package form owns the template data shape for the primary product
// drawer (product-drawer-form.html). Pure types only — no Deps, no
// context.Context, no repository imports.
package form

import (
	"github.com/erniealice/pyeza-golang/types"

	centymo "github.com/erniealice/centymo-golang"
)

// Data is the template data for the product drawer form.
type Data struct {
	FormAction  string
	IsEdit      bool
	ID          string
	Name        string
	Description string
	Price       string
	Currency    string
	Active      bool
	LineID      string
	LineOptions []types.SelectOption
	// Model D — variant configurability and unit-of-measure fields.
	VariantMode string // "none" | "configurable"
	Unit        string
	// CanToggleVariantMode is false when the product already has option or
	// variant rows; the template renders the toggle as disabled + surfaces
	// VariantModeLockedHelp so the user understands why the setting is locked.
	// Defaults to true for the Add flow (no existing children).
	CanToggleVariantMode bool

	// Four-axis product taxonomy — each axis rendered as a <select>. The
	// current value is the stored value (edit) or the mount default (add).
	// The Options slice is narrowed per-mount via Deps.Allowed*: services
	// mount shows only {service}, supplies mount only {consumable}, inventory
	// shows {stocked_good, non_stocked_good}. When len(Options) == 1 the
	// template renders the select disabled so the user still sees the locked
	// classification without being able to change it.
	ProductKind         string
	ProductKindOptions  []types.SelectOption
	DeliveryMode        string
	DeliveryModeOptions []types.SelectOption
	TrackingMode        string
	TrackingModeOptions []types.SelectOption

	Labels       centymo.ProductFormLabels
	CommonLabels any
}
