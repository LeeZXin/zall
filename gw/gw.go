package gw

import (
	"github.com/LeeZXin/zsf/common"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
	"path/filepath"
	"strings"
)

func Init() {
	httpserver.AppendFilters(DefaultRouter())
}

func DefaultRouter() gin.HandlerFunc {
	return func(c *gin.Context) {
		p := c.Request.URL.Path
		if strings.HasPrefix(p, "/api") ||
			c.Request.Method != http.MethodGet ||
			strings.HasSuffix(p, "/info/refs") ||
			strings.HasSuffix(p, "/git-upload-pack") ||
			strings.HasSuffix(p, "/git-receive-pack") ||
			strings.HasPrefix(p, "/httpTask") {
			c.Next()
		} else {
			switch path.Ext(p) {
			case ".css", ".js", ".html", ".jpg", ".jpeg", ".png", ".ico":
				c.File(filepath.Join(common.ResourcesDir, "dist", c.Request.URL.RequestURI()))
			default:
				c.File(filepath.Join(common.ResourcesDir, "dist", "index.html"))
			}
			c.Abort()
		}
	}
}
