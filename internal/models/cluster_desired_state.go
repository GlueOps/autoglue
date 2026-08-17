package models

import (
	"time"

	"github.com/google/uuid"
)

// ClusterDesiredState is the per-cluster generation pointer the sync endpoint
// compares an agent's claim against.
//
// One row per cluster, and Generation only ever increases. It is a counter
// rather than a hash or a timestamp because the agent's question is ordinal —
// "is there anything newer than N" — and because two clusters converging at
// different speeds must not share a sequence.
type ClusterDesiredState struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	OrganizationID uuid.UUID `gorm:"type:uuid;not null;index" json:"organization_id"`

	ClusterID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"cluster_id"`
	Cluster   Cluster   `gorm:"foreignKey:ClusterID;constraint:OnDelete:CASCADE" json:"-"`

	// Generation is bumped by `generation + 1` inside the same transaction that
	// writes the new DesiredResource rows, so a reader can never observe a
	// pointer aimed at a half-written snapshot.
	Generation int64 `gorm:"not null;default:0" json:"generation"`

	PublishedAt time.Time `gorm:"type:timestamptz;not null;default:now()" json:"published_at"`

	CreatedAt time.Time `json:"created_at,omitempty" gorm:"type:timestamptz;column:created_at;not null;default:now()"`
	UpdatedAt time.Time `json:"updated_at,omitempty" gorm:"type:timestamptz;autoUpdateTime;column:updated_at;not null;default:now()"`
}
