package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/glueops/autoglue/internal/api/httpmiddleware"
	"github.com/glueops/autoglue/internal/handlers/dto"
	"github.com/glueops/autoglue/internal/models"
	"github.com/glueops/autoglue/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ListTaints godoc
//
//	@ID				ListTaints
//	@Summary		List node pool taints (org scoped)
//	@Description	Returns node taints for the organization in X-Org-ID. Filters: `key`, `value`, and `q` (key contains). Add `include=node_pools` to include linked node pools.
//	@Tags			Taints
//	@Accept			json
//	@Produce		json
//	@Param			X-Org-ID	header		string	false	"Organization UUID"
//	@Param			key			query		string	false	"Exact key"
//	@Param			value		query		string	false	"Exact value"
//	@Param			q			query		string	false	"key contains (case-insensitive)"
//	@Success		200			{array}		dto.TaintResponse
//	@Failure		401			{string}	string	"Unauthorized"
//	@Failure		403			{string}	string	"organization required"
//	@Failure		500			{string}	string	"failed to list node taints"
//	@Router			/taints [get]
//	@Security		BearerAuth
//	@Security		OrgKeyAuth
//	@Security		OrgSecretAuth
func ListTaints(db *gorm.DB) http.HandlerFunc {
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

		var out []dto.TaintResponse
		if err := q.Model(&models.Taint{}).Order("created_at DESC").Find(&out).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}
		utils.WriteJSON(w, http.StatusOK, out)
	}
}

// GetTaint godoc
//
//	@ID			GetTaint
//	@Summary	Get node taint by ID (org scoped)
//	@Tags		Taints
//	@Accept		json
//	@Produce	json
//	@Param		X-Org-ID	header		string	false	"Organization UUID"
//	@Param		id			path		string	true	"Node Taint ID (UUID)"
//	@Success	200			{object}	dto.TaintResponse
//	@Failure	400			{string}	string	"invalid id"
//	@Failure	401			{string}	string	"Unauthorized"
//	@Failure	403			{string}	string	"organization required"
//	@Failure	404			{string}	string	"not found"
//	@Failure	500			{string}	string	"fetch failed"
//	@Router		/taints/{id} [get]
//	@Security	BearerAuth
//	@Security	OrgKeyAuth
//	@Security	OrgSecretAuth
func GetTaint(db *gorm.DB) http.HandlerFunc {
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

		var out dto.TaintResponse
		if err := db.Model(&models.Taint{}).Where("id = ? AND organization_id = ?", id, orgID).First(&out).Error; err != nil {
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

// CreateTaint godoc
//
//	@ID				CreateTaint
//	@Summary		Create node taint (org scoped)
//	@Description	Creates a taint.
//	@Tags			Taints
//	@Accept			json
//	@Produce		json
//	@Param			X-Org-ID	header		string					false	"Organization UUID"
//	@Param			body		body		dto.CreateTaintRequest	true	"Taint payload"
//	@Success		201			{object}	dto.TaintResponse
//	@Failure		400			{string}	string	"invalid json / missing fields / invalid node_pool_ids"
//	@Failure		401			{string}	string	"Unauthorized"
//	@Failure		403			{string}	string	"organization required"
//	@Failure		500			{string}	string	"create failed"
//	@Router			/taints [post]
//	@Security		BearerAuth
//	@Security		OrgKeyAuth
//	@Security		OrgSecretAuth
func CreateTaint(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}

		var req dto.CreateTaintRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_request", "bad request")
			return
		}

		req.Key = strings.TrimSpace(req.Key)
		req.Value = strings.TrimSpace(req.Value)
		req.Effect = strings.TrimSpace(req.Effect)

		if req.Key == "" || req.Value == "" || req.Effect == "" {
			utils.WriteError(w, http.StatusBadRequest, "bad_request", "missing key/value/effect")
			return
		}

		if _, ok := allowedEffects[req.Effect]; !ok {
			utils.WriteError(w, http.StatusBadRequest, "bad_request", "invalid effect")
			return
		}

		t := models.Taint{
			OrganizationID: orgID,
			Key:            req.Key,
			Value:          req.Value,
			Effect:         req.Effect,
		}
		if err := db.Create(&t).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		out := dto.TaintResponse{
			ID:             t.ID,
			Key:            t.Key,
			Value:          t.Value,
			Effect:         t.Effect,
			OrganizationID: t.OrganizationID,
			CreatedAt:      t.CreatedAt.UTC().Format(time.RFC3339),
			UpdatedAt:      t.UpdatedAt.UTC().Format(time.RFC3339),
		}
		utils.WriteJSON(w, http.StatusCreated, out)
	}
}

