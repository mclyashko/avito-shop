package model

import (
	"github.com/google/uuid"
)

type CoinTransfer struct {
	ID         uuid.UUID
	SenderID   string
	ReceiverID string
	Amount     int64
}
