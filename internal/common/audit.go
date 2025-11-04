package common

import (
	"time"

	"github.com/google/uuid"
)

type AuditFields struct {
	ID             uuid.UUID `json:"id"                   gorm:"type:uuid;default:gen_random_uuid()"`
	OrganizationID uuid.UUID `json:"organization_id"      gorm:"type:uuid;index"`
	CreatedAt      time.Time `json:"created_at,omitempty" gorm:"column:created_at;not null;default:now()"`
	UpdatedAt      time.Time `json:"updated_at,omitempty" gorm:"autoUpdateTime;column:updated_at;not null;default:now()"`
}
