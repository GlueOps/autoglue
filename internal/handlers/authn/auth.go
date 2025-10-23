package authn

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/glueops/autoglue/internal/config"
	"github.com/glueops/autoglue/internal/db"
	"github.com/glueops/autoglue/internal/db/models"
	"github.com/glueops/autoglue/internal/middleware"
	"github.com/glueops/autoglue/internal/response"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Register godoc
// @ID           Register
// @Summary      Register a new user
// @Description  Registers a new user and stores credentials
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      RegisterInput  true  "User registration input"
// @Success      201   {string}  string         "created"
// @Failure      400   {string}  string         "bad request"
// @Router       /api/v1/auth/register [post]
func Register(w http.ResponseWriter, r *http.Request) {
	var input RegisterInput
	json.NewDecoder(r.Body).Decode(&input)

	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "failed to hash password", http.StatusInternalServerError)
		return
	}

	user := models.User{Email: input.Email, Password: string(hashed), Name: input.Name, Role: "user"}
	if err := db.DB.Create(&user).Error; err != nil {
		http.Error(w, "registration failed", 400)
		return
	}

	_ = response.JSON(w, http.StatusCreated, map[string]string{"status": "created"})
}

// Login godoc
// @ID           Login
// @Summary      Authenticate and return a token
// @Description  Authenticates a user and returns a JWT bearer token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      LoginInput     true  "User login input"
// @Success      200   {object}  map[string]string "token"
// @Failure      401   {string}  string         "unauthorized"
// @Router       /api/v1/auth/login [post]
func Login(w http.ResponseWriter, r *http.Request) {
	var input LoginInput
	json.NewDecoder(r.Body).Decode(&input)

	var user models.User
	if err := db.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessStr, _ := accessToken.SignedString(jwtSecret)

	refreshTokenStr := uuid.NewString()

	_ = db.DB.Create(&models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshTokenStr,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		Revoked:   false,
	}).Error

	_ = response.JSON(w, http.StatusOK, map[string]string{
		"access_token":  accessStr,
		"refresh_token": refreshTokenStr,
	})
}

// Refresh godoc
// @ID           Refresh
// @Summary      Refresh access token
// @Description  Use a refresh token to obtain a new access token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      map[string]string  true  "refresh_token"
// @Success      200   {object}  map[string]string "new access token"
// @Failure      401   {string}  string         "unauthorized"
// @Security     BearerAuth
// @Router       /api/v1/auth/refresh [post]
func Refresh(w http.ResponseWriter, r *http.Request) {
	var input struct {
		RefreshToken string `json:"refresh_token"`
	}
	json.NewDecoder(r.Body).Decode(&input)

	var token models.RefreshToken
	if err := db.DB.Where("token = ? AND revoked = false", input.RefreshToken).First(&token).Error; err != nil || token.ExpiresAt.Before(time.Now()) {
		http.Error(w, "invalid or expired refresh token", http.StatusUnauthorized)
		return
	}

	claims := jwt.MapClaims{
		"sub": token.UserID,
		"exp": time.Now().Add(rotatedAccessTTL).Unix(),
	}
	newAccess := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	newToken, _ := newAccess.SignedString(jwtSecret)

	_ = response.JSON(w, http.StatusOK, map[string]string{
		"access_token": newToken,
	})
}

// Logout godoc
// @ID           Logout
// @Summary      Logout user
// @Description  Revoke a refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      map[string]string  true  "refresh_token"
// @Success      204   {string}  string         "no content"
// @Security     BearerAuth
// @Router       /api/v1/auth/logout [post]
func Logout(w http.ResponseWriter, r *http.Request) {
	var input struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil || input.RefreshToken == "" {
		response.Error(w, http.StatusBadRequest, "bad request")
		return
	}

	db.DB.Model(&models.RefreshToken{}).Where("token = ?", input.RefreshToken).Update("revoked", true)
	response.NoContent(w)
}

