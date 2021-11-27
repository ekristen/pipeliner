package workflow

import (
	"context"

	pipelinerv1 "github.com/ekristen/pipeliner/pkg/apis/pipeliner.ekristen.dev/v1"
	pipelinerc "github.com/ekristen/pipeliner/pkg/generated/controllers/pipeliner.ekristen.dev/v1"
	"github.com/rancher/wrangler/pkg/apply"
	"github.com/rancher/wrangler/pkg/condition"
	"github.com/sirupsen/logrus"

	"github.com/ekristen/pipeliner/pkg/gitlab"
)

var validCondition condition.Cond = "valid"

type controller struct {
	ctx           context.Context
	log           *logrus.Entry
	apply         apply.Apply
	workflow      pipelinerc.WorkflowController
	workflowCache pipelinerc.WorkflowCache
}

func Register(ctx context.Context, apply apply.Apply, workflow pipelinerc.WorkflowController) error {
	c := &controller{
		ctx:           ctx,
		log:           logrus.WithField("controller", "workflow"),
		apply:         apply,
		workflow:      workflow,
		workflowCache: workflow.Cache(),
	}

	// workflow.OnChange(ctx, "workflow", c.OnChange)

	pipelinerc.RegisterWorkflowStatusHandler(ctx, c.workflow, validCondition, "valid-status", c.ValidStatus)

	return nil
}

func (c *controller) OnChange(key string, workflow *pipelinerv1.Workflow) (*pipelinerv1.Workflow, error) {
	if workflow == nil {
		return nil, nil
	}

	return workflow, nil
}

func (c *controller) ValidStatus(workflow *pipelinerv1.Workflow, status pipelinerv1.WorkflowStatus) (pipelinerv1.WorkflowStatus, error) {
	log := c.log.WithField("workflow", workflow.GetName())

	v := gitlab.NewYAMLValidator([]byte(workflow.Spec.Raw))
	if err := v.Validate(); err != nil {
		log.WithError(err).Error("unable to validate yaml config")
		validCondition.SetError(&status, "invalid yaml", err)

		return status, err
	}

	validCondition.SetStatus(&status, "True")
	validCondition.Message(&status, "Workflow YAML is valid")

	return status, nil
}
