package httpmiddleware

import (
	"net/http"

	"github.com/glueops/autoglue/internal/utils"
)

// RequireAuthenticatedUser ensures a user principal is present (i.e. not an org/machine key).
func RequireAuthenticatedUser() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if user, ok := UserFrom(r.Context()); !ok || user == nil {
				// No user in context -> probably org/machine principal, or unauthenticated
				utils.WriteError(w, http.StatusUnauthorized, "unauthorized", "user principal required")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// RequirePlatformAdmin requires a user principal with IsAdmin=true.
// This is platform-wide (non-org) admin and does NOT depend on org roles.
func RequirePlatformAdmin() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := UserFrom(r.Context())
			if !ok || user == nil {
				utils.WriteError(w, http.StatusUnauthorized, "unauthorized", "user principal required")
				return
			}
			if !user.IsAdmin {
				utils.WriteError(w, http.StatusForbidden, "forbidden", "platform admin required")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// RequireUserAdmin is an alias for RequirePlatformAdmin for readability at call sites.
func RequireUserAdmin() func(http.Handler) http.Handler {
	return RequirePlatformAdmin()
}
