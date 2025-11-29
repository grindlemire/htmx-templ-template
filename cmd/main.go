package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/grindlemire/gothem-stack/pkg/log"
	"github.com/grindlemire/gothem-stack/pkg/nodes/server"
	"github.com/pkg/errors"

	"github.com/grindlemire/graft"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

func main() {
	log.InitGlobal()

	app := &cli.App{
		Name:  "serve",
		Usage: "serve an htmx api",
		Action: func(c *cli.Context) (err error) {
			ctx, cancel := context.WithCancel(c.Context)
			defer cancel()

			// Create a signal channel
			sigCh := make(chan os.Signal, 1)
			// Register a signal handler for SIGINT
			signal.Notify(sigCh, os.Interrupt)

			go func() {
				<-sigCh
				cancel()
			}()

			// execute the server node which will run all the server dependencies
			// and start the server. It will return an error if it fails to run.
			output, _, err := graft.ExecuteFor[server.Output](ctx)
			if err != nil {
				return errors.Wrap(err, "server execution failed")
			}

			return errors.Wrap(output.Err, "server run failed")
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		// we don't care about context cancellation as that happens if we kill the process
		// while it is waiting for a request to finish
		if errors.Is(err, context.Canceled) {
			os.Exit(1)
		}
		zap.S().Fatal(err)
	}
}
