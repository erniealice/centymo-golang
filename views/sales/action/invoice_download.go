package action

import (
	"context"
	"embed"
	"fmt"
	"log"
	"net/http"

	centymo "github.com/erniealice/centymo-golang"
	"github.com/erniealice/fycha-golang/services/pdfconv"

	revenuepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/revenue"
	revenuelineitempb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/revenue_line_item"
)

//go:embed templates/invoice-template.docx
var invoiceTemplateFS embed.FS

// InvoiceDownloadDeps holds dependencies for the invoice download handler.
type InvoiceDownloadDeps struct {
	Routes centymo.SalesRoutes
	Labels centymo.SalesLabels

	// Revenue operations
	ReadRevenue          func(ctx context.Context, req *revenuepb.ReadRevenueRequest) (*revenuepb.ReadRevenueResponse, error)
	ListRevenueLineItems func(ctx context.Context, req *revenuelineitempb.ListRevenueLineItemsRequest) (*revenuelineitempb.ListRevenueLineItemsResponse, error)

	// Document generation (injected by composition root — wraps fycha.DocumentService.ProcessBytes)
	GenerateDoc func(templateData []byte, data map[string]any) ([]byte, error)

	// Optional: load custom default template from storage (nil = use embedded fallback)
	LoadDefaultTemplate func(ctx context.Context, purpose string) ([]byte, error)
}

// NewInvoiceDownloadHandler creates an http.HandlerFunc that generates and downloads
// an invoice DOCX for a given sale/revenue.
func NewInvoiceDownloadHandler(deps *InvoiceDownloadDeps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := r.PathValue("id")
		if id == "" {
			http.Error(w, "sale ID required", http.StatusBadRequest)
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

		// 6. Convert DOCX to PDF (falls back to DOCX if LibreOffice unavailable)
		outputBytes, isPDF, err := pdfconv.ConvertDocxToPDF(docBytes)
		if err != nil {
			log.Printf("invoice download: PDF conversion failed: %v", err)
			http.Error(w, "failed to convert invoice to PDF", http.StatusInternalServerError)
			return
		}

		// 7. Send as file download
		refNumber := revenue.GetReferenceNumber()
		if refNumber == "" {
			refNumber = id
		}

		if isPDF {
			filename := fmt.Sprintf("invoice-%s.pdf", refNumber)
			w.Header().Set("Content-Type", "application/pdf")
			w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", filename))
		} else {
			filename := fmt.Sprintf("invoice-%s.docx", refNumber)
			w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")
			w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", filename))
		}
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
			"date":             revenue.GetRevenueDateString(),
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
