package models

import (
	"github.com/glueops/autoglue/internal/common"
	"github.com/google/uuid"
)

type ClusterMetadata struct {
	common.AuditFields `gorm:"embedded"`
	ClusterID          uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_cluster_metadata_cluster_key" json:"cluster_id"`
	Cluster            Cluster   `gorm:"foreignKey:ClusterID;constraint:OnDelete:CASCADE" json:"-"`
	Key                string    `gorm:"not null;uniqueIndex:idx_cluster_metadata_cluster_key" json:"key"`
	Value              string    `gorm:"not null" json:"value"`
}
