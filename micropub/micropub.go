package micropub

import (
	"encoding/json"
	// "fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type HTTPError struct {
	resp *http.Response
}

func (e *HTTPError) Error() string {
	return e.resp.Status
}

type Item struct {
	Type       string `json:"type"`
	Properties struct {
		Name       []string    `json:"name"`
		Content    []string    `json:"content"`
		Photo      []string    `json:"photo"`
		PostStatus []string    `json:"post-status"`
		Published  []time.Time `json:"published"`
		UID        []uint64    `json:"uid"`
		URL        []string    `json:"url"`
	} `json:"properties"`
}

type Client struct {
	Endpoint string
	Token    string
}

func NewClient(endpoint, token string) *Client {
	return &Client{Endpoint: endpoint, Token: token}
}

type Config struct {
	Destination []struct {
		MicroblogAudio bool   `json:"microblog-audio"`
		Name           string `json:"name"`
		UID            string `json:"uid"`
	} `json:"destination"`
	MediaEndpoint string `json:"media-endpoint"`
	PostTypes     []struct {
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"post-types"`
}

func (c *Client) GetConfig() (*Config, error) {
	config := Config{}
	if err := c.get("?q=config", &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func (c *Client) GetCategories() ([]string, error) {
	var resp struct {
		Categories []string `json:"categories"`
	}
	if err := c.get("?q=category", &resp); err != nil {
		return nil, err
	}
	return resp.Categories, nil
}

func (c *Client) GetPosts() ([]*Item, error) {
	var resp struct {
		Items []*Item `json:"items"`
	}
	if err := c.get("?q=source", &resp); err != nil {
		return nil, err
	}
	return resp.Items, nil
}

func (c *Client) get(path string, dest interface{}) error {
	h := &http.Client{}

	log.Info("micropub: GET /micropub" + path)

	// FIXME: URL-encode path
	req, _ := http.NewRequest(http.MethodGet, c.Endpoint+path, nil)
	req.Header.Set("Authorization", "Bearer "+c.Token)

	resp, err := h.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return &HTTPError{resp: resp}
	}

	if err := json.NewDecoder(resp.Body).Decode(dest); err != nil {
		return err
	}

	return nil
}
