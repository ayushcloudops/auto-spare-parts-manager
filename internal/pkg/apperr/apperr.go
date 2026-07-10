// Package apperr defines a small, typed error model shared across the backend.
//
// Errors carry a Kind so the binding layer can translate them into a stable
// shape for the frontend (e.g. a NotFound becomes a 404-like UI state, a
// Validation error is shown inline on a form). Wrapping preserves the original
// cause for logging while keeping a clean message for the user.
package apperr

import (
	"errors"
	"fmt"
)

// Kind classifies an error so callers can branch on category, not on string.
type Kind string

const (
	KindValidation Kind = "validation" // bad user input
	KindNotFound   Kind = "not_found"  // entity does not exist
	KindConflict   Kind = "conflict"   // unique constraint / business rule clash
	KindInternal   Kind = "internal"   // unexpected; log it
)

// Error is the application error type.
type Error struct {
	Kind    Kind   // category, for the UI to branch on
	Message string // safe, user-facing message
	cause   error  // wrapped underlying error (for logs), never shown raw
}

// Error implements the error interface.
func (e *Error) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s: %s: %v", e.Kind, e.Message, e.cause)
	}
	return fmt.Sprintf("%s: %s", e.Kind, e.Message)
}

// Unwrap exposes the underlying cause for errors.Is / errors.As.
func (e *Error) Unwrap() error { return e.cause }

// new builds an *Error with optional wrapped cause.
func newErr(kind Kind, cause error, format string, args ...any) *Error {
	return &Error{Kind: kind, Message: fmt.Sprintf(format, args...), cause: cause}
}

// Validation creates a validation error.
func Validation(format string, args ...any) *Error {
	return newErr(KindValidation, nil, format, args...)
}

// NotFound creates a not-found error.
func NotFound(format string, args ...any) *Error {
	return newErr(KindNotFound, nil, format, args...)
}

// Conflict creates a conflict error.
func Conflict(format string, args ...any) *Error {
	return newErr(KindConflict, nil, format, args...)
}

// Internal wraps an unexpected error as an internal error.
func Internal(cause error, format string, args ...any) *Error {
	return newErr(KindInternal, cause, format, args...)
}

// KindOf returns the Kind of err, defaulting to KindInternal for plain errors.
func KindOf(err error) Kind {
	var ae *Error
	if errors.As(err, &ae) {
		return ae.Kind
	}
	return KindInternal
}
