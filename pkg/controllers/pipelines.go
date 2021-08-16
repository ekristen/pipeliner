package controllers

import (
	"time"

	"github.com/ekristen/pipeliner/pkg/models"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// LastState --
type LastState struct {
	State        string
	AllowFailure bool
	IsLastStage  bool
}

// PipelineStats --
type PipelineStats struct {
	Failed          int64
	Successful      int64
	Canceled        int64
	AllowedFailures int64
	OnFailure       int64
	Manual          int64
	Skipped         int64
}

// Pipelines --
type Pipelines struct {
	db  *gorm.DB
	log *logrus.Entry
}

// NewPipelines --
func NewPipelines(db *gorm.DB) *Pipelines {
	return &Pipelines{
		db:  db,
		log: logrus.WithField("component", "controller").WithField("controller", "pipelines"),
	}
}

// Run --
func (p *Pipelines) Run(pipeline *models.Pipeline) error {
	if pipeline.State == "pending" {
		if err := p.Update(pipeline, map[string]interface{}{"state": "running", "started_at": time.Now().UTC()}); err != nil {
			return err
		}
	}
	return nil
}

// Update --
func (p *Pipelines) Update(pipeline *models.Pipeline, updates map[string]interface{}) error {
	p.log.WithField("pipeline", pipeline.ID).Info("Called: Pipelines.Update")
	if err := p.db.Model(pipeline).Where("id = ?", pipeline.ID).Updates(updates).Error; err != nil {
		return err
	}

	return nil
}

// UpdateState --
func (p *Pipelines) UpdateState(pipelineID int64) error {
	log := p.log.WithField("pipeline", pipelineID)

	var lastStage models.PipelineStage
	sql := p.db.
		Model(&models.PipelineStage{}).
		Where("pipeline_id = ?", pipelineID).
		Order("`index` DESC").
		Limit(1).
		First(&lastStage)
	if sql.Error != nil {
		return sql.Error
	}

	var build models.Build
	sql = p.db.
		Model(&models.Build{}).
		Preload(clause.Associations).
		Where("pipeline_id = ? and state IN ?", pipelineID, []string{"pending", "running", "success", "failed", "skipped", "canceled"}).
		Order("updated_at desc").
		Limit(1).
		Select("state, allow_failure, pipeline_id, stage_idx").
		First(&build)
	if sql.Error != nil {
		log.WithError(sql.Error).Error("unable to find build")
		return sql.Error
	}

	log.WithField("pipeline", pipelineID).WithField("builds", sql.RowsAffected).Debug("found builds for pipeline")

	lastState := LastState{State: "skipped", AllowFailure: false, IsLastStage: false}
	if sql.RowsAffected == 1 {
		lastState = LastState{State: build.State, AllowFailure: build.AllowFailure, IsLastStage: false}
	}
	if lastState.State == "skipped" && lastStage.Index == build.StageIdx {
		lastState.IsLastStage = true
	}

	log.WithField("state", lastState.State).Debug("last state")
	switch lastState.State {
	case "canceled":
		return p.Cancel(pipelineID)
	case "success":
		return p.Success(pipelineID)
	case "failed":
		return p.Failed(pipelineID)
	case "skipped":
		if lastState.IsLastStage {
			return p.Success(pipelineID)
		}

		if err := p.Update(&build.Pipeline, map[string]interface{}{"state": "running"}); err != nil {
			log.WithError(err).Error("unable to update pipeline")
			return err
		}
	case "running":
		if err := p.Update(&build.Pipeline, map[string]interface{}{"state": "running"}); err != nil {
			log.WithError(err).Error("unable to update pipeline")
			return err
		}
	}

	return nil
}

// DoCancel --
func (p *Pipelines) DoCancel(pipelineID int64) error {
	var pipeline models.Pipeline

	sql := p.db.
		Model(&models.Pipeline{}).
		Preload(clause.Associations).
		Where("id = ?", pipelineID).
		First(&pipeline)
	if sql.Error != nil {
		return sql.Error
	}

	err := p.db.Transaction(func(tx *gorm.DB) error {

		for _, build := range pipeline.Builds {
			b := build
			b.State = "canceled"
			tx.Model(&b).Save(&b)
		}

		for _, stage := range pipeline.Stages {
			s := stage
			s.State = "canceled"
			tx.Model(&s).Save(&s)
		}

		pipeline.State = "canceled"
		tx.Model(&pipeline).Save(&pipeline)

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// Cancel --
func (p *Pipelines) Cancel(pipelineID int64) error {
	log := p.log.WithField("pipeline", pipelineID).WithField("state", "cancel")

	var pipeline models.Pipeline
	result := p.db.
		Model(&models.Pipeline{}).
		Preload(clause.Associations).
		Where("id = ?", pipelineID).
		Find(&pipeline)
	if result.Error != nil {
		log.WithError(result.Error).Error("unable to find pipeline")
		return result.Error
	}

	for _, stage := range pipeline.Stages {
		stage.State = "canceled"
		sql := p.db.
			Model(&stage).
			Save(&stage)
		if sql.Error != nil {
			log.WithError(sql.Error).Error("unable to update pipeline stage")
			return sql.Error
		}
	}

	updates := map[string]interface{}{
		"state": "canceled",
	}

	if pipeline.StartedAt != nil {
		finishedAt := time.Now().UTC()
		duration := finishedAt.Sub(*pipeline.StartedAt)

		updates = map[string]interface{}{
			"duration":    int(duration.Seconds()),
			"finished_at": finishedAt,
			"state":       "canceled",
		}
	}

	if err := p.Update(&pipeline, updates); err != nil {
		log.WithError(err).Error("unable to update pipeline")
		return err
	}

	return nil
}

// Success marks and handles the success of a pipeline
func (p *Pipelines) Success(pipelineID int64) error {
	log := p.log.WithField("pipeline", pipelineID).WithField("state", "success")

	var pipeline models.Pipeline
	result := p.db.
		Model(&models.Pipeline{}).
		Preload(clause.Associations).
		Where("id = ?", pipelineID).
		Find(&pipeline)
	if result.Error != nil {
		log.WithError(result.Error).Error("unable to find pipeline")
		return result.Error
	}

	buildStats, err := p.BuildsByState(pipelineID)
	if err != nil {
		log.WithError(err).Error("unable to get build stats")
		return err
	}
	totalBuilds := int64(len(pipeline.Builds))

	log.WithField("builds", totalBuilds).Debug("total builds")

	switch {
	case buildStats.Canceled > 0:
		if err := p.Cancel(pipelineID); err != nil {
			log.WithError(err).Error("unable to run cancel pipeline handler")
			return err
		}
	case buildStats.Failed > 0 || pipeline.State == "failed":
		if err := p.Failed(pipelineID); err != nil {
			log.WithError(err).Error("unable to run failed pipeline handler")
			return err
		}
	case buildStats.Successful+buildStats.AllowedFailures+buildStats.OnFailure+buildStats.Manual+buildStats.Skipped == totalBuilds:
		finishedAt := time.Now().UTC()
		duration := finishedAt.Sub(*pipeline.StartedAt)

		if err := p.Update(&pipeline, map[string]interface{}{
			"state":       "success",
			"duration":    int(duration.Seconds()),
			"finished_at": finishedAt,
		}); err != nil {
			log.WithError(err).Error("unable to update pipeline")
			return err
		}
	}

	return nil
}

// Failed marks and handles the failure of a pipline
func (p *Pipelines) Failed(pipelineID int64) error {
	log := logrus.WithField("pipelind", pipelineID).WithField("handler", "failed")

	var pipeline models.Pipeline
	err := p.db.
		Model(&models.Pipeline{}).
		Where("id = ?", pipelineID).
		Find(&pipeline).Error
	if err != nil {
		log.WithError(err).Error("unable to find pipeline")
		return err
	}

	finishedAt := time.Now().UTC()
	duration := finishedAt.Sub(*pipeline.StartedAt)

	if err := p.Update(&pipeline, map[string]interface{}{
		"duration":    int(duration.Seconds()),
		"finished_at": finishedAt,
		"state":       "failed",
	}); err != nil {
		log.WithError(err).Error("unable to update pipeline")
		return err
	}

	return nil
}

// BuildsByState returns stats based on build states
func (p *Pipelines) BuildsByState(pipelineID int64) (*PipelineStats, error) {
	log := p.log.WithField("pipeline", pipelineID)

	var failedBuilds int64
	err := p.db.
		Model(&models.Build{}).
		Where("pipeline_id = ? AND state = ? AND allow_failure = ?", pipelineID, "failed", false).
		Count(&failedBuilds).Error
	if err != nil {
		log.WithError(err).Error("unable to get failed build count")
		return nil, err
	}

	var successfulBuilds int64
	err = p.db.
		Model(&models.Build{}).
		Where("pipeline_id = ? AND state = ?", pipelineID, "success").
		Count(&successfulBuilds).Error
	if err != nil {
		log.WithError(err).Error("unable to get successful build count")
		return nil, err
	}

	var allowedFailures int64
	err = p.db.
		Model(&models.Build{}).
		Where("pipeline_id = ? AND state = ? AND allow_failure = ?", pipelineID, "failed", true).
		Count(&allowedFailures).Error
	if err != nil {
		log.WithError(err).Error("unable to get allowed failures build count")
		return nil, err
	}

	var onFailures int64
	sql := p.db.
		Model(&models.Build{}).
		Where("pipeline_id = ? AND state = \"created\" AND (`when` = \"on_failure\" OR `when` = \"always\")", pipelineID).
		Count(&onFailures)
	if sql.Error != nil {
		log.WithError(err).Error("unable to get failed build count")
		return nil, err
	}

	var skippedBuilds int64
	sql = p.db.
		Model(&models.Build{}).
		Where("pipeline_id = ? AND state = ?", pipelineID, "skipped").
		Count(&skippedBuilds)
	if sql.Error != nil {
		log.WithError(err).Error("unable to get skipped build count")
		return nil, err
	}

	var manualBuilds int64
	sql = p.db.
		Model(&models.Build{}).
		Where("pipeline_id = ? AND state = ? AND `when` = ?", pipelineID, "manual", "manual").
		Count(&manualBuilds)
	if sql.Error != nil {
		log.WithError(err).Error("unable to get manual build count")
		return nil, err
	}

	var canceledBuilds int64
	sql = p.db.
		Model(&models.Build{}).
		Where("pipeline_id = ? AND state = ?", pipelineID, "canceled").
		Count(&canceledBuilds)
	if sql.Error != nil {
		log.WithError(err).Error("unable to get manual build count")
		return nil, err
	}

	log.WithFields(logrus.Fields{
		"failed":           failedBuilds,
		"successful":       successfulBuilds,
		"skipped":          skippedBuilds,
		"allowed_failures": allowedFailures,
		"on_failure":       onFailures,
		"manual":           manualBuilds,
	}).Info("build stats")

	return &PipelineStats{
		Failed:          failedBuilds,
		Successful:      successfulBuilds,
		Canceled:        canceledBuilds,
		AllowedFailures: allowedFailures,
		OnFailure:       onFailures,
		Manual:          manualBuilds,
		Skipped:         skippedBuilds,
	}, nil
}
