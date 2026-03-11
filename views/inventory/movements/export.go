package movements

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"time"
)

// NewExportHandler creates an http.HandlerFunc for CSV export of movements.
// It applies the same filters as the table view.
func NewExportHandler(deps *Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		dateFrom := r.URL.Query().Get("date_from")
		dateTo := r.URL.Query().Get("date_to")
		location := r.URL.Query().Get("location")
		txType := r.URL.Query().Get("type")
		search := r.URL.Query().Get("search")

		// Default date range (same as page view)
		if dateFrom == "" {
			now := time.Now()
			dateFrom = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Format("2006-01-02")
		}
		if dateTo == "" {
			dateTo = time.Now().Format("2006-01-02")
		}

		// Build the table to get rows
		tableConfig := buildFilteredTable(ctx, deps, dateFrom, dateTo, location, txType, search)

		// Set CSV response headers
		filename := fmt.Sprintf("transactions-%s.csv", time.Now().Format("2006-01-02"))
		w.Header().Set("Content-Type", "text/csv; charset=utf-8")
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))

		writer := csv.NewWriter(w)
		defer writer.Flush()

		// Write header row from column labels
		headers := make([]string, len(tableConfig.Columns))
		for i, col := range tableConfig.Columns {
			headers[i] = col.Label
		}
		if err := writer.Write(headers); err != nil {
			log.Printf("Failed to write CSV header: %v", err)
			return
		}

		// Write data rows
		for _, row := range tableConfig.Rows {
			record := make([]string, len(row.Cells))
			for i, cell := range row.Cells {
				record[i] = fmt.Sprintf("%v", cell.Value)
			}
			if err := writer.Write(record); err != nil {
				log.Printf("Failed to write CSV row: %v", err)
				return
			}
		}
	}
}
