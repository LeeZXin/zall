package gitnode

import (
	"context"
	"errors"
	"github.com/LeeZXin/zall/meta/modules/model/gitnodemd"
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
	cache map[string][2]selector
	sync.RWMutex
}

func (n *nodeCache) refreshCache() {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	nodes, err := gitnodemd.GetAll(ctx)
	if err != nil {
		logger.Logger.Error(err)
		return
	}
	n.Lock()
	defer n.Unlock()
	c := make(map[string][2]selector, 8)
	for _, node := range nodes {
		httpNodes := listutil.Shuffle(node.HttpHosts)
		sshNodes := listutil.Shuffle(node.SshHosts)
		httpSelector := newRoundRobinSelector(httpNodes)
		sshSelector := newRoundRobinSelector(sshNodes)
		c[node.NodeId] = [2]selector{
			httpSelector, sshSelector,
		}
	}
	n.cache = c
}

func (n *nodeCache) getHttpSelector(nodeId string) (selector, bool) {
	n.RLock()
	defer n.RUnlock()
	if n.cache == nil {
		return nil, false
	}
	selectors, b := n.cache[nodeId]
	if !b {
		return nil, false
	}
	return selectors[0], true
}

func (n *nodeCache) getSshSelector(nodeId string) (selector, bool) {
	n.RLock()
	defer n.RUnlock()
	if n.cache == nil {
		return nil, false
	}
	selectors, b := n.cache[nodeId]
	if !b {
		return nil, false
	}
	return selectors[1], true
}

func newNodeCache() *nodeCache {
	ret := new(nodeCache)
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
