# JSON to Text Transformation Example

This example demonstrates converting JSON customer data to formatted text output using omniwriter.

## Overview

**Input:** JSON file containing an array of customer objects
**Output:** Text file with formatted customer records

## Files

- `input.json` - Sample JSON customer data (3 customers)
- `schema.json` - Transformation schema defining the mapping and template
- `main.go` - Example application code
- `output.txt` - Generated text output (created when you run the example)

## Schema Features

- **Parser**: omni.2.1 JSON parser
- **Writer**: omni.1.0 text writer
- **Transformations**:
  - Extracts customer fields (id, name, email, account type, balance)
  - Applies custom template for formatting each record
  - Separates records with newlines

## Running the Example

```bash
go run main.go
```

## Expected Output

```text
Customer: CUST001 | Name: John Doe | Email: john@example.com | Type: Premium | Balance: $1250.5
Customer: CUST002 | Name: Jane Smith | Email: jane@example.com | Type: Standard | Balance: $450
Customer: CUST003 | Name: Bob Johnson | Email: bob@example.com | Type: Premium | Balance: $2100.75
```

## Use Cases

- Generating human-readable reports from JSON data
- Creating log files or audit trails
- Exporting data for plain text systems
- Formatting data for email notifications or alerts

## Key Points

1. **Template Formatting**: Uses Go template syntax with `{{field}}` placeholders
2. **Record Iteration**: The `/customers/*` XPath iterates over each customer in the array
3. **Type Preservation**: Float values (balance) are preserved as-is
4. **Custom Separators**: Records are separated by newline characters

## Template Syntax

The `output_declaration.template` supports:
- Field placeholders: `{{field_name}}`
- Static text and formatting characters
- Custom record separators via `record_separator`

Note: The current text emitter uses a simple space-separated format. For custom templating, you may need to post-process the output or enhance the emitter.

## Customization

You can modify `schema.json` to:
- Change the template format
- Add or remove fields from the output
- Modify record separators
- Apply XPath transformations to field values

## Related Examples

- `json_to_csv/` - Converting JSON to CSV
- `json_to_xml/` - Converting JSON to XML
- `csv_to_text/` - Converting CSV to text
