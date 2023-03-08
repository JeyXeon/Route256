package converter

import (
	"route256/loms/internal/model"
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
