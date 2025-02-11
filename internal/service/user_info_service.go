package service

import (
	"context"

	"github.com/mclyashko/avito-shop/internal/db"
	"github.com/mclyashko/avito-shop/internal/model"
)

type UserInfoService interface {
	GetUserInfo(ctx context.Context, username string) (balance *int64, userItems []model.UserItem, recievedTransfers []model.CoinTransfer, sentTransfers []model.CoinTransfer, err error)
}

type UserInfoServiceImp struct {
	Service
	UserAccessor         db.UserAccessor
	UserItemAccessor     db.UserItemAccessor
	CoinTransferAccessor db.CoinTransferAccessor
}

func (s *UserInfoServiceImp) GetUserInfo(ctx context.Context, username string) (balance *int64, userItems []model.UserItem, recievedTransfers []model.CoinTransfer, sentTransfers []model.CoinTransfer, err error) {
	user, err := s.UserAccessor.GetUserByLogin(ctx, username)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	balance = &user.Balance

	userItems, err = s.UserItemAccessor.GetUserItemsByUsername(ctx, username)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	recievedTransfers, sentTransfers, err = s.CoinTransferAccessor.GetUserTransactionHistory(ctx, username)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return balance, userItems, recievedTransfers, sentTransfers, nil
}
