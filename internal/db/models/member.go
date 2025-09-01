package models

import "github.com/google/uuid"

type MemberRole string

const (
	MemberRoleAdmin  MemberRole = "admin"
	MemberRoleMember MemberRole = "member"
	MemberRoleUser   MemberRole = "user"
)

type Member struct {
	ID             uuid.UUID    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID         uuid.UUID    `gorm:"type:uuid;not null" json:"user_id"`
	User           User         `gorm:"foreignKey:UserID" json:"user"`
	OrganizationID uuid.UUID    `gorm:"type:uuid;not null" json:"organization_id"`
	Organization   Organization `gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE" json:"organization"`
	Role           MemberRole   `gorm:"not null;default:member" json:"role"` // e.g. admin, member
	Timestamped
}
