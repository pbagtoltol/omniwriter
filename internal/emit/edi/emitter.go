package edi

import (
	"fmt"
	"strings"
)

// Config holds EDI-specific output configuration.
type Config struct {
	SegmentDelimiter    string
	ElementDelimiter    string
	ComponentDelimiter  string
	RepetitionDelimiter string
	IgnoreCRLF          bool
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
	if cfg.RepetitionDelimiter == "" {
		cfg.RepetitionDelimiter = "^"
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

		segName = maybeStripCRLF(segName, cfg.IgnoreCRLF)
		out.WriteString(segName)

		rawElems, _ := segMap["elements"].([]interface{})
		for _, e := range rawElems {
			out.WriteString(cfg.ElementDelimiter)
			out.WriteString(stringifyElement(e, cfg.ComponentDelimiter, cfg.RepetitionDelimiter, cfg.IgnoreCRLF))
		}

		out.WriteString(cfg.SegmentDelimiter)
	}

	return nil
}

// stringifyElement handles EDI element serialization, including composite and repeating elements.
// Elements can be:
// - Simple values: "ABC"
// - Composite (array of values): ["A", "B", "C"] -> "A:B:C"
// - Repeating (array of arrays): [["A", "B"], ["C", "D"]] -> "A:B^C:D"
func stringifyElement(v interface{}, compDelim, repDelim string, ignoreCRLF bool) string {
	switch t := v.(type) {
	case []interface{}:
		// Check if this is a repeating element (array of arrays)
		if len(t) > 0 {
			if _, isArray := t[0].([]interface{}); isArray {
				// This is a repeating element
				repetitions := make([]string, 0, len(t))
				for _, rep := range t {
					repetitions = append(repetitions, stringifyElement(rep, compDelim, repDelim, ignoreCRLF))
				}
				return strings.Join(repetitions, repDelim)
			}
		}
		// This is a composite element
		parts := make([]string, 0, len(t))
		for _, sub := range t {
			parts = append(parts, stringifyElement(sub, compDelim, repDelim, ignoreCRLF))
		}
		return strings.Join(parts, compDelim)
	default:
		s := stringify(v)
		return maybeStripCRLF(s, ignoreCRLF)
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

// maybeStripCRLF removes CR and LF characters from the string if ignoreCRLF is true.
// This is useful when EDI data contains line breaks that should be ignored.
func maybeStripCRLF(s string, ignoreCRLF bool) string {
	if !ignoreCRLF {
		return s
	}
	s = strings.ReplaceAll(s, "\r", "")
	s = strings.ReplaceAll(s, "\n", "")
	return s
}
