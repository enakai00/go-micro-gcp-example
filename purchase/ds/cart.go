package ds

import (
	"context"
	"encoding/json"
	"os"

	"cloud.google.com/go/datastore"
	"github.com/google/uuid"
	log "github.com/micro/go-micro/v2/logger"
	"google.golang.org/api/iterator"

	"github.com/enakai00/go-micro-gcp-example/purchase/events"
	purchase "github.com/enakai00/go-micro-gcp-example/purchase/proto/purchase"
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
	client, _ = datastore.NewClient(context.Background(), projectID)
)

func getUUID() string {
	uuid, _ := uuid.NewRandom()
	return uuid.String()
}

func getCartStruct(cartid string) (*Cart, bool) {
	query := datastore.NewQuery("Cart").Filter("cartid =", cartid)
	it := client.Run(context.Background(), query)

	var cart Cart
	_, err := it.Next(&cart)
	if err == iterator.Done {
		return nil, false
	}
	return &cart, true
}

func getCartItem(cartid string, itemid string) (*CartItem, bool) {
	cart, ok := getCartStruct(cartid)
	if !ok {
		return nil, false
	}
	ancestor := cart.Key
	query := datastore.NewQuery("CartItem").Ancestor(ancestor).Filter("itemid =", itemid)
	it := client.Run(context.Background(), query)
	cartItem := CartItem{}
	_, err := it.Next(&cartItem)
	if err == iterator.Done {
		return nil, false
	}
	return &cartItem, true
}

func CreateCart(cartid string) *purchase.Cart {
	cart, ok := getCartStruct(cartid)
	if !ok {
		cartKey := datastore.IncompleteKey("Cart", nil)
		cart = &Cart{
			Cartid: cartid,
			Status: "open",
		}
		_, err := client.Put(context.Background(), cartKey, cart)
		if err != nil {
			log.Fatalf("Error stroing data: %v", err)
		}
	}
	purchaseCart := purchase.Cart{
		Cartid: cart.Cartid,
		Status: cart.Status,
	}
	return &purchaseCart
}

func GetCart(cartid string) *purchase.Cart {
	cart, ok := getCartStruct(cartid)
	if !ok {
		return &purchase.Cart{}
	}
	purchaseCart := purchase.Cart{
		Cartid: cart.Cartid,
		Status: cart.Status,
	}
	return &purchaseCart
}

func CloseCart(cartid string) *purchase.Cart {
	cartStructWithoutTx, ok := getCartStruct(cartid)
	if !ok {
		return &purchase.Cart{}
	}
	_, err := client.RunInTransaction(context.Background(),
		func(tx *datastore.Transaction) error {
			var cartStruct Cart
			err := tx.Get(cartStructWithoutTx.Key, &cartStruct)
			if err != nil {
				return err
			}
			if cartStruct.Status != "open" {
				return nil // cannot close
			}
			cartStruct.Status = "closed"
			_, err = tx.Put(cartStruct.Key, &cartStruct)
			if err != nil {
				return err
			}
			return nil
		})
	if err != nil {
		log.Fatalf("Error stroing data: %v", err)
	}
	return GetCart(cartid)
}

func AddItem(cartid string, itemid string, count int64) []*purchase.CartItem {
	cartStructWithoutTx, ok := getCartStruct(cartid)
	if !ok {
		cartItems := make([]*purchase.CartItem, 0)
		return cartItems
	}

	_, err := client.RunInTransaction(context.Background(),
		func(tx *datastore.Transaction) error {

			var cartStruct Cart
			err := tx.Get(cartStructWithoutTx.Key, &cartStruct)
			if err != nil {
				return err
			}
			if cartStruct.Status != "open" {
				return nil // cannot add new items
			}

			var key *datastore.Key
			var data CartItem

			cartItemWithoutTx, ok := getCartItem(cartid, itemid)
			if ok {
				var cartItem CartItem
				err := tx.Get(cartItemWithoutTx.Key, &cartItem)
				if err != nil {
					return err
				}
				key = cartItem.Key
				data = CartItem{
					Itemid: itemid,
					Count:  cartItem.Count + count,
				}
			} else {
				cart, _ := getCartStruct(cartid)
				key = datastore.IncompleteKey("CartItem", cart.Key)
				data = CartItem{
					Itemid: itemid,
					Count:  count,
				}
			}
			_, err = tx.Put(key, &data)
			if err != nil {
				return err
			}
			return nil
		})
	if err != nil {
		log.Fatalf("Error stroing data: %v", err)
	}
	cartItems := GetCartContents(cartid)
	return cartItems
}

func getCartItems(cartid string) []*CartItem {
	cart, ok := getCartStruct(cartid)
	cartItems := []*CartItem{}
	if !ok {
		return cartItems
	}
	query := datastore.NewQuery("CartItem").Ancestor(cart.Key)
	it := client.Run(context.Background(), query)
	for {
		cartItem := CartItem{}
		_, err := it.Next(&cartItem)
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Error fetching next entity: %v", err)
		}
		cartItems = append(cartItems, &cartItem)
	}
	return cartItems
}

func GetCartContents(cartid string) []*purchase.CartItem {
	cartItems := []*purchase.CartItem{}
	for _, cartItem := range getCartItems(cartid) {
		purchaseCartItem := purchase.CartItem{
			Itemid: cartItem.Itemid,
			Count:  cartItem.Count,
		}
		cartItems = append(cartItems, &purchaseCartItem)
	}
	return cartItems
}

func Checkout(cartid string) *purchase.OrderTicket {
	purchaseOrderTicket := purchase.OrderTicket{}

	_, err := client.RunInTransaction(context.Background(),
		func(tx *datastore.Transaction) error {
			cartWithoutTx, ok := getCartStruct(cartid)
			if !ok {
				return nil // non existing cartid
			}

			var cart Cart
			err := tx.Get(cartWithoutTx.Key, &cart) // Read again within a transaction
			if err != nil {
				return err
			}
			if cart.Status != "closed" {
				return nil // can't checkout
			}
			cart.Status = "checked-out"
			_, err = tx.Put(cart.Key, &cart)
			if err != nil {
				return err
			}

			ticketKey := datastore.IncompleteKey("OrderTicket", nil)
			orderTicket := OrderTicket{
				Orderid: getUUID(),
				Cartid:  cartid,
				Status:  "pre-approved",
			}
			_, err = tx.Put(ticketKey, &orderTicket)
			if err != nil {
				return err
			}

			cartContents := GetCartContents(cartid)
			purchaseOrderTicket = purchase.OrderTicket{
				Orderid:   orderTicket.Orderid,
				Status:    orderTicket.Status,
				CartItems: cartContents,
			}
			jsonBytes, _ := json.Marshal(purchaseOrderTicket)
			eventEntity := events.EventEntity{
				Eventid:   getUUID(),
				Type:      "purchase.OrderTicket",
				Sent:      false,
				EventData: jsonBytes,
			}
			eventEntityKey := datastore.IncompleteKey("Event", nil)
			_, err = tx.Put(eventEntityKey, &eventEntity)
			if err != nil {
				return err
			}

			return nil
		})
	if err != nil {
		log.Fatalf("Error stroing data: %v", err)
	}
	events.Sendout()
	return &purchaseOrderTicket
}
