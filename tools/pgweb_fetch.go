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
			SHA256: "",
		},
		{
			Name:   "pgweb-linux-arm64",
			URL:    fmt.Sprintf("https://github.com/sosedoff/pgweb/releases/download/v%s/pgweb_linux_arm64.zip", version),
			SHA256: "",
		},
		{
			Name:   "pgweb-darwin-amd64",
			URL:    fmt.Sprintf("https://github.com/sosedoff/pgweb/releases/download/v%s/pgweb_darwin_amd64.zip", version),
			SHA256: "",
		},
		{
			Name:   "pgweb-darwin-arm64",
			URL:    fmt.Sprintf("https://github.com/sosedoff/pgweb/releases/download/v%s/pgweb_darwin_arm64.zip", version),
			SHA256: "",
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
	// minimal unzip: because pgweb zip has only one binary
	r, err := os.Open(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	// use archive/zip
	stat, err := os.Stat(zipPath)
	if err != nil {
		return err
	}
	return unzipFile(zipPath, outPath, stat.Size())
}

func unzipFile(zipFile, outFile string, _ int64) error {
	r, err := os.Open(zipFile)
	if err != nil {
		return err
	}
	defer r.Close()
	fi, _ := r.Stat()

	// rely on standard zip reader
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	tmpZip := filepath.Join(os.TempDir(), fi.Name())
	if err := os.WriteFile(tmpZip, data, 0o644); err != nil {
		return err
	}
	defer os.Remove(tmpZip)

	zr, err := os.Open(tmpZip)
	if err != nil {
		return err
	}
	defer zr.Close()
	// extract using standard lib
	zr2, err := zip.OpenReader(tmpZip)
	if err != nil {
		return err
	}
	defer zr2.Close()
	for _, f := range zr2.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()
		out, err := os.Create(outFile)
		if err != nil {
			return err
		}
		if _, err := io.Copy(out, rc); err != nil {
			out.Close()
			return err
		}
		out.Close()
		break
	}
	return nil
}
