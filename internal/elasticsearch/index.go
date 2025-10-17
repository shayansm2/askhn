package elasticsearch

import (
	"context"
	"fmt"
	"strconv"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/densevectorsimilarity"
	"github.com/shayansm2/askhn/internal/config"
	"github.com/shayansm2/askhn/internal/embeding"
)

var dim int = 1024
var index bool = true
var ESIndexSchema = &types.TypeMapping{
	Properties: map[string]types.Property{
		"id":    types.IntegerNumberProperty{},
		"score": types.IntegerNumberProperty{},
		"title": types.TextProperty{},
		"text":  types.TextProperty{},
		"type":  types.KeywordProperty{},
		"title_v": types.DenseVectorProperty{
			Type:       "dense_vector",
			Dims:       &dim,
			Index:      &index,
			Similarity: &densevectorsimilarity.Cosine,
		},
		"text_v": types.DenseVectorProperty{
			Type:       "dense_vector",
			Dims:       &dim,
			Index:      &index,
			Similarity: &densevectorsimilarity.Cosine,
		},
		"title_text_v": types.DenseVectorProperty{
			Type:       "dense_vector",
			Dims:       &dim,
			Index:      &index,
			Similarity: &densevectorsimilarity.Cosine,
		},
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

type indexedDocument struct {
	ESDocument
	VectorTitle     []float32 `json:"title_v,omitempty"`
	VectorText      []float32 `json:"text_v,omitempty"`
	VectorTitleText []float32 `json:"title_text_v,omitempty"`
}

func IndexDocument(document ESDocument) error {
	req := indexedDocument{ESDocument: document}
	if document.Title != "" {
		if embeding, err := embeding.Encode(document.Title); err == nil {
			req.VectorTitle = embeding
		} else {
			return fmt.Errorf("embeding error: %w", err)
		}
	}
	if document.Text != "" {
		if embeding, err := embeding.Encode(document.Text); err == nil {
			req.VectorText = embeding
		} else {
			return fmt.Errorf("embeding error: %w", err)
		}
	}
	if document.Text != "" || document.Title != "" {
		if embeding, err := embeding.Encode(document.Title + " " + document.Text); err == nil {
			req.VectorTitleText = embeding
		} else {
			return fmt.Errorf("embeding error: %w", err)
		}
	}
	_, err := GetClient().
		Index(config.Load().ElasticsearchIndexName).
		Id(strconv.Itoa(document.Id)).
		Request(req).
		Do(context.TODO())
	return err
}
