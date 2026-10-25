package mascot

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"
	"unicode"
)

type palette struct{ bg, panel, ink, muted, accent, highlight, line string }

func colors(theme string) palette {
	switch theme {
	case "dark":
		return palette{"#111827", "#1e293b", "#f1f5f9", "#b3c2d6", "#a5b4fc", "#f9cf85", "#34445f"}
	case "neon":
		return palette{"#140f26", "#201736", "#faf5ff", "#c3afd9", "#70f4e5", "#ff80c8", "#49335d"}
	default:
		return palette{"#f5f3ee", "#fffefa", "#29354b", "#596779", "#a4c9c1", "#f3bc78", "#e7e6df"}
	}
}
func escape(s string) string {
	var b bytes.Buffer
	_ = xml.EscapeText(&b, []byte(s))
	return b.String()
}

type drawing struct{ strings.Builder }

func (d *drawing) add(format string, args ...any) { fmt.Fprintf(&d.Builder, format, args...) }

// Render validates before drawing and has no shared mutable state or external assets.
func Render(c Config) ([]byte, error) {
	r, err := Resolve(c)
	if err != nil {
		return nil, err
	}
	p := colors(c.Theme)
	d := &drawing{}
	d.add(`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 360 360" role="img" aria-labelledby="title desc" data-character="%s" data-mood="%s" data-theme="%s">`, c.Size, c.Size, c.Character, r.Mood, c.Theme)
	d.add(`<title id="title">BuildBuddy: %s build</title><desc id="desc">%s %s. %s</desc>`, c.Status, r.Mood, c.Character, escape(r.Message))
	d.add(`<defs><radialGradient id="halo"><stop stop-color="%s" stop-opacity=".2"/><stop offset="1" stop-color="%s" stop-opacity="0"/></radialGradient></defs>`, p.accent, p.accent)
	if c.Animate {
		d.add(`<style>
 .buddy{transform-origin:180px 170px;animation:%s 3.2s ease-in-out infinite}
 .eyes{transform-box:fill-box;transform-origin:center;animation:blink 5.4s infinite}
 .eye-glint,.spark{transform-box:fill-box;transform-origin:center;animation:twinkle 2.8s ease-in-out infinite}
 @keyframes bounce{0%%,100%%{transform:translateY(0) rotate(-2deg)}50%%{transform:translateY(-7px) rotate(2deg)}}
 @keyframes shake{0%%,85%%,100%%{transform:translateX(0)}89%%,97%%{transform:translateX(-2px) rotate(-1deg)}93%%{transform:translateX(2px) rotate(1deg)}}
 @keyframes work{0%%,100%%{transform:rotate(-2deg)}50%%{transform:rotate(2deg) translateY(-2px)}}
 @keyframes float{0%%,100%%{transform:translateY(0)}50%%{transform:translateY(-6px)}}
 @keyframes blink{0%%,43%%,47%%,100%%{transform:scaleY(1)}45%%{transform:scaleY(.12)}}
 @keyframes twinkle{0%%,100%%{opacity:.55;transform:scale(.88)}50%%{opacity:1;transform:scale(1.08)}}
 @media(prefers-reduced-motion:reduce){.buddy,.eyes,.eye-glint,.spark{animation:none!important}}
 </style>`, r.animationName())
	}
	d.add(`<rect width="360" height="360" rx="26" fill="%s"/><rect x="10" y="10" width="340" height="340" rx="21" fill="%s" stroke="%s"/>`, p.bg, p.panel, p.line)
	d.add(`<g font-family="ui-sans-serif,system-ui,sans-serif"><text x="28" y="39" font-size="10" font-weight="700" letter-spacing="2" fill="%s">BUILD BUDDY</text><circle cx="322" cy="35" r="4" fill="%s"/>`, p.muted, statusColor(c.Status))
	d.add(`<ellipse cx="180" cy="162" rx="123" ry="111" fill="url(#halo)"/><circle cx="180" cy="160" r="87" fill="none" stroke="%s" stroke-dasharray="2 8"/><ellipse cx="180" cy="240" rx="57" ry="8" fill="%s" opacity=".13"/>`, p.line, p.ink)
	decorations(d, r, p)
	d.add(`<g class="buddy" stroke="%s" stroke-width="3" stroke-linecap="round" stroke-linejoin="round">`, p.ink)
	switch c.Character {
	case "kappa":
		kappa(d, r)
	case "tanuki":
		tanuki(d, r)
	case "tengu":
		tengu(d, r)
	case "kurama":
		kurama(d, r)
	case "chopper":
		chopper(d, r)
	}
	d.add(`</g>`)
	d.add(`<rect x="%d" y="259" width="%d" height="22" rx="11" fill="%s" fill-opacity=".13"/><circle cx="%d" cy="270" r="3" fill="%s"/><text x="188" y="274" text-anchor="middle" fill="%s" font-size="10" font-weight="750" letter-spacing="1.5">%s</text>`, 180-statusWidth(c.Status)/2, statusWidth(c.Status), statusColor(c.Status), 193-statusWidth(c.Status)/2, statusColor(c.Status), p.ink, strings.ToUpper(c.Status))
	lines := messageLines(r.Message)
	for i, line := range lines {
		fit := ""
		if captionWidth(line, messageFont(r.Message)) > 300 {
			fit = ` textLength="300" lengthAdjust="spacingAndGlyphs"`
		}
		d.add(`<text x="180" y="%d" text-anchor="middle" fill="%s" font-size="%d" font-weight="600"%s>%s</text>`, 302+i*15, p.ink, messageFont(r.Message), fit, escape(line))
	}
	if c.Status == "running" && c.Progress != nil {
		d.add(`<rect x="95" y="331" width="170" height="5" rx="2.5" fill="%s"/><rect x="95" y="331" width="%.1f" height="5" rx="2.5" fill="%s"/>`, p.line, float64(*c.Progress)*1.7, p.accent)
	} else {
		detail := "YOUR BUILD HAS FEELINGS"
		if c.Coverage != nil {
			detail = fmt.Sprintf("COVERAGE  %d%%", *c.Coverage)
			if r.CoverageWarning {
				detail += "  ·  ROOM TO GROW"
			}
		}
		d.add(`<text x="180" y="337" text-anchor="middle" fill="%s" font-size="8" font-weight="600" letter-spacing="1.2">%s</text>`, p.muted, detail)
	}
	d.add(`</g></svg>`)
	return []byte(d.String()), nil
}
func (r Resolved) animationName() string {
	if r.Animation == "blink" {
		return "float"
	}
	return r.Animation
}
func statusWidth(s string) int { return 45 + len(s)*8 }
func statusColor(s string) string {
	switch s {
	case "passing":
		return "#36a885"
	case "failing":
		return "#e27683"
	case "running":
		return "#7b9fea"
	case "queued":
		return "#b599d9"
	default:
		return "#c49b55"
	}
}
func messageFont(s string) int {
	if len([]rune(s)) > 28 {
		return 11
	}
	return 13
}

