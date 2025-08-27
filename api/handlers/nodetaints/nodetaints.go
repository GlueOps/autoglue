package nodetaints

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

type taintResponse struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Value string    `json:"value"`

	NodeGroups []nodeGroupBrief `json:"node_groups,omitempty"`
}

type nodeGroupBrief struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type createTaintRequest struct {
	Name         string   `json:"name"`
	Value        string   `json:"value"`
	NodeGroupIDs []string `json:"node_group_ids,omitempty"`
}

type updateTaintRequest struct {
	Name  *string `json:"name,omitempty"`
	Value *string `json:"value,omitempty"`
}

func toResp(t models.NodeTaint, include bool) taintResponse {
	resp := taintResponse{
		ID:    t.ID,
		Name:  t.Name,
		Value: t.Value,
	}
	if include {
		resp.NodeGroups = make([]nodeGroupBrief, 0, len(t.NodeGroups))
		for _, ng := range t.NodeGroups {
			resp.NodeGroups = append(resp.NodeGroups, nodeGroupBrief{ID: ng.ID, Name: ng.Name})
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

func ensureNodeGroupsBelongToOrg(orgID uuid.UUID, ids []uuid.UUID) error {
	var count int64
	if err := db.DB.Model(&models.NodeGroup{}).
		Where("organization_id = ? AND id IN ?", orgID, ids).
		Count(&count).Error; err != nil {
		return err
	}
	if count != int64(len(ids)) {
		return fmt.Errorf("some node groups do not belong to this organization")
	}
	return nil
}

// ListNodeTaints godoc
// @Summary      List node taints (org scoped)
// @Description  Returns node taints for the organization in X-Org-ID. Filters: `name`, `value`, and `q` (name contains). Add `include=node_groups` to include linked node groups.
// @Tags         node-taints
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        name query string false "Exact name"
// @Param        value query string false "Exact value"
// @Param        q query string false "Name contains (case-insensitive)"
// @Param        include query string false "Optional: node_groups"
// @Security     BearerAuth
// @Success      200 {array}  taintResponse
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      500 {string} string "failed to list node taints"
// @Router       /api/v1/node-taints [get]
func ListNodeTaints(w http.ResponseWriter, r *http.Request) {
	ac := middleware.GetAuthContext(r)
	if ac == nil || ac.OrganizationID == uuid.Nil {
		http.Error(w, "organization required", http.StatusForbidden)
		return
	}
	q := db.DB.Where("organization_id = ?", ac.OrganizationID)
	if name := strings.TrimSpace(r.URL.Query().Get("name")); name != "" {
		q = q.Where("name = ?", name)
	}
	if val := strings.TrimSpace(r.URL.Query().Get("value")); val != "" {
		q = q.Where("value = ?", val)
	}
	if needle := strings.TrimSpace(r.URL.Query().Get("q")); needle != "" {
		q = q.Where("name ILIKE ?", "%"+needle+"%")
	}

	include := strings.EqualFold(strings.TrimSpace(r.URL.Query().Get("include")), "node_groups")
	if include {
		q = q.Preload("NodeGroups")
	}

	var rows []models.NodeTaint
	if err := q.Order("created_at DESC").Find(&rows).Error; err != nil {
		http.Error(w, "failed to list node taints", http.StatusInternalServerError)
		return
	}
	out := make([]taintResponse, 0, len(rows))
	for _, t := range rows {
		out = append(out, toResp(t, include))
	}
	response.JSON(w, http.StatusOK, out)
}

// GetNodeTaint godoc
// @Summary      Get node taint by ID (org scoped)
// @Description  Returns one taint. Add `include=node_groups` to include node groups.
// @Tags         node-taints
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        id path string true "Node Taint ID (UUID)"
// @Param        include query string false "Optional: node_groups"
// @Security     BearerAuth
// @Success      200 {object} taintResponse
// @Failure      400 {string} string "invalid id"
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      404 {string} string "not found"
// @Failure      500 {string} string "fetch failed"
// @Router       /api/v1/node-taints/{id} [get]
func GetNodeTaint(w http.ResponseWriter, r *http.Request) {
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

	include := strings.EqualFold(strings.TrimSpace(r.URL.Query().Get("include")), "node_groups")
	var t models.NodeTaint
	q := db.DB.Where("id = ? AND organization_id = ?", id, ac.OrganizationID)
	if include {
		q = q.Preload("NodeGroups")
	}
	if err := q.First(&t).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "fetch failed", http.StatusInternalServerError)
		return
	}
	response.JSON(w, http.StatusOK, toResp(t, include))
}

// CreateNodeTaint godoc
// @Summary      Create node taint (org scoped)
// @Description  Creates a taint. Optionally link to node groups.
// @Tags         node-taints
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        body body createTaintRequest true "Taint payload"
// @Security     BearerAuth
// @Success      201 {object} taintResponse
// @Failure      400 {string} string "invalid json / missing fields / invalid node_group_ids"
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      500 {string} string "create failed"
// @Router       /api/v1/node-taints [post]
func CreateNodeTaint(w http.ResponseWriter, r *http.Request) {
	ac := middleware.GetAuthContext(r)
	if ac == nil || ac.OrganizationID == uuid.Nil {
		http.Error(w, "organization required", http.StatusForbidden)
		return
	}

	var req createTaintRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" || req.Value == "" {
		http.Error(w, "invalid json or missing name/value", http.StatusBadRequest)
		return
	}

	t := models.NodeTaint{
		OrganizationID: ac.OrganizationID,
		Name:           strings.TrimSpace(req.Name),
		Value:          strings.TrimSpace(req.Value),
	}
	if err := db.DB.Create(&t).Error; err != nil {
		http.Error(w, "create failed", http.StatusInternalServerError)
		return
	}

	if len(req.NodeGroupIDs) > 0 {
		ids, err := parseUUIDs(req.NodeGroupIDs)
		if err != nil {
			http.Error(w, "invalid node_group_ids", http.StatusBadRequest)
			return
		}
		if err := ensureNodeGroupsBelongToOrg(ac.OrganizationID, ids); err != nil {
			http.Error(w, "invalid node_group_ids for this organization", http.StatusBadRequest)
			return
		}
		var ngs []models.NodeGroup
		if err := db.DB.Where("id IN ? AND organization_id = ?", ids, ac.OrganizationID).
			Find(&ngs).Error; err != nil {
			http.Error(w, "attach node groups failed", http.StatusInternalServerError)
			return
		}
		if err := db.DB.Model(&t).Association("NodeGroups").Append(&ngs); err != nil {
			http.Error(w, "attach node groups failed", http.StatusInternalServerError)
			return
		}
	}
	response.JSON(w, http.StatusCreated, toResp(t, false))
}

