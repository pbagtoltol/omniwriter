package edi

import (
	"fmt"
	"strings"
)

// Config holds EDI-specific output configuration.
type Config struct {
	SegmentDelimiter   string
	ElementDelimiter   string
	ComponentDelimiter string
}

// DefaultConfig returns EDI emitter defaults per X12 standard.
func DefaultConfig() Config {
	return Config{
		SegmentDelimiter:   "~",
		ElementDelimiter:   "*",
		ComponentDelimiter: ":",
	}
}

// WriteRecord emits a single EDI record from the canonical payload.
// The payload must contain a "segments" array with segment objects.
func WriteRecord(out *strings.Builder, config Config, payload map[string]interface{}) error {
	cfg := config
	if cfg.SegmentDelimiter == "" {
		cfg.SegmentDelimiter = "~"
	}
	if cfg.ElementDelimiter == "" {
		cfg.ElementDelimiter = "*"
	}
	if cfg.ComponentDelimiter == "" {
		cfg.ComponentDelimiter = ":"
	}

	segs, ok := payload["segments"].([]interface{})
	if !ok || len(segs) == 0 {
		return fmt.Errorf("missing or invalid segments in payload")
	}

	for _, rawSeg := range segs {
		segMap, ok := rawSeg.(map[string]interface{})
		if !ok {
			return fmt.Errorf("segment is not an object")
		}

		segName := stringify(segMap["name"])
		if segName == "" {
			return fmt.Errorf("segment missing name")
		}

		out.WriteString(segName)

		rawElems, _ := segMap["elements"].([]interface{})
		for _, e := range rawElems {
			out.WriteString(cfg.ElementDelimiter)
			out.WriteString(stringifyElement(e, cfg.ComponentDelimiter))
		}

		out.WriteString(cfg.SegmentDelimiter)
	}

	return nil
}

// stringifyElement handles EDI element serialization, including composite elements.
func stringifyElement(v interface{}, compDelim string) string {
	switch t := v.(type) {
	case []interface{}:
		parts := make([]string, 0, len(t))
		for _, sub := range t {
			parts = append(parts, stringifyElement(sub, compDelim))
		}
		return strings.Join(parts, compDelim)
	default:
		return stringify(v)
	}
}

// stringify converts a value to string representation.
func stringify(v interface{}) string {
	switch t := v.(type) {
	case nil:
		return ""
	case string:
		return t
	case float64:
		return fmt.Sprintf("%v", t)
	case bool:
		if t {
			return "true"
		}
		return "false"
	default:
		return fmt.Sprintf("%v", t)
	}
}
