package mascot

// Chopper uses his classic pink hat and red shorts. The antler forks, blue nose,
// split hooves, and cross are drawn as vectors, with no fetched image assets.
func chopper(d *drawing, r Resolved) {
	antler, antlerLight := "#765138", "#c29a66"
	fur, cream, pink, hatShade, hoof := "#bd8e5a", "#f1d4a0", "#db83a5", "#b65f84", "#674438"
	if r.Theme == "dark" {
		fur, cream, pink, hatShade = "#d0a36e", "#f8deb1", "#ec9ebc", "#c87799"
		antler, antlerLight, hoof = "#ad875e", "#e1bf8c", "#9a7259"
	} else if r.Theme == "neon" {
		fur, cream, pink, hatShade = "#d4a774", "#f8deb1", "#f29fc5", "#c77fa9"
		antler, antlerLight, hoof = "#b28d6d", "#e2c4a1", "#a07865"
	}
	d.add(`<g stroke="#49362f" stroke-width="2.6"><path d="m156 211 -5 17 12 5 9 -19 m16 0 9 19 13 -6 -7 -16" fill="%s"/><path d="M150 224 Q143 228 146 237 L168 237 l2 -10 M193 227 l1 10 23 0 Q220 225 207 224" fill="%s"/><path d="m158 229 0 8 m47 -8 0 8" fill="none" stroke="#c19d7a" stroke-width="1.6"/><path d="M156 167 Q180 155 205 169 L210 207 196 219 162 220 149 204Z" fill="%s"/><path d="M166 173 Q180 168 195 175 L197 200 164 201Z" fill="%s" stroke="none"/><path d="M150 198 Q180 206 209 198 L210 221 190 225 180 216 170 225 148 220Z" fill="#b94f51"/><path d="m153 203 53 0 m-26 1 0 12" fill="none" stroke="#eaa28b" stroke-width="1.5"/>`, fur, hoof, fur, cream)
	// His little doctor's bag grounds the build-repair theme without replacing
	// the character's recognizable clothes with a generic costume.
	if r.Status == "queued" || r.Status == "unknown" {
		d.add(`<g data-part="doctor-bag"><path d="M204 201 v-7 q9 -8 17 0 v7" fill="none"/><rect x="199" y="199" width="29" height="26" rx="6" fill="#977258"/><path d="m204 205 19 0" stroke="#cbab83" stroke-width="1.5"/><path d="m212 208 3 0 0 4 4 0 0 3 -4 0 0 4 -3 0 0 -4 -4 0 0 -3 4 0Z" fill="#f6dbb1" stroke="none"/></g>`)
	}
	switch r.Status {
	case "passing":
		d.add(`<path d="M155 178 Q127 176 115 159 L106 152 l7 -8 11 6 Q139 163 164 163 M203 178 Q232 176 245 159 L254 152 l-7 -8 -11 6 Q221 163 197 163" fill="%s"/><path d="m104 151 9 -9 12 9 -10 10Z m132 0 11 -9 10 9 -11 10Z" fill="%s"/><path d="m109 147 10 8 m133 -8 -10 8" stroke="#c19d7a" stroke-width="1.5"/>`, fur, hoof)
	case "failing":
		d.add(`<path d="M155 178 Q129 175 128 145 L134 130 l11 5 -6 16 23 11 M203 178 Q230 174 232 145 L226 130 l-11 5 6 16 -23 11" fill="%s"/><path d="m130 132 7 -9 10 7 -5 12 -13 -1Z m88 -2 8 -7 7 9 -1 9 -12 1Z" fill="%s"/><path d="m136 127 7 7 m81 -7 -3 8" stroke="#c19d7a" stroke-width="1.5"/>`, fur, hoof)
	case "running":
		d.add(`<g data-part="medical-chart" transform="rotate(-8 170 200)"><path d="m143 180 53 0 0 43 -53 0Z" fill="#bb956d"/><path d="m148 185 43 0 0 32 -43 0Z" fill="#fff0cb" stroke="none"/><rect x="158" y="177" width="22" height="8" rx="3" fill="#657c88" stroke-width="1.5"/><path d="m154 195 28 0 m-28 7 15 0 m-15 8 8 0 m13 -5 4 4 7 -9" fill="none" stroke="#9f8b6e" stroke-width="1.8"/></g><path d="M155 175 Q133 179 137 199 L150 205 l6 -8 -9 -7 M203 175 Q223 189 204 207 L190 202 l4 -10 12 0" fill="%s"/><path d="m139 193 13 5 -2 10 -14 -7Z m54 1 11 4 -5 11 -12 -7Z" fill="%s"/><path d="m144 197 -3 6 m56 -6 -4 6" stroke="#c19d7a" stroke-width="1.5"/>`, fur, hoof)
	case "queued":
		d.add(`<path d="M154 178 Q132 196 156 214 L180 213 l3 -11 -22 -1 M203 178 Q222 195 202 213 L178 217 l-5 -10 24 -9" fill="%s"/><path d="m172 200 12 0 2 14 -13 0Z m0 8 13 -1" fill="%s"/>`, fur, hoof)
	default:
		d.add(`<path d="M155 178 Q133 184 140 209 L156 213 l4 -10 -12 -4 M203 176 Q228 169 225 149 L217 144 l-7 9 5 9 -19 4" fill="%s"/><path d="m143 201 16 1 0 13 -15 -3Z m69 -54 9 -5 7 10 -9 8Z" fill="%s"/><path d="m151 204 -1 9 m67 -66 6 9" stroke="#c19d7a" stroke-width="1.5"/>`, fur, hoof)
	}
	angle := map[string]int{"passing": -7, "failing": 4, "running": 4, "queued": 11, "unknown": -10}[r.Status]
	d.add(`<g transform="rotate(%d 180 141)"><g data-part="antlers" fill="none" stroke="%s" stroke-width="9"><path d="M141 115 Q117 103 117 80 L109 70 M118 92 102 89 95 78 M120 100 99 102 90 95 M219 115 Q243 103 243 80 L251 70 M242 92 258 89 265 78 M240 100 261 102 270 95"/><path d="M140 115 Q118 101 118 80 M220 115 Q242 101 242 80" stroke="%s" stroke-width="2"/></g><path d="M139 128 Q113 111 106 123 Q109 142 137 145 M221 128 Q247 111 254 123 Q251 142 223 145" fill="%s"/><path d="m116 125 19 11 m109 -11 -19 11" stroke="#f1d4a0" stroke-width="3"/><path d="M134 132 Q137 113 180 115 Q224 113 227 137 L229 152 Q223 179 180 183 Q137 180 131 155Z" fill="%s"/><path d="M146 153 Q153 143 168 149 Q180 155 192 148 Q210 144 218 159 Q209 178 181 178 Q153 178 141 164Z" fill="%s" stroke="none"/>`, angle, antler, antlerLight, fur, fur, cream)
	chopperFace(d, r)
	// The brim sits in front of the antlers; the white X is a filled polygon.
	d.add(`<g data-part="pink-hat"><path d="M142 119 142 84 Q180 68 219 84 L219 119Z" fill="%s"/><ellipse cx="180" cy="83" rx="38" ry="10" fill="%s"/><path d="M142 111 Q180 122 219 111 L219 122 142 122Z" fill="%s" stroke="none"/><path d="M125 120 Q180 132 235 120 L231 133 Q180 148 129 133Z" fill="%s"/><path d="M139 125 Q179 136 221 125" fill="none" stroke="#f6c4d8" stroke-width="2"/><path d="m169 87 11 11 11 -11 6 6 -11 11 11 11 -6 6 -11 -11 -11 11 -6 -6 11 -11 -11 -11Z" fill="#fff9ee" stroke="none"/><path d="M149 89 v14" stroke="#f6c4d8" stroke-width="3"/></g>`, pink, pink, hatShade, pink)
	if r.Status == "failing" {
		d.add(`<path d="m121 129 9 -7 11 7 -5 13 -14 -1Z m98 0 11 -7 9 7 -1 12 -14 1Z" fill="%s"/><path d="m128 126 8 9 m96 -9 -8 9" fill="none" stroke="#c19d7a" stroke-width="1.5"/>`, hoof)
	}
	d.add(`</g></g>`)
}