// UpdateNodeTaint godoc
// @Summary      Update node taint (org scoped)
// @Description  Partially update taint fields.
// @Tags         node-taints
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        id path string true "Node Taint ID (UUID)"
// @Param        body body updateTaintRequest true "Fields to update"
// @Security     BearerAuth
// @Success      200 {object} taintResponse
// @Failure      400 {string} string "invalid id / invalid json"
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      404 {string} string "not found"
// @Failure      500 {string} string "update failed"
// @Router       /api/v1/node-taints/{id} [patch]
func UpdateNodeTaint(w http.ResponseWriter, r *http.Request) {
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

	var t models.NodeTaint
	if err := db.DB.Where("id = ? AND organization_id = ?", id, ac.OrganizationID).First(&t).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "fetch failed", http.StatusInternalServerError)
		return
	}

	var req updateTaintRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if req.Name != nil {
		t.Name = strings.TrimSpace(*req.Name)
	}
	if req.Value != nil {
		t.Value = strings.TrimSpace(*req.Value)
	}

	if err := db.DB.Save(&t).Error; err != nil {
		http.Error(w, "update failed", http.StatusInternalServerError)
		return
	}
	response.JSON(w, http.StatusOK, toResp(t, false))
}

// DeleteNodeTaint godoc
// @Summary      Delete node taint (org scoped)
// @Description  Permanently deletes the taint.
// @Tags         node-taints
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        id path string true "Node Taint ID (UUID)"
// @Security     BearerAuth
// @Success      204 {string} string "No Content"
// @Failure      400 {string} string "invalid id"
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      500 {string} string "delete failed"
// @Router       /api/v1/node-taints/{id} [delete]
func DeleteNodeTaint(w http.ResponseWriter, r *http.Request) {
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
		Delete(&models.NodeTaint{}).Error; err != nil {
		http.Error(w, "delete failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
