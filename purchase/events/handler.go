package events

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"

	"cloud.google.com/go/datastore"
	"github.com/micro/go-micro/v2/broker"
	log "github.com/micro/go-micro/v2/logger"
	"google.golang.org/api/iterator"

	"github.com/enakai00/go-micro-gcp-example/purchase/ds"
	purchase "github.com/enakai00/go-micro-gcp-example/purchase/proto/purchase"
)

func getUUID() string {
	uuid, _ := uuid.NewRandom()
	return uuid.String()
}

func handlePurchaseOrderTicket(purchaseOrderTicket purchase.OrderTicket) error {
	newStatus := purchaseOrderTicket.Status
	orderid := purchaseOrderTicket.Orderid
	query := datastore.NewQuery("OrderTicket").Filter("orderid =", orderid)
	it := client.Run(context.Background(), query)
	var orderTicket ds.OrderTicket
	_, err := it.Next(&orderTicket)
	if err == iterator.Done {
		log.Errorf("Order with orderid %s doen't exist.", orderid)
		return nil
	} else if err != nil {
		log.Fatalf("Error reading data: %v", err)
	}

	_, err = client.RunInTransaction(context.Background(),
		func(tx *datastore.Transaction) error {
			err := tx.Get(orderTicket.Key, &orderTicket)
			if err != nil {
				return err
			}

			if newStatus == "reserve_failed" && orderTicket.Status == "pre-approved" {
				orderTicket.Status = "reserve_failed"
				_, err := tx.Put(orderTicket.Key, &orderTicket)
				if err != nil {
					log.Fatalf("Error storing data: %v", err)
				}
				return nil
			}

			if newStatus == "reserved" && orderTicket.Status == "pre-approved" {
				orderTicket.Status = "reserved"
				_, err := tx.Put(orderTicket.Key, &orderTicket)
				if err != nil {
					log.Fatalf("Error storing data: %v", err)
				}
				return nil
			}

			if newStatus == "paid" {
				orderTicket.Status = "approved"
				_, err := tx.Put(orderTicket.Key, &orderTicket)
				if err != nil {
					log.Fatalf("Error storing data: %v", err)
				}

				purchaseOrderTicket.Status = "approved"
				jsonBytes, _ := json.Marshal(purchaseOrderTicket)
				eventEntity := EventEntity{
					Eventid:   getUUID(),
					Type:      "purchase.OrderTicket",
					Sent:      false,
					EventData: jsonBytes,
				}
				eventEntityKey := datastore.IncompleteKey(EventPublishTable, nil)
				_, err = tx.Put(eventEntityKey, &eventEntity)
				if err != nil {
					return err
				}
				return nil
			}
			return nil
		})
	if err != nil {
		log.Fatalf("Error stroing data: %v", err)
	}
	Sendout(EventPublishTable)
	return nil
}

func eventHandler(p broker.Event) error {
	header := p.Message().Header
	eventid := header["eventid"]
	eventType := header["type"]

	query := datastore.NewQuery(EventRecordTable).Filter("eventid =", eventid)
	it := client.Run(context.Background(), query)
	var receivedEvent ReceivedEvent
	_, err := it.Next(&receivedEvent)
	if err == nil { // duplicated event
		p.Ack()
		return nil
	}
	if err != iterator.Done {
		log.Fatalf("Error reading datastore: %v", err)
	}

	switch eventType {
	case "purchase.OrderTicket":
		var purchaseOrderTicket purchase.OrderTicket
		err := json.Unmarshal(p.Message().Body, &purchaseOrderTicket)
		if err != nil {
			log.Fatalf("Error unmarshalling eventt: %v", err)
		}
		log.Infof("Handle event purchase.OrderTicket: %v", purchaseOrderTicket)
		err = handlePurchaseOrderTicket(purchaseOrderTicket)
		if err != nil {
			log.Fatalf("Error handling purchaseOrderTicket: %v", err)
		}
		RecordEvent(EventRecordTable, eventid)
		p.Ack()
	default:
		log.Infof("Unknown event type: %s", eventType)
		RecordEvent(EventRecordTable, eventid)
		p.Ack()
	}
	return nil
}
