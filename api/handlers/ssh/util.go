package ssh

import (
	"archive/zip"
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"
	"strings"

	gossh "golang.org/x/crypto/ssh"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func allowedBits(b int) bool {
	return b == 2048 || b == 3072 || b == 4096
}

func generateRSA(bits int, comment string) (privPEM string, pubAuthorized string, err error) {
	if !allowedBits(bits) {
		bits = 4096
	}
	key, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return "", "", fmt.Errorf("generate rsa: %w", err)
	}

	// Private: PEM (PKCS#1)
	privDER := x509.MarshalPKCS1PrivateKey(key)
	privPEMBlock := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: privDER}
	var privBuf bytes.Buffer
	if err := pem.Encode(&privBuf, privPEMBlock); err != nil {
		return "", "", fmt.Errorf("encode pem: %w", err)
	}

	// Public: authorized_keys
	pub, err := gossh.NewPublicKey(&key.PublicKey)
	if err != nil {
		return "", "", fmt.Errorf("ssh pub: %w", err)
	}
	auth := strings.TrimSpace(string(gossh.MarshalAuthorizedKey(pub)))
	if comment = strings.TrimSpace(comment); comment != "" {
		auth = auth + " " + comment
	}

	return privBuf.String(), auth, nil
}

func toZipFile(filename string, content []byte, zw *zip.Writer) error {
	f, err := zw.Create(filename)
	if err != nil {
		return err
	}
	_, err = f.Write(content)
	return err
}
