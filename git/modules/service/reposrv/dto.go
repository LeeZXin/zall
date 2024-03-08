package reposrv

import (
	"github.com/LeeZXin/zall/git/modules/model/actiontaskmd"
	"github.com/LeeZXin/zall/git/modules/model/repomd"
	"github.com/LeeZXin/zall/pkg/apisession"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf-utils/collections/hashset"
	"regexp"
	"strings"
	"time"
)

const (
	UpDirection   = "up"
	DownDirection = "down"
)

var (
	validRepoNamePattern = regexp.MustCompile("^[\\w\\-]{1,32}$")
	validBranchPattern   = regexp.MustCompile("^\\w{1,32}$")
)

type InitRepoReqDTO struct {
	Operator      apisession.UserInfo `json:"operator"`
	TeamId        int64               `json:"teamId"`
	Name          string              `json:"name"`
	Desc          string              `json:"desc"`
	CreateReadme  bool                `json:"createReadme"`
	GitIgnoreName string              `json:"gitIgnoreName"`
	DefaultBranch string              `json:"defaultBranch"`
}

func (r *InitRepoReqDTO) IsValid() error {
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if !validRepoNamePattern.MatchString(r.Name) {
		return util.InvalidArgsError()
	}
	if len(r.Desc) > 255 {
		return util.InvalidArgsError()
	}
	if r.DefaultBranch != "" && !validBranchPattern.MatchString(r.DefaultBranch) {
		return util.InvalidArgsError()
	}
	if r.GitIgnoreName != "" && !gitignoreSet.Contains(r.GitIgnoreName) {
		return util.InvalidArgsError()
	}
	return nil
}

type DeleteRepoReqDTO struct {
	Id       int64               `json:"id"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *DeleteRepoReqDTO) IsValid() error {
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type TreeRepoReqDTO struct {
	Id       int64               `json:"id"`
	Ref      string              `json:"ref"`
	Dir      string              `json:"dir"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *TreeRepoReqDTO) IsValid() error {
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
	Id       int64               `json:"id"`
	Ref      string              `json:"ref"`
	Dir      string              `json:"dir"`
	FileName string              `json:"fileName"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *CatFileReqDTO) IsValid() error {
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if len(r.Ref) > 128 || len(r.Ref) == 0 {
		return util.InvalidArgsError()
	}
	if len(r.Dir) > 128 || len(r.Dir) == 0 {
		return util.InvalidArgsError()
	}
	if !validateFileName(r.FileName) {
		return util.InvalidArgsError()
	}
	if strings.HasSuffix(r.Dir, "/") {
		return util.InvalidArgsError()
	}
	return nil
}

type CatFileRespDTO struct {
	FileMode string `json:"fileMode"`
	ModeName string `json:"modeName"`
	Content  string `json:"content"`
}

type EntriesRepoReqDTO struct {
	Id       int64               `json:"id"`
	Ref      string              `json:"ref"`
	Dir      string              `json:"dir"`
	Offset   int                 `json:"offset"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *EntriesRepoReqDTO) IsValid() error {
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if r.Ref == "" {
		return util.InvalidArgsError()
	}
	if strings.HasSuffix(r.Dir, "/") {
		return util.InvalidArgsError()
	}
	if r.Offset < 0 {
		return util.InvalidArgsError()
	}
	return nil
}

