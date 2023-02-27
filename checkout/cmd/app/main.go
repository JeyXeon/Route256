package main

import (
	"log"
	"net/http"
	"route256/checkout/internal/clients/loms"
	"route256/checkout/internal/clients/productservice"
	"route256/checkout/internal/config"
	"route256/checkout/internal/domain"
	"route256/checkout/internal/handlers/addtocart"
	"route256/checkout/internal/handlers/deletefromcart"
	"route256/checkout/internal/handlers/listcart"
	"route256/checkout/internal/handlers/purchase"
	"route256/libs/requestprocessor"
	"route256/libs/srvwrapper"
)

func main() {
	err := config.Init()
	if err != nil {
		log.Fatal("config init", err)
	}

	lomsRequestProcessor := requestprocessor.New(config.ConfigData.Services.Loms)
	productServiceRequestProcessor := requestprocessor.New(config.ConfigData.Services.ProductService)

	lomsClient := loms.New(lomsRequestProcessor)
	productServiceClient := productservice.New(productServiceRequestProcessor, config.ConfigData.Token)

	businessLogic := domain.New(lomsClient, productServiceClient)

	addToCartHandler := addtocart.New(businessLogic)
	deleteFromCartHandler := deletefromcart.New(businessLogic)
	purchaseHandler := purchase.New(businessLogic)
	listCartHandler := listcart.New(businessLogic)

	http.Handle("/addToCart", srvwrapper.New(addToCartHandler.Handle))
	http.Handle("/deleteFromCart", srvwrapper.New(deleteFromCartHandler.Handle))
	http.Handle("/purchase", srvwrapper.New(purchaseHandler.Handle))
	http.Handle("/listCart", srvwrapper.New(listCartHandler.Handle))

	port := config.ConfigData.Port

	log.Println("listening http at", port)
	err = http.ListenAndServe(port, nil)
	log.Fatal("cannot listen http", err)
}
