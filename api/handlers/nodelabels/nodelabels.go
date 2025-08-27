package nodelabels

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

type labelResponse struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Value string    `json:"value"`

	// Present if include=node_groups
	NodeGroups []nodeGroupBrief `json:"node_groups,omitempty"`
}

type nodeGroupBrief struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type createLabelRequest struct {
	Name         string   `json:"name"`
	Value        string   `json:"value"`
	NodeGroupIDs []string `json:"node_group_ids,omitempty"`
}

type updateLabelRequest struct {
	Name  *string `json:"name,omitempty"`
	Value *string `json:"value,omitempty"`
}

type attachNodeGroupsRequest struct {
	NodeGroupIDs []string `json:"node_group_ids"`
}

func toResp(l models.NodeLabel, include bool) labelResponse {
	resp := labelResponse{
		ID:    l.ID,
		Name:  l.Name,
		Value: l.Value,
	}
	if include {
		resp.NodeGroups = make([]nodeGroupBrief, 0, len(l.NodeGroups))
		for _, ng := range l.NodeGroups {
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

// ListNodeLabels godoc
// @Summary      List node labels (org scoped)
// @Description  Returns node labels for the organization in X-Org-ID. Filters: `name`, `value`, and `q` (name contains). Add `include=node_groups` to include linked node groups.
// @Tags         node-labels
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        name query string false "Exact name"
// @Param        value query string false "Exact value"
// @Param        q query string false "Name contains (case-insensitive)"
// @Param        include query string false "Optional: node_groups"
// @Security     BearerAuth
// @Success      200 {array}  labelResponse
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      500 {string} string "failed to list node labels"
// @Router       /api/v1/node-labels [get]
func ListNodeLabels(w http.ResponseWriter, r *http.Request) {
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

	includeNG := strings.EqualFold(strings.TrimSpace(r.URL.Query().Get("include")), "node_groups")
	if includeNG {
		q = q.Preload("NodeGroups")
	}

	var rows []models.NodeLabel
	if err := q.Order("created_at DESC").Find(&rows).Error; err != nil {
		http.Error(w, "failed to list node labels", http.StatusInternalServerError)
		return
	}

	out := make([]labelResponse, 0, len(rows))
	for _, l := range rows {
		out = append(out, toResp(l, includeNG))
	}
	response.JSON(w, http.StatusOK, out)
}

// GetNodeLabel godoc
// @Summary      Get node label by ID (org scoped)
// @Description  Returns one label. Add `include=node_groups` to include node groups.
// @Tags         node-labels
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        id path string true "Node Label ID (UUID)"
// @Param        include query string false "Optional: node_groups"
// @Security     BearerAuth
// @Success      200 {object} labelResponse
// @Failure      400 {string} string "invalid id"
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      404 {string} string "not found"
// @Failure      500 {string} string "fetch failed"
// @Router       /api/v1/node-labels/{id} [get]
func GetNodeLabel(w http.ResponseWriter, r *http.Request) {
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

	includeNG := strings.EqualFold(strings.TrimSpace(r.URL.Query().Get("include")), "node_groups")

	var l models.NodeLabel
	q := db.DB.Where("id = ? AND organization_id = ?", id, ac.OrganizationID)
	if includeNG {
		q = q.Preload("NodeGroups")
	}
	if err := q.First(&l).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "fetch failed", http.StatusInternalServerError)
		return
	}
	response.JSON(w, http.StatusOK, toResp(l, includeNG))
}

// CreateNodeLabel godoc
// @Summary      Create node label (org scoped)
// @Description  Creates a label. Optionally link to node groups.
// @Tags         node-labels
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        body body createLabelRequest true "Label payload"
// @Security     BearerAuth
// @Success      201 {object} labelResponse
// @Failure      400 {string} string "invalid json / missing fields / invalid node_group_ids"
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      500 {string} string "create failed"
// @Router       /api/v1/node-labels [post]
func CreateNodeLabel(w http.ResponseWriter, r *http.Request) {
	ac := middleware.GetAuthContext(r)
	if ac == nil || ac.OrganizationID == uuid.Nil {
		http.Error(w, "organization required", http.StatusForbidden)
		return
	}

	var req createLabelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" || req.Value == "" {
		http.Error(w, "invalid json or missing name/value", http.StatusBadRequest)
		return
	}

	l := models.NodeLabel{
		OrganizationID: ac.OrganizationID,
		Name:           strings.TrimSpace(req.Name),
		Value:          strings.TrimSpace(req.Value),
	}
	if err := db.DB.Create(&l).Error; err != nil {
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
		if err := db.DB.Model(&l).Association("NodeGroups").Append(&ngs); err != nil {
			http.Error(w, "attach node groups failed", http.StatusInternalServerError)
			return
		}
	}

	response.JSON(w, http.StatusCreated, toResp(l, false))
}

