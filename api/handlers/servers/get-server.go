package servers

import (
	"encoding/json"
	"net/http"

	"github.com/glueops/autoglue/internal/db"
	"github.com/glueops/autoglue/internal/db/models"
	"github.com/gorilla/mux"
)

// GetServer fetches a single server by ID
// @Summary Get server by ID
// @Tags Servers
// @Produce json
// @Param id path string true "Server ID"
// @Success 200 {object} models.Server
// @Failure 404 {string} string "Not Found"
// @Router /api/v1/servers/{id} [get]
// @Security BearerAuth
func GetServer(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var server models.Server
	if err := db.DB.First(&server, "id = ?", id).Error; err != nil {
		http.Error(w, "Server not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(server)
}
