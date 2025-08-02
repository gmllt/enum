package internal

import (
	"fmt"
	"testing"
)

// TestStringToIndex tests the adaptive string-to-index lookup
func TestStringToIndex(t *testing.T) {
	tests := []struct {
		name        string
		labels      []string
		target      string
		expectedVal int
		expectedOk  bool
	}{
		{
			name:        "small enum - found",
			labels:      []string{"a", "b", "c"},
			target:      "b",
			expectedVal: 1,
			expectedOk:  true,
		},
		{
			name:        "small enum - not found",
			labels:      []string{"a", "b", "c"},
			target:      "d",
			expectedVal: 0,
			expectedOk:  false,
		},
		{
			name:        "large enum - found",
			labels:      generateLabels(20),
			target:      "label_10",
			expectedVal: 10,
			expectedOk:  true,
		},
		{
			name:        "large enum - not found",
			labels:      generateLabels(20),
			target:      "nonexistent",
			expectedVal: 0,
			expectedOk:  false,
		},
		{
			name:        "empty enum",
			labels:      []string{},
			target:      "any",
			expectedVal: 0,
			expectedOk:  false,
		},
		{
			name:        "single element - found",
			labels:      []string{"only"},
			target:      "only",
			expectedVal: 0,
			expectedOk:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, ok := StringToIndex[int](tt.labels, tt.target)

			if ok != tt.expectedOk {
				t.Errorf("expected ok=%v, got ok=%v", tt.expectedOk, ok)
			}

			if val != tt.expectedVal {
				t.Errorf("expected val=%d, got val=%d", tt.expectedVal, val)
			}
		})
	}
}

// TestMapLookup tests the map-based lookup function
func TestMapLookup(t *testing.T) {
	labels := []string{"red", "green", "blue", "yellow"}

	tests := []struct {
		name        string
		target      string
		expectedVal int
		expectedOk  bool
	}{
		{
			name:        "first element",
			target:      "red",
			expectedVal: 0,
			expectedOk:  true,
		},
		{
			name:        "middle element",
			target:      "green",
			expectedVal: 1,
			expectedOk:  true,
		},
		{
			name:        "last element",
			target:      "yellow",
			expectedVal: 3,
			expectedOk:  true,
		},
		{
			name:        "not found",
			target:      "purple",
			expectedVal: 0,
			expectedOk:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, ok := mapLookup[int](labels, tt.target)

			if ok != tt.expectedOk {
				t.Errorf("expected ok=%v, got ok=%v", tt.expectedOk, ok)
			}

			if val != tt.expectedVal {
				t.Errorf("expected val=%d, got val=%d", tt.expectedVal, val)
			}
		})
	}
}

// TestLinearLookup tests the linear search function
func TestLinearLookup(t *testing.T) {
	labels := []string{"dog", "cat", "bird"}

	tests := []struct {
		name        string
		target      string
		expectedVal int
		expectedOk  bool
	}{
		{
			name:        "first element",
			target:      "dog",
			expectedVal: 0,
			expectedOk:  true,
		},
		{
			name:        "middle element",
			target:      "cat",
			expectedVal: 1,
			expectedOk:  true,
		},
		{
			name:        "last element",
			target:      "bird",
			expectedVal: 2,
			expectedOk:  true,
		},
		{
			name:        "not found",
			target:      "fish",
			expectedVal: 0,
			expectedOk:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, ok := linearLookup[int](labels, tt.target)

			if ok != tt.expectedOk {
				t.Errorf("expected ok=%v, got ok=%v", tt.expectedOk, ok)
			}

			if val != tt.expectedVal {
				t.Errorf("expected val=%d, got val=%d", tt.expectedVal, val)
			}
		})
	}
}

