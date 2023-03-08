package main

import (
	"context"
	"log"
	"net"
	"route256/libs/dbmanager"
	loms "route256/loms/internal/api"
	"route256/loms/internal/config"
	"route256/loms/internal/repository/postgres"
	"route256/loms/internal/service"
	desc "route256/loms/pkg/loms"

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

	lomsDbUrl := config.ConfigData.LomsDbUrl
	pool, err := pgxpool.Connect(context.Background(), lomsDbUrl)
	if err != nil {
		log.Fatal("db connect", err)
	}

	dbManager := dbmanager.New(pool)

	itemRepository := postgres.NewItemsRepository(dbManager)
	stocksRepository := postgres.NewStocksRepository(dbManager)
	orderRepository := postgres.NewOrderRepository(dbManager)

	lomsService := service.New(dbManager, itemRepository, stocksRepository, orderRepository)

	desc.RegisterLomsServer(s, loms.NewLoms(lomsService))

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
