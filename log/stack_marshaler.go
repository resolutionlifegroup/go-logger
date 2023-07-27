package log

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"

	"github.com/pkg/errors"
)

type causer interface {
	Cause() error
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

// marshalStackForGCP is an adaptor of formatStack for zerolog.ErrorStackMarshaler
func marshalStackForGCP(err error) interface{} {
	formatted := formatStack(err)

	if len(formatted) == 0 {
		return nil
	}

	return string(formatted[:])
}

// formatStack is a best attempt at recreating the output of runtime.Stack
// Google Cloud's operations suite (formerly Stackdriver) currently only supports the stack output from runtime.Stack
// Code Source (ticket now closed): https://github.com/googleapis/google-cloud-go/issues/1084
// Issue on Google's side: https://issuetracker.google.com/issues/138952283
func formatStack(err error) (buffer []byte) {
	if err == nil {
		return
	}

	// find the inner most error with a stack
	inner := err
	for inner != nil {
		if cause, ok := inner.(causer); ok {
			inner = cause.Cause()
			if _, ok := inner.(stackTracer); ok {
				err = inner
			}
		} else {
			break
		}
	}

	if stackTrace, ok := err.(stackTracer); ok {
		buf := bytes.Buffer{}
		// routine id and state aren't available in pure go, so we hard-coded these
		buf.WriteString("goroutine 1 [running]:\n")

		// format each frame of the stack to match runtime.Stack's format
		var lines []string
		for _, frame := range stackTrace.StackTrace() {
			pc := uintptr(frame) - 1
			fn := runtime.FuncForPC(pc)
			if fn != nil {
				file, line := fn.FileLine(pc)
				lines = append(lines, fmt.Sprintf("%s()\n\t%s:%d +%#x", fn.Name(), file, line, fn.Entry()))
			}
		}
		buf.WriteString(strings.Join(lines, "\n"))

		buffer = buf.Bytes()
	}

	return
}
