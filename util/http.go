package util

import (
	"encoding/base64"
	"errors"
	"github.com/LeeZXin/zall/pkg/apicode"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zsf-utils/bizerr"
	"github.com/LeeZXin/zsf-utils/ginutil"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type ValueWithExist[T any] struct {
	Exist bool `json:"exist"`
	Value T    `json:"value"`
}

func HandleApiErr(err error, c *gin.Context) {
	if err != nil {
		berr, ok := err.(*bizerr.Err)
		if !ok {
			c.JSON(http.StatusInternalServerError, ginutil.BaseResp{
				Code:    apicode.InternalErrorCode.Int(),
				Message: i18n.GetByKey(i18n.SystemInternalError),
			})
		} else {
			c.JSON(http.StatusOK, ginutil.BaseResp{
				Code:    berr.Code,
				Message: berr.Message,
			})
		}
	}
}

func ShouldBindJSON(obj any, c *gin.Context) bool {
	err := c.ShouldBindJSON(obj)
	if err != nil {
		ReturnHttpBadRequest(c)
		return false
	}
	return true
}

func ShouldBindQuery(obj any, c *gin.Context) bool {
	err := ginutil.BindQuery(c, obj)
	if err != nil {
		ReturnHttpBadRequest(c)
		return false
	}
	return true
}

func ShouldBindParams(obj any, c *gin.Context) bool {
	err := ginutil.BindParams(c, obj)
	if err != nil {
		ReturnHttpBadRequest(c)
		return false
	}
	return true
}

func ReturnHttpBadRequest(c *gin.Context) {
	c.JSON(http.StatusBadRequest, ginutil.BaseResp{
		Code:    apicode.BadRequestCode.Int(),
		Message: i18n.GetByKey(i18n.SystemInvalidArgs),
	})
}

func DefaultOkResponse(c *gin.Context) {
	c.JSON(http.StatusOK, ginutil.DefaultSuccessResp)
}

// BasicAuthDecode decode basic auth string
func BasicAuthDecode(encoded string) (string, string, error) {
	s, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", "", err
	}
	auth := strings.SplitN(string(s), ":", 2)

	if len(auth) != 2 {
		return "", "", errors.New("invalid basic authentication")
	}
	return auth[0], auth[1], nil
}
