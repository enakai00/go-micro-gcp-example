package handler

import (
	"context"

	"cloud.google.com/go/datastore"
	"github.com/google/uuid"
	log "github.com/micro/go-micro/v2/logger"
	"google.golang.org/api/iterator"

	"github.com/enakai00/go-micro-gcp-example/purchase/ds"
)

func getUUID() string {
	uuid, _ := uuid.NewRandom()
	return uuid.String()
}

func getStructByID(idName, idValue, table string, dst interface{}) bool {
	query := datastore.NewQuery(ds.Kind[table]).Filter(idName+" =", idValue).Limit(1)
	it := ds.Client.Run(context.Background(), query)
	var err error
	switch dst := dst.(type) {
	case *ds.Cart:
		_, err = it.Next(dst)
	case *ds.OrderTicket:
		_, err = it.Next(dst)
	}

	if err == iterator.Done {
		return false
	}
	if err != nil {
		log.Fatalf("Error reading data: %v", err)
	}
	return true
}

func createEntity(table string, src interface{},
	parent *datastore.Key, tx *datastore.Transaction) error {
	key := datastore.IncompleteKey(ds.Kind[table], parent)
	var err error
	if tx == nil {
		_, err = ds.Client.Put(context.Background(), key, src)
	} else {
		_, err = tx.Put(key, src)
	}
	return err
}
