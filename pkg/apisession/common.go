package apisession

import (
	"github.com/LeeZXin/zall/pkg/apicode"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zsf-utils/ginutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	SessionKey          = "session"
	LoginUser           = "loginUser"
	LoginCookie         = "zgit-auth"
	AuthorizationHeader = "Authorization"
)

func CheckLogin(c *gin.Context) {
	sessionId := GetSessionId(c)
	if sessionId == "" {
		c.JSON(http.StatusUnauthorized, ginutil.BaseResp{
			Code:    apicode.NotLoginCode.Int(),
			Message: i18n.GetByKey(i18n.SystemNotLogin),
		})
		c.Abort()
		return
	}
	sessionStore := GetStore()
	session, b, err := sessionStore.GetBySessionId(sessionId)
	if err != nil {
		logger.Logger.WithContext(c).Error(err)
		c.JSON(http.StatusInternalServerError, ginutil.BaseResp{
			Code:    apicode.InternalErrorCode.Int(),
			Message: i18n.GetByKey(i18n.SystemInternalError),
		})
		c.Abort()
		return
	}
	// session不存在
	if !b {
		c.JSON(http.StatusUnauthorized, ginutil.BaseResp{
			Code:    apicode.NotLoginCode.Int(),
			Message: i18n.GetByKey(i18n.SystemNotLogin),
		})
		c.Abort()
		return
	}
	c.Set(LoginUser, session.UserInfo)
	c.Set(SessionKey, session)
	c.Next()
}

func MustGetSession(c *gin.Context) Session {
	return c.MustGet(SessionKey).(Session)
}

func MustGetLoginUser(c *gin.Context) UserInfo {
	return c.MustGet(LoginUser).(UserInfo)
}

func GetSessionId(c *gin.Context) string {
	cookie, _ := c.Cookie(LoginCookie)
	if cookie == "" {
		return c.GetHeader(AuthorizationHeader)
	}
	return cookie
}
