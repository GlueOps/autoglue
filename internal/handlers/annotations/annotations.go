package annotations

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

/* ------------------------------- DTOs ----------------------------------- */

type nodePoolBrief struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type annotationResponse struct {
	ID        uuid.UUID       `json:"id"`
	Key       string          `json:"key"`
	Value     string          `json:"value"`
	NodePools []nodePoolBrief `json:"node_pools,omitempty"`
}

type createAnnotationRequest struct {
	Key         string   `json:"key"`
	Value       string   `json:"value"`
	NodePoolIDs []string `json:"node_pool_ids"`
}

type updateAnnotationRequest struct {
	Key   *string `json:"key,omitempty"`
	Value *string `json:"value,omitempty"`
}

type addAnnotationToNodePool struct {
	NodePoolIDs []string `json:"node_pool_ids"`
}

/* ------------------------------- Helpers -------------------------------- */

func toResp(a models.Annotation, includePools bool) annotationResponse {
	out := annotationResponse{
		ID:    a.ID,
		Key:   a.Key,
		Value: a.Value,
	}
	if includePools {
		for _, p := range a.NodePools {
			out.NodePools = append(out.NodePools, nodePoolBrief{ID: p.ID, Name: p.Name})
		}
	}
	return out
}

func parseUUIDs(in []string) ([]uuid.UUID, error) {
	out := make([]uuid.UUID, 0, len(in))
	for _, s := range in {
		id, err := uuid.Parse(s)
		if err != nil {
			return nil, err
		}
		out = append(out, id)
	}
	return out, nil
}

func ensureNodePoolsBelongToOrg(orgID uuid.UUID, ids []uuid.UUID) error {
	if len(ids) == 0 {
		return nil
	}
	var count int64
	if err := db.DB.Model(&models.NodePool{}).
		Where("id IN ? AND organization_id = ?", ids, orgID).
		Count(&count).Error; err != nil {
		return err
	}
	if count != int64(len(ids)) {
		return errors.New("one or more node pools do not belong to organization")
	}
	return nil
}

/* -------------------------------- Routes -------------------------------- */

// ListAnnotations godoc
// @ID           ListAnnotations
// @Summary      List annotations (org scoped)
// @Description  Returns annotations for the organization in X-Org-ID. Filters: `name`, `value`, and `q` (name contains). Add `include=node_pools` to include linked node pools.
// @Tags         annotations
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        name query string false "Exact name"
// @Param        value query string false "Exact value"
// @Param        q query string false "name contains (case-insensitive)"
// @Param        include query string false "Optional: node_pools"
// @Security     BearerAuth
// @Success      200 {array}  annotationResponse
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      500 {string} string "failed to list annotations"
// @Router       /api/v1/annotations [get]
func ListAnnotations(w http.ResponseWriter, r *http.Request) {
	ac := middleware.GetAuthContext(r)
	if ac == nil || ac.OrganizationID == uuid.Nil {
		http.Error(w, "organization required", http.StatusForbidden)
		return
	}

	q := db.DB.Where("organization_id = ?", ac.OrganizationID)
	if name := strings.TrimSpace(r.URL.Query().Get("name")); name != "" {
		q = q.Where(`name = ?`, name)
	}
	if val := strings.TrimSpace(r.URL.Query().Get("value")); val != "" {
		q = q.Where(`value = ?`, val)
	}
	if needle := strings.TrimSpace(r.URL.Query().Get("q")); needle != "" {
		q = q.Where(`name ILIKE ?`, "%"+needle+"%")
	}

	includePools := strings.EqualFold(strings.TrimSpace(r.URL.Query().Get("include")), "node_pools")
	if includePools {
		q = q.Preload("NodePools")
	}

	var rows []models.Annotation
	if err := q.Order("created_at DESC").Find(&rows).Error; err != nil {
		http.Error(w, "failed to list annotations", http.StatusInternalServerError)
		return
	}

	out := make([]annotationResponse, 0, len(rows))
	for _, a := range rows {
		out = append(out, toResp(a, includePools))
	}
	_ = response.JSON(w, http.StatusOK, out)
}

// GetAnnotation godoc
// @ID           GetAnnotation
// @Summary      Get annotation by ID (org scoped)
// @Description  Returns one annotation. Add `include=node_pools` to include node pools.
// @Tags         annotations
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        id path string true "Annotation ID (UUID)"
// @Param        include query string false "Optional: node_pools"
// @Security     BearerAuth
// @Success      200 {object} annotationResponse
// @Failure      400 {string} string "invalid id"
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      404 {string} string "not found"
// @Failure      500 {string} string "fetch failed"
// @Router       /api/v1/annotations/{id} [get]
func GetAnnotation(w http.ResponseWriter, r *http.Request) {
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

	includePools := strings.EqualFold(strings.TrimSpace(r.URL.Query().Get("include")), "node_pools")

	var a models.Annotation
	q := db.DB.Where("id = ? AND organization_id = ?", id, ac.OrganizationID)
	if includePools {
		q = q.Preload("NodePools")
	}

	if err := q.First(&a).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "fetch failed", http.StatusInternalServerError)
		return
	}

	_ = response.JSON(w, http.StatusOK, toResp(a, includePools))
}

