package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/glueops/autoglue/internal/api/httpmiddleware"
	"github.com/glueops/autoglue/internal/handlers/dto"
	"github.com/glueops/autoglue/internal/models"
	"github.com/glueops/autoglue/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ListAnnotations godoc
// @ID           ListAnnotations
// @Summary      List annotations (org scoped)
// @Description  Returns annotations for the organization in X-Org-ID. Filters: `name`, `value`, and `q` (name contains). Add `include=node_pools` to include linked node pools.
// @Tags         Annotations
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string false "Organization UUID"
// @Param        name query string false "Exact name"
// @Param        value query string false "Exact value"
// @Param        q query string false "name contains (case-insensitive)"
// @Success      200 {array}  dto.AnnotationResponse
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      500 {string} string "failed to list annotations"
// @Router       /annotations [get]
// @Security     BearerAuth
// @Security     OrgKeyAuth
// @Security     OrgSecretAuth
func ListAnnotations(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}

		q := db.Where("organization_id = ?", orgID)

		if key := strings.TrimSpace(r.URL.Query().Get("key")); key != "" {
			q = q.Where(`key = ?`, key)
		}
		if val := strings.TrimSpace(r.URL.Query().Get("value")); val != "" {
			q = q.Where(`value = ?`, val)
		}
		if needle := strings.TrimSpace(r.URL.Query().Get("q")); needle != "" {
			q = q.Where(`key ILIKE ?`, "%"+needle+"%")
		}

		var out []dto.AnnotationResponse
		if err := q.Model(&models.Annotation{}).Order("created_at DESC").Scan(&out).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}
		utils.WriteJSON(w, http.StatusOK, out)
	}
}

// GetAnnotation godoc
// @ID           GetAnnotation
// @Summary      Get annotation by ID (org scoped)
// @Description  Returns one annotation. Add `include=node_pools` to include node pools.
// @Tags         Annotations
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string false "Organization UUID"
// @Param        id path string true "Annotation ID (UUID)"
// @Param        include query string false "Optional: node_pools"
// @Success      200 {object} dto.AnnotationResponse
// @Failure      400 {string} string "invalid id"
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      404 {string} string "not found"
// @Failure      500 {string} string "fetch failed"
// @Router       /annotations/{id} [get]
// @Security     BearerAuth
// @Security     OrgKeyAuth
// @Security     OrgSecretAuth
func GetAnnotation(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_request", "bad request")
			return
		}

		var out dto.AnnotationResponse
		if err := db.Model(&models.Annotation{}).Where("id = ? AND organization_id = ?", id, orgID).First(&out).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "not_found", "not_found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}
		utils.WriteJSON(w, http.StatusOK, out)
	}
}
