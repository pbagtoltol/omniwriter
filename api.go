// Package omniwriter provides backward compatibility wrappers.
// This file provides aliases to the new package structure.
// New code should import github.com/pbagtoltol/omniwriter/pkg/omniwriter directly.
package omniwriter

import (
	api "github.com/pbagtoltol/omniwriter/pkg/omniwriter"
)

// Re-export types from pkg/omniwriter for backward compatibility
type (
	Format          = api.Format
	Options         = api.Options
	Warning         = api.Warning
	Stats           = api.Stats
	TransformRequest = api.TransformRequest
	TransformResult = api.TransformResult
)

// Re-export constants
const (
	FormatCSV    = api.FormatCSV
	FormatJSON   = api.FormatJSON
	FormatXML    = api.FormatXML
	FormatEDI    = api.FormatEDI
	FormatText   = api.FormatText
	FormatCustom = api.FormatCustom
)

// Transform is a wrapper for backward compatibility.
// New code should use pkg/omniwriter.Transform directly.
var Transform = api.Transform
