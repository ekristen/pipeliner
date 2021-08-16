package models

import (
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/ekristen/pipeliner/pkg/common"
	"gorm.io/gorm"
)

// Trace --
type Trace struct {
	ID        int64       `gorm:"primaryKey;autoIncrement:false" json:"id"`
	BuildID   int64       `gorm:"unique" json:"build_id"`
	Build     Build       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"build"`
	CreatedAt *time.Time  `json:"created_at"`
	UpdatedAt *time.Time  `json:"updated_at"`
	Parts     []TracePart `json:"parts,omitempty"`
}

// BeforeCreate will generate a Snowflake ID
func (b *Trace) BeforeCreate(tx *gorm.DB) error {
	if b.ID == 0 {
		node := tx.Statement.Context.Value(common.ContextKeyNode).(*snowflake.Node)
		b.ID = node.Generate().Int64()
	}

	return nil
}

// TracePart --
type TracePart struct {
	ID        int64      `gorm:"primaryKey;autoIncrement:false" json:"id"`
	Start     int        `json:"start"`
	End       int        `json:"end"`
	Size      int        `json:"size"`
	Data      []byte     `json:"data"`
	TraceID   int64      `json:"trace_id"`
	Trace     Trace      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"  json:"trace"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

// BeforeCreate will generate a Snowflake ID
func (b *TracePart) BeforeCreate(tx *gorm.DB) error {
	if b.ID == 0 {
		node := tx.Statement.Context.Value(common.ContextKeyNode).(*snowflake.Node)
		b.ID = node.Generate().Int64()
	}

	return nil
}
