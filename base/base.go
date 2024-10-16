package base

import (
	"context"
	"fmt"
	"log/slog"
)

func PanicIfErr(c context.Context, err error, attrs ...slog.Attr) {
	if err == nil {
		return
	}
	panic(err)
}

func ErrWithCaller(err error) error {
	if err == nil {
		return nil
	}
	return &errWithCaller{
		err: err,
		pc:  GetCaller(1),
	}
}

type errWithCaller struct {
	err error
	pc  uintptr
}

func (e *errWithCaller) Error() string {
	pos := GetFilePosition(e.pc)
	return fmt.Sprintf("%s(%s)",
		e.err.Error(),
		pos,
	)
}
