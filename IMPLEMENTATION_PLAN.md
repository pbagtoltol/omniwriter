# Omniwriter Implementation Plan

## 1. Mission

`omniwriter` extends `omniparser` from a JSON-output parser into a multi-target transformation library.

Design philosophy: KISS (Keep It Simple, Stupid) - implement the smallest useful contract first, then extend with tests.

Core objective:
- support `CSV, JSON, XML, EDI, text, custom` as both input and output formats.
- support passthrough and same-format transforms (e.g. CSV -> CSV, EDI -> EDI) as first-class use cases.
- preserve `omniparser` strengths: streaming, schema-driven behavior, extensibility, context-aware errors.
- preserve schema style symmetry with `omniparser` by introducing:
  - `writer_settings` (parallel to `parser_settings`)
  - `output_declaration` (parallel to `file_declaration` for output physical format options)
  - extending `transform_declarations` for mapping logic

## 2. Schema Contract (Primary Design)

Every omniwriter schema keeps omniparser sections and adds:

- `writer_settings`
  - `version`: writer schema version (start with `omni.1.0`)
  - `file_format_type`: target output format (`csv|json|xml|edi|text|custom`)
- `output_declaration`
  - contains output-format physical options only
  - examples:
    - EDI: `segment_delimiter`, `element_delimiter`, `component_delimiter`, etc.
    - CSV: `delimiter`, `replace_double_quotes`, `header`, etc.
- `transform_declarations` (extended)
  - remains the single mapping DSL
  - supports `xpath`, `const`, `array`, `template`, `custom_func`, `type`, and related transform attributes
  - supports multi-target output mapping based on `writer_settings.file_format_type`
  - contains logical target shape, e.g. EDI `segments` layout.

Rules:
1. One schema defines one output format only.
2. Runtime must enforce `request.TargetFormat == writer_settings.file_format_type` (or infer target from schema if omitted).
3. Passthrough is configured via `writer_settings.file_format_type` + `transform_declarations` rules, not by special runtime flags.
4. Output rendering concerns (delimiters, escaping, indentation, segment separators, etc.) live in `output_declaration`, not in mapping expressions.

## 3. Product Scope

### Phase 1 (Foundational)
- canonical pipeline and mapping system.
- input ingest via `omniparser` adapters for existing formats.
- output emitters for `CSV`, `EDI`, `JSON`.
- passthrough for `CSV`, `JSON`, `EDI`.

### Phase 2 (Coverage Expansion)
- output emitters for `XML`, `text`, and custom formats.
- broaden permutation matrix and profile-based validation.

### Phase 3 (Scale & Hardening)
- performance and memory tuning.
- larger fixtures and benchmark/regression suites.

## 4. Architecture

Pipeline:
`Input Reader -> Ingest Adapter -> Canonical Model -> Mapper -> Emitter -> Output Writer`

Core execution contract:
1. ingest any supported source format through omniparser (`parser_settings` + `file_declaration`).
2. execute `transform_declarations` once as the logical transform engine.
3. serialize transformed records with target emitter selected by `writer_settings.file_format_type` and configured by `output_declaration`.

This is not a JSON-only architecture. JSON is just one possible target emitter.

Package layout:
- `pkg/omniwriter`: stable public API.
- `internal/pipeline`: execution orchestration by stage.
- `internal/model`: canonical document model.
- `internal/ingest/omniparser`: wrappers over `omniparser` schema + transform.
- `internal/mapper`: source-to-canonical and canonical-to-target mapping compiler/runtime.
- `internal/emit/csv`
- `internal/emit/json`
- `internal/emit/edi`
- `internal/emit/xml` (phase 2)
- `internal/emit/text` (phase 2)
- `internal/emit/custom` (phase 2)
- `internal/schema`: mapping file types and validation.
- `internal/errs`: typed, stage-aware errors.

## 5. Transformation Matrix Strategy

Target state is full N x N transforms for supported formats.

Delivery model:
1. Build minimal vertical slices (input + mapping + emitter + tests).
2. Add symmetric and passthrough paths early.
3. Reuse shared canonical contracts to avoid pairwise explosion.

Initial committed permutations:
1. `EDI -> CSV`
2. `CSV -> EDI`
3. `JSON -> EDI`
4. `CSV -> CSV` (passthrough)
5. `EDI -> EDI` (passthrough)
6. `JSON -> JSON` (passthrough)

## 6. Canonical Model Principles

- preserve order for segments/fields/columns.
- preserve source metadata for diagnostics (`input name`, record index, checksum, line/segment hints).
- keep typed scalar representation (`string/int/float/bool/date/time/decimal`).
- allow optional raw payload attachment for roundtrip/passthrough fidelity.
- track warnings per record, not just global.

## 7. Public API (Draft)

```go
package omniwriter

type Format string

const (
    FormatCSV    Format = "csv"
    FormatJSON   Format = "json"
    FormatXML    Format = "xml"
    FormatEDI    Format = "edi"
    FormatText   Format = "text"
    FormatCustom Format = "custom"
)

type TransformRequest struct {
    SourceFormat Format
    TargetFormat Format
    Mapping      []byte
    Input        io.Reader
    Options      Options
}

type TransformResult struct {
    Output   []byte
    Warnings []Warning
    Stats    Stats
}

func Transform(ctx context.Context, req TransformRequest) (*TransformResult, error)
```

