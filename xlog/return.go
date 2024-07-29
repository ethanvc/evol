package xlog

import "context"

type Composer struct {
	c   context.Context
	err error
}

func New(c context.Context, err error) *Composer {
	if c == nil {
		c = context.Background()
	}
	return &Composer{
		c:   c,
		err: err,
	}
}

func (com *Composer) Report() *Composer {
	return com
}

func (com *Composer) Error() error {
	return com.err
}
