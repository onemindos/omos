// Package model is the language-model seam of agentops. Every model call goes through a
// Chain: a primary (Hermes over the agent protocol, or an OpenAI-compatible endpoint), an
// optional local fallback (an OpenAI-compatible server on loopback such as Ollama or
// llama.cpp `llama-server`), and the deterministic Stub as the last resort, so a brief or an
// explanation is still produced with the network down (FEATURES D10).
//
// Safety rules in code:
//   - no API keys are read or sent: a paid hosted API cannot be called (run rule 10);
//   - an OpenAI-compatible endpoint must be loopback or a cluster-internal host unless
//     AllowRemote is set explicitly;
//   - model output is data: callers parse it, validate it against a closed list, and never
//     execute it.
package model

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/nats-io/nats.go"

	"github.com/onemindos/omos/services/agentops/internal/agentproto"
)

// Purposes.
const (
	PurposeExplain = "explain" // a watch fired: explain it and choose at most one allowed proposal
	PurposeBrief   = "brief"   // an AOI brief
)

// Request is one completion request. Structured carries the same facts as Prompt in JSON, so
// the Stub can answer deterministically; real models read System and Prompt.
type Request struct {
	Purpose    string
	System     string
	Prompt     string
	Structured json.RawMessage
	MaxTokens  int
}

// Model completes a request.
type Model interface {
	Name() string
	// Local reports whether the model runs without the wide-area network (loopback or the
	// built-in stub). A brief made by a local model is marked offline-capable.
	Local() bool
	Complete(ctx context.Context, r Request) (string, error)
}

// Attempt records one model call in a chain.
type Attempt struct {
	Model string `json:"model"`
	OK    bool   `json:"ok"`
	Error string `json:"error,omitempty"`
	MS    int64  `json:"ms"`
}

// Chain tries each model in order until one answers.
type Chain struct {
	Models  []Model
	Timeout time.Duration // per model; default 60 s
}

// Result is the chain's answer.
type Result struct {
	Text     string
	Model    string
	Local    bool
	Attempts []Attempt
}

// ErrAllFailed means no model in the chain answered.
var ErrAllFailed = errors.New("every model in the chain failed")

// Complete runs the chain.
func (c Chain) Complete(ctx context.Context, r Request) (Result, error) {
	to := c.Timeout
	if to <= 0 {
		to = 60 * time.Second
	}
	var res Result
	for _, m := range c.Models {
		if m == nil {
			continue
		}
		start := time.Now()
		cctx, cancel := context.WithTimeout(ctx, to)
		text, err := m.Complete(cctx, r)
		cancel()
		a := Attempt{Model: m.Name(), OK: err == nil && strings.TrimSpace(text) != "", MS: time.Since(start).Milliseconds()}
		if err != nil {
			a.Error = truncate(err.Error(), 200)
		} else if !a.OK {
			a.Error = "empty answer"
		}
		res.Attempts = append(res.Attempts, a)
		if a.OK {
			res.Text, res.Model, res.Local = text, m.Name(), m.Local()
			return res, nil
		}
		if ctx.Err() != nil {
			break
		}
	}
	return res, ErrAllFailed
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}

// OpenAICompat calls POST {BaseURL}/chat/completions of an OpenAI-compatible server
// (Ollama `/v1`, llama.cpp `llama-server`, vLLM). No Authorization header is ever sent.
type OpenAICompat struct {
	Label       string
	BaseURL     string // e.g. http://127.0.0.1:11434/v1
	ModelName   string // e.g. qwen2.5:7b-instruct
	AllowRemote bool
	HTTP        *http.Client
}

func (o *OpenAICompat) Name() string {
	if o.Label != "" {
		return o.Label
	}
	return "openai-compat:" + o.ModelName
}

func (o *OpenAICompat) Local() bool { return isLoopback(o.BaseURL) }

func isLoopback(raw string) bool {
	u, err := url.Parse(raw)
	if err != nil {
		return false
	}
	h := u.Hostname()
	if h == "localhost" {
		return true
	}
	ip := net.ParseIP(h)
	return ip != nil && ip.IsLoopback()
}

// allowedHost: loopback, private address, or a cluster-internal name.
func allowedHost(raw string) error {
	u, err := url.Parse(raw)
	if err != nil {
		return err
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("scheme %q not allowed", u.Scheme)
	}
	h := u.Hostname()
	if h == "localhost" || strings.HasSuffix(h, ".svc") || strings.HasSuffix(h, ".svc.cluster.local") || strings.HasSuffix(h, ".ts.net") {
		return nil
	}
	if ip := net.ParseIP(h); ip != nil && (ip.IsLoopback() || ip.IsPrivate()) {
		return nil
	}
	return fmt.Errorf("host %q is not loopback, private or cluster-internal (set AllowRemote to override)", h)
}

func (o *OpenAICompat) Complete(ctx context.Context, r Request) (string, error) {
	if !o.AllowRemote {
		if err := allowedHost(o.BaseURL); err != nil {
			return "", err
		}
	}
	maxTok := r.MaxTokens
	if maxTok <= 0 {
		maxTok = 600
	}
	body, _ := json.Marshal(map[string]any{
		"model": o.ModelName,
		"messages": []map[string]string{
			{"role": "system", "content": r.System},
			{"role": "user", "content": r.Prompt},
		},
		"temperature": 0.2,
		"max_tokens":  maxTok,
		"stream":      false,
	})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(o.BaseURL, "/")+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	hc := o.HTTP
	if hc == nil {
		hc = &http.Client{Timeout: 90 * time.Second}
	}
	resp, err := hc.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return "", err
	}
	if resp.StatusCode/100 != 2 {
		return "", fmt.Errorf("model server status %d", resp.StatusCode)
	}
	var out struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(raw, &out); err != nil {
		return "", err
	}
	if len(out.Choices) == 0 {
		return "", errors.New("model server returned no choices")
	}
	return out.Choices[0].Message.Content, nil
}

// Agent asks an agent (Hermes) over the agent protocol. The prompt carries the system text
// first; the envelope's `onemind` block names the purpose so the agent can route it.
type Agent struct {
	NC         *nats.Conn
	Subject    string // agents.prompt.hermes.zeus.legacy
	Inactivity time.Duration
}

func (a *Agent) Name() string { return "agent:" + a.Subject }
func (a *Agent) Local() bool  { return false }

func (a *Agent) Complete(ctx context.Context, r Request) (string, error) {
	env := map[string]any{
		"prompt":  r.System + "\n\n" + r.Prompt,
		"onemind": map[string]any{"purpose": r.Purpose, "tools": "none", "source": "omos-agentops"},
	}
	return agentproto.Prompt(ctx, a.NC, a.Subject, env, a.Inactivity)
}
