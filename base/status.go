package base

import (
	"bytes"
	"fmt"

	mysqlerrnum "github.com/bombsimon/mysql-error-numbers/v2"
	"github.com/go-json-experiment/json/jsontext"
	"github.com/go-sql-driver/mysql"

	"google.golang.org/grpc/codes"
)

type Status struct {
	// use code to do condition test
	code codes.Code
	// event happened at pc
	pc uintptr
	// event for searching, event cause functionapi return code.
	event string
	// show msg to api user, let them know what happened
	// only return this msg to end user when code equal NotFound,
	// so do not contain sensitive message or unsuitable message
	// in this field.
	msg string
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
		errNum := mysqlerrnum.FromNumber(int(realErr.Number))
		s.event = fmt.Sprintf(
			"MySQLErrorNumber_%d_%s_%s,", realErr.Number, realErr.SQLState,
			errNum.String(),
		) + GetStackPosition(1)
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

func (s *Status) MarshalJSONV2(encoder *jsontext.Encoder, opts jsontext.Options) error {
	if s == nil {
		encoder.WriteToken(jsontext.Null)
		return nil
	}
	encoder.WriteToken(jsontext.ObjectStart)
	encoder.WriteToken(jsontext.String("code"))
	encoder.WriteToken(jsontext.Int(int64(s.code)))
	if s.event != "" {
		encoder.WriteToken(jsontext.String("event"))
		encoder.WriteToken(jsontext.String(s.event))
	}
	if s.rawEvent != "" {
		encoder.WriteToken(jsontext.String("raw_event"))
		encoder.WriteToken(jsontext.String(s.rawEvent))
	}
	if s.msg != "" {
		encoder.WriteToken(jsontext.String("msg"))
		encoder.WriteToken(jsontext.String(s.msg))
	}
	encoder.WriteToken(jsontext.ObjectEnd)
	return nil
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
