package omniwriter

import (
	"context"
	"io"

	"github.com/pbagtoltol/omniwriter/internal/pipeline"
)

func Transform(ctx context.Context, req TransformRequest) (*TransformResult, error) {
	return pipeline.Execute(ctx, req)
}

func TransformToWriter(ctx context.Context, req TransformRequest, out io.Writer) error {
	return pipeline.ExecuteToWriter(ctx, req, out)
}