func chopperFace(d *drawing, r Resolved) {
	d.add(`<g data-part="chopper-expression"><ellipse cx="143" cy="159" rx="7" ry="3" fill="#de9e86" stroke="none"/><ellipse cx="216" cy="159" rx="7" ry="3" fill="#de9e86" stroke="none"/>`)
	switch r.Status {
	case "passing":
		d.add(`<path d="M148 148 q7 -12 14 0 M198 148 q7 -12 14 0" fill="none" stroke-width="3"/><path d="M169 162 Q180 179 193 162Z" fill="#593b35"/><path d="m174 169 q7 -4 13 0" stroke="#e9a28e" stroke-width="3"/>`)
	case "failing":
		d.add(`<ellipse cx="155" cy="149" rx="8" ry="10" fill="#fff8e5"/><ellipse cx="205" cy="149" rx="8" ry="10" fill="#fff8e5"/><path d="m155 146 0 6 m50 -6 0 6 m-57 -16 14 -4 m36 0 14 4" fill="none" stroke-width="3"/><path d="M170 173 Q180 157 192 172 l-6 3 -5 -3 -5 4Z" fill="#593b35"/>`)
	case "running":
		d.add(`<path d="M147 142 163 146 M197 146 213 142" fill="none" stroke-width="3"/><ellipse cx="156" cy="152" rx="4" ry="6" fill="#49362f" stroke="none"/><ellipse cx="204" cy="152" rx="4" ry="6" fill="#49362f" stroke="none"/><path d="m174 169 14 0" fill="none" stroke-width="2"/>`)
	case "queued":
		d.add(`<path d="M148 150 q7 8 14 0 M198 150 q7 8 14 0" fill="none"/><ellipse cx="181" cy="169" rx="4" ry="5" fill="#593b35" stroke="none"/>`)
	default:
		d.add(`<ellipse cx="154" cy="150" rx="5" ry="7" fill="#49362f" stroke="none"/><ellipse cx="206" cy="147" rx="6" ry="9" fill="#49362f" stroke="none"/><circle cx="155" cy="148" r="1.5" fill="#fff8e5" stroke="none"/><circle cx="207" cy="144" r="2" fill="#fff8e5" stroke="none"/><path d="m147 136 14 1 m37 -5 14 -4 M176 170 q5 -4 10 -1" fill="none" stroke-width="2"/>`)
	}
	d.add(`<g data-part="blue-nose"><path d="M172 153 Q181 148 190 153 Q188 162 181 162 Q175 162 172 153Z" fill="#5898b6" stroke-width="1.8"/><path d="m176 154 5 -1" stroke="#b6e0e9" stroke-width="1.5"/></g>`)
	if r.Worried {
		d.add(`<path d="M222 145 q-6 11 0 13 q6 -2 0 -13Z" fill="#a0d8df" stroke-width="1.5"/><path d="m140 139 -5 -5 m9 2 0 -7" fill="none" stroke-width="2"/>`)
	}
	d.add(`</g>`)
}
