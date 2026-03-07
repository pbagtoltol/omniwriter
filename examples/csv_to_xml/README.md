# CSV to XML Transformation Example

This example demonstrates converting CSV order data to XML format using omniwriter.

## Overview

**Input:** CSV file with order records
**Output:** XML file with structured order elements

## Files

- `input.csv` - Sample CSV order data (3 orders)
- `schema.json` - Transformation schema defining the mapping
- `main.go` - Example application code
- `output.xml` - Generated XML output (created when you run the example)

## Schema Features

- **Parser**: omni.2.1 CSV parser with comma delimiter
- **Writer**: omni.1.0 XML writer
- **Transformations**:
  - Direct field mapping from CSV columns to XML elements
  - Preserves order metadata (ID, customer info, product, quantity, total)

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
  <order_id>ORD-001</order_id>
  <product_sku>WIDGET-A</product_sku>
  <quantity>5</quantity>
  <total>99.95</total>
</record>
<record>
  <customer_email>jane@example.com</customer_email>
  <customer_name>Jane Smith</customer_name>
  <order_id>ORD-002</order_id>
  <product_sku>GADGET-B</product_sku>
  <quantity>2</quantity>
  <total>99.98</total>
</record>
<record>
  <customer_email>bob@example.com</customer_email>
  <customer_name>Bob Johnson</customer_name>
  <order_id>ORD-003</order_id>
  <product_sku>WIDGET-C</product_sku>
  <quantity>10</quantity>
  <total>150</total>
</record>
```

## Use Cases

- Converting CSV exports to XML
- Generating XML feeds from spreadsheet data
- Data exchange with XML-based systems
- Order processing and fulfillment systems

## Key Points

1. **CSV Parsing**: Automatically detects header row and maps columns to fields
2. **Field Mapping**: Each CSV column is mapped directly to an XML element via XPath
3. **Type Preservation**: Numeric values (quantity, total) are preserved as-is
4. **Row Iteration**: Each CSV row becomes an XML record element

## Customization

You can modify `schema.json` to:
- Change the delimiter (e.g., pipe `|` or tab `\t`)
- Add or remove field mappings
- Apply XPath transformations to field values
- Configure XML element names

## Related Examples

- `json_to_xml/` - Converting JSON to XML
- `edi_to_xml/` - Converting EDI to XML
- `csv_to_json/` - Converting CSV to JSON
