package cmd

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/glueops/autoglue/internal/app"
	"github.com/glueops/autoglue/internal/models"
	"github.com/spf13/cobra"
)

var rotateMasterCmd = &cobra.Command{
	Use:   "rotate-master",
	Short: "Generate and activate a new master encryption key",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		rt := app.NewRuntime()
		db := rt.DB

		key := make([]byte, 32)
		if _, err := io.ReadFull(rand.Reader, key); err != nil {
			return fmt.Errorf("generating random key: %w", err)
		}

		encoded := base64.StdEncoding.EncodeToString(key)

		if err := db.Model(&models.MasterKey{}).
			Where("is_active = ?", true).
			Update("is_active", false).Error; err != nil {
			return fmt.Errorf("deactivating previous key: %w", err)
		}

		if err := db.Create(&models.MasterKey{
			Key:      encoded,
			IsActive: true,
		}).Error; err != nil {
			return fmt.Errorf("creating new master key: %w", err)
		}

		fmt.Println("Master key rotated successfully")
		return nil
	},
}

var createMasterCmd = &cobra.Command{
	Use:   "create-master",
	Short: "Generate and activate a new master encryption key",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		rt := app.NewRuntime()
		db := rt.DB
		key := make([]byte, 32)
		if _, err := io.ReadFull(rand.Reader, key); err != nil {
			return fmt.Errorf("generating random key: %w", err)
		}

		encoded := base64.StdEncoding.EncodeToString(key)

		if err := db.Create(&models.MasterKey{
			Key:      encoded,
			IsActive: true,
		}).Error; err != nil {
			return fmt.Errorf("creating master key: %w", err)
		}

		fmt.Println("Master key created successfully")
		return nil
	},
}

var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "Manage autoglue encryption keys",
	Long:  "Manage autoglue master encryption keys used for securing data.",
}

func init() {
	encryptCmd.AddCommand(rotateMasterCmd)
	encryptCmd.AddCommand(createMasterCmd)
	rootCmd.AddCommand(encryptCmd)
}
