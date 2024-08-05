package obs

import (
	"context"
	"github.com/ethanvc/evol/base"
	"google.golang.org/grpc/codes"
)

type Composer struct {
	c         context.Context
	originErr error
	s         *base.Status
}

func New(c context.Context, err error) *Composer {
	if c == nil {
		c = context.Background()
	}
	com := &Composer{
		c:         c,
		originErr: err,
	}
	com.convertToStatus()
	return com
}

func (com *Composer) convertToStatus() {
	if com.originErr == nil {
		return
	}
	if s, ok := com.originErr.(*base.Status); ok {
		com.s = s
		return
	}
	com.s = base.New(codes.Internal).SetMsg(com.originErr.Error())
}

func (com *Composer) Report() *Composer {
	return com
}

func (com *Composer) Error() error {
	if com.s == nil {
		return nil
	}
	return com.s
}

func (com *Composer) Status() *base.Status {
	return com.s
}