// Me godoc
// @ID           Me
// @Summary      Get authenticated user info
// @Description  Returns the authenticated user's profile and auth context
// @Tags         auth
// @Produce      json
// @Success      200  {object}  MeResponse
// @Failure      401  {string}  string  "unauthorized"
// @Security     BearerAuth
// @Router       /api/v1/auth/me [get]
func Me(w http.ResponseWriter, r *http.Request) {
	authCtx := middleware.GetAuthContext(r)
	if authCtx == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var user models.User
	if err := db.DB.First(&user, "id = ?", authCtx.UserID).Error; err != nil {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	out := MeResponse{
		User: UserDTO{
			ID:            user.ID,
			Name:          user.Name,
			Email:         user.Email,
			EmailVerified: user.EmailVerified,
			Role:          user.Role,
			CreatedAt:     user.CreatedAt, // from Timestamped
			UpdatedAt:     user.UpdatedAt, // from Timestamped
		},
		OrgRole: authCtx.OrgRole,
	}

	if authCtx.OrganizationID != uuid.Nil {
		s := authCtx.OrganizationID.String()
		out.OrganizationID = &s
	}

	if c := authCtx.Claims; c != nil {
		var exp, iat, nbf int64
		if c.ExpiresAt != nil {
			exp = c.ExpiresAt.Time.Unix()
		}
		if c.IssuedAt != nil {
			iat = c.IssuedAt.Time.Unix()
		}
		if c.NotBefore != nil {
			nbf = c.NotBefore.Time.Unix()
		}

		out.Claims = &AuthClaimsDTO{
			Orgs:      c.Orgs,
			Roles:     c.Roles,
			Issuer:    c.Issuer,
			Subject:   c.Subject,
			Audience:  []string(c.Audience),
			ExpiresAt: exp,
			IssuedAt:  iat,
			NotBefore: nbf,
		}
	}

	_ = response.JSON(w, http.StatusOK, out)
}

// Introspect godoc
// @ID           Introspect
// @Summary      Introspect a token
// @Description  Returns whether the token is active and basic metadata
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body   map[string]string true "token"
// @Success      200   {object}  map[string]any
// @Router       /api/v1/auth/introspect [post]
func Introspect(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil || in.Token == "" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	tok, err := jwt.Parse(in.Token, func(t *jwt.Token) (any, error) { return jwtSecret, nil })
	if err == nil && tok.Valid {
		claims, _ := tok.Claims.(jwt.MapClaims)
		_ = response.JSON(w, http.StatusOK, map[string]any{
			"active": true,
			"type":   "access",
			"sub":    claims["sub"],
			"exp":    claims["exp"],
			"iat":    claims["iat"],
			"nbf":    claims["nbf"],
		})
		return
	}

	var rt models.RefreshToken
	if err := db.DB.Where("token = ? AND revoked = false", in.Token).First(&rt).Error; err == nil && rt.ExpiresAt.After(time.Now()) {
		_ = response.JSON(w, http.StatusOK, map[string]any{
			"active": true,
			"type":   "refresh",
			"sub":    rt.UserID,
			"exp":    rt.ExpiresAt.Unix(),
		})
		return
	}

	_ = response.JSON(w, http.StatusOK, map[string]any{"active": false})
}

// RequestPasswordReset godoc
// @ID           RequestPasswordReset
// @Summary      Request password reset
// @Description  Sends a reset token to the user's email address
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body  map[string]string true "email"
// @Success      204   {string} string "no content"
// @Router       /api/v1/auth/password/forgot [post]
func RequestPasswordReset(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil || in.Email == "" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	// Always return 204 to avoid user enumeration.
	var user models.User
	if err := db.DB.Where("email = ?", in.Email).First(&user).Error; err == nil {
		_ = db.DB.Model(&models.PasswordReset{}).
			Where("user_id = ? AND used = false AND expires_at > ?", user.ID, time.Now()).
			Update("used", true).Error

		if tok, err := issuePasswordReset(user.ID, user.Email); err == nil {
			_ = sendEmail(user.Email, "Password reset", fmt.Sprintf("Your password reset token is: %s", tok))
			err := sendTemplatedEmail(user.Email, "password_reset.tmpl", PasswordResetData{
				Name:     user.Name,
				Email:    user.Email,
				Token:    tok,
				ResetURL: fmt.Sprintf("%s/auth/reset?token=%s", config.FrontendBaseURL(), tok), // e.g. fmt.Sprintf("%s/reset?token=%s", frontendURL, tok)
			})
			if err != nil {
				fmt.Printf("smtp send error: %v\n", err)
			}
		}
	}
	response.NoContent(w)
}

// ConfirmPasswordReset godoc
// @ID           ConfirmPasswordReset
// @Summary      Confirm password reset
// @Description  Resets the password using a valid reset token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body  map[string]string true "token, new_password"
// @Success      204   {string} string "no content"
// @Failure      400   {string} string "bad request"
// @Router       /api/v1/auth/password/reset [post]
func ConfirmPasswordReset(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil || in.Token == "" || in.NewPassword == "" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	var pr models.PasswordReset
	if err := db.DB.Where("token = ? AND used = false", in.Token).First(&pr).Error; err != nil || pr.ExpiresAt.Before(time.Now()) {
		response.Error(w, http.StatusBadRequest, "invalid or expired token")
		return
	}

	var user models.User
	if err := db.DB.First(&user, "id = ?", pr.UserID).Error; err != nil {
		response.Error(w, http.StatusBadRequest, "invalid token")
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(in.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to hash")
		return
	}
	if err := db.DB.Model(&user).Update("password", string(hashed)).Error; err != nil {
		response.Error(w, http.StatusInternalServerError, "update failed")
		return
	}

	_ = db.DB.Model(&models.PasswordReset{}).Where("id = ?", pr.ID).Update("used", true).Error
	_ = db.DB.Model(&models.RefreshToken{}).Where("user_id = ? AND revoked = false", user.ID).Update("revoked", true).Error

	response.NoContent(w)
}

// ChangePassword godoc
// @ID           ChangePassword
// @Summary      Change password
// @Description  Changes the password for the authenticated user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body  map[string]string true "current_password, new_password"
// @Success      204   {string} string "no content"
// @Failure      400   {string} string "bad request"
// @Security     BearerAuth
// @Router       /api/v1/auth/password/change [post]
func ChangePassword(w http.ResponseWriter, r *http.Request) {
	ctx := middleware.GetAuthContext(r)
	if ctx == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var in struct {
		Current string `json:"current_password"`
		New     string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil || in.Current == "" || in.New == "" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := db.DB.First(&user, "id = ?", ctx.UserID).Error; err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(in.Current)); err != nil {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(in.New), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "failed to hash", http.StatusInternalServerError)
		return
	}
	if err := db.DB.Model(&user).Update("password", string(hashed)).Error; err != nil {
		http.Error(w, "update failed", http.StatusInternalServerError)
		return
	}

	// Optional hardening: revoke all refresh tokens after password change
	// _ = db.DB.Model(&models.RefreshToken{}).Where("user_id = ?", user.ID).Update("revoked", true)

	w.WriteHeader(http.StatusNoContent)
}

// VerifyEmail godoc
// @ID           VerifyEmail
// @Summary      Verify email address
// @Description  Verifies the user's email using a token (often from an emailed link)
// @Tags         auth
// @Produce      json
// @Param        token  query  string  true  "verification token"
// @Success      204    {string} string "no content"
// @Failure      400    {string} string "bad request"
// @Router       /api/v1/auth/verify [get]
func VerifyEmail(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	var ev models.EmailVerification
	if err := db.DB.Where("token = ? AND used = false", token).First(&ev).Error; err != nil || ev.ExpiresAt.Before(time.Now()) {
		response.Error(w, http.StatusBadRequest, "invalid or expired token")
		return
	}

	_ = db.DB.Model(&models.User{}).Where("id = ?", ev.UserID).Updates(map[string]any{
		"email_verified":    true,
		"email_verified_at": time.Now(),
	}).Error
	_ = db.DB.Model(&models.EmailVerification{}).Where("id = ?", ev.ID).Update("used", true).Error

	response.NoContent(w)
}

// ResendVerification godoc
// @ID           ResendVerification
// @Summary      Resend email verification
// @Description  Sends a new email verification token if needed
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body  map[string]string true "email"
// @Success      204   {string} string "no content"
// @Router       /api/v1/auth/verify/resend [post]
func ResendVerification(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil || in.Email == "" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := db.DB.Where("email = ?", in.Email).First(&user).Error; err == nil {
		_ = db.DB.Model(&models.EmailVerification{}).
			Where("user_id = ? AND used = false AND expires_at > ?", user.ID, time.Now()).
			Update("used", true).Error

		if tok, err := issueEmailVerification(user.ID, user.Email); err == nil {
			_ = sendEmail(user.Email, "Verify your account", fmt.Sprintf("Your verification token is: %s", tok))
			_ = sendTemplatedEmail(user.Email, "verify_account.tmpl", VerifyEmailData{
				Name:            user.Name,
				Email:           user.Email,
				Token:           tok,
				VerificationURL: "", // e.g. fmt.Sprintf("%s/verify?token=%s", frontendURL, tok)
			})
		}
	}

	response.NoContent(w)
}

// LogoutAll godoc
// @ID           LogoutAll
// @Summary      Logout from all sessions
// @Description  Revokes all active refresh tokens for the authenticated user
// @Tags         auth
// @Produce      json
// @Success      204   {string} string "no content"
// @Security     BearerAuth
// @Router       /api/v1/auth/logout_all [post]
func LogoutAll(w http.ResponseWriter, r *http.Request) {
	ctx := middleware.GetAuthContext(r)
	if ctx == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	db.DB.Model(&models.RefreshToken{}).Where("user_id = ? AND revoked = false", ctx.UserID).Update("revoked", true)
	response.NoContent(w)
}

// RotateRefreshToken godoc
// @ID           RotateRefreshToken
// @Summary      Rotate refresh token
// @Description  Exchanges a valid refresh token for a new access and refresh token, revoking the old one
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body  map[string]string true "refresh_token"
// @Success      200   {object} map[string]string "access_token, refresh_token"
// @Failure      401   {string} string "unauthorized"
// @Security     BearerAuth
// @Router       /api/v1/auth/refresh/rotate [post]
func RotateRefreshToken(w http.ResponseWriter, r *http.Request) {
	var in struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil || in.RefreshToken == "" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	var old models.RefreshToken
	if err := db.DB.Where("token = ? AND revoked = false", in.RefreshToken).First(&old).Error; err != nil || old.ExpiresAt.Before(time.Now()) {
		http.Error(w, "invalid or expired refresh token", http.StatusUnauthorized)
		return
	}

	_ = db.DB.Model(&models.RefreshToken{}).Where("id = ?", old.ID).Update("revoked", true)

	claims := jwt.MapClaims{
		"sub": old.UserID,
		"exp": time.Now().Add(15 * time.Minute).Unix(),
	}
	access := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessStr, _ := access.SignedString(jwtSecret)

	newRefresh := models.RefreshToken{
		UserID:    old.UserID,
		Token:     uuid.NewString(),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		Revoked:   false,
	}
	_ = db.DB.Create(&newRefresh).Error

	_ = response.JSON(w, http.StatusOK, map[string]string{
		"access_token":  accessStr,
		"refresh_token": newRefresh.Token,
	})
}

// AdminListUsers godoc
// @ID           AdminListUsers
// @Summary      Admin: list all users
// @Description  Returns paginated list of users (admin only)
// @Tags         admin
// @Produce      json
// @Param        page      query int false "Page number (1-based)"
// @Param        page_size query int false "Page size (max 200)"
// @Success      200 {object} ListUsersOut
// @Failure      401 {string} string "unauthorized"
// @Failure      403 {string} string "forbidden"
// @Security     BearerAuth
// @Router       /api/v1/admin/users [get]
func AdminListUsers(w http.ResponseWriter, r *http.Request) {
	ctx := middleware.GetAuthContext(r)
	if ctx == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Load current user to check global role
	var me models.User
	if err := db.DB.Select("id, role").First(&me, "id = ?", ctx.UserID).Error; err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if me.Role != "admin" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	// Pagination
	page := mustInt(r.URL.Query().Get("page"), 1)
	if page < 1 {
		page = 1
	}
	pageSize := mustInt(r.URL.Query().Get("page_size"), 50)
	if pageSize < 1 {
		pageSize = 50
	}
	if pageSize > 200 {
		pageSize = 200
	}
	offset := (page - 1) * pageSize

	// Query
	var total int64
	_ = db.DB.Model(&models.User{}).Count(&total).Error

	var users []models.User
	if err := db.DB.
		Model(&models.User{}).
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&users).Error; err != nil {
		http.Error(w, "failed to fetch users", http.StatusInternalServerError)
		return
	}

	out := make([]UserListItem, len(users))
	for i, u := range users {
		out[i] = UserListItem{
			ID:            u.ID,
			Name:          u.Name,
			Email:         u.Email,
			EmailVerified: u.EmailVerified,
			Role:          string(u.Role),
			CreatedAt:     u.CreatedAt,
			UpdatedAt:     u.UpdatedAt,
		}
	}

	_ = response.JSON(w, http.StatusOK, ListUsersOut{
		Users:    out,
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	})
}

