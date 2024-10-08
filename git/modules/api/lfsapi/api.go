package lfsapi

import (
	"fmt"
	"github.com/LeeZXin/zall/git/modules/model/repomd"
	"github.com/LeeZXin/zall/git/modules/service/lfssrv"
	"github.com/LeeZXin/zall/git/modules/service/reposrv"
	"github.com/LeeZXin/zall/meta/modules/model/usermd"
	"github.com/LeeZXin/zall/meta/modules/service/cfgsrv"
	"github.com/LeeZXin/zall/meta/modules/service/usersrv"
	"github.com/LeeZXin/zall/pkg/git/lfs"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/LeeZXin/zsf/logger"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	// MediaType contains the media type for LFS server requests
	MediaType = "application/vnd.git-lfs+json"
)

func InitApi() {
	// 注册lfs api
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		infoLfs := e.Group("/:corpId/:repoName/info/lfs", packRepoPath)
		{
			infoLfs.POST("/objects/batch", checkMediaType, batch)
			infoLfs.PUT("/objects/:oid/:size", upload)
			infoLfs.GET("/objects/:oid/:filename", download)
			infoLfs.GET("/objects/:oid", download)
			infoLfs.POST("/verify", checkMediaType, verify)
			locks := infoLfs.Group("/locks", checkMediaType)
			{
				locks.GET("/", listLock)
				locks.POST("/", lock)
				locks.POST("/verify", listLockVerify)
				locks.POST("/:id/unlock", unlock)
			}
		}
	})
}

// packRepoPath
func packRepoPath(c *gin.Context) {
	ctx := c
	corpId := c.Param("corpId")
	repoName := c.Param("repoName")
	repoPath := filepath.Join(corpId, repoName)
	repo, b := reposrv.GetByRepoPath(ctx, repoPath)
	if !b {
		c.JSON(http.StatusUnauthorized, ErrVO{
			Message: i18n.GetByKey(i18n.SystemInvalidArgs),
		})
		c.Abort()
		return
	}
	authorization := c.GetHeader("Authorization")
	if authorization == "" {
		c.JSON(http.StatusUnauthorized, ErrVO{
			Message: i18n.GetByKey(i18n.SystemUnauthorized),
		})
		c.Abort()
		return
	}
	var userInfo usermd.UserInfo
	if strings.HasPrefix(authorization, "Basic ") {
		// 如果不是jwt 就看看是不是basic
		account, password, ok := c.Request.BasicAuth()
		if !ok {
			c.JSON(http.StatusUnauthorized, ErrVO{
				Message: i18n.GetByKey(i18n.SystemUnauthorized),
			})
			c.Abort()
			return
		}
		userInfo, b = usersrv.CheckAccountAndPassword(ctx, usersrv.CheckAccountAndPasswordReqDTO{
			Account:  account,
			Password: password,
		})
		if !b {
			c.JSON(http.StatusUnauthorized, ErrVO{
				Message: i18n.GetByKey(i18n.SystemUnauthorized),
			})
			c.Abort()
			return
		}
	} else {
		cfg, err := cfgsrv.GetGitCfgFromDB(c)
		if err != nil {
			logger.Logger.WithContext(c).Error(err)
			c.JSON(http.StatusInternalServerError, ErrVO{
				Message: i18n.GetByKey(i18n.SystemInternalError),
			})
			c.Abort()
			return
		}
		token, err := jwt.ParseWithClaims(authorization, new(lfs.Claims), func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return cfg.GetLfsJwtSecretBytes(), nil
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, ErrVO{
				Message: i18n.GetByKey(i18n.SystemUnauthorized),
			})
			c.Abort()
			return
		}
		claims, ok := token.Claims.(*lfs.Claims)
		if !ok {
			c.JSON(http.StatusUnauthorized, ErrVO{
				Message: i18n.GetByKey(i18n.SystemUnauthorized),
			})
			c.Abort()
			return
		}
		userInfo, b = usersrv.GetByAccount(ctx, claims.Account)
		if !b {
			c.JSON(http.StatusUnauthorized, ErrVO{
				Message: i18n.GetByKey(i18n.SystemInvalidArgs),
			})
			c.Abort()
			return
		}
	}
	c.Set("corpId", corpId)
	c.Set("operator", userInfo)
	c.Set("Authorization", authorization)
	c.Set("repo", repo)
	c.Next()
}

