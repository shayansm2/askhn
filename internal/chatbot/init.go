package chatbot

import (
	"context"
	"fmt"

	"github.com/shayansm2/temporallm/internal/utils/elasticsearch"
)

func CreateElasticsearchIndex() error {
	client := elasticsearch.GetClient()
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
