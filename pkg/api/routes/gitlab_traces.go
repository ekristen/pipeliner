package routes

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/ekristen/pipeliner/pkg/models"
)

// TODO: handle patch abort
func (c *Routes) gitlabTraceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	buildID := vars["id"]
	token := r.Header.Get("job-token")

	log := logrus.WithField("handler", "trace").WithField("build", buildID)

	var build models.Build
	err := c.db.Model(&build).
		Where("id = ? AND token = ?", buildID, token).
		Find(&build).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		log.Warn("unable to find build for trace")
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if err != nil {
		log.WithError(err).Error("unable to get build")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	contentLength := r.Header.Get("content-length")
	contentRange := r.Header.Get("content-range")
	ranges := strings.Split(contentRange, "-")
	lower, _ := strconv.Atoi(ranges[0])
	upper, _ := strconv.Atoi(ranges[1])
	size, _ := strconv.Atoi(contentLength)

	log.WithField("range", contentRange).WithField("length", contentLength).Debug("received trace details")

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.WithError(err).Error("unable to read trace http request body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var trace models.Trace

	sql := c.db.
		Model(&models.Trace{}).
		Where("build_id = ?", build.ID).
		Assign(models.Trace{BuildID: build.ID, Build: build}).
		FirstOrCreate(&trace)
	if sql.Error != nil {
		log.WithError(sql.Error).Error("unable to save trace to database")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tracePart := models.TracePart{
		Trace: trace,
		Start: lower,
		End:   upper,
		Size:  size,
		Data:  data,
	}

	sql = c.db.Create(&tracePart)
	if sql.Error != nil {
		log.WithError(sql.Error).Error("unable to save trace part to database")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	/*
		os.MkdirAll(fmt.Sprintf("/tmp/pipeline/%s", buildID), 0755)

		f, err := os.OpenFile("/tmp/pipeline/"+buildID+"/log.txt", os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.WithError(err).Error("unable to open trace file")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer f.Close()

		whence := io.SeekStart
		_, err = f.Seek(int64(lower), whence)
		if err != nil {
			log.WithError(err).Error("unable to seek trace file")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		f.Write(data)
		f.Sync() //flush to disk
	*/

	w.Header().Set("Job-Status", build.State)
	w.Header().Set("Range", fmt.Sprintf("0-%d", upper))
	w.WriteHeader(http.StatusAccepted)
}
