package omniwriter

import (
	"github.com/pbagtoltol/omniwriter/internal/types"
)

// Re-export types from internal/types for public API
type (
	Format           = types.Format
	Options          = types.Options
	Warning          = types.Warning
	Stats            = types.Stats
	TransformRequest = types.TransformRequest
	TransformResult  = types.TransformResult
)

// Re-export format constants
const (
	FormatCSV    = types.FormatCSV
	FormatJSON   = types.FormatJSON
	FormatXML    = types.FormatXML
	FormatEDI    = types.FormatEDI
	FormatText   = types.FormatText
	FormatCustom = types.FormatCustom
)
