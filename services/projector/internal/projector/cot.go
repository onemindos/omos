package projector

// Entity to CoT (onemind-xr docs/contracts/cot-projection.md §4 and §5, robotics-cot-profile.md
// §2). Every `a-` type below is one the contract checked against MITRE CoTtypes.xml.

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	cotTimeFmt  = "2006-01-02T15:04:05.000Z"
	unknownVal  = 9999999.0
	xmlDecl     = `<?xml version="1.0" encoding="UTF-8" standalone="yes" ?>` + "\n"
	flowTagAttr = "omos-tak-projector"
)

// Reasons a record gets no CoT type (counted, cot-projection.md §4.1 rules 4 and 5).
const (
	NoTypeRow     = "no_type_row"
	NoTypeNever   = "type_never"
	NoTypeUnsafe  = "type_forbidden"
	NoTypeNoShape = "zone_shape_unknown"
)

var cotTypeRe = regexp.MustCompile(`^[abtuy](-[A-Za-z0-9.]+)+$`)

// affiliation letters (§4.2). There is no hostile disposition in OneMind; a record that asks
// for one is projected as suspect (civilian check, rule 4).
var dispositionLetter = map[string]string{
	"friendly": "f", "assumed_friendly": "a", "neutral": "n", "unknown": "u",
	"suspicious": "s", "suspect": "s", "pending": "p", "hostile": "s", "joker": "s", "faker": "s",
}

var dimLetter = map[string]string{"air": "A", "surface": "S", "subsurface": "U", "land": "G", "space": "P"}

// civilian applies rule 4: affiliation h/j/k becomes s; strike, engage, all tasking (t-*)
// and weapon atoms are refused. downgraded reports a rewrite (counted).
func civilian(t string) (out string, ok, downgraded bool) {
	if strings.HasPrefix(t, "t-") || strings.HasPrefix(t, "y-") {
		return "", false, false
	}
	parts := strings.Split(t, "-")
	if parts[0] == "a" && len(parts) >= 3 {
		switch parts[1] {
		case "h", "j", "k":
			parts[1] = "s"
			downgraded = true
		}
		if len(parts) >= 5 && ((parts[2] == "G" && parts[3] == "E" && parts[4] == "W") || (parts[2] == "U" && parts[3] == "W")) {
			return "", false, false // weapon atoms a-.-G-E-W-*, a-.-U-W-*
		}
		if len(parts) == 4 && parts[2] == "U" && parts[3] == "W" {
			return "", false, false
		}
		t = strings.Join(parts, "-")
	}
	return t, true, downgraded
}

func (e *Entity) affiliation(def string) string {
	if l, ok := dispositionLetter[strings.ToLower(e.Disposition)]; ok {
		return l
	}
	return def
}

func (e *Entity) dim() string {
	if d, ok := dimLetter[strings.ToLower(e.Environment)]; ok {
		return d
	}
	return "X"
}

// eventAtoms is the §4.1 `event` table; matched against kind and the source's event name.
var eventAtoms = []struct {
	keys []string
	atom string
}{
	{[]string{"earthquake", "quake"}, "a-n-X-i-g-e"},
	{[]string{"hotspot", "hot_spot", "fire_hot"}, "a-n-X-i-f-h"},
	{[]string{"wildfire", "red flag", "red_flag", "fire"}, "a-n-X-i-f-w"},
	{[]string{"tsunami"}, "a-n-X-i-m-n"},
	{[]string{"flood"}, "a-n-X-i-m-f"},
	{[]string{"tornado"}, "a-n-X-i-m-t"},
	{[]string{"hurricane", "cyclone", "tropical", "typhoon"}, "a-n-X-i-m-c"},
	{[]string{"thunderstorm"}, "a-n-X-i-m-u"},
	{[]string{"snow", "winter", "blizzard", "ice storm", "ice_storm"}, "a-n-X-i-m-s"},
	{[]string{"hail"}, "a-n-X-i-m-h"},
	{[]string{"fog"}, "a-n-X-i-m-g"},
	{[]string{"drought"}, "a-n-X-i-m-d"},
	{[]string{"volcano", "volcanic"}, "a-n-X-i-g-v"},
	{[]string{"landslide", "mudslide"}, "a-n-X-i-g-l"},
	{[]string{"avalanche"}, "a-n-X-i-g-a"},
	{[]string{"pollution", "spill"}, "a-n-X-i-e"},
	{[]string{"public_health", "public health", "outbreak", "disease"}, "a-n-X-i-h"},
	{[]string{"transport_accident", "accident", "crash"}, "a-n-X-i-t"},
	{[]string{"conflict", "security", "acled", "protest", "riot"}, "a-n-X-i-l"},
}

