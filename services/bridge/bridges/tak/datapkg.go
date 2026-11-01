// SPDX-License-Identifier: Apache-2.0

package main

// Peer-hosted data package transfer, as Common Commo does it on a mesh
// (commoncommo core/missionpackage-rules.txt, core/impl/cotmessage.cpp):
//
//   - The sender serves the file at http://<ip>:<port>/getfile?file=<id> and
//     sends the recipient a directed CoT file transfer request, type b-f-t-r,
//     with <fileshare filename senderUrl sizeInBytes sha256 senderUid
//     senderCallsign name/> and, optionally, <ackrequest uid ackrequested="true"/>.
//   - The recipient fetches the URL, and replies with a b-f-t-a CoT carrying
//     <ackresponse uid senderUid success tag reason sha256 sizeInBytes/>.
//
// The gateway fetches only requests addressed to it (directed TCP), only over
// http, only from the IP address the request came from (Commo's own fallback
// rule 4b uses that address), only up to DATAPKG_MAX_BYTES, and keeps the file
// only if its size and SHA-256 match the request. At most MaxConcurrent
// fetches run at once, a package already in the inbox (same SHA-256) is acked
// without fetching it again, and a request that would take the inbox over
// MaxInboxBytes is refused with a nack.

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	typeFileRequest = "b-f-t-r"
	typeFileAck     = "b-f-t-a"
	ackStale        = 10 * time.Second // Commo STALE_TIME_FILE_ACK

	defaultMaxConcurrentFetches = 2
)

// FileShare is a parsed CoT file transfer request.
type FileShare struct {
	Filename, SenderURL, SHA256, SenderUID, SenderCallsign, Name string
	Size                                                         int64
	AckUID                                                       string // set when an ack was requested
}

type xmlFileDetail struct {
	FileShare *struct {
		Filename       string `xml:"filename,attr"`
		SenderURL      string `xml:"senderUrl,attr"`
		Size           string `xml:"sizeInBytes,attr"`
		SHA256         string `xml:"sha256,attr"`
		SenderUID      string `xml:"senderUid,attr"`
		SenderCallsign string `xml:"senderCallsign,attr"`
		Name           string `xml:"name,attr"`
	} `xml:"fileshare"`
	AckRequest *struct {
		UID          string `xml:"uid,attr"`
		AckRequested string `xml:"ackrequested,attr"`
	} `xml:"ackrequest"`
}

var sha256Hex = regexp.MustCompile(`^[0-9a-fA-F]{64}$`)

// ParseFileShare extracts the file transfer request from ev.
func ParseFileShare(ev *Event) (*FileShare, error) {
	if ev.Type != typeFileRequest {
		return nil, fmt.Errorf("type %s is not %s", ev.Type, typeFileRequest)
	}
	var d xmlFileDetail
	if err := xml.Unmarshal([]byte("<detail>"+ev.Detail.XML+"</detail>"), &d); err != nil {
		return nil, fmt.Errorf("file request detail: %w", err)
	}
	if d.FileShare == nil {
		return nil, errors.New("file request has no <fileshare>")
	}
	f := d.FileShare
	size, err := strconv.ParseInt(f.Size, 10, 64)
	if err != nil || size < 0 {
		return nil, fmt.Errorf("fileshare sizeInBytes %q", f.Size)
	}
	if !sha256Hex.MatchString(f.SHA256) {
		return nil, fmt.Errorf("fileshare sha256 %q", f.SHA256)
	}
	fs := &FileShare{
		Filename: f.Filename, SenderURL: f.SenderURL, SHA256: strings.ToLower(f.SHA256),
		SenderUID: f.SenderUID, SenderCallsign: f.SenderCallsign, Name: f.Name, Size: size,
	}
	if a := d.AckRequest; a != nil && a.AckRequested == "true" {
		fs.AckUID = a.UID
	}
	return fs, nil
}

// DataPackages receives and serves peer-hosted data packages under Dir:
// received files go to Dir/inbox, files to offer go in Dir/outbox.
type DataPackages struct {
	Dir      string
	MaxBytes int64 // largest single package
	// MaxInboxBytes caps the inbox, counting files and fetches in flight;
	// 0 means no cap. MaxConcurrent caps fetches in flight (default 2).
	MaxInboxBytes int64
	MaxConcurrent int
	SelfUID       string
	Client        *http.Client
	// SendAck delivers an ack to the sender; senderUID and from (the
	// request's network source) let it find the sender's endpoint.
	SendAck func(ack *Event, senderUID, from string) error
	OnErr   func(error)

	Received atomic.Int64
	Failed   atomic.Int64

	mu       sync.Mutex
	inflight map[string]int64 // sha256 -> size, fetches running
}

