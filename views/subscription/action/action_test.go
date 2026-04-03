package action

import (
	"testing"
	"time"
)

func TestFormatDateForInput(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		dateString string
		dateMillis int64
		want       string
	}{
		{
			name:       "date only string",
			dateString: "2026-04-03",
			want:       "2026-04-03",
		},
		{
			name:       "timestamp string trimmed for input",
			dateString: "2026-04-03T09:30:00Z",
			want:       "2026-04-03",
		},
		{
			name:       "falls back to millis",
			dateMillis: time.Date(2026, time.April, 3, 0, 0, 0, 0, time.UTC).UnixMilli(),
			want:       "2026-04-03",
		},
		{
			name: "empty when no date is available",
			want: "",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := formatDateForInput(tt.dateString, tt.dateMillis)
			if got != tt.want {
				t.Fatalf("formatDateForInput(%q, %d) = %q, want %q", tt.dateString, tt.dateMillis, got, tt.want)
			}
		})
	}
}
