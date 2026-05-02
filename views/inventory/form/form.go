// Package form owns the template data shape for the primary inventory item
// drawer (inventory-drawer-form.html). Pure types only — no Deps, no
// context.Context, no repository imports.
package form

// Labels holds i18n labels for the inventory drawer form template.
type Labels struct {
	Product          string
	SKU              string
	SKUPlaceholder   string
	OnHand           string
	Reserved         string
	ReorderLevel     string
	UnitOfMeasure    string
	Notes            string
	NotesPlaceholder string
	Active           string

	// Field-level info text surfaced via an info button beside each label.
	ProductInfo       string
	SKUInfo           string
	OnHandInfo        string
	ReservedInfo      string
	ReorderLevelInfo  string
	UnitOfMeasureInfo string
	NotesInfo         string
	ActiveInfo        string
}

// Data is the template data for the inventory drawer form.
type Data struct {
	FormAction    string
	IsEdit        bool
	ID            string
	Name          string
	SKU           string
	OnHand        string
	Reserved      string
	ReorderLevel  string
	UnitOfMeasure string
	LocationID    string
	Notes         string
	Active        bool
	Labels        Labels
	CommonLabels  any
}
