package authentication

import (
	"encoding/json"
	"github.com/glueops/autoglue/internal/config"
	"github.com/glueops/autoglue/internal/db"
	"github.com/glueops/autoglue/internal/db/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

var jwtSecret = []byte(config.GetAuthSecret())

type RegisterInput struct {
	Email    string `json:"email" example:"me@here.com"`
	Name     string `json:"name" example:"My Name"`
	Password string `json:"password" example:"123456"`
}

type LoginInput struct {
	Email    string `json:"email" example:"me@here.com"`
	Password string `json:"password" example:"123456"`
}

// Register godoc
// @Summary      Register a new user
// @Description  Registers a new user and stores credentials
// @Tags         Auth
// @Accept       json
// @Produce      plain
// @Param        body  body      RegisterInput  true  "User registration input"
// @Success      201   {string}  string         "created"
// @Failure      400   {string}  string         "bad request"
// @Router       /api/v1/authentication/register [post]
func Register(w http.ResponseWriter, r *http.Request) {
	var input RegisterInput
	json.NewDecoder(r.Body).Decode(&input)
	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	if err != nil {
		http.Error(w, "failed to hash password", http.StatusInternalServerError)
		return
	}

	user := models.User{ID: uuid.NewString(), Email: input.Email, Password: string(hashed), Name: input.Name, Role: "user"}
	if err := db.DB.Create(&user).Error; err != nil {
		http.Error(w, "registration failed", 400)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// Login godoc
// @Summary      Authenticate and return a token
// @Description  Authenticates a user and returns a JWT bearer token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      LoginInput     true  "User login input"
// @Success      200   {object}  map[string]string "token"
// @Failure      401   {string}  string         "unauthorized"
// @Router       /api/v1/authentication/login [post]
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

	// Refresh token (long-lived)
	refreshTokenStr := uuid.NewString()

	db.DB.Create(&models.RefreshToken{
		ID:        uuid.NewString(),
		UserID:    user.ID,
		Token:     refreshTokenStr,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		Revoked:   false,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"access_token":  accessStr,
		"refresh_token": refreshTokenStr,
	})
}

// Refresh godoc
// @Summary      Refresh access token
// @Description  Use a refresh token to obtain a new access token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      map[string]string  true  "refresh_token"
// @Success      200   {object}  map[string]string "new access token"
// @Failure      401   {string}  string         "unauthorized"
// @Router       /api/v1/authentication/refresh [post]
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
		"exp": time.Now().Add(15 * time.Minute).Unix(),
	}
	newAccess := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	newToken, _ := newAccess.SignedString(jwtSecret)

	json.NewEncoder(w).Encode(map[string]string{
		"access_token": newToken,
	})
}

// Logout godoc
// @Summary      Logout user
// @Description  Revoke a refresh token
// @Tags         Auth
// @Accept       json
// @Produce      plain
// @Param        body  body      map[string]string  true  "refresh_token"
// @Success      204   {string}  string         "no content"
// @Router       /api/v1/authentication/logout [post]
func Logout(w http.ResponseWriter, r *http.Request) {
	var input struct {
		RefreshToken string `json:"refresh_token"`
	}
	json.NewDecoder(r.Body).Decode(&input)

	db.DB.Model(&models.RefreshToken{}).Where("token = ?", input.RefreshToken).Update("revoked", true)
	w.WriteHeader(http.StatusNoContent)
}