// AdminCreateUser godoc
// @ID           AdminCreateUser
// @Summary      Admin: create user
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        body body AdminCreateUserRequest true "payload"
// @Success      201 {object} userOut
// @Failure      400 {string} string "bad request"
// @Failure      401 {string} string "unauthorized"
// @Failure      403 {string} string "forbidden"
// @Failure      409 {string} string "conflict"
// @Security     BearerAuth
// @Router       /api/v1/admin/users [post]
func AdminCreateUser(w http.ResponseWriter, r *http.Request) {
	if _, ok := requireGlobalAdmin(w, r); !ok {
		return
	}

	var in struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"` // "user" | "admin"
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	in.Email = strings.TrimSpace(strings.ToLower(in.Email))
	in.Role = strings.TrimSpace(in.Role)
	if in.Email == "" || in.Password == "" {
		http.Error(w, "email and password required", http.StatusBadRequest)
		return
	}
	if in.Role == "" {
		in.Role = "user"
	}
	if in.Role != "user" && in.Role != "admin" {
		http.Error(w, "invalid role", http.StatusBadRequest)
		return
	}

	var exists int64
	if err := db.DB.Model(&models.User{}).Where("LOWER(email)=?", in.Email).Count(&exists).Error; err == nil && exists > 0 {
		http.Error(w, "email already in use", http.StatusConflict)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "hash error", http.StatusInternalServerError)
		return
	}

	u := models.User{
		Name:     in.Name,
		Email:    in.Email,
		Password: string(hash),
		Role:     models.Role(in.Role),
	}
	if err := db.DB.Create(&u).Error; err != nil {
		http.Error(w, "create failed", http.StatusInternalServerError)
		return
	}
	_ = response.JSON(w, http.StatusCreated, asUserOut(u))
}

