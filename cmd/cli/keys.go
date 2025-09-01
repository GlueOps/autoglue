package cli

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/glueops/autoglue/internal/db"
	"github.com/glueops/autoglue/internal/db/models"
	"github.com/spf13/cobra"
)

var rotateMasterCmd = &cobra.Command{
	Use:   "rotate-master",
	Short: "Generate and activate a new master encryption key",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		key := make([]byte, 32)
		if _, err := io.ReadFull(rand.Reader, key); err != nil {
			return fmt.Errorf("generating random key: %w", err)
		}

		encoded := base64.StdEncoding.EncodeToString(key)

		if err := db.DB.Model(&models.MasterKey{}).
			Where("is_active = ?", true).
			Update("is_active", false).Error; err != nil {
			return fmt.Errorf("deactivating previous key: %w", err)
		}

		if err := db.DB.Create(&models.MasterKey{
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
		key := make([]byte, 32)
		if _, err := io.ReadFull(rand.Reader, key); err != nil {
			return fmt.Errorf("generating random key: %w", err)
		}

		encoded := base64.StdEncoding.EncodeToString(key)

		if err := db.DB.Create(&models.MasterKey{
			Key:      encoded,
			IsActive: true,
		}).Error; err != nil {
			return fmt.Errorf("creating master key: %w", err)
		}

		fmt.Println("Master key created successfully")
		return nil
	},
}

var keysCmd = &cobra.Command{
	Use:   "keys",
	Short: "Manage autoglue encryption keys",
	Long:  "Manage autoglue master encryption keys used for securing data.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if db.DB != nil {
			return nil
		}
		db.Connect()
		return nil
	},
}

func init() {
	keysCmd.AddCommand(rotateMasterCmd)
	keysCmd.AddCommand(createMasterCmd)
	rootCmd.AddCommand(keysCmd)
}
