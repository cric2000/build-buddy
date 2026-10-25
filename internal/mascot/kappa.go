package mascot

// Kappa is a round seated river creature, with a leaf crown, water dish,
// tucked flippers, a cream belly, and broad orange webbed feet.
func kappa(d *drawing, r Resolved) {
	skin, shadow, leaf, water := "#87d996", "#53ae72", "#398952", "#65cde4"
	if r.Theme == "dark" || r.Theme == "neon" {
		skin, shadow, leaf, water = "#a0e8aa", "#68bc87", "#4b9c66", "#7ee4ed"
	}
	d.add(`<defs><radialGradient id="kappa-skin" cx=".38" cy=".32" r=".8"><stop stop-color="%s"/><stop offset="1" stop-color="%s"/></radialGradient><linearGradient id="kappa-belly" x2="0" y2="1"><stop stop-color="#fff0c8"/><stop offset="1" stop-color="#e6be85"/></linearGradient></defs>`, skin, shadow)
	d.add(`<g stroke="#234c35" stroke-width="2.2"><g data-part="shell"><ellipse cx="180" cy="180" rx="64" ry="57" fill="#a87944"/><ellipse cx="180" cy="182" rx="59" ry="55" fill="#f0cb8e"/><path d="M125 170 Q115 191 126 215 M235 170 Q245 191 234 215" fill="none" stroke="#805d38" stroke-width="1.4"/></g>`)
	d.add(`<path d="M145 145 Q180 130 215 145 C239 164 241 199 230 219 Q211 237 180 236 Q148 237 130 219 C119 196 121 164 145 145Z" fill="url(#kappa-skin)"/><path d="M159 151 Q180 142 201 151 Q218 181 214 217 Q183 233 146 217 Q143 184 159 151Z" fill="url(#kappa-belly)"/><path d="M155 169 Q180 180 207 169 M150 186 Q180 198 212 186 M148 204 Q180 217 213 204" fill="none" stroke="#b99767" stroke-width="1.2"/>`)
	// Large soft webbing sits in front of the shell and belly.
	d.add(`<g data-part="webbed-feet" fill="#e9a45c" stroke="#795d37"><path d="M158 219 Q153 208 143 204 L140 215 Q129 201 127 207 L129 222 Q116 213 114 219 Q123 237 144 238 Q160 237 158 219Z"/><path transform="translate(360 0) scale(-1 1)" d="M158 219 Q153 208 143 204 L140 215 Q129 201 127 207 L129 222 Q116 213 114 219 Q123 237 144 238 Q160 237 158 219Z"/><path d="M137 219 q8 1 13 9 M223 219 q-8 1 -13 9" fill="none" stroke="#f9cf8b" stroke-width="2"/></g>`)
	switch r.Status {
	case "passing":
		d.add(`<path d="M146 157 Q125 162 123 146 Q117 134 111 139 Q108 151 114 166 Q124 182 147 181 M213 157 Q228 170 220 193 Q216 204 200 201 Q194 199 200 190 L204 174" fill="url(#kappa-skin)"/>`)
	case "failing":
		d.add(`<path d="M146 157 Q126 162 133 191 L159 201 Q166 197 159 190 L147 180 M214 157 Q234 162 227 191 L201 201 Q194 197 201 190 L213 180" fill="url(#kappa-skin)"/>`)
	case "running":
		d.add(`<path d="M146 157 Q124 180 139 201 Q149 207 169 201 L164 189 151 188 M214 157 Q236 180 221 201 Q211 207 191 201 L196 189 209 188" fill="url(#kappa-skin)"/><path d="m172 196 6 3 11 -9" fill="none" stroke="#568258" stroke-width="2"/>`)
	case "queued":
		d.add(`<path d="M146 157 Q128 176 143 196 Q160 208 180 201 Q185 197 178 193 L155 186 M214 157 Q232 176 217 196 Q204 207 183 204 Q177 201 183 196 L205 187" fill="url(#kappa-skin)"/>`)
	default:
		d.add(`<path d="M146 157 Q128 171 137 190 Q144 204 160 202 Q169 202 165 194 L153 175 M214 157 Q232 171 223 190 Q216 204 200 202 Q191 202 195 194 L207 175" fill="url(#kappa-skin)"/>`)
	}
	angle := map[string]int{"passing": -3, "failing": 3, "running": 0, "queued": 7, "unknown": -5}[r.Status]
	d.add(`<g transform="rotate(%d 180 129)"><path d="M121 124 C121 97 143 86 180 87 C218 86 240 98 240 125 Q241 167 180 173 Q119 168 121 124Z" fill="url(#kappa-skin)"/><ellipse cx="174" cy="119" rx="35" ry="24" fill="#d1ffd0" fill-opacity=".18" stroke="none"/>`, angle)
	// Rounded crown leaves replace the old serrated outline.
	d.add(`<g data-part="leaf-crown" fill="%s"><path d="M150 91 Q134 82 126 99 Q142 104 157 96Z"/><path d="M153 94 Q135 106 145 119 Q162 114 166 98Z"/><path d="M165 96 Q162 115 180 124 Q198 114 195 96Z"/><path d="M207 94 Q225 106 215 119 Q198 114 194 98Z"/><path d="M210 91 Q226 82 234 99 Q218 104 203 96Z"/></g>`, leaf)
	d.add(`<g data-part="water-dish"><path d="M144 89 Q180 101 216 89 L210 99 Q180 110 150 99Z" fill="#c3d795"/><ellipse cx="180" cy="89" rx="35" ry="9" fill="%s"/><path d="M151 89 Q166 84 183 88 T206 88" fill="none" stroke="#efffff" stroke-width="1.7"/><ellipse cx="195" cy="86" rx="5" ry="1.5" fill="#efffff" stroke="none"/></g>`, water)
	if r.Status == "failing" {
		d.add(`<path d="m213 90 q11 -6 15 2 M229 101 q-5 8 0 10 q5 -2 0 -10Z" fill="%s" stroke="%s" stroke-width="1.2"/>`, water, water)
	}
	d.add(`<g data-part="blossom" transform="translate(127 111)" fill="#dc96c8" stroke="#694a64" stroke-width="1"><ellipse cy="-5" rx="3.5" ry="5"/><ellipse cx="5" cy="-1" rx="5" ry="3.5"/><ellipse cx="3" cy="5" rx="3.5" ry="5" transform="rotate(-20 3 5)"/><ellipse cx="-3" cy="5" rx="3.5" ry="5" transform="rotate(20 -3 5)"/><ellipse cx="-5" cy="-1" rx="5" ry="3.5"/><circle r="3" fill="#f5d675"/></g>`)
	kappaFace(d, r, skin)
	d.add(`</g></g>`)
}

