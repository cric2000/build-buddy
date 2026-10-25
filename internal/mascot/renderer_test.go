package mascot

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
	"sync"
	"testing"
)

func validSVG(t *testing.T, b []byte) {
	t.Helper()
	dec := xml.NewDecoder(bytes.NewReader(b))
	root := false
	ids := map[string]bool{}
	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("invalid XML: %v", err)
		}
		e, ok := tok.(xml.StartElement)
		if !ok {
			continue
		}
		if !root {
			if e.Name.Local != "svg" || e.Name.Space != "http://www.w3.org/2000/svg" {
				t.Fatal("not SVG")
			}
			root = true
		}
		if e.Name.Local == "script" || e.Name.Local == "foreignObject" || e.Name.Local == "image" {
			t.Fatalf("unsafe element: %s", e.Name.Local)
		}
		for _, a := range e.Attr {
			if strings.HasPrefix(strings.ToLower(a.Name.Local), "on") || a.Name.Local == "href" {
				t.Fatalf("unsafe attribute %s", a.Name.Local)
			}
			if a.Name.Local == "id" {
				if ids[a.Value] {
					t.Fatal("duplicate ID")
				}
				ids[a.Value] = true
			}
		}
	}
	if !root {
		t.Fatal("empty SVG")
	}
}
func TestRenderMatrix(t *testing.T) {
	for _, character := range Characters() {
		for _, status := range []string{"passing", "failing", "running", "queued", "unknown"} {
			for _, theme := range Themes() {
				for _, animate := range []bool{true, false} {
					t.Run(fmt.Sprintf("%s/%s/%s/%t", character.Name, status, theme, animate), func(t *testing.T) {
						c := Default()
						c.Character = character.Name
						c.Status = status
						c.Theme = theme
						c.Animate = animate
						b, err := Render(c)
						if err != nil {
							t.Fatal(err)
						}
						validSVG(t, b)
						resolved, _ := Resolve(c)
						for _, s := range []string{`data-character="` + c.Character + `"`, `data-mood="` + resolved.Mood + `"`, `data-theme="` + c.Theme + `"`, colors(theme).bg, strings.ToUpper(status)} {
							if !bytes.Contains(b, []byte(s)) {
								t.Fatalf("missing %s", s)
							}
						}
						if bytes.Contains(b, []byte("<style>")) != animate {
							t.Fatal("animation mismatch")
						}
						if animate && !bytes.Contains(b, []byte("prefers-reduced-motion")) {
							t.Fatal("missing reduced motion")
						}
						if len(b) > 16000 || len(b) < 2000 {
							t.Fatalf("unexpected SVG size %d", len(b))
						}
						again, _ := Render(c)
						if !bytes.Equal(b, again) {
							t.Fatal("nondeterministic output")
						}
					})
				}
			}
		}
	}
}
func TestRenderEscapesText(t *testing.T) {
	c := Default()
	c.Message = `</text><script>alert("x")</script>& 'hello'`
	b, err := Render(c)
	if err != nil {
		t.Fatal(err)
	}
	validSVG(t, b)
	if bytes.Contains(b, []byte("<script>")) || !bytes.Contains(b, []byte("&lt;script&gt;")) {
		t.Fatal("unescaped message")
	}
	dec := xml.NewDecoder(bytes.NewReader(b))
	var texts string
	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		if text, ok := tok.(xml.CharData); ok {
			texts += string(text)
		}
	}
	if !strings.Contains(texts, c.Message) {
		t.Fatal("text did not survive XML escaping")
	}
	c.Message = string([]byte{0xff})
	if _, err := Render(c); err == nil {
		t.Fatal("invalid UTF-8 accepted")
	}
	c.Character = `robot" onload="alert(1)`
	if _, err := Render(c); err == nil {
		t.Fatal("invalid identifier accepted")
	}
}
func TestMetricsAndDimensions(t *testing.T) {
	n := 65
	c := Default()
	c.Status = "running"
	c.Progress = &n
	c.Size = 1000
	b, _ := Render(c)
	validSVG(t, b)
	for _, s := range []string{`width="1000" height="1000"`, `width="110.5"`, `65%`} {
		if !bytes.Contains(b, []byte(s)) {
			t.Fatal("missing", s)
		}
	}
	c = Default()
	c.Status = "passing"
	n = 42
	c.Coverage = &n
	b, _ = Render(c)
	if !bytes.Contains(b, []byte("ROOM TO GROW")) {
		t.Fatal("missing low coverage hint")
	}
	c.Message = strings.Repeat("界", 80)
	b, _ = Render(c)
	validSVG(t, b)
}
func TestConcurrentRender(t *testing.T) {
	c := Default()
	c.Character = "cat"
	want, _ := Render(c)
	var wg sync.WaitGroup
	for i := 0; i < 64; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				b, err := Render(c)
				if err != nil || !bytes.Equal(b, want) {
					t.Error("concurrent render differs", err)
				}
			}
		}()
	}
	wg.Wait()
}
func BenchmarkRender(b *testing.B) {
	for _, char := range Characters() {
		b.Run(char.Name, func(b *testing.B) {
			c := Default()
			c.Character = char.Name
			b.ReportAllocs()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					if _, err := Render(c); err != nil {
						b.Fatal(err)
					}
				}
			})
		})
	}
}
func FuzzRender(f *testing.F) {
	for _, s := range []string{"Hello", `<script>&\"`, strings.Repeat("猫", 80), "\x00"} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, message string) {
		c := Default()
		c.Message = message
		b, err := Render(c)
		if err == nil {
			validSVG(t, b)
		}
	})
}
