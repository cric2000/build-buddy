package app

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"log/slog"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"buildbuddy/internal/api"
	"buildbuddy/internal/mascot"
)

func TestAPICLIRenderParity(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	handler := api.New(logger)
	output := filepath.Join(t.TempDir(), "mascot.svg")
	for _, character := range mascot.Characters() {
		for _, status := range []string{"passing", "failing", "running", "queued", "unknown"} {
			for _, theme := range mascot.Themes() {
				for _, animate := range []string{"true", "false"} {
					t.Run(fmt.Sprintf("%s/%s/%s/%s", character.Name, status, theme, animate), func(t *testing.T) {
						options := []string{"--character", character.Name, "--status", status, "--theme", theme, "--animate", animate, "--size", "500", "--message", `Build <3 & ship safely`}
						query := url.Values{"character": {character.Name}, "status": {status}, "theme": {theme}, "animate": {animate}, "size": {"500"}, "message": {`Build <3 & ship safely`}}
						switch status {
						case "passing":
							options = append(options, "--coverage", "96")
							query.Set("coverage", "96")
						case "failing":
							options = append(options, "--failed_tests", "14")
							query.Set("failed_tests", "14")
						case "running":
							options = append(options, "--progress", "65")
							query.Set("progress", "65")
						}
						options = append(options, "--output", output)
						if err := run(context.Background(), append([]string{"render"}, options...), os.Getenv, io.Discard, io.Discard); err != nil {
							t.Fatal(err)
						}
						local, err := os.ReadFile(output)
						if err != nil {
							t.Fatal(err)
						}
						response := httptest.NewRecorder()
						handler.ServeHTTP(response, httptest.NewRequest("GET", "/v1/mascot.svg?"+query.Encode(), nil))
						if response.Code != 200 || !bytes.Equal(local, response.Body.Bytes()) {
							t.Fatalf("API and CLI differ (HTTP %d)", response.Code)
						}
					})
				}
			}
		}
	}
}

// Exercise the same scripts as Actions against a local Git remote. No GitHub token
// or network service is required, and test results come from real Go test runs.
func TestPublishFromActualTestResults(t *testing.T) {
	for _, tool := range []string{"git", "bash"} {
		if _, err := exec.LookPath(tool); err != nil {
			t.Skipf("publishing integration requires %s", tool)
		}
	}
	renderScript, err := filepath.Abs("../../examples/render-test-result.sh")
	if err != nil {
		t.Fatal(err)
	}
	publishScript, err := filepath.Abs("../../examples/publish-mascot.sh")
	if err != nil {
		t.Fatal(err)
	}
	// Track the shell scripts as Go test-cache inputs, since Bash reads them in a child process.
	for _, script := range []string{renderScript, publishScript} {
		if _, err := os.ReadFile(script); err != nil {
			t.Fatal(err)
		}
	}
	temp := t.TempDir()
	binary := filepath.Join(temp, "buildbuddy")
	remote := filepath.Join(temp, "remote.git")
	work := filepath.Join(temp, "source")
	command(t, ".", nil, "go", "build", "-o", binary, "../../cmd/buildbuddy")
	command(t, temp, nil, "git", "init", "--bare", "--initial-branch=main", remote)
	command(t, temp, nil, "git", "clone", remote, work)
	command(t, work, nil, "git", "config", "user.name", "Integration test")
	command(t, work, nil, "git", "config", "user.email", "test@example.invalid")
	writeFile(t, filepath.Join(work, "go.mod"), "module fixture\n\ngo 1.26.0\n")
	writeFile(t, filepath.Join(work, "result_test.go"), `package fixture
import ("os"; "testing")
func TestBuild(t *testing.T) {
 if os.Getenv("BUILDBUDDY_FIXTURE_FAIL") == "1" { t.Fatal("intentional test failure") }
}
`)
	command(t, work, nil, "git", "add", "go.mod", "result_test.go")
	command(t, work, nil, "git", "commit", "-m", "Initial test project")
	command(t, work, nil, "git", "push", "origin", "main")
	initial := command(t, work, nil, "git", "rev-parse", "HEAD")
	var previous string
	for i, status := range []string{"passing", "failing", "failing", "passing"} {
		t.Run(fmt.Sprintf("%d-%s", i, status), func(t *testing.T) {
			svg := filepath.Join(temp, fmt.Sprintf("mascot-%d.svg", i))
			fail := "0"
			if status == "failing" {
				fail = "1"
			}
			// -count=1 ensures this runs tests, even when their source has not changed.
			testCommand := exec.Command("bash", renderScript, binary, svg, "go", "test", "-count=1", "./...")
			testCommand.Dir = work
			testCommand.Env = append(os.Environ(), "GOWORK=off", "BUILDBUDDY_FIXTURE_FAIL="+fail)
			log, err := testCommand.CombinedOutput()
			if status == "passing" && err != nil {
				t.Fatalf("passing tests failed: %v\n%s", err, log)
			}
			if status == "failing" {
				exit, ok := err.(*exec.ExitError)
				if !ok || exit.ExitCode() != 1 {
					t.Fatalf("failure exit status lost: %v\n%s", err, log)
				}
			}
			command(t, work, nil, "git", "fetch", "origin")
			command(t, work, []string{"SVG_PATH=" + svg, "BADGE_REPO=" + filepath.Join(temp, fmt.Sprintf("publish-%d", i))}, "bash", publishScript)
			reader := filepath.Join(temp, fmt.Sprintf("reader-%d", i))
			command(t, temp, nil, "git", "clone", "--branch", "badge-artifact", remote, reader)
			published, err := os.ReadFile(filepath.Join(reader, "mascot.svg"))
			if err != nil {
				t.Fatal(err)
			}
			local, err := os.ReadFile(svg)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(local, published) {
				t.Fatal("published image differs from generated image")
			}
			var doc struct {
				Title string `xml:"title"`
				Mood  string `xml:"data-mood,attr"`
			}
			if err := xml.Unmarshal(published, &doc); err != nil {
				t.Fatal(err)
			}
			mood := "happy"
			if status == "failing" {
				mood = "stressed"
			}
			if doc.Title != "BuildBuddy: "+status+" build" || doc.Mood != mood {
				t.Fatalf("published wrong build result: %+v", doc)
			}
			if files := command(t, reader, nil, "git", "ls-tree", "--name-only", "HEAD"); files != "mascot.svg" {
				t.Fatalf("badge branch includes source files: %s", files)
			}
			head := command(t, reader, nil, "git", "rev-parse", "HEAD")
			if i == 2 && head != previous {
				t.Fatal("unchanged artwork created a commit")
			}
			if i != 2 && head == previous {
				t.Fatal("build status transition was not published")
			}
			previous = head
			if command(t, work, nil, "git", "config", "user.name") != "Integration test" {
				t.Fatal("publication changed the development Git identity")
			}
			if command(t, work, nil, "git", "rev-parse", "HEAD") != initial || command(t, work, nil, "git", "status", "--porcelain") != "" {
				t.Fatal("publication changed the development checkout")
			}
			if command(t, work, nil, "git", "ls-remote", "origin", "refs/heads/main") != initial+"\trefs/heads/main" {
				t.Fatal("publication changed the remote development branch")
			}
		})
	}
}

func command(t *testing.T, dir string, env []string, args ...string) string {
	t.Helper()
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null")
	cmd.Env = append(cmd.Env, env...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%s: %v\n%s", strings.Join(args, " "), err, output)
	}
	return strings.TrimSpace(string(output))
}
func writeFile(t *testing.T, path, contents string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(contents), 0600); err != nil {
		t.Fatal(err)
	}
}
