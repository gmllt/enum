package enum

import (
	"fmt"
	"reflect"
	"testing"
)

// TestNewEnum tests the creation of new enum instances
func TestNewEnum(t *testing.T) {
	tests := []struct {
		name   string
		labels []string
	}{
		{
			name:   "empty enum",
			labels: []string{},
		},
		{
			name:   "single label",
			labels: []string{"first"},
		},
		{
			name:   "multiple labels",
			labels: []string{"first", "second", "third"},
		},
		{
			name:   "large enum",
			labels: []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enum := NewEnum[int](tt.labels...)

			if enum == nil {
				t.Fatal("NewEnum returned nil")
			}

			if len(enum.labels) != len(tt.labels) {
				t.Errorf("expected %d labels, got %d", len(tt.labels), len(enum.labels))
			}

			if len(enum.labelMap) != len(tt.labels) {
				t.Errorf("expected labelMap size %d, got %d", len(tt.labels), len(enum.labelMap))
			}

			if len(enum.allVals) != len(tt.labels) {
				t.Errorf("expected allVals size %d, got %d", len(tt.labels), len(enum.allVals))
			}
		})
	}
}

// TestEnumString tests the string representation of enum values
func TestEnumString(t *testing.T) {
	enum := NewEnum[int]("zero", "one", "two")

	tests := []struct {
		name     string
		value    int
		expected string
	}{
		{
			name:     "valid first value",
			value:    0,
			expected: "zero",
		},
		{
			name:     "valid middle value",
			value:    1,
			expected: "one",
		},
		{
			name:     "valid last value",
			value:    2,
			expected: "two",
		},
		{
			name:     "invalid negative value",
			value:    -1,
			expected: "Invalid(-1)",
		},
		{
			name:     "invalid high value",
			value:    5,
			expected: "Invalid(5)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := enum.String(tt.value)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// TestEnumFromString tests string to enum value conversion
func TestEnumFromString(t *testing.T) {
	enum := NewEnum[int]("alpha", "beta", "gamma")

	tests := []struct {
		name        string
		input       string
		expectedVal int
		expectError bool
	}{
		{
			name:        "valid first label",
			input:       "alpha",
			expectedVal: 0,
			expectError: false,
		},
		{
			name:        "valid middle label",
			input:       "beta",
			expectedVal: 1,
			expectError: false,
		},
		{
			name:        "valid last label",
			input:       "gamma",
			expectedVal: 2,
			expectError: false,
		},
		{
			name:        "invalid label",
			input:       "delta",
			expectedVal: 0,
			expectError: true,
		},
		{
			name:        "empty string",
			input:       "",
			expectedVal: 0,
			expectError: true,
		},
		{
			name:        "case sensitive",
			input:       "Alpha",
			expectedVal: 0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := enum.FromString(tt.input)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
				if result != tt.expectedVal {
					t.Errorf("expected %d, got %d", tt.expectedVal, result)
				}
			}
		})
	}
}

// TestEnumAll tests getting all enum values
func TestEnumAll(t *testing.T) {
	tests := []struct {
		name     string
		labels   []string
		expected []int
	}{
		{
			name:     "empty enum",
			labels:   []string{},
			expected: []int{},
		},
		{
			name:     "single value",
			labels:   []string{"first"},
			expected: []int{0},
		},
		{
			name:     "multiple values",
			labels:   []string{"a", "b", "c"},
			expected: []int{0, 1, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enum := NewEnum[int](tt.labels...)
			result := enum.All()

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}

			// Test that modifying the returned slice doesn't affect the enum
			if len(result) > 0 {
				original := result[0]
				result[0] = 999
				newResult := enum.All()
				if len(newResult) > 0 && newResult[0] != original {
					t.Error("modifying returned slice affected the enum internal state")
				}
			}
		})
	}
}

// TestEnumLabels tests getting all enum labels
func TestEnumLabels(t *testing.T) {
	labels := []string{"red", "green", "blue"}
	enum := NewEnum[int](labels...)

	result := enum.Labels()

	if !reflect.DeepEqual(result, labels) {
		t.Errorf("expected %v, got %v", labels, result)
	}

	// Test that modifying the returned slice doesn't affect the enum
	result[0] = "modified"
	newResult := enum.Labels()
	if newResult[0] != "red" {
		t.Error("modifying returned slice affected the enum internal state")
	}
}

// TestEnumLabelsReadOnly tests read-only access to labels
func TestEnumLabelsReadOnly(t *testing.T) {
	labels := []string{"monday", "tuesday", "wednesday"}
	enum := NewEnum[int](labels...)

	result := enum.LabelsReadOnly()

	if !reflect.DeepEqual(result, labels) {
		t.Errorf("expected %v, got %v", labels, result)
	}

	// Test that the returned slice shares memory (this is expected behavior)
	if &result[0] != &enum.labels[0] {
		t.Error("LabelsReadOnly should return a slice that shares memory with internal labels")
	}
}

// TestEnumPerformance tests performance characteristics
func TestEnumPerformance(t *testing.T) {
	// Create a large enum to test map-based lookup
	largeLabels := make([]string, 100)
	for i := 0; i < 100; i++ {
		largeLabels[i] = fmt.Sprintf("label_%d", i)
	}

	largeEnum := NewEnum[int](largeLabels...)

	// Test that lookups work correctly even with large enums
	val, err := largeEnum.FromString("label_50")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if val != 50 {
		t.Errorf("expected 50, got %d", val)
	}

	// Test invalid lookup
	_, err = largeEnum.FromString("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent label")
	}
}

// TestEnumWithDifferentValues tests enum with different integer values
func TestEnumWithDifferentValues(t *testing.T) {
	// Test with custom type based on int
	type CustomInt int

	enum := NewEnum[CustomInt]("first", "second")
	val, err := enum.FromString("second")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if val != 1 {
		t.Errorf("expected 1, got %d", val)
	}

	// Test string representation
	str := enum.String(0)
	if str != "first" {
		t.Errorf("expected 'first', got %q", str)
	}
}
