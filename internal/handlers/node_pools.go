package handlers

import (
	"encoding/json"
	"errors"
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

// -- Node Pools Core

// ListNodePools godoc
//
//	@ID				ListNodePools
//	@Summary		List node pools (org scoped)
//	@Description	Returns node pools for the organization in X-Org-ID.
//	@Tags			NodePools
//	@Accept			json
//	@Produce		json
//	@Param			X-Org-ID	header	string	false	"Organization UUID"
//	@Param			q			query	string	false	"Name contains (case-insensitive)"
//	@Success		200	{array}		dto.NodePoolResponse
//	@Failure		401	{string}	string	"Unauthorized"
//	@Failure		403	{string}	string	"organization required"
//	@Failure		500	{string}	string	"failed to list node pools"
//	@Router			/node-pools [get]
//	@Security		BearerAuth
//	@Security		OrgKeyAuth
//	@Security		OrgSecretAuth
func ListNodePools(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}

		q := db.Where("organization_id = ?", orgID)
		if needle := strings.TrimSpace(r.URL.Query().Get("q")); needle != "" {
			q = q.Where("name LIKE ?", "%"+needle+"%")
		}

		var pools []models.NodePool
		if err := q.
			Preload("Servers").
			Preload("Labels").
			Preload("Taints").
			Preload("Annotations").
			Order("created_at DESC").
			Find(&pools).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		out := make([]dto.NodePoolResponse, 0, len(pools))
		for _, p := range pools {
			npr := dto.NodePoolResponse{
				AuditFields: p.AuditFields,
				Name:        p.Name,
				Role:        dto.NodeRole(p.Role),
				Servers:     make([]dto.ServerResponse, 0, len(p.Servers)),
				Labels:      make([]dto.LabelResponse, 0, len(p.Labels)),
				Taints:      make([]dto.TaintResponse, 0, len(p.Taints)),
				Annotations: make([]dto.AnnotationResponse, 0, len(p.Annotations)),
			}
			//Servers
			for _, s := range p.Servers {
				outSrv := dto.ServerResponse{
					ID:               s.ID,
					Hostname:         s.Hostname,
					PublicIPAddress:  s.PublicIPAddress,
					PrivateIPAddress: s.PrivateIPAddress,
					OrganizationID:   s.OrganizationID,
					SshKeyID:         s.SshKeyID,
					SSHUser:          s.SSHUser,
					Role:             s.Role,
					Status:           s.Status,
					CreatedAt:        s.CreatedAt.UTC().Format(time.RFC3339),
					UpdatedAt:        s.UpdatedAt.UTC().Format(time.RFC3339),
					// add more fields as needed
				}
				npr.Servers = append(npr.Servers, outSrv)
			}
			//Labels
			for _, l := range p.Labels {
				outL := dto.LabelResponse{
					AuditFields: common.AuditFields{
						ID:             l.ID,
						OrganizationID: l.OrganizationID,
						CreatedAt:      l.CreatedAt,
						UpdatedAt:      l.UpdatedAt,
					},
					Key:   l.Key,
					Value: l.Value,
				}
				npr.Labels = append(npr.Labels, outL)
			}
			// Taints
			for _, t := range p.Taints {
				outT := dto.TaintResponse{
					ID:             t.ID,
					OrganizationID: t.OrganizationID,
					CreatedAt:      t.CreatedAt.UTC().Format(time.RFC3339),
					UpdatedAt:      t.UpdatedAt.UTC().Format(time.RFC3339),
					Key:            t.Key,
					Value:          t.Value,
					Effect:         t.Effect,
				}
				npr.Taints = append(npr.Taints, outT)
			}
			// Annotations
			for _, a := range p.Annotations {
				outA := dto.AnnotationResponse{
					AuditFields: common.AuditFields{
						ID:             a.ID,
						OrganizationID: a.OrganizationID,
						CreatedAt:      a.CreatedAt,
						UpdatedAt:      a.UpdatedAt,
					},
					Key:   a.Key,
					Value: a.Value,
				}
				npr.Annotations = append(npr.Annotations, outA)
			}

			out = append(out, npr)
		}
		utils.WriteJSON(w, http.StatusOK, out)
	}
}

