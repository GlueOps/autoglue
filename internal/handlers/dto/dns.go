package dto

import (
	"encoding/json"
	"strings"

	"github.com/go-playground/validator/v10"
)

var dnsValidate = validator.New()

func init() {
	_ = dnsValidate.RegisterValidation("fqdn", func(fl validator.FieldLevel) bool {
		s := strings.TrimSpace(fl.Field().String())
		if s == "" || len(s) > 253 {
			return false
		}
		// Minimal: lower-cased, no trailing dot in our API (normalize server-side)
		// You can add stricter checks later.
		return !strings.HasPrefix(s, ".") && !strings.Contains(s, "..")
	})
	_ = dnsValidate.RegisterValidation("rrtype", func(fl validator.FieldLevel) bool {
		switch strings.ToUpper(fl.Field().String()) {
		case "A", "AAAA", "CNAME", "TXT", "MX", "NS", "SRV", "CAA":
			return true
		default:
			return false
		}
	})
}

// ---- Domains ----

type CreateDomainRequest struct {
	DomainName   string `json:"domain_name" validate:"required,fqdn"`
	CredentialID string `json:"credential_id" validate:"required,uuid4"`
	ZoneID       string `json:"zone_id,omitempty" validate:"omitempty,max=128"`
}

type UpdateDomainRequest struct {
	CredentialID *string `json:"credential_id,omitempty" validate:"omitempty,uuid4"`
	ZoneID       *string `json:"zone_id,omitempty" validate:"omitempty,max=128"`
	Status       *string `json:"status,omitempty" validate:"omitempty,oneof=pending provisioning ready failed"`
	DomainName   *string `json:"domain_name,omitempty" validate:"omitempty,fqdn"`
}

type DomainResponse struct {
	ID             string `json:"id"`
	OrganizationID string `json:"organization_id"`
	DomainName     string `json:"domain_name"`
	ZoneID         string `json:"zone_id"`
	Status         string `json:"status"`
	LastError      string `json:"last_error"`
	CredentialID   string `json:"credential_id"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

// ---- Record Sets ----

type AliasTarget struct {
	HostedZoneID         string `json:"hosted_zone_id" validate:"required"`
	DNSName              string `json:"dns_name" validate:"required"`
	EvaluateTargetHealth bool   `json:"evaluate_target_health"`
}

type CreateRecordSetRequest struct {
	// Name relative to domain ("endpoint") OR FQDN ("endpoint.example.com").
	// Server normalizes to relative.
	Name   string   `json:"name" validate:"required,max=253"`
	Type   string   `json:"type" validate:"required,rrtype"`
	TTL    *int     `json:"ttl,omitempty" validate:"omitempty,gte=1,lte=86400"`
	Values []string `json:"values" validate:"omitempty,dive,min=1,max=1024"`
}

type UpdateRecordSetRequest struct {
	// Any change flips status back to pending (worker will UPSERT)
	Name   *string   `json:"name,omitempty" validate:"omitempty,max=253"`
	Type   *string   `json:"type,omitempty" validate:"omitempty,rrtype"`
	TTL    *int      `json:"ttl,omitempty" validate:"omitempty,gte=1,lte=86400"`
	Values *[]string `json:"values,omitempty" validate:"omitempty,dive,min=1,max=1024"`
	Status *string   `json:"status,omitempty" validate:"omitempty,oneof=pending provisioning ready failed"`
}

type RecordSetResponse struct {
	ID          string          `json:"id"`
	DomainID    string          `json:"domain_id"`
	Name        string          `json:"name"`
	Type        string          `json:"type"`
	TTL         *int            `json:"ttl,omitempty"`
	Values      json.RawMessage `json:"values" swaggertype:"object"` // []string JSON
	Fingerprint string          `json:"fingerprint"`
	Status      string          `json:"status"`
	LastError   string          `json:"last_error"`
	Owner       string          `json:"owner"`
	CreatedAt   string          `json:"created_at"`
	UpdatedAt   string          `json:"updated_at"`
}

// DNSValidate Quick helper to validate DTOs in handlers
func DNSValidate(i any) error {
	return dnsValidate.Struct(i)
}