func checkMediaType(c *gin.Context) {
	header := c.GetHeader("Accept")
	accepts := strings.Split(header, ";")
	if len(accepts) == 0 || accepts[0] != MediaType {
		c.JSON(http.StatusUnsupportedMediaType, ErrVO{
			Message: "unsupported media type",
		})
		c.Abort()
		return
	} else {
		c.Next()
	}
}

func batch(c *gin.Context) {
	var req BatchReqVO
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusForbidden, ErrVO{
			Message:   i18n.GetByKey(i18n.SystemInvalidArgs),
			RequestID: logger.GetTraceId(c),
		})
		return
	}
	var isUpload bool
	if req.Operation == "upload" {
		r := getRepo(c)
		// 仓库归档或禁用lfs
		if r.IsArchived || r.GetCfg().DisableLfs {
			c.JSON(http.StatusForbidden, ErrVO{
				Message:   i18n.GetByKey(i18n.SystemForbidden),
				RequestID: logger.GetTraceId(c),
			})
			return
		}
		isUpload = true
	} else if req.Operation == "download" {
		isUpload = false
	} else {
		c.JSON(http.StatusBadRequest, ErrVO{
			Message:   "bad request",
			RequestID: logger.GetTraceId(c),
		})
		return
	}
	cfg, err := cfgsrv.GetGitCfgFromDB(c)
	if err != nil {
		logger.Logger.WithContext(c).Error(err)
		c.JSON(http.StatusInternalServerError, ErrVO{
			Message:   i18n.GetByKey(i18n.SystemInternalError),
			RequestID: logger.GetTraceId(c),
		})
		return
	}
	reqDTO := lfssrv.BatchReqDTO{
		Repo:     getRepo(c),
		Operator: getOperator(c),
		IsUpload: isUpload,
		RefName:  req.Ref.Name,
	}
	reqDTO.Objects = listutil.MapNe(req.Objects, func(t PointerVO) lfssrv.PointerDTO {
		return lfssrv.PointerDTO{
			Oid:  t.Oid,
			Size: t.Size,
		}
	})
	respDTO, err := lfssrv.Batch(c, reqDTO)
	if err != nil {
		c.JSON(http.StatusForbidden, ErrVO{
			Message:   err.Error(),
			RequestID: logger.GetTraceId(c),
		})
		return
	}
	authorization := c.MustGet("Authorization").(string)
	header := map[string]string{
		"Authorization": authorization,
	}
	verifyHeader := map[string]string{
		"Accept":        MediaType,
		"Authorization": authorization,
	}
	var resp BatchRespVO
	repoPath := getRepo(c).Path
	resp.Objects = listutil.MapNe(respDTO.ObjectList, func(t lfssrv.ObjectDTO) ObjectRespVO {
		if t.ErrObjDTO.Code == 0 {
			var actions map[string]LinkVO
			if isUpload {
				actions = map[string]LinkVO{
					"upload": {
						Href:   fmt.Sprintf("%s/%s/info/lfs/objects/%s/%d", cfg.HttpUrl, repoPath, t.Oid, t.Size),
						Header: header,
					},
					"verify": {
						Href:   fmt.Sprintf("%s/%s/info/lfs/verify", cfg.HttpUrl, repoPath),
						Header: verifyHeader,
					},
				}
			} else {
				actions = map[string]LinkVO{
					"download": {
						Href:   fmt.Sprintf("%s/%s/info/lfs/objects/%s", cfg.HttpUrl, repoPath, t.Oid),
						Header: header,
					},
				}
			}
			return ObjectRespVO{
				PointerVO: PointerVO{
					Oid:  t.Oid,
					Size: t.Size,
				},
				Actions: actions,
			}
		} else {
			return ObjectRespVO{
				Error: &ObjectErrVO{
					Code:    t.ErrObjDTO.Code,
					Message: t.ErrObjDTO.Message,
				},
			}
		}
	})
	c.Header("Content-Type", MediaType)
	c.JSON(http.StatusOK, resp)
}

