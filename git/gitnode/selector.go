package gitnode

import (
	"context"
	"errors"
	"github.com/LeeZXin/zall/meta/modules/model/gitnodemd"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf-utils/localcache"
	"github.com/LeeZXin/zsf-utils/selector"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/xorm/mysqlstore"
	"time"
)

var (
	nodeSelector    *localcache.SingleCacheEntry[map[string][2]selector.Selector[string]]
	NodeNotFoundErr = errors.New("node not found")
)

func init() {
	nodeSelector, _ = localcache.NewSingleCacheEntry(func(ctx context.Context) (map[string][2]selector.Selector[string], error) {
		ctx, closer := mysqlstore.Context(ctx)
		defer closer.Close()
		nodes, err := gitnodemd.GetAll(ctx)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return nil, err
		}
		return listutil.CollectToMap(nodes, func(t gitnodemd.GitNodeDTO) (string, error) {
			return t.NodeId, nil
		}, func(t gitnodemd.GitNodeDTO) ([2]selector.Selector[string], error) {
			httpNodes, _ := listutil.Map(listutil.Shuffle(t.HttpHosts), func(t string) (selector.Node[string], error) {
				return selector.Node[string]{
					Data: t,
				}, nil
			})
			sshNodes, _ := listutil.Map(listutil.Shuffle(t.SshHosts), func(t string) (selector.Node[string], error) {
				return selector.Node[string]{
					Data: t,
				}, nil
			})
			httpSelector := selector.NewRoundRobinSelector(httpNodes)
			sshSelector := selector.NewRoundRobinSelector(sshNodes)
			return [2]selector.Selector[string]{
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
	node, err := slr[index].Select()
	if err != nil {
		return "", NodeNotFoundErr
	}
	return node.Data, err
}