// UpdateTaint godoc
//
//	@ID				UpdateTaint
//	@Summary		Update node taint (org scoped)
//	@Description	Partially update taint fields.
//	@Tags			Taints
//	@Accept			json
//	@Produce		json
//	@Param			X-Org-ID	header		string					false	"Organization UUID"
//	@Param			id			path		string					true	"Node Taint ID (UUID)"
//	@Param			body		body		dto.UpdateTaintRequest	true	"Fields to update"
//	@Success		200			{object}	dto.TaintResponse
//	@Failure		400			{string}	string	"invalid id / invalid json"
//	@Failure		401			{string}	string	"Unauthorized"
//	@Failure		403			{string}	string	"organization required"
//	@Failure		404			{string}	string	"not found"
//	@Failure		500			{string}	string	"update failed"
//	@Router			/taints/{id} [patch]
//	@Security		BearerAuth
//	@Security		OrgKeyAuth
//	@Security		OrgSecretAuth
func UpdateTaint(db *gorm.DB) http.HandlerFunc {
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

		var t models.Taint
		if err := db.Where("id = ? AND organization_id = ?", id, orgID).First(&t).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "not_found", "not_found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		var req dto.UpdateTaintRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_request", "bad request")
			return
		}

		next := t

		if req.Key != nil {
			next.Key = strings.TrimSpace(*req.Key)
		}
		if req.Value != nil {
			next.Value = strings.TrimSpace(*req.Value)
		}
		if req.Effect != nil {
			e := strings.TrimSpace(*req.Effect)
			if e == "" {
				utils.WriteError(w, http.StatusBadRequest, "bad_request", "missing effect")
				return
			}
			if _, ok := allowedEffects[e]; !ok {
				utils.WriteError(w, http.StatusBadRequest, "bad_request", "invalid effect")
				return
			}
			next.Effect = e
		}

		if err := db.Save(&next).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		out := dto.TaintResponse{
			ID:             next.ID,
			Key:            next.Key,
			Value:          next.Value,
			Effect:         next.Effect,
			OrganizationID: next.OrganizationID,
			CreatedAt:      next.CreatedAt.UTC().Format(time.RFC3339),
			UpdatedAt:      next.UpdatedAt.UTC().Format(time.RFC3339),
		}
		utils.WriteJSON(w, http.StatusOK, out)
	}
}

// DeleteTaint godoc
//
//	@ID				DeleteTaint
//	@Summary		Delete taint (org scoped)
//	@Description	Permanently deletes the taint.
//	@Tags			Taints
//	@Accept			json
//	@Produce		json
//	@Param			X-Org-ID	header		string	false	"Organization UUID"
//	@Param			id			path		string	true	"Node Taint ID (UUID)"
//	@Success		204			{string}	string	"No Content"
//	@Failure		400			{string}	string	"invalid id"
//	@Failure		401			{string}	string	"Unauthorized"
//	@Failure		403			{string}	string	"organization required"
//	@Failure		500			{string}	string	"delete failed"
//	@Router			/taints/{id} [delete]
//	@Security		BearerAuth
//	@Security		OrgKeyAuth
//	@Security		OrgSecretAuth
func DeleteTaint(db *gorm.DB) http.HandlerFunc {
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

		var row models.Taint
		if err := db.Where("id = ? AND organization_id = ?", id, orgID).First(&row).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "not_found", "not_found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		if err := db.Delete(&row).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// --- Helpers ---
var allowedEffects = map[string]struct{}{
	"NoSchedule":       {},
	"PreferNoSchedule": {},
	"NoExecute":        {},
}
