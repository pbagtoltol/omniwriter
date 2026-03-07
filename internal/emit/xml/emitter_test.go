package xml

import (
	"bytes"
	"strings"
	"testing"
)

func TestXMLEmitter_Basic(t *testing.T) {
	var buf bytes.Buffer
	config := DefaultConfig()
	writer := NewWriter(&buf, config)

	payload := map[string]interface{}{
		"id":   "123",
		"name": "Test Item",
	}

	if err := writer.WriteRecord(payload); err != nil {
		t.Fatalf("WriteRecord failed: %v", err)
	}

	if err := writer.Flush(); err != nil {
		t.Fatalf("Flush failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "<?xml version") {
		t.Error("Expected XML declaration")
	}
	if !strings.Contains(output, "<record>") {
		t.Error("Expected <record> element")
	}
	if !strings.Contains(output, "<id>123</id>") {
		t.Error("Expected <id> element with value 123")
	}
	if !strings.Contains(output, "<name>Test Item</name>") {
		t.Error("Expected <name> element with value 'Test Item'")
	}
}

func TestXMLEmitter_NoDeclaration(t *testing.T) {
	var buf bytes.Buffer
	config := Config{
		Indent:         "  ",
		XMLDeclaration: false,
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
	if strings.Contains(output, "<?xml version") {
		t.Error("Did not expect XML declaration")
	}
	if !strings.Contains(output, "<record>") {
		t.Error("Expected <record> element")
	}
}

func TestXMLEmitter_Array(t *testing.T) {
	var buf bytes.Buffer
	config := DefaultConfig()
	writer := NewWriter(&buf, config)

	payload := map[string]interface{}{
		"items": []interface{}{"apple", "banana", "cherry"},
	}

	if err := writer.WriteRecord(payload); err != nil {
		t.Fatalf("WriteRecord failed: %v", err)
	}

	if err := writer.Flush(); err != nil {
		t.Fatalf("Flush failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "<items>") {
		t.Error("Expected <items> element")
	}
	if !strings.Contains(output, "<item>apple</item>") {
		t.Error("Expected array items")
	}
}
