package routes

import (
	"encoding/json"
	"net/http"

	"github.com/ekristen/pipeliner/pkg/gitlab"
)

// LintInput --
type LintInput struct {
	Data string `json:"data"`
}

func (c *Routes) lintWorkflowHandler(w http.ResponseWriter, r *http.Request) {
	var req LintInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.log.WithError(err).Error("unable to decode json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	v := gitlab.NewYAMLValidator([]byte(req.Data))
	if err := v.Validate(); err != nil {
		c.log.WithError(err).Error("unable to validate yaml config")
		c.writeJSONError(w, http.StatusUnprocessableEntity, err)
		return
	}

	c.writeJSON(w, struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	})
}
