package taints

import (
	"encoding/json"
	"errors"
	"net/http"
	"slices"
	"strings"

	"github.com/glueops/autoglue/internal/db"
	"github.com/glueops/autoglue/internal/db/models"
	"github.com/glueops/autoglue/internal/middleware"
	"github.com/glueops/autoglue/internal/response"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ---------- Handlers ----------

// ListTaints godoc
// @Summary      List node taints (org scoped)
// @Description  Returns node taints for the organization in X-Org-ID. Filters: `key`, `value`, and `q` (key contains). Add `include=node_pools` to include linked node pools.
// @Tags         taints
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        key query string false "Exact key"
// @Param        value query string false "Exact value"
// @Param        q query string false "key contains (case-insensitive)"
// @Param        include query string false "Optional: node_pools"
// @Security     BearerAuth
// @Success      200 {array}  taintResponse
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      500 {string} string "failed to list node taints"
// @Router       /api/v1/taints [get]
func ListTaints(w http.ResponseWriter, r *http.Request) {
	ac := middleware.GetAuthContext(r)
	if ac == nil || ac.OrganizationID == uuid.Nil {
		http.Error(w, "organization required", http.StatusForbidden)
		return
	}

	q := db.DB.Where("organization_id = ?", ac.OrganizationID)

	if key := strings.TrimSpace(r.URL.Query().Get("key")); key != "" {
		q = q.Where(`key = ?`, key)
	}
	if val := strings.TrimSpace(r.URL.Query().Get("value")); val != "" {
		q = q.Where(`value = ?`, val)
	}
	if needle := strings.TrimSpace(r.URL.Query().Get("q")); needle != "" {
		q = q.Where(`key ILIKE ?`, "%"+needle+"%")
	}

	withPools := includeNodePools(r)
	if withPools {
		q = q.Preload("NodePools")
	}

	var rows []models.Taint
	if err := q.Order("created_at DESC").Find(&rows).Error; err != nil {
		http.Error(w, "failed to list node taints", http.StatusInternalServerError)
		return
	}

	out := make([]taintResponse, 0, len(rows))
	for _, t := range rows {
		out = append(out, toResp(t, withPools))
	}
	_ = response.JSON(w, http.StatusOK, out)
}

// GetTaint godoc
// @Summary      Get node taint by ID (org scoped)
// @Description  Returns one taint. Add `include=node_pools` to include node pools.
// @Tags         taints
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        id path string true "Node Taint ID (UUID)"
// @Param        include query string false "Optional: node_pools"
// @Security     BearerAuth
// @Success      200 {object} taintResponse
// @Failure      400 {string} string "invalid id"
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      404 {string} string "not found"
// @Failure      500 {string} string "fetch failed"
// @Router       /api/v1/taints/{id} [get]
func GetTaint(w http.ResponseWriter, r *http.Request) {
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

	withPools := includeNodePools(r)

	var t models.Taint
	q := db.DB.Where("id = ? AND organization_id = ?", id, ac.OrganizationID)
	if withPools {
		q = q.Preload("NodePools")
	}
	if err := q.First(&t).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "fetch failed", http.StatusInternalServerError)
		return
	}

	_ = response.JSON(w, http.StatusOK, toResp(t, withPools))
}

// CreateTaint godoc
// @Summary      Create node taint (org scoped)
// @Description  Creates a taint. Optionally link to node pools.
// @Tags         taints
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        body body createTaintRequest true "Taint payload"
// @Security     BearerAuth
// @Success      201 {object} taintResponse
// @Failure      400 {string} string "invalid json / missing fields / invalid node_pool_ids"
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      500 {string} string "create failed"
// @Router       /api/v1/taints [post]
func CreateTaint(w http.ResponseWriter, r *http.Request) {
	ac := middleware.GetAuthContext(r)
	if ac == nil || ac.OrganizationID == uuid.Nil {
		http.Error(w, "organization required", http.StatusForbidden)
		return
	}

	var req createTaintRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	req.Key = strings.TrimSpace(req.Key)
	req.Value = strings.TrimSpace(req.Value)
	req.Effect = strings.TrimSpace(req.Effect)

	if req.Key == "" || req.Effect == "" {
		http.Error(w, "invalid json or missing key/effect", http.StatusBadRequest)
		return
	}
	if _, ok := allowedEffects[req.Effect]; !ok {
		http.Error(w, "invalid effect", http.StatusBadRequest)
		return
	}

	t := models.Taint{
		OrganizationID: ac.OrganizationID,
		Key:            req.Key,
		Value:          req.Value,
		Effect:         req.Effect,
	}
	if err := db.DB.Create(&t).Error; err != nil {
		http.Error(w, "create failed", http.StatusInternalServerError)
		return
	}

	// optional initial links
	if len(req.NodePoolIDs) > 0 {
		ids, err := parseUUIDs(req.NodePoolIDs)
		if err != nil {
			http.Error(w, "invalid node_pool_ids", http.StatusBadRequest)
			return
		}
		if err := ensureNodePoolsBelongToOrg(ac.OrganizationID, ids); err != nil {
			http.Error(w, "invalid node_pool_ids for this organization", http.StatusBadRequest)
			return
		}

		var pools []models.NodePool
		if err := db.DB.Where("id IN ? AND organization_id = ?", ids, ac.OrganizationID).
			Find(&pools).Error; err != nil {
			http.Error(w, "create failed", http.StatusInternalServerError)
			return
		}
		if len(pools) != len(ids) {
			http.Error(w, "invalid node_pool_ids", http.StatusBadRequest)
			return
		}
		if err := db.DB.Model(&t).Association("NodePools").Append(&pools); err != nil {
			http.Error(w, "create failed", http.StatusInternalServerError)
			return
		}
	}

	_ = response.JSON(w, http.StatusCreated, toResp(t, false))
}

