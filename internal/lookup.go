package internal

// LookupThreshold defines when to use map-based lookup vs linear search
const LookupThreshold = DefaultLookupThreshold

// StringToIndex performs optimized string-to-index lookup.
// Uses map-based lookup for large slices, linear search for small ones.
func StringToIndex[T ~int](labels []string, target string) (T, bool) {
	if len(labels) > LookupThreshold {
		return mapLookup[T](labels, target)
	}
	return linearLookup[T](labels, target)
}

// mapLookup uses a map for O(1) lookup - efficient for large enums
func mapLookup[T ~int](labels []string, target string) (T, bool) {
	labelMap := make(map[string]T, len(labels))
	for i, label := range labels {
		labelMap[label] = T(i)
	}

	if val, ok := labelMap[target]; ok {
		return val, true
	}

	var zero T
	return zero, false
}

// linearLookup uses linear search - efficient for small enums
func linearLookup[T ~int](labels []string, target string) (T, bool) {
	for i, label := range labels {
		if label == target {
			return T(i), true
		}
	}

	var zero T
	return zero, false
}

// BuildLabelMap creates a map for string-to-index lookup.
// Used when the map will be reused multiple times.
func BuildLabelMap[T ~int](labels []string) map[string]T {
	labelMap := make(map[string]T, len(labels))
	for i, label := range labels {
		labelMap[label] = T(i)
	}
	return labelMap
}
