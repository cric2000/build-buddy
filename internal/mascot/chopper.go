package mascot

// Chopper's oversized soft hat and round face sit above a tiny seated body.
// His black button eyes have a pale rim, unlike Kappa's unrimmed glossy eyes.
func chopper(d *drawing, r Resolved) {
	fur, cream, pink, hatShade, hoof := "#d6b27c", "#f4d8a3", "#ea94b0", "#c57494", "#493b36"
	antler := "#936e48"
	if r.Theme == "dark" || r.Theme == "neon" {
		fur, cream, pink, hatShade, antler = "#e5c18c", "#ffe4b3", "#f3aac8", "#cf84aa", "#bb9365"
	}
	d.add(`<g stroke="#574234" stroke-width="2"><path d="M156 170 Q180 157 204 170 Q219 190 211 211 L149 211 Q140 189 156 170Z" fill="%s"/><ellipse cx="180" cy="191" rx="20" ry="23" fill="%s" stroke="none"/><path d="M148 201 Q180 211 212 201 L215 219 Q200 231 180 222 Q159 231 145 219Z" fill="#82556e"/><path d="M150 207 Q180 217 210 207 M180 214 v8" fill="none" stroke="#b88b9d" stroke-width="1.2"/>`, fur, cream)
	d.add(`<ellipse cx="151" cy="220" rx="17" ry="12" transform="rotate(-20 151 220)" fill="%s"/><ellipse cx="209" cy="220" rx="17" ry="12" transform="rotate(20 209 220)" fill="%s"/><ellipse cx="149" cy="228" rx="14" ry="9" fill="%s"/><ellipse cx="211" cy="228" rx="14" ry="9" fill="%s"/><path d="m149 222 0 14 m62 -14 0 14" stroke="#a58c73" stroke-width="1.2"/>`, fur, fur, hoof, hoof)
	switch r.Status {
	case "failing":
		d.add(`<path d="M156 178 Q136 178 130 154 L139 148 Q144 163 165 166 M204 178 Q224 178 230 154 L221 148 Q216 163 195 166" fill="%s"/>`, fur)
	case "running":
		d.add(`<g data-part="medical-chart" transform="rotate(-7 180 204)"><rect x="154" y="187" width="52" height="37" rx="4" fill="#b8906c"/><rect x="158" y="191" width="44" height="29" rx="2" fill="#fff2ce" stroke="none"/><rect x="170" y="184" width="20" height="7" rx="2" fill="#8d9a94"/><path d="m164 199 20 0 m-20 7 14 0 m9 7 4 3 6 -8" fill="none" stroke="#ae9675" stroke-width="1.5"/></g><path d="M156 178 Q140 188 153 206 L164 207 l3 -8 -12 -7 M204 178 Q220 188 207 206 L196 207 l-3 -8 12 -7" fill="%s"/><path d="m156 198 10 1 -1 10 -11 -3Z m39 1 10 -1 1 8 -11 3Z" fill="%s"/>`, fur, hoof)
	case "queued":
		d.add(`<path d="M156 178 Q141 194 158 209 L180 213 l4 -8 -22 -10 M204 178 Q219 194 202 209 L181 213 l-4 -8 21 -10" fill="%s"/><path d="m172 204 15 -1 1 9 -15 1Z" fill="%s"/>`, fur, hoof)
	default:
		d.add(`<path d="M156 178 Q145 181 149 200 L160 217 Q169 218 171 210 L162 192 M204 178 Q215 181 211 200 L200 217 Q191 218 189 210 L198 192" fill="%s"/><path d="m158 207 12 -1 1 10 -9 4 -6 -6Z m44 0 -12 -1 -1 10 9 4 6 -6Z" fill="%s"/><path d="m163 211 1 7 m32 -7 -1 7" stroke="#a58c73" stroke-width="1.2"/>`, fur, hoof)
	}
	angle := map[string]int{"passing": -4, "failing": 3, "running": 2, "queued": 7, "unknown": -6}[r.Status]
	d.add(`<g transform="rotate(%d 180 140)"><g data-part="antlers" fill="none" stroke="%s" stroke-width="8"><path d="M143 114 Q117 104 123 76 Q123 70 119 67 M122 91 Q109 92 104 83 M123 102 Q105 106 99 98 M217 114 Q243 104 237 76 Q237 70 241 67 M238 91 Q251 92 256 83 M237 102 Q255 106 261 98"/><path d="M143 114 Q119 103 124 80 M217 114 Q241 103 236 80" stroke="#d3b28b" stroke-width="1.5"/></g><path d="M139 134 Q117 118 109 130 Q112 144 135 146 M221 134 Q243 118 251 130 Q248 144 225 146" fill="%s"/><path d="M132 140 Q135 117 180 119 Q226 118 229 142 Q234 176 180 183 Q127 178 132 140Z" fill="%s"/><ellipse cx="180" cy="162" rx="33" ry="17" fill="%s" stroke="none"/>`, angle, antler, fur, fur, cream)
	chopperFace(d, r, fur)
	d.add(`<g data-part="pink-hat"><path d="M138 121 139 83 Q180 65 221 83 L223 121Z" fill="%s"/><ellipse cx="180" cy="83" rx="41" ry="10" fill="%s"/><path d="M139 112 Q180 122 222 112 L222 122 139 122Z" fill="%s" stroke="none"/><path d="M118 119 Q180 130 243 119 L239 133 Q180 148 122 133Z" fill="%s"/><path d="M130 124 Q180 137 231 124" fill="none" stroke="#ffd8e5" stroke-width="1.7"/><path d="m171 87 18 18 m0 -18 -18 18" fill="none" stroke="#fffaf1" stroke-width="6"/><path d="m146 89 -1 13" stroke="#ffd8e5" stroke-width="2"/></g>`, pink, pink, hatShade, pink)
	if r.Status == "failing" {
		d.add(`<path d="m122 143 9 -6 9 7 -5 10 -10 -1Z m116 0 -9 -6 -9 7 5 10 10 -1Z" fill="%s"/>`, hoof)
	}
	d.add(`</g></g>`)
}

func chopperFace(d *drawing, r Resolved, fur string) {
	d.add(`<g data-part="chopper-expression"><ellipse cx="141" cy="163" rx="5" ry="2.5" fill="#e7a496" stroke="none"/><ellipse cx="219" cy="163" rx="5" ry="2.5" fill="#e7a496" stroke="none"/>`)
	expressiveEyes(d, r, eyeStyle{left: 156, right: 204, y: 152, width: 8.5, height: 10, ink: "#1c2020", brow: "#654f38", light: "#80847a", lid: fur, shape: "button", rim: "#fff9e8"})
	d.add(`<path d="M173 157 Q181 152 189 157 Q187 164 181 164 Q175 164 173 157Z" fill="#6597ba" stroke-width="1.5"/><path d="m177 157 5 -1" stroke="#d0e6f2" stroke-width="1.2"/>`)
	switch r.Status {
	case "passing":
		d.add(`<path d="M170 167 q5 7 11 0 q6 7 11 0" fill="none" stroke-width="1.5"/>`)
	case "failing":
		d.add(`<path d="M172 173 q9 -7 18 0" fill="none" stroke-width="1.6"/>`)
	case "running":
		d.add(`<path d="m175 171 12 0" fill="none" stroke-width="1.5"/>`)
	case "queued":
		d.add(`<ellipse cx="181" cy="171" rx="3" ry="3.5" fill="#654f38" stroke="none"/>`)
	default:
		d.add(`<path d="M176 172 q5 -4 10 -1" fill="none" stroke-width="1.5"/>`)
	}
	d.add(`</g>`)
}
