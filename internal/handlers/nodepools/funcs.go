package nodepools

import (
	"errors"
	"fmt"
	"strings"

	"github.com/glueops/autoglue/internal/db"
	"github.com/glueops/autoglue/internal/db/models"
	"github.com/google/uuid"
)

func toResp(ng models.NodePool, includeServers bool) nodePoolResponse {
	resp := nodePoolResponse{
		ID:   ng.ID,
		Name: ng.Name,
	}
	if includeServers {
		resp.Servers = make([]serverBrief, 0, len(ng.Servers))
		for _, s := range ng.Servers {
			resp.Servers = append(resp.Servers, serverBrief{
				ID:       s.ID,
				Hostname: s.Hostname,
				IP:       s.IPAddress,
				Role:     s.Role,
				Status:   s.Status,
			})
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

func ensureServersBelongToOrg(orgID uuid.UUID, ids []uuid.UUID) error {
	var count int64
	if err := db.DB.Model(&models.Server{}).
		Where("organization_id = ? AND id IN ?", orgID, ids).
		Count(&count).Error; err != nil {
		return err
	}
	if count != int64(len(ids)) {
		return fmt.Errorf("some servers do not belong to this organization")
	}
	return nil
}

func ensureTaintsBelongToOrg(orgID uuid.UUID, ids []uuid.UUID) error {
	if len(ids) == 0 {
		return nil
	}
	var count int64
	if err := db.DB.Model(&models.Taint{}).
		Where("organization_id = ? AND id IN ?", orgID, ids).
		Count(&count).Error; err != nil {
		return err
	}
	if count != int64(len(ids)) {
		return errors.New("some taints not in organization")
	}
	return nil
}
