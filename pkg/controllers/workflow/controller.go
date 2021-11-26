package workflow

import (
	"context"

	pipelinerv1 "github.com/ekristen/pipeliner/pkg/apis/pipeliner.ekristen.dev/v1"
	pipelinerc "github.com/ekristen/pipeliner/pkg/generated/controllers/pipeliner.ekristen.dev/v1"
	"github.com/rancher/wrangler/pkg/apply"
)

type controller struct {
	ctx           context.Context
	apply         apply.Apply
	workflow      pipelinerc.WorkflowClient
	workflowCache pipelinerc.WorkflowCache
}

func Register(ctx context.Context, apply apply.Apply, workflow pipelinerc.WorkflowController) error {
	c := &controller{
		ctx:           ctx,
		apply:         apply,
		workflow:      workflow,
		workflowCache: workflow.Cache(),
	}

	workflow.OnChange(ctx, "workflow", c.OnChange)

	return nil
}

func (c *controller) OnChange(key string, workflow *pipelinerv1.Workflow) (*pipelinerv1.Workflow, error) {
	if workflow == nil {
		return nil, nil
	}

	return workflow, nil
}
