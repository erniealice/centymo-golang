package form

import (
	centymo "github.com/erniealice/centymo-golang"
)

// Labels holds i18n labels for the price list drawer form template.
type Labels struct {
	Name            string
	Description     string
	DescPlaceholder string
	DateStart       string
	DateEnd         string
	Active          string

	// Field-level info text surfaced via an info button beside each label.
	NameInfo        string
	DescriptionInfo string
	DateStartInfo   string
	DateEndInfo     string
	ActiveInfo      string
}

// Data is the template data for the price list drawer form.
type Data struct {
	FormAction   string
	IsEdit       bool
	ID           string
	Name         string
	Description  string
	DateStart    string
	DateEnd      string
	Active       bool
	Labels       Labels
	CommonLabels any
}

// BuildLabels assembles the drawer's Labels from the translation function and
// the centymo label struct. Real transformation: calls t(key) for most fields,
// pulls *Info strings directly from the centymo struct.
func BuildLabels(t func(string) string, f centymo.PriceListFormLabels) Labels {
	return Labels{
		Name:            t("pricelist.form.name"),
		Description:     t("pricelist.form.description"),
		DescPlaceholder: t("pricelist.form.descriptionPlaceholder"),
		DateStart:       t("pricelist.form.dateStart"),
		DateEnd:         t("pricelist.form.dateEnd"),
		Active:          t("pricelist.form.active"),
		NameInfo:        f.NameInfo,
		DescriptionInfo: f.DescriptionInfo,
		DateStartInfo:   f.DateStartInfo,
		DateEndInfo:     f.DateEndInfo,
		ActiveInfo:      f.ActiveInfo,
	}
}
