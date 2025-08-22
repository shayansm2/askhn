package elasticsearch

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
	"github.com/shayansm2/temporallm/internal/config"
)

var once sync.Once
var client *elasticsearch.TypedClient

func GetClient() *elasticsearch.TypedClient {
	log.Println("Getting client", config.Load().ElasticsearchURL)
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

var ESIndexSchema = &types.TypeMapping{
	Properties: map[string]types.Property{
		"id":    types.IntegerNumberProperty{},
		"score": types.IntegerNumberProperty{},
		"title": types.TextProperty{},
		"text":  types.TextProperty{},
		"type":  types.KeywordProperty{},
	},
}

func CreateElasticsearchIndex() error {
	client := GetClient()
	req := client.Indices.Create(config.Load().ElasticsearchIndexName)
	req.Mappings(ESIndexSchema)
	res, err := req.Do(context.Background())
	if err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}
	if !res.Acknowledged {
		return fmt.Errorf("index creation not acknowledged")
	}
	return nil
}

type ESDocument struct {
	Id    int    `json:"id"`
	Score int    `json:"score"`
	Title string `json:"title"`
	Type  string `json:"type"`
	Text  string `json:"text"`
}

func IndexDocument(document ESDocument) error {
	_, err := GetClient().
		Index(config.Load().ElasticsearchIndexName).
		Id(strconv.Itoa(document.Id)).
		Request(document).
		Do(context.TODO())
	return err
}

func SearchDocuments(query string, size int) ([]ESDocument, error) {
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

	var result []ESDocument
	res, err := GetClient().Search().Index("hacker_news").Request(&req).Do(context.TODO())
	if err != nil {
		return result, err
	}
	for _, hit := range res.Hits.Hits {
		var doc ESDocument
		if err := json.Unmarshal(hit.Source_, &doc); err != nil {
			return result, err
		}
		result = append(result, doc)
	}
	return result, nil
}
