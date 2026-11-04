package adapt

import (
	"encoding/xml"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/onemindos/omos/services/entity/internal/identity"
	"github.com/onemindos/omos/services/entity/internal/model"
)

// cotEvent is the subset of CoT XML the entity service reads (tak-protocols.md §1).
type cotEvent struct {
	XMLName      xml.Name `xml:"event"`
	Type         string   `xml:"type,attr"`
	UID          string   `xml:"uid,attr"`
	How          string   `xml:"how,attr"`
	Time         string   `xml:"time,attr"`
	Start        string   `xml:"start,attr"`
	Stale        string   `xml:"stale,attr"`
	Access       string   `xml:"access,attr"`
	Caveat       string   `xml:"caveat,attr"`
	ReleasableTo string   `xml:"releasableTo,attr"`
	Point        struct {
		Lat string `xml:"lat,attr"`
		Lon string `xml:"lon,attr"`
		HAE string `xml:"hae,attr"`
		CE  string `xml:"ce,attr"`
		LE  string `xml:"le,attr"`
	} `xml:"point"`
	Detail struct {
		Contact *struct {
			Callsign string `xml:"callsign,attr"`
		} `xml:"contact"`
		Track *struct {
			Speed  string `xml:"speed,attr"`
			Course string `xml:"course,attr"`
		} `xml:"track"`
		Takv *struct {
			Device   string `xml:"device,attr"`
			Platform string `xml:"platform,attr"`
		} `xml:"takv"`
		Group *struct {
			Name string `xml:"name,attr"`
			Role string `xml:"role,attr"`
		} `xml:"__group"`
		Status *struct {
			Battery string `xml:"battery,attr"`
		} `xml:"status"`
		Emergency *struct{} `xml:"emergency"`
		FlowTags  *struct {
			Attrs []xml.Attr `xml:",any,attr"`
		} `xml:"_flow-tags_"`
		Remarks string `xml:"remarks"`
	} `xml:"detail"`
}

var (
	reICAO  = regexp.MustCompile(`(?i)\bICAO:([0-9a-f]{6})\b`)
	reNORAD = regexp.MustCompile(`(?i)\bNORAD:(\d{1,9})\b`)
	reMMSI  = regexp.MustCompile(`(?i)\bMMSI:(\d{9})\b`)
	reHex6  = regexp.MustCompile(`^[0-9a-fA-F]{6}$`)
	reDigit = regexp.MustCompile(`^\d{1,9}$`)
)

// ProjectorUIDPrefix marks the projector's own output on `tak.cot.inject`
// (tak-projector.yaml `uid = f"OM.{entity_type}.{entity_id}"`); entity.md §5 "Echo filter".
const ProjectorUIDPrefix = "OM."

// IntelUIDPrefix marks the projector's interim intel CoT (`INTEL.{id}`, cot-projection.md §8.2
// and the `intel.>` row of its source table); it is dropped as an echo too.
const IntelUIDPrefix = "INTEL."

// unknownCE is the CoT "unknown" error value; real feeds send 9999999 (feed-scripts.yaml).
const unknownCE = 999999

func parseCoTTime(s string) (int64, bool) {
	if s == "" {
		return 0, false
	}
	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		return 0, false
	}
	return t.UnixMilli(), true
}

func optFloat(s string) *float64 {
	if s == "" {
		return nil
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil
	}
	return &v
}

// cotPath derives the integration, marking and TAK egress from the subject (cot-projection.md
// §8.1 table "Marking and flags on TAK-sourced entities").
func cotPath(subject, uid string) (integration, tlp string) {
	switch subject {
	case "tak.cot.out":
		return "tak-server", "AMBER+STRICT"
	case "tak.cot.mesh.rx":
		return "tak-mesh", "AMBER+STRICT"
	}
	// tak.cot.inject and anything else under tak.cot.>: CoT the feeds put on the TAK wire.
	switch {
	case strings.HasPrefix(uid, "ADSB-"):
		return "feed-aircraft", "CLEAR"
	case strings.HasPrefix(uid, "SAT-"):
		return "feed-satellites", "CLEAR"
	case strings.HasPrefix(uid, "MTA-"):
		return "feed-mta", "CLEAR"
	case strings.HasPrefix(uid, "AIS-"), strings.HasPrefix(uid, "MMSI-"):
		return "feed-ships", "CLEAR"
	}
	return "tak-inject", "AMBER+STRICT"
}

// atom splits a CoT atom type `a-{affiliation}-{battle dimension}-...`.
func atom(t string) (aff, dim, rest string, ok bool) {
	parts := strings.Split(t, "-")
	if len(parts) < 3 || parts[0] != "a" {
		return "", "", "", false
	}
	return parts[1], parts[2], strings.Join(parts[3:], "-"), true
}

