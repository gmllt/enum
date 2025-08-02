package internal

import (
	"fmt"
)

// ErrInvalidEnumValue is returned when trying to unmarshal an invalid enum value.
type ErrInvalidEnumValue struct {
	Value       string
	ValidValues []string
}

func (e *ErrInvalidEnumValue) Error() string {
	if len(e.ValidValues) == 0 {
		return fmt.Sprintf("invalid enum value: %q (no valid values available)", e.Value)
	}
	return fmt.Sprintf("invalid enum value: %q (valid values: %v)", e.Value, e.ValidValues)
}

// Is implements the errors.Is interface for error comparison.
func (e *ErrInvalidEnumValue) Is(target error) bool {
	_, ok := target.(*ErrInvalidEnumValue)
	return ok
}

// ErrBinaryDataTooShort is returned when binary data is too short to contain valid data.
type ErrBinaryDataTooShort struct {
	Expected int
	Actual   int
}

func (e *ErrBinaryDataTooShort) Error() string {
	return fmt.Sprintf("binary data too short: expected at least %d bytes, got %d", e.Expected, e.Actual)
}

// Is implements the errors.Is interface for error comparison.
func (e *ErrBinaryDataTooShort) Is(target error) bool {
	_, ok := target.(*ErrBinaryDataTooShort)
	return ok
}

// ErrBinaryDataTruncated is returned when binary data is truncated.
type ErrBinaryDataTruncated struct {
	Expected int
	Actual   int
}

func (e *ErrBinaryDataTruncated) Error() string {
	return fmt.Sprintf("binary data truncated: expected %d bytes, got %d", e.Expected, e.Actual)
}

// Is implements the errors.Is interface for error comparison.
func (e *ErrBinaryDataTruncated) Is(target error) bool {
	_, ok := target.(*ErrBinaryDataTruncated)
	return ok
}

// ErrLabelTooLong is returned when a label exceeds the maximum allowed length for binary encoding.
type ErrLabelTooLong struct {
	Length    int
	MaxLength int
}

func (e *ErrLabelTooLong) Error() string {
	return fmt.Sprintf("label too long: %d bytes (max %d)", e.Length, e.MaxLength)
}

// Is implements the errors.Is interface for error comparison.
func (e *ErrLabelTooLong) Is(target error) bool {
	_, ok := target.(*ErrLabelTooLong)
	return ok
}

// NewInvalidEnumValueError creates a new ErrInvalidEnumValue.
func NewInvalidEnumValueError(value string, validValues []string) *ErrInvalidEnumValue {
	// Create a copy of validValues to avoid external modifications
	validValuesCopy := make([]string, len(validValues))
	copy(validValuesCopy, validValues)

	return &ErrInvalidEnumValue{
		Value:       value,
		ValidValues: validValuesCopy,
	}
}

// NewBinaryDataTooShortError creates a new ErrBinaryDataTooShort.
func NewBinaryDataTooShortError(expected, actual int) *ErrBinaryDataTooShort {
	return &ErrBinaryDataTooShort{
		Expected: expected,
		Actual:   actual,
	}
}

// NewBinaryDataTruncatedError creates a new ErrBinaryDataTruncated.
func NewBinaryDataTruncatedError(expected, actual int) *ErrBinaryDataTruncated {
	return &ErrBinaryDataTruncated{
		Expected: expected,
		Actual:   actual,
	}
}

// NewLabelTooLongError creates a new ErrLabelTooLong.
func NewLabelTooLongError(length, maxLength int) *ErrLabelTooLong {
	return &ErrLabelTooLong{
		Length:    length,
		MaxLength: maxLength,
	}
}
