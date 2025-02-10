package model

import (
	"github.com/google/uuid"
)

type UserItem struct {
	ID       uuid.UUID
	UserID   string
	ItemName string
	Quantity int // Количество предметов
}
