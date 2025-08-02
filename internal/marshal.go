package internal

import (
	"database/sql/driver"
	"encoding/binary"
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
	return zero, NewInvalidEnumValueError(s, labels)
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
	return zero, NewInvalidEnumValueError(s, labels)
}

// ToText serializes an enum value into text (for encoding.TextMarshaler).
func ToText[T ~int](labels []string, v T) ([]byte, error) {
	label := SafeGetLabel(labels, v, InvalidLabel)
	return []byte(label), nil
}

// FromText deserializes text into an enum value (for encoding.TextUnmarshaler).
func FromText[T ~int](labels []string, text []byte) (T, error) {
	s := string(text)
	if val, found := StringToIndex[T](labels, s); found {
		return val, nil
	}

	var zero T
	return zero, NewInvalidEnumValueError(s, labels)
}

// ToBinary serializes an enum value into binary (for encoding.BinaryMarshaler).
func ToBinary[T ~int](labels []string, v T) ([]byte, error) {
	label := SafeGetLabel(labels, v, InvalidLabel)
	// Store as length-prefixed string (2 bytes, big-endian) for efficiency
	labelBytes := []byte(label)
	if len(labelBytes) > 65535 {
		return nil, NewLabelTooLongError(len(labelBytes), 65535)
	}
	result := make([]byte, 2+len(labelBytes))
	binary.BigEndian.PutUint16(result[0:2], uint16(len(labelBytes)))
	copy(result[2:], labelBytes)
	return result, nil
}

// FromBinary deserializes binary into an enum value (for encoding.BinaryUnmarshaler).
func FromBinary[T ~int](labels []string, data []byte) (T, error) {
	var zero T

	if len(data) < 2 {
		return zero, NewBinaryDataTooShortError(2, len(data))
	}

	// Read length-prefixed string (2 bytes, big-endian)
	length := int(binary.BigEndian.Uint16(data[0:2]))
	if len(data) < 2+length {
		return zero, NewBinaryDataTruncatedError(2+length, len(data))
	}

	label := string(data[2 : 2+length])
	if val, found := StringToIndex[T](labels, label); found {
		return val, nil
	}

	return zero, NewInvalidEnumValueError(label, labels)
}

// ToSQLValue serializes an enum value for SQL storage (for driver.Valuer).
func ToSQLValue[T ~int](labels []string, v T) (driver.Value, error) {
	if !IsValidIndex(labels, v) {
		return nil, NewInvalidEnumValueError("", labels)
	}
	label := SafeGetLabel(labels, v, InvalidLabel)
	return label, nil
}

// FromSQLValue deserializes an SQL value into an enum value (for sql.Scanner).
func FromSQLValue[T ~int](labels []string, src any) (T, error) {
	var zero T

	if src == nil {
		// SQL NULL maps to zero value
		return zero, nil
	}

	var s string
	switch v := src.(type) {
	case string:
		s = v
	case []byte:
		s = string(v)
	default:
		return zero, NewInvalidEnumValueError("non-string SQL value", labels)
	}

	if val, found := StringToIndex[T](labels, s); found {
		return val, nil
	}

	return zero, NewInvalidEnumValueError(s, labels)
}
