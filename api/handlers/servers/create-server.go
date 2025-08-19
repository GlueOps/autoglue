package servers

import (
	"encoding/json"
	"net/http"

	"github.com/glueops/autoglue/internal/db"
	"github.com/glueops/autoglue/internal/db/models"
	"github.com/google/uuid"
)

type ServerInput struct {
	IPAddress string     `json:"ip_address" validate:"required,ip"`
	SSHUser   string     `json:"ssh_user" validate:"required"`
	SshKeyID  uuid.UUID  `json:"ssh_key_id" validate:"required"`
	Role      string     `json:"role" validate:"required,oneof=master worker"`
	ClusterID *uuid.UUID `json:"cluster_id,omitempty"` // Optional
}

// CreateServer creates a new server
// @Summary Create a new server
// @Tags Servers
// @Accept json
// @Produce json
// @Param server body ServerInput true "Server definition"
// @Success 201 {object} models.Server
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/v1/servers [post]
// @Security BearerAuth
func CreateServer(w http.ResponseWriter, r *http.Request) {
	var input ServerInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	server := models.Server{
		IPAddress: input.IPAddress,
		SSHUser:   input.SSHUser,
		SshKeyID:  input.SshKeyID,
		Role:      input.Role,
		Status:    "pending",
	}

	if err := db.DB.Create(&server).Error; err != nil {
		http.Error(w, "Failed to create server", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(server)
}
