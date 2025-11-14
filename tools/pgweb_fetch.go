//go:build ignore
// +build ignore

package main

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type Target struct {
	Name   string
	URL    string
	SHA256 string
}

const version = "0.16.2"

func main() {
	targets := []Target{
		{
			Name:   "pgweb-linux-amd64",
			URL:    fmt.Sprintf("https://github.com/sosedoff/pgweb/releases/download/v%s/pgweb_linux_amd64.zip", version),
			SHA256: "3d6c2063e1040b8a625eb7c43c9b84f8ed12cfc9a798eacbce85179963ee2554",
		},
		{
			Name:   "pgweb-linux-arm64",
			URL:    fmt.Sprintf("https://github.com/sosedoff/pgweb/releases/download/v%s/pgweb_linux_arm64.zip", version),
			SHA256: "079c698a323ed6431ce7e6343ee5847c7da62afbf45dfb2e78f8289d7b381783",
		},
		{
			Name:   "pgweb-darwin-amd64",
			URL:    fmt.Sprintf("https://github.com/sosedoff/pgweb/releases/download/v%s/pgweb_darwin_amd64.zip", version),
			SHA256: "c0a098e2eb9cf9f7c20161a2947522eb67eacbf2b6c3389c2f8e8c5ed7238957",
		},
		{
			Name:   "pgweb-darwin-arm64",
			URL:    fmt.Sprintf("https://github.com/sosedoff/pgweb/releases/download/v%s/pgweb_darwin_arm64.zip", version),
			SHA256: "c8f5fca847f461ba22a619e2d96cb1656cefdffd8f2aef2340e14fc5b518d3a2",
		},
	}

	outDir := filepath.Join("internal", "web", "pgwebbin")
	_ = os.MkdirAll(outDir, 0o755)

	for _, t := range targets {
		destZip := filepath.Join(outDir, t.Name+".zip")
		fmt.Printf("Downloading %s...\n", t.URL)
		if err := downloadFile(destZip, t.URL); err != nil {
			panic(err)
		}
		binPath := filepath.Join(outDir, t.Name)
		if err := unzipSingle(destZip, binPath); err != nil {
			panic(err)
		}
		_ = os.Remove(destZip)

		// Make executable
		if err := os.Chmod(binPath, 0o755); err != nil {
			panic(err)
		}
		fmt.Printf("Saved %s\n", binPath)

		// Compute checksum
		sum, _ := fileSHA256(binPath)
		fmt.Printf("  SHA256: %s\n", sum)
	}
}

func downloadFile(dest, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("bad status: %s", resp.Status)
	}
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}

func fileSHA256(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func unzipSingle(zipPath, outPath string) error {
	zr, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer zr.Close()

	if len(zr.File) == 0 {
		return fmt.Errorf("zip file %s is empty", zipPath)
	}

	f := zr.File[0]

	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	out, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, rc); err != nil {
		return err
	}

	return nil
}