func getOperator(c *gin.Context) usermd.UserInfo {
	return c.MustGet("operator").(usermd.UserInfo)
}

func getRepo(c *gin.Context) repomd.Repo {
	return c.MustGet("repo").(repomd.Repo)
}

func lock(c *gin.Context) {
	var req PostLockReqVO
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusForbidden, ErrVO{
			Message:   i18n.GetByKey(i18n.SystemInvalidArgs),
			RequestID: logger.GetTraceId(c),
		})
		return
	}
	operator := getOperator(c)
	respDTO, err := lfssrv.Lock(c, lfssrv.LockReqDTO{
		Repo:     getRepo(c),
		RefName:  req.Ref.Name,
		Operator: operator,
		Path:     req.Path,
	})
	if err != nil {
		c.JSON(http.StatusForbidden, ErrVO{
			Message:   err.Error(),
			RequestID: logger.GetTraceId(c),
		})
		return
	}
	// 锁冲突
	if respDTO.AlreadyExists {
		c.JSON(http.StatusConflict, PostLockRespVO{
			Lock: model2LockVO(respDTO.Lock, operator),
			ErrVO: ErrVO{
				Message:   "already created lock",
				RequestID: logger.GetTraceId(c),
			},
		})
	} else {
		c.JSON(http.StatusOK, PostLockRespVO{
			Lock: model2LockVO(respDTO.Lock, operator),
		})
	}
}

func listLock(c *gin.Context) {
	var req ListLockReqVO
	err := c.ShouldBindQuery(&req)
	if err != nil {
		c.JSON(http.StatusForbidden, ErrVO{
			Message:   i18n.GetByKey(i18n.SystemInvalidArgs),
			RequestID: logger.GetTraceId(c),
		})
		return
	}
	operator := getOperator(c)
	cursor, _ := strconv.ParseInt(req.Cursor, 10, 64)
	listResp, err := lfssrv.ListLock(c, lfssrv.ListLockReqDTO{
		Repo:     getRepo(c),
		Operator: operator,
		Cursor:   cursor,
		Limit:    req.Limit,
		RefName:  req.RefSpec,
	})
	if err != nil {
		c.JSON(http.StatusForbidden, ErrVO{
			Message:   err.Error(),
			RequestID: logger.GetTraceId(c),
		})
		return
	}
	listVO := listutil.MapNe(listResp.LockList, func(lock lfssrv.LfsLockDTO) LockVO {
		return model2LockVO(lock, operator)
	})
	resp := ListLockRespVO{
		Locks: listVO,
	}
	if listResp.Next > 0 {
		resp.Next = strconv.FormatInt(listResp.Next, 10)
	}
	c.JSON(http.StatusOK, resp)
}

func unlock(c *gin.Context) {
	var req UnlockReqVO
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusForbidden, ErrVO{
			Message:   i18n.GetByKey(i18n.SystemInvalidArgs),
			RequestID: logger.GetTraceId(c),
		})
		return
	}
	lockId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusForbidden, ErrVO{
			Message:   i18n.GetByKey(i18n.SystemInvalidArgs),
			RequestID: logger.GetTraceId(c),
		})
		return
	}
	operator := getOperator(c)
	singleLock, err := lfssrv.Unlock(c, lfssrv.UnlockReqDTO{
		Repo:     getRepo(c),
		LockId:   lockId,
		Force:    req.Force,
		Operator: operator,
	})
	if err != nil {
		c.JSON(http.StatusForbidden, ErrVO{
			Message:   err.Error(),
			RequestID: logger.GetTraceId(c),
		})
		return
	}
	c.JSON(http.StatusOK, UnlockRespVO{
		Lock: model2LockVO(singleLock, operator),
	})
}

