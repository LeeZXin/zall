package util

import (
	"os"
	"strings"
)

func WriteFile(filePath string, content []byte) error {
	return os.WriteFile(filePath, content, os.ModePerm)
}

func AppendFile(filePath string, content []byte) error {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	_, err = f.Write(content)
	if err1 := f.Close(); err1 != nil && err == nil {
		err = err1
	}
	return err
}

func ContainsParentDirectorySeparator(v string) bool {
	if !strings.Contains(v, "..") {
		return false
	}
	for _, ent := range strings.FieldsFunc(v, isSlashRune) {
		if ent == ".." {
			return true
		}
	}
	return false
}

func isSlashRune(r rune) bool { return r == '/' || r == '\\' }
