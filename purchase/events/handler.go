package events

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"

	"cloud.google.com/go/datastore"
	"github.com/micro/go-micro/v2/broker"
	log "github.com/micro/go-micro/v2/logger"

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
	query = query.Limit(1).KeysOnly()
	keys, err := client.GetAll(context.Background(), query, nil)
	if err != nil {
		return err
	}
	if len(keys) == 0 {
		log.Errorf("Order with orderid %s doen't exist.", orderid)
		return nil
	}

	_, err = client.RunInTransaction(context.Background(),
		func(tx *datastore.Transaction) error {
			var orderTicket ds.OrderTicket
			err := tx.Get(keys[0], &orderTicket)
			if err != nil {
				return err
			}

			if newStatus == "reserve_failed" && orderTicket.Status == "pre-approved" {
				orderTicket.Status = "reserve_failed"
				_, err := tx.Put(orderTicket.Key, &orderTicket)
				if err != nil {
					return err
				}
				return nil
			}

			if newStatus == "reserved" && orderTicket.Status == "pre-approved" {
				orderTicket.Status = "reserved"
				_, err := tx.Put(orderTicket.Key, &orderTicket)
				if err != nil {
					return err
				}
				return nil
			}

			if newStatus == "paid" {
				orderTicket.Status = "approved"
				_, err := tx.Put(orderTicket.Key, &orderTicket)
				if err != nil {
					return err
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
		return err
	}
	Sendout(EventPublishTable)
	return nil
}

func eventHandler(p broker.Event) error {
	header := p.Message().Header
	eventid := header["eventid"]
	eventType := header["type"]

	duplicate, err := isDuplicated(eventid)
	if err != nil {
		log.Fatalf("Error checking duplicated event: %v", err)
	}
	if duplicate {
		p.Ack()
		return nil
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
			log.Warnf("Failed to handle purchaseOrderTicke: %v", err)
		} else {
			RecordEvent(EventRecordTable, eventid)
			p.Ack()
		}
	default:
		log.Infof("Unknown event type: %s", eventType)
		RecordEvent(EventRecordTable, eventid)
		p.Ack()
	}
	return nil
}
