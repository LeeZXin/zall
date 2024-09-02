package cfgsrv

import (
	"encoding/base64"
	"encoding/json"
	"github.com/LeeZXin/zall/pkg/login"
	"github.com/LeeZXin/zall/util"
	"net/url"
	"regexp"
	"time"
)

var (
	DefaultSysCfg = &SysCfg{
		DisableSelfRegisterUser: false,
		AllowUserCreateTeam:     true,
	}
)

type SysCfg struct {
	// 禁用用户自行注册账号
	DisableSelfRegisterUser bool `json:"disableSelfRegisterUser"`
	// 允许用户(非管理员)自行创建项目组
	AllowUserCreateTeam bool `json:"allowUserCreateTeam"`
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
	// HttpUrl smart http url
	HttpUrl string `json:"httpUrl"`
	// SshUrl ssh url
	SshUrl string `json:"sshUrl"`
	// LfsJwtExpiry lfs jwt过期时间 单位秒
	LfsJwtExpiry int64 `json:"lfsJwtExpiry"`
	// LfsJwtSecret lfs 密钥
	LfsJwtSecret string `json:"lfsJwtSecret"`
	// lfsJwtSecretBytes lfs 密钥
	lfsJwtSecretBytes []byte
}

func (c *GitCfg) IsValid() bool {
	httpUrl, err := url.Parse(c.HttpUrl)
	if err != nil {
		return false
	}
	if httpUrl.Scheme != "http" && httpUrl.Scheme != "https" {
		return false
	}
	sshUrl, err := url.Parse(c.SshUrl)
	if err != nil {
		return false
	}
	if sshUrl.Scheme != "ssh" && sshUrl.User.Username() != "git" {
		return false
	}
	if c.LfsJwtExpiry <= 0 {
		return false
	}
	if c.LfsJwtSecret == "" || len(c.LfsJwtSecret) > 1024 {
		return false
	}
	return true
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
		c.LfsJwtExpiry = 600
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

func (c *EnvCfg) IsValid() bool {
	if len(c.Envs) == 0 {
		return false
	}
	pattern := regexp.MustCompile(`^[a-zA-Z]{1,16}$`)
	for _, env := range c.Envs {
		if !pattern.MatchString(env) {
			return false
		}
	}
	return true
}

func (c *EnvCfg) FromStore(val string) error {
	return json.Unmarshal([]byte(val), c)
}

type GitRepoServerCfg struct {
	HttpHost string `json:"httpHost"`
	SshHost  string `json:"sshHost"`
}

func (c *GitRepoServerCfg) IsValid() bool {
	if !util.GenIpPortPattern().MatchString(c.HttpHost) {
		return false
	}
	if !util.GenIpPortPattern().MatchString(c.SshHost) {
		return false
	}
	return true
}

func (*GitRepoServerCfg) Key() string {
	return "git_repo_server_cfg"
}

func (c *GitRepoServerCfg) Val() string {
	ret, _ := json.Marshal(c)
	return string(ret)
}

func (c *GitRepoServerCfg) FromStore(val string) error {
	return json.Unmarshal([]byte(val), c)
}

type LoginCfg struct {
	login.Cfg
}

func (*LoginCfg) Key() string {
	return "login"
}

func (c *LoginCfg) Val() string {
	ret, _ := json.Marshal(c)
	return string(ret)
}

func (c *LoginCfg) FromStore(val string) error {
	return json.Unmarshal([]byte(val), c)
}
