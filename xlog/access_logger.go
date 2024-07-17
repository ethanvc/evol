package xlog

import (
	"context"
	"github.com/ethanvc/evol/base"
	"google.golang.org/grpc/codes"
	"log/slog"
	"time"
)

type AccessLogger struct {
	conf *AccessLoggerConfig
}

func NewAccessLogger(conf *AccessLoggerConfig) *AccessLogger {
	al := &AccessLogger{
		conf: conf,
	}
	return al
}

var accessLogger *AccessLogger = NewAccessLogger(NewDefaultAccessLoggerConfig())

func GetAccessLogger() *AccessLogger {
	return accessLogger
}

func (al *AccessLogger) LogAccess(c context.Context, skip int, err error, req, resp any, extra ...any) {
	lc := GetLogContext(c)
	now := time.Now()
	record := slog.NewRecord(lc.GetStartTime(), al.conf.GetLogLevel(err), "REQ_END", base.GetCaller(skip+1))
	record.Add("method", lc.GetMethod())
	record.Add("err", err)
	record.Add("req", req)
	record.Add("resp", resp)
	record.Add(extra...)
	record.Add("tc_us", now.Sub(now).Microseconds())
	handler := slog.Default().Handler()
	handler.Handle(c, record)
}

func (al *AccessLogger) ReportInfo(c context.Context) {}

func (al *AccessLogger) ReportErr(c context.Context) {}

func (al *AccessLogger) ReportError(c context.Context, err error) {}

func (al *AccessLogger) ReportAccess(c context.Context, err error) {}

func (al *AccessLogger) Report(c context.Context, lvl MonitorLevel, event string) {}

func (al *AccessLogger) ReportDuration(c context.Context) {}

type MonitorLevel string

const (
	MonitorLevelInfo MonitorLevel = "info"
	MonitorLevelErr  MonitorLevel = "err"
)

type AccessLoggerConfig struct {
	GetLogLevel func(err error) slog.Level
}

func NewDefaultAccessLoggerConfig() *AccessLoggerConfig {
	return &AccessLoggerConfig{
		GetLogLevel: func(err error) slog.Level {
			switch base.Code(err) {
			case codes.Canceled, codes.Unknown, codes.DeadlineExceeded,
				codes.ResourceExhausted, codes.Aborted, codes.Internal, codes.Unavailable:
				return slog.LevelError
			default:
				return slog.LevelInfo
			}
		},
	}
}
