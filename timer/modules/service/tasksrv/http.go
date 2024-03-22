package tasksrv

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/httputil"
	"github.com/LeeZXin/zsf/services/discovery"
	"github.com/LeeZXin/zsf/services/lb"
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
	Url         string            `json:"url"`
	Headers     map[string]string `json:"headers"`
	Method      string            `json:"method"`
	BodyStr     string            `json:"bodyStr"`
	ContentType string            `json:"contentType"`
	Zones       []string          `json:"zones"`
}

func (t *HttpTask) IsValid() bool {
	_, err := url.Parse(t.Url)
	if err != nil {
		return false
	}
	if t.Method == "" {
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
	finalUrl := task.Url
	parsedUrl, err := url.Parse(task.Url)
	if err != nil {
		sb.WriteString(fmt.Sprintf("invalid http url err: %v", task.Url))
		return false
	}
	b := strings.HasSuffix(parsedUrl.Host, "-http")
	if b {
		// 跨数据中心请求
		if len(task.Zones) > 0 {
			zoneRet := true
			for _, zone := range task.Zones {
				servers, err := discovery.DiscoverWithZone(context.Background(), zone, parsedUrl.Host)
				if err != nil {
					zoneRet = false
					sb.WriteString(fmt.Sprintf("can not find service: %s with zone: %s err: %v", parsedUrl.Host, zone, err))
					continue
				}
				if len(servers) == 0 {
					zoneRet = false
					sb.WriteString(fmt.Sprintf("can not find service: %s with zone: %s err: %v", parsedUrl.Host, zone, lb.ServerNotFound))
					continue
				}
				server := discovery.ChooseRandomServer(servers)
				finalUrl = parsedUrl.Scheme + "://" + fmt.Sprintf("%s:%d", server.Host, server.Port) + parsedUrl.RequestURI()
				zoneRet = zoneRet && doHttpRequest(sb, finalUrl, &task)
				sb.WriteString("--------- end zone: " + zone)
			}
			return zoneRet
		} else {
			servers, err := discovery.Discover(context.Background(), parsedUrl.Host)
			if err != nil {
				sb.WriteString(fmt.Sprintf("can not find service: %s err: %v", parsedUrl.Host, err))
				return false
			}
			if len(servers) == 0 {
				sb.WriteString(fmt.Sprintf("can not find service: %s err: %v", parsedUrl.Host, lb.ServerNotFound))
				return false
			}
			server := discovery.ChooseRandomServer(servers)
			finalUrl = parsedUrl.Scheme + "://" + fmt.Sprintf("%s:%d", server.Host, server.Port) + parsedUrl.RequestURI()
		}
	}
	return doHttpRequest(sb, finalUrl, &task)
}

func doHttpRequest(sb *util.SimpleLogger, rawUrl string, task *HttpTask) bool {
	sb.WriteString(fmt.Sprintf("do http task url: %s method: %s", rawUrl, task.Method))
	req, err := http.NewRequest(task.Method, rawUrl, strings.NewReader(task.BodyStr))
	if err != nil {
		sb.WriteString(fmt.Sprintf("http request failed: %v", err))
		return false
	}
	if task.ContentType != "" {
		req.Header.Set("Content-Type", task.ContentType)
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
	return true
}
