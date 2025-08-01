package chatbot

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type HackerNewsResponse struct {
	Id    int    `json:"id"`
	Kids  []int  `json:"kids"`
	Score int    `json:"score"`
	Title string `json:"title"`
	Type  string `json:"type"`
	Text  string `json:"text"`
}

func (i *HackerNewsResponse) GetID() string {
	return strconv.Itoa(i.Id)
}

func HackerNewsItem(id int) (*HackerNewsResponse, error) {
	url := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%d.json", id)
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call Hacker News API: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	var result HackerNewsResponse
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}
	return &result, nil
}