func (d *DataPackages) inbox() string  { return filepath.Join(d.Dir, "inbox") }
func (d *DataPackages) outbox() string { return filepath.Join(d.Dir, "outbox") }

// Init creates the inbox and outbox directories.
func (d *DataPackages) Init() error {
	for _, p := range []string{d.inbox(), d.outbox()} {
		if err := os.MkdirAll(p, 0o750); err != nil {
			return err
		}
	}
	return nil
}

// Handle processes one received message. Only directed file transfer requests
// are acted on; it returns at once and fetches in the background.
func (d *DataPackages) Handle(r Received) {
	if r.Event.Type != typeFileRequest || r.Delivery != DeliveryDirect {
		return
	}
	fs, err := ParseFileShare(r.Event)
	if err != nil {
		d.OnErr(fmt.Errorf("data package from %s: %w", r.From, err))
		return
	}
	switch release, err := d.admit(fs); {
	case errors.Is(err, errInFlight):
		// The running fetch of the same package will ack it.
	case errors.Is(err, errHave):
		go d.reply(fs, r.From, nil)
	case err != nil:
		d.Failed.Add(1)
		d.OnErr(fmt.Errorf("data package %q from %s: %w", fs.Name, fs.SenderUID, err))
		go d.reply(fs, r.From, err)
	default:
		go func() {
			defer release()
			d.receive(fs, r.From)
		}()
	}
}

var (
	errInFlight = errors.New("already being fetched")
	errHave     = errors.New("already received")
)

func (d *DataPackages) storedPath(fs *FileShare) string {
	return filepath.Join(d.inbox(), fs.SHA256+"-"+safeFilename(fs.Filename))
}

// admit reserves a fetch slot and inbox space for fs. The caller runs release
// when the fetch ends.
func (d *DataPackages) admit(fs *FileShare) (release func(), err error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.inflight == nil {
		d.inflight = map[string]int64{}
	}
	if _, ok := d.inflight[fs.SHA256]; ok {
		return nil, errInFlight
	}
	if m, _ := filepath.Glob(filepath.Join(d.inbox(), fs.SHA256+"-*")); len(m) > 0 {
		return nil, errHave
	}
	limit := d.MaxConcurrent
	if limit <= 0 {
		limit = defaultMaxConcurrentFetches
	}
	if len(d.inflight) >= limit {
		return nil, fmt.Errorf("busy: %d transfers in progress", len(d.inflight))
	}
	if d.MaxInboxBytes > 0 {
		used, err := dirBytes(d.inbox())
		if err != nil {
			return nil, err
		}
		for _, n := range d.inflight {
			used += n
		}
		if used+fs.Size > d.MaxInboxBytes {
			return nil, fmt.Errorf("inbox full: %d bytes used, limit %d", used, d.MaxInboxBytes)
		}
	}
	d.inflight[fs.SHA256] = fs.Size
	return func() {
		d.mu.Lock()
		delete(d.inflight, fs.SHA256)
		d.mu.Unlock()
	}, nil
}

// dirBytes sums the sizes of the regular files in dir.
func dirBytes(dir string) (int64, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, err
	}
	var n int64
	for _, e := range entries {
		if fi, err := e.Info(); err == nil && fi.Mode().IsRegular() {
			n += fi.Size()
		}
	}
	return n, nil
}

func (d *DataPackages) receive(fs *FileShare, from string) {
	_, err := d.Fetch(fs, from)
	if err != nil {
		d.Failed.Add(1)
		d.OnErr(fmt.Errorf("data package %q from %s: %w", fs.Name, fs.SenderUID, err))
	} else {
		d.Received.Add(1)
	}
	d.reply(fs, from, err)
}

// reply sends the ack the sender asked for, if any.
func (d *DataPackages) reply(fs *FileShare, from string, err error) {
	if fs.AckUID == "" || d.SendAck == nil {
		return
	}
	reason := "Transfer complete"
	if err != nil {
		reason = err.Error()
	}
	if aerr := d.SendAck(FileAck(fs, d.SelfUID, err == nil, reason, time.Now()), fs.SenderUID, from); aerr != nil {
		d.OnErr(fmt.Errorf("data package ack to %s: %w", fs.SenderUID, aerr))
	}
}

