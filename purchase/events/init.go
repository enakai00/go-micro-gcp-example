package events

import (
	"context"
	"os"
	"strings"

	"cloud.google.com/go/datastore"
	"github.com/micro/go-micro/v2/broker"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-plugins/broker/googlepubsub/v2"
)

type EventEntity struct {
	Eventid   string         `datastore:"eventid"`
	Type      string         `datastore:"type"`
	Sent      bool           `datastore:"sent"`
	EventData []byte         `datastore:"event_data"`
	Key       *datastore.Key `datastore:"__key__"`
}

type ReceivedEvent struct {
	Eventid string `datastore:"eventid"`
}

var (
	publishTopic      = os.Getenv("EVENT_PUBLISH_TOPIC")
	subscribeTopics   = os.Getenv("EVENT_SUBSCRIBE_TOPICS")
	projectID         = os.Getenv("GOOGLE_CLOUD_PROJECT")
	EventPublishTable = "PurchasePublishEvent"
	EventRecordTable  = "PurchaseReceivedEvent"
	client, _         = datastore.NewClient(context.Background(), projectID)
	brk               broker.Broker
)

func init() {
	brk = googlepubsub.NewBroker(googlepubsub.ProjectID(projectID))
	err := brk.Connect()
	if err != nil {
		log.Fatalf("Broker Connect error: %v", err)
	}

	for _, subscribeTopic := range strings.Split(subscribeTopics, ",") {
		_, err := brk.Subscribe(subscribeTopic, eventHandler,
			broker.Queue("purchase-service"),
			broker.DisableAutoAck())
		if err != nil {
			log.Fatalf("Broker Subscribe error: %v", err)
		}
	}
}
