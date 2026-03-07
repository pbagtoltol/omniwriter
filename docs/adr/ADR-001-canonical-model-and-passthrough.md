# ADR-001: Canonical Model and Passthrough Semantics

Status: Accepted

## Context

`omniwriter` transforms data between multiple formats (CSV, JSON, XML, EDI, text, custom). We need a clear internal representation that:
1. Enables N × N format transformations without pair-specific logic
2. Preserves essential metadata for diagnostics and error reporting
3. Supports high-fidelity passthrough for same-format operations
4. Maintains streaming-friendly characteristics inherited from omniparser

## Decision

### Canonical Model Structure

The canonical model is an **intermediate JSON representation** produced by omniparser's transform engine. This design leverages existing omniparser infrastructure rather than creating a new internal format.

**Current Phase 1 Model:**
- Records are produced via `omniparser.Transform.Read()` as `[]byte` (JSON)
- Each record is deserialized to `map[string]interface{}` for mapping
- Target emitters serialize from `map[string]interface{}` to output format

**Metadata Preservation (Phase 2+):**
When implementing the full canonical model in `internal/model`, include:
- **Order preservation:** Maintain segment/field/column sequence from source
- **Source metadata:** Input name, record index, line/segment hints
- **Type information:** Preserve typed scalars (string/int/float/bool/date/time/decimal)
- **Raw payload attachment:** Optional for high-fidelity roundtrip
- **Per-record warnings:** Collect transformation warnings at record level, not just globally

**Canonical Model Interface (Future):**
```go
type Record struct {
    // Logical data as nested map/array structure
    Data map[string]interface{}

    // Source tracking
    Metadata RecordMetadata

    // Optional raw source for passthrough fidelity
    RawPayload []byte

    // Record-level warnings
    Warnings []Warning
}

type RecordMetadata struct {
    InputName    string
    RecordIndex  int
    SourceHints  SourceHints  // Line numbers, segment positions, etc.
    Checksum     string       // Optional integrity check
}
```

### Passthrough Semantics

**Design Principle:** Passthrough is a first-class transformation mode, not a special runtime shortcut.

**Phase 1 Implementation (Current):**
Same-format passthrough is detected and handled by simple byte copying:
```go
if req.SourceFormat == target &&
    (target == FormatCSV || target == FormatJSON || target == FormatEDI) {
    return io.ReadAll(req.Input)
}
```

**Normalization Guarantees:**
- **CSV passthrough:** Byte-identical output (no reformatting)
- **JSON passthrough:** Byte-identical output (no parsing/re-serialization)
- **EDI passthrough:** Byte-identical output (no segment re-rendering)

**Fidelity Constraints:**
- Passthrough mode **does not** execute `transform_declarations`
- Passthrough mode **does not** validate against schema structure
- Passthrough mode **only validates** `writer_settings` and basic schema syntax
- Whitespace, formatting, and encoding are preserved exactly

**Phase 2+ Enhancement:**
When full canonical model is implemented, passthrough can optionally:
- Execute `transform_declarations` for filtering/projection
- Normalize formatting (e.g., standardize EDI delimiters)
- Validate structural constraints
- Add metadata enrichment

This will be controlled by `output_declaration.passthrough_mode`:
```json
{
  "output_declaration": {
    "passthrough_mode": "raw",  // or "normalized" or "validated"
    ...
  }
}
```

### Pipeline Architecture

**Current (Phase 1):**
```
Input Reader → Passthrough Check → [if not passthrough]:
  → omniparser Transform → JSON Records →
  → Emitter (EDI/CSV/JSON) → Output Writer
```

**Future (Phase 2+):**
```
Input Reader → Ingest Adapter (omniparser) →
  → Canonical Model → Mapper (transform_declarations) →
  → Emitter (format-specific) → Output Writer
```

**Key Contracts:**
1. **Ingest:** Any supported format → Canonical Model via `omniparser`
2. **Map:** Execute `transform_declarations` once as logical transform
3. **Emit:** Canonical Model → Target format via emitter selected by `writer_settings.file_format_type`

## Consequences

### Positive

1. **Simplicity:** Reuses omniparser's proven JSON transformation engine
2. **Separation of concerns:** Physical format (delimiters, escaping) in `output_declaration`, logical mapping in `transform_declarations`
3. **Extensibility:** New formats only require new emitters, not new mapping logic
4. **Performance:** Passthrough mode avoids unnecessary parse/serialize cycles
5. **Consistency:** Same mental model as omniparser for schema authors

### Tradeoffs

1. **JSON as intermediate format:** All data passes through JSON serialization
   - **Mitigation:** Acceptable for Phase 1; optimize in Phase 2 if benchmarks show bottlenecks
2. **Type fidelity:** JSON's limited type system may lose precision (e.g., decimals, dates)
   - **Mitigation:** Canonical model preserves type hints; emitters handle format-specific serialization
3. **Memory allocation:** Record-by-record deserialization to `map[string]interface{}`
   - **Mitigation:** Streaming architecture limits memory pressure; optimize allocation in Phase 2

### Risks and Mitigations

| Risk | Mitigation |
|------|-----------|
| JSON intermediate format loses type information | Store type metadata in canonical model; emitters use format-specific type coercion |
| Passthrough fidelity issues with whitespace/encoding | Clearly document normalization guarantees; add passthrough mode options in Phase 2 |
| Performance overhead for large datasets | Benchmark early; implement streaming TransformToWriter; optimize hot paths |
| Canonical model becomes too complex | KISS principle: add fields only when tests demonstrate need |

## Implementation Notes

### Phase 1 (Current)
- ✅ Use `omniparser.Transform` to produce JSON records
- ✅ Deserialize each record to `map[string]interface{}`
- ✅ Implement passthrough as byte copy for same-format transforms
- ✅ Keep emitters simple: CSV uses `encoding/csv`, EDI uses string builder

### Phase 2 (Future)
- ⏳ Create `internal/model` package for canonical record structure
- ⏳ Add metadata tracking (source hints, checksums, warnings)
- ⏳ Implement `TransformToWriter(io.Writer)` for streaming output
- ⏳ Add passthrough mode options (`raw`, `normalized`, `validated`)
- ⏳ Optimize allocation patterns based on benchmark data

### Testing Strategy
- **Unit tests:** Canonical model construction, metadata preservation
- **Integration tests:** End-to-end transforms for each format pair
- **Golden snapshots:** Deterministic output verification
- **Passthrough tests:** Byte-identical verification for same-format operations
- **Benchmark tests:** Allocation and throughput budgets

## References

- `IMPLEMENTATION_PLAN.md` Section 6: Canonical Model Principles
- `OMNIPARSER_STYLE_REFERENCE.md` Section 3: Streaming & Performance Orientation
- ADR-002: Schema contract for `writer_settings` and `output_declaration`
