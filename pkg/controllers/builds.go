package controllers

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v1"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	runnerCommon "gitlab.com/gitlab-org/gitlab-runner/common"
	"gitlab.com/gitlab-org/gitlab-runner/helpers/gitlab_ci_yaml_parser"

	"github.com/ekristen/pipeliner/pkg/common"
	"github.com/ekristen/pipeliner/pkg/models"
	"github.com/ekristen/pipeliner/pkg/models/scopes"
	"github.com/ekristen/pipeliner/pkg/utils"
)

// Builds --
type Builds struct {
	db             *gorm.DB
	notifyPipeline chan *models.Pipeline
	pipelines      *Pipelines
	log            *logrus.Entry
}

// NewBuilds --
func NewBuilds(db *gorm.DB, pc *Pipelines, notifyPipeline chan *models.Pipeline) *Builds {
	return &Builds{
		db:             db,
		notifyPipeline: notifyPipeline,
		pipelines:      pc,
		log:            logrus.WithField("component", "worker").WithField("worker", "builds"),
	}
}

// CreateFromPipeline --
func (c *Builds) CreateFromPipeline(pipeline *models.Pipeline) error {
	config := make(gitlab_ci_yaml_parser.DataBag)
	if err := yaml.Unmarshal(pipeline.Data, config); err != nil {
		return err
	}

	if err := config.Sanitize(); err != nil {
		return err
	}

	// TODO: write a help function to reduce the stages lines between 52-72
	defaultStages, _ := config.GetStringSlice("stages")
	if len(defaultStages) == 0 {
		defaultStages = common.GitLabDefaultStages
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

	stages := []string{}
	for _, ds := range defaultStages {
		if utils.StringSliceContains(usedStages, ds) {
			stages = append(stages, ds)
		}
	}

	// TODO: this seems sloppy and prone to not be consistent,
	// look for better way to solve

	for _, createStage := range stages {
		for k := range config {
			// exclude hidden jobs
			if strings.HasPrefix(k, ".") {
				continue
			}

			// exclude reserved keywords
			if utils.StringSliceContains(common.GitLabUnavailableJobNames, k) {
				continue
			}

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

	return nil
}

// ForRunner --
func (c *Builds) ForRunner(runner *models.Runner) (*models.Build, error) {
	var build models.Build

	var runnerTags []string
	//var buildTags []string

	if len(runner.Tags) > 0 {
		for _, tag := range runner.Tags {
			runnerTags = append(runnerTags, tag.Tag)
		}
	}

	/*
		SELECT C.* FROM
		(SELECT test.id,COUNT(1) fullcount
		  FROM test
		  INNER JOIN test_tags as tags ON test.id = tags.test_id
		  GROUP BY test.id) as A
		INNER JOIN
		(SELECT test.id,COUNT(1) goodcount
		  FROM test
		  INNER JOIN test_tags as tags ON test.id = tags.test_id
		  WHERE tags.name IN ('aws','linux','k8s','one','two') GROUP BY test.id) as B USING (id)
		INNER JOIN (SELECT DISTINCT id,name FROM test) C USING (id)
		WHERE fullcount=goodcount;
	*/

	for i := 0; i < 3; i++ {
		expectedSQL := c.db.
			Select("b1.id, count(1) AS expectedCount").
			Table("builds as b1").
			Joins("INNER JOIN build_tags AS t1 ON b1.id = t1.build_id").
			Group("b1.id")

		actualSQL := c.db.
			Select("b2.id, count(1) as actualCount").
			Table("builds as b2").
			Joins("INNER JOIN build_tags AS t2 ON b2.id = t2.build_id").
			Where("t2.tag IN ?", runnerTags).
			Group("b2.id")

		buildsSQL := c.db.
			Model(&models.Build{}).
			Select("*").
			Scopes(scopes.BuildLatest)

		if err := c.db.Transaction(func(tx *gorm.DB) error {
			sql := tx.
				Clauses(clause.Locking{Strength: "UPDATE"}).
				Select("A.*").
				Table("(?) AS A", buildsSQL).
				Unscoped().
				Joins("LEFT JOIN (?) AS B USING(id)", actualSQL).
				Joins("LEFT JOIN (?) AS C USING(id)", expectedSQL).
				Where("A.state = ? AND A.runner_id IS NULL", "pending")

			if runner.RunUntagged {
				sql.Where("C.expectedCount=B.actualCount OR C.expectedCount IS NULL AND B.actualCount IS NULL")
			} else {
				sql.Where("C.expectedCount=B.actualCount")
			}

			sql = sql.
				Order("A.stage_idx, A.id").
				Limit(1).
				First(&build)
			if sql.Error != nil {
				return sql.Error
			}

			sql = tx.
				Model(&models.Build{}).
				Where("id = ? AND runner_id IS NULL AND state = 'pending' AND updated_at = ?", build.ID, build.UpdatedAt).
				Update("runner_id", runner.ID)

			if sql.Error != nil {
				logrus.WithError(sql.Error).Error("unable to update record")
				return sql.Error
			}

			if sql.RowsAffected == 0 {
				logrus.WithField("try", i).Warn("reprocess")
				return errors.New("no rows affected")
			}

			return nil
		}); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, nil
			}

			logrus.WithError(err).Error("transaction error")
			continue
		}

		sql := c.db.
			Model(&models.Build{}).
			Preload(clause.Associations).
			Where("id = ?", build.ID).
			First(&build)
		if sql.Error != nil {
			logrus.WithError(sql.Error).Error("uinable to retrive build data")
		}

		return &build, nil
	}

	return nil, fmt.Errorf("Builds available, but unable to get lock")
}

// StartBuild --
func (c *Builds) StartBuild(build *models.Build, runner *models.Runner) (*models.Build, error) {
	running := "running"

	build, err := c.TransitionState(build, &running, nil, true)
	if err != nil {
		return build, err
	}

	if err := c.pipelines.Run(&build.Pipeline); err != nil {
		return build, err
	}

	return build, nil
}

// Enqueue --
func (c *Builds) Enqueue(build *models.Build, lastStageState string, hadPreviousFailure bool) {
	c.log.WithField("build_id", build.ID).WithField("current_state", build.State).Debug("Enqueue")

	if build.State == "canceled" {
		c.DoUpdateState(build, "canceled", nil)
	} else if build.State != "skipped" && lastStageState == "failed" {
		c.DoUpdateState(build, "skipped", nil)
	} else if lastStageState == "success" && build.When == "on_failure" && !hadPreviousFailure {
		c.DoUpdateState(build, "skipped", nil)
	} else if lastStageState == "skipped" && utils.StringSliceContains([]string{"always", "on_failure"}, build.When) && hadPreviousFailure {
		c.DoUpdateState(build, "pending", nil)
	} else if lastStageState == "skipped" || lastStageState == "canceled" {
		c.DoUpdateState(build, "skipped", nil)
	} else if build.State == "created" && build.When == "manual" {
		c.DoUpdateState(build, "manual", nil)
	} else if build.State == "created" {
		c.DoUpdateState(build, "pending", nil)
	}

	/*
		ORIGINAL
			if build.State == "created" && build.When == "manual" {
				c.DoUpdateState(build, "manual")
			} else if build.State == "created" {
				c.DoUpdateState(build, "pending")
			}
	*/

	/**
			} else if (build.State == "created" || build.State == "pending") && build.When == "manual" {
			c.DoUpdateState(build, "pending")
		} else if (build.State == "created" || build.State == "pending") && build.When != "manual" {
			c.DoUpdateState(build, "pending")
	**/
}

// TransitionState --
func (c *Builds) TransitionState(build *models.Build, state *string, failureReason *string, pipelineNotify bool) (*models.Build, error) {
	var err error

	c.log.WithField("current", build.State).WithField("new", *state).Info("build: transition state")

	if build.State == "canceled" {
		return build, nil
	}

	running := "running"
	success := "success"

	if build.State == "created" || build.State == "pending" {
		if state == nil {
			state = &running
		}

		build, err = c.UpdateState(build, *state, nil)
	} else if build.State == "manual" || build.State == "running" {
		if state == nil {
			state = &success
		}
		build, err = c.UpdateState(build, *state, failureReason)
	}

	if pipelineNotify {
		c.notifyPipeline <- &build.Pipeline
	}

	return build, err
}

// UpdateState --
func (c *Builds) UpdateState(build *models.Build, state string, failureReason *string) (*models.Build, error) {
	c.log.WithField("id", build.ID).WithField("current", build.State).WithField("new", state).Debugf("Called: BuildsUpdateState")

	build, err := c.DoUpdateState(build, state, failureReason)
	if err != nil {
		return nil, err
	}

	return build, nil
}

// DoUpdateState --
func (c *Builds) DoUpdateState(build *models.Build, state string, failureReason *string) (*models.Build, error) {
	c.log.Debugf("Called: BuildsDoUpdateState (state: %s)\n", state)

	var updates map[string]interface{}

	if state == "success" || state == "failed" || state == "canceled" {
		if build.StartedAt != nil {
			finishedAt := time.Now().UTC()
			duration := finishedAt.Sub(*build.StartedAt)

			updates = map[string]interface{}{"state": state, "finished_at": finishedAt, "duration": int(duration.Seconds())}
		} else {
			updates = map[string]interface{}{"state": state}
		}

		if state == "failed" && failureReason != nil {
			updates["failure_reason"] = *failureReason
		}
	} else if state == "pending" {
		updates = map[string]interface{}{"state": state, "queued_at": time.Now().UTC()}
	} else if state == "running" {
		updates = map[string]interface{}{"state": state, "started_at": time.Now().UTC()}
	} else {
		updates = map[string]interface{}{"state": state}
	}

	sql := c.db.Model(build).
		Updates(updates)
	if sql.Error != nil {
		c.log.WithError(sql.Error).Error("unable to update build")
		return build, sql.Error
	}

	sql = c.db.Model(&models.Build{}).Where("id = ?", build.ID).First(&build)
	if sql.Error != nil {
		c.log.WithError(sql.Error).Error("unable to get updated build data")
		return build, sql.Error
	}

	c.log.Debug("successfully performed state update")

	return build, nil
}

// BuildDependencies --
func (c *Builds) BuildDependencies(build *models.Build) runnerCommon.Dependencies {
	deps := runnerCommon.Dependencies{}
	for _, name := range strings.Split(build.Dependencies, ",") {
		dd := c.NamedDependencies(build, name)
		deps = append(deps, dd...)
	}

	return deps
}

// NamedDependencies --
func (c *Builds) NamedDependencies(build *models.Build, name string) runnerCommon.Dependencies {
	var builds []models.Build

	// TODO: add error handling
	c.db.
		Model(&models.Build{}).
		Preload(clause.Associations).
		Where("pipeline_id = ? and stage_idx < ? and job = ?", build.PipelineID, build.StageIdx, name).
		Find(&builds)

	return c.MapDependency(builds)
}

// MapDependency --
func (c *Builds) MapDependency(builds []models.Build) runnerCommon.Dependencies {
	deps := []runnerCommon.Dependency{}

	for _, build := range builds {
		dep := runnerCommon.Dependency{
			ID:    int(build.ID),
			Token: build.Token,
			Name:  build.Name,
		}

		archive := getArtifactArchive(build.Artifacts)
		if archive != nil {
			dep.ArtifactsFile = runnerCommon.DependencyArtifactsFile{
				Filename: archive.File,
				Size:     archive.Size,
			}
		}
		deps = append(deps, dep)
	}

	return deps
}

// BuildsForPipelineAndStage --
func (c *Builds) BuildsForPipelineAndStage(pipelineID, stageIdx int64) []*models.Build {
	var builds []*models.Build

	// TODO: add error handling
	c.db.
		Model(&models.Build{}).
		Where("pipeline_id = ?", pipelineID).
		Where("stage_idx = ?", stageIdx).
		Order("created_at asc").
		Find(&builds)

	return builds
}

func getArtifactArchive(artifacts []*models.Artifact) *models.Artifact {
	for _, artifact := range artifacts {
		if artifact.Type == "archive" {
			return artifact
		}
	}

	return nil
}

// MatrixBuild --
type MatrixBuild struct {
	Name string
	Job  string
	Env  map[string]string
}

func generateBuilds(name string, yaml gitlab_ci_yaml_parser.DataBag) []MatrixBuild {
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
