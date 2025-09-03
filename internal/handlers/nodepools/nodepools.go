package nodepools

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/glueops/autoglue/internal/db"
	"github.com/glueops/autoglue/internal/db/models"
	"github.com/glueops/autoglue/internal/middleware"
	"github.com/glueops/autoglue/internal/response"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ListNodePools godoc
// @Summary      List node pools (org scoped)
// @Description  Returns node pools for the organization in X-Org-ID. Add `include=servers` to include attached servers. Filter by `q` (name contains).
// @Tags         node-pools
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        q query string false "Name contains (case-insensitive)"
// @Param        include query string false "Optional: servers"
// @Security     BearerAuth
// @Success      200 {array}  nodePoolResponse
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      500 {string} string "failed to list node groups"
// @Router       /api/v1/node-pools [get]
func ListNodePools(w http.ResponseWriter, r *http.Request) {
	ac := middleware.GetAuthContext(r)
	if ac == nil || ac.OrganizationID == uuid.Nil {
		http.Error(w, "organization required", http.StatusForbidden)
		return
	}

	q := db.DB.Where("organization_id = ?", ac.OrganizationID)
	if needle := strings.TrimSpace(r.URL.Query().Get("q")); needle != "" {
		q = q.Where("name ILIKE ?", "%"+needle+"%")
	}

	includeServers := strings.EqualFold(strings.TrimSpace(r.URL.Query().Get("include")), "servers")
	if includeServers {
		q = q.Preload("Servers")
	}

	var rows []models.NodePool
	if err := q.Order("created_at DESC").Find(&rows).Error; err != nil {
		http.Error(w, "failed to list node groups", http.StatusInternalServerError)
		return
	}

	out := make([]nodePoolResponse, 0, len(rows))
	for _, ng := range rows {
		out = append(out, toResp(ng, includeServers))
	}
	_ = response.JSON(w, http.StatusOK, out)
}

// GetNodePool godoc
// @Summary      Get node group by ID (org scoped)
// @Description  Returns one node group. Add `include=servers` to include servers.
// @Tags         node-pools
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        id path string true "Node Group ID (UUID)"
// @Param        include query string false "Optional: servers"
// @Security     BearerAuth
// @Success      200 {object} nodePoolResponse
// @Failure      400 {string} string "invalid id"
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      404 {string} string "not found"
// @Failure      500 {string} string "fetch failed"
// @Router       /api/v1/node-pools/{id} [get]
func GetNodePool(w http.ResponseWriter, r *http.Request) {
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

	includeServers := strings.EqualFold(strings.TrimSpace(r.URL.Query().Get("include")), "servers")

	var ng models.NodePool
	q := db.DB.Where("id = ? AND organization_id = ?", id, ac.OrganizationID)
	if includeServers {
		q = q.Preload("Servers")
	}
	if err := q.First(&ng).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "fetch failed", http.StatusInternalServerError)
		return
	}
	_ = response.JSON(w, http.StatusOK, toResp(ng, includeServers))
}

// CreateNodePool godoc
// @Summary      Create node group (org scoped)
// @Description  Creates a node group. Optionally attach initial servers.
// @Tags         node-pools
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        body body createNodePoolRequest true "NodeGroup payload"
// @Security     BearerAuth
// @Success      201 {object} nodePoolResponse
// @Failure      400 {string} string "invalid json / missing fields / invalid server_ids"
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      500 {string} string "create failed"
// @Router       /api/v1/node-pools [post]
func CreateNodePool(w http.ResponseWriter, r *http.Request) {
	ac := middleware.GetAuthContext(r)
	if ac == nil || ac.OrganizationID == uuid.Nil {
		http.Error(w, "organization required", http.StatusForbidden)
		return
	}

	var req createNodePoolRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.Name) == "" {
		http.Error(w, "invalid json or missing name", http.StatusBadRequest)
		return
	}

	ng := models.NodePool{
		OrganizationID: ac.OrganizationID,
		Name:           strings.TrimSpace(req.Name),
	}
	if err := db.DB.Create(&ng).Error; err != nil {
		http.Error(w, "create failed", http.StatusInternalServerError)
		return
	}

	// attach servers if provided
	if len(req.ServerIDs) > 0 {
		ids, err := parseUUIDs(req.ServerIDs)
		if err != nil {
			http.Error(w, "invalid server_ids", http.StatusBadRequest)
			return
		}
		if err := ensureServersBelongToOrg(ac.OrganizationID, ids); err != nil {
			http.Error(w, "invalid server_ids for this organization", http.StatusBadRequest)
			return
		}
		var servers []models.Server
		if err := db.DB.Where("id IN ? AND organization_id = ?", ids, ac.OrganizationID).
			Find(&servers).Error; err != nil {
			http.Error(w, "attach servers failed", http.StatusInternalServerError)
			return
		}
		if err := db.DB.Model(&ng).Association("Servers").Append(&servers); err != nil {
			http.Error(w, "attach servers failed", http.StatusInternalServerError)
			return
		}
	}

	_ = response.JSON(w, http.StatusCreated, toResp(ng, false))
}

