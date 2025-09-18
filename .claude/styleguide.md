# Go Style Guide

## General Principles
- Write idiomatic Go code
- Prioritize clarity and simplicity
- Follow the official Go Code Review Comments
- Use `gofmt` and `goimports` for formatting

## Naming Conventions

### Packages
- Use lowercase, single-word names
- Avoid underscores or mixedCaps
- Package name should be descriptive but concise
```go
package user     // Good
package userMgmt // Bad
package user_mgmt // Bad
```

### Variables and Functions
- Use camelCase for unexported identifiers
- Use PascalCase for exported identifiers
- Use short but descriptive names in small scopes
- Use descriptive names in larger scopes
```go
// Good
func processOrder(o Order) error { ... }
func (s *Service) GetUserByID(id string) (*User, error) { ... }

// In small scopes
for i, v := range values { ... } // OK
for index, value := range values { ... } // Unnecessary
```

### Constants
- Use PascalCase for exported constants
- Use camelCase for unexported constants
- Group related constants with iota when appropriate
```go
const (
    MaxRetries = 3
    DefaultTimeout = 30 * time.Second
)

type Status int
const (
    StatusPending Status = iota
    StatusActive
    StatusClosed
)
```

### Interfaces
- Name interfaces with -er suffix when possible
- Keep interfaces small and focused
```go
type Reader interface {
    Read([]byte) (int, error)
}

type UserRepository interface {
    GetByID(id string) (*User, error)
    Save(user *User) error
}
```

## Error Handling

### Error Messages
- Error messages should be lowercase
- Don't end with punctuation
- Include context when wrapping errors
```go
// Good
return fmt.Errorf("failed to connect to database: %w", err)

// Bad
return fmt.Errorf("Failed to connect to database.")
```

### Error Checking
- Check errors immediately
- Handle errors explicitly
- Don't ignore errors without comment
```go
// Good
result, err := doSomething()
if err != nil {
    return fmt.Errorf("doSomething failed: %w", err)
}

// If genuinely ignoring
_ = file.Close() // Best effort cleanup
```

## Functions and Methods

### Function Design
- Keep functions small and focused
- Prefer returning early over deep nesting
- Order parameters: Context, then important to less important
```go
// Good
func (s *Service) CreateUser(ctx context.Context, user *User) error {
    if err := s.validate(user); err != nil {
        return err
    }
    
    if err := s.repo.Save(ctx, user); err != nil {
        return fmt.Errorf("failed to save user: %w", err)
    }
    
    return nil
}
```

### Method Receivers
- Use pointer receivers when modifying state
- Be consistent within a type
- Use value receivers for small, immutable types
```go
type User struct {
    ID   string
    Name string
}

// Pointer receiver - modifies state
func (u *User) UpdateName(name string) {
    u.Name = name
}

// Value receiver - doesn't modify state
func (u User) String() string {
    return fmt.Sprintf("User{ID: %s, Name: %s}", u.ID, u.Name)
}
```

## Concurrency

### Goroutines
- Always handle goroutine lifecycle
- Use sync.WaitGroup or channels for synchronization
- Avoid goroutine leaks
```go
// Good
var wg sync.WaitGroup
for _, item := range items {
    wg.Add(1)
    go func(item Item) {
        defer wg.Done()
        process(item)
    }(item)
}
wg.Wait()
```

### Channels
- Make channel purpose clear in naming
- Close channels from sender side
- Use directional channels in function signatures
```go
func producer(out chan<- int) {
    defer close(out)
    for i := 0; i < 10; i++ {
        out <- i
    }
}

func consumer(in <-chan int) {
    for val := range in {
        process(val)
    }
}
```

## Testing

### Test Organization
- Test files end with `_test.go`
- Use table-driven tests for multiple cases
- Name test cases descriptively
```go
func TestCalculateTotal(t *testing.T) {
    tests := []struct {
        name     string
        input    []int
        expected int
    }{
        {
            name:     "empty slice",
            input:    []int{},
            expected: 0,
        },
        {
            name:     "single element",
            input:    []int{5},
            expected: 5,
        },
        {
            name:     "multiple elements",
            input:    []int{1, 2, 3},
            expected: 6,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := CalculateTotal(tt.input)
            if result != tt.expected {
                t.Errorf("got %d, want %d", result, tt.expected)
            }
        })
    }
}
```

### Test Helpers
- Accept `*testing.T` as first parameter
- Use `t.Helper()` to improve error reporting
```go
func setupTestDB(t *testing.T) *sql.DB {
    t.Helper()
    db, err := sql.Open("sqlite3", ":memory:")
    if err != nil {
        t.Fatalf("failed to open test database: %v", err)
    }
    return db
}
```

## Comments and Documentation

### Package Comments
- Every package should have a package comment
- For main packages, describe the program
```go
// Package user provides functionality for user management,
// including creation, authentication, and authorization.
package user
```

### Exported Identifiers
- All exported types, functions, and methods need comments
- Start with the identifier name
```go
// User represents a system user with authentication details.
type User struct {
    ID       string
    Username string
    Email    string
}

// Authenticate verifies the user's credentials and returns
// an authentication token if successful.
func (u *User) Authenticate(password string) (string, error) {
    // ...
}
```

## Performance

### Optimization
- Measure before optimizing
- Use benchmarks to validate improvements
- Prefer clear code over premature optimization
```go
func BenchmarkProcess(b *testing.B) {
    data := generateTestData()
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        process(data)
    }
}
```

### Memory Management
- Reuse allocations where sensible
- Use sync.Pool for frequently allocated objects
- Be aware of slice capacity
```go
// Reuse slice
buffer := make([]byte, 0, 1024)
for _, data := range items {
    buffer = buffer[:0] // Reset length, keep capacity
    buffer = append(buffer, data...)
    process(buffer)
}
```