package model

import (
	"github.com/google/uuid"
)

type CoinTransfer struct {
	ID         uuid.UUID
	SenderID   string // Может быть ""
	ReceiverID string
	Amount     int64
}
