package mapper

import (
	"testing"
	"time"

	"github.com/glueops/autoglue/internal/common"
	"github.com/glueops/autoglue/internal/models"
	"github.com/google/uuid"
)

// The cluster response nests labels, annotations, taints and servers. Those
// nested copies used to carry only key/value, so every id and organization_id
// came back as the zero UUID — a consumer reading a cluster could see a label
// but had no id to delete or patch it with, and org-owned rows looked unowned.
func TestNodePoolToDTO_CarriesIdentityIntoNestedItems(t *testing.T) {
	orgID := uuid.New()
	now := time.Now().UTC().Truncate(time.Second)

	labelID, annID, taintID, serverID := uuid.New(), uuid.New(), uuid.New(), uuid.New()

	np := models.NodePool{
		AuditFields: common.AuditFields{
			ID: uuid.New(), OrganizationID: orgID, CreatedAt: now, UpdatedAt: now,
		},
		Name: "worker-pool",
		Role: "worker",
		Labels: []models.Label{{
			AuditFields: common.AuditFields{
				ID: labelID, OrganizationID: orgID, CreatedAt: now, UpdatedAt: now,
			},
			Key: "glueops.dev/role", Value: "glueops-platform",
		}},
		Annotations: []models.Annotation{{
			AuditFields: common.AuditFields{
				ID: annID, OrganizationID: orgID, CreatedAt: now, UpdatedAt: now,
			},
			Key: "note", Value: "hello",
		}},
		// models.Taint carries its own fields rather than embedding AuditFields.
		Taints: []models.Taint{{
			ID: taintID, OrganizationID: orgID, CreatedAt: now, UpdatedAt: now,
			Key: "dedicated", Value: "platform", Effect: "NoSchedule",
		}},
		Servers: []models.Server{{
			ID: serverID, OrganizationID: orgID, Hostname: "worker-0",
			PrivateIPAddress: "10.0.0.1", SSHUser: "cluster", Role: "worker",
			Status: "pending", CreatedAt: now, UpdatedAt: now,
		}},
	}

	got := NodePoolToDTO(np)

	if len(got.Labels) != 1 || got.Labels[0].ID != labelID {
		t.Errorf("label id = %v, want %v", got.Labels[0].ID, labelID)
	}
	if got.Labels[0].OrganizationID != orgID {
		t.Errorf("label org = %v, want %v", got.Labels[0].OrganizationID, orgID)
	}

	if len(got.Annotations) != 1 || got.Annotations[0].ID != annID {
		t.Errorf("annotation id = %v, want %v", got.Annotations[0].ID, annID)
	}
	if got.Annotations[0].OrganizationID != orgID {
		t.Errorf("annotation org = %v, want %v", got.Annotations[0].OrganizationID, orgID)
	}

	if len(got.Taints) != 1 || got.Taints[0].ID != taintID {
		t.Errorf("taint id = %v, want %v", got.Taints[0].ID, taintID)
	}
	if got.Taints[0].OrganizationID != orgID {
		t.Errorf("taint org = %v, want %v", got.Taints[0].OrganizationID, orgID)
	}
	// Taint timestamps are strings on the DTO, unlike the embedded AuditFields
	// used by labels and annotations.
	if got.Taints[0].CreatedAt == "" {
		t.Error("taint created_at is empty")
	}

	if len(got.Servers) != 1 || got.Servers[0].ID != serverID {
		t.Errorf("server id = %v, want %v", got.Servers[0].ID, serverID)
	}
	if got.Servers[0].OrganizationID != orgID {
		t.Errorf("server org = %v, want %v", got.Servers[0].OrganizationID, orgID)
	}
}

func TestNodePoolToDTO_EmptyCollectionsSerializeAsLists(t *testing.T) {
	// These must be non-nil so they marshal as [] rather than null; the SPA
	// maps over them directly.
	got := NodePoolToDTO(models.NodePool{Name: "empty", Role: "worker"})

	if got.Labels == nil || got.Annotations == nil || got.Taints == nil || got.Servers == nil {
		t.Fatalf("nested collections must be empty slices, not nil: %+v", got)
	}
	if len(got.Labels) != 0 || len(got.Servers) != 0 {
		t.Error("expected empty collections for an empty node pool")
	}
}
