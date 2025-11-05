package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/glueops/autoglue/internal/api/httpmiddleware"
	"github.com/glueops/autoglue/internal/common"
	"github.com/glueops/autoglue/internal/handlers/dto"
	"github.com/glueops/autoglue/internal/models"
	"github.com/glueops/autoglue/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ListLabels godoc
// @ID           ListLabels
// @Summary      List node labels (org scoped)
// @Description  Returns node labels for the organization in X-Org-ID. Filters: `key`, `value`, and `q` (key contains). Add `include=node_pools` to include linked node groups.
// @Tags         Labels
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string false "Organization UUID"
// @Param        key query string false "Exact key"
// @Param        value query string false "Exact value"
// @Param        q query string false "Key contains (case-insensitive)"
// @Success      200 {array}  dto.LabelResponse
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      500 {string} string "failed to list node taints"
// @Router       /labels [get]
// @Security    BearerAuth
// @Security    OrgKeyAuth
// @Security    OrgSecretAuth
func ListLabels(db *gorm.DB) http.HandlerFunc {
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
		var out []dto.LabelResponse
		if err := q.Model(&models.Label{}).Order("created_at DESC").Scan(&out).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}
		if out == nil {
			out = []dto.LabelResponse{}
		}

		utils.WriteJSON(w, http.StatusOK, out)
	}
}

// GetLabel godoc
// @ID           GetLabel
// @Summary      Get label by ID (org scoped)
// @Description  Returns one label.
// @Tags         Labels
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string false "Organization UUID"
// @Param        id path string true "Label ID (UUID)"
// @Success      200 {object} dto.LabelResponse
// @Failure      400 {string} string "invalid id"
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      404 {string} string "not found"
// @Failure      500 {string} string "fetch failed"
// @Router       /labels/{id} [get]
// @Security    BearerAuth
// @Security    OrgKeyAuth
// @Security    OrgSecretAuth
func GetLabel(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "id_required", "id required")
			return
		}

		var out dto.LabelResponse
		if err := db.Model(&models.Label{}).Where("id = ? AND organization_id = ?", id, orgID).Limit(1).Scan(&out).Error; err != nil {
			if out.ID == uuid.Nil {
				utils.WriteError(w, http.StatusNotFound, "label_not_found", "label not found")
				return
			}
		}

		utils.WriteJSON(w, http.StatusOK, out)
	}
}

// CreateLabel godoc
// @ID           CreateLabel
// @Summary      Create label (org scoped)
// @Description  Creates a label.
// @Tags         Labels
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string false "Organization UUID"
// @Param        body body dto.CreateLabelRequest true "Label payload"
// @Success      201 {object} dto.LabelResponse
// @Failure      400 {string} string "invalid json / missing fields / invalid node_pool_ids"
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      500 {string} string "create failed"
// @Router       /labels [post]
// @Security     BearerAuth
// @Security     OrgKeyAuth
// @Security     OrgSecretAuth
func CreateLabel(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}

		var req dto.CreateLabelRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_request", "bad request")
			return
		}

		req.Key = strings.TrimSpace(req.Key)
		req.Value = strings.TrimSpace(req.Value)

		if req.Key == "" || req.Value == "" {
			utils.WriteError(w, http.StatusBadRequest, "bad_request", "missing key/value")
			return
		}

		l := models.Label{
			AuditFields: common.AuditFields{OrganizationID: orgID},
			Key:         req.Key,
			Value:       req.Value,
		}
		if err := db.Create(&l).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		out := dto.LabelResponse{
			AuditFields: l.AuditFields,
			Key:         l.Key,
			Value:       l.Value,
		}
		utils.WriteJSON(w, http.StatusCreated, out)
	}
}

// UpdateLabel godoc
// UpdateLabel godoc
// @ID           UpdateLabel
// @Summary      Update label (org scoped)
// @Description  Partially update label fields.
// @Tags         Labels
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string false "Organization UUID"
// @Param        id path string true "Label ID (UUID)"
// @Param        body body dto.UpdateLabelRequest true "Fields to update"
// @Success      200 {object} dto.LabelResponse
// @Failure      400 {string} string "invalid id / invalid json"
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      404 {string} string "not found"
// @Failure      500 {string} string "update failed"
// @Router       /labels/{id} [patch]
// @Security     BearerAuth
// @Security     OrgKeyAuth
// @Security     OrgSecretAuth
func UpdateLabel(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "id_required", "id required")
			return
		}

		var l models.Label
		if err := db.Where("id = ? AND organization_id = ?", id, orgID).First(&l).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "label_not_found", "label not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		var req dto.UpdateLabelRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_request", "bad request")
			return
		}

		if req.Key != nil {
			l.Key = strings.TrimSpace(*req.Key)
		}
		if req.Value != nil {
			l.Value = strings.TrimSpace(*req.Value)
		}

		if err := db.Save(&l).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		out := dto.LabelResponse{
			AuditFields: l.AuditFields,
			Key:         l.Key,
			Value:       l.Value,
		}
		utils.WriteJSON(w, http.StatusOK, out)
	}
}

// DeleteLabel godoc
// @ID           DeleteLabel
// @Summary      Delete label (org scoped)
// @Description  Permanently deletes the label.
// @Tags         Labels
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string false "Organization UUID"
// @Param        id path string true "Label ID (UUID)"
// @Success      204 {string} string "No Content"
// @Failure      400 {string} string "invalid id"
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      500 {string} string "delete failed"
// @Router       /labels/{id} [delete]
// @Security     BearerAuth
// @Security     OrgKeyAuth
// @Security     OrgSecretAuth
func DeleteLabel(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "id_required", "id required")
			return
		}

		if err := db.Where("id = ? AND organization_id = ?", id, orgID).Delete(&models.Label{}).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
