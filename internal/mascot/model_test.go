package mascot

import (
	"net/url"
	"slices"
	"strings"
	"testing"
)

func TestCharacterCatalog(t *testing.T) {
	want := []string{"kappa", "tanuki", "tengu", "kurama", "chopper"}
	if !slices.Equal(CharacterNames(), want) {
		t.Fatalf("unexpected catalog: %v", CharacterNames())
	}
	if Default().Character != "kappa" {
		t.Fatal("Kappa should be the default companion")
	}
}

func TestResolve(t *testing.T) {
	for _, tt := range []struct {
		query, mood                   string
		celebration, warning, worried bool
	}{
		{"status=passing", "happy", false, false, false}, {"status=passing&coverage=90", "happy", true, false, false}, {"status=passing&coverage=49", "happy", false, true, false},
		{"status=failing&failed_tests=12", "stressed", false, false, true}, {"status=running&progress=65", "focused", false, false, false}, {"status=queued", "sleepy", false, false, false}, {"", "curious", false, false, false},
	} {
		t.Run(tt.query, func(t *testing.T) {
			q, _ := url.ParseQuery(tt.query)
			c, err := Parse(q)
			if err != nil {
				t.Fatal(err)
			}
			r, err := Resolve(c)
			if err != nil || r.Mood != tt.mood || r.Celebration != tt.celebration || r.CoverageWarning != tt.warning || r.Worried != tt.worried {
				t.Fatalf("%+v %v", r, err)
			}
		})
	}
	c := Default()
	c.Message = "Custom"
	c.Animate = false
	r, _ := Resolve(c)
	if r.Message != "Custom" || r.Animation != "none" {
		t.Fatal(r)
	}
}
func TestInvalid(t *testing.T) {
	for _, raw := range []string{"status=bad", "character=dog", "theme=blue", "coverage=-1", "coverage=101", "coverage=1.5", "coverage=", "coverage=1&coverage=2", "size=99", "size=1001", "failed_tests=100001", "status=passing&failed_tests=4", "status=failing&failed_tests=0", "progress=30", "status=running&progress=101", "animate=1", "status=", "wat=1", "message=%00", "message=%FF", "message=" + strings.Repeat("a", 81)} {
		t.Run(raw, func(t *testing.T) {
			q, err := url.ParseQuery(raw)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := Parse(q); err == nil {
				t.Fatal("accepted invalid input")
			}
		})
	}
}
func TestBoundaries(t *testing.T) {
	for _, raw := range []string{"size=100&coverage=0", "size=1000&coverage=100", "status=failing&failed_tests=100000", "status=running&progress=0", "status=running&progress=100", "animate=false", "message=" + strings.Repeat("森", 80)} {
		q, _ := url.ParseQuery(raw)
		if _, err := Parse(q); err != nil {
			t.Fatal(raw, err)
		}
	}
}
