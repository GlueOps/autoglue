package db

import (
	"fmt"
	"log"

	"github.com/glueops/autoglue/internal/db/models"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := viper.GetString("database.dsn")
	log.Println("DB DSN:", dsn)
	if dsn == "" {
		log.Fatal("AUTOGLUE_DATABASE_DSN is not set")
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}

	err = DB.AutoMigrate(
		&models.Cluster{},
		&models.Credential{},
		&models.Invitation{},
		&models.MasterKey{},
		&models.Member{},
		&models.NodeGroup{},
		&models.NodeLabel{},
		&models.NodeTaint{},
		&models.Organization{},
		&models.OrganizationKey{},
		&models.RefreshToken{},
		&models.Server{},
		&models.SshKey{},
		&models.User{},
	)

	if err != nil {
		log.Fatalf("auto migration failed: %v", err)
	}

	fmt.Println("database connected and migrated")
}
