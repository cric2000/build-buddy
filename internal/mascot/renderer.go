package mascot

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"
	"unicode"
)

type palette struct{ bg, panel, ink, muted, accent, soft, highlight, line string }

func colors(theme string) palette {
	switch theme {
	case "dark":
		return palette{"#111827", "#1e293b", "#f1f5f9", "#b3c2d6", "#a5b4fc", "#5b66a5", "#f9cf85", "#34445f"}
	case "neon":
		return palette{"#140f26", "#201736", "#faf5ff", "#c3afd9", "#70f4e5", "#297d85", "#ff80c8", "#49335d"}
	default:
		return palette{"#f5f3ee", "#fffefa", "#29354b", "#596779", "#a4c9c1", "#659e99", "#f3bc78", "#e7e6df"}
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
	d.add(`<defs><linearGradient id="shell" x1="0" y1="0" x2="0.8" y2="1"><stop stop-color="%s"/><stop offset="1" stop-color="%s"/></linearGradient><radialGradient id="halo"><stop stop-color="%s" stop-opacity=".2"/><stop offset="1" stop-color="%s" stop-opacity="0"/></radialGradient></defs>`, p.accent, p.soft, p.accent, p.accent)
	if c.Animate {
		d.add(`<style>
 .buddy{transform-origin:180px 170px;animation:%s 3.2s ease-in-out infinite}
 .eyes{transform-box:fill-box;transform-origin:center;animation:blink 5.4s infinite}
 .tail{transform-origin:229px 206px;animation:tail 3.2s ease-in-out infinite}
 .spark{transform-box:fill-box;transform-origin:center;animation:twinkle 2.8s ease-in-out infinite}
 @keyframes bounce{0%%,100%%{transform:translateY(0) rotate(-2deg)}50%%{transform:translateY(-7px) rotate(2deg)}}
 @keyframes shake{0%%,85%%,100%%{transform:translateX(0)}89%%,97%%{transform:translateX(-2px) rotate(-1deg)}93%%{transform:translateX(2px) rotate(1deg)}}
 @keyframes work{0%%,100%%{transform:rotate(-2deg)}50%%{transform:rotate(2deg) translateY(-2px)}}
 @keyframes float{0%%,100%%{transform:translateY(0)}50%%{transform:translateY(-6px)}}
 @keyframes blink{0%%,43%%,47%%,100%%{transform:scaleY(1)}45%%{transform:scaleY(.12)}}
 @keyframes tail{0%%,100%%{transform:rotate(-5deg)}50%%{transform:rotate(8deg)}}
 @keyframes twinkle{0%%,100%%{opacity:.55;transform:scale(.88)}50%%{opacity:1;transform:scale(1.08)}}
 @media(prefers-reduced-motion:reduce){.buddy,.eyes,.tail,.spark{animation:none!important}}
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
	case "robot":
		robot(d, r, p)
	case "cat":
		cat(d, r, p)
	case "ghost":
		ghost(d, r, p)
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
func robot(d *drawing, r Resolved, p palette) {
	d.add(`<path d="M155 222 v13 h-19 q-5 -12 4 -15 M204 222 v13 h19 q5 -12 -4 -15" fill="%s"/><rect x="146" y="181" width="68" height="48" rx="17" fill="url(#shell)"/><path d="M148 194 Q123 191 121 211 M212 194 Q235 190 240 177" fill="none"/><circle cx="121" cy="213" r="6" fill="%s"/><circle cx="241" cy="174" r="6" fill="%s"/>`, p.soft, p.accent, p.accent)
	d.add(`<path d="M180 93 v-15"/><circle cx="180" cy="73" r="7" fill="%s"/><rect x="109" y="124" width="15" height="34" rx="6" fill="%s"/><rect x="236" y="124" width="15" height="34" rx="6" fill="%s"/><rect x="119" y="95" width="122" height="94" rx="28" fill="url(#shell)"/><path d="M137 110 h30" stroke="%s" stroke-width="5" opacity=".6"/><rect x="130" y="119" width="100" height="55" rx="18" fill="%s" stroke="none"/>`, p.highlight, p.soft, p.soft, p.panel, p.panel)
	face(d, r, p, 180, 147, 1)
	d.add(`<path d="M173 206 q7 -12 14 0 q0 5 -7 10 q-7 -5 -7 -10Z" fill="%s" stroke="none"/><circle cx="201" cy="215" r="2" fill="%s" stroke="none"/>`, p.highlight, p.panel)
}
func cat(d *drawing, r Resolved, p palette) {
	d.add(`<path class="tail" d="M223 213 C269 221 271 188 253 180 C240 175 240 193 251 193" fill="none" stroke="%s" stroke-width="14"/><ellipse cx="180" cy="202" rx="44" ry="30" fill="url(#shell)"/><ellipse cx="180" cy="209" rx="24" ry="20" fill="%s" stroke="none"/><path d="M133 121 L131 82 Q149 82 163 105 M197 105 Q214 83 229 82 L227 125" fill="url(#shell)"/><path d="m139 109 -1 -16 15 12 m54 0 15 -12 -1 16" fill="%s" stroke="none"/><path d="M124 142 Q120 102 180 101 Q240 102 236 142 Q239 185 180 189 Q121 185 124 142Z" fill="url(#shell)"/><path d="m170 103 3 13 m7 -14 0 17 m10 -16 -3 13" stroke="%s" stroke-width="4"/>`, p.soft, p.panel, p.highlight, p.soft)
	face(d, r, p, 180, 144, 1.05)
	d.add(`<path d="m173 149 7 4 7 -4Z" fill="%s" stroke="none"/><path d="m133 153 -25 -4 m26 13 -24 5 m117 -14 25 -4 m-26 13 24 5" fill="none" stroke-width="2"/><ellipse cx="153" cy="228" rx="16" ry="9" fill="%s"/><ellipse cx="207" cy="228" rx="16" ry="9" fill="%s"/><path d="m151 225 0 5 m6 -5 0 5 m48 -5 0 5 m6 -5 0 5" stroke-width="1.5"/>`, p.ink, p.accent, p.accent)
}
func ghost(d *drawing, r Resolved, p palette) {
	d.add(`<path d="M124 166 Q103 166 107 183 Q112 191 125 183 M236 166 Q257 166 253 183 Q248 191 235 183" fill="url(#shell)"/><path d="M122 157 C122 74 238 74 238 157 L240 225 Q227 211 218 228 Q208 240 195 225 Q180 212 169 228 Q158 240 146 225 Q135 214 120 229Z" fill="url(#shell)"/><path d="M138 132 Q146 109 162 108" fill="none" stroke="%s" stroke-width="6" opacity=".55"/>`, p.panel)
	face(d, r, p, 180, 152, 1.2)
	d.add(`<path d="m166 196 q14 10 28 0" fill="none" stroke="%s" stroke-width="2" opacity=".5"/>`, p.panel)
}
func face(d *drawing, r Resolved, p palette, x, y int, scale float64) {
	d.add(`<g transform="translate(%d %d) scale(%.2f)" stroke="%s">`, x, y, scale, p.ink)
	d.add(`<ellipse cx="-33" cy="10" rx="8" ry="4" fill="%s" opacity=".65" stroke="none"/><ellipse cx="33" cy="10" rx="8" ry="4" fill="%s" opacity=".65" stroke="none"/>`, p.highlight, p.highlight)
	switch r.Mood {
	case "happy":
		d.add(`<g class="eyes"><path d="M-29 0 q7 -13 14 0 M15 0 q7 -13 14 0" fill="none"/></g><path d="M-9 10 Q0 27 9 10Z" fill="%s"/><path d="M-4 16 h8" stroke="%s" stroke-width="2"/>`, p.ink, p.highlight)
	case "stressed":
		d.add(`<g class="eyes"><ellipse cx="-22" cy="0" rx="5" ry="7" fill="%s" stroke="none"/><ellipse cx="22" cy="0" rx="5" ry="7" fill="%s" stroke="none"/></g><path d="m-30 -14 14 -5 m32 0 14 5 M-8 18 Q0 8 8 18" fill="none"/>`, p.ink, p.ink)
		if r.Worried {
			d.add(`<path d="m-37 -16 -4 -7 m78 7 4 -7" stroke-width="2"/><path d="M33 3 Q26 15 33 17 Q40 15 33 3" fill="%s" stroke="none"/>`, p.soft)
		}
	case "focused":
		d.add(`<path d="m-31 -12 18 3 m26 0 18 -3" fill="none"/><g class="eyes"><rect x="-27" y="-4" width="9" height="12" rx="4" fill="%s" stroke="none"/><rect x="18" y="-4" width="9" height="12" rx="4" fill="%s" stroke="none"/></g><path d="M-6 17 h12" fill="none"/>`, p.ink, p.ink)
	case "sleepy":
		d.add(`<path d="M-30 -1 q8 10 16 0 M14 -1 q8 10 16 0" fill="none"/><ellipse cx="0" cy="17" rx="4" ry="5" fill="%s" stroke="none"/>`, p.ink)
	default:
		d.add(`<path d="m-30 -14 14 0 m31 -3 14 -4" fill="none"/><g class="eyes"><ellipse cx="-21" cy="0" rx="5" ry="7" fill="%s" stroke="none"/><ellipse cx="23" cy="-2" rx="6" ry="8" fill="%s" stroke="none"/><circle cx="-20" cy="-2" r="1.5" fill="%s" stroke="none"/><circle cx="24" cy="-4" r="2" fill="%s" stroke="none"/></g><path d="M-4 17 q5 -4 9 0" fill="none"/>`, p.ink, p.ink, p.panel, p.panel)
	}
	d.add(`</g>`)
}
