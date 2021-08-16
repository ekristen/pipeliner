package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm/clause"

	"github.com/ekristen/pipeliner/pkg/models"
	"github.com/ekristen/pipeliner/pkg/models/scopes"
)

func (c *Routes) pipelinesListHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	var pipelines []models.Pipeline

	sql := c.db.
		Model(&models.Pipeline{}).
		Preload(clause.Associations).
		Order("created_at DESC").
		Limit(20)

	if query.Get("since") != "" {
		sql.Where("updated_at >= ?", query.Get("since"))
	}

	if err := sql.Find(&pipelines).Error; err != nil {
		c.log.WithError(err).Error("unable to get workflows")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("content-type", "application/json")

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(pipelines); err != nil {
		c.log.WithError(err).Error("unable to write pipelines to http request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (c *Routes) pipelineStagesListHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pipelineID := vars["id"]

	var stages []models.PipelineStage
	if err := c.db.
		Model(&models.PipelineStage{}).
		Where("pipeline_id = ?", pipelineID).
		Order("`index` ASC").
		Find(&stages).Error; err != nil {
		c.log.WithError(err).Error("unable to get stages")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("content-type", "application/json")

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(stages); err != nil {
		c.log.WithError(err).Error("unable to write pipelines to http request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (c *Routes) pipelineStageBuildsListHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pipelineID := vars["id"]
	stageID := vars["stage_id"]

	var stage models.PipelineStage
	if err := c.db.
		Model(&models.PipelineStage{}).
		Where("pipeline_id = ? AND (id = ? OR `index` = ?)", pipelineID, stageID, stageID).
		First(&stage).Error; err != nil {
		c.log.WithError(err).Error("unable to get stages")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var builds []models.Build
	if err := c.db.
		Model(&models.Build{}).
		Where("pipeline_id = ? AND stage = ? AND stage_idx = ?", stage.PipelineID, stage.Name, stage.Index).
		Find(&builds).Error; err != nil {
		c.log.WithError(err).Error("unable to get stage builds")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.writeJSON(w, builds)
}

func (c *Routes) pipelineBuildsListHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	query := r.URL.Query()

	pipelineIDs := vars["id"]
	pipelineID, _ := strconv.Atoi(pipelineIDs)

	var builds []models.Build

	/*
		sql := c.db.
			Model(&models.Build{}).
			//Preload(clause.Associations).
			Joins("LEFT JOIN builds AS current ON builds.pipeline_id = current.pipeline_id AND builds.stage = current.stage AND builds.stage_idx = current.stage_idx AND builds.name = current.name AND builds.instance_idx > current.instance_idx").
			Where("current.instance_idx IS NULL").
			Where("builds.pipeline_id = ?", pipelineID).
			//Where("builds.retried = ?", false).
			Order("builds.stage_idx ASC").
			Order("builds.id ASC")
	*/

	sql := c.db.
		Model(&models.Build{}).
		Scopes(
			scopes.ForPipeline(int64(pipelineID)),
			//scopes.BuildLatest,
			scopes.OrderedByStage,
			scopes.Ordered,
			//scopes.OrderedByID,
		)

	if query.Get("since") != "" {
		sql.Where("updated_at >= ?", query.Get("since"))
	}

	if err := sql.Find(&builds).Error; err != nil {
		c.log.WithError(err).Error("unable to get stages")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("content-type", "application/json")

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(builds); err != nil {
		c.log.WithError(err).Error("unable to write pipelines to http request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