// GetNodePool godoc
//
//	@ID				GetNodePool
//	@Summary		Get node pool by ID (org scoped)
//	@Description	Returns one node pool. Add `include=servers` to include servers.
//	@Tags			NodePools
//	@Accept			json
//	@Produce		json
//	@Param			X-Org-ID	header	string	false	"Organization UUID"
//	@Param			id			path	string	true	"Node Pool ID (UUID)"
//	@Success		200	{object}	dto.NodePoolResponse
//	@Failure		400	{string}	string	"invalid id"
//	@Failure		401	{string}	string	"Unauthorized"
//	@Failure		403	{string}	string	"organization required"
//	@Failure		404	{string}	string	"not found"
//	@Failure		500	{string}	string	"fetch failed"
//	@Router			/node-pools/{id} [get]
//	@Security		BearerAuth
//	@Security		OrgKeyAuth
//	@Security		OrgSecretAuth
func GetNodePool(db *gorm.DB) http.HandlerFunc {
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

		var out dto.NodePoolResponse
		if err := db.Model(&models.NodePool{}).Preload("Servers").Where("id = ? AND organization_id = ?", id, orgID).First(&out, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "node_pool_not_found", "node pool not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}
		utils.WriteJSON(w, http.StatusOK, out)
	}
}

// CreateNodePool godoc
//
//	@ID				CreateNodePool
//	@Summary		Create node pool (org scoped)
//	@Description	Creates a node pool. Optionally attach initial servers.
//	@Tags			NodePools
//	@Accept			json
//	@Produce		json
//	@Param			X-Org-ID	header	string					false	"Organization UUID"
//	@Param			body		body	dto.CreateNodePoolRequest	true	"NodePool payload"
//	@Success		201	{object}	dto.NodePoolResponse
//	@Failure		400	{string}	string	"invalid json / missing fields / invalid server_ids"
//	@Failure		401	{string}	string	"Unauthorized"
//	@Failure		403	{string}	string	"organization required"
//	@Failure		500	{string}	string	"create failed"
//	@Router			/node-pools [post]
//	@Security		BearerAuth
//	@Security		OrgKeyAuth
//	@Security		OrgSecretAuth
func CreateNodePool(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}

		var req dto.CreateNodePoolRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_request", "bad request")
			return
		}

		req.Name = strings.TrimSpace(req.Name)
		req.Role = dto.NodeRole(strings.TrimSpace(string(req.Role)))

		if req.Name == "" || req.Role == "" {
			utils.WriteError(w, http.StatusBadRequest, "bad_request", "missing name/role")
			return
		}

		n := models.NodePool{
			AuditFields: common.AuditFields{
				OrganizationID: orgID,
			},
			Name: req.Name,
			Role: string(req.Role),
		}

		if err := db.Create(&n).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		out := dto.NodePoolResponse{
			AuditFields: n.AuditFields,
			Name:        n.Name,
			Role:        dto.NodeRole(n.Role),
		}
		utils.WriteJSON(w, http.StatusCreated, out)
	}
}

// UpdateNodePool godoc
//
//	@ID				UpdateNodePool
//	@Summary		Update node pool (org scoped)
//	@Description	Partially update node pool fields.
//	@Tags			NodePools
//	@Accept			json
//	@Produce		json
//	@Param			X-Org-ID	header	string					false	"Organization UUID"
//	@Param			id			path	string					true	"Node Pool ID (UUID)"
//	@Param			body		body	dto.UpdateNodePoolRequest	true	"Fields to update"
//	@Success		200	{object}	dto.NodePoolResponse
//	@Failure		400	{string}	string	"invalid id / invalid json"
//	@Failure		401	{string}	string	"Unauthorized"
//	@Failure		403	{string}	string	"organization required"
//	@Failure		404	{string}	string	"not found"
//	@Failure		500	{string}	string	"update failed"
//	@Router			/node-pools/{id} [patch]
//	@Security		BearerAuth
//	@Security		OrgKeyAuth
//	@Security		OrgSecretAuth
func UpdateNodePool(db *gorm.DB) http.HandlerFunc {
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

		var n models.NodePool
		if err := db.Where("id = ? AND organization_id = ?", id, orgID).First(&n).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "node_pool_not_found", "node pool not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		var req dto.UpdateNodePoolRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_request", "bad request")
			return
		}

		if req.Name != nil {
			n.Name = strings.TrimSpace(*req.Name)
		}
		if req.Role != nil {
			v := dto.NodeRole(strings.TrimSpace(string(*req.Role)))
			n.Role = string(v)
		}

		if err := db.Save(&n).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}
		out := dto.NodePoolResponse{
			AuditFields: n.AuditFields,
			Name:        n.Name,
			Role:        dto.NodeRole(n.Role),
		}
		utils.WriteJSON(w, http.StatusOK, out)
	}
}

