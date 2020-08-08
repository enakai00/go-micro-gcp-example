// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/purchase/purchase.proto

package com_example_service_purchase

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	api "github.com/micro/go-micro/v2/api"
	client "github.com/micro/go-micro/v2/client"
	server "github.com/micro/go-micro/v2/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for Purchase service

func NewPurchaseEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for Purchase service

type PurchaseService interface {
	CreateCart(ctx context.Context, in *CreateCartRequest, opts ...client.CallOption) (*CreateCartResponse, error)
	GetCart(ctx context.Context, in *GetCartRequest, opts ...client.CallOption) (*GetCartResponse, error)
	AddItem(ctx context.Context, in *AddItemRequest, opts ...client.CallOption) (*AddItemResponse, error)
	GetCartContents(ctx context.Context, in *GetCartContentsRequest, opts ...client.CallOption) (*GetCartContentsResponse, error)
	CloseCart(ctx context.Context, in *CloseCartRequest, opts ...client.CallOption) (*CloseCartResponse, error)
	Checkout(ctx context.Context, in *CheckoutRequest, opts ...client.CallOption) (*CheckoutResponse, error)
	GetOrderTicket(ctx context.Context, in *GetOrderTicketRequest, opts ...client.CallOption) (*GetOrderTicketResponse, error)
}

type purchaseService struct {
	c    client.Client
	name string
}

func NewPurchaseService(name string, c client.Client) PurchaseService {
	return &purchaseService{
		c:    c,
		name: name,
	}
}

func (c *purchaseService) CreateCart(ctx context.Context, in *CreateCartRequest, opts ...client.CallOption) (*CreateCartResponse, error) {
	req := c.c.NewRequest(c.name, "Purchase.CreateCart", in)
	out := new(CreateCartResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *purchaseService) GetCart(ctx context.Context, in *GetCartRequest, opts ...client.CallOption) (*GetCartResponse, error) {
	req := c.c.NewRequest(c.name, "Purchase.GetCart", in)
	out := new(GetCartResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *purchaseService) AddItem(ctx context.Context, in *AddItemRequest, opts ...client.CallOption) (*AddItemResponse, error) {
	req := c.c.NewRequest(c.name, "Purchase.AddItem", in)
	out := new(AddItemResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *purchaseService) GetCartContents(ctx context.Context, in *GetCartContentsRequest, opts ...client.CallOption) (*GetCartContentsResponse, error) {
	req := c.c.NewRequest(c.name, "Purchase.GetCartContents", in)
	out := new(GetCartContentsResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *purchaseService) CloseCart(ctx context.Context, in *CloseCartRequest, opts ...client.CallOption) (*CloseCartResponse, error) {
	req := c.c.NewRequest(c.name, "Purchase.CloseCart", in)
	out := new(CloseCartResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *purchaseService) Checkout(ctx context.Context, in *CheckoutRequest, opts ...client.CallOption) (*CheckoutResponse, error) {
	req := c.c.NewRequest(c.name, "Purchase.Checkout", in)
	out := new(CheckoutResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *purchaseService) GetOrderTicket(ctx context.Context, in *GetOrderTicketRequest, opts ...client.CallOption) (*GetOrderTicketResponse, error) {
	req := c.c.NewRequest(c.name, "Purchase.GetOrderTicket", in)
	out := new(GetOrderTicketResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Purchase service

type PurchaseHandler interface {
	CreateCart(context.Context, *CreateCartRequest, *CreateCartResponse) error
	GetCart(context.Context, *GetCartRequest, *GetCartResponse) error
	AddItem(context.Context, *AddItemRequest, *AddItemResponse) error
	GetCartContents(context.Context, *GetCartContentsRequest, *GetCartContentsResponse) error
	CloseCart(context.Context, *CloseCartRequest, *CloseCartResponse) error
	Checkout(context.Context, *CheckoutRequest, *CheckoutResponse) error
	GetOrderTicket(context.Context, *GetOrderTicketRequest, *GetOrderTicketResponse) error
}

func RegisterPurchaseHandler(s server.Server, hdlr PurchaseHandler, opts ...server.HandlerOption) error {
	type purchase interface {
		CreateCart(ctx context.Context, in *CreateCartRequest, out *CreateCartResponse) error
		GetCart(ctx context.Context, in *GetCartRequest, out *GetCartResponse) error
		AddItem(ctx context.Context, in *AddItemRequest, out *AddItemResponse) error
		GetCartContents(ctx context.Context, in *GetCartContentsRequest, out *GetCartContentsResponse) error
		CloseCart(ctx context.Context, in *CloseCartRequest, out *CloseCartResponse) error
		Checkout(ctx context.Context, in *CheckoutRequest, out *CheckoutResponse) error
		GetOrderTicket(ctx context.Context, in *GetOrderTicketRequest, out *GetOrderTicketResponse) error
	}
	type Purchase struct {
		purchase
	}
	h := &purchaseHandler{hdlr}
	return s.Handle(s.NewHandler(&Purchase{h}, opts...))
}

type purchaseHandler struct {
	PurchaseHandler
}

func (h *purchaseHandler) CreateCart(ctx context.Context, in *CreateCartRequest, out *CreateCartResponse) error {
	return h.PurchaseHandler.CreateCart(ctx, in, out)
}

func (h *purchaseHandler) GetCart(ctx context.Context, in *GetCartRequest, out *GetCartResponse) error {
	return h.PurchaseHandler.GetCart(ctx, in, out)
}

func (h *purchaseHandler) AddItem(ctx context.Context, in *AddItemRequest, out *AddItemResponse) error {
	return h.PurchaseHandler.AddItem(ctx, in, out)
}

func (h *purchaseHandler) GetCartContents(ctx context.Context, in *GetCartContentsRequest, out *GetCartContentsResponse) error {
	return h.PurchaseHandler.GetCartContents(ctx, in, out)
}

func (h *purchaseHandler) CloseCart(ctx context.Context, in *CloseCartRequest, out *CloseCartResponse) error {
	return h.PurchaseHandler.CloseCart(ctx, in, out)
}

func (h *purchaseHandler) Checkout(ctx context.Context, in *CheckoutRequest, out *CheckoutResponse) error {
	return h.PurchaseHandler.Checkout(ctx, in, out)
}

func (h *purchaseHandler) GetOrderTicket(ctx context.Context, in *GetOrderTicketRequest, out *GetOrderTicketResponse) error {
	return h.PurchaseHandler.GetOrderTicket(ctx, in, out)
}
