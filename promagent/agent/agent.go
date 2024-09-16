package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/LeeZXin/zall/promagent/modules/model/prommd"
	"github.com/LeeZXin/zall/util"
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
	endpoint    string
	lastHashSum uint32
	filePath    string
	env         string
)

func StartAgent() {
	env = static.GetString("prom.agent.env")
	if env == "" {
		logger.Logger.Fatal("prom agent started with empty env")
	}
	endpoint = static.GetString("prom.agent.endpoint")
	if endpoint == "" {
		logger.Logger.Fatal("empty prom agent endpoint")
	}
	filePath = static.GetString("prom.agent.filePath")
	if filePath == "" || !filepath.IsAbs(filePath) {
		logger.Logger.Fatalf("wrong prom.agent.filesd.path: %v", filePath)
	}
	logger.Logger.Infof("prom agent started with endpoint: %v filePath: %v env: %v", endpoint, filePath, env)
	stopFunc, _ := taskutil.RunPeriodicalTask(0, 30*time.Second, updateFileSd)
	quit.AddShutdownHook(quit.ShutdownHook(stopFunc), true)
}

func updateFileSd(context.Context) {
	ctx, closer := xormstore.Context(context.Background())
	defer closer.Close()
	scrapes, err := prommd.GetAllScrape(ctx, prommd.GetAllScrapeReqDTO{
		Endpoint: endpoint,
		Env:      env,
		Cols:     []string{"app_id", "target", "target_type"},
	})
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
					targets := listutil.MapNe(servers, func(t lb.Server) string {
						return fmt.Sprintf("%s:%d", t.Host, t.Port)
					})
					appTargets = append(appTargets, targets...)
				}
			case prommd.HostTargetType:
				if scrape.Target != "" {
					targets := strings.Split(scrape.Target, ";")
					targets = listutil.FilterNe(targets, func(t string) bool {
						return len(t) > 0
					})
					appTargets = append(appTargets, targets...)
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
	targets = listutil.Distinct(targets...)
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
