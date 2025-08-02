package models

type OrganizationKey struct {
	ID             string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrganizationID string `gorm:"uniqueIndex"`
	EncryptedKey   string `gorm:"not null"`
	IV             string `gorm:"not null"`
	Tag            string `gorm:"not null"`
	Timestamped
}

type Credential struct {
	ID             string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrganizationID string `gorm:"uniqueIndex"`
	Provider       string `gorm:"type:varchar(50);not null"`
	EncryptedData  string `gorm:"not null"`
	IV             string `gorm:"not null"`
	Tag            string `gorm:"not null"`
	Timestamped
}

type SshKey struct {
	ID             string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrganizationID string `gorm:"uniqueIndex"`
	PublicKey      string `gorm:"not null"`
	PrivateKey     string `gorm:"not null"`
	Timestamped
}
