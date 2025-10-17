package elasticsearch

import (
	"context"
	"encoding/json"

	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/textquerytype"
	"github.com/shayansm2/askhn/internal/config"
	"github.com/shayansm2/askhn/internal/embeding"
)

func TextSearch(query string, size int) ([]ESDocument, error) {
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
		Source_: []string{"id", "score", "title", "type", "text"},
	}
	return searchAndGetDocuments(req)
}

func VectorSearch(query string, size int) ([]ESDocument, error) {
	queryVector, err := embeding.Encode(query)
	if err != nil {
		return nil, err
	}
	req := search.Request{
		Size:    &size,
		Source_: []string{"id", "score", "title", "type", "text"},
		Knn: []types.KnnSearch{{
			Field:       "title_text_v",
			K:           &size,
			QueryVector: queryVector,
		}},
	}
	return searchAndGetDocuments(req)
}

func HybridSearch(query string, size int) ([]ESDocument, error) {
	queryVector, err := embeding.Encode(query)
	if err != nil {
		return nil, err
	}
	boost := float32(0.5)
	req := search.Request{
		Size:    &size,
		Source_: []string{"id", "score", "title", "type", "text"},
		Knn: []types.KnnSearch{{
			Field:       "title_text_v",
			K:           &size,
			QueryVector: queryVector,
			Boost:       &boost,
		}},
		Query: &types.Query{
			Bool: &types.BoolQuery{
				Must: []types.Query{
					{
						MultiMatch: &types.MultiMatchQuery{
							Query:  query,
							Fields: []string{"title^2", "text"},
							Type:   &textquerytype.Bestfields,
							Boost:  &boost,
						},
					},
				},
			},
		},
	}
	return searchAndGetDocuments(req)
}

func GetAllStories() ([]string, error) {
	cnf := config.Load()
	req := search.Request{
		Query: &types.Query{
			Bool: &types.BoolQuery{
				Filter: []types.Query{
					{
						Term: map[string]types.TermQuery{
							"type": {
								Value: "story",
							},
						},
					},
				},
			},
		},
	}

	var result []string
	res, err := GetClient().Search().Index(cnf.ElasticsearchIndexName).Request(&req).Do(context.TODO())
	if err != nil {
		return result, err
	}
	for _, hit := range res.Hits.Hits {
		var doc ESDocument
		if err := json.Unmarshal(hit.Source_, &doc); err != nil {
			return result, err
		}
		result = append(result, doc.Title)
	}
	return result, nil
}

func GetAllComments(title string) ([]string, error) {
	cnf := config.Load()
	req := search.Request{
		Query: &types.Query{
			Bool: &types.BoolQuery{
				Must: []types.Query{{
					MatchPhrase: map[string]types.MatchPhraseQuery{
						"title": {
							Query: title,
						},
					},
				}},
				Filter: []types.Query{{
					Term: map[string]types.TermQuery{
						"type": {
							Value: "comment",
						},
					},
				}},
			},
		},
	}

	var result []string
	res, err := GetClient().Search().Index(cnf.ElasticsearchIndexName).Request(&req).Do(context.TODO())
	if err != nil {
		return result, err
	}
	for _, hit := range res.Hits.Hits {
		var doc ESDocument
		if err := json.Unmarshal(hit.Source_, &doc); err != nil {
			return result, err
		}
		result = append(result, doc.Title)
	}
	return result, nil
}

func searchAndGetDocuments(request search.Request) ([]ESDocument, error) {
	cnf := config.Load()
	var result []ESDocument
	res, err := GetClient().Search().Index(cnf.ElasticsearchIndexName).Request(&request).Do(context.TODO())
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
