package chatbot

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/textquerytype"
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

func CreateElasticsearchIndex() error {
	client := GetClient()
	req := client.Indices.Create(ElasticsearchIndexName)
	req.Mappings(ElasticsearchIndexSchema)
	res, err := req.Do(context.Background())
	if err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}
	if !res.Acknowledged {
		return fmt.Errorf("index creation not acknowledged")
	}
	return nil
}

func IndexDocument(document ElasticSearchDocument) error {
	_, err := GetClient().
		Index(ElasticsearchIndexName).
		Id(strconv.Itoa(document.Id)).
		Request(document).
		Do(context.TODO())
	return err
}

func SearchDocuments(query string, size int) ([]ElasticSearchDocument, error) {
	req := search.Request{
		Size: &size,
		Query: &types.Query{
			Bool: &types.BoolQuery{
				Must: []types.Query{
					{
						MultiMatch: &types.MultiMatchQuery{
							Query:  query,
							Fields: []string{"title^2", "text"},
							Type:   &textquerytype.Bestfields,
						},
					},
				},
			},
		},
	}

	var result []ElasticSearchDocument
	res, err := GetClient().Search().Index("hacker_news").Request(&req).Do(context.TODO())
	if err != nil {
		return result, err
	}
	for _, hit := range res.Hits.Hits {
		var doc ElasticSearchDocument
		if err := json.Unmarshal(hit.Source_, &doc); err != nil {
			return result, err
		}
		result = append(result, doc)
	}
	return result, nil
}