var vesselAtoms = map[string]string{
	"cargo": "S-X-M-C", "tanker": "S-X-M-O", "passenger": "S-X-M-P", "ferry": "S-X-M-F",
	"fishing": "S-X-F", "tug": "S-X-M-T-U", "towing": "S-X-M-T-O", "leisure": "S-X-R",
	"pleasure": "S-X-R", "sailing": "S-X-R", "hazmat": "S-X-M-H",
}

// zoneTypes maps geo.shape to the TAK drawing type (§5.5).
var zoneTypes = map[string]string{
	"circle": "u-d-c-c", "polygon": "u-d-f", "rectangle": "u-d-r", "line": "b-m-r", "corridor": "u-d-f",
}

// tableType builds the type from §4.1 and the robotics profile; "" with a reason when none.
func tableType(e *Entity) (string, string) {
	kind := strings.ToLower(e.Kind)
	switch e.EntityType {
	case "person", "team", "agent", "node", "station", "alert":
		return "", NoTypeNever
	case "drone":
		if !e.IsOwnAsset() {
			return "", NoTypeRow
		}
		if kind == "multirotor" {
			return "a-f-A-M-H-Q", ""
		}
		return "a-f-A-M-F-Q", "" // fixed_wing, vtol, unknown airframe
	case "robot":
		if !e.IsOwnAsset() {
			return "", NoTypeRow
		}
		switch strings.ToLower(e.Environment) {
		case "surface":
			return "a-f-S-X", ""
		case "subsurface":
			return "a-f-U-S-U", ""
		}
		return "a-f-G-E-V", ""
	case "vehicle":
		if !e.IsOwnAsset() {
			return "", NoTypeRow
		}
		if strings.ToLower(e.Environment) == "surface" {
			return "a-f-S-X", ""
		}
		return "a-f-G-E-V-C", ""
	case "aircraft":
		a := e.affiliation("n")
		switch kind {
		case "fixed_wing", "airliner", "jet", "light", "small", "large", "heavy", "turboprop":
			return "a-" + a + "-A-C-F", ""
		case "rotary", "rotorcraft", "helicopter":
			return "a-" + a + "-A-C-H", ""
		case "balloon", "blimp", "airship":
			return "a-" + a + "-A-C-L", ""
		}
		return "a-" + a + "-A", ""
	case "vessel":
		a := e.affiliation("n")
		if s, ok := vesselAtoms[kind]; ok {
			return "a-" + a + "-" + s, ""
		}
		return "a-" + a + "-S-X", ""
	case "track":
		if e.IsPersonClass() {
			return "", NoTypeNever
		}
		return "a-" + e.affiliation("u") + "-" + e.dim(), ""
	case "sensor":
		if kind == "emplaced" {
			return "a-" + e.affiliation("n") + "-G-E-S-E", ""
		}
		return "a-" + e.affiliation("n") + "-G-E-S", ""
	case "camera":
		return "a-" + e.affiliation("n") + "-G-E-S-E", ""
	case "structure":
		return "a-" + e.affiliation("n") + "-G-I-c", ""
	case "place":
		return "b-m-p-s-m", ""
	case "zone":
		if e.Geo == nil {
			return "", NoTypeNoShape
		}
		if t, ok := zoneTypes[strings.ToLower(e.Geo.Shape)]; ok {
			return t, ""
		}
		return "", NoTypeNoShape
	case "event":
		return eventAtom(e), ""
	case "satellite":
		return "a-n-P", ""
	}
	return "", NoTypeRow
}

func eventAtom(e *Entity) string {
	hay := strings.ToLower(e.Kind)
	if ev, ok := e.Meta["event"].(string); ok { // NWS CAP event name, e.g. "Tornado Warning"
		hay += " " + strings.ToLower(ev)
	}
	return matchEvent(hay)
}

func matchEvent(hay string) string {
	for _, row := range eventAtoms {
		for _, k := range row.keys {
			if strings.Contains(hay, k) {
				return row.atom
			}
		}
	}
	return "a-n-X-i"
}

