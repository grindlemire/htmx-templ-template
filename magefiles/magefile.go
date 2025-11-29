package main

import (
	"context"
	"os"
	"time"

	"github.com/grindlemire/gothem-stack/pkg/log"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func init() {
	log.InitGlobal()
}

// Config is the identifying information pulled out of the environment to execute
// different mage commands
type Config struct {
	Env  string   `envconfig:"env" default:"local"`
	Args []string `envconfig:"args"    default:""`
}

// withBoilerplate wraps mage commands with panic recovery, timing, and context setup
func withBoilerplate(fn func(context.Context) error) (err error) {
	start := time.Now()
	defer func() {
		if r := recover(); r != nil {
			err = errors.Errorf("%s", r)
		}
		finish(start, err)
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	return fn(WithConfig(ctx, os.Args[2:]...))
}

// Install will install all dependencies
func Install() error { return withBoilerplate(install) }

// Tidy will run go mod tidy
func Tidy() error { return withBoilerplate(tidy) }

// Build will build a new binary
func Build() error { return withBoilerplate(build) }

// Run will run a local dev server and UI
func Run() error { return withBoilerplate(run) }

// Templ will run the templ command and generate the go code for the templates
func Templ() error { return withBoilerplate(templ) }

// Deploy will deploy the static site to Firebase hosting
func Deploy() error { return withBoilerplate(deploy) }

// Auth will authenticate with cloud services (gcloud or firebase)
func Auth() error { return withBoilerplate(auth) }

// Release will deploy a specific version of the service to Cloud Run
func Release() error { return withBoilerplate(release) }

// Destroy will delete all remote cloud infrastructure created during deploy
func Destroy() error { return withBoilerplate(destroy) }

func finish(start time.Time, err error) {
	zap.S().Infof("elapsed time: %s", time.Since(start))
	// This is a hack to get around the fact that mage treats command line args
	// as other targets. In a run we just want to interpret them as arugments to the binary,
	// not as other targets. So we just short circuit and tell mage to stop
	if err != nil {
		zap.S().Fatal(err)
	}
	os.Exit(0)
}
