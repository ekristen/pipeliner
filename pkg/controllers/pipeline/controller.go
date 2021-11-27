package pipeline

import (
	"context"
	"fmt"

	pipelinerv1 "github.com/ekristen/pipeliner/pkg/apis/pipeliner.ekristen.dev/v1"
	pipelinerc "github.com/ekristen/pipeliner/pkg/generated/controllers/pipeliner.ekristen.dev/v1"
	"github.com/rancher/wrangler/pkg/apply"
)

type controller struct {
	ctx           context.Context
	apply         apply.Apply
	pipeline      pipelinerc.PipelineController
	pipelineCache pipelinerc.PipelineCache
}

func Register(ctx context.Context, apply apply.Apply, pipeline pipelinerc.PipelineController) error {
	c := &controller{
		ctx:           ctx,
		apply:         apply,
		pipeline:      pipeline,
		pipelineCache: pipeline.Cache(),
	}

	// pipeline.OnChange(ctx, "pipeline", c.OnChange)

	pipelinerc.RegisterPipelineStatusHandler(ctx, c.pipeline, "state", "state-changes", c.onStateChanges)

	return nil
}

func (c *controller) OnChange(key string, pipeline *pipelinerv1.Pipeline) (*pipelinerv1.Pipeline, error) {
	if pipeline == nil {
		return nil, nil
	}

	return pipeline, nil
}

func (c *controller) onStateChanges(pipeline *pipelinerv1.Pipeline, status pipelinerv1.PipelineStatus) (pipelinerv1.PipelineStatus, error) {
	// TODO: restore state from persisted storage

	switch status.State {
	case pipelinerv1.EMPTY:
		status.State = pipelinerv1.INITIALIZING

		// TODO: parse and create jobs frome pipeline
		// also build dag graph of job execution order
		// from dependencies or needs
	case pipelinerv1.INITIALIZING:
		c.doInitialize(pipeline, status)

	}

	return status, nil
}

func (c *controller) doInitialize(pipeline *pipelinerv1.Pipeline, status pipelinerv1.PipelineStatus) (pipelinerv1.PipelineStatus, error) {
	if pipeline.Spec.SourceRef.Kind != "Workflow" {
		return status, fmt.Errorf("kind: %s not yet implemented", pipeline.Spec.SourceRef.Kind)
	}

	return status, nil
}
