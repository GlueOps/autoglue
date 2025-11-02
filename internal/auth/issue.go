package auth

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/glueops/autoglue/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func randomToken(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	// URL-safe, no padding
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// IssueUserAPIKey creates a single-token user API key (X-API-KEY)
func IssueUserAPIKey(db *gorm.DB, userID uuid.UUID, name string, ttl *time.Duration) (plaintext string, rec models.APIKey, err error) {
	plaintext, err = randomToken(32)
	if err != nil {
		return "", models.APIKey{}, err
	}
	rec = models.APIKey{
		Name:    name,
		Scope:   "user",
		UserID:  &userID,
		KeyHash: SHA256Hex(plaintext), // deterministic lookup
	}
	if ttl != nil {
		ex := time.Now().Add(*ttl)
		rec.ExpiresAt = &ex
	}
	if err = db.Create(&rec).Error; err != nil {
		return "", models.APIKey{}, err
	}
	return plaintext, rec, nil
}
