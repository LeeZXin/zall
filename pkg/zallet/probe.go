package zallet

import (
	"encoding/json"
	"net/url"
	"regexp"
	"strings"
)

var (
	IpPortPattern = regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}:\d+$`)
)

type ProbeType string

const (
	HttpProbeType ProbeType = "http"
	TcpProbeType  ProbeType = "tcp"
)

type TcpProbe struct {
	Host string `json:"host" yaml:"host"`
}

func (t *TcpProbe) IsValid() bool {
	return IpPortPattern.MatchString(t.Host)
}

type HttpProbe struct {
	Url string `json:"url" yaml:"url"`
}

func (t *HttpProbe) IsValid() bool {
	parsed, err := url.Parse(t.Url)
	if err != nil {
		return false
	}
	return strings.HasPrefix(parsed.Scheme, "http")
}

type ProbeFail struct {
	Times  int    `json:"times" yaml:"times"`
	Action string `json:"action" yaml:"action"`
}

func (f *ProbeFail) IsValid() bool {
	return f.Times > 0 && len(f.Action) > 0
}

type Probe struct {
	Type   ProbeType  `json:"type" yaml:"type"`
	Tcp    *TcpProbe  `json:"tcp,omitempty" yaml:"tcp,omitempty"`
	Http   *HttpProbe `json:"http,omitempty" yaml:"http,omitempty"`
	OnFail *ProbeFail `json:"onFail,omitempty" yaml:"onFail,omitempty"`
}

func (p *Probe) IsValid() bool {
	if p.OnFail != nil && !p.OnFail.IsValid() {
		return false
	}
	switch p.Type {
	case HttpProbeType:
		return p.Http != nil && p.Http.IsValid()
	case TcpProbeType:
		return p.Tcp != nil && p.Tcp.IsValid()
	default:
		return false
	}
}

func (p *Probe) FromDB(content []byte) error {
	if p == nil {
		*p = Probe{}
	}
	return json.Unmarshal(content, p)
}

func (p *Probe) ToDB() ([]byte, error) {
	return json.Marshal(p)
}
