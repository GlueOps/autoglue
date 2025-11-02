package auth

import (
	"github.com/glueops/autoglue/internal/config"
	"github.com/glueops/autoglue/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ValidateJWT verifies RS256/RS384/RS512/EdDSA tokens using the in-memory key cache.
// It honors kid when present, and falls back to any active key.
func ValidateJWT(tokenStr string, db *gorm.DB) *models.User {
	cfg, _ := config.Load()

	parser := jwt.NewParser(
		jwt.WithIssuer(cfg.JWTIssuer),
		jwt.WithAudience(cfg.JWTAudience),
		jwt.WithValidMethods([]string{"RS256", "RS384", "RS512", "EdDSA"}),
	)

	token, err := parser.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		// Resolve by kid first
		kid, _ := t.Header["kid"].(string)

		kc.mu.RLock()
		defer kc.mu.RUnlock()

		if kid != "" {
			if k, ok := kc.pub[kid]; ok {
				return k, nil
			}
		}
		// Fallback: try first active key
		for _, k := range kc.pub {
			return k, nil
		}
		return nil, jwt.ErrTokenUnverifiable
	})
	if err != nil || !token.Valid {
		return nil
	}

	claims, _ := token.Claims.(jwt.MapClaims)
	sub, _ := claims["sub"].(string)
	uid, err := uuid.Parse(sub)
	if err != nil {
		return nil
	}

	var u models.User
	if err := db.First(&u, "id = ? AND is_disabled = false", uid).Error; err != nil {
		return nil
	}
	return &u
}
