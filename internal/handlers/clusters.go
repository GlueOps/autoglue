package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/glueops/autoglue/internal/api/httpmiddleware"
	"github.com/glueops/autoglue/internal/common"
	"github.com/glueops/autoglue/internal/handlers/dto"
	"github.com/glueops/autoglue/internal/models"
	"github.com/glueops/autoglue/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ListClusters godoc
//
//	@ID				ListClusters
//	@Summary		List clusters (org scoped)
//	@Description	Returns clusters for the organization in X-Org-ID. Filter by `q` (name contains).
//	@Tags			Clusters
//	@Produce		json
//	@Param			X-Org-ID	header		string	false	"Organization UUID"
//	@Param			q			query		string	false	"Name contains (case-insensitive)"
//	@Success		200			{array}		dto.ClusterResponse
//	@Failure		401			{string}	string	"Unauthorized"
//	@Failure		403			{string}	string	"organization required"
//	@Failure		500			{string}	string	"failed to list clusters"
//	@Router			/clusters [get]
//	@Security		BearerAuth
//	@Security		OrgKeyAuth
//	@Security		OrgSecretAuth
func ListClusters(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}

		q := db.Where("organization_id = ?", orgID)
		if needle := strings.TrimSpace(r.URL.Query().Get("q")); needle != "" {
			q = q.Where(`name ILIKE ?`, "%"+needle+"%")
		}

		var rows []models.Cluster
		if err := q.
			Preload("CaptainDomain").
			Preload("ControlPlaneRecordSet").
			Preload("AppsLoadBalancer").
			Preload("GlueOpsLoadBalancer").
			Preload("BastionServer").
			Preload("NodePools").
			Preload("NodePools.Labels").
			Preload("NodePools.Annotations").
			Preload("NodePools.Taints").
			Preload("NodePools.Servers").
			Find(&rows).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		out := make([]dto.ClusterResponse, 0, len(rows))
		for _, row := range rows {
			out = append(out, clusterToDTO(row))
		}
		utils.WriteJSON(w, http.StatusOK, out)
	}
}

// GetCluster godoc
//
//	@ID				GetCluster
//	@Summary		Get a single cluster by ID (org scoped)
//	@Description	Returns a cluster with all related resources (domain, record set, load balancers, bastion, node pools).
//	@Tags			Clusters
//	@Produce		json
//	@Param			X-Org-ID	header		string	false	"Organization UUID"
//	@Param			clusterID	path		string	true	"Cluster ID"
//	@Success		200			{object}	dto.ClusterResponse
//	@Failure		400			{string}	string	"bad request"
//	@Failure		401			{string}	string	"Unauthorized"
//	@Failure		403			{string}	string	"organization required"
//	@Failure		404			{string}	string	"cluster not found"
//	@Failure		500			{string}	string	"db error"
//	@Router			/clusters/{clusterID} [get]
//	@Security		BearerAuth
//	@Security		OrgKeyAuth
//	@Security		OrgSecretAuth
func GetCluster(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}

		clusterID, err := uuid.Parse(chi.URLParam(r, "clusterID"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_cluster_id", "invalid cluster id")
			return
		}

		var cluster models.Cluster
		if err := db.
			Where("id = ? AND organization_id = ?", clusterID, orgID).
			Preload("CaptainDomain").
			Preload("ControlPlaneRecordSet").
			Preload("AppsLoadBalancer").
			Preload("GlueOpsLoadBalancer").
			Preload("BastionServer").
			Preload("NodePools").
			Preload("NodePools.Labels").
			Preload("NodePools.Annotations").
			Preload("NodePools.Taints").
			Preload("NodePools.Servers").
			First(&cluster).Error; err != nil {

			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "not_found", "cluster not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		utils.WriteJSON(w, http.StatusOK, clusterToDTO(cluster))
	}
}

// CreateCluster godoc
//
//	@ID				CreateCluster
//	@Summary		Create cluster (org scoped)
//	@Description	Creates a cluster. Status is managed by the system and starts as `pre_pending` for validation.
//	@Tags			Clusters
//	@Accept			json
//	@Produce		json
//	@Param			X-Org-ID	header		string						false	"Organization UUID"
//	@Param			body		body		dto.CreateClusterRequest	true	"payload"
//	@Success		201			{object}	dto.ClusterResponse
//	@Failure		400			{string}	string	"invalid json"
//	@Failure		401			{string}	string	"Unauthorized"
//	@Failure		403			{string}	string	"organization required"
//	@Failure		500			{string}	string	"create failed"
//	@Router			/clusters [post]
//	@Security		BearerAuth
//	@Security		OrgKeyAuth
//	@Security		OrgSecretAuth
func CreateCluster(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}

		var in dto.CreateClusterRequest
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_json", err.Error())
			return
		}

		certificateKey, err := GenerateSecureHex(32)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "internal_error", "failed to generate certificate key")
			return
		}

		randomToken, err := GenerateFormattedToken()
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "internal_error", "failed to generate random token")
			return
		}

		c := models.Cluster{
			OrganizationID: orgID,
			Name:           in.Name,
			Provider:       in.Provider,
			Region:         in.Region,
			Status:         models.ClusterStatusPrePending,
			LastError:      "",
			CertificateKey: certificateKey,
			RandomToken:    randomToken,
		}

		if err := db.Create(&c).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		utils.WriteJSON(w, http.StatusCreated, clusterToDTO(c))
	}
}

