package feishu

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type httpResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func SendMessage(webhookUrl string, signKey string, msg Message) error {
	if err := msg.IsValid(); err != nil {
		return err
	}
	var (
		msgContent []byte
		err        error
	)
	if signKey != "" {
		ts := strconv.FormatInt(time.Now().Unix(), 10)
		stringToSign := ts + "\n" + signKey
		var data []byte
		h := hmac.New(sha256.New, []byte(stringToSign))
		_, err = h.Write(data)
		if err != nil {
			return err
		}
		msgContent, err = json.Marshal(signedMessage{
			Message:   msg,
			Timestamp: ts,
			Sign:      base64.StdEncoding.EncodeToString(h.Sum(nil)),
		})
	} else {
		msgContent, err = json.Marshal(msg)
	}
	if err != nil {
		return err
	}
	payload := strings.NewReader(string(msgContent))
	httpReq, err := http.NewRequest(http.MethodPost, webhookUrl, payload)
	if err != nil {
		return err
	}
	httpReq.Header.Add("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("webhook server return status code: %v", response.StatusCode)
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	var resp httpResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return err
	}
	if resp.Code != 0 {
		return fmt.Errorf("webhook server return err: %v", resp)
	}
	return nil
}
