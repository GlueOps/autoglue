package clusters

import (
	"net/http"
	"strings"

	"github.com/glueops/autoglue/internal/db"
	"github.com/glueops/autoglue/internal/db/models"
	"github.com/glueops/autoglue/internal/middleware"
	"github.com/glueops/autoglue/internal/response"
	"github.com/glueops/autoglue/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ListClusters godoc
// @Summary      List clusters (org scoped)
// @Description  Returns clusters for the organization in X-Org-ID. Add `include=node_pools,bastion` to expand. Filter by `q` (name contains).
// @Tags         clusters
// @Security     BearerAuth
// @Produce      json
// @Param        X-Org-ID  header  string  true  "Organization UUID"
// @Param        q         query   string  false "Name contains (case-insensitive)"
// @Param        include   query   string  false "Optional: node_pools,bastion"
// @Success      200       {array} clusterResponse
// @Failure      401       {string} string "Unauthorized"
// @Failure      403       {string} string "organization required"
// @Failure      500       {string} string "failed to list clusters"
// @Router       /api/v1/clusters [get]
func ListClusters(w http.ResponseWriter, r *http.Request) {
	ac := middleware.GetAuthContext(r)
	if ac == nil || ac.OrganizationID == uuid.Nil {
		http.Error(w, "organization required", http.StatusForbidden)
		return
	}
	include := strings.Split(strings.ToLower(r.URL.Query().Get("include")), ",")
	withPools := contains(include, "node_pools")
	withBastion := contains(include, "bastion")
	q := strings.TrimSpace(r.URL.Query().Get("q"))

	var rows []models.Cluster
	tx := db.DB.Where("organization_id = ?", ac.OrganizationID)
	if q != "" {
		tx = tx.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(q)+"%")
	}
	if withPools {
		tx = tx.
			Preload("NodePools").
			Preload("NodePools.Labels").
			Preload("NodePools.Annotations").
			Preload("NodePools.Taints").
			Preload("NodePools.Servers")
	}
	if withBastion {
		tx = tx.Preload("BastionServer")
	}
	if err := tx.Find(&rows).Error; err != nil {
		http.Error(w, "failed to list clusters", http.StatusInternalServerError)
		return
	}

	out := make([]clusterResponse, 0, len(rows))
	for _, c := range rows {
		out = append(out, toResp(c, withPools, withBastion))
	}
	_ = response.JSON(w, http.StatusOK, out)
}

