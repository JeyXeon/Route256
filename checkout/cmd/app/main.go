package main

import (
	"context"
	"log"
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
	"route256/libs/ratelimiter"

	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	err := config.Init()
	if err != nil {
		log.Fatal("config init", err)
	}

	port := config.ConfigData.Port

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)

	checkoutDbUrl := config.ConfigData.CheckoutDbUrl
	pool, err := pgxpool.Connect(context.Background(), checkoutDbUrl)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	dbManager := dbmanager.New(pool)
	orderItemRepository := postgres.NewCartItemRepository(dbManager)

	lomsConn, err := clientconnwrapper.GetConn(config.ConfigData.Services.Loms.Url)
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer lomsConn.Close()

	productServiceConn, err := clientconnwrapper.GetConn(config.ConfigData.Services.ProductService.Url)
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer productServiceConn.Close()

	token := config.ConfigData.Token
	lomsClient := loms.New(lomsapi.NewLomsClient(lomsConn))
	productServiceClient := productservice.New(productserviceapi.NewProductServiceClient(productServiceConn), token)

	productServiceRateSeconds := config.ConfigData.Services.ProductService.RateSeconds
	productServiceTokens := config.ConfigData.Services.ProductService.Tokens
	productServiceLimiter := ratelimiter.NewLimiter(
		context.Background(),
		productServiceRateSeconds,
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

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
