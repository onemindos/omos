// chaos-verify reads every message the collectors stored and compares it with the ids the
// stub emitted: duplicates (a Nats-Msg-Id stored twice), missing (emitted, valid, never
// stored), dead letters (malformed ids must land on evt.sys.{source}.rejected, never on
// the data subject), and the last feed status. It writes a JSON report and exits 1 on
// any duplicate, missing or misrouted record.
package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type emitted struct {
	T   time.Time `json:"t"`
	Src string    `json:"src"`
	ID  string    `json:"id"`
	Bad bool      `json:"bad"`
}

type sourceReport struct {
	Subject          string   `json:"subject"`
	EmittedValid     int      `json:"emittedValid"`
	EmittedMalformed int      `json:"emittedMalformed"`
	Stored           int      `json:"stored"`
	Duplicates       int      `json:"duplicates"`
	DuplicateIDs     []string `json:"duplicateIds,omitempty"`
	Missing          int      `json:"missing"`
	MissingIDs       []string `json:"missingIds,omitempty"`
	MalformedStored  int      `json:"malformedStoredOnDataSubject"`
	DeadLetters      int      `json:"deadLetters"`
	IgnoredTail      int      `json:"ignoredTail"`
}

func main() {
	url := flag.String("nats", "nats://127.0.0.1:34322", "NATS URL")
	logPath := flag.String("log", "emitted.jsonl", "stub emitted log")
	tail := flag.Duration("tail", 15*time.Second, "ignore ids emitted in the last tail before the newest")
	out := flag.String("out", "verify.json", "report path")
	flag.Parse()

	nc, err := nats.Connect(*url)
	must(err)
	defer nc.Close()
	js, err := jetstream.New(nc)
	must(err)

	em := map[string][]emitted{}
	var newest time.Time
	f, err := os.Open(*logPath)
	must(err)
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		var e emitted
		if json.Unmarshal(sc.Bytes(), &e) == nil {
			em[e.Src] = append(em[e.Src], e)
			if e.T.After(newest) {
				newest = e.T
			}
		}
	}
	cutoff := newest.Add(-*tail)

	report := map[string]any{"cutoff": cutoff, "generatedAt": time.Now().UTC()}
	sources := map[string]*sourceReport{}
	ok := true
	for src, subj := range map[string]string{"nws": "intel.weather.noaa.alert", "sse": "intel.test.chaos-sse.event"} {
		stored := map[string]int{}
		for _, m := range read(js, "INTEL", subj) {
			stored[m.Headers().Get(jetstream.MsgIDHeader)]++
		}
		feed := map[string]string{"nws": "chaos-nws", "sse": "chaos-sse"}[src]
		r := &sourceReport{Subject: subj, DeadLetters: len(read(js, "SYS", "evt.sys."+feed+".rejected"))}
		for id, n := range stored {
			r.Stored++
			if n > 1 {
				r.Duplicates += n - 1
				r.DuplicateIDs = append(r.DuplicateIDs, id)
			}
		}
		for _, e := range em[src] {
			if e.T.After(cutoff) {
				r.IgnoredTail++
				continue
			}
			if e.Bad {
				r.EmittedMalformed++
				if stored[e.ID] > 0 {
					r.MalformedStored++
				}
				continue
			}
			r.EmittedValid++
			if stored[e.ID] == 0 {
				r.Missing++
				if len(r.MissingIDs) < 50 {
					r.MissingIDs = append(r.MissingIDs, e.ID)
				}
			}
		}
		sort.Strings(r.DuplicateIDs)
		if len(r.DuplicateIDs) > 50 {
			r.DuplicateIDs = r.DuplicateIDs[:50]
		}
		if r.Duplicates > 0 || r.Missing > 0 || r.MalformedStored > 0 || (r.EmittedMalformed > 0 && r.DeadLetters < r.EmittedMalformed) {
			ok = false
		}
		sources[src] = r
	}
	report["sources"] = sources

	status := map[string]any{}
	for _, m := range read(js, "TEL", "tel.feed.>") {
		var v map[string]any
		if json.Unmarshal(m.Data(), &v) == nil {
			status[strings.TrimSuffix(strings.TrimPrefix(m.Subject(), "tel.feed."), ".status")] = v
		}
	}
	report["lastStatus"] = status
	report["pass"] = ok
	b, _ := json.MarshalIndent(report, "", "  ")
	must(os.WriteFile(*out, b, 0o644))
	fmt.Println(string(b))
	if !ok {
		os.Exit(1)
	}
}

func read(js jetstream.JetStream, stream, filter string) []jetstream.Msg {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	st, err := js.Stream(ctx, stream)
	must(err)
	c, err := st.OrderedConsumer(ctx, jetstream.OrderedConsumerConfig{FilterSubjects: []string{filter}})
	must(err)
	info, err := c.Info(ctx)
	must(err)
	var out []jetstream.Msg
	for pending := info.NumPending; pending > 0; {
		b, err := c.Fetch(int(min(pending, 5000)), jetstream.FetchMaxWait(5*time.Second))
		must(err)
		n := 0
		for m := range b.Messages() {
			out = append(out, m)
			n++
		}
		if n == 0 {
			break
		}
		pending -= uint64(n)
	}
	return out
}

func must(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
}
