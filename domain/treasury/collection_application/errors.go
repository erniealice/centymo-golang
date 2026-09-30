package collection_application

import "errors"

// Error codes carried by a use-case refusal (`ErrorCode()`); each maps to the
// Lyngua key `collection_application.errors.<code>` (Labels.ErrorMessage).
const (
	ErrKindNone                = ""
	ErrKindCurrencyMismatch    = "currency_mismatch"
	ErrKindNotFound            = "not_found"
	ErrKindAmountInvalid       = "amount_invalid"
	ErrKindClientRequired      = "client_required"
	ErrKindAlreadyReversed     = "already_reversed"
	ErrKindMissingPosting      = "missing_posting"
	ErrKindTransactionRequired = "transaction_required"
	ErrKindDateInvalid         = "date_invalid"
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

// ErrorMessage maps a collection_application error code to its Lyngua message;
// an unmapped code returns "" (the caller falls back to the shared general error).
func (l Labels) ErrorMessage(code string) string {
	switch code {
	case ErrKindCurrencyMismatch:
		return l.Errors.CurrencyMismatch
	case ErrKindNotFound:
		return l.Errors.NotFound
	case ErrKindAmountInvalid:
		return l.Errors.AmountInvalid
	case ErrKindClientRequired:
		return l.Errors.ClientRequired
	case ErrKindAlreadyReversed:
		return l.Errors.AlreadyReversed
	case ErrKindMissingPosting:
		return l.Errors.MissingPosting
	case ErrKindTransactionRequired:
		return l.Errors.TransactionRequired
	case ErrKindDateInvalid:
		return l.Errors.DateInvalid
	default:
		return ""
	}
}