// UpdateCluster godoc
//
//	@ID				UpdateCluster
//	@Summary		Update basic cluster details (org scoped)
//	@Description	Updates the cluster name, provider, and/or region. Status is managed by the system.
//	@Tags			Clusters
//	@Accept			json
//	@Produce		json
//	@Param			X-Org-ID	header		string						false	"Organization UUID"
//	@Param			clusterID	path		string						true	"Cluster ID"
//	@Param			body		body		dto.UpdateClusterRequest	true	"payload"
//	@Success		200			{object}	dto.ClusterResponse
//	@Failure		400			{string}	string	"bad request"
//	@Failure		401			{string}	string	"Unauthorized"
//	@Failure		403			{string}	string	"organization required"
//	@Failure		404			{string}	string	"cluster not found"
//	@Failure		500			{string}	string	"db error"
//	@Router			/clusters/{clusterID} [patch]
//	@Security		BearerAuth
//	@Security		OrgKeyAuth
//	@Security		OrgSecretAuth
func UpdateCluster(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}

		clusterID, err := uuid.Parse(chi.URLParam(r, "clusterID"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_cluster_id", "invalid cluster id")
			return
		}

		var in dto.UpdateClusterRequest
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_json", err.Error())
			return
		}

		var cluster models.Cluster
		if err := db.Where("id = ? AND organization_id = ?", clusterID, orgID).First(&cluster).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "not_found", "cluster not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		// Apply only provided fields
		if in.Name != nil {
			cluster.Name = *in.Name
		}
		if in.Provider != nil {
			cluster.Provider = *in.Provider
		}
		if in.Region != nil {
			cluster.Region = *in.Region
		}

		if err := db.Save(&cluster).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		// Any change to the cluster config may require re-validation.
		_ = markClusterNeedsValidation(db, cluster.ID)

		// Preload for a rich response
		if err := db.Preload("CaptainDomain").
			Preload("ControlPlaneRecordSet").
			Preload("AppsLoadBalancer").
			Preload("GlueOpsLoadBalancer").
			Preload("BastionServer").
			Preload("NodePools").
			Preload("NodePools.Labels").
			Preload("NodePools.Annotations").
			Preload("NodePools.Taints").
			Preload("NodePools.Servers").
			First(&cluster, "id = ?", cluster.ID).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		utils.WriteJSON(w, http.StatusOK, clusterToDTO(cluster))
	}
}

// DeleteCluster godoc
//
//	@ID				DeleteCluster
//	@Summary		Delete a cluster (org scoped)
//	@Description	Deletes the cluster. Related resources are cleaned up via DB constraints (e.g. CASCADE).
//	@Tags			Clusters
//	@Produce		json
//	@Param			X-Org-ID	header		string	false	"Organization UUID"
//	@Param			clusterID	path		string	true	"Cluster ID"
//	@Success		204			{string}	string	"deleted"
//	@Failure		400			{string}	string	"bad request"
//	@Failure		401			{string}	string	"Unauthorized"
//	@Failure		403			{string}	string	"organization required"
//	@Failure		404			{string}	string	"cluster not found"
//	@Failure		500			{string}	string	"db error"
//	@Router			/clusters/{clusterID} [delete]
//	@Security		BearerAuth
//	@Security		OrgKeyAuth
//	@Security		OrgSecretAuth
func DeleteCluster(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}

		clusterID, err := uuid.Parse(chi.URLParam(r, "clusterID"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_cluster_id", "invalid cluster id")
			return
		}

		tx := db.Where("id = ? AND organization_id = ?", clusterID, orgID).Delete(&models.Cluster{})
		if tx.Error != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}
		if tx.RowsAffected == 0 {
			utils.WriteError(w, http.StatusNotFound, "not_found", "cluster not found")
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// AttachCaptainDomain godoc
//
//	@ID				AttachCaptainDomain
//	@Summary		Attach a captain domain to a cluster
//	@Description	Sets captain_domain_id on the cluster. Validation of shape happens asynchronously.
//	@Tags			Clusters
//	@Accept			json
//	@Produce		json
//	@Param			X-Org-ID	header		string							false	"Organization UUID"
//	@Param			clusterID	path		string							true	"Cluster ID"
//	@Param			body		body		dto.AttachCaptainDomainRequest	true	"payload"
//	@Success		200			{object}	dto.ClusterResponse
//	@Failure		400			{string}	string	"bad request"
//	@Failure		401			{string}	string	"Unauthorized"
//	@Failure		403			{string}	string	"organization required"
//	@Failure		404			{string}	string	"cluster or domain not found"
//	@Failure		500			{string}	string	"db error"
//	@Router			/clusters/{clusterID}/captain-domain [post]
//	@Security		BearerAuth
//	@Security		OrgKeyAuth
//	@Security		OrgSecretAuth
func AttachCaptainDomain(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}

		clusterID, err := uuid.Parse(chi.URLParam(r, "clusterID"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_cluster_id", "invalid cluster id")
			return
		}

		var in dto.AttachCaptainDomainRequest
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_json", err.Error())
			return
		}

		var cluster models.Cluster
		if err := db.Where("id = ? AND organization_id = ?", clusterID, orgID).First(&cluster).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "not_found", "cluster not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		// Ensure domain exists and belongs to the org
		var domain models.Domain
		if err := db.Where("id = ? AND organization_id = ?", in.DomainID, orgID).First(&domain).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "domain_not_found", "domain not found for organization")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		cluster.CaptainDomainID = &domain.ID
		if err := db.Save(&cluster).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		if err := markClusterNeedsValidation(db, cluster.ID); err != nil {
			// Don't fail the request, just log if you have logging.
		}

		// Preload domain for response
		if err := db.Preload("CaptainDomain").
			Preload("ControlPlaneRecordSet").
			Preload("AppsLoadBalancer").
			Preload("GlueOpsLoadBalancer").
			Preload("BastionServer").
			Preload("NodePools").
			Preload("NodePools.Labels").
			Preload("NodePools.Annotations").
			Preload("NodePools.Taints").
			Preload("NodePools.Servers").
			First(&cluster, "id = ?", cluster.ID).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}
		utils.WriteJSON(w, http.StatusOK, clusterToDTO(cluster))
	}
}

