# Omniparser Style Reference for Omniwriter

This document captures observed architecture and engineering style from `omniparser` so `omniwriter` can stay consistent.

Guiding principle for omniwriter implementation: KISS. Keep the design minimal, explicit, and profile-driven.

## 1. Architecture Patterns To Reuse

- Interface-first core contracts (`Schema`, `Transform`, `SchemaHandler`, `Ingester`).
- Extension hooks for behavior injection (`Extension`, custom schema handlers, custom funcs).
- Clear staged lifecycle:
  1. parse/validate schema
  2. build transform/ingester
  3. stream records via `Read()`
- Backward-compatible entry points with optional extensibility parameters.
- Symmetric schema sections and naming conventions.

## 2. Error Model Characteristics

- Strong distinction between recoverable record-level errors and fatal operation errors.
- Central sentinel/type errors (`ErrSchemaNotSupported`, `ErrTransformFailed`).
- Context-aware formatting pattern (input/schema names, position context when possible).
- Errors are practical and direct, often used as stable contract in tests.

Implication for `omniwriter`:
- keep stage-aware typed errors and recovery semantics explicit.
- expose predictable behavior for continuable vs fatal failures.
- include precise schema-path error locations for future GUI editing workflows.

## 3. Streaming & Performance Orientation

- Reader-based APIs (`io.Reader`) with repeated `Read()` consumption.
- Designed to avoid loading whole input where possible.
- Includes benchmarks in performance-sensitive packages.
- Uses allocation reuse patterns internally in parser components.

Implication for `omniwriter`:
- streaming-friendly pipeline by default.
- provide writer-based output APIs as soon as feasible.
- benchmark early for EDI-heavy workloads.

## 4. Test Style & Organization

Observed traits:
- table-driven tests for behavior matrices.
- snapshot-based sample/integration tests (cupaloy snapshots).
- explicit assertions on exact error messages in many cases.
- dedicated benchmark tests adjacent to format readers.
- sample fixtures grouped by format and scenario.

Implication for `omniwriter`:
- use table tests for mapper and emitter logic.
- use golden snapshots for end-to-end format permutations.
- keep fixture directories grouped by source-target pair.
- add schema validation tests for `writer_settings`, `output_declaration`, and extended `transform_declarations` profiles.

## 5. Package & Naming Conventions

- packages are focused and modular (`errs`, `schemahandler`, `transformctx`, `validation`).
- exported APIs are concise; details are internal to subpackages.
- comments explain contract/behavior around interfaces and public methods.

Implication for `omniwriter`:
- preserve small, purpose-specific packages.
- keep public surface minimal; hide internals under `internal/`.
- document behavioral contracts where callers depend on them.
- keep section symmetry for consistency:
  - `parser_settings` <-> `writer_settings`
  - `file_declaration` <-> `output_declaration`
  - keep `transform_declarations` as the mapping DSL

## 6. Pragmatic Design Quirks Worth Preserving

- Contract clarity over abstraction complexity.
- Extensibility favored over hardcoded branching.
- Deterministic outputs used heavily in snapshots/tests.
- Schema-driven behavior emphasized over custom imperative logic.

Implication for `omniwriter`:
- treat mappings/schemas as primary source of behavior.
- avoid embedding format-pair special cases in orchestration.
- keep a single logical transform engine (`transform_declarations`) and separate target emitters for serialization.

## 7. Where Omniwriter Should Intentionally Diverge

`omniparser` currently centers on parse-to-JSON output. `omniwriter` should keep style consistency but extend capabilities:
- multi-target emitter layer (`CSV/JSON/XML/EDI/text/custom`).
- explicit passthrough semantics as first-class behavior.
- broader permutation test matrix with strict TDD iteration gates.
- modern Go toolchain/lint/race hygiene from project start.
- single-output-per-schema via:
  - `writer_settings` containing `version` and target `file_format_type`
  - `output_declaration` containing output physical format options
  - target-specific validation of `transform_declarations`

## 8. Working Rules For Contributors

Before implementing a feature:
1. add failing tests first.
2. verify failure reason is correct.
3. implement minimal passing code.
4. run full relevant tests and race/vet checks.
5. refactor only with tests green.
6. prefer the simplest solution that satisfies current tests and contracts.

No next iteration begins until current tests pass.
