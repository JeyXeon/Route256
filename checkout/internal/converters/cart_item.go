package converters

import (
	"route256/checkout/internal/model"
	"route256/checkout/internal/repository/schema"
	lomsapi "route256/checkout/pkg/loms"
)

func ModelToOrderItemListLomsApi(orderItems []*model.CartItem) []*lomsapi.OrderItem {
	if orderItems == nil {
		return nil
	}

	result := make([]*lomsapi.OrderItem, 0, len(orderItems))
	for _, i := range orderItems {
		item := ModelToOrderItemLomsApi(i)
		result = append(result, item)
	}

	return result
}

func ModelToOrderItemLomsApi(orderItem *model.CartItem) *lomsapi.OrderItem {
	if orderItem == nil {
		return nil
	}

	return &lomsapi.OrderItem{
		Sku:   orderItem.SKU,
		Count: orderItem.Count,
	}
}

func SchemaToOrderItemModel(cartItem *schema.CartItem) *model.CartItem {
	if cartItem == nil {
		return nil
	}

	return &model.CartItem{
		SKU:   cartItem.Sku,
		Count: cartItem.Count,
	}
}

func SchemaToOrderItemsModel(cartItems []*schema.CartItem) []*model.CartItem {
	if cartItems == nil {
		return nil
	}

	result := make([]*model.CartItem, 0, len(cartItems))
	for _, cartItem := range cartItems {
		result = append(result, SchemaToOrderItemModel(cartItem))
	}

	return result
}
