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

type UserListItem struct {
	ID            any    `json:"id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Role          string `json:"role"`
	CreatedAt     any    `json:"created_at"`
	UpdatedAt     any    `json:"updated_at"`
}

type ListUsersOut struct {
	Users    []UserListItem `json:"users"`
	Page     int            `json:"page"`
	PageSize int            `json:"page_size"`
	Total    int64          `json:"total"`
}

type userOut struct {
	ID            any    `json:"id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Role          string `json:"role"`
	CreatedAt     any    `json:"created_at"`
	UpdatedAt     any    `json:"updated_at"`
}

type AdminCreateUserRequest struct {
	Name     string `json:"name" example:"Jane Doe"`
	Email    string `json:"email" example:"jane@example.com"`
	Password string `json:"password" example:"Secret123!"`
	// Role allowed values: "user" or "admin"
	Role string `json:"role" example:"user" enums:"user,admin"`
}

type AdminUpdateUserRequest struct {
	Name     *string `json:"name,omitempty" example:"Jane Doe"`
	Email    *string `json:"email,omitempty" example:"jane@example.com"`
	Password *string `json:"password,omitempty" example:"NewSecret123!"`
	Role     *string `json:"role,omitempty" example:"admin" enums:"user,admin"`
}

type AdminUserResponse struct {
	ID            uuid.UUID `json:"id" example:"6aa012bc-ce8a-4cd9-9971-58f3917037f8"`
	Name          string    `json:"name" example:"Jane Doe"`
	Email         string    `json:"email" example:"jane@example.com"`
	EmailVerified bool      `json:"email_verified" example:"false"`
	Role          string    `json:"role" example:"user"`
	CreatedAt     string    `json:"created_at" example:"2025-09-01T08:38:12Z"`
	UpdatedAt     string    `json:"updated_at" example:"2025-09-01T17:02:36Z"`
}

type AdminListUsersResponse struct {
	Users    []AdminUserResponse `json:"users"`
	Page     int                 `json:"page" example:"1"`
	PageSize int                 `json:"page_size" example:"50"`
	Total    int64               `json:"total" example:"123"`
}
