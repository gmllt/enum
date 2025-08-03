package enum

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"encoding/json"

	"github.com/gmllt/enum/internal"
)

// Wrapper wraps an Enum and provides JSON/YAML serialization.
type Wrapper[T Value] struct {
	Enum    *Enum[T]
	Current T
	labels  []string
}

// Ensure Wrapper implements the necessary interfaces.
var (
	_ json.Marshaler             = (*Wrapper[int])(nil)
	_ json.Unmarshaler           = (*Wrapper[int])(nil)
	_ encoding.TextMarshaler     = (*Wrapper[int])(nil)
	_ encoding.TextUnmarshaler   = (*Wrapper[int])(nil)
	_ encoding.BinaryMarshaler   = (*Wrapper[int])(nil)
	_ encoding.BinaryUnmarshaler = (*Wrapper[int])(nil)
	_ driver.Valuer              = (*Wrapper[int])(nil)
	_ sql.Scanner                = (*Wrapper[int])(nil)
)

// NewWrapper creates a new Wrapper with the given labels.
func NewWrapper[T Value](labels ...string) Wrapper[T] {
	e := NewEnum[T](labels...)
	return Wrapper[T]{
		Enum:   e,
		labels: labels,
	}
}

// String returns the string representation of the wrapped value.
func (w Wrapper[T]) String() string {
	return w.Enum.String(w.Current)
}

// All returns all values of the wrapped enum.
func (w Wrapper[T]) All() []T {
	return w.Enum.All()
}

// Labels returns all labels of the wrapped enum.
func (w Wrapper[T]) Labels() []string {
	return w.Enum.Labels()
}

// ensureEnum initializes the Enum if it is nil and labels are provided.
func (w *Wrapper[T]) ensureEnum() {
	if w.Enum == nil && w.labels != nil {
		w.Enum = NewEnum[T](w.labels...)
	}
}

// MarshalJSON implements json.Marshaler.
func (w Wrapper[T]) MarshalJSON() ([]byte, error) {
	return internal.ToJSON[T](w.Enum.labels, w.Current)
}

// UnmarshalJSON implements json.Unmarshaler.
func (w *Wrapper[T]) UnmarshalJSON(data []byte) error {
	w.ensureEnum()
	val, err := internal.FromJSON[T](w.Enum.labels, data)
	if err != nil {
		return err
	}
	w.Current = val
	return nil
}

// MarshalYAML implements yaml.Marshaler.
func (w Wrapper[T]) MarshalYAML() (any, error) {
	return internal.ToYAML[T](w.Enum.labels, w.Current)
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (w *Wrapper[T]) UnmarshalYAML(unmarshal func(any) error) error {
	w.ensureEnum()
	val, err := internal.FromYAML[T](w.Enum.labels, unmarshal)
	if err != nil {
		return err
	}
	w.Current = val
	return nil
}

// Get returns the current value.
func (w Wrapper[T]) Get() T {
	return w.Current
}

// Set sets the current value.
func (w *Wrapper[T]) Set(v T) {
	w.Current = v
}

// MarshalText implements encoding.TextMarshaler.
func (w Wrapper[T]) MarshalText() ([]byte, error) {
	return internal.ToText[T](w.Enum.labels, w.Current)
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (w *Wrapper[T]) UnmarshalText(text []byte) error {
	w.ensureEnum()
	val, err := internal.FromText[T](w.Enum.labels, text)
	if err != nil {
		return err
	}
	w.Current = val
	return nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (w Wrapper[T]) MarshalBinary() ([]byte, error) {
	return internal.ToBinary[T](w.Enum.labels, w.Current)
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (w *Wrapper[T]) UnmarshalBinary(data []byte) error {
	w.ensureEnum()
	val, err := internal.FromBinary[T](w.Enum.labels, data)
	if err != nil {
		return err
	}
	w.Current = val
	return nil
}

// Value implements driver.Valuer for SQL integration.
func (w Wrapper[T]) Value() (driver.Value, error) {
	return internal.ToSQLValue[T](w.Enum.labels, w.Current)
}

// Scan implements sql.Scanner for SQL integration.
func (w *Wrapper[T]) Scan(src any) error {
	w.ensureEnum()
	val, err := internal.FromSQLValue[T](w.Enum.labels, src)
	if err != nil {
		return err
	}
	w.Current = val
	return nil
}
