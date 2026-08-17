package runner

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
)

// Commander runs an external command. Injected so the docker logic is testable
// without a docker daemon — the argument construction and output parsing are
// where the bugs live, not in exec itself.
type Commander interface {
	Output(ctx context.Context, name string, args ...string) (string, error)
	Stream(ctx context.Context, name string, args ...string) (io.ReadCloser, func() error, error)
}

// Docker drives the docker CLI rather than the Engine API. The CLI is what is
// already on a bastion, and the SDK would pull a large dependency tree into a
// module whose whole purpose is to stay small enough to ship anywhere.
type Docker struct {
	cmd Commander
	bin string
}

func NewDocker(cmd Commander) *Docker {
	return &Docker{cmd: cmd, bin: "docker"}
}

func (d *Docker) Find(ctx context.Context, taskID string) (Container, error) {
	// -a so an exited container is found too: a task whose container has
	// finished but whose result was never delivered is exactly what adoption
	// exists to recover.
	out, err := d.cmd.Output(ctx, d.bin, "ps", "-a",
		"--filter", "label="+LabelTask+"="+taskID,
		"--format", "{{.ID}} {{.State}}")
	if err != nil {
		return Container{}, fmt.Errorf("docker ps for task %s: %w", taskID, err)
	}

	line := strings.TrimSpace(out)
	if line == "" {
		return Container{Phase: PhaseAbsent}, nil
	}
	// More than one container for a task means something reused the id. Refuse
	// rather than pick: guessing which one is authoritative is how a half-done
	// run gets reported as someone else's success.
	if strings.Contains(line, "\n") {
		return Container{}, fmt.Errorf("task %s has %d containers; refusing to guess",
			taskID, strings.Count(line, "\n")+1)
	}

	id, state, ok := strings.Cut(line, " ")
	if !ok {
		return Container{}, fmt.Errorf("unparsable docker ps output for %s: %q", taskID, line)
	}

	c := Container{ID: id}
	switch state {
	case "running", "created", "restarting":
		c.Phase = PhaseRunning
	default:
		c.Phase = PhaseExited
		code, err := d.exitCode(ctx, id)
		if err != nil {
			return Container{}, err
		}
		c.ExitCode = code
	}
	return c, nil
}

func (d *Docker) exitCode(ctx context.Context, containerID string) (int, error) {
	out, err := d.cmd.Output(ctx, d.bin, "inspect", "-f", "{{.State.ExitCode}}", containerID)
	if err != nil {
		return 0, fmt.Errorf("inspect exit code of %s: %w", containerID, err)
	}
	code, err := strconv.Atoi(strings.TrimSpace(out))
	if err != nil {
		return 0, fmt.Errorf("exit code of %s is not a number (%q): %w", containerID, out, err)
	}
	return code, nil
}

func (d *Docker) Start(ctx context.Context, spec Spec) (string, error) {
	args := []string{
		"run", "--detach",
		"--label", LabelCluster + "=" + spec.ClusterID,
		"--label", LabelRun + "=" + spec.RunID,
		"--label", LabelTask + "=" + spec.TaskID,
	}
	// Deliberately no --rm. Once the container outlives the agent, it is the
	// only record of what happened after the stream was lost, and --rm would
	// delete exactly the evidence adoption came back for.
	if spec.WorkDir != "" {
		args = append(args, "--workdir", spec.WorkDir)
	}
	for _, m := range spec.Mounts {
		args = append(args, "--volume", m.Source+":"+m.Target)
	}
	args = append(args, spec.Image+":"+spec.Tag, "make", spec.Target)

	out, err := d.cmd.Output(ctx, d.bin, args...)
	if err != nil {
		return "", fmt.Errorf("docker run for task %s: %w", spec.TaskID, err)
	}
	id := strings.TrimSpace(out)
	if id == "" {
		return "", fmt.Errorf("docker run for task %s returned no container id", spec.TaskID)
	}
	return id, nil
}

func (d *Docker) Follow(ctx context.Context, containerID, since string, onLog LogFunc) (int, error) {
	args := []string{"logs", "--follow", "--timestamps"}
	if since != "" {
		// Resume rather than replay. The boundary is only as precise as docker's
		// timestamps, so a line may repeat at the seam — duplicated output is a
		// far better failure than a silent hole in a transcript someone is
		// using to work out what broke.
		args = append(args, "--since", since)
	}
	args = append(args, containerID)

	stream, wait, err := d.cmd.Stream(ctx, d.bin, args...)
	if err != nil {
		return 0, fmt.Errorf("docker logs for %s: %w", containerID, err)
	}

	scanner := bufio.NewScanner(stream)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		ts, line, _ := strings.Cut(scanner.Text(), " ")
		if err := onLog("stdout", ts, line); err != nil {
			_ = stream.Close()
			_ = wait()
			return 0, err
		}
	}
	scanErr := scanner.Err()
	_ = stream.Close()
	// The log stream ending is not the container ending, and its exit status
	// says nothing about the work. Ask docker for the real code.
	_ = wait()
	if scanErr != nil {
		return 0, fmt.Errorf("reading logs of %s: %w", containerID, scanErr)
	}

	out, err := d.cmd.Output(ctx, d.bin, "wait", containerID)
	if err != nil {
		return 0, fmt.Errorf("docker wait %s: %w", containerID, err)
	}
	code, err := strconv.Atoi(strings.TrimSpace(out))
	if err != nil {
		return 0, fmt.Errorf("docker wait %s returned %q: %w", containerID, out, err)
	}
	return code, nil
}

func (d *Docker) Remove(ctx context.Context, containerID string) error {
	if _, err := d.cmd.Output(ctx, d.bin, "rm", containerID); err != nil {
		return fmt.Errorf("docker rm %s: %w", containerID, err)
	}
	return nil
}

// ExecCommander is the real Commander, shelling out to the docker CLI.
type ExecCommander struct{}

func (ExecCommander) Output(ctx context.Context, name string, args ...string) (string, error) {
	out, err := exec.CommandContext(ctx, name, args...).Output()
	if err != nil {
		// docker puts the useful part ("No such container", "permission
		// denied") on stderr, and bare *exec.ExitError renders as "exit status
		// 1", which tells a reader nothing.
		var ee *exec.ExitError
		if errors.As(err, &ee) && len(ee.Stderr) > 0 {
			return "", fmt.Errorf("%w: %s", err, strings.TrimSpace(string(ee.Stderr)))
		}
		return "", err
	}
	return string(out), nil
}

func (ExecCommander) Stream(ctx context.Context, name string, args ...string) (io.ReadCloser, func() error, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	// Container output arrives on both, and a transcript that drops stderr
	// loses precisely the lines someone is looking for.
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}
	cmd.Stderr = cmd.Stdout
	if err := cmd.Start(); err != nil {
		return nil, nil, err
	}
	return stdout, cmd.Wait, nil
}
