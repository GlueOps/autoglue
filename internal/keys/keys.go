package keys

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"time"

	"github.com/glueops/autoglue/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GenOpts struct {
	Alg  string // "RS256"|"RS384"|"RS512"|"EdDSA"
	Bits int    // RSA bits (2048/3072/4096). ignored for EdDSA
	KID  string // optional; if empty we generate one
	NBF  *time.Time
	EXP  *time.Time
}

func GenerateAndStore(db *gorm.DB, encKeyB64 string, opts GenOpts) (*models.SigningKey, error) {
	if opts.KID == "" {
		opts.KID = uuid.NewString()
	}

	var pubPEM, privPEM []byte
	var alg = opts.Alg

	switch alg {
	case "RS256", "RS384", "RS512":
		if opts.Bits == 0 {
			opts.Bits = 3072
		}
		priv, err := rsa.GenerateKey(rand.Reader, opts.Bits)
		if err != nil {
			return nil, err
		}
		privDER := x509.MarshalPKCS1PrivateKey(priv)
		privPEM = pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: privDER})

		pubDER := x509.MarshalPKCS1PublicKey(&priv.PublicKey)
		pubPEM = pem.EncodeToMemory(&pem.Block{Type: "RSA PUBLIC KEY", Bytes: pubDER})

	case "EdDSA":
		pub, priv, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return nil, err
		}
		privDER, err := x509.MarshalPKCS8PrivateKey(priv)
		if err != nil {
			return nil, err
		}
		privPEM = pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privDER})

		pubDER, err := x509.MarshalPKIXPublicKey(pub)
		if err != nil {
			return nil, err
		}
		pubPEM = pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubDER})

	default:
		return nil, fmt.Errorf("unsupported alg: %s", alg)
	}

	privateOut := string(privPEM)
	if encKeyB64 != "" {
		enc, err := encryptAESGCM(encKeyB64, privPEM)
		if err != nil {
			return nil, err
		}
		privateOut = enc
	}

	rec := models.SigningKey{
		Kid:        opts.KID,
		Alg:        alg,
		Use:        "sig",
		IsActive:   true,
		PublicPEM:  string(pubPEM),
		PrivatePEM: privateOut,
		NotBefore:  opts.NBF,
		ExpiresAt:  opts.EXP,
	}
	if err := db.Create(&rec).Error; err != nil {
		return nil, err
	}
	return &rec, nil
}

func encryptAESGCM(b64 string, plaintext []byte) (string, error) {
	key, err := decode32ByteKey(b64)
	if err != nil {
		return "", err
	}
	if len(key) != 32 {
		return "", errors.New("JWT_PRIVATE_ENC_KEY must be 32 bytes (base64url)")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, aead.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return "", err
	}
	out := aead.Seal(nonce, nonce, plaintext, nil)
	return "enc:aesgcm:" + base64.RawStdEncoding.EncodeToString(out), nil
}

func decryptAESGCM(b64 string, enc string) ([]byte, error) {
	if !bytes.HasPrefix([]byte(enc), []byte("enc:aesgcm:")) {
		return nil, errors.New("not encrypted")
	}
	key, err := decode32ByteKey(b64)
	if err != nil {
		return nil, err
	}
	blob, err := base64.RawStdEncoding.DecodeString(enc[len("enc:aesgcm:"):])
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := aead.NonceSize()
	if len(blob) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}
	nonce, ct := blob[:nonceSize], blob[nonceSize:]
	return aead.Open(nil, nonce, ct, nil)
}