// DetachCaptainDomain godoc
//
//	@ID				DetachCaptainDomain
//	@Summary		Detach the captain domain from a cluster
//	@Description	Clears captain_domain_id on the cluster. This will likely cause the cluster to become incomplete.
//	@Tags			Clusters
//	@Produce		json
//	@Param			X-Org-ID	header		string	false	"Organization UUID"
//	@Param			clusterID	path		string	true	"Cluster ID"
//	@Success		200			{object}	dto.ClusterResponse
//	@Failure		400			{string}	string	"bad request"
//	@Failure		401			{string}	string	"Unauthorized"
//	@Failure		403			{string}	string	"organization required"
//	@Failure		404			{string}	string	"cluster not found"
//	@Failure		500			{string}	string	"db error"
//	@Router			/clusters/{clusterID}/captain-domain [delete]
//	@Security		BearerAuth
//	@Security		OrgKeyAuth
//	@Security		OrgSecretAuth
func DetachCaptainDomain(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}

		clusterID, err := uuid.Parse(chi.URLParam(r, "clusterID"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_cluster_id", "invalid cluster id")
			return
		}

		var cluster models.Cluster
		if err := db.Where("id = ? AND organization_id = ?", clusterID, orgID).First(&cluster).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "not_found", "cluster not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		cluster.CaptainDomainID = nil
		if err := db.Save(&cluster).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		_ = markClusterNeedsValidation(db, cluster.ID)

		if err := db.Preload("CaptainDomain").
			Preload("ControlPlaneRecordSet").
			Preload("AppsLoadBalancer").
			Preload("GlueOpsLoadBalancer").
			Preload("BastionServer").
			Preload("NodePools").
			Preload("NodePools.Labels").
			Preload("NodePools.Annotations").
			Preload("NodePools.Taints").
			Preload("NodePools.Servers").
			First(&cluster, "id = ?", cluster.ID).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		utils.WriteJSON(w, http.StatusOK, clusterToDTO(cluster))
	}
}

// AttachControlPlaneRecordSet godoc
//
//	@ID				AttachControlPlaneRecordSet
//	@Summary		Attach a control plane record set to a cluster
//	@Description	Sets control_plane_record_set_id on the cluster.
//	@Tags			Clusters
//	@Accept			json
//	@Produce		json
//	@Param			X-Org-ID	header		string						false	"Organization UUID"
//	@Param			clusterID	path		string						true	"Cluster ID"
//	@Param			body		body		dto.AttachRecordSetRequest	true	"payload"
//	@Success		200			{object}	dto.ClusterResponse
//	@Failure		400			{string}	string	"bad request"
//	@Failure		401			{string}	string	"Unauthorized"
//	@Failure		403			{string}	string	"organization required"
//	@Failure		404			{string}	string	"cluster or record set not found"
//	@Failure		500			{string}	string	"db error"
//	@Router			/clusters/{clusterID}/control-plane-record-set [post]
//	@Security		BearerAuth
//	@Security		OrgKeyAuth
//	@Security		OrgSecretAuth
func AttachControlPlaneRecordSet(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}

		clusterID, err := uuid.Parse(chi.URLParam(r, "clusterID"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_cluster_id", "invalid cluster id")
			return
		}

		var in dto.AttachRecordSetRequest
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_json", err.Error())
			return
		}

		var cluster models.Cluster
		if err := db.Where("id = ? AND organization_id = ?", clusterID, orgID).First(&cluster).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "not_found", "cluster not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		// record sets are indirectly org-scoped via their domain
		var rs models.RecordSet
		if err := db.
			Joins("JOIN domains d ON d.id = record_sets.domain_id").
			Where("record_sets.id = ? AND d.organization_id = ?", in.RecordSetID, orgID).
			First(&rs).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "recordset_not_found", "record set not found for organization")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		cluster.ControlPlaneRecordSetID = &rs.ID
		if err := db.Save(&cluster).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		_ = markClusterNeedsValidation(db, cluster.ID)

		if err := db.Preload("CaptainDomain").
			Preload("ControlPlaneRecordSet").
			Preload("AppsLoadBalancer").
			Preload("GlueOpsLoadBalancer").
			Preload("BastionServer").
			Preload("NodePools").
			Preload("NodePools.Labels").
			Preload("NodePools.Annotations").
			Preload("NodePools.Taints").
			Preload("NodePools.Servers").
			First(&cluster, "id = ?", cluster.ID).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		utils.WriteJSON(w, http.StatusOK, clusterToDTO(cluster))
	}
}

// DetachControlPlaneRecordSet godoc
//
//	@ID				DetachControlPlaneRecordSet
//	@Summary		Detach the control plane record set from a cluster
//	@Description	Clears control_plane_record_set_id on the cluster.
//	@Tags			Clusters
//	@Produce		json
//	@Param			X-Org-ID	header		string	false	"Organization UUID"
//	@Param			clusterID	path		string	true	"Cluster ID"
//	@Success		200			{object}	dto.ClusterResponse
//	@Failure		400			{string}	string	"bad request"
//	@Failure		401			{string}	string	"Unauthorized"
//	@Failure		403			{string}	string	"organization required"
//	@Failure		404			{string}	string	"cluster not found"
//	@Failure		500			{string}	string	"db error"
//	@Router			/clusters/{clusterID}/control-plane-record-set [delete]
//	@Security		BearerAuth
//	@Security		OrgKeyAuth
//	@Security		OrgSecretAuth
func DetachControlPlaneRecordSet(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}

		clusterID, err := uuid.Parse(chi.URLParam(r, "clusterID"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_cluster_id", "invalid cluster id")
			return
		}

		var cluster models.Cluster
		if err := db.Where("id = ? AND organization_id = ?", clusterID, orgID).First(&cluster).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "not_found", "cluster not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		cluster.ControlPlaneRecordSetID = nil
		if err := db.Save(&cluster).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		_ = markClusterNeedsValidation(db, cluster.ID)

		if err := db.Preload("CaptainDomain").
			Preload("ControlPlaneRecordSet").
			Preload("AppsLoadBalancer").
			Preload("GlueOpsLoadBalancer").
			Preload("BastionServer").
			Preload("NodePools").
			Preload("NodePools.Labels").
			Preload("NodePools.Annotations").
			Preload("NodePools.Taints").
			Preload("NodePools.Servers").
			First(&cluster, "id = ?", cluster.ID).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		utils.WriteJSON(w, http.StatusOK, clusterToDTO(cluster))
	}
}

