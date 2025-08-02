package enum

import "github.com/gmllt/enum/internal"

// Wrapper wraps an Enum and provides JSON/YAML serialization.
type Wrapper[T Value] struct {
	Enum  *Enum[T]
	Value T
}

// NewWrapper creates a new Wrapper with the given labels.
func NewWrapper[T Value](labels ...string) Wrapper[T] {
	return Wrapper[T]{Enum: NewEnum[T](labels...)}
}

// String returns the string representation of the wrapped value.
func (w Wrapper[T]) String() string {
	return w.Enum.String(w.Value)
}

// All returns all values of the wrapped enum.
func (w Wrapper[T]) All() []T {
	return w.Enum.All()
}

// Labels returns all labels of the wrapped enum.
func (w Wrapper[T]) Labels() []string {
	return w.Enum.Labels()
}

// MarshalJSON implements json.Marshaler.
func (w Wrapper[T]) MarshalJSON() ([]byte, error) {
	return internal.ToJSON[T](w.Enum.labels, w.Value)
}

// UnmarshalJSON implements json.Unmarshaler.
func (w *Wrapper[T]) UnmarshalJSON(b []byte) error {
	val, err := internal.FromJSON[T](w.Enum.labels, b)
	if err != nil {
		return err
	}
	w.Value = val
	return nil
}

// MarshalYAML implements yaml.Marshaler.
func (w Wrapper[T]) MarshalYAML() (any, error) {
	return internal.ToYAML[T](w.Enum.labels, w.Value)
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (w *Wrapper[T]) UnmarshalYAML(unmarshal func(any) error) error {
	val, err := internal.FromYAML[T](w.Enum.labels, unmarshal)
	if err != nil {
		return err
	}
	w.Value = val
	return nil
}

// Get returns the wrapped value.
func (w Wrapper[T]) Get() T {
	return w.Value
}

// Set sets the wrapped value.
func (w *Wrapper[T]) Set(v T) {
	w.Value = v
}
