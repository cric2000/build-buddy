// Package app shares command handling and server lifecycle across entry points.
package app

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/cric2000/build-buddy/internal/api"
	"github.com/cric2000/build-buddy/internal/config"
	"github.com/cric2000/build-buddy/internal/mascot"
)

// Main runs the CLI or HTTP server and exits nonzero on failure.
func Main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := run(ctx, os.Args[1:], os.Getenv, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, "buildbuddy:", err)
		os.Exit(1)
	}
}
func run(ctx context.Context, args []string, getenv func(string) string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		args = []string{"serve"}
	}
	switch args[0] {
	case "help", "-h", "--help":
		fmt.Fprintln(stdout, "BuildBuddy — Your Build Has Feelings\n\nUsage: buildbuddy [serve | render | help]\n\n  serve   Start HTTP API (HOST=0.0.0.0 PORT=8080)\n  render  Generate SVG locally; use render --help for flags")
		return nil
	case "render":
		return render(args[1:], stdout, stderr)
	case "serve":
		fs := flag.NewFlagSet("serve", flag.ContinueOnError)
		fs.SetOutput(stderr)
		fs.Usage = func() { fmt.Fprintln(stderr, "Usage: buildbuddy serve\nEnvironment: HOST=0.0.0.0 PORT=8080") }
		if err := fs.Parse(args[1:]); err != nil {
			if errors.Is(err, flag.ErrHelp) {
				return nil
			}
			return err
		}
		if fs.NArg() != 0 {
			return fmt.Errorf("serve accepts no positional arguments")
		}
		cfg, err := config.Load(getenv)
		if err != nil {
			return err
		}
		logger := slog.New(slog.NewJSONHandler(stderr, nil))
		listener, err := net.Listen("tcp", cfg.Address)
		if err != nil {
			return fmt.Errorf("listen: %w", err)
		}
		server := api.Server(cfg.Address, logger)
		logger.Info("BuildBuddy listening", "address", listener.Addr().String())
		return serve(ctx, server, listener, logger)
	default:
		return fmt.Errorf("unknown command %q; use buildbuddy --help", args[0])
	}
}
func serve(ctx context.Context, server *http.Server, listener net.Listener, logger *slog.Logger) error {
	done := make(chan error, 1)
	go func() { done <- server.Serve(listener) }()
	select {
	case err := <-done:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
		logger.Info("shutting down")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()
			<-done
			return fmt.Errorf("shutdown: %w", err)
		}
		err := <-done
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
}
func render(args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("render", flag.ContinueOnError)
	fs.SetOutput(stderr)
	values := map[string]*singleFlag{}
	for _, f := range []struct{ name, def, help string }{
		{"status", "unknown", "passing, failing, running, queued, unknown"}, {"character", mascot.Default().Character, strings.Join(mascot.CharacterNames(), ", ")}, {"theme", "light", "light, dark, neon"},
		{"coverage", "", "coverage percentage (0–100)"}, {"failed_tests", "", "failed tests (0–100000)"}, {"progress", "", "running progress (0–100)"}, {"animate", "true", "true or false; use --animate=false"}, {"size", "300", "width and height in pixels (100–1000)"}, {"message", "", "custom caption (at most 80 characters)"},
	} {
		value := &singleFlag{value: f.def}
		values[f.name] = value
		fs.Var(value, f.name, f.help)
	}
	output := &singleFlag{value: "mascot.svg"}
	fs.Var(output, "output", "output file, or - for stdout")
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}
	if fs.NArg() != 0 {
		return fmt.Errorf("render accepts flags only")
	}
	q := url.Values{}
	fs.Visit(func(f *flag.Flag) {
		if f.Name != "output" {
			q.Set(f.Name, values[f.Name].value)
		}
	})
	c, err := mascot.Parse(q)
	if err != nil {
		return err
	}
	svg, err := mascot.Render(c)
	if err != nil {
		return err
	}
	if output.value == "-" {
		_, err = stdout.Write(svg)
		return err
	}
	if output.value == "" {
		return fmt.Errorf("output path must not be empty")
	}
	if err := os.WriteFile(output.value, svg, 0644); err != nil {
		return fmt.Errorf("write output: %w", err)
	}
	fmt.Fprintf(stdout, "Wrote %s\n", output.value)
	return nil
}

// singleFlag makes CLI duplicate handling identical to the HTTP API.
type singleFlag struct {
	value string
	seen  bool
}

func (f *singleFlag) String() string { return f.value }
func (f *singleFlag) Set(value string) error {
	if f.seen {
		return fmt.Errorf("option must occur exactly once")
	}
	f.value = value
	f.seen = true
	return nil
}
