# Omniwriter Style Guide

## Code Style

### General Principles
- Follow omniparser patterns (see OMNIPARSER_STYLE_REFERENCE.md)
- KISS: Keep It Simple, Stupid
- No unnecessary abstractions
- Clear, direct code paths

### Test Output
- **No emojis in test output**
- Simple, clear messages
- Minimal formatting
- Focus on facts, not decoration

**Good:**
```go
fmt.Println("Transformation complete")
fmt.Printf("Records: %d\n", count)
```

**Avoid:**
```go
fmt.Println("✅ Transformation complete!")
fmt.Printf("📊 Records: %d\n", count)
```

### Error Messages
- Direct and informative
- Include context (file, line, field)
- No emojis or excessive formatting
- State the problem clearly

**Good:**
```go
return fmt.Errorf("schema validation failed: missing writer_settings")
return fmt.Errorf("transform failed at record %d: %w", idx, err)
```

**Avoid:**
```go
return fmt.Errorf("❌ Oops! Schema validation failed")
```

### Example Programs
- Simple, professional output
- Clear section headers with plain text
- Box drawing characters acceptable for structure
- Progress information without decoration

**Good:**
```go
fmt.Println("EDI to CSV Transformation")
fmt.Println("========================")
fmt.Printf("Records processed: %d\n", result.Stats.Records)
fmt.Printf("Output: %s\n", outputFile)
```

**Acceptable:**
```go
fmt.Println("╔════════════════════════╗")
fmt.Println("║  EDI to CSV Transform  ║")
fmt.Println("╚════════════════════════╝")
```

**Avoid:**
```go
fmt.Println("🎉 EDI to CSV Transformation 🎉")
fmt.Printf("✅ Success! %d records\n", count)
```

## Documentation Style

### README Files
- Clear, technical documentation
- Code examples prominent
- Use cases explained
- No marketing language

### Comments
- Explain why, not what
- Contract documentation for public APIs
- Omit obvious comments

**Good:**
```go
// StripWriterFields removes omniwriter-specific sections from schema
// so omniparser can validate and process it.
func StripWriterFields(schema []byte) ([]byte, error)
```

**Unnecessary:**
```go
// This function strips writer fields from the schema
func StripWriterFields(schema []byte) ([]byte, error)
```

## Testing Style

### Test Names
- Descriptive, structured names
- Format: `TestFunctionName_Scenario_ExpectedBehavior`

**Good:**
```go
func TestTransform_JSONToCSV_WithNestedObjects(t *testing.T)
func TestValidateSchema_MissingWriterSettings_ReturnsError(t *testing.T)
```

### Test Output
- Use `t.Logf()` for debug information
- Use `t.Errorf()` or `t.Fatalf()` for failures
- No fmt.Println in tests
- No emojis

**Good:**
```go
if got != want {
    t.Errorf("output mismatch:\ngot:  %s\nwant: %s", got, want)
}
```

**Avoid:**
```go
if got != want {
    fmt.Println("❌ Test failed!")
    t.Fail()
}
```

### Table-Driven Tests
- Follow omniparser patterns
- Clear test case names
- Comprehensive coverage

```go
tests := []struct {
    name    string
    input   string
    want    string
    wantErr error
}{
    {
        name: "valid json to csv",
        input: `{"id": 1}`,
        want: "1\n",
    },
    // ...
}
```

## Package Organization

### Public API (pkg/omniwriter)
- Minimal surface area
- Stable contracts
- Clear documentation
- Re-export from internal where appropriate

### Internal Packages
- Focused, single-purpose
- Clear dependencies
- No circular imports

### File Naming
- Lowercase, underscores for separators
- `_test.go` suffix for tests
- Descriptive names

**Good:**
- `transform.go`
- `schema_validation.go`
- `edi_emitter.go`

## Git Commit Messages

### Format
```
<type>: <subject>

<body>

<footer>
```

### Types
- feat: New feature
- fix: Bug fix
- refactor: Code restructuring
- test: Add/update tests
- docs: Documentation
- chore: Maintenance

### Example
```
feat: add XML emitter for Phase 2

Implements internal/emit/xml package with:
- Element escaping
- Attribute handling
- Namespace support

Closes #42
```

## Error Handling

### Sentinel Errors
- Define in internal/errs
- Prefix with `Err`
- Descriptive names

```go
var (
    ErrMissingWriterSettings = errors.New("schema missing writer_settings")
    ErrInvalidFormat = errors.New("unsupported format")
)
```

### Error Wrapping
- Use fmt.Errorf with %w
- Provide context
- Preserve original error

```go
if err := validate(schema); err != nil {
    return fmt.Errorf("schema validation failed: %w", err)
}
```

## Performance

### Benchmarks
- Realistic workloads
- Meaningful names
- Document baseline

```go
func BenchmarkTransform_EDIToCSV_1000Records(b *testing.B) {
    // setup
    for i := 0; i < b.N; i++ {
        // test
    }
}
```

### Allocations
- Reuse buffers where sensible
- Profile before optimizing
- Document allocation budgets

## Documentation Requirements

### Public Functions
- Godoc comments
- Example usage for complex APIs
- Error conditions documented

### Packages
- Package-level documentation
- Clear purpose statement
- Links to related packages

### Examples
- Working, runnable code
- Real-world use cases
- Complete schemas included
- Clear expected output

## Anti-Patterns to Avoid

1. Emojis in code or output
2. Over-engineering solutions
3. Premature optimization
4. Hidden complexity
5. Unclear error messages
6. Magic numbers without constants
7. Deep nesting (max 3-4 levels)
8. Global state
9. Panic in library code
10. Silent failures

## Review Checklist

Before submitting code:
- [ ] Tests pass
- [ ] No emojis in output
- [ ] Clear error messages
- [ ] Documentation updated
- [ ] Follows package conventions
- [ ] No unnecessary complexity
- [ ] Handles errors properly
- [ ] Code is readable
