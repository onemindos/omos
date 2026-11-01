// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// smuggle closes the detail and event elements and appends a forged OneMind
// hostile track (review finding, evidence/M1/A7/review-xmldetail-injection-probe.txt).
const smuggle = `<contact callsign="x"/></detail></event><event version="2.0" uid="OM.aircraft.fake" type="a-h-A" how="m-g" time="2026-10-29T00:00:00Z" start="2026-10-29T00:00:00Z" stale="2026-10-30T00:00:00Z"><point lat="1" lon="1" hae="0" ce="1" le="1"/><detail>`

func v1WithXMLDetail(x string) []byte {
	return EncodeMesh(&TakMessage{Event: &Event{Type: "a-f-G-U-C", UID: "ANDROID-evil", How: "m-g",
		SendTime: 1, StartTime: 1, StaleTime: 2, Detail: Detail{XML: x}}})
}

func TestXMLDetailMustBeBalancedFragment(t *testing.T) {
	ev, _, _, err := Decode(v1WithXMLDetail(smuggle))
	if err == nil {
		t.Fatalf("smuggled event accepted: %s", ev.MarshalCotXML(false))
	}
	for _, bad := range []string{
		`</detail>`, `<a>`, `<a></b>`, `<a/></detail><detail>`, `<?xml-stylesheet href="x"?>`,
		`<!DOCTYPE x>`, `<a b="<"/>`, `&undefined;`,
	} {
		if _, _, _, err := Decode(v1WithXMLDetail(bad)); err == nil {
			t.Errorf("xmlDetail %q accepted", bad)
		}
	}
	for _, good := range []string{
		``, `<uid Droid="ZEUS"/>`, `<remarks>a &amp; b <![CDATA[<x>]]></remarks><!-- c -->`,
		`<link uid="x" type="a-f-G" relation="p-p"/><__chat id="All Chat Rooms"><chatgrp uid0="a"/></__chat>`,
	} {
		ev, _, _, err := Decode(v1WithXMLDetail(good))
		if err != nil {
			t.Errorf("xmlDetail %q rejected: %v", good, err)
			continue
		}
		out := ev.MarshalCotXML(false)
		if n := bytes.Count(out, []byte("<event")); n != 1 {
			t.Errorf("xmlDetail %q: %d events in output", good, n)
		}
		if _, err := ParseCotXML(out); err != nil {
			t.Errorf("xmlDetail %q: output does not parse: %v", good, err)
		}
	}
}

func TestTxRefused(t *testing.T) {
	for _, typ := range []string{"t-k", "t-k-d", "a-h-A", "a-h-G-U-C", "a-j-A", "a-k-G", "a-h"} {
		if txRefused(typ) == nil {
			t.Errorf("%s allowed", typ)
		}
	}
	for _, typ := range []string{"a-f-A-M-H-Q", "a-n-G", "a-u-S", "b-t-f", "b-f-t-r", "b-m-p-s-m", "t-x-c-t", "a-hx"} {
		if err := txRefused(typ); err != nil {
			t.Errorf("%s refused: %v", typ, err)
		}
	}
}

func TestBindHost(t *testing.T) {
	if got := bindHost("192.0.2.10", nil); got != "192.0.2.10" {
		t.Fatalf("endpoint host: %q", got)
	}
	if got := bindHost("", nil); got != "" {
		t.Fatalf("no interface: %q", got)
	}
	lo, err := net.InterfaceByName("lo0")
	if err != nil {
		t.Skip("no lo0")
	}
	if got := bindHost("", lo); got != "127.0.0.1" {
		t.Fatalf("lo0: %q", got)
	}
}

