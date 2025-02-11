package model

import (
	"github.com/google/uuid"
)

type UserItem struct {
	ID        uuid.UUID
	UserLogin string
	ItemName  string
	Quantity  int
}