// UpdateNodePool godoc
// @Summary      Update node pool (org scoped)
// @Description  Partially update node pool fields.
// @Tags         node-pools
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        id path string true "Node Pool ID (UUID)"
// @Param        body body updateNodePoolRequest true "Fields to update"
// @Security     BearerAuth
// @Success      200 {object} nodePoolResponse
// @Failure      400 {string} string "invalid id / invalid json"
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      404 {string} string "not found"
// @Failure      500 {string} string "update failed"
// @Router       /api/v1/node-pools/{id} [patch]
func UpdateNodePool(w http.ResponseWriter, r *http.Request) {
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

	var ng models.NodePool
	if err := db.DB.Where("id = ? AND organization_id = ?", id, ac.OrganizationID).First(&ng).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "fetch failed", http.StatusInternalServerError)
		return
	}

	var req updateNodePoolRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if req.Name != nil {
		ng.Name = strings.TrimSpace(*req.Name)
	}

	if err := db.DB.Save(&ng).Error; err != nil {
		http.Error(w, "update failed", http.StatusInternalServerError)
		return
	}
	_ = response.JSON(w, http.StatusOK, toResp(ng, false))
}

// DeleteNodePool godoc
// @Summary      Delete node pool (org scoped)
// @Description  Permanently deletes the node pool.
// @Tags         node-pools
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        id path string true "Node Group ID (UUID)"
// @Security     BearerAuth
// @Success      204 {string} string "No Content"
// @Failure      400 {string} string "invalid id"
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      500 {string} string "delete failed"
// @Router       /api/v1/node-pools/{id} [delete]
func DeleteNodePool(w http.ResponseWriter, r *http.Request) {
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

	if err := db.DB.Where("id = ? AND organization_id = ?", id, ac.OrganizationID).
		Delete(&models.NodePool{}).Error; err != nil {
		http.Error(w, "delete failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ListNodePoolServers godoc
// @Summary      List servers attached to a node pool (org scoped)
// @Tags         node-pools
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        id path string true "Node Group ID (UUID)"
// @Security     BearerAuth
// @Success      200 {array}  serverBrief
// @Failure      400 {string} string "invalid id"
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      404 {string} string "not found"
// @Failure      500 {string} string "fetch failed"
// @Router       /api/v1/node-pools/{id}/servers [get]
func ListNodePoolServers(w http.ResponseWriter, r *http.Request) {
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

	var ng models.NodePool
	if err := db.DB.Where("id = ? AND organization_id = ?", id, ac.OrganizationID).
		Preload("Servers").First(&ng).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "fetch failed", http.StatusInternalServerError)
		return
	}

	out := make([]serverBrief, 0, len(ng.Servers))
	for _, s := range ng.Servers {
		out = append(out, serverBrief{
			ID:       s.ID,
			Hostname: s.Hostname,
			IP:       s.IPAddress,
			Role:     s.Role,
			Status:   s.Status,
		})
	}
	_ = response.JSON(w, http.StatusOK, out)
}

// AttachNodePoolServers godoc
// @Summary      Attach servers to a node pool (org scoped)
// @Tags         node-pools
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        id path string true "Node Group ID (UUID)"
// @Param        body body attachServersRequest true "Server IDs to attach"
// @Security     BearerAuth
// @Success      204 {string} string "No Content"
// @Failure      400 {string} string "invalid id / invalid server_ids"
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      404 {string} string "not found"
// @Failure      500 {string} string "attach failed"
// @Router       /api/v1/node-pools/{id}/servers [post]
func AttachNodePoolServers(w http.ResponseWriter, r *http.Request) {
	ac := middleware.GetAuthContext(r)
	if ac == nil || ac.OrganizationID == uuid.Nil {
		http.Error(w, "organization required", http.StatusForbidden)
		return
	}
	ngID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var ng models.NodePool
	if err := db.DB.Where("id = ? AND organization_id = ?", ngID, ac.OrganizationID).First(&ng).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "fetch failed", http.StatusInternalServerError)
		return
	}

	var body struct {
		ServerIDs []string `json:"server_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || len(body.ServerIDs) == 0 {
		http.Error(w, "invalid server_ids", http.StatusBadRequest)
		return
	}

	ids, err := parseUUIDs(body.ServerIDs)
	if err != nil {
		http.Error(w, "invalid server_ids", http.StatusBadRequest)
		return
	}
	if err := ensureServersBelongToOrg(ac.OrganizationID, ids); err != nil {
		http.Error(w, "invalid server_ids for this organization", http.StatusBadRequest)
		return
	}

	var servers []models.Server
	if err := db.DB.Where("id IN ? AND organization_id = ?", ids, ac.OrganizationID).
		Find(&servers).Error; err != nil {
		http.Error(w, "attach failed", http.StatusInternalServerError)
		return
	}
	if err := db.DB.Model(&ng).Association("Servers").Append(&servers); err != nil {
		http.Error(w, "attach failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// DetachNodePoolServer godoc
// @Summary      Detach one server from a node pool (org scoped)
// @Tags         node-pools
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        id path string true "Node Pool ID (UUID)"
// @Param        serverId path string true "Server ID (UUID)"
// @Security     BearerAuth
// @Success      204 {string} string "No Content"
// @Failure      400 {string} string "invalid id"
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      404 {string} string "not found"
// @Failure      500 {string} string "detach failed"
// @Router       /api/v1/node-pools/{id}/servers/{serverId} [delete]
func DetachNodePoolServer(w http.ResponseWriter, r *http.Request) {
	ac := middleware.GetAuthContext(r)
	if ac == nil || ac.OrganizationID == uuid.Nil {
		http.Error(w, "organization required", http.StatusForbidden)
		return
	}
	ngID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	sid, err := uuid.Parse(chi.URLParam(r, "serverId"))
	if err != nil {
		http.Error(w, "invalid serverId", http.StatusBadRequest)
		return
	}

	var ng models.NodePool
	if err := db.DB.Where("id = ? AND organization_id = ?", ngID, ac.OrganizationID).
		Preload("Servers").First(&ng).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "fetch failed", http.StatusInternalServerError)
		return
	}

	var s models.Server
	if err := db.DB.Where("id = ? AND organization_id = ?", sid, ac.OrganizationID).First(&s).Error; err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	if err := db.DB.Model(&ng).Association("Servers").Delete(&s); err != nil {
		http.Error(w, "detach failed", http.StatusInternalServerError)
		return
	}
	response.NoContent(w)
}

// ListNodePoolTaints godoc
// @Summary      List taints attached to a node pool (org scoped)
// @Tags         node-pools
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        id path string true "Node Pool ID (UUID)"
// @Security     BearerAuth
// @Success      200 {array}  taintBrief
// @Failure      400 {string} string "invalid id"
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      404 {string} string "not found"
// @Failure      500 {string} string "fetch failed"
// @Router       /api/v1/node-pools/{id}/taints [get]
func ListNodePoolTaints(w http.ResponseWriter, r *http.Request) {
	ac := middleware.GetAuthContext(r)
	if ac == nil || ac.OrganizationID == uuid.Nil {
		http.Error(w, "organization required", http.StatusForbidden)
		return
	}

	ngID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var ng models.NodePool
	if err := db.DB.Where("id = ? AND organization_id = ?", ngID, ac.OrganizationID).
		Preload("Taints").
		First(&ng).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "fetch failed", http.StatusInternalServerError)
		return
	}

	out := make([]taintBrief, 0, len(ng.Taints))
	for _, t := range ng.Taints {
		out = append(out, taintBrief{
			ID:     t.ID,
			Key:    t.Key,
			Value:  t.Value,
			Effect: t.Effect,
		})
	}
	_ = response.JSON(w, http.StatusOK, out)
}

// AttachNodePoolTaints godoc
// @Summary      Attach taints to a node pool (org scoped)
// @Tags         node-pools
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        id path string true "Node Pool ID (UUID)"
// @Param        body body attachTaintsRequest true "Taint IDs to attach"
// @Security     BearerAuth
// @Success      204 {string} string "No Content"
// @Failure      400 {string} string "invalid id / invalid taint_ids"
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      404 {string} string "not found"
// @Failure      500 {string} string "attach failed"
// @Router       /api/v1/node-pools/{id}/taints [post]
func AttachNodePoolTaints(w http.ResponseWriter, r *http.Request) {
	ac := middleware.GetAuthContext(r)
	if ac == nil || ac.OrganizationID == uuid.Nil {
		http.Error(w, "organization required", http.StatusForbidden)
		return
	}

	ngID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var ng models.NodePool
	if err := db.DB.Where("id = ? AND organization_id = ?", ngID, ac.OrganizationID).
		First(&ng).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "fetch failed", http.StatusInternalServerError)
		return
	}

	var body struct {
		TaintIDs []string `json:"taint_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || len(body.TaintIDs) == 0 {
		http.Error(w, "invalid taint_ids", http.StatusBadRequest)
		return
	}

	ids, err := parseUUIDs(body.TaintIDs)
	if err != nil {
		http.Error(w, "invalid taint_ids", http.StatusBadRequest)
		return
	}
	if err := ensureTaintsBelongToOrg(ac.OrganizationID, ids); err != nil {
		http.Error(w, "invalid taint_ids for this organization", http.StatusBadRequest)
		return
	}

	var taints []models.Taint
	if err := db.DB.Where("id IN ? AND organization_id = ?", ids, ac.OrganizationID).
		Find(&taints).Error; err != nil {
		http.Error(w, "attach failed", http.StatusInternalServerError)
		return
	}
	if err := db.DB.Model(&ng).Association("Taints").Append(&taints); err != nil {
		http.Error(w, "attach failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// DetachNodePoolTaint godoc
// @Summary      Detach one taint from a node pool (org scoped)
// @Tags         node-pools
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        id path string true "Node Pool ID (UUID)"
// @Param        taintId path string true "Taint ID (UUID)"
// @Security     BearerAuth
// @Success      204 {string} string "No Content"
// @Failure      400 {string} string "invalid id"
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      404 {string} string "not found"
// @Failure      500 {string} string "detach failed"
// @Router       /api/v1/node-pools/{id}/taints/{taintId} [delete]
func DetachNodePoolTaint(w http.ResponseWriter, r *http.Request) {
	ac := middleware.GetAuthContext(r)
	if ac == nil || ac.OrganizationID == uuid.Nil {
		http.Error(w, "organization required", http.StatusForbidden)
		return
	}

	ngID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	tid, err := uuid.Parse(chi.URLParam(r, "taintId"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var ng models.NodePool
	if err := db.DB.Where("id = ? AND organization_id = ?", ngID, ac.OrganizationID).
		First(&ng).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "fetch failed", http.StatusInternalServerError)
		return
	}

	var t models.Taint
	if err := db.DB.Where("id = ? AND organization_id = ?", tid, ac.OrganizationID).
		First(&t).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "fetch failed", http.StatusInternalServerError)
		return
	}

	if err := db.DB.Model(&ng).Association("Taints").Delete(&t); err != nil {
		http.Error(w, "detach failed", http.StatusInternalServerError)
		return
	}
	response.NoContent(w)
}

