package events

import (
	"context"

	"cloud.google.com/go/datastore"
	"github.com/micro/go-micro/v2/broker"
	log "github.com/micro/go-micro/v2/logger"
	"google.golang.org/api/iterator"
)

func RecordEvent(table string, eventid string) {
	eventKey := datastore.IncompleteKey(table, nil)
	receivedEvent := ReceivedEvent{Eventid: eventid}
	_, err := client.Put(context.Background(), eventKey, &receivedEvent)
	if err != nil {
		log.Fatalf("Error stroing data: %v", err)
	}
}

func Sendout(table string) int {
	query := datastore.NewQuery(table).Filter("sent =", false)
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
