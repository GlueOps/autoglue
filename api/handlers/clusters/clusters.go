package clusters

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/glueops/autoglue/api/middleware"
	"github.com/glueops/autoglue/internal/db"
	"github.com/glueops/autoglue/internal/db/models"
	"github.com/glueops/autoglue/internal/utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// ListClusters godoc
// @Summary      List clusters (org scoped)
// @Description  List clusters for the organization in X-Org-ID. Use `provider`, `region`, `status`, and `q` (name contains). Add `include=node_groups` to include attached servers.
// @Tags         clusters
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        provider query string false "Filter by provider"
// @Param        region query string false "Filter by region"
// @Param        status query string false "Filter by status (provisioning|ready|failed)"
// @Param        q query string false "Name contains (case-insensitive)"
// @Param        include query string false "Optional: servers"
// @Security     BearerAuth
// @Success      200 {array}  clusterResponse
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      500 {string} string "failed to list clusters"
// @Router       /api/v1/clusters [get]
func ListClusters(w http.ResponseWriter, r *http.Request) {
	ac := middleware.GetAuthContext(r)
	if ac == nil || ac.OrganizationID == uuid.Nil {
		http.Error(w, "organization required", http.StatusForbidden)
		return
	}

	q := db.DB.Where("organization_id = ?", ac.OrganizationID)
	if p := strings.TrimSpace(r.URL.Query().Get("provider")); p != "" {
		q = q.Where("provider = ?", p)
	}
	if region := strings.TrimSpace(r.URL.Query().Get("region")); region != "" {
		q = q.Where("region = ?", region)
	}
	if s := strings.TrimSpace(r.URL.Query().Get("status")); s != "" {
		if !validClusterStatus(s) {
			http.Error(w, "invalid status", http.StatusBadRequest)
			return
		}
		q = q.Where("status = ?", strings.ToLower(s))
	}
	if needle := strings.TrimSpace(r.URL.Query().Get("q")); needle != "" {
		q = q.Where("name ILIKE ?", "%"+needle+"%")
	}

	includeNodeGroups := strings.EqualFold(strings.TrimSpace(r.URL.Query().Get("include")), "node_groups")
	if includeNodeGroups {
		q = q.Preload("NodeGroups")
	}

	var rows []models.Cluster
	if err := q.Order("created_at DESC").Find(&rows).Error; err != nil {
		http.Error(w, "failed to list clusters", http.StatusInternalServerError)
		return
	}

	out := make([]clusterResponse, 0, len(rows))
	for _, c := range rows {
		out = append(out, clusterToResp(c, includeNodeGroups))
	}
	writeJSON(w, http.StatusOK, out)
}

// GetCluster godoc
// @Summary      Get cluster by ID (org scoped)
// @Description  Returns one cluster. Add `include=servers` to include servers. Add `reveal_kubeconfig=true` to include decrypted kubeconfig.
// @Tags         clusters
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        id path string true "Cluster ID (UUID)"
// @Param        include query string false "Optional: servers"
// @Param        reveal_kubeconfig query bool false "Reveal decrypted kubeconfig"
// @Security     BearerAuth
// @Success      200 {object} clusterResponse
// @Failure      400 {string} string "invalid id"
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      404 {string} string "not found"
// @Failure      500 {string} string "fetch/decrypt failed"
// @Router       /api/v1/clusters/{id} [get]
func GetCluster(w http.ResponseWriter, r *http.Request) {
	ac := middleware.GetAuthContext(r)
	if ac == nil || ac.OrganizationID == uuid.Nil {
		http.Error(w, "organization required", http.StatusForbidden)
		return
	}

	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	includeNodeGroups := strings.EqualFold(strings.TrimSpace(r.URL.Query().Get("include")), "node_groups")
	reveal := strings.EqualFold(strings.TrimSpace(r.URL.Query().Get("reveal_kubeconfig")), "true")

	var c models.Cluster
	q := db.DB.Where("id = ? AND organization_id = ?", id, ac.OrganizationID)
	if includeNodeGroups {
		q = q.Preload("NodeGroups")
	}
	if err := q.First(&c).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "fetch failed", http.StatusInternalServerError)
		return
	}

	resp := clusterToResp(c, includeNodeGroups)
	if reveal && c.EncryptedKubeconfig != "" {
		plain, err := utils.DecryptForOrg(ac.OrganizationID, c.EncryptedKubeconfig, c.KubeIV, c.KubeTag)
		if err != nil {
			http.Error(w, "decrypt failed", http.StatusInternalServerError)
			return
		}
		resp.Kubeconfig = plain
	}
	writeJSON(w, http.StatusOK, resp)
}

