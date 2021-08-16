package routes

import (
	"context"
	"net/http"

	"github.com/go-http-utils/etag"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/ekristen/pipeliner/pkg/controllers"
	"github.com/ekristen/pipeliner/pkg/models"
	"github.com/ekristen/pipeliner/pkg/store"
)

// Routes --
type Routes struct {
	ctx                 context.Context
	db                  *gorm.DB
	r                   *mux.Router
	log                 *logrus.Entry
	storage             *store.Uploader
	notifyPiplines      chan *models.Pipeline
	notifyBuilds        chan *models.Pipeline
	controllerPipelines *controllers.Pipelines
	controllerBuilds    *controllers.Builds
}

// NewRoutes returns a Routes instance --
func NewRoutes(
	ctx context.Context,
	db *gorm.DB,
	r *mux.Router,
	log *logrus.Entry,
	storage *store.Uploader,
	notifyPipelines chan *models.Pipeline,
	notifyBuilds chan *models.Pipeline,
	controllerPipelines *controllers.Pipelines,
	controllerBuilds *controllers.Builds,
) *Routes {
	routes := &Routes{
		ctx:                 ctx,
		db:                  db,
		r:                   r,
		log:                 log,
		storage:             storage,
		notifyPiplines:      notifyPipelines,
		notifyBuilds:        notifyBuilds,
		controllerPipelines: controllerPipelines,
		controllerBuilds:    controllerBuilds,
	}

	routes.r.Path("/").HandlerFunc(routes.defaultHandler)

	// UI Router
	ui := routes.r.PathPrefix("/ui").Subrouter()
	ui.Use(func(next http.Handler) http.Handler {
		return etag.Handler(next, true)
	})
	ui.PathPrefix("/").Handler(http.StripPrefix("/ui/", routes.uiHandler()))

	// GitLab API Router
	// These re-implement all the apis needed by the
	// GitLab Runner to operate properly.
	gitlabAPI := routes.r.PathPrefix("/api/v4").Subrouter()
	gitlabAPI.Path("/jobs/request").HandlerFunc(routes.jobsRequestHandler)
	gitlabAPI.Path("/jobs/{id}").HandlerFunc(routes.jobHandler)
	gitlabAPI.Path("/jobs/{id}/trace").HandlerFunc(routes.gitlabTraceHandler)
	gitlabAPI.Path("/jobs/{id}/artifacts").Methods("POST").HandlerFunc(routes.gitlabArtifactsCreateHandler)
	gitlabAPI.Path("/jobs/{id}/artifacts").Methods("GET").HandlerFunc(routes.gitlabArtifactsGetHandler)
	gitlabAPI.Path("/runners").Methods("POST").HandlerFunc(routes.gitlabRunnersCreateHandler)
	gitlabAPI.Path("/runners/verify").HandlerFunc(routes.runnersVerifyHandler)

	// Pipeliner API
	// These implement all the APIs available and needed
	// by the pipeliner system
	api := routes.r.PathPrefix("/v1").Subrouter()
	/*
		api.Use(func(next http.Handler) http.Handler {
			return etag.Handler(next, true)
		})
	*/

	api.Path("/version").Methods("GET").HandlerFunc(routes.versionHandler)
	api.Path("/lint").Methods("POST").HandlerFunc(routes.lintWorkflowHandler)

	api.Path("/workflows").Methods("GET").HandlerFunc(routes.workflowsListHandler)
	api.Path("/workflows").Methods("POST").HandlerFunc(routes.workflowsAddHandler)
	api.Path("/workflows/{id}").Methods("GET").HandlerFunc(routes.workflowsGetHandler)
	api.Path("/workflows/{id}").Methods("DELETE").HandlerFunc(routes.crudWorkflowsDeleteHandler)
	api.Path("/workflows/{id}").Methods("PATCH").HandlerFunc(routes.crudWorkflowsPatchHandler)
	api.Path("/workflows/{id}/pipeline").Methods("POST").HandlerFunc(routes.workflowsCreatePipelineHandler)

	api.Path("/workflows/{workflow_id}/variables").Methods("GET").HandlerFunc(routes.crudVariablesListHandler)
	api.Path("/workflows/{workflow_id}/variables").Methods("POST").HandlerFunc(routes.crudVariablesCreateHandler)
	api.Path("/workflows/{workflow_id}/variables/{id}").Methods("GET").HandlerFunc(routes.crudVariablesGetHandler)
	api.Path("/workflows/{workflow_id}/variables/{id}").Methods("DELETE").HandlerFunc(routes.crudVariablesDeleteHandler)

	api.Path("/pipelines").Methods("GET").HandlerFunc(routes.pipelinesListHandler)
	api.Path("/pipelines").Methods("POST").HandlerFunc(routes.crudPipelinesCreateHandler)
	api.Path("/pipelines/{id}").Methods("GET").HandlerFunc(routes.crudPipelinesGetHandler)
	api.Path("/pipelines/{id}/artifacts").HandlerFunc(routes.crudPipelineArtifactsListHandler)
	api.Path("/pipelines/{id}/builds").HandlerFunc(routes.pipelineBuildsListHandler)
	api.Path("/pipelines/{id}/stages").HandlerFunc(routes.pipelineStagesListHandler)
	api.Path("/pipelines/{id}/stages/{stage_id}/builds").HandlerFunc(routes.pipelineStageBuildsListHandler)
	api.Path("/pipelines/{id}/cancel").HandlerFunc(routes.crudPipelinesCancelHandler)

	api.Path("/pipelines/{pipeline_id}/variables").Methods("GET").HandlerFunc(routes.crudVariablesListHandler)
	api.Path("/pipelines/{pipeline_id}/variables").Methods("POST").HandlerFunc(routes.crudVariablesCreateHandler)
	api.Path("/pipelines/{pipeline_id}/variables/{id}").Methods("GET").HandlerFunc(routes.crudVariablesGetHandler)
	api.Path("/pipelines/{pipeline_id}/variables/{id}").Methods("DELETE").HandlerFunc(routes.crudVariablesDeleteHandler)

	api.Path("/builds").HandlerFunc(routes.buildListHandler)
	api.Path("/builds/{id}").HandlerFunc(routes.buildGetHandler)
	api.Path("/builds/{id}/run").HandlerFunc(routes.buildRunHandler)
	api.Path("/builds/{id}/retry").HandlerFunc(routes.buildRetryHandler)
	api.Path("/builds/{id}/cancel").HandlerFunc(routes.buildCancelHandler)
	api.Path("/builds/{id}/trace").HandlerFunc(routes.traceCRUDGetHandler)
	api.Path("/builds/{id}/artifact/{type}").HandlerFunc(routes.crudBuildArtifactArchiveGetHandler)
	api.Path("/builds/{id}/tags").HandlerFunc(routes.crudBuildsListTagsHandler)

	api.Path("/traces").Methods("GET").HandlerFunc(routes.defaultHandler)

	api.Path("/artifacts/{id}/download").Methods("GET").HandlerFunc(routes.crudBuildArtifactDownload)
	api.Path("/artifacts/{id}/files").Methods("GET").HandlerFunc(routes.crudBuildArtifactContents)

	api.Path("/variables").Methods("GET").HandlerFunc(routes.crudVariablesListHandler)
	api.Path("/variables").Methods("POST").HandlerFunc(routes.crudVariablesCreateHandler)
	api.Path("/variables/{name}").Methods("GET").HandlerFunc(routes.crudVariablesGetHandler)
	api.Path("/variables/{name}").Methods("DELETE").HandlerFunc(routes.crudVariablesDeleteHandler)

	api.Path("/runners").Methods("GET").HandlerFunc(routes.crudRunnersListHandler)
	api.Path("/runners/{id}").Methods("GET").HandlerFunc(routes.crudRunnersGetHandler)
	api.Path("/runners/{id}").Methods("PATCH").HandlerFunc(routes.crudRunnersPatchHandler)

	api.Path("/tokens").Methods("GET").HandlerFunc(routes.crudTokensListHandler)
	api.Path("/tokens").Methods("POST").HandlerFunc(routes.crudTokensAddHandler)

	api.Path("/stats").Methods("GET").HandlerFunc(routes.statsHandler)

	return routes
}
