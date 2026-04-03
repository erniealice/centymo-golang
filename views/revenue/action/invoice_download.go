package action

import (
	"context"
	"embed"
	"fmt"
	"log"
	"net/http"
	"time"

	centymo "github.com/erniealice/centymo-golang"
	"github.com/erniealice/fycha-golang/services/pdfconv"

	revenuepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/revenue"
	revenuelineitempb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/revenue_line_item"
)

//go:embed templates/invoice-template.docx
var invoiceTemplateFS embed.FS

// InvoiceDownloadDeps holds dependencies for the invoice download handler.
type InvoiceDownloadDeps struct {
	Routes centymo.RevenueRoutes
	Labels centymo.RevenueLabels

	// Revenue operations
	ReadRevenue          func(ctx context.Context, req *revenuepb.ReadRevenueRequest) (*revenuepb.ReadRevenueResponse, error)
	ListRevenueLineItems func(ctx context.Context, req *revenuelineitempb.ListRevenueLineItemsRequest) (*revenuelineitempb.ListRevenueLineItemsResponse, error)

	// Document generation (injected by composition root — wraps fycha.DocumentService.ProcessBytes)
	GenerateDoc func(templateData []byte, data map[string]any) ([]byte, error)

	// Optional: load custom default template from storage (nil = use embedded fallback)
	LoadDefaultTemplate func(ctx context.Context, purpose string) ([]byte, error)
}

// NewInvoiceDownloadHandler creates an http.HandlerFunc that generates and downloads
// an invoice for a given sale/revenue.
//
// Query parameters:
//   - format: "pdf" (default) or "docx" — controls the output file format.
//     Example: /action/sales/detail/{id}/invoice/download?format=docx
func NewInvoiceDownloadHandler(deps *InvoiceDownloadDeps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := r.PathValue("id")
		if id == "" {
			http.Error(w, "sale ID required", http.StatusBadRequest)
			return
		}

		// Determine output format from query param (default: pdf)
		format := r.URL.Query().Get("format")
		if format == "" {
			format = "pdf"
		}
		if format != "pdf" && format != "docx" {
			http.Error(w, "invalid format: must be \"pdf\" or \"docx\"", http.StatusBadRequest)
			return
		}

		// 1. Read revenue
		resp, err := deps.ReadRevenue(ctx, &revenuepb.ReadRevenueRequest{
			Data: &revenuepb.Revenue{Id: id},
		})
		if err != nil {
			log.Printf("invoice download: failed to read revenue %s: %v", id, err)
			http.Error(w, "failed to load sale", http.StatusInternalServerError)
			return
		}
		data := resp.GetData()
		if len(data) == 0 {
			http.Error(w, "sale not found", http.StatusNotFound)
			return
		}
		revenue := data[0]

		// 2. Read line items
		lineItemResp, err := deps.ListRevenueLineItems(ctx, &revenuelineitempb.ListRevenueLineItemsRequest{
			RevenueId: &id,
		})
		if err != nil {
			log.Printf("invoice download: failed to list line items for %s: %v", id, err)
			http.Error(w, "failed to load line items", http.StatusInternalServerError)
			return
		}

		// Filter line items belonging to this revenue
		var lineItems []*revenuelineitempb.RevenueLineItem
		for _, item := range lineItemResp.GetData() {
			if item.GetRevenueId() == id {
				lineItems = append(lineItems, item)
			}
		}

		// 3. Build invoice data map (matches doctemplate placeholder format)
		invoiceData := buildInvoiceData(revenue, lineItems)

		// 4. Load template (prefer custom default, fall back to embedded)
		templateBytes, err := loadTemplate(ctx, deps.LoadDefaultTemplate)
		if err != nil {
			log.Printf("invoice download: failed to load template: %v", err)
			http.Error(w, "failed to load template", http.StatusInternalServerError)
			return
		}

		// 5. Generate DOCX
		docBytes, err := deps.GenerateDoc(templateBytes, invoiceData)
		if err != nil {
			log.Printf("invoice download: failed to generate document: %v", err)
			http.Error(w, "failed to generate invoice", http.StatusInternalServerError)
			return
		}

		// 6. Prepare output based on requested format
		refNumber := revenue.GetReferenceNumber()
		if refNumber == "" {
			refNumber = id
		}
		ts := time.Now().Unix()

		var outputBytes []byte
		var contentType, filename string

		if format == "pdf" {
			pdfBytes, ok, convErr := pdfconv.ConvertDocxToPDF(docBytes)
			if convErr != nil {
				log.Printf("invoice download: PDF conversion failed: %v", convErr)
				http.Error(w, "failed to convert invoice to PDF", http.StatusInternalServerError)
				return
			}
			if !ok {
				log.Printf("invoice download: LibreOffice not installed — cannot generate PDF")
				http.Error(w, "PDF generation unavailable: LibreOffice is not installed on the server", http.StatusServiceUnavailable)
				return
			}
			outputBytes = pdfBytes
			contentType = "application/pdf"
			filename = fmt.Sprintf("invoice-%s-%d.pdf", refNumber, ts)
		} else {
			outputBytes = docBytes
			contentType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
			filename = fmt.Sprintf("invoice-%s-%d.docx", refNumber, ts)
		}

		// 7. Send as file download
		w.Header().Set("Content-Type", contentType)
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
		w.Write(outputBytes)
	}
}

// loadTemplate loads the invoice template. Tries custom default first, falls back to embedded.
func loadTemplate(ctx context.Context, loadDefault func(context.Context, string) ([]byte, error)) ([]byte, error) {
	// Try custom default template if available
	if loadDefault != nil {
		templateBytes, err := loadDefault(ctx, "invoice")
		if err == nil && len(templateBytes) > 0 {
			return templateBytes, nil
		}
		// Fall through to embedded if custom not found or errored
		if err != nil {
			log.Printf("invoice download: custom template load failed (using fallback): %v", err)
		}
	}

	// Use embedded fallback
	return invoiceTemplateFS.ReadFile("templates/invoice-template.docx")
}

// buildInvoiceData assembles the template data map from revenue + line items.
// This matches the doctemplate placeholder format: {{invoice.reference_number}}, {{#items}}, etc.
func buildInvoiceData(revenue *revenuepb.Revenue, lineItems []*revenuelineitempb.RevenueLineItem) map[string]any {
	// Build line item array
	items := make([]any, 0, len(lineItems))
	for _, item := range lineItems {
		items = append(items, map[string]any{
			"description": item.GetDescription(),
			"quantity":    fmt.Sprintf("%.0f", item.GetQuantity()),
			"unit_price":  fmt.Sprintf("%.2f", item.GetUnitPrice()),
			"total":       fmt.Sprintf("%.2f", item.GetTotalPrice()),
		})
	}

	return map[string]any{
		"invoice": map[string]any{
			"reference_number": revenue.GetReferenceNumber(),
			"date":             revenue.GetRevenueDate(),
			"status":           revenue.GetStatus(),
			"currency":         revenue.GetCurrency(),
			"total_amount":     fmt.Sprintf("%.2f", revenue.GetTotalAmount()),
			"notes":            revenue.GetNotes(),
		},
		"customer": map[string]any{
			"name": revenue.GetName(),
		},
		"items":    items,
		"total":    fmt.Sprintf("%.2f", revenue.GetTotalAmount()),
		"currency": revenue.GetCurrency(),
	}
}
