package converters

import (
	"route256/checkout/internal/model"
	lomsapi "route256/checkout/pkg/loms"
)

func ToOrderItemListLomsApi(orderItems []*model.OrderItem) []*lomsapi.Item {
	if orderItems == nil {
		return nil
	}

	result := make([]*lomsapi.Item, 0, len(orderItems))
	for _, i := range orderItems {
		item := ToOrderItemLomsApi(i)
		result = append(result, item)
	}

	return result
}

func ToOrderItemLomsApi(orderItem *model.OrderItem) *lomsapi.Item {
	if orderItem == nil {
		return nil
	}

	return &lomsapi.Item{
		Sku:   orderItem.SKU,
		Count: orderItem.Count,
	}
}