// CotType returns the CoT type for an entity. downgraded is true when the civilian check
// rewrote a hostile-looking affiliation; reason is set when the entity has no type.
func CotType(e *Entity) (t string, downgraded bool, reason string) {
	if base, why := tableType(e); why != "" {
		return "", false, why
	} else if e.CotType != "" && cotTypeRe.MatchString(e.CotType) && !e.IsOwnAsset() && e.EntityType != "zone" {
		// Rule 1: a precomputed type from the entity service, when it passes the civilian check.
		if out, ok, down := civilian(e.CotType); ok {
			return out, down, ""
		}
		downgraded = true // refused the precomputed type; fall back to the table
		t = base
	} else {
		t = base
	}
	out, ok, down := civilian(t)
	if !ok {
		return "", false, NoTypeUnsafe
	}
	switch strings.ToLower(e.Disposition) {
	case "hostile", "joker", "faker": // asked for an adversary affiliation; sent as suspect
		down = true
	}
	return out, downgraded || down, ""
}

// staleDefault is the §5.3 type default (seconds) when the entity has no stale.
func staleDefault(e *Entity) int {
	if e.IsOwnAsset() {
		return 60
	}
	switch e.EntityType {
	case "aircraft":
		return 60
	case "vessel", "track":
		return 120
	case "sensor", "camera":
		return 600
	case "event", "satellite":
		return 3600
	case "place", "zone", "structure":
		return 86400
	}
	return 120
}

// StaleSeconds applies the §5.3 clamp.
func StaleSeconds(e *Entity, interval time.Duration) int {
	s := e.Stale
	if s <= 0 {
		s = staleDefault(e)
	}
	lower := int(math.Ceil(2*interval.Seconds())) + 5
	upper := 3600
	if e.IsPinnedType() {
		upper = 86400
	}
	if s < lower {
		s = lower
	}
	if s > upper {
		s = upper
	}
	return s
}

var allowedHow = map[string]bool{"m-g": true, "m-f": true, "m-s": true, "m-p": true, "m-a": true, "h-e": true, "h-c": true}

func (e *Entity) how() string {
	if e.Indicators.Simulated {
		return "m-s" // SL7: simulated entities carry how="m-s" (JADC2-MAP §5.2)
	}
	if allowedHow[e.Provenance.How] {
		return e.Provenance.How
	}
	if len(e.Sources) > 1 {
		return "m-f"
	}
	return "m-p"
}

// xb is a tiny XML writer with escaped attributes.
type xb struct{ bytes.Buffer }

func (b *xb) attr(k, v string) {
	b.WriteByte(' ')
	b.WriteString(k)
	b.WriteString(`="`)
	_ = xml.EscapeText(&b.Buffer, []byte(v))
	b.WriteByte('"')
}

func ff(v float64, prec int) string { return strconv.FormatFloat(v, 'f', prec, 64) }

func fnum(v float64) string {
	if v >= 999999 || math.IsNaN(v) || math.IsInf(v, 0) {
		return ff(unknownVal, 1)
	}
	return ff(v, 1)
}

// Built is one CoT ready to publish.
type Built struct {
	UID  string
	Type string
	XML  []byte
}

type buildOpts struct {
	target      Target
	interval    time.Duration
	maxVertices int
	maxBytes    int
	removal     bool
}

// zone colours (§5.5): stroke, fill (0 = none), as signed 32-bit ARGB.
var zoneColours = map[string][2]int32{
	"aoi":          {argb(0xFF, 0xFF, 0xFF, 0x00), 0},
	"keep_in":      {argb(0xFF, 0x00, 0xFF, 0x00), 0},
	"keep_out":     {argb(0xFF, 0xFF, 0x00, 0x00), argb(0x40, 0xFF, 0x00, 0x00)},
	"hazard":       {argb(0xFF, 0xFF, 0xA5, 0x00), argb(0x40, 0xFF, 0xA5, 0x00)},
	"control_area": {argb(0xFF, 0x00, 0x00, 0xFF), 0},
	"footprint":    {argb(0xFF, 0x00, 0xFF, 0xFF), argb(0x26, 0x00, 0xFF, 0xFF)},
	"route":        {argb(0xFF, 0xFF, 0xFF, 0xFF), 0},
}

func argb(a, r, g, b uint32) int32 { return int32(a<<24 | r<<16 | g<<8 | b) }

