package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/glueops/autoglue/internal/db"
	"github.com/glueops/autoglue/internal/db/models"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthClaims struct {
	Orgs  []string `json:"orgs"`
	Roles []string `json:"roles"`
	jwt.RegisteredClaims
}

type AuthContext struct {
	UserID         uuid.UUID
	OrganizationID uuid.UUID
	OrgRole        string      // Role in the org
	Claims         *AuthClaims `json:"claims,omitempty" swaggerignore:"true"`
}

type contextKey struct{}

var authContextKey = contextKey{}

func AuthMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
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

			orgIDStr := r.Header.Get("X-Org-ID")
			if orgIDStr == "" {
				if rc := chi.RouteContext(r.Context()); rc != nil {
					if v := rc.URLParam("orgId"); v != "" {
						orgIDStr = v
					} else if v := rc.URLParam("organizationId"); v != "" {
						orgIDStr = v
					}
				}
			}

			if orgIDStr != "" {
				orgUUID, err := uuid.Parse(orgIDStr)
				if err != nil {
					http.Error(w, "invalid organization id", http.StatusBadRequest)
					return
				}

				var member models.Member
				if err := db.DB.
					Where("user_id = ? AND organization_id = ?", userUUID, orgUUID).
					First(&member).Error; err != nil {
					http.Error(w, "forbidden", http.StatusForbidden)
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

func GetAuthContext(r *http.Request) *AuthContext {
	if ac, ok := r.Context().Value(authContextKey).(*AuthContext); ok {
		return ac
	}
	return nil
}
