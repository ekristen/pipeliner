package routes

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	yamlParser "gitlab.com/gitlab-org/gitlab-runner/helpers/gitlab_ci_yaml_parser"
	"gopkg.in/yaml.v1"
	"gorm.io/gorm"

	"github.com/ekristen/pipeliner/pkg/models"
	"github.com/ekristen/pipeliner/pkg/utils"
)

func (c *Routes) workflowsListHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	var workflows []models.Workflow

	sql := c.db.Model(&models.Workflow{})

	if rawFields := query.Get("fields"); rawFields != "" {
		fields := strings.Split(rawFields, ",")
		if !utils.StringSliceContains(fields, "id") {
			fields = append(fields, "id")
		}
		sql = sql.Select(fields)
	}

	sql = sql.Find(&workflows)
	if sql.Error != nil {
		c.log.WithError(sql.Error).Error("unable to get workflows")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("content-type", "application/json")

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(workflows); err != nil {
		c.log.WithError(err).Error("unable to write workflows to http request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (c *Routes) workflowsAddHandler(w http.ResponseWriter, r *http.Request) {
	var workflow models.Workflow
	if err := json.NewDecoder(r.Body).Decode(&workflow); err != nil {
		c.log.WithError(err).Error("unable to save workflow")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := c.db.Model(&models.Workflow{}).Create(&workflow).Error; err != nil {
		c.log.WithError(err).Error("unable to get workflows")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.writeJSON(w, &workflow)
}

func (c *Routes) workflowsGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var workflow models.Workflow

	if err := c.db.Model(&models.Workflow{}).Where("id = ?", vars["id"]).First(&workflow).Error; err != nil {
		c.log.WithError(err).Error("unable to get workflows")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(workflow); err != nil {
		c.log.WithError(err).Error("unable to write workflows to http request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type createPipelineRequest struct {
	Variables []models.PipelineVariable `json:"variables"`
}

func (c *Routes) workflowsCreatePipelineHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]

	var req createPipelineRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.log.WithError(err).Error("unable to save workflow")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var workflow models.Workflow
	err := c.db.
		Model(&models.Workflow{}).
		Where("id = ?", id).
		First(&workflow).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.log.Warn("workflow not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		c.log.WithError(err).Error("error occurred")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	config := make(yamlParser.DataBag)
	if err := yaml.Unmarshal([]byte(workflow.Data), config); err != nil {
		c.log.WithError(err).Error("unable to unmarshal yaml")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := config.Sanitize(); err != nil {
		c.log.WithError(err).Error("unable to sanitize yaml config")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	pipeline := &models.Pipeline{
		State:    "pending",
		Workflow: workflow,
		Data:     []byte(workflow.Data),
	}

	if err := c.db.Create(pipeline).Error; err != nil {
		c.log.WithError(err).Error("unable to create pipeline")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := c.db.Model(pipeline).Association("Variables").Append(req.Variables); err != nil {
		c.log.WithError(err).Error("unable to create pipeline")
		c.writeJSONError(w, http.StatusInternalServerError, err)
		return
	}

	c.notifyBuilds <- pipeline

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(pipeline); err != nil {
		c.log.WithError(err).Error("unable to create pipeline")
		c.writeJSONError(w, http.StatusInternalServerError, err)
		return
	}
}
