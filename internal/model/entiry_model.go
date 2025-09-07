package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Entity struct
type Entity struct {
	gorm.Model
	ID          uuid.UUID `gorm:"primaryKey;uniqueIndex;not null;type:uuid;"`
	UserID      uuid.UUID `gorm:"type:uuid;not null;" json:"user_id"`
	Description *string   `gorm:"type:text;" json:"description"`
}

func (user *Entity) BeforeCreate(tx *gorm.DB) error {
	user.ID = uuid.New()

	return nil
}

// TableName sets the table name for the Entity model
func (Entity) TableName() string {
	return "entitys"
}