// fetchURL rewrites the request URL to the request's source address.
func fetchURL(senderURL, from string) (string, error) {
	u, err := url.Parse(senderURL)
	if err != nil {
		return "", fmt.Errorf("senderUrl: %w", err)
	}
	if u.Scheme != "http" {
		return "", fmt.Errorf("senderUrl scheme %q: only peer-hosted http is supported", u.Scheme)
	}
	host, _, err := net.SplitHostPort(from)
	if err != nil {
		return "", fmt.Errorf("request source %q: %w", from, err)
	}
	port := u.Port()
	if port == "" {
		port = "80"
	}
	u.Host = net.JoinHostPort(host, port)
	u.User = nil
	return u.String(), nil
}

var unsafeFileChars = regexp.MustCompile(`[^a-zA-Z0-9._-]`)

func safeFilename(name string) string {
	name = unsafeFileChars.ReplaceAllString(filepath.Base(strings.ReplaceAll(name, `\`, "/")), "_")
	name = strings.TrimLeft(name, ".")
	if len(name) > 100 {
		name = name[len(name)-100:]
	}
	if name == "" {
		name = "package.zip"
	}
	return name
}

// Fetch downloads fs from the sender and stores it in the inbox if its size
// and hash match. It returns the stored path.
func (d *DataPackages) Fetch(fs *FileShare, from string) (string, error) {
	if fs.Size > d.MaxBytes {
		return "", fmt.Errorf("size %d over limit %d", fs.Size, d.MaxBytes)
	}
	u, err := fetchURL(fs.SenderURL, from)
	if err != nil {
		return "", err
	}
	client := d.Client
	if client == nil {
		client = &http.Client{Timeout: 2 * time.Minute}
	}
	resp, err := client.Get(u)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GET %s: %s", u, resp.Status)
	}
	tmp, err := os.CreateTemp(d.inbox(), ".partial-*")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmp.Name()) // no-op after a successful rename
	h := sha256.New()
	n, err := io.Copy(io.MultiWriter(tmp, h), io.LimitReader(resp.Body, d.MaxBytes+1))
	if cerr := tmp.Close(); err == nil {
		err = cerr
	}
	if err != nil {
		return "", err
	}
	if n > d.MaxBytes {
		return "", fmt.Errorf("body over limit %d", d.MaxBytes)
	}
	if n != fs.Size {
		return "", fmt.Errorf("size %d, request said %d", n, fs.Size)
	}
	if got := hex.EncodeToString(h.Sum(nil)); got != fs.SHA256 {
		return "", fmt.Errorf("sha256 %s, request said %s", got, fs.SHA256)
	}
	dst := d.storedPath(fs)
	if err := os.Rename(tmp.Name(), dst); err != nil {
		return "", err
	}
	return dst, nil
}

func randomID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// FileAck builds the b-f-t-a reply Commo expects (cotmessage.cpp, ackresponse).
func FileAck(fs *FileShare, selfUID string, ok bool, reason string, now time.Time) *Event {
	var b strings.Builder
	b.WriteString("<ackresponse")
	for _, a := range [][2]string{
		{"uid", fs.AckUID}, {"senderUid", selfUID}, {"success", strconv.FormatBool(ok)},
		{"tag", fs.Name}, {"reason", reason}, {"sha256", fs.SHA256},
		{"sizeInBytes", strconv.FormatInt(fs.Size, 10)},
	} {
		b.WriteString(" " + a[0] + `="`)
		_ = xml.EscapeText(&b, []byte(a[1]))
		b.WriteString(`"`)
	}
	b.WriteString("/>")
	ms := uint64(now.UnixMilli())
	return &Event{
		Type: typeFileAck, How: "m-g", UID: selfUID + "-ack-" + randomID(),
		SendTime: ms, StartTime: ms, StaleTime: ms + uint64(ackStale.Milliseconds()),
		HAE: 9999999, CE: 9999999, LE: 9999999,
		Detail: Detail{XML: b.String()},
	}
}

// ServeHTTP answers GET /getfile?file=<name> from the outbox, the URL form
// Commo receivers fetch. Names are plain file names; nothing else is served.
func (d *DataPackages) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/getfile" || (r.Method != http.MethodGet && r.Method != http.MethodHead) {
		http.NotFound(w, r)
		return
	}
	name := r.URL.Query().Get("file")
	if name == "" || name != filepath.Base(name) || strings.HasPrefix(name, ".") ||
		strings.ContainsAny(name, `/\`) {
		http.NotFound(w, r)
		return
	}
	f, err := os.Open(filepath.Join(d.outbox(), name))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer f.Close()
	st, err := f.Stat()
	if err != nil || !st.Mode().IsRegular() {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeContent(w, r, name, st.ModTime(), f)
}
