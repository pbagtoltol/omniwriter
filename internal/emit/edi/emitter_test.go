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

func TestEDIEmitter_RepetitionDelimiter(t *testing.T) {
	var out strings.Builder
	config := Config{
		SegmentDelimiter:    "~",
		ElementDelimiter:    "*",
		ComponentDelimiter:  ":",
		RepetitionDelimiter: "^",
	}

	payload := map[string]interface{}{
		"segments": []interface{}{
			map[string]interface{}{
				"name": "REF",
				"elements": []interface{}{
					"BM",
					// Repeating element: multiple occurrences
					[]interface{}{
						[]interface{}{"123", "456"},
						[]interface{}{"789", "ABC"},
					},
				},
			},
		},
	}

	if err := WriteRecord(&out, config, payload); err != nil {
		t.Fatalf("WriteRecord failed: %v", err)
	}

	output := out.String()
	// REF*BM*123:456^789:ABC~
	expected := "REF*BM*123:456^789:ABC~"
	if output != expected {
		t.Errorf("Expected %q, got %q", expected, output)
	}
}

func TestEDIEmitter_IgnoreCRLF(t *testing.T) {
	var out strings.Builder
	config := Config{
		SegmentDelimiter:   "~",
		ElementDelimiter:   "*",
		ComponentDelimiter: ":",
		IgnoreCRLF:         true,
	}

	payload := map[string]interface{}{
		"segments": []interface{}{
			map[string]interface{}{
				"name": "TEST",
				"elements": []interface{}{
					"Line1\nLine2",
					"Part1\r\nPart2",
					"Normal",
				},
			},
		},
	}

	if err := WriteRecord(&out, config, payload); err != nil {
		t.Fatalf("WriteRecord failed: %v", err)
	}

	output := out.String()
	// Should strip all \r and \n
	expected := "TEST*Line1Line2*Part1Part2*Normal~"
	if output != expected {
		t.Errorf("Expected %q, got %q", expected, output)
	}
}

func TestEDIEmitter_IgnoreCRLFDisabled(t *testing.T) {
	var out strings.Builder
	config := Config{
		SegmentDelimiter:   "~",
		ElementDelimiter:   "*",
		ComponentDelimiter: ":",
		IgnoreCRLF:         false,
	}

	payload := map[string]interface{}{
		"segments": []interface{}{
			map[string]interface{}{
				"name": "TEST",
				"elements": []interface{}{
					"Line1\nLine2",
				},
			},
		},
	}

	if err := WriteRecord(&out, config, payload); err != nil {
		t.Fatalf("WriteRecord failed: %v", err)
	}

	output := out.String()
	// Should keep \n when IgnoreCRLF is false
	expected := "TEST*Line1\nLine2~"
	if output != expected {
		t.Errorf("Expected %q, got %q", expected, output)
	}
}
