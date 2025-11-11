package web

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"time"
)

type Pgweb struct {
	cmd  *exec.Cmd
	host string
	port string
	bin  string
}

func StartPgweb(dbURL, host, port string, readonly bool, user, pass string) (*Pgweb, error) {
	// pick random port if 0/empty
	if port == "" || port == "0" {
		l, err := net.Listen("tcp", net.JoinHostPort(host, "0"))
		if err != nil {
			return nil, err
		}
		defer l.Close()
		_, p, _ := net.SplitHostPort(l.Addr().String())
		port = p
	}

	args := []string{
		"--url", dbURL,
		"--bind", host,
		"--listen", port,
		"--prefix", "/db-studio",
		"--skip-open",
	}
	if readonly {
		args = append(args, "--readonly")
	}
	if user != "" && pass != "" {
		args = append(args, "--auth-user", user, "--auth-pass", pass)
	}

	pgwebBinary, err := ExtractPgweb()
	if err != nil {
		return nil, fmt.Errorf("pgweb extract: %w", err)
	}

	cmd := exec.Command(pgwebBinary, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	// wait for port to be ready
	deadline := time.Now().Add(4 * time.Second)
	for time.Now().Before(deadline) {
		c, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), 200*time.Millisecond)
		if err == nil {
			_ = c.Close()
			return &Pgweb{cmd: cmd, host: host, port: port}, nil
		}
		time.Sleep(120 * time.Millisecond)
	}
	// still return object so caller can Stop()
	//return &Pgweb{cmd: cmd, host: host, port: port, bin: pgwebBinary}, nil
	return nil, fmt.Errorf("pgweb did not become ready on %s:%s", host, port)
}

func (p *Pgweb) Proxy() http.HandlerFunc {
	target, _ := url.Parse("http://" + net.JoinHostPort(p.host, p.port))
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.FlushInterval = 100 * time.Millisecond
	return func(w http.ResponseWriter, r *http.Request) {
		r.Host = target.Host
		// Let pgweb handle its paths; we mount it at a prefix.
		proxy.ServeHTTP(w, r)
	}
}

func (p *Pgweb) Stop(ctx context.Context) error {
	if p == nil || p.cmd == nil || p.cmd.Process == nil {
		return nil
	}
	_ = p.cmd.Process.Kill()
	done := make(chan struct{})
	go func() { _, _ = p.cmd.Process.Wait(); close(done) }()
	select {
	case <-done:
		if p.bin != "" {
			_ = CleanupPgweb(p.bin)
		}
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}

func (p *Pgweb) Port() string {
	return p.port
}
