package handler

import (
	"context"
	"net/http"

	"github.com/labstack/echo"

	purchase "github.com/enakai00/go-micro-gcp-example/api-gateway/proto/purchase"
)

var client purchase.PurchaseService
var serverError = echo.NewHTTPError(http.StatusInternalServerError, "Server Error")

func Register(e *echo.Echo, purchaseClient purchase.PurchaseService) {
	client = purchaseClient
	e.POST("/purchase/api/create_cart", createCart)
	e.POST("/purchase/api/get_cart", getCart)
	e.POST("/purchase/api/add_item", addItem)
	e.POST("/purchase/api/get_cart_contents", getCartContents)
	e.POST("/purchase/api/close_cart", closeCart)
	e.POST("/purchase/api/checkout", checkout)
	e.POST("/purchase/api/get_order_ticket", getOrderTicket)
}

//	e.POST("/purcahse/api/create_cart", createCart)
func createCart(c echo.Context) error {
	var request purchase.CreateCartRequest
	err := c.Bind(&request)
	if err != nil {
		return serverError
	}

	resp, err := client.CreateCart(context.Background(), &request)
	if err != nil {
		return serverError
	}

	return c.JSON(http.StatusCreated, resp)
}

//	e.POST("/purchase/api/get_cart", getCart)
func getCart(c echo.Context) error {
	var request purchase.GetCartRequest
	err := c.Bind(&request)
	if err != nil {
		return serverError
	}

	resp, err := client.GetCart(context.Background(), &request)
	if err != nil {
		return serverError
	}

	return c.JSON(http.StatusCreated, resp)
}

//	e.POST("/purcahse/api/add_item", addItem)
func addItem(c echo.Context) error {
	var request purchase.AddItemRequest
	err := c.Bind(&request)
	if err != nil {
		return serverError
	}

	resp, err := client.AddItem(context.Background(), &request)
	if err != nil {
		return serverError
	}

	return c.JSON(http.StatusCreated, resp)
}

//	e.POST("/purcahse/api/get_cart_contents", getCartContents)
func getCartContents(c echo.Context) error {
	var request purchase.GetCartContentsRequest
	err := c.Bind(&request)
	if err != nil {
		return serverError
	}

	resp, err := client.GetCartContents(context.Background(), &request)
	if err != nil {
		return serverError
	}

	return c.JSON(http.StatusCreated, resp)
}

//	e.POST("/purchase/api/close_cart", closeCart)
func closeCart(c echo.Context) error {
	var request purchase.CloseCartRequest
	err := c.Bind(&request)
	if err != nil {
		return serverError
	}

	resp, err := client.CloseCart(context.Background(), &request)
	if err != nil {
		return serverError
	}

	return c.JSON(http.StatusCreated, resp)
}

//	e.POST("/purchase/api/checkout", checkout)
func checkout(c echo.Context) error {
	var request purchase.CheckoutRequest
	err := c.Bind(&request)
	if err != nil {
		return serverError
	}

	resp, err := client.Checkout(context.Background(), &request)
	if err != nil {
		return serverError
	}

	return c.JSON(http.StatusCreated, resp)
}

//	e.POST("/purchase/api/get_order_ticket", getOrderTicket)
func getOrderTicket(c echo.Context) error {
	var request purchase.GetOrderTicketRequest
	err := c.Bind(&request)
	if err != nil {
		return serverError
	}

	resp, err := client.GetOrderTicket(context.Background(), &request)
	if err != nil {
		return serverError
	}

	return c.JSON(http.StatusCreated, resp)
}
