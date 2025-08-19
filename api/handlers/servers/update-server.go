package servers

import (
	"encoding/json"
	"net/http"

	"github.com/glueops/autoglue/internal/db"
	"github.com/glueops/autoglue/internal/db/models"
	"github.com/gorilla/mux"
)

// UpdateServer updates a server by ID
// @Summary Update a server
// @Tags Servers
// @Accept json
// @Produce json
// @Param id path string true "Server ID"
// @Param server body ServerInput true "Server fields"
// @Success 200 {object} models.Server
// @Failure 404 {string} string "Not Found"
// @Router /api/v1/servers/{id} [patch]
// @Security BearerAuth
func UpdateServer(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var server models.Server
	if err := db.DB.First(&server, "id = ?", id).Error; err != nil {
		http.Error(w, "Server not found", http.StatusNotFound)
		return
	}

	var input ServerInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	server.IPAddress = input.IPAddress
	server.SSHUser = input.SSHUser
	server.SshKeyID = input.SshKeyID
	server.Role = input.Role

	db.DB.Save(&server)
	json.NewEncoder(w).Encode(server)
}
