package gitlab

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"

	"github.com/ekristen/pipeliner/pkg/utils"
	"gitlab.com/gitlab-org/gitlab-runner/helpers/gitlab_ci_yaml_parser"
	"gopkg.in/yaml.v2"
)

func parse(data []byte) error {
	config := make(gitlab_ci_yaml_parser.DataBag)
	if err := yaml.Unmarshal(data, config); err != nil {
		return err
	}

	if err := config.Sanitize(); err != nil {
		return err
	}

	stages, ok := config.GetStringSlice("stages")
	if ok && len(stages) == 0 {
		stages = DefaultStages
	}

	/*
		defaults, ok := config.GetSubOptions("default")
		if ok {
		} else {
			fmt.Println("defaults do not exist")
			// TODO: assemble defaults from non-default section
		}
	*/

	pipeline := &Pipeline{}

	for k := range config {
		if IsReserved(k) {
			fmt.Println("reserved: ", k)
			continue
		}

		if IsHidden(k) {
			fmt.Println("hidden: ", k)
			continue
		}

		usedStages := []string{}
		for k := range config {
			d, _ := config.GetSubOptions(k)
			stage, ok := d.GetString("stage")
			if !ok {
				stage = "build"
			}
			usedStages = append(usedStages, stage)
		}

		stages = []string{}
		for _, ds := range stages {
			if utils.StringSliceContains(usedStages, ds) {
				stages = append(stages, ds)
			}
		}

		// Set resolved states on the pipeline definition
		pipeline.Stages = stages

		// Process Job information now
		jobData, _ := config.GetSubOptions(k)

		stage, ok := jobData.GetString("stage")
		if !ok {
			stage = "build"
		}

		dependencies, _ := jobData.GetStringSlice("dependencies")
		tags, _ := jobData.GetStringSlice("tags")
		//services, _ := jobData.GetStringSlice("services")
		// TODO: get timeout value!

		allowFailure := false
		when := "on_success"
		if v, ok := jobData["allow_failure"]; ok {
			allowFailure = v.(bool)
		}
		if v, ok := jobData["when"]; ok {
			when = v.(string)
		}

		matrixJobs := generateJobs(k, jobData)

		var jobs []Job

		for i, b := range matrixJobs {

			var jobVars []Variable
			for k, v := range b.Env {
				jobVars = append(jobVars, Variable{
					Name:   k,
					Value:  v,
					Masked: false,
					File:   false,
				})
			}

			jobs = append(jobs, Job{
				Name:         b.Name,
				Stage:        stage,
				State:        "created",
				Token:        utils.RandomString(16),
				When:         when,
				Variables:    jobVars,
				Dependencies: dependencies,

				Data: data,

				Default: Default{
					Tags: tags,
				},

				Parallel: Parallel{
					Count: i,
				},
				AllowFailure: AllowFailure{
					Enabled: allowFailure,
				},
			})
		}

		fmt.Println(jobs)
	}

	/*

		// TODO: this seems sloppy and prone to not be consistent,
		// look for better way to solve

		for _, createStage := range stages {
			for k := range config {

				defaultTags, defaultTagsOk := config.GetStringSlice("tags")

				d, _ := config.GetSubOptions(k)

				stage, ok := d.GetString("stage")
				if !ok {
					stage = "build"
				}

				if stage != createStage {
					continue
				}

				dependencies, _ := d.GetStringSlice("dependencies")

				tags, ok := d.GetStringSlice("tags")
				if !ok && defaultTagsOk {
					tags = defaultTags
				}

				// TODO: get timeout value!

				//services, _ := d.GetStringSlice("services")
				allowFailure := false
				when := "on_success"
				if v, ok := d["allow_failure"]; ok {
					allowFailure = v.(bool)
				}
				if v, ok := d["when"]; ok {
					when = v.(string)
				}

				c.log.WithField("stages", stages).Debug("stages")

				var buildTags []models.BuildTag
				for _, tag := range tags {
					buildTags = append(buildTags, models.BuildTag{
						Tag: tag,
					})
				}

				matrixBuilds := generateBuilds(k, d)

				var builds []*models.Build

				for i, b := range matrixBuilds {

					var buildVariables []*models.BuildVariable
					for k, v := range b.Env {
						buildVariables = append(buildVariables, &models.BuildVariable{
							Variable: &models.Variable{
								Name:   k,
								Value:  v,
								Masked: false,
								File:   false,
							},
						})
					}

					builds = append(builds, &models.Build{
						State:        "created",
						AllowFailure: allowFailure,
						Name:         b.Name,
						Job:          b.Job,
						Stage:        stage,
						StageIdx:     int64(utils.StringSlicePosition(stages, stage)),
						Token:        utils.RandomString(16),
						When:         when,
						Data:         pipeline.Data,
						Dependencies: strings.Join(dependencies, ","),
						Parallel:     int64(i),

						PipelineID: pipeline.ID,

						Tags:      buildTags,
						Variables: buildVariables,
					})

				}

				if err := c.db.Transaction(func(tx *gorm.DB) error {
					if err := tx.Create(builds).Error; err != nil {
						return err
					}

					if err := tx.Model(&pipeline).Association("Stages").Append(&models.PipelineStage{
						Index: int64(utils.StringSlicePosition(stages, stage)),
						Name:  stage,
						State: "created",
					}); err != nil {
						return err
					}

					return nil
				}); err != nil {
					c.log.WithError(err).Error("error with transaction")
					return err
				}
			}
		}
	*/

	return nil
}