// CreateAnnotation godoc
// @ID           CreateAnnotation
// @Summary      Create annotation (org scoped)
// @Description  Creates an annotation. Optionally link to node pools.
// @Tags         annotations
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        body body createAnnotationRequest true "Annotation payload"
// @Security     BearerAuth
// @Success      201 {object} annotationResponse
// @Failure      400 {string} string "invalid json / missing fields / invalid node_pool_ids"
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      500 {string} string "create failed"
// @Router       /api/v1/annotations [post]
func CreateAnnotation(w http.ResponseWriter, r *http.Request) {
	ac := middleware.GetAuthContext(r)
	if ac == nil || ac.OrganizationID == uuid.Nil {
		http.Error(w, "organization required", http.StatusForbidden)
		return
	}

	var req createAnnotationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.Key) == "" || strings.TrimSpace(req.Value) == "" {
		http.Error(w, "invalid json or missing key/value", http.StatusBadRequest)
		return
	}

	a := models.Annotation{
		OrganizationID: ac.OrganizationID,
		Key:            strings.TrimSpace(req.Key),
		Value:          strings.TrimSpace(req.Value),
	}

	if err := db.DB.Create(&a).Error; err != nil {
		http.Error(w, "create failed", http.StatusInternalServerError)
		return
	}

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
		if err := db.DB.Where("id IN ? AND organization_id = ?", ids, ac.OrganizationID).Find(&pools).Error; err != nil {
			http.Error(w, "attach failed", http.StatusInternalServerError)
			return
		}
		if err := db.DB.Model(&a).Association("NodePools").Append(&pools); err != nil {
			http.Error(w, "attach failed", http.StatusInternalServerError)
			return
		}
	}

	_ = response.JSON(w, http.StatusCreated, toResp(a, false))
}

// UpdateAnnotation godoc
// @ID           UpdateAnnotation
// @Summary      Update annotation (org scoped)
// @Description  Partially update annotation fields.
// @Tags         annotations
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        id path string true "Annotation ID (UUID)"
// @Param        body body updateAnnotationRequest true "Fields to update"
// @Security     BearerAuth
// @Success      200 {object} annotationResponse
// @Failure      400 {string} string "invalid id / invalid json"
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      404 {string} string "not found"
// @Failure      500 {string} string "update failed"
// @Router       /api/v1/annotations/{id} [patch]
func UpdateAnnotation(w http.ResponseWriter, r *http.Request) {
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

	var a models.Annotation
	if err := db.DB.Where("id = ? AND organization_id = ?", id, ac.OrganizationID).First(&a).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "fetch failed", http.StatusInternalServerError)
		return
	}

	var req updateAnnotationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if req.Key != nil {
		a.Key = strings.TrimSpace(*req.Key)
	}
	if req.Value != nil {
		a.Value = strings.TrimSpace(*req.Value)
	}

	if err := db.DB.Save(&a).Error; err != nil {
		http.Error(w, "update failed", http.StatusInternalServerError)
		return
	}
	_ = response.JSON(w, http.StatusOK, toResp(a, false))
}

// DeleteAnnotation godoc
// @ID           DeleteAnnotation
// @Summary      Delete annotation (org scoped)
// @Description  Permanently deletes the annotation.
// @Tags         annotations
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        id path string true "Annotation ID (UUID)"
// @Security     BearerAuth
// @Success      204 {string} string "No Content"
// @Failure      400 {string} string "invalid id"
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      500 {string} string "delete failed"
// @Router       /api/v1/annotations/{id} [delete]
func DeleteAnnotation(w http.ResponseWriter, r *http.Request) {
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

	if err := db.DB.Where("id = ? AND organization_id = ?", id, ac.OrganizationID).Delete(&models.Annotation{}).Error; err != nil {
		http.Error(w, "delete failed", http.StatusInternalServerError)
		return
	}

	response.NoContent(w)
}

