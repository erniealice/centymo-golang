package form

import (
	centymo "github.com/erniealice/centymo-golang"
)

// Data is the template data for the resource drawer form.
type Data struct {
	FormAction   string
	IsEdit       bool
	ID           string
	Name         string
	Description  string
	ProductId    string
	UserId       string
	Labels       centymo.ResourceFormLabels
	CommonLabels any
}
