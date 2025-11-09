package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Credential struct {
	ID               uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrganizationID   uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex:idx_credentials_org_provider" json:"organization_id"`
	Organization     Organization   `gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE" json:"organization"`
	Provider         string         `gorm:"type:varchar(50);not null;index"`
	Kind             string         `gorm:"type:varchar(50);not null;index"` // "aws_access_key", "api_token", "basic_auth", ...
	SchemaVersion    int            `gorm:"not null;default:1"`
	Name             string         `gorm:"type:varchar(100);not null;default:''"` // human label, lets you have multiple for same service
	ScopeKind        string         `gorm:"type:varchar(20);not null"`             // "provider" | "service" | "resource"
	Scope            datatypes.JSON `gorm:"type:jsonb;not null;default:'{}'"`      // e.g. {"service":"route53"} or {"arn":"arn:aws:s3:::my-bucket"}
	ScopeVersion     int            `gorm:"not null;default:1"`
	AccountID        string         `gorm:"type:varchar(32)"` // AWS account ID if applicable
	Region           string         `gorm:"type:varchar(32)"` // default region (non-secret)
	ScopeFingerprint string         `gorm:"type:char(64);not null;index"`
	EncryptedData    string         `gorm:"not null"`
	IV               string         `gorm:"not null"`
	Tag              string         `gorm:"not null"`
	CreatedAt        time.Time      `json:"created_at,omitempty" gorm:"type:timestamptz;column:created_at;not null;default:now()" format:"date-time"`
	UpdatedAt        time.Time      `json:"updated_at,omitempty" gorm:"type:timestamptz;autoUpdateTime;column:updated_at;not null;default:now()" format:"date-time"`
}