// AttachAppsLoadBalancer godoc
//
//	@ID				AttachAppsLoadBalancer
//	@Summary		Attach an apps load balancer to a cluster
//	@Description	Sets apps_load_balancer_id on the cluster.
//	@Tags			Clusters
//	@Accept			json
//	@Produce		json
//	@Param			X-Org-ID	header		string							false	"Organization UUID"
//	@Param			clusterID	path		string							true	"Cluster ID"
//	@Param			body		body		dto.AttachLoadBalancerRequest	true	"payload"
//	@Success		200			{object}	dto.ClusterResponse
//	@Failure		400			{string}	string	"bad request"
//	@Failure		401			{string}	string	"Unauthorized"
//	@Failure		403			{string}	string	"organization required"
//	@Failure		404			{string}	string	"cluster or load balancer not found"
//	@Failure		500			{string}	string	"db error"
//	@Router			/clusters/{clusterID}/apps-load-balancer [post]
//	@Security		BearerAuth
//	@Security		OrgKeyAuth
//	@Security		OrgSecretAuth
func AttachAppsLoadBalancer(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}

		clusterID, err := uuid.Parse(chi.URLParam(r, "clusterID"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_cluster_id", "invalid cluster id")
			return
		}

		var in dto.AttachLoadBalancerRequest
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_json", err.Error())
			return
		}

		var cluster models.Cluster
		if err := db.Where("id = ? AND organization_id = ?", clusterID, orgID).First(&cluster).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "not_found", "cluster not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		var lb models.LoadBalancer
		if err := db.Where("id = ? AND organization_id = ?", in.LoadBalancerID, orgID).First(&lb).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "lb_not_found", "load balancer not found for organization")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		cluster.AppsLoadBalancerID = &lb.ID
		if err := db.Save(&cluster).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		_ = markClusterNeedsValidation(db, cluster.ID)

		if err := db.Preload("CaptainDomain").
			Preload("ControlPlaneRecordSet").
			Preload("AppsLoadBalancer").
			Preload("GlueOpsLoadBalancer").
			Preload("BastionServer").
			Preload("NodePools").
			Preload("NodePools.Labels").
			Preload("NodePools.Annotations").
			Preload("NodePools.Taints").
			Preload("NodePools.Servers").
			First(&cluster, "id = ?", cluster.ID).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		utils.WriteJSON(w, http.StatusOK, clusterToDTO(cluster))
	}
}

// DetachAppsLoadBalancer godoc
//
//	@ID				DetachAppsLoadBalancer
//	@Summary		Detach the apps load balancer from a cluster
//	@Description	Clears apps_load_balancer_id on the cluster.
//	@Tags			Clusters
//	@Produce		json
//	@Param			X-Org-ID	header		string	false	"Organization UUID"
//	@Param			clusterID	path		string	true	"Cluster ID"
//	@Success		200			{object}	dto.ClusterResponse
//	@Failure		400			{string}	string	"bad request"
//	@Failure		401			{string}	string	"Unauthorized"
//	@Failure		403			{string}	string	"organization required"
//	@Failure		404			{string}	string	"cluster not found"
//	@Failure		500			{string}	string	"db error"
//	@Router			/clusters/{clusterID}/apps-load-balancer [delete]
//	@Security		BearerAuth
//	@Security		OrgKeyAuth
//	@Security		OrgSecretAuth
func DetachAppsLoadBalancer(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}

		clusterID, err := uuid.Parse(chi.URLParam(r, "clusterID"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_cluster_id", "invalid cluster id")
			return
		}

		var cluster models.Cluster
		if err := db.Where("id = ? AND organization_id = ?", clusterID, orgID).First(&cluster).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "not_found", "cluster not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		cluster.AppsLoadBalancerID = nil
		if err := db.Save(&cluster).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		_ = markClusterNeedsValidation(db, cluster.ID)

		if err := db.Preload("CaptainDomain").
			Preload("ControlPlaneRecordSet").
			Preload("AppsLoadBalancer").
			Preload("GlueOpsLoadBalancer").
			Preload("BastionServer").
			Preload("NodePools").
			Preload("NodePools.Labels").
			Preload("NodePools.Annotations").
			Preload("NodePools.Taints").
			Preload("NodePools.Servers").
			First(&cluster, "id = ?", cluster.ID).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		utils.WriteJSON(w, http.StatusOK, clusterToDTO(cluster))
	}
}

