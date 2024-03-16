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
	methodMap = make(map[string]func(*FlowContext) (map[string]any, error))
)

type FlowContext struct {
	context.Context
	FlowId  int64
	BizId   string
	Process *Process
}

// RegisterMethod 注册方法 no thread-safe
func RegisterMethod(name string, fn func(*FlowContext) (map[string]any, error)) {
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

func (a *Api) DoRequest(ctx *FlowContext) (map[string]any, error) {
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
		if parseUrl.RawQuery != "" {
			finalUrl = finalUrl + fmt.Sprintf("&bizId=%s", ctx.BizId)
		} else {
			finalUrl = finalUrl + fmt.Sprintf("?bizId=%s", ctx.BizId)
		}
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

func (m *Method) DoMethod(ctx *FlowContext) (map[string]any, error) {
	fn, b := methodMap[m.Name]
	if !b {
		return nil, fmt.Errorf("unknown process method: %s", m.Name)
	}
	return fn(ctx)
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

type ConditionalNodeCfg struct {
	Node      *NodeCfg `json:"node"`
	Condition string   `json:"condition"`
}

type Kv struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

type NodeCfg struct {
	NodeType   NodeType              `json:"nodeType"`
	Api        *Api                  `json:"api"`
	Method     *Method               `json:"method"`
	Accounts   []string              `json:"accounts"`
	AtLeastNum int                   `json:"atLeastNum"`
	Title      string                `json:"title"`
	Content    string                `json:"content"`
	Vars       map[string]string     `json:"vars"`
	Next       []*ConditionalNodeCfg `json:"next"`
}

func (n *NodeCfg) IsValid() bool {
	if n.Title == "" {
		return false
	}
	switch n.NodeType {
	case PeopleNode:
		accounts := hashset.NewHashSet(n.Accounts...)
		if accounts.Size() > 100 || accounts.Size() == 0 || len(n.Accounts) != accounts.Size() || n.AtLeastNum < 1 || n.AtLeastNum > accounts.Size() {
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

func (n *NodeCfg) Convert() *Node {
	var c nodeConverter
	return c.ConvertNode(n)
}

type nodeConverter struct {
	nodeId int
}

func (c *nodeConverter) ConvertNode(n *NodeCfg) *Node {
	c.nodeId += 1
	ret := &Node{
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
	ret.Next, _ = listutil.Map(n.Next, func(t *ConditionalNodeCfg) (*ConditionalNode, error) {
		return c.ConvertConditionalNode(t), nil
	})
	return ret
}

func (c *nodeConverter) ConvertConditionalNode(n *ConditionalNodeCfg) *ConditionalNode {
	return &ConditionalNode{
		Node:      c.ConvertNode(n.Node),
		Condition: n.Condition,
	}
}

type Node struct {
	NodeId     int                `json:"nodeId"`
	NodeType   NodeType           `json:"nodeType"`
	Api        *Api               `json:"api"`
	Method     *Method            `json:"method"`
	Accounts   []string           `json:"accounts"`
	AtLeastNum int                `json:"atLeastNum"`
	Title      string             `json:"title"`
	Content    string             `json:"content"`
	Vars       map[string]string  `json:"vars"`
	Next       []*ConditionalNode `json:"next"`
}

type ConditionalNode struct {
	Node      *Node  `json:"node"`
	Condition string `json:"condition"`
}

type ProcessCfg struct {
	Kvs  []Kv     `json:"kvs"`
	Node *NodeCfg `json:"node"`
}

func (c *ProcessCfg) IsValid() bool {
	if c.Node == nil {
		return false
	}
	return c.Node.IsValid()
}

func (c *ProcessCfg) Convert() Process {
	ret := Process{}
	ret.Kvs = c.Kvs
	if c.Node != nil {
		ret.Node = c.Node.Convert()
	}
	return ret
}

type Process struct {
	Kvs  []Kv  `json:"kvs"`
	Node *Node `json:"node"`
}

func (p *Process) Find(nodeId int) *Node {
	if p.Node == nil {
		return nil
	}
	return p.Node.Find(nodeId)
}

func (p *Process) FindAndDo(nodeId int, fnMap map[NodeType]func(*Node)) {
	if p.Node == nil {
		return
	}
	p.Node.FindAndDo(nodeId, fnMap)
}

func findByNodeId(node *Node, nodeId int) *Node {
	if node == nil {
		return nil
	}
	if node.NodeId == nodeId {
		return node
	}
	for _, next := range node.Next {
		p := findByNodeId(next.Node, nodeId)
		if p != nil {
			return p
		}
	}
	return nil
}

func (n *Node) Find(nodeId int) *Node {
	return findByNodeId(n, nodeId)
}

func (n *Node) FindAndDo(nodeId int, fnMap map[NodeType]func(*Node)) {
	node := findByNodeId(n, nodeId)
	if node == nil {
		return
	}
	fn, b := fnMap[node.NodeType]
	if !b {
		return
	}
	fn(node)
}
