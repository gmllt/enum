# Contributing to Enum Library

Thank you for your interest in contributing to the Enum library! This document provides guidelines and instructions for contributing to the project.

---

## Getting Started

### Prerequisites

- Go (latest stable version recommended)
- Git

### Setting Up Your Development Environment

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/your-username/enum.git
   cd enum
   ```
3. Add the original repository as a remote:
   ```bash
   git remote add upstream https://github.com/gmllt/enum.git
   ```
4. Create a new branch for your changes:
   ```bash
   git checkout -b feature-name
   ```

---

## Running Tests

The project has comprehensive test coverage across all packages. Here are the different ways to run tests:

### Run All Tests

```bash
# Run all tests with verbose output
go test -v ./...

# Run tests without verbose output
go test ./...
```

### Run Tests for Specific Package

```bash
# Run tests for main package only
go test -v .

# Run tests for internal package only
go test -v ./internal
```

### Run Specific Tests

```bash
# Run tests matching a pattern
go test -v -run TestEnumString

# Run tests for a specific function
go test -v -run TestNewEnum
```

### Test Coverage

```bash
# Generate coverage report
go test -cover ./...

# Generate detailed coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

## Running Benchmarks

The project includes performance benchmarks to ensure optimal performance:

### Run All Benchmarks

```bash
# Run all benchmarks
go test -bench=.

# Run benchmarks with memory allocation stats
go test -bench=. -benchmem

# Run benchmarks multiple times for stable results
go test -bench=. -count=5
```

### Run Specific Benchmarks

```bash
# Run benchmarks for enum operations
go test -bench=BenchmarkEnum

# Run benchmarks for wrapper operations
go test -bench=BenchmarkWrapper

# Run benchmarks for specific functions
go test -bench=BenchmarkEnum*
```

### Benchmark Results Interpretation

The benchmark output shows:
- Function name and CPU cores used (e.g., `BenchmarkEnumString-24`)
- Number of iterations performed
- Nanoseconds per operation (ns/op)
- Bytes allocated per operation (B/op)
- Number of allocations per operation (allocs/op)

---

## Code Quality

### Linting and Formatting

```bash
# Format code
go fmt ./...

# Run go vet for static analysis
go vet ./...

# Run golangci-lint (if installed)
golangci-lint run
```

### Code Style Guidelines

- Follow standard Go formatting (`go fmt`)
- Write clear, descriptive variable and function names
- Add comments for exported functions and types
- Keep functions focused and small
- Use meaningful test names that describe what is being tested

---

## Testing Guidelines

### Writing Tests

- Write tests for all new functionality
- Use table-driven tests for multiple test cases
- Include edge cases and error conditions
- Test both success and failure paths
- Use descriptive test names

Example test structure:
```go
func TestNewFeature(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected int
        wantErr  bool
    }{
        {
            name:     "valid input",
            input:    "test",
            expected: 0,
            wantErr:  false,
        },
        // Add more test cases...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### Benchmark Guidelines

When adding benchmarks:
- Benchmark critical paths and performance-sensitive operations
- Use `b.ResetTimer()` before the benchmark loop if setup is required
- Avoid allocations in benchmark loops when possible
- Compare before and after performance when making optimizations

Example benchmark:
```go
func BenchmarkNewFeature(b *testing.B) {
    // Setup code here
    input := "test data"
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        // Code to benchmark
        _ = processInput(input)
    }
}
```

---

## Submitting Changes

### Before Submitting

1. **Run all tests**: Make sure all tests pass
   ```bash
   go test ./...
   ```

2. **Run benchmarks**: Ensure no performance regressions
   ```bash
   go test -bench=. -benchmem
   ```

3. **Check formatting**: Format your code
   ```bash
   go fmt ./...
   ```

4. **Update documentation**: Update README.md if needed

### Commit Guidelines

- Write clear, descriptive commit messages
- Keep commits focused on a single change
- Use present tense ("Add feature" not "Added feature")
- Reference issues when applicable

### Pull Request Process

1. Push your changes to your fork
2. Create a pull request against the main branch
3. Provide a clear description of your changes
4. Include any relevant issue numbers
5. Wait for review and address feedback

---

## Performance Considerations

The library is designed for high performance:

- **Lookup Optimization**: Automatically switches between linear search (small enums) and map-based lookup (large enums)
- **Memory Efficiency**: Minimizes allocations in hot paths
- **Caching**: Pre-computes commonly used data structures

When contributing:
- Be mindful of memory allocations
- Run benchmarks to check performance impact
- Consider both small and large enum use cases

---

## Common Tasks

### Adding New Functionality

1. Add the functionality to the appropriate file
2. Write comprehensive tests
3. Add benchmarks if performance-critical
4. Update documentation
5. Run all tests and benchmarks

### Fixing Bugs

1. Write a test that reproduces the bug
2. Fix the bug
3. Ensure the test passes
4. Run all tests to avoid regressions

### Performance Improvements

1. Write benchmarks to measure current performance
2. Implement improvements
3. Run benchmarks to verify improvements
4. Ensure no functionality regressions

---

## Getting Help

- Check existing issues and discussions
- Look at test files for usage examples
- Review the internal package for implementation details

---

## Code of Conduct

- Be respectful and inclusive
- Focus on constructive feedback
- Help others learn and improve
- Follow Go community standards

Thank you for contributing to the Enum library!
