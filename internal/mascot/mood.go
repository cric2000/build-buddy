package mascot

import "fmt"

type Resolved struct {
	Config
	Mood            string `json:"mood"`
	Animation       string `json:"animation"`
	Celebration     bool   `json:"celebration"`
	CoverageWarning bool   `json:"coverage_warning"`
	Worried         bool   `json:"worried"`
}

// Resolve preserves the explicit build status. Optional metrics only add detail.
func Resolve(c Config) (Resolved, error) {
	if err := c.Validate(); err != nil {
		return Resolved{}, err
	}
	r := Resolved{Config: c}
	switch c.Status {
	case "passing":
		r.Mood = "happy"
		r.Animation = "bounce"
		r.Message = "All systems smiling"
	case "failing":
		r.Mood = "stressed"
		r.Animation = "shake"
		r.Message = "A little help, please"
	case "running":
		r.Mood = "focused"
		r.Animation = "work"
		r.Message = "Making things happen"
	case "queued":
		r.Mood = "sleepy"
		r.Animation = "float"
		r.Message = "Dreaming of the next build"
	default:
		r.Mood = "curious"
		r.Animation = "blink"
		r.Message = "Wonder what comes next?"
	}
	if !c.Animate {
		r.Animation = "none"
	}
	if c.Status == "passing" && c.Coverage != nil {
		r.Celebration = *c.Coverage >= 90
		r.CoverageWarning = *c.Coverage < 50
	}
	if c.Status == "failing" && c.FailedTests != nil {
		r.Message = fmt.Sprintf("%d failing tests", *c.FailedTests)
		if *c.FailedTests == 1 {
			r.Message = "1 failing test"
		}
		r.Worried = *c.FailedTests >= 10
	}
	if c.Status == "running" && c.Progress != nil {
		r.Message = fmt.Sprintf("Build in progress · %d%%", *c.Progress)
	}
	if c.Message != "" {
		r.Message = c.Message
	}
	return r, nil
}
