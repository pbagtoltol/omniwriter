package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	// Read the base omniparser schema
	base, err := os.ReadFile("base_schema.json")
	if err != nil {
		panic(err)
	}

	var schema map[string]interface{}
	if err := json.Unmarshal(base, &schema); err != nil {
		panic(err)
	}

	// Add writer_settings
	schema["writer_settings"] = map[string]interface{}{
		"version":          "omni.1.0",
		"file_format_type": "csv",
	}

	// Add output_declaration
	schema["output_declaration"] = map[string]interface{}{
		"delimiter": ",",
		"columns": []map[string]string{
			{"name": "tracking_number", "path": "tracking_number"},
			{"name": "weight", "path": "weight"},
			{"name": "weight_uom", "path": "weight_uom"},
			{"name": "pickup_date", "path": "pickup_date"},
			{"name": "delivery_date", "path": "delivery_date"},
		},
	}

	// Write the complete schema
	output, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		panic(err)
	}

	if err := os.WriteFile("schema.json", output, 0644)	; err != nil {
		panic(err)
	}

	fmt.Println("✅ Complete schema created: schema.json")
	fmt.Println("\nSchema includes:")
	fmt.Println("  - parser_settings (from omniparser)")
	fmt.Println("  - file_declaration (EDI structure)")
	fmt.Println("  - transform_declarations (field mappings)")
	fmt.Println("  - writer_settings (omniwriter - CSV output)")
	fmt.Println("  - output_declaration (CSV columns)")
}
