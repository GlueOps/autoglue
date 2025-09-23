package taints

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/glueops/autoglue/internal/db"
	"github.com/glueops/autoglue/internal/db/models"
	"github.com/google/uuid"
)

var allowedEffects = map[string]struct{}{
	"NoSchedule":       {},
	"PreferNoSchedule": {},
	"NoExecute":        {},
}

// includeNodePools returns true when the query param requests linked pools.
// Accepts both "node_pools" and "node_groups" for compatibility.
func includeNodePools(r *http.Request) bool {
	inc := strings.TrimSpace(r.URL.Query().Get("include"))
	return strings.EqualFold(inc, "node_pools") || strings.EqualFold(inc, "node_groups")
}

func toResp(t models.Taint, include bool) taintResponse {
	resp := taintResponse{
		ID:     t.ID,
		Key:    t.Key,
		Value:  t.Value,
		Effect: t.Effect,
	}
	if include {
		resp.NodeGroups = make([]nodePoolBrief, 0, len(t.NodePools))
		for _, np := range t.NodePools {
			resp.NodeGroups = append(resp.NodeGroups, nodePoolBrief{ID: np.ID, Name: np.Name})
		}
	}
	return resp
}

func parseUUIDs(ids []string) ([]uuid.UUID, error) {
	out := make([]uuid.UUID, 0, len(ids))
	for _, s := range ids {
		u, err := uuid.Parse(strings.TrimSpace(s))
		if err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, nil
}

func ensureNodePoolsBelongToOrg(orgID uuid.UUID, ids []uuid.UUID) error {
	var count int64
	if err := db.DB.Model(&models.NodePool{}).
		Where("organization_id = ? AND id IN ?", orgID, ids).
		Count(&count).Error; err != nil {
		return err
	}
	if count != int64(len(ids)) {
		return fmt.Errorf("some node groups do not belong to this organization")
	}
	return nil
}
