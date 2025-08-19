package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Cluster struct {
	ID             string         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	OrganizationID uuid.UUID      `gorm:"type:uuid;not null" json:"organization_id"`
	Organization   Organization   `gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE" json:"organization"`
	ServerID       uuid.UUID      `gorm:"type:uuid;not null" json:"server_id"`
	Server         Server         `gorm:"foreignKey:ServerID;constraint:OnDelete:CASCADE" json:"server"`
	Name           string         `json:"name"`
	Provider       string         `json:"provider"`
	Region         string         `json:"region"`
	Status         string         `json:"status"`
	Kubeconfig     string         `gorm:"type:text" json:"kubeconfig"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}
