package routes

import (
	"encoding/json"
	"net/http"

	"github.com/ekristen/pipeliner/pkg/models"
	"github.com/ekristen/pipeliner/pkg/models/scopes"
	"github.com/ekristen/pipeliner/pkg/utils"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// This allows a job that is not in a final state to be started (aka ran)
func (c *Routes) buildRunHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]

	var build models.Build

	validStates := []string{"manual"}

	err := c.db.Model(&models.Build{}).Where("id = ? AND state IN ?", id, validStates).Update("state", "pending").First(&build).Error
	if err != nil && err.Error() == "record not found" {
		w.WriteHeader(http.StatusConflict)
		return
	}
	if err != nil && err.Error() != "record not found" {
		c.log.WithError(err).Error("unable to get workflows")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(build); err != nil {
		c.log.WithError(err).Error("unable to write pipelines to http request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (c *Routes) buildCancelHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var build models.Build
	sql := c.db.
		Preload(clause.Associations).
		Model(&models.Build{}).
		Where("id = ?", id).
		First(&build)
	if sql.Error != nil {
		c.log.WithError(sql.Error).Error("unable to cancel build")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	state := "canceled"
	tBuild, err := c.controllerBuilds.TransitionState(&build, &state, &state, true)
	if err != nil {
		c.log.WithError(err).Error("unable to transition build")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.writeJSON(w, tBuild)
}

func (c *Routes) buildRetryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var build models.Build
	sql := c.db.
		Model(&models.Build{}).
		Preload(clause.Associations).
		Where("id = ?", id).
		First(&build)

	if sql.Error != nil {
		c.log.WithError(sql.Error).Error("unable to get build")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	newBuild, err := build.Clone()
	if err != nil {
		c.log.WithError(sql.Error).Error("unable to clone build")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := c.db.Transaction(func(tx *gorm.DB) error {
		sql = c.db.
			Model(&models.Build{}).
			Where("id = ?", build.ID).
			Updates(map[string]interface{}{"retried": true})
		if sql.Error != nil {
			c.log.WithError(sql.Error).Error("unable to update a build")
			return sql.Error
		}

		sql := c.db.Create(&newBuild)
		if sql.Error != nil {
			c.log.WithError(sql.Error).Error("unable to create a build")
			return sql.Error
		}

		return nil
	}); err != nil {
		c.log.WithError(sql.Error).Error("unable to update builds")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var pipeline models.Pipeline
	sql = c.db.
		Model(&models.Pipeline{}).
		Where("id = ?", build.PipelineID).
		First(&pipeline)
	if sql.Error != nil {
		c.log.WithError(sql.Error).Error("unable to get the pipeline")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.notifyPiplines <- &pipeline

	c.writeJSON(w, newBuild)
}

// This gets a list of jobs out of the database
func (c *Routes) buildListHandler(w http.ResponseWriter, r *http.Request) {
	var builds []models.Build

	var buildCount int64
	/*
		sql := c.db.
			Model(&models.Build{}).
			Joins("LEFT JOIN builds AS current ON builds.pipeline_id = current.pipeline_id AND builds.stage = current.stage AND builds.stage_idx = current.stage_idx AND builds.name = current.name AND builds.instance_idx < current.instance_idx").
			Where("current.instance_idx IS NULL").
			Order("builds.id DESC").
			Count(&buildCount)
	*/

	sql := c.db.
		Model(&models.Build{}).
		Scopes(scopes.BuildLatest).
		Order("id DESC").
		Count(&buildCount)
	if sql.Error != nil {
		c.log.WithError(sql.Error).Error("unable to get workflows")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	pagination, _ := utils.Paginate(r, buildCount)

	/*
		sql = c.db.
			Model(&models.Build{}).
			Preload(clause.Associations).
			Joins("LEFT JOIN builds AS current ON builds.pipeline_id = current.pipeline_id AND builds.stage = current.stage AND builds.stage_idx = current.stage_idx AND builds.name = current.name AND builds.instance_idx < current.instance_idx").
			Where("current.instance_idx IS NULL").
			Order("builds.id DESC").
			Scopes(scopes.BuildLatest).
			Limit(pagination.Limit).
			Offset(pagination.Offset).
			Find(&builds)
	*/

	sql = c.db.
		Model(&models.Build{}).
		Preload(clause.Associations).
		Scopes(scopes.BuildLatest).
		Order("id DESC").
		Limit(pagination.Limit).
		Offset(pagination.Offset).
		Find(&builds)
	if sql.Error != nil {
		c.log.WithError(sql.Error).Error("unable to get workflows")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.writeJSON(w, builds)
}

// This gets a single build from the database
func (c *Routes) buildGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var build models.Build

	sql := c.db.
		Model(&models.Build{}).
		Where("id = ?", id).
		First(&build)

	if sql.Error != nil {
		c.log.WithError(sql.Error).Error("unable to get workflows")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.writeJSON(w, build)
}

func (c *Routes) crudBuildsListTagsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var buildTags []models.BuildTag

	sql := c.db.
		Model(&models.BuildTag{}).
		Where("build_id = ?", id).
		Find(&buildTags)

	if sql.Error != nil {
		c.log.WithError(sql.Error).Error("unable to get build tags")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.writeJSON(w, buildTags)
}
