package model

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type failing struct{ name string }

func (f failing) Name() string { return f.name }
func (f failing) Local() bool  { return false }
func (f failing) Complete(context.Context, Request) (string, error) {
	return "", errors.New("network is down")
}

func TestChainFallsBackToStub(t *testing.T) {
	in, _ := json.Marshal(ExplainInput{WatchID: "w", RuleKind: "zone", Title: "SIM-X1 left zone", About: "drone/sim-x1",
		AllowedTaskTypes: []string{"onemind.task.rtl.v1"}})
	res, err := Chain{Models: []Model{failing{"primary"}, failing{"local"}, Stub{}}}.Complete(context.Background(),
		Request{Purpose: PurposeExplain, Structured: in})
	if err != nil || res.Model != "stub" || !res.Local || len(res.Attempts) != 3 || res.Attempts[0].OK {
		t.Fatalf("%+v %v", res, err)
	}
	var out ExplainOutput
	if err := json.Unmarshal([]byte(res.Text), &out); err != nil || out.Proposal == nil || out.Proposal.TaskType != "onemind.task.rtl.v1" {
		t.Fatalf("stub output %s", res.Text)
	}
	if _, err := (Chain{Models: []Model{failing{"a"}}}).Complete(context.Background(), Request{}); err != ErrAllFailed {
		t.Fatal("all failed must error")
	}
}

func TestOpenAICompatLoopback(t *testing.T) {
	var auth string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth = r.Header.Get("Authorization")
		if r.URL.Path != "/v1/chat/completions" {
			http.NotFound(w, r)
			return
		}
		var body struct {
			Model    string `json:"model"`
			Messages []struct{ Role, Content string }
		}
		_ = json.NewDecoder(r.Body).Decode(&body)
		_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"local says [sim-x1] ` + body.Model + `"}}]}`))
	}))
	defer srv.Close()
	m := &OpenAICompat{BaseURL: srv.URL + "/v1", ModelName: "qwen2.5:7b-instruct"}
	if !m.Local() {
		t.Fatal("httptest server is loopback")
	}
	out, err := m.Complete(context.Background(), Request{System: "s", Prompt: "p"})
	if err != nil || !strings.Contains(out, "qwen2.5:7b-instruct") || auth != "" {
		t.Fatalf("%q %v auth=%q", out, err, auth)
	}
}

func TestOpenAICompatRefusesRemote(t *testing.T) {
	m := &OpenAICompat{BaseURL: "https://api.example.com/v1", ModelName: "x"}
	if _, err := m.Complete(context.Background(), Request{}); err == nil || !strings.Contains(err.Error(), "not loopback") {
		t.Fatalf("remote host must be refused: %v", err)
	}
	for _, ok := range []string{"http://127.0.0.1:11434/v1", "http://ollama.ai.svc.cluster.local:11434/v1", "http://10.0.0.5:8080/v1"} {
		if err := allowedHost(ok); err != nil {
			t.Errorf("%s: %v", ok, err)
		}
	}
}

func TestStubBrief(t *testing.T) {
	la, lo := 30.2, -97.7
	in, _ := json.Marshal(BriefInput{AOIID: "home", AOIName: "Home", AtUTC: "2026-10-29T12:00:00Z", Counts: map[string]int{"drone": 1, "aircraft": 2},
		Entities: []BriefEntity{{ID: "sim-x1", Type: "drone", Callsign: "SIM-X1", Lat: &la, Lon: &lo}, {ID: "icao-a1", Type: "aircraft"}, {ID: "icao-b2", Type: "aircraft"}},
		Gaps:     []BriefGap{{Source: "nws-alerts", Reason: "stale", Layers: []string{"weather"}}}})
	out, err := Stub{}.Complete(context.Background(), Request{Purpose: PurposeBrief, Structured: in})
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"[sim-x1]", "[icao-a1]", "[icao-b2]", "nws-alerts", "weather", "3 entities"} {
		if !strings.Contains(out, want) {
			t.Errorf("stub brief lacks %q:\n%s", want, out)
		}
	}
}
