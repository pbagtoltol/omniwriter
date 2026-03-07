# Schema Requirements Reference

## Required vs Optional Sections

### Always Required
1. **`parser_settings`** - Defines input format
2. **`writer_settings`** - Defines output format (omniwriter addition)
3. **`transform_declarations`** - Defines field mappings

### Format-Specific Requirements

#### `file_declaration` (Omniparser Input Parsing)

**Required for:**
- EDI input (must define segment structure)
- CSV with specific delimiters or headers
- Fixed-length files

**Optional/Auto-detected for:**
- JSON (auto-detected structure)
- XML (auto-detected structure)

**Why JSON/XML don't need it:**
Omniparser can automatically parse JSON and XML without explicit structure definitions because they are self-describing formats. The parser can navigate the hierarchy using XPath directly.

**Why EDI/CSV need it:**
EDI and CSV are positional/delimited formats that require explicit structure definitions:
- EDI: Which segments exist, element positions, delimiters
- CSV: Column definitions, delimiter, header presence

#### `output_declaration` (Omniwriter Output Formatting)

**Required for:**
- EDI output (segment/element/component delimiters)
- CSV output (delimiter, column definitions)

**Optional for:**
- JSON output (uses defaults)
- XML output (uses defaults)
- Text output (format-specific)

## Complete Schema Examples

### JSON Input (No file_declaration needed)
```json
{
  "parser_settings": {
    "version": "omni.2.1",
    "file_format_type": "json"
  },
  "writer_settings": {
    "version": "omni.1.0",
    "file_format_type": "csv"
  },
  "output_declaration": {
    "delimiter": ",",
    "columns": [...]
  },
  "transform_declarations": {
    "FINAL_OUTPUT": {...}
  }
}
```

### EDI Input (file_declaration required)
```json
{
  "parser_settings": {
    "version": "omni.2.1",
    "file_format_type": "edi"
  },
  "file_declaration": {
    "segment_delimiter": "~",
    "element_delimiter": "*",
    "segment_declarations": [...]
  },
  "writer_settings": {
    "version": "omni.1.0",
    "file_format_type": "csv"
  },
  "output_declaration": {
    "delimiter": ",",
    "columns": [...]
  },
  "transform_declarations": {
    "FINAL_OUTPUT": {...}
  }
}
```

### XML Input (No file_declaration needed)
```json
{
  "parser_settings": {
    "version": "omni.2.1",
    "file_format_type": "xml"
  },
  "writer_settings": {
    "version": "omni.1.0",
    "file_format_type": "edi"
  },
  "output_declaration": {
    "segment_delimiter": "'",
    "element_delimiter": "+"
  },
  "transform_declarations": {
    "FINAL_OUTPUT": {...}
  }
}
```

## Minimal Schema by Format Combination

| Input → Output | Requires file_declaration | Requires output_declaration |
|---------------|---------------------------|----------------------------|
| JSON → CSV | No | Yes (columns) |
| JSON → EDI | No | Yes (delimiters) |
| JSON → JSON | No | No |
| XML → CSV | No | Yes (columns) |
| XML → EDI | No | Yes (delimiters) |
| XML → JSON | No | No |
| EDI → CSV | Yes (segments) | Yes (columns) |
| EDI → EDI | Yes (segments) | Yes (delimiters) |
| EDI → JSON | Yes (segments) | No |
| CSV → * | Optional | Depends on output |

## Summary

**Rule of thumb:**
- `file_declaration` describes input structure when format is not self-describing
- `output_declaration` describes output formatting options
- `transform_declarations` always required for mapping logic
