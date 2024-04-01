package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/LeeZXin/zall/promagent/modules/model/prommd"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/collections/hashset"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf-utils/quit"
	"github.com/LeeZXin/zsf-utils/taskutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/property/static"
	"github.com/LeeZXin/zsf/services/discovery"
	"github.com/LeeZXin/zsf/services/lb"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"hash/crc32"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

var (
	watchTask   *taskutil.PeriodicalTask
	serverUrl   string
	agentEnv    string
	lastHashSum uint32
	filePath    string
)

func StartAgent() {
	agentEnv = static.GetString("prom.agent.env")
	if agentEnv == "" {
		logger.Logger.Fatal("empty prom agent env")
	}
	serverUrl = static.GetString("prom.agent.serverUrl")
	if serverUrl == "" {
		logger.Logger.Fatal("empty prom agent serverUrl")
	}
	filePath = static.GetString("prom.agent.filePath")
	if filePath == "" || !filepath.IsAbs(filePath) {
		logger.Logger.Fatalf("wrong prom.agent.filesd.path: %v", filePath)
	}
	logger.Logger.Infof("prom agent started with env: %v serverUrl: %v filePath: %v", agentEnv, serverUrl, filePath)
	watchTask, _ = taskutil.NewPeriodicalTask(5*time.Second, updateFileSd)
	watchTask.Start()
	quit.AddShutdownHook(watchTask.Stop, true)
}

func updateFileSd() {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	scrapes, err := prommd.GetAllScrapeByServerUrl(ctx, serverUrl, agentEnv)
	if err != nil {
		logger.Logger.Error(err)
		return
	}
	fileContent := packFileContent(scrapes)
	hashSum := crc32.ChecksumIEEE(fileContent)
	if hashSum != lastHashSum {
		// save file
		err = util.WriteFile(filePath, fileContent)
		if err != nil {
			logger.Logger.Error(err)
		}
		lastHashSum = hashSum
	}
}

type Config struct {
	Targets []string          `json:"targets,omitempty"`
	Labels  map[string]string `json:"labels,omitempty"`
}

func packFileContent(scrapes []prommd.Scrape) []byte {
	if len(scrapes) == 0 {
		return []byte{}
	}
	group := make(map[string][]prommd.Scrape)
	for _, scrape := range scrapes {
		list, b := group[scrape.AppId]
		if !b {
			list = make([]prommd.Scrape, 0)
		}
		list = append(list, scrape)
		group[scrape.AppId] = list
	}
	configs := make([]Config, 0)
	ctx := context.Background()
	for appId, list := range group {
		appTargets := make([]string, 0)
		for _, scrape := range list {
			switch scrape.TargetType {
			case prommd.DiscoveryTargetType:
				servers, err := discovery.Discover(ctx, scrape.Target)
				if err != nil {
					if err != lb.ServerNotFound {
						logger.Logger.Error(err)
					}
				} else if len(servers) > 0 {
					targets, _ := listutil.Map(servers, func(t lb.Server) (string, error) {
						return fmt.Sprintf("%s:%d", t.Host, 16054), nil
					})
					appTargets = append(appTargets, targets...)
				}
			case prommd.HostTargetType:
				if scrape.Target != "" {
					appTargets = append(appTargets, scrape.Target)
				}
			}
		}
		if len(appTargets) > 0 {
			configs = append(configs, NewConfig(appTargets, appId))
		}
	}
	ret, _ := json.MarshalIndent(configs, "", "\t")
	return ret
}

func NewConfig(targets []string, appId string) Config {
	// 去重
	targets = hashset.NewHashSet(targets...).AllKeys()
	sort.SliceStable(targets, func(i, j int) bool {
		return strings.Compare(targets[i], targets[j]) > 0
	})
	return Config{
		Targets: targets,
		Labels: map[string]string{
			"app_id": appId,
		},
	}
}
