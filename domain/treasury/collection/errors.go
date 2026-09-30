package collection

import "errors"

// Error codes carried by a use-case refusal (`ErrorCode()`); each maps to the
// Lyngua key `collection.errors.<code>` (Labels.ErrorMessage).
const (
	ErrKindNone = ""
	// ErrKindReceiptHasApplications: a receipt referenced by an active collection
	// application is immutable (edit, status change and delete are refused).
	ErrKindReceiptHasApplications = "receipt_has_applications"
	ErrKindUnknown                = "unknown"
)

// CollectionTypeReceipt is the collection_type a Receive & apply receipt carries.
const CollectionTypeReceipt = "receipt"

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

// ErrorMessage maps a collection error code to its Lyngua message; an unmapped
// code returns "" (the caller falls back to the error text).
func (l Labels) ErrorMessage(code string) string {
	switch code {
	case ErrKindReceiptHasApplications:
		return l.Errors.ReceiptHasApplications
	default:
		return ""
	}
}
