package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Credential struct {
	ID               uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrganizationID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"organization_id"`
	Organization     Organization   `gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE" json:"organization"`
	Provider         string         `gorm:"type:varchar(50);not null;uniqueIndex:uniq_org_provider_scopekind_scope,priority:2;index:idx_provider_kind"`
	Kind             string         `gorm:"type:varchar(50);not null;index:idx_provider_kind;index:idx_kind_scope"`
	ScopeKind        string         `gorm:"type:varchar(20);not null;uniqueIndex:uniq_org_provider_scopekind_scope,priority:3"`
	Scope            datatypes.JSON `gorm:"type:jsonb;not null;default:'{}';index:idx_kind_scope"`
	ScopeFingerprint string         `gorm:"type:char(64);not null;uniqueIndex:uniq_org_provider_scopekind_scope,priority:4;index"`
	SchemaVersion    int            `gorm:"not null;default:1"`
	Name             string         `gorm:"type:varchar(100);not null;default:''"`
	ScopeVersion     int            `gorm:"not null;default:1"`
	AccountID        string         `gorm:"type:varchar(32)"`
	Region           string         `gorm:"type:varchar(32)"`
	EncryptedData    string         `gorm:"not null"`
	IV               string         `gorm:"not null"`
	Tag              string         `gorm:"not null"`
	CreatedAt        time.Time      `json:"created_at,omitempty" gorm:"type:timestamptz;column:created_at;not null;default:now()" format:"date-time"`
	UpdatedAt        time.Time      `json:"updated_at,omitempty" gorm:"type:timestamptz;autoUpdateTime;column:updated_at;not null;default:now()" format:"date-time"`
}
