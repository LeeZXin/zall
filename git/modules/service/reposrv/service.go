package reposrv

import (
	"context"
	"fmt"
	"github.com/LeeZXin/zall/git/modules/model/branchmd"
	"github.com/LeeZXin/zall/git/modules/model/gpgkeymd"
	"github.com/LeeZXin/zall/git/modules/model/pullrequestmd"
	"github.com/LeeZXin/zall/git/modules/model/repomd"
	"github.com/LeeZXin/zall/git/modules/model/webhookmd"
	"github.com/LeeZXin/zall/git/modules/model/workflowmd"
	"github.com/LeeZXin/zall/git/modules/service/gpgkeysrv"
	"github.com/LeeZXin/zall/git/modules/service/sshkeysrv"
	"github.com/LeeZXin/zall/git/repo/client"
	"github.com/LeeZXin/zall/git/repo/reqvo"
	"github.com/LeeZXin/zall/meta/modules/model/teammd"
	"github.com/LeeZXin/zall/meta/modules/service/cfgsrv"
	"github.com/LeeZXin/zall/meta/modules/service/usersrv"
	"github.com/LeeZXin/zall/pkg/apicode"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/event"
	"github.com/LeeZXin/zall/pkg/git"
	"github.com/LeeZXin/zall/pkg/git/signature"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zall/pkg/limiter"
	"github.com/LeeZXin/zall/pkg/perm"
	"github.com/LeeZXin/zall/pkg/teamhook"
	"github.com/LeeZXin/zall/pkg/webhook"
	"github.com/LeeZXin/zall/teamhook/modules/service/teamhooksrv"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/bizerr"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf-utils/psub"
	"github.com/LeeZXin/zsf/logger"
	"github.com/LeeZXin/zsf/property/static"
	"github.com/LeeZXin/zsf/xorm/xormstore"
	"github.com/keybase/go-crypto/openpgp/packet"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	accessRepo = iota
	updateRepo
)

// GetByRepoPath 通过相对路径获取仓库信息
func GetByRepoPath(ctx context.Context, path string) (repomd.Repo, bool) {
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	r, b, err := repomd.GetByPath(ctx, path)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
	}
	return r, b
}

var (
	createArchiveLimiter     limiter.Limiter
	createArchiveLimiterOnce = sync.Once{}
	initPsubOnce             = sync.Once{}
)

func initCreateArchiveLimiter() {
	createArchiveLimiterOnce.Do(func() {
		limit := static.GetInt64("createArchiveLimit")
		if limit <= 0 {
			limit = 10
		}
		createArchiveLimiter = limiter.NewCountLimiter(limit)
	})
}

func initPsub() {
	initPsubOnce.Do(func() {
		psub.Subscribe(event.GitRepoTopic, func(data any) {
			req, ok := data.(event.GitRepoEvent)
			if ok {
				ctx, closer := xormstore.Context(context.Background())
				// 触发webhook
				hookList, err := webhookmd.ListWebhookByRepoId(ctx, req.RepoId)
				closer.Close()
				if err == nil && len(hookList) > 0 {
					for _, hook := range hookList {
						if hook.GetEvents().GitRepo {
							webhook.TriggerWebhook(hook.HookUrl, hook.Secret, &req)
						}
					}
				}
				teamhooksrv.TriggerTeamHook(&req, req.TeamId, func(events *teamhook.Events) bool {
					switch req.Action {
					case event.RepoEventCreateAction:
						return events.GitRepo.Create
					case event.RepoEventUpdateAction:
						return events.GitRepo.Update
					case event.RepoEventDeleteTemporarilyAction:
						return events.GitRepo.DeleteTemporarily
					case event.RepoEventDeletePermanentlyAction:
						return events.GitRepo.DeletePermanently
					case event.RepoEventArchivedAction:
						return events.GitRepo.Archived
					case event.RepoEventUnArchivedAction:
						return events.GitRepo.UnArchived
					case event.RepoEventRecoverFromRecycleAction:
						return events.GitRepo.RecoverFromRecycle
					default:
						return false
					}
				})
			}
		})
		psub.Subscribe(event.GitPushTopic+"-remote", func(data any) {
			req, ok := data.(event.GitPushEvent)
			if ok {
				teamhooksrv.TriggerTeamHook(&req, req.TeamId, func(events *teamhook.Events) bool {
					switch req.Action {
					case event.GitPushEventDeleteAction:
						return events.GitPush.Delete
					default:
						return false
					}
				})
			}
		})
	})
}

