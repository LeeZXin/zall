package util

import "regexp"

func GenIpPortPattern() *regexp.Regexp {
	return regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}:\d+$`)
}
