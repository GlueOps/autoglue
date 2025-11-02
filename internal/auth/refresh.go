package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/glueops/autoglue/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// random opaque token (returned to client once)
func generateOpaqueToken(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

type RefreshPair struct {
	Plain  string
	Record models.RefreshToken
}

// Issue a new refresh token (new family if familyID == nil)
func IssueRefreshToken(db *gorm.DB, userID uuid.UUID, ttl time.Duration, familyID *uuid.UUID) (RefreshPair, error) {
	plain, err := generateOpaqueToken(32)
	if err != nil {
		return RefreshPair{}, err
	}
	hash, err := HashSecretArgon2id(plain)
	if err != nil {
		return RefreshPair{}, err
	}

	fid := uuid.New()
	if familyID != nil {
		fid = *familyID
	}

	rec := models.RefreshToken{
		UserID:    userID,
		FamilyID:  fid,
		TokenHash: hash,
		ExpiresAt: time.Now().Add(ttl),
	}
	if err := db.Create(&rec).Error; err != nil {
		return RefreshPair{}, err
	}
	return RefreshPair{Plain: plain, Record: rec}, nil
}

// ValidateRefreshToken refresh token; returns record if valid & not revoked/expired
func ValidateRefreshToken(db *gorm.DB, plain string) (*models.RefreshToken, error) {
	if plain == "" {
		return nil, errors.New("empty")
	}
	// var rec models.RefreshToken
	// We can't query by hash w/ Argon; scan candidates by expiry window. Keep small TTL (e.g. 30d).
	if err := db.Where("expires_at > ? AND revoked_at IS NULL", time.Now()).
		Find(&[]models.RefreshToken{}).Error; err != nil {
		return nil, err
	}
	// Better: add a prefix column to narrow scan; omitted for brevity.

	// Pragmatic approach: single SELECT per token:
	// Add a TokenHashSHA256 column for deterministic lookup if you want O(1). (Optional)

	// Minimal: iterate limited set; for simplicity we fetch by created window:
	var recs []models.RefreshToken
	if err := db.Where("expires_at > ? AND revoked_at IS NULL", time.Now()).
		Order("created_at desc").Limit(500).Find(&recs).Error; err != nil {
		return nil, err
	}
	for _, r := range recs {
		ok, _ := VerifySecretArgon2id(r.TokenHash, plain)
		if ok {
			return &r, nil
		}
	}
	return nil, errors.New("invalid")
}

// RevokeFamily revokes all tokens in a family (logout everywhere)
func RevokeFamily(db *gorm.DB, familyID uuid.UUID) error {
	now := time.Now()
	return db.Model(&models.RefreshToken{}).
		Where("family_id = ? AND revoked_at IS NULL", familyID).
		Update("revoked_at", &now).Error
}

// RotateRefreshToken replaces one token with a fresh one within the same family
func RotateRefreshToken(db *gorm.DB, used *models.RefreshToken, ttl time.Duration) (RefreshPair, error) {
	// revoke the used token (one-time use)
	now := time.Now()
	if err := db.Model(&models.RefreshToken{}).
		Where("id = ? AND revoked_at IS NULL", used.ID).
		Update("revoked_at", &now).Error; err != nil {
		return RefreshPair{}, err
	}
	return IssueRefreshToken(db, used.UserID, ttl, &used.FamilyID)
}
