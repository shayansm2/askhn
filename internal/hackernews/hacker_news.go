package hackernews

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Item struct {
	Id    int    `json:"id"`
	Kids  []int  `json:"kids"`
	Score int    `json:"score"`
	Title string `json:"title"`
	Type  string `json:"type"`
	Text  string `json:"text"`
}

func GetItem(id int) (*Item, error) {
	url := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%d.json", id)
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call Hacker News API: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	var result Item
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}
	return &result, nil
}

type SearchResult struct {
	Hits []struct {
		ObjectID    string `json:"objectID"`
		Title       string `json:"title"`
		URL         string `json:"url"`
		Points      int    `json:"points"`
		NumComments int    `json:"num_comments"`
	} `json:"hits"`
}

func Search(query string) ([]int, error) {
	baseURL := "https://hn.algolia.com/api/v1/search"

	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	params := url.Values{}
	params.Add("query", query)
	params.Add("tags", "story")
	params.Add("hitsPerPage", "20")
	u.RawQuery = params.Encode()

	response, err := http.Get(u.String())
	if err != nil {
		return nil, fmt.Errorf("failed to call search API: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	var result SearchResult
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	var ids []int
	for _, hit := range result.Hits {
		var id int
		if _, err := fmt.Sscanf(hit.ObjectID, "%d", &id); err == nil {
			ids = append(ids, id)
		}
	}

	return ids, nil
}
