package gitrepository

import (
	"context"

	pipelinerv1 "github.com/ekristen/pipeliner/pkg/apis/pipeliner.ekristen.dev/v1"
	pipelinerc "github.com/ekristen/pipeliner/pkg/generated/controllers/pipeliner.ekristen.dev/v1"
	"github.com/rancher/wrangler/pkg/apply"
)

type controller struct {
	ctx                context.Context
	apply              apply.Apply
	gitrepository      pipelinerc.GitRepositoryClient
	gitrepositoryCache pipelinerc.GitRepositoryCache
}

func Register(ctx context.Context, apply apply.Apply, gitrepository pipelinerc.GitRepositoryController) error {
	c := &controller{
		ctx:                ctx,
		apply:              apply,
		gitrepository:      gitrepository,
		gitrepositoryCache: gitrepository.Cache(),
	}

	gitrepository.OnChange(ctx, "gitrepository", c.OnChange)

	return nil
}

func (c *controller) OnChange(key string, gitrepository *pipelinerv1.GitRepository) (*pipelinerv1.GitRepository, error) {
	if gitrepository == nil {
		return nil, nil
	}

	return gitrepository, nil
}
