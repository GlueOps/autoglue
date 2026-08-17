package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// DesiredResource is one entry of one cluster's desired-state snapshot.
//
// Rows are written per generation and never mutated, so a generation is
// immutable once published and the sync endpoint can serve it by plain
// selection. That is what makes a snapshot deliverable as a whole rather than
// as a diff: an agent that was switched off for a week has no history to apply
// a diff against, and a snapshot is correct from any starting state including a
// wiped disk.
//
// Cluster shape and credentials live here. They are config, not tasks — a
// bastion that has lost its kubeconfig should get it back by converging, not by
// someone remembering to enqueue a run.
type DesiredResource struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	OrganizationID uuid.UUID `gorm:"type:uuid;not null;index" json:"organization_id"`

	// The unique key mirrors the agent's own primary key on desired_resources,
	// so re-publishing a generation is an upsert on both sides rather than a
	// duplicate on either.
	ClusterID  uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:uniq_desired_resources_key,priority:1;index:idx_desired_resources_gen,priority:1" json:"cluster_id"`
	Cluster    Cluster   `gorm:"foreignKey:ClusterID;constraint:OnDelete:CASCADE" json:"-"`
	Generation int64     `gorm:"not null;uniqueIndex:uniq_desired_resources_key,priority:2;index:idx_desired_resources_gen,priority:2" json:"generation"`

	ResourceType string `gorm:"type:text;not null;uniqueIndex:uniq_desired_resources_key,priority:3" json:"resource_type"`
	ResourceID   string `gorm:"type:text;not null;uniqueIndex:uniq_desired_resources_key,priority:4" json:"resource_id"`
	ResourceName string `gorm:"type:text;not null;default:''" json:"resource_name"`

	// Phase orders application within a generation; DependsOn is a JSON array of
	// resource ids for ordering inside a phase. Both are advisory data the agent
	// interprets, which is what lets a new resource type ship without the
	// control plane learning how to sequence it.
	Phase     int            `gorm:"not null;default:0" json:"phase"`
	DependsOn datatypes.JSON `gorm:"type:jsonb;not null;default:'[]'" json:"depends_on"`

	Spec datatypes.JSON `gorm:"type:jsonb;not null;default:'{}'" json:"spec"`

	// SpecHash lets the agent skip work whose spec has not moved between two
	// generations. Computed control-plane side so both ends agree on what was
	// hashed rather than on how it happened to be serialised.
	SpecHash string `gorm:"type:text;not null;default:''" json:"spec_hash"`

	// Required distinguishes "this generation cannot be called applied without
	// it" from best-effort. An agent that silently skips a required resource
	// and still reports convergence is the failure this flag exists to make
	// impossible.
	Required bool `gorm:"not null;default:true" json:"required"`

	CreatedAt time.Time `gorm:"type:timestamptz;not null;default:now()" json:"created_at"`
}
