package dto

import (
	"encoding/json"

	"github.com/go-playground/validator/v10"
)

// RawJSON is a swagger-friendly wrapper for json.RawMessage.
type RawJSON json.RawMessage

var Validate = validator.New()

func init() {
	_ = Validate.RegisterValidation("awsarn", func(fl validator.FieldLevel) bool {
		v := fl.Field().String()
		return len(v) > 10 && len(v) < 2048 && len(v) >= 4 && v[:4] == "arn:"
	})
}

/*** Shapes for secrets ***/

type AWSCredential struct {
	AccessKeyID     string `json:"access_key_id" validate:"required,alphanum,len=20"`
	SecretAccessKey string `json:"secret_access_key" validate:"required"`
	Region          string `json:"region" validate:"omitempty"`
}

type BasicAuth struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type APIToken struct {
	Token string `json:"token" validate:"required"`
}

type OAuth2Credential struct {
	ClientID     string `json:"client_id" validate:"required"`
	ClientSecret string `json:"client_secret" validate:"required"`
	RefreshToken string `json:"refresh_token" validate:"required"`
}

/*** Shapes for scopes ***/

type AWSProviderScope struct{}

type AWSServiceScope struct {
	Service string `json:"service" validate:"required,oneof=route53 s3 ec2 iam rds dynamodb"`
}

type AWSResourceScope struct {
	ARN string `json:"arn" validate:"required,awsarn"`
}

/*** Registries ***/

type ProviderDef struct {
	New      func() any
	Validate func(any) error
}

type ScopeDef struct {
	New         func() any
	Validate    func(any) error
	Specificity int // 0=provider, 1=service, 2=resource
}

// Secret shapes per provider/kind/version

var CredentialRegistry = map[string]map[string]map[int]ProviderDef{
	"aws": {
		"aws_access_key": {
			1: {New: func() any { return &AWSCredential{} }, Validate: func(x any) error { return Validate.Struct(x) }},
		},
	},
	"cloudflare":   {"api_token": {1: {New: func() any { return &APIToken{} }, Validate: func(x any) error { return Validate.Struct(x) }}}},
	"hetzner":      {"api_token": {1: {New: func() any { return &APIToken{} }, Validate: func(x any) error { return Validate.Struct(x) }}}},
	"digitalocean": {"api_token": {1: {New: func() any { return &APIToken{} }, Validate: func(x any) error { return Validate.Struct(x) }}}},
	"generic": {
		"basic_auth": {1: {New: func() any { return &BasicAuth{} }, Validate: func(x any) error { return Validate.Struct(x) }}},
		"oauth2":     {1: {New: func() any { return &OAuth2Credential{} }, Validate: func(x any) error { return Validate.Struct(x) }}},
	},
}

// Scope shapes per provider/scopeKind/version

var ScopeRegistry = map[string]map[string]map[int]ScopeDef{
	"aws": {
		"provider": {1: {New: func() any { return &AWSProviderScope{} }, Validate: func(any) error { return nil }, Specificity: 0}},
		"service":  {1: {New: func() any { return &AWSServiceScope{} }, Validate: func(x any) error { return Validate.Struct(x) }, Specificity: 1}},
		"resource": {1: {New: func() any { return &AWSResourceScope{} }, Validate: func(x any) error { return Validate.Struct(x) }, Specificity: 2}},
	},
}

/*** API DTOs used by swagger ***/

// CreateCredentialRequest represents the POST /credentials payload
type CreateCredentialRequest struct {
	Provider      string  `json:"provider" validate:"required,oneof=aws cloudflare hetzner digitalocean generic"`
	Kind          string  `json:"kind" validate:"required"`                 // aws_access_key, api_token, basic_auth, oauth2
	SchemaVersion int     `json:"schema_version" validate:"required,gte=1"` // secret schema version
	Name          string  `json:"name" validate:"omitempty,max=100"`        // human label
	ScopeKind     string  `json:"scope_kind" validate:"required,oneof=provider service resource"`
	ScopeVersion  int     `json:"scope_version" validate:"required,gte=1"`        // scope schema version
	Scope         RawJSON `json:"scope" validate:"required" swaggertype:"object"` // {"service":"route53"} or {"arn":"..."}
	AccountID     string  `json:"account_id,omitempty" validate:"omitempty,max=32"`
	Region        string  `json:"region,omitempty" validate:"omitempty,max=32"`
	Secret        RawJSON `json:"secret" validate:"required" swaggertype:"object"` // encrypted later
}

// UpdateCredentialRequest represents PATCH /credentials/{id}
type UpdateCredentialRequest struct {
	Name         *string  `json:"name,omitempty"`
	AccountID    *string  `json:"account_id,omitempty"`
	Region       *string  `json:"region,omitempty"`
	ScopeKind    *string  `json:"scope_kind,omitempty"`
	ScopeVersion *int     `json:"scope_version,omitempty"`
	Scope        *RawJSON `json:"scope,omitempty" swaggertype:"object"`
	Secret       *RawJSON `json:"secret,omitempty" swaggertype:"object"` // set if rotating

}

// CredentialOut is what we return (no secrets)
type CredentialOut struct {
	ID            string  `json:"id"`
	Provider      string  `json:"provider"`
	Kind          string  `json:"kind"`
	SchemaVersion int     `json:"schema_version"`
	Name          string  `json:"name"`
	ScopeKind     string  `json:"scope_kind"`
	ScopeVersion  int     `json:"scope_version"`
	Scope         RawJSON `json:"scope" swaggertype:"object"`
	AccountID     string  `json:"account_id,omitempty"`
	Region        string  `json:"region,omitempty"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}
