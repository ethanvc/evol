package xlog

import (
	"fmt"
	"runtime"
	"strings"
)

func GetCaller(skip int) uintptr {
	var pcs [1]uintptr
	const ParentSkipCount = 2
	runtime.Callers(skip+ParentSkipCount, pcs[:])
	return pcs[0]
}

func GetCallerFrame(pc uintptr) runtime.Frame {
	fs := runtime.CallersFrames([]uintptr{pc})
	frame, _ := fs.Next()
	return frame
}

func GetFilePathTailPart(filePath string, count int) string {
	currentCnt := 0
	for i := len(filePath) - 1; i >= 0; i-- {
		if filePath[i] == '/' || filePath[i] == '\\' {
			currentCnt++
			if currentCnt >= count {
				return filePath[i+1:]
			}
		}
	}
	return filePath
}

func GetStackPosition(skip int) string {
	pc := GetCaller(skip + 1)
	frame := GetCallerFrame(pc)
	f := GetFilePathTailPart(frame.File, 2)
	return fmt.Sprintf("%s:%d", f, frame.Line)
}

// GetPanicReportEvent only call this function in recover code
func GetPanicReportEvent(skip int, r any) string {
	var errStr string
	switch realR := r.(type) {
	case error:
		errStr = realR.Error()
	case string:
		errStr = realR
	}

	var pcs [10]uintptr
	const ParentSkipCount = 4
	cnt := runtime.Callers(skip+ParentSkipCount, pcs[:])
	if cnt == 0 {
		return "RuntimePanic;CallersReturnNothing;" + ConvertToEventString(errStr, 0)
	}
	realPcs := pcs[0:cnt]
	frames := runtime.CallersFrames(realPcs)

	more := true
	var frame runtime.Frame
	for {
		if !more {
			return "RuntimePanic;NotFoundBusinessCode;" + ConvertToEventString(errStr, 0)
		}
		frame, more = frames.Next()
		if strings.Contains(frame.File, "/src/runtime/") {
			continue
		}
		break
	}
	const keepTailPart = 2
	s := GetFilePathTailPart(frame.File, keepTailPart)
	return fmt.Sprintf("RuntimePanic;%s:%d;", s, frame.Line) + ConvertToEventString(errStr, 0)
}
