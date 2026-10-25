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
	kuramaFace(d, r)
	d.add(`</g></g>`)
}

func kuramaFace(d *drawing, r Resolved) {
	d.add(`<g data-part="kurama-expression"><path d="M141 118 158 124 166 142 153 136Z M219 117 201 124 194 142 208 136Z" fill="#3f2b2b" stroke="none"/>`)
	switch r.Status {
	case "passing":
		d.add(`<path d="M147 132 q8 -10 16 1 M198 132 q8 -10 16 -1" fill="none" stroke-width="3"/><path d="M160 160 Q179 184 200 159 L194 159 191 167 184 161 177 169 170 162 166 166Z" fill="#472d2b"/><path d="m164 161 6 2 -3 6Z m25 2 6 -3 -1 9Z" fill="#fff1d4" stroke="none"/>`)
	case "failing":
		d.add(`<path d="m145 121 18 -5 m34 -1 18 5" fill="none" stroke-width="3"/><path d="M147 132 q8 -6 16 2 l-3 9 -10 -2Z M198 134 q8 -7 16 -3 l-2 10 -11 2Z" fill="#fff1d4"/><path d="m156 133 0 6 m49 -6 0 6" stroke-width="3"/><path d="M161 170 Q180 156 199 169 l-5 7 -8 -5 -8 4 -8 -5 -6 5Z" fill="#472d2b"/><path d="m165 169 5 1 -3 5 m23 -5 5 -1 -2 6" fill="#fff1d4" stroke="none"/>`)
	case "running":
		d.add(`<path d="M145 126 164 134 158 142 150 138Z M197 134 216 126 211 138 203 142Z" fill="#fff1d4"/><path d="m156 132 0 6 m49 -6 0 6 M146 123 l18 6 m33 0 18 -6" fill="none" stroke-width="3"/><path d="m162 165 17 3 20 -4" fill="none"/><path d="m168 166 3 6 3 -5 m14 0 3 4 3 -6" fill="#fff1d4" stroke-width="1.5"/>`)
	case "queued":
		d.add(`<path d="M147 133 q8 8 16 1 M198 134 q8 8 16 -1" fill="none" stroke-width="3"/><path d="M165 164 Q180 172 196 163" fill="none" stroke-width="2"/><path d="m168 166 3 6 3 -4 m14 -1 3 4 3 -7" fill="#fff1d4" stroke-width="1.5"/>`)
	default:
		d.add(`<path d="M146 131 q8 -5 16 2 l-2 9 -12 -1Z M197 127 q9 -8 19 -2 l-2 15 -14 2Z" fill="#fff1d4"/><path d="m155 132 0 6 m52 -11 0 10 M146 120 161 120 M197 115 213 110" fill="none" stroke-width="3"/><path d="M164 166 q16 5 31 -3" fill="none"/><path d="m170 167 3 6 3 -5" fill="#fff1d4" stroke-width="1.5"/>`)
	}
	d.add(`<path d="M171 151 Q180 147 189 151 L181 158Z" fill="#3f2b2b" stroke-width="1.5"/><path d="m180 158 0 5" fill="none" stroke-width="2"/>`)
	if r.Worried {
		d.add(`<path d="M220 133 q-7 12 0 14 q7 -2 0 -14Z" fill="#a0d8df" stroke-width="1.5"/><path d="m137 113 -5 -6 m11 3 -1 -8" fill="none" stroke-width="2"/>`)
	}
	d.add(`</g>`)
}