// TestBuildLabelMap tests the label map building function
func TestBuildLabelMap(t *testing.T) {
	labels := []string{"apple", "banana", "cherry"}
	labelMap := BuildLabelMap[int](labels)

	// Test that all labels are correctly mapped
	for i, label := range labels {
		if val, ok := labelMap[label]; !ok {
			t.Errorf("label %q not found in map", label)
		} else if val != i {
			t.Errorf("expected %d for label %q, got %d", i, label, val)
		}
	}

	// Test that map has correct size
	if len(labelMap) != len(labels) {
		t.Errorf("expected map size %d, got %d", len(labels), len(labelMap))
	}

	// Test that non-existent label returns false
	if _, ok := labelMap["nonexistent"]; ok {
		t.Error("expected false for nonexistent label")
	}
}

// TestLookupThreshold tests that the correct lookup method is used
func TestLookupThreshold(t *testing.T) {
	// Test with small enum (should use linear lookup)
	smallLabels := generateLabels(5)
	val, ok := StringToIndex[int](smallLabels, "label_2")
	if !ok {
		t.Error("expected to find label_2 in small enum")
	}
	if val != 2 {
		t.Errorf("expected 2, got %d", val)
	}

	// Test with large enum (should use map lookup)
	largeLabels := generateLabels(20)
	val, ok = StringToIndex[int](largeLabels, "label_15")
	if !ok {
		t.Error("expected to find label_15 in large enum")
	}
	if val != 15 {
		t.Errorf("expected 15, got %d", val)
	}
}

// TestEdgeCases tests edge cases for lookup functions
func TestEdgeCases(t *testing.T) {
	// Empty labels
	val, ok := StringToIndex[int]([]string{}, "anything")
	if ok {
		t.Error("expected false for empty labels")
	}
	if val != 0 {
		t.Errorf("expected 0 for empty labels, got %d", val)
	}

	// Empty target string
	labels := []string{"", "not empty"}
	val, ok = StringToIndex[int](labels, "")
	if !ok {
		t.Error("expected true for empty string target")
	}
	if val != 0 {
		t.Errorf("expected 0 for empty string target, got %d", val)
	}

	// Duplicates in labels (should return first occurrence)
	duplicateLabels := []string{"duplicate", "unique", "duplicate"}
	val, ok = StringToIndex[int](duplicateLabels, "duplicate")
	if !ok {
		t.Error("expected true for duplicate label")
	}
	if val != 0 {
		t.Errorf("expected 0 (first occurrence) for duplicate label, got %d", val)
	}
}

// TestLookupConsistency ensures both lookup methods give the same results
func TestLookupConsistency(t *testing.T) {
	labels := []string{"alpha", "beta", "gamma", "delta", "epsilon"}

	for i, label := range labels {
		// Test both lookup methods
		mapVal, mapOk := mapLookup[int](labels, label)
		linearVal, linearOk := linearLookup[int](labels, label)

		if mapOk != linearOk {
			t.Errorf("inconsistent ok values for %q: map=%v, linear=%v", label, mapOk, linearOk)
		}

		if mapVal != linearVal {
			t.Errorf("inconsistent values for %q: map=%d, linear=%d", label, mapVal, linearVal)
		}

		if mapVal != i {
			t.Errorf("expected %d for %q, got %d", i, label, mapVal)
		}
	}

	// Test non-existent label
	mapVal, mapOk := mapLookup[int](labels, "nonexistent")
	linearVal, linearOk := linearLookup[int](labels, "nonexistent")

	if mapOk != linearOk || mapOk != false {
		t.Errorf("inconsistent results for nonexistent label: map=%v, linear=%v", mapOk, linearOk)
	}

	if mapVal != linearVal || mapVal != 0 {
		t.Errorf("inconsistent values for nonexistent label: map=%d, linear=%d", mapVal, linearVal)
	}
}

// generateLabels is a helper function to create test labels
func generateLabels(count int) []string {
	labels := make([]string, count)
	for i := 0; i < count; i++ {
		labels[i] = fmt.Sprintf("label_%d", i)
	}
	return labels
}
