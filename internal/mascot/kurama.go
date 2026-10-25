package mascot

// A compact fan-art depiction of Kurama. Nine separate tails surround a fox
// silhouette; the eye markings, long ears, muzzle, and claws stay readable small.
func kurama(d *drawing, r Resolved) {
	fur, tail, shade, muzzle := "#df873d", "#ed9d50", "#bc642f", "#f3c484"
	if r.Theme == "dark" {
		fur, tail, shade = "#eea05b", "#f7b56d", "#cb7e43"
	} else if r.Theme == "neon" {
		fur, tail, shade = "#f2a15c", "#ffc578", "#ca7947"
	}
	d.add(`<g stroke="#452d2b" stroke-width="2.6"><g data-part="nine-tails"`)
	if r.Status == "queued" {
		d.add(` transform="translate(18 42) scale(.9 .8)"`)
	} else if r.Status == "failing" {
		d.add(` transform="translate(0 8)"`)
	}
	d.add(`>`)
	// Separate closed paths, rather than one scalloped shape, preserve all nine
	// tails even when a client disables animation or renders a small thumbnail.
	for i, angle := range []int{-84, -63, -42, -21, 0, 21, 42, 63, 84} {
		fill := tail
		if i%2 == 0 {
			fill = fur
		}
		d.add(`<g data-tail="%d" transform="rotate(%d 180 206)"><path d="M173 210 C144 182 142 119 180 52 C172 112 207 175 189 210Z" fill="%s"/><path d="M180 198 Q158 141 177 83" fill="none" stroke="%s" stroke-width="1.8"/></g>`, i+1, angle, fill, shade)
	}
	d.add(`</g>`)
	// A hunched fox torso and haunches avoid reusing the round humanoid bodies.
	d.add(`<path d="M154 164 Q180 151 207 166 L214 210 Q206 231 181 228 Q150 232 144 212Z" fill="%s"/><path d="m161 174 9 11 9 -13 10 13 10 -11 -2 39 -34 0Z" fill="%s" stroke="none"/><path d="M149 201 Q126 200 130 222 L150 230 165 216 M209 201 Q232 201 229 222 L209 231 194 217" fill="%s"/>`, fur, shade, fur)
	if r.Status == "passing" {
		d.add(`<path d="M151 176 Q130 177 119 154 L112 152 116 146 113 138 121 140 124 134 130 149 158 158 M207 176 Q228 175 239 151 L246 149 242 143 245 135 237 137 234 131 228 147 201 159" fill="%s"/><path d="m116 146 5 2 m119 -5 -5 2" fill="none" stroke="#f8dbac" stroke-width="2"/>`, fur)
	} else if r.Status == "failing" {
		d.add(`<path d="M152 175 Q133 184 140 206 L160 212 168 203 155 198 160 185 M206 175 Q225 185 218 207 L198 212 190 202 203 198 198 185" fill="%s"/><path d="m150 203 3 6 m4 -5 3 6 m46 -6 -3 6 m-4 -6 -3 6" fill="none" stroke="#f8dbac" stroke-width="2"/>`, fur)
	} else if r.Status == "running" {
		d.add(`<path d="M153 173 Q130 183 141 220 L161 224 l7 -8 -14 -4 6 -29 M207 173 Q230 182 218 220 L198 224 l-7 -8 14 -4 -6 -29" fill="%s"/><path d="m146 218 1 7 m7 -5 1 6 m54 -8 -1 7 m-7 -5 -1 6" fill="none" stroke="#f8dbac" stroke-width="2"/>`, fur)
	} else if r.Status == "queued" {
		d.add(`<path d="M153 177 Q133 197 154 218 L179 219 l5 -9 -24 -7 4 -17 M206 177 Q227 196 207 218 L180 222 l-5 -9 24 -10 -5 -17" fill="%s"/>`, fur)
	} else {
		d.add(`<path d="M152 174 Q133 185 143 216 L164 221 l5 -8 -15 -7 6 -22 M208 175 Q230 175 228 150 L221 143 213 146 l8 11 -14 8" fill="%s"/>`, fur)
	}
	d.add(`<path d="M151 219 Q136 218 135 233 L169 234 l3 -10 M193 224 l3 10 30 -1 Q225 219 210 219" fill="%s"/><path d="m143 229 -2 6 m11 -6 -1 6 m11 -6 0 6 m38 -5 0 5 m9 -6 1 6 m8 -6 2 5" fill="none" stroke="#f8dbac" stroke-width="2"/>`, fur)
	angle := map[string]int{"passing": -4, "failing": 5, "running": 0, "queued": 11, "unknown": -10}[r.Status]
	d.add(`<g transform="rotate(%d 180 140)"><path d="M136 121 Q119 94 129 65 Q153 81 157 112 M202 111 Q209 78 232 66 Q239 98 222 122" fill="%s"/><path d="M137 107 132 80 148 102 M213 104 228 80 224 109" fill="#4c302c" stroke="none"/><path d="M135 118 Q150 102 180 105 Q212 103 226 121 L234 137 224 138 235 150 222 152 225 164 208 165 Q198 181 180 184 Q161 181 151 167 L135 168 139 155 125 151 137 140 128 137Z" fill="%s"/><path d="m144 148 16 -4 20 9 19 -9 18 4 -10 22 -27 13 -27 -13Z" fill="%s" stroke="none"/><path d="m139 145 12 2 m-10 7 10 1 m59 -8 12 -3 m-12 10 11 -2" fill="none" stroke-width="2"/>`, angle, fur, fur, muzzle)
	kuramaFace(d, r, fur)
	d.add(`</g></g>`)
}

