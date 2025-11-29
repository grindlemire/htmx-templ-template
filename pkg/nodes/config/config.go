package config

import (
	"context"

	"github.com/grindlemire/graft"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

const ID = graft.ID("config")

type Output struct {
	Server ServerConfig
}

// ServerConfig is configuration for the server parsed from the env.
type ServerConfig struct {
	Port       int  `envconfig:"PORT"              default:"4433"`
	LocalCerts bool `envconfig:"LOCAL_CERTS"       default:"false" split_words:"true"`
}

func init() {
	graft.Register(graft.Node[Output]{
		ID: ID,
		// cache the config so multiple nodes can use it without reloading
		Cacheable: true,
		Run:       run,
	})
}

// run will load the server config from the environment and return it.
func run(ctx context.Context) (Output, error) {
	var config ServerConfig
	err := envconfig.Process("", &config)
	if err != nil {
		return Output{}, errors.Wrap(err, "loading environment")
	}

	return Output{
		Server: config,
	}, nil
}
