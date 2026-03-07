# Complete EDI to CSV Transformation Example

This example demonstrates a **full end-to-end transformation** from EDI (Electronic Data Interchange) to CSV format using omniwriter, showcasing integration with omniparser.

## 🎯 What This Example Shows

1. **Complete Schema Integration**: All five schema sections working together
2. **Real EDI Data**: Canada Post EDI 214 (Transportation Carrier Shipment Status) format
3. **Production-Ready Code**: Error handling, context management, statistics
4. **Omniparser → Omniwriter Pipeline**: Seamless integration between the libraries

## 📁 Files

| File | Purpose |
|------|---------|
| `base_schema.json` | Original omniparser EDI schema (from omniparser samples) |
| `create_schema.go` | Helper to add omniwriter sections to base schema |
| `schema.json` | **Complete schema** with all 5 sections (generated) |
| `input.edi` | Sample EDI 214 shipment status data (19 records) |
| `main.go` | Transformation executable |
| `output.csv` | Generated CSV output |

## 🚀 Quick Start

```bash
# 1. Generate the complete schema (adds writer_settings + output_declaration)
go run create_schema.go

# 2. Run the transformation
go run main.go
```

**Expected Output:**
```
✅ Transformation complete!
Records processed: 19
Output size: 459 bytes
Output written to: output.csv
```

## 📋 Schema Structure Deep Dive

The complete schema (`schema.json`) contains **5 sections**:

### 1. `parser_settings` (Omniparser)
```json
{
  "version": "omni.2.1",
  "file_format_type": "edi"
}
```

### 2. `file_declaration` (Omniparser)
Defines the physical structure of EDI input with envelope and segment declarations, plus delimiters.

### 3. `transform_declarations` (Omniparser)
```json
{
  "FINAL_OUTPUT": {
    "xpath": "ISA/GS/ST",
    "object": {
      "tracking_number": {"xpath": "B10[1]/reference_number"},
      "weight": {"xpath": "W07[1]/weight"}
    }
  }
}
```

### 4. `writer_settings` (Omniwriter)
```json
{
  "version": "omni.1.0",
  "file_format_type": "csv"
}
```

### 5. `output_declaration` (Omniwriter)
```json
{
  "delimiter": ",",
  "columns": [
    {"name": "tracking_number", "path": "tracking_number"},
    {"name": "weight", "path": "weight"}
  ]
}
```

## 🔄 Data Flow

```
EDI Input → [Omniparser Parse] → Parsed Structure →
[Transform] → Canonical JSON → [Omniwriter Emit] → CSV Output
```

## 💡 Key Takeaway

Same `transform_declarations` works with any output format - just change `writer_settings` and `output_declaration`!
