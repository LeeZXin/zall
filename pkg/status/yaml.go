package status

import (
	"encoding/json"
)

type ProbeType string

type TcpProbe struct {
	Host string `json:"host" yaml:"host"`
}

type HttpProbe struct {
	Url string `json:"url" yaml:"url"`
}

type Probe struct {
	Type ProbeType  `json:"type" yaml:"type"`
	Tcp  *TcpProbe  `json:"tcp,omitempty" yaml:"tcp,omitempty"`
	Http *HttpProbe `json:"http,omitempty" yaml:"http,omitempty"`
}

type Yaml struct {
	Env   string            `json:"env" yaml:"env"`
	App   string            `json:"app" yaml:"app"`
	Start string            `json:"start" yaml:"start"`
	With  map[string]string `json:"with" yaml:"with"`
	Probe *Probe            `json:"probe" yaml:"probe"`
}

func (f *Yaml) FromDB(content []byte) error {
	return json.Unmarshal(content, f)
}

func (f *Yaml) ToDB() ([]byte, error) {
	return json.Marshal(f)
}
