package ssh

import (
	"archive/zip"
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
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

// GenerateRSA returns a new RSA private key.
func GenerateRSA(bits int) (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, bits)
}

// RSAPrivateToPEMAndAuthorized encodes the private key to PEM and the public key to authorized_keys,
// appending an optional comment to the authorized key.
func RSAPrivateToPEMAndAuthorized(priv *rsa.PrivateKey, comment string) (privPEM string, authorized string, err error) {
	// Private (PKCS#1) to PEM
	der := x509.MarshalPKCS1PrivateKey(priv)
	block := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: der}
	var buf bytes.Buffer
	if err = pem.Encode(&buf, block); err != nil {
		return "", "", err
	}

	// Public to authorized_keys
	pub, err := gossh.NewPublicKey(&priv.PublicKey)
	if err != nil {
		return "", "", err
	}
	auth := strings.TrimSpace(string(gossh.MarshalAuthorizedKey(pub)))
	comment = strings.TrimSpace(comment)
	if comment != "" {
		auth += " " + comment
	}
	return buf.String(), auth, nil
}

// Convenience wrapper used by CreateSSHKey
func GenerateRSAPEMAndAuthorized(bits int, comment string) (string, string, error) {
	priv, err := GenerateRSA(bits)
	if err != nil {
		return "", "", err
	}
	return RSAPrivateToPEMAndAuthorized(priv, comment)
}
func toZipFile(filename string, content []byte, zw *zip.Writer) error {
	f, err := zw.Create(filename)
	if err != nil {
		return err
	}
	_, err = f.Write(content)
	return err
}
