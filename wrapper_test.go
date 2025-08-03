package enum

import (
	"database/sql/driver"
	"encoding/binary"
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
	if wrapper.Current != 0 {
		t.Errorf("expected default value 0, got %d", wrapper.Current)
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
			wrapper.Current = tt.value
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

	// Test that internal Current field is also updated
	if wrapper.Current != 2 {
		t.Errorf("expected internal Current 2, got %d", wrapper.Current)
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
			expectError: true, // Now returns error for invalid values
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

// TestTextMarshalling tests encoding.TextMarshaler/TextUnmarshaler interfaces
func TestTextMarshalling(t *testing.T) {
	wrapper := NewWrapper[int]("first", "second", "third")
	wrapper.Set(1) // "second"

	// Test MarshalText
	textBytes, err := wrapper.MarshalText()
	if err != nil {
		t.Errorf("MarshalText failed: %v", err)
		return
	}

	expected := "second"
	if string(textBytes) != expected {
		t.Errorf("expected %q, got %q", expected, string(textBytes))
	}

	// Test UnmarshalText
	newWrapper := NewWrapper[int]("first", "second", "third")
	err = newWrapper.UnmarshalText([]byte("third"))
	if err != nil {
		t.Errorf("UnmarshalText failed: %v", err)
		return
	}

	if newWrapper.Get() != 2 {
		t.Errorf("expected value 2, got %d", newWrapper.Get())
	}

	// Test round trip
	textBytes2, err := newWrapper.MarshalText()
	if err != nil {
		t.Errorf("MarshalText round trip failed: %v", err)
		return
	}

	if string(textBytes2) != "third" {
		t.Errorf("expected %q, got %q", "third", string(textBytes2))
	}
}

// TestBinaryMarshalling tests encoding.BinaryMarshaler/BinaryUnmarshaler interfaces
func TestBinaryMarshalling(t *testing.T) {
	wrapper := NewWrapper[int]("alpha", "beta", "gamma")
	wrapper.Set(2) // "gamma"

	// Test MarshalBinary
	binaryData, err := wrapper.MarshalBinary()
	if err != nil {
		t.Errorf("MarshalBinary failed: %v", err)
		return
	}

	// Verify binary format: 2-byte length prefix (big-endian) + string
	expectedLen := uint16(len("gamma"))
	if len(binaryData) < 2 {
		t.Errorf("binary data too short: expected at least 2 bytes, got %d", len(binaryData))
		return
	}

	actualLen := binary.BigEndian.Uint16(binaryData[0:2])
	if actualLen != expectedLen {
		t.Errorf("expected length prefix %d, got %d", expectedLen, actualLen)
		return
	}

	actualLabel := string(binaryData[2:])
	if actualLabel != "gamma" {
		t.Errorf("expected label %q, got %q", "gamma", actualLabel)
	}

	// Test UnmarshalBinary
	newWrapper := NewWrapper[int]("alpha", "beta", "gamma")

	// Create binary data for "beta" (index 1) - 2-byte length prefix (big-endian)
	testData := []byte{0, 4, 'b', 'e', 't', 'a'} // length 4 (big-endian) + "beta"
	err = newWrapper.UnmarshalBinary(testData)
	if err != nil {
		t.Errorf("UnmarshalBinary failed: %v", err)
		return
	}

	if newWrapper.Get() != 1 {
		t.Errorf("expected value 1, got %d", newWrapper.Get())
	}

	// Test round trip
	binaryData2, err := newWrapper.MarshalBinary()
	if err != nil {
		t.Errorf("MarshalBinary round trip failed: %v", err)
		return
	}

	// Should match our test data
	if len(binaryData2) != len(testData) {
		t.Errorf("expected binary length %d, got %d", len(testData), len(binaryData2))
		return
	}

	for i, b := range testData {
		if binaryData2[i] != b {
			t.Errorf("binary data mismatch at index %d: expected %d, got %d", i, b, binaryData2[i])
		}
	}
}

// TestTextUnmarshalInvalidValue tests text unmarshalling with invalid values
func TestTextUnmarshalInvalidValue(t *testing.T) {
	wrapper := NewWrapper[int]("valid1", "valid2")

	err := wrapper.UnmarshalText([]byte("invalid"))
	if err == nil {
		t.Errorf("UnmarshalText should fail with invalid value")
		return
	}

	// Error is expected, test passes
}

// TestBinaryUnmarshalInvalidData tests binary unmarshalling with invalid data
func TestBinaryUnmarshalInvalidData(t *testing.T) {
	wrapper := NewWrapper[int]("valid1", "valid2")

	tests := []struct {
		name string
		data []byte
	}{
		{
			name: "empty data",
			data: []byte{},
		},
		{
			name: "truncated data",
			data: []byte{0, 10, 'a', 'b'}, // claims length 10 (2-byte big-endian) but only has 2 chars
		},
		{
			name: "invalid label",
			data: []byte{0, 7, 'i', 'n', 'v', 'a', 'l', 'i', 'd'}, // 2-byte length prefix + "invalid"
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := wrapper.UnmarshalBinary(tt.data)
			if err == nil {
				t.Errorf("UnmarshalBinary should fail with invalid data for test case: %s", tt.name)
				return
			}
			// Error is expected for invalid data, test passes
		})
	}
}

// TestAllMarshallingInterfaces tests that wrapper implements all expected interfaces
func TestAllMarshallingInterfaces(t *testing.T) {
	wrapper := NewWrapper[int]("test")

	// Test that wrapper implements all marshalling interfaces
	var _ json.Marshaler = wrapper
	var _ json.Unmarshaler = &wrapper
	var _ driver.Valuer = wrapper
	var _ driver.Valuer = &wrapper // Valuer works on both value and pointer
	// Note: sql.Scanner requires pointer receiver, so test with pointer
	// Note: yaml interfaces would need yaml package import to test, but they're implemented

	// These should compile if wrapper implements the interfaces
	_, err1 := wrapper.MarshalText()
	_, err2 := wrapper.MarshalBinary()
	_, err3 := wrapper.Value()
	err4 := wrapper.UnmarshalText([]byte("test"))
	err5 := wrapper.UnmarshalBinary([]byte{0, 4, 't', 'e', 's', 't'}) // 2-byte length prefix + "test"
	err6 := wrapper.Scan("test")

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil || err6 != nil {
		t.Errorf("marshalling interfaces should work without errors")
	}
}

// TestWrapperSQLValue tests SQL Value implementation
func TestWrapperSQLValue(t *testing.T) {
	wrapper := NewWrapper[int]("red", "green", "blue")

	tests := []struct {
		name     string
		value    int
		expected driver.Value
		hasError bool
	}{
		{
			name:     "valid first value",
			value:    0,
			expected: "red",
			hasError: false,
		},
		{
			name:     "valid middle value",
			value:    1,
			expected: "green",
			hasError: false,
		},
		{
			name:     "valid last value",
			value:    2,
			expected: "blue",
			hasError: false,
		},
		{
			name:     "invalid value",
			value:    5,
			expected: nil,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrapper.Current = tt.value
			result, err := wrapper.Value()

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

// TestWrapperSQLScan tests SQL Scan implementation
func TestWrapperSQLScan(t *testing.T) {
	wrapper := NewWrapper[int]("alpha", "beta", "gamma")

	tests := []struct {
		name          string
		src           any
		expectedValue int
		hasError      bool
	}{
		{
			name:          "scan string first",
			src:           "alpha",
			expectedValue: 0,
			hasError:      false,
		},
		{
			name:          "scan string middle",
			src:           "beta",
			expectedValue: 1,
			hasError:      false,
		},
		{
			name:          "scan string last",
			src:           "gamma",
			expectedValue: 2,
			hasError:      false,
		},
		{
			name:          "scan bytes",
			src:           []byte("alpha"),
			expectedValue: 0,
			hasError:      false,
		},
		{
			name:          "scan nil",
			src:           nil,
			expectedValue: 0,
			hasError:      false,
		},
		{
			name:     "scan invalid string",
			src:      "invalid",
			hasError: true,
		},
		{
			name:     "scan invalid type",
			src:      123,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrapper.Current = 999 // Set to invalid value to ensure scan works
			err := wrapper.Scan(tt.src)

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

			if wrapper.Current != tt.expectedValue {
				t.Errorf("expected value %d, got %d", tt.expectedValue, wrapper.Current)
			}
		})
	}
}

// TestWrapperSQLRoundTrip tests SQL Value/Scan round trip
func TestWrapperSQLRoundTrip(t *testing.T) {
	wrapper1 := NewWrapper[int]("one", "two", "three")
	wrapper2 := NewWrapper[int]("one", "two", "three")

	testValues := []int{0, 1, 2}

	for _, value := range testValues {
		t.Run("round trip", func(t *testing.T) {
			// Set original value
			wrapper1.Current = value

			// Convert to SQL value
			sqlValue, err := wrapper1.Value()
			if err != nil {
				t.Fatalf("Value() failed: %v", err)
			}

			// Scan SQL value back
			err = wrapper2.Scan(sqlValue)
			if err != nil {
				t.Fatalf("Scan() failed: %v", err)
			}

			// Verify round trip
			if wrapper1.Current != wrapper2.Current {
				t.Errorf("round trip failed: original %d, got %d", wrapper1.Current, wrapper2.Current)
			}
		})
	}
}

// TestWrapperEnsureEnum tests the lazy initialization of Enum through ensureEnum()
func TestWrapperEnsureEnum(t *testing.T) {
	// Create a wrapper with labels but nil Enum (simulating deserialization scenario)
	labels := []string{"red", "green", "blue"}
	wrapper := Wrapper[int]{
		Enum:    nil, // Simulate a deserialized state where Enum might be nil
		Current: 1,
		labels:  labels,
	}

	// Test that ensureEnum() works through UnmarshalJSON
	jsonData := []byte(`"green"`)
	err := wrapper.UnmarshalJSON(jsonData)
	if err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}

	// Verify that the Enum was properly initialized
	if wrapper.Enum == nil {
		t.Fatal("ensureEnum() did not initialize the Enum")
	}

	// Verify that the labels were properly set
	if !reflect.DeepEqual(wrapper.Enum.labels, labels) {
		t.Errorf("expected labels %v, got %v", labels, wrapper.Enum.labels)
	}

	// Verify that the value was correctly unmarshaled
	if wrapper.Current != 1 {
		t.Errorf("expected current value 1, got %d", wrapper.Current)
	}
}

// TestWrapperEnsureEnumWithAllUnmarshalMethods tests ensureEnum through all unmarshal methods
func TestWrapperEnsureEnumWithAllUnmarshalMethods(t *testing.T) {
	labels := []string{"small", "medium", "large"}

	testCases := []struct {
		name   string
		testFn func(*Wrapper[int]) error
	}{
		{
			name: "UnmarshalJSON",
			testFn: func(w *Wrapper[int]) error {
				return w.UnmarshalJSON([]byte(`"medium"`))
			},
		},
		{
			name: "UnmarshalText",
			testFn: func(w *Wrapper[int]) error {
				return w.UnmarshalText([]byte("large"))
			},
		},
		{
			name: "UnmarshalBinary",
			testFn: func(w *Wrapper[int]) error {
				// Create binary data for "small" (2-byte length + "small")
				data := make([]byte, 2+5)
				binary.BigEndian.PutUint16(data[:2], 5) // length of "small"
				copy(data[2:], "small")
				return w.UnmarshalBinary(data)
			},
		},
		{
			name: "Scan",
			testFn: func(w *Wrapper[int]) error {
				return w.Scan("medium")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create wrapper with nil Enum but valid labels
			wrapper := Wrapper[int]{
				Enum:    nil,
				Current: 0,
				labels:  labels,
			}

			// Test the unmarshal method
			err := tc.testFn(&wrapper)
			if err != nil {
				t.Fatalf("%s failed: %v", tc.name, err)
			}

			// Verify that ensureEnum() worked
			if wrapper.Enum == nil {
				t.Fatalf("ensureEnum() did not initialize the Enum in %s", tc.name)
			}

			// Verify labels are correct
			if !reflect.DeepEqual(wrapper.Enum.labels, labels) {
				t.Errorf("expected labels %v, got %v", labels, wrapper.Enum.labels)
			}
		})
	}
}
