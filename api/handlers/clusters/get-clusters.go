package clusters

import (
	"encoding/json"
	"net/http"

	"github.com/glueops/autoglue/api/middleware"
	"github.com/glueops/autoglue/internal/db"
	"github.com/glueops/autoglue/internal/db/models"
	"github.com/google/uuid"
)

// GetClusters godoc
// @Summary      List all clusters
// @Description  Returns a list of all clusters in the database
// @Tags         Clusters
// @Param X-Org-ID header string true "Organization context"
// @Security 	 BearerAuth
// @Accept       json
// @Produce      json
// @Success      200  {array}  models.Cluster
// @Failure      500  {string}  string  "Internal Server Error"
// @Router       /api/v1/clusters [get]
func GetClusters(w http.ResponseWriter, r *http.Request) {
	authCtx := middleware.GetAuthContext(r)
	if authCtx == nil || authCtx.OrganizationID == uuid.Nil {
		http.Error(w, "missing org context", http.StatusUnauthorized)
		return
	}

	var clusters []models.Cluster
	if err := db.DB.Where("organization_id = ?", authCtx.OrganizationID).Find(&clusters).Error; err != nil {
		http.Error(w, "failed to list clusters", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(clusters)
}