// Constrain unusually wide captions (e.g. 80 CJK characters) without stretching normal text.
func captionWidth(s string, font int) float64 {
	width := 0.0
	for _, r := range s {
		switch {
		case r == 'W' || r == 'M' || r == 'w' || r == 'm':
			width += 1
		case r == 'i' || r == 'l' || r == ' ' || r == '.' || r == '!':
			width += 0.35
		case r < 128 && unicode.IsUpper(r):
			width += 0.8
		case r < 128:
			width += 0.65
		default:
			width += 1.2
		}
	}
	return width * float64(font)
}
func messageLines(s string) []string {
	runes := []rune(s)
	if len(runes) <= 40 {
		return []string{s}
	}
	split := 40
	for i := 40; i >= 24; i-- {
		if runes[i] == ' ' {
			split = i
			break
		}
	}
	// Balance at most 80 runes into two lines to keep the footer clear.
	rest := strings.TrimSpace(string(runes[split:]))
	lines := []string{strings.TrimSpace(string(runes[:split]))}
	if len([]rune(rest)) > 40 { // Balance two lines to keep the card footer clear.
		return []string{string(runes[:40]), string(runes[40:])}
	}
	return append(lines, rest)
}
func decorations(d *drawing, r Resolved, p palette) {
	spark := p.highlight
	if r.Theme == "neon" {
		spark = "#ffe58a"
	}
	if r.Celebration || r.Status == "passing" {
		for _, pt := range [][2]int{{76, 100}, {285, 150}, {263, 78}} {
			d.add(`<path class="spark" d="M%d %d l3 8 8 3 -8 3 -3 8 -3 -8 -8 -3 8 -3Z" fill="%s"/>`, pt[0], pt[1], spark)
		}
		if r.Celebration {
			d.add(`<path d="m91 193 4 7 m176 -91 5 -5 m-173 -26 5 3" stroke="%s" stroke-width="3" stroke-linecap="round"/>`, p.accent)
		}
	} else if r.Status == "queued" {
		d.add(`<g fill="%s" font-weight="700"><text x="251" y="118" font-size="18">z</text><text x="272" y="96" font-size="12">z</text></g>`, p.muted)
	} else if r.Status == "running" {
		d.add(`<g stroke="%s" stroke-width="3" stroke-linecap="round"><path d="m76 135 8 -8 -8 -8 m200 68 8 -8 -8 -8" fill="none"/><path d="m72 184 10 0"/></g>`, p.accent)
	} else if r.Status == "failing" {
		d.add(`<path d="M269 109 C263 119 262 121 263 124 C265 131 275 128 274 122Z" fill="%s"/>`, p.accent)
	} else {
		d.add(`<text x="264" y="114" fill="%s" font-size="25" font-weight="600">?</text>`, spark)
	}
}
