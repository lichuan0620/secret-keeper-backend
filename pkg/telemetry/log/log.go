package log

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/klog/v2/klogr"
)

// Log level
// from: https://github.com/kubernetes/community/blob/master/contributors/devel/sig-instrumentation/logging.md
const (
	// LevelWarning represents warning or low-priority errors
	LevelWarning = 0
	// LevelDefault represents important non-error information (default level if you are unsure)
	LevelDefault = 1
	// LevelState represents resource state change or request log
	LevelState = 2
	// LevelExtended represents verbose system information
	LevelExtended = 3
	// LevelDebug represents logging in particularly thorny parts of code where you may want to come back later and check it
	LevelDebug = 4
	// LevelTrace represents more information for troubleshooting reported issues
	LevelTrace = 5
)

var contextKeyLogger interface{} = new(byte)

// FromContext returns the pre-built logger from the given context if it exists, and construct
// a new one if it does not.
func FromContext(ctx context.Context) logr.Logger {
	v := ctx.Value(contextKeyLogger)
	if v == nil {
		return New()
	}
	return v.(logr.Logger)
}

// SetContext returns a copy of the parent context that contains the given logger. The logger
// can be retrieved later by calling FromContext.
func SetContext(ctx context.Context, logger logr.Logger) context.Context {
	return context.WithValue(ctx, contextKeyLogger, logger)
}

// New returns a new logger.
func New() logr.Logger {
	return klogr.New().V(LevelDefault)
}
