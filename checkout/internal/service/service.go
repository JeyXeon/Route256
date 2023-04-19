package service

import (
	"context"
	"route256/checkout/internal/model"
)

type TransactionManager interface {
	RunRepeatableRead(ctx context.Context, f func(ctxTX context.Context) error) error
}

type Limiter interface {
	Wait(ctx context.Context) (err error)
}

type ItemRepository interface {
	AddItem(ctx context.Context, userId int64, item *model.CartItem) error
	DeleteItem(ctx context.Context, userId int64, item *model.CartItem) error
	GetItems(ctx context.Context, userId int64) ([]*model.CartItem, error)
	GetItem(ctx context.Context, userId int64, sku uint32) (*model.CartItem, error)
	RemoveItems(ctx context.Context, userId int64, item *model.CartItem) error
}

type LomsClient interface {
	Stocks(ctx context.Context, sku uint32) ([]*model.Stock, error)
	CreateOrder(ctx context.Context, user int64, items []*model.CartItem) (int64, error)
}

type ProductServiceClient interface {
	GetProduct(ctx context.Context, sku uint32) (*model.Product, error)
}

type ProductServiceCache interface {
	Get(sku uint32) (*model.Product, bool)
	Set(sku uint32, product model.Product)
}

type Service struct {
	transactionManager    TransactionManager
	itemsRepository       ItemRepository
	lomsClient            LomsClient
	productServiceClient  ProductServiceClient
	productServiceLimiter Limiter
	productServiceCache   ProductServiceCache
}

func New(
	transactionManager TransactionManager,
	itemRepository ItemRepository,
	stocksChecker LomsClient,
	productServiceClient ProductServiceClient,
	productServiceLimiter Limiter,
	productServiceCache ProductServiceCache,
) *Service {
	return &Service{
		transactionManager:    transactionManager,
		itemsRepository:       itemRepository,
		lomsClient:            stocksChecker,
		productServiceClient:  productServiceClient,
		productServiceLimiter: productServiceLimiter,
		productServiceCache:   productServiceCache,
	}
}
