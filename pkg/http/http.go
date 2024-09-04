package http

import (
	"fmt"
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

func (t *Task) DoRequest() error {
	req, err := http.NewRequest(t.Method, t.Url, strings.NewReader(t.BodyStr))
	if err != nil {
		return fmt.Errorf("http request failed: %v", err)
	}
	if t.ContentType != "" {
		req.Header.Set("Content-Type", t.ContentType)
	}
	for k, v := range t.Headers {
		req.Header.Set(k, v)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("http request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http request return code: %v", resp.StatusCode)
	}
	return nil
}
