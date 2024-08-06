package base

import (
	"bytes"
	"fmt"

	"google.golang.org/grpc/codes"
)

type Status struct {
	code     codes.Code
	event    string
	rawEvent string
	msg      string
}

func New(code codes.Code) *Status {
	return &Status{
		code: code,
	}
}

func (s *Status) SetEvent(event string) *Status {
	s.event = event
	return s
}

func (s *Status) SetErrEvent(err error) *Status {
	if realErr, ok := err.(*Status); ok {
		s.event = realErr.GetEvent()
		s.rawEvent = realErr.GetRawEvent()
	}
	s.event = ToEventString(err.Error(), 0)
	s.event += ";" + GetStackPosition(1)
	s.rawEvent = err.Error()
	return s
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

func (s *Status) GetRawEvent() string {
	if s == nil {
		return ""
	}
	return s.rawEvent
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

func Code(err error) codes.Code {
	if err == nil {
		return codes.OK
	}
	status, ok := err.(*Status)
	if !ok {
		return codes.Unknown
	}
	return status.GetCode()
}
