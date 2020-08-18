package handler

import (
	"context"

	"cloud.google.com/go/datastore"
	log "github.com/micro/go-micro/v2/logger"
	"google.golang.org/api/iterator"

	"github.com/enakai00/go-micro-gcp-example/purchase/ds"
	"github.com/enakai00/go-micro-gcp-example/purchase/events"
	purchase "github.com/enakai00/go-micro-gcp-example/purchase/proto/purchase"
)

func getCartStruct(cartid string) (*ds.Cart, bool) {
	var cart ds.Cart
	ok := getStructByID("cartid", cartid, "Cart", &cart)
	return &cart, ok
}

func getOrderTicketStruct(orderid string) (*ds.OrderTicket, bool) {
	var orderTicket ds.OrderTicket
	ok := getStructByID("orderid", orderid, "OrderTicket", &orderTicket)
	return &orderTicket, ok
}

func getCartItem(cartid string, itemid string) (*ds.CartItem, bool) {
	cart, ok := getCartStruct(cartid)
	if !ok {
		return nil, false
	}
	ancestor := cart.Key
	query := datastore.NewQuery(ds.Kind["CartItem"]).Ancestor(ancestor).Filter("itemid =", itemid).Limit(1)
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
		cart = &ds.Cart{
			Cartid: cartid,
			Status: "open",
		}
		err := createEntity("Cart", cart, nil, nil)
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

	if err == datastore.ErrNoSuchEntity {
		log.Warnf("Cart has been deleted: cartid = ", cartid)
	} else if err != nil {
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
				data = ds.CartItem{
					Itemid: itemid,
					Count:  cartItem.Count + count,
				}
				key = cartItem.Key
				_, err = tx.Put(key, &data)
			} else {
				data = ds.CartItem{
					Itemid: itemid,
					Count:  count,
				}
				cart, _ := getCartStruct(cartid)
				err = createEntity("CartItem", &data, cart.Key, nil)
			}
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
	cartWithoutTx, ok := getCartStruct(cartid)
	if !ok {
		log.Warnf("Cart doesn't exist: cartid = ", cartid)
		return &purchaseOrderTicket
	}
	_, err := ds.Client.RunInTransaction(context.Background(),
		func(tx *datastore.Transaction) error {
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

			orderTicket := ds.OrderTicket{
				Orderid: getUUID(),
				Cartid:  cartid,
				Status:  "pre-approved",
			}
			err = createEntity("OrderTicket", &orderTicket, nil, tx)
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
	if err == datastore.ErrNoSuchEntity {
		log.Warnf("Cart doesn't exist: cartid = ", cartid)
	} else if err != nil {
		log.Fatalf("Error stroing data: %v", err)
	}
	events.PublishEvents()
	return &purchaseOrderTicket
}
