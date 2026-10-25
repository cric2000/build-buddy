package mascot

// Kappa is a river mechanic: webbed feet, a repaired turtle shell, a water dish,
// and a cucumber sharpened into a carpenter's pencil. Its poses work without CSS.
func kappa(d *drawing, r Resolved) {
	skin, shade, water, scarf := "#91b866", "#537e50", "#83d9dc", "#dc8650"
	if r.Theme == "dark" {
		skin, shade = "#a4c77a", "#78985e"
	}
	if r.Theme == "neon" {
		skin, shade, water, scarf = "#a3d998", "#397c70", "#83f5e8", "#f29fbc"
	}
	d.add(`<g stroke="#283e35" stroke-width="2.8"><g data-part="shell"><path d="M197 157 C242 139 260 177 242 214 Q228 237 201 219Z" fill="%s"/><path d="M206 163 Q241 152 243 184 Q247 212 218 222" fill="none" stroke="#b6c58a"/><path d="m219 163 11 16 -4 23 -17 9 m21 -32 13 -1 m-17 24 13 9" fill="none" stroke="#354f39" stroke-width="2"/><path d="m232 193 9 4 -3 10 -10 -4Z" fill="#e9c78f"/><path d="m233 196 5 2 m-6 2 5 2" stroke-width="1.4"/></g>`, shade)
	// Broad toes and an asymmetric stance give the small-size silhouette weight.
	if r.Status == "passing" {
		d.add(`<path d="M154 214 136 227 124 222 l3 13 37 -1 12 -13 M198 215 213 226 226 221 l6 13 -36 1 -10 -14" fill="%s"/>`, skin)
	} else {
		d.add(`<path d="M152 216 143 228 129 230 l5 8 34 -1 7 -16 M193 218 199 235 l32 1 -5 -8 -14 -1 -6 -14" fill="%s"/>`, skin)
	}
	d.add(`<path d="M147 174 Q169 159 199 172 Q220 188 211 216 Q186 234 150 219 Q134 201 147 174Z" fill="%s"/><path d="M163 180 Q184 173 197 187 L199 216 Q182 225 160 214 Q151 196 163 180Z" fill="#e2d59e" stroke-width="2"/><path d="m163 196 31 0 m-29 9 29 0" fill="none" stroke="#bdba85" stroke-width="1.5"/><path d="m144 173 q31 15 63 -1 l-3 12 q-28 13 -57 -2Z" fill="%s"/><path d="m200 182 19 11 -15 4 -5 14 -9 -28" fill="%s"/><path d="m205 192 6 4" stroke="#f7c584"/>`, skin, scarf, scarf)
	// Hands carry the story: celebrating, protecting the dish, working, or dozing.
	switch r.Status {
	case "passing":
		d.add(`<path d="M146 183 Q119 174 111 147 L101 140 Q96 133 102 131 L111 135 108 123 Q110 117 115 123 L122 139 Q133 148 155 166 M208 183 Q227 180 239 164" fill="%s"/><path d="m113 142 8 -4" fill="none"/>`, skin)
		cucumber(d, 241, 155, -23)
	case "failing":
		d.add(`<path d="M148 185 Q124 175 129 135 L135 117 Q140 113 142 120 L140 139 146 130 Q152 130 150 138 L144 153 159 170 M207 182 Q236 165 225 126 L218 115 Q212 113 212 121 L216 139 210 131 Q205 133 208 140 L215 153 198 168" fill="%s"/>`, skin)
	case "running":
		// A tiny work ledger and cucumber pencil replace generic code symbols.
		d.add(`<path d="m138 190 55 -4 5 38 -55 4Z" fill="#f5e8bf"/><path d="m148 200 21 -2 m-20 10 30 -3 m-29 12 15 -1" stroke="#829477" stroke-width="2"/><path d="M145 180 Q122 184 135 208 L148 211 Q156 209 151 202 L144 197 M208 181 Q226 191 208 207 L193 204 Q189 198 197 197 L203 194" fill="%s"/>`, skin)
		cucumber(d, 191, 189, 33)
	case "queued":
		d.add(`<path d="M145 184 Q123 202 149 220 L169 219 Q176 213 168 209 L153 207 M210 183 Q229 204 208 218 L180 220 Q173 214 181 209 L201 207" fill="%s"/>`, skin)
		cucumber(d, 178, 216, 85)
	default:
		d.add(`<path d="M146 182 Q124 185 128 206 L145 212 Q153 208 145 203 L140 200 M210 182 Q233 178 225 151 L216 145 Q207 149 216 155 L214 167 201 170" fill="%s"/>`, skin)
	}
	angle := map[string]int{"passing": -6, "failing": -3, "running": 5, "queued": 9, "unknown": -10}[r.Status]
	d.add(`<g transform="rotate(%d 180 142)"><path d="m130 124 -12 -11 14 -5 -3 -13 17 3 8 -18 16 10 14 -15 14 12 18 -7 6 15 18 0 -2 15 14 9 -16 13Z" fill="#365942"/><path d="M129 126 Q129 100 155 99 L207 101 Q233 104 234 132 L230 154 Q222 178 182 179 Q143 178 132 159Z" fill="%s"/><path d="M139 139 Q145 116 166 122 Q177 124 180 134 Q189 119 207 126 Q225 131 223 149 Q219 165 203 164 L154 162 Q138 158 139 139Z" fill="#e9e4b8" stroke="none"/><path d="m133 119 12 7 -3 -12 17 7 -1 -15 16 12 9 -12 12 13 14 -8 1 14 14 -3" fill="#365942" stroke="none"/>`, angle, skin)
	// The water remains visible in every mood; a failing build spills it.
	d.add(`<g data-part="water-dish"><path d="M151 94 Q180 106 210 94 L204 107 Q180 117 156 106Z" fill="#d8ce9c"/><ellipse cx="181" cy="94" rx="30" ry="9" fill="%s"/><path d="M160 94 Q171 89 183 93 T202 94" fill="none" stroke="#e0ffff" stroke-width="2"/></g>`, water)
	if r.Status == "failing" {
		d.add(`<path d="M209 95 q15 -9 18 1 m-13 -7 5 -8" fill="none" stroke="%s"/><path d="M233 106 q-8 10 0 12 q8 -2 0 -12Z" fill="%s" stroke-width="1.8"/>`, water, water)
	}
	kappaFace(d, r)
	d.add(`<path d="m139 153 4 1 m-2 6 4 -1 m69 -5 4 -2 m-3 7 4 -1" stroke="#72935c" stroke-width="2"/></g></g>`)
}

