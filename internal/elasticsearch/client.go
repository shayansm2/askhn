package elasticsearch

import (
	"log"
	"sync"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/shayansm2/temporallm/internal/config"
)

var once sync.Once
var client *elasticsearch.TypedClient

func GetClient() *elasticsearch.TypedClient {
	once.Do(func() {
		var err error
		client, err = elasticsearch.NewTypedClient(elasticsearch.Config{
			Addresses: []string{config.Load().ElasticsearchURL},
			Username:  config.Load().ElasticsearchUser,
			Password:  config.Load().ElasticsearchPass,
		})
		if err != nil {
			log.Fatalf("Error creating the client: %s", err)
		}
	})
	return client
}
