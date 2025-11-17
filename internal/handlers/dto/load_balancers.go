package dto

import (
	"time"

	"github.com/google/uuid"
)

type LoadBalancerResponse struct {
	ID               uuid.UUID `json:"id"`
	OrganizationID   uuid.UUID `json:"organization_id"`
	Name             string    `json:"name"`
	Kind             string    `json:"kind"`
	PublicIPAddress  string    `json:"public_ip_address"`
	PrivateIPAddress string    `json:"private_ip_address"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type CreateLoadBalancerRequest struct {
	Name             string `json:"name" example:"glueops"`
	Kind             string `json:"kind" example:"public" enums:"glueops,public"`
	PublicIPAddress  string `json:"public_ip_address" example:"8.8.8.8"`
	PrivateIPAddress string `json:"private_ip_address" example:"192.168.0.2"`
}

type UpdateLoadBalancerRequest struct {
	Name             *string `json:"name" example:"glue"`
	Kind             *string `json:"kind" example:"public" enums:"glueops,public"`
	PublicIPAddress  *string `json:"public_ip_address" example:"8.8.8.8"`
	PrivateIPAddress *string `json:"private_ip_address" example:"192.168.0.2"`
}
