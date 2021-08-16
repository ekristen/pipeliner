package workers

import (
	"github.com/ekristen/pipeliner/pkg/models"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// BuildsWorker --
type BuildsWorker struct {
	db             *gorm.DB
	PipelineWorker *PipelineWorker
	notify         chan *models.Pipeline
}

// NewBuildsWorker --
func NewBuildsWorker(db *gorm.DB, PipelineWorker *PipelineWorker, notify chan *models.Pipeline) *BuildsWorker {
	return &BuildsWorker{
		db:             db,
		PipelineWorker: PipelineWorker,
		notify:         notify,
	}
}

// Perform --
// -- Called after all builds are created
// -- Called after every build status update
func (w *BuildsWorker) Perform(pipeline *models.Pipeline) {
	log := logrus.WithField("component", "builds-worker").WithField("id", pipeline.ID)

	if err := w.PipelineWorker.buildController.CreateFromPipeline(pipeline); err != nil {
		log.WithError(err).Error("unable to create builds")
	}

	log.Info("Builds created successfully")

	// Notify the PipelineWorker Channel
	w.notify <- pipeline
}
