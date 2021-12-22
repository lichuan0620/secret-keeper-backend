package middlewares

import (
	"context"
	"net/http"
	"time"

	"github.com/go-logr/logr"
	"github.com/lichuan0620/secret-keeper-backend/pkg/service"
	"github.com/lichuan0620/secret-keeper-backend/pkg/service/model"
	"github.com/lichuan0620/secret-keeper-backend/pkg/telemetry/log"
)

// WithLogger adds the given logger to the request context.
func WithLogger(logger logr.Logger) model.Middleware {
	return func(ctx context.Context, f func(context.Context)) {
		f(log.SetContext(ctx, logger))
	}
}

// RequestLog uses the given logger to log every request that passes through this middleware.
func RequestLog(logger logr.Logger) model.Middleware {
	return func(ctx context.Context, f func(context.Context)) {
		start := time.Now()
		f(ctx)
		info := service.GetHandlingInfo(ctx)
		httpStatusCode := http.StatusOK
		logger = logger.V(log.LevelDefault)
		if info.Error != nil {
			logger = logger.WithValues("error_code", info.Error.GetCode()).V(log.LevelWarning)
			httpStatusCode = int(info.Error.GetHTTPCode())
		}
		logger.WithValues(
			"action", info.Action,
			"version", info.Version,
			"http_status_code", httpStatusCode,
			"duration_seconds", time.Since(start).Seconds(),
		).Info("request handled")
	}
}