// DeleteNodePool godoc
//
//	@ID				DeleteNodePool
//	@Summary		Delete node pool (org scoped)
//	@Description	Permanently deletes the node pool.
//	@Tags			NodePools
//	@Accept			json
//	@Produce		json
//	@Param			X-Org-ID	header	string	false	"Organization UUID"
//	@Param			id			path	string	true	"Node Pool ID (UUID)"
//	@Success		204	{string}	string	"No Content"
//	@Failure		400	{string}	string	"invalid id"
//	@Failure		401	{string}	string	"Unauthorized"
//	@Failure		403	{string}	string	"organization required"
//	@Failure		500	{string}	string	"delete failed"
//	@Router			/node-pools/{id} [delete]
//	@Security		BearerAuth
//	@Security		OrgKeyAuth
//	@Security		OrgSecretAuth
func DeleteNodePool(db *gorm.DB) http.HandlerFunc {
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

		if err := db.Where("id = ? AND organization_id = ?", id, orgID).Delete(&models.NodePool{}).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// -- Node Pools Servers

// ListNodePoolServers godoc
//
//	@ID			ListNodePoolServers
//	@Summary	List servers attached to a node pool (org scoped)
//	@Tags		NodePools
//	@Accept		json
//	@Produce	json
//	@Param		X-Org-ID	header	string	false	"Organization UUID"
//	@Param		id			path	string	true	"Node Pool ID (UUID)"
//	@Success	200	{array}		dto.ServerResponse
//	@Failure	400	{string}	string	"invalid id"
//	@Failure	401	{string}	string	"Unauthorized"
//	@Failure	403	{string}	string	"organization required"
//	@Failure	404	{string}	string	"not found"
//	@Failure	500	{string}	string	"fetch failed"
//	@Router		/node-pools/{id}/servers [get]
//	@Security	BearerAuth
//	@Security	OrgKeyAuth
//	@Security	OrgSecretAuth
func ListNodePoolServers(db *gorm.DB) http.HandlerFunc {
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

		var np models.NodePool
		if err := db.Where("id = ? AND organization_id = ?", id, orgID).Preload("Servers").First(&np).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "node_pool_not_found", "node pool not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		out := make([]dto.ServerResponse, 0, len(np.Servers))
		for _, server := range np.Servers {
			out = append(out, dto.ServerResponse{
				ID:               server.ID,
				OrganizationID:   server.OrganizationID,
				Hostname:         server.Hostname,
				PrivateIPAddress: server.PrivateIPAddress,
				PublicIPAddress:  server.PublicIPAddress,
				Role:             server.Role,
				SshKeyID:         server.SshKeyID,
				SSHUser:          server.SSHUser,
				Status:           server.Status,
				CreatedAt:        server.CreatedAt.UTC().Format(time.RFC3339),
				UpdatedAt:        server.UpdatedAt.UTC().Format(time.RFC3339),
			})
		}
		utils.WriteJSON(w, http.StatusOK, out)
	}
}

