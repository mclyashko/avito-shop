package model

import (
	"github.com/google/uuid"
)

type CoinTransfer struct {
	ID            uuid.UUID
	SenderLogin   string
	ReceiverLogin string
	Amount        int64
}
