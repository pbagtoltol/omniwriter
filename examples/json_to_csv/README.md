# JSON to CSV Transformation Example

This example demonstrates transforming JSON data to CSV format, useful for:
- Exporting API responses to spreadsheets
- Converting JSON logs to CSV for analysis
- Data migration from NoSQL to relational databases

## Quick Start

```bash
go run main.go
```

## Files

- `schema.json` - Transformation schema
- `input.json` - Sample JSON data (customer orders)
- `main.go` - Executable
- `output.csv` - Generated CSV

## Schema Highlights

### Input (JSON)
```json
{
  "order_id": "ORD-001",
  "customer": {
    "name": "John Doe",
    "email": "john@example.com"
  },
  "items": [
    {"sku": "WIDGET-A", "quantity": 5, "price": 19.99}
  ],
  "total": 99.95
}
```

### Output (CSV)
```csv
order_id,customer_name,customer_email,item_count,total
ORD-001,John Doe,john@example.com,1,99.95
```

## Key Features

- Nested object access via XPath (`customer/name`)
- Array aggregation (`count(items)`)
- Automatic type conversion
- Column ordering control
