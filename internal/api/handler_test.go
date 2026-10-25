package api

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func quiet() *slog.Logger { return slog.New(slog.NewTextHandler(io.Discard, nil)) }
func request(h http.Handler, method, path string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest(method, path, nil))
	return w
}
func TestEndpoints(t *testing.T) {
	h := New(quiet())
	for _, tt := range []struct{ path, mime, contains string }{
		{"/healthz", "application/json", `"status":"ok"`}, {"/v1/characters", "application/json", `"cat"`}, {"/v1/themes", "application/json", `"neon"`},
		{"/v1/mascot.svg?status=passing", "image/svg+xml", `data-mood="happy"`}, {"/v1/mascot.json?status=failing&failed_tests=4", "application/json", `"message":"4 failing tests"`},
	} {
		t.Run(tt.path, func(t *testing.T) {
			w := request(h, "GET", tt.path)
			if w.Code != 200 || !strings.HasPrefix(w.Header().Get("Content-Type"), tt.mime) || !strings.Contains(w.Body.String(), tt.contains) {
				t.Fatalf("%d %v %s", w.Code, w.Header(), w.Body)
			}
			if strings.HasPrefix(tt.mime, "application/json") && !json.Valid(w.Body.Bytes()) {
				t.Fatal("invalid JSON")
			}
			if w.Header().Get("X-Content-Type-Options") != "nosniff" {
				t.Fatal("missing nosniff")
			}
			head := request(h, "HEAD", tt.path)
			if head.Code != 200 || head.Body.Len() != 0 || head.Header().Get("Content-Length") != w.Header().Get("Content-Length") {
				t.Fatal("bad HEAD response")
			}
		})
	}
	w := request(h, "GET", "/v1/mascot.json")
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	for _, k := range []string{"coverage", "failed_tests", "progress"} {
		if v, ok := body[k]; !ok || v != nil {
			t.Fatal("expected consistent null field", k)
		}
	}
}
func TestErrors(t *testing.T) {
	h := New(quiet())
	for _, tt := range []struct {
		method, path string
		status       int
		code         string
	}{
		{"GET", "/", 404, "NOT_FOUND"}, {"POST", "/v1/mascot.svg", 405, "METHOD_NOT_ALLOWED"}, {"OPTIONS", "/healthz", 405, "METHOD_NOT_ALLOWED"},
		{"GET", "/v1/mascot.svg?status=x", 400, "INVALID_PARAMETER"}, {"GET", "/v1/mascot.json?coverage=no", 400, "INVALID_PARAMETER"},
		{"GET", "/v1/mascot.svg?status=passing&status=failing", 400, "INVALID_PARAMETER"}, {"GET", "/v1/mascot.svg?character=", 400, "INVALID_PARAMETER"},
		{"GET", "/v1/mascot.svg?theme=bad", 400, "INVALID_PARAMETER"}, {"GET", "/v1/mascot.svg?animate=TRUE", 400, "INVALID_PARAMETER"},
		{"GET", "/v1/mascot.svg?status=passing&failed_tests=4", 400, "INVALID_PARAMETER"}, {"GET", "/v1/mascot.svg?message=%00", 400, "INVALID_PARAMETER"},
		{"GET", "/v1/mascot.svg?message=%FF", 400, "INVALID_PARAMETER"}, {"GET", "/v1/mascot.svg?message=%ZZ", 400, "INVALID_PARAMETER"},
		{"GET", "/v1/mascot.svg?message=a;b", 400, "INVALID_PARAMETER"}, {"GET", "/v1/mascot.svg?x=" + strings.Repeat("a", 4097), 414, "URI_TOO_LONG"},
	} {
		t.Run(tt.path, func(t *testing.T) {
			w := request(h, tt.method, tt.path)
			if w.Code != tt.status || !strings.Contains(w.Body.String(), `"code":"`+tt.code+`"`) || !json.Valid(w.Body.Bytes()) {
				t.Fatalf("%d %s", w.Code, w.Body)
			}
			if w.Header().Get("Cache-Control") != "no-store" {
				t.Fatal("error cached")
			}
			if tt.status == 405 && w.Header().Get("Allow") != "GET, HEAD" {
				t.Fatal("missing Allow")
			}
		})
	}
	r := httptest.NewRequest("GET", "/healthz", strings.NewReader("body"))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != 400 {
		t.Fatal("accepted body")
	}
}
func TestCaching(t *testing.T) {
	h := New(quiet())
	path := "/v1/mascot.svg?status=passing&character=cat"
	first := request(h, "GET", path)
	etag := first.Header().Get("ETag")
	if etag != fmt.Sprintf(`"%x"`, sha256.Sum256(first.Body.Bytes())) || !strings.Contains(first.Header().Get("Cache-Control"), "max-age=60") {
		t.Fatal("bad cache headers")
	}
	for _, header := range []string{etag, "W/" + etag, `"other", W/` + etag, "*"} {
		r := httptest.NewRequest("GET", path, nil)
		r.Header.Set("If-None-Match", header)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		if w.Code != 304 || w.Body.Len() != 0 || w.Header().Get("ETag") != etag {
			t.Fatal("conditional request", header, w.Code)
		}
	}
	r := httptest.NewRequest("GET", path, nil)
	r.Header.Set("If-None-Match", `"stale"`)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != 200 {
		t.Fatal("stale returned 304")
	}
	normalized := request(h, "GET", "/v1/mascot.svg?character=cat&status=passing&size=0300&animate=true&theme=light")
	if !bytes.Equal(first.Body.Bytes(), normalized.Body.Bytes()) {
		t.Fatal("normalization mismatch")
	}
	if request(h, "GET", "/healthz").Header().Get("ETag") != "" {
		t.Fatal("health cached")
	}
	if request(h, "GET", "/v1/themes").Header().Get("Cache-Control") != "public, max-age=3600" {
		t.Fatal("metadata cache")
	}
}
func TestSVGAndJSONAgree(t *testing.T) {
	h := New(quiet())
	q := "?status=running&progress=65&character=ghost&animate=false"
	j := request(h, "GET", "/v1/mascot.json"+q)
	var resolved struct{ Mood, Message, Animation string }
	if err := json.Unmarshal(j.Body.Bytes(), &resolved); err != nil {
		t.Fatal(err)
	}
	w := request(h, "GET", "/v1/mascot.svg"+q)
	dec := xml.NewDecoder(w.Body)
	var doc struct {
		Mood string `xml:"data-mood,attr"`
		Desc string `xml:"desc"`
	}
	if err := dec.Decode(&doc); err != nil {
		t.Fatal(err)
	}
	if doc.Mood != resolved.Mood || !strings.Contains(doc.Desc, resolved.Message) || resolved.Animation != "none" {
		t.Fatal("inconsistent output")
	}
}
func TestLoggingDoesNotLeakInput(t *testing.T) {
	var logs bytes.Buffer
	h := New(slog.New(slog.NewJSONHandler(&logs, nil)))
	request(h, "GET", "/v1/mascot.svg?message=private-secret")
	request(h, "GET", "/private-secret")
	if strings.Contains(logs.String(), "private-secret") {
		t.Fatal("input leaked to logs")
	}
	for _, field := range []string{`"method"`, `"route"`, `"status"`, `"duration"`} {
		if !strings.Contains(logs.String(), field) {
			t.Fatal("missing log field", field)
		}
	}
}
func TestBusyAndTimeouts(t *testing.T) {
	h := New(quiet()).(*Handler)
	for i := 0; i < cap(h.slots); i++ {
		h.slots <- struct{}{}
	}
	w := request(h, "GET", "/v1/mascot.svg")
	if w.Code != 503 || w.Header().Get("Retry-After") != "1" {
		t.Fatal("unbounded concurrency")
	}
	s := Server("127.0.0.1:8080", quiet())
	if s.ReadHeaderTimeout <= 0 || s.ReadTimeout <= 0 || s.WriteTimeout <= 0 || s.IdleTimeout <= 0 || s.MaxHeaderBytes <= 0 {
		t.Fatal("missing server limits")
	}
}

func TestInternalErrorsAreGeneric(t *testing.T) {
	var logs bytes.Buffer
	h := New(slog.New(slog.NewJSONHandler(&logs, nil))).(*Handler)
	r := httptest.NewRequest("GET", "/v1/mascot.json", nil)
	w := httptest.NewRecorder()
	h.json(w, r, make(chan int), "public, max-age=60")
	if w.Code != 500 || !strings.Contains(w.Body.String(), `"code":"INTERNAL_ERROR"`) || strings.Contains(w.Body.String(), "chan") || !json.Valid(w.Body.Bytes()) {
		t.Fatalf("unsafe internal error: %d %s", w.Code, w.Body)
	}
	if !strings.Contains(logs.String(), "response failed") {
		t.Fatal("internal failure not logged")
	}
}
