package omniwriter

import (
	"errors"
	"testing"

	"github.com/pbagtoltol/omniwriter/internal/errs"
	"github.com/pbagtoltol/omniwriter/internal/schema"
)

func TestValidateSchema_TargetProfiles(t *testing.T) {
	tests := []struct {
		name    string
		schema  string
		wantErr error
	}{
		{
			name: "json profile valid with FINAL_OUTPUT",
			schema: `{
				"parser_settings": {"version": "omni.2.1", "file_format_type": "json"},
				"writer_settings": {"version": "omni.1.0", "file_format_type": "json"},
				"output_declaration": {},
				"transform_declarations": {"FINAL_OUTPUT": {"object": {"id": {"const": "1"}}}}
			}`,
			wantErr: nil,
		},
		{
			name: "edi profile requires output_declaration.edi",
			schema: `{
				"parser_settings": {"version": "omni.2.1", "file_format_type": "json"},
				"writer_settings": {"version": "omni.1.0", "file_format_type": "edi"},
				"output_declaration": {},
				"transform_declarations": {"FINAL_OUTPUT": {"object": {"id": {"const": "1"}}}}
			}`,
			wantErr: errs.ErrMissingSegmentDecls,
		},
		{
			name: "edi profile requires non-empty segment name",
			schema: `{
				"parser_settings": {"version": "omni.2.1", "file_format_type": "json"},
				"writer_settings": {
					"version": "omni.1.0",
					"file_format_type": "edi"
				},
				"output_declaration": {
					"segment_delimiter": "~",
					"element_delimiter": "*"
				},
				"transform_declarations": {"FINAL_OUTPUT": {"object": {
					"segments": {"array": [{"object": {
						"name": {"const": ""},
						"elements": {"array": [{"const": "x"}]}
					}}]}
				}}}
			}`,
			wantErr: errs.ErrMissingSegmentDeclName,
		},
		{
			name: "csv profile requires output_declaration.csv",
			schema: `{
				"parser_settings": {"version": "omni.2.1", "file_format_type": "json"},
				"writer_settings": {"version": "omni.1.0", "file_format_type": "csv"},
				"output_declaration": {},
				"transform_declarations": {"FINAL_OUTPUT": {"object": {"id": {"const": "1"}}}}
			}`,
			wantErr: errs.ErrMissingWriterCSVDecl,
		},
		{
			name: "csv profile requires columns",
			schema: `{
				"parser_settings": {"version": "omni.2.1", "file_format_type": "json"},
				"writer_settings": {
					"version": "omni.1.0",
					"file_format_type": "csv"
				},
				"output_declaration": {"delimiter": ",", "columns": []},
				"transform_declarations": {"FINAL_OUTPUT": {"object": {"id": {"const": "1"}}}}
			}`,
			wantErr: errs.ErrMissingCSVColumns,
		},
		{
			name: "missing FINAL_OUTPUT fails",
			schema: `{
				"parser_settings": {"version": "omni.2.1", "file_format_type": "json"},
				"writer_settings": {"version": "omni.1.0", "file_format_type": "json"},
				"output_declaration": {},
				"transform_declarations": {"OTHER": {"object": {"id": {"const": "1"}}}}
			}`,
			wantErr: errs.ErrMissingFinalOutput,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := schema.Validate([]byte(tc.schema))
			if tc.wantErr == nil {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("expected %v, got %v", tc.wantErr, err)
			}
		})
	}
}
