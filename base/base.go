package base

import (
	"context"
	"log/slog"
)

func PanicIfErr(c context.Context, err error, attrs ...slog.Attr) {
	if err == nil {
		return
	}
	panic(err)
}