// UpdateNodeLabel godoc
// @Summary      Update node label (org scoped)
// @Description  Partially update label fields.
// @Tags         node-labels
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        id path string true "Node Label ID (UUID)"
// @Param        body body updateLabelRequest true "Fields to update"
// @Security     BearerAuth
// @Success      200 {object} labelResponse
// @Failure      400 {string} string "invalid id / invalid json"
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      404 {string} string "not found"
// @Failure      500 {string} string "update failed"
// @Router       /api/v1/node-labels/{id} [patch]
func UpdateNodeLabel(w http.ResponseWriter, r *http.Request) {
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

	var l models.NodeLabel
	if err := db.DB.Where("id = ? AND organization_id = ?", id, ac.OrganizationID).First(&l).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "fetch failed", http.StatusInternalServerError)
		return
	}

	var req updateLabelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if req.Name != nil {
		l.Name = strings.TrimSpace(*req.Name)
	}
	if req.Value != nil {
		l.Value = strings.TrimSpace(*req.Value)
	}

	if err := db.DB.Save(&l).Error; err != nil {
		http.Error(w, "update failed", http.StatusInternalServerError)
		return
	}
	response.JSON(w, http.StatusOK, toResp(l, false))
}

// DeleteNodeLabel godoc
// @Summary      Delete node label (org scoped)
// @Description  Permanently deletes the label.
// @Tags         node-labels
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        id path string true "Node Label ID (UUID)"
// @Security     BearerAuth
// @Success      204 {string} string "No Content"
// @Failure      400 {string} string "invalid id"
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      500 {string} string "delete failed"
// @Router       /api/v1/node-labels/{id} [delete]
func DeleteNodeLabel(w http.ResponseWriter, r *http.Request) {
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
		Delete(&models.NodeLabel{}).Error; err != nil {
		http.Error(w, "delete failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Optional: attach/detach node groups to a label

// ListLabelNodeGroups godoc
// @Summary      List node groups linked to a label (org scoped)
// @Tags         node-labels
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        id path string true "Node Label ID (UUID)"
// @Security     BearerAuth
// @Success      200 {array}  nodeGroupBrief
// @Failure      400 {string} string "invalid id"
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      404 {string} string "not found"
// @Failure      500 {string} string "fetch failed"
// @Router       /api/v1/node-labels/{id}/node-groups [get]
func ListLabelNodeGroups(w http.ResponseWriter, r *http.Request) {
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

	var l models.NodeLabel
	if err := db.DB.Where("id = ? AND organization_id = ?", id, ac.OrganizationID).
		Preload("NodeGroups").First(&l).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "fetch failed", http.StatusInternalServerError)
		return
	}

	out := make([]nodeGroupBrief, 0, len(l.NodeGroups))
	for _, ng := range l.NodeGroups {
		out = append(out, nodeGroupBrief{ID: ng.ID, Name: ng.Name})
	}
	response.JSON(w, http.StatusOK, out)
}

// AttachLabelNodeGroups godoc
// @Summary      Attach node groups to a label (org scoped)
// @Tags         node-labels
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        id path string true "Node Label ID (UUID)"
// @Param        body body attachNodeGroupsRequest true "Node Group IDs to attach"
// @Security     BearerAuth
// @Success      204 {string} string "No Content"
// @Failure      400 {string} string "invalid id / invalid node_group_ids"
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      404 {string} string "not found"
// @Failure      500 {string} string "attach failed"
// @Router       /api/v1/node-labels/{id}/node-groups [post]
func AttachLabelNodeGroups(w http.ResponseWriter, r *http.Request) {
	ac := middleware.GetAuthContext(r)
	if ac == nil || ac.OrganizationID == uuid.Nil {
		http.Error(w, "organization required", http.StatusForbidden)
		return
	}
	lid, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var l models.NodeLabel
	if err := db.DB.Where("id = ? AND organization_id = ?", lid, ac.OrganizationID).First(&l).Error; err != nil {
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
		http.Error(w, "attach failed", http.StatusInternalServerError)
		return
	}
	if err := db.DB.Model(&l).Association("NodeGroups").Append(&ngs); err != nil {
		http.Error(w, "attach failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// DetachLabelNodeGroup godoc
// @Summary      Detach one node group from a label (org scoped)
// @Tags         node-labels
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        id path string true "Node Label ID (UUID)"
// @Param        nodeGroupId path string true "Node Group ID (UUID)"
// @Security     BearerAuth
// @Success      204 {string} string "No Content"
// @Failure      400 {string} string "invalid id"
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      404 {string} string "not found"
// @Failure      500 {string} string "detach failed"
// @Router       /api/v1/node-labels/{id}/node-groups/{nodeGroupId} [delete]
func DetachLabelNodeGroup(w http.ResponseWriter, r *http.Request) {
	ac := middleware.GetAuthContext(r)
	if ac == nil || ac.OrganizationID == uuid.Nil {
		http.Error(w, "organization required", http.StatusForbidden)
		return
	}
	vars := mux.Vars(r)
	lid, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	ngid, err := uuid.Parse(vars["nodeGroupId"])
	if err != nil {
		http.Error(w, "invalid nodeGroupId", http.StatusBadRequest)
		return
	}

	var l models.NodeLabel
	if err := db.DB.Where("id = ? AND organization_id = ?", lid, ac.OrganizationID).
		Preload("NodeGroups").First(&l).Error; err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	var ng models.NodeGroup
	if err := db.DB.Where("id = ? AND organization_id = ?", ngid, ac.OrganizationID).
		First(&ng).Error; err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if err := db.DB.Model(&l).Association("NodeGroups").Delete(&ng); err != nil {
		http.Error(w, "detach failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
