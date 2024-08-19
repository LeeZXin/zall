package wework

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type httpResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func SendMessage(webhookUrl string, msg Message) error {
	if err := msg.IsValid(); err != nil {
		return err
	}
	msgContent, err := json.Marshal(msg)
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
	if resp.ErrCode != 0 {
		return fmt.Errorf("webhook server return err: %v", resp)
	}
	return nil
}
