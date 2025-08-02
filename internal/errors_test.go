package internal

import (
	"errors"
	"testing"
)

// TestErrInvalidEnumValue tests the ErrInvalidEnumValue error type.
func TestErrInvalidEnumValue(t *testing.T) {
	validValues := []string{"red", "green", "blue"}
	err := NewInvalidEnumValueError("yellow", validValues)

	// Test error message
	expectedMsg := `invalid enum value: "yellow" (valid values: [red green blue])`
	if err.Error() != expectedMsg {
		t.Errorf("expected error message %q, got %q", expectedMsg, err.Error())
	}

	// Test Is method
	var target *ErrInvalidEnumValue
	if !errors.As(err, &target) {
		t.Error("expected error to be of type *ErrInvalidEnumValue")
	}

	// Test that modifying original slice doesn't affect error
	validValues[0] = "modified"
	if err.ValidValues[0] != "red" {
		t.Error("error's ValidValues should not be affected by external modifications")
	}
}

// TestErrInvalidEnumValueEmpty tests ErrInvalidEnumValue with empty valid values.
func TestErrInvalidEnumValueEmpty(t *testing.T) {
	err := NewInvalidEnumValueError("test", []string{})

	expectedMsg := `invalid enum value: "test" (no valid values available)`
	if err.Error() != expectedMsg {
		t.Errorf("expected error message %q, got %q", expectedMsg, err.Error())
	}
}

// TestErrBinaryDataTooShort tests the ErrBinaryDataTooShort error type.
func TestErrBinaryDataTooShort(t *testing.T) {
	err := NewBinaryDataTooShortError(4, 2)

	expectedMsg := "binary data too short: expected at least 4 bytes, got 2"
	if err.Error() != expectedMsg {
		t.Errorf("expected error message %q, got %q", expectedMsg, err.Error())
	}

	// Test Is method
	var target *ErrBinaryDataTooShort
	if !errors.As(err, &target) {
		t.Error("expected error to be of type *ErrBinaryDataTooShort")
	}

	if target.Expected != 4 || target.Actual != 2 {
		t.Errorf("expected Expected=4, Actual=2, got Expected=%d, Actual=%d", target.Expected, target.Actual)
	}
}

// TestErrBinaryDataTruncated tests the ErrBinaryDataTruncated error type.
func TestErrBinaryDataTruncated(t *testing.T) {
	err := NewBinaryDataTruncatedError(10, 5)

	expectedMsg := "binary data truncated: expected 10 bytes, got 5"
	if err.Error() != expectedMsg {
		t.Errorf("expected error message %q, got %q", expectedMsg, err.Error())
	}

	// Test Is method
	var target *ErrBinaryDataTruncated
	if !errors.As(err, &target) {
		t.Error("expected error to be of type *ErrBinaryDataTruncated")
	}

	if target.Expected != 10 || target.Actual != 5 {
		t.Errorf("expected Expected=10, Actual=5, got Expected=%d, Actual=%d", target.Expected, target.Actual)
	}
}

// TestErrLabelTooLong tests the ErrLabelTooLong error type.
func TestErrLabelTooLong(t *testing.T) {
	err := NewLabelTooLongError(70000, 65535)

	expectedMsg := "label too long: 70000 bytes (max 65535)"
	if err.Error() != expectedMsg {
		t.Errorf("expected error message %q, got %q", expectedMsg, err.Error())
	}

	// Test Is method
	var target *ErrLabelTooLong
	if !errors.As(err, &target) {
		t.Error("expected error to be of type *ErrLabelTooLong")
	}

	if target.Length != 70000 || target.MaxLength != 65535 {
		t.Errorf("expected Length=70000, MaxLength=65535, got Length=%d, MaxLength=%d", target.Length, target.MaxLength)
	}
}

// TestErrorTypeComparison tests that different error types can be distinguished.
func TestErrorTypeComparison(t *testing.T) {
	invalidValueErr := NewInvalidEnumValueError("test", []string{"valid"})
	binaryShortErr := NewBinaryDataTooShortError(4, 2)
	binaryTruncatedErr := NewBinaryDataTruncatedError(10, 5)
	labelTooLongErr := NewLabelTooLongError(70000, 65535)

	// Test that each error type is distinct
	if errors.Is(invalidValueErr, binaryShortErr) {
		t.Error("ErrInvalidEnumValue should not be equal to ErrBinaryDataTooShort")
	}

	if errors.Is(binaryShortErr, binaryTruncatedErr) {
		t.Error("ErrBinaryDataTooShort should not be equal to ErrBinaryDataTruncated")
	}

	if errors.Is(binaryTruncatedErr, labelTooLongErr) {
		t.Error("ErrBinaryDataTruncated should not be equal to ErrLabelTooLong")
	}

	// Test that same error types are equal
	anotherInvalidValueErr := NewInvalidEnumValueError("other", []string{"different"})
	if !errors.Is(invalidValueErr, anotherInvalidValueErr) {
		t.Error("Two ErrInvalidEnumValue instances should be considered equal via errors.Is")
	}

	// Test that same error types are equal for ErrLabelTooLong
	anotherLabelTooLongErr := NewLabelTooLongError(50000, 32767)
	if !errors.Is(labelTooLongErr, anotherLabelTooLongErr) {
		t.Error("Two ErrLabelTooLong instances should be considered equal via errors.Is")
	}
}
