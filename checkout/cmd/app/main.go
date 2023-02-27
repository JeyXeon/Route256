package main

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	checkout "route256/checkout/internal/api"
	"route256/checkout/internal/clients/loms"
	"route256/checkout/internal/clients/productservice"
	"route256/checkout/internal/config"
	"route256/checkout/internal/service"
	desc "route256/checkout/pkg"
	"route256/libs/requestprocessor"
)

func main() {
	err := config.Init()
	if err != nil {
		log.Fatal("config init", err)
	}

	port := config.ConfigData.Port

	lis, err := net.Listen("tcp", fmt.Sprintf("%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)

	lomsRequestProcessor := requestprocessor.New(config.ConfigData.Services.Loms)
	productServiceRequestProcessor := requestprocessor.New(config.ConfigData.Services.ProductService)

	lomsClient := loms.New(lomsRequestProcessor)
	productServiceClient := productservice.New(productServiceRequestProcessor, config.ConfigData.Token)

	checkoutService := service.New(lomsClient, productServiceClient)

	desc.RegisterCheckoutServer(s, checkout.NewCheckout(checkoutService))

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
