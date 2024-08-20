package http

import (
	"context"
	"fmt"
	"github.com/LeeZXin/zsf/services/discovery"
	"github.com/LeeZXin/zsf/services/lb"
	"net/http"
	"net/url"
	"strings"
)

type Task struct {
	Url         string            `json:"url"`
	Headers     map[string]string `json:"headers"`
	Method      string            `json:"method"`
	BodyStr     string            `json:"bodyStr"`
	ContentType string            `json:"contentType"`
}

func (t *Task) IsValid() bool {
	parsedUrl, err := url.Parse(t.Url)
	if err != nil || !strings.HasPrefix(parsedUrl.Scheme, "http") {
		return false
	}
	if t.Method == "" {
		return false
	}
	return true
}

func (t *Task) DoRequest(httpClient *http.Client) error {
	httpUrl := t.Url
	parsedUrl, err := url.Parse(httpUrl)
	if err != nil {
		return err
	}
	// 走服务发现
	if strings.HasSuffix(parsedUrl.Host, "-http") {
		servers, err := discovery.Discover(context.Background(), parsedUrl.Host)
		if err != nil {
			return err
		}
		if len(servers) == 0 {
			return lb.ServerNotFound
		}
		server := discovery.ChooseRandomServer(servers)
		httpUrl = fmt.Sprintf("%s://%s:%d/%s", parsedUrl.Scheme, server.Host, server.Port, parsedUrl.RequestURI())
	}
	req, err := http.NewRequest(t.Method, httpUrl, strings.NewReader(t.BodyStr))
	if err != nil {
		return fmt.Errorf("http request failed: %v", err)
	}
	if t.ContentType != "" {
		req.Header.Set("Content-Type", t.ContentType)
	}
	for k, v := range t.Headers {
		req.Header.Set(k, v)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("http request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("http request return code: %v", resp.StatusCode)
	}
	return nil
}
