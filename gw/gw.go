package gw

import (
	"github.com/LeeZXin/zsf/common"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/gin-gonic/gin"
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
		if strings.HasPrefix(p, "/api") {
			c.Next()
		} else {
			switch path.Ext(p) {
			case ".css", ".js", ".html", ".jpg", ".jpeg", ".png", ".ico":
				c.File(filepath.Join(common.ResourcesDir, "dist", c.Request.URL.RequestURI()))
			default:
				c.File(filepath.Join(common.ResourcesDir, "dist", "index.html"))
			}
		}
	}
}
