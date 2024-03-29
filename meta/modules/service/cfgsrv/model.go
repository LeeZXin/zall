package cfgsrv

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/idutil"
	"github.com/LeeZXin/zsf/common"
	"time"
)

var (
	DefaultSysCfg = &SysCfg{
		DisableSelfRegisterUser: false,
		AllowUserCreateTeam:     true,
		DisableSmartHttp:        false,
	}
	DefaultGitCfg = &GitCfg{
		AppUrl:       fmt.Sprintf("http://%s:%d", common.GetLocalIP(), common.HttpServerPort()),
		LfsEnabled:   true,
		LfsJwtExpiry: 3600,
		LfsJwtSecret: idutil.RandomUuid(),
		ActionToken:  idutil.RandomUuid(),
		RepoToken:    idutil.RandomUuid(),
	}
	DefaultEnvCfg = &EnvCfg{
		Envs: []string{
			"prd",
		},
	}
)

type SysCfg struct {
	// 禁用用户自行注册账号
	DisableSelfRegisterUser bool `json:"disableSelfRegisterUser"`
	// 允许用户(非管理员)自行创建项目组
	AllowUserCreateTeam bool `json:"allowUserCreateTeam"`
	// 禁用http smart协议
	DisableSmartHttp bool `json:"disableSmartHttp"`
}

func (c *SysCfg) Key() string {
	return "sys_cfg"
}

func (c *SysCfg) Val() string {
	ret, _ := json.Marshal(c)
	return string(ret)
}

func (c *SysCfg) FromStore(val string) error {
	return json.Unmarshal([]byte(val), c)
}

type GitCfg struct {
	// AppUrl 域名地址
	AppUrl string `json:"appUrl"`
	// LfsEnabled 是否启用lfs
	LfsEnabled bool `json:"lfsEnabled"`
	// LfsJwtExpiry lfs jwt过期时间 单位秒
	LfsJwtExpiry int64 `json:"lfsJwtExpiry"`
	// LfsJwtSecret lfs 密钥
	LfsJwtSecret string `json:"lfsJwtSecret"`
	// lfsJwtSecretBytes lfs 密钥
	lfsJwtSecretBytes []byte
	// ActionToken action 执行token
	ActionToken string `json:"actionToken"`
	// RepoToken git服务token
	RepoToken string `json:"repoToken"`
}

func (c *GitCfg) GetLfsJwtExpiry() time.Duration {
	return time.Duration(c.LfsJwtExpiry) * time.Second
}

func (c *GitCfg) GetLfsJwtSecretBytes() []byte {
	return c.lfsJwtSecretBytes
}

func (c *GitCfg) Key() string {
	return "git_cfg"
}

func (c *GitCfg) Val() string {
	ret, _ := json.Marshal(c)
	return string(ret)
}

func (c *GitCfg) FromStore(val string) error {
	err := json.Unmarshal([]byte(val), c)
	if err != nil {
		return err
	}
	if c.LfsJwtExpiry == 0 {
		c.LfsJwtExpiry = 3600
	}
	c.lfsJwtSecretBytes = make([]byte, 32)
	n, err := base64.RawURLEncoding.Decode(c.lfsJwtSecretBytes, []byte(c.LfsJwtSecret))
	if err != nil || n != 32 {
		c.lfsJwtSecretBytes, err = util.NewRandomJwtSecret()
		if err != nil {
			return err
		}
	}
	return nil
}

type EnvCfg struct {
	Envs []string `json:"envs"`
}

func (c *EnvCfg) Key() string {
	return "env_cfg"
}

func (c *EnvCfg) Val() string {
	ret, _ := json.Marshal(c)
	return string(ret)
}

func (c *EnvCfg) FromStore(val string) error {
	return json.Unmarshal([]byte(val), c)
}
