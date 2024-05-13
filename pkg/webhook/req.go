package webhook

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/LeeZXin/zsf-utils/httputil"
	"net/http"
)

const (
	XEventType = "x-zgit-event-type"
	XSignature = "x-zgit-signature"
)

type PullRequestAction string

const (
	SubmitAction PullRequestAction = "submit"
	CloseAction  PullRequestAction = "close"
	MergeAction  PullRequestAction = "merge"
	ReviewAction PullRequestAction = "review"
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
