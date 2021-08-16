package gitlab

import (
	"fmt"
	"strings"

	"github.com/ekristen/pipeliner/pkg/common"
	"github.com/ekristen/pipeliner/pkg/utils"
	"gitlab.com/gitlab-org/gitlab-runner/helpers/gitlab_ci_yaml_parser"
	"gopkg.in/yaml.v2"
)

// YAMLLinter --
type YAMLLinter struct {
	yaml   []byte
	config gitlab_ci_yaml_parser.DataBag
}

// NewYAMLValidator --
func NewYAMLValidator(yaml []byte) *YAMLLinter {
	return &YAMLLinter{
		yaml: yaml,
	}
}

// TODO: switch to yaml.v3 for anchor detection

// Parse --
func (c *YAMLLinter) Parse() (err error) {
	config := make(gitlab_ci_yaml_parser.DataBag)
	err = yaml.Unmarshal(c.yaml, config)
	if err != nil {
		return err
	}

	err = config.Sanitize()
	if err != nil {
		return err
	}

	c.config = config

	return nil
}

// Validate --
func (c *YAMLLinter) Validate() error {
	if err := c.Parse(); err != nil {
		return err
	}

	if err := c.ValidateJobsStages(); err != nil {
		return err
	}

	if err := c.ValidateJobsRequirements(); err != nil {
		return err
	}

	if err := c.ValidateDependencies(); err != nil {
		return err
	}

	return nil
}

// ValidateJobsStages --
func (c *YAMLLinter) ValidateJobsStages() error {
	stages, ok := c.config.GetStringSlice("stages")
	if !ok {
		stages = common.GitLabDefaultStages
	}

	for k := range c.getJobConfigs() {
		opts, _ := c.config.GetSubOptions(k)
		stage, _ := opts.GetString("stage")

		if !utils.StringSliceContains(stages, stage) {
			return fmt.Errorf("%s job: chosen stage (%s) does not exist; available stages are %s", k, stage, strings.Join(stages, ","))
		}
	}

	return nil
}

// ValidateJobsRequirements --
func (c *YAMLLinter) ValidateJobsRequirements() error {
	for k := range c.getJobConfigs() {
		opts, _ := c.config.GetSubOptions(k)
		if _, ok := opts.GetString("script"); !ok {
			if _, ok := opts.GetStringSlice("script"); !ok {
				return fmt.Errorf("%s job: missing script directive or is wrong format", k)
			}
		}
	}

	return nil
}

// ValidateDependencies --
func (c *YAMLLinter) ValidateDependencies() error {
	jobs := []string{}
	for k := range c.getJobConfigs() {
		jobs = append(jobs, k)
	}

	for k := range c.getJobConfigs() {
		opts, _ := c.config.GetSubOptions(k)

		stage, ok := opts.GetString("stage")
		if !ok {
			stage = "test"
		}

		if _, ok := opts.Get("dependencies"); ok {
			if deps, ok := opts.GetStringSlice("dependencies"); ok {
				for _, dep := range deps {
					if !utils.StringSliceContains(jobs, dep) {
						return fmt.Errorf("%s job: dependency %s not found", k, dep)
					}

					if !c.isDepInPreviousStage(stage, dep) {
						return fmt.Errorf("%s job: dependency %s not in previous stage", k, dep)
					}
				}
			}
		}
	}

	return nil
}

func (c *YAMLLinter) isDepInPreviousStage(stage, dep string) bool {
	stages, ok := c.config.GetStringSlice("stages")
	if !ok {
		stages = common.GitLabDefaultStages
	}

	depStageIndex := 0
	for k := range c.getJobConfigs() {
		if k == dep {
			opts, _ := c.config.GetSubOptions(k)

			stage, ok := opts.GetString("stage")
			if !ok {
				stage = "test"
			}

			depStageIndex = utils.StringSlicePosition(stages, stage)
		}
	}

	currentIndex := utils.StringSlicePosition(stages, stage)

	return currentIndex > depStageIndex
}

func (c *YAMLLinter) getJobConfigs() gitlab_ci_yaml_parser.DataBag {
	jobs := gitlab_ci_yaml_parser.DataBag{}
	for k, v := range c.config {
		if strings.HasPrefix(k, ".") {
			continue
		}
		if utils.StringSliceContains(common.GitLabUnavailableJobNames, k) {
			continue
		}

		jobs[k] = v
	}

	return jobs
}
