package zallet

type AppFile struct {
	App   string            `json:"app" yaml:"app"`
	Run   string            `json:"run" yaml:"run"`
	Stop  string            `json:"stop" yaml:"stop"`
	With  map[string]string `json:"with" yaml:"with"`
	Probe *Probe            `json:"probe" yaml:"probe"`
}
