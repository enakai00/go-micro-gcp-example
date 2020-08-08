package main

import (
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"

	"github.com/enakai00/go-micro-gcp-example/api-gateway/handler"
	purchase "github.com/enakai00/go-micro-gcp-example/api-gateway/proto/purchase"

	_ "github.com/micro/go-plugins/registry/kubernetes/v2"
)

func createMux() *echo.Echo {
	e := echo.New()
	http.Handle("/", e)
	return e
}

// e.GET("/", home)
func home(c echo.Context) error {
	return c.String(http.StatusOK, "api-gateway")
}

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("com.example.service.api-gateway"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	purchaseClient := purchase.NewPurchaseService(
		"com.example.service.purchase", service.Client())

	e := createMux()
	e.GET("/", home)
	handler.Register(e, purchaseClient)

	// Run HTTP server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Infof("Defaulting to port %s", port)
	}
	log.Infof("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
