package billable_charge

import (
	"errors"
	"strconv"
	"strings"
)

// Error codes carried by a use-case refusal (`ErrorCode()`); each maps to the
// Lyngua key `billable_charge.errors.<code>` (Labels.ErrorMessage).
const (
	ErrKindNone                = ""
	ErrKindValidation          = "validation"
	ErrKindNotFound            = "not_found"
	ErrKindNotIssued           = "not_issued"
	ErrKindAdjustNotDownward   = "adjust_not_downward"
	ErrKindTransactionRequired = "transaction_required"
	ErrKindLockUnavailable     = "lock_unavailable"
	ErrKindObligationConflict  = "obligation_conflict"
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

// ErrorMessage maps a billable_charge error code to its Lyngua message; an
// unmapped code returns "" (the caller falls back to the shared general error).
func (l Labels) ErrorMessage(code string) string {
	switch code {
	case ErrKindValidation:
		return l.Errors.Validation
	case ErrKindNotFound:
		return l.Errors.NotFound
	case ErrKindNotIssued:
		return l.Errors.NotIssued
	case ErrKindAdjustNotDownward:
		return l.Errors.AdjustNotDownward
	case ErrKindTransactionRequired:
		return l.Errors.TransactionRequired
	case ErrKindLockUnavailable:
		return l.Errors.LockUnavailable
	case ErrKindObligationConflict:
		return l.Errors.ObligationConflict
	default:
		return ""
	}
}

// Format substitutes {0}, {1}, ... in a Lyngua template with args.
func Format(tmpl string, args ...string) string {
	for i, a := range args {
		tmpl = strings.ReplaceAll(tmpl, "{"+itoa(i)+"}", a)
	}
	return tmpl
}

func itoa(i int) string { return strconv.Itoa(i) }
