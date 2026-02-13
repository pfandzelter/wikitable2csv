package wiki

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type WikiResponse struct {
	Parse WikiParse `json:"parse"`
}

type WikiParse struct {
	PageID int      `json:"pageid"`
	Title  string   `json:"title"`
	Text   WikiText `json:"text"`
}

type WikiText struct {
	Content string `json:"*"`
}

func Fetch(queryUrl string, userAgent string) (*WikiResponse, error) {

	if queryUrl == "" {
		return nil, fmt.Errorf("empty url given")
	}

	req, err := http.NewRequest("GET", queryUrl, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("User-Agent", userAgent)

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var wikiRes WikiResponse
	if err := json.NewDecoder(resp.Body).Decode(&wikiRes); err != nil {
		return nil, err
	}

	return &wikiRes, nil
}
