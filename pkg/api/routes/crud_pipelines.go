package routes

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ekristen/pipeliner/pkg/models"
	"github.com/gorilla/mux"
	yamlParser "gitlab.com/gitlab-org/gitlab-runner/helpers/gitlab_ci_yaml_parser"
	"gopkg.in/yaml.v1"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (c *Routes) crudPipelinesGetHandler(w http.ResponseWriter, r *http.Request) {
	// swagger:operation GET /pipelines/{id} pipeline getPipeline
	// ---
	// summary: Get pipeline
	// parameters:
	// - name: id
	//   in: path
	//   description: id of the pipeline
	//   type: integer
	//   required: true
	// produces:
	// - application/json
	// responses:
	//   "200":
	//     "$ref": "#/responses/Pipeline"

	vars := mux.Vars(r)
	query := r.URL.Query()

	var pipeline models.Pipeline

	sql := c.db.
		Model(&models.Pipeline{}).
		Where("id = ?", vars["id"])

	if query.Get("since") != "" {
		sql.Where("updated_at >= ?", query.Get("since"))
	}

	if err := sql.Find(&pipeline).Error; err != nil {
		c.log.WithError(err).Error("unable to get workflows")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("content-type", "application/json")

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(pipeline); err != nil {
		c.log.WithError(err).Error("unable to write pipelines to http request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// swagger:parameters createPipeline
type pipelineRequest struct {
	WorkflowName *string           `json:"workflow_name,omitempty"`
	WorkflowID   *int64            `json:"workflow_id,omitempty"`
	Variables    []models.Variable `json:"variables,omitempty"`
}

func (c *Routes) crudPipelinesCreateHandler(w http.ResponseWriter, r *http.Request) {
	// swagger:operation POST /pipelines pipeline createPipeline
	// ---
	// summary: Create a pipeline
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// responses:
	//   "200":
	//     "$ref": "#/responses/Pipeline"

	if r.Header.Get("content-type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// TODO: validate unique variables
	// TODO: validate yaml

	var req pipelineRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.log.WithError(err).Error("unable to write pipelines to http request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if req.WorkflowID == nil && req.WorkflowName == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var workflow models.Workflow
	sql := c.db.Model(&models.Workflow{})
	if req.WorkflowID != nil {
		sql.Where("id = ?", req.WorkflowID)
	} else if req.WorkflowName != nil {
		sql.Where("name = ?", req.WorkflowName)
	}
	err := sql.First(&workflow).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.log.Warn("workflow not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		c.log.WithError(err).Error("an error occurred")
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

	if req.Variables != nil {
		for _, v := range req.Variables {
			if err := c.db.Create(&models.PipelineVariable{
				Variable:   &v,
				PipelineID: pipeline.ID,
			}).Error; err != nil {
				c.log.WithError(err).Error("unable to create pipeline variable")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	}

	c.notifyBuilds <- pipeline

	w.WriteHeader(http.StatusOK)
}

func (c *Routes) crudPipelinesCancelHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var builds []models.Build

	sql := c.db.
		Preload(clause.Associations).
		Model(&models.Build{}).
		Where("pipeline_id = ?", vars["id"]).
		Find(&builds)

	if sql.Error != nil {
		c.log.WithError(sql.Error).Error("unable to get pipeline")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	canceled := "canceled"
	for i, build := range builds {
		c.controllerBuilds.TransitionState(&build, &canceled, &canceled, i == len(builds)-1)
	}

	w.Header().Add("content-type", "application/json")

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(builds[0].Pipeline); err != nil {
		c.log.WithError(err).Error("unable to write pipeline to http request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
