package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/glueops/autoglue/internal/auth"
	"github.com/glueops/autoglue/internal/config"
	"github.com/glueops/autoglue/internal/handlers/dto"
	"github.com/glueops/autoglue/internal/models"
	"github.com/glueops/autoglue/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

type oauthProvider struct {
	Name     string
	Issuer   string
	Scopes   []string
	ClientID string
	Secret   string
}

func providerConfig(cfg config.Config, name string) (oauthProvider, bool) {
	switch strings.ToLower(name) {
	case "google":
		return oauthProvider{
			Name:     "google",
			Issuer:   "https://accounts.google.com",
			Scopes:   []string{oidc.ScopeOpenID, "email", "profile"},
			ClientID: cfg.GoogleClientID,
			Secret:   cfg.GoogleClientSecret,
		}, true
	case "github":
		// GitHub is not a pure OIDC provider; we use OAuth2 + user email API
		return oauthProvider{
			Name:     "github",
			Issuer:   "github",
			Scopes:   []string{"read:user", "user:email"},
			ClientID: cfg.GithubClientID, Secret: cfg.GithubClientSecret,
		}, true
	}
	return oauthProvider{}, false
}

// AuthStart godoc
//
//	@ID				AuthStart
//	@Summary		Begin social login
//	@Description	Returns provider authorization URL for the frontend to redirect
//	@Tags			Auth
//	@Param			provider	path	string	true	"google|github"
//	@Produce		json
//	@Success		200	{object}	dto.AuthStartResponse
//	@Router			/auth/{provider}/start [post]
func AuthStart(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cfg, _ := config.Load()
		provider := strings.ToLower(chi.URLParam(r, "provider"))

		p, ok := providerConfig(cfg, provider)
		if !ok || p.ClientID == "" || p.Secret == "" {
			utils.WriteError(w, http.StatusBadRequest, "unsupported_provider", "provider not configured")
			return
		}

		redirect := cfg.OAuthRedirectBase + "/api/v1/auth/" + p.Name + "/callback"

		// Optional SPA hints to be embedded into state
		mode := r.URL.Query().Get("mode")     // "spa" enables postMessage callback page
		origin := r.URL.Query().Get("origin") // e.g. http://localhost:5173

		state := uuid.NewString()
		if mode == "spa" && origin != "" {
			state = state + "|mode=spa|origin=" + url.QueryEscape(origin)
		}

		var authURL string

		if p.Issuer == "github" {
			o := &oauth2.Config{
				ClientID:     p.ClientID,
				ClientSecret: p.Secret,
				RedirectURL:  redirect,
				Scopes:       p.Scopes,
				Endpoint: oauth2.Endpoint{
					AuthURL:  "https://github.com/login/oauth/authorize",
					TokenURL: "https://github.com/login/oauth/access_token",
				},
			}
			authURL = o.AuthCodeURL(state, oauth2.AccessTypeOffline)
		} else {
			// Google OIDC
			ctx := context.Background()
			prov, err := oidc.NewProvider(ctx, p.Issuer)
			if err != nil {
				utils.WriteError(w, http.StatusInternalServerError, "oidc_discovery_failed", err.Error())
				return
			}
			o := &oauth2.Config{
				ClientID:     p.ClientID,
				ClientSecret: p.Secret,
				RedirectURL:  redirect,
				Endpoint:     prov.Endpoint(),
				Scopes:       p.Scopes,
			}
			authURL = o.AuthCodeURL(state, oauth2.AccessTypeOffline)
		}

		utils.WriteJSON(w, http.StatusOK, dto.AuthStartResponse{AuthURL: authURL})
	}
}

