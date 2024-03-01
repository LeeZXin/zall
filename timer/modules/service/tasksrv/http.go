package tasksrv

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/httputil"
	"github.com/LeeZXin/zsf/services/discovery"
	"net/http"
	"net/url"
	"strings"
)

var (
	httpClient = httputil.NewRetryableHttpClient()
)

type TaskObj struct {
	TaskType string
	Content  string
}

type HttpTask struct {
	Url      string            `json:"url"`
	Headers  map[string]string `json:"headers"`
	Method   string            `json:"method"`
	PostJson string            `json:"postJson"`
}

func (t *HttpTask) IsValid() bool {
	_, err := url.Parse(t.Url)
	if err != nil {
		return false
	}
	if !util.FindInSlice([]string{"GET", "POST"}, t.Method) {
		return false
	}
	return true
}

func handleHttpTask(content string, sb *util.SimpleLogger) bool {
	var task HttpTask
	err := json.Unmarshal([]byte(content), &task)
	if err != nil {
		sb.WriteString(fmt.Sprintf("unmarshal task content err: %v", err))
		return false
	}
	rawUrl := task.Url
	parsedUrl, err := url.Parse(task.Url)
	if err != nil {
		sb.WriteString(fmt.Sprintf("invalid http url err: %v", task.Url))
		return false
	}
	b := strings.HasSuffix(parsedUrl.Host, "-http")
	if b {
		host, err := discovery.PickOneHost(context.Background(), parsedUrl.Host)
		if err != nil {
			sb.WriteString(fmt.Sprintf("can not find service: %s err: %v", parsedUrl.Host, task.Url))
			return false
		}
		rawUrl = parsedUrl.Scheme + "://" + host + parsedUrl.RequestURI()
	}
	sb.WriteString(fmt.Sprintf("do http task url: %s method: %s", rawUrl, task.Method))
	switch task.Method {
	case "GET":
		req, err := http.NewRequest("GET", rawUrl, nil)
		if err != nil {
			sb.WriteString(fmt.Sprintf("http request failed: %v", err))
			return false
		}
		for k, v := range task.Headers {
			req.Header.Set(k, v)
		}
		resp, err := httpClient.Do(req)
		if err != nil {
			sb.WriteString(fmt.Sprintf("http request failed: %v", err))
			return false
		}
		sb.WriteString(fmt.Sprintf("http request return code: %v", resp.StatusCode))
	case "POST":
		req, err := http.NewRequest("POST", rawUrl, strings.NewReader(task.PostJson))
		if err != nil {
			sb.WriteString(fmt.Sprintf("http request failed: %v", err))
			return false
		}
		req.Header.Set("Content-Type", httputil.JsonContentType)
		for k, v := range task.Headers {
			req.Header.Set(k, v)
		}
		resp, err := httpClient.Do(req)
		if err != nil {
			sb.WriteString(fmt.Sprintf("http request failed: %v", err))
			return false
		}
		sb.WriteString(fmt.Sprintf("http request return code: %v", resp.StatusCode))
	default:
		sb.WriteString(fmt.Sprintf("unsupported method: %s", task.Method))
		return false
	}
	return true
}