var affToDisposition = map[string]string{
	"f": "friendly", "a": "assumed_friendly", "n": "neutral", "u": "unknown",
	"s": "suspicious", "p": "pending", "j": "suspicious", "k": "suspicious", "o": "unknown",
	"h": "hostile", // never published: the engine rejects persons and downgrades assets (SL9)
}

var dimToEnv = map[string]string{"A": "air", "S": "surface", "U": "subsurface", "G": "land", "P": "space"}

// CoT parses one CoT XML event into an observation.
func CoT(data []byte, subject string) (*model.Observation, string, error) {
	trim := strings.TrimSpace(string(data))
	if !strings.HasPrefix(trim, "<") {
		return nil, "", fmt.Errorf("unsupported CoT encoding (not XML)")
	}
	var ev cotEvent
	if err := xml.Unmarshal([]byte(trim), &ev); err != nil {
		return nil, "", fmt.Errorf("cot xml: %w", err)
	}
	if ev.UID == "" {
		return nil, "", fmt.Errorf("cot event without uid")
	}
	// Echo filter (cot-projection.md §8.2): the projector's own output (uid "OM." or
	// "INTEL.", or its _flow-tags_ marker).
	if strings.HasPrefix(ev.UID, ProjectorUIDPrefix) || strings.HasPrefix(ev.UID, IntelUIDPrefix) {
		return nil, SkipEcho, nil
	}
	if ev.Detail.FlowTags != nil {
		for _, a := range ev.Detail.FlowTags.Attrs {
			if a.Name.Local == "omos-tak-projector" {
				return nil, SkipEcho, nil
			}
		}
	}
	aff, dim, rest, isAtom := atom(ev.Type)
	if !isAtom {
		return nil, SkipUnsupported, nil // b-*, u-d-*, t-*: markers, shapes and tasks come later
	}
	srcT, ok := parseCoTTime(ev.Time)
	if !ok {
		return nil, "", fmt.Errorf("cot time %q", ev.Time)
	}
	lat, err1 := strconv.ParseFloat(ev.Point.Lat, 64)
	lon, err2 := strconv.ParseFloat(ev.Point.Lon, 64)
	if err1 != nil || err2 != nil || lat < -90 || lat > 90 || lon < -180 || lon > 180 {
		return nil, "", fmt.Errorf("cot point lat=%q lon=%q", ev.Point.Lat, ev.Point.Lon)
	}

	integ, tlp := cotPath(subject, ev.UID)
	o := &model.Observation{
		Integration: integ, SourceID: ev.UID, DataType: "cot", Subject: subject,
		CotType: ev.Type, How: ev.How, HasPos: true, Lat: lat, Lon: lon,
		HAE: optFloat(ev.Point.HAE), SourceT: srcT, Marking: tlp,
		Disposition: affToDisposition[aff], Environment: dimToEnv[dim],
		Components: []string{"location"},
	}
	if ce := optFloat(ev.Point.CE); ce != nil && *ce > 0 && *ce < unknownCE {
		o.CE = *ce
	}
	if le := optFloat(ev.Point.LE); le != nil && *le > 0 && *le < unknownCE {
		o.LE = *le
	}
	if st, ok := parseCoTTime(ev.Stale); ok && st > srcT {
		o.StaleS = int((st - srcT + 999) / 1000)
	}
	if ev.Detail.Contact != nil {
		o.Callsign = ev.Detail.Contact.Callsign
	}
	if tr := ev.Detail.Track; tr != nil {
		o.Speed, o.Course = optFloat(tr.Speed), optFloat(tr.Course)
		if strings.HasPrefix(ev.UID, "ADSB-") && o.Speed != nil {
			v := *o.Speed * 0.514444 // feed-aircraft writes knots (deployment.yaml spd_kt)
			o.Speed = &v
		}
	}
	if ev.Detail.Emergency != nil {
		o.Emergency = true
	}
	if ev.Access != "" || ev.Caveat != "" || ev.ReleasableTo != "" {
		o.CotSecurity = map[string]any{"access": ev.Access, "caveat": ev.Caveat, "releasableTo": ev.ReleasableTo}
		if !strings.EqualFold(ev.Access, "unclassified") {
			o.Marking = "AMBER+STRICT"
		}
	}
	if ev.Detail.Status != nil && ev.Detail.Status.Battery != "" {
		if b := optFloat(ev.Detail.Status.Battery); b != nil {
			o.Health = map[string]any{"battery_pct": *b}
		}
	}
	if strings.HasPrefix(ev.How, "h-") {
		o.Tier = model.TierHuman
	} else {
		o.Tier = model.TierSelfReport
	}
	egressFalse := false
	o.EgressTAK = &egressFalse // already on the TAK wire (cot-projection.md §3.1 G5)

	uidKey := model.Alias{Type: identity.KeyCotUID, ID: ev.UID}
	low := strings.ToLower(ev.UID)
	switch {
	case strings.HasPrefix(low, "sim-"):
		o.Keys = []model.Alias{{Type: identity.KeySim, ID: low}, uidKey}
		o.Simulated = true
		o.EntityType, o.Kind, o.DocType, o.Template = "drone", "multirotor", "fleet.drone", "asset"
		o.Traits = []string{"locatable", "trackable", "taskable", "observable"}
		if dim == "G" {
			o.EntityType, o.Kind, o.DocType = "robot", "", "fleet.robot"
			o.Traits = []string{"locatable", "trackable", "taskable"}
		}
		t := true
		o.EgressTAK = &t
		o.Marking = "AMBER+STRICT"
	case matchKey(ev.UID, "ADSB-", reHex6) != "" || reICAO.MatchString(ev.Detail.Remarks):
		hex := matchKey(ev.UID, "ADSB-", reHex6)
		if hex == "" {
			hex = reICAO.FindStringSubmatch(ev.Detail.Remarks)[1]
		}
		o.Keys = []model.Alias{{Type: identity.KeyICAO, ID: strings.ToLower(hex)}, uidKey}
		o.EntityType, o.DocType, o.Template = "aircraft", "air.unknown", "asset"
		o.Traits = []string{"locatable", "trackable"}
	case matchKey(ev.UID, "SAT-", reDigit) != "" || reNORAD.MatchString(ev.Detail.Remarks):
		n := matchKey(ev.UID, "SAT-", reDigit)
		if n == "" {
			n = reNORAD.FindStringSubmatch(ev.Detail.Remarks)[1]
		}
		o.Keys = []model.Alias{{Type: identity.KeyNORAD, ID: n}, uidKey}
		// No `satellite` token in entity.md §3: a vehicle of kind satellite in space (ONT-P8).
		o.EntityType, o.Kind, o.DocType, o.Template = "vehicle", "satellite", "onemind.vehicle", "asset"
		o.Environment = "space"
		o.Traits = []string{"locatable", "trackable"}
	case reMMSI.MatchString(ev.Detail.Remarks):
		o.Keys = []model.Alias{{Type: identity.KeyMMSI, ID: reMMSI.FindStringSubmatch(ev.Detail.Remarks)[1]}, uidKey}
		o.EntityType, o.DocType, o.Template = "vessel", "sea.unknown", "asset"
		o.Traits = []string{"locatable", "trackable"}
	case strings.HasPrefix(ev.UID, "MTA-"):
		// MTA-MNR-{trip} / MTA-LIRR-{trip} (redpanda-connect feed-mta): gtfs-{agency}-{vehicle}.
		agency, trip, _ := strings.Cut(strings.TrimPrefix(ev.UID, "MTA-"), "-")
		o.Keys = []model.Alias{{Type: identity.KeyGTFS, ID: strings.ToLower(agency) + "-" + trip}, uidKey}
		o.EntityType, o.Kind, o.DocType, o.Template = "vehicle", "train", "onemind.vehicle", "asset"
		o.Traits = []string{"locatable", "trackable"}
	case ev.Detail.Takv != nil || (dim == "G" && strings.HasPrefix(rest, "U-C") && ev.Detail.Group != nil):
		// A TAK end-user device: a consenting team member (ONTOLOGY §3.3 person).
		o.Keys = []model.Alias{uidKey}
		if o.Callsign != "" {
			o.Keys = append(o.Keys, model.Alias{Type: identity.KeyTakCallsign, ID: o.Callsign})
		}
		o.Person = true
		o.EntityType, o.DocType, o.Template = "person", "tak.atak", "asset"
		o.Traits = []string{"locatable"}
		o.Marking = "AMBER+STRICT"
		if ev.Detail.Group != nil {
			o.Status = map[string]any{"role": ev.Detail.Group.Role, "activity": "", "mode": ""}
			o.Meta = map[string]any{"team": ev.Detail.Group.Name}
		}
	default:
		o.Keys = []model.Alias{uidKey}
		o.Template, o.Traits = "track", []string{"locatable", "trackable"}
		switch {
		case dim == "A":
			o.EntityType, o.DocType = "aircraft", "air.unknown"
		case dim == "S" || dim == "U":
			o.EntityType, o.DocType = "vessel", "sea.unknown"
		case dim == "G" && strings.HasPrefix(rest, "E-V"):
			o.EntityType, o.DocType = "vehicle", "onemind.vehicle"
		default:
			o.EntityType, o.DocType = "track", "onemind.track"
		}
	}
	if strings.HasPrefix(ev.How, "m-s") {
		o.Simulated = true
	}
	return o, "", nil
}

// matchKey returns the part of uid after prefix if it matches re, else "".
func matchKey(uid, prefix string, re *regexp.Regexp) string {
	if !strings.HasPrefix(uid, prefix) {
		return ""
	}
	if s := strings.TrimPrefix(uid, prefix); re.MatchString(s) {
		return s
	}
	return ""
}
