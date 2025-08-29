package elasticsearch

import (
	"context"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/shayansm2/askhn/internal/config"
)

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