func cucumber(d *drawing, x, y, angle int) {
	d.add(`<g data-part="cucumber-pencil" transform="translate(%d %d) rotate(%d)"><path d="M-5 15 -6 -17 Q0 -25 6 -17 L5 15 0 24Z" fill="#426d43"/><path d="m-2 -16 1 29" stroke="#b8ce79" stroke-width="2"/><path d="m-5 15 10 0 -5 9Z" fill="#f2d39a"/><path d="m-1 21 1 3 1 -3" fill="#283e35" stroke-width="1"/></g>`, x, y, angle)
}

func kappaFace(d *drawing, r Resolved) {
	d.add(`<g data-part="kappa-expression">`)
	expressiveEyes(d, r, eyeStyle{left: 157, right: 203, y: 139, width: 11, height: 13, ink: "#283e35", brow: "#283e35", iris: "#527d55", light: "#b8dda0", lid: "#e9e4b8"})
	switch r.Status {
	case "passing":
		d.add(`<path d="M166 154 Q181 148 198 154 L191 165 Q180 175 170 164Z" fill="#c59348"/><path d="m172 158 q9 5 19 0" fill="none" stroke-width="2"/>`)
	case "failing":
		d.add(`<path d="M166 159 Q180 149 196 159 L190 166 Q180 162 171 166Z" fill="#c59348"/><path d="m174 162 q7 -4 14 0" fill="none" stroke-width="1.5"/>`)
	case "running":
		d.add(`<path d="M165 156 Q181 150 197 155 L191 163 170 165Z" fill="#c59348"/><path d="m174 159 15 -1" fill="none" stroke-width="2"/>`)
	case "queued":
		d.add(`<path d="M166 156 Q181 148 197 156 L190 163 172 163Z" fill="#c59348"/><ellipse cx="182" cy="160" rx="4" ry="3" fill="#283e35" stroke="none"/>`)
	default:
		d.add(`<path d="M166 158 Q181 151 198 156 L188 164 173 165Z" fill="#c59348"/><path d="m178 160 9 0" stroke-width="1.5"/>`)
	}
	d.add(`</g>`)
}
