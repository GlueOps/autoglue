package nodepools

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/glueops/autoglue/api/middleware"
	"github.com/glueops/autoglue/internal/db"
	"github.com/glueops/autoglue/internal/db/models"
	"github.com/glueops/autoglue/internal/response"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type nodePoolResponse struct {
	ID      uuid.UUID     `json:"id"`
	Name    string        `json:"name"`
	Servers []serverBrief `json:"servers,omitempty"`
}

type serverBrief struct {
	ID       uuid.UUID `json:"id"`
	Hostname string    `json:"hostname"`
	IP       string    `json:"ip"`
	Role     string    `json:"role"`
	Status   string    `json:"status"`
}

type createNodePoolRequest struct {
	Name      string   `json:"name"`
	ServerIDs []string `json:"server_ids,omitempty"` // optional initial servers
}

type updateNodePoolRequest struct {
	Name *string `json:"name,omitempty"`
}

type attachServersRequest struct {
	ServerIDs []string `json:"server_ids"`
}

func toResp(ng models.NodeGroup, includeServers bool) nodePoolResponse {
	resp := nodePoolResponse{
		ID:   ng.ID,
		Name: ng.Name,
	}
	if includeServers {
		resp.Servers = make([]serverBrief, 0, len(ng.Servers))
		for _, s := range ng.Servers {
			resp.Servers = append(resp.Servers, serverBrief{
				ID:       s.ID,
				Hostname: s.Hostname,
				IP:       s.IPAddress,
				Role:     s.Role,
				Status:   s.Status,
			})
		}
	}
	return resp
}

func parseUUIDs(ids []string) ([]uuid.UUID, error) {
	out := make([]uuid.UUID, 0, len(ids))
	for _, s := range ids {
		u, err := uuid.Parse(strings.TrimSpace(s))
		if err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, nil
}

func ensureServersBelongToOrg(orgID uuid.UUID, ids []uuid.UUID) error {
	var count int64
	if err := db.DB.Model(&models.Server{}).
		Where("organization_id = ? AND id IN ?", orgID, ids).
		Count(&count).Error; err != nil {
		return err
	}
	if count != int64(len(ids)) {
		return fmt.Errorf("some servers do not belong to this organization")
	}
	return nil
}

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
// @Router       /api/v1/node-groups [get]
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

	var rows []models.NodeGroup
	if err := q.Order("created_at DESC").Find(&rows).Error; err != nil {
		http.Error(w, "failed to list node groups", http.StatusInternalServerError)
		return
	}

	out := make([]nodePoolResponse, 0, len(rows))
	for _, ng := range rows {
		out = append(out, toResp(ng, includeServers))
	}
	response.JSON(w, http.StatusOK, out)
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
// @Router       /api/v1/node-groups/{id} [get]
func GetNodePool(w http.ResponseWriter, r *http.Request) {
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

	includeServers := strings.EqualFold(strings.TrimSpace(r.URL.Query().Get("include")), "servers")

	var ng models.NodeGroup
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
	response.JSON(w, http.StatusOK, toResp(ng, includeServers))
}

// CreateNodeGroup godoc
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
// @Router       /api/v1/node-groups [post]
func CreateNodeGroup(w http.ResponseWriter, r *http.Request) {
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

	ng := models.NodeGroup{
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

	response.JSON(w, http.StatusCreated, toResp(ng, false))
}

// UpdateNodeGroup godoc
// @Summary      Update node group (org scoped)
// @Description  Partially update node group fields.
// @Tags         node-pools
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        id path string true "Node Group ID (UUID)"
// @Param        body body updateNodePoolRequest true "Fields to update"
// @Security     BearerAuth
// @Success      200 {object} nodePoolResponse
// @Failure      400 {string} string "invalid id / invalid json"
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      404 {string} string "not found"
// @Failure      500 {string} string "update failed"
// @Router       /api/v1/node-groups/{id} [patch]
func UpdateNodeGroup(w http.ResponseWriter, r *http.Request) {
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

	var ng models.NodeGroup
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
	response.JSON(w, http.StatusOK, toResp(ng, false))
}

// DeleteNodeGroup godoc
// @Summary      Delete node group (org scoped)
// @Description  Permanently deletes the node group.
// @Tags         node-groups
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
// @Router       /api/v1/node-groups/{id} [delete]
func DeleteNodeGroup(w http.ResponseWriter, r *http.Request) {
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
		Delete(&models.NodeGroup{}).Error; err != nil {
		http.Error(w, "delete failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// --- Server linking under a node group ---

// ListNodeGroupServers godoc
// @Summary      List servers attached to a node group (org scoped)
// @Tags         node-groups
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
// @Router       /api/v1/node-groups/{id}/servers [get]
func ListNodeGroupServers(w http.ResponseWriter, r *http.Request) {
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

	var ng models.NodeGroup
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
	response.JSON(w, http.StatusOK, out)
}

// AttachNodeGroupServers godoc
// @Summary      Attach servers to a node group (org scoped)
// @Tags         node-groups
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
// @Router       /api/v1/node-groups/{id}/servers [post]
func AttachNodeGroupServers(w http.ResponseWriter, r *http.Request) {
	ac := middleware.GetAuthContext(r)
	if ac == nil || ac.OrganizationID == uuid.Nil {
		http.Error(w, "organization required", http.StatusForbidden)
		return
	}
	ngID, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var ng models.NodeGroup
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

// DetachNodeGroupServer godoc
// @Summary      Detach one server from a node group (org scoped)
// @Tags         node-groups
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        id path string true "Node Group ID (UUID)"
// @Param        serverId path string true "Server ID (UUID)"
// @Security     BearerAuth
// @Success      204 {string} string "No Content"
// @Failure      400 {string} string "invalid id"
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      404 {string} string "not found"
// @Failure      500 {string} string "detach failed"
// @Router       /api/v1/node-groups/{id}/servers/{serverId} [delete]
func DetachNodeGroupServer(w http.ResponseWriter, r *http.Request) {
	ac := middleware.GetAuthContext(r)
	if ac == nil || ac.OrganizationID == uuid.Nil {
		http.Error(w, "organization required", http.StatusForbidden)
		return
	}
	vars := mux.Vars(r)
	ngID, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	sid, err := uuid.Parse(vars["serverId"])
	if err != nil {
		http.Error(w, "invalid serverId", http.StatusBadRequest)
		return
	}

	var ng models.NodeGroup
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
	w.WriteHeader(http.StatusNoContent)
}
