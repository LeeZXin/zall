package workmode

type Mode int

const (
	StandaloneMode Mode = iota
	ShardingMode
	ProxyMode
	MetadataMode
)

var (
	currentMode Mode
)

func SetCurrentMode(mode Mode) {
	currentMode = mode
}

func IsStandaloneMode() bool {
	return currentMode == StandaloneMode
}

func IsShardingMode() bool {
	return currentMode == ShardingMode
}

func IsProxyMode() bool {
	return currentMode == ProxyMode
}

func IsMetadataMode() bool {
	return currentMode == MetadataMode
}
