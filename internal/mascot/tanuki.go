package mascot

// Tanuki is built from a large round head, little folded paws, a cream belly,
// and one soft striped tail. Its leaf cap keeps a woodland identity.
func tanuki(d *drawing, r Resolved) {
	fur, shade, cream, mask := "#bf955b", "#997040", "#f8dfaf", "#765333"
	if r.Theme == "dark" || r.Theme == "neon" {
		fur, shade, cream, mask = "#d5ae74", "#ae8955", "#ffe8bc", "#896747"
	}
	tail := `M204 209 C231 234 266 225 271 199 Q277 169 255 167 Q245 170 251 181 Q247 192 233 188 Q213 187 204 209Z`
	d.add(`<defs><clipPath id="tanuki-tail-cutout"><path d="%s"/></clipPath></defs><g stroke="#57412b" stroke-width="2"><g data-part="brush-tail"><path d="%s" fill="%s"/><g clip-path="url(#tanuki-tail-cutout)" fill="none" stroke="%s" stroke-width="13"><path d="M240 184 Q254 199 248 223 M222 194 Q233 205 227 226"/></g></g>`, tail, tail, fur, mask)
	d.add(`<path d="M151 161 Q180 150 208 163 Q231 181 220 214 Q210 233 180 233 Q147 233 138 214 Q128 182 151 161Z" fill="%s"/><ellipse cx="179" cy="198" rx="32" ry="30" fill="%s" stroke="none"/><ellipse cx="150" cy="223" rx="21" ry="13" transform="rotate(-15 150 223)" fill="%s"/><ellipse cx="210" cy="223" rx="21" ry="13" transform="rotate(15 210 223)" fill="%s"/><ellipse cx="150" cy="227" rx="7" ry="4" fill="%s" stroke="none"/><ellipse cx="210" cy="227" rx="7" ry="4" fill="%s" stroke="none"/>`, fur, cream, fur, fur, cream, cream)
	switch r.Status {
	case "passing":
		d.add(`<path d="M146 173 Q125 171 119 156 Q112 147 109 156 Q109 180 139 192 M211 173 Q227 179 218 196 Q209 206 194 199 L201 186" fill="%s"/>`, fur)
	case "failing":
		d.add(`<path d="M146 174 Q128 185 143 201 L168 209 l5 -8 -22 -12 M211 174 Q229 185 214 201 L189 209 l-5 -8 22 -12" fill="%s"/>`, fur)
	case "running":
		d.add(`<path d="M146 174 Q132 184 141 199 Q153 208 169 199 L164 190 151 189 M211 174 Q225 184 216 199 Q204 208 188 199 L193 190 206 189" fill="%s"/><path d="m174 201 5 4 9 -10" fill="none" stroke="#7d8b57" stroke-width="1.8"/>`, fur)
	case "queued":
		d.add(`<path d="M146 174 Q131 192 149 204 L179 208 l4 -8 -25 -11 M211 174 Q226 192 208 204 L180 208 l-4 -8 24 -11" fill="%s"/>`, fur)
	default:
		d.add(`<path d="M146 174 Q132 185 142 199 L159 207 Q169 205 165 197 L152 186 M211 174 Q225 185 215 199 L198 207 Q188 205 192 197 L205 186" fill="%s"/>`, fur)
	}
	angle := map[string]int{"passing": -3, "failing": 3, "running": 0, "queued": 7, "unknown": -6}[r.Status]
	d.add(`<g transform="rotate(%d 180 133)"><circle cx="139" cy="98" r="16" fill="%s"/><circle cx="221" cy="98" r="16" fill="%s"/><circle cx="139" cy="98" r="9" fill="%s" stroke="none"/><circle cx="221" cy="98" r="9" fill="%s" stroke="none"/><path d="M126 125 Q125 99 151 94 Q177 83 208 96 Q235 102 234 127 Q245 147 229 160 Q213 178 180 178 Q145 178 130 160 Q116 148 126 125Z" fill="%s"/><path d="M135 132 Q148 116 166 126 Q180 135 194 126 Q213 116 225 133 L225 149 Q213 161 194 151 L180 146 166 153 Q147 162 135 149Z" fill="%s" stroke="none"/><ellipse cx="180" cy="158" rx="22" ry="16" fill="%s" stroke="none"/>`, angle, shade, shade, cream, cream, fur, mask, cream)
	tanukiFace(d, r, cream, mask)
	d.add(`<g data-part="leaf-cap"><path d="M162 97 Q151 82 167 73 Q183 66 202 73 Q200 93 179 98Z" fill="#79a563"/><path d="M163 97 Q177 85 191 77 M177 86 171 79 m7 7 12 2" fill="none" stroke="#487447" stroke-width="1.3"/><path d="m162 97 -6 3" fill="none" stroke="#487447" stroke-width="2"/></g></g></g>`)
}

func tanukiFace(d *drawing, r Resolved, cream, mask string) {
	d.add(`<g data-part="tanuki-expression">`)
	expressiveEyes(d, r, eyeStyle{left: 154, right: 206, y: 138, width: 8.5, height: 10, ink: "#30251a", brow: cream, light: "#b9924b", lid: mask, shape: "button"})
	d.add(`<ellipse cx="139" cy="151" rx="5" ry="2.5" fill="#e7ae8f" stroke="none"/><ellipse cx="221" cy="151" rx="5" ry="2.5" fill="#e7ae8f" stroke="none"/><path d="M173 152 Q180 148 187 152 L181 159 Q177 160 173 152Z" fill="#3f2c20" stroke="none"/>`)
	switch r.Status {
	case "passing":
		d.add(`<path d="M171 163 q4 5 9 0 q5 5 9 0" fill="none" stroke-width="1.5"/>`)
	case "failing":
		d.add(`<path d="M172 168 q8 -6 16 0" fill="none" stroke-width="1.6"/>`)
	case "running":
		d.add(`<path d="m174 166 12 0" fill="none" stroke-width="1.5"/>`)
	case "queued":
		d.add(`<ellipse cx="180" cy="165" rx="2.5" ry="3" fill="#3f2c20" stroke="none"/>`)
	default:
		d.add(`<path d="M176 166 q4 -4 8 0" fill="none" stroke-width="1.5"/>`)
	}
	d.add(`</g>`)
}
