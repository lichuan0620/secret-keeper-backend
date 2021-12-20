package server

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/lichuan0620/secret-keeper-backend/cmd/server/service"
	"github.com/lichuan0620/secret-keeper-backend/internal/queueclient"
	"github.com/lichuan0620/secret-keeper-backend/pkg/mongo"
	"github.com/lichuan0620/secret-keeper-backend/pkg/telemetry"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

const component = "server"

func Command(ctx context.Context) *cobra.Command {
	var (
		MongoEndpoint          string
		QueueEndpoint          string
		ListenAddress          string
		TelemetryListenAddress string
	)
	cmd := cobra.Command{
		Use:   component,
		Short: "Run a secret-keeper server instance",
	}
	flags := cmd.PersistentFlags()
	flags.StringVar(&MongoEndpoint, "mongodb-endpoint", os.Getenv("MONGODB_ENDPOINT"), "address to the MongoDB service")
	flags.StringVar(&QueueEndpoint, "queue-endpoint", os.Getenv("QUEUE_ENDPOINT"), "address to the secret-keeper queue service")
	flags.StringVar(&ListenAddress, "listen-address", os.Getenv("LISTEN_ADDRESS"), "address to listen to for HTTP requests")
	flags.StringVar(&TelemetryListenAddress, "telemetry-listen-address", os.Getenv("TELEMETRY_LISTEN_ADDRESS"), "address to listen to for telemetry requests")
	cmd.RunE = func(_ *cobra.Command, _ []string) error {
		if err := mongo.Init(MongoEndpoint); err != nil {
			return errors.Wrap(err, "initialize MongoDB connection")
		}
		qc := queueclient.New(QueueEndpoint)
		handler, err := service.Build(qc)
		if err != nil {
			return errors.Wrap(err, "build service handler")
		}
		mux := http.NewServeMux()
		mux.Handle("/", handler)
		server := http.Server{
			Addr:         ListenAddress,
			Handler:      mux,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  3 * time.Minute,
		}
		telemetryServer := telemetry.NewServer(&telemetry.ServerOptions{ListenAddress: TelemetryListenAddress})
		eg, egCtx := errgroup.WithContext(ctx)
		eg.Go(func() error {
			return errors.Wrap(telemetryServer.Start(egCtx), "serve telemetry")
		})
		eg.Go(func() error {
			if err = server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				return errors.Wrap(err, "serve HTTP")
			}
			return nil
		})
		eg.Go(func() error {
			<-egCtx.Done()
			gracefulCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			return errors.Wrap(server.Shutdown(gracefulCtx), "HTTP server graceful shutdown")
		})
		return eg.Wait()
	}
	return &cmd
}
