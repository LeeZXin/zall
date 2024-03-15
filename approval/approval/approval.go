package approval

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/LeeZXin/zsf-utils/collections/hashset"
	"github.com/LeeZXin/zsf-utils/httputil"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/services/discovery"
	"github.com/LeeZXin/zsf/services/lb"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
)

var (
	httpClient = httputil.NewRetryableHttpClient()
	// methodMap 存储方法
	methodMap = make(map[string]func() (map[string]any, error))
)

// RegisterMethod 注册方法 no thread-safe
func RegisterMethod(name string, fn func() (map[string]any, error)) {
	if fn == nil {
		logger.Logger.Fatalf("nil process method func: %s", name)
	}
	_, b := methodMap[name]
	if b {
		logger.Logger.Fatalf("duplicated process method name: %s", name)
	}
	methodMap[name] = fn
}

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

func (a *Api) DoRequest() (map[string]any, error) {
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
			return nil, lb.ServerNotFound
		}
		server := servers[rand.Int()%len(servers)]
		finalUrl = parseUrl.Scheme + "://" + fmt.Sprintf("%s:%d", server.Host, server.Port) + parseUrl.RequestURI()
	}
	request, err := http.NewRequest(a.Method, finalUrl, strings.NewReader(a.BodyStr))
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
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("http response return code: %v", response.StatusCode)
	}
	// 限制返回字节3MB大小
	bodyAll, err := io.ReadAll(io.LimitReader(response.Body, 3*1024*1024))
	if err != nil {
		return nil, fmt.Errorf("read http response body err: %v", err)
	}
	ret := make(map[string]any)
	// 忽略json异常
	err = json.Unmarshal(bodyAll, &ret)
	if err != nil {
		logger.Logger.Error(err)
	}
	return ret, nil
}

type Method struct {
	Name string `json:"name"`
}

func (m *Method) DoMethod() (map[string]any, error) {
	fn, b := methodMap[m.Name]
	if !b {
		return nil, fmt.Errorf("unknown process method: %s", m.Name)
	}
	return fn()
}

func (m *Method) IsValid() bool {
	if m.Name == "" {
		return false
	}
	_, b := methodMap[m.Name]
	return b
}

type NodeType int

const (
	PeopleNode NodeType = iota + 1
	ApiNode
	MethodNode
	DisagreeNode
	AgreeNode
)

type ConditionNode struct {
	Node      *Node  `json:"node"`
	Condition string `json:"condition"`
}

type Kv struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

type Node struct {
	NodeType   NodeType          `json:"nodeType"`
	Api        *Api              `json:"api"`
	Method     *Method           `json:"method"`
	Accounts   []string          `json:"accounts"`
	AtLeastNum int               `json:"atLeastNum"`
	Title      string            `json:"title"`
	Content    string            `json:"content"`
	Vars       map[string]string `json:"vars"`
	Next       []*ConditionNode  `json:"next"`
}

func (n *Node) IsValid() bool {
	if n.Title == "" {
		return false
	}
	switch n.NodeType {
	case PeopleNode:
		accounts := hashset.NewHashSet(n.Accounts...)
		if accounts.Size() == 0 || len(n.Accounts) != accounts.Size() || n.AtLeastNum < 1 || n.AtLeastNum > accounts.Size() {
			return false
		}
	case ApiNode:
		if n.Api == nil || !n.Api.IsValid() {
			return false
		}
	case MethodNode:
		if n.Method == nil || !n.Method.IsValid() {
			return false
		}
	case DisagreeNode:
		if n.Content == "" {
			return false
		}
	case AgreeNode:
		if n.Content == "" {
			return false
		}
	default:
		return false
	}
	for _, cn := range n.Next {
		if cn.Node == nil || !cn.Node.IsValid() {
			return false
		}
	}
	return true
}

func (n *Node) ToApproval() *Approval {
	var c convert
	return c.ToApproval(n)
}

type convert struct {
	nodeId int
}

func (c *convert) ToApproval(n *Node) *Approval {
	c.nodeId += 1
	ret := &Approval{
		NodeId:     c.nodeId,
		NodeType:   n.NodeType,
		Api:        n.Api,
		Method:     n.Method,
		Accounts:   n.Accounts,
		AtLeastNum: n.AtLeastNum,
		Title:      n.Title,
		Content:    n.Content,
		Vars:       n.Vars,
	}
	ret.Next, _ = listutil.Map(n.Next, func(t *ConditionNode) (*ConditionApproval, error) {
		return c.ToConditionApproval(t), nil
	})
	return ret
}

func (c *convert) ToConditionApproval(n *ConditionNode) *ConditionApproval {
	return &ConditionApproval{
		Node:      c.ToApproval(n.Node),
		Condition: n.Condition,
	}
}

type Approval struct {
	NodeId     int                  `json:"nodeId"`
	NodeType   NodeType             `json:"nodeType"`
	Api        *Api                 `json:"api"`
	Method     *Method              `json:"method"`
	Accounts   []string             `json:"accounts"`
	AtLeastNum int                  `json:"atLeastNum"`
	Title      string               `json:"title"`
	Content    string               `json:"content"`
	Vars       map[string]string    `json:"vars"`
	Next       []*ConditionApproval `json:"next"`
}

type ConditionApproval struct {
	Node      *Approval `json:"node"`
	Condition string    `json:"condition"`
}

func findByNodeId(a *Approval, nodeId int) *Approval {
	if a.NodeId == nodeId {
		return a
	}
	for _, n := range a.Next {
		p := findByNodeId(n.Node, nodeId)
		if p != nil {
			return p
		}
	}
	return nil
}
func (a *Approval) Find(nodeId int) *Approval {
	return findByNodeId(a, nodeId)
}

func (a *Approval) FindAndDo(nodeId int, fnMap map[NodeType]func(*Approval)) {
	n := findByNodeId(a, nodeId)
	if n == nil {
		return
	}
	fn, b := fnMap[a.NodeType]
	if !b {
		return
	}
	fn(n)
}
