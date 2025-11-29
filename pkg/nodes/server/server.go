package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/grindlemire/gothem-stack/pkg/nodes/config"
	"github.com/grindlemire/gothem-stack/pkg/nodes/router"
	"go.uber.org/zap"

	"github.com/grindlemire/graft"
	"github.com/pkg/errors"
)

const ID = graft.ID("server")

type Output struct {
	Err error
}

func init() {
	graft.Register(graft.Node[Output]{
		ID:        ID,
		Cacheable: true,
		DependsOn: []graft.ID{config.ID, router.ID},
		Run:       run,
	})
}

func run(ctx context.Context) (Output, error) {
	config, err := graft.Dep[config.Output](ctx)
	if err != nil {
		return Output{}, errors.Wrap(err, "getting config")
	}

	router, err := graft.Dep[router.Output](ctx)
	if err != nil {
		return Output{}, errors.Wrap(err, "getting router")
	}

	if config.Server.LocalCerts {
		if !hasCerts() {
			_, err := generateCerts()
			if err != nil {
				return Output{}, errors.Wrap(err, "generating certs")
			}
		}
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Server.Port),
		Handler: router,
	}

	err = start(ctx, server, config)
	if err != nil {
		return Output{
			Err: errors.Wrap(err, "starting server"),
		}, nil
	}

	return Output{
		Err: nil,
	}, nil
}

func start(ctx context.Context, server *http.Server, config config.Output) (err error) {
	// run the listeners in their own goroutine, this is so we can properly propagate signals
	// and cleanup everything since there may be other signals that need to be cleaned up.
	errCh := make(chan error, 1)
	go func() {
		zap.S().Infof("started listening on :%d", config.Server.Port)
		if config.Server.LocalCerts {
			zap.S().Debug(ctx, "listening with tls")
			err := server.ListenAndServeTLS(publicKeyFile, privateKeyFile)
			errCh <- errors.Wrap(err, "starting server")
			return
		}
		err := server.ListenAndServe()
		errCh <- errors.Wrap(err, "starting server")
	}()

	// wait for either the context to be cancelled indicating we should shutdown
	// or for the servers to fail
	for {
		select {
		case <-ctx.Done():
			shutdownCTX, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			err = server.Shutdown(shutdownCTX)
			if err != nil {
				return err
			}
			return ctx.Err()
		case err := <-errCh:
			return err
		}
	}
}
