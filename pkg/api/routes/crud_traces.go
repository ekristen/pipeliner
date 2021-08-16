package routes

import (
	"encoding/json"
	"net/http"

	terminal "github.com/buildkite/terminal-to-html/v3"
	"github.com/ekristen/pipeliner/pkg/models"
	"github.com/gorilla/mux"
	"gorm.io/gorm/clause"
)

type traceResponse struct {
	BuildID string `json:"build_id"`
	Data    []byte `json:"data"`
}

func (c *Routes) traceCRUDGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	buildID := vars["id"]

	var trace models.Trace
	sql := c.db.
		Preload(clause.Associations).
		Model(&models.Trace{}).
		Where("build_id = ?", buildID).
		First(&trace)
	if sql.Error != nil {
		if sql.Error.Error() == "record not found" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		c.log.WithError(sql.Error).Error("unable to open trace file")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var dat []byte
	for _, tp := range trace.Parts {
		dat = append(dat, tp.Data...)
	}

	tRes := traceResponse{
		BuildID: buildID,
		Data:    terminal.Render(dat),
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(tRes); err != nil {
		c.log.WithError(err).Error("unable to write traces")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
