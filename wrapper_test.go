package enum

import (
	"encoding/json"
	"reflect"
	"testing"
)

// TestNewWrapper tests wrapper creation
func TestNewWrapper(t *testing.T) {
	labels := []string{"first", "second", "third"}
	wrapper := NewWrapper[int](labels...)

	if wrapper.Enum == nil {
		t.Fatal("NewWrapper created wrapper with nil Enum")
	}

	if !reflect.DeepEqual(wrapper.Enum.labels, labels) {
		t.Errorf("expected labels %v, got %v", labels, wrapper.Enum.labels)
	}

	// Default value should be zero
	if wrapper.Value != 0 {
		t.Errorf("expected default value 0, got %d", wrapper.Value)
	}
}

// TestWrapperString tests string representation
func TestWrapperString(t *testing.T) {
	wrapper := NewWrapper[int]("alpha", "beta", "gamma")

	tests := []struct {
		name     string
		value    int
		expected string
	}{
		{
			name:     "first value",
			value:    0,
			expected: "alpha",
		},
		{
			name:     "middle value",
			value:    1,
			expected: "beta",
		},
		{
			name:     "last value",
			value:    2,
			expected: "gamma",
		},
		{
			name:     "invalid value",
			value:    5,
			expected: "Invalid(5)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrapper.Value = tt.value
			result := wrapper.String()

			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// TestWrapperAll tests getting all values
func TestWrapperAll(t *testing.T) {
	wrapper := NewWrapper[int]("red", "green", "blue")
	result := wrapper.All()
	expected := []int{0, 1, 2}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

// TestWrapperLabels tests getting all labels
func TestWrapperLabels(t *testing.T) {
	labels := []string{"monday", "tuesday", "wednesday"}
	wrapper := NewWrapper[int](labels...)
	result := wrapper.Labels()

	if !reflect.DeepEqual(result, labels) {
		t.Errorf("expected %v, got %v", labels, result)
	}

	// Test that modifying returned slice doesn't affect wrapper
	result[0] = "modified"
	newResult := wrapper.Labels()
	if newResult[0] != "monday" {
		t.Error("modifying returned slice affected wrapper internal state")
	}
}

// TestWrapperGetSet tests getter and setter
func TestWrapperGetSet(t *testing.T) {
	wrapper := NewWrapper[int]("one", "two", "three")

	// Test initial value
	if wrapper.Get() != 0 {
		t.Errorf("expected initial value 0, got %d", wrapper.Get())
	}

	// Test setting value
	wrapper.Set(2)
	if wrapper.Get() != 2 {
		t.Errorf("expected value 2, got %d", wrapper.Get())
	}

	// Test that internal Value field is also updated
	if wrapper.Value != 2 {
		t.Errorf("expected internal Value 2, got %d", wrapper.Value)
	}
}

// TestWrapperJSONMarshal tests JSON marshaling
func TestWrapperJSONMarshal(t *testing.T) {
	wrapper := NewWrapper[int]("spring", "summer", "autumn", "winter")

	tests := []struct {
		name     string
		value    int
		expected string
	}{
		{
			name:     "first season",
			value:    0,
			expected: `"spring"`,
		},
		{
			name:     "middle season",
			value:    2,
			expected: `"autumn"`,
		},
		{
			name:     "last season",
			value:    3,
			expected: `"winter"`,
		},
		{
			name:     "invalid season",
			value:    10,
			expected: `"Invalid"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrapper.Set(tt.value)
			result, err := wrapper.MarshalJSON()
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if string(result) != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, string(result))
			}
		})
	}
}

// TestWrapperJSONUnmarshal tests JSON unmarshaling
func TestWrapperJSONUnmarshal(t *testing.T) {
	wrapper := NewWrapper[int]("dog", "cat", "bird", "fish")

	tests := []struct {
		name        string
		input       string
		expectedVal int
		expectError bool
	}{
		{
			name:        "valid first animal",
			input:       `"dog"`,
			expectedVal: 0,
			expectError: false,
		},
		{
			name:        "valid middle animal",
			input:       `"bird"`,
			expectedVal: 2,
			expectError: false,
		},
		{
			name:        "valid last animal",
			input:       `"fish"`,
			expectedVal: 3,
			expectError: false,
		},
		{
			name:        "invalid animal",
			input:       `"elephant"`,
			expectedVal: 0,
			expectError: false, // Returns zero value, no error
		},
		{
			name:        "invalid JSON",
			input:       `invalid json`,
			expectedVal: 0,
			expectError: true,
		},
		{
			name:        "non-string JSON",
			input:       `123`,
			expectedVal: 0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset wrapper value
			wrapper.Set(99)

			err := wrapper.UnmarshalJSON([]byte(tt.input))

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}

				if wrapper.Get() != tt.expectedVal {
					t.Errorf("expected value %d, got %d", tt.expectedVal, wrapper.Get())
				}
			}
		})
	}
}

// TestWrapperYAMLMarshal tests YAML marshaling
func TestWrapperYAMLMarshal(t *testing.T) {
	wrapper := NewWrapper[int]("north", "south", "east", "west")

	tests := []struct {
		name     string
		value    int
		expected any
	}{
		{
			name:     "first direction",
			value:    0,
			expected: "north",
		},
		{
			name:     "middle direction",
			value:    2,
			expected: "east",
		},
		{
			name:     "last direction",
			value:    3,
			expected: "west",
		},
		{
			name:     "invalid direction",
			value:    10,
			expected: "Invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrapper.Set(tt.value)
			result, err := wrapper.MarshalYAML()
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestWrapperYAMLUnmarshal tests YAML unmarshaling
func TestWrapperYAMLUnmarshal(t *testing.T) {
	wrapper := NewWrapper[int]("small", "medium", "large")

	tests := []struct {
		name        string
		input       any
		expectedVal int
		expectError bool
	}{
		{
			name:        "valid first size",
			input:       "small",
			expectedVal: 0,
			expectError: false,
		},
		{
			name:        "valid middle size",
			input:       "medium",
			expectedVal: 1,
			expectError: false,
		},
		{
			name:        "valid last size",
			input:       "large",
			expectedVal: 2,
			expectError: false,
		},
		{
			name:        "invalid size",
			input:       "huge",
			expectedVal: 0,
			expectError: false, // Returns zero value, no error
		},
		{
			name:        "non-string input",
			input:       123,
			expectedVal: 0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset wrapper value
			wrapper.Set(99)

			// Create unmarshal function
			unmarshal := func(v any) error {
				if ptr, ok := v.(*string); ok {
					if str, ok := tt.input.(string); ok {
						*ptr = str
						return nil
					}
					return &json.UnsupportedTypeError{Type: reflect.TypeOf(tt.input)}
				}
				return &json.UnsupportedTypeError{Type: reflect.TypeOf(v)}
			}

			err := wrapper.UnmarshalYAML(unmarshal)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}

				if wrapper.Get() != tt.expectedVal {
					t.Errorf("expected value %d, got %d", tt.expectedVal, wrapper.Get())
				}
			}
		})
	}
}

// TestWrapperJSONRoundTrip tests JSON serialization and deserialization consistency
func TestWrapperJSONRoundTrip(t *testing.T) {
	wrapper := NewWrapper[int]("red", "green", "blue", "yellow", "orange")

	for i := 0; i < 5; i++ {
		// Set value
		wrapper.Set(i)

		// Marshal to JSON
		jsonBytes, err := wrapper.MarshalJSON()
		if err != nil {
			t.Fatalf("MarshalJSON failed for value %d: %v", i, err)
		}

		// Create new wrapper and unmarshal
		newWrapper := NewWrapper[int]("red", "green", "blue", "yellow", "orange")
		err = newWrapper.UnmarshalJSON(jsonBytes)
		if err != nil {
			t.Fatalf("UnmarshalJSON failed for value %d: %v", i, err)
		}

		// Check that values match
		if newWrapper.Get() != i {
			t.Errorf("round trip failed: expected %d, got %d", i, newWrapper.Get())
		}
	}
}

// TestWrapperWithCustomTypes tests wrapper with custom integer types
func TestWrapperWithCustomTypes(t *testing.T) {
	type CustomInt int

	wrapper := NewWrapper[CustomInt]("first", "second")

	// Test setting and getting
	wrapper.Set(CustomInt(1))
	if wrapper.Get() != CustomInt(1) {
		t.Errorf("expected CustomInt(1), got %v", wrapper.Get())
	}

	// Test string representation
	result := wrapper.String()
	if result != "second" {
		t.Errorf("expected 'second', got %q", result)
	}

	// Test JSON marshaling
	jsonBytes, err := wrapper.MarshalJSON()
	if err != nil {
		t.Errorf("JSON marshal failed: %v", err)
	}

	if string(jsonBytes) != `"second"` {
		t.Errorf("expected '\"second\"', got %s", string(jsonBytes))
	}
}

// TestWrapperConcurrency tests that wrapper operations are safe for concurrent read access
func TestWrapperConcurrency(t *testing.T) {
	wrapper := NewWrapper[int]("concurrent1", "concurrent2", "concurrent3")
	wrapper.Set(1)

	// Test concurrent reads (these should be safe)
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func() {
			defer func() { done <- true }()

			// These operations should be safe for concurrent read access
			_ = wrapper.String()
			_ = wrapper.Get()
			_ = wrapper.All()
			_ = wrapper.Labels()
			_, _ = wrapper.MarshalJSON()
			_, _ = wrapper.MarshalYAML()
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify wrapper state is still consistent
	if wrapper.Get() != 1 {
		t.Errorf("expected value 1 after concurrent operations, got %d", wrapper.Get())
	}
}
