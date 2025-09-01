package cli

import (
	"fmt"
	"log"

	"github.com/glueops/autoglue/internal/db"
	"github.com/glueops/autoglue/internal/db/models"
	"github.com/spf13/cobra"
)

var (
	userEmail string
)

var setAdminCmd = &cobra.Command{
	Use:   "set-admin",
	Short: "Set an existing user to admin role",
	Long:  "Set an existing user to admin role, looked up by email address",
	Run: func(cmd *cobra.Command, args []string) {
		if userEmail == "" {
			log.Fatal("email is required (use --email)")
		}

		db.Connect()

		var user models.User
		if err := db.DB.Where("email = ?", userEmail).First(&user).Error; err != nil {
			log.Fatalf("could not find user with email %s: %v", userEmail, err)
		}

		if err := db.DB.Model(&user).Update("role", models.RoleAdmin).Error; err != nil {
			log.Fatalf("failed to update user role: %v", err)
		}

		fmt.Printf("User %s (%s) set to admin role\n", user.Name, user.Email)
	},
}

func init() {
	setAdminCmd.Flags().StringVarP(&userEmail, "email", "e", "", "Email of the user to promote to admin")
	rootCmd.AddCommand(setAdminCmd)
}
