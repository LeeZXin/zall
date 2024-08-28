package prom

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
)

func NewPromHttpClient(address string, client *http.Client) (*HttpClient, error) {
	u, err := url.Parse(address)
	if err != nil {
		return nil, err
	}
	u.Path = strings.TrimRight(u.Path, "/")
	return &HttpClient{
		endpoint: u,
		client:   client,
	}, nil
}

type HttpClient struct {
	endpoint *url.URL
	client   *http.Client
}

func (c *HttpClient) URL(ep string, args map[string]string) *url.URL {
	p := path.Join(c.endpoint.Path, ep)
	for arg, val := range args {
		p = strings.Replace(p, ":"+arg, val, -1)
	}
	u := *c.endpoint
	u.Path = p
	return &u
}

func (c *HttpClient) Do(ctx context.Context, req *http.Request) (*http.Response, []byte, error) {
	if ctx != nil {
		req = req.WithContext(ctx)
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	return resp, body, err
}
