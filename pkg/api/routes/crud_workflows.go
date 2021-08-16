package routes

import (
	"encoding/json"
	"net/http"

	"github.com/ekristen/pipeliner/pkg/models"
	"github.com/gorilla/mux"
)

func (c *Routes) crudWorkflowsCreateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var workflow models.Workflow

	if err := c.db.
		Model(&models.Workflow{}).
		Where("id = ?", vars["id"]).
		Find(&workflow).Error; err != nil {
		c.log.WithError(err).Error("unable to get workflows")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("content-type", "application/json")

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(workflow); err != nil {
		c.log.WithError(err).Error("unable to write pipelines to http request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (c *Routes) crudWorkflowsDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var workflow models.Workflow

	sql := c.db.
		Model(&models.Workflow{}).
		Where("id = ?", vars["id"]).
		First(&workflow)
	if sql.Error != nil {
		c.log.WithError(sql.Error).Error("unable to get workflows")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sql = c.db.Delete(&workflow)
	if sql.Error != nil {
		c.log.WithError(sql.Error).Error("unable to delete workflows")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.writeJSON(w, workflow)
}

func (c *Routes) crudWorkflowsPatchHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var patchWorkflow models.Workflow
	if err := json.NewDecoder(r.Body).Decode(&patchWorkflow); err != nil {
		c.log.WithError(err).Error("unable to save workflow")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var workflow models.Workflow

	sql := c.db.
		Model(&models.Workflow{}).
		Where("id = ?", vars["id"]).
		First(&workflow)
	if sql.Error != nil {
		c.log.WithError(sql.Error).Error("unable to get workflows")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sql = c.db.Model(&workflow).Select("data").Updates(&patchWorkflow)
	if sql.Error != nil {
		c.log.WithError(sql.Error).Error("unable to update workflows")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.writeJSON(w, workflow)
}
