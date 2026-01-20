package kyt

import (
	"errors"
	"fmt"
)

// Sentinel errors for common error conditions.
var (
	// ErrInvalidConfig indicates the provider configuration is invalid.
	ErrInvalidConfig = errors.New("invalid provider configuration")

	// ErrProviderNotFound indicates the requested provider does not exist.
	ErrProviderNotFound = errors.New("provider not found")

	// ErrUnsupportedChain indicates the blockchain is not supported by the provider.
	ErrUnsupportedChain = errors.New("unsupported blockchain")

	// ErrInvalidAddress indicates the address format is invalid.
	ErrInvalidAddress = errors.New("invalid address format")

	// ErrInvalidTxHash indicates the transaction hash format is invalid.
	ErrInvalidTxHash = errors.New("invalid transaction hash")

	// ErrAssessmentPending indicates the assessment is still in progress.
	ErrAssessmentPending = errors.New("assessment is still pending")

	// ErrRateLimited indicates the API rate limit has been exceeded.
	ErrRateLimited = errors.New("rate limit exceeded")
)

// ErrorType categorizes the type of error for handling decisions.
type ErrorType int

const (
	// ErrorTypeUnknown indicates an unclassified error.
	ErrorTypeUnknown ErrorType = iota

	// ErrorTypeValidation indicates a validation error (bad input).
	ErrorTypeValidation

	// ErrorTypeProvider indicates a provider-specific error.
	ErrorTypeProvider

	// ErrorTypeRetryable indicates an error that may succeed on retry.
	ErrorTypeRetryable

	// ErrorTypeRateLimit indicates a rate limiting error.
	ErrorTypeRateLimit
)

// Error represents a unified KYT error with additional context.
type Error struct {
	// Type categorizes the error for handling decisions.
	Type ErrorType

	// Message is the human-readable error message.
	Message string

	// Provider is the name of the provider that generated this error (if applicable).
	Provider string

	// Cause is the underlying error that caused this error.
	Cause error
}

// Error implements the error interface.
func (e *Error) Error() string {
	if e.Provider != "" {
		return fmt.Sprintf("[%s] %s", e.Provider, e.Message)
	}
	return e.Message
}

// Unwrap returns the underlying error for errors.Is/As support.
func (e *Error) Unwrap() error {
	return e.Cause
}

// Is reports whether the target matches this error's type.
func (e *Error) Is(target error) bool {
	if t, ok := target.(*Error); ok {
		return e.Type == t.Type
	}
	return false
}

// IsRetryable returns true if the error might succeed on retry.
func (e *Error) IsRetryable() bool {
	return e.Type == ErrorTypeRetryable
}

// IsValidationError returns true if this is a validation error (bad input).
func (e *Error) IsValidationError() bool {
	return e.Type == ErrorTypeValidation
}

// IsProviderError returns true if this is a provider-specific error.
func (e *Error) IsProviderError() bool {
	return e.Type == ErrorTypeProvider
}

// IsRateLimitError returns true if this is a rate limit error.
func (e *Error) IsRateLimitError() bool {
	return e.Type == ErrorTypeRateLimit
}

// Error constructors

// NewValidationError creates a validation error.
func NewValidationError(message string, cause error) *Error {
	return &Error{
		Type:    ErrorTypeValidation,
		Message: message,
		Cause:   cause,
	}
}

// NewProviderError creates a provider-specific error.
func NewProviderError(provider, message string, cause error) *Error {
	return &Error{
		Type:     ErrorTypeProvider,
		Provider: provider,
		Message:  message,
		Cause:    cause,
	}
}

// NewRetryableError creates a retryable error.
func NewRetryableError(message string, cause error) *Error {
	return &Error{
		Type:    ErrorTypeRetryable,
		Message: message,
		Cause:   cause,
	}
}

// NewRateLimitError creates a rate limit error.
func NewRateLimitError(provider string, cause error) *Error {
	return &Error{
		Type:     ErrorTypeRateLimit,
		Provider: provider,
		Message:  "rate limit exceeded",
		Cause:    cause,
	}
}

// Helper functions for error checking

// IsRetryable checks if an error is retryable.
func IsRetryable(err error) bool {
	var kytErr *Error
	if errors.As(err, &kytErr) {
		return kytErr.IsRetryable()
	}
	return false
}

// IsValidation checks if an error is a validation error.
func IsValidation(err error) bool {
	var kytErr *Error
	if errors.As(err, &kytErr) {
		return kytErr.IsValidationError()
	}
	return false
}

// IsRateLimit checks if an error is a rate limit error.
func IsRateLimit(err error) bool {
	var kytErr *Error
	if errors.As(err, &kytErr) {
		return kytErr.IsRateLimitError()
	}
	return false
}

// GetProvider extracts the provider name from an error if available.
func GetProvider(err error) string {
	var kytErr *Error
	if errors.As(err, &kytErr) {
		return kytErr.Provider
	}
	return ""
}
