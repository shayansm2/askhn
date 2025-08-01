package elasticsearch

import (
	"log"
	"sync"

	"github.com/elastic/go-elasticsearch/v8"
)

var once sync.Once
var client *elasticsearch.TypedClient

func GetClient() *elasticsearch.TypedClient {
	once.Do(func() {
		var err error
		client, err = elasticsearch.NewTypedClient(elasticsearch.Config{
			Addresses: []string{"http://localhost:9200"},
		})
		if err != nil {
			log.Fatalf("Error creating the client: %s", err)
		}
	})
	return client
}
