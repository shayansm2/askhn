package chatbot

import (
	"strconv"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

type HackerNewsResponse struct {
	Id    int    `json:"id"`
	Kids  []int  `json:"kids"`
	Score int    `json:"score"`
	Title string `json:"title"`
	Type  string `json:"type"`
	Text  string `json:"text"`
}

type ElasticSearchDocument struct {
	Id    int    `json:"id"`
	Score int    `json:"score"`
	Title string `json:"title"`
	Type  string `json:"type"`
	Text  string `json:"text"`
}

func (i *ElasticSearchDocument) GetID() string {
	return strconv.Itoa(i.Id)
}

var ElasticsearchIndexSchema = &types.TypeMapping{
	Properties: map[string]types.Property{
		"id":    types.IntegerNumberProperty{},
		"score": types.IntegerNumberProperty{},
		"title": types.TextProperty{},
		"text":  types.TextProperty{},
		"type":  types.KeywordProperty{},
	},
}
