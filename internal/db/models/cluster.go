package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Cluster struct {
	ID         string         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	OrgID      uuid.UUID      `gorm:"type:uuid;index" json:"org_id"`
	Name       string         `gorm:"uniqueIndex" json:"name"`
	Provider   string         `json:"provider"`
	Region     string         `json:"region"`
	Status     string         `json:"status"`
	Kubeconfig string         `gorm:"type:text" json:"kubeconfig"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}
