package app

import (
	"log"

	"github.com/glueops/autoglue/internal/config"
	"github.com/glueops/autoglue/internal/db"
	"github.com/glueops/autoglue/internal/models"
	"gorm.io/gorm"
)

type Runtime struct {
	Cfg config.Config
	DB  *gorm.DB
}

func NewRuntime() *Runtime {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	d := db.Open(cfg.DbURL)

	err = db.Run(d,
		&models.Job{},
		&models.MasterKey{},
		&models.SigningKey{},
		&models.User{},
		&models.Organization{},
		&models.Account{},
		&models.Membership{},
		&models.APIKey{},
		&models.UserEmail{},
		&models.RefreshToken{},
		&models.OrganizationKey{},
		&models.SshKey{},
		&models.Server{},
		&models.Taint{},
		&models.Label{},
		&models.Annotation{},
		&models.NodePool{},
		&models.Credential{},
		&models.Domain{},
		&models.RecordSet{},
		&models.LoadBalancer{},
		&models.Cluster{},
	)

	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	return &Runtime{
		Cfg: cfg,
		DB:  d,
	}
}
