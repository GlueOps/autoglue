package web

import (
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
)

//go:embed pgwebbin/*
var pgwebFS embed.FS

type pgwebAsset struct {
	Path   string
	SHA256 string
}

var pgwebIndex = map[string]pgwebAsset{
	"linux/amd64":  {Path: "pgwebbin/pgweb-linux-amd64", SHA256: ""},
	"linux/arm64":  {Path: "pgwebbin/pgweb-linux-arm64", SHA256: ""},
	"darwin/amd64": {Path: "pgwebbin/pgweb-darwin-amd64", SHA256: ""},
	"darwin/arm64": {Path: "pgwebbin/pgweb-darwin-arm64", SHA256: ""},
}

func ExtractPgweb() (string, error) {
	key := runtime.GOOS + "/" + runtime.GOARCH
	as, ok := pgwebIndex[key]
	if !ok {
		return "", fmt.Errorf("pgweb not embedded for %s", key)
	}
	f, err := pgwebFS.Open(as.Path)
	if err != nil {
		return "", fmt.Errorf("embedded pgweb missing: %w", err)
	}
	defer f.Close()

	tmpDir, err := os.MkdirTemp("", "pgweb-*")
	if err != nil {
		return "", err
	}

	filename := "pgweb"
	if runtime.GOOS == "windows" {
		filename += ".exe"
	}
	outPath := filepath.Join(tmpDir, filename)

	out, err := os.OpenFile(outPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o700)
	if err != nil {
		return "", err
	}
	defer out.Close()

	h := sha256.New()
	if _, err = io.Copy(io.MultiWriter(out, h), f); err != nil {
		return "", err
	}

	if as.SHA256 != "" {
		got := hex.EncodeToString(h.Sum(nil))
		if got != as.SHA256 {
			return "", fmt.Errorf("pgweb checksum mismatch: got=%s want=%s", got, as.SHA256)
		}
	}

	// Make sure it’s executable on Unix; Windows ignores this.
	_ = os.Chmod(outPath, 0o700)
	return outPath, nil
}

func CleanupPgweb(pgwebPath string) error {
	if pgwebPath == "" {
		return nil
	}
	dir := filepath.Dir(pgwebPath)
	if dir == "" || dir == "/" || dir == "." {
		return errors.New("refusing to remove suspicious directory")
	}
	return os.RemoveAll(dir)
}
