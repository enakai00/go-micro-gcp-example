package events

import (
	"context"
	"encoding/json"

	"cloud.google.com/go/datastore"
	"github.com/micro/go-micro/v2/broker"
	log "github.com/micro/go-micro/v2/logger"
	"google.golang.org/api/iterator"
)

func isDuplicated(eventid string) (bool, error) {
	query := datastore.NewQuery(eventRecordTable).Filter("eventid =", eventid)
	query = query.Limit(1).KeysOnly()
	keys, err := client.GetAll(context.Background(), query, nil)
	if err != nil {
		return false, err
	}
	if len(keys) == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

func RecordEvent(eventid string) {
	eventKey := datastore.IncompleteKey(eventRecordTable, nil)
	receivedEvent := ReceivedEvent{Eventid: eventid}
	_, err := client.Put(context.Background(), eventKey, &receivedEvent)
	if err != nil {
		log.Fatalf("Error stroing data: %v", err)
	}
}

func PublishEvents() int {
	query := datastore.NewQuery(eventPublishTable).Filter("sent =", false)
	it := client.Run(context.Background(), query)

	i := 0
	for {
		var eventEntity EventEntity
		_, err := it.Next(&eventEntity)
		if err == iterator.Done {
			break
		}

		msg := &broker.Message{
			Header: map[string]string{
				"eventid": eventEntity.Eventid,
				"type":    eventEntity.Type,
			},
			Body: eventEntity.EventData,
		}
		err = brk.Publish(publishTopic, msg)
		if err != nil {
			log.Fatalf("Error publishing event: %v", err)
		}

		eventEntity.Sent = true
		_, err = client.Put(context.Background(), eventEntity.Key, &eventEntity)
		if err != nil {
			log.Fatalf("Error stroing data: %v", err)
		}
		i += 1
	}
	return i
}

func RegisterEvent(e interface{}, eventType string, tx *datastore.Transaction) error {
	jsonBytes, _ := json.Marshal(e)
	eventEntity := EventEntity{
		Eventid:   getUUID(),
		Type:      eventType,
		Sent:      false,
		EventData: jsonBytes,
	}
	eventEntityKey := datastore.IncompleteKey(eventPublishTable, nil)
	_, err := tx.Put(eventEntityKey, &eventEntity)
	return err
}
