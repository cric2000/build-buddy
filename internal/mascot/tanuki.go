package mascot

// Tanuki is a woodland courier. Its broad tail, dark eye mask, leaf cap, and
// overstuffed satchel distinguish it from the upright river and mountain folk.
func tanuki(d *drawing, r Resolved) {
	fur, dark, cream, leaf, cloth := "#b88758", "#514337", "#f1dbaf", "#729353", "#7f999c"
	if r.Theme == "dark" {
		fur, dark, leaf = "#c99b6d", "#68584b", "#91ab6d"
	}
	if r.Theme == "neon" {
		fur, dark, cream, leaf, cloth = "#cba47a", "#65566f", "#f8e0bb", "#93c293", "#a391bf"
	}
	d.add(`<g stroke="#382f2b" stroke-width="2.8"><g data-part="brush-tail"`)
	if r.Status == "passing" {
		d.add(` transform="rotate(-16 213 211)"`)
	} else if r.Status == "failing" {
		d.add(` transform="rotate(9 213 211)"`)
	}
	d.add(`><path d="M204 212 C226 234 272 222 270 183 Q269 160 251 159 Q257 177 241 181 Q224 181 213 193Z" fill="%s"/><path d="M249 160 Q270 157 270 183 Q270 201 260 211 Q263 193 247 187 Q238 184 233 187 L241 181 Q257 177 249 160Z" fill="%s" stroke="none"/><path d="M226 189 Q245 188 252 200 M223 200 Q235 201 239 212" fill="none" stroke="%s" stroke-width="7"/></g>`, fur, dark, dark)
	if r.Status == "queued" {
		d.add(`<path d="M143 215 Q120 220 127 233 L158 235 166 221 M196 220 Q213 215 225 231 Q223 238 198 237 L184 224" fill="%s"/>`, dark)
	} else {
		d.add(`<path d="m148 214 -12 16 q-1 7 25 5 l12 -16 M191 217 l8 18 q24 5 26 -2 l-15 -19" fill="%s"/>`, dark)
	}
	d.add(`<path d="M148 160 Q122 180 131 213 Q137 233 175 231 Q218 233 228 211 Q232 179 211 162Z" fill="%s"/><path d="M159 169 Q179 155 200 173 Q216 191 207 218 Q174 233 146 213 Q143 186 159 169Z" fill="%s" stroke="none"/><path d="m149 221 4 3 m8 -1 5 2 m29 -3 5 -3" fill="none" stroke="#c6aa7d" stroke-width="1.5"/>`, fur, cream)
	// A stitched cross-body bag is part of the courier's silhouette in every pose.
	d.add(`<path d="M147 162 Q169 187 210 203" fill="none" stroke="#775846" stroke-width="9"/><path d="M147 162 Q169 187 210 203" fill="none" stroke="#dab27c" stroke-width="3"/><g data-part="satchel"><path d="m193 194 29 -1 9 28 q-17 13 -39 4Z" fill="%s"/><path d="m193 194 29 -1 4 12 -16 7 -16 -6Z" fill="#b6b8a0"/><path d="m199 219 4 2 m7 1 4 -1 m7 -3 3 -2" fill="none" stroke="#e8d5ad" stroke-width="1.5"/><rect x="207" y="204" width="6" height="8" rx="2" fill="#e7bd69"/></g>`, cloth)
	switch r.Status {
	case "passing":
		d.add(`<path d="M145 177 Q124 170 112 150 L100 149 Q95 145 101 141 L110 142 108 132 Q112 126 116 133 L123 145 155 159" fill="%s"/><path d="M212 172 Q233 170 239 148 L249 143 Q256 142 254 149 L248 156 Q246 181 219 191" fill="%s"/>`, fur, fur)
	case "failing":
		d.add(`<path d="M145 176 Q122 168 134 148 L145 142 Q155 144 150 151 L144 155 157 161 M213 173 Q232 163 218 146 L210 140 Q203 142 207 149 L213 155 200 163" fill="%s"/>`, fur)
	case "running":
		d.add(`<g data-part="delivery-scroll" transform="rotate(-8 168 195)"><path d="m138 181 56 0 0 36 -55 0Z" fill="#fff0c8"/><path d="M136 181 q-8 0 -8 7 q0 7 11 7 v22 M194 217 q9 0 9 -7 q0 -7 -9 -7" fill="#e3c493"/><path d="m148 191 31 0 m-31 7 17 0 m-17 7 25 0" stroke="#b39870" stroke-width="2"/><path d="m180 204 4 4 8 -10" fill="none" stroke="#5e8054"/></g><path d="M144 172 Q125 183 140 199 L151 197 Q154 192 146 190 M212 171 Q227 187 208 202 L195 198 Q192 190 201 191 L204 185" fill="%s"/>`, fur)
	case "queued":
		d.add(`<path d="M142 174 Q125 198 151 210 L175 204 Q179 197 171 196 L153 195 M212 174 Q231 195 207 209 L177 211 Q170 205 179 201 L201 197" fill="%s"/><path d="M173 208 Q189 211 206 205" fill="none" stroke="#e6c797" stroke-width="2"/>`, fur)
	default:
		d.add(`<path d="M144 174 Q123 187 136 206 L150 210 Q157 204 148 201 L145 196 M211 173 Q235 164 227 145 L215 142 Q209 146 218 152 L212 160 201 162" fill="%s"/>`, fur)
	}
	angle := map[string]int{"passing": 7, "failing": -5, "running": 3, "queued": -11, "unknown": 12}[r.Status]
	d.add(`<g transform="rotate(%d 178 134)"><path d="M131 114 Q107 90 119 77 Q142 69 155 100 M197 99 Q216 73 235 83 Q249 100 224 119" fill="%s"/><path d="M130 99 Q119 84 123 84 Q134 81 143 101 M209 102 Q225 85 229 91 Q236 98 222 108" fill="#c8aa85" stroke="none"/><path d="M126 122 Q124 99 155 96 Q180 85 207 100 Q231 104 234 133 L241 146 230 148 233 158 221 159 Q209 181 178 181 Q144 180 130 160 L119 157 125 149 116 143Z" fill="%s"/><path d="M134 133 Q151 115 173 132 L182 140 192 129 Q215 118 226 139 L225 155 Q205 170 188 155 L178 146 Q159 170 136 155Z" fill="%s" stroke="none"/><path d="M157 154 Q179 138 202 155 Q214 165 201 174 Q178 186 154 172 Q143 163 157 154Z" fill="%s" stroke="none"/>`, angle, dark, fur, dark, cream)
	tanukiFace(d, r, cream, dark)
	// The leaf has an asymmetric bite and a bent stem, rather than a stock icon.
	d.add(`<g data-part="leaf-cap"><path d="M155 103 Q145 81 164 67 Q180 61 203 65 L199 73 209 76 Q206 95 180 103Z" fill="%s"/><path d="M159 101 Q178 87 193 73 M176 87 166 80 M180 83 195 84" fill="none" stroke="#425e3b" stroke-width="1.8"/><path d="M157 104 q-7 3 -12 0" fill="none" stroke-width="3"/></g></g></g>`, leaf)
}

