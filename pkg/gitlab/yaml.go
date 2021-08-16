package gitlab

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"gitlab.com/gitlab-org/gitlab-runner/common"
	"gitlab.com/gitlab-org/gitlab-runner/helpers/gitlab_ci_yaml_parser"

	"gopkg.in/yaml.v2"

	common2 "github.com/ekristen/pipeliner/pkg/common"
)

// YAMLParser --
type YAMLParser struct {
	filename  string
	jobName   string
	yaml      []byte
	config    gitlab_ci_yaml_parser.DataBag
	jobConfig gitlab_ci_yaml_parser.DataBag
}

func (c *YAMLParser) parseFile() (err error) {
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

	return
}

func (c *YAMLParser) loadJob() (err error) {
	jobConfig, ok := c.config.GetSubOptions(c.jobName)
	if !ok {
		return fmt.Errorf("no job named %s", c.jobName)
	}

	c.jobConfig = jobConfig

	return
}

func (c *YAMLParser) prepareJobInfo(job *common.JobResponse) (err error) {
	//job.JobInfo = common.JobInfo{
	//	Name: c.jobName,
	//}

	if stage, ok := c.jobConfig.GetString("stage"); ok {
		job.JobInfo.Stage = stage
	} else {
		job.JobInfo.Stage = "test"
	}

	return
}

func (c *YAMLParser) getCommands(commands interface{}) (common.StepScript, error) {
	switch t := commands.(type) {
	case []interface{}:
		var steps common.StepScript
		for _, line := range t {
			if lineText, ok := line.(string); ok {
				steps = append(steps, lineText)
			} else {
				return common.StepScript{}, errors.New("unsupported script")
			}
		}
		return steps, nil
	case string:
		return strings.Split(t, "\n"), nil
	default:
		if commands != nil {
			return common.StepScript{}, errors.New("unsupported script")
		}
	}

	return common.StepScript{}, nil
}

func (c *YAMLParser) prepareSteps(job *common.JobResponse) error {
	if c.jobConfig["script"] == nil {
		return fmt.Errorf("missing 'script' for job")
	}

	var scriptCommands, afterScriptCommands common.StepScript

	// get before_script
	beforeScript, err := c.getCommands(c.config["before_script"])
	if err != nil {
		return err
	}

	// get job before_script
	jobBeforeScript, err := c.getCommands(c.jobConfig["before_script"])
	if err != nil {
		return err
	}

	if len(jobBeforeScript) < 1 {
		scriptCommands = beforeScript
	} else {
		scriptCommands = jobBeforeScript
	}

	// get script
	script, err := c.getCommands(c.jobConfig["script"])
	if err != nil {
		return err
	}
	scriptCommands = append(scriptCommands, script...)

	afterScriptCommands, err = c.getCommands(c.jobConfig["after_script"])
	if err != nil {
		return err
	}

	job.Steps = common.Steps{
		common.Step{
			Name:         common.StepNameScript,
			Script:       scriptCommands,
			Timeout:      c.GetTimeout(),
			When:         common.StepWhenOnSuccess,
			AllowFailure: false,
		},
		common.Step{
			Name:         common.StepNameAfterScript,
			Script:       afterScriptCommands,
			Timeout:      c.GetTimeout(),
			When:         common.StepWhenAlways,
			AllowFailure: false,
		},
	}

	return nil
}

// GetTimeout --
func (c *YAMLParser) GetTimeout() int {
	stimeout, ok := c.jobConfig.Get("timeout")
	if !ok {
		stimeout, ok = c.config.Get("timeout")
		if !ok {
			return 3600 // return default
		}
	}

	return stimeout.(int)
}

