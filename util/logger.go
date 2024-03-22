package util

import (
	"strings"
	"time"
)

type SimpleLogger struct {
	sb *strings.Builder
}

func NewSimpleLogger() *SimpleLogger {
	return &SimpleLogger{
		sb: new(strings.Builder),
	}
}

func (l *SimpleLogger) WriteString(content string) {
	l.sb.WriteString(time.Now().Format(time.DateTime) + " " + content + "\n")
}

func (l *SimpleLogger) Clear() {
	l.sb.Reset()
}

func (l *SimpleLogger) ToString() string {
	return l.sb.String()
}
