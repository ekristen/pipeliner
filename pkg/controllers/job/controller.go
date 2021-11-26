package job

import (
	"context"

	pipelinerv1 "github.com/ekristen/pipeliner/pkg/apis/pipeliner.ekristen.dev/v1"
	pipelinerc "github.com/ekristen/pipeliner/pkg/generated/controllers/pipeliner.ekristen.dev/v1"
	"github.com/rancher/wrangler/pkg/apply"
)

type controller struct {
	ctx      context.Context
	apply    apply.Apply
	job      pipelinerc.JobClient
	jobCache pipelinerc.JobCache
}

func Register(ctx context.Context, apply apply.Apply, job pipelinerc.JobController) error {
	c := &controller{
		ctx:      ctx,
		apply:    apply,
		job:      job,
		jobCache: job.Cache(),
	}

	job.OnChange(ctx, "job", c.OnChange)

	return nil
}

func (c *controller) OnChange(key string, job *pipelinerv1.Job) (*pipelinerv1.Job, error) {
	if job == nil {
		return nil, nil
	}

	return job, nil
}
