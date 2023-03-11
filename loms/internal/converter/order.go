package converter

import (
	"route256/loms/internal/model"
	"route256/loms/internal/repository/schema"

	"github.com/pkg/errors"
)

var (
	ErrUnknownOrderStatus = errors.New("unknown error status")
)

func ToOrderStatusModel(orderStatus schema.OrderStatus) (model.OrderStatus, error) {
	switch orderStatus {
	case schema.New:
		return model.New, nil
	case schema.AwaitingPayment:
		return model.AwaitingPayment, nil
	case schema.Cancelled:
		return model.Cancelled, nil
	case schema.Failed:
		return model.Failed, nil
	case schema.Payed:
		return model.Payed, nil
	default:
		return "", ErrUnknownOrderStatus
	}
}

func ToOrderModel(order *schema.Order) (*model.Order, error) {
	if order == nil {
		return nil, nil
	}

	orderStatus, err := ToOrderStatusModel(order.Status)
	if err != nil {
		return nil, err
	}

	return &model.Order{
		ID:     order.Id,
		User:   order.UserId,
		Status: orderStatus,
	}, nil
}

func ToOrderStatusSchema(orderStatus model.OrderStatus) (schema.OrderStatus, error) {
	switch orderStatus {
	case model.New:
		return schema.New, nil
	case model.AwaitingPayment:
		return schema.AwaitingPayment, nil
	case model.Cancelled:
		return schema.Cancelled, nil
	case model.Failed:
		return schema.Failed, nil
	case model.Payed:
		return schema.Payed, nil
	default:
		return "", ErrUnknownOrderStatus
	}
}
