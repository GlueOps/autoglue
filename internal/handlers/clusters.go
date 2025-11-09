package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/glueops/autoglue/internal/api/httpmiddleware"
	"github.com/glueops/autoglue/internal/common"
	"github.com/glueops/autoglue/internal/handlers/dto"
	"github.com/glueops/autoglue/internal/models"
	"github.com/glueops/autoglue/internal/utils"
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
			Preload("NodePools").
			Preload("NodePools.Labels").
			Preload("NodePools.Annotations").
			Preload("NodePools.Labels").
			Preload("NodePools.Taints").
			Preload("NodePools.Servers").
			Preload("BastionServer").
			Find(&rows).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		out := make([]dto.ClusterResponse, 0, len(rows))
		for _, row := range rows {
			out = append(out, clusterToDTO(row))
		}

	}

}

// CreateCluster godoc
//
//	@ID				CreateCluster
//	@Summary		Create cluster (org scoped)
//	@Description	Creates a cluster. If `kubeconfig` is provided, it will be encrypted per-organization and stored securely (never returned).
//	@Tags			Clusters
//	@Accept			json
//	@Produce		json
//	@Param			X-Org-ID	header	string	false	"Organization UUID"
//	@Param			body		body	dto.CreateClusterRequest	true	"payload"
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
		w.WriteHeader(http.StatusNoContent)
	}
}

// -- Helpers

func clusterToDTO(c models.Cluster) dto.ClusterResponse {
	var bastion *dto.ServerResponse
	if c.BastionServer != nil {
		b := serverToDTO(*c.BastionServer)
		bastion = &b
	}

	nps := make([]dto.NodePoolResponse, 0, len(c.NodePools))
	for _, np := range c.NodePools {
		nps = append(nps, nodePoolToDTO(np))
	}

	return dto.ClusterResponse{
		ID:                  c.ID,
		Name:                c.Name,
		Provider:            c.Provider,
		Region:              c.Region,
		Status:              c.Status,
		CaptainDomain:       c.CaptainDomain,
		ClusterLoadBalancer: c.ClusterLoadBalancer,
		RandomToken:         c.RandomToken,
		CertificateKey:      c.CertificateKey,
		ControlLoadBalancer: c.ControlLoadBalancer,
		NodePools:           nps,
		BastionServer:       bastion,
		CreatedAt:           c.CreatedAt,
		UpdatedAt:           c.UpdatedAt,
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
