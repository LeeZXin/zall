package repoapi

import (
	"github.com/LeeZXin/zall/git/modules/model/pullrequestmd"
	"github.com/LeeZXin/zall/pkg/git"
	"github.com/LeeZXin/zall/pkg/perm"
	"github.com/LeeZXin/zsf-utils/ginutil"
)

type CreateRepoReqVO struct {
	Name          string `json:"name"`
	Desc          string `json:"desc"`
	AddReadme     bool   `json:"addReadme"`
	TeamId        int64  `json:"teamId"`
	GitIgnoreName string `json:"gitIgnoreName"`
	DefaultBranch string `json:"defaultBranch"`
}

type IndexRepoReqVO struct {
	RepoId  int64       `json:"repoId"`
	Ref     string      `json:"ref"`
	RefType git.RefType `json:"refType"`
}

type EntriesRepoReqVO struct {
	RepoId  int64       `json:"repoId"`
	Ref     string      `json:"ref"`
	Dir     string      `json:"dir"`
	RefType git.RefType `json:"refType"`
}

type UserVO struct {
	Account string `json:"account"`
	Email   string `json:"email"`
}

type CommitVO struct {
	Parent        []string  `json:"parent"`
	Author        UserVO    `json:"author"`
	Committer     UserVO    `json:"committer"`
	AuthoredTime  string    `json:"authoredTime"`
	CommittedTime string    `json:"committedTime"`
	CommitMsg     string    `json:"commitMsg"`
	CommitId      string    `json:"commitId"`
	ShortId       string    `json:"shortId"`
	Verified      bool      `json:"verified"`
	Tagger        *UserVO   `json:"tagger,omitempty"`
	TaggerTime    *string   `json:"taggerTime,omitempty"`
	ShortTagId    *string   `json:"shortTagId,omitempty"`
	TagCommitMsg  *string   `json:"tagCommitMsg,omitempty"`
	Signer        *SignerVO `json:"signer,omitempty"`
}

type SignerVO struct {
	Account   string `json:"account"`
	Name      string `json:"name"`
	AvatarUrl string `json:"avatarUrl"`
	Key       string `json:"key"`
	Type      string `json:"type"`
}

type FileVO struct {
	Mode    string   `json:"mode"`
	RawPath string   `json:"rawPath"`
	Path    string   `json:"path"`
	Commit  CommitVO `json:"commit"`
}

type TreeVO struct {
	Files []FileVO `json:"files"`
}

type BlobVO struct {
	Mode    string `json:"mode"`
	RawPath string `json:"rawPath"`
	Path    string `json:"path"`
}

type IndexRepoRespVO struct {
	ginutil.BaseResp
	HasReadme    bool     `json:"hasReadme"`
	ReadmeText   string   `json:"readmeText"`
	LatestCommit CommitVO `json:"latestCommit"`
	Tree         TreeVO   `json:"tree"`
}

type RepoVO struct {
	RepoId       int64  `json:"repoId"`
	Name         string `json:"name"`
	Path         string `json:"path"`
	RepoDesc     string `json:"repoDesc"`
	GitSize      int64  `json:"gitSize"`
	LfsSize      int64  `json:"lfsSize"`
	Created      string `json:"created"`
	TeamId       int64  `json:"teamId"`
	LastOperated string `json:"lastOperated"`
	DisableLfs   bool   `json:"disableLfs"`
	LfsLimitSize int64  `json:"lfsLimitSize"`
	GitLimitSize int64  `json:"gitLimitSize"`
	IsArchived   bool   `json:"isArchived"`
}

type DeletedRepoVO struct {
	RepoVO
	Deleted string `json:"deleted"`
}

type CatFileReqVO struct {
	RepoId   int64       `json:"repoId"`
	Ref      string      `json:"ref"`
	FilePath string      `json:"filePath"`
	RefType  git.RefType `json:"refType"`
}

type BlameReqVO struct {
	RepoId   int64       `json:"repoId"`
	Ref      string      `json:"ref"`
	FilePath string      `json:"filePath"`
	RefType  git.RefType `json:"refType"`
}

type CatFileVO struct {
	FileMode string   `json:"fileMode"`
	Content  string   `json:"content"`
	Size     string   `json:"size"`
	Commit   CommitVO `json:"commit"`
}

type SimpleInfoVO struct {
	Branches     []string `json:"branches"`
	Tags         []string `json:"tags"`
	CloneHttpUrl string   `json:"cloneHttpUrl"`
	CloneSshUrl  string   `json:"cloneSshUrl"`
}

type DiffRefsReqVO struct {
	RepoId     int64       `json:"repoId"`
	Target     string      `json:"target"`
	TargetType git.RefType `json:"targetType"`
	Head       string      `json:"head"`
	HeadType   git.RefType `json:"headType"`
}

