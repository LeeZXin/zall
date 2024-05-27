package pullrequestmd

import "encoding/json"

type ActionType int

const (
	CommentType ActionType = iota + 1
	ReplyType
	PrType
	ReviewType
)

func (t ActionType) IsRelatedToComment() bool {
	switch t {
	case CommentType, ReplyType:
		return true
	default:
		return false
	}
}

type PrAction struct {
	Id     int64    `json:"id"`
	Status PrStatus `json:"status"`
}

type ReplyAction struct {
	FromId       int64  `json:"fromId"`
	FromComment  string `json:"fromComment"`
	FromAccount  string `json:"fromAccount"`
	ReplyComment string `json:"replyComment"`
}

type CommentAction struct {
	Comment string `json:"comment"`
}

type ReviewAction struct {
	ReviewId int64 `json:"reviewId"`
}

type Action struct {
	Comment    *CommentAction `json:"comment,omitempty"`
	Reply      *ReplyAction   `json:"reply,omitempty"`
	Pr         *PrAction      `json:"pr,omitempty"`
	Review     *ReviewAction  `json:"review,omitempty"`
	ActionType ActionType     `json:"actionType"`
}

func (c *Action) GetCommentText() string {
	switch c.ActionType {
	case CommentType:
		if c.Comment != nil {
			return c.Comment.Comment
		}
	case ReplyType:
		if c.Reply != nil {
			return c.Reply.ReplyComment
		}
	}
	return ""
}

func (c *Action) FromDB(content []byte) error {
	if c == nil {
		*c = Action{}
	}
	return json.Unmarshal(content, c)
}

func (c *Action) ToDB() ([]byte, error) {
	return json.Marshal(c)
}

func NewCommentAction(comment string) Action {
	return Action{
		Comment: &CommentAction{
			Comment: comment,
		},
		ActionType: CommentType,
	}
}

func NewReplyAction(fromId int64, fromAccount, fromComment, replyComment string) Action {
	return Action{
		Reply: &ReplyAction{
			FromId:       fromId,
			FromComment:  fromComment,
			FromAccount:  fromAccount,
			ReplyComment: replyComment,
		},
		ActionType: ReplyType,
	}
}

func NewPrAction(prId int64, status PrStatus) Action {
	return Action{
		Pr: &PrAction{
			Id:     prId,
			Status: status,
		},
		ActionType: PrType,
	}
}

func NewReviewAction(reviewId int64) Action {
	return Action{
		Review: &ReviewAction{
			ReviewId: reviewId,
		},
		ActionType: ReviewType,
	}
}