// AdminUpdateUser godoc
// @ID           AdminUpdateUser
// @Summary      Admin: update user
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        userId path string true "User ID"
// @Param        body body AdminUpdateUserRequest true "payload"
// @Success      200 {object} userOut
// @Failure      400 {string} string "bad request"
// @Failure      401 {string} string "unauthorized"
// @Failure      403 {string} string "forbidden"
// @Failure      404 {string} string "not found"
// @Failure      409 {string} string "conflict"
// @Security     BearerAuth
// @Router       /api/v1/admin/users/{userId} [patch]
func AdminUpdateUser(w http.ResponseWriter, r *http.Request) {
	_, ok := requireGlobalAdmin(w, r)
	if !ok {
		return
	}

	idStr := chi.URLParam(r, "userId")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "bad user id", http.StatusBadRequest)
		return
	}

	var u models.User
	if err := db.DB.First(&u, "id = ?", id).Error; err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	var in struct {
		Name     *string `json:"name"`
		Email    *string `json:"email"`
		Password *string `json:"password"`
		Role     *string `json:"role"` // "user" | "admin"
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	updates := map[string]any{}
	if in.Name != nil {
		updates["name"] = *in.Name
	}
	if in.Email != nil {
		email := strings.TrimSpace(strings.ToLower(*in.Email))
		if email == "" {
			http.Error(w, "email required", http.StatusBadRequest)
			return
		}
		var exists int64
		_ = db.DB.Model(&models.User{}).Where("LOWER(email)=? AND id <> ?", email, u.ID).Count(&exists).Error
		if exists > 0 {
			http.Error(w, "email already in use", http.StatusConflict)
			return
		}
		updates["email"] = email
	}
	if in.Password != nil && *in.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(*in.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "hash error", http.StatusInternalServerError)
			return
		}
		updates["password"] = string(hash)
	}
	if in.Role != nil {
		role := strings.TrimSpace(*in.Role)
		if role != "user" && role != "admin" {
			http.Error(w, "invalid role", http.StatusBadRequest)
			return
		}
		// prevent demoting the last admin
		if u.Role == "admin" && role == "user" {
			n, _ := adminCount(&u.ID)
			if n == 0 {
				http.Error(w, "cannot demote last admin", http.StatusConflict)
				return
			}
		}
		updates["role"] = role
	}
	if len(updates) == 0 {
		_ = response.JSON(w, http.StatusOK, asUserOut(u))
		return
	}
	if err := db.DB.Model(&u).Updates(updates).Error; err != nil {
		http.Error(w, "update failed", http.StatusInternalServerError)
		return
	}
	_ = response.JSON(w, http.StatusOK, asUserOut(u))
}

