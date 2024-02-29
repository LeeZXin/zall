package util

import "regexp"

var (
	IpPortPattern = regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}:\d$`)
)
