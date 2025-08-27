package servers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/glueops/autoglue/internal/db"
	"github.com/glueops/autoglue/internal/db/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func validStatus(s string) bool {
	switch strings.ToLower(s) {
	case "pending", "provisioning", "ready", "failed", "":
		return true
	default:
		return false
	}
}

// ensureKeyBelongsToOrg loads the key and ensures it’s in the same org.
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

// toResponse maps a models.Server to serverResponse.
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
