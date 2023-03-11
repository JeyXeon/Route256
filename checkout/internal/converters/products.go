package converters

import (
	"route256/checkout/internal/model"
	desc "route256/checkout/pkg/checkout"
)

func ModelToProductListDesc(products []*model.Product) []*desc.Product {
	if products == nil {
		return nil
	}

	result := make([]*desc.Product, 0, len(products))
	for _, p := range products {
		product := ModelToProductDesc(p)
		result = append(result, product)
	}

	return result
}

func ModelToProductDesc(product *model.Product) *desc.Product {
	if product == nil {
		return nil
	}

	return &desc.Product{
		Sku:   product.SKU,
		Count: product.Count,
		Name:  product.Name,
		Price: product.Price,
	}
}
