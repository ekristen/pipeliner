package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/ekristen/pipeliner/pkg/common"
	"gorm.io/gorm"
)

// Pipeline --
type Pipeline struct {
	ID         int64              `gorm:"primaryKey;autoIncrement:false" json:"id"`
	Duration   int                `json:"duration,omitempty"`
	FinishedAt *time.Time         `json:"finished_at,omitempty"`
	StartedAt  *time.Time         `json:"started_at,omitempty"`
	State      string             `json:"state"`
	Data       []byte             `json:"data"`
	WorkflowID int64              `json:"workflow_id"`
	Workflow   Workflow           `json:"workflow" gorm:"foreignKey:WorkflowID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt  *time.Time         `json:"created_at"`
	UpdatedAt  *time.Time         `json:"updated_at"`
	Builds     []Build            `json:"builds"`
	Stages     []PipelineStage    `json:"stages"`
	Variables  []PipelineVariable `json:"variables,omitempty"`
}

// BeforeCreate --
func (p *Pipeline) BeforeCreate(tx *gorm.DB) error {
	if p.ID == 0 {
		node := tx.Statement.Context.Value(common.ContextKeyNode).(*snowflake.Node)
		p.ID = node.Generate().Int64()
	}

	return nil
}

// AfterSave - sends channel notification
func (p *Pipeline) AfterSave(tx *gorm.DB) error {
	go func() {
		notify := tx.Statement.Context.Value(common.ContextKeyWebsocket).(chan []byte)
		change := VueWebsocketMessage{
			Action: "addPipeline",
			Data:   p,
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

// PipelineStage --
type PipelineStage struct {
	ID         int64      `gorm:"primaryKey;autoIncrement:false" json:"id"`
	Name       string     `json:"name" gorm:"size:16,uniqueIndex:stage"`
	Index      int64      `json:"index" gorm:"uniqueIndex:stage,sort:asc"`
	State      string     `json:"state"`
	PipelineID int64      `json:"pipeline_id" gorm:"uniqueIndex:stage"`
	CreatedAt  *time.Time `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
}

// BeforeCreate --
func (p *PipelineStage) BeforeCreate(tx *gorm.DB) error {
	if p.ID == 0 {
		node := tx.Statement.Context.Value(common.ContextKeyNode).(*snowflake.Node)
		p.ID = node.Generate().Int64()
	}

	return nil
}

// AfterSave - sends channel notification
func (p *PipelineStage) AfterSave(tx *gorm.DB) error {
	if p.ID == 0 {
		return nil
	}

	go func() {
		notify := tx.Statement.Context.Value(common.ContextKeyWebsocket).(chan []byte)
		change := VueWebsocketMessage{
			Action: "addStage",
			Data:   p,
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