// AttachNodePoolServers godoc
//
//	@ID			AttachNodePoolServers
//	@Summary	Attach servers to a node pool (org scoped)
//	@Tags		NodePools
//	@Accept		json
//	@Produce	json
//	@Param		X-Org-ID	header	string					false	"Organization UUID"
//	@Param		id			path	string					true	"Node Pool ID (UUID)"
//	@Param		body		body	dto.AttachServersRequest	true	"Server IDs to attach"
//	@Success	204	{string}	string	"No Content"
//	@Failure	400	{string}	string	"invalid id / invalid server_ids"
//	@Failure	401	{string}	string	"Unauthorized"
//	@Failure	403	{string}	string	"organization required"
//	@Failure	404	{string}	string	"not found"
//	@Failure	500	{string}	string	"attach failed"
//	@Router		/node-pools/{id}/servers [post]
//	@Security	BearerAuth
//	@Security	OrgKeyAuth
//	@Security	OrgSecretAuth
func AttachNodePoolServers(db *gorm.DB) http.HandlerFunc {
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

		var np models.NodePool
		if err := db.Where("id = ? AND organization_id = ?", id, orgID).First(&np).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "node_pool_not_found", "node pool not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		var req dto.AttachServersRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_request", "bad request")
			return
		}

		ids, err := parseUUIDs(req.ServerIDs)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_request", "invalid server_ids")
			return
		}

		if len(ids) == 0 {
			// nothing to attach
			utils.WriteError(w, http.StatusBadRequest, "bad_request", "nothing to attach")
			return
		}

		// validate IDs belong to org
		if err := ensureServersBelongToOrg(orgID, ids, db); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_request", "invalid server_ids for this organization")
			return
		}

		// fetch only the requested servers
		var servers []models.Server
		if err := db.Where("organization_id = ? AND id IN ?", orgID, ids).Find(&servers).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "attach db error")
			return
		}

		if len(servers) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		if err := db.Model(&np).Association("Servers").Append(&servers); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "attach failed")
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// DetachNodePoolServer godoc
//
//	@ID			DetachNodePoolServer
//	@Summary	Detach one server from a node pool (org scoped)
//	@Tags		NodePools
//	@Accept		json
//	@Produce	json
//	@Param		X-Org-ID	header	string	false	"Organization UUID"
//	@Param		id			path	string	true	"Node Pool ID (UUID)"
//	@Param		serverId	path	string	true	"Server ID (UUID)"
//	@Success	204	{string}	string	"No Content"
//	@Failure	400	{string}	string	"invalid id"
//	@Failure	401	{string}	string	"Unauthorized"
//	@Failure	403	{string}	string	"organization required"
//	@Failure	404	{string}	string	"not found"
//	@Failure	500	{string}	string	"detach failed"
//	@Router		/node-pools/{id}/servers/{serverId} [delete]
//	@Security	BearerAuth
//	@Security	OrgKeyAuth
//	@Security	OrgSecretAuth
func DetachNodePoolServer(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "id_required", "node pool id required")
			return
		}
		serverId, err := uuid.Parse(chi.URLParam(r, "serverId"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "id_required", "server id required")
			return
		}
		var np models.NodePool
		if err := db.Where("id = ? AND organization_id = ?", id, orgID).First(&np).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "node_pool_not_found", "node pool not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		var s models.Server
		if err := db.Where("id = ? AND organization_id = ?", serverId, orgID).First(&s).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "server_not_found", "server not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		if err := db.Model(&np).Association("Servers").Delete(&s); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "detach error")
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// -- Node Pools Taints

// ListNodePoolTaints godoc
//
//	@ID			ListNodePoolTaints
//	@Summary	List taints attached to a node pool (org scoped)
//	@Tags		NodePools
//	@Accept		json
//	@Produce	json
//	@Param		X-Org-ID	header	string	false	"Organization UUID"
//	@Param		id			path	string	true	"Node Pool ID (UUID)"
//	@Success	200	{array}		dto.TaintResponse
//	@Failure	400	{string}	string	"invalid id"
//	@Failure	401	{string}	string	"Unauthorized"
//	@Failure	403	{string}	string	"organization required"
//	@Failure	404	{string}	string	"not found"
//	@Failure	500	{string}	string	"fetch failed"
//	@Router		/node-pools/{id}/taints [get]
//	@Security	BearerAuth
//	@Security	OrgKeyAuth
//	@Security	OrgSecretAuth
func ListNodePoolTaints(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "id_required", "node pool id required")
			return
		}

		var np models.NodePool
		if err := db.Where("id = ? AND organization_id = ?", id, orgID).Preload("Taints").First(&np).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "node_pool_not_found", "node pool not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		out := make([]dto.TaintResponse, 0, len(np.Taints))
		for _, t := range np.Taints {
			out = append(out, dto.TaintResponse{
				ID:     t.ID,
				Key:    t.Key,
				Value:  t.Value,
				Effect: t.Effect,
			})
		}
		utils.WriteJSON(w, http.StatusOK, out)
	}
}

