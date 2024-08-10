package approval

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/collections/hashset"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/logger"
)

var (
	// methodMap 存储方法
	methodMap = make(map[string]func(*FlowContext) (map[string]any, error))
)

type FlowContext struct {
	context.Context
	FlowId  int64
	BizId   string
	Kvs     Kvs
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

type KvCfg struct {
	Key      string `json:"key"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Required bool   `json:"required"`
}

type Kvs []Kv

func (k *Kvs) FromDB(content []byte) error {
	if k == nil {
		*k = make(Kvs, 0)
	}
	return json.Unmarshal(content, k)
}

func (k *Kvs) ToDB() ([]byte, error) {
	return json.Marshal(k)
}

func (k *Kvs) ToMap() map[string]string {
	if len(*k) == 0 {
		return map[string]string{}
	}
	ret := make(map[string]string)
	for _, kv := range *k {
		ret[kv.Key] = kv.Value
	}
	return ret
}

type Kv struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type NodeCfg struct {
	NodeType   NodeType              `json:"nodeType"`
	Api        *util.Api             `json:"api"`
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
	Api        *util.Api          `json:"api"`
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
	KvCfgs []KvCfg  `json:"kvCfgs"`
	Node   *NodeCfg `json:"node"`
}

func (c *ProcessCfg) IsValid() bool {
	if c.Node == nil || len(c.KvCfgs) > 1000 {
		return false
	}
	return c.Node.IsValid()
}

func (c *ProcessCfg) Convert() Process {
	ret := Process{}
	ret.KvCfgs = c.KvCfgs
	if c.Node != nil {
		ret.Node = c.Node.Convert()
	}
	return ret
}

type Process struct {
	KvCfgs []KvCfg `json:"kvCfgs"`
	Node   *Node   `json:"node"`
}

func (p *Process) Find(nodeId int) *Node {
	if p.Node == nil {
		return nil
	}
	return p.Node.Find(nodeId)
}

func (p *Process) CheckKvCfgs(kvs Kvs) []string {
	errKeys := make([]string, 0)
	for _, cfg := range p.KvCfgs {
		find := false
		for _, kv := range kvs {
			if cfg.Key == kv.Key {
				find = true
				if cfg.Required && kv.Value == "" {
					errKeys = append(errKeys, kv.Key)
				}
				break
			}
		}
		if cfg.Required && !find {
			errKeys = append(errKeys, cfg.Key)
		}
	}
	return errKeys
}

func (p *Process) FindAndDo(nodeId int, fnMap map[NodeType]func(*Node)) {
	if p.Node == nil {
		return
	}
	p.Node.FindAndDo(nodeId, fnMap)
}

func (p *Process) FromDB(content []byte) error {
	if p == nil {
		*p = Process{}
	}
	return json.Unmarshal(content, p)
}

func (p *Process) ToDB() ([]byte, error) {
	return json.Marshal(p)
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
