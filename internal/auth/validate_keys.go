package auth

import (
	"time"

	"github.com/glueops/autoglue/internal/models"
	"gorm.io/gorm"
)

// ValidateAPIKey validates a single-token user API key sent via X-API-KEY.
func ValidateAPIKey(rawKey string, db *gorm.DB) *models.User {
	if rawKey == "" {
		return nil
	}
	digest := SHA256Hex(rawKey)

	var k models.APIKey
	if err := db.
		Where("key_hash = ? AND scope = ? AND (expires_at IS NULL OR expires_at > ?)", digest, "user", time.Now()).
		First(&k).Error; err != nil {
		return nil
	}
	if k.UserID == nil {
		return nil
	}
	var u models.User
	if err := db.First(&u, "id = ? AND is_disabled = false", *k.UserID).Error; err != nil {
		return nil
	}
	// Optional: touch last_used_at here if you've added it on the model.
	return &u
}

// ValidateAppKeyPair validates a user key/secret pair via X-APP-KEY / X-APP-SECRET.
func ValidateAppKeyPair(appKey, secret string, db *gorm.DB) *models.User {
	if appKey == "" || secret == "" {
		return nil
	}
	digest := SHA256Hex(appKey)

	var k models.APIKey
	if err := db.
		Where("key_hash = ? AND scope = ? AND (expires_at IS NULL OR expires_at > ?)", digest, "user", time.Now()).
		First(&k).Error; err != nil {
		return nil
	}
	ok, _ := VerifySecretArgon2id(zeroIfNil(k.SecretHash), secret)
	if !ok || k.UserID == nil {
		return nil
	}
	var u models.User
	if err := db.First(&u, "id = ? AND is_disabled = false", *k.UserID).Error; err != nil {
		return nil
	}
	return &u
}

// ValidateOrgKeyPair validates an org key/secret via X-ORG-KEY / X-ORG-SECRET.
func ValidateOrgKeyPair(orgKey, secret string, db *gorm.DB) *models.Organization {
	if orgKey == "" || secret == "" {
		return nil
	}
	digest := SHA256Hex(orgKey)

	var k models.APIKey
	if err := db.
		Where("key_hash = ? AND scope = ? AND (expires_at IS NULL OR expires_at > ?)", digest, "org", time.Now()).
		First(&k).Error; err != nil {
		return nil
	}
	ok, _ := VerifySecretArgon2id(zeroIfNil(k.SecretHash), secret)
	if !ok || k.OrgID == nil {
		return nil
	}
	var o models.Organization
	if err := db.First(&o, "id = ?", *k.OrgID).Error; err != nil {
		return nil
	}
	return &o
}

// local helper; avoids nil-deref when comparing secrets
func zeroIfNil(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
