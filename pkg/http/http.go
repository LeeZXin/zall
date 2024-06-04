package http

import (
	"context"
	"fmt"
	"github.com/LeeZXin/zsf/services/discovery"
	"github.com/LeeZXin/zsf/services/lb"
	"go.uber.org/multierr"
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
	Zones       []string          `json:"zones"`
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
	parsedUrl, err := url.Parse(t.Url)
	if err != nil {
		return fmt.Errorf("invalid http url err: %v", t.Url)
	}
	b := strings.HasSuffix(parsedUrl.Host, "-http")
	if !b {
		// 不需要走服务发现
		return t.doHttpRequest(httpClient, t.Url)
	}
	handleDiscovery := func(servers []lb.Server, err error) error {
		if err != nil {
			return fmt.Errorf("can not find service: %s err: %v", parsedUrl.Host, err)
		}
		if len(servers) == 0 {
			return fmt.Errorf("can not find service: %s", parsedUrl.Host)
		}
		server := discovery.ChooseRandomServer(servers)
		finalUrl := parsedUrl.Scheme + "://" + fmt.Sprintf("%s:%d", server.Host, server.Port) + parsedUrl.RequestURI()
		return t.doHttpRequest(httpClient, finalUrl)
	}
	// 跨数据中心请求
	if len(t.Zones) > 0 {
		zoneErr := make([]error, 0)
		for _, zone := range t.Zones {
			servers, err := discovery.DiscoverWithZone(context.Background(), zone, parsedUrl.Host)
			err = handleDiscovery(servers, err)
			if err != nil {
				zoneErr = append(zoneErr, err)
			}
		}
		return multierr.Combine(zoneErr...)
	}
	servers, err := discovery.Discover(context.Background(), parsedUrl.Host)
	return handleDiscovery(servers, err)
}

func (t *Task) doHttpRequest(httpClient *http.Client, rawUrl string) error {
	req, err := http.NewRequest(t.Method, rawUrl, strings.NewReader(t.BodyStr))
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