// CreateCluster godoc
// @Summary      Create cluster (org scoped)
// @Description  Creates a cluster and optionally attaches initial servers. If kubeconfig is provided, it is encrypted at rest.
// @Tags         clusters
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        body body createClusterRequest true "Cluster payload"
// @Security     BearerAuth
// @Success      201 {object} clusterResponse
// @Failure      400 {string} string "invalid json / invalid status / invalid server_ids"
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      500 {string} string "create failed"
// @Router       /api/v1/clusters [post]
func CreateCluster(w http.ResponseWriter, r *http.Request) {
	ac := middleware.GetAuthContext(r)
	if ac == nil || ac.OrganizationID == uuid.Nil {
		http.Error(w, "organization required", http.StatusForbidden)
		return
	}

	var req createClusterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if req.Name == "" || req.Provider == "" || req.Region == "" {
		http.Error(w, "name, provider, region are required", http.StatusBadRequest)
		return
	}
	if req.Status != "" && !validClusterStatus(req.Status) {
		http.Error(w, "invalid status", http.StatusBadRequest)
		return
	}

	c := models.Cluster{
		OrganizationID: ac.OrganizationID,
		Name:           req.Name,
		Provider:       req.Provider,
		Region:         req.Region,
		Status:         "provisioning",
	}
	if req.Status != "" {
		c.Status = strings.ToLower(req.Status)
	}

	// Encrypt kubeconfig if provided
	if strings.TrimSpace(req.Kubeconfig) != "" {
		cipher, iv, tag, err := utils.EncryptForOrg(ac.OrganizationID, []byte(req.Kubeconfig))
		if err != nil {
			http.Error(w, "kubeconfig encryption failed", http.StatusInternalServerError)
			return
		}
		c.EncryptedKubeconfig = cipher
		c.KubeIV = iv
		c.KubeTag = tag
	}

	// Create cluster
	if err := db.DB.Create(&c).Error; err != nil {
		http.Error(w, "create failed", http.StatusInternalServerError)
		return
	}

	// Attach servers if provided (unchanged)
	if len(req.NodeGroupIDs) > 0 {
		ids, err := parseUUIDs(req.NodeGroupIDs)
		if err != nil {
			http.Error(w, "invalid node_group_ids", http.StatusBadRequest)
			return
		}
		var groups []models.NodeGroup
		if err := db.DB.Where("id IN ? AND organization_id = ?", ids, ac.OrganizationID).
			Find(&groups).Error; err != nil {
			http.Error(w, "attach node groups failed", http.StatusInternalServerError)
			return
		}
		if err := db.DB.Model(&c).Association("NodeGroups").Append(&groups); err != nil {
			http.Error(w, "attach node groups failed", http.StatusInternalServerError)
			return
		}
	}

	writeJSON(w, http.StatusCreated, clusterToResp(c, false))
}

// UpdateCluster godoc
// @Summary      Update cluster (org scoped)
// @Description  Partially update cluster fields. If kubeconfig is provided, it is encrypted at rest. Provide empty string to clear stored kubeconfig.
// @Tags         clusters
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        id path string true "Cluster ID (UUID)"
// @Param        body body updateClusterRequest true "Fields to update"
// @Security     BearerAuth
// @Success      200 {object} clusterResponse
// @Failure      400 {string} string "invalid id / invalid json / invalid status"
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      404 {string} string "not found"
// @Failure      500 {string} string "update failed"
// @Router       /api/v1/clusters/{id} [patch]
func UpdateCluster(w http.ResponseWriter, r *http.Request) {
	ac := middleware.GetAuthContext(r)
	if ac == nil || ac.OrganizationID == uuid.Nil {
		http.Error(w, "organization required", http.StatusForbidden)
		return
	}
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var c models.Cluster
	if err := db.DB.Where("id = ? AND organization_id = ?", id, ac.OrganizationID).First(&c).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "fetch failed", http.StatusInternalServerError)
		return
	}

	var req updateClusterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if req.Name != nil {
		c.Name = *req.Name
	}
	if req.Provider != nil {
		c.Provider = *req.Provider
	}
	if req.Region != nil {
		c.Region = *req.Region
	}
	if req.Status != nil {
		if !validClusterStatus(*req.Status) {
			http.Error(w, "invalid status", http.StatusBadRequest)
			return
		}
		c.Status = strings.ToLower(*req.Status)
	}
	if req.Kubeconfig != nil {
		k := strings.TrimSpace(*req.Kubeconfig)
		if k == "" {
			// Clear any stored kubeconfig
			c.EncryptedKubeconfig = ""
			c.KubeIV = ""
			c.KubeTag = ""
		} else {
			cipher, iv, tag, err := utils.EncryptForOrg(ac.OrganizationID, []byte(k))
			if err != nil {
				http.Error(w, "kubeconfig encryption failed", http.StatusInternalServerError)
				return
			}
			c.EncryptedKubeconfig = cipher
			c.KubeIV = iv
			c.KubeTag = tag
		}
	}

	if err := db.DB.Save(&c).Error; err != nil {
		http.Error(w, "update failed", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, clusterToResp(c, false))
}