// UpdateTaint godoc
// @Summary      Update node taint (org scoped)
// @Description  Partially update taint fields.
// @Tags         taints
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
// @Router       /api/v1/taints/{id} [patch]
func UpdateTaint(w http.ResponseWriter, r *http.Request) {
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

	var t models.Taint
	if err := db.DB.Where("id = ? AND organization_id = ?", id, ac.OrganizationID).
		First(&t).Error; err != nil {
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

	if req.Key != nil {
		t.Key = strings.TrimSpace(*req.Key)
	}
	if req.Value != nil {
		t.Value = strings.TrimSpace(*req.Value)
	}
	if req.Effect != nil {
		e := strings.TrimSpace(*req.Effect)
		if e == "" {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		if _, ok := allowedEffects[e]; !ok {
			http.Error(w, "invalid effect", http.StatusBadRequest)
			return
		}
		t.Effect = e
	}

	if err := db.DB.Save(&t).Error; err != nil {
		http.Error(w, "update failed", http.StatusInternalServerError)
		return
	}
	_ = response.JSON(w, http.StatusOK, toResp(t, false))
}

// DeleteTaint godoc
// @Summary      Delete taint (org scoped)
// @Description  Permanently deletes the taint.
// @Tags         taints
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
// @Router       /api/v1/taints/{id} [delete]
func DeleteTaint(w http.ResponseWriter, r *http.Request) {
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
		Delete(&models.Taint{}).Error; err != nil {
		http.Error(w, "delete failed", http.StatusInternalServerError)
		return
	}
	response.NoContent(w)
}

// AddTaintToNodePool godoc
// @Summary      Attach taint to node pools (org scoped)
// @Description  Links the taint to one or more node pools in the same organization.
// @Tags         taints
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        id path string true "Taint ID (UUID)"
// @Param        body body addTaintToPoolRequest true "IDs to attach"
// @Param        include query string false "Optional: node_pools"
// @Security     BearerAuth
// @Success      200 {object} taintResponse
// @Failure      400 {string} string "invalid id / invalid json / invalid node_pool_ids"
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      404 {string} string "not found"
// @Failure      500 {string} string "attach failed"
// @Router       /api/v1/taints/{id}/node_pools [post]
func AddTaintToNodePool(w http.ResponseWriter, r *http.Request) {
	ac := middleware.GetAuthContext(r)
	if ac == nil || ac.OrganizationID == uuid.Nil {
		http.Error(w, "organization required", http.StatusForbidden)
		return
	}
	taintID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var t models.Taint
	if err := db.DB.Where("id = ? AND organization_id = ?", taintID, ac.OrganizationID).
		First(&t).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "fetch failed", http.StatusInternalServerError)
		return
	}

	var in addTaintToPoolRequest
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil || len(in.NodePoolIDs) == 0 {
		http.Error(w, "invalid json or empty node_pool_ids", http.StatusBadRequest)
		return
	}

	ids, err := parseUUIDs(in.NodePoolIDs)
	if err != nil {
		http.Error(w, "invalid node_pool_ids", http.StatusBadRequest)
		return
	}
	if err := ensureNodePoolsBelongToOrg(ac.OrganizationID, ids); err != nil {
		http.Error(w, "invalid node_pool_ids for this organization", http.StatusBadRequest)
		return
	}

	// Fetch existing links to avoid duplicates
	var existing []models.NodePool
	if err := db.DB.Model(&t).Association("NodePools").Find(&existing); err != nil {
		http.Error(w, "attach failed", http.StatusInternalServerError)
		return
	}
	existingIDs := make([]uuid.UUID, 0, len(existing))
	for _, p := range existing {
		existingIDs = append(existingIDs, p.ID)
	}

	toFetch := make([]uuid.UUID, 0, len(ids))
	for _, id := range ids {
		if !slices.Contains(existingIDs, id) {
			toFetch = append(toFetch, id)
		}
	}

	if len(toFetch) > 0 {
		var toAttach []models.NodePool
		if err := db.DB.Where("id IN ? AND organization_id = ?", toFetch, ac.OrganizationID).
			Find(&toAttach).Error; err != nil {
			http.Error(w, "attach failed", http.StatusInternalServerError)
			return
		}
		if len(toAttach) != len(toFetch) {
			http.Error(w, "invalid node_pool_ids", http.StatusBadRequest)
			return
		}
		if err := db.DB.Model(&t).Association("NodePools").Append(&toAttach); err != nil {
			http.Error(w, "attach failed", http.StatusInternalServerError)
			return
		}
	}

	withPools := includeNodePools(r)
	if withPools {
		if err := db.DB.Preload("NodePools").
			First(&t, "id = ? AND organization_id = ?", taintID, ac.OrganizationID).Error; err != nil {
			http.Error(w, "fetch failed", http.StatusInternalServerError)
			return
		}
	}

	_ = response.JSON(w, http.StatusOK, toResp(t, withPools))
}

