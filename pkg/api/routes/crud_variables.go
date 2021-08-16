package routes

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ekristen/pipeliner/pkg/models"
	"github.com/gorilla/mux"
)

func (c *Routes) crudVariablesListHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	isGlobal := true
	isWorkflow := false
	isPipeline := false

	if _, ok := vars["workflow_id"]; ok {
		isGlobal = false
		isWorkflow = true
	}

	if _, ok := vars["pipeline_id"]; ok {
		isGlobal = false
		isPipeline = true
	}

	if isGlobal {
		var variables []*models.GlobalVariable
		if err := c.db.Model(&models.GlobalVariable{}).Find(&variables).Error; err != nil {
			c.log.WithError(err).Error("unable to get variables")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Add("content-type", "application/json")

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(variables); err != nil {
			c.log.WithError(err).Error("unable to write pipelines to http request")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else if isWorkflow {
		var variables []*models.WorkflowVariable
		if err := c.db.Model(&models.WorkflowVariable{}).Where("workflow_id = ?", vars["workflow_id"]).Find(&variables).Error; err != nil {
			c.log.WithError(err).Error("unable to get workflow variables")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Add("content-type", "application/json")

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(variables); err != nil {
			c.log.WithError(err).Error("unable to write pipelines to http request")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else if isPipeline {
		var variables []*models.PipelineVariable
		if err := c.db.Model(&models.PipelineVariable{}).Where("pipeline_id = ?", vars["pipeline_id"]).Find(&variables).Error; err != nil {
			c.log.WithError(err).Error("unable to get pipeline variables")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Add("content-type", "application/json")

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(variables); err != nil {
			c.log.WithError(err).Error("unable to write pipelines to http request")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		// this should not be possible
	}
}

func (c *Routes) crudVariablesGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	isGlobal := true
	isWorkflow := false
	isPipeline := false

	if _, ok := vars["workflow_id"]; ok {
		isGlobal = false
		isWorkflow = true
	}

	if _, ok := vars["pipeline_id"]; ok {
		isGlobal = false
		isPipeline = true
	}

	if _, ok := vars["id"]; !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`missing id in request`))
		return
	}

	id := vars["id"]

	if isGlobal {
		var variable models.GlobalVariable
		if err := c.db.Model(&models.GlobalVariable{}).Where("id = ?", id).Find(&variable).Error; err != nil {
			c.log.WithError(err).Error("unable to write pipelines to http request")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Add("content-type", "application/json")

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(variable); err != nil {
			c.log.WithError(err).Error("unable to write pipelines to http request")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else if isWorkflow {
		var variable models.WorkflowVariable
		if err := c.db.Model(&models.WorkflowVariable{}).Where("id = ?", id).Find(&variable).Error; err != nil {
			c.log.WithError(err).Error("unable to write pipelines to http request")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Add("content-type", "application/json")

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(variable); err != nil {
			c.log.WithError(err).Error("unable to write pipelines to http request")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else if isPipeline {
		var variable models.PipelineVariable
		if err := c.db.Model(&models.PipelineVariable{}).Where("id = ?", id).Find(&variable).Error; err != nil {
			c.log.WithError(err).Error("unable to write pipelines to http request")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Add("content-type", "application/json")

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(variable); err != nil {
			c.log.WithError(err).Error("unable to write pipelines to http request")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		// this is bad
	}
}

func (c *Routes) crudVariablesCreateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	isGlobal := true
	isWorkflow := false
	isPipeline := false

	if _, ok := vars["workflow_id"]; ok {
		isGlobal = false
		isWorkflow = true
	}

	if _, ok := vars["pipeline_id"]; ok {
		isGlobal = false
		isPipeline = true
	}

	variable := &models.Variable{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := json.NewDecoder(r.Body).Decode(&variable); err != nil {
		c.log.WithError(err).Error("unable to write pipelines to http request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if isGlobal {
		var variable = models.GlobalVariable{
			Variable: variable,
		}
		if err := c.db.Model(&models.GlobalVariable{}).Create(&variable).Error; err != nil {
			if strings.Contains(err.Error(), "Duplicate") {
				c.log.WithError(err).Error("duplicate entry detected")
				w.WriteHeader(http.StatusConflict)
				return
			}
			c.log.WithError(err).Error("unable to write pipelines to http request")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Add("content-type", "application/json")

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(variable); err != nil {
			c.log.WithError(err).Error("unable to write pipelines to http request")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	} else if isWorkflow {
		i, err := strconv.Atoi(vars["workflow_id"])
		if err != nil {
			c.log.WithError(err).Error("unable to write pipelines to http request")
			w.WriteHeader(http.StatusInternalServerError)
		}

		var variable = models.WorkflowVariable{
			Variable:   variable,
			WorkflowID: int64(i),
		}

		if err := c.db.Model(&models.WorkflowVariable{}).Create(&variable).Error; err != nil {
			if strings.Contains(err.Error(), "Duplicate") {
				c.log.WithError(err).Error("duplicate entry detected")
				w.WriteHeader(http.StatusConflict)
				return
			}
			c.log.WithError(err).Error("unable to write pipelines to http request")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Add("content-type", "application/json")

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(variable); err != nil {
			c.log.WithError(err).Error("unable to write pipelines to http request")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	} else if isPipeline {
		i, err := strconv.Atoi(vars["pipeline_id"])
		if err != nil {
			c.log.WithError(err).Error("unable to write pipelines to http request")
			w.WriteHeader(http.StatusInternalServerError)
		}

		var variable = models.PipelineVariable{
			Variable:   variable,
			PipelineID: int64(i),
		}

		if err := c.db.Model(&models.PipelineVariable{}).Create(&variable).Error; err != nil {
			if strings.Contains(err.Error(), "Duplicate") {
				c.log.WithError(err).Error("duplicate entry detected")
				w.WriteHeader(http.StatusConflict)
				return
			}
			c.log.WithError(err).Error("unable to write pipelines to http request")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Add("content-type", "application/json")

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(variable); err != nil {
			c.log.WithError(err).Error("unable to write pipelines to http request")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	} else {
		// this is bad
	}
}

func (c *Routes) crudVariablesDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	isGlobal := true
	isWorkflow := false
	isPipeline := false

	if _, ok := vars["workflow_id"]; ok {
		isGlobal = false
		isWorkflow = true
	}

	if _, ok := vars["pipeline_id"]; ok {
		isGlobal = false
		isPipeline = true
	}

	if _, ok := vars["name"]; !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`missing name in request`))
		return
	}

	name := vars["name"]

	if isGlobal {
		var variable models.GlobalVariable
		if err := c.db.Model(&models.GlobalVariable{}).Where("name = ?", name).First(&variable).Error; err != nil {
			c.log.WithError(err).Error("unable to write pipelines to http request")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := c.db.Delete(&variable).Error; err != nil {
			c.log.WithError(err).Error("unable to delete global variable")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		c.writeJSON(w, variable)
	} else if isWorkflow {
		if err := c.db.Model(&models.WorkflowVariable{}).Delete(&models.WorkflowVariable{}, name).Error; err != nil {
			c.log.WithError(err).Error("unable to write pipelines to http request")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	} else if isPipeline {
		if err := c.db.Model(&models.PipelineVariable{}).Delete(&models.PipelineVariable{}, name).Error; err != nil {
			c.log.WithError(err).Error("unable to write pipelines to http request")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	} else {
		// this is bad
	}
}

func (c *Routes) writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Add("content-type", "application/json")

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		c.log.WithError(err).Error("unable to write pipelines to http request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (c *Routes) writeJSONError(w http.ResponseWriter, statusCode int, err error) {
	data := struct {
		Status string   `json:"status"`
		Errors []string `json:"errors"`
	}{
		Status: "error",
		Errors: []string{err.Error()},
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		c.log.WithError(err).Error("unable to write pipelines to http request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
