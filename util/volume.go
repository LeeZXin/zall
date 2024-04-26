package util

import "fmt"

const (
	Gib = 1024 * 1024 * 1024
	Mib = 1024 * 1024
	Kib = 1024
)

func VolumeReadable(b int64) string {
	if b > Gib {
		return fmt.Sprintf("%dGB", b/Gib)
	}
	if b > Mib {
		return fmt.Sprintf("%dMB", b/Mib)
	}
	if b > Kib {
		return fmt.Sprintf("%dKB", b/Kib)
	}
	return fmt.Sprintf("%dB", b)
}
