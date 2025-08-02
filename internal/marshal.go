package internal

import (
	"encoding/json"
)

// ToJSON serializes an enum value into JSON.
func ToJSON[T ~int](labels []string, v T) ([]byte, error) {
	label := SafeGetLabel(labels, v, InvalidLabel)
	return json.Marshal(label)
}

// FromJSON deserializes JSON into an enum value.
func FromJSON[T ~int](labels []string, b []byte) (T, error) {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		var zero T
		return zero, err
	}

	if val, found := StringToIndex[T](labels, s); found {
		return val, nil
	}

	var zero T
	return zero, nil
} // ToYAML serializes an enum value into YAML.
func ToYAML[T ~int](labels []string, v T) (any, error) {
	return SafeGetLabel(labels, v, InvalidLabel), nil
}

// FromYAML deserializes YAML into an enum value.
func FromYAML[T ~int](labels []string, unmarshal func(any) error) (T, error) {
	var s string
	if err := unmarshal(&s); err != nil {
		var zero T
		return zero, err
	}

	if val, found := StringToIndex[T](labels, s); found {
		return val, nil
	}

	var zero T
	return zero, nil
}
