package app

import (
	"errors"
	"log"

	"autoshop/internal/pkg/apperr"
)

// bindError converts an internal error into one safe to surface to the UI.
// Validation/NotFound/Conflict messages are user-facing and passed through;
// unexpected internal errors are logged and replaced with a generic message so
// implementation details never leak to the shopkeeper.
func bindError(err error) error {
	if err == nil {
		return nil
	}
	var ae *apperr.Error
	if errors.As(err, &ae) {
		if ae.Kind == apperr.KindInternal {
			log.Printf("internal error: %v", ae)
			return errors.New("Something went wrong. Please try again.")
		}
		return errors.New(ae.Message)
	}
	log.Printf("unexpected error: %v", err)
	return errors.New("Unexpected error occurred.")
}
