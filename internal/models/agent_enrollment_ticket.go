package models

import (
	"time"

	"github.com/google/uuid"
)

// AgentEnrollmentTicket is a single-use bearer token that buys exactly one
// agent credential for one cluster.
//
// It exists because enrolment cannot be authenticated by the thing it is
// issuing. The ticket is minted control-plane side and delivered over the
// channel that already proves possession of the bastion — SSH, alongside the
// rest of the cluster assets — so redeeming it demonstrates the caller reached
// that host. Short TTL and single use are what stop a ticket read out of a
// backup from minting an agent months later.
type AgentEnrollmentTicket struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	OrganizationID uuid.UUID `gorm:"type:uuid;not null;index" json:"organization_id"`

	ClusterID uuid.UUID `gorm:"type:uuid;not null;index" json:"cluster_id"`
	Cluster   Cluster   `gorm:"foreignKey:ClusterID;constraint:OnDelete:CASCADE" json:"-"`

	ServerID uuid.UUID `gorm:"type:uuid;not null;index" json:"server_id"`

	// TicketHash is SHA256 of the plaintext. The plaintext is returned to the
	// minting caller once, written to the bastion, and never persisted here —
	// same rule as the agent secret, for the same reason.
	TicketHash string `gorm:"type:text;not null;uniqueIndex" json:"-"`
	Prefix     string `gorm:"type:text;not null;default:''" json:"prefix"`

	ExpiresAt time.Time `gorm:"type:timestamptz;not null;index" json:"expires_at"`

	// RedeemedAt is the single-use latch. Nullable rather than a bool so the
	// redemption is also its own audit record, and so the check and the claim
	// can be one conditional UPDATE.
	RedeemedAt        *time.Time `gorm:"type:timestamptz" json:"redeemed_at,omitempty"`
	RedeemedByAgentID *uuid.UUID `gorm:"type:uuid" json:"redeemed_by_agent_id,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty" gorm:"type:timestamptz;column:created_at;not null;default:now()"`
	UpdatedAt time.Time `json:"updated_at,omitempty" gorm:"type:timestamptz;autoUpdateTime;column:updated_at;not null;default:now()"`
}
