package pipeline

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/jf-tech/omniparser"
	"github.com/jf-tech/omniparser/transformctx"

	csvemit "github.com/pbagtoltol/omniwriter/internal/emit/csv"
	ediemit "github.com/pbagtoltol/omniwriter/internal/emit/edi"
	jsonemit "github.com/pbagtoltol/omniwriter/internal/emit/json"
	textemit "github.com/pbagtoltol/omniwriter/internal/emit/text"
	xmlemit "github.com/pbagtoltol/omniwriter/internal/emit/xml"
	"github.com/pbagtoltol/omniwriter/internal/errs"
	"github.com/pbagtoltol/omniwriter/internal/schema"
	"github.com/pbagtoltol/omniwriter/internal/types"
)

func Execute(ctx context.Context, req types.TransformRequest) (*types.TransformResult, error) {
	if req.Input == nil || len(req.Mapping) == 0 {
		return nil, errs.ErrInvalidRequest
	}

	ws, err := schema.ParseWriterSettings(req.Mapping)
	if err != nil {
		return nil, err
	}

	target := req.TargetFormat
	if target == "" {
		target = types.Format(ws.FileFormatType)
	}
	if string(target) != ws.FileFormatType {
		return nil, errs.ErrTargetFormatMismatch
	}

	if req.SourceFormat == target &&
		(target == types.FormatCSV || target == types.FormatJSON || target == types.FormatEDI || target == types.FormatXML) {
		b, err := io.ReadAll(req.Input)
		if err != nil {
			return nil, err
		}
		return &types.TransformResult{Output: b}, nil
	}

	if err := schema.Validate(req.Mapping); err != nil {
		return nil, err
	}

	od, err := schema.ParseOutputDeclaration(req.Mapping)
	if err != nil {
		return nil, err
	}

	switch target {
	case types.FormatEDI:
		return transformToEDI(ctx, req, od)
	case types.FormatJSON:
		return transformToJSON(ctx, req)
	case types.FormatCSV:
		return transformToCSV(ctx, req, od)
	case types.FormatXML:
		return transformToXML(ctx, req, od)
	case types.FormatText:
		return transformToText(ctx, req, od)
	default:
		return nil, fmt.Errorf("%w: %s->%s", errs.ErrUnsupportedTransform, req.SourceFormat, target)
	}
}

func ExecuteToWriter(ctx context.Context, req types.TransformRequest, out io.Writer) error {
	result, err := Execute(ctx, req)
	if err != nil {
		return err
	}
	_, err = out.Write(result.Output)
	return err
}

func transformToEDI(ctx context.Context, req types.TransformRequest, od *schema.OutputDeclaration) (*types.TransformResult, error) {
	if od == nil {
		return nil, errs.ErrMissingWriterEDIDecl
	}

	schemaWithoutWriter, err := schema.StripWriterFields(req.Mapping)
	if err != nil {
		return nil, err
	}

	omniSchema, err := omniparser.NewSchema("omniwriter-schema", bytes.NewReader(schemaWithoutWriter))
	if err != nil {
		return nil, err
	}

	tfm, err := omniSchema.NewTransform("input", req.Input, &transformctx.Ctx{})
	if err != nil {
		return nil, err
	}

	config := ediemit.Config{
		SegmentDelimiter:    od.SegmentDelimiter,
		ElementDelimiter:    od.ElementDelimiter,
		ComponentDelimiter:  od.ComponentDelimiter,
		RepetitionDelimiter: od.RepetitionDelimiter,
		IgnoreCRLF:          od.IgnoreCRLF,
	}

	var out strings.Builder
	records := 0

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		rec, err := tfm.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		var payload map[string]interface{}
		if err := json.Unmarshal(rec, &payload); err != nil {
			return nil, err
		}

		if err := ediemit.WriteRecord(&out, config, payload); err != nil {
			return nil, err
		}

		records++
	}

	return &types.TransformResult{
		Output: []byte(out.String()),
		Stats:  types.Stats{Records: records},
	}, nil
}

