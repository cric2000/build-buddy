package mascot

// Eye geometry is shared by HTTP and CLI rendering. Placement, proportions,
// colors, and pupil shape preserve each character's facial identity.
type eyeStyle struct {
	left, right, y, width, height float64
	ink, brow, iris, light, lid   string
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

func expressiveEyes(d *drawing, r Resolved, s eyeStyle) {
	expression := eyeExpression(r)
	d.add(`<g data-part="expressive-eyes" data-expression="%s" stroke="%s" stroke-width="1.8">`, expression, s.ink)
	for i, x := range []float64{s.left, s.right} {
		w, h, y := s.width, s.height, s.y
		if expression == "curious" && i == 1 {
			w, h, y = w*1.08, h*1.08, y-2
		}
		d.add(`<g data-eye="%d" transform="translate(%.1f %.1f)">`, i+1, x, y)
		if expression == "sleepy" {
			d.add(`<path d="M%.1f -2 Q0 %.1f %.1f -2" fill="none" stroke="%s" stroke-width="2.8"/><path d="m%.1f 0 -2 3 M%.1f 0 l2 3" fill="none" stroke="%s" stroke-width="1.5"/><path d="M%.1f -7 Q0 -4 %.1f -7" fill="none" stroke="%s" stroke-opacity=".35" stroke-width="1.5"/>`, -w, h*.75, w, s.brow, -w+2, w-2, s.brow, -w+2, w-2, s.brow)
		} else {
			d.add(`<g class="eyes"><ellipse rx="%.1f" ry="%.1f" fill="#fff9ed"/>`, w, h)
			irisY := 1.0
			if expression == "focused" || expression == "disappointed" {
				irisY = 3
			}
			d.add(`<ellipse cy="%.1f" rx="%.1f" ry="%.1f" fill="%s" stroke="none"/><path d="M%.1f %.1f Q0 %.1f %.1f %.1f Q0 %.1f %.1f %.1fZ" fill="%s" stroke="none"/>`, irisY, w*.76, h*.83, s.iris, -w*.68, h*.35, h*.85, w*.68, h*.35, h*1.1, -w*.68, h*.35, s.light)
			pupilW := w * .43
			if s.slit {
				pupilW = w * .24
			}
			d.add(`<ellipse cy="%.1f" rx="%.1f" ry="%.1f" fill="%s" stroke="none"/>`, irisY-1, pupilW, h*.66, s.ink)
			switch expression {
			case "sparkly":
				// Large star reflections plus round and lower glints remain visible
				// when animation is off, rather than depending on a twinkle effect.
				d.add(`<g data-part="catchlights" fill="#ffffff" stroke="none"><path class="eye-glint" d="m-3 -8 1.5 3.5 3.5 1.5 -3.5 1.5 -1.5 3.5 -1.5 -3.5 -3.5 -1.5 3.5 -1.5Z"/><circle cx="4" cy="-5" r="2.1"/><ellipse cx="3" cy="5" rx="2.2" ry="1.6"/></g>`)
			case "sad", "teary":
				d.add(`<g data-part="catchlights" fill="#ffffff" stroke="none"><ellipse cx="-3" cy="-5" rx="2.7" ry="3.3"/><circle cx="3" cy="1" r="1.5"/></g>`)
			case "disappointed":
				d.add(`<path d="M%.1f %.1f Q0 %.1f %.1f %.1f L%.1f %.1f Q0 %.1f %.1f %.1fZ" fill="%s" stroke="none"/><path d="M%.1f -3 Q0 0 %.1f -3" fill="none" stroke="%s" stroke-width="2.2"/><circle cx="-3" cy="3" r="1.8" fill="#fff9ed" stroke="none"/>`, -w, -h, -h*1.3, w, -h, w, -3.0, 0.0, -w, -3.0, s.lid, -w, w, s.brow)
			case "focused":
				d.add(`<path d="M%.1f %.1f Q0 %.1f %.1f %.1f L%.1f -4 Q0 -1 %.1f -4Z" fill="%s" stroke="none"/><path d="M%.1f -4 Q0 -1 %.1f -4" fill="none" stroke="%s" stroke-width="2.5"/><circle cx="-2" cy="2" r="1.8" fill="#ffffff" stroke="none"/>`, -w, -h, -h*1.3, w, -h, w, -w, s.lid, -w, w, s.brow)
			default:
				d.add(`<g data-part="catchlights" fill="#ffffff" stroke="none"><ellipse cx="-3" cy="-5" rx="2.8" ry="3.5"/><circle cx="3" cy="3" r="1.8"/></g>`)
			}
			d.add(`</g>`)
			if expression == "sad" || expression == "teary" || expression == "disappointed" {
				// Water gathers along the lower lid; larger failures overflow it.
				d.add(`<g data-part="tears"><path d="M%.1f %.1f Q0 %.1f %.1f %.1f" fill="none" stroke="#85cbd9" stroke-width="3.2"/><path d="M%.1f %.1f Q0 %.1f %.1f %.1f" fill="none" stroke="#e1fbff" stroke-width="1.2"/>`, -w*.7, h*.65, h*1.15, w*.7, h*.65, -w*.5, h*.75, h, w*.5, h*.75)
				if expression == "teary" {
					x := w * .65
					if i == 0 {
						x = -x
					}
					d.add(`<path d="M%.1f %.1f q-3 6 -3 12 q3 5 6 0 q0 -6 -3 -12Z" fill="#98dce5" stroke="#558f9e" stroke-width="1"/><path d="m%.1f %.1f -1 7" fill="none" stroke="#eaffff" stroke-width="1.3"/>`, x, h*.8, x, h+2)
				}
				d.add(`</g>`)
			}
		}
		// Brows are outside the blinking eye group so their emotional pose holds.
		browY := -h - 4
		sign := 1.0
		if i == 1 {
			sign = -1
		}
		switch expression {
		case "sad", "teary", "disappointed":
			d.add(`<path d="M%.1f %.1f Q0 %.1f %.1f %.1f" fill="none" stroke="%s" stroke-width="2.6"/>`, -w, browY+sign*3, browY-2, w, browY-sign*3, s.brow)
		case "sparkly":
			d.add(`<path d="M%.1f %.1f Q0 %.1f %.1f %.1f" fill="none" stroke="%s" stroke-width="2.2"/>`, -w+1, browY+2, browY-4, w-1, browY+2, s.brow)
		case "focused":
			d.add(`<path d="M%.1f %.1f %.1f %.1f" fill="none" stroke="%s" stroke-width="2.5"/>`, -w, browY-sign*2, w, browY+sign*2, s.brow)
		case "curious":
			d.add(`<path d="M%.1f %.1f Q0 %.1f %.1f %.1f" fill="none" stroke="%s" stroke-width="2.3"/>`, -w, browY+sign*2, browY-2, w, browY-sign*2, s.brow)
		}
		d.add(`</g>`)
	}
	d.add(`</g>`)
}
