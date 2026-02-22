package detail

import (
	"fmt"
	"net/http"

	"github.com/erniealice/pyeza-golang/view"
)

// HtmxSuccess returns a view result that triggers sheet close and table refresh via HTMX events.
func HtmxSuccess(tableID string) view.ViewResult {
	return view.ViewResult{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"HX-Trigger": fmt.Sprintf(`{"formSuccess":true,"refreshTable":"%s"}`, tableID),
		},
	}
}

// HtmxError returns a view result with an error message header for HTMX error handling.
func HtmxError(message string) view.ViewResult {
	return view.ViewResult{
		StatusCode: http.StatusUnprocessableEntity,
		Headers: map[string]string{
			"HX-Error-Message": message,
		},
	}
}
