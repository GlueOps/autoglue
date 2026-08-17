// Command agent runs the autoglue bastion agent.
//
// Scaffold: the two plane loops and the durable store are wired up, but the
// API client and the task runner are not implemented yet — the control-plane
// endpoints do not exist. Starting it now will fail at client construction,
// deliberately and loudly, rather than pretending to run.
package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/glueops/autoglue/agent/internal/runner"
	"github.com/glueops/autoglue/agent/internal/store"
	"github.com/glueops/autoglue/agent/internal/supervisor"
)

func main() {
	if err := run(); err != nil {
		slog.Error("agent exited", "error", err)
		os.Exit(1)
	}
}

func run() error {
	var (
		stateDir = flag.String("state-dir", "/var/lib/autoglue-agent",
			"durable state directory; must survive re-provisioning, or an unreported result is lost with it")
		configInterval = flag.Duration("config-interval", 30*time.Second,
			"how often to sync desired state; this is also the liveness heartbeat")
		taskInterval = flag.Duration("task-interval", 10*time.Second,
			"how often to ask what task to be working on")
		outboxInterval = flag.Duration("outbox-interval", 5*time.Second,
			"how often to drain undelivered results and logs")
		debug = flag.Bool("debug", false, "verbose logging")
	)
	flag.Parse()

	level := slog.LevelInfo
	if *debug {
		level = slog.LevelDebug
	}
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level}))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := os.MkdirAll(*stateDir, 0o700); err != nil {
		return fmt.Errorf("create state dir: %w", err)
	}

	st, err := store.Open(ctx, filepath.Join(*stateDir, "agent.db"))
	if err != nil {
		return err
	}
	defer func() { _ = st.Close() }()

	sup := supervisor.New(st, nil, nil, supervisor.Options{
		ConfigInterval: *configInterval,
		TaskInterval:   *taskInterval,
		OutboxInterval: *outboxInterval,
		Log:            log,
	}).WithRunner(runner.NewDocker(runner.ExecCommander{}))

	// Adoption runs before the loops, so a container that survived the last
	// shutdown is picked up immediately — and, more importantly, so one that
	// vanished is dead-lettered before anything can be assigned on top of it.
	if err := sup.AdoptOnStart(ctx); err != nil {
		return fmt.Errorf("adopt on start: %w", err)
	}

	// Not implemented: the control-plane endpoints the supervisor calls do not
	// exist yet. Failing here is the honest outcome — a no-op client would let
	// the agent report healthy convergence to nothing at all.
	return fmt.Errorf("api client not implemented yet; store opened at %s, "+
		"intervals config=%s task=%s outbox=%s",
		*stateDir, *configInterval, *taskInterval, *outboxInterval)
}
