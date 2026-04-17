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

// ListClusterMetadata godoc
//
//	@ID				ListClusterMetadata
//	@Summary		List metadata for a cluster (org scoped)
//	@Description	Returns all metadata key-value pairs attached to the cluster.
//	@Tags			Cluster Metadata
//	@Produce		json
//	@Param			X-Org-ID	header		string	false	"Organization UUID"
//	@Param			clusterID	path		string	true	"Cluster ID"
//	@Success		200			{array}		dto.ClusterMetadataResponse
//	@Failure		400			{string}	string	"invalid cluster id"
//	@Failure		401			{string}	string	"Unauthorized"
//	@Failure		403			{string}	string	"organization required"
//	@Failure		500			{string}	string	"db error"
//	@Router			/clusters/{clusterID}/metadata [get]
//	@Security		BearerAuth
//	@Security		OrgKeyAuth
//	@Security		OrgSecretAuth
func ListClusterMetadata(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}

		clusterID, err := uuid.Parse(chi.URLParam(r, "clusterID"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_request", "invalid cluster id")
			return
		}

		// Ensure cluster belongs to org
		var cluster models.Cluster
		if err := db.Where("id = ? AND organization_id = ?", clusterID, orgID).First(&cluster).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "not_found", "cluster not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		var rows []models.ClusterMetadata
		if err := db.Where("cluster_id = ?", clusterID).Order("created_at ASC").Find(&rows).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		out := make([]dto.ClusterMetadataResponse, 0, len(rows))
		for _, m := range rows {
			out = append(out, clusterMetadataToDTO(m))
		}
		utils.WriteJSON(w, http.StatusOK, out)
	}
}

// GetClusterMetadata godoc
//
//	@ID				GetClusterMetadata
//	@Summary		Get a single cluster metadata entry (org scoped)
//	@Description	Returns one metadata key-value pair by ID.
//	@Tags			Cluster Metadata
//	@Produce		json
//	@Param			X-Org-ID	header		string	false	"Organization UUID"
//	@Param			clusterID	path		string	true	"Cluster ID"
//	@Param			metadataID	path		string	true	"Metadata ID (UUID)"
//	@Success		200			{object}	dto.ClusterMetadataResponse
//	@Failure		400			{string}	string	"invalid id"
//	@Failure		401			{string}	string	"Unauthorized"
//	@Failure		403			{string}	string	"organization required"
//	@Failure		404			{string}	string	"not found"
//	@Failure		500			{string}	string	"db error"
//	@Router			/clusters/{clusterID}/metadata/{metadataID} [get]
//	@Security		BearerAuth
//	@Security		OrgKeyAuth
//	@Security		OrgSecretAuth
func GetClusterMetadata(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}

		clusterID, err := uuid.Parse(chi.URLParam(r, "clusterID"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_request", "invalid cluster id")
			return
		}

		metadataID, err := uuid.Parse(chi.URLParam(r, "metadataID"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_request", "invalid metadata id")
			return
		}

		// Ensure cluster belongs to org
		var cluster models.Cluster
		if err := db.Where("id = ? AND organization_id = ?", clusterID, orgID).First(&cluster).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "not_found", "cluster not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		var m models.ClusterMetadata
		if err := db.Where("id = ? AND cluster_id = ?", metadataID, clusterID).First(&m).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "not_found", "not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		utils.WriteJSON(w, http.StatusOK, clusterMetadataToDTO(m))
	}
}

// CreateClusterMetadata godoc
//
//	@ID				CreateClusterMetadata
//	@Summary		Create cluster metadata (org scoped)
//	@Description	Adds a new key-value metadata entry to a cluster. Keys are forced to lowercase; values preserve case.
//	@Tags			Cluster Metadata
//	@Accept			json
//	@Produce		json
//	@Param			X-Org-ID	header		string								false	"Organization UUID"
//	@Param			clusterID	path		string								true	"Cluster ID"
//	@Param			body		body		dto.CreateClusterMetadataRequest	true	"Key-value pair"
//	@Success		201			{object}	dto.ClusterMetadataResponse
//	@Failure		400			{string}	string	"invalid id / invalid json / missing key"
//	@Failure		401			{string}	string	"Unauthorized"
//	@Failure		403			{string}	string	"organization required"
//	@Failure		404			{string}	string	"cluster not found"
//	@Failure		500			{string}	string	"db error"
//	@Router			/clusters/{clusterID}/metadata [post]
//	@Security		BearerAuth
//	@Security		OrgKeyAuth
//	@Security		OrgSecretAuth
func CreateClusterMetadata(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}

		clusterID, err := uuid.Parse(chi.URLParam(r, "clusterID"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_request", "invalid cluster id")
			return
		}

		// Ensure cluster belongs to org
		var cluster models.Cluster
		if err := db.Where("id = ? AND organization_id = ?", clusterID, orgID).First(&cluster).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "not_found", "cluster not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		var req dto.CreateClusterMetadataRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_request", "invalid json")
			return
		}

		key := strings.ToLower(strings.TrimSpace(req.Key))
		if key == "" {
			utils.WriteError(w, http.StatusBadRequest, "bad_request", "key is required")
			return
		}
		value := strings.TrimSpace(req.Value)
		if value == "" {
			utils.WriteError(w, http.StatusBadRequest, "bad_request", "value is required")
			return
		}

		m := models.ClusterMetadata{
			ClusterID: clusterID,
			Key:       key,
			Value:     value, // value case preserved
		}
		m.OrganizationID = orgID

		if err := db.Create(&m).Error; err != nil {
			if isUniqueConstraintViolation(err) {
				utils.WriteError(w, http.StatusConflict, "conflict", "metadata key already exists for this cluster")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		utils.WriteJSON(w, http.StatusCreated, clusterMetadataToDTO(m))
	}
}

