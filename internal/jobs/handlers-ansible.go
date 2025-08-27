package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
)

type AnsiblePayload struct {
	Image         string   `json:"image"`           // e.g. "ansible-debian:latest"
	PlaybookPath  string   `json:"playbook_path"`   // inside your host repo (mounted)
	InventoryHost string   `json:"inventory_host"`  // "user@host,"
	SSHKeyPath    string   `json:"ssh_key_path"`    // host path to key (mounted read-only)
	WorkspaceHost string   `json:"workspace_host"`  // host path to playbooks to mount
	ExtraVars     []string `json:"extra_vars"`      // ["foo=bar","env=prod"]
	SSHCommonArgs []string `json:"ssh_common_args"` // e.g. ["-o","StrictHostKeyChecking=no"]
}

func handleAnsible(ctx context.Context, payloadRaw []byte) error {
	var p AnsiblePayload
	if err := json.Unmarshal(payloadRaw, &p); err != nil {
		return err
	}

	args := []string{
		"run", "--rm",
		"-v", p.WorkspaceHost + ":/workspace",
		"-v", p.SSHKeyPath + ":/root/.ssh/id_rsa:ro",
		"-e", "ANSIBLE_HOST_KEY_CHECKING=False",
		p.Image,
		"bash", "-lc",
	}

	extras := ""
	for _, kv := range p.ExtraVars {
		extras += fmt.Sprintf(" -e %q", kv)
	}
	sshArgs := ""
	if len(p.SSHCommonArgs) > 0 {
		for _, a := range p.SSHCommonArgs {
			sshArgs += " " + a
		}
	}

	cmdStr := fmt.Sprintf(
		`chmod 600 /root/.ssh/id_rsa && ansible-playbook -i %q %s --ssh-common-args=%q%s`,
		p.InventoryHost,
		p.PlaybookPath,
		joinArgs(p.SSHCommonArgs), // produces a single string; same as sshArgs but quoted
		extras,
	)
	args = append(args, cmdStr)

	c := exec.CommandContext(ctx, "docker", args...)
	out, err := c.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ansible run failed: %v\n%s", err, string(out))
	}
	return nil
}

func joinArgs(a []string) string {
	// turn ["-o","StrictHostKeyChecking=no"] into "-o StrictHostKeyChecking=no"
	out := ""
	for i, s := range a {
		if i > 0 {
			out += " "
		}
		out += s
	}
	return out
}
