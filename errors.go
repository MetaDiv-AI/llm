package llm

import (
	"errors"
	"fmt"
)

// ErrUnknownProvider is returned when the provider is not supported.
type ErrUnknownProvider struct {
	Provider string
}

func (e *ErrUnknownProvider) Error() string {
	return fmt.Sprintf("llm: unknown provider %q", e.Provider)
}

// Is supports errors.Is for ErrUnknownProvider.
func (e *ErrUnknownProvider) Is(target error) bool {
	t, ok := target.(*ErrUnknownProvider)
	return ok && (t == nil || t.Provider == "" || t.Provider == e.Provider)
}

// ErrInvalidRequest is returned when a request fails validation.
var ErrInvalidRequest = errors.New("llm: invalid request")

// ValidationError represents a validation failure with field and message.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("llm: %s: %s", e.Field, e.Message)
}

// Is supports errors.Is for ValidationError and ErrInvalidRequest.
func (e *ValidationError) Is(target error) bool {
	if errors.Is(target, ErrInvalidRequest) {
		return true
	}
	t, ok := target.(*ValidationError)
	return ok && (t == nil || (t.Field == e.Field && t.Message == e.Message))
}