func transformToCSV(ctx context.Context, req types.TransformRequest, od *schema.OutputDeclaration) (*types.TransformResult, error) {
	if od == nil || len(od.Columns) == 0 {
		return nil, errs.ErrMissingWriterCSVDecl
	}

	schemaWithoutWriter, err := schema.StripWriterFields(req.Mapping)
	if err != nil {
		return nil, err
	}

	omniSchema, err := omniparser.NewSchema("omniwriter-schema", bytes.NewReader(schemaWithoutWriter))
	if err != nil {
		return nil, err
	}

	tfm, err := omniSchema.NewTransform("input", req.Input, &transformctx.Ctx{})
	if err != nil {
		return nil, err
	}

	var out bytes.Buffer
	delimiter := ','
	if od.Delimiter != "" {
		delimiter = []rune(od.Delimiter)[0]
	}

	writer := csvemit.NewWriter(&out, csvemit.Config{
		Delimiter: delimiter,
		Columns:   od.Columns,
	})

	records := 0

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		rec, err := tfm.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		var payload map[string]interface{}
		if err := json.Unmarshal(rec, &payload); err != nil {
			return nil, err
		}

		if err := writer.WriteRecord(payload); err != nil {
			return nil, err
		}

		records++
	}

	if err := writer.Flush(); err != nil {
		return nil, err
	}

	return &types.TransformResult{
		Output: out.Bytes(),
		Stats:  types.Stats{Records: records},
	}, nil
}

func transformToJSON(ctx context.Context, req types.TransformRequest) (*types.TransformResult, error) {
	schemaWithoutWriter, err := schema.StripWriterFields(req.Mapping)
	if err != nil {
		return nil, err
	}

	omniSchema, err := omniparser.NewSchema("omniwriter-schema", bytes.NewReader(schemaWithoutWriter))
	if err != nil {
		return nil, err
	}

	tfm, err := omniSchema.NewTransform("input", req.Input, &transformctx.Ctx{})
	if err != nil {
		return nil, err
	}

	var out bytes.Buffer
	writer := jsonemit.NewWriter(&out, jsonemit.DefaultConfig())
	records := 0

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		rec, err := tfm.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if err := writer.WriteRecord(rec); err != nil {
			return nil, err
		}

		records++
	}

	return &types.TransformResult{
		Output: out.Bytes(),
		Stats:  types.Stats{Records: records},
	}, nil
}

func transformToXML(ctx context.Context, req types.TransformRequest, od *schema.OutputDeclaration) (*types.TransformResult, error) {
	schemaWithoutWriter, err := schema.StripWriterFields(req.Mapping)
	if err != nil {
		return nil, err
	}

	omniSchema, err := omniparser.NewSchema("omniwriter-schema", bytes.NewReader(schemaWithoutWriter))
	if err != nil {
		return nil, err
	}

	tfm, err := omniSchema.NewTransform("input", req.Input, &transformctx.Ctx{})
	if err != nil {
		return nil, err
	}

	var out bytes.Buffer
	config := xmlemit.DefaultConfig()
	writer := xmlemit.NewWriter(&out, config)
	records := 0

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		rec, err := tfm.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		var payload map[string]interface{}
		if err := json.Unmarshal(rec, &payload); err != nil {
			return nil, err
		}

		if err := writer.WriteRecord(payload); err != nil {
			return nil, err
		}

		records++
	}

	if err := writer.Flush(); err != nil {
		return nil, err
	}

	return &types.TransformResult{
		Output: out.Bytes(),
		Stats:  types.Stats{Records: records},
	}, nil
}

func transformToText(ctx context.Context, req types.TransformRequest, od *schema.OutputDeclaration) (*types.TransformResult, error) {
	schemaWithoutWriter, err := schema.StripWriterFields(req.Mapping)
	if err != nil {
		return nil, err
	}

	omniSchema, err := omniparser.NewSchema("omniwriter-schema", bytes.NewReader(schemaWithoutWriter))
	if err != nil {
		return nil, err
	}

	tfm, err := omniSchema.NewTransform("input", req.Input, &transformctx.Ctx{})
	if err != nil {
		return nil, err
	}

	var out bytes.Buffer
	config := textemit.DefaultConfig()
	writer := textemit.NewWriter(&out, config)
	records := 0

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		rec, err := tfm.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		var payload map[string]interface{}
		if err := json.Unmarshal(rec, &payload); err != nil {
			return nil, err
		}

		if err := writer.WriteRecord(payload); err != nil {
			return nil, err
		}

		records++
	}

	if err := writer.Flush(); err != nil {
		return nil, err
	}

	return &types.TransformResult{
		Output: out.Bytes(),
		Stats:  types.Stats{Records: records},
	}, nil
}
