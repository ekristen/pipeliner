package routes

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	runnerCommon "gitlab.com/gitlab-org/gitlab-runner/common"

	"github.com/ekristen/pipeliner/pkg/models"
	"github.com/ekristen/pipeliner/pkg/utils"
)

// runnerCommon.RegisterRunnerRequest
// StatusCreated -- ok
// StatusForbidden -- not ok
// Method -- POST
// TODO: generate token?!
func (c *Routes) gitlabRunnersCreateHandler(w http.ResponseWriter, r *http.Request) {
	var runnerRequest runnerCommon.RegisterRunnerRequest
	if err := json.NewDecoder(r.Body).Decode(&runnerRequest); err != nil {
		c.log.WithError(err).Error("unable to decode json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// logrus.Info(runnerRequest)
	// TODO: compare register token

	runUntagged := true
	if runnerRequest.Tags != "" {
		runUntagged = false
	}

	var tags []models.RunnerTag
	if runnerRequest.Tags != "" {
		for _, tag := range strings.Split(runnerRequest.Tags, ",") {
			tags = append(tags, models.RunnerTag{
				Tag: tag,
			})
		}
	}

	runner := &models.Runner{
		Active:       true,
		Architecture: runnerRequest.Info.Architecture,
		Description:  runnerRequest.Description,
		Global:       true,
		Locked:       false,
		Name:         runnerRequest.Info.Name,
		Platform:     runnerRequest.Info.Platform,
		Token:        utils.RandomString(16),
		Tags:         tags,
		RunUntagged:  runUntagged,
	}

	if err := c.db.Create(runner).Error; err != nil {
		c.log.WithError(err).Error("unable to create runner")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := runnerCommon.RegisterRunnerResponse{
		Token: runner.Token,
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		c.log.WithError(err).Error("unable to encode json response")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (c *Routes) runnersVerifyHandler(w http.ResponseWriter, r *http.Request) {
	all, _ := ioutil.ReadAll(r.Body)
	w.WriteHeader(http.StatusOK)
	w.Write(all)
}
