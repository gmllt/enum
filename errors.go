package enum

import (
	"github.com/gmllt/enum/internal"
)

// Re-export internal error types for public use.
// This allows users of the library to check for specific error types
// using errors.Is() or errors.As().

// ErrInvalidEnumValue is returned when trying to unmarshal an invalid enum value.
// It provides information about the invalid value and the list of valid values.
type ErrInvalidEnumValue = internal.ErrInvalidEnumValue

// ErrBinaryDataTooShort is returned when binary data is too short to contain valid data.
type ErrBinaryDataTooShort = internal.ErrBinaryDataTooShort

// ErrBinaryDataTruncated is returned when binary data is truncated.
type ErrBinaryDataTruncated = internal.ErrBinaryDataTruncated

// ErrLabelTooLong is returned when a label exceeds the maximum allowed length for binary encoding.
type ErrLabelTooLong = internal.ErrLabelTooLong

// Helper functions for creating error instances (optional, for convenience).

// NewInvalidEnumValueError creates a new ErrInvalidEnumValue.
// This is provided for convenience, though users typically won't need to create these errors manually.
func NewInvalidEnumValueError(value string, validValues []string) *ErrInvalidEnumValue {
	return internal.NewInvalidEnumValueError(value, validValues)
}

// NewBinaryDataTooShortError creates a new ErrBinaryDataTooShort.
func NewBinaryDataTooShortError(expected, actual int) *ErrBinaryDataTooShort {
	return internal.NewBinaryDataTooShortError(expected, actual)
}

// NewBinaryDataTruncatedError creates a new ErrBinaryDataTruncated.
func NewBinaryDataTruncatedError(expected, actual int) *ErrBinaryDataTruncated {
	return internal.NewBinaryDataTruncatedError(expected, actual)
}

// NewLabelTooLongError creates a new ErrLabelTooLong.
func NewLabelTooLongError(length, maxLength int) *ErrLabelTooLong {
	return internal.NewLabelTooLongError(length, maxLength)
}
