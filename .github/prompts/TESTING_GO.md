# Go Testing Best Practices for Frank

## Core Testing Philosophy

- **Test behavior, not implementation**: Focus on _what_ the code does, not _how_ it does it
- **Aim for high coverage, but prioritize critical paths**
- **Integration tests complement unit tests** but should be kept separate
- **Test edge cases and error conditions thoroughly**

<!-- REF: https://go.dev/doc/tutorial/add-a-test -->
<!-- REF: https://github.com/golang/go/wiki/TestComments -->

## Project-Specific Testing Structure

```bash


```

## Test Organization

### Unit Test Structure

```go
func TestFunctionName_Scenario_ExpectedBehavior(t *testing.T) {
    // ARRANGE: Set up test data and expectations
    input := "sample"
    expected := "result"

    // ACT: Call the function being tested
    actual := FunctionBeingTested(input)

    // ASSERT: Verify the results
    if actual != expected {
        t.Errorf("Expected %v but got %v", expected, actual)
    }
}
```

### Table-Driven Tests

Preferred for functions with multiple input/output scenarios:

```go
func TestIsValidEmail(t *testing.T) {
    tests := []struct {
        name  string
        email string
        want  bool
    }{
        {"Valid email", "user@example.com", true},
        {"Missing @", "userexample.com", false},
        {"Missing domain", "user@", false},
        {"Invalid TLD", "user@example.123", false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := pst.IsValidEmail(tt.email)
            if got != tt.want {
                t.Errorf("IsValidEmail(%q) = %v, want %v", tt.email, got, tt.want)
            }
        })
    }
}
```

<!-- REF: https://dave.cheney.net/2019/05/07/prefer-table-driven-tests -->
<!-- REF: https://blog.golang.org/subtests -->

## Mocking and Dependency Injection

Use interfaces for testability:

```go
// Make dependencies explicit through interfaces
type PSTextractor interface {
    ExtractContacts(filePath string, extractAttachments bool, debug bool) ([]Contact, error)
}

// Mock implementation for testing
type MockPSTextractor struct {
    Contacts []Contact
    Err      error
}

func (m *MockPSTextractor) ExtractContacts(filePath string, extractAttachments bool, debug bool) ([]Contact, error) {
    return m.Contacts, m.Err
}

// Test using the mock
func TestProcessFile_WithValidFile(t *testing.T) {
    mockExtractor := &MockPSTextractor{
        Contacts: []Contact{
            {Email: "test@example.com", DisplayName: "Test User"},
        },
        Err: nil,
    }

    processor := NewProcessor(config, mockExtractor)
    // Test processor with the mock
}
```

<!-- REF: https://github.com/golang/mock -->
<!-- REF: https://quii.gitbook.io/learn-go-with-tests/go-fundamentals/dependency-injection -->

## Testing Concurrency

1. **Deterministic Tests**: Avoid timing-dependent tests

```go
// AVOID
func TestConcurrent_Flaky(t *testing.T) {
    go doSomething()
    time.Sleep(100 * time.Millisecond) // Unreliable!
    // Assert something
}

// BETTER
func TestConcurrent_Reliable(t *testing.T) {
    done := make(chan struct{})
    go func() {
        doSomething()
        close(done)
    }()

    select {
    case <-done:
        // Success case
    case <-time.After(1 * time.Second):
        t.Fatal("Test timed out")
    }
}
```

2. **Testing Worker Pools**: Verify behavior, not exact timing

```go
func TestWorkerPool(t *testing.T) {
    // Arrange: Create sample jobs
    jobs := []string{"job1", "job2", "job3"}
    results := make(chan string, len(jobs))

    // Act: Run the worker pool
    processJobs(jobs, results)

    // Assert: Verify all jobs were processed
    var processed []string
    for i := 0; i < len(jobs); i++ {
        processed = append(processed, <-results)
    }

    // Use a helper like testify/assert for more readable assertions
    if len(processed) != len(jobs) {
        t.Errorf("Expected %d processed jobs, got %d", len(jobs), len(processed))
    }
}
```

<!-- REF: https://blog.golang.org/race-detector -->
<!-- REF: https://golang.org/pkg/testing/#hdr-Main -->

## Testing File Operations

Use `io/fs` testing utilities and temporary directories:

```go
func TestFileProcessing(t *testing.T) {
    // Create temp directory that's automatically cleaned up
    tempDir := t.TempDir()

    // Create test files
    testFile := filepath.Join(tempDir, "test.pst")
    err := os.WriteFile(testFile, []byte("test data"), 0644)
    if err != nil {
        t.Fatalf("Failed to create test file: %v", err)
    }

    // Test file processing
    processor := NewProcessor(Config{FolderPath: tempDir})
    results, err := processor.Process()

    // Assertions
    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }
}
```

<!-- REF: https://pkg.go.dev/os#TempDir -->
<!-- REF: https://pkg.go.dev/testing#T.TempDir -->

## Testing Resource Management

Focus on behavior change rather than absolute values:

```go
func TestResourceManager_AdjustsToLowMemory(t *testing.T) {
    // Arrange
    mgr := NewResourceManager(ResourceLimits{
        MemoryPercent: 50.0,
        // Other settings
    })

    // Act - simulate low memory condition
    workersBefore := mgr.CalculateWorkerCount(10)

    // Simulate memory pressure
    mgr.memoryAvailable = mgr.memoryAvailable / 2

    workersAfter := mgr.CalculateWorkerCount(10)

    // Assert - worker count should decrease
    if workersAfter >= workersBefore {
        t.Errorf("Expected worker count to decrease under memory pressure")
    }
}
```

## Recommended Testing Libraries

- **Standard library**: `testing` for most needs
- **Assertion helpers**: `github.com/stretchr/testify/assert` for readable assertions
- **HTTP testing**: `net/http/httptest` for API testing
- **Mock generation**: `github.com/golang/mock/gomock` for interface mocking

<!-- REF: https://pkg.go.dev/github.com/stretchr/testify/assert -->
<!-- REF: https://pkg.go.dev/net/http/httptest -->
<!-- REF: https://pkg.go.dev/github.com/golang/mock/gomock -->

## Common Testing Anti-patterns to Avoid

1. **Testing implementation details** rather than behavior
2. **Brittle tests** that break when internal implementation changes
3. **Slow tests** that access real resources unnecessarily
4. **Global state** that makes tests non-deterministic
5. **Partial assertions** that don't verify complete function behavior
6. **Exposing fields or methods** solely for testing purposes

<!-- REF: https://dave.cheney.net/2016/08/20/solid-go-design -->

## Test Coverage

Use the built-in Go coverage tool:

```bash
# Run tests with coverage
go test -coverprofile=coverage.out ./...

# View coverage in browser
go tool cover -html=coverage.out

# Check coverage percentage
go tool cover -func=coverage.out
```

Aim for high coverage (>80%) of critical components:

- Contact extraction logic
- Error handling paths
- Resource management
- Progress tracking

<!-- REF: https://blog.golang.org/cover -->

## References
