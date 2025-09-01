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

	if dsn == "" {
		log.Fatal("DRAGON_DATABASE_DSN is not set")
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}

	err = DB.AutoMigrate(
		&models.Credential{},
		&models.EmailVerification{},
		&models.Invitation{},
		&models.MasterKey{},
		&models.Member{},
		&models.Organization{},
		&models.OrganizationKey{},
		&models.PasswordReset{},
		&models.RefreshToken{},
		&models.SshKey{},
		&models.User{},
	)
	if err != nil {
		log.Fatalf("auto migration failed: %v", err)
	}

	fmt.Println("database connected and migrated")
}
