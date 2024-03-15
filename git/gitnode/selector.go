package gitnode

import (
	"context"
	"errors"
	"github.com/LeeZXin/zall/meta/modules/model/gitnodemd"
	"github.com/LeeZXin/zsf-utils/collections/hashset"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf-utils/quit"
	"github.com/LeeZXin/zsf-utils/taskutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"sync"
	"sync/atomic"
	"time"
)

type selector interface {
	Select() string
}

type roundRobinSelector struct {
	targets []string
	index   *atomic.Uint64
}

func (r *roundRobinSelector) Select() string {
	return r.targets[r.index.Add(1)%uint64(len(r.targets))]
}

func newRoundRobinSelector(targets []string) selector {
	if len(targets) == 0 {
		return new(emptyTargetsSelector)
	}
	index := atomic.Uint64{}
	index.Store(0)
	return &roundRobinSelector{
		targets: targets,
		index:   &index,
	}
}

type emptyTargetsSelector struct {
}

func (r *emptyTargetsSelector) Select() string {
	return ""
}

type nodeCache struct {
	versions  map[string]int64
	selectors map[string][2]selector
	sync.RWMutex
}

func (n *nodeCache) getVersion(nodeId string) (int64, bool) {
	n.RLock()
	defer n.RUnlock()
	ret, b := n.versions[nodeId]
	return ret, b
}

func (n *nodeCache) putCache(nodeId string, httpSelector, sshSelector selector, version int64) {
	n.Lock()
	defer n.Unlock()
	n.selectors[nodeId] = [2]selector{httpSelector, sshSelector}
	n.versions[nodeId] = version
}

func (n *nodeCache) delCache(nodeId string) {
	n.Lock()
	defer n.Unlock()
	delete(n.selectors, nodeId)
	delete(n.versions, nodeId)
}

func (n *nodeCache) allNodeIds() *hashset.HashSet[string] {
	n.RLock()
	defer n.RUnlock()
	ret := make([]string, 0, len(n.selectors))
	for nodeId := range n.selectors {
		ret = append(ret, nodeId)
	}
	return hashset.NewHashSet(ret...)
}

func (n *nodeCache) refreshCache() {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	nodes, err := gitnodemd.GetAll(ctx)
	if err != nil {
		logger.Logger.Error(err)
		return
	}
	tmp := make(map[string]gitnodemd.GitNodeDTO, len(nodes))
	for i := range nodes {
		tmp[nodes[i].NodeId] = nodes[i]
	}
	// 检查新增的或编辑过的
	for nodeId, node := range tmp {
		version, b := n.getVersion(nodeId)
		// 不存在或版本号不一致
		if !b || version != node.Version {
			httpNodes := listutil.Shuffle(node.HttpHosts)
			sshNodes := listutil.Shuffle(node.SshHosts)
			httpSelector := newRoundRobinSelector(httpNodes)
			sshSelector := newRoundRobinSelector(sshNodes)
			n.putCache(nodeId, httpSelector, sshSelector, node.Version)
		}
	}
	// 检查删除的
	n.allNodeIds().Range(func(nodeId string) {
		_, b := tmp[nodeId]
		if !b {
			n.delCache(nodeId)
		}
	})
}

func (n *nodeCache) getHttpSelector(nodeId string) (selector, bool) {
	n.RLock()
	defer n.RUnlock()
	selectors, b := n.selectors[nodeId]
	if !b {
		return nil, false
	}
	return selectors[0], true
}

func (n *nodeCache) getSshSelector(nodeId string) (selector, bool) {
	n.RLock()
	defer n.RUnlock()
	selectors, b := n.selectors[nodeId]
	if !b {
		return nil, false
	}
	return selectors[1], true
}

func newNodeCache() *nodeCache {
	ret := &nodeCache{
		selectors: make(map[string][2]selector),
		versions:  make(map[string]int64),
	}
	task, _ := taskutil.NewPeriodicalTask(10*time.Second, ret.refreshCache)
	task.Start()
	quit.AddShutdownHook(task.Stop)
	return ret
}

var (
	nodeSelector    = newNodeCache()
	NodeNotFoundErr = errors.New("node not found")
)

func PickHttpHost(nodeId string) (string, error) {
	s, b := nodeSelector.getHttpSelector(nodeId)
	if !b {
		return "", NodeNotFoundErr
	}
	ret := s.Select()
	if ret == "" {
		return "", NodeNotFoundErr
	}
	return ret, nil
}

func PickSshHost(nodeId string) (string, error) {
	s, b := nodeSelector.getSshSelector(nodeId)
	if !b {
		return "", NodeNotFoundErr
	}
	ret := s.Select()
	if ret == "" {
		return "", NodeNotFoundErr
	}
	return ret, nil
}
