package utils

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/glueops/autoglue/internal/db"
	"github.com/glueops/autoglue/internal/db/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/ssh"
)

func generateSSHKeys(orgID string) (string, string, error) {
	var key models.SshKey
	db.DB.Where("organization_id = ?", orgID).First(&key)
	if key.ID != uuid.Nil {
		return key.PublicKey, key.PrivateKey, nil
	}

	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)

	privBuf := new(bytes.Buffer)
	pem.Encode(privBuf, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	pubKey, _ := ssh.NewPublicKey(&privateKey.PublicKey)
	pubKeyString := string(ssh.MarshalAuthorizedKey(pubKey))

	orgUUID, err := uuid.Parse(orgID)
	if err != nil {
		return "", "", fmt.Errorf("invalid organization ID: %w", err)
	}

	newKey := models.SshKey{
		OrganizationID: orgUUID,
		PublicKey:      pubKeyString,
		PrivateKey:     privBuf.String(),
	}
	db.DB.Create(&newKey)

	return pubKeyString, privBuf.String(), nil
}
