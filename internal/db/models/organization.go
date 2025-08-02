package models

import "github.com/google/uuid"

type Organization struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name     string    `gorm:"not null" json:"name"`
	Slug     string    `gorm:"unique" json:"slug"`
	Logo     string    `json:"logo"`
	Metadata string    `json:"metadata"`
	Timestamped
}
