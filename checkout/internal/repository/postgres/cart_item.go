package postgres

import (
	"context"
	"route256/checkout/internal/converters"
	"route256/checkout/internal/model"
	"route256/checkout/internal/repository/schema"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
)

type CartItemRepository interface {
	AddItem(ctx context.Context, userId int64, item *model.CartItem) error
	DeleteItem(ctx context.Context, userId int64, item *model.CartItem) error
	GetItems(ctx context.Context, userId int64) ([]*model.CartItem, error)
	GetItem(ctx context.Context, userId int64, sku uint32) (*model.CartItem, error)
	RemoveItems(ctx context.Context, userId int64, item *model.CartItem) error
}

type cartItemRepository struct {
	queryEngineProvider QueryEngineProvider
}

func NewCartItemRepository(queryEngineProvider QueryEngineProvider) CartItemRepository {
	return &cartItemRepository{
		queryEngineProvider: queryEngineProvider,
	}
}

const (
	cartItemTable = "cart_item"

	userIdColumn = "user_id"
	skuColumn    = "sku"
	countColumn  = "count"
)

func (r *cartItemRepository) AddItem(ctx context.Context, userId int64, item *model.CartItem) error {
	db := r.queryEngineProvider.GetQueryEngine(ctx)

	query, args, err := queryBuilder().
		Insert(cartItemTable).
		Columns(userIdColumn, skuColumn, countColumn).
		Values(userId, item.SKU, item.Count).
		Suffix("ON CONFLICT (user_id, sku) DO UPDATE SET count = EXCLUDED.count").
		ToSql()
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *cartItemRepository) DeleteItem(ctx context.Context, userId int64, item *model.CartItem) error {
	db := r.queryEngineProvider.GetQueryEngine(ctx)

	query, args, err := queryBuilder().
		Delete(cartItemTable).
		Where(sq.Eq{userIdColumn: userId, skuColumn: item.SKU}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *cartItemRepository) GetItems(ctx context.Context, userId int64) ([]*model.CartItem, error) {
	db := r.queryEngineProvider.GetQueryEngine(ctx)

	query, args, err := queryBuilder().
		Select(userIdColumn, skuColumn, countColumn).
		From(cartItemTable).
		Where(sq.Eq{userIdColumn: userId}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var cartItems []*schema.CartItem
	if err := pgxscan.Select(ctx, db, &cartItems, query, args...); err != nil {
		return nil, err
	}

	return converters.SchemaToOrderItemsModel(cartItems), nil
}

func (r *cartItemRepository) GetItem(ctx context.Context, userId int64, sku uint32) (*model.CartItem, error) {
	db := r.queryEngineProvider.GetQueryEngine(ctx)

	query, args, err := queryBuilder().
		Select(userIdColumn, skuColumn, countColumn).
		From(cartItemTable).
		Where(sq.Eq{userIdColumn: userId, skuColumn: sku}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var cartItem schema.CartItem
	if err := pgxscan.Get(ctx, db, &cartItem, query, args...); err != nil {
		return nil, err
	}

	return converters.SchemaToOrderItemModel(&cartItem), nil
}

func (r *cartItemRepository) RemoveItems(ctx context.Context, userId int64, item *model.CartItem) error {
	db := r.queryEngineProvider.GetQueryEngine(ctx)

	query, args, err := queryBuilder().
		Update(cartItemTable).
		Set(countColumn, sq.ConcatExpr(countColumn, sq.Expr(" - ?", item.Count))).
		Where(sq.Eq{userIdColumn: userId, skuColumn: item.SKU}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
