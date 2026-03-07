package json

import (
	"bytes"
	"strings"
	"testing"
)

func TestJSONEmitter_Basic(t *testing.T) {
	var buf bytes.Buffer
	config := DefaultConfig()
	writer := NewWriter(&buf, config)

	payload := []byte(`{"id": "123", "name": "Test"}`)

	if err := writer.WriteRecord(payload); err != nil {
		t.Fatalf("WriteRecord failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, `"id"`) {
		t.Error("Expected id field in output")
	}
	if !strings.Contains(output, `"123"`) {
		t.Error("Expected id value in output")
	}
}

func TestJSONEmitter_Multiple(t *testing.T) {
	var buf bytes.Buffer
	config := DefaultConfig()
	writer := NewWriter(&buf, config)

	records := [][]byte{
		[]byte(`{"id": "1"}`),
		[]byte(`{"id": "2"}`),
	}

	for _, rec := range records {
		if err := writer.WriteRecord(rec); err != nil {
			t.Fatalf("WriteRecord failed: %v", err)
		}
	}

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) != 2 {
		t.Errorf("Expected 2 lines, got %d", len(lines))
	}
}

func TestJSONEmitter_Newline(t *testing.T) {
	var buf bytes.Buffer
	config := DefaultConfig()
	writer := NewWriter(&buf, config)

	payload := []byte(`{"id": "456"}`)

	if err := writer.WriteRecord(payload); err != nil {
		t.Fatalf("WriteRecord failed: %v", err)
	}

	output := buf.String()
	if !strings.HasSuffix(output, "\n") {
		t.Error("Expected output to end with newline")
	}
}