// Build renders the CoT event for an entity (or its removal, §5.3).
func Build(e *Entity, cotType string, now time.Time, o buildOpts) (Built, error) {
	lat, lon, ok := e.Position()
	if !ok {
		return Built{}, fmt.Errorf("no position")
	}
	staleS := StaleSeconds(e, o.interval)
	stale := now.Add(time.Duration(staleS) * time.Second)
	if o.removal {
		stale = now.Add(time.Second)
	}
	start := now.Add(-time.Duration(staleS) * time.Second)
	if et := time.UnixMilli(e.T); e.T > 0 && et.After(start) && !et.After(now) {
		start = et
	}
	if o.removal {
		start = now
	}
	nowS := now.UTC().Format(cotTimeFmt)
	uid := e.UID()

	shape := shapePoints(e, o.maxVertices)
	if e.EntityType == "zone" && len(shape) > 0 {
		if s := strings.ToLower(e.Geo.Shape); s == "line" {
			lat, lon = shape[0][1], shape[0][0] // route: the first vertex
		} else if s != "circle" {
			lat, lon, _ = centroid(shape)
		}
	}
	remarks := remarksText(e, o.removal)
	for attempt := 0; ; attempt++ {
		var b xb
		b.WriteString(xmlDecl)
		b.WriteString(`<event`)
		b.attr("version", "2.0")
		b.attr("uid", uid)
		b.attr("type", cotType)
		b.attr("how", e.how())
		b.attr("time", nowS)
		b.attr("start", start.UTC().Format(cotTimeFmt))
		b.attr("stale", stale.UTC().Format(cotTimeFmt))
		if e.Indicators.Simulated {
			b.attr("opex", "s")
		} else if e.Indicators.Exercise {
			b.attr("opex", "e")
		}
		b.WriteString(`><point`)
		b.attr("lat", ff(lat, 6))
		b.attr("lon", ff(lon, 6))
		hae, ce, le := unknownVal, unknownVal, unknownVal
		if a, ok := e.Altitude(); ok {
			hae = a
		}
		if e.Uncertainty != nil {
			if e.Uncertainty.CEM > 0 {
				ce = e.Uncertainty.CEM
			}
			if e.Uncertainty.LEM > 0 {
				le = e.Uncertainty.LEM
			}
		}
		b.attr("hae", fnum(hae))
		b.attr("ce", fnum(ce))
		b.attr("le", fnum(le))
		b.WriteString(`/><detail><contact`)
		b.attr("callsign", e.callsignText())
		b.WriteString(`/>`)
		if e.hasTrait("trackable") {
			if c, s, ok := e.motion(); ok {
				b.WriteString(`<track`)
				b.attr("course", ff(c, 1))
				b.attr("speed", ff(s, 2))
				b.WriteString(`/>`)
			}
		}
		if bp, ok := e.Health["battery_pct"].(float64); ok {
			b.WriteString(`<status`)
			b.attr("battery", strconv.Itoa(int(math.Round(bp))))
			b.WriteString(`/>`)
		}
		if remarks != "" {
			r := remarks
			if o.target == TargetMesh && len(r) > 80 {
				r = r[:80]
			}
			b.WriteString(`<remarks>`)
			_ = xml.EscapeText(&b.Buffer, []byte(r))
			b.WriteString(`</remarks>`)
		}
		if e.EntityType == "zone" && !o.removal {
			writeShape(&b, e, uid, shape)
		}
		b.WriteString(`<_flow-tags_`)
		b.attr(flowTagAttr, nowS)
		b.WriteString(`/>`)
		if o.target == TargetServer {
			b.WriteString(`<onemind`)
			b.attr("v", "1")
			b.attr("id", e.EntityID)
			b.attr("type", e.EntityType)
			b.attr("seq", strconv.FormatUint(e.Seq, 10))
			b.attr("conf", ff(e.Confidence, 2))
			b.attr("lc", e.Lifecycle)
			b.attr("tlp", strings.ToUpper(e.Marking.TLP))
			sim := "0"
			if e.Indicators.Simulated {
				sim = "1"
			}
			b.attr("sim", sim)
			b.WriteString(`/>`)
		}
		b.WriteString(`</detail></event>`)
		if o.target != TargetMesh || o.maxBytes <= 0 || b.Len() <= o.maxBytes {
			return Built{UID: uid, Type: cotType, XML: b.Bytes()}, nil
		}
		// Mesh size cap (§6.2): fewer vertices first, then no remarks.
		switch {
		case len(shape) > 4:
			shape = simplify(shape, len(shape)/2)
		case remarks != "":
			remarks = ""
		default:
			return Built{}, fmt.Errorf("cot is %d bytes, over MESH_MAX_BYTES %d", b.Len(), o.maxBytes)
		}
		_ = attempt
	}
}