// DeleteCluster godoc
// @Summary      Delete cluster (org scoped)
// @Description  Permanently deletes the cluster (associated links are removed; servers remain).
// @Tags         clusters
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        id path string true "Cluster ID (UUID)"
// @Security     BearerAuth
// @Success      204 {string} string "No Content"
// @Failure      400 {string} string "invalid id"
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      500 {string} string "delete failed"
// @Router       /api/v1/clusters/{id} [delete]
func DeleteCluster(w http.ResponseWriter, r *http.Request) {
	ac := middleware.GetAuthContext(r)
	if ac == nil || ac.OrganizationID == uuid.Nil {
		http.Error(w, "organization required", http.StatusForbidden)
		return
	}
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := db.DB.Where("id = ? AND organization_id = ?", id, ac.OrganizationID).
		Delete(&models.Cluster{}).Error; err != nil {
		http.Error(w, "delete failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ListClusterNodeGroups godoc
// @Summary      List node groups attached to a cluster (org scoped)
// @Tags         clusters
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        id path string true "Cluster ID (UUID)"
// @Security     BearerAuth
// @Success      200 {array}  nodeGroupBrief
// @Failure      400,401,403,404,500 {string} string
// @Router       /api/v1/clusters/{id}/node-groups [get]
func ListClusterNodeGroups(w http.ResponseWriter, r *http.Request) {
	ac := middleware.GetAuthContext(r)
	if ac == nil || ac.OrganizationID == uuid.Nil {
		http.Error(w, "organization required", http.StatusForbidden)
		return
	}
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	var c models.Cluster
	if err := db.DB.Where("id = ? AND organization_id = ?", id, ac.OrganizationID).
		Preload("NodeGroups").First(&c).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "fetch failed", http.StatusInternalServerError)
		return
	}
	out := make([]nodeGroupBrief, 0, len(c.NodeGroups))
	for _, ng := range c.NodeGroups {
		out = append(out, nodeGroupBrief{ID: ng.ID, Name: ng.Name})
	}
	writeJSON(w, http.StatusOK, out)
}

// AttachClusterNodeGroups godoc
// @Summary      Attach node groups to a cluster (org scoped)
// @Tags         clusters
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        id path string true "Cluster ID (UUID)"
// @Param        body body nodeGroupIds true "Node Group IDs to attach"
// @Security     BearerAuth
// @Success      204 {string} string "No Content"
// @Failure      400,401,403,404,500 {string} string
// @Router       /api/v1/clusters/{id}/node-groups [post]
func AttachClusterNodeGroups(w http.ResponseWriter, r *http.Request) {
	ac := middleware.GetAuthContext(r)
	if ac == nil || ac.OrganizationID == uuid.Nil {
		http.Error(w, "organization required", http.StatusForbidden)
		return
	}

	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var c models.Cluster
	if err := db.DB.Where("id = ? AND organization_id = ?", id, ac.OrganizationID).First(&c).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "fetch failed", http.StatusInternalServerError)
		return
	}

	var body struct {
		NodeGroupIDs []string `json:"node_group_ids"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || len(body.NodeGroupIDs) == 0 {
		http.Error(w, "invalid node_group_ids", http.StatusBadRequest)
		return
	}

	ids, err := parseUUIDs(body.NodeGroupIDs)
	if err != nil || ensureNodeGroupsBelongToOrg(ac.OrganizationID, ids) != nil {
		http.Error(w, "invalid node_group_ids", http.StatusBadRequest)
		return
	}

	var groups []models.NodeGroup
	if err := db.DB.Where("id IN ? AND organization_id = ?", ids, ac.OrganizationID).
		Find(&groups).Error; err != nil {
		http.Error(w, "attach failed", http.StatusInternalServerError)
		return
	}

	if err := db.DB.Model(&c).Association("NodeGroups").Append(&groups); err != nil {
		http.Error(w, "attach failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// DetachClusterNodeGroup godoc
// @Summary      Detach one node group from a cluster (org scoped)
// @Tags         clusters
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        id path string true "Cluster ID (UUID)"
// @Param        nodeGroupId path string true "Node Group ID (UUID)"
// @Security     BearerAuth
// @Success      204 {string} string "No Content"
// @Failure      400,401,403,404,500 {string} string
// @Router       /api/v1/clusters/{id}/node-groups/{nodeGroupId} [delete]
func DetachClusterNodeGroup(w http.ResponseWriter, r *http.Request) {
	ac := middleware.GetAuthContext(r)
	if ac == nil || ac.OrganizationID == uuid.Nil {
		http.Error(w, "organization required", http.StatusForbidden)
		return
	}
	vars := mux.Vars(r)
	cid, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	ngid, err := uuid.Parse(vars["nodeGroupId"])
	if err != nil {
		http.Error(w, "invalid nodeGroupId", http.StatusBadRequest)
		return
	}

	var c models.Cluster
	if err := db.DB.Where("id = ? AND organization_id = ?", cid, ac.OrganizationID).
		Preload("NodeGroups").First(&c).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "fetch failed", http.StatusInternalServerError)
		return
	}

	var ng models.NodeGroup
	if err := db.DB.Where("id = ? AND organization_id = ?", ngid, ac.OrganizationID).First(&ng).Error; err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	if err := db.DB.Model(&c).Association("NodeGroups").Delete(&ng); err != nil {
		http.Error(w, "detach failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
