package mascot

// Tengu is a mountain lookout: angular red face and nose, a small cap, broad
// crow wings, knotted cord, high sandals, and a fan that changes with its task.
func tengu(d *drawing, r Resolved) {
	red, feather, robe, paper := "#ce624c", "#35434f", "#d6d2bb", "#edb36f"
	if r.Theme == "dark" {
		red, feather = "#e47c64", "#738492"
	}
	if r.Theme == "neon" {
		red, feather, robe, paper = "#ea8193", "#81739f", "#dad0df", "#e6ca8a"
	}
	d.add(`<g stroke="#302f38" stroke-width="2.8"><path d="m162 210 -5 22 -11 1 -2 7 29 0 3 -28 M190 210 l3 22 20 1 3 7 -29 0 -9 -28" fill="#b9a07a"/><path d="m150 230 20 0 m23 0 16 0 m-56 10 0 4 m13 -4 0 4 m29 -4 0 4 m12 -4 0 4" fill="none" stroke-width="2"/>`)
	// Wings change the entire outline rather than relying only on eye changes.
	d.add(`<g data-part="crow-wings" fill="%s">`, feather)
	switch r.Status {
	case "passing":
		d.add(`<path d="M151 174 Q126 159 111 128 L103 133 101 114 Q86 124 90 148 Q92 166 119 192 L108 197 Q131 210 151 192 M207 173 Q231 149 243 117 L250 129 258 121 Q270 146 243 172 L252 176 Q240 196 216 195Z"/><path d="m103 143 26 31 m-17 -8 13 17 m119 -43 -20 31 m20 -15 -15 22" fill="none" stroke="#74848b" stroke-width="2"/>`)
	case "failing":
		d.add(`<path d="M150 173 Q123 170 112 195 L102 207 115 207 110 217 Q140 220 154 194 M210 174 Q239 179 246 201 L257 213 243 211 246 221 Q219 221 205 196Z"/><path d="m121 199 19 -12 m93 15 -16 -16" fill="none" stroke="#74848b" stroke-width="2"/>`)
	case "running":
		d.add(`<path d="M150 172 Q116 168 107 185 L98 195 110 197 104 204 Q130 216 155 193 M208 174 Q231 170 242 143 L255 136 252 154 263 150 Q262 183 219 195Z"/><path d="m113 191 25 -10 m105 -25 -20 24" fill="none" stroke="#74848b" stroke-width="2"/>`)
	case "queued":
		d.add(`<path d="M150 172 Q128 179 130 211 L139 227 148 219 155 225 171 195 M205 173 Q230 183 224 211 L215 226 208 218 200 225 185 195Z"/><path d="m140 201 6 15 m68 -15 -6 15" fill="none" stroke="#74848b" stroke-width="2"/>`)
	default:
		d.add(`<path d="M150 173 Q121 170 108 155 L98 160 107 176 98 180 Q116 204 153 194 M208 173 Q236 159 234 140 L242 134 249 145 258 141 Q264 177 217 195Z"/><path d="m112 174 22 11 m111 -29 -18 23" fill="none" stroke="#74848b" stroke-width="2"/>`)
	}
	d.add(`</g><path d="M154 161 204 161 212 218 Q182 233 145 217Z" fill="%s"/><path d="m158 166 21 30 20 -29 -4 43 -34 2Z" fill="#f5e7c5"/><path d="m180 196 -2 21 m-29 -2 57 0" fill="none" stroke="#9d9e91" stroke-width="2"/><path d="M149 166 Q176 184 207 165" fill="none" stroke="#bb8a51" stroke-width="5"/>`, robe)
	for _, pt := range [][2]int{{154, 171}, {171, 178}, {190, 177}, {205, 170}} {
		d.add(`<circle cx="%d" cy="%d" r="5" fill="#f1e7cb" stroke-width="1.8"/>`, pt[0], pt[1])
	}
	// Each fan pose has a distinct job: triumphant, drooping, checking, tucked, inspecting.
	x, y, angle := 229, 184, -25
	switch r.Status {
	case "passing":
		x, y, angle = 244, 143, 18
	case "failing":
		x, y, angle = 222, 213, 135
	case "running":
		x, y, angle = 139, 193, -65
	case "queued":
		x, y, angle = 181, 205, 80
	}
	d.add(`<g data-part="maple-fan" transform="translate(%d %d) rotate(%d)"><path d="M0 14 V-4" fill="none" stroke="#a88058" stroke-width="5"/><path d="M0 1 Q-17 -2 -21 -15 L-11 -13 -13 -27 -3 -20 1 -35 7 -21 18 -28 15 -14 25 -16 Q20 -1 0 1Z" fill="%s"/><path d="M0 0 1 -25 M0 -1 -13 -16 M1 -2 16 -17" fill="none" stroke="#9a5749" stroke-width="1.6"/></g>`, x, y, angle, paper)
	angle = map[string]int{"passing": -8, "failing": 7, "running": 5, "queued": 13, "unknown": -12}[r.Status]
	d.add(`<g transform="rotate(%d 178 131)"><path d="M132 136 119 142 125 125 115 123 130 110 125 101 145 98 150 83 166 90 181 76 192 89 212 84 214 101 235 106 226 118 237 132 224 137 229 151 208 156 146 159Z" fill="%s"/><path d="M137 115 Q143 93 178 96 Q211 97 218 120 L213 150 Q207 174 179 179 Q148 172 140 150Z" fill="%s"/><path d="m144 156 13 4 -3 11 14 -5 10 13 9 -12 15 5 -2 -13 14 -5 -4 18 -18 10 -22 0 -23 -15Z" fill="%s" stroke="none"/><path d="M149 101 Q177 89 207 104" fill="none" stroke="#f1a17b" stroke-width="2"/>`, angle, feather, red, feather)
	tenguFace(d, r)
	// The long nose is drawn after the eyes to maintain a clear three-quarter face.
	d.add(`<g data-part="long-nose"><path d="M174 136 Q185 128 194 135 L225 141 Q233 144 227 149 L190 149 Q179 153 174 146Z" fill="%s"/><path d="m192 140 28 4" fill="none" stroke="#f4ac84" stroke-width="2"/><path d="m177 155 q9 4 16 -1" fill="none" stroke-width="2"/></g><g data-part="mountain-cap"><path d="m162 95 3 -20 25 -3 9 19 -19 9Z" fill="#41424c"/><path d="m165 75 13 10 12 -13 M178 85 l2 15" fill="none" stroke="#d3b077" stroke-width="1.8"/><path d="m166 95 -11 14 m39 -15 13 14" fill="none" stroke="#bb8a51" stroke-width="2"/></g></g></g>`, red)
}

