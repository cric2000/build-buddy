package mascot

// Each character keeps its own eye silhouette and materials while sharing the
// same build-state rules and SVG construction in the API and CLI.
type eyeStyle struct {
	left, right, y, width, height float64
	ink, brow, iris, light, lid   string
	shape, rim                    string
	slit                          bool
}

func eyeExpression(r Resolved) string {
	switch r.Status {
	case "passing":
		return "sparkly"
	case "failing":
		if r.Worried {
			return "teary"
		}
		if r.FailedTests != nil {
			return "disappointed"
		}
		return "sad"
	case "running":
		return "focused"
	case "queued":
		return "sleepy"
	default:
		return "curious"
	}
}

func eyeContour(d *drawing, w, h float64, shape string, side int) {
	if shape == "fox" {
		direction := 1
		if side == 1 {
			direction = -1
		}
		d.add(`<path transform="scale(%d 1)" d="M%.1f 0 C%.1f %.1f %.1f %.1f %.1f %.1f Q%.1f %.1f %.1f %.1f Q%.1f %.1f %.1f 0Z"/>`, direction, -w, -w*.4, -h*1.1, w*.6, -h*1.3, w*.9, -h*.35, w*1.15, h*.25, w, h, -w*.35, h*1.1, -w)
		return
	}
	d.add(`<ellipse rx="%.1f" ry="%.1f"/>`, w, h)
}

