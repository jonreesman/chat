package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	User      Client
	UserID    uuid.UUID `gorm:"constraint:OnDelete:CASCADE;"`
	Timestamp int64
	Content   string
	Room      Room      `gorm:"constraint:OnDelete:CASCADE;"`
	RoomID    uuid.UUID `gorm:"constraint:OnDelete:CASCADE;"`
}
