package reqvo

import (
	"github.com/LeeZXin/zall/pkg/git"
	"github.com/gin-gonic/gin"
)

type InitRepoReq struct {
	UserAccount   string `json:"userAccount"`
	UserEmail     string `json:"userEmail"`
	RepoName      string `json:"repoName"`
	RepoPath      string `json:"repoPath"`
	AddReadme     bool   `json:"addReadme"`
	GitIgnoreName string `json:"gitIgnoreName"`
	DefaultBranch string `json:"defaultBranch"`
}

type DeleteRepoReq struct {
	RepoPath string `json:"repoPath"`
}

type GetAllBranchesReq struct {
	RepoPath string `json:"repoPath"`
}

type ListRefCommitsReq struct {
	RepoPath string `json:"repoPath"`
	PageNum  int    `json:"pageNum"`
}

type DeleteBranchReq struct {
	RepoPath string `json:"repoPath"`
	Branch   string `json:"branch"`
}

type GetAllTagsReq struct {
	RepoPath string `json:"repoPath"`
}

type GcReq struct {
	RepoPath string `json:"repoPath"`
}

type GcResp struct {
	GitSize int64 `json:"gitSize"`
}

type DiffRefsReq struct {
	RepoPath   string      `json:"repoPath"`
	Target     string      `json:"target"`
	TargetType git.RefType `json:"targetType"`
	Head       string      `json:"head"`
	HeadType   git.RefType `json:"headType"`
}

type DiffCommitsReq struct {
	RepoPath string `json:"repoPath"`
	CommitId string `json:"commitId"`
}

type CanMergeReq struct {
	RepoPath   string      `json:"repoPath"`
	Target     string      `json:"target"`
	TargetType git.RefType `json:"targetType"`
	Head       string      `json:"head"`
	HeadType   git.RefType `json:"headType"`
}

type DiffRefsResp struct {
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

type DiffCommitsResp struct {
	Commit        CommitVO           `json:"commit"`
	NumFiles      int                `json:"numFiles"`
	DiffNumsStats DiffNumsStatInfoVO `json:"diffNumsStats"`
}

type UserVO struct {
	Account string `json:"account"`
	Email   string `json:"email"`
}

type CommitVO struct {
	Parent        []string `json:"parent"`
	Author        UserVO   `json:"author"`
	Committer     UserVO   `json:"committer"`
	AuthoredTime  int64    `json:"authoredTime"`
	CommittedTime int64    `json:"committedTime"`
	CommitMsg     string   `json:"commitMsg"`
	CommitId      string   `json:"commitId"`
	ShortId       string   `json:"shortId"`
	CommitSig     string   `json:"commitSig"`
	Payload       string   `json:"payload"`
	Tagger        UserVO   `json:"tagger"`
	TaggerTime    int64    `json:"taggerTime"`
	TagCommitMsg  string   `json:"tagCommitMsg"`
	TagSig        string   `json:"tagSig"`
	TagPayload    string   `json:"tagPayload"`
	ShortTagId    string   `json:"shortTagId"`
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

type DiffFileReq struct {
	RepoPath string `json:"repoPath"`
	Target   string `json:"target"`
	Head     string `json:"head"`
	FilePath string `json:"filePath"`
}

type DiffFileResp struct {
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

type GetRepoSizeReq struct {
	RepoPath string `json:"repoPath"`
}

type ShowDiffTextContentReq struct {
	RepoPath  string `json:"repoPath"`
	CommitId  string `json:"commitId"`
	FileName  string `json:"fileName"`
	StartLine int    `json:"startLine"`
	Limit     int    `json:"limit"`
}

type HistoryCommitsReq struct {
	RepoPath string `json:"repoPath"`
	Ref      string `json:"ref"`
	Offset   int    `json:"offset"`
}

type HistoryCommitsResp struct {
	Data   []CommitVO `json:"data"`
	Cursor int        `json:"cursor"`
}

type InitRepoHookReq struct {
	RepoPath string `json:"repoPath"`
}

type EntriesRepoReq struct {
	RepoPath string      `json:"repoPath"`
	Ref      string      `json:"ref"`
	RefType  git.RefType `json:"refType"`
	Dir      string      `json:"dir"`
}

type FileVO struct {
	Mode    string   `json:"mode"`
	RawPath string   `json:"rawPath"`
	Path    string   `json:"path"`
	Commit  CommitVO `json:"commit"`
}

type BlobVO struct {
	Mode    string `json:"mode"`
	RawPath string `json:"rawPath"`
	Path    string `json:"path"`
}

type TreeVO struct {
	Files []FileVO `json:"files"`
}

type BlameLineVO struct {
	Number int      `json:"number"`
	Commit CommitVO `json:"commit"`
}

type IndexRepoReq struct {
	RepoPath string      `json:"repoPath"`
	Ref      string      `json:"ref"`
	RefType  git.RefType `json:"refType"`
	Dir      string      `json:"dir"`
}

type IndexRepoResp struct {
	ReadmeText   string   `json:"readmeText"`
	HasReadme    bool     `json:"hasReadme"`
	LatestCommit CommitVO `json:"latestCommit"`
	Tree         TreeVO   `json:"tree"`
}

type CatFileReq struct {
	RepoPath string      `json:"repoPath"`
	Ref      string      `json:"ref"`
	RefType  git.RefType `json:"refType"`
	FilePath string      `json:"filePath"`
}

type CatFileResp struct {
	FileMode string   `json:"fileMode"`
	ModeName string   `json:"modeName"`
	Content  string   `json:"content"`
	Size     int64    `json:"size"`
	IsText   bool     `json:"isText"`
	Commit   CommitVO `json:"commit"`
}

type UploadPackReq struct {
	RepoPath string       `json:"repoPath"`
	C        *gin.Context `json:"-"`
}

type ReceivePackReq struct {
	RepoPath string       `json:"repoPath"`
	C        *gin.Context `json:"-"`
}

type CreateArchiveReq struct {
	RepoPath string       `json:"repoPath"`
	FileName string       `json:"fileName"`
	C        *gin.Context `json:"-"`
}

type InfoRefsReq struct {
	Service  string       `json:"service"`
	RepoPath string       `json:"repoPath"`
	C        *gin.Context `json:"-"`
}

type MergeReq struct {
	RepoPath  string `json:"repoPath"`
	Target    string `json:"target"`
	Head      string `json:"head"`
	MergeOpts struct {
		RepoId        int64  `json:"repoId"`
		PrId          int64  `json:"prId"`
		PusherAccount string `json:"pusherAccount"`
		PusherName    string `json:"pusherName"`
		PusherEmail   string `json:"pusherEmail"`
		Message       string `json:"message"`
		AppUrl        string `json:"appUrl"`
	} `json:"mergeOpts"`
}

type BlameReq struct {
	RepoPath string      `json:"repoPath"`
	Ref      string      `json:"ref"`
	FilePath string      `json:"filePath"`
	RefType  git.RefType `json:"refType"`
}

type RefVO struct {
	LastCommitId string `json:"lastCommitId"`
	Name         string `json:"name"`
}

type RefCommitVO struct {
	Name   string   `json:"name"`
	Commit CommitVO `json:"commit"`
}

type DeleteTagReqVO struct {
	RepoPath string `json:"repoPath"`
	Tag      string `json:"tag"`
}
