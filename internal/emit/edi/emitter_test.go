package edi

import (
	"strings"
	"testing"
)

func TestEDIEmitter_Basic(t *testing.T) {
	var out strings.Builder
	config := Config{
		SegmentDelimiter:   "'",
		ElementDelimiter:   "+",
		ComponentDelimiter: ":",
	}

	payload := map[string]interface{}{
		"segments": []interface{}{
			map[string]interface{}{
				"name": "UNH",
				"elements": []interface{}{
					"1",
					[]interface{}{"ORDERS", "D", "96A", "UN"},
				},
			},
		},
	}

	if err := WriteRecord(&out, config, payload); err != nil {
		t.Fatalf("WriteRecord failed: %v", err)
	}

	output := out.String()
	expected := "UNH+1+ORDERS:D:96A:UN'"
	if output != expected {
		t.Errorf("Expected %q, got %q", expected, output)
	}
}

func TestEDIEmitter_MultipleSegments(t *testing.T) {
	var out strings.Builder
	config := Config{
		SegmentDelimiter:   "'",
		ElementDelimiter:   "+",
		ComponentDelimiter: ":",
	}

	payload := map[string]interface{}{
		"segments": []interface{}{
			map[string]interface{}{
				"name":     "UNH",
				"elements": []interface{}{"1"},
			},
			map[string]interface{}{
				"name":     "BGM",
				"elements": []interface{}{"220"},
			},
		},
	}

	if err := WriteRecord(&out, config, payload); err != nil {
		t.Fatalf("WriteRecord failed: %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "UNH+1'") {
		t.Error("Expected UNH segment")
	}
	if !strings.Contains(output, "BGM+220'") {
		t.Error("Expected BGM segment")
	}
}

func TestEDIEmitter_CompositeElements(t *testing.T) {
	var out strings.Builder
	config := Config{
		SegmentDelimiter:   "'",
		ElementDelimiter:   "+",
		ComponentDelimiter: ":",
	}

	payload := map[string]interface{}{
		"segments": []interface{}{
			map[string]interface{}{
				"name": "DTM",
				"elements": []interface{}{
					[]interface{}{"137", "20240101", "102"},
				},
			},
		},
	}

	if err := WriteRecord(&out, config, payload); err != nil {
		t.Fatalf("WriteRecord failed: %v", err)
	}

	output := out.String()
	expected := "DTM+137:20240101:102'"
	if output != expected {
		t.Errorf("Expected %q, got %q", expected, output)
	}
}

func TestEDIEmitter_EmptyElements(t *testing.T) {
	var out strings.Builder
	config := Config{
		SegmentDelimiter:   "'",
		ElementDelimiter:   "+",
		ComponentDelimiter: ":",
	}

	payload := map[string]interface{}{
		"segments": []interface{}{
			map[string]interface{}{
				"name": "TEST",
				"elements": []interface{}{
					"A",
					"",
					"C",
				},
			},
		},
	}

	if err := WriteRecord(&out, config, payload); err != nil {
		t.Fatalf("WriteRecord failed: %v", err)
	}

	output := out.String()
	expected := "TEST+A++C'"
	if output != expected {
		t.Errorf("Expected %q, got %q", expected, output)
	}
}
