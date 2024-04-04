package deploy

import (
	"encoding/json"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/util"
	"net/url"
	"strings"
)

type ServiceType int

const (
	ProcessServiceType ServiceType = iota + 1
	K8sServiceType
	DockerSwarmType
)

func (t ServiceType) IsValid() bool {
	switch t {
	case ProcessServiceType, K8sServiceType, DockerSwarmType:
		return true
	default:
		return false
	}
}

func (t ServiceType) Readable() string {
	switch t {
	case ProcessServiceType:
		return i18n.GetByKey(i18n.DeployProcessServiceType)
	case K8sServiceType:
		return i18n.GetByKey(i18n.DeployK8sServiceType)
	case DockerSwarmType:
		return i18n.GetByKey(i18n.DeployDockerSwarmType)
	default:
		return i18n.GetByKey(i18n.DeployUnknownServiceType)
	}
}

type DetectType int

const (
	TcpDetectType DetectType = iota + 1
	HttpGetDetectType
)

func (t DetectType) IsValid() bool {
	switch t {
	case TcpDetectType, HttpGetDetectType:
		return true
	default:
		return false
	}
}

type HttpGetDetect struct {
	Url string `json:"url"`
}

func (c *HttpGetDetect) IsValid() bool {
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
	DetectType    DetectType     `json:"detectType"`
	HttpGetDetect *HttpGetDetect `json:"httpGetDetect,omitempty"`
	TcpDetect     *TcpDetect     `json:"tcpDetect,omitempty"`
}

func (c *DetectConfig) IsValid() bool {
	switch c.DetectType {
	case TcpDetectType:
		return c.TcpDetect != nil && c.TcpDetect.IsValid()
	case HttpGetDetectType:
		return c.HttpGetDetect != nil && c.HttpGetDetect.IsValid()
	default:
		return false
	}
}

func (c *DetectConfig) FromDB(content []byte) error {
	if c == nil {
		*c = DetectConfig{}
	}
	return json.Unmarshal(content, c)
}

func (c *DetectConfig) ToDB() ([]byte, error) {
	return json.Marshal(c)
}

type ProcessConfig struct {
	Host          string       `json:"host"`
	AgentHost     string       `json:"agentHost"`
	AgentToken    string       `json:"agentToken"`
	SshHost       string       `json:"sshHost"`
	SshPassword   string       `json:"sshPassword"`
	DetectConfig  DetectConfig `json:"detectConfig"`
	DeployScript  string       `json:"deployScript"`
	StopScript    string       `json:"stopScript"`
	RestartScript string       `json:"restartScript"`
}

func (c *ProcessConfig) IsValid() bool {
	b := util.IpPattern.MatchString(c.Host)
	if !b {
		return false
	}
	if !util.IpPortPattern.MatchString(c.AgentHost) {
		return false
	}
	if len(c.AgentToken) > 32 {
		return false
	}
	if !c.DetectConfig.IsValid() {
		return false
	}
	if c.DeployScript == "" {
		return false
	}
	if c.StopScript == "" {
		return false
	}
	if c.RestartScript == "" {
		return false
	}
	return true
}

type K8sConfig struct {
}

func (c *K8sConfig) IsValid() bool {
	return true
}

type NormalConfig struct {
	// DeployPlanApprovers 发布计划审批人
	DeployPlanApprovers []string `json:"deployPlanApprovers"`
	// DisallowDeployServiceWithoutPlan  禁止不通过发布计划部署服务
	DisallowDeployServiceWithoutPlan bool `json:"disallowDeployServiceWithoutPlan"`
}

func (c *NormalConfig) FromDB(content []byte) error {
	if c == nil {
		*c = NormalConfig{}
	}
	return json.Unmarshal(content, c)
}

func (c *NormalConfig) ToDB() ([]byte, error) {
	return json.Marshal(c)
}
