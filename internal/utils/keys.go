package utils

import (
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/glueops/autoglue/internal/db"
	"github.com/glueops/autoglue/internal/db/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func getMasterKey() ([]byte, error) {
	var mk models.MasterKey
	if err := db.DB.Where("is_active = ?", true).Order("created_at DESC").First(&mk).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNoActiveMasterKey
		}
		return nil, fmt.Errorf("querying master key: %w", err)
	}

	keyBytes, err := base64.StdEncoding.DecodeString(mk.Key)
	if err != nil {
		return nil, fmt.Errorf("decoding master key: %w", err)
	}
	if len(keyBytes) != 32 {
		return nil, fmt.Errorf("%w: got %d, want 32", ErrInvalidMasterKeyLen, len(keyBytes))
	}
	return keyBytes, nil
}

func getOrCreateTenantKey(orgID string) ([]byte, error) {
	var orgKey models.OrganizationKey
	err := db.DB.Where("organization_id = ?", orgID).First(&orgKey).Error
	if err == nil {
		encKeyB64 := orgKey.EncryptedKey
		ivB64 := orgKey.IV
		tagB64 := orgKey.Tag

		encryptedKey, err := decodeB64(encKeyB64)
		if err != nil {
			return nil, fmt.Errorf("decode enc key: %w", err)
		}
		iv, err := decodeB64(ivB64)
		if err != nil {
			return nil, fmt.Errorf("decode iv: %w", err)
		}
		tag, err := decodeB64(tagB64)
		if err != nil {
			return nil, fmt.Errorf("decode tag: %w", err)
		}
		masterKey, err := getMasterKey()
		if err != nil {
			return nil, err
		}
		return decryptAESGCM(encryptedKey, masterKey, iv, tag)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Create new tenant key and wrap with the current master key
	orgUUID, err := uuid.Parse(orgID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidOrgID, err)
	}

	tenantKey, err := randomBytes(32)
	if err != nil {
		return nil, fmt.Errorf("tenant key gen: %w", err)
	}
	masterKey, err := getMasterKey()
	if err != nil {
		return nil, err
	}
	encrypted, iv, tag, err := encryptAESGCM(tenantKey, masterKey)
	if err != nil {
		return nil, fmt.Errorf("wrap tenant key: %w", err)
	}

	var mk models.MasterKey
	if err := db.DB.Where("is_active = ?", true).Order("created_at DESC").First(&mk).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNoActiveMasterKey
		}
		return nil, fmt.Errorf("querying master key: %w", err)
	}

	orgKey = models.OrganizationKey{
		OrganizationID: orgUUID,
		MasterKeyID:    mk.ID,
		EncryptedKey:   encodeB64(encrypted),
		IV:             encodeB64(iv),
		Tag:            encodeB64(tag),
	}
	if err := db.DB.Create(&orgKey).Error; err != nil {
		return nil, fmt.Errorf("persist org key: %w", err)
	}
	return tenantKey, nil
}
