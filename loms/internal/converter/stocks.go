package converter

import (
	"route256/loms/internal/model"
	"route256/loms/internal/repository/schema"
	desc "route256/loms/pkg/loms"
)

func ToStockListLomsApi(stocks []*model.Stock) []*desc.Stock {
	if stocks == nil {
		return nil
	}

	result := make([]*desc.Stock, 0, len(stocks))
	for _, stock := range stocks {
		result = append(result, ToStockLomsApi(stock))
	}
	return result
}

func ToStockLomsApi(stock *model.Stock) *desc.Stock {
	if stock == nil {
		return nil
	}

	return &desc.Stock{
		WarehouseId: stock.WareHouseId,
		Count:       stock.Count,
	}
}

func ToStockListModel(stocks []*schema.Stock) []*model.Stock {
	if stocks == nil {
		return nil
	}

	result := make([]*model.Stock, 0, len(stocks))
	for _, stock := range stocks {
		result = append(result, ToStockModel(stock))
	}
	return result
}

func ToStockModel(stock *schema.Stock) *model.Stock {
	if stock == nil {
		return nil
	}

	return &model.Stock{
		Sku:         stock.Sku,
		WareHouseId: stock.WarehouseId,
		Count:       stock.Count,
	}
}
