package telemetry

import (
	"context"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/lichuan0620/secret-keeper-backend/pkg/telemetry/healthz"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// for various short operations
const timeout = 3 * time.Second

const (
	metricsEndpoint = "/metrics"
	healthzEndpoint = "/healthz"
	pprofEndpoint   = "/debug/pprof"
)

// Server is used to manage and serve telemetry APIs.
type Server struct {
	base         *http.Server
	mux          *http.ServeMux
	healthChecks map[string]healthz.Checker
}

// ServerOptions is used to build a Server.
type ServerOptions struct {
	ListenAddress string
}

// DefaultServerOptions returns a ServerOptions with default values.
func DefaultServerOptions() *ServerOptions {
	return &ServerOptions{
		ListenAddress: ":9090",
	}
}

// NewServer builds a Server. A new Server can be modified further until it had started running.
func NewServer(options *ServerOptions) *Server {
	if options == nil {
		options = DefaultServerOptions()
	}
	mux := http.NewServeMux()
	healthChecks := map[string]healthz.Checker{
		"ping": healthz.Ping,
	}
	mux.HandleFunc(healthzEndpoint, func(writer http.ResponseWriter, request *http.Request) {
		for _, check := range healthChecks {
			if err := check(request); err != nil {
				writer.WriteHeader(http.StatusInternalServerError)
				_, _ = writer.Write([]byte(err.Error()))
				return
			}
		}
		writer.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc(pprofEndpoint, pprof.Index)
	mux.HandleFunc(pprofEndpoint+"/cmdline", pprof.Cmdline)
	mux.HandleFunc(pprofEndpoint+"/profile", pprof.Profile)
	mux.HandleFunc(pprofEndpoint+"/symbol", pprof.Symbol)
	mux.HandleFunc(pprofEndpoint+"/trace", pprof.Trace)
	mux.Handle("/metrics", promhttp.Handler())
	server := &http.Server{
		Addr:              options.ListenAddress,
		Handler:           mux,
		IdleTimeout:       90 * time.Second,
		ReadHeaderTimeout: timeout,
	}
	return &Server{
		base:         server,
		mux:          mux,
		healthChecks: healthChecks,
	}
}

// SetHealthCheck adds or updates a named health checker. This should not be used after the server
// has started.
func (server *Server) SetHealthCheck(name string, checker healthz.Checker) {
	server.healthChecks[name] = checker
}

// Start runs the telemetry HTTP server and blocks until it has completed the shutdown process.
func (server *Server) Start(ctx context.Context) error {
	serveErrChan := make(chan error)
	go func() {
		if err := server.base.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serveErrChan <- err
		}
	}()
	select {
	case <-ctx.Done():
	case err := <-serveErrChan:
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return server.base.Shutdown(ctx)
}
