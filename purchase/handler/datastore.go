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

type entityID struct {
	name  string
	value string
	table string
}

func getStructByID(id entityID, dst interface{}, parent *datastore.Key) bool {
	query := datastore.NewQuery(ds.Kind[id.table])
	if parent != nil {
		query = query.Ancestor(parent)
	}
	query = query.Filter(id.name+" =", id.value).Limit(1)
	it := ds.Client.Run(context.Background(), query)
	_, err := it.Next(dst)
	if err == iterator.Done {
		return false
	}
	if err != nil {
		log.Fatalf("Error reading data: %v", err)
	}
	return true
}

func getStructsByParentKey(table string, dst interface{}, parent *datastore.Key) {
	query := datastore.NewQuery(ds.Kind[table]).Ancestor(parent)
	_, err := ds.Client.GetAll(context.Background(), query, dst)
	if err != nil {
		log.Fatalf("Error reading data: %v", err)
	}
	return
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
