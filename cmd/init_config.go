package cmd

import (
	"fmt"
	"os"

	"github.com/glueops/autoglue/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var initConfigCmd = &cobra.Command{
	Use:   "init-config",
	Short: "Create a default config.yaml file",
	Run: func(cmd *cobra.Command, args []string) {
		file := "config.yaml"
		if _, err := os.Stat(file); err == nil {
			fmt.Println("config.yaml already exists")
			return
		}

		defaultSecret := config.GenerateSecureSecret()

		defaultConfig := map[string]interface{}{
			"bind_address": "127.0.0.1",
			"bind_port":    "8080",
			"database": map[string]string{
				"dsn": "postgres://user:pass@localhost:5432/autoglue?sslmode=disable",
			},
			"authentication": map[string]string{
				"secret": defaultSecret,
			},
		}

		data, err := yaml.Marshal(defaultConfig)
		if err != nil {
			fmt.Println("Error marshalling YAML:", err)
			return
		}

		err = os.WriteFile(file, data, 0644)
		if err != nil {
			fmt.Println("Error writing config.yaml:", err)
			return
		}

		fmt.Println("Created config.yaml")
	},
}

func init() {
	rootCmd.AddCommand(initConfigCmd)
	rootCmd.PersistentFlags().String("dsn", "", "Database DSN")
	rootCmd.PersistentFlags().String("authentication-secret", "", "Authentication secret")

	_ = viper.BindPFlag("database.dsn", rootCmd.PersistentFlags().Lookup("dsn"))
	_ = viper.BindPFlag("authentication.secret", rootCmd.PersistentFlags().Lookup("authentication-secret"))
}
