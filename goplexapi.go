package goplexapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type PlexClient struct {
	BaseURL string
	Token   string
	Client  *http.Client
}

func NewPlexClient(baseURL, token string) *PlexClient {
	return &PlexClient{
		BaseURL: baseURL,
		Token:   token,
		Client:  &http.Client{},
	}
}

func (p *PlexClient) makeRequest(method, endpoint string, payload interface{}) ([]byte, error) {
	url := fmt.Sprintf("%s%s", p.BaseURL, endpoint)
	var req *http.Request
	var err error

	if method == "POST" {
		req, err = http.NewRequest(method, url, strings.NewReader(payload.(string)))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Plex-Token", p.Token)
	resp, err := p.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			panic(err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (p *PlexClient) GetLibrarySections() ([]map[string]interface{}, error) {
	data, err := p.makeRequest("GET", "/library/sections", nil)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result["MediaContainer"].(map[string]interface{})["Directory"].([]map[string]interface{}), nil
}

func (p *PlexClient) GetCurrentPlayingSong(playerClientName string) (string, error) {
	data, err := p.makeRequest("GET", "/status/sessions", nil)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return "", err
	}

	sessions := result["MediaContainer"].(map[string]interface{})["Metadata"].([]interface{})

	for _, session := range sessions {
		sessionMap := session.(map[string]interface{})
		client := sessionMap["Player"].(map[string]interface{})["title"].(string)
		if client == playerClientName {
			track := sessionMap["title"].(string)
			return track, nil
		}
	}

	return "", fmt.Errorf("No song currently playing on %s", playerClientName)
}
