package servers

import (
	"net/http"

	"github.com/glueops/autoglue/internal/db"
	"github.com/glueops/autoglue/internal/db/models"
	"github.com/gorilla/mux"
)

// DeleteServer deletes a server
// @Summary Delete a server
// @Tags Servers
// @Produce plain
// @Param id path string true "Server ID"
// @Success 204 {string} string "Deleted"
// @Router /api/v1/servers/{id} [delete]
// @Security BearerAuth
func DeleteServer(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := db.DB.Delete(&models.Server{}, "id = ?", id).Error; err != nil {
		http.Error(w, "Failed to delete", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
