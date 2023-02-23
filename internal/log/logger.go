package log

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"
)

type TraceLogger interface{}

type traceLoggerInternal interface {
	// Internal setup/teardown functions
	logStart(ctx context.Context) context.Context
	logExit(ctx context.Context, exitCode int)

	TraceLogger
}

func WithTraceLogger(
	ctx context.Context,
	mainFunc func(context.Context, TraceLogger),
) {
	logger := NewTrace2()

	// Set up the program-level context
	ctx = logger.logStart(ctx)
	defer func() {
		if panicInfo := recover(); panicInfo != nil {
			// Panicking - log, print panic info, then exit
			logger.logExit(ctx, 1)
			os.Stderr.WriteString(fmt.Sprintf("panic: %s\n\n", panicInfo))
			debug.PrintStack()
			os.Exit(1)
		} else {
			// Just log the exit (but don't os.Exit()) so we can exit normally
			logger.logExit(ctx, 0)
		}
	}()

	mainFunc(ctx, logger)
}
