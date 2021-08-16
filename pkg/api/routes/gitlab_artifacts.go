package routes

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/ekristen/pipeliner/pkg/models"
	"github.com/ekristen/pipeliner/pkg/store"
)

// params: id, file, token, expire_id, artifact_type, artifact_format, metadata
func (c *Routes) gitlabArtifactsCreateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]
	token := r.Header.Get("job-token")

	// fmt.Println("id", r.FormValue("id"))
	// fmt.Println("expire_in", r.FormValue("expire_in"))
	// fmt.Println("artifact_type", r.FormValue("artifact_type"))
	// fmt.Println("artifact_format", r.FormValue("artifact_format"))
	// fmt.Println("metadata", r.FormValue("metadata"))

	log := logrus.WithField("build", id)

	if err := os.MkdirAll(fmt.Sprintf("/tmp/pipeline/%s", id), 0755); err != nil {
		log.WithError(err).Error("unable to mkdir")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		log.WithError(err).Error("unable to get form value")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	filename := fmt.Sprintf("%s/%s", id, handler.Filename)

	var buf bytes.Buffer
	_, err = io.Copy(&buf, file)
	if err != nil {
		log.WithError(err).Error("unable to copy to file")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := c.storage.Put(filename, buf.Bytes()); err != nil {
		log.WithError(err).Error("unable to write file")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var build models.Build
	result := c.db.
		Model(&models.Build{}).
		Where("id = ? and token = ?", id, token).
		First(&build)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Warn("unable to find build")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		log.WithError(err).Error("unable retrieve build from database")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	artifact := models.Artifact{
		File:       handler.Filename,
		Type:       r.FormValue("artifact_type"),
		Format:     r.FormValue("artifact_format"),
		Size:       int64(len(buf.Bytes())),
		BuildID:    build.ID,
		PipelineID: build.PipelineID,
	}
	sql := c.db.Create(&artifact)
	if sql.Error != nil {
		log.WithError(err).Error("unable to save artifact to the database")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := parseArtifact(c.db, c.storage, &artifact); err != nil {
		log.WithError(err).Error("unable to parse artifact to the database")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (c *Routes) gitlabArtifactsGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	buildID := vars["id"]
	buildToken := r.Header.Get("job-token")

	var artifact models.Artifact
	sql := c.db.
		Model(&models.Artifact{}).
		Joins("LEFT JOIN builds ON builds.id = artifacts.build_id").
		Where("artifacts.build_id = ? AND artifacts.type = ? AND builds.token = ?", buildID, "archive", buildToken).
		Find(&artifact)
	if sql.Error != nil {
		c.log.WithError(sql.Error).Error("unable to retrieve artifact from database")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if sql.RowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("content-type", "application/octet-stream")
	w.Header().Set("content-transfer-encoding", "binary")
	w.Header().Set("content-disposition", fmt.Sprintf("attachment; filename=%s", artifact.File))

	filename := fmt.Sprintf("%s/%s", buildID, artifact.File)

	contents, err := c.storage.Get(filename)
	if err != nil {
		c.log.WithError(err).Error("unable to open file")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = io.Copy(w, bytes.NewReader(contents))
}

func parseArtifact(db *gorm.DB, storage *store.Uploader, artifact *models.Artifact) error {
	if artifact.Type == "dotenv" {
		filename := fmt.Sprintf("%d/%s", artifact.BuildID, artifact.File)

		contents, err := storage.Get(filename)
		if err != nil {
			return err
		}

		gr, err := gzip.NewReader(bytes.NewReader(contents))
		if err != nil {
			return err
		}

		scanner := bufio.NewScanner(gr)
		for scanner.Scan() {
			parts := strings.Split(scanner.Text(), "=")
			if len(parts) == 2 {
				key := parts[0]
				value := parts[1]

				db.Create(&models.BuildVariable{
					Variable: &models.Variable{
						Name:   key,
						Value:  value,
						Masked: false,
						File:   false,
						// TODO: do we want all variables to be internal by default from a dotenv  file?
						Internal: true,
					},
					BuildID: artifact.BuildID,
				})
			}
		}

		if err := scanner.Err(); err != nil {
			return err
		}
	}

	return nil
}
