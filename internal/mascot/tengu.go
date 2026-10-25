package mascot

import "math"

// Tengu has a soft red face, cloudlike pale hair, a compact robe, small crow
// wings, and a maple fan. Rounded sleeves and hakama give it a seated silhouette.
func tengu(d *drawing, r Resolved) {
	red, wing, vest, hair := "#dc786c", "#53657a", "#4d4c55", "#f9efd6"
	if r.Theme == "dark" || r.Theme == "neon" {
		red, wing, vest, hair = "#ed9589", "#8493a4", "#727280", "#fff4df"
	}
	d.add(`<g stroke="#56413a" stroke-width="2"><g data-part="crow-wings" fill="%s"><path d="M149 171 Q125 159 115 170 Q104 177 116 183 Q104 194 120 195 Q122 209 139 199 L159 186Z"/><path transform="translate(360 0) scale(-1 1)" d="M149 171 Q125 159 115 170 Q104 177 116 183 Q104 194 120 195 Q122 209 139 199 L159 186Z"/><path d="m119 183 19 3 m-16 10 17 -4 m102 -9 -19 3 m16 10 -17 -4" fill="none" stroke="#afbac5" stroke-width="1.2"/></g>`, wing)
	d.add(`<path d="M157 195 Q140 208 149 224 Q160 234 179 226 Q198 234 212 224 Q222 208 202 195Z" fill="#9481a2"/><path d="m180 201 -1 24 M158 206 q-5 8 -1 16 m45 -16 q5 8 1 16" fill="none" stroke="#c0acc9" stroke-width="1.3"/><path d="M151 225 q12 -4 20 2 l-1 9 -24 0Z M189 227 q9 -6 20 -2 l5 11 -24 0Z" fill="#ac8d6b"/><path d="m150 231 17 0 m26 0 17 0" stroke="#f0dab1" stroke-width="1.5"/><path d="M153 161 Q180 151 207 162 L216 202 Q180 216 143 202Z" fill="#f0ddb4"/><path d="M154 161 173 176 168 205 144 201Z M205 162 188 176 193 205 216 201Z" fill="%s"/><path d="m173 175 7 14 9 -14 m-9 14 -1 17" fill="none" stroke="#a99474" stroke-width="1.4"/>`, vest)
	d.add(`<path d="M154 165 Q132 165 130 188 L145 201 161 186 M206 165 Q228 165 230 188 L215 201 199 186" fill="#f0ddb4"/><path d="M137 187 q-6 0 -7 6 q1 7 9 7 l7 -8 M223 187 q6 0 7 6 q-1 7 -9 7 l-7 -8" fill="%s"/>`, red)
	for _, pt := range [][2]int{{158, 174}, {169, 183}, {202, 174}, {191, 183}} {
		d.add(`<circle cx="%d" cy="%d" r="4.2" fill="#fff7e5" stroke-width="1.1"/>`, pt[0], pt[1])
	}
	x, y, fanAngle := 232, 183, -15
	switch r.Status {
	case "passing":
		x, y, fanAngle = 238, 158, 10
	case "failing":
		x, y, fanAngle = 231, 207, 110
	case "running":
		x, y, fanAngle = 220, 192, -45
	case "queued":
		x, y, fanAngle = 180, 210, 75
	}
	gripX := float64(x) - 12*math.Sin(float64(fanAngle)*math.Pi/180)
	gripY := float64(y) + 12*math.Cos(float64(fanAngle)*math.Pi/180)
	if r.Status != "queued" {
		d.add(`<path d="M223 193 Q238 191 %.1f %.1f" fill="none" stroke="%s" stroke-width="7"/><ellipse cx="%.1f" cy="%.1f" rx="4" ry="4.5" fill="%s"/>`, gripX, gripY, red, gripX, gripY, red)
	}
	d.add(`<g data-part="maple-fan" transform="translate(%d %d) rotate(%d)"><path d="M0 13 V-2" stroke="#a17a52" stroke-width="4"/><path d="M0 1 Q-17 -2 -19 -14 L-10 -12 -11 -24 -3 -18 1 -31 7 -19 17 -24 14 -12 23 -14 Q19 -1 0 1Z" fill="#e5a56f"/><path d="M0 0 1 -23 M0 -1 -11 -14 M1 -2 14 -15" fill="none" stroke="#ad6a4f" stroke-width="1.1"/></g>`, x, y, fanAngle)
	angle := map[string]int{"passing": -4, "failing": 3, "running": 0, "queued": 7, "unknown": -6}[r.Status]
	d.add(`<g transform="rotate(%d 180 130)"><path d="M137 111 Q122 104 130 94 Q125 85 144 84 Q149 73 168 78 Q183 72 196 79 Q218 76 222 88 Q240 90 230 105 Q243 117 230 130 Q242 145 223 152 Q222 166 204 160 L155 160 Q137 167 137 152 Q117 143 128 129 Q115 119 137 111Z" fill="%s"/><ellipse cx="144" cy="130" rx="7" ry="10" fill="%s"/><ellipse cx="216" cy="130" rx="7" ry="10" fill="%s"/><path d="M145 119 Q147 94 180 94 Q213 95 215 121 L213 143 Q208 167 180 174 Q151 167 147 145Z" fill="%s"/>`, angle, hair, red, red, red)
	tenguFace(d, r, red)
	d.add(`<g data-part="long-nose"><path d="M177 139 Q185 132 193 138 L223 141 Q233 142 229 148 Q226 153 192 150 Q182 153 177 147Z" fill="%s"/><path d="m192 141 27 3" stroke="#ffd0ad" stroke-width="1.4"/></g><g data-part="mountain-cap"><path d="m166 92 2 -17 23 -1 6 18 -16 6Z" fill="%s"/><path d="m169 76 10 8 11 -9 m-11 9 2 13" fill="none" stroke="#d4b888" stroke-width="1.1"/></g></g></g>`, red, vest)
}

func tenguFace(d *drawing, r Resolved, red string) {
	d.add(`<g data-part="tengu-expression"><ellipse cx="153" cy="148" rx="4" ry="2.4" fill="#f4baa8" stroke="none"/>`)
	expressiveEyes(d, r, eyeStyle{left: 159, right: 195, y: 128, width: 7.5, height: 10, ink: "#3b2829", brow: "#71443f", iris: "#713746", light: "#d895a0", lid: red})
	switch r.Status {
	case "passing":
		d.add(`<path d="M165 156 q8 8 16 0" fill="none" stroke-width="1.7"/>`)
	case "failing":
		d.add(`<path d="M165 161 q8 -7 16 -1" fill="none" stroke-width="1.6"/>`)
	case "running":
		d.add(`<path d="m168 158 11 0" fill="none" stroke-width="1.6"/>`)
	case "queued":
		d.add(`<ellipse cx="173" cy="158" rx="2.5" ry="3" fill="#71443f" stroke="none"/>`)
	default:
		d.add(`<path d="M169 159 q5 -4 10 -1" fill="none" stroke-width="1.5"/>`)
	}
	d.add(`</g>`)
}
