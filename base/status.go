package base

import (
	"bytes"
	"fmt"
	"github.com/go-sql-driver/mysql"

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
	if err == nil {
		return s
	}

	switch realErr := err.(type) {
	case *Status:
		s.event = realErr.GetEvent()
		s.rawEvent = realErr.GetRawEvent()
		return s
	case *mysql.MySQLError:
		// message returned by mysql may contain insert data, which is not good as monitor event.
		s.event = fmt.Sprintf("MySQLErrorNumber_%d_%s;", realErr.Number, realErr.SQLState) + GetStackPosition(1)
		s.rawEvent = realErr.Error()
		return s
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
