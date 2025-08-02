package enum

import (
	"errors"
	"testing"
)

// TestPublicErrorTypes tests that the public error types work correctly.
func TestPublicErrorTypes(t *testing.T) {
	// Test that we can create and use public error types
	invalidErr := NewInvalidEnumValueError("invalid", []string{"valid1", "valid2"})
	binaryShortErr := NewBinaryDataTooShortError(4, 2)
	binaryTruncatedErr := NewBinaryDataTruncatedError(10, 5)
	labelTooLongErr := NewLabelTooLongError(70000, 65535)

	// Test that errors.Is works with public types
	var targetInvalid *ErrInvalidEnumValue
	if !errors.As(invalidErr, &targetInvalid) {
		t.Error("should be able to cast to public ErrInvalidEnumValue type")
	}

	var targetBinaryShort *ErrBinaryDataTooShort
	if !errors.As(binaryShortErr, &targetBinaryShort) {
		t.Error("should be able to cast to public ErrBinaryDataTooShort type")
	}

	var targetBinaryTruncated *ErrBinaryDataTruncated
	if !errors.As(binaryTruncatedErr, &targetBinaryTruncated) {
		t.Error("should be able to cast to public ErrBinaryDataTruncated type")
	}

	var targetLabelTooLong *ErrLabelTooLong
	if !errors.As(labelTooLongErr, &targetLabelTooLong) {
		t.Error("should be able to cast to public ErrLabelTooLong type")
	}

	// Test error messages
	expectedInvalidMsg := `invalid enum value: "invalid" (valid values: [valid1 valid2])`
	if invalidErr.Error() != expectedInvalidMsg {
		t.Errorf("expected %q, got %q", expectedInvalidMsg, invalidErr.Error())
	}

	expectedBinaryShortMsg := "binary data too short: expected at least 4 bytes, got 2"
	if binaryShortErr.Error() != expectedBinaryShortMsg {
		t.Errorf("expected %q, got %q", expectedBinaryShortMsg, binaryShortErr.Error())
	}
}

// TestErrorIntegrationWithWrapper tests that wrapper methods return the correct error types.
func TestErrorIntegrationWithWrapper(t *testing.T) {
	wrapper := NewWrapper[int]("red", "green", "blue")

	// Test JSON unmarshalling with invalid value
	err := wrapper.UnmarshalJSON([]byte(`"yellow"`))
	if err == nil {
		t.Fatal("expected error for invalid JSON value")
	}

	var invalidErr *ErrInvalidEnumValue
	if !errors.As(err, &invalidErr) {
		t.Errorf("expected ErrInvalidEnumValue, got %T: %v", err, err)
	} else {
		if invalidErr.Value != "yellow" {
			t.Errorf("expected invalid value 'yellow', got %q", invalidErr.Value)
		}
		if len(invalidErr.ValidValues) != 3 {
			t.Errorf("expected 3 valid values, got %d", len(invalidErr.ValidValues))
		}
	}

	// Test binary unmarshalling with too short data
	err = wrapper.UnmarshalBinary([]byte{})
	if err == nil {
		t.Fatal("expected error for empty binary data")
	}

	var binaryShortErr *ErrBinaryDataTooShort
	if !errors.As(err, &binaryShortErr) {
		t.Errorf("expected ErrBinaryDataTooShort, got %T: %v", err, err)
	} else {
		if binaryShortErr.Expected != 2 || binaryShortErr.Actual != 0 {
			t.Errorf("expected Expected=2, Actual=0, got Expected=%d, Actual=%d",
				binaryShortErr.Expected, binaryShortErr.Actual)
		}
	}

	// Test binary unmarshalling with truncated data
	err = wrapper.UnmarshalBinary([]byte{0, 10, 'a', 'b'}) // claims 10 bytes but only has 2
	if err == nil {
		t.Fatal("expected error for truncated binary data")
	}

	var binaryTruncatedErr *ErrBinaryDataTruncated
	if !errors.As(err, &binaryTruncatedErr) {
		t.Errorf("expected ErrBinaryDataTruncated, got %T: %v", err, err)
	} else {
		if binaryTruncatedErr.Expected != 12 || binaryTruncatedErr.Actual != 4 {
			t.Errorf("expected Expected=12, Actual=4, got Expected=%d, Actual=%d",
				binaryTruncatedErr.Expected, binaryTruncatedErr.Actual)
		}
	}

	// Test text unmarshalling with invalid value
	err = wrapper.UnmarshalText([]byte("purple"))
	if err == nil {
		t.Fatal("expected error for invalid text value")
	}

	var textInvalidErr *ErrInvalidEnumValue
	if !errors.As(err, &textInvalidErr) {
		t.Errorf("expected ErrInvalidEnumValue for text, got %T: %v", err, err)
	} else {
		if textInvalidErr.Value != "purple" {
			t.Errorf("expected invalid value 'purple', got %q", textInvalidErr.Value)
		}
	}
}
