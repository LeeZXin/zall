package server

import (
	"github.com/LeeZXin/zall/git/repo/server/githook"
	"github.com/LeeZXin/zall/git/repo/server/lfs"
	"github.com/LeeZXin/zall/git/repo/server/store"
)

func InitHttpApi() {
	lfs.InitApi()
	store.InitApi()
	githook.InitApi()
}