// RemoveTaintFromNodePool godoc
// @Summary      Detach taint from a node pool (org scoped)
// @Description  Unlinks the taint from the specified node pool.
// @Tags         taints
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        id path string true "Taint ID (UUID)"
// @Param        poolId path string true "Node Pool ID (UUID)"
// @Security     BearerAuth
// @Success      204 {string} string "No Content"
// @Failure      400 {string} string "invalid id"
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      404 {string} string "not found"
// @Failure      500 {string} string "detach failed"
// @Router       /api/v1/taints/{id}/node_pools/{poolId} [delete]
func RemoveTaintFromNodePool(w http.ResponseWriter, r *http.Request) {
	ac := middleware.GetAuthContext(r)
	if ac == nil || ac.OrganizationID == uuid.Nil {
		http.Error(w, "organization required", http.StatusForbidden)
		return
	}

	taintID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	poolID, err := uuid.Parse(chi.URLParam(r, "poolId"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var t models.Taint
	if err := db.DB.Where("id = ? AND organization_id = ?", taintID, ac.OrganizationID).
		First(&t).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "fetch failed", http.StatusInternalServerError)
		return
	}

	var p models.NodePool
	if err := db.DB.Where("id = ? AND organization_id = ?", poolID, ac.OrganizationID).
		First(&p).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "fetch failed", http.StatusInternalServerError)
		return
	}

	if err := db.DB.Model(&t).Association("NodePools").Delete(&p); err != nil {
		http.Error(w, "detach failed", http.StatusInternalServerError)
		return
	}
	response.NoContent(w)
}

// ListNodePoolsWithTaint godoc
// @Summary      List node pools linked to a taint (org scoped)
// @Description  Returns node pools attached to the taint. Supports `q` (name contains, case-insensitive).
// @Tags         taints
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        id path string true "Taint ID (UUID)"
// @Param        q query string false "Name contains (case-insensitive)"
// @Security     BearerAuth
// @Success      200 {array}  nodePoolResponse
// @Failure      400 {string} string "invalid id"
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      404 {string} string "not found"
// @Failure      500 {string} string "fetch failed"
// @Router       /api/v1/taints/{id}/node_pools [get]
func ListNodePoolsWithTaint(w http.ResponseWriter, r *http.Request) {
	ac := middleware.GetAuthContext(r)
	if ac == nil || ac.OrganizationID == uuid.Nil {
		http.Error(w, "organization required", http.StatusForbidden)
		return
	}

	taintID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	// Load the taint and its pools using GORM's mapping (avoids guessing join table name)
	var t models.Taint
	if err := db.DB.Where("id = ? AND organization_id = ?", taintID, ac.OrganizationID).
		Preload("NodePools").
		First(&t).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "fetch failed", http.StatusInternalServerError)
		return
	}

	needle := strings.TrimSpace(r.URL.Query().Get("q"))
	out := make([]nodePoolResponse, 0, len(t.NodePools))
	for _, p := range t.NodePools {
		if needle != "" && !strings.Contains(strings.ToLower(p.Name), strings.ToLower(needle)) {
			continue
		}
		out = append(out, nodePoolResponse{
			ID:   p.ID,
			Name: p.Name,
			// Servers intentionally omitted here; this endpoint doesn't include them by default.
		})
	}

	_ = response.JSON(w, http.StatusOK, out)
}
