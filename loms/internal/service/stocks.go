package service

import (
	"context"
	"route256/loms/internal/model"

	"github.com/pkg/errors"
)

var (
	ErrGettingStocks = errors.New("getting stocks failed")
)

func (s *Service) GetStocks(ctx context.Context, sku uint32) ([]*model.Stock, error) {
	stocks, err := s.stocksRepository.GetStocks(ctx, []uint32{sku})
	if err != nil {
		return nil, ErrGettingStocks
	}

	return stocks, nil
}
