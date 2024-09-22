package gw

import (
	"fmt"
	"github.com/LeeZXin/zsf/common"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
	"path/filepath"
	"strings"
	"time"
)

func Init() {
	httpserver.AppendFilters(DefaultRouter())
}

func DefaultRouter() gin.HandlerFunc {
	return func(c *gin.Context) {
		p := c.Request.URL.Path
		if strings.HasPrefix(p, "/api") ||
			strings.HasPrefix(p, "/zgit") ||
			strings.HasPrefix(p, "/actuator") ||
			strings.HasPrefix(p, "/debug") ||
			c.Request.Method != http.MethodGet ||
			strings.HasSuffix(p, "/info/refs") ||
			strings.HasSuffix(p, "/git-upload-pack") ||
			strings.HasSuffix(p, "/git-receive-pack") ||
			strings.HasPrefix(p, "/httpTask") {
			c.Next()
		} else {
			switch path.Ext(p) {
			case ".css", ".js", ".html", ".jpg", ".jpeg", ".png", ".ico":
				// 设置浏览器缓存策略
				now := time.Now().Unix()
				expires := now + 31536000
				c.Header("Date", fmt.Sprintf("%d", now))
				c.Header("Expires", fmt.Sprintf("%d", expires))
				c.Header("Cache-Control", "public, max-age=31536000")
				c.File(filepath.Join(common.ResourcesDir, "dist", c.Request.URL.RequestURI()))
			default:
				// 不缓存
				c.Header("Expires", "Fri, 01 Jan 1980 00:00:00 GMT")
				c.Header("Pragma", "no-cache")
				c.Header("Cache-Control", "no-cache, max-age=0, must-revalidate")
				c.File(filepath.Join(common.ResourcesDir, "dist", "index.html"))
			}
			c.Abort()
		}
	}
}
