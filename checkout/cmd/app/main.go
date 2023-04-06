package main

import (
	"context"
	"net"
	checkout "route256/checkout/internal/api"
	"route256/checkout/internal/clients/loms"
	"route256/checkout/internal/clients/productservice"
	"route256/checkout/internal/config"
	"route256/checkout/internal/repository/postgres"
	"route256/checkout/internal/service"
	desc "route256/checkout/pkg/checkout"
	lomsapi "route256/checkout/pkg/loms"
	productserviceapi "route256/checkout/pkg/productservice"
	"route256/libs/clientconnwrapper"
	"route256/libs/dbmanager"
	"route256/libs/logger"
	"route256/libs/ratelimiter"
	"time"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	logger.Init(false)

	err := config.Init()
	if err != nil {
		logger.Fatal("config init", zap.Error(err))
	}

	port := config.ConfigData.Port

	lis, err := net.Listen("tcp", port)
	if err != nil {
		logger.Fatal("failed to listen", zap.Error(err))
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpcMiddleware.ChainUnaryServer(
				logger.LoggingInterceptor,
			),
		),
	)
	reflection.Register(s)

	checkoutDbUrl := config.ConfigData.CheckoutDbUrl
	pool, err := pgxpool.Connect(context.Background(), checkoutDbUrl)
	if err != nil {
		logger.Fatal("failed to connect to db", zap.Error(err))
	}

	dbManager := dbmanager.New(pool)
	orderItemRepository := postgres.NewCartItemRepository(dbManager)

	lomsConn, err := clientconnwrapper.GetConn(config.ConfigData.Services.Loms.Url)
	if err != nil {
		logger.Fatal("failed to connect to loms server", zap.Error(err))
	}
	defer lomsConn.Close()

	productServiceConn, err := clientconnwrapper.GetConn(config.ConfigData.Services.ProductService.Url)
	if err != nil {
		logger.Fatal("failed to connect to productService server", zap.Error(err))
	}
	defer productServiceConn.Close()

	token := config.ConfigData.Token
	lomsClient := loms.New(lomsapi.NewLomsClient(lomsConn))
	productServiceClient := productservice.New(productserviceapi.NewProductServiceClient(productServiceConn), token)

	productServiceRateSeconds := config.ConfigData.Services.ProductService.RateSeconds
	productServiceTokens := config.ConfigData.Services.ProductService.Tokens
	productServiceLimiter := ratelimiter.NewLimiter(
		context.Background(),
		time.Duration(productServiceRateSeconds),
		productServiceTokens,
	)

	checkoutService := service.New(
		dbManager,
		orderItemRepository,
		lomsClient,
		productServiceClient,
		productServiceLimiter,
	)

	desc.RegisterCheckoutServer(s, checkout.NewCheckout(checkoutService))

	logger.Info("server listening", zap.Any("address", lis.Addr()))

	if err = s.Serve(lis); err != nil {
		logger.Fatal("failed to serve", zap.Error(err))
	}

}