func listLockVerify(c *gin.Context) {
	var req ListLockVerifyReqVO
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusForbidden, ErrVO{
			Message:   i18n.GetByKey(i18n.SystemInvalidArgs),
			RequestID: logger.GetTraceId(c),
		})
		return
	}
	cursor, _ := strconv.ParseInt(req.Cursor, 10, 64)
	operator := getOperator(c)
	listResp, err := lfssrv.ListLock(c, lfssrv.ListLockReqDTO{
		Repo:     getRepo(c),
		Operator: operator,
		Cursor:   cursor,
		Limit:    req.Limit,
		RefName:  req.Ref.Name,
	})
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, ErrVO{
			Message: err.Error(),
		})
		return
	}
	voList := listResp.LockList
	ours := listutil.FilterNe(voList, func(l lfssrv.LfsLockDTO) bool {
		return l.Owner == operator.Account
	})
	oursRet := listutil.MapNe(ours, func(l lfssrv.LfsLockDTO) LockVO {
		return model2LockVO(l, operator)
	})
	theirs := listutil.FilterNe(voList, func(l lfssrv.LfsLockDTO) bool {
		return l.Owner != operator.Account
	})
	theirsRet := listutil.MapNe(theirs, func(l lfssrv.LfsLockDTO) LockVO {
		return model2LockVO(l, operator)
	})
	respVO := ListLockVerifyRespVO{
		Ours:   oursRet,
		Theirs: theirsRet,
	}
	if listResp.Next > 0 {
		respVO.Next = strconv.FormatInt(listResp.Next, 10)
	}
	c.JSON(http.StatusOK, respVO)
}

func verify(c *gin.Context) {
	var req PointerVO
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusForbidden, ErrVO{
			Message:   i18n.GetByKey(i18n.SystemInvalidArgs),
			RequestID: logger.GetTraceId(c),
		})
		return
	}
	exists, validate, err := lfssrv.Verify(c, lfssrv.VerifyReqDTO{
		PointerDTO: lfssrv.PointerDTO{
			Oid:  req.Oid,
			Size: req.Size,
		},
		Repo:     getRepo(c),
		Operator: getOperator(c),
	})
	if err != nil {
		c.JSON(http.StatusForbidden, ErrVO{
			Message:   err.Error(),
			RequestID: logger.GetTraceId(c),
		})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, ErrVO{
			Message:   "not found",
			RequestID: logger.GetTraceId(c),
		})
		return
	}
	if !validate {
		c.JSON(http.StatusUnprocessableEntity, ErrVO{
			Message:   "validate error",
			RequestID: logger.GetTraceId(c),
		})
		return
	}
	writeRespMessage(c, http.StatusOK, "")
}

func download(c *gin.Context) {
	oid := c.Param("oid")
	ctx := c
	err := lfssrv.Download(ctx, lfssrv.DownloadReqDTO{
		Oid:      oid,
		Repo:     getRepo(c),
		Operator: getOperator(c),
		C:        c,
	})
	if err != nil {
		c.JSON(http.StatusForbidden, ErrVO{
			Message:   err.Error(),
			RequestID: logger.GetTraceId(c),
		})
		return
	}
}

func upload(c *gin.Context) {
	body := c.Request.Body
	defer body.Close()
	size, err := strconv.ParseInt(c.Param("size"), 10, 64)
	if err != nil {
		c.JSON(http.StatusForbidden, ErrVO{
			Message:   i18n.GetByKey(i18n.SystemInvalidArgs),
			RequestID: logger.GetTraceId(c),
		})
		return
	}
	repo := getRepo(c)
	// 仓库归档或禁用lfs
	if repo.IsArchived || repo.GetCfg().DisableLfs {
		c.JSON(http.StatusForbidden, ErrVO{
			Message:   i18n.GetByKey(i18n.SystemForbidden),
			RequestID: logger.GetTraceId(c),
		})
		return
	}
	oid := c.Param("oid")
	err = lfssrv.Upload(c, lfssrv.UploadReqDTO{
		Oid:      oid,
		Size:     size,
		Repo:     repo,
		Operator: getOperator(c),
		C:        c,
	})
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, ErrVO{
			Message:   err.Error(),
			RequestID: logger.GetTraceId(c),
		})
		return
	}
}

func model2LockVO(l lfssrv.LfsLockDTO, locker usermd.UserInfo) LockVO {
	return LockVO{
		Id:       strconv.FormatInt(l.LockId, 10),
		Path:     l.Path,
		LockedAt: l.Created.Round(time.Second),
		Owner: &LockOwnerVO{
			Name: locker.Name,
		},
	}
}

func writeRespMessage(c *gin.Context, code int, message string) {
	c.Data(code, MediaType, []byte(message))
}
