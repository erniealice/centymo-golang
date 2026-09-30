package recovery_document

import (
	"errors"
	"strconv"
	"strings"
)

// Error codes carried by a use-case refusal (`ErrorCode()`); each maps to the
// Lyngua key `recovery_document.errors.<code>` (Labels.ErrorMessage).
const (
	ErrKindNone                  = ""
	ErrKindNothingSelected       = "nothing_selected"
	ErrKindNotOpen               = "not_open"
	ErrKindNotFound              = "not_found"
	ErrKindNotIssued             = "not_issued"
	ErrKindCurrencyMismatch      = "currency_mismatch"
	ErrKindMissingPosting        = "missing_posting"
	ErrKindSeriesRetired         = "series_retired"
	ErrKindSeriesKindMismatch    = "series_kind_mismatch"
	ErrKindComponentKindMismatch = "component_kind_mismatch"
	ErrKindIssuanceConflict      = "issuance_conflict"
	ErrKindHasApplications       = "has_applications"
	ErrKindHasCreditNotes        = "has_credit_notes"
	ErrKindAlreadyVoid           = "already_void"
	ErrKindVoidReasonRequired    = "void_reason_required"
	ErrKindValidation            = "validation"
	ErrKindTransactionRequired   = "transaction_required"
	// ErrKindPermissionDenied is the strict-gate denial (actiongate `permission_denied`);
	// views map it to the shared PermissionDenied message, not the generic error.
	ErrKindPermissionDenied = "permission_denied"
	ErrKindUnknown          = "unknown"
)

// ErrorKind classifies a use-case error through its `ErrorCode()` (no espyna
// import). Unknown / uncoded errors classify as ErrKindUnknown.
func ErrorKind(err error) string {
	if err == nil {
		return ErrKindNone
	}
	var coded interface{ ErrorCode() string }
	if errors.As(err, &coded) {
		if c := coded.ErrorCode(); c != "" {
			return c
		}
	}
	return ErrKindUnknown
}

// ErrorMessage maps a recovery_document error code to its Lyngua message; an
// unmapped code returns "" (the caller falls back to the shared general error).
func (l Labels) ErrorMessage(code string) string {
	switch code {
	case ErrKindNothingSelected:
		return l.Errors.NothingSelected
	case ErrKindNotOpen:
		return l.Errors.NotOpen
	case ErrKindNotFound:
		return l.Errors.NotFound
	case ErrKindNotIssued:
		return l.Errors.NotIssued
	case ErrKindCurrencyMismatch:
		return l.Errors.CurrencyMismatch
	case ErrKindMissingPosting:
		return l.Errors.MissingPosting
	case ErrKindSeriesRetired:
		return l.Errors.SeriesRetired
	case ErrKindSeriesKindMismatch:
		return l.Errors.SeriesKindMismatch
	case ErrKindComponentKindMismatch:
		return l.Errors.ComponentKindMismatch
	case ErrKindIssuanceConflict:
		return l.Errors.IssuanceConflict
	case ErrKindHasApplications:
		return l.Errors.HasApplications
	case ErrKindHasCreditNotes:
		return l.Errors.HasCreditNotes
	case ErrKindAlreadyVoid:
		return l.Errors.AlreadyVoid
	case ErrKindVoidReasonRequired:
		return l.Errors.VoidReasonRequired
	case ErrKindValidation:
		return l.Errors.Validation
	case ErrKindTransactionRequired:
		return l.Errors.TransactionRequired
	default:
		return ""
	}
}

// Format substitutes {0}, {1}, ... in a Lyngua template with args.
func Format(tmpl string, args ...string) string {
	for i, a := range args {
		tmpl = strings.ReplaceAll(tmpl, "{"+strconv.Itoa(i)+"}", a)
	}
	return tmpl
}
