package util

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/LeeZXin/zsf/services/discovery"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Api struct {
	Headers     map[string]string `json:"headers"`
	Url         string            `json:"url"`
	Method      string            `json:"method"`
	ContentType string            `json:"contentType"`
	BodyStr     string            `json:"bodyStr"`
}

func (a *Api) IsValid() bool {
	_, err := url.Parse(a.Url)
	if err != nil {
		return false
	}
	if a.Method == "" {
		return false
	}
	return true
}

func (a *Api) replaceStr(args map[string]string, str string) string {
	for k, v := range args {
		str = strings.ReplaceAll(str, "{{"+k+"}}", v)
	}
	return str
}

func (a *Api) formatQuery(query map[string]string) string {
	if query == nil || len(query) == 0 {
		return ""
	}
	sb := strings.Builder{}
	// a=1&b=2&c=3
	for k, v := range query {
		sb.WriteString(k)
		sb.WriteString("=")
		sb.WriteString(base64.StdEncoding.EncodeToString([]byte(v)))
		sb.WriteString("&")
	}
	return sb.String()[:sb.Len()-1]
}

func (a *Api) DoRequest(httpClient *http.Client, query, args map[string]string) (map[string]any, error) {
	finalUrl := a.Url
	parseUrl, err := url.Parse(finalUrl)
	if err != nil {
		return nil, err
	}
	if strings.HasSuffix(parseUrl.Host, "-http") {
		servers, err := discovery.Discover(context.Background(), parseUrl.Host)
		if err != nil {
			return nil, err
		}
		if len(servers) == 0 {
			return nil, fmt.Errorf("host: %v not found", parseUrl.Host)
		}
		server := discovery.ChooseRandomServer(servers)
		finalUrl = parseUrl.Scheme + "://" + fmt.Sprintf("%s:%d", server.Host, server.Port) + parseUrl.RequestURI()
	}
	if len(query) > 0 {
		if parseUrl.RawQuery != "" {
			finalUrl = finalUrl + "&" + a.formatQuery(query)
		} else {
			finalUrl = finalUrl + "?" + a.formatQuery(query)
		}
	}
	bodyStr := a.replaceStr(args, a.BodyStr)
	request, err := http.NewRequest(a.Method, finalUrl, strings.NewReader(bodyStr))
	if err != nil {
		return nil, err
	}
	if a.ContentType != "" {
		request.Header.Set("Content-Type", a.ContentType)
	}
	for k, v := range a.Headers {
		request.Header.Set(k, v)
	}
	response, err := httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("url: %v http failed with err: %v", finalUrl, err)
	}
	defer response.Body.Close()
	if response.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("url: %v http response return code: %v", finalUrl, response.StatusCode)
	}
	// 限制返回字节10MB大小
	bodyAll, err := io.ReadAll(io.LimitReader(response.Body, 10*1024*1024))
	if err != nil {
		return nil, fmt.Errorf("url: %v  read http response body err: %v", finalUrl, err)
	}
	ret := make(map[string]any)
	// 忽略json异常
	err = json.Unmarshal(bodyAll, &ret)
	if err != nil {
		return nil, fmt.Errorf("url: %v do http request failed with Err: %v", finalUrl, err)
	}
	return ret, nil
}
