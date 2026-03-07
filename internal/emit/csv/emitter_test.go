package csv

import (
	"bytes"
	"strings"
	"testing"

	"github.com/pbagtoltol/omniwriter/internal/schema"
)

func TestCSVEmitter_Basic(t *testing.T) {
	var buf bytes.Buffer
	config := Config{
		Delimiter: ',',
		Columns: []schema.CSVColumn{
			{Name: "id", Path: "id"},
			{Name: "name", Path: "name"},
		},
	}
	writer := NewWriter(&buf, config)

	payload := map[string]interface{}{
		"id":   "123",
		"name": "Test",
	}

	if err := writer.WriteRecord(payload); err != nil {
		t.Fatalf("WriteRecord failed: %v", err)
	}

	if err := writer.Flush(); err != nil {
		t.Fatalf("Flush failed: %v", err)
	}

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")

	if len(lines) != 1 {
		t.Fatalf("Expected 1 line (data only), got %d", len(lines))
	}

	if lines[0] != "123,Test" {
		t.Errorf("Expected data '123,Test', got: %s", lines[0])
	}
}

func TestCSVEmitter_NoHeader(t *testing.T) {
	var buf bytes.Buffer
	config := Config{
		Delimiter: ',',
		Columns: []schema.CSVColumn{
			{Name: "id", Path: "id"},
		},
	}
	writer := NewWriter(&buf, config)

	payload := map[string]interface{}{
		"id": "456",
	}

	if err := writer.WriteRecord(payload); err != nil {
		t.Fatalf("WriteRecord failed: %v", err)
	}

	if err := writer.Flush(); err != nil {
		t.Fatalf("Flush failed: %v", err)
	}

	output := buf.String()
	if strings.Contains(output, "id") && !strings.HasPrefix(output, "456") {
		t.Error("Output should not contain header when Header is false")
	}
	if !strings.Contains(output, "456") {
		t.Error("Expected data value in output")
	}
}

func TestCSVEmitter_CustomDelimiter(t *testing.T) {
	var buf bytes.Buffer
	config := Config{
		Delimiter: '|',
		Columns: []schema.CSVColumn{
			{Name: "a", Path: "a"},
			{Name: "b", Path: "b"},
		},
	}
	writer := NewWriter(&buf, config)

	payload := map[string]interface{}{
		"a": "1",
		"b": "2",
	}

	if err := writer.WriteRecord(payload); err != nil {
		t.Fatalf("WriteRecord failed: %v", err)
	}

	if err := writer.Flush(); err != nil {
		t.Fatalf("Flush failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "|") {
		t.Error("Expected pipe delimiter in output")
	}
	if !strings.Contains(output, "1|2") {
		t.Errorf("Expected '1|2' in output, got: %s", output)
	}
}
