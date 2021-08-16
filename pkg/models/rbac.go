package models

import (
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/ekristen/pipeliner/pkg/common"
	"gorm.io/gorm"
)

// Group --
type Group struct {
	ID        int64 `gorm:"primaryKey;autoIncrement:false" json:"id"`
	Name      string
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	Users     []*User    `gorm:"many2many:user_groups"`
}

// BeforeCreate --
func (g *Group) BeforeCreate(tx *gorm.DB) error {
	if g.ID == 0 {
		node := tx.Statement.Context.Value(common.ContextKeyNode).(*snowflake.Node)
		g.ID = node.Generate().Int64()
	}

	return nil
}

// User --
type User struct {
	ID        int64 `gorm:"primaryKey;autoIncrement:false" json:"id"`
	Name      string
	Username  string
	Password  string
	Groups    []*Group   `gorm:"many2many:user_groups"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

// BeforeCreate --
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == 0 {
		node := tx.Statement.Context.Value(common.ContextKeyNode).(*snowflake.Node)
		u.ID = node.Generate().Int64()
	}

	return nil
}
