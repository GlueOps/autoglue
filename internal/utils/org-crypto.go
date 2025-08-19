package utils

import (
	"fmt"

	"github.com/glueops/autoglue/internal/db"
	"github.com/glueops/autoglue/internal/db/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// EncryptForOrg encrypts plaintext for an org using its TenantKey.
// Returns base64-encoded ciphertext, iv, and tag.
func EncryptForOrg(orgID uuid.UUID, plaintext []byte) (cipherB64, ivB64, tagB64 string, err error) {
	tenantKey, err := getOrCreateTenantKey(orgID.String())
	if err != nil {
		return "", "", "", err
	}
	ct, iv, tag, err := encryptAESGCM(plaintext, tenantKey)
	if err != nil {
		return "", "", "", err
	}
	return encodeB64(ct), encodeB64(iv), encodeB64(tag), nil
}

// DecryptForOrg decrypts b64 cipher/iv/tag using the org's TenantKey and returns plaintext.
func DecryptForOrg(orgID uuid.UUID, cipherB64, ivB64, tagB64 string) (string, error) {
	tenantKey, err := getOrCreateTenantKey(orgID.String())
	if err != nil {
		return "", err
	}
	ct, err := decodeB64(cipherB64)
	if err != nil {
		return "", fmt.Errorf("decode cipher: %w", err)
	}
	iv, err := decodeB64(ivB64)
	if err != nil {
		return "", fmt.Errorf("decode iv: %w", err)
	}
	tag, err := decodeB64(tagB64)
	if err != nil {
		return "", fmt.Errorf("decode tag: %w", err)
	}
	plain, err := decryptAESGCM(ct, tenantKey, iv, tag)
	if err != nil {
		return "", err
	}
	return string(plain), nil
}

// Optional convenience used in some flows (upsert by provider).
func EncryptAndUpsertCredential(orgID uuid.UUID, provider, plaintext string) error {
	data, iv, tag, err := EncryptForOrg(orgID, []byte(plaintext))
	if err != nil {
		return err
	}

	var cred models.Credential
	err = db.DB.Where("organization_id = ? AND provider = ?", orgID, provider).
		First(&cred).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	cred.OrganizationID = orgID
	cred.Provider = provider
	cred.EncryptedData = data
	cred.IV = iv
	cred.Tag = tag

	if cred.ID != uuid.Nil {
		return db.DB.Save(&cred).Error
	}
	return db.DB.Create(&cred).Error
}
