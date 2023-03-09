package loms

import (
	"context"
	"route256/checkout/internal/converters"
	"route256/checkout/internal/model"
	lomsapi "route256/checkout/pkg/loms"
)

func (c *client) CreateOrder(ctx context.Context, user int64, items []*model.CartItem) (int64, error) {
	orderItems := converters.ToOrderItemListLomsApi(items)
	req := &lomsapi.CreateOrderRequest{
		User:  user,
		Items: orderItems,
	}

	res, err := c.lomsClient.CreateOrder(ctx, req)
	if err != nil {
		return 0, err
	}

	return res.GetOrderID(), nil

}