API behavior rules:
- deterministic output for equivalent input + mapping.
- strict stage attribution for errors (`ingest`, `map`, `emit`).
- `TransformToWriter(...)` streaming variant introduced in phase 2.
- `TransformRequest.TargetFormat` must match `writer_settings.file_format_type` when provided.

## 8. Strict TDD Protocol (Mandatory)

No implementation step proceeds without tests first.

Per-iteration workflow:
1. pick one smallest behavior slice.
2. write/adjust failing tests first.
3. run focused tests and verify failure is expected.
4. implement minimal code to pass tests.
5. run full affected suite (`go test ./...` and targeted race checks).
6. refactor without behavior change.
7. commit only when all tests pass.

Hard gate:
- Do not start the next iteration until current iteration tests pass.

Test pyramid:
- unit: mappers, emitters, validators, errors.
- integration: end-to-end transform by format pair.
- golden snapshots: deterministic outputs for representative fixtures.
- benchmark regression: allocations and throughput checkpoints.

KISS + TDD rule:
- any new abstraction must be justified by a failing test or repeated duplication.

## 9. Omniparser-Consistency Rules

`omniwriter` should mirror `omniparser` engineering style where practical:
- interface-first seams for handlers/adapters/ingesters/emitters.
- explicit context-aware errors with useful location context.
- streaming-oriented internals and low allocation pressure.
- extension-driven architecture instead of hardcoding format logic.
- table-driven tests and snapshot-heavy integration fixtures.
- clear package boundaries and pragmatic comments.
- schema naming and section design should remain symmetric with omniparser:
  - `parser_settings` <-> `writer_settings`
  - `file_declaration` <-> `output_declaration`
  - `transform_declarations` remains the mapping DSL.

See `OMNIPARSER_STYLE_REFERENCE.md` for detailed derived guidance.

## 10. Go Version & Best Practices

- target latest stable Go in your environment (update `go.mod` accordingly).
- enforce modern toolchain checks in CI:
  - `go test ./...`
  - `go test -race ./...` (integration-focused or full based on runtime budget)
  - `go vet ./...`
  - static analysis (`staticcheck`) where available.
- use context propagation, sentinel/type errors judiciously, and avoid hidden panics.
- keep public API minimal and backward-compatible once released.

## 11. GUI Builder Readiness Requirements

The schema must be easy to generate and validate in a future GUI.

Requirements:
1. `transform_declarations` output shape must be explicit and strongly typed by `writer_settings.file_format_type`.
2. all nodes should be machine-discoverable for form generation (no ambiguous polymorphism where avoidable).
3. provide JSON Schema definitions for `writer_settings`, `output_declaration`, and extended `transform_declarations`.
4. validation errors must point to exact schema path for GUI highlighting.
5. keep transform primitives consistent with omniparser mental model to reduce GUI complexity.

## 12. Milestones

### Milestone 0: Bootstrap
- initialize module and folder layout.
- establish CI with mandatory green test gate.
- add ADR-001 canonical model.
- add style reference doc from omniparser analysis.
- add ADR-002 schema contract for `writer_settings` + `output_declaration` + extended `transform_declarations`.

Exit criteria:
- CI green with empty baseline tests and lint/vet checks.

### Milestone 1: Core Engine + First Failing Tests
- define `TransformRequest/Result` and stage error model.
- define and validate `writer_settings`, `output_declaration`, and format-specific `transform_declarations` profiles.
- write failing integration tests for:
  - `JSON -> EDI`
  - `CSV -> CSV` passthrough
- implement minimal pipeline to pass.

Exit criteria:
- both scenarios green in CI.

### Milestone 2: EDI/CSV Bidirectional Slices
- write failing tests for `EDI -> CSV` and `CSV -> EDI`.
- implement canonical mapping + emitters for those paths.
- add deterministic golden snapshots.

Exit criteria:
- all four active permutations green.

### Milestone 3: Same-Format Roundtrip Behavior
- failing tests for `EDI -> EDI` and `JSON -> JSON` passthrough.
- implement normalization policy and fidelity constraints.

Exit criteria:
- passthrough tests green with documented normalization guarantees.

### Milestone 4: Expand Outputs
- failing tests first for `XML` and `text` emitters.
- implement and integrate with canonical model.

Exit criteria:
- representative transforms to XML/text green.

### Milestone 5: Custom Output Plugins
- define emitter plugin interfaces.
- tests for custom format registration and execution.

Exit criteria:
- custom output demo fixture green.

### Milestone 6: Performance & Reliability
- benchmark suites and allocation budgets.
- tighten error diagnostics and docs.

Exit criteria:
- benchmark baseline documented; no regression against agreed threshold.

## 13. First Sprint Task List

1. Create module skeleton and package directories.
2. Add CI pipeline with strict green-test gating.
3. Add `OMNIPARSER_STYLE_REFERENCE.md`.
4. Add ADR-001 for canonical model and passthrough semantics.
5. Add ADR-002 for `writer_settings`, `output_declaration`, and extended `transform_declarations`.
6. Write first failing tests (`JSON -> EDI`, `CSV -> CSV` passthrough).
7. Implement minimal pipeline until tests pass.

## 14. Definition Of Done (Per Story)

A story is complete only if:
- tests were written first and initially failed.
- implementation passes unit + integration tests.
- race/vet checks pass for affected areas.
- docs and fixtures updated.
- no known regression in existing permutation tests.