// AddAnnotationToNodePools godoc
// @ID           AddAnnotationToNodePools
// @Summary      Attach annotation to node pools (org scoped)
// @Description  Links the annotation to one or more node pools in the same organization.
// @Tags         annotations
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        id path string true "Annotation ID (UUID)"
// @Param        body body addAnnotationToNodePool true "IDs to attach"
// @Param        include query string false "Optional: node_pools"
// @Security     BearerAuth
// @Success      200 {object} annotationResponse
// @Failure      400 {string} string "invalid id / invalid json / invalid node_pool_ids"
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      404 {string} string "not found"
// @Failure      500 {string} string "attach failed"
// @Router       /api/v1/annotations/{id}/node_pools [post]
func AddAnnotationToNodePools(w http.ResponseWriter, r *http.Request) {
	ac := middleware.GetAuthContext(r)
	if ac == nil || ac.OrganizationID == uuid.Nil {
		http.Error(w, "organization required", http.StatusForbidden)
		return
	}

	annID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var a models.Annotation
	if err := db.DB.Where("id = ? AND organization_id = ?", annID, ac.OrganizationID).First(&a).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "fetch failed", http.StatusInternalServerError)
		return
	}

	var in struct {
		NodePoolIDs []string `json:"node_pool_ids"`
	}
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

	var pools []models.NodePool
	if err := db.DB.Where("id IN ? AND organization_id = ?", ids, ac.OrganizationID).Find(&pools).Error; err != nil {
		http.Error(w, "attach failed", http.StatusInternalServerError)
		return
	}
	if err := db.DB.Model(&a).Association("NodePools").Append(&pools); err != nil {
		http.Error(w, "attach failed", http.StatusInternalServerError)
		return
	}

	includePools := strings.EqualFold(strings.TrimSpace(r.URL.Query().Get("include")), "node_pools")
	if includePools {
		if err := db.DB.Preload("NodePools").First(&a, "id = ? AND organization_id = ?", annID, ac.OrganizationID).Error; err != nil {
			http.Error(w, "fetch failed", http.StatusInternalServerError)
			return
		}
	}

	_ = response.JSON(w, http.StatusOK, toResp(a, includePools))
}

// RemoveAnnotationFromNodePool godoc
// @ID           RemoveAnnotationFromNodePool
// @Summary      Detach annotation from a node pool (org scoped)
// @Description  Unlinks the annotation from the specified node pool.
// @Tags         annotations
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        id path string true "Annotation ID (UUID)"
// @Param        poolId path string true "Node Pool ID (UUID)"
// @Security     BearerAuth
// @Success      204 {string} string "No Content"
// @Failure      400 {string} string "invalid id"
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      404 {string} string "not found"
// @Failure      500 {string} string "detach failed"
// @Router       /api/v1/annotations/{id}/node_pools/{poolId} [delete]
func RemoveAnnotationFromNodePool(w http.ResponseWriter, r *http.Request) {
	ac := middleware.GetAuthContext(r)
	if ac == nil || ac.OrganizationID == uuid.Nil {
		http.Error(w, "organization required", http.StatusForbidden)
		return
	}

	annID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	poolID, err := uuid.Parse(chi.URLParam(r, "poolId"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var a models.Annotation
	if err := db.DB.Where("id = ? AND organization_id = ?", annID, ac.OrganizationID).First(&a).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "fetch failed", http.StatusInternalServerError)
		return
	}

	var p models.NodePool
	if err := db.DB.Where("id = ? AND organization_id = ?", poolID, ac.OrganizationID).First(&p).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "fetch failed", http.StatusInternalServerError)
		return
	}

	if err := db.DB.Model(&a).Association("NodePools").Delete(&p); err != nil {
		http.Error(w, "detach failed", http.StatusInternalServerError)
		return
	}

	response.NoContent(w)
}

// ListNodePoolsWithAnnotation godoc
// @ID           ListNodePoolsWithAnnotation
// @Summary      List node pools linked to an annotation (org scoped)
// @Description  Returns node pools attached to the annotation. Supports `q` (name contains, case-insensitive).
// @Tags         annotations
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        id path string true "Annotation ID (UUID)"
// @Param        q query string false "Name contains (case-insensitive)"
// @Security     BearerAuth
// @Success      200 {array}  nodePoolBrief
// @Failure      400 {string} string "invalid id"
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      404 {string} string "not found"
// @Failure      500 {string} string "fetch failed"
// @Router       /api/v1/annotations/{id}/node_pools [get]
func ListNodePoolsWithAnnotation(w http.ResponseWriter, r *http.Request) {
	ac := middleware.GetAuthContext(r)
	if ac == nil || ac.OrganizationID == uuid.Nil {
		http.Error(w, "organization required", http.StatusForbidden)
		return
	}

	annID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	// Ensure the annotation exists within this org
	var a models.Annotation
	if err := db.DB.Where("id = ? AND organization_id = ?", annID, ac.OrganizationID).First(&a).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "fetch failed", http.StatusInternalServerError)
		return
	}

	// Find pools joined via the M2M table "node_annotations"
	q := db.DB.Model(&models.NodePool{}).
		Joins("JOIN node_annotations na ON na.node_pool_id = node_pools.id").
		Where("na.annotation_id = ? AND node_pools.organization_id = ?", annID, ac.OrganizationID)

	if needle := strings.TrimSpace(r.URL.Query().Get("q")); needle != "" {
		q = q.Where("node_pools.name ILIKE ?", "%"+needle+"%")
	}

	var pools []models.NodePool
	if err := q.Order("node_pools.created_at DESC").Find(&pools).Error; err != nil {
		http.Error(w, "fetch failed", http.StatusInternalServerError)
		return
	}

	out := make([]nodePoolBrief, 0, len(pools))
	for _, p := range pools {
		out = append(out, nodePoolBrief{ID: p.ID, Name: p.Name})
	}

	_ = response.JSON(w, http.StatusOK, out)
}