func tenguFace(d *drawing, r Resolved) {
	d.add(`<g data-part="tengu-expression">`)
	switch r.Status {
	case "passing":
		d.add(`<path d="M145 130 q8 -13 17 -2 M186 127 q8 -13 17 0" fill="none" stroke-width="3.5"/><path d="M160 158 Q170 174 182 158" fill="#382f39"/><path d="m167 163 9 0" stroke="#e6b799" stroke-width="2"/>`)
	case "failing":
		d.add(`<path d="m144 116 17 -5 m24 -1 18 6" fill="none" stroke-width="4"/><path d="M145 127 q8 -5 17 0 l-2 10 -13 -1Z M186 125 q8 -5 16 1 l-2 9 -12 0Z" fill="#f8e8c7"/><path d="m154 128 0 6 m40 -6 0 5" stroke-width="3.5"/><path d="m159 164 q10 -10 21 -3" fill="none"/>`)
	case "running":
		d.add(`<path d="m144 119 19 7 m22 -1 19 -7" fill="none" stroke-width="5"/><path d="m146 129 16 0 -4 8 -10 -1Z m41 -1 14 -1 -3 8 -9 -1Z" fill="#fff0cc"/><path d="m155 130 0 4 m40 -5 0 4 m-34 30 18 0" fill="none" stroke-width="3"/>`)
	case "queued":
		d.add(`<path d="M144 130 q9 10 18 0 M186 129 q9 10 18 0" fill="none"/><ellipse cx="170" cy="162" rx="4" ry="5" fill="#382f39" stroke="none"/>`)
	default:
		d.add(`<path d="m144 114 17 2 m25 -5 17 -7" fill="none" stroke-width="4"/><path d="M145 126 q8 -5 16 0 v11 h-14Z M185 124 q9 -10 19 0 v11 h-17Z" fill="#fff0cc"/><path d="m155 128 0 6 m42 -8 0 7 M163 164 q8 -5 14 -1" fill="none" stroke-width="3"/>`)
	}
	if r.Worried {
		d.add(`<path d="M211 119 q-7 11 0 13 q7 -2 0 -13Z" fill="#9bd2d6" stroke-width="1.5"/><path d="m134 109 -5 -6 m10 3 -2 -7" fill="none" stroke-width="2"/>`)
	}
	d.add(`</g>`)
}
