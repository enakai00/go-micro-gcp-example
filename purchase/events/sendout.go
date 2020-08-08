package events

import (
	"context"
	"os"

	"cloud.google.com/go/datastore"
	"github.com/micro/go-micro/v2/broker"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-plugins/broker/googlepubsub/v2"
	"google.golang.org/api/iterator"
)

var (
	topic     = os.Getenv("EVENT_PUBLISH_TOPIC")
	projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")
	client, _ = datastore.NewClient(context.Background(), projectID)
	brk       broker.Broker
)

func init() {
	brk = googlepubsub.NewBroker(googlepubsub.ProjectID(projectID))
	err := brk.Connect()
	if err != nil {
		log.Fatalf("Broker Connect error: %v", err)
	}
}

func Sendout() int {
	query := datastore.NewQuery("Event").Filter("sent =", false)
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
				"type:":   eventEntity.Type,
			},
			Body: eventEntity.EventData,
		}
		err = brk.Publish(topic, msg)
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
