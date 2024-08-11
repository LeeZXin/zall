package reposrv

import (
	"github.com/LeeZXin/zall/git/modules/model/pullrequestmd"
	"github.com/LeeZXin/zall/git/modules/model/repomd"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/pkg/git"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/collections/hashset"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

type CreateRepoReqDTO struct {
	Operator      apisession.UserInfo `json:"operator"`
	TeamId        int64               `json:"teamId"`
	Name          string              `json:"name"`
	Desc          string              `json:"desc"`
	AddReadme     bool                `json:"addReadme"`
	GitIgnoreName string              `json:"gitIgnoreName"`
	DefaultBranch string              `json:"defaultBranch"`
}

func (r *CreateRepoReqDTO) IsValid() error {
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if !repomd.IsRepoNameValid(r.Name) {
		return util.InvalidArgsError()
	}
	if !repomd.IsBranchValid(r.DefaultBranch) {
		return util.InvalidArgsError()
	}
	if r.GitIgnoreName != "" && !gitignoreSet.Contains(r.GitIgnoreName) {
		return util.InvalidArgsError()
	}
	return nil
}

type GetRepoAndPermReqDTO struct {
	RepoId   int64               `json:"repoId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *GetRepoAndPermReqDTO) IsValid() error {
	if r.RepoId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type DeleteRepoReqDTO struct {
	RepoId   int64               `json:"repoId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *DeleteRepoReqDTO) IsValid() error {
	if r.RepoId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type RecoverFromRecycleReqDTO struct {
	RepoId   int64               `json:"repoId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *RecoverFromRecycleReqDTO) IsValid() error {
	if r.RepoId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type IndexRepoReqDTO struct {
	RepoId   int64               `json:"repoId"`
	Ref      string              `json:"ref"`
	RefType  git.RefType         `json:"refType"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *IndexRepoReqDTO) IsValid() error {
	if !r.RefType.IsValid() {
		return util.InvalidArgsError()
	}
	if r.RepoId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if len(r.Ref) > 128 || len(r.Ref) == 0 {
		return util.InvalidArgsError()
	}
	return nil
}

type LsTreeRepoReqDTO struct {
	RepoId   int64               `json:"repoId"`
	Ref      string              `json:"ref"`
	Dir      string              `json:"dir"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *LsTreeRepoReqDTO) IsValid() error {
	if r.RepoId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if len(r.Ref) > 128 || len(r.Ref) == 0 {
		return util.InvalidArgsError()
	}
	if strings.HasSuffix(r.Dir, "/") {
		return util.InvalidArgsError()
	}
	return nil
}

type CatFileReqDTO struct {
	RepoId   int64               `json:"repoId"`
	Ref      string              `json:"ref"`
	RefType  git.RefType         `json:"refType"`
	FilePath string              `json:"filePath"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *CatFileReqDTO) IsValid() error {
	if !r.RefType.IsValid() {
		return util.InvalidArgsError()
	}
	if r.RepoId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if len(r.Ref) > 128 || len(r.Ref) == 0 {
		return util.InvalidArgsError()
	}
	return nil
}

type BlameReqDTO struct {
	RepoId   int64               `json:"repoId"`
	Ref      string              `json:"ref"`
	RefType  git.RefType         `json:"refType"`
	FilePath string              `json:"filePath"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *BlameReqDTO) IsValid() error {
	if !r.RefType.IsValid() {
		return util.InvalidArgsError()
	}
	if r.RepoId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if len(r.Ref) > 128 || len(r.Ref) == 0 {
		return util.InvalidArgsError()
	}
	return nil
}

type CatFileRespDTO struct {
	FileMode string    `json:"fileMode"`
	ModeName string    `json:"modeName"`
	Content  string    `json:"content"`
	Size     int64     `json:"size"`
	Commit   CommitDTO `json:"commit"`
}

type EntriesRepoReqDTO struct {
	RepoId   int64               `json:"repoId"`
	Ref      string              `json:"ref"`
	RefType  git.RefType         `json:"refType"`
	Dir      string              `json:"dir"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *EntriesRepoReqDTO) IsValid() error {
	if !r.RefType.IsValid() {
		return util.InvalidArgsError()
	}
	if r.RepoId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if r.Ref == "" {
		return util.InvalidArgsError()
	}
	return nil
}

type ListRepoReqDTO struct {
	TeamId   int64               `json:"teamId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListRepoReqDTO) IsValid() error {
	if r.TeamId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type UserDTO struct {
	Account string `json:"account"`
	Email   string `json:"email"`
}

type CommitDTO struct {
	Parent        []string  `json:"parent"`
	Author        UserDTO   `json:"author"`
	Committer     UserDTO   `json:"committer"`
	AuthoredTime  int64     `json:"authoredTime"`
	CommittedTime int64     `json:"committedTime"`
	CommitMsg     string    `json:"commitMsg"`
	CommitId      string    `json:"commitId"`
	ShortId       string    `json:"shortId"`
	Tagger        UserDTO   `json:"tagger"`
	TaggerTime    int64     `json:"taggerTime"`
	ShortTagId    string    `json:"shortTagId"`
	TagCommitMsg  string    `json:"tagCommitMsg"`
	Verified      bool      `json:"verified"`
	Signer        SignerDTO `json:"signer"`
}

type SignerDTO struct {
	Account   string
	Name      string
	AvatarUrl string
	Key       string
	Type      string
}

type FileDTO struct {
	Mode    string    `json:"mode"`
	RawPath string    `json:"rawPath"`
	Path    string    `json:"path"`
	Commit  CommitDTO `json:"commit"`
}

type BlobDTO struct {
	Mode    string `json:"mode"`
	RawPath string `json:"rawPath"`
	Path    string `json:"path"`
}

type TreeDTO struct {
	Files []FileDTO `json:"files"`
}

type IndexRepoRespDTO struct {
	ReadmeText   string    `json:"readmeText"`
	HasReadme    bool      `json:"hasReadme"`
	LatestCommit CommitDTO `json:"latestCommit"`
	Tree         TreeDTO   `json:"tree"`
}

type RepoTypeDTO struct {
	Option int    `json:"option"`
	Name   string `json:"name"`
}

type AllBranchesReqDTO struct {
	RepoId   int64               `json:"repoId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *AllBranchesReqDTO) IsValid() error {
	if r.RepoId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type AllTagsReqDTO struct {
	RepoId   int64               `json:"repoId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *AllTagsReqDTO) IsValid() error {
	if r.RepoId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type GcReqDTO struct {
	RepoId   int64               `json:"repoId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *GcReqDTO) IsValid() error {
	if r.RepoId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type GetRepoSizeReqDTO struct {
	RepoId   int64               `json:"repoId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *GetRepoSizeReqDTO) IsValid() error {
	if r.RepoId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

var gitignoreSet = hashset.NewHashSet(
	"AL", "Actionscript", "Ada", "Agda", "AltiumDesigner", "Android", "Anjuta", "Ansible", "AppEngine",
	"AppceleratorTitanium", "ArchLinuxPackages", "Archives", "AtmelStudio", "AutoIt", "Autotools", "B4X", "Backup",
	"Bazaar", "Bazel", "Beef", "Bitrix", "BricxCC", "C", "C++", "CDK", "CFWheels", "CMake", "CUDA", "CakePHP",
	"Calabash", "ChefCookbook", "Clojure", "Cloud9", "CodeIgniter", "CodeKit", "CodeSniffer", "CommonLisp",
	"Composer", "Concrete5", "Coq", "Cordova", "CraftCMS", "D", "DM", "Dart", "DartEditor", "Delphi", "Diff",
	"Dreamweaver", "Dropbox", "Drupal", "Drupal7", "EPiServer", "Eagle", "Eclipse", "EiffelStudio", "Elisp",
	"Elixir", "Elm", "Emacs", "Ensime", "Erlang", "Espresso", "Exercism", "ExpressionEngine", "ExtJs", "Fancy",
	"Finale", "FlaxEngine", "FlexBuilder", "ForceDotCom", "Fortran", "FuelPHP", "GNOMEShellExtension", "GPG",
	"GWT", "Gcov", "GitBook", "Go", "Go.AllowList", "Godot", "Gradle", "Grails", "Gretl", "Haskell", "Hugo",
	"IAR_EWARM", "IGORPro", "Idris", "Images", "InforCMS", "JBoss", "JBoss4", "JBoss6", "JDeveloper", "JENKINS_HOME",
	"JEnv", "Java", "Jekyll", "JetBrains", "Jigsaw", "Joomla", "Julia", "JupyterNotebooks", "KDevelop4", "Kate",
	"Kentico", "KiCad", "Kohana", "Kotlin", "LabVIEW", "Laravel", "Lazarus", "Leiningen", "LemonStand", "LensStudio",
	"LibreOffice", "Lilypond", "Linux", "Lithium", "Logtalk", "Lua", "LyX", "MATLAB", "Magento", "Magento1", "Magento2",
	"Maven", "Mercurial", "Mercury", "MetaProgrammingSystem", "Metals", "Meteor", "MicrosoftOffice", "ModelSim",
	"Momentics", "MonoDevelop", "NWjs", "Nanoc", "NasaSpecsIntact", "NetBeans", "Nikola", "Nim", "Ninja", "Nix",
	"Node", "NotepadPP", "OCaml", "Objective-C", "Octave", "Opa", "OpenCart", "OpenSSL", "OracleForms", "Otto",
	"PSoCCreator", "Packer", "Patch", "Perl", "Perl6", "Phalcon", "Phoenix", "Pimcore", "PlayFramework", "Plone",
	"Prestashop", "Processing", "PuTTY", "Puppet", "PureScript", "Python", "Qooxdoo", "Qt", "R", "ROS", "ROS2",
	"Racket", "Rails", "Raku", "Red", "Redcar", "Redis", "RhodesRhomobile", "Ruby", "Rust", "SAM", "SBT", "SCons",
	"SPFx", "SVN", "Sass", "Scala", "Scheme", "Scrivener", "Sdcc", "SeamGen", "SketchUp", "SlickEdit", "Smalltalk",
	"Snap", "Splunk", "Stata", "Stella", "Strapi", "SublimeText", "SugarCRM", "Swift", "Symfony", "SymphonyCMS",
	"Syncthing", "SynopsysVCS", "Tags", "TeX", "Terraform", "TextMate", "Textpattern", "ThinkPHP", "Toit", "TortoiseGit",
	"TurboGears2", "TwinCAT3", "Typo3", "Umbraco", "Unity", "UnrealEngine", "V", "VVVV", "Vagrant", "Vim", "VirtualEnv",
	"Virtuoso", "VisualStudio", "VisualStudioCode", "Vue", "Waf", "WebMethods", "Windows", "WordPress", "Xcode", "Xilinx",
	"XilinxISE", "Xojo", "Yeoman", "Yii", "ZendFramework", "Zephir", "core", "esp-idf", "macOS", "uVision",
)

type DiffRefsReqDTO struct {
	RepoId     int64               `json:"repoId"`
	Target     string              `json:"target"`
	TargetType git.RefType         `json:"targetType"`
	Head       string              `json:"head"`
	HeadType   git.RefType         `json:"headType"`
	Operator   apisession.UserInfo `json:"operator"`
}

func (r *DiffRefsReqDTO) IsValid() error {
	if r.RepoId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if !r.TargetType.IsValid() {
		return util.InvalidArgsError()
	}
	if !r.HeadType.IsValid() {
		return util.InvalidArgsError()
	}
	if !util.ValidateRef(r.Target) {
		return util.InvalidArgsError()
	}
	if !util.ValidateRef(r.Head) {
		return util.InvalidArgsError()
	}
	return nil
}

type DiffCommitsReqDTO struct {
	RepoId   int64               `json:"repoId"`
	CommitId string              `json:"commitId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *DiffCommitsReqDTO) IsValid() error {
	if r.RepoId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type DiffRefsRespDTO struct {
	Target        string              `json:"target"`
	Head          string              `json:"head"`
	TargetCommit  CommitDTO           `json:"targetCommit"`
	HeadCommit    CommitDTO           `json:"headCommit"`
	Commits       []CommitDTO         `json:"commits"`
	NumFiles      int                 `json:"numFiles"`
	DiffNumsStats DiffNumsStatInfoDTO `json:"diffNumsStats"`
	ConflictFiles []string            `json:"conflictFiles"`
	CanMerge      bool                `json:"canMerge"`
}

type DiffCommitsRespDTO struct {
	Commit        CommitDTO           `json:"commit"`
	NumFiles      int                 `json:"numFiles"`
	DiffNumsStats DiffNumsStatInfoDTO `json:"diffNumsStats"`
}

type DiffNumsStatInfoDTO struct {
	FileChangeNums int `json:"fileChangeNums"`
	InsertNums     int `json:"insertNums"`
	DeleteNums     int `json:"deleteNums"`
	Stats          []DiffNumsStatDTO
}

type DiffNumsStatDTO struct {
	RawPath    string `json:"rawPath"`
	Path       string `json:"path"`
	InsertNums int    `json:"insertNums"`
	DeleteNums int    `json:"deleteNums"`
}

type DiffFileRespDTO struct {
	FilePath    string        `json:"filePath"`
	OldMode     string        `json:"oldMode"`
	Mode        string        `json:"mode"`
	IsSubModule bool          `json:"isSubModule"`
	FileType    string        `json:"fileType"`
	IsBinary    bool          `json:"isBinary"`
	RenameFrom  string        `json:"renameFrom"`
	RenameTo    string        `json:"renameTo"`
	CopyFrom    string        `json:"copyFrom"`
	CopyTo      string        `json:"copyTo"`
	Lines       []DiffLineDTO `json:"lines"`
}

type DiffLineDTO struct {
	LeftNo  int    `json:"leftNo"`
	Prefix  string `json:"prefix"`
	RightNo int    `json:"rightNo"`
	Text    string `json:"text"`
}

type DiffFileReqDTO struct {
	RepoId   int64               `json:"repoId"`
	Target   string              `json:"target"`
	Head     string              `json:"head"`
	FilePath string              `json:"filePath"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *DiffFileReqDTO) IsValid() error {
	if r.RepoId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if !util.ValidateRef(r.Target) {
		return util.InvalidArgsError()
	}
	if r.Head != "" && !util.ValidateRef(r.Head) {
		return util.InvalidArgsError()
	}
	return nil
}

type HistoryCommitsReqDTO struct {
	RepoId   int64               `json:"repoId"`
	Ref      string              `json:"ref"`
	Cursor   int                 `json:"cursor"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *HistoryCommitsReqDTO) IsValid() error {
	if r.RepoId <= 0 {
		return util.InvalidArgsError()
	}
	if len(r.Ref) == 0 {
		return util.InvalidArgsError()
	}
	if r.Cursor < 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type HistoryCommitsRespDTO struct {
	Data   []CommitDTO
	Cursor int
}

type TransferTeamReqDTO struct {
	RepoId   int64               `json:"repoId"`
	TeamId   int64               `json:"teamId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *TransferTeamReqDTO) IsValid() error {
	if r.RepoId <= 0 || r.TeamId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type GetSimpleInfoReqDTO struct {
	RepoId   int64               `json:"repoId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *GetSimpleInfoReqDTO) IsValid() error {
	if r.RepoId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type GetDetailInfoReqDTO struct {
	RepoId   int64               `json:"repoId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *GetDetailInfoReqDTO) IsValid() error {
	if r.RepoId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type SimpleInfoDTO struct {
	Branches     []string `json:"branches"`
	Tags         []string `json:"tags"`
	CloneHttpUrl string   `json:"cloneHttpUrl"`
	CloneSshUrl  string   `json:"cloneSshUrl"`
}

type BlameLineDTO struct {
	Number int       `json:"number"`
	Commit CommitDTO `json:"commit"`
}

type PageRefCommitsReqDTO struct {
	RepoId   int64               `json:"repoId"`
	PageNum  int                 `json:"pageNum"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *PageRefCommitsReqDTO) IsValid() error {
	if r.RepoId <= 0 {
		return util.InvalidArgsError()
	}
	if r.PageNum <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type BranchCommitDTO struct {
	Name              string
	IsProtectedBranch bool
	LastCommit        CommitDTO
	LastPullRequest   *PullRequestDTO
}

type TagCommitDTO struct {
	Name   string
	Commit CommitDTO
}

type PullRequestDTO struct {
	Id       int64                  `json:"id"`
	PrStatus pullrequestmd.PrStatus `json:"prStatus"`
	PrTitle  string                 `json:"prTitle"`
	Created  time.Time              `json:"created"`
}

type DeleteBranchReqDTO struct {
	RepoId   int64               `json:"repoId"`
	Branch   string              `json:"branch"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *DeleteBranchReqDTO) IsValid() error {
	if r.RepoId <= 0 {
		return util.InvalidArgsError()
	}
	if r.Branch == "" {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type DeleteTagReqDTO struct {
	RepoId   int64               `json:"repoId"`
	Tag      string              `json:"tag"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *DeleteTagReqDTO) IsValid() error {
	if r.RepoId <= 0 {
		return util.InvalidArgsError()
	}
	if r.Tag == "" {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type CreateArchiveReqDTO struct {
	RepoId   int64               `json:"repoId"`
	FileName string              `json:"fileName"`
	C        *gin.Context        `json:"-"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *CreateArchiveReqDTO) IsValid() error {
	if r.RepoId <= 0 {
		return util.InvalidArgsError()
	}
	if r.FileName == "" {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if r.C == nil {
		return util.InvalidArgsError()
	}
	return nil
}

type RepoDTO struct {
	Id            int64
	Path          string
	Name          string
	TeamId        int64
	RepoDesc      string
	DefaultBranch string
	GitSize       int64
	LfsSize       int64
	DisableLfs    bool
	LfsLimitSize  int64
	GitLimitSize  int64
	LastOperated  time.Time
	IsArchived    bool
	Created       time.Time
}

type UpdateRepoReqDTO struct {
	RepoId       int64               `json:"repoId"`
	Desc         string              `json:"desc"`
	DisableLfs   bool                `json:"disableLfs"`
	LfsLimitSize int64               `json:"lfsLimitSize"`
	GitLimitSize int64               `json:"gitLimitSize"`
	Operator     apisession.UserInfo `json:"operator"`
}

func (r *UpdateRepoReqDTO) IsValid() error {
	if r.RepoId <= 0 {
		return util.InvalidArgsError()
	}
	if !repomd.IsRepoDescValid(r.Desc) {
		return util.InvalidArgsError()
	}
	if r.GitLimitSize < 0 || r.LfsLimitSize < 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type SetRepoArchivedStatusReqDTO struct {
	RepoId     int64               `json:"repoId"`
	IsArchived bool                `json:"isArchived"`
	Operator   apisession.UserInfo `json:"operator"`
}

func (r *SetRepoArchivedStatusReqDTO) IsValid() error {
	if r.RepoId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ListDeletedRepoReqDTO struct {
	TeamId   int64               `json:"teamId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListDeletedRepoReqDTO) IsValid() error {
	if r.TeamId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type DeletedRepoDTO struct {
	RepoDTO
	Deleted time.Time
}

type ListRepoByAdminReqDTO struct {
	TeamId   int64               `json:"teamId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListRepoByAdminReqDTO) IsValid() error {
	if r.TeamId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type SimpleRepoDTO struct {
	RepoId int64
	Name   string
	TeamId int64
}
