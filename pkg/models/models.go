package models

import (
	"time"

	"github.com/bwmarrin/snowflake"
	"gorm.io/gorm"

	"github.com/ekristen/pipeliner/pkg/common"
	"github.com/ekristen/pipeliner/pkg/utils"
)

// VueWebsocketMessage --
type VueWebsocketMessage struct {
	Namespace string      `json:"namespace,omitempty"`
	Mutation  string      `json:"mutation,omitempty"`
	Action    string      `json:"action,omitempty"`
	Data      interface{} `json:"data"`
}

/*
// Base --
type Base struct {
	ID        int64      `gorm:"primaryKey;autoIncrement:false" json:"id"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	//DeletedAt *gorm.DeletedAt `json:"-"`
}

// BeforeCreate will generate a Snowflake ID
func (b *Base) BeforeCreate(tx *gorm.DB) error {
	if b.ID == 0 {
		node := tx.Statement.Context.Value(common.ContextKeyNode).(*snowflake.Node)
		b.ID = node.Generate().Int64()
	}

	return nil
}
*/
// AfterSave - sends channel notification
/*
func (b *Base) AfterSave(tx *gorm.DB) error {
	go func() {
		notify := tx.Statement.Context.Value(common.ContextKeyWebsocket).(chan []byte)
		change := VueWebsocketMessage{
			Data: b,
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
*/

// RegisterToken --
type RegisterToken struct {
	ID        int64      `json:"id" gorm:"primaryKey;autoIncrement:false"`
	Token     string     `json:"token"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

// BeforeCreate will generate a Snowflake ID
func (r *RegisterToken) BeforeCreate(tx *gorm.DB) error {
	if r.ID == 0 {
		node := tx.Statement.Context.Value(common.ContextKeyNode).(*snowflake.Node)
		r.ID = node.Generate().Int64()
	}

	r.Token = utils.RandomString(16)

	return nil
}
