package converters

import (
	"route256/checkout/internal/model"
	lomsapi "route256/checkout/pkg/loms"
)

func ToOrderItemListLomsApi(orderItems []*model.CartItem) []*lomsapi.OrderItem {
	if orderItems == nil {
		return nil
	}

	result := make([]*lomsapi.OrderItem, 0, len(orderItems))
	for _, i := range orderItems {
		item := ToOrderItemLomsApi(i)
		result = append(result, item)
	}

	return result
}

func ToOrderItemLomsApi(orderItem *model.CartItem) *lomsapi.OrderItem {
	if orderItem == nil {
		return nil
	}

	return &lomsapi.OrderItem{
		Sku:   orderItem.SKU,
		Count: orderItem.Count,
	}
}
