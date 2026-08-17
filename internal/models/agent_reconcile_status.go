package models

import (
	"time"

	"github.com/google/uuid"
)

// Reconcile states for a single resource, as the agent reports them.
const (
	AgentReconcileStatePending = "pending"
	AgentReconcileStateApplied = "applied"
	AgentReconcileStateFailed  = "failed"
)

// AgentReconcileStatus is what one agent last said about one resource.
//
// Current status only, not history: there is one truth per resource and it
// carries the generation it was achieved at, which is what makes a stale
// DesiredGeneration readable as "this resource has not caught up" without
// diffing anything.
//
// The control plane never writes these on the agent's behalf. A row that says
// applied means an agent claimed it, and the gap between Agent.AppliedGeneration
// and ClusterDesiredState.Generation is the only convergence signal an operator
// should trust.
type AgentReconcileStatus struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	OrganizationID uuid.UUID `gorm:"type:uuid;not null;index" json:"organization_id"`
	ClusterID      uuid.UUID `gorm:"type:uuid;not null;index" json:"cluster_id"`

	AgentID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:uniq_agent_reconcile_status,priority:1" json:"agent_id"`
	Agent   Agent     `gorm:"foreignKey:AgentID;constraint:OnDelete:CASCADE" json:"-"`

	ResourceType string `gorm:"type:text;not null;uniqueIndex:uniq_agent_reconcile_status,priority:2" json:"resource_type"`
	ResourceID   string `gorm:"type:text;not null;uniqueIndex:uniq_agent_reconcile_status,priority:3" json:"resource_id"`

	DesiredGeneration int64  `gorm:"not null;default:0" json:"desired_generation"`
	State             string `gorm:"type:text;not null;default:'pending'" json:"state"`
	LastError         string `gorm:"type:text;not null;default:''" json:"last_error"`

	ReportedAt time.Time `gorm:"type:timestamptz;not null;default:now()" json:"reported_at"`

	CreatedAt time.Time `json:"created_at,omitempty" gorm:"type:timestamptz;column:created_at;not null;default:now()"`
	UpdatedAt time.Time `json:"updated_at,omitempty" gorm:"type:timestamptz;autoUpdateTime;column:updated_at;not null;default:now()"`
}
