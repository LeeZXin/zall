package util

import (
	"crypto/rand"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
)

func LfsRet(c *gin.Context, statusCode int, obj any) {
	ret, err := json.Marshal(obj)
	if err != nil {
		c.Error(err)
		c.Abort()
	} else {
		c.Data(statusCode, "application/vnd.git-lfs+json", ret)
	}
}

// NewRandomJwtSecret generates a new value intended to be used for JWT secrets.
func NewRandomJwtSecret() ([]byte, error) {
	bytes := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, bytes)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
