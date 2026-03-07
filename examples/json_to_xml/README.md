# JSON to XML Transformation Example

This example demonstrates converting JSON order data to XML format using omniwriter.

## Overview

**Input:** JSON file containing an array of order objects
**Output:** XML file with structured order records

## Files

- `input.json` - Sample JSON order data (3 orders)
- `schema.json` - Transformation schema defining the mapping
- `main.go` - Example application code
- `output.xml` - Generated XML output (created when you run the example)

## Schema Features

- **Parser**: omni.2.1 JSON parser
- **Writer**: omni.1.0 XML writer
- **Transformations**:
  - Flattens nested customer data (customer.name → customer_name)
  - Calculates item count using XPath `count()` function
  - Preserves order metadata (ID, date, status, total)

## Running the Example

```bash
go run main.go
```

## Expected Output

```xml
<?xml version="1.0" encoding="UTF-8"?>
<record>
  <customer_email>john@example.com</customer_email>
  <customer_name>John Doe</customer_name>
  <customer_phone>+1-555-0100</customer_phone>
  <item_count>2</item_count>
  <order_date>2024-01-15</order_date>
  <order_id>ORD-001</order_id>
  <status>shipped</status>
  <total>199.93</total>
</record>
<record>
  <customer_email>jane@example.com</customer_email>
  <customer_name>Jane Smith</customer_name>
  <customer_phone>+1-555-0200</customer_phone>
  <item_count>1</item_count>
  <order_date>2024-01-16</order_date>
  <order_id>ORD-002</order_id>
  <status>pending</status>
  <total>150</total>
</record>
<record>
  <customer_email>bob@example.com</customer_email>
  <customer_name>Bob Johnson</customer_name>
  <customer_phone>+1-555-0300</customer_phone>
  <item_count>2</item_count>
  <order_date>2024-01-17</order_date>
  <order_id>ORD-003</order_id>
  <status>delivered</status>
  <total>177.47</total>
</record>
```

## Use Cases

- Converting REST API responses to XML
- Generating XML feeds from JSON data
- Data exchange with XML-based systems
- Order processing and fulfillment systems

## Key Points

1. **Nested Data Handling**: The schema flattens nested JSON structures (e.g., `customer.name` becomes `customer_name`)
2. **XPath Functions**: Uses `count()` to calculate the number of items per order
3. **Array Iteration**: The `/orders/*` XPath iterates over each order in the array
4. **Type Preservation**: Numeric values (total, item_count) are preserved

## Customization

You can modify `schema.json` to:
- Add or remove fields from the output
- Change XPath expressions to extract different data
- Add conditional logic with XPath predicates
- Include item details (currently summarized as count)

## Related Examples

- `csv_to_xml/` - Converting CSV to XML
- `edi_to_xml/` - Converting EDI to XML
- `json_to_csv/` - Converting JSON to CSV
