package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ekristen/pipeliner/pkg/common"
	"gorm.io/gorm"
)

// Variable --
type Variable struct {
	Name     string `gorm:"primaryKey;autoIncrement:false" json:"name"`
	Value    string `json:"value"`
	Masked   bool   `json:"masked" gorm:"default:0"`
	File     bool   `json:"file" gorm:"default:0"`
	Internal bool   `json:"internal" gorm:"default:0"`
	Public   bool   `json:"public" gorm:"default:1"`

	// TODO: switch to pointers?
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// GlobalVariable define variables that are merged into all jobs
type GlobalVariable struct {
	*Variable
}

// AfterSave - sends channel notification
func (e *GlobalVariable) AfterSave(tx *gorm.DB) error {
	go func() {
		notify := tx.Statement.Context.Value(common.ContextKeyWebsocket).(chan []byte)
		change := VueWebsocketMessage{
			Action: "addGlobalVariable",
			Data:   e,
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

// WorkflowVariable define variables that are merged into all pipelines for a workflow
type WorkflowVariable struct {
	*Variable
	WorkflowID int64    `gorm:"primaryKey" json:"workflow_id"`
	Workflow   Workflow `gorm:"foreignKey:WorkflowID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}

// PipelineVariable define variables that are merged into all pipeline builds
type PipelineVariable struct {
	*Variable
	PipelineID int64    `gorm:"primaryKey" json:"pipeline_id"`
	Pipeline   Pipeline `gorm:"foreignKey:PipelineID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}

// BuildVariable defines variables that are merged into specific builds and dependent builds
type BuildVariable struct {
	*Variable
	BuildID int64 `gorm:"primaryKey" json:"build_id"`
	Build   Build `gorm:"foreignKey:BuildID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}

// ScheduleVariable defines variables that are merged into pipelines when run by a schedule
type ScheduleVariable struct {
	*Variable
	BuildID int64 `gorm:"primaryKey" json:"build_id"`
	Build   Build `gorm:"foreignKey:BuildID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}