func (c *YAMLParser) buildDefaultVariables(job *common.JobResponse) common.JobVariables {
	return common.JobVariables{
		{Key: "CI", Value: "true", Public: true, Internal: true, File: false},
		{Key: "GITLAB_CI", Value: "true", Public: true, Internal: true, File: false},
		{Key: "CI_SERVER_NAME", Value: common2.AppVersion.Name, Public: true, Internal: true, File: false},
		{Key: "CI_SERVER_VERSION", Value: common2.AppVersion.Version, Public: true, Internal: true, File: false},
		{Key: "CI_SERVER_REVISION", Value: common2.AppVersion.Summary, Public: true, Internal: true, File: false},
		{Key: "CI_PROJECT_ID", Value: strconv.Itoa(job.JobInfo.ProjectID), Public: true, Internal: true, File: false},
		{Key: "CI_JOB_ID", Value: strconv.Itoa(job.ID), Public: true, Internal: true, File: false},
		{Key: "CI_JOB_NAME", Value: job.JobInfo.Name, Public: true, Internal: true, File: false},
		{Key: "CI_JOB_STAGE", Value: job.JobInfo.Stage, Public: true, Internal: true, File: false},
		{Key: "CI_JOB_TOKEN", Value: job.Token, Public: true, Internal: true, File: false},
		{Key: "CI_REPOSITORY_URL", Value: job.GitInfo.RepoURL, Public: true, Internal: true, File: false},
		{Key: "CI_COMMIT_SHA", Value: job.GitInfo.Sha, Public: true, Internal: true, File: false},
		{Key: "CI_COMMIT_BEFORE_SHA", Value: job.GitInfo.BeforeSha, Public: true, Internal: true, File: false},
		{Key: "CI_COMMIT_REF_NAME", Value: job.GitInfo.Ref, Public: true, Internal: true, File: false},
	}
}

func (c *YAMLParser) buildVariables(
	configVariables interface{},
) (buildVariables common.JobVariables, err error) {
	if variables, ok := configVariables.(map[string]interface{}); ok {
		for key, value := range variables {
			if valueText, ok := value.(string); ok {
				buildVariables = append(buildVariables, common.JobVariable{
					Key:    key,
					Value:  valueText,
					Public: true,
				})
			} else {
				err = fmt.Errorf("invalid value for variable %q", key)
			}
		}
	} else if configVariables != nil {
		err = errors.New("unsupported variables")
	}

	return
}

func (c *YAMLParser) prepareVariables(job *common.JobResponse) error {
	job.Variables = common.JobVariables{}

	defaultVariables := c.buildDefaultVariables(job)
	job.Variables = append(job.Variables, defaultVariables...)

	globalVariables, err := c.buildVariables(c.config["variables"])
	if err != nil {
		return err
	}

	job.Variables = append(job.Variables, globalVariables...)

	jobVariables, err := c.buildVariables(c.jobConfig["variables"])
	if err != nil {
		return err
	}

	job.Variables = append(job.Variables, jobVariables...)

	return nil
}

func (c *YAMLParser) prepareImage(job *common.JobResponse) error {
	job.Image = common.Image{}

	if imageName, ok := c.jobConfig.GetString("image"); ok {
		job.Image.Name = imageName
		return nil
	}

	if imageDefinition, ok := c.jobConfig.GetSubOptions("image"); ok {
		job.Image.Name, _ = imageDefinition.GetString("name")
		job.Image.Entrypoint, _ = imageDefinition.GetStringSlice("entrypoint")
		return nil
	}

	if imageName, ok := c.config.GetString("image"); ok {
		job.Image.Name = imageName
		return nil
	}

	if imageDefinition, ok := c.config.GetSubOptions("image"); ok {
		job.Image.Name, _ = imageDefinition.GetString("name")
		job.Image.Entrypoint, _ = imageDefinition.GetStringSlice("entrypoint")
		return nil
	}

	return nil
}

func parseExtendedServiceDefinitionMap(serviceDefinition map[interface{}]interface{}) (image common.Image) {
	service := make(gitlab_ci_yaml_parser.DataBag)
	for key, value := range serviceDefinition {
		service[key.(string)] = value
	}

	image.Name, _ = service.GetString("name")
	image.Alias, _ = service.GetString("alias")
	image.Command, _ = service.GetStringSlice("command")
	image.Entrypoint, _ = service.GetStringSlice("entrypoint")
	return
}

func (c *YAMLParser) prepareServices(job *common.JobResponse) (err error) {
	job.Services = common.Services{}

	if servicesMap, ok := getOptions("services", c.jobConfig, c.config); ok {
		for _, service := range servicesMap {
			if serviceName, ok := service.(string); ok {
				job.Services = append(job.Services, common.Image{
					Name: serviceName,
				})
				continue
			}

			if serviceDefinition, ok := service.(map[interface{}]interface{}); ok {
				job.Services = append(job.Services, parseExtendedServiceDefinitionMap(serviceDefinition))
			}
		}
	}

	return
}

