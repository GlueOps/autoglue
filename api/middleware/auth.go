package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/glueops/autoglue/internal/db"
	"github.com/glueops/autoglue/internal/db/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthClaims struct {
	Orgs  []string `json:"orgs"`
	Roles []string `json:"roles"` // optional global roles
	jwt.RegisteredClaims
}

type AuthContext struct {
	UserID         uuid.UUID
	OrganizationID uuid.UUID
	OrgRole        string      // Role in the org
	Claims         *AuthClaims `json:"claims,omitempty" swaggerignore:"true"`
}

type contextKey string

const authContextKey contextKey = "authentication"

func AuthMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			claims := &AuthClaims{}

			token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			userUUID, err := uuid.Parse(claims.Subject)
			if err != nil {
				http.Error(w, "invalid user id", http.StatusUnauthorized)
				return
			}

			authCtx := &AuthContext{
				UserID: userUUID,
				Claims: claims,
			}

			orgID := r.Header.Get("X-Org-ID")
			orgUUID, _ := uuid.Parse(orgID)

			if orgID != "" {
				var member models.Member
				if err := db.DB.Where("user_id = ? AND organization_id = ?", claims.Subject, orgID).First(&member).Error; err != nil {
					http.Error(w, "User not a member of the organization", http.StatusForbidden)
					return
				}
				authCtx.OrganizationID = orgUUID
				authCtx.OrgRole = string(member.Role)
			}

			ctx := context.WithValue(r.Context(), authContextKey, authCtx)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetAuthContext retrieves the full authenticated context from the request
func GetAuthContext(r *http.Request) *AuthContext {
	val := r.Context().Value(authContextKey)
	if ac, ok := val.(*AuthContext); ok {
		return ac
	}
	return nil
}
