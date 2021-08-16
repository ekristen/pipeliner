package routes

import (
	"archive/zip"
	"bytes"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ekristen/pipeliner/pkg/models"
	"github.com/gorilla/mux"
)

// ArtifactFile --
type ArtifactFile struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Size       int64  `json:"size"`
	ArtifactID int64  `json:"artifact_id"`
}

func (c *Routes) crudPipelineArtifactsListHandler(w http.ResponseWriter, r *http.Request) {
	// swagger:operation GET /pipelines/{id}/artifacts pipeline listPipelineArtifacts
	// ---
	// summary: List the artifacts for a pipeline
	// parameters:
	// - name: id
	//   in: path
	//   description: id of the pipeline
	//   type: integer
	//   required: true
	// - name: page
	//   in: query
	//   description: page number of results to return (1-based)
	//   type: integer
	// - name: limit
	//   in: query
	//   description: page size of results
	//   type: integer
	// produces:
	// - application/json
	// responses:
	//   "200":
	//     "$ref": "#/responses/ArtifactList"

	vars := mux.Vars(r)
	query := r.URL.Query()

	pipelineID := vars["id"]

	var artifacts []models.Artifact

	sql := c.db.
		Model(&models.Artifact{}).
		Where("pipeline_id = ?", pipelineID).
		Order("created_at ASC")

	if query.Get("since") != "" {
		sql.Where("updated_at >= ?", query.Get("since"))
	}

	if err := sql.Find(&artifacts).Error; err != nil {
		c.log.WithError(err).Error("unable to get stages")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.writeJSON(w, artifacts)
}

func (c *Routes) crudBuildArtifactArchiveGetHandler(w http.ResponseWriter, r *http.Request) {
	// swagger:operation GET /builds/{id}/artifact/{type} build getBuildArtifactByType
	// ---
	// summary: Retrieve build artifact by type
	// parameters:
	// - name: id
	//   in: path
	//   description: id of the build/job
	//   type: integer
	//   required: true
	// - name: type
	//   in: path
	//   description: type of artifact to retrieve
	//   type: string
	//   required: true
	// produces:
	// - application/json
	// responses:
	//   "200":
	//     "$ref": "#/responses/Artifact"

	vars := mux.Vars(r)
	buildID := vars["id"]
	artifactType := vars["type"]

	var artifact models.Artifact

	sql := c.db.
		Model(&models.Artifact{}).
		Where("build_id = ? AND type = ?", buildID, artifactType)

	if err := sql.Find(&artifact).Error; err != nil {
		c.log.WithError(err).Error("unable to get stages")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.writeJSON(w, artifact)
}

func (c *Routes) crudBuildArtifactDownload(w http.ResponseWriter, r *http.Request) {
	// swagger:operation GET /artifacts/{id}/download artifact downloadArtifact
	// ---
	// summary: Download an artifact
	// parameters:
	// - name: id
	//   in: path
	//   description: id of the build/job
	//   type: integer
	//   required: true
	// produces:
	// - application/octet-stream
	// - application/tar+gzip
	// - application/zip
	// responses:
	//   "200":
	//     "$ref": "#/responses/Artifact"

	vars := mux.Vars(r)
	artifactID := vars["id"]

	var artifact models.Artifact

	sql := c.db.
		Model(&models.Artifact{}).
		Where("id = ?", artifactID)
	if err := sql.Find(&artifact).Error; err != nil {
		c.log.WithError(err).Error("unable to get artifact")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	contentType := "application/octet-stream"
	if artifact.Format == "gzip" {
		contentType = "application/tar+gzip"
	} else if artifact.Format == "zip" {
		contentType = "application/zip"
	}

	w.Header().Set("content-type", contentType)
	w.Header().Set("content-length", strconv.Itoa(int(artifact.Size)))
	w.Header().Set("content-transfer-encoding", "binary")
	w.Header().Set("content-disposition", fmt.Sprintf("attachment; filename=%s", artifact.File))

	filename := fmt.Sprintf("%d/%s", artifact.BuildID, artifact.File)
	contents, err := c.storage.Get(filename)
	if err != nil {
		c.log.WithError(err).Error("unable to read artifact")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(contents)
}

func (c *Routes) crudBuildArtifactContents(w http.ResponseWriter, r *http.Request) {
	// swagger:operation GET /artifacts/{id}/files artifact listArtifactFiles
	// ---
	// summary: List files in archive artifacts
	// parameters:
	// - name: id
	//   in: path
	//   description: id of the build/job
	//   type: integer
	//   required: true
	// produces:
	// - application/json
	// responses:
	//   "200":
	//     "$ref": "#/responses/ArtifactFiles"

	vars := mux.Vars(r)
	artifactID := vars["id"]

	var artifact models.Artifact

	sql := c.db.
		Model(&models.Artifact{}).
		Where("id = ?", artifactID)
	if err := sql.Find(&artifact).Error; err != nil {
		c.log.WithError(err).Error("unable to get artifact")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	filename := fmt.Sprintf("%d/%s", artifact.BuildID, artifact.File)
	contents, err := c.storage.Get(filename)
	if err != nil {
		c.log.WithError(err).Error("unable to read artifact")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// TODO: might be useful to stream instead of read in the whole file
	/*
		files := []*zip.FileHeader{}
		zr := zipstream.NewReader(f)
		//var ir io.Reader
		for {
			fh, err := zr.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				c.log.WithError(err).Error("unable to get artifact")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			files = append(files, fh)
		}
	*/

	br := bytes.NewReader(contents)
	zr, err := zip.NewReader(br, br.Size())
	if err != nil {
		c.log.WithError(err).Error("unable to get artifact")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	files := []ArtifactFile{}
	for i, f := range zr.File {
		files = append(files, ArtifactFile{
			ID:         int64(i),
			Name:       f.Name,
			Size:       int64(f.UncompressedSize64),
			ArtifactID: artifact.ID,
		})
	}

	c.writeJSON(w, files)
}
