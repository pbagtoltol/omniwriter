# EDI to XML Transformation Example

This example demonstrates converting EDI X12 214 (shipment status) data to XML format using omniwriter.

## Overview

**Input:** EDI X12 214 file with shipment tracking data
**Output:** XML file with structured shipment records

## Files

- `input.edi` - Sample EDI 214 shipment data (1 shipment)
- `schema.json` - Transformation schema defining the mapping
- `main.go` - Example application code
- `output.xml` - Generated XML output (created when you run the example)

## Schema Features

- **Parser**: omni.2.1 EDI parser with segment/element parsing
- **Writer**: omni.1.0 XML writer
- **Transformations**:
  - Extracts shipment tracking number from B10 segment
  - Maps destination address from N4 segment (city, state, zip, country)
  - Captures event data from AT7 and MS1 segments (date, time, location)
  - Extracts weight information from AT8 segment

## Running the Example

```bash
go run main.go
```

## Expected Output

```xml
<?xml version="1.0" encoding="UTF-8"?>
<record>
  <destination_city>New York</destination_city>
  <destination_country>US</destination_country>
  <destination_state>NY</destination_state>
  <destination_zip>10001</destination_zip>
  <event_city>New York</event_city>
  <event_date>20241103</event_date>
  <event_state>NY</event_state>
  <event_time>1200</event_time>
  <tracking_number>TRK001</tracking_number>
  <weight>25.5</weight>
  <weight_unit>KG</weight_unit>
</record>
```

## Use Cases

- Converting EDI shipment status to XML
- Integrating EDI data with XML-based systems
- Tracking shipment events and status updates
- Logistics and supply chain data transformation

## Key Points

1. **Segment Parsing**: EDI segments (ISA, GS, ST, B10, N4, AT7, MS1, AT8) are parsed hierarchically
2. **Element Extraction**: Individual data elements are extracted by index from segments
3. **Target Group**: The `scanInfo` segment group is marked as `is_target` to capture shipment records
4. **Type Conversion**: Weight value is converted to float for proper numeric handling

## EDI Structure

The EDI 214 transaction set follows this hierarchy:
- ISA (Interchange Control Header)
  - GS (Functional Group Header)
    - ST (Transaction Set Header)
    - B10 (Shipment Identification)
    - N1/N3/N4 (Name/Address segments)
    - AT7 (Shipment Status Details)
    - MS1 (Equipment, Shipment or Real Property Location)
    - AT8 (Shipment Weight)
    - SE (Transaction Set Trailer)
  - GE (Functional Group Trailer)
- IEA (Interchange Control Trailer)

## Customization

You can modify `schema.json` to:
- Add more EDI segments and elements
- Change field mappings and names
- Apply custom transformations (date formatting, unit conversion)
- Handle different EDI transaction sets (810, 850, 856, etc.)

## Related Examples

- `json_to_xml/` - Converting JSON to XML
- `csv_to_xml/` - Converting CSV to XML
- `edi_to_csv/` - Converting EDI to CSV
