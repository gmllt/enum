package enum

import (
	"fmt"
	"testing"
)

// BenchmarkEnumString benchmarks string conversion
func BenchmarkEnumString(b *testing.B) {
	enum := NewEnum[int]("first", "second", "third", "fourth", "fifth")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = enum.String(i % 5)
	}
}

// BenchmarkEnumFromString benchmarks string to value conversion
func BenchmarkEnumFromString(b *testing.B) {
	labels := []string{"first", "second", "third", "fourth", "fifth"}
	enum := NewEnum[int](labels...)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = enum.FromString(labels[i%len(labels)])
	}
}

// BenchmarkEnumFromStringSmall benchmarks small enum lookup
func BenchmarkEnumFromStringSmall(b *testing.B) {
	labels := []string{"a", "b", "c"}
	enum := NewEnum[int](labels...)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = enum.FromString(labels[i%len(labels)])
	}
}

// BenchmarkEnumFromStringLarge benchmarks large enum lookup
func BenchmarkEnumFromStringLarge(b *testing.B) {
	// Create large enum to test map-based lookup
	labels := make([]string, 100)
	for i := 0; i < 100; i++ {
		labels[i] = fmt.Sprintf("label_%d", i)
	}
	enum := NewEnum[int](labels...)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = enum.FromString(labels[i%len(labels)])
	}
}

// BenchmarkEnumAll benchmarks getting all values
func BenchmarkEnumAll(b *testing.B) {
	enum := NewEnum[int]("red", "green", "blue", "yellow", "orange", "purple", "pink", "brown")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = enum.All()
	}
}

// BenchmarkEnumLabels benchmarks getting all labels
func BenchmarkEnumLabels(b *testing.B) {
	enum := NewEnum[int]("monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = enum.Labels()
	}
}

// BenchmarkEnumLabelsReadOnly benchmarks read-only label access
func BenchmarkEnumLabelsReadOnly(b *testing.B) {
	enum := NewEnum[int]("monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = enum.LabelsReadOnly()
	}
}

// BenchmarkWrapperMarshalJSON benchmarks JSON marshaling
func BenchmarkWrapperMarshalJSON(b *testing.B) {
	wrapper := NewWrapper[int]("alpha", "beta", "gamma", "delta", "epsilon")
	wrapper.Set(2)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = wrapper.MarshalJSON()
	}
}

// BenchmarkWrapperUnmarshalJSON benchmarks JSON unmarshaling
func BenchmarkWrapperUnmarshalJSON(b *testing.B) {
	wrapper := NewWrapper[int]("alpha", "beta", "gamma", "delta", "epsilon")
	jsonData := []byte(`"gamma"`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = wrapper.UnmarshalJSON(jsonData)
	}
}

// BenchmarkNewEnum benchmarks enum creation
func BenchmarkNewEnum(b *testing.B) {
	labels := []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewEnum[int](labels...)
	}
}

// BenchmarkNewEnumLarge benchmarks large enum creation
func BenchmarkNewEnumLarge(b *testing.B) {
	labels := make([]string, 50)
	for i := 0; i < 50; i++ {
		labels[i] = fmt.Sprintf("item_%d", i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewEnum[int](labels...)
	}
}

// BenchmarkCompareSmallVsLarge compares performance characteristics
func BenchmarkCompareSmallVsLarge(b *testing.B) {
	// Small enum
	b.Run("Small/FromString", func(b *testing.B) {
		enum := NewEnum[int]("a", "b", "c")
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = enum.FromString("b")
		}
	})

	// Large enum
	b.Run("Large/FromString", func(b *testing.B) {
		labels := make([]string, 50)
		for i := 0; i < 50; i++ {
			labels[i] = fmt.Sprintf("label_%d", i)
		}
		enum := NewEnum[int](labels...)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = enum.FromString("label_25")
		}
	})
}
