package repoapi

import (
	"github.com/LeeZXin/zsf-utils/ginutil"
)

type InitRepoReqVO struct {
	Name          string `json:"name"`
	Desc          string `json:"Desc"`
	RepoType      int    `json:"repoType"`
	CreateReadme  bool   `json:"createReadme"`
	TeamId        int64  `json:"teamId"`
	NodeId        int64  `json:"nodeId"`
	GitIgnoreName string `json:"gitIgnoreName"`
	DefaultBranch string `json:"defaultBranch"`
}

type DeleteRepoReqVO struct {
	Id int64 `json:"id"`
}

type TreeRepoReqVO struct {
	Id  int64  `json:"id"`
	Ref string `json:"ref"`
	Dir string `json:"dir"`
}

type EntriesRepoReqVO struct {
	Id     int64  `json:"id"`
	Ref    string `json:"ref"`
	Dir    string `json:"dir"`
	Offset int    `json:"offset"`
}

type ListRepoReqVO struct {
	TeamId int64 `json:"teamId"`
}

type ListRepoRespVO struct {
	ginutil.BaseResp
	RepoList []RepoVO `json:"repoList"`
	Cursor   int64    `json:"cursor"`
	Limit    int      `json:"limit"`
}

type UserVO struct {
	Account string `json:"account"`
	Email   string `json:"email"`
}

type CommitVO struct {
	Author        UserVO `json:"author"`
	Committer     UserVO `json:"committer"`
	AuthoredTime  string `json:"authoredTime"`
	CommittedTime string `json:"committedTime"`
	CommitMsg     string `json:"commitMsg"`
	CommitId      string `json:"commitId"`
	ShortId       string `json:"shortId"`
	Verified      bool   `json:"verified"`
}

type FileVO struct {
	Mode    string   `json:"mode"`
	RawPath string   `json:"rawPath"`
	Path    string   `json:"path"`
	Commit  CommitVO `json:"commit"`
}

type TreeVO struct {
	Files   []FileVO `json:"files"`
	Limit   int      `json:"limit"`
	Offset  int      `json:"offset"`
	HasMore bool     `json:"hasMore"`
}

type TreeRepoRespVO struct {
	ginutil.BaseResp
	ReadmeText   string   `json:"readmeText"`
	LatestCommit CommitVO `json:"latestCommit"`
	Tree         TreeVO   `json:"tree"`
}

type RepoVO struct {
	Id      int64  `json:"id"`
	Name    string `json:"name"`
	Path    string `json:"path"`
	Author  string `json:"author"`
	TeamId  int64  `json:"teamId"`
	GitSize int64  `json:"gitSize"`
	LfsSize int64  `json:"lfsSize"`
	Created string `json:"created"`
}

type CatFileReqVO struct {
	Id       int64  `json:"id"`
	Ref      string `json:"ref"`
	Dir      string `json:"dir"`
	FileName string `json:"fileName"`
}

type CatFileRespVO struct {
	ginutil.BaseResp
	Mode    string `json:"mode"`
	Content string `json:"content"`
}

type RepoTypeVO struct {
	Option int    `json:"option"`
	Name   string `json:"name"`
}

type AllBranchesReqVO struct {
	Id int64 `json:"id"`
}

type AllTagsReqVO struct {
	Id int64 `json:"id"`
}

type GcReqVO struct {
	Id int64 `json:"id"`
}

type PrepareMergeReqVO struct {
	Id     int64  `json:"id"`
	Target string `json:"target"`
	Head   string `json:"head"`
}

type DiffFileReqVO struct {
	Id       int64  `json:"id"`
	Target   string `json:"target"`
	Head     string `json:"head"`
	FileName string `json:"fileName"`
}

type PrepareMergeRespVO struct {
	ginutil.BaseResp
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

type DiffNumsStatInfoVO struct {
	FileChangeNums int              `json:"fileChangeNums"`
	InsertNums     int              `json:"insertNums"`
	DeleteNums     int              `json:"deleteNums"`
	Stats          []DiffNumsStatVO `json:"stats"`
}

type DiffNumsStatVO struct {
	RawPath    string `json:"rawPath"`
	Path       string `json:"path"`
	TotalNums  int    `json:"totalNums"`
	InsertNums int    `json:"insertNums"`
	DeleteNums int    `json:"deleteNums"`
}

type DiffFileRespVO struct {
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
	Index   int    `json:"index"`
	LeftNo  int    `json:"leftNo"`
	Prefix  string `json:"prefix"`
	RightNo int    `json:"rightNo"`
	Text    string `json:"text"`
}

type ShowDiffTextContentReqVO struct {
	Id        int64  `json:"id"`
	CommitId  string `json:"commitId"`
	FileName  string `json:"fileName"`
	Offset    int    `json:"offset"`
	Limit     int    `json:"limit"`
	Direction string `json:"direction"`
}

type HistoryCommitsReqVO struct {
	Id     int64  `json:"id"`
	Ref    string `json:"ref"`
	Cursor int    `json:"cursor"`
}

type HistoryCommitsRespVO struct {
	ginutil.BaseResp
	Data   []CommitVO `json:"data"`
	Cursor int        `json:"cursor"`
}

type RepoTokenVO struct {
	Id      int64  `json:"id"`
	Account string `json:"account"`
	Token   string `json:"token"`
	Created string `json:"created"`
}

type CreateRepoTokenReqVO struct {
	Id int64 `json:"id"`
}

type DeleteRepoTokenReqVO struct {
	Id int64 `json:"id"`
}

type ListRepoTokenReqVO struct {
	Id int64 `json:"id"`
}

type TransferTeam struct {
	Id     int64 `json:"id"`
	TeamId int64 `json:"teamId"`
}
