package service

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/mclyashko/avito-shop/internal/db"
)

type BuyService interface {
	BuyItem(ctx context.Context, username string, itemName string) error
}

type BuyServiceImpl struct {
	Service
	UserAccessor     db.UserAccessor
	ItemAccessor     db.ItemAccessor
	UserItemAccessor db.UserItemAccessor
}

func (s *BuyServiceImpl) BuyItem(ctx context.Context, username string, itemName string) error {
	return s.RunWithTx(ctx, func(tx pgx.Tx) error {
		user, err := s.UserAccessor.GetUserByLoginTx(ctx, tx, username)
		if err != nil {
			return err
		}

		item, err := s.ItemAccessor.GetItemByNameTx(ctx, tx, itemName)
		if err != nil {
			return err
		}

		if user.Balance < item.Price {
			return ErrInsufficientFunds
		}

		err = s.UserAccessor.UpdateUserBalanceTx(ctx, tx, username, -item.Price)
		if err != nil {
			return err
		}

		err = s.UserItemAccessor.InsertUserItemTx(ctx, tx, username, itemName)
		if err != nil {
			return err
		}

		return nil
	})
}