type ListRepoReqDTO struct {
	TeamId   int64               `json:"teamId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListRepoReqDTO) IsValid() error {
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
	Author        UserDTO `json:"author"`
	Committer     UserDTO `json:"committer"`
	AuthoredTime  int64   `json:"authoredTime"`
	CommittedTime int64   `json:"committedTime"`
	CommitMsg     string  `json:"commitMsg"`
	CommitId      string  `json:"commitId"`
	ShortId       string  `json:"shortId"`
	Verified      bool    `json:"verified"`
}

type FileDTO struct {
	Mode    string    `json:"mode"`
	RawPath string    `json:"rawPath"`
	Path    string    `json:"path"`
	Commit  CommitDTO `json:"commit"`
}

type TreeDTO struct {
	Files   []FileDTO `json:"files"`
	Limit   int       `json:"limit"`
	Offset  int       `json:"offset"`
	HasMore bool      `json:"hasMore"`
}

type TreeRepoRespDTO struct {
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
	Id       int64               `json:"id"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *AllBranchesReqDTO) IsValid() error {
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type AllTagsReqDTO struct {
	Id       int64               `json:"id"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *AllTagsReqDTO) IsValid() error {
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type GcReqDTO struct {
	Id       int64               `json:"id"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *GcReqDTO) IsValid() error {
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

type DiffCommitsReqDTO struct {
	Id       int64               `json:"id"`
	Target   string              `json:"target"`
	Head     string              `json:"head"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *DiffCommitsReqDTO) IsValid() error {
	if !r.Operator.IsValid() {
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

type DiffCommitsRespDTO struct {
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

type DiffNumsStatInfoDTO struct {
	FileChangeNums int `json:"fileChangeNums"`
	InsertNums     int `json:"insertNums"`
	DeleteNums     int `json:"deleteNums"`
	Stats          []DiffNumsStatDTO
}

type DiffNumsStatDTO struct {
	RawPath    string `json:"rawPath"`
	Path       string `json:"path"`
	TotalNums  int    `json:"totalNums"`
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
	Index   int    `json:"index"`
	LeftNo  int    `json:"leftNo"`
	Prefix  string `json:"prefix"`
	RightNo int    `json:"rightNo"`
	Text    string `json:"text"`
}

type ShowDiffTextContentReqDTO struct {
	Id        int64               `json:"id"`
	CommitId  string              `json:"commitId"`
	FileName  string              `json:"fileName"`
	Offset    int                 `json:"offset"`
	Limit     int                 `json:"limit"`
	Direction string              `json:"direction"`
	Operator  apisession.UserInfo `json:"operator"`
}

func (r *ShowDiffTextContentReqDTO) IsValid() error {
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if !util.ValidateCommitId(r.CommitId) {
		return util.InvalidArgsError()
	}
	if r.Offset < 0 {
		return util.InvalidArgsError()
	}
	if len(r.Direction) == 0 || len(r.Direction) > 10 {
		return util.InvalidArgsError()
	}
	if r.Direction != UpDirection && r.Direction != DownDirection {
		return util.InvalidArgsError()
	}
	if !validateFileName(r.FileName) {
		return util.InvalidArgsError()
	}
	return nil
}

type DiffFileReqDTO struct {
	Id       int64               `json:"id"`
	Target   string              `json:"target"`
	Head     string              `json:"head"`
	FileName string              `json:"fileName"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *DiffFileReqDTO) IsValid() error {
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	if !util.ValidateRef(r.Target) {
		return util.InvalidArgsError()
	}
	if !util.ValidateRef(r.Head) {
		return util.InvalidArgsError()
	}
	if !validateFileName(r.FileName) {
		return util.InvalidArgsError()
	}
	return nil
}

type HistoryCommitsReqDTO struct {
	Id       int64               `json:"id"`
	Ref      string              `json:"ref"`
	Cursor   int                 `json:"cursor"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *HistoryCommitsReqDTO) IsValid() error {
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

func validateFileName(name string) bool {
	return len(name) <= 255 && len(name) > 0
}

type InsertAccessTokenReqDTO struct {
	Id       int64               `json:"id"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *InsertAccessTokenReqDTO) IsValid() error {
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type DeleteAccessTokenReqDTO struct {
	Id       int64               `json:"id"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *DeleteAccessTokenReqDTO) IsValid() error {
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ListAccessTokenReqDTO struct {
	Id       int64               `json:"id"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListAccessTokenReqDTO) IsValid() error {
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type AccessTokenDTO struct {
	Id      int64     `json:"id"`
	Account string    `json:"account"`
	Token   string    `json:"token"`
	Created time.Time `json:"created"`
}

type CheckAccessTokenReqDTO struct {
	Id      int64  `json:"id"`
	Account string `json:"account"`
	Token   string `json:"token"`
}

func (r *CheckAccessTokenReqDTO) IsValid() error {
	if !repomd.IsAccessTokenAccountValid(r.Account) {
		return util.InvalidArgsError()
	}
	if !repomd.IsAccessTokenTokenValid(r.Token) {
		return util.InvalidArgsError()
	}
	return nil
}

type InsertActionReqDTO struct {
	Id             int64               `json:"id"`
	ActionContent  string              `json:"actionContent"`
	AssignInstance string              `json:"assignInstance"`
	Operator       apisession.UserInfo `json:"operator"`
}

func (r *InsertActionReqDTO) IsValid() error {
	if r.ActionContent == "" {
		return util.InvalidArgsError()
	}
	if r.AssignInstance != "" && !actiontaskmd.IsInstanceIdValid(r.AssignInstance) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type DeleteActionReqDTO struct {
	Id       int64               `json:"id"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *DeleteActionReqDTO) IsValid() error {
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type ListActionReqDTO struct {
	Id       int64               `json:"id"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *ListActionReqDTO) IsValid() error {
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type UpdateActionReqDTO struct {
	Id             int64               `json:"id"`
	ActionContent  string              `json:"actionContent"`
	AssignInstance string              `json:"assignInstance"`
	Operator       apisession.UserInfo `json:"operator"`
}

func (r *UpdateActionReqDTO) IsValid() error {
	if r.ActionContent == "" {
		return util.InvalidArgsError()
	}
	if r.AssignInstance != "" && !actiontaskmd.IsInstanceIdValid(r.AssignInstance) {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type RefreshAllGitHooksReqDTO struct {
	Operator apisession.UserInfo `json:"operator"`
}

func (r *RefreshAllGitHooksReqDTO) IsValid() error {
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type TriggerActionReqDTO struct {
	Id       int64               `json:"id"`
	Ref      string              `json:"ref"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *TriggerActionReqDTO) IsValid() error {
	if len(r.Ref) > 255 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}

type TransferTeamReqDTO struct {
	Id       int64               `json:"id"`
	TeamId   int64               `json:"teamId"`
	Operator apisession.UserInfo `json:"operator"`
}

func (r *TransferTeamReqDTO) IsValid() error {
	if r.Id <= 0 || r.TeamId <= 0 {
		return util.InvalidArgsError()
	}
	if !r.Operator.IsValid() {
		return util.InvalidArgsError()
	}
	return nil
}
