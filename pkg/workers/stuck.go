package workers

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/ekristen/pipeliner/pkg/common"
	"github.com/ekristen/pipeliner/pkg/models"
)

// StuckWorker --
type StuckWorker struct {
	db *gorm.DB
}

// NewStuckWorker --
func NewStuckWorker(db *gorm.DB) *StuckWorker {
	return &StuckWorker{
		db: db,
	}
}

// Perform --
func (w *StuckWorker) Perform() error {
	if !w.obtainLease() {
		return nil
	}

	stmt := &gorm.Statement{DB: w.db}
	stmt.Parse(&models.Build{})

	w.drop("running", common.BuildRunningOutdatedTimeout, fmt.Sprintf("%s.updated_at < ?", stmt.Schema.Table), "stuck_or_timeout_failure")
	w.drop("pending", common.BuildPendingOutdatedTimeout, fmt.Sprintf("%s.updated_at < ?", stmt.Schema.Table), "stuck_or_timeout_failure")
	w.drop("scheduled", common.BuildScheduledOutdatedTimeout, fmt.Sprintf("%s.scheduled_at IS NOT NULL AND scheduled_at < ?", stmt.Schema.Table), "stale_schedule")
	w.dropStuck("pending", common.BuildPendingStuckTimeout, fmt.Sprintf("%s.updated_at < ?", stmt.Schema.Table), "stuck_or_timeout_failure")

	return nil
}

func (w *StuckWorker) obtainLease() bool {

	return false
}

func (w *StuckWorker) dropBuild(dropType string, build *models.Build, state string, timeout time.Duration, reason string) error {
	//build.Drop(reason)
	return nil
}

func (w *StuckWorker) drop(state string, timeout time.Duration, condition string, reason string) error {
	builds, err := w.search(state, timeout, condition)
	if err != nil {
		return err
	}

	for _, build := range builds {
		if err := w.dropBuild("outdated", build, state, timeout, condition); err != nil {
			return err
		}
	}

	return nil
}

func (w *StuckWorker) dropStuck(state string, timeout time.Duration, condition string, reason string) error {
	builds, err := w.search(state, timeout, condition)
	if err != nil {
		return err
	}

	for _, build := range builds {
		// check if build.Stuck

		if err := w.dropBuild("stuck", build, state, timeout, condition); err != nil {
			return err
		}
	}

	return nil
}

func (w *StuckWorker) search(state string, timeout time.Duration, condition string) ([]*models.Build, error) {

	return nil, nil
}