// AttachGlueOpsLoadBalancer godoc
//
//	@ID				AttachGlueOpsLoadBalancer
//	@Summary		Attach a GlueOps/control-plane load balancer to a cluster
//	@Description	Sets glueops_load_balancer_id on the cluster.
//	@Tags			Clusters
//	@Accept			json
//	@Produce		json
//	@Param			X-Org-ID	header		string							false	"Organization UUID"
//	@Param			clusterID	path		string							true	"Cluster ID"
//	@Param			body		body		dto.AttachLoadBalancerRequest	true	"payload"
//	@Success		200			{object}	dto.ClusterResponse
//	@Failure		400			{string}	string	"bad request"
//	@Failure		401			{string}	string	"Unauthorized"
//	@Failure		403			{string}	string	"organization required"
//	@Failure		404			{string}	string	"cluster or load balancer not found"
//	@Failure		500			{string}	string	"db error"
//	@Router			/clusters/{clusterID}/glueops-load-balancer [post]
//	@Security		BearerAuth
//	@Security		OrgKeyAuth
//	@Security		OrgSecretAuth
func AttachGlueOpsLoadBalancer(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}

		clusterID, err := uuid.Parse(chi.URLParam(r, "clusterID"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_cluster_id", "invalid cluster id")
			return
		}

		var in dto.AttachLoadBalancerRequest
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_json", err.Error())
			return
		}

		var cluster models.Cluster
		if err := db.Where("id = ? AND organization_id = ?", clusterID, orgID).First(&cluster).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "not_found", "cluster not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		var lb models.LoadBalancer
		if err := db.Where("id = ? AND organization_id = ?", in.LoadBalancerID, orgID).First(&lb).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "lb_not_found", "load balancer not found for organization")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		cluster.GlueOpsLoadBalancerID = &lb.ID
		if err := db.Save(&cluster).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		_ = markClusterNeedsValidation(db, cluster.ID)

		if err := db.Preload("CaptainDomain").
			Preload("ControlPlaneRecordSet").
			Preload("AppsLoadBalancer").
			Preload("GlueOpsLoadBalancer").
			Preload("BastionServer").
			Preload("NodePools").
			Preload("NodePools.Labels").
			Preload("NodePools.Annotations").
			Preload("NodePools.Taints").
			Preload("NodePools.Servers").
			First(&cluster, "id = ?", cluster.ID).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		utils.WriteJSON(w, http.StatusOK, clusterToDTO(cluster))
	}
}

// DetachGlueOpsLoadBalancer godoc
//
//	@ID				DetachGlueOpsLoadBalancer
//	@Summary		Detach the GlueOps/control-plane load balancer from a cluster
//	@Description	Clears glueops_load_balancer_id on the cluster.
//	@Tags			Clusters
//	@Produce		json
//	@Param			X-Org-ID	header		string	false	"Organization UUID"
//	@Param			clusterID	path		string	true	"Cluster ID"
//	@Success		200			{object}	dto.ClusterResponse
//	@Failure		400			{string}	string	"bad request"
//	@Failure		401			{string}	string	"Unauthorized"
//	@Failure		403			{string}	string	"organization required"
//	@Failure		404			{string}	string	"cluster not found"
//	@Failure		500			{string}	string	"db error"
//	@Router			/clusters/{clusterID}/glueops-load-balancer [delete]
//	@Security		BearerAuth
//	@Security		OrgKeyAuth
//	@Security		OrgSecretAuth
func DetachGlueOpsLoadBalancer(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}

		clusterID, err := uuid.Parse(chi.URLParam(r, "clusterID"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_cluster_id", "invalid cluster id")
			return
		}

		var cluster models.Cluster
		if err := db.Where("id = ? AND organization_id = ?", clusterID, orgID).First(&cluster).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "not_found", "cluster not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		cluster.GlueOpsLoadBalancerID = nil
		if err := db.Save(&cluster).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		_ = markClusterNeedsValidation(db, cluster.ID)

		if err := db.Preload("CaptainDomain").
			Preload("ControlPlaneRecordSet").
			Preload("AppsLoadBalancer").
			Preload("GlueOpsLoadBalancer").
			Preload("BastionServer").
			Preload("NodePools").
			Preload("NodePools.Labels").
			Preload("NodePools.Annotations").
			Preload("NodePools.Taints").
			Preload("NodePools.Servers").
			First(&cluster, "id = ?", cluster.ID).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		utils.WriteJSON(w, http.StatusOK, clusterToDTO(cluster))
	}
}

// AttachBastionServer godoc
//
//	@ID				AttachBastionServer
//	@Summary		Attach a bastion server to a cluster
//	@Description	Sets bastion_server_id on the cluster.
//	@Tags			Clusters
//	@Accept			json
//	@Produce		json
//	@Param			X-Org-ID	header		string						false	"Organization UUID"
//	@Param			clusterID	path		string						true	"Cluster ID"
//	@Param			body		body		dto.AttachBastionRequest	true	"payload"
//	@Success		200			{object}	dto.ClusterResponse
//	@Failure		400			{string}	string	"bad request"
//	@Failure		401			{string}	string	"Unauthorized"
//	@Failure		403			{string}	string	"organization required"
//	@Failure		404			{string}	string	"cluster or server not found"
//	@Failure		500			{string}	string	"db error"
//	@Router			/clusters/{clusterID}/bastion [post]
//	@Security		BearerAuth
//	@Security		OrgKeyAuth
//	@Security		OrgSecretAuth
func AttachBastionServer(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}

		clusterID, err := uuid.Parse(chi.URLParam(r, "clusterID"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_cluster_id", "invalid cluster id")
			return
		}

		var in dto.AttachBastionRequest
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_json", err.Error())
			return
		}

		var cluster models.Cluster
		if err := db.Where("id = ? AND organization_id = ?", clusterID, orgID).First(&cluster).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "not_found", "cluster not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		var server models.Server
		if err := db.Where("id = ? AND organization_id = ?", in.ServerID, orgID).First(&server).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "server_not_found", "server not found for organization")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		cluster.BastionServerID = &server.ID
		if err := db.Save(&cluster).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		_ = markClusterNeedsValidation(db, cluster.ID)

		if err := db.Preload("CaptainDomain").
			Preload("ControlPlaneRecordSet").
			Preload("AppsLoadBalancer").
			Preload("GlueOpsLoadBalancer").
			Preload("BastionServer").
			Preload("NodePools").
			Preload("NodePools.Labels").
			Preload("NodePools.Annotations").
			Preload("NodePools.Taints").
			Preload("NodePools.Servers").
			First(&cluster, "id = ?", cluster.ID).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		utils.WriteJSON(w, http.StatusOK, clusterToDTO(cluster))
	}
}

