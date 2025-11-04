package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	"github.com/glueops/autoglue/internal/api/httpmiddleware"
	"github.com/glueops/autoglue/internal/auth"
	"github.com/glueops/autoglue/internal/models"
	"github.com/glueops/autoglue/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userAPIKeyOut struct {
	ID         uuid.UUID  `json:"id" format:"uuid"`
	Name       *string    `json:"name,omitempty"`
	Scope      string     `json:"scope"` // "user"
	CreatedAt  time.Time  `json:"created_at"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
	Plain      *string    `json:"plain,omitempty"` // Shown only on create:
}

// ListUserAPIKeys godoc
// @ID ListUserAPIKeys
// @Summary List my API keys
// @Tags MeAPIKeys
// @Produce json
// @Success 200 {array} userAPIKeyOut
// @Router /me/api-keys [get]
// @Security BearerAuth
// @Security ApiKeyAuth
func ListUserAPIKeys(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, ok := httpmiddleware.UserFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusUnauthorized, "unauthorized", "not signed in")
			return
		}
		var rows []models.APIKey
		if err := db.
			Where("scope = ? AND user_id = ?", "user", u.ID).
			Order("created_at desc").
			Find(&rows).Error; err != nil {
			utils.WriteError(w, 500, "db_error", err.Error())
			return
		}
		out := make([]userAPIKeyOut, 0, len(rows))
		for _, k := range rows {
			out = append(out, toUserKeyOut(k, nil))
		}
		utils.WriteJSON(w, 200, out)
	}
}

type createUserKeyRequest struct {
	Name           string `json:"name,omitempty"`
	ExpiresInHours *int   `json:"expires_in_hours,omitempty"` // optional TTL
}

// CreateUserAPIKey godoc
// @ID CreateUserAPIKey
// @Summary Create a new user API key
// @Description Returns the plaintext key once. Store it securely on the client side.
// @Tags MeAPIKeys
// @Accept json
// @Produce json
// @Param body body createUserKeyRequest true "Key options"
// @Success 201 {object} userAPIKeyOut
// @Router /me/api-keys [post]
// @Security BearerAuth
// @Security ApiKeyAuth
func CreateUserAPIKey(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, ok := httpmiddleware.UserFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusUnauthorized, "unauthorized", "not signed in")
			return
		}
		var req createUserKeyRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteError(w, 400, "invalid_json", err.Error())
			return
		}

		plain, err := generateUserAPIKey()
		if err != nil {
			utils.WriteError(w, 500, "gen_failed", err.Error())
			return
		}
		hash := auth.SHA256Hex(plain)

		var exp *time.Time
		if req.ExpiresInHours != nil && *req.ExpiresInHours > 0 {
			t := time.Now().Add(time.Duration(*req.ExpiresInHours) * time.Hour)
			exp = &t
		}

		rec := models.APIKey{
			Scope:     "user",
			UserID:    &u.ID,
			KeyHash:   hash,
			Name:      req.Name, // if field exists
			ExpiresAt: exp,
			// SecretHash: nil (not used for user keys)
		}
		if err := db.Create(&rec).Error; err != nil {
			utils.WriteError(w, 500, "db_error", err.Error())
			return
		}
		utils.WriteJSON(w, http.StatusCreated, toUserKeyOut(rec, &plain))
	}
}

// DeleteUserAPIKey godoc
// @ID DeleteUserAPIKey
// @Summary Delete a user API key
// @Tags MeAPIKeys
// @Produce json
// @Param id path string true "Key ID (UUID)"
// @Success 204 "No Content"
// @Router /me/api-keys/{id} [delete]
// @Security BearerAuth
func DeleteUserAPIKey(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, ok := httpmiddleware.UserFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusUnauthorized, "unauthorized", "not signed in")
			return
		}
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			utils.WriteError(w, 400, "invalid_id", "must be uuid")
			return
		}
		tx := db.Where("id = ? AND scope = ? AND user_id = ?", id, "user", u.ID).
			Delete(&models.APIKey{})
		if tx.Error != nil {
			utils.WriteError(w, 500, "db_error", tx.Error.Error())
			return
		}
		if tx.RowsAffected == 0 {
			utils.WriteError(w, 404, "not_found", "key not found")
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func toUserKeyOut(k models.APIKey, plain *string) userAPIKeyOut {
	return userAPIKeyOut{
		ID:         k.ID,
		Name:       &k.Name, // if your model has it; else remove
		Scope:      k.Scope,
		CreatedAt:  k.CreatedAt,
		ExpiresAt:  k.ExpiresAt,
		LastUsedAt: k.LastUsedAt, // if present; else remove
		Plain:      plain,
	}
}

func generateUserAPIKey() (string, error) {
	// 24 random bytes → base64url (no padding), with "u_" prefix
	b := make([]byte, 24)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	s := base64.RawURLEncoding.EncodeToString(b)
	return "u_" + s, nil
}