func (c *YAMLParser) prepareArtifacts(job *common.JobResponse) error {
	var ok bool

	artifactsMap := getOptionsMap("artifacts", c.jobConfig, c.config)

	artifactsPaths, _ := artifactsMap.GetSlice("paths")
	paths := common.ArtifactPaths{}
	for _, path := range artifactsPaths {
		paths = append(paths, path.(string))
	}

	var artifactsName string
	if artifactsName, ok = artifactsMap.GetString("name"); !ok {
		artifactsName = ""
	}

	var artifactsUntracked interface{}
	if artifactsUntracked, ok = artifactsMap.Get("untracked"); !ok {
		artifactsUntracked = false
	}

	var artifactsWhen string
	if artifactsWhen, ok = artifactsMap.GetString("when"); !ok {
		artifactsWhen = string(common.ArtifactWhenOnSuccess)
	}

	var artifactsExpireIn string
	if artifactsExpireIn, ok = artifactsMap.GetString("expireIn"); !ok {
		artifactsExpireIn = ""
	}

	job.Artifacts = make(common.Artifacts, 0)
	job.Artifacts = append(job.Artifacts, common.Artifact{
		Name:      artifactsName,
		Untracked: artifactsUntracked.(bool),
		Paths:     paths,
		Type:      "archive", // default type
		Format:    "zip",
		When:      common.ArtifactWhen(artifactsWhen),
		ExpireIn:  artifactsExpireIn,
	})

	if reportsMap, ok := artifactsMap.GetSubOptions("reports"); ok {
		if dotenv, ok := reportsMap.GetString("dotenv"); ok {
			job.Artifacts = append(job.Artifacts, common.Artifact{
				Name:     "dotenv",
				Type:     "dotenv",
				Format:   common.ArtifactFormatGzip,
				Paths:    common.ArtifactPaths{dotenv},
				When:     common.ArtifactWhenAlways,
				ExpireIn: "",
			})
		}
	}

	return nil
}

func (c *YAMLParser) prepareCache(job *common.JobResponse) error {
	var ok bool

	cacheMap := getOptionsMap("cache", c.jobConfig, c.config)

	cachePaths, _ := cacheMap.GetSlice("paths")
	paths := common.ArtifactPaths{}
	for _, path := range cachePaths {
		paths = append(paths, path.(string))
	}

	var cacheKey string
	if cacheKey, ok = cacheMap.GetString("key"); !ok {
		cacheKey = ""
	}

	var cacheUntracked interface{}
	if cacheUntracked, ok = cacheMap.Get("untracked"); !ok {
		cacheUntracked = false
	}

	job.Cache = make(common.Caches, 1)
	job.Cache[0] = common.Cache{
		Key:       cacheKey,
		Untracked: cacheUntracked.(bool),
		Paths:     paths,
	}

	return nil
}

// ParseYaml --
func (c *YAMLParser) ParseYaml(job *common.JobResponse) (err error) {
	err = c.parseFile()
	if err != nil {
		return err
	}

	err = c.loadJob()
	if err != nil {
		return err
	}

	parsers := []struct {
		method func(job *common.JobResponse) error
	}{
		{c.prepareJobInfo},
		{c.prepareSteps},
		{c.prepareVariables},
		{c.prepareImage},
		{c.prepareServices},
		{c.prepareArtifacts},
		{c.prepareCache},
	}

	for _, parser := range parsers {
		err = parser.method(job)
		if err != nil {
			return err
		}
	}

	return nil
}

// NewYAMLParser --
func NewYAMLParser(jobName string, yaml []byte) *YAMLParser {
	return &YAMLParser{
		filename: ".gitlab-ci.yml",
		jobName:  jobName,
		yaml:     yaml,
	}
}

func getOptionsMap(optionKey string, primary, secondary gitlab_ci_yaml_parser.DataBag) (value gitlab_ci_yaml_parser.DataBag) {
	value, ok := primary.GetSubOptions(optionKey)
	if !ok {
		value, _ = secondary.GetSubOptions(optionKey)
	}

	return
}

func getOptions(optionKey string, primary, secondary gitlab_ci_yaml_parser.DataBag) (value []interface{}, ok bool) {
	value, ok = primary.GetSlice(optionKey)
	if !ok {
		value, ok = secondary.GetSlice(optionKey)
	}

	return
}
