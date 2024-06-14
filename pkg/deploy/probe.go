package deploy

import (
	"net/url"
	"strings"
)

type ProbeType string

const (
	HttpProbeType ProbeType = "http"
	TcpProbeType  ProbeType = "tcp"
)

type TcpProbe struct {
	Addr string `json:"addr" yaml:"addr"`
}

func (t *TcpProbe) isValid() bool {
	return IpPortPattern.MatchString(t.Addr)
}

type HttpProbe struct {
	Url string `json:"url" yaml:"url"`
}

func (t *HttpProbe) isValid() bool {
	parsed, err := url.Parse(t.Url)
	if err != nil {
		return false
	}
	return strings.HasPrefix(parsed.Scheme, "http")
}

type Probe struct {
	Type ProbeType  `json:"type" yaml:"type"`
	Tcp  *TcpProbe  `json:"tcp,omitempty" yaml:"tcp,omitempty"`
	Http *HttpProbe `json:"http,omitempty" yaml:"http,omitempty"`
}

func (p *Probe) IsValid() bool {
	switch p.Type {
	case HttpProbeType:
		return p.Http != nil && p.Http.isValid()
	case TcpProbeType:
		return p.Tcp != nil && p.Tcp.isValid()
	default:
		return false
	}
}
