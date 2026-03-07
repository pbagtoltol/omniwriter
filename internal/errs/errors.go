package errs

import "errors"

// Request and validation errors
var (
	ErrInvalidRequest           = errors.New("invalid transform request")
	ErrUnsupportedTransform     = errors.New("unsupported transform")
	ErrUnsupportedWriterFormat  = errors.New("unsupported writer_settings.file_format_type")
	ErrTargetFormatMismatch     = errors.New("target format does not match writer_settings.file_format_type")
	ErrMissingWriterSettings    = errors.New("schema missing writer_settings")
	ErrMissingOutputDeclaration = errors.New("schema missing output_declaration")
	ErrMissingWriterEDIDecl     = errors.New("output_declaration is required for edi output")
	ErrMissingWriterCSVDecl     = errors.New("output_declaration.columns is required for csv output")
	ErrMissingCSVColumns        = errors.New("output_declaration.columns is required for csv output")
	ErrMissingSegmentDecls      = errors.New("transform_declarations.FINAL_OUTPUT.segments is required for edi output")
	ErrMissingFinalOutput       = errors.New("transform_declarations.FINAL_OUTPUT is required")
	ErrMissingSegmentDeclName   = errors.New("edi segment name is required")
)
