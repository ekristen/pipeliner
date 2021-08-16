package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/ekristen/pipeliner/pkg/common"
	"gorm.io/gorm"
)

// Workflow --
type Workflow struct {
	ID        int64      `json:"id" gorm:"primaryKey;autoIncrement:false"`
	Name      string     `json:"name" gorm:"uniqueIndex;size:256"`
	Data      string     `json:"data,omitempty"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

// BeforeCreate Hook
func (w *Workflow) BeforeCreate(tx *gorm.DB) error {
	if w.ID == 0 {
		node := tx.Statement.Context.Value(common.ContextKeyNode).(*snowflake.Node)
		w.ID = node.Generate().Int64()
	}

	return nil
}

// AfterSave - sends channel notification
func (w *Workflow) AfterSave(tx *gorm.DB) error {
	go func() {
		notify := tx.Statement.Context.Value(common.ContextKeyWebsocket).(chan []byte)
		change := VueWebsocketMessage{
			Action: "addWorkflow",
			Data:   w,
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
