package deploy

import (
	"encoding/json"
	"github.com/LeeZXin/zall/util"
	"net/url"
	"strings"
)

type DetectType int

const (
	TcpDetectType DetectType = iota
	HttpDetectType
)

func (t DetectType) IsValid() bool {
	switch t {
	case TcpDetectType, HttpDetectType:
		return true
	default:
		return false
	}
}

type HttpDetect struct {
	Url     string            `json:"url"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers"`
}

func (c *HttpDetect) IsValid() bool {
	parsedUrl, err := url.Parse(c.Url)
	return err == nil && strings.HasPrefix(parsedUrl.Scheme, "http")
}

type TcpDetect struct {
	Ip   string `json:"ip"`
	Port int    `json:"port"`
}

func (c *TcpDetect) IsValid() bool {
	return util.IpPattern.MatchString(c.Ip) && c.Port > 0
}

type DetectConfig struct {
	DetectType DetectType  `json:"detectType"`
	HttpDetect *HttpDetect `json:"httpDetect"`
	TcpDetect  *TcpDetect  `json:"tcpDetect"`
}

func (c *DetectConfig) IsValid() bool {
	switch c.DetectType {
	case TcpDetectType:
		return c.TcpDetect != nil && c.TcpDetect.IsValid()
	case HttpDetectType:
		return c.HttpDetect != nil && c.HttpDetect.IsValid()
	default:
		return false
	}
}

type ProcessConfig struct {
	Host         string       `json:"host"`
	AgentUrl     string       `json:"agentUrl"`
	AgentToken   string       `json:"agentToken"`
	SshUrl       string       `json:"sshUrl"`
	SshPassword  string       `json:"sshPassword"`
	DetectConfig DetectConfig `json:"detectConfig"`
	DeployScript string       `json:"deployScript"`
}

func (c *ProcessConfig) IsValid() bool {
	b := util.IpPattern.MatchString(c.Host)
	if !b {
		return false
	}
	parsedUrl, err := url.Parse(c.AgentUrl)
	if err != nil || !strings.HasPrefix(parsedUrl.Scheme, "http") {
		return false
	}
	if len(c.AgentToken) > 32 {
		return false
	}
	if !c.DetectConfig.IsValid() {
		return false
	}
	return true
}

type K8sConfig struct {
}

func (c *K8sConfig) IsValid() bool {
	return true
}

type Config struct {
	ProcessConfigList []ProcessConfig `json:"processConfigList"`
	K8sConfigList     []K8sConfig     `json:"k8SConfigList"`
}

func (c *Config) IsValid() bool {
	if len(c.ProcessConfigList) > 100 || len(c.K8sConfigList) > 100 {
		return false
	}
	for _, cfg := range c.ProcessConfigList {
		if !cfg.IsValid() {
			return false
		}
	}
	for _, cfg := range c.K8sConfigList {
		if !cfg.IsValid() {
			return false
		}
	}
	return true
}

func (c *Config) FromDB(content []byte) error {
	if c == nil {
		*c = Config{}
	}
	return json.Unmarshal(content, c)
}

func (c *Config) ToDB() ([]byte, error) {
	return json.Marshal(c)
}
