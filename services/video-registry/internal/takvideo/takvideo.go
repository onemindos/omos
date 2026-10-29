// Package takvideo renders a stream registry entry for TAK (FEATURES K3): the CoT <__video>
// detail with a nested <ConnectionEntry> (robotics-cot-profile.md section 5) and the TAK Server
// video feed body for /Marti/api/video. It enforces the profile's URL hygiene: a TAK URL is never
// loopback, link-local, a pod or cluster address, and never carries credentials or a query.
//
// The hub projector (FEATURES A6) imports this package; until it exists the registry uses it to
// decide which URL it publishes as tak_url.
package takvideo

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
)

// Audience selects the hygiene rule set.
type Audience int

const (
	// Hub is a TAK user on the tailnet (TAK Server target): no private or cluster addresses.
	Hub Audience = iota
	// Site is a LAN mesh peer (site projector, tak.cot.mesh.tx): LAN addresses are fine.
	Site
)

var defaultPorts = map[string]int{"rtsp": 554, "rtsps": 443, "srt": 0, "rtmp": 1935, "http": 80, "https": 443}

// Check returns an error when raw is not a URL TAK users can open for the audience.
func Check(raw string, aud Audience) error {
	u, err := url.Parse(raw)
	if err != nil {
		return err
	}
	if _, ok := defaultPorts[u.Scheme]; !ok {
		return fmt.Errorf("scheme %q not usable by TAK clients", u.Scheme)
	}
	if u.User != nil || strings.Contains(raw, "@") {
		return errors.New("credentials in URL")
	}
	if u.RawQuery != "" || u.ForceQuery || u.Fragment != "" {
		return errors.New("query or fragment in URL")
	}
	host := strings.ToLower(u.Hostname())
	if host == "" {
		return errors.New("no host")
	}
	if host == "localhost" || strings.HasSuffix(host, ".localhost") {
		return errors.New("loopback host")
	}
	if aud == Hub && (strings.HasSuffix(host, ".svc") || strings.HasSuffix(host, ".cluster.local") || strings.HasSuffix(host, ".local")) {
		return errors.New("cluster or mDNS host is not reachable from the tailnet")
	}
	if ip := net.ParseIP(host); ip != nil {
		switch {
		case ip.IsLoopback(), ip.IsUnspecified(), ip.IsLinkLocalUnicast(), ip.IsLinkLocalMulticast(), ip.IsMulticast():
			return errors.New("loopback, link-local or unspecified address")
		case aud == Hub && ip.IsPrivate():
			return errors.New("private (pod or LAN) address for a hub audience")
		}
	}
	return nil
}

// Select returns the first URL that passes Check: the registry's TAK URL first, then the
// entity's media[] URLs (robotics-cot-profile.md section 5 "URL source"). "" means omit <__video>.
func Select(aud Audience, registryURL string, mediaURLs ...string) string {
	for _, c := range append([]string{registryURL}, mediaURLs...) {
		if c != "" && Check(c, aud) == nil {
			return c
		}
	}
	return ""
}

type connectionEntry struct {
	XMLName           xml.Name `xml:"ConnectionEntry"`
	UID               string   `xml:"uid,attr"`
	Alias             string   `xml:"alias,attr"`
	Protocol          string   `xml:"protocol,attr"`
	Address           string   `xml:"address,attr"`
	Port              int      `xml:"port,attr"`
	Path              string   `xml:"path,attr"`
	NetworkTimeout    int      `xml:"networkTimeout,attr"`
	BufferTime        int      `xml:"bufferTime,attr"`
	RTSPReliable      int      `xml:"rtspReliable,attr"`
	IgnoreEmbeddedKLV bool     `xml:"ignoreEmbeddedKLV,attr"`
	RoverPort         int      `xml:"roverPort,attr"`
}

type videoDetail struct {
	XMLName xml.Name        `xml:"__video"`
	UID     string          `xml:"uid,attr"`
	URL     string          `xml:"url,attr"`
	Entry   connectionEntry `xml:"ConnectionEntry"`
}

// UID is the video uid of an entity, OM.{entity_type}.{entity_id}.video.
func UID(entityType, entityID string) string { return "OM." + entityType + "." + entityID + ".video" }

// Detail renders <__video> for an entity, or "" when no URL passes the audience's hygiene.
// alias is the human label (e.g. "SIM-X1 cam0").
func Detail(aud Audience, entityType, entityID, alias string, candidates ...string) (string, error) {
	if len(candidates) == 0 {
		return "", nil
	}
	chosen := Select(aud, candidates[0], candidates[1:]...)
	if chosen == "" {
		return "", nil
	}
	u, err := url.Parse(chosen)
	if err != nil {
		return "", err
	}
	port := defaultPorts[u.Scheme]
	if p := u.Port(); p != "" {
		if port, err = strconv.Atoi(p); err != nil {
			return "", err
		}
	}
	uid := UID(entityType, entityID)
	d := videoDetail{UID: uid, URL: chosen, Entry: connectionEntry{
		UID: uid, Alias: alias, Protocol: u.Scheme, Address: u.Hostname(), Port: port, Path: u.EscapedPath(),
		NetworkTimeout: 5000, BufferTime: -1, RTSPReliable: 1, IgnoreEmbeddedKLV: false, RoverPort: -1,
	}}
	b, err := xml.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Feed is the TAK Server video feed (FeedV2) body for POST /Marti/api/video. The field list is
// the one read from the TAK Server 5.7-RELEASE-14 source (tak-protocols.md section 4); TAK Server
// 5.8 and the wrapping collection format are UNVERIFIED, so the planned tak-client must check
// them against a live server (owner step) before posting.
type Feed struct {
	UID               string   `json:"uid"`
	Alias             string   `json:"alias"`
	URL               string   `json:"url"`
	Latitude          *float64 `json:"latitude,omitempty"`
	Longitude         *float64 `json:"longitude,omitempty"`
	FOV               *float64 `json:"fov,omitempty"`
	Heading           *float64 `json:"heading,omitempty"`
	Range             *float64 `json:"range,omitempty"`
	IgnoreEmbeddedKLV bool     `json:"ignoreEmbeddedKLV"`
	RTSPReliable      int      `json:"rtspReliable"`
	BufferTime        int      `json:"bufferTime"`
}

// NewFeed builds the feed for an entity, or returns an error when the URL fails hub hygiene.
func NewFeed(entityType, entityID, alias, takURL string) (Feed, error) {
	if err := Check(takURL, Hub); err != nil {
		return Feed{}, err
	}
	return Feed{UID: UID(entityType, entityID), Alias: alias, URL: takURL, RTSPReliable: 1, BufferTime: -1}, nil
}
