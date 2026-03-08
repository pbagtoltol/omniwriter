# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**Omniwriter** is a multi-format data transformation library for Go that extends [omniparser](https://github.com/jf-tech/omniparser) from a JSON-output parser into a comprehensive transformation library supporting conversions between JSON, CSV, EDI, XML, and text formats. It implements a full 5×5 transformation matrix (25 combinations) with schema-driven behavior.

## Core Architecture

### Pipeline Flow
```
Input Reader → Ingest Adapter → Canonical Model → Mapper → Emitter → Output Writer
```

### Key Packages

- **Root `api.go`**: Backward compatibility layer (legacy import path)
- **`internal/types`**: Core request/response types and format constants
- **`internal/pipeline`**: Main orchestration - `Execute()` and `ExecuteToWriter()`
- **`internal/schema`**: Schema parsing and validation (parser_settings, writer_settings, output_declaration)
- **`internal/emit/*`**: Format-specific emitters (csv, json, edi, xml, text)
- **`internal/ingest`**: Wrappers around omniparser for input parsing
- **`internal/mapper`**: Transform logic execution
- **`internal/model`**: Canonical document representation
- **`internal/errs`**: Typed, stage-aware errors

### Schema Structure

Every omniwriter schema has three main sections:

1. **`parser_settings`** (omniparser) - Defines input format
   - `version`: e.g., "omni.2.1"
   - `file_format_type`: csv, json, xml, edi, text

2. **`writer_settings`** (omniwriter addition) - Defines output format
   - `version`: e.g., "omni.1.0"
   - `file_format_type`: csv, json, xml, edi, text

3. **`transform_declarations`** - Field mapping logic using XPath, constants, templates, custom functions

4. **`file_declaration`** (optional) - Required for EDI/CSV input structure
5. **`output_declaration`** (optional) - Required for EDI/CSV output formatting

**Important**: JSON and XML inputs don't require `file_declaration` (self-describing formats). EDI and CSV outputs require `output_declaration` for delimiters/columns.

### Passthrough Optimization

Same-format transforms (CSV→CSV, JSON→JSON, EDI→EDI, XML→XML) use a fast passthrough path that bypasses the full pipeline when no transformation is needed.

## Development Commands

### Testing
```bash
# Run all tests
go test -v ./...

# Run tests without cache
go test -v -count=1 ./...

# Run with race detector
go test -race -short ./...

# Run specific test
go test -v -run TestTransform_JSONToEDI ./...

# Generate coverage report
go test -v -coverprofile=coverage.out -covermode=atomic ./...
go tool cover -func=coverage.out

# View coverage in browser
go tool cover -html=coverage.out
```

### Code Quality
```bash
# Run go vet
go vet ./...

# Check formatting
gofmt -s -l .

# Apply formatting
gofmt -s -w .

# Run staticcheck (if installed)
staticcheck ./...
```

### Dependencies
```bash
# Download dependencies
go mod download

# Verify dependencies
go mod verify

# Tidy dependencies
go mod tidy

# Update dependencies
go get -u ./...
```

### Examples
```bash
# Run an example
cd examples/json_to_csv
go run main.go

# All examples follow same pattern
cd examples/<example_name>
go run main.go
```

## Development Guidelines

### TDD Protocol (Mandatory)

This project follows strict TDD:
1. Write failing test first
2. Implement minimal code to pass
3. Run full test suite
4. Refactor without behavior change
5. Commit only when all tests pass

**Never implement features without tests first.**

### Omniparser Consistency

Mirror omniparser engineering style:
- Interface-first seams for handlers/adapters/emitters
- Explicit context-aware errors with location context
- Streaming-oriented internals, low allocation pressure
- Table-driven tests with snapshot fixtures
- Symmetric schema naming: `parser_settings` ↔ `writer_settings`, `file_declaration` ↔ `output_declaration`

### Style Requirements

- **No emojis** in code, tests, or error messages
- Simple, professional output in examples
- Direct error messages with context
- KISS principle - no unnecessary abstractions

See `docs/STYLE_GUIDE.md` and `docs/OMNIPARSER_STYLE_REFERENCE.md` for details.

## Important Module Dependencies

### Omniparser Version

Currently using `github.com/jf-tech/omniparser v1.0.5` (latest stable release).

**Critical**: Never add `replace` directives for omniparser in `go.mod`. This breaks CI/CD as the local path doesn't exist on other machines.

### Example Module Structure

Examples use `replace` directives pointing to parent omniwriter module (`replace github.com/pbagtoltol/omniwriter => ../..`). This is correct and expected for local example development.

## Common Patterns

### Adding a New Transform

1. Add test case to `transform_test.go` with test schema and input data
2. Create test fixtures in `testdata/<format>/` if needed
3. Ensure test fails initially
4. Implement emitter logic if new output format
5. Update `internal/pipeline/transform.go` switch statement if needed
6. Run tests until passing
7. Add working example in `examples/`

### Schema Validation

All schemas are validated against format-specific requirements. The validation enforces:
- Required sections based on input/output formats
- Target format matches `writer_settings.file_format_type`
- Format-specific declarations (e.g., EDI requires segment definitions)

### Error Handling

Use typed errors from `internal/errs`:
- `ErrInvalidRequest` - Bad request parameters
- `ErrTargetFormatMismatch` - Target format doesn't match schema
- `ErrUnsupportedTransform` - Transform combination not implemented

Errors should include stage context (ingest/map/emit) and location information.

## Testing Strategy

- **Unit tests**: Emitters, validators, mappers
- **Integration tests**: End-to-end transform by format pair in `transform_test.go`
- **Golden snapshots**: Deterministic outputs in `testdata/`
- **Schema validation tests**: In `schema_validation_test.go`

Test files are located at repository root level alongside implementation.

## Documentation References

- `docs/IMPLEMENTATION_PLAN.md` - Full architecture and roadmap
- `docs/SCHEMA_REQUIREMENTS.md` - Schema section requirements by format
- `docs/PROJECT_SUMMARY.md` - Current status and capabilities
- `docs/IMPLEMENTATION_CHECKLIST.md` - Detailed progress tracking
- `README.md` - User-facing documentation

## CI/CD

GitHub Actions runs on push/PR to main/develop:
- Test job: Tests on Go 1.21, 1.22, 1.23 with race detector
- Lint job: go vet, go fmt check, staticcheck
- Coverage job: Coverage report uploaded to codecov

Run CI checks locally before pushing (see Development Commands above).
