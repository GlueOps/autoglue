package nodepools

import (
	"errors"

	"github.com/glueops/autoglue/internal/db"
	"github.com/glueops/autoglue/internal/db/models"
	"github.com/google/uuid"
)

func toResp(ng models.NodePool, includeServers bool) nodePoolResponse {
	out := nodePoolResponse{
		ID:   ng.ID,
		Name: ng.Name,
	}
	if includeServers {
		out.Servers = make([]serverBrief, 0, len(ng.Servers))
		for _, s := range ng.Servers {
			out.Servers = append(out.Servers, serverBrief{
				ID:       s.ID,
				Hostname: s.Hostname,
				IP:       s.IPAddress,
				Role:     s.Role,
				Status:   s.Status,
			})
		}
	}
	return out
}

func parseUUIDs(ids []string) ([]uuid.UUID, error) {
	out := make([]uuid.UUID, 0, len(ids))
	for _, raw := range ids {
		id, err := uuid.Parse(raw)
		if err != nil {
			return nil, err
		}
		out = append(out, id)
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
		return errors.New("some servers do not belong to org")
	}
	return nil
}

func ensureLabelsBelongToOrg(orgID uuid.UUID, ids []uuid.UUID) error {
	var count int64
	if err := db.DB.Model(&models.Label{}).
		Where("organization_id = ? AND id IN ?", orgID, ids).
		Count(&count).Error; err != nil {
		return err
	}
	if count != int64(len(ids)) {
		return errors.New("some labels do not belong to org")
	}
	return nil
}

func ensureTaintsBelongToOrg(orgID uuid.UUID, ids []uuid.UUID) error {
	var count int64
	if err := db.DB.Model(&models.Taint{}).
		Where("organization_id = ? AND id IN ?", orgID, ids).
		Count(&count).Error; err != nil {
		return err
	}
	if count != int64(len(ids)) {
		return errors.New("some taints do not belong to org")
	}
	return nil
}