// DetachBastionServer godoc
//
//	@ID				DetachBastionServer
//	@Summary		Detach the bastion server from a cluster
//	@Description	Clears bastion_server_id on the cluster.
//	@Tags			Clusters
//	@Produce		json
//	@Param			X-Org-ID	header		string	false	"Organization UUID"
//	@Param			clusterID	path		string	true	"Cluster ID"
//	@Success		200			{object}	dto.ClusterResponse
//	@Failure		400			{string}	string	"bad request"
//	@Failure		401			{string}	string	"Unauthorized"
//	@Failure		403			{string}	string	"organization required"
//	@Failure		404			{string}	string	"cluster not found"
//	@Failure		500			{string}	string	"db error"
//	@Router			/clusters/{clusterID}/bastion [delete]
//	@Security		BearerAuth
//	@Security		OrgKeyAuth
//	@Security		OrgSecretAuth
func DetachBastionServer(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}

		clusterID, err := uuid.Parse(chi.URLParam(r, "clusterID"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_cluster_id", "invalid cluster id")
			return
		}

		var cluster models.Cluster
		if err := db.Where("id = ? AND organization_id = ?", clusterID, orgID).First(&cluster).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "not_found", "cluster not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		cluster.BastionServerID = nil
		if err := db.Save(&cluster).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		_ = markClusterNeedsValidation(db, cluster.ID)

		if err := db.Preload("CaptainDomain").
			Preload("ControlPlaneRecordSet").
			Preload("AppsLoadBalancer").
			Preload("GlueOpsLoadBalancer").
			Preload("BastionServer").
			Preload("NodePools").
			Preload("NodePools.Labels").
			Preload("NodePools.Annotations").
			Preload("NodePools.Taints").
			Preload("NodePools.Servers").
			First(&cluster, "id = ?", cluster.ID).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		utils.WriteJSON(w, http.StatusOK, clusterToDTO(cluster))
	}
}

// SetClusterKubeconfig godoc
//
//	@ID				SetClusterKubeconfig
//	@Summary		Set (or replace) the kubeconfig for a cluster
//	@Description	Stores the kubeconfig encrypted per organization. The kubeconfig is never returned in responses.
//	@Tags			Clusters
//	@Accept			json
//	@Produce		json
//	@Param			X-Org-ID	header		string						false	"Organization UUID"
//	@Param			clusterID	path		string						true	"Cluster ID"
//	@Param			body		body		dto.SetKubeconfigRequest	true	"payload"
//	@Success		200			{object}	dto.ClusterResponse
//	@Failure		400			{string}	string	"bad request"
//	@Failure		401			{string}	string	"Unauthorized"
//	@Failure		403			{string}	string	"organization required"
//	@Failure		404			{string}	string	"cluster not found"
//	@Failure		500			{string}	string	"db error"
//	@Router			/clusters/{clusterID}/kubeconfig [post]
//	@Security		BearerAuth
//	@Security		OrgKeyAuth
//	@Security		OrgSecretAuth
func SetClusterKubeconfig(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}

		clusterID, err := uuid.Parse(chi.URLParam(r, "clusterID"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_cluster_id", "invalid cluster id")
			return
		}

		var in dto.SetKubeconfigRequest
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_json", err.Error())
			return
		}

		var cluster models.Cluster
		if err := db.Where("id = ? AND organization_id = ?", clusterID, orgID).First(&cluster).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "not_found", "cluster not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		ct, iv, tag, err := utils.EncryptForOrg(orgID, []byte(in.Kubeconfig), db)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "encryption_error", "failed to encrypt kubeconfig")
			return
		}

		cluster.EncryptedKubeconfig = ct
		cluster.KubeIV = iv
		cluster.KubeTag = tag

		if err := db.Save(&cluster).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		_ = markClusterNeedsValidation(db, cluster.ID)

		if err := db.Preload("CaptainDomain").
			Preload("ControlPlaneRecordSet").
			Preload("AppsLoadBalancer").
			Preload("GlueOpsLoadBalancer").
			Preload("BastionServer").
			Preload("NodePools").
			Preload("NodePools.Labels").
			Preload("NodePools.Annotations").
			Preload("NodePools.Taints").
			Preload("NodePools.Servers").
			First(&cluster, "id = ?", cluster.ID).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		utils.WriteJSON(w, http.StatusOK, clusterToDTO(cluster))
	}
}

// ClearClusterKubeconfig godoc
//
//	@ID				ClearClusterKubeconfig
//	@Summary		Clear the kubeconfig for a cluster
//	@Description	Removes the encrypted kubeconfig, IV, and tag from the cluster record.
//	@Tags			Clusters
//	@Produce		json
//	@Param			X-Org-ID	header		string	false	"Organization UUID"
//	@Param			clusterID	path		string	true	"Cluster ID"
//	@Success		200			{object}	dto.ClusterResponse
//	@Failure		400			{string}	string	"bad request"
//	@Failure		401			{string}	string	"Unauthorized"
//	@Failure		403			{string}	string	"organization required"
//	@Failure		404			{string}	string	"cluster not found"
//	@Failure		500			{string}	string	"db error"
//	@Router			/clusters/{clusterID}/kubeconfig [delete]
//	@Security		BearerAuth
//	@Security		OrgKeyAuth
//	@Security		OrgSecretAuth
func ClearClusterKubeconfig(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}

		clusterID, err := uuid.Parse(chi.URLParam(r, "clusterID"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_cluster_id", "invalid cluster id")
			return
		}

		var cluster models.Cluster
		if err := db.Where("id = ? AND organization_id = ?", clusterID, orgID).First(&cluster).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "not_found", "cluster not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		cluster.EncryptedKubeconfig = ""
		cluster.KubeIV = ""
		cluster.KubeTag = ""

		if err := db.Save(&cluster).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		_ = markClusterNeedsValidation(db, cluster.ID)

		if err := db.Preload("CaptainDomain").
			Preload("ControlPlaneRecordSet").
			Preload("AppsLoadBalancer").
			Preload("GlueOpsLoadBalancer").
			Preload("BastionServer").
			Preload("NodePools").
			Preload("NodePools.Labels").
			Preload("NodePools.Annotations").
			Preload("NodePools.Taints").
			Preload("NodePools.Servers").
			First(&cluster, "id = ?", cluster.ID).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		utils.WriteJSON(w, http.StatusOK, clusterToDTO(cluster))
	}
}

