package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/ekristen/pipeliner/pkg/common"
	"gorm.io/gorm"
)

// Runner --
type Runner struct {
	ID           int64       `gorm:"primaryKey;autoIncrement:false" json:"id"`
	Active       bool        `json:"active"`
	Architecture string      `json:"architecture"`
	ContactedAt  *time.Time  `json:"contacted_at"`
	Description  string      `json:"description"`
	Global       bool        `json:"global"`
	Locked       bool        `json:"locked"`
	Name         string      `json:"name"`
	Platform     string      `json:"platform"`
	Token        string      `json:"token"`
	Tags         []RunnerTag `json:"tags,omitempty"` // []string
	Version      string      `json:"version"`
	RunUntagged  bool        `json:"run_untagged"`
	CreatedAt    *time.Time  `json:"created_at"`
	UpdatedAt    *time.Time  `json:"updated_at"`
}

// BeforeCreate --
func (r *Runner) BeforeCreate(tx *gorm.DB) error {
	if r.ID == 0 {
		node := tx.Statement.Context.Value(common.ContextKeyNode).(*snowflake.Node)
		r.ID = node.Generate().Int64()
	}

	return nil
}

// AfterSave - sends channel notification
func (r *Runner) AfterSave(tx *gorm.DB) error {
	if r.ID == 0 {
		return nil
	}

	go func() {
		notify := tx.Statement.Context.Value(common.ContextKeyWebsocket).(chan []byte)
		change := VueWebsocketMessage{
			Action: "addRunner",
			Data:   r,
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

// RunnerTag --
type RunnerTag struct {
	Tag      string `json:"tag" gorm:"primaryKey;size:64"`
	RunnerID int64  `json:"runner_id" gorm:"primaryKey"`
	Runner   Runner `json:"runner" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

// AfterSave - sends channel notification
func (e *RunnerTag) AfterSave(tx *gorm.DB) error {
	go func() {
		notify := tx.Statement.Context.Value(common.ContextKeyWebsocket).(chan []byte)
		change := VueWebsocketMessage{
			Action: "addRunnerTag",
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
