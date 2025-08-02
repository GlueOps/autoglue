package db

import (
	"fmt"
	"github.com/glueops/autoglue/internal/db/models"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
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
		&models.User{},
		&models.Organization{},
		&models.Member{},
		&models.RefreshToken{},
		&models.OrganizationKey{},
		&models.SshKey{},
		&models.Credential{},
		&models.Invitation{},
	)
	if err != nil {
		log.Fatalf("auto migration failed: %v", err)
	}

	fmt.Println("database connected and migrated")
}
