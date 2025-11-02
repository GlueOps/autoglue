package auth

import (
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"sync"
	"time"

	"github.com/glueops/autoglue/internal/keys"
	"github.com/glueops/autoglue/internal/models"
	"gorm.io/gorm"
)

type keyCache struct {
	mu      sync.RWMutex
	pub     map[string]interface{} // kid -> public key object
	meta    map[string]models.SigningKey
	selKid  string
	selAlg  string
	selPriv any
}

var kc keyCache

// Refresh loads active keys into memory. Call on startup and periodically (ticker/cron).
func Refresh(db *gorm.DB, encKeyB64 string) error {
	var rows []models.SigningKey
	if err := db.Where("is_active = true AND (expires_at IS NULL OR expires_at > ?)", time.Now()).
		Order("created_at desc").Find(&rows).Error; err != nil {
		return err
	}

	pub := make(map[string]interface{}, len(rows))
	meta := make(map[string]models.SigningKey, len(rows))
	var selKid string
	var selAlg string
	var selPriv any

	for i, r := range rows {
		// parse public
		block, _ := pem.Decode([]byte(r.PublicPEM))
		if block == nil {
			continue
		}
		var pubKey any
		switch r.Alg {
		case "RS256", "RS384", "RS512":
			pubKey, _ = x509.ParsePKCS1PublicKey(block.Bytes)
			if pubKey == nil {
				// also allow PKIX format
				if k, err := x509.ParsePKIXPublicKey(block.Bytes); err == nil {
					pubKey = k
				}
			}
		case "EdDSA":
			k, err := x509.ParsePKIXPublicKey(block.Bytes)
			if err == nil {
				if edk, ok := k.(ed25519.PublicKey); ok {
					pubKey = edk
				}
			}
		}
		if pubKey == nil {
			continue
		}
		pub[r.Kid] = pubKey
		meta[r.Kid] = r

		// pick first row as current signer (most recent because of order desc)
		if i == 0 {
			privPEM := r.PrivatePEM
			// decrypt if necessary
			if len(privPEM) > 10 && privPEM[:10] == "enc:aesgcm" {
				pt, err := keysDecrypt(encKeyB64, privPEM)
				if err != nil {
					continue
				}
				privPEM = string(pt)
			}
			blockPriv, _ := pem.Decode([]byte(privPEM))
			if blockPriv == nil {
				continue
			}
			switch r.Alg {
			case "RS256", "RS384", "RS512":
				if k, err := x509.ParsePKCS1PrivateKey(blockPriv.Bytes); err == nil {
					selPriv = k
					selAlg = r.Alg
					selKid = r.Kid
				} else if kAny, err := x509.ParsePKCS8PrivateKey(blockPriv.Bytes); err == nil {
					if k, ok := kAny.(*rsa.PrivateKey); ok {
						selPriv = k
						selAlg = r.Alg
						selKid = r.Kid
					}
				}
			case "EdDSA":
				if kAny, err := x509.ParsePKCS8PrivateKey(blockPriv.Bytes); err == nil {
					if k, ok := kAny.(ed25519.PrivateKey); ok {
						selPriv = k
						selAlg = r.Alg
						selKid = r.Kid
					}
				}
			}
		}
	}

	kc.mu.Lock()
	defer kc.mu.Unlock()
	kc.pub = pub
	kc.meta = meta
	kc.selKid = selKid
	kc.selAlg = selAlg
	kc.selPriv = selPriv
	return nil
}

func keysDecrypt(encKey, enc string) ([]byte, error) {
	return keysDecryptImpl(encKey, enc)
}

// indirection for same package
var keysDecryptImpl = func(encKey, enc string) ([]byte, error) {
	return nil, errors.New("not wired")
}

// Wire up from keys package
func init() {
	keysDecryptImpl = keysDecryptShim
}

func keysDecryptShim(encKey, enc string) ([]byte, error) {
	return keys.Decrypt(encKey, enc)
}
