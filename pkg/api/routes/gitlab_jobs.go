package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ekristen/pipeliner/pkg/gitlab"
	"github.com/ekristen/pipeliner/pkg/models"
	"github.com/gorilla/mux"
	runnerCommon "gitlab.com/gitlab-org/gitlab-runner/common"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Receives runnerCommon.UpdateJobRequest
func (c *Routes) jobHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	buildID := vars["id"]

	log := c.log.WithField("build", buildID)

	var jobUpdate runnerCommon.UpdateJobRequest
	if err := json.NewDecoder(r.Body).Decode(&jobUpdate); err != nil {
		log.WithError(err).Error("unable to decode job request body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var build models.Build

	result := c.db.
		Model(&models.Build{}).
		Preload(clause.Associations).
		Where("id = ? and token = ?", buildID, jobUpdate.Token).
		First(&build)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.WithError(result.Error).Error("unable to get build from database")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sql := c.db.Model(&models.Runner{}).Where("id = ?", build.RunnerID).Updates(models.Runner{
		Version:      jobUpdate.Info.Version,
		Architecture: jobUpdate.Info.Architecture,
		Platform:     jobUpdate.Info.Platform,
	})
	if sql.Error != nil {
		log.WithError(result.Error).Error("unable to update runner metadata")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// TODO: update runner info from jobUpdate??
	// fmt.Println("RECEIVED:", jobUpdate.State)

	state := string(jobUpdate.State)
	failureReason := string(jobUpdate.FailureReason)

	log.WithField("state", state).WithField("failure_reason", failureReason).Debug("recevied job state")

	c.controllerBuilds.TransitionState(&build, &state, &failureReason, true)

	w.WriteHeader(http.StatusOK)
}

// Recevies runnerCommon.JobRequest
// See https://github.com/AlloyCI/alloy_ci/blob/6732662cb94e78290b784e90dd59adb4339c1c6a/lib/alloy_ci/web/controllers/api/builds_event_controller.ex
func (c *Routes) jobsRequestHandler(w http.ResponseWriter, r *http.Request) {
	var jobRequest runnerCommon.JobRequest
	if err := json.NewDecoder(r.Body).Decode(&jobRequest); err != nil {
		c.log.WithError(err).Error("unable to parse job request json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var runner models.Runner

	result := c.db.
		Model(&models.Runner{}).
		Preload(clause.Associations).
		Where("token = ? AND active = ? AND locked = ?", jobRequest.Token, true, false).
		First(&runner)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.log.WithField("token", jobRequest.Token).Warn("unable to find runner")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if result.Error != nil {
		c.log.WithError(result.Error).Error("unable to retrieve runner record from database")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log := c.log.WithField("runner_id", runner.ID)

	// Update runner contacted_at field
	if err := c.db.Model(&runner).Update("contacted_at", time.Now().UTC()).Error; err != nil {
		c.log.WithError(err).Error("unable to update runner info")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	build, err := c.controllerBuilds.ForRunner(&runner)
	if err != nil && err.Error() == "Builds available, but unable to get lock" {
		log.WithError(err).Warn("informing runner that builds are available but are in conflict")
		w.WriteHeader(http.StatusConflict)
		return
	}
	if err != nil {
		log.WithError(err).Error("unable to get build")
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	if build == nil {
		lastUpdate := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
		// log.WithField("last_update", lastUpdate).Debug("no builds found")
		w.Header().Set("X-GitLab-Last-Update", lastUpdate)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	build, err = c.controllerBuilds.StartBuild(build, &runner)
	if err != nil {
		log.WithError(err).Error("unable to get build")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.WithField("build_id", build.ID).Info("found a build")

	var workflow models.Workflow
	var pipeline models.Pipeline
	if err := c.db.Model(&models.Pipeline{}).Where("id = ?", build.PipelineID).First(&pipeline).Error; err != nil {
		log.WithError(err).Error("unable to get pipeline")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := c.db.Model(&models.Workflow{}).Where("id = ?", pipeline.WorkflowID).First(&workflow).Error; err != nil {
		log.WithError(err).Error("unable to get pipeline")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	job := runnerCommon.JobResponse{
		ID:    int(build.ID),
		Token: build.Token,
		JobInfo: runnerCommon.JobInfo{
			Name:        build.Name,
			Stage:       build.Stage,
			ProjectID:   int(workflow.ID),
			ProjectName: workflow.Name,
		},
		GitInfo: runnerCommon.GitInfo{
			RepoURL:   fmt.Sprintf("pipeliner://%d/build/%d", build.PipelineID, build.ID),
			Ref:       "master",
			Sha:       "000000000000000000000000000",
			BeforeSha: "000000000000000000000000000",
			RefType:   runnerCommon.GitInfoRefType("branch"),
		},
	}

	parser := gitlab.NewYAMLParser(build.Job, build.Data)
	if err := parser.ParseYaml(&job); err != nil {
		log.WithError(err).Error("unable to parse yaml")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// parser.ParseYaml stomps on the original job
	// just reset this information here for now
	job.JobInfo = runnerCommon.JobInfo{
		Name:        build.Name,
		Stage:       build.Stage,
		ProjectID:   int(workflow.ID),
		ProjectName: workflow.Name,
	}

	// This gets the timeout for the JOB!
	job.RunnerInfo = runnerCommon.RunnerInfo{
		Timeout: parser.GetTimeout(),
	}

	// Build Dependencies
	job.Dependencies = c.controllerBuilds.BuildDependencies(build)

	job.Variables = append(job.Variables, c.globalVariables()...)
	job.Variables = append(job.Variables, c.workflowVariables(workflow.ID)...)
	job.Variables = append(job.Variables, c.pipelineVariables(pipeline.ID)...)
	job.Variables = append(job.Variables, c.buildVariables(build.ID)...)

	// Loop through deps and append any variables generated
	// during the run
	for _, d := range job.Dependencies {
		job.Variables = append(job.Variables, c.buildVariables(int64(d.ID))...)
	}

	// Add Additional Variables
	job.Variables = append(job.Variables, runnerCommon.JobVariable{Key: "CI_TOKEN", Value: build.Token, Masked: true})

	// Add GIT_STRATEGY always to prevent git clone, we don't need it!
	job.Variables = append(job.Variables, runnerCommon.JobVariable{Key: "GIT_STRATEGY", Value: "none"})

	//b, _ := json.Marshal(job)
	//fmt.Println(string(b))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(job); err != nil {
		log.WithError(err).Error("unable to encode json")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (c *Routes) globalVariables() []runnerCommon.JobVariable {
	var variables []models.GlobalVariable
	c.db.Model(&models.GlobalVariable{}).Find(&variables)
	var jobVariables = []runnerCommon.JobVariable{}
	for _, variable := range variables {
		jobVariables = append(jobVariables, runnerCommon.JobVariable{
			Key:      variable.Name,
			Value:    variable.Value,
			Masked:   variable.Masked,
			File:     variable.File,
			Internal: variable.Internal,
			Public:   variable.Public,
		})
	}
	return jobVariables
}

func (c *Routes) workflowVariables(id int64) []runnerCommon.JobVariable {
	var variables []models.WorkflowVariable
	c.db.Model(&models.WorkflowVariable{}).Where("workflow_id = ?", id).Find(&variables)
	var jobVariables = []runnerCommon.JobVariable{}
	for _, variable := range variables {
		jobVariables = append(jobVariables, runnerCommon.JobVariable{
			Key:      variable.Name,
			Value:    variable.Value,
			Masked:   variable.Masked,
			File:     variable.File,
			Internal: variable.Internal,
			Public:   variable.Public,
		})
	}
	return jobVariables
}

func (c *Routes) pipelineVariables(id int64) []runnerCommon.JobVariable {
	var variables []models.PipelineVariable
	c.db.Model(&models.PipelineVariable{}).Where("pipeline_id = ?", id).Find(&variables)
	var jobVariables = []runnerCommon.JobVariable{}
	for _, variable := range variables {
		jobVariables = append(jobVariables, runnerCommon.JobVariable{
			Key:      variable.Name,
			Value:    variable.Value,
			Masked:   variable.Masked,
			File:     variable.File,
			Internal: variable.Internal,
			Public:   variable.Public,
		})
	}
	return jobVariables
}

func (c *Routes) buildVariables(id int64) []runnerCommon.JobVariable {
	var variables []models.BuildVariable
	c.db.Model(&models.BuildVariable{}).Where("build_id = ?", id).Find(&variables)
	var jobVariables = []runnerCommon.JobVariable{}
	for _, variable := range variables {
		jobVariables = append(jobVariables, runnerCommon.JobVariable{
			Key:      variable.Name,
			Value:    variable.Value,
			Masked:   variable.Masked,
			File:     variable.File,
			Internal: variable.Internal,
			Public:   variable.Public,
		})
	}
	return jobVariables
}
