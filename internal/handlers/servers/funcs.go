package servers

import (
	"errors"
	"strings"
	"time"

	"github.com/glueops/autoglue/internal/db"
	"github.com/glueops/autoglue/internal/db/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func toResponse(s models.Server) serverResponse {
	return serverResponse{
		ID:             s.ID,
		OrganizationID: s.OrganizationID,
		Hostname:       s.Hostname,
		IPAddress:      s.IPAddress,
		SSHUser:        s.SSHUser,
		SshKeyID:       s.SshKeyID,
		Role:           s.Role,
		Status:         s.Status,
		CreatedAt:      s.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:      s.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func validStatus(s string) bool {
	switch strings.ToLower(s) {
	case "pending", "provisioning", "ready", "failed", "":
		return true
	default:
		return false
	}
}

func ensureKeyBelongsToOrg(orgID uuid.UUID, keyID uuid.UUID) error {
	var k models.SshKey
	if err := db.DB.Where("id = ? AND organization_id = ?", keyID, orgID).First(&k).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("ssh key not found for this organization")
		}
		return err
	}
	return nil
}
