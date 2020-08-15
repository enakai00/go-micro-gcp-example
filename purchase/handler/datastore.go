package handler

import (
	"context"

	"cloud.google.com/go/datastore"
	"github.com/google/uuid"
	log "github.com/micro/go-micro/v2/logger"
	"google.golang.org/api/iterator"

	"github.com/enakai00/go-micro-gcp-example/purchase/ds"
	"github.com/enakai00/go-micro-gcp-example/purchase/events"
	purchase "github.com/enakai00/go-micro-gcp-example/purchase/proto/purchase"
)

func getUUID() string {
	uuid, _ := uuid.NewRandom()
	return uuid.String()
}

func getCartStruct(cartid string) (*ds.Cart, bool) {
	query := datastore.NewQuery(ds.Kind["Cart"]).Filter("cartid =", cartid)
	it := ds.Client.Run(context.Background(), query)

	var cart ds.Cart
	_, err := it.Next(&cart)
	if err == iterator.Done {
		return nil, false
	}
	if err != nil {
		log.Fatalf("Error reading data: %v", err)
	}
	return &cart, true
}

func getOrderTicketStruct(orderid string) (*ds.OrderTicket, bool) {
	query := datastore.NewQuery(ds.Kind["OrderTicket"]).Filter("orderid =", orderid)
	it := ds.Client.Run(context.Background(), query)

	var orderTicket ds.OrderTicket
	_, err := it.Next(&orderTicket)
	if err == iterator.Done {
		return nil, false
	}
	if err != nil {
		log.Fatalf("Error reading data: %v", err)
	}
	return &orderTicket, true
}

func getCartItem(cartid string, itemid string) (*ds.CartItem, bool) {
	cart, ok := getCartStruct(cartid)
	if !ok {
		return nil, false
	}
	ancestor := cart.Key
	query := datastore.NewQuery(ds.Kind["CartItem"]).Ancestor(ancestor).Filter("itemid =", itemid)
	it := ds.Client.Run(context.Background(), query)
	cartItem := ds.CartItem{}
	_, err := it.Next(&cartItem)
	if err == iterator.Done {
		return nil, false
	}
	if err != nil {
		log.Fatalf("Error reading data: %v", err)
	}
	return &cartItem, true
}

func CreateCart(cartid string) *purchase.Cart {
	cart, ok := getCartStruct(cartid)
	if !ok {
		cartKey := datastore.IncompleteKey(ds.Kind["Cart"], nil)
		cart = &ds.Cart{
			Cartid: cartid,
			Status: "open",
		}
		_, err := ds.Client.Put(context.Background(), cartKey, cart)
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
	_, err := ds.Client.RunInTransaction(context.Background(),
		func(tx *datastore.Transaction) error {
			var cartStruct ds.Cart
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

	_, err := ds.Client.RunInTransaction(context.Background(),
		func(tx *datastore.Transaction) error {

			var cartStruct ds.Cart
			err := tx.Get(cartStructWithoutTx.Key, &cartStruct)
			if err != nil {
				return err
			}
			if cartStruct.Status != "open" {
				return nil // cannot add new items
			}

			var key *datastore.Key
			var data ds.CartItem

			cartItemWithoutTx, ok := getCartItem(cartid, itemid)
			if ok {
				var cartItem ds.CartItem
				err := tx.Get(cartItemWithoutTx.Key, &cartItem)
				if err != nil {
					return err
				}
				key = cartItem.Key
				data = ds.CartItem{
					Itemid: itemid,
					Count:  cartItem.Count + count,
				}
			} else {
				cart, _ := getCartStruct(cartid)
				key = datastore.IncompleteKey(ds.Kind["CartItem"], cart.Key)
				data = ds.CartItem{
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

func getCartItems(cartid string) []*ds.CartItem {
	cart, ok := getCartStruct(cartid)
	cartItems := []*ds.CartItem{}
	if !ok {
		return cartItems
	}
	query := datastore.NewQuery(ds.Kind["CartItem"]).Ancestor(cart.Key)
	it := ds.Client.Run(context.Background(), query)
	for {
		cartItem := ds.CartItem{}
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

func GetOrderTicket(orderid string) *purchase.OrderTicket {
	orderTicket, ok := getOrderTicketStruct(orderid)
	if !ok {
		return &purchase.OrderTicket{}
	}
	purchaseCartItems := GetCartContents(orderTicket.Cartid)
	purchaseOrderTicket := purchase.OrderTicket{
		Orderid:   orderTicket.Orderid,
		Status:    orderTicket.Status,
		CartItems: purchaseCartItems,
	}
	return &purchaseOrderTicket
}

func Checkout(cartid string) *purchase.OrderTicket {
	purchaseOrderTicket := purchase.OrderTicket{}

	_, err := ds.Client.RunInTransaction(context.Background(),
		func(tx *datastore.Transaction) error {
			cartWithoutTx, ok := getCartStruct(cartid)
			if !ok {
				return nil // non existing cartid
			}

			var cart ds.Cart
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

			ticketKey := datastore.IncompleteKey(ds.Kind["OrderTicket"], nil)
			orderTicket := ds.OrderTicket{
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
			err = events.RegisterEvent(purchaseOrderTicket,
				"purchase.OrderTicket",
				tx)
			if err != nil {
				return err
			}

			return nil
		})
	if err != nil {
		log.Fatalf("Error stroing data: %v", err)
	}
	events.PublishEvents()
	return &purchaseOrderTicket
}
