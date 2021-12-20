package main

import (
	"flag"
	"os"

	"github.com/lichuan0620/secret-keeper-backend/cmd/queue"
	"github.com/lichuan0620/secret-keeper-backend/cmd/server"
	"github.com/lichuan0620/secret-keeper-backend/pkg/signal"
	"github.com/spf13/cobra"
)

var command *cobra.Command

func main() {
	ctx := signal.SetupBackgroundContext()
	command = &cobra.Command{
		Use:   "secret-keeper",
		Short: "Commands of Secret Keeper backend.",
	}
	command.AddCommand(
		queue.Command(ctx),
		server.Command(ctx),
	)
	flags := command.PersistentFlags()
	flags.AddGoFlagSet(flag.CommandLine)
	if command.Execute() != nil {
		os.Exit(1)
	}
}
