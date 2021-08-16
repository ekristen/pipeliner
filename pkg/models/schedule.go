package models

import (
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/ekristen/pipeliner/pkg/common"
	"github.com/gorhill/cronexpr"
	"gorm.io/gorm"
)

// Schedule --
type Schedule struct {
	ID         int64      `gorm:"primaryKey;autoIncrement:false" json:"id"`
	Name       string     `json:"name"`
	Cron       string     `json:"cron"`
	Last       *time.Time `json:"last"`
	Next       *time.Time `json:"next"`
	WorkflowID int64      `json:"workflow_id"`
	Workflow   Workflow   `gorm:"foreignKey:WorkflowID" json:"workflow"`
	CreatedAt  *time.Time `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
}

// BeforeCreate will generate a Snowflake ID
func (b *Schedule) BeforeCreate(tx *gorm.DB) error {
	if b.ID == 0 {
		node := tx.Statement.Context.Value(common.ContextKeyNode).(*snowflake.Node)
		b.ID = node.Generate().Int64()
	}

	e, err := cronexpr.Parse(b.Cron)
	if err != nil {
		return err
	}

	if b.Next == nil {
		nextTime := e.Next(time.Now())
		b.Next = &nextTime
	}

	return nil
}

// BeforeUpdate --
func (b *Schedule) BeforeUpdate(tx *gorm.DB) error {
	e, err := cronexpr.Parse(b.Cron)
	if err != nil {
		return err
	}

	if tx.Statement.Changed("Last") {
		nextTime := e.Next(time.Now())
		b.Next = &nextTime
	}

	return nil
}