// AttachNodePool godoc
//
//	@ID				AttachNodePool
//	@Summary		Attach a node pool to a cluster
//	@Description	Adds an entry in the cluster_node_pools join table.
//	@Tags			Clusters
//	@Accept			json
//	@Produce		json
//	@Param			X-Org-ID	header		string						false	"Organization UUID"
//	@Param			clusterID	path		string						true	"Cluster ID"
//	@Param			body		body		dto.AttachNodePoolRequest	true	"payload"
//	@Success		200			{object}	dto.ClusterResponse
//	@Failure		400			{string}	string	"bad request"
//	@Failure		401			{string}	string	"Unauthorized"
//	@Failure		403			{string}	string	"organization required"
//	@Failure		404			{string}	string	"cluster or node pool not found"
//	@Failure		500			{string}	string	"db error"
//	@Router			/clusters/{clusterID}/node-pools [post]
//	@Security		BearerAuth
//	@Security		OrgKeyAuth
//	@Security		OrgSecretAuth
func AttachNodePool(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}

		clusterID, err := uuid.Parse(chi.URLParam(r, "clusterID"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_cluster_id", "invalid cluster id")
			return
		}

		var in dto.AttachNodePoolRequest
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_json", err.Error())
			return
		}

		// Load cluster (org scoped)
		var cluster models.Cluster
		if err := db.
			Where("id = ? AND organization_id = ?", clusterID, orgID).
			First(&cluster).Error; err != nil {

			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "not_found", "cluster not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		// Load node pool (org scoped)
		var np models.NodePool
		if err := db.
			Where("id = ? AND organization_id = ?", in.NodePoolID, orgID).
			First(&np).Error; err != nil {

			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "nodepool_not_found", "node pool not found for organization")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		// Create association in join table
		if err := db.Model(&cluster).Association("NodePools").Append(&np); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "failed to attach node pool")
			return
		}

		_ = markClusterNeedsValidation(db, cluster.ID)

		// Reload for rich response
		if err := db.
			Preload("CaptainDomain").
			Preload("ControlPlaneRecordSet").
			Preload("AppsLoadBalancer").
			Preload("GlueOpsLoadBalancer").
			Preload("BastionServer").
			Preload("NodePools").
			Preload("NodePools.Labels").
			Preload("NodePools.Annotations").
			Preload("NodePools.Taints").
			Preload("NodePools.Servers").
			First(&cluster, "id = ?", cluster.ID).Error; err != nil {

			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		utils.WriteJSON(w, http.StatusOK, clusterToDTO(cluster))
	}
}

// DetachNodePool godoc
//
//	@ID				DetachNodePool
//	@Summary		Detach a node pool from a cluster
//	@Description	Removes an entry from the cluster_node_pools join table.
//	@Tags			Clusters
//	@Produce		json
//	@Param			X-Org-ID	header		string	false	"Organization UUID"
//	@Param			clusterID	path		string	true	"Cluster ID"
//	@Param			nodePoolID	path		string	true	"Node Pool ID"
//	@Success		200			{object}	dto.ClusterResponse
//	@Failure		400			{string}	string	"bad request"
//	@Failure		401			{string}	string	"Unauthorized"
//	@Failure		403			{string}	string	"organization required"
//	@Failure		404			{string}	string	"cluster or node pool not found"
//	@Failure		500			{string}	string	"db error"
//	@Router			/clusters/{clusterID}/node-pools/{nodePoolID} [delete]
//	@Security		BearerAuth
//	@Security		OrgKeyAuth
//	@Security		OrgSecretAuth
func DetachNodePool(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}

		clusterID, err := uuid.Parse(chi.URLParam(r, "clusterID"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_cluster_id", "invalid cluster id")
			return
		}

		nodePoolID, err := uuid.Parse(chi.URLParam(r, "nodePoolID"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_nodepool_id", "invalid node pool id")
			return
		}

		var cluster models.Cluster
		if err := db.
			Where("id = ? AND organization_id = ?", clusterID, orgID).
			First(&cluster).Error; err != nil {

			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "not_found", "cluster not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		var np models.NodePool
		if err := db.
			Where("id = ? AND organization_id = ?", nodePoolID, orgID).
			First(&np).Error; err != nil {

			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "nodepool_not_found", "node pool not found for organization")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		if err := db.Model(&cluster).Association("NodePools").Delete(&np); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "failed to detach node pool")
			return
		}

		_ = markClusterNeedsValidation(db, cluster.ID)

		if err := db.
			Preload("CaptainDomain").
			Preload("ControlPlaneRecordSet").
			Preload("AppsLoadBalancer").
			Preload("GlueOpsLoadBalancer").
			Preload("BastionServer").
			Preload("NodePools").
			Preload("NodePools.Labels").
			Preload("NodePools.Annotations").
			Preload("NodePools.Taints").
			Preload("NodePools.Servers").
			First(&cluster, "id = ?", cluster.ID).Error; err != nil {

			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		utils.WriteJSON(w, http.StatusOK, clusterToDTO(cluster))
	}
}

// -- Helpers

