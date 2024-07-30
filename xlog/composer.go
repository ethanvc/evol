package xlog

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
	event := ConvertToEventString(com.originErr.Error(), 80) + ";" + GetStackPosition(2)
	com.s = base.New(codes.Internal, event).SetMsg(com.originErr.Error())
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
