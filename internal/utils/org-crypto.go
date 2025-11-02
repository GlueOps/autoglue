package utils

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func EncryptForOrg(orgID uuid.UUID, plaintext []byte, db *gorm.DB) (cipherB64, ivB64, tagB64 string, err error) {
	tenantKey, err := getOrCreateTenantKey(orgID.String(), db)
	if err != nil {
		return "", "", "", err
	}
	ct, iv, tag, err := encryptAESGCM(plaintext, tenantKey)
	if err != nil {
		return "", "", "", err
	}
	return EncodeB64(ct), EncodeB64(iv), EncodeB64(tag), nil
}

func DecryptForOrg(orgID uuid.UUID, cipherB64, ivB64, tagB64 string, db *gorm.DB) (string, error) {
	tenantKey, err := getOrCreateTenantKey(orgID.String(), db)
	if err != nil {
		return "", err
	}
	ct, err := DecodeB64(cipherB64)
	if err != nil {
		return "", fmt.Errorf("decode cipher: %w", err)
	}
	iv, err := DecodeB64(ivB64)
	if err != nil {
		return "", fmt.Errorf("decode iv: %w", err)
	}
	tag, err := DecodeB64(tagB64)
	if err != nil {
		return "", fmt.Errorf("decode tag: %w", err)
	}
	plain, err := decryptAESGCM(ct, tenantKey, iv, tag)
	if err != nil {
		return "", err
	}
	return string(plain), nil
}
