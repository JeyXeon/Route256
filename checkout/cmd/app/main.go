package main

import (
	"log"
	"net"
	checkout "route256/checkout/internal/api"
	"route256/checkout/internal/clients/loms"
	"route256/checkout/internal/clients/productservice"
	"route256/checkout/internal/config"
	"route256/checkout/internal/service"
	desc "route256/checkout/pkg/checkout"
	lomsapi "route256/checkout/pkg/loms"
	productserviceapi "route256/checkout/pkg/productservice"
	"route256/libs/clientconnwrapper"

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

	lomsConn, err := clientconnwrapper.GetConn(config.ConfigData.Services.Loms)
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer lomsConn.Close()

	productServiceConn, err := clientconnwrapper.GetConn(config.ConfigData.Services.ProductService)
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer productServiceConn.Close()

	token := config.ConfigData.Token
	lomsClient := loms.New(lomsapi.NewLomsClient(lomsConn))
	productServiceClient := productservice.New(productserviceapi.NewProductServiceClient(productServiceConn), token)

	checkoutService := service.New(lomsClient, productServiceClient)

	desc.RegisterCheckoutServer(s, checkout.NewCheckout(checkoutService))

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
