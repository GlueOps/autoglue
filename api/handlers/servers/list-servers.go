package servers

import (
	"encoding/json"
	"net/http"

	"github.com/glueops/autoglue/internal/db"
	"github.com/glueops/autoglue/internal/db/models"
)

// ListServers lists all servers
// @Summary List all servers
// @Tags Servers
// @Produce json
// @Success 200 {array} models.Server
// @Router /api/v1/servers [get]
// @Security BearerAuth
func ListServers(w http.ResponseWriter, r *http.Request) {
	var servers []models.Server
	if err := db.DB.Find(&servers).Error; err != nil {
		http.Error(w, "Failed to fetch servers", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(servers)
}
