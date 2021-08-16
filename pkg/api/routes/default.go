package routes

import (
	"encoding/json"
	"net/http"

	"github.com/ekristen/pipeliner/pkg/common"
)

type versionResponse struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func (c *Routes) defaultHandler(w http.ResponseWriter, r *http.Request) {
	if IsBrowser(r, true) {
		http.Redirect(w, r, "/ui/", 302)
		return
	}

	c.versionHandler(w, r)
}

func (c *Routes) versionHandler(w http.ResponseWriter, r *http.Request) {
	data := versionResponse{
		Name:    common.AppVersion.Name,
		Version: common.AppVersion.Summary,
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(200)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		c.log.WithError(err).Error("unable to write pipelines to http request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type statistics struct {
	Workflows int64 `json:"workflows"`
	Pipelines int64 `json:"pipelines"`
	Jobs      int64 `json:"jobs"`
	Artifacts int64 `json:"artifacts"`
	Runners   int64 `json:"runners"`
}

func (c *Routes) statsHandler(w http.ResponseWriter, r *http.Request) {
	var workflowsCount int64
	sql := c.db.Table("workflows").Count(&workflowsCount)
	if sql.Error != nil {
		c.log.WithError(sql.Error).Error("unable to get workflows")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var pipelinesCount int64
	sql = c.db.Table("pipelines").Count(&pipelinesCount)
	if sql.Error != nil {
		c.log.WithError(sql.Error).Error("unable to get pipelines")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var jobsCount int64
	sql = c.db.Table("builds").Count(&jobsCount)
	if sql.Error != nil {
		c.log.WithError(sql.Error).Error("unable to get jobs")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var artifactsCount int64
	sql = c.db.Table("artifacts").Count(&artifactsCount)
	if sql.Error != nil {
		c.log.WithError(sql.Error).Error("unable to get artifacts")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var runnersCount int64
	sql = c.db.Table("runners").Count(&runnersCount)
	if sql.Error != nil {
		c.log.WithError(sql.Error).Error("unable to get runners")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	stats := statistics{
		Workflows: workflowsCount,
		Pipelines: pipelinesCount,
		Jobs:      jobsCount,
		Artifacts: artifactsCount,
		Runners:   runnersCount,
	}

	c.writeJSON(w, stats)
}
