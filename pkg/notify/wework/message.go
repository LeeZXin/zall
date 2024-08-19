package wework

import (
	"fmt"
	"github.com/pingcap/errors"
)

const (
	TextType         = "text"
	MarkdownType     = "markdown"
	ImageType        = "image"
	NewsType         = "news"
	TemplateCardType = "template_card"
)

type TextMessage struct {
	Content           string   `json:"content"`
	MentionedList     []string `json:"mentioned_list"`
	MentionMobileList []string `json:"mentioned_mobile_list"`
}

type MarkdownMessage struct {
	Content string `json:"content"`
}

type ImageMessage struct {
	Base64 string `json:"base64"`
	Md5    string `json:"md5"`
}

type NewsMessage struct {
	Articles []struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Url         string `json:"url"`
		PicUrl      string `json:"picurl"`
	} `json:"articles"`
}

type TemplateCardMessage struct {
	CardType string `json:"card_type"`
	Source   struct {
		IconUrl   string `json:"icon_url"`
		Desc      string `json:"desc"`
		DescColor int    `json:"desc_color"`
	} `json:"source"`
	MainTitle struct {
		Title string `json:"title"`
		Desc  string `json:"desc"`
	} `json:"main_title"`
	CardImage struct {
		Url         string  `json:"url"`
		AspectRatio float64 `json:"aspect_ratio"`
	} `json:"card_image"`
	ImageTextArea struct {
		Type     int    `json:"type"`
		Url      string `json:"url"`
		Title    string `json:"title"`
		Desc     string `json:"desc"`
		ImageUrl string `json:"image_url"`
	} `json:"image_text_area"`
	QuoteArea struct {
		Type      int    `json:"type"`
		Url       string `json:"url"`
		Appid     string `json:"appid"`
		Pagepath  string `json:"pagepath"`
		Title     string `json:"title"`
		QuoteText string `json:"quote_text"`
	} `json:"quote_area"`
	VerticalContentList []struct {
		Title string `json:"title"`
		Desc  string `json:"desc"`
	} `json:"vertical_content_list"`
	HorizontalContentList []struct {
		Keyname string `json:"keyname"`
		Value   string `json:"value"`
		Type    int    `json:"type"`
		Url     string `json:"url,omitempty"`
		MediaId string `json:"media_id,omitempty"`
	} `json:"horizontal_content_list"`
	JumpList []struct {
		Type     int    `json:"type"`
		Url      string `json:"url,omitempty"`
		Title    string `json:"title"`
		Appid    string `json:"appid,omitempty"`
		Pagepath string `json:"pagepath,omitempty"`
	} `json:"jump_list"`
	CardAction struct {
		Type     int    `json:"type"`
		Url      string `json:"url"`
		Appid    string `json:"appid"`
		Pagepath string `json:"pagepath"`
	} `json:"card_action"`
	EmphasisContent *struct {
		Title string `json:"title"`
		Desc  string `json:"desc"`
	} `json:"emphasis_content"`
	SubTitleText string `json:"sub_title_text,omitempty"`
}

type Message struct {
	Msgtype      string               `json:"msgtype"`
	Text         *TextMessage         `json:"text,omitempty"`
	Image        *ImageMessage        `json:"image,omitempty"`
	Markdown     *MarkdownMessage     `json:"markdown,omitempty"`
	News         *NewsMessage         `json:"news,omitempty"`
	TemplateCard *TemplateCardMessage `json:"template_card,omitempty"`
}

func (m *Message) IsValid() error {
	switch m.Msgtype {
	case TextType:
		if m.Text == nil {
			return errors.New("empty text message")
		}
		return nil
	case ImageType:
		if m.Image == nil {
			return errors.New("empty image message")
		}
		return nil
	case MarkdownType:
		if m.Markdown == nil {
			return errors.New("empty markdown message")
		}
		return nil
	case NewsType:
		if m.News == nil {
			return errors.New("empty news message")
		}
		return nil
	case TemplateCardType:
		if m.TemplateCard == nil {
			return errors.New("empty template card message")
		}
		return nil
	default:
		return fmt.Errorf("unsupport type: %s", m.Msgtype)
	}
}
