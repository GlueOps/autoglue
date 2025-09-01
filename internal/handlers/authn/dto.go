package authn

import (
	"sync"
	"time"

	"github.com/glueops/autoglue/internal/config"
	"github.com/glueops/autoglue/internal/db/models"
	appsmtp "github.com/glueops/autoglue/internal/smtp"
	"github.com/google/uuid"
)

var jwtSecret = []byte(config.GetAuthSecret())
var (
	mailerOnce sync.Once
	mailer     *appsmtp.Mailer
	mailerErr  error
)

const (
	resetTTL         = 1 * time.Hour  // password reset token validity
	verifyTTL        = 48 * time.Hour // email verification token validity
	refreshTTL       = 7 * 24 * time.Hour
	accessTTL        = 72 * time.Hour
	rotatedAccessTTL = 15 * time.Minute
)

type RegisterInput struct {
	Email    string `json:"email" example:"me@here.com"`
	Name     string `json:"name" example:"My Name"`
	Password string `json:"password" example:"123456"`
}

type LoginInput struct {
	Email    string `json:"email" example:"me@here.com"`
	Password string `json:"password" example:"123456"`
}

type UserDTO struct {
	ID            uuid.UUID   `json:"id"`
	Name          string      `json:"name"`
	Email         string      `json:"email"`
	EmailVerified bool        `json:"email_verified"`
	Role          models.Role `json:"role"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
}

type AuthClaimsDTO struct {
	Orgs      []string `json:"orgs,omitempty"`
	Roles     []string `json:"roles,omitempty"`
	Issuer    string   `json:"iss,omitempty"`
	Subject   string   `json:"sub,omitempty"`
	Audience  []string `json:"aud,omitempty"`
	ExpiresAt int64    `json:"exp,omitempty"`
	IssuedAt  int64    `json:"iat,omitempty"`
	NotBefore int64    `json:"nbf,omitempty"`
}

type MeResponse struct {
	User           UserDTO        `json:"user_id"`
	OrganizationID *string        `json:"organization_id,omitempty"`
	OrgRole        string         `json:"org_role,omitempty"`
	Claims         *AuthClaimsDTO `json:"claims,omitempty"`
}

type VerifyEmailData struct {
	Name            string
	Email           string
	Token           string
	VerificationURL string
}

type PasswordResetData struct {
	Name     string
	Email    string
	Token    string
	ResetURL string
}
