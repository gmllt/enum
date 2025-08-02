package internal

import (
	"encoding/binary"
	"encoding/json"
	"errors"
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
			expectError: true, // Now returns error for invalid values
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
			expectError: true, // Now returns error for invalid values
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
			expectError: true, // Now returns error for invalid values
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
			expectError: true, // Now returns error for invalid values
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

// TestToText tests text marshalling
func TestToText(t *testing.T) {
	labels := []string{"first", "second", "third"}

	tests := []struct {
		name     string
		value    int
		expected string
	}{
		{
			name:     "valid first value",
			value:    0,
			expected: "first",
		},
		{
			name:     "valid second value",
			value:    1,
			expected: "second",
		},
		{
			name:     "invalid value",
			value:    10,
			expected: InvalidLabel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ToText(labels, tt.value)
			if err != nil {
				t.Errorf("ToText failed: %v", err)
				return
			}

			if string(result) != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, string(result))
			}
		})
	}
}

// TestFromText tests text unmarshalling
func TestFromText(t *testing.T) {
	labels := []string{"first", "second", "third"}

	tests := []struct {
		name      string
		text      string
		expected  int
		wantError bool
	}{
		{
			name:      "valid first label",
			text:      "first",
			expected:  0,
			wantError: false,
		},
		{
			name:      "valid second label",
			text:      "second",
			expected:  1,
			wantError: false,
		},
		{
			name:      "invalid label",
			text:      "invalid",
			expected:  0,    // zero value
			wantError: true, // now expects error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := FromText[int](labels, []byte(tt.text))

			if tt.wantError {
				if err == nil {
					t.Errorf("expected error for invalid label, but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("FromText failed: %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}

// TestToBinary tests binary marshalling
func TestToBinary(t *testing.T) {
	labels := []string{"first", "second", "third"}

	tests := []struct {
		name  string
		value int
		label string
	}{
		{
			name:  "valid first value",
			value: 0,
			label: "first",
		},
		{
			name:  "valid second value",
			value: 1,
			label: "second",
		},
		{
			name:  "invalid value",
			value: 10,
			label: InvalidLabel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ToBinary(labels, tt.value)
			if err != nil {
				t.Errorf("ToBinary failed: %v", err)
				return
			}

			// Check format: first 2 bytes are length (big-endian), followed by string
			expectedLen := uint16(len(tt.label))
			if len(result) < 2 {
				t.Errorf("result too short: expected at least 2 bytes, got %d", len(result))
				return
			}

			actualLen := binary.BigEndian.Uint16(result[0:2])
			if actualLen != expectedLen {
				t.Errorf("expected length prefix %d, got %d", expectedLen, actualLen)
				return
			}

			actualLabel := string(result[2:])
			if actualLabel != tt.label {
				t.Errorf("expected label %q, got %q", tt.label, actualLabel)
			}
		})
	}
}

// TestFromBinary tests binary unmarshalling
func TestFromBinary(t *testing.T) {
	labels := []string{"first", "second", "third"}

	tests := []struct {
		name     string
		data     []byte
		expected int
		wantErr  bool
	}{
		{
			name:     "valid first label",
			data:     []byte{0, 5, 'f', 'i', 'r', 's', 't'}, // 2-byte length prefix (big-endian) + "first"
			expected: 0,
			wantErr:  false,
		},
		{
			name:     "valid second label",
			data:     []byte{0, 6, 's', 'e', 'c', 'o', 'n', 'd'}, // 2-byte length prefix + "second"
			expected: 1,
			wantErr:  false,
		},
		{
			name:     "invalid label",
			data:     []byte{0, 7, 'i', 'n', 'v', 'a', 'l', 'i', 'd'}, // 2-byte length prefix + "invalid"
			expected: 0,                                               // zero value
			wantErr:  true,                                            // should error for invalid enum value
		},
		{
			name:     "empty data",
			data:     []byte{},
			expected: 0,    // zero value
			wantErr:  true, // should error for too short data
		},
		{
			name:     "truncated data",
			data:     []byte{0, 10, 'a', 'b'}, // claims length 10 but only has 2 chars
			expected: 0,                       // zero value
			wantErr:  true,                    // should error for truncated data
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := FromBinary[int](labels, tt.data)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("FromBinary failed: %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}

// TestTextBinaryRoundTrip tests that marshalling and unmarshalling preserves values
func TestTextBinaryRoundTrip(t *testing.T) {
	labels := []string{"alpha", "beta", "gamma", "delta"}

	for i, label := range labels {
		t.Run(label, func(t *testing.T) {
			// Test Text round trip
			textBytes, err := ToText(labels, i)
			if err != nil {
				t.Errorf("ToText failed: %v", err)
				return
			}

			textResult, err := FromText[int](labels, textBytes)
			if err != nil {
				t.Errorf("FromText failed: %v", err)
				return
			}

			if textResult != i {
				t.Errorf("Text round trip failed: expected %d, got %d", i, textResult)
			}

			// Test Binary round trip
			binaryBytes, err := ToBinary(labels, i)
			if err != nil {
				t.Errorf("ToBinary failed: %v", err)
				return
			}

			binaryResult, err := FromBinary[int](labels, binaryBytes)
			if err != nil {
				t.Errorf("FromBinary failed: %v", err)
				return
			}

			if binaryResult != i {
				t.Errorf("Binary round trip failed: expected %d, got %d", i, binaryResult)
			}
		})
	}
}

// TestToBinaryLabelTooLong tests binary marshalling with very long labels
func TestToBinaryLabelTooLong(t *testing.T) {
	// Create a label that's longer than 65535 bytes (maximum for uint16)
	longLabel := string(make([]byte, 65536))
	labels := []string{longLabel}

	_, err := ToBinary[int](labels, 0)
	if err == nil {
		t.Error("expected error for label too long, got nil")
		return
	}

	// Verify it's the right type of error
	var labelTooLongErr *ErrLabelTooLong
	if !errors.As(err, &labelTooLongErr) {
		t.Errorf("expected ErrLabelTooLong, got %T", err)
		return
	}

	if labelTooLongErr.Length != 65536 {
		t.Errorf("expected Length=65536, got %d", labelTooLongErr.Length)
	}

	if labelTooLongErr.MaxLength != 65535 {
		t.Errorf("expected MaxLength=65535, got %d", labelTooLongErr.MaxLength)
	}
}

// TestToSQLValue tests SQL value marshalling
func TestToSQLValue(t *testing.T) {
	labels := []string{"first", "second", "third"}

	tests := []struct {
		name     string
		value    int
		expected string
		hasError bool
	}{
		{
			name:     "valid first value",
			value:    0,
			expected: "first",
			hasError: false,
		},
		{
			name:     "valid middle value",
			value:    1,
			expected: "second",
			hasError: false,
		},
		{
			name:     "valid last value",
			value:    2,
			expected: "third",
			hasError: false,
		},
		{
			name:     "invalid value",
			value:    5,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ToSQLValue[int](labels, tt.value)

			if tt.hasError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestFromSQLValue tests SQL value unmarshalling
func TestFromSQLValue(t *testing.T) {
	labels := []string{"alpha", "beta", "gamma"}

	tests := []struct {
		name          string
		src           any
		expectedValue int
		hasError      bool
	}{
		{
			name:          "string first",
			src:           "alpha",
			expectedValue: 0,
			hasError:      false,
		},
		{
			name:          "string middle",
			src:           "beta",
			expectedValue: 1,
			hasError:      false,
		},
		{
			name:          "string last",
			src:           "gamma",
			expectedValue: 2,
			hasError:      false,
		},
		{
			name:          "bytes",
			src:           []byte("alpha"),
			expectedValue: 0,
			hasError:      false,
		},
		{
			name:          "nil value",
			src:           nil,
			expectedValue: 0,
			hasError:      false,
		},
		{
			name:     "invalid string",
			src:      "invalid",
			hasError: true,
		},
		{
			name:     "invalid type",
			src:      123,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := FromSQLValue[int](labels, tt.src)

			if tt.hasError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result != tt.expectedValue {
				t.Errorf("expected value %d, got %d", tt.expectedValue, result)
			}
		})
	}
}

// TestSQLRoundTrip tests SQL marshalling round trip
func TestSQLRoundTrip(t *testing.T) {
	labels := []string{"one", "two", "three"}
	testValues := []int{0, 1, 2}

	for _, value := range testValues {
		t.Run("round trip", func(t *testing.T) {
			// Convert to SQL value
			sqlValue, err := ToSQLValue[int](labels, value)
			if err != nil {
				t.Fatalf("ToSQLValue failed: %v", err)
			}

			// Convert back from SQL value
			result, err := FromSQLValue[int](labels, sqlValue)
			if err != nil {
				t.Fatalf("FromSQLValue failed: %v", err)
			}

			// Verify round trip
			if result != value {
				t.Errorf("round trip failed: original %d, got %d", value, result)
			}
		})
	}
}
