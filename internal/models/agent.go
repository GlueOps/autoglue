package models

import (
	"time"

	"github.com/google/uuid"
)

// Agent lifecycle. Re-enrolment revokes rather than deletes: an agent's tasks
// are the only record of what ever ran on a bastion, and a dead-lettered task
// has to outlive the credential that produced it.
const (
	AgentStatusActive  = "active"
	AgentStatusRevoked = "revoked"
)

// Agent is one bastion's identity in the control plane.
//
// The credential is a tuple — id, key, secret — so that no part of it is
// derivable: the id names the row, the key is the indexed lookup, and the
// secret is the only thing that proves possession. It is cluster-scoped and
// short-lived by design. The org-wide 24h key currently pushed to bastions is
// what this model exists to retire: an agent that leaks must reach exactly one
// cluster, not the whole organization.
type Agent struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	// OrganizationID is denormalised so agent-plane writes can stamp job_logs
	// and authorize without joining back through the cluster on every poll.
	OrganizationID uuid.UUID `gorm:"type:uuid;not null;index" json:"organization_id"`

	// ClusterID is the entire authorization scope. Every agent-plane query
	// filters on it rather than on the org.
	//
	// The unique index has to be partial: bastions are rebuilt and re-enrol,
	// so revoked rows accumulate, and a plain unique index would make the
	// second enrolment fail. Partial expresses the actual invariant — one
	// *live* agent per cluster — and enforces it in the database rather than
	// in whichever handler remembers to check.
	ClusterID uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:uniq_agents_live_cluster,where:status = 'active'" json:"cluster_id"`
	Cluster   Cluster   `gorm:"foreignKey:ClusterID;constraint:OnDelete:CASCADE" json:"-"`

	// ServerID is the bastion this agent runs on. Carried for operator
	// diagnosis; authorization never consults it, because a cluster may be
	// re-bastioned without the agent's scope changing.
	ServerID uuid.UUID `gorm:"type:uuid;not null;index" json:"server_id"`
	Server   Server    `gorm:"foreignKey:ServerID;constraint:OnDelete:CASCADE" json:"-"`

	// KeyHash is SHA256 of the public half, which makes the lookup a plain
	// indexed equality. That is safe precisely because the key is not the
	// secret; the secret is verified with argon2id, which is constant time.
	KeyHash    string `gorm:"type:text;not null;uniqueIndex" json:"-"`
	SecretHash string `gorm:"type:text;not null" json:"-"`
	Prefix     string `gorm:"type:text;not null;default:''" json:"prefix"`

	Status string `gorm:"type:text;not null;default:'active';index" json:"status"`

	// Version is whatever the agent reported at enrolment. An agent behind the
	// control plane is the normal state of deployed software, so version skew
	// is data to read, not a condition to reject.
	Version string `gorm:"type:text;not null;default:''" json:"version"`

	EnrolledAt time.Time  `gorm:"type:timestamptz;not null;default:now()" json:"enrolled_at"`
	LastSeenAt *time.Time `gorm:"type:timestamptz;index" json:"last_seen_at,omitempty"`
	ExpiresAt  *time.Time `gorm:"type:timestamptz;index" json:"expires_at,omitempty"`
	RevokedAt  *time.Time `gorm:"type:timestamptz" json:"revoked_at,omitempty"`

	// Config-plane observation. These are what the agent last *said*, never
	// what the control plane wishes were true — the desired generation lives on
	// ClusterDesiredState, and the gap between the two is the only honest
	// measure of convergence.
	ReportedGeneration int64      `gorm:"not null;default:0" json:"reported_generation"`
	AppliedGeneration  int64      `gorm:"not null;default:0" json:"applied_generation"`
	Healthy            bool       `gorm:"not null;default:false" json:"healthy"`
	LastReconcileAt    *time.Time `gorm:"type:timestamptz" json:"last_reconcile_at,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty" gorm:"type:timestamptz;column:created_at;not null;default:now()"`
	UpdatedAt time.Time `json:"updated_at,omitempty" gorm:"type:timestamptz;autoUpdateTime;column:updated_at;not null;default:now()"`
}