// AttachNodePoolTaints godoc
//
//	@ID			AttachNodePoolTaints
//	@Summary	Attach taints to a node pool (org scoped)
//	@Tags		NodePools
//	@Accept		json
//	@Produce	json
//	@Param		X-Org-ID	header	string					false	"Organization UUID"
//	@Param		id			path	string					true	"Node Pool ID (UUID)"
//	@Param		body		body	dto.AttachTaintsRequest	true	"Taint IDs to attach"
//	@Success	204	{string}	string	"No Content"
//	@Failure	400	{string}	string	"invalid id / invalid taint_ids"
//	@Failure	401	{string}	string	"Unauthorized"
//	@Failure	403	{string}	string	"organization required"
//	@Failure	404	{string}	string	"not found"
//	@Failure	500	{string}	string	"attach failed"
//	@Router		/node-pools/{id}/taints [post]
//	@Security	BearerAuth
//	@Security	OrgKeyAuth
//	@Security	OrgSecretAuth
func AttachNodePoolTaints(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "id_required", "node pool id required")
			return
		}

		var np models.NodePool
		if err := db.Where("id = ? AND organization_id = ?", id, orgID).First(&np).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "node_pool_not_found", "node pool not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		var req dto.AttachTaintsRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_request", "bad request")
			return
		}

		ids, err := parseUUIDs(req.TaintIDs)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_request", "invalid taint_ids")
			return
		}

		if len(ids) == 0 {
			utils.WriteError(w, http.StatusBadRequest, "bad_request", "nothing to attach")
			return
		}

		// validate IDs belong to org
		if err := ensureTaintsBelongToOrg(orgID, ids, db); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_request", "invalid taint_ids for this organization")
			return
		}

		var taints []models.Taint
		if err := db.Where("organization_id = ? AND id IN ?", orgID, ids).Find(&taints).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "attach db error")
			return
		}

		if len(taints) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		if err := db.Model(&np).Association("Taints").Append(&taints); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "attach db error")
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// DetachNodePoolTaint godoc
//
//	@ID			DetachNodePoolTaint
//	@Summary	Detach one taint from a node pool (org scoped)
//	@Tags		NodePools
//	@Accept		json
//	@Produce	json
//	@Param		X-Org-ID	header	string	false	"Organization UUID"
//	@Param		id			path	string	true	"Node Pool ID (UUID)"
//	@Param		taintId	path	string	true	"Taint ID (UUID)"
//	@Success	204	{string}	string	"No Content"
//	@Failure	400	{string}	string	"invalid id"
//	@Failure	401	{string}	string	"Unauthorized"
//	@Failure	403	{string}	string	"organization required"
//	@Failure	404	{string}	string	"not found"
//	@Failure	500	{string}	string	"detach failed"
//	@Router		/node-pools/{id}/taints/{taintId} [delete]
//	@Security	BearerAuth
//	@Security	OrgKeyAuth
//	@Security	OrgSecretAuth
func DetachNodePoolTaint(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "id_required", "node pool id required")
			return
		}
		taintId, err := uuid.Parse(chi.URLParam(r, "taintId"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "taintId_required", "taint id required")
			return
		}

		var np models.NodePool
		if err := db.Where("id = ? AND organization_id = ?", id, orgID).First(&np).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "node_pool_not_found", "node pool not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		var t models.Taint
		if err := db.Where("id = ? AND organization_id = ?", taintId, orgID).First(&t).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "taint_not_found", "taint not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		if err := db.Model(&np).Association("Taints").Delete(&t); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// -- Node Pools Labels

