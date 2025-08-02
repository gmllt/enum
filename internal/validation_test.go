package internal

import (
	"testing"
)

// TestValidateIndex tests index validation
func TestValidateIndex(t *testing.T) {
	labels := []string{"first", "second", "third"}

	tests := []struct {
		name        string
		index       int
		expectError bool
	}{
		{
			name:        "valid first index",
			index:       0,
			expectError: false,
		},
		{
			name:        "valid middle index",
			index:       1,
			expectError: false,
		},
		{
			name:        "valid last index",
			index:       2,
			expectError: false,
		},
		{
			name:        "negative index",
			index:       -1,
			expectError: true,
		},
		{
			name:        "too high index",
			index:       3,
			expectError: true,
		},
		{
			name:        "way too high index",
			index:       100,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateIndex(labels, tt.index)

			if tt.expectError && err == nil {
				t.Error("expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("expected no error but got: %v", err)
			}
		})
	}
}

// TestIsValidIndex tests boolean index validation
func TestIsValidIndex(t *testing.T) {
	labels := []string{"a", "b"}

	tests := []struct {
		name     string
		index    int
		expected bool
	}{
		{
			name:     "valid first index",
			index:    0,
			expected: true,
		},
		{
			name:     "valid second index",
			index:    1,
			expected: true,
		},
		{
			name:     "negative index",
			index:    -1,
			expected: false,
		},
		{
			name:     "too high index",
			index:    2,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidIndex(labels, tt.index)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestSafeGetLabel tests safe label retrieval
func TestSafeGetLabel(t *testing.T) {
	labels := []string{"red", "green", "blue"}
	defaultLabel := "unknown"

	tests := []struct {
		name     string
		index    int
		expected string
	}{
		{
			name:     "valid first index",
			index:    0,
			expected: "red",
		},
		{
			name:     "valid middle index",
			index:    1,
			expected: "green",
		},
		{
			name:     "valid last index",
			index:    2,
			expected: "blue",
		},
		{
			name:     "negative index",
			index:    -1,
			expected: defaultLabel,
		},
		{
			name:     "too high index",
			index:    3,
			expected: defaultLabel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SafeGetLabel(labels, tt.index, defaultLabel)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// TestSafeGetLabelWithError tests safe label retrieval with error
func TestSafeGetLabelWithError(t *testing.T) {
	labels := []string{"monday", "tuesday"}

	tests := []struct {
		name          string
		index         int
		expectedLabel string
		expectError   bool
	}{
		{
			name:          "valid first index",
			index:         0,
			expectedLabel: "monday",
			expectError:   false,
		},
		{
			name:          "valid second index",
			index:         1,
			expectedLabel: "tuesday",
			expectError:   false,
		},
		{
			name:          "negative index",
			index:         -1,
			expectedLabel: "",
			expectError:   true,
		},
		{
			name:          "too high index",
			index:         2,
			expectedLabel: "",
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			label, err := SafeGetLabelWithError(labels, tt.index)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
				if label != tt.expectedLabel {
					t.Errorf("expected %q, got %q", tt.expectedLabel, label)
				}
			}
		})
	}
}

// TestEmptyLabels tests validation with empty labels
func TestEmptyLabels(t *testing.T) {
	labels := []string{}

	// Test ValidateIndex
	err := ValidateIndex(labels, 0)
	if err == nil {
		t.Error("expected error for index 0 with empty labels")
	}

	// Test IsValidIndex
	if IsValidIndex(labels, 0) {
		t.Error("expected false for index 0 with empty labels")
	}

	// Test SafeGetLabel
	result := SafeGetLabel(labels, 0, "default")
	if result != "default" {
		t.Errorf("expected 'default', got %q", result)
	}

	// Test SafeGetLabelWithError
	_, err = SafeGetLabelWithError(labels, 0)
	if err == nil {
		t.Error("expected error for index 0 with empty labels")
	}
}

// TestValidationConsistency ensures validation functions are consistent
func TestValidationConsistency(t *testing.T) {
	labels := []string{"x", "y", "z"}

	for i := -2; i <= 5; i++ {
		// ValidateIndex and IsValidIndex should be consistent
		err := ValidateIndex(labels, i)
		isValid := IsValidIndex(labels, i)

		if (err == nil) != isValid {
			t.Errorf("inconsistent validation for index %d: ValidateIndex error=%v, IsValidIndex=%v",
				i, err, isValid)
		}

		// SafeGetLabelWithError should be consistent with validation
		_, err2 := SafeGetLabelWithError(labels, i)
		if (err == nil) != (err2 == nil) {
			t.Errorf("inconsistent error handling for index %d: ValidateIndex error=%v, SafeGetLabelWithError error=%v",
				i, err, err2)
		}

		// SafeGetLabel should return the actual label for valid indices
		if isValid {
			safeLabel := SafeGetLabel(labels, i, "default")
			expectedLabel := labels[i]
			if safeLabel != expectedLabel {
				t.Errorf("expected %q for valid index %d, got %q", expectedLabel, i, safeLabel)
			}
		}
	}
}

// TestCustomTypes tests validation with custom integer types
func TestCustomTypes(t *testing.T) {
	type CustomInt int

	labels := []string{"custom1", "custom2"}

	// Test with custom type
	err := ValidateIndex(labels, CustomInt(0))
	if err != nil {
		t.Errorf("expected no error for valid custom type index, got: %v", err)
	}

	err = ValidateIndex(labels, CustomInt(-1))
	if err == nil {
		t.Error("expected error for invalid custom type index")
	}

	// Test SafeGetLabel with custom type
	result := SafeGetLabel(labels, CustomInt(1), "default")
	if result != "custom2" {
		t.Errorf("expected 'custom2', got %q", result)
	}
}
