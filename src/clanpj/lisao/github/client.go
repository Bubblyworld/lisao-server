package github

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
)

var apiHost = "https://api.github.com/"
var ErrInternalServerError = errors.New("github: internal server error")

// Client provides utilities for querying refs/tags in the given repo.
type Client struct {
	user   string
	repo   string
	client *http.Client
}

func NewClient(user, repo string) *Client {
	return &Client{
		repo:   repo,
		user:   user,
		client: &http.Client{},
	}
}

func (c *Client) NewRequest(method, apiUrl string, params url.Values) (*http.Request, error) {
	body := strings.NewReader(params.Encode())
	url := apiHost + strings.Trim(apiUrl, "/")
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == 500 {
		return nil, ErrInternalServerError
	}

	return res, nil
}

func (c *Client) repoApiUrl() string {
	return "/repos/" + c.user + "/" + c.repo
}
