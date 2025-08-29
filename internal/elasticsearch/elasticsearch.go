package elasticsearch

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/textquerytype"
	"github.com/shayansm2/askhn/internal/config"
)

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

func GetAllStories() ([]string, error) {
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
	res, err := GetClient().Search().Index("hacker_news").Request(&req).Do(context.TODO())
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
	res, err := GetClient().Search().Index("hacker_news").Request(&req).Do(context.TODO())
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

// func getFieldsFromResult(res *search.Response, field string) ([]string, error) {
// 	var result []string
// 	for _, hit := range res.Hits.Hits {
// 		var doc ESDocument
// 		if err := json.Unmarshal(hit.Source_, &doc); err != nil {
// 			return result, err
// 		}
// 		result = append(result, doc.Title)
// 	}
// 	return result, nil
// }