func tanukiFace(d *drawing, r Resolved, cream, dark string) {
	d.add(`<g data-part="tanuki-expression" stroke="%s">`, cream)
	switch r.Status {
	case "passing":
		d.add(`<path d="M143 143 q8 -13 16 0 M197 143 q8 -13 16 0" fill="none" stroke-width="3.5"/><path d="M165 164 Q178 180 191 164" fill="%s" stroke="%s" stroke-width="2"/><path d="m173 169 10 0" stroke="#df9a80" stroke-width="3"/>`, dark, dark)
	case "failing":
		d.add(`<path d="m140 132 20 -6 m35 0 20 6" fill="none"/><ellipse cx="153" cy="144" rx="4" ry="7" fill="%s" stroke="none"/><ellipse cx="203" cy="144" rx="4" ry="7" fill="%s" stroke="none"/><path d="M166 173 Q178 159 190 173" fill="none" stroke="%s"/>`, cream, cream, dark)
	case "running":
		d.add(`<path d="m141 136 20 5 m34 0 20 -5" fill="none" stroke-width="3.5"/><path d="m152 144 0 5 m51 -5 0 5" stroke-width="4"/><path d="m169 170 18 -1" fill="none" stroke="%s"/>`, dark)
	case "queued":
		d.add(`<path d="M141 143 q10 9 20 0 M195 143 q10 9 20 0" fill="none"/><ellipse cx="179" cy="171" rx="4" ry="5" fill="%s" stroke="none"/>`, dark)
	default:
		d.add(`<ellipse cx="151" cy="143" rx="4" ry="6" fill="%s" stroke="none"/><ellipse cx="204" cy="141" rx="6" ry="8" fill="%s" stroke="none"/><path d="m141 129 17 1 m36 -4 20 -5" fill="none"/><path d="M172 172 q7 -7 13 -2" fill="none" stroke="%s"/>`, cream, cream, dark)
	}
	d.add(`<path d="M168 155 Q178 149 188 155 L181 162 Q176 165 168 155Z" fill="%s" stroke="%s" stroke-width="2"/><path d="M178 162 v5 M146 158 l2 0 m5 2 2 0 m45 -1 2 0 m5 -3 2 0" stroke="%s" stroke-width="1.5"/>`, dark, dark, dark)
	if r.Worried {
		d.add(`<path d="M218 132 q-7 11 0 13 q7 -2 0 -13Z" fill="#a6ced1" stroke="%s" stroke-width="1.5"/><path d="m131 121 -5 -6 m10 4 -2 -8" fill="none" stroke="%s"/>`, dark, dark)
	}
	d.add(`</g>`)
}
