package feishu

import (
	"fmt"
	"github.com/pingcap/errors"
)

const (
	TextMsgType        = "text"
	PostMsgType        = "post"
	InteractiveMsgType = "interactive"
)

type TextMessage string

func FormatAtUser(userId string, name string) string {
	return fmt.Sprintf(`<at user_id="%s">%s</at>`, userId, name)
}

const (
	TextPostParagraphType = "text"
	HrefPostParagraphType = "a"
	AtPostParagraphType   = "at"
)

type PostParagraph struct {
	Tag      string `json:"tag"`
	Text     string `json:"text,omitempty"`
	Href     string `json:"href,omitempty"`
	UserId   string `json:"user_id,omitempty"`
	UserName string `json:"user_name,omitempty"`
}

func (p *PostParagraph) IsValid() error {
	switch p.Tag {
	case TextPostParagraphType:
		if p.Text == "" {
			return errors.New("empty post text")
		}
	case HrefPostParagraphType:
		if p.Href == "" {
			return errors.New("empty post href")
		}
	case AtPostParagraphType:
		if p.UserId == "" {
			return errors.New("empty post at")
		}
	default:
		return fmt.Errorf("unsupported post paragraph type: %s", p.Tag)
	}
	return nil
}

type PostLangMessage struct {
	Title   string            `json:"title"`
	Content [][]PostParagraph `json:"content"`
}

func (p *PostLangMessage) IsValid() error {
	if len(p.Content) == 0 {
		return errors.New("empty post lang message")
	}
	for _, content := range p.Content {
		for _, paragraph := range content {
			if err := paragraph.IsValid(); err != nil {
				return err
			}
		}
	}
	return nil
}

type PostMessage struct {
	ZhCn *PostLangMessage `json:"zh_cn,omitempty"`
	EnUs *PostLangMessage `json:"en_us,omitempty"`
}

type InteractiveCard struct {
	Elements InteractiveCardElements `json:"elements"`
	Header   InteractiveCardHeader   `json:"header"`
}

type InteractiveCardElements []struct {
	Tag     string                         `json:"tag"`
	Text    InteractiveCardElementsText    `json:"text,omitempty"`
	Actions InteractiveCardElementsActions `json:"actions,omitempty"`
}

type InteractiveCardElementsText struct {
	Content string `json:"content"`
	Tag     string `json:"tag"`
}

type InteractiveCardElementsActions []struct {
	Tag   string                             `json:"tag"`
	Text  InteractiveCardElementsActionsText `json:"text"`
	Url   string                             `json:"url"`
	Type  string                             `json:"type"`
	Value struct {
	} `json:"value"`
}

type InteractiveCardElementsActionsText struct {
	Content string `json:"content"`
	Tag     string `json:"tag"`
}

type InteractiveCardHeader struct {
	Title InteractiveCardHeaderTitle `json:"title"`
}

type InteractiveCardHeaderTitle struct {
	Content string `json:"content"`
	Tag     string `json:"tag"`
}

type Content struct {
	Text TextMessage  `json:"text,omitempty"`
	Post *PostMessage `json:"post,omitempty"`
}

type Message struct {
	MsgType string           `json:"msg_type"`
	Content *Content         `json:"content,omitempty"`
	Card    *InteractiveCard `json:"card,omitempty"`
}

type signedMessage struct {
	Message
	Timestamp string `json:"timestamp"`
	Sign      string `json:"sign"`
}

func (m *Message) IsValid() error {
	switch m.MsgType {
	case TextMsgType:
		if m.Content == nil {
			return errors.New("empty content")
		}
	case PostMsgType:
		if m.Content == nil || m.Content.Post == nil {
			return errors.New("empty post")
		}
		if m.Content.Post.ZhCn == nil && m.Content.Post.EnUs == nil {
			return errors.New("empty post content")
		}
		if m.Content.Post.ZhCn != nil {
			if err := m.Content.Post.ZhCn.IsValid(); err != nil {
				return err
			}
		}
		if m.Content.Post.EnUs != nil {
			if err := m.Content.Post.EnUs.IsValid(); err != nil {
				return err
			}
		}
	case InteractiveMsgType:
		if m.Card == nil {
			return fmt.Errorf("empty card")
		}
	default:
		return fmt.Errorf("unsupported msg type: %s", m.MsgType)
	}
	return nil
}
