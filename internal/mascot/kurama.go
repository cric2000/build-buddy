package mascot

// Kurama sits on soft haunches behind two small forepaws. Broad overlapping
// tails and ear-to-eye markings keep the silhouette close to a little fox.
func kurama(d *drawing, r Resolved) {
	fur, tail, shade := "#f58d22", "#f3952e", "#ce711b"
	if r.Theme == "dark" || r.Theme == "neon" {
		fur, tail, shade = "#ff9a38", "#f9a343", "#d77d29"
	}
	d.add(`<g stroke="#25201d" stroke-width="2" fill="%s"><g data-part="nine-tails">`, fur)
	d.add(`<path data-tail="1" d="M160 198 C139 145 151 81 180 69 C210 82 221 144 199 198Z" fill="%s"/>`, tail)
	tails := []string{
		`M175 211 C132 194 102 151 101 109 Q132 99 153 127 C174 150 186 184 185 203Z`,
		`M173 214 C124 208 81 171 79 133 Q109 117 145 151 Q175 181 185 210Z`,
		`M174 217 C122 232 65 214 56 175 Q81 160 113 178 Q145 190 185 211Z`,
		`M177 220 C125 244 69 236 52 200 Q67 188 100 196 Q137 201 188 217Z`,
	}
	for i, path := range tails {
		fill := tail
		if i%2 == 1 {
			fill = fur
		}
		d.add(`<path data-tail="%d" d="%s" fill="%s"/><path data-tail="%d" transform="translate(360 0) scale(-1 1)" d="%s" fill="%s"/>`, i*2+2, path, fill, i*2+3, path, fill)
	}
	d.add(`</g><path d="M157 161 Q180 150 203 162 Q218 185 215 214 Q209 234 180 233 Q150 234 145 214 Q143 184 157 161Z"/><path d="M164 178 Q180 169 197 179 L199 221 Q182 234 161 221Z" fill="%s" fill-opacity=".15" stroke="none"/>`, shade)
	d.add(`<ellipse cx="145" cy="217" rx="15" ry="22" transform="rotate(-28 145 217)"/><ellipse cx="215" cy="217" rx="15" ry="22" transform="rotate(28 215 217)"/><path d="m133 219 3 5 m83 -5 -3 5" fill="none" stroke-width="1.2"/>`)
	switch r.Status {
	case "failing":
		d.add(`<path d="M158 180 Q146 185 151 207 L170 216 Q181 211 171 204 L162 199 M202 180 Q214 185 209 207 L190 216 Q179 211 189 204 L198 199"/>`)
	case "queued":
		d.add(`<path d="M157 182 Q144 200 159 219 L180 222 l4 -7 -21 -9 M203 182 Q216 201 201 219 L181 222 l-4 -7 21 -9"/>`)
	default:
		d.add(`<path d="M161 179 C150 180 153 210 151 228 Q153 239 170 235 Q177 231 174 221 L176 189 M199 179 C210 180 207 210 209 228 Q207 239 190 235 Q183 231 186 221 L184 189"/><path d="m158 229 -1 5 m7 -5 0 6 m32 -6 0 6 m7 -6 1 5" fill="none" stroke-width="1.1"/>`)
	}
	angle := map[string]int{"passing": -2, "failing": 3, "running": 0, "queued": 6, "unknown": -6}[r.Status]
	d.add(`<g transform="rotate(%d 180 133)"><path d="M144 111 C115 105 96 78 97 52 Q132 51 162 97 M216 111 C245 105 264 78 263 52 Q228 51 198 97"/><path d="M124 110 C129 90 153 84 180 84 Q220 84 236 111 L242 148 Q246 161 222 171 Q180 187 138 171 Q114 161 118 148Z"/>`, angle)
	// The dark shapes are continuous from the inner ears to the eyes.
	d.add(`<path data-part="fox-markings" d="M104 64 C119 91 147 99 161 116 Q174 132 171 151 Q149 154 133 137 C112 116 105 89 104 64Z M256 64 C241 91 213 99 199 116 Q186 132 189 151 Q211 154 227 137 C248 116 255 89 256 64Z" fill="#171817" stroke="none"/>`)
	kuramaFace(d, r)
	d.add(`</g></g>`)
}

func kuramaFace(d *drawing, r Resolved) {
	d.add(`<g data-part="kurama-expression">`)
	eyes := eyeStyle{left: 153, right: 207, y: 134, width: 15, height: 17, ink: "#121413", brow: "#171817", iris: "#c82f30", light: "#f45b43", lid: "#171817", shape: "fox", slit: true}
	if r.Status == "queued" {
		eyes.brow = "#fff0dc"
	}
	expressiveEyes(d, r, eyes)
	d.add(`<ellipse cx="180" cy="156" rx="5.5" ry="2.8" fill="#171817" stroke="none"/>`)
	switch r.Status {
	case "passing":
		d.add(`<path d="M169 164 q5 6 11 0 q6 6 11 0" fill="none" stroke-width="1.8"/>`)
	case "failing":
		d.add(`<path d="M170 168 q10 -8 20 0" fill="none" stroke-width="1.8"/>`)
	case "running":
		d.add(`<path d="m174 166 12 0" fill="none" stroke-width="1.8"/>`)
	case "queued":
		d.add(`<path d="M172 164 q8 5 16 0" fill="none" stroke-width="1.5"/>`)
	default:
		d.add(`<ellipse cx="181" cy="166" rx="3" ry="3.5" fill="none" stroke-width="1.5"/>`)
	}
	d.add(`</g>`)
}