type DiffCommitsReqVO struct {
	RepoId   int64  `json:"repoId"`
	CommitId string `json:"commitId"`
}

type DiffFileReqVO struct {
	RepoId   int64  `json:"repoId"`
	Target   string `json:"target"`
	Head     string `json:"head"`
	FilePath string `json:"filePath"`
}

type DiffRefsVO struct {
	Target        string             `json:"target"`
	Head          string             `json:"head"`
	TargetCommit  CommitVO           `json:"targetCommit"`
	HeadCommit    CommitVO           `json:"headCommit"`
	Commits       []CommitVO         `json:"commits"`
	NumFiles      int                `json:"numFiles"`
	DiffNumsStats DiffNumsStatInfoVO `json:"diffNumsStats"`
	ConflictFiles []string           `json:"conflictFiles"`
	CanMerge      bool               `json:"canMerge"`
}

type DiffCommitsVO struct {
	Commit        CommitVO           `json:"commit"`
	NumFiles      int                `json:"numFiles"`
	DiffNumsStats DiffNumsStatInfoVO `json:"diffNumsStats"`
}

type DiffNumsStatInfoVO struct {
	FileChangeNums int              `json:"fileChangeNums"`
	InsertNums     int              `json:"insertNums"`
	DeleteNums     int              `json:"deleteNums"`
	Stats          []DiffNumsStatVO `json:"stats"`
}

type DiffNumsStatVO struct {
	RawPath    string `json:"rawPath"`
	Path       string `json:"path"`
	InsertNums int    `json:"insertNums"`
	DeleteNums int    `json:"deleteNums"`
}

type DiffFileVO struct {
	FilePath    string       `json:"filePath"`
	OldMode     string       `json:"oldMode"`
	Mode        string       `json:"mode"`
	IsSubModule bool         `json:"isSubModule"`
	FileType    string       `json:"fileType"`
	IsBinary    bool         `json:"isBinary"`
	RenameFrom  string       `json:"renameFrom"`
	RenameTo    string       `json:"renameTo"`
	CopyFrom    string       `json:"copyFrom"`
	CopyTo      string       `json:"copyTo"`
	Lines       []DiffLineVO `json:"lines"`
}

type DiffLineVO struct {
	LeftNo  int    `json:"leftNo"`
	Prefix  string `json:"prefix"`
	RightNo int    `json:"rightNo"`
	Text    string `json:"text"`
}

type HistoryCommitsReqVO struct {
	RepoId int64  `json:"repoId"`
	Ref    string `json:"ref"`
	Cursor int    `json:"cursor"`
}

type HistoryCommitsRespVO struct {
	ginutil.BaseResp
	Data   []CommitVO `json:"data"`
	Cursor int        `json:"cursor"`
}

type TransferTeamReqVO struct {
	RepoId int64 `json:"repoId"`
	TeamId int64 `json:"teamId"`
}

type BlameLineVO struct {
	Number int      `json:"number"`
	Commit CommitVO `json:"commit"`
}

type PullRequestVO struct {
	Id       int64                  `json:"id"`
	PrStatus pullrequestmd.PrStatus `json:"prStatus"`
	PrTitle  string                 `json:"prTitle"`
	Created  string                 `json:"created"`
}

type BranchCommitVO struct {
	Name              string         `json:"name"`
	IsProtectedBranch bool           `json:"isProtectedBranch"`
	LastCommit        CommitVO       `json:"lastCommit"`
	LastPullRequest   *PullRequestVO `json:"lastPullRequest,omitempty"`
}

type TagCommitVO struct {
	Name   string   `json:"name"`
	Commit CommitVO `json:"commit"`
}

type DeleteBranchReqVO struct {
	RepoId int64  `json:"repoId"`
	Branch string `json:"branch"`
}

type PageRefCommitsReqVO struct {
	RepoId  int64 `json:"repoId"`
	PageNum int   `json:"pageNum"`
}

type CreateArchiveReqVO struct {
	RepoId   int64  `json:"repoId"`
	FileName string `json:"fileName"`
}

type DeleteTagReqVO struct {
	RepoId int64  `json:"repoId"`
	Tag    string `json:"tag"`
}

type GetRepoSizeRespVO struct {
	GitSize int64 `json:"gitSize"`
	LfsSize int64 `json:"lfsSize"`
}

type UpdateRepoReqVO struct {
	RepoId       int64  `json:"repoId"`
	Desc         string `json:"desc"`
	DisableLfs   bool   `json:"disableLfs"`
	LfsLimitSize int64  `json:"lfsLimitSize"`
	GitLimitSize int64  `json:"gitLimitSize"`
}

type SimpleRepoVO struct {
	RepoId int64  `json:"repoId"`
	Name   string `json:"name"`
	TeamId int64  `json:"teamId"`
}

type RepoWithPermVO struct {
	SimpleRepoVO
	Perm perm.RepoPerm `json:"perm"`
}
