package natstest

import (
	"io"
	"net"
	"sync"
	"testing"
)

// Proxy is a TCP relay in front of the server. Cut drops every relayed connection and
// refuses new ones until Restore, so one client is cut off from the bus while every other
// client stays connected (a network split around one instance).
type Proxy struct {
	ln     net.Listener
	target string
	mu     sync.Mutex
	cut    bool
	conns  map[net.Conn]struct{}
}

// Proxy starts a relay to the server; it stops when the test ends.
func (s *Server) Proxy(t testing.TB) *Proxy {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	p := &Proxy{ln: ln, target: s.URL()[len("nats://"):], conns: map[net.Conn]struct{}{}}
	go p.accept()
	t.Cleanup(func() { _ = ln.Close(); p.Cut() })
	return p
}

func (p *Proxy) URL() string { return "nats://" + p.ln.Addr().String() }

func (p *Proxy) accept() {
	for {
		c, err := p.ln.Accept()
		if err != nil {
			return
		}
		p.mu.Lock()
		cut := p.cut
		p.mu.Unlock()
		if cut {
			_ = c.Close()
			continue
		}
		u, err := net.Dial("tcp", p.target)
		if err != nil {
			_ = c.Close()
			continue
		}
		p.mu.Lock()
		p.conns[c], p.conns[u] = struct{}{}, struct{}{}
		p.mu.Unlock()
		go p.pipe(c, u)
		go p.pipe(u, c)
	}
}

func (p *Proxy) pipe(dst, src net.Conn) {
	_, _ = io.Copy(dst, src)
	_ = dst.Close()
	_ = src.Close()
}

// Cut drops all relayed connections and refuses new ones.
func (p *Proxy) Cut() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.cut = true
	for c := range p.conns {
		_ = c.Close()
	}
	p.conns = map[net.Conn]struct{}{}
}

// Restore accepts connections again.
func (p *Proxy) Restore() {
	p.mu.Lock()
	p.cut = false
	p.mu.Unlock()
}
