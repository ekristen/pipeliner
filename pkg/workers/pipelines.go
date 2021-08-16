package workers

import (
	"github.com/ekristen/pipeliner/pkg/controllers"
	"github.com/ekristen/pipeliner/pkg/models"
	"github.com/ekristen/pipeliner/pkg/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Stats --
type Stats struct {
	Created                 int64
	Active                  int64
	Successful              int64
	Skipped                 int64
	Failed                  int64
	Canceled                int64
	Total                   int64
	AllowedFailures         int64
	NonBlockingManualBuilds int64
	BlockingManualBuilds    int64
}

// PipelineWorker --
type PipelineWorker struct {
	db                 *gorm.DB
	pipelineController *controllers.Pipelines
	buildController    *controllers.Builds
	log                *logrus.Entry
}

// NewPipelineWorker creates a pipeline worker
func NewPipelineWorker(db *gorm.DB, pipeline *controllers.Pipelines, build *controllers.Builds) *PipelineWorker {
	return &PipelineWorker{
		db:                 db,
		pipelineController: pipeline,
		buildController:    build,
		log:                logrus.WithField("component", "worker").WithField("worker", "pipline"),
	}
}

// Perform is called after all jobs are created, after all job updates
func (w *PipelineWorker) Perform(pipelineID int64) error {
	log := logrus.WithField("component", "platform-worker").WithField("id", pipelineID)

	log.Info("Processing builds for pipeline")

	stages, err := w.getStageIndexes(pipelineID)
	if err != nil {
		w.log.WithError(err).Error("unable to get pipeline stages")
		return err
	}

	for _, stage := range stages {
		if err := w.processStage(pipelineID, stage); err != nil {
			w.log.WithError(err).Error("unable to process stage")
			return err
		}
	}

	log.Info("Updating states for pipeline")
	if err := w.pipelineController.UpdateState(pipelineID); err != nil {
		log.WithError(err).Error("unable to update pipeline state")
	}

	return nil
}

func (w *PipelineWorker) getStageIndexes(pipelineID int64) ([]int64, error) {
	var stages []*models.PipelineStage

	sql := w.db.
		Model(&models.PipelineStage{}).
		Select("`index`").
		Distinct("`index`").
		Where("pipeline_id = ?", pipelineID).
		Order("`index` ASC").
		Find(&stages)
	if sql.Error != nil {
		return nil, sql.Error
	}

	var indexes = []int64{}
	for _, r := range stages {
		indexes = append(indexes, r.Index)
	}

	return indexes, nil
}

func (w *PipelineWorker) processStage(pipelineID, stageIdx int64) error {
	// Current Stage State
	state := w.stageState(pipelineID, stageIdx)
	var stage models.PipelineStage
	sql := w.db.
		Model(&models.PipelineStage{}).
		Where("pipeline_id = ? AND `index` = ?", pipelineID, stageIdx).
		First(&stage)
	if sql.Error != nil {
		return sql.Error
	}

	stage.State = state
	sql = w.db.
		Model(&stage).
		Save(&stage)
	if sql.Error != nil {
		return sql.Error
	}

	// TODO: verify we need to leave this here
	if state == "success" || state == "failed" || state == "skipped" {
		return nil
	}

	// Get state of the last stage
	hadPreviousFailure, _ := w.hadPreviousFailure(pipelineID, stageIdx)
	lastStageState := w.stageState(pipelineID, stageIdx-1)
	builds := w.buildController.BuildsForPipelineAndStage(pipelineID, stageIdx)
	for _, build := range builds {
		err := w.processBuild(build, lastStageState, hadPreviousFailure)
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *PipelineWorker) processBuild(build *models.Build, lastStageState string, hadPreviousFailure bool) error {
	log := w.log.WithField("handler", "processBuild")

	log.WithField("stage", 1).WithField("build", build.ID).WithField("state", lastStageState).Debug("processBuild")
	if build.State == "canceled" {
		w.buildController.Enqueue(build, "canceled", false)
	} else if w.validState(build.When, lastStageState, build.ID) {
		log.WithField("stage", 2).WithField("build", build.ID).WithField("state", lastStageState).Debug("processBuild")
		w.buildController.Enqueue(build, lastStageState, hadPreviousFailure)
	} else if lastStageState == "failed" && build.State != "skipped" {
		w.buildController.Enqueue(build, lastStageState, hadPreviousFailure)
	} else if lastStageState == "success" && build.When == "on_failure" {
		w.buildController.Enqueue(build, lastStageState, hadPreviousFailure)
	} else if lastStageState == "skipped" || lastStageState == "canceled" {
		w.buildController.Enqueue(build, lastStageState, hadPreviousFailure)
	}
	return nil
}

func (w *PipelineWorker) hadPreviousFailure(pipelineID int64, currentStage int64) (bool, error) {
	var build models.Build
	sql := w.db.
		Model(&models.Build{}).
		Where("pipeline_id = ? AND stage_idx < ? AND state = ?", pipelineID, currentStage, "failed").
		First(&build)
	if sql.Error != nil {
		return false, sql.Error
	}

	return sql.RowsAffected > 0, nil
}

func (w *PipelineWorker) stageState(pipelineID, stageIdx int64) string {
	stats, _ := w.stageBuildStats(pipelineID, stageIdx)

	switch {
	case stats.Canceled > 0:
		return "canceled"
	case stats.Created > 0:
		return "created"
	case stats.Total == 0:
		return "success"
	case stats.Total == stats.Successful+stats.AllowedFailures+stats.NonBlockingManualBuilds:
		return "success"
	case stats.Active > 0 || stats.BlockingManualBuilds > 0:
		return "running"
	case stats.Skipped == stats.Total:
		return "skipped"
	//case stats.Total == stats.Created:
	//	return "created"
	//case stats.Active == 0:
	//	return "pending"
	default:
		return "failed"
	}
}

// stageBuildStates obtains job states per stage
func (w *PipelineWorker) stageBuildStats(pipelineID, stageIdx int64) (*Stats, error) {
	var createdBuilds int64
	sql := w.db.Model(&models.Build{}).
		Where("pipeline_id = ? and stage_idx = ? and state IN ?", pipelineID, stageIdx, []string{"created"}).
		Count(&createdBuilds)
	if sql.Error != nil {
		return nil, sql.Error
	}

	var activeBuilds int64
	sql = w.db.Model(&models.Build{}).
		Where("pipeline_id = ? and stage_idx = ? and state IN ?", pipelineID, stageIdx, []string{"pending", "running"}).
		Count(&activeBuilds)
	if sql.Error != nil {
		return nil, sql.Error
	}

	var successfulBuilds int64
	sql = w.db.Model(&models.Build{}).
		Where("pipeline_id = ? and stage_idx = ? and state = ?", pipelineID, stageIdx, "success").
		Count(&successfulBuilds)
	if sql.Error != nil {
		return nil, sql.Error
	}

	var skippedBuilds int64
	sql = w.db.Model(&models.Build{}).
		Where("pipeline_id = ? and stage_idx = ? and state = ?", pipelineID, stageIdx, "skipped").
		Count(&skippedBuilds)
	if sql.Error != nil {
		return nil, sql.Error
	}

	var failedBuilds int64
	sql = w.db.Model(&models.Build{}).
		Where("pipeline_id = ? and stage_idx = ? and state = ?", pipelineID, stageIdx, "failed").
		Count(&failedBuilds)
	if sql.Error != nil {
		return nil, sql.Error
	}

	var canceledBuilds int64
	sql = w.db.Model(&models.Build{}).
		Where("pipeline_id = ? and stage_idx = ? and state = ?", pipelineID, stageIdx, "canceled").
		Count(&canceledBuilds)
	if sql.Error != nil {
		return nil, sql.Error
	}

	var totalBuilds int64
	sql = w.db.Model(&models.Build{}).
		Where("pipeline_id = ? and stage_idx = ?", pipelineID, stageIdx).
		Count(&totalBuilds)
	if sql.Error != nil {
		return nil, sql.Error
	}

	allowedFailures, err := w.stateCounter(pipelineID, stageIdx, "failed", true)
	if err != nil {
		return nil, err
	}
	nonBlockingManual, err := w.stateCounter(pipelineID, stageIdx, "manual", true)
	if err != nil {
		return nil, err
	}
	blockingManual, err := w.stateCounter(pipelineID, stageIdx, "manual", false)
	if err != nil {
		return nil, err
	}

	stats := &Stats{
		Created:                 createdBuilds,
		Active:                  activeBuilds,
		Failed:                  failedBuilds,
		Successful:              successfulBuilds,
		Skipped:                 skippedBuilds,
		Canceled:                canceledBuilds,
		Total:                   totalBuilds,
		AllowedFailures:         *allowedFailures,
		NonBlockingManualBuilds: *nonBlockingManual,
		BlockingManualBuilds:    *blockingManual,
	}

	logrus.WithFields(logrus.Fields{
		"pipeline":            pipelineID,
		"stage":               stageIdx,
		"created":             stats.Created,
		"failed":              stats.Failed,
		"skipped":             stats.Skipped,
		"active":              stats.Active,
		"success":             stats.Successful,
		"canceled":            stats.Canceled,
		"total":               stats.Total,
		"allow_failures":      stats.AllowedFailures,
		"non_blocking_manual": stats.NonBlockingManualBuilds,
		"blocking_manual":     stats.BlockingManualBuilds,
	}).Info("pipeline stats")

	return stats, nil
}

// stateCounter returns the number of builds in a specific state for a given stage
func (w *PipelineWorker) stateCounter(pipelineID, stageIdx int64, state string, allowFailure bool) (*int64, error) {
	var count int64

	sql := w.db.
		Model(&models.Build{}).
		Where("pipeline_id = ? and stage_idx = ? and state = ? and allow_failure = ?", pipelineID, stageIdx, state, allowFailure).
		Count(&count)
	if sql.Error != nil {
		return nil, sql.Error
	}

	return &count, nil
}

// validState returns true/false depending on the last stage state
func (w *PipelineWorker) validState(buildWhen string, lastStageState string, buildID int64) bool {
	w.log.WithField("build", buildID).WithField("when", buildWhen).WithField("last_stage_state", lastStageState).Debug("called: validState")
	switch buildWhen {
	case "on_success":
		return utils.StringSliceContains([]string{"success"}, lastStageState)
	case "on_failure":
		return utils.StringSliceContains([]string{"failed"}, lastStageState)
	case "always":
		// we need to make sure that the previous stage has reached a completed
		// state before allowing the build to be a valid state
		if !utils.StringSliceContains([]string{"success", "failed", "skipped"}, lastStageState) {
			return false
		}

		return true
	case "manual":
		return lastStageState == "success"
	default:
		return false
	}
}
