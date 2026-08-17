package runner

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"
)

type call struct {
	name string
	args []string
}

type fakeCmd struct {
	calls   []call
	outputs map[string]string // matched by a substring of the joined args
	stream  string
	failOn  string
}

func (f *fakeCmd) record(name string, args []string) string {
	f.calls = append(f.calls, call{name: name, args: args})
	return strings.Join(args, " ")
}

func (f *fakeCmd) Output(_ context.Context, name string, args ...string) (string, error) {
	joined := f.record(name, args)
	if f.failOn != "" && strings.Contains(joined, f.failOn) {
		return "", errors.New("boom")
	}
	for key, out := range f.outputs {
		if strings.Contains(joined, key) {
			return out, nil
		}
	}
	return "", nil
}

func (f *fakeCmd) Stream(_ context.Context, name string, args ...string) (io.ReadCloser, func() error, error) {
	f.record(name, args)
	return io.NopCloser(strings.NewReader(f.stream)), func() error { return nil }, nil
}

func (f *fakeCmd) argsFor(sub string) []string {
	for _, c := range f.calls {
		if strings.Contains(strings.Join(c.args, " "), sub) {
			return c.args
		}
	}
	return nil
}

func has(args []string, want ...string) bool {
	joined := " " + strings.Join(args, " ") + " "
	for _, w := range want {
		if !strings.Contains(joined, " "+w+" ") {
			return false
		}
	}
	return true
}

// Detached is the decision the whole runner rests on: nothing is attached, so
// the agent dying or being upgraded cannot touch work in flight, and adopting a
// container is identical to following one we started.
func TestStartIsDetachedAndLabelled(t *testing.T) {
	f := &fakeCmd{outputs: map[string]string{"run": "container-abc\n"}}
	d := NewDocker(f)

	id, err := d.Start(context.Background(), Spec{
		TaskID: "task-1", ClusterID: "cl-1", RunID: "run-1",
		Image: "ghcr.io/glueops/gluekube", Tag: "v1", Target: "ping-servers",
		WorkDir: "/opt/cluster",
		Mounts:  []Mount{{Source: "/home/deploy/.ssh", Target: "/root/.ssh"}},
	})
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	if id != "container-abc" {
		t.Fatalf("container id = %q", id)
	}

	args := f.argsFor("run")
	if !has(args, "--detach") {
		t.Fatalf("container not detached: %v", args)
	}
	if !has(args, LabelCluster+"=cl-1", LabelRun+"=run-1", LabelTask+"=task-1") {
		t.Fatalf("labels missing: %v", args)
	}
	// No --rm: once a container outlives the agent it is the only record of what
	// happened after the stream was lost.
	if has(args, "--rm") {
		t.Fatalf("--rm would delete the evidence adoption exists to recover: %v", args)
	}
	if !has(args, "--volume", "/home/deploy/.ssh:/root/.ssh") {
		t.Fatalf("mount missing: %v", args)
	}
	if got := strings.Join(args, " "); !strings.HasSuffix(got, "ghcr.io/glueops/gluekube:v1 make ping-servers") {
		t.Fatalf("image and target must come last: %v", args)
	}
}

func TestFindReportsPhase(t *testing.T) {
	cases := []struct {
		name      string
		psOutput  string
		wantPhase Phase
		wantCode  int
	}{
		{"running", "abc running\n", PhaseRunning, 0},
		{"created counts as running", "abc created\n", PhaseRunning, 0},
		{"exited", "abc exited\n", PhaseExited, 2},
		{"absent", "\n", PhaseAbsent, 0},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := &fakeCmd{outputs: map[string]string{
				"ps":      tc.psOutput,
				"inspect": "2\n",
			}}
			got, err := NewDocker(f).Find(context.Background(), "task-1")
			if err != nil {
				t.Fatalf("find: %v", err)
			}
			if got.Phase != tc.wantPhase {
				t.Fatalf("phase = %q, want %q", got.Phase, tc.wantPhase)
			}
			if tc.wantPhase == PhaseExited && got.ExitCode != tc.wantCode {
				t.Fatalf("exit code = %d, want %d", got.ExitCode, tc.wantCode)
			}
		})
	}
}

