package events

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"

	"cloud.google.com/go/datastore"
	"github.com/micro/go-micro/v2/broker"
	log "github.com/micro/go-micro/v2/logger"

	purchase "github.com/enakai00/go-micro-gcp-example/stock/proto/purchase"
)

func getUUID() string {
	uuid, _ := uuid.NewRandom()
	return uuid.String()
}

func handlePurchaseOrderTicket(purchaseOrderTicket purchase.OrderTicket) error {
	switch purchaseOrderTicket.Status {
	case "pre-approved":
		_, err := client.RunInTransaction(context.Background(),
			func(tx *datastore.Transaction) error {
				// Reserve stocks in the order
				purchaseOrderTicket.Status = "reserved"
				for _, cartItem := range purchaseOrderTicket.CartItems {
					if cartItem.Itemid == "yellow" { // yellow is sold out.
						purchaseOrderTicket.Status = "reserve_failed"
					}
					if cartItem.Itemid == "red" { // red is free!
						purchaseOrderTicket.Status = "paid"
					}
				}

				err := RegisterEvent(purchaseOrderTicket,
					"purchase.OrderTicket",
					tx)
				if err != nil {
					return err
				}
				return nil
			})
		if err != nil {
			return err
		}
		PublishEvents()
	}
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
			recordEvent(eventid)
			p.Ack()
		}
	default:
		log.Infof("Unknown event type: %s", eventType)
		recordEvent(eventid)
		p.Ack()
	}
	return nil
}
