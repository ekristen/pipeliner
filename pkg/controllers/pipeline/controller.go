package pipeline

import (
	"context"

	pipelinerv1 "github.com/ekristen/pipeliner/pkg/apis/pipeliner.ekristen.dev/v1"
	pipelinerc "github.com/ekristen/pipeliner/pkg/generated/controllers/pipeliner.ekristen.dev/v1"
	"github.com/rancher/wrangler/pkg/apply"
)

type controller struct {
	ctx           context.Context
	apply         apply.Apply
	pipeline      pipelinerc.PipelineClient
	pipelineCache pipelinerc.PipelineCache
}

func Register(ctx context.Context, apply apply.Apply, pipeline pipelinerc.PipelineController) error {
	c := &controller{
		ctx:           ctx,
		apply:         apply,
		pipeline:      pipeline,
		pipelineCache: pipeline.Cache(),
	}

	pipeline.OnChange(ctx, "pipeline", c.OnChange)

	return nil
}

func (c *controller) OnChange(key string, pipeline *pipelinerv1.Pipeline) (*pipelinerv1.Pipeline, error) {
	if pipeline == nil {
		return nil, nil
	}

	return pipeline, nil
}