// UpdateClusterMetadata godoc
//
//	@ID				UpdateClusterMetadata
//	@Summary		Update cluster metadata (org scoped)
//	@Description	Partially updates a metadata entry. Key is forced to lowercase if provided; value case is preserved.
//	@Tags			Cluster Metadata
//	@Accept			json
//	@Produce		json
//	@Param			X-Org-ID	header		string								false	"Organization UUID"
//	@Param			clusterID	path		string								true	"Cluster ID"
//	@Param			metadataID	path		string								true	"Metadata ID (UUID)"
//	@Param			body		body		dto.UpdateClusterMetadataRequest	true	"Fields to update"
//	@Success		200			{object}	dto.ClusterMetadataResponse
//	@Failure		400			{string}	string	"invalid id / invalid json"
//	@Failure		401			{string}	string	"Unauthorized"
//	@Failure		403			{string}	string	"organization required"
//	@Failure		404			{string}	string	"not found"
//	@Failure		500			{string}	string	"db error"
//	@Router			/clusters/{clusterID}/metadata/{metadataID} [patch]
//	@Security		BearerAuth
//	@Security		OrgKeyAuth
//	@Security		OrgSecretAuth
func UpdateClusterMetadata(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}

		clusterID, err := uuid.Parse(chi.URLParam(r, "clusterID"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_request", "invalid cluster id")
			return
		}

		metadataID, err := uuid.Parse(chi.URLParam(r, "metadataID"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_request", "invalid metadata id")
			return
		}

		// Ensure cluster belongs to org
		var cluster models.Cluster
		if err := db.Where("id = ? AND organization_id = ?", clusterID, orgID).First(&cluster).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "not_found", "cluster not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		var m models.ClusterMetadata
		if err := db.Where("id = ? AND cluster_id = ?", metadataID, clusterID).First(&m).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "not_found", "not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		var req dto.UpdateClusterMetadataRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_request", "invalid json")
			return
		}

		if req.Key != nil {
			normalizedKey := strings.TrimSpace(*req.Key)
			if normalizedKey == "" {
				utils.WriteError(w, http.StatusBadRequest, "bad_request", "key cannot be empty")
				return
			}
			m.Key = strings.ToLower(normalizedKey)
		}
		if req.Value != nil {
			value := strings.TrimSpace(*req.Value)
			if value == "" {
				utils.WriteError(w, http.StatusBadRequest, "bad_request", "value cannot be empty")
				return
			}
			m.Value = value // value case preserved
		}

		if err := db.Save(&m).Error; err != nil {
			if isUniqueConstraintViolation(err) {
				utils.WriteError(w, http.StatusConflict, "conflict", "metadata key already exists for this cluster")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		utils.WriteJSON(w, http.StatusOK, clusterMetadataToDTO(m))
	}
}

// DeleteClusterMetadata godoc
//
//	@ID				DeleteClusterMetadata
//	@Summary		Delete cluster metadata (org scoped)
//	@Description	Permanently deletes a metadata entry from a cluster.
//	@Tags			Cluster Metadata
//	@Param			X-Org-ID	header	string	false	"Organization UUID"
//	@Param			clusterID	path	string	true	"Cluster ID"
//	@Param			metadataID	path	string	true	"Metadata ID (UUID)"
//	@Success		204			"No Content"
//	@Failure		400			{string}	string	"invalid id"
//	@Failure		401			{string}	string	"Unauthorized"
//	@Failure		403			{string}	string	"organization required"
//	@Failure		404			{string}	string	"cluster not found"
//	@Failure		500			{string}	string	"db error"
//	@Router			/clusters/{clusterID}/metadata/{metadataID} [delete]
//	@Security		BearerAuth
//	@Security		OrgKeyAuth
//	@Security		OrgSecretAuth
func DeleteClusterMetadata(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}

		clusterID, err := uuid.Parse(chi.URLParam(r, "clusterID"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_request", "invalid cluster id")
			return
		}

		metadataID, err := uuid.Parse(chi.URLParam(r, "metadataID"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_request", "invalid metadata id")
			return
		}

		// Ensure cluster belongs to org
		var cluster models.Cluster
		if err := db.Where("id = ? AND organization_id = ?", clusterID, orgID).First(&cluster).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "not_found", "cluster not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		if err := db.Where("id = ? AND cluster_id = ?", metadataID, clusterID).Delete(&models.ClusterMetadata{}).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func clusterMetadataToDTO(m models.ClusterMetadata) dto.ClusterMetadataResponse {
	return dto.ClusterMetadataResponse{
		AuditFields: m.AuditFields,
		ClusterID:   m.ClusterID.String(),
		Key:         m.Key,
		Value:       m.Value,
	}
}

func isUniqueConstraintViolation(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "duplicate key value") || strings.Contains(msg, "unique constraint")
}
