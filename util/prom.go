package util

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
)

func NewPromHttpClient(address string, client *http.Client) (*PromHttpClient, error) {
	u, err := url.Parse(address)
	if err != nil {
		return nil, err
	}
	u.Path = strings.TrimRight(u.Path, "/")

	return &PromHttpClient{
		endpoint: u,
		client:   client,
	}, nil
}

type PromHttpClient struct {
	endpoint *url.URL
	client   *http.Client
}

func (c *PromHttpClient) URL(ep string, args map[string]string) *url.URL {
	p := path.Join(c.endpoint.Path, ep)

	for arg, val := range args {
		arg = ":" + arg
		p = strings.Replace(p, arg, val, -1)
	}

	u := *c.endpoint
	u.Path = p

	return &u
}

func (c *PromHttpClient) Do(ctx context.Context, req *http.Request) (*http.Response, []byte, error) {
	if ctx != nil {
		req = req.WithContext(ctx)
	}
	resp, err := c.client.Do(req)
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()

	if err != nil {
		return nil, nil, err
	}

	var body []byte
	done := make(chan struct{})
	go func() {
		body, err = io.ReadAll(resp.Body)
		close(done)
	}()

	select {
	case <-ctx.Done():
		<-done
		err = resp.Body.Close()
		if err == nil {
			err = ctx.Err()
		}
	case <-done:
	}

	return resp, body, err
}
