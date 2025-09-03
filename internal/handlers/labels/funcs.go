package labels

import (
	"fmt"
	"strings"

	"github.com/glueops/autoglue/internal/db"
	"github.com/glueops/autoglue/internal/db/models"
	"github.com/google/uuid"
)

func toResp(l models.Label, include bool) labelResponse {
	resp := labelResponse{
		ID:    l.ID,
		Key:   l.Key,
		Value: l.Value,
	}
	if include {
		resp.NodeGroups = make([]nodePoolBrief, 0, len(l.NodePools))
		for _, np := range l.NodePools {
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
