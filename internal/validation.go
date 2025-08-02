package internal

import "fmt"

// ValidateIndex checks if an index is within bounds for the given labels
func ValidateIndex[T ~int](labels []string, index T) error {
	i := int(index)
	if i < 0 || i >= len(labels) {
		return fmt.Errorf("index %d out of bounds for enum with %d labels", i, len(labels))
	}
	return nil
}

// IsValidIndex checks if an index is within bounds (returns bool instead of error)
func IsValidIndex[T ~int](labels []string, index T) bool {
	i := int(index)
	return i >= 0 && i < len(labels)
}

// SafeGetLabel returns the label for an index, or a default value if invalid
func SafeGetLabel[T ~int](labels []string, index T, defaultLabel string) string {
	if IsValidIndex(labels, index) {
		return labels[int(index)]
	}
	return defaultLabel
}

// SafeGetLabelWithError returns the label for an index, or an error if invalid
func SafeGetLabelWithError[T ~int](labels []string, index T) (string, error) {
	if err := ValidateIndex(labels, index); err != nil {
		return "", err
	}
	return labels[int(index)], nil
}
