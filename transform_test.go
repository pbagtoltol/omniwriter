package omniwriter

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTransform_CSVToCSVPassthrough_UsesOmniparserSample(t *testing.T) {
	t.Helper()

	samplePath := filepath.Join("..", "omniparser", "extensions", "omniv21", "samples", "csv2", "1_single_row.input.csv")
	input, err := os.ReadFile(samplePath)
	if err != nil {
		t.Fatalf("read sample input: %v", err)
	}

	schema := []byte(`{
		"parser_settings": {"version": "omni.2.1", "file_format_type": "csv2"},
		"writer_settings": {
			"version": "omni.1.0",
			"file_format_type": "csv"
		},
		"output_declaration": {"delimiter": ","},
		"transform_declarations": {"FINAL_OUTPUT": {"object": {"noop": {"const": "x"}}}}
	}`)

	res, err := Transform(context.Background(), TransformRequest{
		SourceFormat: FormatCSV,
		TargetFormat: FormatCSV,
		Mapping:      schema,
		Input:        bytes.NewReader(input),
	})
	if err != nil {
		t.Fatalf("transform: %v", err)
	}

	if got, want := string(res.Output), string(input); got != want {
		t.Fatalf("passthrough mismatch\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

func TestTransform_JSONToEDI_WithWriterSettings(t *testing.T) {
	t.Helper()

	inputPath := filepath.Join("..", "omniparser", "extensions", "omniv21", "samples", "json", "1_single_object.input.json")
	input, err := os.ReadFile(inputPath)
	if err != nil {
		t.Fatalf("read sample input: %v", err)
	}

	schema := []byte(`{
		"parser_settings": {"version": "omni.2.1", "file_format_type": "json"},
		"writer_settings": {
			"version": "omni.1.0",
			"file_format_type": "edi"
		},
		"output_declaration": {
			"segment_delimiter": "~",
			"element_delimiter": "*"
		},
		"transform_declarations": {
			"FINAL_OUTPUT": {"object": {
				"segments": { "array": [
					{ "object": {
						"name": { "const": "HDR" },
						"elements": { "array": [
							{ "xpath": "order_id" },
							{ "xpath": "tracking_number" }
						]}
					}},
					{ "object": {
						"name": { "const": "CNT" },
						"elements": { "array": [ { "const": "2" } ] }
					}}
				]}
			}}
		}
	}`)

	res, err := Transform(context.Background(), TransformRequest{
		SourceFormat: FormatJSON,
		TargetFormat: FormatEDI,
		Mapping:      schema,
		Input:        bytes.NewReader(input),
	})
	if err != nil {
		t.Fatalf("transform: %v", err)
	}

	const want = "HDR*1234567*1z9999999999999999~CNT*2~"
	if got := string(res.Output); got != want {
		t.Fatalf("edi output mismatch\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

func TestTransform_EDIToCSV_UsesOmniparserSample(t *testing.T) {
	t.Helper()

	schemaPath := filepath.Join("..", "omniparser", "extensions", "omniv21", "samples", "edi", "1_canadapost_edi_214.schema.json")
	inputPath := filepath.Join("..", "omniparser", "extensions", "omniv21", "samples", "edi", "1_canadapost_edi_214.input.txt")
	schemaBase, err := os.ReadFile(schemaPath)
	if err != nil {
		t.Fatalf("read sample schema: %v", err)
	}
	input, err := os.ReadFile(inputPath)
	if err != nil {
		t.Fatalf("read sample input: %v", err)
	}

	schema, err := withWriterBlocks(schemaBase, map[string]any{
		"writer_settings": map[string]any{
			"version":          "omni.1.0",
			"file_format_type": "csv",
		},
		"output_declaration": map[string]any{
			"delimiter": "|",
			"columns": []map[string]any{
				{"name": "tracking_number", "path": "tracking_number"},
				{"name": "weight", "path": "weight"},
				{"name": "weight_uom", "path": "weight_uom"},
			},
		},
	})
	if err != nil {
		t.Fatalf("compose schema: %v", err)
	}

	res, err := Transform(context.Background(), TransformRequest{
		SourceFormat: FormatEDI,
		TargetFormat: FormatCSV,
		Mapping:      schema,
		Input:        bytes.NewReader(input),
	})
	if err != nil {
		t.Fatalf("transform: %v", err)
	}

	r := csv.NewReader(strings.NewReader(string(res.Output)))
	r.Comma = '|'
	rows, err := r.ReadAll()
	if err != nil {
		t.Fatalf("read csv output: %v", err)
	}
	if len(rows) < 1 {
		t.Fatalf("expected at least one row, got %d", len(rows))
	}
	gotFirst := rows[0]
	wantFirst := []string{"4343638097845589", "4", "KG"}
	if len(gotFirst) != len(wantFirst) {
		t.Fatalf("unexpected first row width: got %d want %d", len(gotFirst), len(wantFirst))
	}
	for i := range wantFirst {
		if gotFirst[i] != wantFirst[i] {
			t.Fatalf("first row col %d mismatch: got %q want %q", i, gotFirst[i], wantFirst[i])
		}
	}
}

func TestTransform_CSVToEDI_UsesOmniparserSample(t *testing.T) {
	t.Helper()

	schemaPath := filepath.Join("..", "omniparser", "extensions", "omniv21", "samples", "csv2", "1_single_row.schema.json")
	inputPath := filepath.Join("..", "omniparser", "extensions", "omniv21", "samples", "csv2", "1_single_row.input.csv")
	schemaBase, err := os.ReadFile(schemaPath)
	if err != nil {
		t.Fatalf("read sample schema: %v", err)
	}
	input, err := os.ReadFile(inputPath)
	if err != nil {
		t.Fatalf("read sample input: %v", err)
	}

	schema, err := withWriterBlocks(schemaBase, map[string]any{
		"writer_settings": map[string]any{
			"version":          "omni.1.0",
			"file_format_type": "edi",
		},
		"output_declaration": map[string]any{
			"segment_delimiter": "~",
			"element_delimiter": "*",
		},
		"transform_declarations": map[string]any{
			"FINAL_OUTPUT": map[string]any{
				"xpath": ".[DATE != 'N/A']",
				"object": map[string]any{
					"segments": map[string]any{
						"array": []any{
							map[string]any{
								"object": map[string]any{
									"name": map[string]any{"const": "WX"},
									"elements": map[string]any{
										"array": []any{
											map[string]any{"xpath": "DATE"},
											map[string]any{"xpath": "WIND_DIR"},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("compose schema: %v", err)
	}

	res, err := Transform(context.Background(), TransformRequest{
		SourceFormat: FormatCSV,
		TargetFormat: FormatEDI,
		Mapping:      schema,
		Input:        bytes.NewReader(input),
	})
	if err != nil {
		t.Fatalf("transform: %v", err)
	}

	const want = "WX*2019/01/31T12:34:56-0800*N~WX*2020/07/31T01:23:45-0500*SE~WX*2030/11/22T20:18:00-0500*X~"
	if got := string(res.Output); got != want {
		t.Fatalf("edi output mismatch\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

func TestTransform_XMLToEDI_ComplexCompositeElements(t *testing.T) {
	t.Helper()

	inputPath := filepath.Join("..", "omniparser", "extensions", "omniv21", "samples", "xml", "2_multiple_objects.input.xml")
	input, err := os.ReadFile(inputPath)
	if err != nil {
		t.Fatalf("read sample input: %v", err)
	}

	schema := []byte(`{
		"parser_settings": {"version": "omni.2.1", "file_format_type": "xml"},
		"writer_settings": {
			"version": "omni.1.0",
			"file_format_type": "edi"
		},
		"output_declaration": {
			"segment_delimiter": "~",
			"element_delimiter": "*",
			"component_delimiter": ">"
		},
		"transform_declarations": {
			"FINAL_OUTPUT": {"xpath": "lb0:library/lb0:books", "object": {
				"segments": { "array": [
					{ "object": {
						"name": { "const": "HDR" },
						"elements": { "array": [ { "xpath": "header/publisher" } ] }
					}},
					{ "xpath": "book", "object": {
						"name": { "const": "BK" },
						"elements": { "array": [
							{ "xpath": "@title", "keep_empty_or_null": true },
							{ "xpath": "author" },
							{ "custom_func": {
								"name": "javascript",
								"args": [
									{ "const": "[price, year]" },
									{ "const": "price" }, { "xpath": "@price" },
									{ "const": "year" }, { "xpath": "year" }
								]
							}}
						]}
					}},
					{ "object": {
						"name": { "const": "FTR" },
						"elements": { "array": [ { "xpath": "footer" } ] }
					}}
				]}
			}}
		}
	}`)

	res, err := Transform(context.Background(), TransformRequest{
		SourceFormat: FormatXML,
		TargetFormat: FormatEDI,
		Mapping:      schema,
		Input:        bytes.NewReader(input),
	})
	if err != nil {
		t.Fatalf("transform: %v", err)
	}

	const want = "HDR*Scholastic Press~" +
		"BK*Harry Potter and the Philosopher's Stone*J. K. Rowling*9.99>1997~" +
		"BK*Harry Potter and the Chamber of Secrets*J. K. Rowling*10.99>1998~" +
		"FTR*Harry Potter Collection~" +
		"HDR*Harper & Brothers~" +
		"BK*Goodnight Moon*Margaret Wise Brown*5.99>1947~" +
		"BK*Unknown*3.99>1900~" +
		"FTR*Kids Reading Collection~"
	if got := string(res.Output); got != want {
		t.Fatalf("edi output mismatch\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

func TestTransform_EDIToJSON_UsesOmniparserSample(t *testing.T) {
	t.Helper()

	schemaPath := filepath.Join("..", "omniparser", "extensions", "omniv21", "samples", "edi", "1_canadapost_edi_214.schema.json")
	inputPath := filepath.Join("..", "omniparser", "extensions", "omniv21", "samples", "edi", "1_canadapost_edi_214.input.txt")
	schemaBase, err := os.ReadFile(schemaPath)
	if err != nil {
		t.Fatalf("read sample schema: %v", err)
	}
	input, err := os.ReadFile(inputPath)
	if err != nil {
		t.Fatalf("read sample input: %v", err)
	}

	schema, err := withWriterBlocks(schemaBase, map[string]any{
		"writer_settings": map[string]any{
			"version":          "omni.1.0",
			"file_format_type": "json",
		},
		"output_declaration": map[string]any{},
	})
	if err != nil {
		t.Fatalf("compose schema: %v", err)
	}

	res, err := Transform(context.Background(), TransformRequest{
		SourceFormat: FormatEDI,
		TargetFormat: FormatJSON,
		Mapping:      schema,
		Input:        bytes.NewReader(input),
	})
	if err != nil {
		t.Fatalf("transform: %v", err)
	}

	dec := json.NewDecoder(bytes.NewReader(res.Output))
	records := 0
	var first map[string]any
	for {
		var rec map[string]any
		err := dec.Decode(&rec)
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("decode json stream: %v", err)
		}
		if records == 0 {
			first = rec
		}
		records++
	}
	if records != 19 {
		t.Fatalf("expected 19 records, got %d", records)
	}
	if got, ok := first["tracking_number"].(string); !ok || got != "4343638097845589" {
		t.Fatalf("first tracking_number mismatch: got %v", first["tracking_number"])
	}
}

func TestTransform_JSONToJSON_Passthrough(t *testing.T) {
	t.Helper()

	inputPath := filepath.Join("..", "omniparser", "extensions", "omniv21", "samples", "json", "1_single_object.input.json")
	input, err := os.ReadFile(inputPath)
	if err != nil {
		t.Fatalf("read sample input: %v", err)
	}

	schema := []byte(`{
		"parser_settings": {"version": "omni.2.1", "file_format_type": "json"},
		"writer_settings": {"version": "omni.1.0", "file_format_type": "json"}
	}`)

	res, err := Transform(context.Background(), TransformRequest{
		SourceFormat: FormatJSON,
		TargetFormat: FormatJSON,
		Mapping:      schema,
		Input:        bytes.NewReader(input),
	})
	if err != nil {
		t.Fatalf("transform: %v", err)
	}
	if got, want := string(res.Output), string(input); got != want {
		t.Fatalf("passthrough mismatch\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

func TestTransform_EDIToEDI_Passthrough(t *testing.T) {
	t.Helper()

	inputPath := filepath.Join("..", "omniparser", "extensions", "omniv21", "samples", "edi", "1_canadapost_edi_214.input.txt")
	input, err := os.ReadFile(inputPath)
	if err != nil {
		t.Fatalf("read sample input: %v", err)
	}

	schema := []byte(`{
		"parser_settings": {"version": "omni.2.1", "file_format_type": "edi"},
		"writer_settings": {"version": "omni.1.0", "file_format_type": "edi"}
	}`)

	res, err := Transform(context.Background(), TransformRequest{
		SourceFormat: FormatEDI,
		TargetFormat: FormatEDI,
		Mapping:      schema,
		Input:        bytes.NewReader(input),
	})
	if err != nil {
		t.Fatalf("transform: %v", err)
	}
	if got, want := string(res.Output), string(input); got != want {
		t.Fatalf("passthrough mismatch\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

func TestTransform_XMLToXML_Passthrough(t *testing.T) {
	t.Helper()

	inputPath := filepath.Join("..", "omniparser", "extensions", "omniv21", "samples", "xml", "2_multiple_objects.input.xml")
	input, err := os.ReadFile(inputPath)
	if err != nil {
		t.Fatalf("read sample input: %v", err)
	}

	schema := []byte(`{
		"parser_settings": {"version": "omni.2.1", "file_format_type": "xml"},
		"writer_settings": {"version": "omni.1.0", "file_format_type": "xml"}
	}`)

	res, err := Transform(context.Background(), TransformRequest{
		SourceFormat: FormatXML,
		TargetFormat: FormatXML,
		Mapping:      schema,
		Input:        bytes.NewReader(input),
	})
	if err != nil {
		t.Fatalf("transform: %v", err)
	}
	if got, want := string(res.Output), string(input); got != want {
		t.Fatalf("passthrough mismatch\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

func TestTransform_JSONToCSV_UsesOmniparserSample(t *testing.T) {
	t.Helper()

	// Use JSON object with array field (omniparser expects object wrapping)
	input := []byte(`{
		"orders": [
			{"order_id": "1234567", "tracking_number": "1z9999999999999999"},
			{"order_id": "7654321", "tracking_number": "1z8888888888888888"}
		]
	}`)

	schema := []byte(`{
		"parser_settings": {"version": "omni.2.1", "file_format_type": "json"},
		"writer_settings": {
			"version": "omni.1.0",
			"file_format_type": "csv"
		},
		"output_declaration": {
			"delimiter": ",",
			"columns": [
				{"name": "order_id", "path": "order_id"},
				{"name": "tracking_number", "path": "tracking_number"}
			]
		},
		"transform_declarations": {
			"FINAL_OUTPUT": {"xpath": "/orders/*", "object": {
				"order_id": {"xpath": "order_id"},
				"tracking_number": {"xpath": "tracking_number"}
			}}
		}
	}`)

	res, err := Transform(context.Background(), TransformRequest{
		SourceFormat: FormatJSON,
		TargetFormat: FormatCSV,
		Mapping:      schema,
		Input:        bytes.NewReader(input),
	})
	if err != nil {
		t.Fatalf("transform: %v", err)
	}

	r := csv.NewReader(strings.NewReader(string(res.Output)))
	rows, err := r.ReadAll()
	if err != nil {
		t.Fatalf("read csv output: %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(rows))
	}
	want := []string{"1234567", "1z9999999999999999"}
	if len(rows[0]) != len(want) {
		t.Fatalf("unexpected row width: got %d want %d", len(rows[0]), len(want))
	}
	for i := range want {
		if rows[0][i] != want[i] {
			t.Fatalf("row col %d mismatch: got %q want %q", i, rows[0][i], want[i])
		}
	}
}

func TestTransform_CSVToJSON_UsesOmniparserSample(t *testing.T) {
	t.Helper()

	schemaPath := filepath.Join("..", "omniparser", "extensions", "omniv21", "samples", "csv2", "1_single_row.schema.json")
	inputPath := filepath.Join("..", "omniparser", "extensions", "omniv21", "samples", "csv2", "1_single_row.input.csv")
	schemaBase, err := os.ReadFile(schemaPath)
	if err != nil {
		t.Fatalf("read sample schema: %v", err)
	}
	input, err := os.ReadFile(inputPath)
	if err != nil {
		t.Fatalf("read sample input: %v", err)
	}

	schema, err := withWriterBlocks(schemaBase, map[string]any{
		"writer_settings": map[string]any{
			"version":          "omni.1.0",
			"file_format_type": "json",
		},
		"output_declaration": map[string]any{},
	})
	if err != nil {
		t.Fatalf("compose schema: %v", err)
	}

	res, err := Transform(context.Background(), TransformRequest{
		SourceFormat: FormatCSV,
		TargetFormat: FormatJSON,
		Mapping:      schema,
		Input:        bytes.NewReader(input),
	})
	if err != nil {
		t.Fatalf("transform: %v", err)
	}

	dec := json.NewDecoder(bytes.NewReader(res.Output))
	records := 0
	var first map[string]any
	for {
		var rec map[string]any
		err := dec.Decode(&rec)
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("decode json stream: %v", err)
		}
		if records == 0 {
			first = rec
		}
		records++
	}
	if records != 3 {
		t.Fatalf("expected 3 records, got %d", records)
	}
	// Just check that we have some fields (schema dependent)
	if len(first) == 0 {
		t.Fatalf("first record is empty")
	}
}

func TestTransform_XMLToJSON_UsesOmniparserSample(t *testing.T) {
	t.Helper()

	schemaPath := filepath.Join("..", "omniparser", "extensions", "omniv21", "samples", "xml", "2_multiple_objects.schema.json")
	inputPath := filepath.Join("..", "omniparser", "extensions", "omniv21", "samples", "xml", "2_multiple_objects.input.xml")
	schemaBase, err := os.ReadFile(schemaPath)
	if err != nil {
		t.Fatalf("read sample schema: %v", err)
	}
	input, err := os.ReadFile(inputPath)
	if err != nil {
		t.Fatalf("read sample input: %v", err)
	}

	schema, err := withWriterBlocks(schemaBase, map[string]any{
		"writer_settings": map[string]any{
			"version":          "omni.1.0",
			"file_format_type": "json",
		},
		"output_declaration": map[string]any{},
	})
	if err != nil {
		t.Fatalf("compose schema: %v", err)
	}

	res, err := Transform(context.Background(), TransformRequest{
		SourceFormat: FormatXML,
		TargetFormat: FormatJSON,
		Mapping:      schema,
		Input:        bytes.NewReader(input),
	})
	if err != nil {
		t.Fatalf("transform: %v", err)
	}

	dec := json.NewDecoder(bytes.NewReader(res.Output))
	records := 0
	var first map[string]any
	for {
		var rec map[string]any
		err := dec.Decode(&rec)
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("decode json stream: %v", err)
		}
		if records == 0 {
			first = rec
		}
		records++
	}
	if records != 2 {
		t.Fatalf("expected 2 records, got %d", records)
	}
	// Check for nested header.publisher field
	if header, ok := first["header"].(map[string]any); ok {
		if pub, ok := header["publisher"].(string); !ok || pub == "" {
			t.Fatalf("first header.publisher missing or invalid: got %v", header["publisher"])
		}
	} else {
		t.Fatalf("first record missing header field: got %v", first)
	}
}

func TestTransform_XMLToCSV_UsesOmniparserSample(t *testing.T) {
	t.Helper()

	inputPath := filepath.Join("..", "omniparser", "extensions", "omniv21", "samples", "xml", "2_multiple_objects.input.xml")
	input, err := os.ReadFile(inputPath)
	if err != nil {
		t.Fatalf("read sample input: %v", err)
	}

	// Create a simpler schema that iterates books
	schema := []byte(`{
		"parser_settings": {"version": "omni.2.1", "file_format_type": "xml"},
		"writer_settings": {
			"version": "omni.1.0",
			"file_format_type": "csv"
		},
		"output_declaration": {
			"delimiter": ",",
			"columns": [
				{"name": "publisher", "path": "publisher"},
				{"name": "title", "path": "title"},
				{"name": "author", "path": "author"}
			]
		},
		"transform_declarations": {
			"FINAL_OUTPUT": {"xpath": "lb0:library/lb0:books/book", "object": {
				"publisher": {"xpath": "../header/publisher"},
				"title": {"xpath": "@title"},
				"author": {"xpath": "author"}
			}}
		}
	}`)

	res, err := Transform(context.Background(), TransformRequest{
		SourceFormat: FormatXML,
		TargetFormat: FormatCSV,
		Mapping:      schema,
		Input:        bytes.NewReader(input),
	})
	if err != nil {
		t.Fatalf("transform: %v", err)
	}

	r := csv.NewReader(strings.NewReader(string(res.Output)))
	rows, err := r.ReadAll()
	if err != nil {
		t.Fatalf("read csv output: %v", err)
	}
	// XML file has 2 books sections with multiple books each
	if len(rows) < 2 {
		t.Fatalf("expected at least 2 rows, got %d", len(rows))
	}
	// Check that we have data (first row should have publisher)
	if rows[0][0] == "" {
		t.Fatalf("first row publisher is empty, got: %v", rows[0])
	}
}

func TestTransform_JSONToXML_UsesOmniparserSample(t *testing.T) {
	t.Helper()

	input := []byte(`{
		"orders": [
			{"order_id": "1234567", "customer_name": "John Doe", "amount": "99.99"},
			{"order_id": "7654321", "customer_name": "Jane Smith", "amount": "149.99"}
		]
	}`)

	schema := []byte(`{
		"parser_settings": {"version": "omni.2.1", "file_format_type": "json"},
		"writer_settings": {
			"version": "omni.1.0",
			"file_format_type": "xml"
		},
		"output_declaration": {
			"root_element": "orders",
			"record_element": "order"
		},
		"transform_declarations": {
			"FINAL_OUTPUT": {"xpath": "/orders/*", "object": {
				"order_id": {"xpath": "order_id"},
				"customer_name": {"xpath": "customer_name"},
				"amount": {"xpath": "amount"}
			}}
		}
	}`)

	res, err := Transform(context.Background(), TransformRequest{
		SourceFormat: FormatJSON,
		TargetFormat: FormatXML,
		Mapping:      schema,
		Input:        bytes.NewReader(input),
	})
	if err != nil {
		t.Fatalf("transform: %v", err)
	}

	// Verify XML output contains expected elements
	output := string(res.Output)
	if !strings.Contains(output, "<record>") {
		t.Fatalf("expected <record> element in output: %s", output)
	}
	if !strings.Contains(output, "<order_id>1234567</order_id>") {
		t.Fatalf("expected order_id in output: %s", output)
	}
	if !strings.Contains(output, "<customer_name>John Doe</customer_name>") {
		t.Fatalf("expected customer_name in output: %s", output)
	}
}

func TestTransform_CSVToXML_UsesOmniparserSample(t *testing.T) {
	t.Helper()

	schemaPath := filepath.Join("..", "omniparser", "extensions", "omniv21", "samples", "csv2", "1_single_row.schema.json")
	inputPath := filepath.Join("..", "omniparser", "extensions", "omniv21", "samples", "csv2", "1_single_row.input.csv")
	schemaBase, err := os.ReadFile(schemaPath)
	if err != nil {
		t.Fatalf("read sample schema: %v", err)
	}
	input, err := os.ReadFile(inputPath)
	if err != nil {
		t.Fatalf("read sample input: %v", err)
	}

	schema, err := withWriterBlocks(schemaBase, map[string]any{
		"writer_settings": map[string]any{
			"version":          "omni.1.0",
			"file_format_type": "xml",
		},
		"output_declaration": map[string]any{
			"root_element":   "records",
			"record_element": "record",
		},
	})
	if err != nil {
		t.Fatalf("compose schema: %v", err)
	}

	res, err := Transform(context.Background(), TransformRequest{
		SourceFormat: FormatCSV,
		TargetFormat: FormatXML,
		Mapping:      schema,
		Input:        bytes.NewReader(input),
	})
	if err != nil {
		t.Fatalf("transform: %v", err)
	}

	output := string(res.Output)
	if !strings.Contains(output, "<record>") {
		t.Fatalf("expected <record> element in output: %s", output)
	}
	// Check for some expected fields from CSV data
	if !strings.Contains(output, "<date>") {
		t.Fatalf("expected <date> field in output: %s", output)
	}
}

func TestTransform_EDIToXML_UsesOmniparserSample(t *testing.T) {
	t.Helper()

	schemaPath := filepath.Join("..", "omniparser", "extensions", "omniv21", "samples", "edi", "1_canadapost_edi_214.schema.json")
	inputPath := filepath.Join("..", "omniparser", "extensions", "omniv21", "samples", "edi", "1_canadapost_edi_214.input.txt")
	schemaBase, err := os.ReadFile(schemaPath)
	if err != nil {
		t.Fatalf("read sample schema: %v", err)
	}
	input, err := os.ReadFile(inputPath)
	if err != nil {
		t.Fatalf("read sample input: %v", err)
	}

	schema, err := withWriterBlocks(schemaBase, map[string]any{
		"writer_settings": map[string]any{
			"version":          "omni.1.0",
			"file_format_type": "xml",
		},
		"output_declaration": map[string]any{
			"root_element":   "shipments",
			"record_element": "shipment",
		},
	})
	if err != nil {
		t.Fatalf("compose schema: %v", err)
	}

	res, err := Transform(context.Background(), TransformRequest{
		SourceFormat: FormatEDI,
		TargetFormat: FormatXML,
		Mapping:      schema,
		Input:        bytes.NewReader(input),
	})
	if err != nil {
		t.Fatalf("transform: %v", err)
	}

	output := string(res.Output)
	if !strings.Contains(output, "<record>") {
		t.Fatalf("expected <record> element in output: %s", output)
	}
	if !strings.Contains(output, "<tracking_number>") {
		t.Fatalf("expected <tracking_number> element in output: %s", output)
	}
	if !strings.Contains(output, "<weight>") {
		t.Fatalf("expected <weight> element in output: %s", output)
	}
}

func TestTransform_JSONToText_UsesOmniparserSample(t *testing.T) {
	t.Helper()

	input := []byte(`{
		"orders": [
			{"order_id": "1234567", "customer_name": "John Doe"},
			{"order_id": "7654321", "customer_name": "Jane Smith"}
		]
	}`)

	schema := []byte(`{
		"parser_settings": {"version": "omni.2.1", "file_format_type": "json"},
		"writer_settings": {
			"version": "omni.1.0",
			"file_format_type": "text"
		},
		"output_declaration": {
			"template": "Order: {{.order_id}} - Customer: {{.customer_name}}"
		},
		"transform_declarations": {
			"FINAL_OUTPUT": {"xpath": "/orders/*", "object": {
				"order_id": {"xpath": "order_id"},
				"customer_name": {"xpath": "customer_name"}
			}}
		}
	}`)

	res, err := Transform(context.Background(), TransformRequest{
		SourceFormat: FormatJSON,
		TargetFormat: FormatText,
		Mapping:      schema,
		Input:        bytes.NewReader(input),
	})
	if err != nil {
		t.Fatalf("transform: %v", err)
	}

	output := string(res.Output)
	// Text output should contain the data (format may vary based on config)
	if !strings.Contains(output, "1234567") {
		t.Fatalf("expected order_id in output: %s", output)
	}
	if !strings.Contains(output, "John Doe") {
		t.Fatalf("expected customer name in output: %s", output)
	}
}

func TestTransform_CSVToText_UsesOmniparserSample(t *testing.T) {
	t.Helper()

	schemaPath := filepath.Join("..", "omniparser", "extensions", "omniv21", "samples", "csv2", "1_single_row.schema.json")
	inputPath := filepath.Join("..", "omniparser", "extensions", "omniv21", "samples", "csv2", "1_single_row.input.csv")
	schemaBase, err := os.ReadFile(schemaPath)
	if err != nil {
		t.Fatalf("read sample schema: %v", err)
	}
	input, err := os.ReadFile(inputPath)
	if err != nil {
		t.Fatalf("read sample input: %v", err)
	}

	schema, err := withWriterBlocks(schemaBase, map[string]any{
		"writer_settings": map[string]any{
			"version":          "omni.1.0",
			"file_format_type": "text",
		},
		"output_declaration": map[string]any{
			"include_keys": true,
		},
	})
	if err != nil {
		t.Fatalf("compose schema: %v", err)
	}

	res, err := Transform(context.Background(), TransformRequest{
		SourceFormat: FormatCSV,
		TargetFormat: FormatText,
		Mapping:      schema,
		Input:        bytes.NewReader(input),
	})
	if err != nil {
		t.Fatalf("transform: %v", err)
	}

	output := string(res.Output)
	// Text output should contain some data
	if len(output) == 0 {
		t.Fatalf("expected non-empty text output")
	}
}

func TestTransform_EDIToText_UsesOmniparserSample(t *testing.T) {
	t.Helper()

	schemaPath := filepath.Join("..", "omniparser", "extensions", "omniv21", "samples", "edi", "1_canadapost_edi_214.schema.json")
	inputPath := filepath.Join("..", "omniparser", "extensions", "omniv21", "samples", "edi", "1_canadapost_edi_214.input.txt")
	schemaBase, err := os.ReadFile(schemaPath)
	if err != nil {
		t.Fatalf("read sample schema: %v", err)
	}
	input, err := os.ReadFile(inputPath)
	if err != nil {
		t.Fatalf("read sample input: %v", err)
	}

	schema, err := withWriterBlocks(schemaBase, map[string]any{
		"writer_settings": map[string]any{
			"version":          "omni.1.0",
			"file_format_type": "text",
		},
		"output_declaration": map[string]any{
			"template": "Tracking: {{.tracking_number}} Weight: {{.weight}}{{.weight_uom}}",
		},
	})
	if err != nil {
		t.Fatalf("compose schema: %v", err)
	}

	res, err := Transform(context.Background(), TransformRequest{
		SourceFormat: FormatEDI,
		TargetFormat: FormatText,
		Mapping:      schema,
		Input:        bytes.NewReader(input),
	})
	if err != nil {
		t.Fatalf("transform: %v", err)
	}

	output := string(res.Output)
	// Text output should contain tracking numbers
	if !strings.Contains(output, "4343638097845589") {
		t.Fatalf("expected tracking number in output: %s", output)
	}
	// Check output is non-empty
	if len(output) < 10 {
		t.Fatalf("expected substantial text output, got: %s", output)
	}
}

func TestTransform_XMLToText_UsesOmniparserSample(t *testing.T) {
	t.Helper()

	schemaPath := filepath.Join("..", "omniparser", "extensions", "omniv21", "samples", "xml", "2_multiple_objects.schema.json")
	inputPath := filepath.Join("..", "omniparser", "extensions", "omniv21", "samples", "xml", "2_multiple_objects.input.xml")
	schemaBase, err := os.ReadFile(schemaPath)
	if err != nil {
		t.Fatalf("read sample schema: %v", err)
	}
	input, err := os.ReadFile(inputPath)
	if err != nil {
		t.Fatalf("read sample input: %v", err)
	}

	schema, err := withWriterBlocks(schemaBase, map[string]any{
		"writer_settings": map[string]any{
			"version":          "omni.1.0",
			"file_format_type": "text",
		},
		"output_declaration": map[string]any{
			"include_keys": true,
		},
	})
	if err != nil {
		t.Fatalf("compose schema: %v", err)
	}

	res, err := Transform(context.Background(), TransformRequest{
		SourceFormat: FormatXML,
		TargetFormat: FormatText,
		Mapping:      schema,
		Input:        bytes.NewReader(input),
	})
	if err != nil {
		t.Fatalf("transform: %v", err)
	}

	output := string(res.Output)
	// Text output should contain some data
	if len(output) == 0 {
		t.Fatalf("expected non-empty text output")
	}
}

func withWriterBlocks(baseSchema []byte, overlays map[string]any) ([]byte, error) {
	var m map[string]any
	if err := json.Unmarshal(baseSchema, &m); err != nil {
		return nil, err
	}
	for k, v := range overlays {
		m[k] = v
	}
	return json.Marshal(m)
}