func expressiveEyes(d *drawing, r Resolved, s eyeStyle) {
	expression := eyeExpression(r)
	d.add(`<g data-part="expressive-eyes" data-expression="%s" stroke="%s" stroke-width="1.6">`, expression, s.ink)
	for i, x := range []float64{s.left, s.right} {
		w, h, y := s.width, s.height, s.y
		if expression == "curious" && i == 1 {
			w, h, y = w*1.06, h*1.06, y-1.5
		}
		d.add(`<g data-eye="%d" transform="translate(%.1f %.1f)">`, i+1, x, y)
		if expression == "sleepy" {
			d.add(`<path d="M%.1f -1 Q0 %.1f %.1f -1" fill="none" stroke="%s" stroke-width="2.4"/><path d="m%.1f 0 -1.5 2 M%.1f 0 l1.5 2" fill="none" stroke="%s" stroke-width="1.2"/>`, -w, h*.65, w, s.brow, -w+1, w-1, s.brow)
		} else {
			d.add(`<g class="eyes">`)
			d.add(`<defs><clipPath id="eye-cutout-%d" clipPathUnits="userSpaceOnUse">`, i+1)
			eyeContour(d, w, h, s.shape, i)
			d.add(`</clipPath></defs>`)
			if s.shape == "button" {
				rim := s.ink
				if s.rim != "" {
					rim = s.rim
				}
				d.add(`<ellipse rx="%.1f" ry="%.1f" fill="%s" stroke="%s" stroke-width="1.4"/><ellipse cy="%.1f" rx="%.1f" ry="%.1f" fill="%s" fill-opacity=".3" stroke="none"/>`, w, h, s.ink, rim, h*.4, w*.6, h*.3, s.light)
			} else {
				d.add(`<g fill="#fffdf8">`)
				eyeContour(d, w, h, s.shape, i)
				d.add(`</g><g clip-path="url(#eye-cutout-%d)" stroke="none">`, i+1)
				irisY := 1.0
				if expression == "focused" || expression == "disappointed" {
					irisY = 3
				}
				d.add(`<ellipse cy="%.1f" rx="%.1f" ry="%.1f" fill="%s"/><path d="M%.1f %.1f Q0 %.1f %.1f %.1f Q0 %.1f %.1f %.1fZ" fill="%s"/>`, irisY, w*.76, h*.91, s.iris, -w*.66, h*.25, h*.7, w*.66, h*.25, h*1.15, -w*.66, h*.25, s.light)
				pupilW := w * .43
				if s.slit {
					pupilW = w * .18
				}
				d.add(`<ellipse cy="%.1f" rx="%.1f" ry="%.1f" fill="%s"/></g>`, irisY-1, pupilW, h*.68, s.ink)
			}
			switch expression {
			case "sparkly":
				d.add(`<g data-part="catchlights" fill="#ffffff" stroke="none"><ellipse class="eye-glint" cx="%.1f" cy="%.1f" rx="%.1f" ry="%.1f"/><circle cx="%.1f" cy="%.1f" r="1.3"/></g>`, -w*.3, -h*.4, w*.3, h*.28, w*.35, h*.35)
				if r.Celebration && s.shape != "button" {
					d.add(`<path d="m%.1f %.1f 1 2.5 2.5 1 -2.5 1 -1 2.5 -1 -2.5 -2.5 -1 2.5 -1Z" fill="#fffdf4" stroke="none"/>`, w*.32, -h*.6)
				}
			case "sad", "teary":
				d.add(`<g data-part="catchlights" fill="#ffffff" stroke="none"><ellipse cx="%.1f" cy="%.1f" rx="%.1f" ry="%.1f"/><circle cx="%.1f" cy="%.1f" r="1.2"/></g>`, -w*.3, -h*.35, w*.28, h*.32, w*.3, h*.2)
			case "disappointed", "focused":
				lidY := -h * .15
				if expression == "focused" {
					lidY = -h * .3
				}
				d.add(`<g clip-path="url(#eye-cutout-%d)">`, i+1)
				d.add(`<path d="M%.1f %.1f Q0 %.1f %.1f %.1f L%.1f %.1f Q0 %.1f %.1f %.1fZ" fill="%s" stroke="none"/><path d="M%.1f %.1f Q0 %.1f %.1f %.1f" fill="none" stroke="%s" stroke-width="1.8"/><circle cx="%.1f" cy="%.1f" r="1.4" fill="#fffdf8" stroke="none"/>`, -w-1, -h, -h*1.4, w+1, -h, w+1, lidY, lidY+2, -w-1, lidY, s.lid, -w, lidY, lidY+2, w, lidY, s.brow, -w*.3, h*.3)
				d.add(`</g>`)
			default:
				d.add(`<g data-part="catchlights" fill="#ffffff" stroke="none"><ellipse cx="%.1f" cy="%.1f" rx="%.1f" ry="%.1f"/><circle cx="%.1f" cy="%.1f" r="1.1"/></g>`, -w*.3, -h*.35, w*.27, h*.3, w*.3, h*.3)
			}
			d.add(`</g>`)
			if expression == "sad" || expression == "teary" || expression == "disappointed" {
				d.add(`<g data-part="tears"><path d="M%.1f %.1f Q0 %.1f %.1f %.1f" fill="none" stroke="#87cedb" stroke-width="2.8"/><path d="M%.1f %.1f Q0 %.1f %.1f %.1f" fill="none" stroke="#edffff" stroke-width="1"/>`, -w*.65, h*.65, h*1.1, w*.65, h*.65, -w*.45, h*.75, h*.95, w*.45, h*.75)
				if expression == "teary" {
					x := w * .65
					if i == 0 {
						x = -x
					}
					d.add(`<path d="M%.1f %.1f q-2.5 5 -2.5 10 q2.5 4 5 0 q0 -5 -2.5 -10Z" fill="#a2e0e9" stroke="#6398a6" stroke-width=".8"/><path d="m%.1f %.1f -.6 6" fill="none" stroke="#f3ffff" stroke-width="1.1"/>`, x, h*.8, x, h+1)
				}
				d.add(`</g>`)
			}
		}
		browY := -h - 4
		sign := 1.0
		if i == 1 {
			sign = -1
		}
		switch expression {
		case "sad", "teary", "disappointed":
			d.add(`<path d="M%.1f %.1f Q0 %.1f %.1f %.1f" fill="none" stroke="%s" stroke-width="1.6"/>`, -w, browY+sign*2, browY-1, w, browY-sign*2, s.brow)
		case "focused":
			d.add(`<path d="M%.1f %.1f %.1f %.1f" fill="none" stroke="%s" stroke-width="1.8"/>`, -w, browY-sign*1.5, w, browY+sign*1.5, s.brow)
		case "curious":
			d.add(`<path d="M%.1f %.1f Q0 %.1f %.1f %.1f" fill="none" stroke="%s" stroke-width="1.5"/>`, -w*.7, browY+sign, browY-1, w*.7, browY-sign, s.brow)
		}
		d.add(`</g>`)
	}
	d.add(`</g>`)
}
