package xlog

import "log/slog"

func Error(err error) slog.Attr {
	if err == nil {
		return slog.Any("err", nil)
	}
	return slog.String("err", err.Error())
}
