package util

import (
	"encoding/base64"
	"errors"
	"github.com/LeeZXin/zall/pkg/apicode"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zsf-utils/bizerr"
	"github.com/LeeZXin/zsf-utils/ginutil"
	"github.com/LeeZXin/zsf/http/httpclient"
	"github.com/LeeZXin/zsf/property/static"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

var (
	HttpTokenOpts = httpclient.WithHeader(map[string]string{
		"z-token": static.GetString("http.token"),
	})
)

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
		c.JSON(http.StatusBadRequest, ginutil.BaseResp{
			Code:    apicode.BadRequestCode.Int(),
			Message: i18n.GetByKey(i18n.SystemInvalidArgs),
		})
		return false
	}
	return true
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
