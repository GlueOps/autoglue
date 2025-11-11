package httpmiddleware

import (
	"net/http"
	"strings"

	"github.com/glueops/autoglue/internal/auth"
	"github.com/glueops/autoglue/internal/models"
	"github.com/glueops/autoglue/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AuthMiddleware authenticates either a user principal (JWT, user API key, app key/secret)
// or an org principal (org key/secret). If requireOrg is true, the request must have
// an organization resolved; otherwise org is optional.
//
// Org resolution order for user principals (when requireOrg == true):
//  1. X-Org-ID header (UUID)
//  2. chi URL param {id} (useful under /orgs/{id}/... routers)
//  3. single-membership fallback (exactly one membership)
//
// If none resolves, respond with org_required.
func AuthMiddleware(db *gorm.DB, requireOrg bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var user *models.User
			var org *models.Organization
			var roles []string

			// --- 1) Authenticate principal ---
			// Prefer org principal if explicit machine access is provided.
			if orgKey := r.Header.Get("X-ORG-KEY"); orgKey != "" {
				secret := r.Header.Get("X-ORG-SECRET")
				org = auth.ValidateOrgKeyPair(orgKey, secret, db)
				if org == nil {
					utils.WriteError(w, http.StatusUnauthorized, "unauthorized", "invalid org credentials")
					return
				}
				// org principal implies machine role
				roles = []string{"org:machine"}
			} else {
				// User principals
				if ah := r.Header.Get("Authorization"); strings.HasPrefix(ah, "Bearer ") {
					user = auth.ValidateJWT(ah[7:], db)
				} else if apiKey := r.Header.Get("X-API-KEY"); apiKey != "" {
					user = auth.ValidateAPIKey(apiKey, db)
				} else if appKey := r.Header.Get("X-APP-KEY"); appKey != "" {
					secret := r.Header.Get("X-APP-SECRET")
					user = auth.ValidateAppKeyPair(appKey, secret, db)
				} else if c, err := r.Cookie("ag_jwt"); err == nil {
					tok := strings.TrimSpace(c.Value)
					if strings.HasPrefix(strings.ToLower(tok), "bearer ") {
						tok = tok[7:]
					}
					if tok != "" {
						user = auth.ValidateJWT(tok, db)
					}
				}

				if user == nil {
					utils.WriteError(w, http.StatusUnauthorized, "unauthorized", "invalid credentials")
					return
				}

				// --- 2) Resolve organization (user principal) ---
				// A) Try X-Org-ID if present
				if s := r.Header.Get("X-Org-ID"); s != "" {
					oid, err := uuid.Parse(s)
					if err != nil {
						utils.WriteError(w, http.StatusBadRequest, "invalid_org_id", "X-Org-ID must be a UUID")
						return
					}
					var o models.Organization
					if err := db.First(&o, "id = ?", oid).Error; err != nil {
						// Header provided but org not found
						utils.WriteError(w, http.StatusUnauthorized, "org_forbidden", "organization not found")
						return
					}
					// Verify membership
					if !userIsMember(db, user.ID, o.ID) {
						utils.WriteError(w, http.StatusUnauthorized, "org_forbidden", "user is not a member of specified org")
						return
					}
					org = &o
				}

				// B) If still no org and requireOrg==true, try chi URL param {id}
				if org == nil && requireOrg {
					if sid := chi.URLParam(r, "id"); sid != "" {
						if oid, err := uuid.Parse(sid); err == nil {
							var o models.Organization
							if err := db.First(&o, "id = ?", oid).Error; err == nil && userIsMember(db, user.ID, o.ID) {
								org = &o
							} else {
								utils.WriteError(w, http.StatusUnauthorized, "org_forbidden", "user is not a member of specified org")
								return
							}
						}
					}
				}

				// C) Single-membership fallback (only if requireOrg==true and still nil)
				if org == nil && requireOrg {
					var ms []models.Membership
					if err := db.Where("user_id = ?", user.ID).Find(&ms).Error; err == nil && len(ms) == 1 {
						var o models.Organization
						if err := db.First(&o, "id = ?", ms[0].OrganizationID).Error; err == nil {
							org = &o
						}
					}
				}

				// D) Final check
				if requireOrg && org == nil {
					utils.WriteError(w, http.StatusUnauthorized, "org_required", "specify X-Org-ID or use an endpoint that does not require org")
					return
				}

				// Populate roles if an org was resolved (optional for org-optional endpoints)
				if org != nil {
					roles = userRolesInOrg(db, user.ID, org.ID)
					if len(roles) == 0 {
						utils.WriteError(w, http.StatusForbidden, "forbidden", "no roles in organization")
						return
					}
				}
			}

			// --- 3) Attach to context and proceed ---
			ctx := r.Context()
			if user != nil {
				ctx = WithUser(ctx, user)
			}
			if org != nil {
				ctx = WithOrg(ctx, org)
			}
			if roles != nil {
				ctx = WithRoles(ctx, roles)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func userIsMember(db *gorm.DB, userID, orgID uuid.UUID) bool {
	var count int64
	db.Model(&models.Membership{}).
		Where("user_id = ? AND organization_id = ?", userID, orgID).
		Count(&count)
	return count > 0
}

func userRolesInOrg(db *gorm.DB, userID, orgID uuid.UUID) []string {
	var m models.Membership
	if err := db.Where("user_id = ? AND organization_id = ?", userID, orgID).First(&m).Error; err == nil {
		switch m.Role {
		case "owner":
			return []string{"role:owner", "role:admin", "role:member"}
		case "admin":
			return []string{"role:admin", "role:member"}
		default:
			return []string{"role:member"}
		}
	}
	return nil
}
