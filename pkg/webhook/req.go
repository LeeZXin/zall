package webhook

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/LeeZXin/zall/pkg/event"
	"github.com/LeeZXin/zsf-utils/httputil"
	"net/http"
)

const (
	XEventType = "x-zgit-event-type"
	XSignature = "x-zgit-signature"
)

func Post(ctx context.Context, url, secret string, req event.Event) error {
	return post(ctx, url, req.EventType(), secret, req)
}

func post(ctx context.Context, url, eventType, secret string, req any) error {
	initTrigger()
	reqJson, err := json.Marshal(req)
	if err != nil {
		return err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(reqJson))
	if err != nil {
		return err
	}
	request.Header.Set(XEventType, eventType)
	request.Header.Set(XSignature, CreateSignature(reqJson, secret))
	request.Header.Set("Content-Type", httputil.JsonContentType)
	resp, err := httpClient.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("http request return code: %v", resp.StatusCode)
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
