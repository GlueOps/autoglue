package mapper

import (
	"fmt"
	"time"

	"github.com/glueops/autoglue/internal/common"
	"github.com/glueops/autoglue/internal/handlers/dto"
	"github.com/glueops/autoglue/internal/models"
	"github.com/google/uuid"
)

func ClusterToDTO(c models.Cluster) dto.ClusterResponse {
	var bastion *dto.ServerResponse
	if c.BastionServer != nil {
		b := ServerToDTO(*c.BastionServer)
		bastion = &b
	}

	var captainDomain *dto.DomainResponse
	if c.CaptainDomainID != nil && c.CaptainDomain.ID != uuid.Nil {
		dr := DomainToDTO(c.CaptainDomain)
		captainDomain = &dr
	}

	var controlPlane *dto.RecordSetResponse
	if c.ControlPlaneRecordSet != nil {
		rr := RecordSetToDTO(*c.ControlPlaneRecordSet)
		controlPlane = &rr
	}

	var cfqdn *string
	if captainDomain != nil && controlPlane != nil {
		fq := fmt.Sprintf("%s.%s", controlPlane.Name, captainDomain.DomainName)
		cfqdn = &fq
	}

	var appsLB *dto.LoadBalancerResponse
	if c.AppsLoadBalancer != nil {
		lr := LoadBalancerToDTO(*c.AppsLoadBalancer)
		appsLB = &lr
	}

	var glueOpsLB *dto.LoadBalancerResponse
	if c.GlueOpsLoadBalancer != nil {
		lr := LoadBalancerToDTO(*c.GlueOpsLoadBalancer)
		glueOpsLB = &lr
	}

	nps := make([]dto.NodePoolResponse, 0, len(c.NodePools))
	for _, np := range c.NodePools {
		nps = append(nps, NodePoolToDTO(np))
	}

	return dto.ClusterResponse{
		ID:                    c.ID,
		Name:                  c.Name,
		CaptainDomain:         captainDomain,
		ControlPlaneRecordSet: controlPlane,
		ControlPlaneFQDN:      cfqdn,
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
		DockerImage:           c.DockerImage,
		DockerTag:             c.DockerTag,
		CreatedAt:             c.CreatedAt,
		UpdatedAt:             c.UpdatedAt,
	}
}

func NodePoolToDTO(np models.NodePool) dto.NodePoolResponse {
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
		servers = append(servers, ServerToDTO(s))
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

func ServerToDTO(s models.Server) dto.ServerResponse {
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

func DomainToDTO(d models.Domain) dto.DomainResponse {
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

func RecordSetToDTO(rs models.RecordSet) dto.RecordSetResponse {
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

func LoadBalancerToDTO(lb models.LoadBalancer) dto.LoadBalancerResponse {
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
