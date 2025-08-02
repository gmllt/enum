# Enum Library for Go

A lightweight, generic enumeration library for Go that provides type-safe enums with marshalling support for JSON and YAML.

## Features

- Generic type-safe enumerations
- Automatic JSON and YAML marshalling/unmarshalling
- Optimized lookup performance (linear search for small enums, map-based for large ones)
- String conversion support
- Zero dependencies for core functionality

---

## Basic Usage

### Creating an Enum

```go
package main

import (
    "fmt"
    "github.com/gmllt/enum"
)

func main() {
    // Create an enum with string labels
    colors := enum.NewEnum[int]("red", "green", "blue")
    
    // Use enum values (0, 1, 2 correspond to the labels)
    fmt.Println(colors.String(0)) // Output: "red"
    fmt.Println(colors.String(1)) // Output: "green"
    fmt.Println(colors.String(2)) // Output: "blue"
    
    // Convert string to enum value
    value, err := colors.FromString("green")
    if err == nil {
        fmt.Println(value) // Output: 1
    }
    
    // Get all values and labels
    fmt.Println(colors.All())    // Output: [0 1 2]
    fmt.Println(colors.Labels()) // Output: [red green blue]
}
```

### Using the Wrapper for Marshalling

The `Wrapper` type provides automatic JSON and YAML marshalling:

```go
package main

import (
    "encoding/json"
    "fmt"
    "github.com/gmllt/enum"
)

type Status int

func main() {
    // Create a wrapper with status labels
    status := enum.NewWrapper[Status]("pending", "active", "inactive")
    status.Set(1) // Set to "active"
    
    // JSON marshalling
    jsonData, err := json.Marshal(status)
    if err == nil {
        fmt.Println(string(jsonData)) // Output: "active"
    }
    
    // JSON unmarshalling
    var newStatus enum.Wrapper[Status]
    newStatus.Enum = enum.NewEnum[Status]("pending", "active", "inactive")
    
    err = json.Unmarshal([]byte(`"inactive"`), &newStatus)
    if err == nil {
        fmt.Println(newStatus.Get()) // Output: 2
        fmt.Println(newStatus.String()) // Output: "inactive"
    }
}
```

---

## Advanced Usage

### Custom Types

You can use any integer-based type for your enums:

```go
type Priority int8
type Level uint

priorities := enum.NewEnum[Priority]("low", "medium", "high")
levels := enum.NewEnum[Level]("beginner", "intermediate", "advanced")
```

### Performance Optimization

The library automatically optimizes lookup performance:
- Small enums: Uses linear search for fast access
- Large enums: Uses map-based lookup for efficiency

### Error Handling

```go
colors := enum.NewEnum[int]("red", "green", "blue")

// Invalid string conversion returns error
value, err := colors.FromString("purple")
if err != nil {
    fmt.Println(err) // Output: "invalid value: purple"
}

// Invalid index returns "Invalid(N)" format
fmt.Println(colors.String(99)) // Output: "Invalid(99)"
```

---

## Extending Marshalling

You can easily extend the marshalling functionality by creating custom types that embed the `Wrapper`. The `Wrapper` exposes its `Enum` and `Value` fields publicly, making it simple to implement additional marshalling formats.

### Custom Marshalling with Embedded Wrapper

### Custom Marshalling with Embedded Wrapper

You can create custom types that embed the `Wrapper` to add support for additional formats:

```go
package main

import (
    "encoding/xml"
    "fmt"
    "strings"
    "github.com/gmllt/enum"
)

// CustomStatus embeds the Wrapper and adds XML marshalling
type CustomStatus struct {
    enum.Wrapper[int]
}

// NewCustomStatus creates a new CustomStatus
func NewCustomStatus(labels ...string) CustomStatus {
    return CustomStatus{
        Wrapper: enum.NewWrapper[int](labels...),
    }
}

// MarshalXML implements XML marshalling
func (c CustomStatus) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    return e.EncodeElement(c.String(), start)
}

// UnmarshalXML implements XML unmarshalling
func (c *CustomStatus) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
    var s string
    if err := d.DecodeElement(&s, &start); err != nil {
        return err
    }
    
    value, err := c.Enum.FromString(s)
    if err != nil {
        return err
    }
    
    c.Set(value)
    return nil
}

func main() {
    status := NewCustomStatus("pending", "active", "inactive")
    status.Set(1) // Set to "active"
    
    // XML marshalling works automatically
    xmlData, _ := xml.Marshal(status)
    fmt.Println(string(xmlData)) // Output: <CustomStatus>active</CustomStatus>
}
```

### Custom Format Example

You can also implement completely custom formats:

```go
// StatusWithCustomFormat adds custom marshalling methods
type StatusWithCustomFormat struct {
    enum.Wrapper[int]
}

// MarshalCustom implements custom format marshalling
func (s StatusWithCustomFormat) MarshalCustom() ([]byte, error) {
    // Access the string representation directly
    label := s.String()
    return []byte(fmt.Sprintf("STATUS:%s", strings.ToUpper(label))), nil
}

// UnmarshalCustom implements custom format unmarshalling
func (s *StatusWithCustomFormat) UnmarshalCustom(data []byte) error {
    str := string(data)
    if !strings.HasPrefix(str, "STATUS:") {
        return fmt.Errorf("invalid format")
    }
    
    label := strings.ToLower(strings.TrimPrefix(str, "STATUS:"))
    value, err := s.Enum.FromString(label)
    if err != nil {
        return err
    }
    
    s.Set(value)
    return nil
}
```

### Direct Access Pattern

Since `Wrapper` exposes its fields publicly, you can also work directly with the enum data:

```go
func CustomMarshal[T enum.Value](w enum.Wrapper[T]) ([]byte, error) {
    // Access enum labels and current value directly
    labels := w.Enum.LabelsReadOnly()  // Get read-only access to labels
    currentValue := w.Value            // Get current enum value
    
    if int(currentValue) >= len(labels) {
        return nil, fmt.Errorf("invalid enum value: %d", currentValue)
    }
    
    label := labels[int(currentValue)]
    // Implement your custom format
    return []byte("CUSTOM_" + label), nil
}

func CustomUnmarshal[T enum.Value](w *enum.Wrapper[T], data []byte) error {
    str := string(data)
    if !strings.HasPrefix(str, "CUSTOM_") {
        return fmt.Errorf("invalid format")
    }
    
    label := strings.TrimPrefix(str, "CUSTOM_")
    value, err := w.Enum.FromString(label)
    if err != nil {
        return err
    }
    
    w.Set(value)
    return nil
}
```

This approach gives you full flexibility to implement any marshalling format while maintaining type safety and leveraging the existing enum functionality.

---

## API Reference

### Enum[T] Methods

- `String(v T) string` - Convert enum value to string
- `FromString(s string) (T, error)` - Convert string to enum value
- `All() []T` - Get all enum values
- `Labels() []string` - Get all labels (copy)
- `LabelsReadOnly() []string` - Get all labels (read-only view)

### Wrapper[T] Methods

- `String() string` - Get string representation of current value
- `Get() T` - Get current enum value
- `Set(v T)` - Set enum value
- `All() []T` - Get all enum values
- `Labels() []string` - Get all labels
- `MarshalJSON() ([]byte, error)` - JSON marshalling
- `UnmarshalJSON(b []byte) error` - JSON unmarshalling
- `MarshalYAML() (any, error)` - YAML marshalling
- `UnmarshalYAML(unmarshal func(any) error) error` - YAML unmarshalling

---

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for detailed guidelines on how to contribute to this project, including how to run tests, benchmarks, and submit pull requests.

---

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
