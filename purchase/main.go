package main

import (
	"github.com/enakai00/go-micro-gcp-example/purchase/handler"

	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"

	purchase "github.com/enakai00/go-micro-gcp-example/purchase/proto/purchase"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("com.example.service.purchase"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	purchase.RegisterPurchaseHandler(service.Server(), new(handler.Purchase))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
