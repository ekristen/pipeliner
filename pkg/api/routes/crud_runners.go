package routes

import (
	"encoding/json"
	"net/http"

	"github.com/ekristen/pipeliner/pkg/models"
	"github.com/gorilla/mux"
	"gorm.io/gorm/clause"
)

func (c *Routes) crudRunnersListHandler(w http.ResponseWriter, r *http.Request) {
	var runners []models.Runner

	sql := c.db.
		Model(&models.Runner{}).
		Preload(clause.Associations).
		Find(&runners)

	if sql.Error != nil {
		c.log.WithError(sql.Error).Error("unable to encode json response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("content-type", "application/json")

	if err := json.NewEncoder(w).Encode(&runners); err != nil {
		c.log.WithError(err).Error("unable to encode json response")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (c *Routes) crudRunnersGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]

	var runner models.Runner

	sql := c.db.
		Model(&models.Runner{}).
		Preload(clause.Associations).
		Where("id = ?", id).
		Find(&runner)

	if sql.Error != nil {
		c.log.WithError(sql.Error).Error("unable to encode json response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.writeJSON(w, runner)
}

// runnerPatchRequest --
type runnerPatchRequest struct {
	Description *string   `json:"description,omitempty"`
	Tags        []*string `json:"tags,omitempty"`
	RunUntagged *bool     `json:"run_untagged,omitempty"`
}

func (c *Routes) crudRunnersPatchHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var patch runnerPatchRequest
	if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
		c.log.WithError(err).Error("unable to decode runner patch http request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var tags []models.RunnerTag
	var runner models.Runner

	sql := c.db.
		Model(&models.Runner{}).
		Preload(clause.Associations).
		Where("id = ?", id).
		Find(&runner)

	if sql.Error != nil {
		c.log.WithError(sql.Error).Error("unable to encode json response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	updates := map[string]interface{}{}
	if patch.Description != nil {
		updates["description"] = patch.Description
	}
	if patch.RunUntagged != nil {
		updates["run_untagged"] = patch.RunUntagged
	}
	if patch.Tags != nil {
		for _, tag := range patch.Tags {
			tags = append(tags, models.RunnerTag{
				Tag:      *tag,
				RunnerID: runner.ID,
			})
		}
	}

	sql = c.db.
		Model(&runner).
		Updates(updates)
	if sql.Error != nil {
		c.log.WithError(sql.Error).Error("unable to update runner")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err := c.db.Debug().
		Model(&runner).
		Association("Tags").
		Replace(tags)
	if err != nil {
		c.log.WithError(err).Error("unable to associations for runner")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.writeJSON(w, runner)
}
