package schema

import (
	"encoding/json"

	"github.com/pbagtoltol/omniwriter/internal/errs"
)

// Validate performs minimal schema validation for writer_settings and target-specific profiles.
func Validate(schemaBytes []byte) error {
	ws, err := ParseWriterSettings(schemaBytes)
	if err != nil {
		return err
	}
	od, err := ParseOutputDeclaration(schemaBytes)
	if err != nil {
		return err
	}
	if err := validateFinalOutput(schemaBytes); err != nil {
		return err
	}

	switch ws.FileFormatType {
	case "edi":
		if od == nil {
			return errs.ErrMissingWriterEDIDecl
		}
		segs, err := validateEDISegmentsInTransform(schemaBytes)
		if err != nil {
			return err
		}
		for _, seg := range segs {
			if seg.Name == "" {
				return errs.ErrMissingSegmentDeclName
			}
		}
		return nil
	case "csv":
		hasCols, colCount, err := csvColumnsPresence(schemaBytes)
		if err != nil {
			return err
		}
		if !hasCols {
			return errs.ErrMissingWriterCSVDecl
		}
		if colCount == 0 || len(od.Columns) == 0 {
			return errs.ErrMissingCSVColumns
		}
		return nil
	case "json", "xml", "text", "custom":
		return nil
	default:
		return errs.ErrUnsupportedWriterFormat
	}
}

// ParseWriterSettings extracts writer_settings from the schema.
func ParseWriterSettings(schema []byte) (*WriterSettings, error) {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(schema, &raw); err != nil {
		return nil, err
	}
	v, ok := raw["writer_settings"]
	if !ok {
		return nil, errs.ErrMissingWriterSettings
	}
	var ws WriterSettings
	if err := json.Unmarshal(v, &ws); err != nil {
		return nil, err
	}
	if ws.FileFormatType == "" {
		return nil, errs.ErrMissingWriterSettings
	}
	return &ws, nil
}

// ParseOutputDeclaration extracts output_declaration from the schema.
func ParseOutputDeclaration(schema []byte) (*OutputDeclaration, error) {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(schema, &raw); err != nil {
		return nil, err
	}
	v, ok := raw["output_declaration"]
	if !ok {
		return nil, errs.ErrMissingOutputDeclaration
	}
	var od OutputDeclaration
	if err := json.Unmarshal(v, &od); err != nil {
		return nil, err
	}
	return &od, nil
}

// StripWriterFields removes writer_settings and output_declaration for omniparser.
func StripWriterFields(schema []byte) ([]byte, error) {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(schema, &raw); err != nil {
		return nil, err
	}
	delete(raw, "writer_settings")
	delete(raw, "output_declaration")
	return json.Marshal(raw)
}

func validateFinalOutput(schema []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(schema, &raw); err != nil {
		return err
	}
	declBytes, ok := raw["transform_declarations"]
	if !ok {
		return errs.ErrMissingFinalOutput
	}
	var decls map[string]json.RawMessage
	if err := json.Unmarshal(declBytes, &decls); err != nil {
		return err
	}
	if _, ok := decls["FINAL_OUTPUT"]; !ok {
		return errs.ErrMissingFinalOutput
	}
	return nil
}

func validateEDISegmentsInTransform(schema []byte) ([]EDISegment, error) {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(schema, &raw); err != nil {
		return nil, err
	}
	declBytes, ok := raw["transform_declarations"]
	if !ok {
		return nil, errs.ErrMissingFinalOutput
	}
	var decls map[string]json.RawMessage
	if err := json.Unmarshal(declBytes, &decls); err != nil {
		return nil, err
	}
	finalBytes, ok := decls["FINAL_OUTPUT"]
	if !ok {
		return nil, errs.ErrMissingFinalOutput
	}
	var final map[string]json.RawMessage
	if err := json.Unmarshal(finalBytes, &final); err != nil {
		return nil, err
	}
	objBytes, ok := final["object"]
	if !ok {
		return nil, errs.ErrMissingSegmentDecls
	}
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(objBytes, &obj); err != nil {
		return nil, err
	}
	segsField, ok := obj["segments"]
	if !ok {
		return nil, errs.ErrMissingSegmentDecls
	}
	var segContainer map[string]json.RawMessage
	if err := json.Unmarshal(segsField, &segContainer); err != nil {
		return nil, err
	}
	arrBytes, ok := segContainer["array"]
	if !ok {
		return nil, errs.ErrMissingSegmentDecls
	}
	var segEntries []map[string]json.RawMessage
	if err := json.Unmarshal(arrBytes, &segEntries); err != nil {
		return nil, errs.ErrMissingSegmentDecls
	}
	if len(segEntries) == 0 {
		return nil, errs.ErrMissingSegmentDecls
	}
	segs := make([]EDISegment, 0, len(segEntries))
	for _, entry := range segEntries {
		objEntry, ok := entry["object"]
		if !ok {
			return nil, errs.ErrMissingSegmentDecls
		}
		var segObj map[string]json.RawMessage
		if err := json.Unmarshal(objEntry, &segObj); err != nil {
			return nil, errs.ErrMissingSegmentDecls
		}
		nameField, ok := segObj["name"]
		if !ok {
			return nil, errs.ErrMissingSegmentDeclName
		}
		var nameExpr map[string]json.RawMessage
		if err := json.Unmarshal(nameField, &nameExpr); err != nil {
			return nil, errs.ErrMissingSegmentDeclName
		}
		nameConst, ok := nameExpr["const"]
		if !ok {
			return nil, errs.ErrMissingSegmentDeclName
		}
		var name string
		if err := json.Unmarshal(nameConst, &name); err != nil {
			return nil, errs.ErrMissingSegmentDeclName
		}
		segs = append(segs, EDISegment{Name: name})
	}
	return segs, nil
}

func csvColumnsPresence(schema []byte) (bool, int, error) {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(schema, &raw); err != nil {
		return false, 0, err
	}
	odBytes, ok := raw["output_declaration"]
	if !ok {
		return false, 0, errs.ErrMissingOutputDeclaration
	}
	var odMap map[string]json.RawMessage
	if err := json.Unmarshal(odBytes, &odMap); err != nil {
		return false, 0, err
	}
	colsBytes, ok := odMap["columns"]
	if !ok {
		return false, 0, nil
	}
	var cols []json.RawMessage
	if err := json.Unmarshal(colsBytes, &cols); err != nil {
		return true, 0, err
	}
	return true, len(cols), nil
}
