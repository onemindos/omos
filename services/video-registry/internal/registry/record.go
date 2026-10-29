package registry

import (
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/onemindos/omos/services/video-registry/internal/takvideo"
)

// Endpoints are the reader-facing addresses of the MediaMTX this registry describes. Host is
// the tailnet name of the onemind-media Tailscale Service; SiteHost (optional) is the LAN name
// of a site edge MediaMTX for mesh peers. Ports default to the MediaMTX defaults.
type Endpoints struct {
	Host     string
	SiteHost string
	RTSPPort int
	HLSPort  int
	WebRTC   int
	SRTPort  int
	Playback int
}

func (e Endpoints) withDefaults() Endpoints {
	def := func(v, d int) int {
		if v == 0 {
			return d
		}
		return v
	}
	e.RTSPPort = def(e.RTSPPort, 8554)
	e.HLSPort = def(e.HLSPort, 8888)
	e.WebRTC = def(e.WebRTC, 8889)
	e.SRTPort = def(e.SRTPort, 8890)
	e.Playback = def(e.Playback, 9996)
	return e
}

// URLs are reader URLs per protocol. None carries credentials: readers authenticate with
// their own bus identity (RTSP userinfo, HTTP Basic, SRT streamid).
type URLs struct {
	RTSP     string `json:"rtsp"`
	HLS      string `json:"hls"`
	WHEP     string `json:"whep"`
	SRT      string `json:"srt"`
	Playback string `json:"playback"`
}

// LatencyClass per transport (video-and-plugins.md "How it applies" 12).
type LatencyClass struct {
	WHEP string `json:"whep"`
	RTSP string `json:"rtsp"`
	SRT  string `json:"srt"`
	HLS  string `json:"hls"`
}

var latency = LatencyClass{WHEP: "control", RTSP: "tactical", SRT: "tactical", HLS: "situational"}

// Media is the entry the entity service merges into the entity's media[] (entity.md section 3).
type Media struct {
	Kind     string `json:"kind"`
	URL      string `json:"url"`
	Role     string `json:"role"`
	StreamID string `json:"stream_id"`
}

// Record is one value in KV OVERWATCH_VIDEO, key {entity_type}.{entity_id}.
type Record struct {
	V             int          `json:"v"`
	StreamID      string       `json:"stream_id"`
	EntityType    string       `json:"entity_type"`
	EntityID      string       `json:"entity_id"`
	EntitySubject string       `json:"entity_subject"`
	Live          bool         `json:"live"`
	Reason        string       `json:"reason,omitempty"`
	SourceType    string       `json:"source_type,omitempty"`
	Tracks        []string     `json:"tracks"`
	Codec         string       `json:"codec,omitempty"`
	Width         int          `json:"width,omitempty"`
	Height        int          `json:"height,omitempty"`
	KLV           bool         `json:"klv"`
	LatencyClass  string       `json:"latency_class"`
	Latency       LatencyClass `json:"latency_by_transport"`
	URLs          URLs         `json:"urls"`
	TAKURL        string       `json:"tak_url,omitempty"`
	SiteURL       string       `json:"site_url,omitempty"`
	Permissions   Permissions  `json:"permissions"`
	Recording     Recording    `json:"recording"`
	Media         Media        `json:"media"`
	ReadySince    *time.Time   `json:"ready_since,omitempty"`
	LastSeen      time.Time    `json:"last_seen"`
	UpdatedAt     time.Time    `json:"updated_at"`
	Origin        string       `json:"origin"`
	Seq           uint64       `json:"seq"`
}

func hostPort(h string, p int) string { return net.JoinHostPort(h, strconv.Itoa(p)) }

// BuildURLs returns the reader URLs of a stream on an endpoint set.
func BuildURLs(e Endpoints, streamID string) URLs {
	e = e.withDefaults()
	if e.Host == "" {
		return URLs{}
	}
	return URLs{
		RTSP:     "rtsp://" + hostPort(e.Host, e.RTSPPort) + "/" + streamID,
		HLS:      "http://" + hostPort(e.Host, e.HLSPort) + "/" + streamID + "/index.m3u8",
		WHEP:     "http://" + hostPort(e.Host, e.WebRTC) + "/" + streamID + "/whep",
		SRT:      "srt://" + hostPort(e.Host, e.SRTPort) + "?streamid=read:" + streamID,
		Playback: "http://" + hostPort(e.Host, e.Playback) + "/list?path=" + streamID,
	}
}

// Build makes the registry record of a live MediaMTX path.
func Build(ref StreamRef, item PathItem, cfg *MediaConfig, e Endpoints, origin string, now time.Time) Record {
	r := Record{
		V: 1, StreamID: ref.StreamID, EntityType: ref.EntityType, EntityID: ref.EntityID,
		EntitySubject: ref.EntitySubject(), Live: item.Ready, Tracks: append([]string{}, item.Tracks...),
		Latency: latency, URLs: BuildURLs(e, ref.StreamID), LastSeen: now.UTC(), UpdatedAt: now.UTC(), Origin: origin,
	}
	if r.Tracks == nil {
		r.Tracks = []string{}
	}
	if item.Source != nil {
		r.SourceType = item.Source.Type
	}
	if item.ReadyTime != nil {
		t := item.ReadyTime.UTC()
		r.ReadySince = &t
	}
	for _, t := range item.Tracks {
		if strings.EqualFold(t, "KLV") {
			r.KLV = true
		}
	}
	for _, t := range item.Tracks2 {
		if r.Codec == "" && isVideoCodec(t.Codec) {
			r.Codec = t.Codec
			r.Width = intProp(t.CodecProps["width"])
			r.Height = intProp(t.CodecProps["height"])
		}
	}
	if r.Codec == "" {
		for _, t := range item.Tracks {
			if isVideoCodec(t) {
				r.Codec = t
				break
			}
		}
	}
	// Best class a TAK or native reader gets: RTSP/SRT (tactical). WHEP (control) depends on
	// the WebRTC ICE path, which is not proven on the tailnet yet.
	r.LatencyClass = latency.RTSP
	if cfg != nil {
		r.Permissions = cfg.PermissionsFor(ref.StreamID)
		r.Recording = cfg.RecordingFor(ref.StreamID)
	}
	r.TAKURL = takvideo.Select(takvideo.Hub, r.URLs.RTSP)
	if e.SiteHost != "" {
		site := e
		site.Host = e.SiteHost
		r.SiteURL = takvideo.Select(takvideo.Site, BuildURLs(site, ref.StreamID).RTSP)
	}
	r.Media = Media{Kind: "video", URL: r.URLs.RTSP, Role: "primary", StreamID: ref.StreamID}
	return r
}

func isVideoCodec(c string) bool {
	switch strings.ToUpper(c) {
	case "H264", "H265", "AV1", "VP8", "VP9", "MPEG-4 VIDEO", "MPEG-1/2 VIDEO", "M-JPEG":
		return true
	}
	return false
}

func intProp(v any) int {
	switch n := v.(type) {
	case float64:
		return int(n)
	case int:
		return n
	}
	return 0
}

// sameContent reports whether two records differ only in bookkeeping fields.
func sameContent(a, b Record) bool {
	a.LastSeen, b.LastSeen = time.Time{}, time.Time{}
	a.UpdatedAt, b.UpdatedAt = time.Time{}, time.Time{}
	a.Seq, b.Seq = 0, 0
	if a.ReadySince != nil && b.ReadySince != nil && a.ReadySince.Equal(*b.ReadySince) {
		b.ReadySince = a.ReadySince
	}
	return equalJSON(a, b)
}
