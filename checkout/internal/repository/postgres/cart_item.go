package postgres

import (
	"context"
	"route256/checkout/internal/converters"
	"route256/checkout/internal/model"
	"route256/checkout/internal/repository/schema"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
)

type CartItemRepository struct {
	queryEngineProvider QueryEngineProvider
}

func NewCartItemRepository(queryEngineProvider QueryEngineProvider) *CartItemRepository {
	return &CartItemRepository{
		queryEngineProvider: queryEngineProvider,
	}
}

const (
	cartItemTable = "cart_item"

	userIdColumn = "user_id"
	skuColumn    = "sku"
	countColumn  = "count"
)

func (r *CartItemRepository) AddItem(ctx context.Context, userId int64, item *model.CartItem) error {
	db := r.queryEngineProvider.GetQueryEngine(ctx)

	query, args, err := sq.
		Insert(cartItemTable).
		Columns(userIdColumn, skuColumn, countColumn).
		Values(userId, item.SKU, item.Count).
		Suffix("ON CONFLICT (user_id, sku) DO UPDATE SET count = EXCLUDED.count").
		PlaceholderFormat(sq.Dollar).
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

func (r *CartItemRepository) DeleteItem(ctx context.Context, userId int64, item *model.CartItem) error {
	db := r.queryEngineProvider.GetQueryEngine(ctx)

	query, args, err := sq.
		Delete(cartItemTable).
		Where(sq.Eq{userIdColumn: userId, skuColumn: item.SKU}).
		PlaceholderFormat(sq.Dollar).
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

func (r *CartItemRepository) GetItems(ctx context.Context, userId int64) ([]*model.CartItem, error) {
	db := r.queryEngineProvider.GetQueryEngine(ctx)

	query, args, err := sq.
		Select(userIdColumn, skuColumn, countColumn).
		From(cartItemTable).
		Where(sq.Eq{userIdColumn: userId}).
		PlaceholderFormat(sq.Dollar).
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

func (r *CartItemRepository) GetItem(ctx context.Context, userId int64, sku uint32) (*model.CartItem, error) {
	db := r.queryEngineProvider.GetQueryEngine(ctx)

	query, args, err := sq.
		Select(userIdColumn, skuColumn, countColumn).
		From(cartItemTable).
		Where(sq.Eq{userIdColumn: userId, skuColumn: sku}).
		PlaceholderFormat(sq.Dollar).
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

func (r *CartItemRepository) RemoveItems(ctx context.Context, userId int64, item *model.CartItem) error {
	db := r.queryEngineProvider.GetQueryEngine(ctx)

	query, args, err := sq.
		Update(cartItemTable).
		Set(countColumn, sq.ConcatExpr(countColumn, sq.Expr(" - ?", item.Count))).
		Where(sq.Eq{userIdColumn: userId, skuColumn: item.SKU}).
		PlaceholderFormat(sq.Dollar).
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
