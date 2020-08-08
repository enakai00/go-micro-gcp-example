package handler

import (
	"context"
	"net/http"

	"github.com/labstack/echo"

	purchase "github.com/enakai00/go-micro-gcp-example/purchase/proto/purchase"
)

var client purchase.PurchaseService

func Register(e *echo.Echo, purchaseClient purchase.PurchaseService) {
	client = purchaseClient
	e.POST("/purchase/api/create_cart", createCart)
}

//	e.POST("/purcahse/api/create_cart", createCart)
func createCart(c echo.Context) error {
	var request purchase.CreateCartRequest
	err := c.Bind(&request)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Server Error")
	}

	rsp, err := client.CreateCart(context.Background(), &request)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Server Error")
	}

	return c.JSON(http.StatusCreated, rsp)
}
