package xml

import (
	"encoding/xml"
	"io"
)

// Config holds XML-specific output configuration.
type Config struct {
	Indent         string // Indentation string (e.g., "  " for 2 spaces)
	Prefix         string // Prefix for each line
	RootElement    string // Root element name (optional)
	XMLDeclaration bool   // Include <?xml version="1.0"?>
}

// DefaultConfig returns XML emitter defaults.
func DefaultConfig() Config {
	return Config{
		Indent:         "  ",
		Prefix:         "",
		RootElement:    "",
		XMLDeclaration: true,
	}
}

// Writer wraps xml.Encoder for omniwriter-specific configuration.
type Writer struct {
	w      io.Writer
	enc    *xml.Encoder
	config Config
}

// NewWriter creates a new XML writer with the given config.
func NewWriter(w io.Writer, config Config) *Writer {
	enc := xml.NewEncoder(w)
	enc.Indent(config.Prefix, config.Indent)

	writer := &Writer{
		w:      w,
		enc:    enc,
		config: config,
	}

	// Write XML declaration if requested
	if config.XMLDeclaration {
		w.Write([]byte(xml.Header))
	}

	return writer
}

// WriteRecord writes a single XML record from the canonical payload.
// The payload is expected to be a map that will be converted to XML.
func (w *Writer) WriteRecord(payload map[string]interface{}) error {
	// Convert map to generic XML structure
	elem := mapToXMLElement("record", payload)
	return w.enc.Encode(elem)
}

// Flush ensures all data is written.
func (w *Writer) Flush() error {
	return w.enc.Flush()
}

// mapToXMLElement converts a map to an xml element structure.
func mapToXMLElement(name string, data interface{}) interface{} {
	switch v := data.(type) {
	case map[string]interface{}:
		// Create a generic XML element
		type XMLElement struct {
			XMLName xml.Name
			Attrs   []xml.Attr    `xml:",attr"`
			Content []interface{} `xml:",any"`
		}

		elem := XMLElement{
			XMLName: xml.Name{Local: name},
		}

		// Convert map entries to child elements
		for key, value := range v {
			child := mapToXMLElement(key, value)
			elem.Content = append(elem.Content, child)
		}

		return elem
	case []interface{}:
		// For arrays, wrap each item
		type XMLArray struct {
			XMLName xml.Name
			Items   []interface{} `xml:",any"`
		}

		arr := XMLArray{
			XMLName: xml.Name{Local: name},
		}

		for _, item := range v {
			child := mapToXMLElement("item", item)
			arr.Items = append(arr.Items, child)
		}

		return arr
	default:
		// Primitive value
		type XMLValue struct {
			XMLName xml.Name
			Value   interface{} `xml:",chardata"`
		}

		return XMLValue{
			XMLName: xml.Name{Local: name},
			Value:   v,
		}
	}
}