// AuthCallback godoc
//
//	@ID			AuthCallback
//	@Summary	Handle social login callback
//	@Tags		Auth
//	@Param		provider	path	string	true	"google|github"
//	@Produce	json
//	@Success	200	{object}	dto.TokenPair
//	@Router		/auth/{provider}/callback [get]
func AuthCallback(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cfg, _ := config.Load()
		provider := strings.ToLower(chi.URLParam(r, "provider"))

		p, ok := providerConfig(cfg, provider)
		if !ok {
			utils.WriteError(w, http.StatusBadRequest, "unsupported_provider", "provider not configured")
			return
		}

		code := r.URL.Query().Get("code")
		if code == "" {
			utils.WriteError(w, http.StatusBadRequest, "invalid_request", "missing code")
			return
		}
		redirect := cfg.OAuthRedirectBase + "/api/v1/auth/" + p.Name + "/callback"

		var email, display, subject string

		if p.Issuer == "github" {
			// OAuth2 code exchange
			o := &oauth2.Config{
				ClientID:     p.ClientID,
				ClientSecret: p.Secret,
				RedirectURL:  redirect,
				Scopes:       p.Scopes,
				Endpoint: oauth2.Endpoint{
					AuthURL:  "https://github.com/login/oauth/authorize",
					TokenURL: "https://github.com/login/oauth/access_token",
				},
			}
			tok, err := o.Exchange(r.Context(), code)
			if err != nil {
				utils.WriteError(w, http.StatusUnauthorized, "exchange_failed", err.Error())
				return
			}
			// Fetch user primary email
			req, _ := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
			req.Header.Set("Authorization", "token "+tok.AccessToken)
			resp, err := http.DefaultClient.Do(req)
			if err != nil || resp.StatusCode != 200 {
				utils.WriteError(w, http.StatusUnauthorized, "email_fetch_failed", "github user/emails")
				return
			}
			defer resp.Body.Close()
			var emails []struct {
				Email    string `json:"email"`
				Primary  bool   `json:"primary"`
				Verified bool   `json:"verified"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil || len(emails) == 0 {
				utils.WriteError(w, http.StatusUnauthorized, "email_parse_failed", err.Error())
				return
			}
			email = emails[0].Email
			for _, e := range emails {
				if e.Primary {
					email = e.Email
					break
				}
			}
			subject = "github:" + email
			display = strings.Split(email, "@")[0]
		} else {
			// Google OIDC
			oidcProv, err := oidc.NewProvider(r.Context(), p.Issuer)
			if err != nil {
				utils.WriteError(w, 500, "oidc_discovery_failed", err.Error())
				return
			}
			o := &oauth2.Config{
				ClientID:     p.ClientID,
				ClientSecret: p.Secret,
				RedirectURL:  redirect,
				Endpoint:     oidcProv.Endpoint(),
				Scopes:       p.Scopes,
			}
			tok, err := o.Exchange(r.Context(), code)
			if err != nil {
				utils.WriteError(w, 401, "exchange_failed", err.Error())
				return
			}

			verifier := oidcProv.Verifier(&oidc.Config{ClientID: p.ClientID})
			rawIDToken, ok := tok.Extra("id_token").(string)
			if !ok {
				utils.WriteError(w, 401, "no_id_token", "")
				return
			}
			idt, err := verifier.Verify(r.Context(), rawIDToken)
			if err != nil {
				utils.WriteError(w, 401, "id_token_invalid", err.Error())
				return
			}

			var claims struct {
				Email         string `json:"email"`
				EmailVerified bool   `json:"email_verified"`
				Name          string `json:"name"`
				Sub           string `json:"sub"`
			}
			if err := idt.Claims(&claims); err != nil {
				utils.WriteError(w, 401, "claims_parse_error", err.Error())
				return
			}
			email = strings.ToLower(claims.Email)
			display = claims.Name
			subject = "google:" + claims.Sub
		}

		// Upsert Account + User; domain auto-join (member)
		user, err := upsertAccountAndUser(db, p.Name, subject, email, display)
		if err != nil {
			utils.WriteError(w, 500, "account_upsert_failed", err.Error())
			return
		}

		// Org auto-join: Organization.Domain == email domain
		_ = ensureAutoMembership(db, user.ID, email)

		// Issue tokens
		accessTTL := 1 * time.Hour
		refreshTTL := 30 * 24 * time.Hour

		access, err := auth.IssueAccessToken(auth.IssueOpts{
			Subject:  user.ID.String(),
			Issuer:   cfg.JWTIssuer,
			Audience: cfg.JWTAudience,
			TTL:      accessTTL,
			Claims: map[string]any{
				"email": email,
				"name":  display,
			},
		})
		if err != nil {
			utils.WriteError(w, 500, "issue_access_failed", err.Error())
			return
		}

		rp, err := auth.IssueRefreshToken(db, user.ID, refreshTTL, nil)
		if err != nil {
			utils.WriteError(w, 500, "issue_refresh_failed", err.Error())
			return
		}

		// If the state indicates SPA popup mode, postMessage tokens to the opener and close
		state := r.URL.Query().Get("state")
		if strings.Contains(state, "mode=spa") {
			origin := ""
			for _, part := range strings.Split(state, "|") {
				if strings.HasPrefix(part, "origin=") {
					origin, _ = url.QueryUnescape(strings.TrimPrefix(part, "origin="))
					break
				}
			}
			// fallback: restrict to backend origin if none supplied
			if origin == "" {
				origin = cfg.OAuthRedirectBase
			}
			payload := dto.TokenPair{
				AccessToken:  access,
				RefreshToken: rp.Plain,
				TokenType:    "Bearer",
				ExpiresIn:    int64(accessTTL.Seconds()),
			}
			writePostMessageHTML(w, origin, payload)
			return
		}

		// Default JSON response
		utils.WriteJSON(w, http.StatusOK, dto.TokenPair{
			AccessToken:  access,
			RefreshToken: rp.Plain,
			TokenType:    "Bearer",
			ExpiresIn:    int64(accessTTL.Seconds()),
		})
	}
}

// Refresh godoc
//
//	@ID			Refresh
//	@Summary	Rotate refresh token
//	@Tags		Auth
//	@Accept		json
//	@Produce	json
//	@Param		body	body		dto.RefreshRequest	true	"Refresh token"
//	@Success	200		{object}	dto.TokenPair
//	@Router		/auth/refresh [post]
func Refresh(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cfg, _ := config.Load()
		var req dto.RefreshRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteError(w, 400, "invalid_json", err.Error())
			return
		}
		rec, err := auth.ValidateRefreshToken(db, req.RefreshToken)
		if err != nil {
			utils.WriteError(w, 401, "invalid_refresh", "")
			return
		}

		var u models.User
		if err := db.First(&u, "id = ? AND is_disabled = false", rec.UserID).Error; err != nil {
			utils.WriteError(w, 401, "user_disabled", "")
			return
		}

		// rotate
		newPair, err := auth.RotateRefreshToken(db, rec, 30*24*time.Hour)
		if err != nil {
			utils.WriteError(w, 500, "rotate_failed", err.Error())
			return
		}

		// new access
		access, err := auth.IssueAccessToken(auth.IssueOpts{
			Subject:  u.ID.String(),
			Issuer:   cfg.JWTIssuer,
			Audience: cfg.JWTAudience,
			TTL:      1 * time.Hour,
		})
		if err != nil {
			utils.WriteError(w, 500, "issue_access_failed", err.Error())
			return
		}

		utils.WriteJSON(w, 200, dto.TokenPair{
			AccessToken:  access,
			RefreshToken: newPair.Plain,
			TokenType:    "Bearer",
			ExpiresIn:    3600,
		})
	}
}

// Logout godoc
//
//	@ID			Logout
//	@Summary	Revoke refresh token family (logout everywhere)
//	@Tags		Auth
//	@Accept		json
//	@Produce	json
//	@Param		body	body	dto.LogoutRequest	true	"Refresh token"
//	@Success	204		"No Content"
//	@Router		/auth/logout [post]
func Logout(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.LogoutRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteError(w, 400, "invalid_json", err.Error())
			return
		}
		rec, err := auth.ValidateRefreshToken(db, req.RefreshToken)
		if err != nil {
			w.WriteHeader(204) // already invalid/revoked
			return
		}
		if err := auth.RevokeFamily(db, rec.FamilyID); err != nil {
			utils.WriteError(w, 500, "revoke_failed", err.Error())
			return
		}
		w.WriteHeader(204)
	}
}

// Helpers

func upsertAccountAndUser(db *gorm.DB, provider, subject, email, display string) (*models.User, error) {
	email = strings.ToLower(email)
	var acc models.Account
	if err := db.Where("provider = ? AND subject = ?", provider, subject).First(&acc).Error; err == nil {
		var u models.User
		if err := db.First(&u, "id = ?", acc.UserID).Error; err != nil {
			return nil, err
		}
		return &u, nil
	}
	// Link by email if exists
	var ue models.UserEmail
	if err := db.Where("LOWER(email) = ?", email).First(&ue).Error; err == nil {
		acc = models.Account{
			UserID:        ue.UserID,
			Provider:      provider,
			Subject:       subject,
			Email:         &email,
			EmailVerified: true,
		}
		if err := db.Create(&acc).Error; err != nil {
			return nil, err
		}
		var u models.User
		if err := db.First(&u, "id = ?", ue.UserID).Error; err != nil {
			return nil, err
		}
		return &u, nil
	}
	// Create user
	u := models.User{DisplayName: &display, PrimaryEmail: &email}
	if err := db.Create(&u).Error; err != nil {
		return nil, err
	}
	ue = models.UserEmail{UserID: u.ID, Email: email, IsVerified: true, IsPrimary: true}
	_ = db.Create(&ue).Error
	acc = models.Account{UserID: u.ID, Provider: provider, Subject: subject, Email: &email, EmailVerified: true}
	_ = db.Create(&acc).Error
	return &u, nil
}

func ensureAutoMembership(db *gorm.DB, userID uuid.UUID, email string) error {
	parts := strings.SplitN(strings.ToLower(email), "@", 2)
	if len(parts) != 2 {
		return nil
	}
	domain := parts[1]
	var org models.Organization
	if err := db.Where("LOWER(domain) = ?", domain).First(&org).Error; err != nil {
		return nil
	}
	// if already member, done
	var c int64
	db.Model(&models.Membership{}).
		Where("user_id = ? AND organization_id = ?", userID, org.ID).
		Count(&c)
	if c > 0 {
		return nil
	}
	return db.Create(&models.Membership{
		UserID: userID, OrganizationID: org.ID, Role: "member",
	}).Error
}

// writePostMessageHTML sends a tiny HTML page that posts tokens to the SPA and closes the window.
func writePostMessageHTML(w http.ResponseWriter, origin string, payload dto.TokenPair) {
	b, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`<!doctype html><html><body><script>
(function(){
  try {
    var data = ` + string(b) + `;
    if (window.opener) {
      window.opener.postMessage({ type: 'autoglue:auth', payload: data }, '` + origin + `');
    }
  } catch (e) {}
  window.close();
})();
</script></body></html>`))
}