// ListNodePoolLabels godoc
// @Summary      List labels attached to a node pool (org scoped)
// @Tags         node-pools
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        id path string true "Node Pool ID (UUID)"
// @Security     BearerAuth
// @Success      200 {array}  labelBrief
// @Failure      400 {string} string "invalid id"
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      404 {string} string "not found"
// @Failure      500 {string} string "fetch failed"
// @Router       /api/v1/node-pools/{id}/labels [get]
func ListNodePoolLabels(w http.ResponseWriter, r *http.Request) {
	ac := middleware.GetAuthContext(r)
	if ac == nil || ac.OrganizationID == uuid.Nil {
		http.Error(w, "organization required", http.StatusForbidden)
		return
	}

	ngID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var ng models.NodePool
	if err := db.DB.Where("id = ? AND organization_id = ?", ngID, ac.OrganizationID).
		Preload("Labels").First(&ng).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "fetch failed", http.StatusInternalServerError)
		return
	}

	out := make([]labelBrief, 0, len(ng.Labels))
	for _, l := range ng.Labels {
		out = append(out, labelBrief{
			ID:    l.ID,
			Key:   l.Key,
			Value: l.Value,
		})
	}
	_ = response.JSON(w, http.StatusOK, out)
}

// AttachNodePoolLabels godoc
// @Summary      Attach labels to a node pool (org scoped)
// @Tags         node-pools
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        id path string true "Node Pool ID (UUID)"
// @Param        body body attachLabelsRequest true "Label IDs to attach"
// @Security     BearerAuth
// @Success      204 {string} string "No Content"
// @Failure      400 {string} string "invalid id / invalid label_ids"
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      404 {string} string "not found"
// @Failure      500 {string} string "attach failed"
// @Router       /api/v1/node-pools/{id}/labels [post]
func AttachNodePoolLabels(w http.ResponseWriter, r *http.Request) {
	ac := middleware.GetAuthContext(r)
	if ac == nil || ac.OrganizationID == uuid.Nil {
		http.Error(w, "organization required", http.StatusForbidden)
		return
	}

	ngID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var ng models.NodePool
	if err := db.DB.Where("id = ? AND organization_id = ?", ngID, ac.OrganizationID).
		First(&ng).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "fetch failed", http.StatusInternalServerError)
		return
	}

	var body attachLabelsRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || len(body.LabelIDs) == 0 {
		http.Error(w, "invalid label_ids", http.StatusBadRequest)
		return
	}

	ids, err := parseUUIDs(body.LabelIDs) // already used in this package for servers/taints
	if err != nil {
		http.Error(w, "invalid label_ids", http.StatusBadRequest)
		return
	}
	if err := ensureLabelsBelongToOrg(ac.OrganizationID, ids); err != nil {
		http.Error(w, "invalid label_ids for this organization", http.StatusBadRequest)
		return
	}

	var labels []models.Label
	if err := db.DB.Where("id IN ? AND organization_id = ?", ids, ac.OrganizationID).
		Find(&labels).Error; err != nil {
		http.Error(w, "attach failed", http.StatusInternalServerError)
		return
	}
	if err := db.DB.Model(&ng).Association("Labels").Append(&labels); err != nil {
		http.Error(w, "attach failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// DetachNodePoolLabel godoc
// @Summary      Detach one label from a node pool (org scoped)
// @Tags         node-pools
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        id path string true "Node Pool ID (UUID)"
// @Param        labelId path string true "Label ID (UUID)"
// @Security     BearerAuth
// @Success      204 {string} string "No Content"
// @Failure      400 {string} string "invalid id"
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      404 {string} string "not found"
// @Failure      500 {string} string "detach failed"
// @Router       /api/v1/node-pools/{id}/labels/{labelId} [delete]
func DetachNodePoolLabel(w http.ResponseWriter, r *http.Request) {
	ac := middleware.GetAuthContext(r)
	if ac == nil || ac.OrganizationID == uuid.Nil {
		http.Error(w, "organization required", http.StatusForbidden)
		return
	}

	ngID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	lid, err := uuid.Parse(chi.URLParam(r, "labelId"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var ng models.NodePool
	if err := db.DB.Where("id = ? AND organization_id = ?", ngID, ac.OrganizationID).
		First(&ng).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "fetch failed", http.StatusInternalServerError)
		return
	}

	var l models.Label
	if err := db.DB.Where("id = ? AND organization_id = ?", lid, ac.OrganizationID).
		First(&l).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "fetch failed", http.StatusInternalServerError)
		return
	}

	if err := db.DB.Model(&ng).Association("Labels").Delete(&l); err != nil {
		http.Error(w, "detach failed", http.StatusInternalServerError)
		return
	}
	response.NoContent(w)
}
