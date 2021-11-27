package artifact

import (
	"context"

	pipelinerv1 "github.com/ekristen/pipeliner/pkg/apis/pipeliner.ekristen.dev/v1"
	pipelinerc "github.com/ekristen/pipeliner/pkg/generated/controllers/pipeliner.ekristen.dev/v1"
	"github.com/rancher/wrangler/pkg/apply"
)

type controller struct {
	ctx           context.Context
	apply         apply.Apply
	artifact      pipelinerc.ArtifactClient
	artifactCache pipelinerc.ArtifactCache
}

func Register(ctx context.Context, apply apply.Apply, artifact pipelinerc.ArtifactController) error {
	c := &controller{
		ctx:           ctx,
		apply:         apply,
		artifact:      artifact,
		artifactCache: artifact.Cache(),
	}

	artifact.OnChange(ctx, "artifact", c.OnChange)

	return nil
}

func (c *controller) OnChange(key string, artifact *pipelinerv1.Artifact) (*pipelinerv1.Artifact, error) {
	if artifact == nil {
		return nil, nil
	}

	return artifact, nil
}