func kappaFace(d *drawing, r Resolved, skin string) {
	d.add(`<g data-part="kappa-expression"><ellipse cx="143" cy="148" rx="6" ry="2.8" fill="#f2b3c6" stroke="#527556" stroke-width="1"/><ellipse cx="217" cy="148" rx="6" ry="2.8" fill="#f2b3c6" stroke="#527556" stroke-width="1"/>`)
	expressiveEyes(d, r, eyeStyle{left: 155, right: 205, y: 134, width: 8, height: 9, ink: "#12271f", brow: "#315c40", light: "#789888", lid: skin, shape: "button"})
	switch r.Status {
	case "passing":
		d.add(`<path d="M163 151 180 142 197 151 181 162Z" fill="#e9a16a" stroke="#775b3d" stroke-width="1.6"/><path d="M164 152 Q180 163 196 152" fill="none" stroke="#986244" stroke-width="1.2"/>`)
	case "failing":
		d.add(`<path d="M164 154 180 145 196 154 181 163Z" fill="#e9a16a" stroke="#775b3d" stroke-width="1.6"/><path d="M169 157 Q180 152 191 157" fill="none" stroke="#986244" stroke-width="1.2"/>`)
	case "queued":
		d.add(`<path d="M164 152 180 143 196 152 181 161Z" fill="#e9a16a" stroke="#775b3d" stroke-width="1.6"/><ellipse cx="181" cy="156" rx="3" ry="2" fill="#775b3d" stroke="none"/>`)
	default:
		d.add(`<path d="M164 152 180 143 196 152 181 161Z" fill="#e9a16a" stroke="#775b3d" stroke-width="1.6"/><path d="m166 153 28 0" stroke="#986244" stroke-width="1.2"/>`)
	}
	d.add(`</g>`)
}
