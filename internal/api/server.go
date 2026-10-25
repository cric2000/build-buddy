package api

import (
	"log/slog"
	"net/http"
	"time"
)

func Server(address string, logger *slog.Logger) *http.Server {
	if logger == nil {
		logger = slog.Default()
	}
	return &http.Server{
		Addr: address, Handler: New(logger), ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout: 10 * time.Second, WriteTimeout: 15 * time.Second, IdleTimeout: 60 * time.Second,
		MaxHeaderBytes: 16 << 10, ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}
}
