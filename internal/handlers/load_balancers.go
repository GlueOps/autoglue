package handlers

import (
	"encoding/json"
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

// ListLoadBalancers godoc
//
//	@ID				ListLoadBalancers
//	@Summary		List load balancers (org scoped)
//	@Description	Returns load balancers for the organization in X-Org-ID.
//	@Tags			LoadBalancers
//	@Produce		json
//	@Param			X-Org-ID	header		string	false	"Organization UUID"
//	@Success		200			{array}		dto.LoadBalancerResponse
//	@Failure		401			{string}	string	"Unauthorized"
//	@Failure		403			{string}	string	"organization required"
//	@Failure		500			{string}	string	"failed to list clusters"
//	@Router			/load-balancers [get]
//	@Security		BearerAuth
//	@Security		OrgKeyAuth
//	@Security		OrgSecretAuth
func ListLoadBalancers(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}

		var rows []models.LoadBalancer
		if err := db.Where("organization_id = ?", orgID).Find(&rows).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}
		out := make([]dto.LoadBalancerResponse, 0, len(rows))
		for _, row := range rows {
			out = append(out, loadBalancerOut(&row))
		}
		utils.WriteJSON(w, http.StatusOK, out)
	}
}

// GetLoadBalancer godoc
//
//	@ID				GetLoadBalancers
//	@Summary		Get a load balancer (org scoped)
//	@Description	Returns load balancer for the organization in X-Org-ID.
//	@Tags			LoadBalancers
//	@Produce		json
//	@Param			X-Org-ID	header		string	false	"Organization UUID"
//	@Param			id			path		string	true	"LoadBalancer ID (UUID)"
//	@Success		200			{array}		dto.LoadBalancerResponse
//	@Failure		401			{string}	string	"Unauthorized"
//	@Failure		403			{string}	string	"organization required"
//	@Failure		500			{string}	string	"failed to list clusters"
//	@Router			/load-balancers/{id} [get]
//	@Security		BearerAuth
//	@Security		OrgKeyAuth
//	@Security		OrgSecretAuth
func GetLoadBalancer(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_id", "invalid UUID")
			return
		}

		var row models.LoadBalancer
		if err := db.Where("id = ? AND organization_id = ?", id, orgID).First(&row).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "not_found", "load balancer not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}
		out := loadBalancerOut(&row)
		utils.WriteJSON(w, http.StatusOK, out)
	}
}

// CreateLoadBalancer godoc
//
//	@ID			CreateLoadBalancer
//	@Summary	Create a load balancer
//	@Tags		LoadBalancers
//	@Accept		json
//	@Produce	json
//	@Param		X-Org-ID	header		string							false	"Organization UUID"
//	@Param		body		body		dto.CreateLoadBalancerRequest	true	"Record set payload"
//	@Success	201			{object}	dto.LoadBalancerResponse
//	@Failure	400			{string}	string	"validation error"
//	@Failure	403			{string}	string	"organization required"
//	@Failure	404			{string}	string	"domain not found"
//	@Router		/load-balancers [post]
//	@Security	BearerAuth
//	@Security	OrgKeyAuth
//	@Security	OrgSecretAuth
func CreateLoadBalancer(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}

		var in dto.CreateLoadBalancerRequest
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_json", err.Error())
			return
		}

		if strings.ToLower(in.Kind) != "glueops" || strings.ToLower(in.Kind) != "public" {
			utils.WriteError(w, http.StatusBadRequest, "bad_kind", "invalid kind only 'glueops' or 'public'")
			return
		}

		row := &models.LoadBalancer{
			OrganizationID:   orgID,
			Name:             in.Name,
			Kind:             strings.ToLower(in.Kind),
			PublicIPAddress:  in.PublicIPAddress,
			PrivateIPAddress: in.PrivateIPAddress,
		}
		if err := db.Create(row).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", err.Error())
			return
		}
		utils.WriteJSON(w, http.StatusCreated, loadBalancerOut(row))
	}
}

// UpdateLoadBalancer godoc
//
//	@ID			UpdateLoadBalancer
//	@Summary	Update a load balancer (org scoped)
//	@Tags		LoadBalancers
//	@Accept		json
//	@Produce	json
//	@Param		X-Org-ID	header		string							false	"Organization UUID"
//	@Param		id			path		string							true	"Load Balancer ID (UUID)"
//	@Param		body		body		dto.UpdateLoadBalancerRequest	true	"Fields to update"
//	@Success	200			{object}	dto.LoadBalancerResponse
//	@Failure	400			{string}	string	"validation error"
//	@Failure	403			{string}	string	"organization required"
//	@Failure	404			{string}	string	"not found"
//	@Router		/load-balancers/{id} [patch]
//	@Security	BearerAuth
//	@Security	OrgKeyAuth
//	@Security	OrgSecretAuth
func UpdateLoadBalancer(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_id", "invalid UUID")
			return
		}

		row := &models.LoadBalancer{}
		if err := db.Where("organization_id = ? AND id = ?", orgID, id).First(&row).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "not_found", "load balancer not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", err.Error())
			return
		}

		var in dto.UpdateLoadBalancerRequest
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_json", err.Error())
			return
		}
		if in.Name != nil {
			row.Name = *in.Name
		}
		if in.Kind != nil {
			if strings.ToLower(*in.Kind) != "glueops" || strings.ToLower(*in.Kind) != "public" {
				utils.WriteError(w, http.StatusBadRequest, "bad_kind", "invalid kind only 'glueops' or 'public'")
				return
			}
			row.Kind = strings.ToLower(*in.Kind)
		}
		if in.PublicIPAddress != nil {
			row.PublicIPAddress = *in.PublicIPAddress
		}
		if in.PrivateIPAddress != nil {
			row.PrivateIPAddress = *in.PrivateIPAddress
		}
		if err := db.Save(row).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", err.Error())
			return
		}
		utils.WriteJSON(w, http.StatusOK, loadBalancerOut(row))

	}
}

// DeleteLoadBalancer godoc
//
//	@ID			DeleteLoadBalancer
//	@Summary	Delete a load balancer
//	@Tags		LoadBalancers
//	@Produce	json
//	@Param		X-Org-ID	header	string	false	"Organization UUID"
//	@Param		id			path	string	true	"Load Balancer ID (UUID)"
//	@Success	204
//	@Failure	403	{string}	string	"organization required"
//	@Failure	404	{string}	string	"not found"
//	@Router		/load-balancers/{id} [delete]
//	@Security	BearerAuth
//	@Security	OrgKeyAuth
//	@Security	OrgSecretAuth
func DeleteLoadBalancer(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_id", "invalid UUID")
			return
		}

		row := &models.LoadBalancer{}
		if err := db.Where("organization_id = ? AND id = ?", orgID, id).First(&row).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "not_found", "load balancer not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", err.Error())
			return
		}

		if err := db.Delete(row).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", err.Error())
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// ---------- Out mappers ----------

func loadBalancerOut(m *models.LoadBalancer) dto.LoadBalancerResponse {
	return dto.LoadBalancerResponse{
		ID:               m.ID,
		OrganizationID:   m.OrganizationID,
		Name:             m.Name,
		Kind:             m.Kind,
		PublicIPAddress:  m.PublicIPAddress,
		PrivateIPAddress: m.PrivateIPAddress,
		CreatedAt:        m.CreatedAt.UTC(),
		UpdatedAt:        m.UpdatedAt.UTC(),
	}
}
