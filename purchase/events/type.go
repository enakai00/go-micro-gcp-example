package events

import (
	"cloud.google.com/go/datastore"
)

type EventEntity struct {
	Eventid   string         `datastore:"eventid"`
	Type      string         `datastore:"type"`
	Sent      bool           `datastore:"sent"`
	EventData []byte         `datastore:"event_data"`
	Key       *datastore.Key `datastore:"__key__"`
}
