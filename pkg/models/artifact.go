package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/ekristen/pipeliner/pkg/common"
	"gorm.io/gorm"
)

// Artifact --
type Artifact struct {
	ID         int64      `gorm:"primaryKey;autoIncrement:false" json:"id"`
	Type       string     `json:"type"`
	Format     string     `json:"format"`
	File       string     `json:"file"`
	Size       int64      `json:"size"`
	ExpiresAt  *time.Time `json:"expires_at"`
	BuildID    int64      `json:"build_id"`
	Build      Build      `gorm:"foreignKey:BuildID" json:"-"`
	PipelineID int64      `json:"pipeline_id"`
	Pipeline   Pipeline   `gorm:"foreignKey:PipelineID" json:"-"`
	CreatedAt  *time.Time `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
}

// BeforeCreate --
func (a *Artifact) BeforeCreate(tx *gorm.DB) error {
	if a.ID == 0 {
		node := tx.Statement.Context.Value(common.ContextKeyNode).(*snowflake.Node)
		a.ID = node.Generate().Int64()
	}

	return nil
}

// AfterSave - sends channel notification
func (a *Artifact) AfterSave(tx *gorm.DB) error {
	go func() {
		notify := tx.Statement.Context.Value(common.ContextKeyWebsocket).(chan []byte)
		change := VueWebsocketMessage{
			Action: "addArtifact",
			Data:   a,
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