func kuramaFace(d *drawing, r Resolved, fur string) {
	d.add(`<g data-part="kurama-expression">`)
	if r.Status == "queued" {
		d.add(`<path d="m140 128 6 3 -3 7 -5 -4Z m80 0 -6 3 3 7 5 -4Z" fill="#3f2b2b" stroke="none"/>`)
	} else {
		d.add(`<path d="M141 118 158 124 166 142 153 136Z M219 117 201 124 194 142 208 136Z" fill="#3f2b2b" stroke="none"/>`)
	}
	expressiveEyes(d, r, eyeStyle{left: 154, right: 206, y: 132, width: 10.5, height: 12.5, ink: "#3f2b2b", brow: "#3f2b2b", iris: "#c07a30", light: "#f7cf70", lid: fur, slit: true})
	switch r.Status {
	case "passing":
		d.add(`<path d="M160 162 Q179 185 200 161 L194 161 191 169 184 163 177 171 170 164 166 168Z" fill="#472d2b"/><path d="m164 163 6 2 -3 6Z m25 2 6 -3 -1 9Z" fill="#fff1d4" stroke="none"/>`)
	case "failing":
		d.add(`<path d="M163 173 Q180 161 197 172" fill="none" stroke-width="2.5"/><path d="m166 172 3 5 3 -7 m16 0 4 5 2 -5" fill="#fff1d4" stroke-width="1.5"/>`)
	case "running":
		d.add(`<path d="m162 167 17 3 20 -4" fill="none"/><path d="m168 168 3 6 3 -5 m14 0 3 4 3 -6" fill="#fff1d4" stroke-width="1.5"/>`)
	case "queued":
		d.add(`<path d="M165 164 Q180 172 196 163" fill="none" stroke-width="2"/><path d="m168 166 3 6 3 -4 m14 -1 3 4 3 -7" fill="#fff1d4" stroke-width="1.5"/>`)
	default:
		d.add(`<path d="M164 168 q16 5 31 -3" fill="none"/><path d="m170 169 3 6 3 -5" fill="#fff1d4" stroke-width="1.5"/>`)
	}
	d.add(`<path d="M171 152 Q180 148 189 152 L181 159Z" fill="#3f2b2b" stroke-width="1.5"/><path d="m180 159 0 5" fill="none" stroke-width="2"/></g>`)
}
