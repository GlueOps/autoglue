package ssh

import (
	"archive/zip"
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"strings"

	gossh "golang.org/x/crypto/ssh"
)

func allowedBits(b int) bool {
	return b == 2048 || b == 3072 || b == 4096
}

func GenerateRSA(bits int) (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, bits)
}

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