// ListNodePoolLabels godoc
//
//	@ID			ListNodePoolLabels
//	@Summary	List labels attached to a node pool (org scoped)
//	@Tags		NodePools
//	@Accept		json
//	@Produce	json
//	@Param		X-Org-ID	header	string	false	"Organization UUID"
//	@Param		id			path	string	true	"Label Pool ID (UUID)"
//	@Success	200	{array}		dto.LabelResponse
//	@Failure	400	{string}	string	"invalid id"
//	@Failure	401	{string}	string	"Unauthorized"
//	@Failure	403	{string}	string	"organization required"
//	@Failure	404	{string}	string	"not found"
//	@Failure	500	{string}	string	"fetch failed"
//	@Router		/node-pools/{id}/labels [get]
//	@Security	BearerAuth
//	@Security	OrgKeyAuth
//	@Security	OrgSecretAuth
func ListNodePoolLabels(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "id_required", "node pool id required")
			return
		}

		var np models.NodePool
		if err := db.Where("id = ? AND organization_id = ?", id, orgID).Preload("Labels").First(&np).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "node_pool_not_found", "node pool not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		out := make([]dto.LabelResponse, 0, len(np.Taints))
		for _, taint := range np.Taints {
			out = append(out, dto.LabelResponse{
				AuditFields: common.AuditFields{
					ID:             taint.ID,
					OrganizationID: taint.OrganizationID,
					CreatedAt:      taint.CreatedAt,
					UpdatedAt:      taint.UpdatedAt,
				},
				Key:   taint.Key,
				Value: taint.Value,
			})
		}
		utils.WriteJSON(w, http.StatusOK, out)
	}
}

// AttachNodePoolLabels godoc
//
//	@ID			AttachNodePoolLabels
//	@Summary	Attach labels to a node pool (org scoped)
//	@Tags		NodePools
//	@Accept		json
//	@Produce	json
//	@Param		X-Org-ID	header	string					false	"Organization UUID"
//	@Param		id			path	string					true	"Node Pool ID (UUID)"
//	@Param		body		body	dto.AttachLabelsRequest	true	"Label IDs to attach"
//	@Success	204	{string}	string	"No Content"
//	@Failure	400	{string}	string	"invalid id / invalid server_ids"
//	@Failure	401	{string}	string	"Unauthorized"
//	@Failure	403	{string}	string	"organization required"
//	@Failure	404	{string}	string	"not found"
//	@Failure	500	{string}	string	"attach failed"
//	@Router		/node-pools/{id}/labels [post]
//	@Security	BearerAuth
//	@Security	OrgKeyAuth
//	@Security	OrgSecretAuth
func AttachNodePoolLabels(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "id_required", "node pool id required")
			return
		}

		var np models.NodePool
		if err := db.Where("id = ? AND organization_id = ?", id, orgID).First(&np).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "node_pool_not_found", "node pool not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		var req dto.AttachLabelsRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_request", "bad request")
			return
		}

		ids, err := parseUUIDs(req.LabelIDs)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_request", "invalid label_ids")
			return
		}

		if len(ids) == 0 {
			utils.WriteError(w, http.StatusBadRequest, "bad_request", "nothing to attach")
			return
		}

		if err := ensureLabelsBelongToOrg(orgID, ids, db); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_request", "invalid label_ids for this organization")
		}

		var labels []models.Label
		if err := db.Where("organization_id = ? AND id IN ?", orgID, ids).Find(&labels).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "attach db error")
			return
		}

		if len(labels) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		if err := db.Model(&np).Association("Labels").Append(&labels); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "attach failed")
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// DetachNodePoolLabel godoc
//
//	@ID			DetachNodePoolLabel
//	@Summary	Detach one label from a node pool (org scoped)
//	@Tags		NodePools
//	@Accept		json
//	@Produce	json
//	@Param		X-Org-ID	header	string	false	"Organization UUID"
//	@Param		id			path	string	true	"Node Pool ID (UUID)"
//	@Param		labelId	path	string	true	"Label ID (UUID)"
//	@Success	204	{string}	string	"No Content"
//	@Failure	400	{string}	string	"invalid id"
//	@Failure	401	{string}	string	"Unauthorized"
//	@Failure	403	{string}	string	"organization required"
//	@Failure	404	{string}	string	"not found"
//	@Failure	500	{string}	string	"detach failed"
//	@Router		/node-pools/{id}/labels/{labelId} [delete]
//	@Security	BearerAuth
//	@Security	OrgKeyAuth
//	@Security	OrgSecretAuth
func DetachNodePoolLabel(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "id_required", "node pool id required")
			return
		}
		labelId, err := uuid.Parse(chi.URLParam(r, "labelId"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "id_required", "labelId required")
			return
		}

		var np models.NodePool
		if err := db.Where("id = ? AND organization_id = ?", id, orgID).First(&np).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "node_pool_not_found", "node pool not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		var l models.Label
		if err := db.Where("id = ? AND organization_id = ?", labelId, orgID).First(&l).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "label_not_found", "label not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		if err := db.Model(&np).Association("Labels").Delete(&l); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "detach error")
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// -- Node Pools Annotations

