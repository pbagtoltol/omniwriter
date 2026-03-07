package types

import "io"

// Format is the source or target file format.
type Format string

const (
	FormatCSV    Format = "csv"
	FormatJSON   Format = "json"
	FormatXML    Format = "xml"
	FormatEDI    Format = "edi"
	FormatText   Format = "text"
	FormatCustom Format = "custom"
)

// Options contains optional parameters for transformation.
type Options struct{}

// Warning represents a non-fatal transformation warning.
type Warning struct {
	Message string
}

// Stats contains transformation statistics.
type Stats struct {
	Records int
}

// TransformRequest defines one transformation operation.
type TransformRequest struct {
	SourceFormat Format
	TargetFormat Format
	Mapping      []byte
	Input        io.Reader
	Options      Options
}

// TransformResult contains transformed output and metadata.
type TransformResult struct {
	Output   []byte
	Warnings []Warning
	Stats    Stats
}
