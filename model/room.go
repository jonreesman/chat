package model

import "github.com/google/uuid"

type Room struct {
	Name string
	ID   uuid.UUID
}
