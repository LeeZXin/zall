package repoapi

import (
	"github.com/LeeZXin/zall/git/modules/model/pullrequestmd"
	"github.com/LeeZXin/zall/pkg/git"
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

type DeleteRepoReqVO struct {
	RepoId int64 `json:"repoId"`
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
	Parent        []string `json:"parent"`
	Author        UserVO   `json:"author"`
	Committer     UserVO   `json:"committer"`
	AuthoredTime  string   `json:"authoredTime"`
	CommittedTime string   `json:"committedTime"`
	CommitMsg     string   `json:"commitMsg"`
	CommitId      string   `json:"commitId"`
	ShortId       string   `json:"shortId"`
	Verified      bool     `json:"verified"`
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
	Author       string `json:"author"`
	RepoDesc     string `json:"repoDesc"`
	GitSize      int64  `json:"gitSize"`
	LfsSize      int64  `json:"lfsSize"`
	Created      string `json:"created"`
	Updated      string `json:"updated"`
	TeamId       int64  `json:"teamId"`
	LastOperated string `json:"lastOperated"`
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

type ShowDiffTextContentReqVO struct {
	RepoId    int64  `json:"repoId"`
	CommitId  string `json:"commitId"`
	FileName  string `json:"fileName"`
	Offset    int    `json:"offset"`
	Limit     int    `json:"limit"`
	Direction string `json:"direction"`
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

type RepoTokenVO struct {
	TokenId int64  `json:"tokenId"`
	Account string `json:"account"`
	Token   string `json:"token"`
	Created string `json:"created"`
}

type CreateRepoTokenReqVO struct {
	RepoId int64 `json:"repoId"`
}

type DeleteRepoTokenReqVO struct {
	TokenId int64 `json:"tokenId"`
}

type ListRepoTokenReqVO struct {
	RepoId int64 `json:"repoId"`
}

type TransferTeam struct {
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
	Name            string         `json:"name"`
	LastCommit      CommitVO       `json:"lastCommit"`
	LastPullRequest *PullRequestVO `json:"lastPullRequest,omitempty"`
}

type DeleteBranchReqVO struct {
	RepoId int64  `json:"repoId"`
	Branch string `json:"branch"`
}
