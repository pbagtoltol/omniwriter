# ADR-002: writer_settings + output_declaration + Extended transform_declarations Profiles

Status: Accepted

## Context

`omniwriter` extends `omniparser` to support multiple output formats while keeping schema authoring familiar.
We need a simple, consistent schema contract that also supports a future GUI builder.

## Decision

Use:
- `writer_settings` to declare output intent:
  - `version` (start with `omni.1.0`)
  - `file_format_type` (`csv|json|xml|edi|text|custom`)
- `output_declaration` to declare output physical format options (delimiters, escaping, formatting flags).
- `transform_declarations` as the only mapping DSL (extended across output formats).
  - for EDI output, segment layout is declared in `transform_declarations.FINAL_OUTPUT.segments`.

Do not introduce `writer_declaration`.

Execution model:
- source parsing is format-specific via omniparser.
- transform logic is unified in `transform_declarations`.
- output serialization is format-specific via emitters (`csv/json/xml/edi/text/custom`) configured by `output_declaration`.
- this enables full format permutations (e.g. `EDI->CSV`, `XML->EDI`, `CSV->fixed`, `EDI->EDI`) without a separate pair-specific transform DSL.

## Validation Profile Rules (v1)

Global:
- `writer_settings` is required.
- `output_declaration` is required.
- `transform_declarations.FINAL_OUTPUT` is required.

EDI target (`writer_settings.file_format_type = "edi"`):
- `output_declaration` must include required EDI writer options.
- `transform_declarations.FINAL_OUTPUT.segments` must exist.
- each EDI segment must have a non-empty `name`.
- if a segment element value resolves to an array, emit it as a composite element using `output_declaration.component_delimiter` (default `:`).

CSV target (`writer_settings.file_format_type = "csv"`):
- `output_declaration` must include required CSV writer options.
- CSV column declaration must be non-empty.

JSON/XML/TEXT/CUSTOM target:
- no additional v1 writer profile constraints beyond global requirements.

## Consequences

Positive:
- KISS schema design: one mapping DSL and one output selector.
- Strong consistency with omniparser mental model.
- GUI builder can branch forms by `writer_settings.file_format_type`.

Tradeoff:
- format-specific interpretation can make `transform_declarations` semantics broader.

Mitigation:
- keep profile-specific validation strict and explicit.
- keep rendering options in `output_declaration`, not in expression nodes.

## KISS Guardrails

1. Prefer minimum viable validation per target profile.
2. Add rules only when a failing test demonstrates need.
3. Avoid pair-specific runtime branches if profile-driven behavior can cover it.
4. Keep exported API small and stable.
