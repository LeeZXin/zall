package webhook

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/LeeZXin/zall/pkg/branch"
	"github.com/LeeZXin/zsf-utils/httputil"
	"net/http"
)

const (
	XEventType = "x-zgit-event-type"
	XSignature = "x-zgit-signature"
)

type PullRequestAction string

const (
	PrSubmitAction PullRequestAction = "submit"
	PrCloseAction  PullRequestAction = "close"
	PrMergeAction  PullRequestAction = "merge"
	PrReviewAction PullRequestAction = "review"
)

type GitRepoAction string

const (
	RepoDeleteTemporarilyAction GitRepoAction = "deleteTemporarily"
	RepoDeletePermanentlyAction GitRepoAction = "deletePermanently"
	RepoArchivedAction          GitRepoAction = "archived"
	RepoUnArchivedAction        GitRepoAction = "unArchived"
	RepoRecoverFromRecycle      GitRepoAction = "recoverFromRecycle"
)

type ProtectedBranchAction string

const (
	PbCreateAction ProtectedBranchAction = "create"
	PbUpdateAction ProtectedBranchAction = "update"
	PbDeleteAction ProtectedBranchAction = "delete"
)

type EventReq interface {
	EventType() Event
}

type BaseRepoReq struct {
	RepoId    int64  `json:"repoId"`
	RepoName  string `json:"repoName"`
	Account   string `json:"account"`
	EventTime int64  `json:"eventTime"`
}

type ProtectedBranchObj struct {
	Pattern string `json:"pattern"`
	branch.ProtectedBranchCfg
}

type ProtectedBranchEventReq struct {
	BaseRepoReq
	Action ProtectedBranchAction `json:"action"`
	Before *ProtectedBranchObj   `json:"before,omitempty"`
	After  *ProtectedBranchObj   `json:"after,omitempty"`
}

func (*ProtectedBranchEventReq) EventType() Event {
	return ProtectedBranchEvent
}

type GitRepoEventReq struct {
	BaseRepoReq
	Action GitRepoAction `json:"action"`
}

func (*GitRepoEventReq) EventType() Event {
	return GitRepoEvent
}

type GitPushEventReq struct {
	RefType     string `json:"refType"`
	Ref         string `json:"ref"`
	OldCommitId string `json:"oldCommitId"`
	NewCommitId string `json:"newCommitId"`
	BaseRepoReq
}

func (*GitPushEventReq) EventType() Event {
	return GitPushEvent
}

type PullRequestEventReq struct {
	PrId    int64             `json:"prId"`
	PrTitle string            `json:"prTitle"`
	Action  PullRequestAction `json:"action"`
	BaseRepoReq
}

func (*PullRequestEventReq) EventType() Event {
	return PullRequestEvent
}

type PingEventReq struct {
	EventTime int64 `json:"eventTime"`
}

func (*PingEventReq) EventType() Event {
	return PingEvent
}

func Post(ctx context.Context, url, secret string, req EventReq) error {
	return post(ctx, url, req.EventType().String(), secret, req)
}

func post(ctx context.Context, url, eventType, secret string, req any) error {
	initTrigger()
	reqJson, err := json.Marshal(req)
	if err != nil {
		return err
	}
	request, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(reqJson))
	if err != nil {
		return err
	}
	request.Header.Set(XEventType, eventType)
	request.Header.Set(XSignature, CreateSignature(reqJson, secret))
	request.Header.Set("Content-Type", httputil.JsonContentType)
	post, err := httpClient.Do(request)
	if err != nil {
		return err
	}
	defer post.Body.Close()
	if post.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("http request return code: %v", post.StatusCode)
	}
	return nil
}

func CreateSignature(req []byte, secret string) string {
	hash := hmac.New(sha256.New, []byte(secret))
	_, err := hash.Write(req)
	if err != nil {
		return "failToCreateHmacSignature"
	}
	return hex.EncodeToString(hash.Sum(nil))
}
