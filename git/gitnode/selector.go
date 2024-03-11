package gitnode

import (
	"context"
	"errors"
	"github.com/LeeZXin/zall/meta/modules/model/gitnodemd"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf-utils/localcache"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/xormstore"
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

var (
	nodeSelector    *localcache.SingleCacheEntry[map[string][2]selector]
	NodeNotFoundErr = errors.New("node not found")
)

func init() {
	nodeSelector, _ = localcache.NewSingleCacheEntry(func(ctx context.Context) (map[string][2]selector, error) {
		ctx, closer := xormstore.Context(ctx)
		defer closer.Close()
		nodes, err := gitnodemd.GetAll(ctx)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return nil, err
		}
		return listutil.CollectToMap(nodes, func(t gitnodemd.GitNodeDTO) (string, error) {
			return t.NodeId, nil
		}, func(t gitnodemd.GitNodeDTO) ([2]selector, error) {
			httpNodes := listutil.Shuffle(t.HttpHosts)
			sshNodes := listutil.Shuffle(t.SshHosts)
			httpSelector := newRoundRobinSelector(httpNodes)
			sshSelector := newRoundRobinSelector(sshNodes)
			return [2]selector{
				httpSelector, sshSelector,
			}, nil
		})
	}, 30*time.Second)
}

func PickHttpHost(ctx context.Context, nodeId string) (string, error) {
	return pickHost(ctx, nodeId, 0)
}

func PickSshHost(ctx context.Context, nodeId string) (string, error) {
	return pickHost(ctx, nodeId, 1)
}

func pickHost(ctx context.Context, nodeId string, index int) (string, error) {
	if nodeId == "" {
		return "", errors.New("empty nodeId")
	}
	ret, err := nodeSelector.LoadData(ctx)
	if err != nil {
		return "", err
	}
	slr, b := ret[nodeId]
	if !b {
		return "", NodeNotFoundErr
	}
	node := slr[index].Select()
	if err != nil {
		return "", NodeNotFoundErr
	}
	return node, nil
}
