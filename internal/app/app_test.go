package app

import (
	"bytes"
	"context"
	"encoding/xml"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/cric2000/build-buddy/internal/api"
)

func TestRenderCommand(t *testing.T) {
	path := filepath.Join(t.TempDir(), "mascot.svg")
	var out, errs bytes.Buffer
	err := run(context.Background(), []string{"render", "--status", "passing", "--character", "cat", "--theme", "dark", "--animate=false", "--output", path}, os.Getenv, &out, &errs)
	if err != nil {
		t.Fatal(err, errs.String())
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var doc struct {
		XMLName   xml.Name
		Character string `xml:"data-character,attr"`
	}
	if err := xml.Unmarshal(b, &doc); err != nil || doc.Character != "cat" || bytes.Contains(b, []byte("<style>")) {
		t.Fatalf("invalid SVG %v", err)
	}
	out.Reset()
	if err := run(context.Background(), []string{"render", "--output", "-"}, os.Getenv, &out, &errs); err != nil || !strings.HasPrefix(out.String(), "<svg") {
		t.Fatal("stdout render failed", err)
	}
}
func TestCommandErrors(t *testing.T) {
	for _, args := range [][]string{
		{"nonsense"}, {"render", "--status", "bad"}, {"render", "--coverage", "no"}, {"render", "--status", "passing", "--status", "failing"},
		{"render", "--output", "a", "--output", "b"}, {"render", "--animate", "yes"}, {"render", "--output", ""}, {"render", "--unknown", "x"}, {"render", "extra"},
		{"render", "--output", filepath.Join(t.TempDir(), "absent", "mascot.svg")}, {"render", "--output", t.TempDir()},
		{"serve", "extra"}, {"serve", "--wat"},
	} {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			if err := run(context.Background(), args, os.Getenv, io.Discard, io.Discard); err == nil {
				t.Fatal("accepted invalid command")
			}
		})
	}
}
func TestHelp(t *testing.T) {
	for _, args := range [][]string{{"--help"}, {"help"}, {"render", "--help"}, {"serve", "-h"}} {
		var out bytes.Buffer
		err := run(context.Background(), args, os.Getenv, &out, &out)
		if err != nil || out.Len() == 0 {
			t.Fatal("help", err)
		}
	}
}
func TestServerConfiguration(t *testing.T) {
	for _, args := range [][]string{nil, {"serve"}} {
		err := run(context.Background(), args, func(k string) string {
			if k == "PORT" {
				return "invalid"
			}
			return ""
		}, io.Discard, io.Discard)
		if err == nil || !strings.Contains(err.Error(), "PORT") {
			t.Fatal("bad environment accepted", err)
		}
	}
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	_, port, _ := net.SplitHostPort(listener.Addr().String())
	err = run(context.Background(), []string{"serve"}, func(k string) string {
		if k == "HOST" {
			return "127.0.0.1"
		}
		return port
	}, io.Discard, io.Discard)
	if err == nil || !strings.Contains(err.Error(), "listen") {
		t.Fatal("bind failure lost", err)
	}
}
func TestServeAndGracefulShutdown(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	server := api.Server(listener.Addr().String(), logger)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() { done <- serve(ctx, server, listener, logger) }()
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get("http://" + listener.Addr().String() + "/healthz")
	if err != nil {
		t.Fatal(err)
	}
	b, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil || resp.StatusCode != 200 || !bytes.Contains(b, []byte(`"ok"`)) {
		t.Fatal("health request failed", err)
	}
	cancel()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("shutdown hung")
	}
	if _, err := client.Get("http://" + listener.Addr().String() + "/healthz"); err == nil {
		t.Fatal("server still listening")
	}
}