// CreateCluster godoc
// @Summary      Create cluster (org scoped)
// @Description  Creates a cluster and optionally links node pools and a bastion server. If `kubeconfig` is provided, it will be encrypted per-organization and stored securely (never returned).
// @Tags         clusters
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        X-Org-ID  header  string                           true  "Organization UUID"
// @Param        body      body    clusters.createClusterRequest     true  "payload"
// @Success      201       {object} clusters.clusterResponse
// @Failure      400       {string} string "invalid json / invalid node_pool_ids / invalid bastion_server_id"
// @Failure      401       {string} string "Unauthorized"
// @Failure      403       {string} string "organization required"
// @Failure      500       {string} string "create failed"
// @Router       /api/v1/clusters [post]
func CreateCluster(w http.ResponseWriter, r *http.Request) {
	ac := middleware.GetAuthContext(r)
	if ac == nil || ac.OrganizationID == uuid.Nil {
		http.Error(w, "organization required", http.StatusForbidden)
		return
	}
	var in createClusterRequest
	if !readJSON(w, r, &in) {
		return
	}

	var poolIDs []uuid.UUID
	var err error
	if len(in.NodePoolIDs) > 0 {
		poolIDs, err = parseUUIDs(in.NodePoolIDs)
		if err != nil {
			http.Error(w, "invalid node_pool_ids", http.StatusBadRequest)
			return
		}
		if err := ensureNodePoolsBelongToOrg(ac.OrganizationID, poolIDs); err != nil {
			http.Error(w, "invalid node_pool_ids", http.StatusBadRequest)
			return
		}
	}

	var bastionID *uuid.UUID
	if in.BastionServerID != nil && *in.BastionServerID != "" {
		bid, err := uuid.Parse(*in.BastionServerID)
		if err != nil {
			http.Error(w, "invalid bastion_server_id", http.StatusBadRequest)
			return
		}
		if err := ensureServerBelongsToOrgWithRole(ac.OrganizationID, bid, "bastion"); err != nil {
			http.Error(w, "invalid bastion_server_id", http.StatusBadRequest)
			return
		}
		bastionID = &bid
	}

	c := models.Cluster{
		OrganizationID:  ac.OrganizationID,
		Name:            in.Name,
		Provider:        in.Provider,
		Region:          in.Region,
		Status:          "pending",
		BastionServerID: bastionID,
	}

	if in.ClusterLoadBalancer != nil {
		c.ClusterLoadBalancer = *in.ClusterLoadBalancer
	}

	if in.ControlLoadBalancer != nil {
		c.ControlLoadBalancer = *in.ControlLoadBalancer
	}

	if in.Kubeconfig != nil {
		kc := strings.TrimSpace(*in.Kubeconfig)
		if kc != "" {
			ct, iv, tag, err := utils.EncryptForOrg(ac.OrganizationID, []byte(kc))
			if err != nil {
				http.Error(w, "kubeconfig encrypt failed", http.StatusInternalServerError)
				return
			}
			c.EncryptedKubeconfig = ct
			c.KubeIV = iv
			c.KubeTag = tag
		}
	}

	err = db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&c).Error; err != nil {
			return err
		}
		if len(poolIDs) > 0 {
			var pools []models.NodePool
			if err := tx.Where("id IN ?", poolIDs).Find(&pools).Error; err != nil {
				return err
			}
			if err := tx.Model(&c).Association("NodePools").Replace(&pools); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		http.Error(w, "create failed", http.StatusInternalServerError)
		return
	}

	tx := db.DB.Preload("NodePools").Preload("BastionServer")
	if err := tx.First(&c, "id = ?", c.ID).Error; err != nil {
		http.Error(w, "fetch failed", http.StatusInternalServerError)
		return
	}
	_ = response.JSON(w, http.StatusCreated, toResp(c, true, true))
}

