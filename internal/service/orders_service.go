package service

import (
	"context"
	"strings"

	"github.com/klyakssa/go-diplom.git/internal/domain/orders"
	"github.com/klyakssa/go-diplom.git/pkg/luhn"
)

type OrdersService struct {
	repo orders.Repository
}

func NewOrdersService(repo orders.Repository) *OrdersService {
	return &OrdersService{
		repo: repo,
	}
}

func (o *OrdersService) CreateOrder(ctx context.Context, number, userID string) error {

	if !luhn.Valid(number) {
		return orders.ErrIncorrectOrderNumberFormat
	}

	order, err := o.repo.GetOrderByNumber(ctx, number)
	if err == nil {
		if strings.Compare(order.UserID, userID) == 0 {
			return orders.ErrOrderAlreadyAddedByThisUser
		}
		return orders.ErrOrderAlreadyAddedByOtherUser
	}

	err = o.repo.CreateOrder(ctx, &orders.Order{
		Number: number,
		UserID: userID,
	})
	if err != nil {
		return err
	}

	return nil
}

func (o *OrdersService) GetOrders(ctx context.Context, userID string) ([]orders.Order, error) {

	orders, err := o.repo.GetOrdersByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return orders, nil
}