// ListNodePoolAnnotations godoc
//
//	@ID			ListNodePoolAnnotations
//	@Summary	List annotations attached to a node pool (org scoped)
//	@Tags		NodePools
//	@Accept		json
//	@Produce	json
//	@Param		X-Org-ID	header	string	false	"Organization UUID"
//	@Param		id			path	string	true	"Node Pool ID (UUID)"
//	@Success	200	{array}		dto.AnnotationResponse
//	@Failure	400	{string}	string	"invalid id"
//	@Failure	401	{string}	string	"Unauthorized"
//	@Failure	403	{string}	string	"organization required"
//	@Failure	404	{string}	string	"not found"
//	@Failure	500	{string}	string	"fetch failed"
//	@Router		/node-pools/{id}/annotations [get]
//	@Security	BearerAuth
//	@Security	OrgKeyAuth
//	@Security	OrgSecretAuth
func ListNodePoolAnnotations(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "id_required", "node pool id required")
			return
		}

		var np models.NodePool
		if err := db.Where("id = ? AND organization_id = ?", id, orgID).Preload("Labels").First(&np).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "node_pool_not_found", "node pool not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		out := make([]dto.AnnotationResponse, 0, len(np.Annotations))
		for _, ann := range np.Annotations {
			out = append(out, dto.AnnotationResponse{
				AuditFields: common.AuditFields{
					ID:             ann.ID,
					OrganizationID: ann.OrganizationID,
					CreatedAt:      ann.CreatedAt,
					UpdatedAt:      ann.UpdatedAt,
				},
				Key:   ann.Key,
				Value: ann.Value,
			})
		}
		utils.WriteJSON(w, http.StatusOK, out)
	}
}

