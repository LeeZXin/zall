package pullrequestmd

import "encoding/json"

type ActionType int

const (
	CommentType ActionType = iota + 1
	ReplyType
	PrType
)

type Pr struct {
	Id     int64    `json:"id"`
	Status PrStatus `json:"status"`
}

type Reply struct {
	FromId       int64  `json:"fromId"`
	FromComment  string `json:"fromComment"`
	FromAccount  string `json:"fromAccount"`
	ReplyComment string `json:"replyComment"`
}

type Comment struct {
	Comment string `json:"comment"`
}

type Action struct {
	Comment    *Comment   `json:"comment,omitempty"`
	Reply      *Reply     `json:"reply,omitempty"`
	Pr         *Pr        `json:"pr,omitempty"`
	ActionType ActionType `json:"actionType"`
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
		Comment: &Comment{
			Comment: comment,
		},
		ActionType: CommentType,
	}
}

func NewReplyAction(fromId int64, fromAccount, fromComment, replyComment string) Action {
	return Action{
		Reply: &Reply{
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
		Pr: &Pr{
			Id:     prId,
			Status: status,
		},
		ActionType: PrType,
	}
}
