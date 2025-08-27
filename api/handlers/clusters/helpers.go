package clusters

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/glueops/autoglue/internal/db"
	"github.com/glueops/autoglue/internal/db/models"
	"github.com/google/uuid"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func validClusterStatus(s string) bool {
	switch strings.ToLower(s) {
	case "provisioning", "ready", "failed", "":
		return true
	default:
		return false
	}
}

func clusterToResp(c models.Cluster, includeNodeGroups bool) clusterResponse {
	resp := clusterResponse{
		ID:        c.ID,
		Name:      c.Name,
		Provider:  c.Provider,
		Region:    c.Region,
		Status:    c.Status,
		CreatedAt: c.CreatedAt.Format(time.RFC3339),
		UpdatedAt: c.UpdatedAt.Format(time.RFC3339),
	}
	if includeNodeGroups {
		resp.NodeGroups = make([]nodeGroupBrief, 0, len(c.NodeGroups))
		for _, ng := range c.NodeGroups {
			resp.NodeGroups = append(resp.NodeGroups, nodeGroupBrief{
				ID:   ng.ID,
				Name: ng.Name,
			})
		}
	}
	return resp
}

func ensureNodeGroupsBelongToOrg(orgID uuid.UUID, ids []uuid.UUID) error {
	if len(ids) == 0 {
		return nil
	}
	var count int64
	if err := db.DB.Model(&models.Server{}).
		Where("organization_id = ? AND id IN ?", orgID, ids).
		Count(&count).Error; err != nil {
		return err
	}
	if count != int64(len(ids)) {
		return errors.New("one or more node group ids are invalid or not in this organization")
	}
	return nil
}

func parseUUIDs(strIDs []string) ([]uuid.UUID, error) {
	out := make([]uuid.UUID, 0, len(strIDs))
	for _, s := range strIDs {
		id, err := uuid.Parse(strings.TrimSpace(s))
		if err != nil {
			return nil, err
		}
		out = append(out, id)
	}
	return out, nil
}
