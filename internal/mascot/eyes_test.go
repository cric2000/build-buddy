package mascot

import (
	"bytes"
	"encoding/xml"
	"io"
	"strings"
	"testing"
)

// Compare only the eye artwork: poses, mouths, labels, and expression metadata
// cannot make identical eyes pass this check.
func TestStaticEyeExpressions(t *testing.T) {
	for _, character := range Characters() {
		for _, theme := range Themes() {
			t.Run(character.Name+"/"+theme, func(t *testing.T) {
				seen := map[string]string{}
				for _, tt := range []struct {
					status, expression string
					failures           int
				}{
					{"passing", "sparkly", 0}, {"failing", "sad", 0},
					{"failing", "disappointed", 4}, {"failing", "teary", 12},
					{"running", "focused", 0}, {"queued", "sleepy", 0}, {"unknown", "curious", 0},
				} {
					c := Default()
					c.Character, c.Theme, c.Status, c.Animate = character.Name, theme, tt.status, false
					if tt.failures > 0 {
						c.FailedTests = &tt.failures
					}
					svg, err := Render(c)
					if err != nil {
						t.Fatal(err)
					}
					expression, artwork := renderedEyes(t, svg)
					if expression != tt.expression {
						t.Fatalf("%s: got %s", tt.expression, expression)
					}
					if previous, ok := seen[artwork]; ok {
						t.Fatalf("%s and %s have identical eye artwork", previous, expression)
					}
					seen[artwork] = expression
					if strings.Count(artwork, `data-eye=`) != 2 {
						t.Fatal("expected two drawn eyes")
					}
					if expression == "sparkly" && strings.Count(artwork, `class="eye-glint"`) != 2 {
						t.Fatal("sparkly eyes need visible catchlights without animation")
					}
					if expression == "sleepy" && strings.Contains(artwork, "<ellipse") {
						t.Fatal("sleepy eyes should have closed lids")
					}
					if expression == "teary" && strings.Count(artwork, `data-part="tears"`) != 2 {
						t.Fatal("tears must come from both eyes")
					}
				}
			})
		}
	}
}

func renderedEyes(t *testing.T, svg []byte) (string, string) {
	t.Helper()
	decoder := xml.NewDecoder(bytes.NewReader(svg))
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			t.Fatal("missing eye artwork")
		}
		if err != nil {
			t.Fatal(err)
		}
		start, ok := token.(xml.StartElement)
		if !ok || start.Name.Local != "g" {
			continue
		}
		for _, attribute := range start.Attr {
			if attribute.Name.Local == "data-part" && attribute.Value == "expressive-eyes" {
				var group struct {
					Expression string `xml:"data-expression,attr"`
					Artwork    string `xml:",innerxml"`
				}
				if err := decoder.DecodeElement(&group, &start); err != nil {
					t.Fatal(err)
				}
				return group.Expression, group.Artwork
			}
		}
	}
}
