// Package mascot validates configuration and renders original, deterministic SVG art.
package mascot

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"unicode/utf8"
)

type Character struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func Characters() []Character {
	return []Character{
		{"kappa", "A river tinkerer with a water dish, a patched shell, and a cucumber pencil."},
		{"tanuki", "A leaf-capped woodland courier with an oversized tail and a lucky satchel."},
		{"tengu", "A mountain lookout with a long red nose, feathered wings, and a maple fan."},
		{"robot", "An energetic little tinkerer with a heart of gold."},
		{"cat", "A playful companion with big feelings and a restless tail."},
		{"ghost", "A sleepy little spirit cheering on your next build."},
	}
}

func CharacterNames() []string {
	characters := Characters()
	names := make([]string, len(characters))
	for i, character := range characters {
		names[i] = character.Name
	}
	return names
}
func Themes() []string { return []string{"light", "dark", "neon"} }

type Config struct {
	Character   string `json:"character"`
	Status      string `json:"status"`
	Theme       string `json:"theme"`
	Coverage    *int   `json:"coverage"`
	FailedTests *int   `json:"failed_tests"`
	Progress    *int   `json:"progress"`
	Animate     bool   `json:"animate"`
	Size        int    `json:"size"`
	Message     string `json:"message"`
}

func Default() Config {
	return Config{Character: "robot", Status: "unknown", Theme: "light", Animate: true, Size: 300}
}

// Parse rejects unknown, repeated, empty, malformed, and contradictory parameters.
func Parse(q url.Values) (Config, error) {
	c := Default()
	for key, values := range q {
		if len(values) != 1 {
			return Config{}, fmt.Errorf("%s must occur exactly once", key)
		}
		v := values[0]
		switch key {
		case "character":
			c.Character = v
		case "status":
			c.Status = v
		case "theme":
			c.Theme = v
		case "message":
			c.Message = v
		case "animate":
			if v != "true" && v != "false" {
				return Config{}, fmt.Errorf("animate must be true or false")
			}
			c.Animate = v == "true"
		case "coverage", "failed_tests", "progress", "size":
			if v == "" || strings.Trim(v, "0123456789") != "" {
				return Config{}, fmt.Errorf("%s must be an unsigned decimal integer", key)
			}
			n, err := strconv.Atoi(v)
			if err != nil {
				return Config{}, fmt.Errorf("%s is out of range", key)
			}
			switch key {
			case "coverage":
				c.Coverage = &n
			case "failed_tests":
				c.FailedTests = &n
			case "progress":
				c.Progress = &n
			case "size":
				c.Size = n
			}
		default:
			return Config{}, fmt.Errorf("unsupported parameter: %s", key)
		}
	}
	if err := c.Validate(); err != nil {
		return Config{}, err
	}
	return c, nil
}

func (c Config) Validate() error {
	if !oneOf(c.Character, CharacterNames()...) {
		return fmt.Errorf("character must be one of: %s", strings.Join(CharacterNames(), ", "))
	}
	if !oneOf(c.Status, "passing", "failing", "running", "queued", "unknown") {
		return fmt.Errorf("status must be passing, failing, running, queued, or unknown")
	}
	if !oneOf(c.Theme, "light", "dark", "neon") {
		return fmt.Errorf("theme must be light, dark, or neon")
	}
	for _, field := range []struct {
		name  string
		value *int
		max   int
	}{{"coverage", c.Coverage, 100}, {"failed_tests", c.FailedTests, 100000}, {"progress", c.Progress, 100}} {
		if field.value != nil && (*field.value < 0 || *field.value > field.max) {
			return fmt.Errorf("%s must be between 0 and %d", field.name, field.max)
		}
	}
	if c.Size < 100 || c.Size > 1000 {
		return fmt.Errorf("size must be between 100 and 1000")
	}
	if !utf8.ValidString(c.Message) || utf8.RuneCountInString(c.Message) > 80 {
		return fmt.Errorf("message must be valid UTF-8 and at most 80 characters")
	}
	for _, r := range c.Message {
		if r < 0x20 || r == 0x7f || r == 0xfffe || r == 0xffff {
			return fmt.Errorf("message must contain printable XML characters")
		}
	}
	if c.FailedTests != nil && *c.FailedTests > 0 && c.Status != "failing" {
		return fmt.Errorf("positive failed_tests requires status=failing")
	}
	if c.FailedTests != nil && *c.FailedTests == 0 && c.Status == "failing" {
		return fmt.Errorf("status=failing contradicts failed_tests=0; omit the count for non-test failures")
	}
	if c.Progress != nil && c.Status != "running" {
		return fmt.Errorf("progress requires status=running")
	}
	return nil
}
func oneOf(v string, values ...string) bool {
	for _, s := range values {
		if v == s {
			return true
		}
	}
	return false
}
