package xlog

import (
	"github.com/ethanvc/logjson/slogjson"
	"log/slog"
	"os"
)

func InitDefaultLogger() error {
	opt := &slogjson.HandlerOption{
		Writer: os.Stderr,
	}
	handler := slogjson.NewHandler(opt)
	logger := slog.New(handler)
	slog.SetDefault(logger)
	return nil
}
