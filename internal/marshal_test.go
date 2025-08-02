package internal

import (
	"encoding/json"
	"reflect"
	"testing"
)

// TestToJSON tests JSON serialization
func TestToJSON(t *testing.T) {
	labels := []string{"first", "second", "third"}

	tests := []struct {
		name     string
		value    int
		expected string
	}{
		{
			name:     "valid first value",
			value:    0,
			expected: `"first"`,
		},
		{
			name:     "valid middle value",
			value:    1,
			expected: `"second"`,
		},
		{
			name:     "valid last value",
			value:    2,
			expected: `"third"`,
		},
		{
			name:     "invalid negative value",
			value:    -1,
			expected: `"Invalid"`,
		},
		{
			name:     "invalid high value",
			value:    5,
			expected: `"Invalid"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ToJSON(labels, tt.value)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if string(result) != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, string(result))
			}
		})
	}
}

// TestFromJSON tests JSON deserialization
func TestFromJSON(t *testing.T) {
	labels := []string{"red", "green", "blue"}

	tests := []struct {
		name        string
		input       string
		expectedVal int
		expectError bool
	}{
		{
			name:        "valid first label",
			input:       `"red"`,
			expectedVal: 0,
			expectError: false,
		},
		{
			name:        "valid middle label",
			input:       `"green"`,
			expectedVal: 1,
			expectError: false,
		},
		{
			name:        "valid last label",
			input:       `"blue"`,
			expectedVal: 2,
			expectError: false,
		},
		{
			name:        "invalid label",
			input:       `"yellow"`,
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
		{
			name:        "empty string",
			input:       `""`,
			expectedVal: 0,
			expectError: false, // Returns zero value for not found
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := FromJSON[int](labels, []byte(tt.input))

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

// TestToYAML tests YAML serialization
func TestToYAML(t *testing.T) {
	labels := []string{"alpha", "beta", "gamma"}

	tests := []struct {
		name     string
		value    int
		expected any
	}{
		{
			name:     "valid first value",
			value:    0,
			expected: "alpha",
		},
		{
			name:     "valid middle value",
			value:    1,
			expected: "beta",
		},
		{
			name:     "valid last value",
			value:    2,
			expected: "gamma",
		},
		{
			name:     "invalid negative value",
			value:    -1,
			expected: "Invalid",
		},
		{
			name:     "invalid high value",
			value:    10,
			expected: "Invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ToYAML(labels, tt.value)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestFromYAML tests YAML deserialization
func TestFromYAML(t *testing.T) {
	labels := []string{"dog", "cat", "bird"}

	tests := []struct {
		name        string
		input       any
		expectedVal int
		expectError bool
	}{
		{
			name:        "valid first label",
			input:       "dog",
			expectedVal: 0,
			expectError: false,
		},
		{
			name:        "valid middle label",
			input:       "cat",
			expectedVal: 1,
			expectError: false,
		},
		{
			name:        "valid last label",
			input:       "bird",
			expectedVal: 2,
			expectError: false,
		},
		{
			name:        "invalid label",
			input:       "fish",
			expectedVal: 0,
			expectError: false, // Returns zero value, no error
		},
		{
			name:        "non-string input",
			input:       123,
			expectedVal: 0,
			expectError: true,
		},
		{
			name:        "empty string",
			input:       "",
			expectedVal: 0,
			expectError: false, // Returns zero value for not found
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create unmarshal function that sets the input value
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

			result, err := FromYAML[int](labels, unmarshal)

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

// TestJSONRoundTrip tests JSON serialization and deserialization consistency
func TestJSONRoundTrip(t *testing.T) {
	labels := []string{"monday", "tuesday", "wednesday", "thursday", "friday"}

	for i, expectedLabel := range labels {
		// Serialize
		jsonBytes, err := ToJSON(labels, i)
		if err != nil {
			t.Fatalf("ToJSON failed for index %d: %v", i, err)
		}

		// Deserialize
		result, err := FromJSON[int](labels, jsonBytes)
		if err != nil {
			t.Fatalf("FromJSON failed for index %d: %v", i, err)
		}

		// Check that we got back the original value
		if result != i {
			t.Errorf("round trip failed: expected %d, got %d", i, result)
		}

		// Verify the JSON contains the expected label
		var labelFromJSON string
		if err := json.Unmarshal(jsonBytes, &labelFromJSON); err != nil {
			t.Fatalf("failed to unmarshal JSON: %v", err)
		}

		if labelFromJSON != expectedLabel {
			t.Errorf("expected label %q in JSON, got %q", expectedLabel, labelFromJSON)
		}
	}
}

// TestYAMLRoundTrip tests YAML serialization and deserialization consistency
func TestYAMLRoundTrip(t *testing.T) {
	labels := []string{"spring", "summer", "fall", "winter"}

	for i, expectedLabel := range labels {
		// Serialize
		yamlValue, err := ToYAML(labels, i)
		if err != nil {
			t.Fatalf("ToYAML failed for index %d: %v", i, err)
		}

		// Check that YAML value is correct
		if yamlValue != expectedLabel {
			t.Errorf("expected YAML value %q, got %v", expectedLabel, yamlValue)
		}

		// Deserialize using unmarshal function
		unmarshal := func(v any) error {
			if ptr, ok := v.(*string); ok {
				*ptr = expectedLabel
				return nil
			}
			return &json.UnsupportedTypeError{Type: reflect.TypeOf(v)}
		}

		result, err := FromYAML[int](labels, unmarshal)
		if err != nil {
			t.Fatalf("FromYAML failed for index %d: %v", i, err)
		}

		// Check that we got back the original value
		if result != i {
			t.Errorf("round trip failed: expected %d, got %d", i, result)
		}
	}
}

// TestMarshalWithCustomTypes tests marshaling with custom integer types
func TestMarshalWithCustomTypes(t *testing.T) {
	type CustomInt int

	labels := []string{"custom1", "custom2"}

	// Test JSON
	jsonBytes, err := ToJSON(labels, CustomInt(1))
	if err != nil {
		t.Errorf("ToJSON failed with custom type: %v", err)
	}

	result, err := FromJSON[CustomInt](labels, jsonBytes)
	if err != nil {
		t.Errorf("FromJSON failed with custom type: %v", err)
	}

	if result != CustomInt(1) {
		t.Errorf("expected CustomInt(1), got %v", result)
	}

	// Test YAML
	yamlValue, err := ToYAML(labels, CustomInt(0))
	if err != nil {
		t.Errorf("ToYAML failed with custom type: %v", err)
	}

	if yamlValue != "custom1" {
		t.Errorf("expected 'custom1', got %v", yamlValue)
	}
}

// TestMarshalEmptyLabels tests marshaling with empty labels
func TestMarshalEmptyLabels(t *testing.T) {
	labels := []string{}

	// Test JSON
	jsonBytes, err := ToJSON(labels, 0)
	if err != nil {
		t.Errorf("ToJSON failed with empty labels: %v", err)
	}

	if string(jsonBytes) != `"Invalid"` {
		t.Errorf("expected '\"Invalid\"', got %s", string(jsonBytes))
	}

	// Test YAML
	yamlValue, err := ToYAML(labels, 0)
	if err != nil {
		t.Errorf("ToYAML failed with empty labels: %v", err)
	}

	if yamlValue != "Invalid" {
		t.Errorf("expected 'Invalid', got %v", yamlValue)
	}
}

// TestMarshalLargeEnum tests marshaling with large enums (tests performance paths)
func TestMarshalLargeEnum(t *testing.T) {
	// Create large enum to test map-based lookup path
	labels := make([]string, 50)
	for i := 0; i < 50; i++ {
		labels[i] = string(rune('A'+(i%26))) + string(rune('a'+(i/26)))
	}

	// Test JSON with large enum
	jsonBytes, err := ToJSON(labels, 25)
	if err != nil {
		t.Errorf("ToJSON failed with large enum: %v", err)
	}

	result, err := FromJSON[int](labels, jsonBytes)
	if err != nil {
		t.Errorf("FromJSON failed with large enum: %v", err)
	}

	if result != 25 {
		t.Errorf("expected 25, got %d", result)
	}

	// Test YAML with large enum
	yamlValue, err := ToYAML(labels, 30)
	if err != nil {
		t.Errorf("ToYAML failed with large enum: %v", err)
	}

	expectedLabel := labels[30]
	if yamlValue != expectedLabel {
		t.Errorf("expected %q, got %v", expectedLabel, yamlValue)
	}
}
