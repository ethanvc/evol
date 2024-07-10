package base

import (
	"bytes"
	"fmt"
	"google.golang.org/grpc/codes"
)

type Status struct {
	code  codes.Code
	event string
	msg   string
}

func New(code codes.Code, event string) *Status {
	return &Status{
		code:  code,
		event: event,
	}
}

func (s *Status) GetCode() codes.Code {
	if s == nil {
		return codes.OK
	}
	return s.code
}

func (s *Status) GetEvent() string {
	if s == nil {
		return ""
	}
	return s.event
}

func (s *Status) GetMsg() string {
	if s == nil {
		return ""
	}
	return s.msg
}

func (s *Status) SetMsg(msg string, args ...any) *Status {
	s.msg = fmt.Sprintf(msg, args...)
	return s
}

func (s *Status) Error() string {
	if s == nil {
		return codes.OK.String()
	}
	var buf bytes.Buffer
	buf.WriteString(s.code.String())
	if s.event != "" {
		buf.WriteByte(';')
		buf.WriteString(s.event)
	}
	if s.msg != "" {
		buf.WriteByte(';')
		buf.WriteString(s.msg)
	}
	return buf.String()
}
