package converter

import (
	"route256/loms/internal/model"
	desc "route256/loms/pkg/loms"
)

func ToOrderItemsListModel(orderItems []*desc.OrderItem) []*model.OrderItem {
	if orderItems == nil {
		return nil
	}

	result := make([]*model.OrderItem, 0, len(orderItems))
	for _, orderItem := range orderItems {
		result = append(result, ToOrderItemModel(orderItem))
	}
	return result
}

func ToOrderItemModel(orderItem *desc.OrderItem) *model.OrderItem {
	if orderItem == nil {
		return nil
	}

	return &model.OrderItem{
		Sku:   orderItem.Sku,
		Count: orderItem.Count,
	}
}

func ToOrderItemsListLomsApi(orderItems []*model.OrderItem) []*desc.OrderItem {
	if orderItems == nil {
		return nil
	}

	result := make([]*desc.OrderItem, 0, len(orderItems))
	for _, orderItem := range orderItems {
		result = append(result, ToOrderItemLomsApi(orderItem))
	}
	return result
}

func ToOrderItemLomsApi(orderItem *model.OrderItem) *desc.OrderItem {
	if orderItem == nil {
		return nil
	}

	return &desc.OrderItem{
		Sku:   orderItem.Sku,
		Count: orderItem.Count,
	}
}
