package clusters

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/glueops/autoglue/api/middleware"
	"github.com/glueops/autoglue/internal/db"
	"github.com/glueops/autoglue/internal/db/models"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type ClusterInput struct {
	Name     string `json:"name" validate:"required,min=3"`
	Provider string `json:"provider" validate:"required,oneof=aws hetzner linode digitalocean"`
	Region   string `json:"region" validate:"required"`
}

// CreateCluster handles POST /clusters
// @Summary Create a new cluster
// @Tags Clusters
// @Accept json
// @Produce json
// @Param cluster body ClusterInput true "Cluster definition"
// @Success 201 {object} models.Cluster
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/v1/clusters [post]
// @Param X-Org-ID header string true "Organization context"
// @Security BearerAuth
func CreateCluster(w http.ResponseWriter, r *http.Request) {
	authCtx := middleware.GetAuthContext(r)
	if authCtx == nil || authCtx.OrganizationID == uuid.Nil {
		http.Error(w, "missing org context", http.StatusUnauthorized)
		return
	}

	var input ClusterInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	cluster := models.Cluster{
		Name:           input.Name,
		Provider:       input.Provider,
		Region:         input.Region,
		Status:         "provisioning",
		Kubeconfig:     "",
		OrganizationID: authCtx.OrganizationID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := db.DB.Create(&cluster).Error; err != nil {
		http.Error(w, "Failed to create cluster", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cluster)
}
