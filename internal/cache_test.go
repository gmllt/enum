package internal

import (
	"reflect"
	"testing"
)

// TestNewCacheBuilder tests cache builder creation
func TestNewCacheBuilder(t *testing.T) {
	labels := []string{"first", "second", "third"}
	builder := NewCacheBuilder[int](labels)

	if builder == nil {
		t.Fatal("NewCacheBuilder returned nil")
	}

	if !reflect.DeepEqual(builder.labels, labels) {
		t.Errorf("expected labels %v, got %v", labels, builder.labels)
	}
}

// TestBuildAllValues tests building all values cache
func TestBuildAllValues(t *testing.T) {
	tests := []struct {
		name     string
		labels   []string
		expected []int
	}{
		{
			name:     "empty labels",
			labels:   []string{},
			expected: []int{},
		},
		{
			name:     "single label",
			labels:   []string{"only"},
			expected: []int{0},
		},
		{
			name:     "multiple labels",
			labels:   []string{"a", "b", "c", "d"},
			expected: []int{0, 1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewCacheBuilder[int](tt.labels)
			result := builder.BuildAllValues()

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestBuildLookupMap tests building lookup map cache
func TestBuildLookupMap(t *testing.T) {
	labels := []string{"apple", "banana", "cherry"}
	builder := NewCacheBuilder[int](labels)

	lookupMap := builder.BuildLookupMap()

	// Test that all labels are correctly mapped
	for i, label := range labels {
		if val, ok := lookupMap[label]; !ok {
			t.Errorf("label %q not found in lookup map", label)
		} else if val != i {
			t.Errorf("expected %d for label %q, got %d", i, label, val)
		}
	}

	// Test that map has correct size
	if len(lookupMap) != len(labels) {
		t.Errorf("expected map size %d, got %d", len(labels), len(lookupMap))
	}

	// Test that non-existent label is not in map
	if _, ok := lookupMap["nonexistent"]; ok {
		t.Error("unexpected entry for nonexistent label")
	}
}

// TestShouldUseCachedLookup tests lookup method decision
func TestShouldUseCachedLookup(t *testing.T) {
	tests := []struct {
		name     string
		labels   []string
		expected bool
	}{
		{
			name:     "empty labels",
			labels:   []string{},
			expected: false,
		},
		{
			name:     "small enum",
			labels:   []string{"a", "b", "c"},
			expected: false,
		},
		{
			name:     "at threshold",
			labels:   generateTestLabels(DefaultLookupThreshold),
			expected: false,
		},
		{
			name:     "above threshold",
			labels:   generateTestLabels(DefaultLookupThreshold + 1),
			expected: true,
		},
		{
			name:     "large enum",
			labels:   generateTestLabels(50),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewCacheBuilder[int](tt.labels)
			result := builder.ShouldUseCachedLookup()

			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestCacheBuilderWithCustomTypes tests cache builder with custom integer types
func TestCacheBuilderWithCustomTypes(t *testing.T) {
	type CustomInt int

	labels := []string{"custom1", "custom2", "custom3"}
	builder := NewCacheBuilder[CustomInt](labels)

	// Test BuildAllValues with custom type
	allVals := builder.BuildAllValues()
	expected := []CustomInt{0, 1, 2}

	if !reflect.DeepEqual(allVals, expected) {
		t.Errorf("expected %v, got %v", expected, allVals)
	}

	// Test BuildLookupMap with custom type
	lookupMap := builder.BuildLookupMap()

	for i, label := range labels {
		if val, ok := lookupMap[label]; !ok {
			t.Errorf("label %q not found in lookup map", label)
		} else if val != CustomInt(i) {
			t.Errorf("expected %d for label %q, got %d", i, label, val)
		}
	}
}

// TestCacheBuilderConsistency tests that cache builder produces consistent results
func TestCacheBuilderConsistency(t *testing.T) {
	labels := []string{"red", "green", "blue", "yellow", "orange"}

	// Create multiple builders with same labels
	builder1 := NewCacheBuilder[int](labels)
	builder2 := NewCacheBuilder[int](labels)

	// Test AllValues consistency
	allVals1 := builder1.BuildAllValues()
	allVals2 := builder2.BuildAllValues()

	if !reflect.DeepEqual(allVals1, allVals2) {
		t.Errorf("inconsistent AllValues: %v vs %v", allVals1, allVals2)
	}

	// Test LookupMap consistency
	map1 := builder1.BuildLookupMap()
	map2 := builder2.BuildLookupMap()

	if !reflect.DeepEqual(map1, map2) {
		t.Errorf("inconsistent LookupMaps: %v vs %v", map1, map2)
	}

	// Test ShouldUseCachedLookup consistency
	should1 := builder1.ShouldUseCachedLookup()
	should2 := builder2.ShouldUseCachedLookup()

	if should1 != should2 {
		t.Errorf("inconsistent ShouldUseCachedLookup: %v vs %v", should1, should2)
	}
}

// TestCacheBuilderMemoryEfficiency tests that cache builder doesn't modify original labels
func TestCacheBuilderMemoryEfficiency(t *testing.T) {
	originalLabels := []string{"original1", "original2"}

	// Make a copy to ensure we can detect modifications
	labelsCopy := make([]string, len(originalLabels))
	copy(labelsCopy, originalLabels)

	builder := NewCacheBuilder[int](labelsCopy)

	// Build caches
	_ = builder.BuildAllValues()
	_ = builder.BuildLookupMap()
	_ = builder.ShouldUseCachedLookup()

	// Verify original labels weren't modified
	if !reflect.DeepEqual(labelsCopy, originalLabels) {
		t.Errorf("cache builder modified original labels: %v vs %v", labelsCopy, originalLabels)
	}
}

// TestEmptyLabelsCacheBuilder tests cache builder with empty labels
func TestEmptyLabelsCacheBuilder(t *testing.T) {
	builder := NewCacheBuilder[int]([]string{})

	// Test BuildAllValues
	allVals := builder.BuildAllValues()
	if len(allVals) != 0 {
		t.Errorf("expected empty slice, got %v", allVals)
	}

	// Test BuildLookupMap
	lookupMap := builder.BuildLookupMap()
	if len(lookupMap) != 0 {
		t.Errorf("expected empty map, got %v", lookupMap)
	}

	// Test ShouldUseCachedLookup
	if builder.ShouldUseCachedLookup() {
		t.Error("expected false for empty labels")
	}
}

// generateTestLabels is a helper function to create test labels
func generateTestLabels(count int) []string {
	labels := make([]string, count)
	for i := 0; i < count; i++ {
		labels[i] = string(rune('a' + (i % 26)))
		if i >= 26 {
			labels[i] += string(rune('0' + (i / 26)))
		}
	}
	return labels
}
