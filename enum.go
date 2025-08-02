package enum

import (
	"fmt"

	"github.com/gmllt/enum/internal"
)

// Value is a type constraint for integer values used in the Enum type.
type Value interface {
	~int
}

// Enum is a generic enumeration type that maps integer values to string labels.
type Enum[T Value] struct {
	labels   []string
	labelMap map[string]T
	allVals  []T
}

// NewEnum creates a new Enum instance with the provided labels.
func NewEnum[T Value](labels ...string) *Enum[T] {
	cacheBuilder := internal.NewCacheBuilder[T](labels)

	return &Enum[T]{
		labels:   labels,
		labelMap: cacheBuilder.BuildLookupMap(),
		allVals:  cacheBuilder.BuildAllValues(),
	}
}

// String returns the string representation of the enumeration value.
func (e *Enum[T]) String(v T) string {
	return internal.SafeGetLabel(e.labels, v, fmt.Sprintf("Invalid(%d)", v))
}

// FromString converts a string to the corresponding enumeration value.
func (e *Enum[T]) FromString(s string) (T, error) {
	if val, ok := e.labelMap[s]; ok {
		return val, nil
	}
	var zero T
	return zero, fmt.Errorf("invalid value: %s", s)
}

// All returns all values of the enum.
func (e *Enum[T]) All() []T {
	res := make([]T, len(e.allVals))
	copy(res, e.allVals)
	return res
}

// Labels returns all labels of the enum.
// Note: The returned slice is a copy to prevent modification of internal state.
// For read-only access, consider using LabelsReadOnly() for better performance.
func (e *Enum[T]) Labels() []string {
	cp := make([]string, len(e.labels))
	copy(cp, e.labels)
	return cp
}

// LabelsReadOnly returns a read-only view of all labels.
// WARNING: Do not modify the returned slice as it shares memory with the enum.
func (e *Enum[T]) LabelsReadOnly() []string {
	return e.labels
}
