package reqvo

import (
	"github.com/gin-gonic/gin"
)

type InitRepoReq struct {
	UserAccount   string `json:"userAccount"`
	UserEmail     string `json:"userEmail"`
	RepoName      string `json:"repoName"`
	RepoPath      string `json:"repoPath"`
	CreateReadme  bool   `json:"createReadme"`
	GitIgnoreName string `json:"gitIgnoreName"`
	DefaultBranch string `json:"defaultBranch"`
}

type DeleteRepoReq struct {
	RepoPath string `json:"repoPath"`
}

type GetAllBranchesReq struct {
	RepoPath string `json:"repoPath"`
}

type GetAllTagsReq struct {
	RepoPath string `json:"repoPath"`
}

type GcReq struct {
	RepoPath string `json:"repoPath"`
}

type DiffRefsReq struct {
	RepoPath string `json:"repoPath"`
	Target   string `json:"target"`
	Head     string `json:"head"`
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

type UserVO struct {
	Account string `json:"account"`
	Email   string `json:"email"`
}

type CommitVO struct {
	Author        UserVO `json:"author"`
	Committer     UserVO `json:"committer"`
	AuthoredTime  int64  `json:"authoredTime"`
	CommittedTime int64  `json:"committedTime"`
	CommitMsg     string `json:"commitMsg"`
	CommitId      string `json:"commitId"`
	ShortId       string `json:"shortId"`
	CommitSig     string `json:"commitSig"`
	Payload       string `json:"payload"`
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

type DiffFileReq struct {
	RepoPath string `json:"repoPath"`
	Target   string `json:"target"`
	Head     string `json:"head"`
	FileName string `json:"fileName"`
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
	Index   int    `json:"index"`
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
	RepoPath string `json:"repoPath"`
	Ref      string `json:"ref"`
	Dir      string `json:"dir"`
	Offset   int    `json:"offset"`
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

type TreeRepoReq struct {
	RepoPath string `json:"repoPath"`
	Ref      string `json:"ref"`
	Dir      string `json:"dir"`
}

type TreeRepoResp struct {
	IsEmpty      bool     `json:"isEmpty"`
	ReadmeText   string   `json:"readmeText"`
	HasReadme    bool     `json:"hasReadme"`
	LatestCommit CommitVO `json:"latestCommit"`
	Tree         TreeVO   `json:"tree"`
}

type CatFileReq struct {
	RepoPath string `json:"repoPath"`
	Ref      string `json:"ref"`
	Dir      string `json:"dir"`
	FileName string `json:"fileName"`
}

type CatFileResp struct {
	FileMode string `json:"fileMode"`
	ModeName string `json:"modeName"`
	Content  string `json:"content"`
}

type UploadPackReq struct {
	RepoPath string       `json:"repoPath"`
	C        *gin.Context `json:"-"`
}

type ReceivePackReq struct {
	RepoPath string       `json:"repoPath"`
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
		PusherEmail   string `json:"pusherEmail"`
		Message       string `json:"message"`
		AppUrl        string `json:"appUrl"`
	} `json:"mergeOpts"`
}
