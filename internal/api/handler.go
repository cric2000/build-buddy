// Package api exposes BuildBuddy's read-only HTTP API.
package api

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/cric2000/build-buddy/internal/mascot"
)

type Handler struct {
	logger *slog.Logger
	slots  chan struct{}
}

func New(logger *slog.Logger) http.Handler {
	if logger == nil {
		logger = slog.Default()
	}
	return &Handler{logger: logger, slots: make(chan struct{}, 64)}
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (w *responseWriter) WriteHeader(status int) {
	if w.status != 0 {
		return
	}
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}
func (w *responseWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.WriteHeader(200)
	}
	return w.ResponseWriter.Write(b)
}
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	started := time.Now()
	rw := &responseWriter{ResponseWriter: w}
	route := "unmatched"
	switch r.URL.Path {
	case "/healthz", "/v1/mascot.svg", "/v1/mascot.json", "/v1/characters", "/v1/themes":
		route = r.URL.Path
	}
	defer func() {
		if recovered := recover(); recovered != nil {
			h.logger.Error("request panic", "route", route)
			if rw.status == 0 {
				fail(rw, r, 500, "INTERNAL_ERROR", "internal server error")
			}
		}
		h.logger.Info("request", "method", r.Method, "route", route, "status", rw.status, "duration", time.Since(started).String())
	}()
	w = rw
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Referrer-Policy", "no-referrer")
	w.Header().Set("Content-Security-Policy", "default-src 'none'; style-src 'unsafe-inline'; sandbox")
	if route == "unmatched" {
		fail(w, r, 404, "NOT_FOUND", "route not found")
		return
	}
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		w.Header().Set("Allow", "GET, HEAD")
		fail(w, r, 405, "METHOD_NOT_ALLOWED", "use GET or HEAD")
		return
	}
	if len(r.URL.RawQuery) > 4096 {
		fail(w, r, 414, "URI_TOO_LONG", "query string exceeds 4096 bytes")
		return
	}
	if r.ContentLength > 0 || len(r.TransferEncoding) > 0 {
		fail(w, r, 400, "INVALID_REQUEST", "request body is not supported")
		return
	}
	select {
	case h.slots <- struct{}{}:
		defer func() { <-h.slots }()
	default:
		w.Header().Set("Retry-After", "1")
		fail(w, r, 503, "BUSY", "server is busy; retry shortly")
		return
	}
	switch route {
	case "/healthz":
		respond(w, r, 200, "application/json; charset=utf-8", []byte("{\"status\":\"ok\"}\n"), "no-store")
	case "/v1/characters":
		h.json(w, r, mascot.Characters(), "public, max-age=3600")
	case "/v1/themes":
		h.json(w, r, mascot.Themes(), "public, max-age=3600")
	default:
		q, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			fail(w, r, 400, "INVALID_PARAMETER", "malformed query encoding")
			return
		}
		c, err := mascot.Parse(q)
		if err != nil {
			fail(w, r, 400, "INVALID_PARAMETER", err.Error())
			return
		}
		if route == "/v1/mascot.json" {
			resolved, err := mascot.Resolve(c)
			if err != nil {
				h.internal(w, r, err)
				return
			}
			h.json(w, r, resolved, "public, max-age=60, must-revalidate")
			return
		}
		svg, err := mascot.Render(c)
		if err != nil {
			h.internal(w, r, err)
			return
		}
		respond(w, r, 200, "image/svg+xml; charset=utf-8", svg, "public, max-age=60, must-revalidate")
	}
}
func (h *Handler) json(w http.ResponseWriter, r *http.Request, v any, cache string) {
	b, err := json.Marshal(v)
	if err != nil {
		h.internal(w, r, err)
		return
	}
	respond(w, r, 200, "application/json; charset=utf-8", append(b, '\n'), cache)
}
func (h *Handler) internal(w http.ResponseWriter, r *http.Request, err error) {
	h.logger.Error("response failed", "error", err)
	fail(w, r, 500, "INTERNAL_ERROR", "internal server error")
}
func fail(w http.ResponseWriter, r *http.Request, status int, code, message string) {
	body := struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}{}
	body.Error.Code = code
	body.Error.Message = message
	b, _ := json.Marshal(body)
	respond(w, r, status, "application/json; charset=utf-8", append(b, '\n'), "no-store")
}
func respond(w http.ResponseWriter, r *http.Request, status int, mime string, body []byte, cache string) {
	w.Header().Set("Content-Type", mime)
	w.Header().Set("Cache-Control", cache)
	if status == 200 && cache != "no-store" {
		etag := fmt.Sprintf(`"%x"`, sha256.Sum256(body))
		w.Header().Set("ETag", etag)
		for _, candidate := range strings.Split(strings.Join(r.Header.Values("If-None-Match"), ","), ",") {
			candidate = strings.TrimSpace(candidate)
			if candidate == "*" || strings.TrimPrefix(candidate, "W/") == etag {
				w.WriteHeader(304)
				return
			}
		}
	}
	w.Header().Set("Content-Length", fmt.Sprint(len(body)))
	w.WriteHeader(status)
	if r.Method != http.MethodHead {
		_, _ = w.Write(body)
	}
}
