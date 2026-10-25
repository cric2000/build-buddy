// Package config validates the HTTP server's environment configuration.
package config

import (
	"fmt"
	"net"
	"strconv"
)

type Config struct{ Address string }

func Load(getenv func(string) string) (Config, error) {
	host, port := getenv("HOST"), getenv("PORT")
	if host == "" {
		host = "0.0.0.0"
	}
	if port == "" {
		port = "8080"
	}
	n, err := strconv.Atoi(port)
	if err != nil || n < 1 || n > 65535 {
		return Config{}, fmt.Errorf("PORT must be an integer between 1 and 65535")
	}
	if host != "localhost" && net.ParseIP(host) == nil {
		return Config{}, fmt.Errorf("HOST must be an IP address or localhost")
	}
	return Config{Address: net.JoinHostPort(host, strconv.Itoa(n))}, nil
}