// GetCluster godoc
// @Summary      Get cluster by ID (org scoped)
// @Description  Returns one cluster. Add `include=node_pools,bastion` to expand.
// @Tags         clusters
// @Security     BearerAuth
// @Produce      json
// @Param        X-Org-ID  header  string  true  "Organization UUID"
// @Param        id        path    string  true  "Cluster ID (UUID)"
// @Param        include   query   string  false "Optional: node_pools,bastion"
// @Success      200       {object} clusters.clusterResponse
// @Failure      400       {string} string "invalid id"
// @Failure      401       {string} string "Unauthorized"
// @Failure      403       {string} string "organization required"
// @Failure      404       {string} string "not found"
// @Failure      500       {string} string "fetch failed"
// @Router       /api/v1/clusters/{id} [get]
func GetCluster(w http.ResponseWriter, r *http.Request) {
	ac := middleware.GetAuthContext(r)
	if ac == nil || ac.OrganizationID == uuid.Nil {
		http.Error(w, "organization required", http.StatusForbidden)
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	include := strings.Split(strings.ToLower(r.URL.Query().Get("include")), ",")
	withPools := contains(include, "node_pools")
	withBastion := contains(include, "bastion")

	tx := db.DB.Where("id = ? AND organization_id = ?", id, ac.OrganizationID)
	if withPools {
		tx = tx.Preload("NodePools").
			Preload("NodePools.Taints").
			Preload("NodePools.Annotations").
			Preload("NodePools.Labels").
			Preload("NodePools.Servers")
	}
	if withBastion {
		tx = tx.Preload("BastionServer")
	}

	var c models.Cluster
	if err := tx.First(&c).Error; err != nil {
		if errorsIsNotFound(err) {
			http.Error(w, "not found", http.StatusNotFound)
		} else {
			http.Error(w, "fetch failed", http.StatusInternalServerError)
		}
		return
	}
	_ = response.JSON(w, http.StatusOK, toResp(c, withPools, withBastion))
}

// UpdateCluster godoc
// @Summary      Update cluster (org scoped). If `kubeconfig` is provided and non-empty, it will be encrypted per-organization and stored (never returned). Sending an empty string for `kubeconfig` is ignored (no change).
// @Tags         clusters
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        X-Org-ID  header  string                           true  "Organization UUID"
// @Param        id        path    string                           true  "Cluster ID (UUID)"
// @Param        body      body    clusters.updateClusterRequest     true  "payload"
// @Success      200       {object} clusters.clusterResponse
// @Failure      400       {string} string "invalid id / invalid json / invalid bastion_server_id"
// @Failure      401       {string} string "Unauthorized"
// @Failure      403       {string} string "organization required"
// @Failure      404       {string} string "not found"
// @Failure      500       {string} string "update failed"
// @Router       /api/v1/clusters/{id} [patch]
func UpdateCluster(w http.ResponseWriter, r *http.Request) {
	ac := middleware.GetAuthContext(r)
	if ac == nil || ac.OrganizationID == uuid.Nil {
		http.Error(w, "organization required", http.StatusForbidden)
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var c models.Cluster
	if err := db.DB.Where("id = ? AND organization_id = ?", id, ac.OrganizationID).First(&c).Error; err != nil {
		if errorsIsNotFound(err) {
			http.Error(w, "not found", http.StatusNotFound)
		} else {
			http.Error(w, "fetch failed", http.StatusInternalServerError)
		}
		return
	}

	var in updateClusterRequest
	if !readJSON(w, r, &in) {
		return
	}

	if in.Name != nil {
		c.Name = *in.Name
	}
	if in.Provider != nil {
		c.Provider = *in.Provider
	}
	if in.Region != nil {
		c.Region = *in.Region
	}
	if in.Status != nil {
		c.Status = *in.Status
	}
	if in.ClusterLoadBalancer != nil {
		c.ClusterLoadBalancer = *in.ClusterLoadBalancer
	}
	if in.ControlLoadBalancer != nil {
		c.ControlLoadBalancer = *in.ControlLoadBalancer
	}
	if in.Kubeconfig != nil {
		kc := strings.TrimSpace(*in.Kubeconfig)
		if kc != "" {
			ct, iv, tag, err := utils.EncryptForOrg(ac.OrganizationID, []byte(kc))
			if err != nil {
				http.Error(w, "kubeconfig encrypt failed", http.StatusInternalServerError)
				return
			}
			c.EncryptedKubeconfig = ct
			c.KubeIV = iv
			c.KubeTag = tag
		}
	}
	if in.BastionServerID != nil {
		if *in.BastionServerID == "" {
			c.BastionServerID = nil
		} else {
			bid, err := uuid.Parse(*in.BastionServerID)
			if err != nil {
				http.Error(w, "invalid bastion_server_id", http.StatusBadRequest)
				return
			}
			if err := ensureServerBelongsToOrgWithRole(ac.OrganizationID, bid, "bastion"); err != nil {
				http.Error(w, "invalid bastion_server_id", http.StatusBadRequest)
				return
			}
			c.BastionServerID = &bid
		}
	}

	if err := db.DB.Save(&c).Error; err != nil {
		http.Error(w, "update failed", http.StatusInternalServerError)
		return
	}

	db.DB.Preload("NodePools").Preload("BastionServer").First(&c, "id = ?", c.ID)
	_ = response.JSON(w, http.StatusOK, toResp(c, true, true))
}

// DeleteCluster godoc
// @Summary      Delete cluster (org scoped)
// @Tags         clusters
// @Security     BearerAuth
// @Param        X-Org-ID  header  string  true  "Organization UUID"
// @Param        id        path    string  true  "Cluster ID (UUID)"
// @Success      204       {string} string "No Content"
// @Failure      400       {string} string "invalid id"
// @Failure      401       {string} string "Unauthorized"
// @Failure      403       {string} string "organization required"
// @Failure      500       {string} string "delete failed"
// @Router       /api/v1/clusters/{id} [delete]
func DeleteCluster(w http.ResponseWriter, r *http.Request) {
	ac := middleware.GetAuthContext(r)
	if ac == nil || ac.OrganizationID == uuid.Nil {
		http.Error(w, "organization required", http.StatusForbidden)
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := db.DB.Where("id = ? AND organization_id = ?", id, ac.OrganizationID).Delete(&models.Cluster{}).Error; err != nil {
		http.Error(w, "delete failed", http.StatusInternalServerError)
		return
	}
	response.NoContent(w)
}

// ListClusterNodePools godoc
// @Summary      List node pools attached to a cluster (org scoped)
// @Tags         clusters
// @Security     BearerAuth
// @Produce      json
// @Param        X-Org-ID  header  string  true  "Organization UUID"
// @Param        id        path    string  true  "Cluster ID (UUID)"
// @Param        q         query   string  false "Name contains (case-insensitive)"
// @Success      200       {array}  clusters.nodePoolBrief
// @Failure      400       {string} string "invalid id"
// @Failure      401       {string} string "Unauthorized"
// @Failure      403       {string} string "organization required"
// @Failure      404       {string} string "not found"
// @Failure      500       {string} string "fetch failed"
// @Router       /api/v1/clusters/{id}/node_pools [get]
func ListClusterNodePools(w http.ResponseWriter, r *http.Request) {
	ac := middleware.GetAuthContext(r)
	if ac == nil || ac.OrganizationID == uuid.Nil {
		http.Error(w, "organization required", http.StatusForbidden)
		return
	}
	cid, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	q := strings.TrimSpace(r.URL.Query().Get("q"))

	// ensure cluster exists and belongs to org
	var exists int64
	if err := db.DB.Model(&models.Cluster{}).
		Where("id = ? AND organization_id = ?", cid, ac.OrganizationID).
		Count(&exists).Error; err != nil {
		http.Error(w, "fetch failed", http.StatusInternalServerError)
		return
	}
	if exists == 0 {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	var pools []models.NodePool
	tx := db.DB.
		Model(&models.NodePool{}).
		Joins("JOIN cluster_node_pools cnp ON cnp.node_pool_id = node_pools.id").
		Where("cnp.cluster_id = ? AND node_pools.organization_id = ?", cid, ac.OrganizationID)
	if q != "" {
		tx = tx.Where("LOWER(node_pools.name) LIKE ?", "%"+strings.ToLower(q)+"%")
	}
	if err := tx.Find(&pools).Error; err != nil {
		http.Error(w, "fetch failed", http.StatusInternalServerError)
		return
	}

	out := make([]nodePoolBrief, 0, len(pools))
	for _, p := range pools {
		out = append(out, nodePoolBrief{ID: p.ID, Name: p.Name})
	}
	_ = response.JSON(w, http.StatusOK, out)
}

// Attach/Detach NodePools

// AttachNodePools godoc
// @Summary      Attach node pools to cluster (org scoped)
// @Tags         clusters
// @Security     BearerAuth
// @Accept       json
// @Param        X-Org-ID  header  string                           true  "Organization UUID"
// @Param        id        path    string                           true  "Cluster ID (UUID)"
// @Param        body      body    clusters.attachNodePoolsRequest   true  "node_pool_ids"
// @Success      204       {string} string "No Content"
// @Failure      400       {string} string "invalid id / invalid node_pool_ids"
// @Failure      401       {string} string "Unauthorized"
// @Failure      403       {string} string "organization required"
// @Failure      404       {string} string "not found"
// @Failure      500       {string} string "attach failed"
// @Router       /api/v1/clusters/{id}/node_pools [post]
func AttachNodePools(w http.ResponseWriter, r *http.Request) {
	ac := middleware.GetAuthContext(r)
	if ac == nil || ac.OrganizationID == uuid.Nil {
		http.Error(w, "organization required", http.StatusForbidden)
		return
	}
	cid, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var c models.Cluster
	if err := db.DB.Where("id = ? AND organization_id = ?", cid, ac.OrganizationID).First(&c).Error; err != nil {
		if errorsIsNotFound(err) {
			http.Error(w, "not found", http.StatusNotFound)
		} else {
			http.Error(w, "fetch failed", http.StatusInternalServerError)
		}
		return
	}

	var in attachNodePoolsRequest
	if !readJSON(w, r, &in) {
		return
	}
	ids, err := parseUUIDs(in.NodePoolIDs)
	if err != nil {
		http.Error(w, "invalid node_pool_ids", http.StatusBadRequest)
		return
	}
	if err := ensureNodePoolsBelongToOrg(ac.OrganizationID, ids); err != nil {
		http.Error(w, "invalid node_pool_ids", http.StatusBadRequest)
		return
	}

	var pools []models.NodePool
	if err := db.DB.Where("id IN ?", ids).Find(&pools).Error; err != nil {
		http.Error(w, "attach failed", http.StatusInternalServerError)
		return
	}
	if err := db.DB.Model(&c).Association("NodePools").Append(&pools); err != nil {
		http.Error(w, "attach failed", http.StatusInternalServerError)
		return
	}
	response.NoContent(w)
}

// DetachNodePool godoc
// @Summary      Detach one node pool from a cluster (org scoped)
// @Tags         clusters
// @Security     BearerAuth
// @Param        X-Org-ID  header  string  true  "Organization UUID"
// @Param        id        path    string  true  "Cluster ID (UUID)"
// @Param        poolId    path    string  true  "Node Pool ID (UUID)"
// @Success      204       {string} string "No Content"
// @Failure      400       {string} string "invalid id"
// @Failure      401       {string} string "Unauthorized"
// @Failure      403       {string} string "organization required"
// @Failure      404       {string} string "not found"
// @Failure      500       {string} string "detach failed"
// @Router       /api/v1/clusters/{id}/node_pools/{poolId} [delete]
func DetachNodePool(w http.ResponseWriter, r *http.Request) {
	ac := middleware.GetAuthContext(r)
	if ac == nil || ac.OrganizationID == uuid.Nil {
		http.Error(w, "organization required", http.StatusForbidden)
		return
	}
	cid, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	pid, err := uuid.Parse(chi.URLParam(r, "poolId"))
	if err != nil {
		http.Error(w, "invalid poolId", http.StatusBadRequest)
		return
	}

	var c models.Cluster
	if err := db.DB.Where("id = ? AND organization_id = ?", cid, ac.OrganizationID).First(&c).Error; err != nil {
		if errorsIsNotFound(err) {
			http.Error(w, "not found", http.StatusNotFound)
		} else {
			http.Error(w, "fetch failed", http.StatusInternalServerError)
		}
		return
	}
	var p models.NodePool
	if err := db.DB.Where("id = ? AND organization_id = ?", pid, ac.OrganizationID).First(&p).Error; err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if err := db.DB.Model(&c).Association("NodePools").Delete(&p); err != nil {
		http.Error(w, "detach failed", http.StatusInternalServerError)
		return
	}
	response.NoContent(w)
}

// Bastion subresource

// GetBastion godoc
// @Summary      Get cluster bastion (org scoped)
// @Tags         clusters
// @Security     BearerAuth
// @Produce      json
// @Param        X-Org-ID  header  string  true  "Organization UUID"
// @Param        id        path    string  true  "Cluster ID (UUID)"
// @Success      200       {object} clusters.serverBrief
// @Success      204       {string} string "No Content (no bastion set)"
// @Failure      400       {string} string "invalid id"
// @Failure      401       {string} string "Unauthorized"
// @Failure      403       {string} string "organization required"
// @Failure      404       {string} string "not found"
// @Failure      500       {string} string "fetch failed"
// @Router       /api/v1/clusters/{id}/bastion [get]
func GetBastion(w http.ResponseWriter, r *http.Request) {
	ac := middleware.GetAuthContext(r)
	if ac == nil || ac.OrganizationID == uuid.Nil {
		http.Error(w, "organization required", http.StatusForbidden)
		return
	}
	cid, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var c models.Cluster
	if err := db.DB.Preload("BastionServer").
		Where("id = ? AND organization_id = ?", cid, ac.OrganizationID).
		First(&c).Error; err != nil {
		if errorsIsNotFound(err) {
			http.Error(w, "not found", http.StatusNotFound)
		} else {
			http.Error(w, "fetch failed", http.StatusInternalServerError)
		}
		return
	}
	if c.BastionServer == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	_ = response.JSON(w, http.StatusOK, serverBrief{
		ID: c.BastionServer.ID, Hostname: c.BastionServer.Hostname,
		IP: c.BastionServer.IPAddress, Role: c.BastionServer.Role, Status: c.BastionServer.Status,
	})
}

// PutBastion godoc
// @Summary      Set/replace cluster bastion (org scoped)
// @Tags         clusters
// @Security     BearerAuth
// @Accept       json
// @Param        X-Org-ID  header  string                       true  "Organization UUID"
// @Param        id        path    string                       true  "Cluster ID (UUID)"
// @Param        body      body    clusters.setBastionRequest   true  "server_id with role=bastion"
// @Success      204       {string} string "No Content"
// @Failure      400       {string} string "invalid id / invalid server_id / server not bastion"
// @Failure      401       {string} string "Unauthorized"
// @Failure      403       {string} string "organization required"
// @Failure      404       {string} string "cluster or server not found"
// @Failure      500       {string} string "update failed"
// @Router       /api/v1/clusters/{id}/bastion [post]
func PutBastion(w http.ResponseWriter, r *http.Request) {
	ac := middleware.GetAuthContext(r)
	if ac == nil || ac.OrganizationID == uuid.Nil {
		http.Error(w, "organization required", http.StatusForbidden)
		return
	}
	cid, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var in setBastionRequest
	if !readJSON(w, r, &in) {
		return
	}
	sid, err := uuid.Parse(in.ServerID)
	if err != nil {
		http.Error(w, "invalid server_id", http.StatusBadRequest)
		return
	}
	if err := ensureServerBelongsToOrgWithRole(ac.OrganizationID, sid, "bastion"); err != nil {
		http.Error(w, "server must exist in org and have role=bastion", http.StatusBadRequest)
		return
	}

	if err := db.DB.Model(&models.Cluster{}).
		Where("id = ? AND organization_id = ?", cid, ac.OrganizationID).
		Updates(map[string]any{"bastion_server_id": sid}).Error; err != nil {
		http.Error(w, "update failed", http.StatusInternalServerError)
		return
	}
	response.NoContent(w)
}

// DeleteBastion godoc
// @Summary      Clear cluster bastion (org scoped)
// @Tags         clusters
// @Security     BearerAuth
// @Param        X-Org-ID  header  string  true  "Organization UUID"
// @Param        id        path    string  true  "Cluster ID (UUID)"
// @Success      204       {string} string "No Content"
// @Failure      400       {string} string "invalid id"
// @Failure      401       {string} string "Unauthorized"
// @Failure      403       {string} string "organization required"
// @Failure      404       {string} string "not found"
// @Failure      500       {string} string "update failed"
// @Router       /api/v1/clusters/{id}/bastion [delete]
func DeleteBastion(w http.ResponseWriter, r *http.Request) {
	ac := middleware.GetAuthContext(r)
	if ac == nil || ac.OrganizationID == uuid.Nil {
		http.Error(w, "organization required", http.StatusForbidden)
		return
	}
	cid, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := db.DB.Model(&models.Cluster{}).
		Where("id = ? AND organization_id = ?", cid, ac.OrganizationID).
		Updates(map[string]any{"bastion_server_id": nil}).Error; err != nil {
		http.Error(w, "update failed", http.StatusInternalServerError)
		return
	}
	response.NoContent(w)
}
