package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Domain struct {
	ID             uuid.UUID    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrganizationID uuid.UUID    `gorm:"type:uuid;not null;index;uniqueIndex:uniq_org_domain,priority:1"`
	Organization   Organization `gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE" json:"organization"`
	DomainName     string       `gorm:"type:varchar(253);not null;uniqueIndex:uniq_org_domain,priority:2"`
	ZoneID         string       `gorm:"type:varchar(128);not null;default:''"`       // backfilled for R53 (e.g. "/hostedzone/Z123...")
	Status         string       `gorm:"type:varchar(20);not null;default:'pending'"` // pending, provisioning, ready, failed
	LastError      string       `gorm:"type:text;not null;default:''"`
	CredentialID   uuid.UUID    `gorm:"type:uuid;not null" json:"credential_id"`
	Credential     Credential   `gorm:"foreignKey:CredentialID" json:"credential,omitempty"`
	CreatedAt      time.Time    `json:"created_at,omitempty" gorm:"type:timestamptz;column:created_at;not null;default:now()"`
	UpdatedAt      time.Time    `json:"updated_at,omitempty" gorm:"type:timestamptz;autoUpdateTime;column:updated_at;not null;default:now()"`
}

type RecordSet struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	DomainID    uuid.UUID      `gorm:"type:uuid;not null;index"`
	Domain      Domain         `gorm:"foreignKey:DomainID;constraint:OnDelete:CASCADE"`
	Name        string         `gorm:"type:varchar(253);not null"`      // e.g. "endpoint" (relative to DomainName)
	Type        string         `gorm:"type:varchar(10);not null;index"` // A, AAAA, CNAME, TXT, MX, SRV, NS, CAA...
	TTL         *int           `gorm:""`                                // nil for alias targets (Route 53 ignores TTL for alias)
	Values      datatypes.JSON `gorm:"type:jsonb;not null;default:'[]'"`
	Fingerprint string         `gorm:"type:char(64);not null;index"` // sha256 of canonical(name,type,ttl,values|alias)
	Status      string         `gorm:"type:varchar(20);not null;default:'pending'"`
	Owner       string         `gorm:"type:varchar(16);not null;default:'unknown'"` // 'autoglue' | 'external' | 'unknown'
	LastError   string         `gorm:"type:text;not null;default:''"`
	_           struct{}       `gorm:"uniqueIndex:uniq_domain_name_type,priority:1"` // tag holder
	_           struct{}       `gorm:"uniqueIndex:uniq_domain_name_type,priority:2"`
	_           struct{}       `gorm:"uniqueIndex:uniq_domain_name_type,priority:3"`
	CreatedAt   time.Time      `json:"created_at,omitempty" gorm:"type:timestamptz;column:created_at;not null;default:now()"`
	UpdatedAt   time.Time      `json:"updated_at,omitempty" gorm:"type:timestamptz;autoUpdateTime;column:updated_at;not null;default:now()"`
}
