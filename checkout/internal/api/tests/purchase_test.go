package checkout

import (
	"context"
	checkout "route256/checkout/internal/api"
	"route256/checkout/internal/clients/loms"
	lomsClientMocks "route256/checkout/internal/clients/loms/mocks"
	"route256/checkout/internal/model"
	"route256/checkout/internal/repository/postgres"
	repositoryMocks "route256/checkout/internal/repository/postgres/mocks"
	"route256/checkout/internal/service"
	desc "route256/checkout/pkg/checkout"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestPurchase(t *testing.T) {
	type cartItemRepositoryMockFunc func(mc *minimock.Controller) postgres.CartItemRepository
	type lomsClientMockFunc func(mc *minimock.Controller) loms.Client

	type args struct {
		ctx context.Context
		req *desc.PurchaseRequest
	}

	var (
		mc  = minimock.NewController(t)
		ctx = context.Background()

		cartItemRepositoryErr = errors.New("cart item repository error")
		lomsClientErr         = errors.New("loms client error")
		errInsufficientOrder  = errors.New("insufficient order")

		cartItemRepositoryRes = []*model.CartItem{{
			SKU:   2,
			Count: 1,
		}}
	)
	t.Cleanup(mc.Finish)

	type purchaseTestCase struct {
		name string
		args args

		cartItemRepositoryMock cartItemRepositoryMockFunc
		lomsClientMock         lomsClientMockFunc

		want *desc.PurchaseResponse
		err  error
	}

	testCases := []purchaseTestCase{
		{
			name: "positive case",
			args: args{
				ctx: ctx,
				req: &desc.PurchaseRequest{User: 1},
			},

			cartItemRepositoryMock: func(mc *minimock.Controller) postgres.CartItemRepository {
				mock := repositoryMocks.NewCartItemRepositoryMock(mc)
				mock.GetItemsMock.Expect(ctx, 1).Return(cartItemRepositoryRes, nil)
				return mock
			},
			lomsClientMock: func(mc *minimock.Controller) loms.Client {
				mock := lomsClientMocks.NewClientMock(mc)
				mock.CreateOrderMock.Expect(ctx, 1, cartItemRepositoryRes).Return(3, nil)
				return mock
			},

			want: &desc.PurchaseResponse{OrderID: 3},
			err:  nil,
		},
		{
			name: "negative case - empty user",
			args: args{
				ctx: ctx,
				req: &desc.PurchaseRequest{User: 0},
			},

			cartItemRepositoryMock: func(mc *minimock.Controller) postgres.CartItemRepository {
				return repositoryMocks.NewCartItemRepositoryMock(mc)
			},
			lomsClientMock: func(mc *minimock.Controller) loms.Client {
				mock := lomsClientMocks.NewClientMock(mc)
				return mock
			},

			want: nil,
			err:  checkout.ErrPurchaseEmptyUser,
		},
		{
			name: "negative case - cartItemRepository error",
			args: args{
				ctx: ctx,
				req: &desc.PurchaseRequest{User: 1},
			},

			cartItemRepositoryMock: func(mc *minimock.Controller) postgres.CartItemRepository {
				mock := repositoryMocks.NewCartItemRepositoryMock(mc)
				mock.GetItemsMock.Expect(ctx, 1).Return(nil, cartItemRepositoryErr)
				return mock
			},
			lomsClientMock: func(mc *minimock.Controller) loms.Client {
				mock := lomsClientMocks.NewClientMock(mc)
				return mock
			},

			want: nil,
			err:  cartItemRepositoryErr,
		},
		{
			name: "negative case - lomsClient error",
			args: args{
				ctx: ctx,
				req: &desc.PurchaseRequest{User: 1},
			},

			cartItemRepositoryMock: func(mc *minimock.Controller) postgres.CartItemRepository {
				mock := repositoryMocks.NewCartItemRepositoryMock(mc)
				mock.GetItemsMock.Expect(ctx, 1).Return(cartItemRepositoryRes, nil)
				return mock
			},
			lomsClientMock: func(mc *minimock.Controller) loms.Client {
				mock := lomsClientMocks.NewClientMock(mc)
				mock.CreateOrderMock.Expect(ctx, 1, cartItemRepositoryRes).Return(0, lomsClientErr)
				return mock
			},

			want: nil,
			err:  lomsClientErr,
		},
		{
			name: "negative case - insufficient order error",
			args: args{
				ctx: ctx,
				req: &desc.PurchaseRequest{User: 1},
			},

			cartItemRepositoryMock: func(mc *minimock.Controller) postgres.CartItemRepository {
				mock := repositoryMocks.NewCartItemRepositoryMock(mc)
				mock.GetItemsMock.Expect(ctx, 1).Return(cartItemRepositoryRes, nil)
				return mock
			},
			lomsClientMock: func(mc *minimock.Controller) loms.Client {
				mock := lomsClientMocks.NewClientMock(mc)
				mock.CreateOrderMock.Expect(ctx, 1, cartItemRepositoryRes).Return(0, nil)
				return mock
			},

			want: nil,
			err:  errInsufficientOrder,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			api := checkout.NewCheckout(service.New(
				nil,
				testCase.cartItemRepositoryMock(mc),
				testCase.lomsClientMock(mc),
				nil,
				nil,
				nil,
			))

			res, err := api.Purchase(testCase.args.ctx, testCase.args.req)

			require.Equal(t, testCase.want, res)

			if testCase.err != nil {
				require.ErrorContains(t, err, testCase.err.Error())
			} else {
				require.Equal(t, testCase.err, err)
			}
		})
	}
}
