package wraperr

import (
	"github.com/LeeZXin/zsf-utils/bizerr"
)

type Err struct {
	Internal error `json:"-"`
	*bizerr.Err
}

func WrapErr(err error, bizErr *bizerr.Err) *Err {
	return &Err{
		Internal: err,
		Err:      bizErr,
	}
}
