package text

import (
	"bytes"
	"strings"
	"testing"
)

func TestTextEmitter_Basic(t *testing.T) {
	var buf bytes.Buffer
	config := DefaultConfig()
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
	if !strings.Contains(output, "123") {
		t.Error("Expected id value in output")
	}
	if !strings.Contains(output, "Test") {
		t.Error("Expected name value in output")
	}
}

func TestTextEmitter_WithKeys(t *testing.T) {
	var buf bytes.Buffer
	config := Config{
		RecordSeparator: "\n",
		FieldSeparator:  ", ",
		IncludeKeys:     true,
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
	if !strings.Contains(output, "id=456") {
		t.Errorf("Expected 'id=456' in output, got: %s", output)
	}
}

func TestTextEmitter_MultipleRecords(t *testing.T) {
	var buf bytes.Buffer
	config := Config{
		RecordSeparator: "\n",
		FieldSeparator:  " ",
		IncludeKeys:     false,
	}
	writer := NewWriter(&buf, config)

	records := []map[string]interface{}{
		{"value": "first"},
		{"value": "second"},
	}

	for _, rec := range records {
		if err := writer.WriteRecord(rec); err != nil {
			t.Fatalf("WriteRecord failed: %v", err)
		}
	}

	if err := writer.Flush(); err != nil {
		t.Fatalf("Flush failed: %v", err)
	}

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) != 2 {
		t.Errorf("Expected 2 lines, got %d", len(lines))
	}
}

func TestTextEmitter_Template(t *testing.T) {
	var buf bytes.Buffer
	config := Config{
		RecordSeparator: "\n",
		Template:        "ID: {{.id}}, Name: {{.name}}",
	}
	writer := NewWriter(&buf, config)

	payload := map[string]interface{}{
		"id":   "789",
		"name": "Sample",
	}

	if err := writer.WriteRecord(payload); err != nil {
		t.Fatalf("WriteRecord failed: %v", err)
	}

	if err := writer.Flush(); err != nil {
		t.Fatalf("Flush failed: %v", err)
	}

	output := buf.String()
	expected := "ID: 789, Name: Sample"
	if !strings.Contains(output, expected) {
		t.Errorf("Expected %q in output, got: %s", expected, output)
	}
}