// AdminDeleteUser godoc
// @ID           AdminDeleteUser
// @Summary      Admin: delete user
// @Tags         admin
// @Produce      json
// @Param        userId path string true "User ID"
// @Success      204 {string} string "no content"
// @Failure      400 {string} string "bad request"
// @Failure      401 {string} string "unauthorized"
// @Failure      403 {string} string "forbidden"
// @Failure      404 {string} string "not found"
// @Failure      409 {string} string "conflict"
// @Security     BearerAuth
// @Router       /api/v1/admin/users/{userId} [delete]
func AdminDeleteUser(w http.ResponseWriter, r *http.Request) {
	me, ok := requireGlobalAdmin(w, r)
	if !ok {
		return
	}

	idStr := chi.URLParam(r, "userId")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "bad user id", http.StatusBadRequest)
		return
	}

	if me.ID == id {
		http.Error(w, "cannot delete self", http.StatusBadRequest)
		return
	}

	var u models.User
	if err := db.DB.First(&u, "id = ?", id).Error; err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	if u.Role == "admin" {
		n, _ := adminCount(&u.ID)
		if n == 0 {
			http.Error(w, "cannot delete last admin", http.StatusConflict)
			return
		}
	}

	if err := db.DB.Delete(&models.User{}, "id = ?", id).Error; err != nil {
		http.Error(w, "delete failed", http.StatusInternalServerError)
		return
	}
	response.NoContent(w)
}
