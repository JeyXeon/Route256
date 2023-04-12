package checkout

import (
	"context"
	checkout "route256/checkout/internal/api"
	"route256/checkout/internal/clients/productservice"
	productServiceClientMocks "route256/checkout/internal/clients/productservice/mocks"
	"route256/checkout/internal/converters"
	"route256/checkout/internal/model"
	"route256/checkout/internal/repository/postgres"
	repositoryMocks "route256/checkout/internal/repository/postgres/mocks"
	"route256/checkout/internal/service"
	serviceMocks "route256/checkout/internal/service/mocks"
	desc "route256/checkout/pkg/checkout"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestListCart(t *testing.T) {
	type cartItemRepositoryMockFunc func(mc *minimock.Controller) postgres.CartItemRepository
	type limiterMockFunc func(mc *minimock.Controller) service.Limiter
	type productServiceClientMockFunc func(mc *minimock.Controller) productservice.Client

	type args struct {
		ctx context.Context
		req *desc.ListCartRequest
	}

	var (
		mc  = minimock.NewController(t)
		ctx = context.Background()

		cartItemRepositoryErr   = errors.New("cart item repository error")
		limiterErr              = errors.New("limiter error")
		productServiceClientErr = errors.New("product service client error")

		cartItemRepositoryRes = []*model.CartItem{
			{
				SKU:   2,
				Count: 2,
			},
			{
				SKU:   5,
				Count: 1,
			},
		}

		productServiceClientRes = []*model.Product{
			{
				SKU:   2,
				Count: 1,
				Name:  "test product 1",
				Price: 1000,
			},
			{
				SKU:   5,
				Count: 1,
				Name:  "test product 2",
				Price: 2000,
			},
		}

		cartRes = &model.Cart{
			Items: []*model.Product{
				{
					SKU:   2,
					Count: 2,
					Name:  "test product 1",
					Price: 1000,
				},
				{
					SKU:   5,
					Count: 1,
					Name:  "test product 2",
					Price: 2000,
				},
			},
			TotalPrice: 4000,
		}
	)
	t.Cleanup(mc.Finish)

	type purchaseTestCase struct {
		name string
		args args

		cartItemRepositoryMock   cartItemRepositoryMockFunc
		limiterMock              limiterMockFunc
		productServiceClientMock productServiceClientMockFunc

		want *desc.ListCartResponse
		err  error
	}

	testCases := []purchaseTestCase{
		{
			name: "positive case",
			args: args{
				ctx: ctx,
				req: &desc.ListCartRequest{User: 1},
			},

			cartItemRepositoryMock: func(mc *minimock.Controller) postgres.CartItemRepository {
				mock := repositoryMocks.NewCartItemRepositoryMock(mc)
				mock.GetItemsMock.Expect(ctx, 1).Return(cartItemRepositoryRes, nil)
				return mock
			},
			limiterMock: func(mc *minimock.Controller) service.Limiter {
				mock := serviceMocks.NewLimiterMock(mc)
				mock.WaitMock.Expect(ctx).Return(nil)
				return mock
			},
			productServiceClientMock: func(mc *minimock.Controller) productservice.Client {
				mock := productServiceClientMocks.NewClientMock(mc)
				for _, res := range productServiceClientRes {
					mock.GetProductMock.When(ctx, res.SKU).Then(res, nil)
				}
				return mock
			},

			want: &desc.ListCartResponse{
				Items:      converters.ModelToProductListDesc(cartRes.Items),
				TotalPrice: cartRes.TotalPrice,
			},
			err: nil,
		},
		{
			name: "negative case - empty user",
			args: args{
				ctx: ctx,
				req: &desc.ListCartRequest{User: 0},
			},

			cartItemRepositoryMock: func(mc *minimock.Controller) postgres.CartItemRepository {
				return repositoryMocks.NewCartItemRepositoryMock(mc)

			},
			limiterMock: func(mc *minimock.Controller) service.Limiter {
				return serviceMocks.NewLimiterMock(mc)
			},
			productServiceClientMock: func(mc *minimock.Controller) productservice.Client {
				return productServiceClientMocks.NewClientMock(mc)
			},

			want: nil,
			err:  checkout.ErrListCartEmptyUser,
		},
		{
			name: "negative case - cartItemRepository",
			args: args{
				ctx: ctx,
				req: &desc.ListCartRequest{User: 1},
			},

			cartItemRepositoryMock: func(mc *minimock.Controller) postgres.CartItemRepository {
				mock := repositoryMocks.NewCartItemRepositoryMock(mc)
				mock.GetItemsMock.Expect(ctx, 1).Return(nil, cartItemRepositoryErr)
				return mock
			},
			limiterMock: func(mc *minimock.Controller) service.Limiter {
				return serviceMocks.NewLimiterMock(mc)
			},
			productServiceClientMock: func(mc *minimock.Controller) productservice.Client {
				return productServiceClientMocks.NewClientMock(mc)
			},

			want: nil,
			err:  cartItemRepositoryErr,
		},
		{
			name: "negative case - limiter error",
			args: args{
				ctx: ctx,
				req: &desc.ListCartRequest{User: 1},
			},

			cartItemRepositoryMock: func(mc *minimock.Controller) postgres.CartItemRepository {
				mock := repositoryMocks.NewCartItemRepositoryMock(mc)
				mock.GetItemsMock.Expect(ctx, 1).Return(cartItemRepositoryRes, nil)
				return mock
			},
			limiterMock: func(mc *minimock.Controller) service.Limiter {
				mock := serviceMocks.NewLimiterMock(mc)
				mock.WaitMock.Expect(ctx).Return(limiterErr)
				return mock
			},
			productServiceClientMock: func(mc *minimock.Controller) productservice.Client {
				return productServiceClientMocks.NewClientMock(mc)
			},

			want: nil,
			err:  limiterErr,
		},
		{
			name: "negative case - productServiceClient error",
			args: args{
				ctx: ctx,
				req: &desc.ListCartRequest{User: 1},
			},

			cartItemRepositoryMock: func(mc *minimock.Controller) postgres.CartItemRepository {
				mock := repositoryMocks.NewCartItemRepositoryMock(mc)
				mock.GetItemsMock.Expect(ctx, 1).Return(cartItemRepositoryRes, nil)
				return mock
			},
			limiterMock: func(mc *minimock.Controller) service.Limiter {
				mock := serviceMocks.NewLimiterMock(mc)
				mock.WaitMock.Expect(ctx).Return(nil)
				return mock
			},
			productServiceClientMock: func(mc *minimock.Controller) productservice.Client {
				mock := productServiceClientMocks.NewClientMock(mc)
				for _, res := range productServiceClientRes {
					mock.GetProductMock.When(ctx, res.SKU).Then(nil, productServiceClientErr)
				}
				return mock
			},

			want: nil,
			err:  productServiceClientErr,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			api := checkout.NewCheckout(service.New(
				nil,
				testCase.cartItemRepositoryMock(mc),
				nil,
				testCase.productServiceClientMock(mc),
				testCase.limiterMock(mc),
				nil,
			))

			res, err := api.ListCart(testCase.args.ctx, testCase.args.req)

			if testCase.want != nil {
				require.Equal(t, len(testCase.want.Items), len(res.Items))

				resProductsMap := toProductsMap(testCase.want)
				for _, product := range res.Items {
					require.Equal(t, product, resProductsMap[product.Sku])
				}
			}

			if testCase.err != nil {
				require.ErrorContains(t, err, testCase.err.Error())
			} else {
				require.Equal(t, testCase.err, err)
			}
		})
	}
}

func toProductsMap(cart *desc.ListCartResponse) map[uint32]*desc.Product {
	result := make(map[uint32]*desc.Product, len(cart.Items))
	for _, item := range cart.Items {
		result[item.Sku] = item
	}

	return result
}
