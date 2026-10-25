package config

import "testing"

func TestLoad(t *testing.T) {
	for _, tt := range []struct {
		host, port, want string
		bad              bool
	}{
		{"", "", "0.0.0.0:8080", false}, {"::1", "9090", "[::1]:9090", false}, {"localhost", "80", "localhost:80", false},
		{"example.com", "80", "", true}, {"127.0.0.1", "0", "", true}, {"", "65536", "", true}, {"", "abc", "", true},
	} {
		t.Run(tt.host+tt.port, func(t *testing.T) {
			c, err := Load(func(k string) string {
				if k == "HOST" {
					return tt.host
				}
				return tt.port
			})
			if (err != nil) != tt.bad || (!tt.bad && c.Address != tt.want) {
				t.Fatalf("got %+v, %v", c, err)
			}
		})
	}
}
