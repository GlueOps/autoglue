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

// ListAnnotations godoc
//
//	@ID				ListAnnotations
//	@Summary		List annotations (org scoped)
//	@Description	Returns annotations for the organization in X-Org-ID. Filters: `key`, `value`, and `q` (key contains). Add `include=node_pools` to include linked node pools.
//	@Tags			Annotations
//	@Produce		json
//	@Param			X-Org-ID	header		string	false	"Organization UUID"
//	@Param			key			query		string	false	"Exact key"
//	@Param			value		query		string	false	"Exact value"
//	@Param			q			query		string	false	"key contains (case-insensitive)"
//	@Success		200			{array}		dto.AnnotationResponse
//	@Failure		401			{string}	string	"Unauthorized"
//	@Failure		403			{string}	string	"organization required"
//	@Failure		500			{string}	string	"failed to list annotations"
//	@Router			/annotations [get]
//	@Security		BearerAuth
//	@Security		OrgKeyAuth
//	@Security		OrgSecretAuth
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

		if out == nil {
			out = []dto.AnnotationResponse{}
		}
		utils.WriteJSON(w, http.StatusOK, out)
	}
}

// GetAnnotation godoc
//
//	@ID				GetAnnotation
//	@Summary		Get annotation by ID (org scoped)
//	@Description	Returns one annotation. Add `include=node_pools` to include node pools.
//	@Tags			Annotations
//	@Produce		json
//	@Param			X-Org-ID	header		string	false	"Organization UUID"
//	@Param			id			path		string	true	"Annotation ID (UUID)"
//	@Success		200			{object}	dto.AnnotationResponse
//	@Failure		400			{string}	string	"invalid id"
//	@Failure		401			{string}	string	"Unauthorized"
//	@Failure		403			{string}	string	"organization required"
//	@Failure		404			{string}	string	"not found"
//	@Failure		500			{string}	string	"fetch failed"
//	@Router			/annotations/{id} [get]
//	@Security		BearerAuth
//	@Security		OrgKeyAuth
//	@Security		OrgSecretAuth
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

// CreateAnnotation godoc
//
//	@ID				CreateAnnotation
//	@Summary		Create annotation (org scoped)
//	@Description	Creates an annotation.
//	@Tags			Annotations
//	@Accept			json
//	@Produce		json
//	@Param			X-Org-ID	header		string						false	"Organization UUID"
//	@Param			body		body		dto.CreateAnnotationRequest	true	"Annotation payload"
//	@Success		201			{object}	dto.AnnotationResponse
//	@Failure		400			{string}	string	"invalid json / missing fields"
//	@Failure		401			{string}	string	"Unauthorized"
//	@Failure		403			{string}	string	"organization required"
//	@Failure		500			{string}	string	"create failed"
//	@Router			/annotations [post]
//	@Security		BearerAuth
//	@Security		OrgKeyAuth
//	@Security		OrgSecretAuth
func CreateAnnotation(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}

		var req dto.CreateAnnotationRequest
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

		a := models.Annotation{
			AuditFields: common.AuditFields{OrganizationID: orgID},
			Key:         req.Key,
			Value:       req.Value,
		}

		if err := db.Create(&a).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		out := dto.AnnotationResponse{
			AuditFields: a.AuditFields,
			Key:         a.Key,
			Value:       a.Value,
		}
		utils.WriteJSON(w, http.StatusCreated, out)
	}
}

// UpdateAnnotation godoc
//
//	@ID				UpdateAnnotation
//	@Summary		Update annotation (org scoped)
//	@Description	Partially update annotation fields.
//	@Tags			Annotations
//	@Accept			json
//	@Produce		json
//	@Param			X-Org-ID	header		string						false	"Organization UUID"
//	@Param			id			path		string						true	"Annotation ID (UUID)"
//	@Param			body		body		dto.UpdateAnnotationRequest	true	"Fields to update"
//	@Success		200			{object}	dto.AnnotationResponse
//	@Failure		400			{string}	string	"invalid id / invalid json"
//	@Failure		401			{string}	string	"Unauthorized"
//	@Failure		403			{string}	string	"organization required"
//	@Failure		404			{string}	string	"not found"
//	@Failure		500			{string}	string	"update failed"
//	@Router			/annotations/{id} [patch]
//	@Security		BearerAuth
//	@Security		OrgKeyAuth
//	@Security		OrgSecretAuth
func UpdateAnnotation(db *gorm.DB) http.HandlerFunc {
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

		var a models.Annotation
		if err := db.Where("id = ? AND organization_id = ?", id, orgID).First(&a).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "not_found", "not_found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		var req dto.UpdateAnnotationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_request", "bad request")
			return
		}

		if req.Key != nil {
			a.Key = strings.TrimSpace(*req.Key)
		}
		if req.Value != nil {
			a.Value = strings.TrimSpace(*req.Value)
		}

		if err := db.Save(&a).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		out := dto.AnnotationResponse{
			AuditFields: a.AuditFields,
			Key:         a.Key,
			Value:       a.Value,
		}
		utils.WriteJSON(w, http.StatusOK, out)
	}
}

// DeleteAnnotation godoc
//
//	@ID				DeleteAnnotation
//	@Summary		Delete annotation (org scoped)
//	@Description	Permanently deletes the annotation.
//	@Tags			Annotations
//	@Produce		json
//	@Param			X-Org-ID	header	string	false	"Organization UUID"
//	@Param			id			path	string	true	"Annotation ID (UUID)"
//	@Success		204			"No Content"
//	@Failure		400			{string}	string	"invalid id"
//	@Failure		401			{string}	string	"Unauthorized"
//	@Failure		403			{string}	string	"organization required"
//	@Failure		500			{string}	string	"delete failed"
//	@Router			/annotations/{id} [delete]
//	@Security		BearerAuth
//	@Security		OrgKeyAuth
//	@Security		OrgSecretAuth
func DeleteAnnotation(db *gorm.DB) http.HandlerFunc {
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

		if err := db.Where("id = ? AND organization_id = ?", id, orgID).Delete(&models.Annotation{}).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
