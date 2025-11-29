package main

import (
	"context"
	"io"
	"os"
	"strings"

	"github.com/grindlemire/gothem-stack/magefiles/cmd"
	"github.com/grindlemire/gothem-stack/web/pages/home"

	"github.com/magefile/mage/mg"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func build(ctx context.Context) error {
	config, err := GetConfig(ctx)
	if err != nil {
		return err
	}
	zap.S().Infof("mage running with configs: %+v", config)

	// Run dependencies in order
	mg.SerialCtxDeps(ctx, tidy, templ)

	// Ensure dist directory exists
	if err := os.MkdirAll("dist", 0755); err != nil {
		return errors.Wrap(err, "creating dist directory")
	}

	// Clean previous build
	if err := cmd.Run(ctx, cmd.WithCMD("rm", "-f", "dist/server")); err != nil {
		return errors.Wrap(err, "cleaning previous build")
	}

	// Build the server
	err = cmd.Run(ctx,
		cmd.WithCMD(
			"go",
			"build",
			"-o", "dist/server",
			"cmd/main.go",
		),
	)
	if err != nil {
		return errors.Wrap(err, "building server")
	}

	// Copy static assets to dist
	if err := generateStaticAssets(ctx); err != nil {
		return errors.Wrap(err, "handling static assets")
	}

	zap.S().Info("Build completed successfully")
	return nil
}

// generateStaticAssets generates CSS and copies static files to dist
func generateStaticAssets(ctx context.Context) error {
	// Generate CSS using TailwindCSS
	err := cmd.Run(ctx,
		cmd.WithDir("./web"),
		cmd.WithCMD(
			"node_modules/.bin/tailwindcss",
			"-i", "tailwind.css",
			"-o", "public/styles.min.css",
		),
	)
	if err != nil {
		return errors.Wrap(err, "generating CSS")
	}

	// Create dist/public/dist directory
	if err := os.MkdirAll("dist/public/dist", 0755); err != nil {
		return errors.Wrap(err, "creating dist/public/dist directory")
	}

	// Generate the static HTML
	if err := generateHTML("dist/public/index.html"); err != nil {
		return errors.Wrap(err, "generating HTML")
	}

	// Copy static assets
	files := map[string]string{
		"./web/public/favicon.ico":    "dist/public/dist/favicon.ico",
		"./web/public/styles.min.css": "dist/public/dist/styles.min.css",
	}

	for src, dst := range files {
		if err := copyFile(src, dst); err != nil {
			return errors.Wrapf(err, "copying %s to %s", src, dst)
		}
	}

	return nil
}

// generateHTML renders the home page and writes it to the specified path
func generateHTML(outputPath string) error {
	var s strings.Builder
	if err := home.Page("").Render(context.Background(), &s); err != nil {
		return errors.Wrap(err, "rendering static page")
	}

	if err := os.WriteFile(outputPath, []byte(s.String()), 0644); err != nil {
		return errors.Wrap(err, "writing HTML to file")
	}

	return nil
}

func copyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return errors.Wrapf(err, "opening source file %s", src)
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return errors.Wrapf(err, "creating destination file %s", dst)
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}
