package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"golang.org/x/crypto/ssh"
)

type BootstrapPayload struct {
	Host       string `json:"host"`         // e.g. "10.0.0.10:22" or "host.example.com:22"
	User       string `json:"user"`         // SSH user
	PrivKeyPEM string `json:"priv_key_pem"` // (or use agent; store securely!)
	AddUser    string `json:"add_user"`     // e.g. "deploy" (optional)
	NonInt     bool   `json:"non_interactive"`
	Script     string `json:"script"` // the content of bootstrap_host.sh (embed or load from disk)
}

func handleBootstrap(ctx context.Context, payloadRaw []byte) error {
	var p BootstrapPayload
	if err := json.Unmarshal(payloadRaw, &p); err != nil {
		return err
	}

	signer, err := ssh.ParsePrivateKey([]byte(p.PrivKeyPEM))
	if err != nil {
		return fmt.Errorf("parse key: %w", err)
	}

	cfg := &ssh.ClientConfig{
		User:            p.User,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // or verify!
		Timeout:         30 * time.Second,
	}

	client, err := ssh.Dial("tcp", p.Host, cfg)
	if err != nil {
		return fmt.Errorf("ssh dial: %w", err)
	}
	defer client.Close()

	sess, err := client.NewSession()
	if err != nil {
		return err
	}
	defer sess.Close()

	stdin, _ := sess.StdinPipe()
	go func() {
		defer stdin.Close()
		io.WriteString(stdin, p.Script)
	}()

	args := "-s --"
	if p.AddUser != "" {
		args += " -u " + p.AddUser
	}
	if p.NonInt {
		args += " -y"
	}
	cmd := "sudo bash " + args

	if err := sess.Run(cmd); err != nil {
		return fmt.Errorf("bootstrap failed: %w", err)
	}
	return nil
}