func generateJobs(name string, yaml gitlab_ci_yaml_parser.DataBag) []MatrixBuild {
	buildsToCreate := []MatrixBuild{}

	if parallel, ok := yaml.GetSubOptions("parallel"); ok {
		if matrix, ok := parallel.GetSlice("matrix"); ok {
			matrixVariables := formatMatrix(matrix)

			// TODO: redo
			// Note: not happy with this to generate the matrix builds :/
			for _, v := range matrixVariables {
				varNames := []string{}
				varValues := [][]string{}
				for n, vv := range v {
					varNames = append(varNames, n)
					varValues = append(varValues, vv)
				}

				possible := utils.PermuteStrings(varValues...)
				//fmt.Println(possible)
				var buildEnvVars map[string]string = make(map[string]string)
				for j, v := range possible {
					buildEnvVars = make(map[string]string)
					for i, x := range v {
						buildEnvVars[varNames[i]] = x
					}

					buildEnvVars["CI_NODE_INDEX"] = fmt.Sprintf("%d", j)
					buildEnvVars["CI_NODE_TOTAL"] = fmt.Sprintf("%d", len(possible))

					buildsToCreate = append(buildsToCreate, MatrixBuild{
						Name: generateMatrixName(name, v),
						Job:  name,
						Env:  buildEnvVars,
					})
				}
			}
		}
	} else if parallel, ok := yaml.Get("parallel"); ok {
		parallelInt := parallel.(int)
		if parallelInt > 1 {
			for i := 0; i < parallelInt; i++ {
				buildsToCreate = append(buildsToCreate, MatrixBuild{
					Job:  name,
					Name: fmt.Sprintf("%s %d/%d", name, i+1, parallelInt),
					Env: map[string]string{
						"CI_NODE_INDEX": fmt.Sprintf("%d", i),
						"CI_NODE_TOTAL": fmt.Sprintf("%d", parallelInt),
					},
				})
			}
		}
	} else {
		buildsToCreate = append(buildsToCreate, MatrixBuild{
			Job:  name,
			Name: name,
		})
	}

	return buildsToCreate
}

func generateMatrixName(name string, env []string) string {
	return fmt.Sprintf("%s: [%s]", name, strings.Join(env, ", "))
}

func formatMatrix(matrix []interface{}) []map[string][]string {
	vals := make([]map[string][]string, len(matrix))
	for i, x := range matrix {
		//fmt.Println("----------------------------------------------")
		//fmt.Println("i:", i)
		j := 0
		for a, y := range x.(map[interface{}]interface{}) {
			//fmt.Println(a)
			if vals[i] == nil {
				vals[i] = make(map[string][]string)
			}

			if vals[i][a.(string)] == nil {
				vals[i][a.(string)] = make([]string, 0)
			}

			if reflect.TypeOf(y).String() == "string" {
				//fmt.Println(y)
				//fmt.Println(j, "-", 0)
				vals[i][a.(string)] = append(vals[i][a.(string)], y.(string))
			} else {
				k := 0
				for _, z := range y.([]interface{}) {
					//fmt.Println(z)
					//fmt.Println(j, "-", k)
					vals[i][a.(string)] = append(vals[i][a.(string)], z.(string))
					k++
				}
			}
			j++
		}
	}

	return vals
}

func TestOne(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/docs-default-example.yaml")
	if err != nil {
		t.Error(err)
	}

	parse(data)

	validator := NewYAMLValidator(data)
	validator.Parse()
	if err := validator.Validate(); err != nil {
		t.Error(err)
	}
}
