# XML to EDIFACT D96A INVOIC Example

This example demonstrates a complex transformation from XML to EDIFACT D96A INVOIC (Invoice message).

## Schema Structure

Note: This example has **NO file_declaration** because XML is self-describing. Omniparser can automatically parse XML structure and navigate it using XPath.

### Complete Schema Sections

1. **parser_settings** - Input format (XML)
2. **writer_settings** - Output format (EDI)
3. **output_declaration** - EDI delimiters
4. **transform_declarations** - Field mappings with composite elements

file_declaration is **not required** for XML or JSON input.

## Running

```bash
go run main.go
```

## Output

```
UNH+1+INVOIC:D:96A:UN'
BGM+380+INV-2024-001+9'
DTM+137:20240115:102'
NAD+BY+CUST001::91+Acme Corporation'
NAD+SU+SUPP001::91+Global Suppliers Inc'
LIN+1++WIDGET-A:SA'
QTY+47:100'
PRI+AAA:19.99'
MOA+203:1999.00'
```

## Key Features

- Composite elements via JavaScript custom functions
- Segment loops for line items
- EDIFACT D.96A standard compliance
- Qualifiers (NAD+BY, NAD+SU, etc.)