func clusterToDTO(c models.Cluster) dto.ClusterResponse {
	var bastion *dto.ServerResponse
	if c.BastionServer != nil {
		b := serverToDTO(*c.BastionServer)
		bastion = &b
	}

	var captainDomain *dto.DomainResponse
	if c.CaptainDomainID != nil && c.CaptainDomain.ID != uuid.Nil {
		dr := domainToDTO(c.CaptainDomain)
		captainDomain = &dr
	}

	var controlPlane *dto.RecordSetResponse
	if c.ControlPlaneRecordSet != nil {
		rr := recordSetToDTO(*c.ControlPlaneRecordSet)
		controlPlane = &rr
	}

	var appsLB *dto.LoadBalancerResponse
	if c.AppsLoadBalancer != nil {
		lr := loadBalancerToDTO(*c.AppsLoadBalancer)
		appsLB = &lr
	}

	var glueOpsLB *dto.LoadBalancerResponse
	if c.GlueOpsLoadBalancer != nil {
		lr := loadBalancerToDTO(*c.GlueOpsLoadBalancer)
		glueOpsLB = &lr
	}

	nps := make([]dto.NodePoolResponse, 0, len(c.NodePools))
	for _, np := range c.NodePools {
		nps = append(nps, nodePoolToDTO(np))
	}

	return dto.ClusterResponse{
		ID:                    c.ID,
		Name:                  c.Name,
		CaptainDomain:         captainDomain,
		ControlPlaneRecordSet: controlPlane,
		AppsLoadBalancer:      appsLB,
		GlueOpsLoadBalancer:   glueOpsLB,
		BastionServer:         bastion,
		Provider:              c.Provider,
		Region:                c.Region,
		Status:                c.Status,
		LastError:             c.LastError,
		RandomToken:           c.RandomToken,
		CertificateKey:        c.CertificateKey,
		NodePools:             nps,
		CreatedAt:             c.CreatedAt,
		UpdatedAt:             c.UpdatedAt,
	}
}

func nodePoolToDTO(np models.NodePool) dto.NodePoolResponse {
	labels := make([]dto.LabelResponse, 0, len(np.Labels))
	for _, l := range np.Labels {
		labels = append(labels, dto.LabelResponse{
			Key:   l.Key,
			Value: l.Value,
		})
	}

	annotations := make([]dto.AnnotationResponse, 0, len(np.Annotations))
	for _, a := range np.Annotations {
		annotations = append(annotations, dto.AnnotationResponse{
			Key:   a.Key,
			Value: a.Value,
		})
	}

	taints := make([]dto.TaintResponse, 0, len(np.Taints))
	for _, t := range np.Taints {
		taints = append(taints, dto.TaintResponse{
			Key:    t.Key,
			Value:  t.Value,
			Effect: t.Effect,
		})
	}

	servers := make([]dto.ServerResponse, 0, len(np.Servers))
	for _, s := range np.Servers {
		servers = append(servers, serverToDTO(s))
	}

	return dto.NodePoolResponse{
		AuditFields: common.AuditFields{
			ID:             np.ID,
			OrganizationID: np.OrganizationID,
			CreatedAt:      np.CreatedAt,
			UpdatedAt:      np.UpdatedAt,
		},
		Name:        np.Name,
		Role:        dto.NodeRole(np.Role),
		Labels:      labels,
		Annotations: annotations,
		Taints:      taints,
		Servers:     servers,
	}
}

func serverToDTO(s models.Server) dto.ServerResponse {
	return dto.ServerResponse{
		ID:               s.ID,
		Hostname:         s.Hostname,
		PrivateIPAddress: s.PrivateIPAddress,
		PublicIPAddress:  s.PublicIPAddress,
		Role:             s.Role,
		Status:           s.Status,
		SSHUser:          s.SSHUser,
		SshKeyID:         s.SshKeyID,
		CreatedAt:        s.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:        s.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func domainToDTO(d models.Domain) dto.DomainResponse {
	return dto.DomainResponse{
		ID:             d.ID.String(),
		OrganizationID: d.OrganizationID.String(),
		DomainName:     d.DomainName,
		ZoneID:         d.ZoneID,
		Status:         d.Status,
		LastError:      d.LastError,
		CredentialID:   d.CredentialID.String(),
		CreatedAt:      d.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:      d.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func recordSetToDTO(rs models.RecordSet) dto.RecordSetResponse {
	return dto.RecordSetResponse{
		ID:          rs.ID.String(),
		DomainID:    rs.DomainID.String(),
		Name:        rs.Name,
		Type:        rs.Type,
		TTL:         rs.TTL,
		Values:      []byte(rs.Values),
		Fingerprint: rs.Fingerprint,
		Status:      rs.Status,
		Owner:       rs.Owner,
		LastError:   rs.LastError,
		CreatedAt:   rs.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:   rs.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func loadBalancerToDTO(lb models.LoadBalancer) dto.LoadBalancerResponse {
	return dto.LoadBalancerResponse{
		ID:               lb.ID,
		OrganizationID:   lb.OrganizationID,
		Name:             lb.Name,
		Kind:             lb.Kind,
		PublicIPAddress:  lb.PublicIPAddress,
		PrivateIPAddress: lb.PrivateIPAddress,
		CreatedAt:        lb.CreatedAt,
		UpdatedAt:        lb.UpdatedAt,
	}
}

func GenerateSecureHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

func GenerateFormattedToken() (string, error) {
	part1, err := GenerateSecureHex(3)
	if err != nil {
		return "", fmt.Errorf("failed to generate token part 1: %w", err)
	}
	part2, err := GenerateSecureHex(8)
	if err != nil {
		return "", fmt.Errorf("failed to generate token part 2: %w", err)
	}
	return fmt.Sprintf("%s.%s", part1, part2), nil
}

func markClusterNeedsValidation(db *gorm.DB, clusterID uuid.UUID) error {
	return db.Model(&models.Cluster{}).Where("id = ?", clusterID).Updates(map[string]any{
		"status":     models.ClusterStatusPrePending,
		"last_error": "",
	}).Error
}
