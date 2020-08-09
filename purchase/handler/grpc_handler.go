package handler

import (
	"context"

	purchase "github.com/enakai00/go-micro-gcp-example/purchase/proto/purchase"
)

type Purchase struct{}

func (e *Purchase) CreateCart(ctx context.Context,
	req *purchase.CreateCartRequest,
	rsp *purchase.CreateCartResponse) error {
	rsp.Cart = CreateCart(req.Cartid)
	return nil
}

func (e *Purchase) GetCart(ctx context.Context,
	req *purchase.GetCartRequest,
	rsp *purchase.GetCartResponse) error {
	rsp.Cart = GetCart(req.Cartid)
	return nil
}

func (e *Purchase) CloseCart(ctx context.Context,
	req *purchase.CloseCartRequest,
	rsp *purchase.CloseCartResponse) error {
	rsp.Cart = CloseCart(req.Cartid)
	return nil
}

func (e *Purchase) AddItem(ctx context.Context,
	req *purchase.AddItemRequest,
	rsp *purchase.AddItemResponse) error {
	rsp.CartItems = AddItem(req.Cartid, req.Itemid, req.Count)
	return nil
}

func (e *Purchase) GetCartContents(ctx context.Context,
	req *purchase.GetCartContentsRequest,
	rsp *purchase.GetCartContentsResponse) error {
	rsp.CartItems = GetCartContents(req.Cartid)
	return nil
}

func (e *Purchase) Checkout(ctx context.Context,
	req *purchase.CheckoutRequest,
	rsp *purchase.CheckoutResponse) error {
	rsp.OrderTicket = Checkout(req.Cartid)
	return nil
}

func (e *Purchase) GetOrderTicket(ctx context.Context,
	req *purchase.GetOrderTicketRequest,
	rsp *purchase.GetOrderTicketResponse) error {
	rsp.OrderTicket = GetOrderTicket(req.Orderid)
	return nil
}