// GetRepoAndPerm 获取仓库信息和权限信息
func GetRepoAndPerm(ctx context.Context, reqDTO GetRepoAndPermReqDTO) (SimpleRepoDTO, perm.RepoPerm, error) {
	if err := reqDTO.IsValid(); err != nil {
		return SimpleRepoDTO{}, perm.RepoPerm{}, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	repo, b, err := repomd.GetByRepoId(ctx, reqDTO.RepoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return SimpleRepoDTO{}, perm.RepoPerm{}, util.InternalError(err)
	}
	if !b {
		return SimpleRepoDTO{}, perm.RepoPerm{}, util.InvalidArgsError()
	}
	if reqDTO.Operator.IsAdmin {
		return SimpleRepoDTO{
			RepoId: repo.Id,
			Name:   repo.Name,
			TeamId: repo.TeamId,
		}, perm.DefaultRepoPerm, nil
	}
	// 校验权限
	p, b, err := teammd.GetUserPermDetail(ctx, repo.TeamId, reqDTO.Operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return SimpleRepoDTO{}, perm.RepoPerm{}, util.InternalError(err)
	}
	if !b || (!p.IsAdmin && !p.PermDetail.GetRepoPerm(reqDTO.RepoId).CanAccessRepo) {
		return SimpleRepoDTO{}, perm.RepoPerm{}, util.UnauthorizedError()
	}
	return SimpleRepoDTO{
		RepoId: repo.Id,
		Name:   repo.Name,
		TeamId: repo.TeamId,
	}, p.PermDetail.GetRepoPerm(reqDTO.RepoId), nil
}

func repo2Dto(t repomd.Repo) RepoDTO {
	ret := RepoDTO{
		Id:            t.Id,
		Path:          t.Path,
		Name:          t.Name,
		TeamId:        t.TeamId,
		RepoDesc:      t.RepoDesc,
		DefaultBranch: t.DefaultBranch,
		GitSize:       t.GitSize,
		LfsSize:       t.LfsSize,
		LastOperated:  t.LastOperated,
		IsArchived:    t.IsArchived,
		Created:       t.Created,
	}
	if t.Cfg != nil {
		ret.DisableLfs = t.Cfg.DisableLfs
		ret.LfsLimitSize = t.Cfg.LfsLimitSize
		ret.GitLimitSize = t.Cfg.GitLimitSize
	}
	return ret
}

func EntriesRepo(ctx context.Context, reqDTO EntriesRepoReqDTO) ([]BlobDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	repo, b, err := repomd.GetByRepoId(ctx, reqDTO.RepoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	if !b {
		return nil, util.InvalidArgsError()
	}
	err = checkPermByRepo(ctx, repo, reqDTO.Operator, accessRepo)
	if err != nil {
		return nil, err
	}
	resp, err := client.EntriesRepo(ctx, reqvo.EntriesRepoReq{
		RepoPath: repo.Path,
		Ref:      reqDTO.Ref,
		Dir:      reqDTO.Dir,
		RefType:  reqDTO.RefType,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.Map(resp, func(t reqvo.BlobVO) (BlobDTO, error) {
		return BlobDTO{
			Mode:    t.Mode,
			RawPath: t.RawPath,
			Path:    t.Path,
		}, nil
	})
}

// ListRepo 展示仓库列表
func ListRepo(ctx context.Context, reqDTO ListRepoReqDTO) ([]RepoDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	var (
		repoList []repomd.Repo
		err      error
	)
	if reqDTO.Operator.IsAdmin {
		repoList, err = repomd.ListRepoByTeamId(ctx, reqDTO.TeamId)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return nil, util.InternalError(err)
		}
	} else {
		p, b, err := teammd.GetUserPermDetail(ctx, reqDTO.TeamId, reqDTO.Operator.Account)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return nil, util.InternalError(err)
		}
		if !b {
			return nil, util.UnauthorizedError()
		}
		if len(p.PermDetail.RepoPermList) > 0 {
			// 访问部分仓库
			repoPermList := listutil.FilterNe(p.PermDetail.RepoPermList, func(p perm.RepoPermWithId) bool {
				return p.CanAccessRepo
			})
			// 存在部分可读仓库权限
			if len(repoPermList) > 0 {
				idList := listutil.MapNe(repoPermList, func(t perm.RepoPermWithId) int64 {
					return t.RepoId
				})
				repoList, err = repomd.GetRepoByIdList(ctx, idList, nil)
				if err != nil {
					logger.Logger.WithContext(ctx).Error(err)
					return nil, util.InternalError(err)
				}
			} else {
				repoList = []repomd.Repo{}
			}
		} else if p.PermDetail.DefaultRepoPerm.CanAccessRepo {
			repoList, err = repomd.ListRepoByTeamId(ctx, reqDTO.TeamId)
			if err != nil {
				logger.Logger.WithContext(ctx).Error(err)
				return nil, util.InternalError(err)
			}
		} else {
			repoList = []repomd.Repo{}
		}
	}
	sort.SliceStable(repoList, func(i, j int) bool {
		return repoList[i].LastOperated.After(repoList[j].LastOperated)
	})
	return listutil.Map(repoList, func(t repomd.Repo) (RepoDTO, error) {
		return repo2Dto(t), nil
	})
}

// ListDeletedRepo 展示已删除仓库
func ListDeletedRepo(ctx context.Context, reqDTO ListDeletedRepoReqDTO) ([]DeletedRepoDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, err := checkTeamAdminByTeamId(ctx, reqDTO.TeamId, reqDTO.Operator)
	if err != nil {
		return nil, err
	}
	repoList, err := repomd.GetDeletedRepoListByTeamId(ctx, reqDTO.TeamId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, err
	}
	return listutil.Map(repoList, func(t repomd.Repo) (DeletedRepoDTO, error) {
		return DeletedRepoDTO{
			RepoDTO: repo2Dto(t),
			Deleted: t.Deleted,
		}, nil
	})
}

// CatFile 展示文件内容
func CatFile(ctx context.Context, reqDTO CatFileReqDTO) (CatFileRespDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return CatFileRespDTO{}, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	repo, b, err := repomd.GetByRepoId(ctx, reqDTO.RepoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return CatFileRespDTO{}, util.InternalError(err)
	}
	if !b {
		return CatFileRespDTO{}, util.InvalidArgsError()
	}
	err = checkPermByRepo(ctx, repo, reqDTO.Operator, accessRepo)
	if err != nil {
		return CatFileRespDTO{}, err
	}
	resp, err := client.CatFile(ctx, reqvo.CatFileReq{
		RepoPath: repo.Path,
		Ref:      reqDTO.Ref,
		FilePath: reqDTO.FilePath,
		RefType:  reqDTO.RefType,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return CatFileRespDTO{}, util.InternalError(err)
	}
	commit := commit2Dto(resp.Commit)
	if commit.Committer.Email != "" {
		userMap, err := usersrv.GetUsersNameAndAvatarMapByEmails(ctx, commit.Committer.Email)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return CatFileRespDTO{}, util.InternalError(err)
		}
		commit.Committer.AvatarUrl = userMap[commit.Committer.Email].AvatarUrl
		commit.Committer.Name = userMap[commit.Committer.Email].Name
		commit.Committer.Account = userMap[commit.Committer.Email].Account
	}
	return CatFileRespDTO{
		FileMode: resp.FileMode,
		ModeName: resp.ModeName,
		Content:  resp.Content,
		Size:     resp.Size,
		Commit:   commit,
	}, nil
}

// IndexRepo 代码首页
func IndexRepo(ctx context.Context, reqDTO IndexRepoReqDTO) (IndexRepoRespDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return IndexRepoRespDTO{}, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	repo, b, err := repomd.GetByRepoId(ctx, reqDTO.RepoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return IndexRepoRespDTO{}, util.InternalError(err)
	}
	if !b {
		return IndexRepoRespDTO{}, util.InvalidArgsError()
	}
	err = checkPermByRepo(ctx, repo, reqDTO.Operator, accessRepo)
	if err != nil {
		return IndexRepoRespDTO{}, err
	}
	resp, err := client.IndexRepo(ctx, reqvo.IndexRepoReq{
		RepoPath: repo.Path,
		Ref:      reqDTO.Ref,
		RefType:  reqDTO.RefType,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return IndexRepoRespDTO{}, util.InternalError(err)
	}
	lastestCommit := commit2Dto(resp.LatestCommit)
	if resp.LatestCommit.Committer.Email != "" {
		userMap, err := usersrv.GetUsersNameAndAvatarMapByEmails(ctx, resp.LatestCommit.Committer.Email)
		if err != nil {
			return IndexRepoRespDTO{}, util.InternalError(err)
		}
		lastestCommit.Committer.AvatarUrl = userMap[resp.LatestCommit.Committer.Email].AvatarUrl
		lastestCommit.Committer.Name = userMap[resp.LatestCommit.Committer.Email].Name
		lastestCommit.Committer.Account = userMap[resp.LatestCommit.Committer.Email].Account
	}
	return IndexRepoRespDTO{
		ReadmeText:   resp.ReadmeText,
		HasReadme:    resp.HasReadme,
		LatestCommit: lastestCommit,
		Tree:         tree2Dto(resp.Tree),
	}, nil
}

// GetBaseInfo 基本信息
func GetBaseInfo(ctx context.Context, reqDTO GetBaseInfoReqDTO) (BaseInfoDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return BaseInfoDTO{}, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	repo, b, err := repomd.GetByRepoId(ctx, reqDTO.RepoId)
	if err != nil {
		return BaseInfoDTO{}, util.InternalError(err)
	}
	if !b {
		return BaseInfoDTO{}, util.InvalidArgsError()
	}
	err = checkPermByRepo(ctx, repo, reqDTO.Operator, accessRepo)
	if err != nil {
		return BaseInfoDTO{}, err
	}
	branches, err := client.GetAllBranches(ctx, reqvo.GetAllBranchesReq{
		RepoPath: repo.Path,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return BaseInfoDTO{}, util.InternalError(err)
	}
	tags, err := client.GetAllTags(ctx, reqvo.GetAllTagsReq{
		RepoPath: repo.Path,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return BaseInfoDTO{}, util.InternalError(err)
	}
	ret := BaseInfoDTO{}
	ret.Branches = listutil.MapNe(branches, func(t reqvo.RefVO) string {
		return t.Name
	})
	ret.Tags = listutil.MapNe(tags, func(t reqvo.RefVO) string {
		return t.Name
	})
	cfg, err := cfgsrv.GetGitCfgFromDB(ctx)
	if err == nil {
		if cfg.HttpUrl != "" {
			ret.CloneHttpUrl = strings.TrimSuffix(cfg.HttpUrl, "/") + "/" + repo.Path
		}
		if cfg.SshUrl != "" {
			ret.CloneSshUrl = strings.TrimSuffix(cfg.SshUrl, "/") + "/" + repo.Path
		}
	}
	ret.DefaultBranch = repo.DefaultBranch
	return ret, nil
}

// GetDetailInfo 基本信息
func GetDetailInfo(ctx context.Context, reqDTO GetDetailInfoReqDTO) (RepoDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return RepoDTO{}, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	repo, b, err := repomd.GetByRepoId(ctx, reqDTO.RepoId)
	if err != nil {
		return RepoDTO{}, util.InternalError(err)
	}
	if !b {
		return RepoDTO{}, util.InvalidArgsError()
	}
	err = checkPermByRepo(ctx, repo, reqDTO.Operator, accessRepo)
	if err != nil {
		return RepoDTO{}, err
	}
	return repo2Dto(repo), nil
}

func tree2Dto(vo reqvo.TreeVO) TreeDTO {
	ret := TreeDTO{}
	ret.Files = listutil.MapNe(vo.Files, func(t reqvo.FileVO) FileDTO {
		return FileDTO{
			Mode:    t.Mode,
			RawPath: t.RawPath,
			Path:    t.Path,
			Commit:  commit2Dto(t.Commit),
		}
	})
	return ret
}

// CreateRepo 初始化仓库
func CreateRepo(ctx context.Context, reqDTO CreateRepoReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	// 校验权限
	team, err := checkTeamAdminByTeamId(ctx, reqDTO.TeamId, reqDTO.Operator)
	if err != nil {
		return err
	}
	var b bool
	// 相对路径
	relativePath := filepath.Join("zgit", reqDTO.Name+".git")
	_, b, err = repomd.GetByPathWithoutJudgingDeleted(ctx, relativePath)
	// 数据库异常
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	// 仓库已存在 不能添加
	if b {
		return util.NewBizErr(apicode.InvalidArgsCode, i18n.RepoAlreadyExists)
	}
	var repo repomd.Repo
	// 添加数据
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		var err2 error
		// 插入数据库
		repo, err2 = repomd.InsertRepo(ctx, repomd.InsertRepoReqDTO{
			Name:          reqDTO.Name,
			Path:          relativePath,
			TeamId:        reqDTO.TeamId,
			RepoDesc:      reqDTO.Desc,
			DefaultBranch: reqDTO.DefaultBranch,
			LastOperated:  time.Now(),
		})
		if err2 != nil {
			return err2
		}
		// 调用store
		gitSize, err2 := client.InitRepo(ctx, reqvo.InitRepoReq{
			UserAccount:   reqDTO.Operator.Account,
			UserEmail:     reqDTO.Operator.Email,
			RepoName:      reqDTO.Name,
			RepoPath:      relativePath,
			AddReadme:     reqDTO.AddReadme,
			GitIgnoreName: reqDTO.GitIgnoreName,
			DefaultBranch: reqDTO.DefaultBranch,
		})
		if err2 == nil {
			repomd.UpdateGitSize(ctx, repo.Id, gitSize)
		}
		return err2
	})
	if err != nil {
		if bizerr.IsBizErr(err) {
			return err
		}
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	notifyEvent(repo, team, reqDTO.Operator, event.RepoEventCreateAction)
	return nil
}

// AllGitIgnoreTemplateList 所有gitignore模版名称
func AllGitIgnoreTemplateList() []string {
	return gitignoreSet.AllKeys()
}

// DeleteRepo 删除仓库
func DeleteRepo(ctx context.Context, reqDTO DeleteRepoReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	repo, team, err := checkTeamAdminByRepoId(ctx, reqDTO.RepoId, reqDTO.Operator)
	if err != nil {
		return err
	}
	_, err = repomd.SetRepoDeleted(ctx, reqDTO.RepoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	notifyEvent(repo, team, reqDTO.Operator, event.RepoEventDeleteTemporarilyAction)
	return nil
}

// RecoverFromRecycle 恢复仓库
func RecoverFromRecycle(ctx context.Context, reqDTO RecoverFromRecycleReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	repo, b, err := repomd.GetByRepoIdWithoutJudgingDeleted(ctx, reqDTO.RepoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b || !repo.IsDeleted {
		return util.InvalidArgsError()
	}
	team, err := checkTeamAdminByTeamId(ctx, repo.TeamId, reqDTO.Operator)
	if err != nil {
		return err
	}
	_, err = repomd.SetRepoUnDeleted(ctx, reqDTO.RepoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	notifyEvent(repo, team, reqDTO.Operator, event.RepoEventRecoverFromRecycleAction)
	return nil
}

// DeleteRepoPermanently 永久删除仓库
func DeleteRepoPermanently(ctx context.Context, reqDTO DeleteRepoReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	repo, b, err := repomd.GetByRepoIdWithoutJudgingDeleted(ctx, reqDTO.RepoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b || !repo.IsDeleted {
		return util.InvalidArgsError()
	}
	team, err := checkTeamAdminByTeamId(ctx, repo.TeamId, reqDTO.Operator)
	if err != nil {
		return err
	}
	err = xormstore.WithTx(ctx, func(ctx context.Context) error {
		// 删除db
		_, err2 := repomd.DeleteRepo(ctx, reqDTO.RepoId)
		if err2 != nil {
			return err2
		}
		// 删除远程仓库
		err2 = client.DeleteRepo(ctx, reqvo.DeleteRepoReq{
			RepoPath: repo.Path,
		})
		if err2 != nil {
			return err2
		}
		// 删除合并请求
		err2 = pullrequestmd.DeletePullRequestByRepoId(ctx, reqDTO.RepoId)
		if err2 != nil {
			return err2
		}
		// 删除工作流任务
		err2 = workflowmd.DeleteTaskByRepoId(ctx, reqDTO.RepoId)
		if err2 != nil {
			return err2
		}
		// 删除工作流
		err2 = workflowmd.DeleteWorkflowsByRepoId(ctx, reqDTO.RepoId)
		if err2 != nil {
			return err2
		}
		// 删除工作流密钥
		err2 = workflowmd.DeleteVarsByRepoId(ctx, reqDTO.RepoId)
		if err2 != nil {
			return err2
		}
		// 删除保护分支
		err2 = branchmd.DeleteProtectedBranchByRepoId(ctx, reqDTO.RepoId)
		if err2 != nil {
			return err2
		}
		return nil
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	notifyEvent(repo, team, reqDTO.Operator, event.RepoEventDeletePermanentlyAction)
	return nil
}

func notifyEvent(repo repomd.Repo, team teammd.Team, operator apisession.UserInfo, action event.RepoEventAction) {
	initPsub()
	psub.Publish(event.GitRepoTopic, event.GitRepoEvent{
		BaseTeam: event.BaseTeam{
			TeamId:   team.Id,
			TeamName: team.Name,
		},
		BaseRepo: event.BaseRepo{
			RepoPath: repo.Path,
			RepoId:   repo.Id,
			RepoName: repo.Name,
		},
		BaseEvent: event.BaseEvent{
			Operator:     operator.Account,
			OperatorName: operator.Name,
			EventTime:    time.Now().Format(time.DateTime),
			ActionName:   i18n.GetByLangAndValue(i18n.ZH_CN, action.GetI18nValue()),
			ActionNameEn: i18n.GetByLangAndValue(i18n.EN_US, action.GetI18nValue()),
		},
		Action: action,
	})
}

func notifyGitPushDeleteEvent(repo repomd.Repo, team teammd.Team, operator apisession.UserInfo, ref, refType string) {
	initPsub()
	action := event.GitPushEventDeleteAction
	psub.Publish(event.GitPushTopic+"-remote", event.GitPushEvent{
		RefType: refType,
		Ref:     ref,
		Action:  action,
		BaseTeam: event.BaseTeam{
			TeamId:   team.Id,
			TeamName: team.Name,
		},
		BaseRepo: event.BaseRepo{
			RepoPath: repo.Path,
			RepoId:   repo.Id,
			RepoName: repo.Name,
		},
		BaseEvent: event.BaseEvent{
			Operator:     operator.Account,
			OperatorName: operator.Name,
			EventTime:    time.Now().Format(time.DateTime),
			ActionName:   i18n.GetByLangAndValue(i18n.ZH_CN, action.GetI18nValue()),
			ActionNameEn: i18n.GetByLangAndValue(i18n.EN_US, action.GetI18nValue()),
		},
	})
}

// AllBranches 仓库所有分支
func AllBranches(ctx context.Context, reqDTO AllBranchesReqDTO) ([]string, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	repo, b, err := repomd.GetByRepoId(ctx, reqDTO.RepoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	if !b {
		return nil, util.InvalidArgsError()
	}
	// 校验权限
	err = checkPermByRepo(ctx, repo, reqDTO.Operator, accessRepo)
	if err != nil {
		return nil, err
	}
	branches, err := client.GetAllBranches(ctx, reqvo.GetAllBranchesReq{
		RepoPath: repo.Path,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.Map(branches, func(t reqvo.RefVO) (string, error) {
		return t.Name, nil
	})
}

// DeleteBranch 删除分支
func DeleteBranch(ctx context.Context, reqDTO DeleteBranchReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	repo, b, err := repomd.GetByRepoId(ctx, reqDTO.RepoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.InvalidArgsError()
	}
	team, b, err := teammd.GetByTeamId(ctx, repo.TeamId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.ThereHasBugErr()
	}
	err = checkPermByRepo(ctx, repo, reqDTO.Operator, updateRepo)
	if err != nil {
		return err
	}
	branches, err := branchmd.ListProtectedBranch(ctx, reqDTO.RepoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	isProtectedBranch, _ := branches.IsProtectedBranch(reqDTO.Branch)
	if isProtectedBranch {
		return util.InvalidArgsError()
	}
	err = client.DeleteBranch(ctx, reqvo.DeleteBranchReq{
		RepoPath: repo.Path,
		Branch:   reqDTO.Branch,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	notifyGitPushDeleteEvent(repo, team, reqDTO.Operator, git.BranchPrefix+reqDTO.Branch, "branch")
	return nil
}

// DeleteTag 删除tag
func DeleteTag(ctx context.Context, reqDTO DeleteTagReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	repo, b, err := repomd.GetByRepoId(ctx, reqDTO.RepoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.InvalidArgsError()
	}
	team, b, err := teammd.GetByTeamId(ctx, repo.TeamId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.ThereHasBugErr()
	}
	err = checkPermByRepo(ctx, repo, reqDTO.Operator, updateRepo)
	if err != nil {
		return err
	}
	err = client.DeleteTag(ctx, reqvo.DeleteTagReqVO{
		RepoPath: repo.Path,
		Tag:      reqDTO.Tag,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	notifyGitPushDeleteEvent(repo, team, reqDTO.Operator, git.TagPrefix+reqDTO.Tag, "tag")
	return nil
}

// AllTags 仓库所有tag
func AllTags(ctx context.Context, reqDTO AllTagsReqDTO) ([]string, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	repo, b, err := repomd.GetByRepoId(ctx, reqDTO.RepoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	if !b {
		return nil, util.InvalidArgsError()
	}
	// 校验权限
	err = checkPermByRepo(ctx, repo, reqDTO.Operator, accessRepo)
	if err != nil {
		return nil, err
	}
	tags, err := client.GetAllTags(ctx, reqvo.GetAllTagsReq{
		RepoPath: repo.Path,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.Map(tags, func(t reqvo.RefVO) (string, error) {
		return t.Name, nil
	})
}

// Gc git gc
func Gc(ctx context.Context, reqDTO GcReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	repo, b, err := repomd.GetByRepoId(ctx, reqDTO.RepoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.InvalidArgsError()
	}
	// 管理员权限
	_, err = checkTeamAdminByTeamId(ctx, repo.TeamId, reqDTO.Operator)
	if err != nil {
		return err
	}
	gc, err := client.Gc(ctx, reqvo.GcReq{
		RepoPath: repo.Path,
	})
	if err != nil {
		if bizerr.IsBizErr(err) {
			return err
		}
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	repomd.UpdateGitSize(ctx, reqDTO.RepoId, gc.GitSize)
	return nil
}

// DiffRefs 比较分支或tag的不同
func DiffRefs(ctx context.Context, reqDTO DiffRefsReqDTO) (DiffRefsRespDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return DiffRefsRespDTO{}, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	repo, b, err := repomd.GetByRepoId(ctx, reqDTO.RepoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return DiffRefsRespDTO{}, util.InternalError(err)
	}
	if !b {
		return DiffRefsRespDTO{}, util.InvalidArgsError()
	}
	// 校验权限
	err = checkPermByRepo(ctx, repo, reqDTO.Operator, accessRepo)
	if err != nil {
		return DiffRefsRespDTO{}, err
	}
	refs, err := client.DiffRefs(ctx, reqvo.DiffRefsReq{
		RepoPath:   repo.Path,
		Target:     reqDTO.Target,
		TargetType: reqDTO.TargetType,
		Head:       reqDTO.Head,
		HeadType:   reqDTO.HeadType,
	})
	if err != nil {
		if bizerr.IsBizErr(err) {
			return DiffRefsRespDTO{}, err
		}
		logger.Logger.WithContext(ctx).Error(err)
		return DiffRefsRespDTO{}, util.InternalError(err)
	}
	ret := DiffRefsRespDTO{
		Target:       refs.Target,
		Head:         refs.Head,
		TargetCommit: commit2Dto(refs.TargetCommit),
		HeadCommit:   commit2Dto(refs.HeadCommit),
		NumFiles:     refs.NumFiles,
		DiffNumsStats: DiffNumsStatInfoDTO{
			FileChangeNums: refs.DiffNumsStats.FileChangeNums,
			InsertNums:     refs.DiffNumsStats.InsertNums,
			DeleteNums:     refs.DiffNumsStats.DeleteNums,
		},
		ConflictFiles: refs.ConflictFiles,
	}
	ret.DiffNumsStats.Stats = listutil.MapNe(refs.DiffNumsStats.Stats, func(t reqvo.DiffNumsStatVO) DiffNumsStatDTO {
		return DiffNumsStatDTO{
			RawPath:    t.RawPath,
			Path:       t.Path,
			InsertNums: t.InsertNums,
			DeleteNums: t.DeleteNums,
		}
	})
	ret.Commits = listutil.MapNe(refs.Commits, commit2Dto)
	emails := make([]string, 0)
	emails = append(emails, ret.TargetCommit.Committer.Email)
	emails = append(emails, ret.HeadCommit.Committer.Email)
	for _, commit := range ret.Commits {
		emails = append(emails, commit.Committer.Email)
	}
	userMap, err := usersrv.GetUsersNameAndAvatarMapByEmails(ctx, emails...)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return DiffRefsRespDTO{}, util.InternalError(err)
	}
	// target
	ret.TargetCommit.Committer.AvatarUrl = userMap[ret.TargetCommit.Committer.Email].AvatarUrl
	ret.TargetCommit.Committer.Name = userMap[ret.TargetCommit.Committer.Email].Name
	if userMap[ret.TargetCommit.Committer.Email].Account != "" {
		ret.TargetCommit.Committer.Account = userMap[ret.TargetCommit.Committer.Email].Account
	}
	// head
	ret.HeadCommit.Committer.AvatarUrl = userMap[ret.HeadCommit.Committer.Email].AvatarUrl
	ret.HeadCommit.Committer.Name = userMap[ret.HeadCommit.Committer.Email].Name
	if userMap[ret.HeadCommit.Committer.Email].Account != "" {
		ret.HeadCommit.Committer.Account = userMap[ret.HeadCommit.Committer.Email].Account
	}
	// commits
	for i := range ret.Commits {
		ret.Commits[i].Committer.AvatarUrl = userMap[ret.Commits[i].Committer.Email].AvatarUrl
		ret.Commits[i].Committer.Name = userMap[ret.Commits[i].Committer.Email].Name
		if userMap[ret.Commits[i].Committer.Email].Account != "" {
			ret.Commits[i].Committer.Account = userMap[ret.Commits[i].Committer.Email].Account
		}
	}
	ret.CanMerge = refs.CanMerge
	return ret, nil
}

// DiffCommits 比较commits不同
func DiffCommits(ctx context.Context, reqDTO DiffCommitsReqDTO) (DiffCommitsRespDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return DiffCommitsRespDTO{}, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	repo, b, err := repomd.GetByRepoId(ctx, reqDTO.RepoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return DiffCommitsRespDTO{}, util.InternalError(err)
	}
	if !b {
		return DiffCommitsRespDTO{}, util.InvalidArgsError()
	}
	// 校验权限
	err = checkPermByRepo(ctx, repo, reqDTO.Operator, accessRepo)
	if err != nil {
		return DiffCommitsRespDTO{}, err
	}
	refs, err := client.DiffCommits(ctx, reqvo.DiffCommitsReq{
		RepoPath: repo.Path,
		CommitId: reqDTO.CommitId,
	})
	if err != nil {
		if bizerr.IsBizErr(err) {
			return DiffCommitsRespDTO{}, err
		}
		logger.Logger.WithContext(ctx).Error(err)
		return DiffCommitsRespDTO{}, util.InternalError(err)
	}
	userMap, err := usersrv.GetUsersNameAndAvatarMapByEmails(ctx, refs.Commit.Committer.Email)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return DiffCommitsRespDTO{}, util.InternalError(err)
	}
	commit := commit2Dto(refs.Commit)
	commit.Committer.AvatarUrl = userMap[refs.Commit.Committer.Email].AvatarUrl
	commit.Committer.Name = userMap[refs.Commit.Committer.Email].Name
	commit.Committer.Account = userMap[refs.Commit.Committer.Email].Account
	ret := DiffCommitsRespDTO{
		Commit:   commit,
		NumFiles: refs.NumFiles,
		DiffNumsStats: DiffNumsStatInfoDTO{
			FileChangeNums: refs.DiffNumsStats.FileChangeNums,
			InsertNums:     refs.DiffNumsStats.InsertNums,
			DeleteNums:     refs.DiffNumsStats.DeleteNums,
		},
	}
	ret.DiffNumsStats.Stats = listutil.MapNe(refs.DiffNumsStats.Stats, func(t reqvo.DiffNumsStatVO) DiffNumsStatDTO {
		return DiffNumsStatDTO{
			RawPath:    t.RawPath,
			Path:       t.Path,
			InsertNums: t.InsertNums,
			DeleteNums: t.DeleteNums,
		}
	})
	return ret, nil
}

// Blame 获取每一行提交信息
func Blame(ctx context.Context, reqDTO BlameReqDTO) ([]BlameLineDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	repo, b, err := repomd.GetByRepoId(ctx, reqDTO.RepoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	if !b {
		return nil, util.InvalidArgsError()
	}
	// 校验权限
	err = checkPermByRepo(ctx, repo, reqDTO.Operator, accessRepo)
	if err != nil {
		return nil, err
	}
	lines, err := client.Blame(ctx, reqvo.BlameReq{
		RepoPath: repo.Path,
		Ref:      reqDTO.Ref,
		FilePath: reqDTO.FilePath,
		RefType:  reqDTO.RefType,
	})
	if err != nil {
		if bizerr.IsBizErr(err) {
			return nil, err
		}
		logger.Logger.WithContext(ctx).Error(err)
		return nil, util.InternalError(err)
	}
	return listutil.Map(lines, func(t reqvo.BlameLineVO) (BlameLineDTO, error) {
		return BlameLineDTO{
			Number: t.Number,
			Commit: commit2Dto(t.Commit),
		}, nil
	})
}

func DiffFile(ctx context.Context, reqDTO DiffFileReqDTO) (DiffFileRespDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return DiffFileRespDTO{}, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	repo, b, err := repomd.GetByRepoId(ctx, reqDTO.RepoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return DiffFileRespDTO{}, util.InternalError(err)
	}
	if !b {
		return DiffFileRespDTO{}, util.InvalidArgsError()
	}
	// 校验权限
	err = checkPermByRepo(ctx, repo, reqDTO.Operator, accessRepo)
	if err != nil {
		return DiffFileRespDTO{}, err
	}
	resp, err := client.DiffFile(ctx, reqvo.DiffFileReq{
		RepoPath: repo.Path,
		Target:   reqDTO.Target,
		Head:     reqDTO.Head,
		FilePath: reqDTO.FilePath,
	})
	if err != nil {
		if bizerr.IsBizErr(err) {
			return DiffFileRespDTO{}, err
		}
		logger.Logger.WithContext(ctx).Error(err)
		return DiffFileRespDTO{}, util.InternalError(err)
	}
	ret := DiffFileRespDTO{
		FilePath:    resp.FilePath,
		OldMode:     resp.OldMode,
		Mode:        resp.Mode,
		IsSubModule: resp.IsSubModule,
		FileType:    resp.FileType,
		IsBinary:    resp.IsBinary,
		RenameFrom:  resp.RenameFrom,
		RenameTo:    resp.RenameTo,
		CopyFrom:    resp.CopyFrom,
		CopyTo:      resp.CopyTo,
	}
	ret.Lines = listutil.MapNe(resp.Lines, func(t reqvo.DiffLineVO) DiffLineDTO {
		return DiffLineDTO{
			LeftNo:  t.LeftNo,
			Prefix:  t.Prefix,
			RightNo: t.RightNo,
			Text:    t.Text,
		}
	})
	return ret, nil
}

func commit2Dto(c reqvo.CommitVO) CommitDTO {
	return CommitDTO{
		Parent: c.Parent,
		Author: UserDTO{
			Account: c.Author.Account,
			Email:   c.Author.Email,
		},
		Committer: UserDTO{
			Account: c.Committer.Account,
			Email:   c.Committer.Email,
		},
		AuthoredTime:  c.AuthoredTime,
		CommittedTime: c.CommittedTime,
		CommitMsg:     c.CommitMsg,
		CommitId:      c.CommitId,
		ShortId:       util.LongCommitId2ShortId(c.CommitId),
		Tagger: UserDTO{
			Account: c.Tagger.Account,
			Email:   c.Tagger.Email,
		},
		TaggerTime:   c.TaggerTime,
		ShortTagId:   c.ShortTagId,
		TagCommitMsg: c.TagCommitMsg,
	}
}

// HistoryCommits 获取提交历史
func HistoryCommits(ctx context.Context, reqDTO HistoryCommitsReqDTO) (HistoryCommitsRespDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return HistoryCommitsRespDTO{}, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	repo, b, err := repomd.GetByRepoId(ctx, reqDTO.RepoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return HistoryCommitsRespDTO{}, util.InternalError(err)
	}
	if !b {
		return HistoryCommitsRespDTO{}, util.InvalidArgsError()
	}
	// 校验权限
	err = checkPermByRepo(ctx, repo, reqDTO.Operator, accessRepo)
	if err != nil {
		return HistoryCommitsRespDTO{}, err
	}
	resp, err := client.HistoryCommits(ctx, reqvo.HistoryCommitsReq{
		RepoPath: repo.Path,
		Ref:      reqDTO.Ref,
		Offset:   reqDTO.Cursor,
	})
	if err != nil {
		if bizerr.IsBizErr(err) {
			return HistoryCommitsRespDTO{}, err
		}
		logger.Logger.WithContext(ctx).Error(err)
		return HistoryCommitsRespDTO{}, err
	}
	// 缓存 用于校验签名
	gpgMap := make(map[string][]gpgkeymd.GpgKey)
	gpgIdMap := make(map[string]gpgkeymd.GpgKey)
	sshMap := make(map[string][]sshkeysrv.InnerSshKeyDTO)
	ret := HistoryCommitsRespDTO{
		Cursor: resp.Cursor,
	}
	ret.Data = listutil.MapNe(resp.Data, func(t reqvo.CommitVO) CommitDTO {
		r := commit2Dto(t)
		if t.CommitSig != "" {
			sig := signature.Sig(t.CommitSig)
			if sig.IsSSHSig() {
				r.Verified, r.Signer.Account, r.Signer.Key, r.Signer.Type = verifyPayloadWithSshKeys(ctx, t.Committer.Account, t.CommitSig, t.Payload, sshMap)
			} else if sig.IsGPGSig() {
				r.Verified, r.Signer.Account, r.Signer.Key, r.Signer.Type = verifyPayloadWithGpgKeys(ctx, t.Committer.Account, t.CommitSig, t.Payload, gpgMap, gpgIdMap)
			}
		}
		return r
	})
	// 查找头像和姓名
	accountList := make([]string, 0)
	for _, commit := range ret.Data {
		if commit.Signer.Account != "" {
			accountList = append(accountList, commit.Signer.Account)
		}
	}
	if len(accountList) > 0 {
		userMap, err := usersrv.GetUsersNameAndAvatarMap(ctx, accountList...)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return HistoryCommitsRespDTO{}, err
		}
		for i := range ret.Data {
			if ret.Data[i].Signer.Account != "" {
				user := userMap[ret.Data[i].Signer.Account]
				ret.Data[i].Signer.Name = user.Name
				ret.Data[i].Signer.AvatarUrl = user.AvatarUrl
			}
		}
	}
	emailList := listutil.MapNe(ret.Data, func(t CommitDTO) string {
		return t.Committer.Email
	})
	if len(emailList) > 0 {
		userMap, err := usersrv.GetUsersNameAndAvatarMapByEmails(ctx, emailList...)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return HistoryCommitsRespDTO{}, err
		}
		for i := range ret.Data {
			ret.Data[i].Committer.AvatarUrl = userMap[ret.Data[i].Committer.Email].AvatarUrl
			ret.Data[i].Committer.Name = userMap[ret.Data[i].Committer.Email].Name
			if userMap[ret.Data[i].Committer.Email].Account != "" {
				ret.Data[i].Committer.Account = userMap[ret.Data[i].Committer.Email].Account
			}
		}
	}
	return ret, nil
}

func verifyPayloadWithSshKeys(ctx context.Context, account, sign, payload string, sshMap map[string][]sshkeysrv.InnerSshKeyDTO) (bool, string, string, string) {
	keys, b := sshMap[account]
	if !b {
		keys = sshkeysrv.ListAllPubKeyByAccount(ctx, account)
		sshMap[account] = keys
	}
	for _, key := range keys {
		err := signature.VerifySshSignature(sign, payload, key.Content)
		if err == nil {
			return true, account, key.Fingerprint, "SSH"
		}
	}
	return false, "", "", ""
}

func verifyPayloadWithGpgKeys(ctx context.Context, account, sign, payload string, gpgKeysMap map[string][]gpgkeymd.GpgKey, gpgKeyIdMap map[string]gpgkeymd.GpgKey) (bool, string, string, string) {
	sig, err := signature.ExtractGpgSignature(sign)
	if err != nil {
		return false, "", "", ""
	}
	// 从gpgKeys匹配keyId
	{
		keyId := ""
		if sig.IssuerKeyId != nil && (*sig.IssuerKeyId) != 0 {
			keyId = fmt.Sprintf("%X", *sig.IssuerKeyId)
		}
		if keyId == "" && sig.IssuerFingerprint != nil && len(sig.IssuerFingerprint) > 0 {
			keyId = fmt.Sprintf("%X", sig.IssuerFingerprint[12:20])
		}
		if keyId != "" {
			key, b := gpgKeyIdMap[keyId]
			if !b {
				key, _ = gpgkeysrv.GetByKeyId(ctx, keyId)
				gpgKeyIdMap[keyId] = key
			}
			return verifyPayloadWithGpgKey(&key, sig, payload), key.Account, key.KeyId, "GPG"
		}
	}
	// 匹配committer
	{
		if account != "" {
			keys, b := gpgKeysMap[account]
			if !b {
				keys = gpgkeysrv.ListValidByAccount(ctx, account)
				gpgKeysMap[account] = keys
			}
			for _, key := range keys {
				if verifyPayloadWithGpgKey(&key, sig, payload) {
					return true, key.Account, key.KeyId, "GPG"
				}
			}
		}
	}
	return false, "", "", ""
}

func verifyPayloadWithGpgKey(gpgKey *gpgkeymd.GpgKey, sig *packet.Signature, payload string) bool {
	if gpgKey.KeyId == "" {
		return false
	}
	pk, err := signature.Base64DecGPGPubKey(gpgKey.Content)
	if err == nil && pk.CanSign() {
		hash := sig.Hash.New()
		_, err = hash.Write([]byte(payload))
		if err != nil {
			return false
		}
		err = pk.VerifySignature(hash, sig)
		if err == nil {
			return true
		}
	}
	for _, subKey := range gpgKey.SubKeys {
		pk, err = signature.Base64DecGPGPubKey(subKey.Content)
		if err == nil && pk.CanSign() {
			hash := sig.Hash.New()
			_, err = hash.Write([]byte(payload))
			if err != nil {
				return false
			}
			err = pk.VerifySignature(hash, sig)
			if err == nil {
				return true
			}
		}
	}
	return false
}

func checkTeamAdminByTeamId(ctx context.Context, teamId int64, operator apisession.UserInfo) (teammd.Team, error) {
	team, b, err := teammd.GetByTeamId(ctx, teamId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return teammd.Team{}, util.InternalError(err)
	}
	if !b {
		return teammd.Team{}, util.InvalidArgsError()
	}
	if operator.IsAdmin {
		return team, nil
	}
	p, b, err := teammd.GetUserPermDetail(ctx, teamId, operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return team, util.InternalError(err)
	}
	if !b || !p.IsAdmin {
		return team, util.UnauthorizedError()
	}
	return team, nil
}

func checkTeamAdminByRepoId(ctx context.Context, repoId int64, operator apisession.UserInfo) (repomd.Repo, teammd.Team, error) {
	repo, b, err := repomd.GetByRepoId(ctx, repoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return repomd.Repo{}, teammd.Team{}, util.InternalError(err)
	}
	if !b {
		return repomd.Repo{}, teammd.Team{}, util.InvalidArgsError()
	}
	team, err := checkTeamAdminByTeamId(ctx, repo.TeamId, operator)
	return repo, team, err
}

func checkPermByRepo(ctx context.Context, repo repomd.Repo, operator apisession.UserInfo, permCode int) error {
	if operator.IsAdmin {
		return nil
	}
	p, b, err := teammd.GetUserPermDetail(ctx, repo.TeamId, operator.Account)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.UnauthorizedError()
	}
	if p.IsAdmin {
		return nil
	}
	pass := false
	switch permCode {
	case accessRepo:
		pass = p.PermDetail.GetRepoPerm(repo.Id).CanAccessRepo
	case updateRepo:
		pass = p.PermDetail.GetRepoPerm(repo.Id).CanPushRepo
	}
	if pass {
		return nil
	}
	return util.UnauthorizedError()
}

// TransferTeam 迁移团队
func TransferTeam(ctx context.Context, reqDTO TransferTeamReqDTO) (err error) {
	if err = reqDTO.IsValid(); err != nil {
		return
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	repo, b, err := repomd.GetByRepoId(ctx, reqDTO.RepoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.InvalidArgsError()
	}
	// 系统管理员权限
	if !reqDTO.Operator.IsAdmin {
		return util.InvalidArgsError()
	}
	if repo.TeamId == reqDTO.TeamId {
		return nil
	}
	b, err = teammd.ExistByTeamId(ctx, reqDTO.TeamId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.InvalidArgsError()
	}
	_, err = repomd.TransferTeam(ctx, reqDTO.RepoId, reqDTO.TeamId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		err = util.InternalError(err)
		return
	}
	return
}

// ListBranchCommits 分页获取分支+提交信息
func ListBranchCommits(ctx context.Context, reqDTO ListRefCommitsReqDTO) ([]BranchCommitDTO, int64, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, 0, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	repo, b, err := repomd.GetByRepoId(ctx, reqDTO.RepoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	if !b {
		return nil, 0, util.InvalidArgsError()
	}
	// 校验权限
	err = checkPermByRepo(ctx, repo, reqDTO.Operator, accessRepo)
	if err != nil {
		return nil, 0, err
	}
	branches, totalCount, err := client.ListBranchAndLastCommit(ctx, reqvo.ListRefCommitsReq{
		RepoPath: repo.Path,
		PageNum:  reqDTO.PageNum,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	var prMap map[string]pullrequestmd.PullRequest
	if len(branches) > 0 {
		heads := listutil.MapNe(branches, func(t reqvo.RefCommitVO) string {
			return t.Name
		})
		pullRequests, err := pullrequestmd.GetLastPullRequestByRepoIdAndHead(ctx, reqDTO.RepoId, heads)
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return nil, 0, util.InternalError(err)
		}
		prMap, _ = listutil.CollectToMap(pullRequests, func(t pullrequestmd.PullRequest) (string, error) {
			return t.Head, nil
		}, func(t pullrequestmd.PullRequest) (pullrequestmd.PullRequest, error) {
			return t, nil
		})
	} else {
		prMap = map[string]pullrequestmd.PullRequest{}
	}
	pbList, err := branchmd.ListProtectedBranch(ctx, reqDTO.RepoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	data := listutil.MapNe(branches, func(t reqvo.RefCommitVO) BranchCommitDTO {
		pr, b := prMap[t.Name]
		isProtectedBranch, _ := pbList.IsProtectedBranch(t.Name)
		ret := BranchCommitDTO{
			Name:              t.Name,
			LastCommit:        commit2Dto(t.Commit),
			IsProtectedBranch: isProtectedBranch,
		}
		if b {
			ret.LastPullRequest = &PullRequestDTO{
				Id:       pr.Id,
				PrStatus: pr.PrStatus,
				PrTitle:  pr.PrTitle,
				Created:  pr.Created,
				PrIndex:  pr.PrIndex,
			}
		}
		return ret
	})
	return data, totalCount, nil
}

// ListTagCommits 分页获取tag+提交信息
func ListTagCommits(ctx context.Context, reqDTO ListRefCommitsReqDTO) ([]TagCommitDTO, int64, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, 0, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	repo, b, err := repomd.GetByRepoId(ctx, reqDTO.RepoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	if !b {
		return nil, 0, util.InvalidArgsError()
	}
	// 校验权限
	err = checkPermByRepo(ctx, repo, reqDTO.Operator, accessRepo)
	if err != nil {
		return nil, 0, err
	}
	tags, totalCount, err := client.ListTagAndCommit(ctx, reqvo.ListRefCommitsReq{
		RepoPath: repo.Path,
		PageNum:  reqDTO.PageNum,
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	// 缓存 用于校验签名
	gpgMap := make(map[string][]gpgkeymd.GpgKey)
	gpgIdMap := make(map[string]gpgkeymd.GpgKey)
	sshMap := make(map[string][]sshkeysrv.InnerSshKeyDTO)
	data := listutil.MapNe(tags, func(t reqvo.RefCommitVO) TagCommitDTO {
		r := commit2Dto(t.Commit)
		if t.Commit.TagSig != "" {
			sig := signature.Sig(t.Commit.TagSig)
			if sig.IsSSHSig() {
				r.Verified, r.Signer.Account, r.Signer.Key, r.Signer.Type = verifyPayloadWithSshKeys(ctx, t.Commit.Tagger.Account, t.Commit.TagSig, t.Commit.TagPayload, sshMap)
			} else if sig.IsGPGSig() {
				r.Verified, r.Signer.Account, r.Signer.Key, r.Signer.Type = verifyPayloadWithGpgKeys(ctx, t.Commit.Tagger.Account, t.Commit.TagSig, t.Commit.TagPayload, gpgMap, gpgIdMap)
			}
		}
		return TagCommitDTO{
			Name:   t.Name,
			Commit: r,
		}
	})
	// 查找头像和姓名
	accountList := make([]string, 0)
	for _, tag := range data {
		if tag.Commit.Signer.Account != "" {
			accountList = append(accountList, tag.Commit.Signer.Account)
		}
	}
	avatarMap, err := usersrv.GetUsersNameAndAvatarMap(ctx, accountList...)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, 0, util.InternalError(err)
	}
	for i := range data {
		if data[i].Commit.Signer.Account != "" {
			user := avatarMap[data[i].Commit.Signer.Account]
			data[i].Commit.Signer.Name = user.Name
			data[i].Commit.Signer.AvatarUrl = user.AvatarUrl
		}
	}
	return data, totalCount, nil
}

// CreateArchive 下载代码
func CreateArchive(ctx context.Context, reqDTO CreateArchiveReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	repo, b, err := repomd.GetByRepoId(ctx, reqDTO.RepoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.InvalidArgsError()
	}
	// 校验权限
	err = checkPermByRepo(ctx, repo, reqDTO.Operator, accessRepo)
	if err != nil {
		return err
	}
	initCreateArchiveLimiter()
	if createArchiveLimiter.Borrow() {
		defer createArchiveLimiter.Return()
		err = client.CreateArchive(reqvo.CreateArchiveReq{
			RepoPath: repo.Path,
			FileName: reqDTO.FileName,
			C:        reqDTO.C,
		})
		if err != nil {
			logger.Logger.WithContext(ctx).Error(err)
			return util.InternalError(err)
		}
		return nil
	}
	return util.NewBizErrWithMsg(apicode.TooManyOperationCode, i18n.GetByKey(i18n.SystemTooManyOperation))
}

// UpdateRepo 更新仓库配置
func UpdateRepo(ctx context.Context, reqDTO UpdateRepoReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	repo, b, err := repomd.GetByRepoId(ctx, reqDTO.RepoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.InvalidArgsError()
	}
	team, err := checkTeamAdminByTeamId(ctx, repo.TeamId, reqDTO.Operator)
	if err != nil {
		return err
	}
	_, err = repomd.UpdateRepo(ctx, repomd.UpdateRepoReqDTO{
		Id:       reqDTO.RepoId,
		RepoDesc: reqDTO.Desc,
		Cfg: repomd.RepoCfg{
			DisableLfs:   reqDTO.DisableLfs,
			LfsLimitSize: reqDTO.LfsLimitSize,
			GitLimitSize: reqDTO.GitLimitSize,
		},
	})
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	notifyEvent(repo, team, reqDTO.Operator, event.RepoEventUpdateAction)
	return nil
}

// SetRepoArchivedStatus 归档或非归档仓库
func SetRepoArchivedStatus(ctx context.Context, reqDTO SetRepoArchivedStatusReqDTO) error {
	if err := reqDTO.IsValid(); err != nil {
		return err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	repo, b, err := repomd.GetByRepoId(ctx, reqDTO.RepoId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if !b {
		return util.InvalidArgsError()
	}
	team, err := checkTeamAdminByTeamId(ctx, repo.TeamId, reqDTO.Operator)
	if err != nil {
		return err
	}
	_, err = repomd.UpdateRepoIsArchived(ctx, reqDTO.RepoId, reqDTO.IsArchived)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return util.InternalError(err)
	}
	if reqDTO.IsArchived {
		notifyEvent(repo, team, reqDTO.Operator, event.RepoEventArchivedAction)
	} else {
		notifyEvent(repo, team, reqDTO.Operator, event.RepoEventUnArchivedAction)
	}
	return nil
}

// ListRepoByAdmin 管理员展示仓库列表
func ListRepoByAdmin(ctx context.Context, reqDTO ListRepoByAdminReqDTO) ([]SimpleRepoDTO, error) {
	if err := reqDTO.IsValid(); err != nil {
		return nil, err
	}
	ctx, closer := xormstore.Context(ctx)
	defer closer.Close()
	_, err := checkTeamAdminByTeamId(ctx, reqDTO.TeamId, reqDTO.Operator)
	if err != nil {
		return nil, err
	}
	repos, err := repomd.ListRepoByTeamId(ctx, reqDTO.TeamId)
	if err != nil {
		logger.Logger.WithContext(ctx).Error(err)
		return nil, err
	}
	return listutil.Map(repos, func(t repomd.Repo) (SimpleRepoDTO, error) {
		return SimpleRepoDTO{
			RepoId: t.Id,
			Name:   t.Name,
		}, nil
	})
}
