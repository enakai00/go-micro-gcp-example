package ds

import (
	"context"
	"os"

	"cloud.google.com/go/datastore"
)

type Cart struct {
	Cartid string         `datastore:"cartid"`
	Status string         `datastore:"status"` // open, closed, checked-out
	Key    *datastore.Key `datastore:"__key__"`
}

type CartItem struct {
	Itemid string         `datastore:"itemid"`
	Count  int64          `datastore:"count"`
	Key    *datastore.Key `datastore:"__key__"`
}

type OrderTicket struct {
	Orderid string         `datastore:"orderid"`
	Cartid  string         `datastore:"cartid"`
	Status  string         `datastore:"status"`
	Key     *datastore.Key `datastore:"__key__"`
}

var (
	projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")
	Client, _ = datastore.NewClient(context.Background(), projectID)

	Kind = map[string]string{
		"Cart":        "Cart",
		"CartItem":    "CartItem",
		"OrderTicket": "OrderTicket",
	}
)