// Find must search stopped containers too. An exited container whose result was
// never delivered is the single most valuable thing adoption recovers — if the
// lookup missed it, that outcome would be dead-lettered as unknowable when it
// is sitting right there.
func TestFindSearchesStoppedContainers(t *testing.T) {
	f := &fakeCmd{outputs: map[string]string{"ps": "abc exited\n", "inspect": "0\n"}}
	if _, err := NewDocker(f).Find(context.Background(), "task-1"); err != nil {
		t.Fatalf("find: %v", err)
	}
	if !has(f.argsFor("ps"), "-a") {
		t.Fatalf("docker ps must use -a: %v", f.argsFor("ps"))
	}
	if !has(f.argsFor("ps"), "label="+LabelTask+"=task-1") {
		t.Fatalf("lookup must filter by task label: %v", f.argsFor("ps"))
	}
}

// Two containers for one task means something reused an id. Picking one would
// let a half-finished run be reported as another's success.
func TestFindRefusesAmbiguity(t *testing.T) {
	f := &fakeCmd{outputs: map[string]string{"ps": "abc running\ndef exited\n"}}
	if _, err := NewDocker(f).Find(context.Background(), "task-1"); err == nil {
		t.Fatal("two containers for one task should be refused, not resolved by guessing")
	}
}

func TestFollowStreamsAndReturnsExitCode(t *testing.T) {
	f := &fakeCmd{
		stream:  "2026-08-13T10:00:00.000000000Z first line\n2026-08-13T10:00:01.000000000Z second line\n",
		outputs: map[string]string{"wait": "0\n"},
	}

	var got []string
	var lastTS string
	code, err := NewDocker(f).Follow(context.Background(), "abc", "", func(_, ts, line string) error {
		got = append(got, line)
		lastTS = ts
		return nil
	})
	if err != nil {
		t.Fatalf("follow: %v", err)
	}
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	if len(got) != 2 || got[0] != "first line" || got[1] != "second line" {
		t.Fatalf("lines = %v", got)
	}
	// The timestamp is what a resumed follow uses as --since, so it must be
	// separated from the line rather than left glued to it.
	if lastTS != "2026-08-13T10:00:01.000000000Z" {
		t.Fatalf("timestamp = %q", lastTS)
	}
	if !has(f.argsFor("logs"), "--follow", "--timestamps") {
		t.Fatalf("logs args: %v", f.argsFor("logs"))
	}
}

// Adoption resumes from where the previous agent stopped ingesting rather than
// replaying the whole container log into the outbox.
func TestFollowResumesFromCursor(t *testing.T) {
	f := &fakeCmd{stream: "", outputs: map[string]string{"wait": "0\n"}}
	_, err := NewDocker(f).Follow(context.Background(), "abc", "2026-08-13T10:00:00Z",
		func(string, string, string) error { return nil })
	if err != nil {
		t.Fatalf("follow: %v", err)
	}
	if !has(f.argsFor("logs"), "--since", "2026-08-13T10:00:00Z") {
		t.Fatalf("resume cursor not passed: %v", f.argsFor("logs"))
	}
}

// The log stream ending is not the container ending. Trusting the stream's own
// exit status would report success for a container that is still running or
// that failed after the pipe closed.
func TestFollowAsksDockerForTheRealExitCode(t *testing.T) {
	f := &fakeCmd{stream: "2026-08-13T10:00:00Z x\n", outputs: map[string]string{"wait": "137\n"}}
	code, err := NewDocker(f).Follow(context.Background(), "abc", "", func(string, string, string) error { return nil })
	if err != nil {
		t.Fatalf("follow: %v", err)
	}
	if code != 137 {
		t.Fatalf("exit code = %d, want 137 from docker wait", code)
	}
	if f.argsFor("wait") == nil {
		t.Fatal("docker wait was never called")
	}
}

// A full outbox stops the pump instead of growing without bound.
func TestFollowStopsWhenTheSinkRefuses(t *testing.T) {
	f := &fakeCmd{stream: "2026-08-13T10:00:00Z a\n2026-08-13T10:00:01Z b\n", outputs: map[string]string{"wait": "0\n"}}
	stop := errors.New("outbox full")

	seen := 0
	_, err := NewDocker(f).Follow(context.Background(), "abc", "", func(string, string, string) error {
		seen++
		return stop
	})
	if !errors.Is(err, stop) {
		t.Fatalf("err = %v, want the sink's error", err)
	}
	if seen != 1 {
		t.Fatalf("kept pumping after the sink refused: %d lines", seen)
	}
}