// AttachNodePoolAnnotations godoc
//
//	@ID			AttachNodePoolAnnotations
//	@Summary	Attach annotation to a node pool (org scoped)
//	@Tags		NodePools
//	@Accept		json
//	@Produce	json
//	@Param		X-Org-ID	header	string					false	"Organization UUID"
//	@Param		id			path	string					true	"Node Group ID (UUID)"
//	@Param		body		body	dto.AttachAnnotationsRequest	true	"Annotation IDs to attach"
//	@Success	204	{string}	string	"No Content"
//	@Failure	400	{string}	string	"invalid id / invalid server_ids"
//	@Failure	401	{string}	string	"Unauthorized"
//	@Failure	403	{string}	string	"organization required"
//	@Failure	404	{string}	string	"not found"
//	@Failure	500	{string}	string	"attach failed"
//	@Router		/node-pools/{id}/annotations [post]
//	@Security	BearerAuth
//	@Security	OrgKeyAuth
//	@Security	OrgSecretAuth
func AttachNodePoolAnnotations(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "id_required", "node pool id required")
			return
		}

		var np models.NodePool
		if err := db.Where("id = ? AND organization_id = ?", id, orgID).First(&np).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "node_pool_not_found", "node pool not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		var req dto.AttachAnnotationsRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_request", "bad request")
			return
		}

		ids, err := parseUUIDs(req.AnnotationIDs)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_request", "invalid annotation ids")
			return
		}

		if len(ids) == 0 {
			utils.WriteError(w, http.StatusBadRequest, "bad_request", "nothing to attach")
			return
		}

		if err := ensureAnnotaionsBelongToOrg(orgID, ids, db); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_request", "invalid annotation ids for this organization")
			return
		}

		var ann []models.Annotation
		if err := db.Where("organization_id = ? AND id IN ?", orgID, ids).Find(&ann).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		if len(ann) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		if err := db.Model(&np).Association("Annotations").Append(&ann); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "attach failed")
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// DetachNodePoolAnnotation godoc
//
//	@ID			DetachNodePoolAnnotation
//	@Summary	Detach one annotation from a node pool (org scoped)
//	@Tags		NodePools
//	@Accept		json
//	@Produce	json
//	@Param		X-Org-ID	header	string	false	"Organization UUID"
//	@Param		id			path	string	true	"Node Pool ID (UUID)"
//	@Param		annotationId	path	string	true	"Annotation ID (UUID)"
//	@Success	204	{string}	string	"No Content"
//	@Failure	400	{string}	string	"invalid id"
//	@Failure	401	{string}	string	"Unauthorized"
//	@Failure	403	{string}	string	"organization required"
//	@Failure	404	{string}	string	"not found"
//	@Failure	500	{string}	string	"detach failed"
//	@Router		/node-pools/{id}/annotations/{annotationId} [delete]
//	@Security	BearerAuth
//	@Security	OrgKeyAuth
//	@Security	OrgSecretAuth
func DetachNodePoolAnnotation(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "id_required", "node pool id required")
			return
		}
		annotationId, err := uuid.Parse(chi.URLParam(r, "annotationId"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "id_required", "node pool annotation id required")
			return
		}

		var np models.NodePool
		if err := db.Where("id = ? AND organization_id = ?", id, orgID).First(&np).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "node_pool_not_found", "node pool not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		var ann []models.Annotation
		if err := db.Where("id = ? AND organization_id = ?", annotationId, orgID).First(&ann).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "annotation_not_found", "annotation not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		if err := db.Model(&np).Association("Annotations").Delete(&ann); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// -- Helpers
func parseUUIDs(ids []string) ([]uuid.UUID, error) {
	out := make([]uuid.UUID, 0, len(ids))
	for _, id := range ids {
		u, err := uuid.Parse(id)
		if err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, nil
}

func ensureServersBelongToOrg(orgID uuid.UUID, ids []uuid.UUID, db *gorm.DB) error {
	var count int64
	if err := db.Model(&models.Server{}).Where("organization_id = ? AND id IN ?", orgID, ids).Count(&count).Error; err != nil {
		return err
	}
	if count != int64(len(ids)) {
		return errors.New("some servers do not belong to this org")
	}
	return nil
}

func ensureTaintsBelongToOrg(orgID uuid.UUID, ids []uuid.UUID, db *gorm.DB) error {
	var count int64
	if err := db.Model(&models.Taint{}).Where("organization_id = ? AND id IN ?", orgID, ids).Count(&count).Error; err != nil {
		return err
	}
	if count != int64(len(ids)) {
		return errors.New("some taints do not belong to this org")
	}
	return nil
}

func ensureLabelsBelongToOrg(orgID uuid.UUID, ids []uuid.UUID, db *gorm.DB) error {
	var count int64
	if err := db.Model(&models.Label{}).Where("organization_id = ? AND id IN ?", orgID, ids).Count(&count).Error; err != nil {
		return err
	}
	if count != int64(len(ids)) {
		return errors.New("some labels do not belong to this org")
	}
	return nil
}

func ensureAnnotaionsBelongToOrg(orgID uuid.UUID, ids []uuid.UUID, db *gorm.DB) error {
	var count int64
	if err := db.Model(&models.Annotation{}).Where("organization_id = ? AND id IN ?", orgID, ids).Count(&count).Error; err != nil {
		return err
	}
	if count != int64(len(ids)) {
		return errors.New("some annotations do not belong to this org")
	}
	return nil
}
