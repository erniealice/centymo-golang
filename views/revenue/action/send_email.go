package action

import (
	"context"
	"fmt"
	"log"
	"net/http"

	centymo "github.com/erniealice/centymo-golang"
	"github.com/erniealice/fycha-golang/services/pdfconv"

	revenuepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/revenue"
	revenuelineitempb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/revenue_line_item"
)

// SendEmailDeps holds dependencies for the send-email handler.
type SendEmailDeps struct {
	Routes centymo.RevenueRoutes
	Labels centymo.RevenueLabels

	// Revenue operations
	ReadRevenue          func(ctx context.Context, req *revenuepb.ReadRevenueRequest) (*revenuepb.ReadRevenueResponse, error)
	ListRevenueLineItems func(ctx context.Context, req *revenuelineitempb.ListRevenueLineItemsRequest) (*revenuelineitempb.ListRevenueLineItemsResponse, error)

	// Document generation
	GenerateDoc         func(templateData []byte, data map[string]any) ([]byte, error)
	LoadDefaultTemplate func(ctx context.Context, purpose string) ([]byte, error)

	// Email sending function (injected from espyna email adapter)
	SendEmail func(ctx context.Context, to []string, subject, htmlBody, textBody string, attachmentName string, attachmentData []byte) error
}

// NewSendEmailHandler creates an http.HandlerFunc that generates an invoice and emails it.
//
// Query parameters:
//   - format: "pdf" (default) or "docx" — controls the attachment file format.
func NewSendEmailHandler(deps *SendEmailDeps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := r.PathValue("id")
		if id == "" {
			http.Error(w, "sale ID required", http.StatusBadRequest)
			return
		}

		// Determine attachment format from query param (default: pdf)
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
			log.Printf("send-email: failed to read revenue %s: %v", id, err)
			http.Error(w, "failed to load sale", http.StatusInternalServerError)
			return
		}
		data := resp.GetData()
		if len(data) == 0 {
			http.Error(w, "sale not found", http.StatusNotFound)
			return
		}
		revenue := data[0]

		// 2. Get customer email from revenue → client → user chain
		var customerEmail string
		if client := revenue.GetClient(); client != nil {
			if user := client.GetUser(); user != nil {
				customerEmail = user.GetEmailAddress()
			}
		}
		customerName := revenue.GetName()
		refNumber := revenue.GetReferenceNumber()
		if refNumber == "" {
			refNumber = id
		}

		if customerEmail == "" {
			log.Printf("send-email: no email address for customer %s on revenue %s", customerName, id)
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "No email address found for customer: %s", customerName)
			return
		}

		// 3. Read line items
		lineItemResp, err := deps.ListRevenueLineItems(ctx, &revenuelineitempb.ListRevenueLineItemsRequest{
			RevenueId: &id,
		})
		if err != nil {
			log.Printf("send-email: failed to list line items for %s: %v", id, err)
			http.Error(w, "failed to load line items", http.StatusInternalServerError)
			return
		}
		var lineItems []*revenuelineitempb.RevenueLineItem
		for _, item := range lineItemResp.GetData() {
			if item.GetRevenueId() == id {
				lineItems = append(lineItems, item)
			}
		}

		// 4. Build invoice data and generate document
		invoiceData := buildInvoiceData(revenue, lineItems)
		templateBytes, err := loadTemplate(ctx, deps.LoadDefaultTemplate)
		if err != nil {
			log.Printf("send-email: failed to load template: %v", err)
			http.Error(w, "failed to load template", http.StatusInternalServerError)
			return
		}
		docBytes, err := deps.GenerateDoc(templateBytes, invoiceData)
		if err != nil {
			log.Printf("send-email: failed to generate document: %v", err)
			http.Error(w, "failed to generate invoice", http.StatusInternalServerError)
			return
		}

		// 5. Prepare attachment in requested format
		var attachmentBytes []byte
		var attachmentName string

		if format == "pdf" {
			pdfBytes, ok, convErr := pdfconv.ConvertDocxToPDF(docBytes)
			if convErr != nil {
				log.Printf("send-email: PDF conversion failed, attaching DOCX: %v", convErr)
				attachmentBytes = docBytes
				attachmentName = fmt.Sprintf("invoice-%s.docx", refNumber)
			} else if ok {
				attachmentBytes = pdfBytes
				attachmentName = fmt.Sprintf("invoice-%s.pdf", refNumber)
			} else {
				log.Printf("send-email: LibreOffice not installed, attaching DOCX")
				attachmentBytes = docBytes
				attachmentName = fmt.Sprintf("invoice-%s.docx", refNumber)
			}
		} else {
			attachmentBytes = docBytes
			attachmentName = fmt.Sprintf("invoice-%s.docx", refNumber)
		}

		// 6. Send email with invoice attachment
		subject := fmt.Sprintf("Invoice %s", refNumber)
		textBody := fmt.Sprintf("Dear %s,\n\nPlease find attached your invoice %s.\n\nThank you for your business.", customerName, refNumber)
		htmlBody := fmt.Sprintf("<p>Dear %s,</p><p>Please find attached your invoice <strong>%s</strong>.</p><p>Thank you for your business.</p>", customerName, refNumber)

		err = deps.SendEmail(ctx, []string{customerEmail}, subject, htmlBody, textBody, attachmentName, attachmentBytes)
		if err != nil {
			log.Printf("send-email: failed to send email for revenue %s: %v", id, err)
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Failed to send email: %v", err)
			return
		}

		log.Printf("send-email: invoice %s sent to %s", refNumber, customerEmail)

		// Return HTMX success — refresh table
		w.Header().Set("HX-Trigger", `{"showToast":"Invoice email sent successfully"}`)
		w.WriteHeader(http.StatusOK)
	}
}