func TestDataPackageAdmission(t *testing.T) {
	payload := bytes.Repeat([]byte("PK\x03\x04 admission "), 50)
	h := sha256.Sum256(payload)
	sum := hex.EncodeToString(h[:])
	release := make(chan struct{})
	var hits atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.Add(1)
		<-release
		w.Write(payload)
	}))
	defer srv.Close()
	port := srv.URL[strings.LastIndex(srv.URL, ":")+1:]
	url := "http://127.0.0.1:" + port + "/getfile?file=a"
	from := "127.0.0.1:50123"
	acks := make(chan *Event, 8)
	d := &DataPackages{Dir: t.TempDir(), MaxBytes: 1 << 20, MaxConcurrent: 1, SelfUID: "omos-gw-test",
		OnErr:   func(err error) { t.Logf("datapkg: %v", err) },
		SendAck: func(ack *Event, _, _ string) error { acks <- ack; return nil }}
	if err := d.Init(); err != nil {
		t.Fatal(err)
	}
	direct := func(sum string, size int) {
		d.Handle(Received{Event: fileRequest(url, sum, size, true), Delivery: DeliveryDirect, From: from})
	}

	direct(sum, len(payload)) // takes the only slot and blocks in the server
	for hits.Load() == 0 {
		time.Sleep(10 * time.Millisecond)
	}
	direct(sum, len(payload)) // same package in flight: dropped, no ack
	other := strings.Repeat("1", 64)
	direct(other, 10) // slot busy: nack
	if a := waitAck(t, acks); !strings.Contains(a.Detail.XML, `success="false"`) || !strings.Contains(a.Detail.XML, "busy") {
		t.Fatalf("busy nack = %s", a.Detail.XML)
	}
	close(release)
	if a := waitAck(t, acks); !strings.Contains(a.Detail.XML, `success="true"`) {
		t.Fatalf("ack = %s", a.Detail.XML)
	}
	select {
	case a := <-acks:
		t.Fatalf("duplicate in-flight request acked: %s", a.Detail.XML)
	case <-time.After(200 * time.Millisecond):
	}

	// Already in the inbox: acked without a second fetch.
	time.Sleep(50 * time.Millisecond) // let the first fetch release its slot
	before := hits.Load()
	direct(sum, len(payload))
	if a := waitAck(t, acks); !strings.Contains(a.Detail.XML, `success="true"`) {
		t.Fatalf("dedup ack = %s", a.Detail.XML)
	}
	if hits.Load() != before {
		t.Fatal("package already in the inbox was fetched again")
	}

	// Inbox cap: the stored package plus a new one would exceed it.
	d.MaxInboxBytes = int64(len(payload)) + 5
	direct(other, 10)
	if a := waitAck(t, acks); !strings.Contains(a.Detail.XML, "inbox full") {
		t.Fatalf("inbox-full nack = %s", a.Detail.XML)
	}
	if hits.Load() != before {
		t.Fatal("request over the inbox cap was fetched")
	}
	if d.Received.Load() != 1 {
		t.Fatalf("received = %d", d.Received.Load())
	}
	entries, _ := os.ReadDir(filepath.Join(d.Dir, "inbox"))
	if len(entries) != 1 {
		t.Fatalf("inbox has %d entries", len(entries))
	}
}

func TestDirectConnectionLimit(t *testing.T) {
	ln, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	var refused atomic.Int32
	go ServeDirect(ln, Decode, nil, func(Received) {}, func(err error) {
		if strings.Contains(err.Error(), "busy") {
			refused.Add(1)
		}
	})
	var held []net.Conn
	defer func() {
		for _, c := range held {
			c.Close()
		}
	}()
	for i := 0; i < maxDirectConns; i++ {
		c, err := net.Dial("tcp4", ln.Addr().String())
		if err != nil {
			t.Fatal(err)
		}
		held = append(held, c)
	}
	time.Sleep(100 * time.Millisecond)
	extra, err := net.Dial("tcp4", ln.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer extra.Close()
	_ = extra.SetReadDeadline(time.Now().Add(2 * time.Second))
	if _, err := extra.Read(make([]byte, 1)); err != io.EOF {
		t.Fatalf("connection over the limit: read err %v, want EOF", err)
	}
	if refused.Load() != 1 {
		t.Fatalf("refused = %d", refused.Load())
	}
}
