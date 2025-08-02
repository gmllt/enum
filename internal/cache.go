package internal

// CacheBuilder helps build cached data structures for enum optimization
type CacheBuilder[T ~int] struct {
	labels []string
}

// NewCacheBuilder creates a new cache builder
func NewCacheBuilder[T ~int](labels []string) *CacheBuilder[T] {
	return &CacheBuilder[T]{labels: labels}
}

// BuildAllValues creates a pre-computed slice of all enum values
func (cb *CacheBuilder[T]) BuildAllValues() []T {
	allVals := make([]T, len(cb.labels))
	for i := range cb.labels {
		allVals[i] = T(i)
	}
	return allVals
}

// BuildLookupMap creates a lookup map for string-to-value conversion
func (cb *CacheBuilder[T]) BuildLookupMap() map[string]T {
	return BuildLabelMap[T](cb.labels)
}

// ShouldUseCachedLookup determines if a cached lookup map should be used
// based on the number of labels and expected usage patterns
func (cb *CacheBuilder[T]) ShouldUseCachedLookup() bool {
	return len(cb.labels) > LookupThreshold
}