func (e *Entity) callsignText() string {
	c := e.Callsign
	if c == "" && e.Aliases.Name != nil {
		c = *e.Aliases.Name
	}
	if c == "" {
		c = e.EntityID
	}
	if r := []rune(c); len(r) > 64 {
		c = string(r[:64])
	}
	return c
}

func (e *Entity) motion() (course, speed float64, ok bool) {
	if e.Location != nil && e.Location.CourseDeg != nil && e.Location.SpeedMps != nil {
		return *e.Location.CourseDeg, *e.Location.SpeedMps, true
	}
	if e.Course != nil && e.Speed != nil {
		return *e.Course, *e.Speed, true
	}
	return 0, 0, false
}

func remarksText(e *Entity, removal bool) string {
	if removal {
		return "removed"
	}
	k := e.Kind
	if k == "" {
		k = e.EntityType
	}
	s := fmt.Sprintf("OneMind %s · conf %.2f · %s", k, e.Confidence, e.Lifecycle)
	if e.Indicators.Simulated {
		s += " · SIMULATED"
	}
	seen := map[string]bool{}
	var src []string
	for _, x := range e.Sources {
		if x.Integration != "" && !seen[x.Integration] {
			seen[x.Integration] = true
			src = append(src, x.Integration)
		}
	}
	if len(src) > 0 {
		s += " · " + strings.Join(src, ",")
	}
	return s
}

// shapePoints returns the zone vertices to draw (corridor buffered, simplified to max).
func shapePoints(e *Entity, max int) []pt {
	if e.EntityType != "zone" || e.Geo == nil {
		return nil
	}
	ps := geoPoints(e.Geo)
	switch strings.ToLower(e.Geo.Shape) {
	case "corridor":
		if len(ps) >= 2 && e.Geo.RadiusM > 0 {
			ps = bufferLine(ps, e.Geo.RadiusM)
		}
	case "polygon", "rectangle":
		if len(ps) > 0 && ps[0] != ps[len(ps)-1] {
			ps = append(ps, ps[0])
		}
	}
	return simplify(ps, max)
}

func writeShape(b *xb, e *Entity, uid string, ps []pt) {
	shape := strings.ToLower(e.Geo.Shape)
	switch shape {
	case "circle":
		r := e.Geo.RadiusM
		b.WriteString(`<shape><ellipse`) // attribute names unverified (§5.5)
		b.attr("major", ff(r, 1))
		b.attr("minor", ff(r, 1))
		b.attr("angle", "0")
		b.WriteString(`/></shape>`)
	case "line":
		for i, p := range ps {
			b.WriteString(`<link`)
			b.attr("uid", fmt.Sprintf("%s.wp%d", uid, i))
			b.attr("type", "b-m-p-w")
			b.attr("relation", "c")
			b.attr("point", ff(p[1], 6)+","+ff(p[0], 6))
			b.WriteString(`/>`)
		}
	default:
		for _, p := range ps {
			b.WriteString(`<link`)
			b.attr("point", ff(p[1], 6)+","+ff(p[0], 6))
			b.WriteString(`/>`)
		}
	}
	kind := strings.ToLower(e.Geo.Kind)
	if shape == "line" && kind == "" {
		kind = "route"
	}
	c, ok := zoneColours[kind]
	if !ok {
		c = zoneColours["control_area"]
	}
	b.WriteString(`<strokeColor`)
	b.attr("value", strconv.Itoa(int(c[0])))
	b.WriteString(`/><strokeWeight value="3.0"/>`)
	if c[1] != 0 {
		b.WriteString(`<fillColor`)
		b.attr("value", strconv.Itoa(int(c[1])))
		b.WriteString(`/>`)
	}
	b.WriteString(`<color`)
	b.attr("argb", strconv.Itoa(int(c[0])))
	b.WriteString(`/><labels_on value="true"/>`)
}
