package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/ekristen/pipeliner/pkg/common"
	"github.com/ekristen/pipeliner/pkg/utils"
	"github.com/jinzhu/copier"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Build --
type Build struct {
	ID            int64            `gorm:"primaryKey;autoIncrement:false" json:"id"`
	State         string           `json:"state"`
	AllowFailure  bool             `json:"allow_failure"`
	Dependencies  string           `json:"dependencies,omitempty"`
	FinishedAt    *time.Time       `json:"finished_at,omitempty"`
	Name          string           `json:"name"`
	Job           string           `json:"job" gorm:"uniqueIndex:pipeline_stage_instance,priority:4;size:256"`
	Options       datatypes.JSON   `json:"-"` // map[string]string
	QueuedAt      *time.Time       `json:"queued_at,omitempty"`
	Stage         string           `json:"stage"`
	StageIdx      int64            `json:"stage_idx" gorm:"uniqueIndex:pipeline_stage_instance,priority:2"`
	StartedAt     *time.Time       `json:"started_at,omitempty"`
	Token         string           `json:"token"`
	When          string           `json:"when"`
	Data          []byte           `json:"data,omitempty"`
	InstanceIdx   int64            `json:"instance_idx" gorm:"default:0;uniqueIndex:pipeline_stage_instance,priority:3"`
	Duration      *int             `json:"duration"`
	Parallel      int64            `json:"parallel" gorm:"default:0;uniqueIndex:pipeline_stage_instance,priority:5;size:256"`
	Timeout       *int             `json:"timeout" gorm:"default:3600"`
	FailureReason *string          `json:"failure_reason,omitempty"`
	Retried       bool             `json:"retried" gorm:"default:0"`
	PipelineID    int64            `json:"pipeline_id" gorm:"index"`
	Pipeline      Pipeline         `gorm:"foreignKey:PipelineID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;column:pipeline_id;uniqueIndex:pipeline_stage_instance,priority:1" json:"-"`
	RunnerID      *int64           `json:"runner_id"`
	Runner        *Runner          `gorm:"foreignKey:RunnerID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	Tags          []BuildTag       `json:"tags,omitempty"`
	Variables     []*BuildVariable `json:"variables,omitempty"`
	Artifacts     []*Artifact      `json:"artifacts,omitempty"`
	CreatedAt     *time.Time       `json:"created_at"`
	UpdatedAt     *time.Time       `json:"updated_at"`
}

// BeforeCreate --
func (b *Build) BeforeCreate(tx *gorm.DB) error {
	if b.ID == 0 {
		node := tx.Statement.Context.Value(common.ContextKeyNode).(*snowflake.Node)
		b.ID = node.Generate().Int64()
	}

	return nil
}

// AfterSave - sends channel notification
func (b *Build) AfterSave(tx *gorm.DB) error {
	if b.ID == 0 {
		return nil
	}

	go func() {
		notify := tx.Statement.Context.Value(common.ContextKeyWebsocket).(chan []byte)
		change := VueWebsocketMessage{
			Action: "addBuild",
			Data:   b,
		}
		data, err := json.Marshal(change)
		if err != nil {
			fmt.Println(err)
			return
		}
		notify <- data
	}()

	return nil
}

// Clone the build to a new build minus important information
func (b *Build) Clone() (*Build, error) {
	newBuild := &Build{}

	if err := copier.Copy(newBuild, b); err != nil {
		return nil, err
	}

	var tags []BuildTag
	for _, t := range b.Tags {
		tags = append(tags, BuildTag{
			Tag: t.Tag,
		})
	}

	newBuild.ID = 0
	newBuild.FinishedAt = nil
	newBuild.QueuedAt = nil
	newBuild.StartedAt = nil
	newBuild.Duration = nil
	newBuild.CreatedAt = nil
	newBuild.UpdatedAt = nil
	newBuild.RunnerID = nil
	newBuild.Runner = nil
	newBuild.InstanceIdx = newBuild.InstanceIdx + 1
	newBuild.Token = utils.RandomString(16)
	newBuild.State = "created"
	newBuild.Tags = tags
	newBuild.Variables = nil

	return newBuild, nil
}

// BuildTag --
type BuildTag struct {
	Tag     string `json:"tag" gorm:"size:64;uniqueIndex:build_tag"`
	BuildID int64  `json:"-" gorm:"uniqueIndex:build_tag"`
	Build   Build  `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`

	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// AfterSave - sends channel notification
func (b *BuildTag) AfterSave(tx *gorm.DB) error {
	go func() {
		notify := tx.Statement.Context.Value(common.ContextKeyWebsocket).(chan []byte)
		change := VueWebsocketMessage{
			Action: "addBuildTag",
			Data:   b,
		}
		data, err := json.Marshal(change)
		if err != nil {
			fmt.Println(err)
			return
		}
		notify <- data
	}()

	return nil
}
