package schema

// WriterSettings defines the output target configuration.
type WriterSettings struct {
	Version        string `json:"version"`
	FileFormatType string `json:"file_format_type"` // Format type as string to avoid import cycle
}

// OutputDeclaration contains format-specific output options.
type OutputDeclaration struct {
	// EDI options
	SegmentDelimiter    string `json:"segment_delimiter,omitempty"`
	ElementDelimiter    string `json:"element_delimiter,omitempty"`
	ComponentDelimiter  string `json:"component_delimiter,omitempty"`
	RepetitionDelimiter string `json:"repetition_delimiter,omitempty"`
	IgnoreCRLF          bool   `json:"ignore_crlf,omitempty"`

	// CSV options
	Delimiter string      `json:"delimiter,omitempty"`
	Columns   []CSVColumn `json:"columns,omitempty"`
}

// CSVColumn defines a CSV column mapping.
type CSVColumn struct {
	Name  string `json:"name"`
	Path  string `json:"path,omitempty"`
	Const string `json:"const,omitempty"`
}

// EDISegment represents an EDI segment declaration.
type EDISegment struct {
	Name     string       `json:"name"`
	Elements []EDIElement `json:"elements"`
}

// EDIElement represents an EDI element declaration.
type EDIElement struct {
	Path  string `json:"path,omitempty"`
	Const string `json:"const,omitempty"`
}
