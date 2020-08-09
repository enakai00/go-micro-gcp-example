package main

import (
//	"github.com/enakai00/go-micro-gcp-example/stock/handler"
	_ "github.com/enakai00/go-micro-gcp-example/stock/events"

	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2"

//	stock "github.com/enakai00/go-micro-gcp-example/stock/proto/stock"

        _ "github.com/micro/go-plugins/registry/kubernetes/v2"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("com.example.service.stock"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
//	stock.RegisterStockHandler(service.Server(), new(handler.Stock))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
