package labels

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

// ListLabels godoc
// @Summary      List node labels (org scoped)
// @Description  Returns node labels for the organization in X-Org-ID. Filters: `name`, `value`, and `q` (name contains). Add `include=node_pools` to include linked node groups.
// @Tags         labels
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        name query string false "Exact name"
// @Param        value query string false "Exact value"
// @Param        q query string false "Name contains (case-insensitive)"
// @Param        include query string false "Optional: node_pools"
// @Security     BearerAuth
// @Success      200 {array}  labelResponse
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      500 {string} string "failed to list node taints"
// @Router       /api/v1/labels [get]
func ListLabels(w http.ResponseWriter, r *http.Request) {
	ac := middleware.GetAuthContext(r)
	if ac == nil || ac.OrganizationID == uuid.Nil {
		http.Error(w, "organization required", http.StatusForbidden)
		return
	}

	q := db.DB.Where("organization_id = ?", ac.OrganizationID)
	if needle := strings.TrimSpace(r.URL.Query().Get("q")); needle != "" {
		q = q.Where("name ILIKE ?", "%"+needle+"%")
	}

	includePools := strings.EqualFold(strings.TrimSpace(r.URL.Query().Get("include")), "node_pools")
	if includePools {
		q.Preload("NodePools")
	}

	var rows []models.Label
	if err := q.Order("created_at DESC").Find(&rows).Error; err != nil {
		http.Error(w, "failed to list taints", http.StatusInternalServerError)
		return
	}

	out := make([]labelResponse, 0, len(rows))
	for _, np := range rows {
		out = append(out, toResp(np, includePools))
	}
	_ = response.JSON(w, http.StatusOK, out)
}

// GetLabel godoc
// @Summary      Get label by ID (org scoped)
// @Description  Returns one label. Add `include=node_pools` to include node groups.
// @Tags         labels
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        id path string true "Label ID (UUID)"
// @Param        include query string false "Optional: node_pools"
// @Security     BearerAuth
// @Success      200 {object} labelResponse
// @Failure      400 {string} string "invalid id"
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      404 {string} string "not found"
// @Failure      500 {string} string "fetch failed"
// @Router       /api/v1/labels/{id} [get]
func GetLabel(w http.ResponseWriter, r *http.Request) {
	ac := middleware.GetAuthContext(r)
	if ac == nil || ac.OrganizationID == uuid.Nil {
		http.Error(w, "organization required", http.StatusForbidden)
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid label id", http.StatusBadRequest)
		return
	}

	include := strings.EqualFold(strings.TrimSpace(r.URL.Query().Get("include")), "node_pools")

	var l models.Label
	q := db.DB.Where("id = ? AND organization_id = ?", id, ac.OrganizationID)
	if include {
		q = q.Preload("NodePools")
	}

	if err := q.First(&l).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "label not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to find label", http.StatusInternalServerError)
		return
	}

	_ = response.JSON(w, http.StatusOK, toResp(l, include))
}

// CreateLabel godoc
// @Summary      Create label (org scoped)
// @Description  Creates a label. Optionally link to node pools.
// @Tags         labels
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        body body createLabelRequest true "Label payload"
// @Security     BearerAuth
// @Success      201 {object} labelResponse
// @Failure      400 {string} string "invalid json / missing fields / invalid node_pool_ids"
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      500 {string} string "create failed"
// @Router       /api/v1/labels [post]
func CreateLabel(w http.ResponseWriter, r *http.Request) {
	ac := middleware.GetAuthContext(r)
	if ac == nil || ac.OrganizationID == uuid.Nil {
		http.Error(w, "organization required", http.StatusForbidden)
		return
	}

	var req createLabelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Key == "" || req.Value == "" {
		http.Error(w, "invalid json or missing key/value", http.StatusBadRequest)
		return
	}

	t := models.Label{
		OrganizationID: ac.OrganizationID,
		Key:            req.Key,
		Value:          req.Value,
	}

	if err := db.DB.Create(&t).Error; err != nil {
		http.Error(w, "failed to create label", http.StatusInternalServerError)
		return
	}

	if len(req.NodePoolIDs) > 0 {
		ids, err := parseUUIDs(req.NodePoolIDs)
		if err != nil {
			http.Error(w, "invalid node pool IDs", http.StatusBadRequest)
			return
		}

		if err := ensureNodePoolsBelongToOrg(ac.OrganizationID, ids); err != nil {
			http.Error(w, "invalid node pool IDs for this organization", http.StatusBadRequest)
			return
		}

		var nps []models.NodePool
		if err := db.DB.Where("id in ? AND organization_id = ?", ids, ac.OrganizationID).Find(&nps).Error; err != nil {
			http.Error(w, "node pools not found for this organization", http.StatusInternalServerError)
			return
		}
		if err := db.DB.Model(&t).Association("NodePools").Append(&nps); err != nil {
			http.Error(w, "attach node pools failed", http.StatusInternalServerError)
			return
		}
	}

	_ = response.JSON(w, http.StatusCreated, toResp(t, false))
}

// UpdateLabel godoc
// @Summary      Update label (org scoped)
// @Description  Partially update label fields.
// @Tags         labels
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        id path string true "Label ID (UUID)"
// @Param        body body updateLabelRequest true "Fields to update"
// @Security     BearerAuth
// @Success      200 {object} labelResponse
// @Failure      400 {string} string "invalid id / invalid json"
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      404 {string} string "not found"
// @Failure      500 {string} string "update failed"
// @Router       /api/v1/labels/{id} [patch]
func UpdateLabel(w http.ResponseWriter, r *http.Request) {
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

	var t models.Label
	if err := db.DB.Where("id = ? AND organization_id = ?", id, ac.OrganizationID).First(&t).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "fetch failed", http.StatusInternalServerError)
		return
	}

	var req updateLabelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json or missing key/value", http.StatusBadRequest)
		return
	}
	if req.Key != nil {
		t.Key = strings.TrimSpace(*req.Key)
	}
	if req.Value != nil {
		t.Value = strings.TrimSpace(*req.Value)
	}

	if err := db.DB.Save(&t).Error; err != nil {
		http.Error(w, "update failed", http.StatusInternalServerError)
		return
	}
	_ = response.JSON(w, http.StatusOK, toResp(t, false))
}

// DeleteLabel godoc
// @Summary      Delete label (org scoped)
// @Description  Permanently deletes the label.
// @Tags         labels
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        id path string true "Label ID (UUID)"
// @Security     BearerAuth
// @Success      204 {string} string "No Content"
// @Failure      400 {string} string "invalid id"
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      500 {string} string "delete failed"
// @Router       /api/v1/labels/{id} [delete]
func DeleteLabel(w http.ResponseWriter, r *http.Request) {
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

	if err := db.DB.Where("id = ? AND organization_id = ?", id, ac.OrganizationID).Delete(&models.Taint{}).Error; err != nil {
		http.Error(w, "delete failed", http.StatusInternalServerError)
		return
	}

	response.NoContent(w)
}
