package cli

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
	Short: "Initialize config",
	Long:  "Initialize configuration file",
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
				"dsn": "postgres://user:pass@localhost:5432/db?sslmode=disable",
			},
			"authentication": map[string]string{
				"secret": defaultSecret,
			},
			"smtp": map[string]interface{}{
				"enabled":  false,
				"host":     "smtp.example.com",
				"port":     587,
				"username": "",
				"password": "",
				"from":     "no-reply@example.com",
			},
			"frontend": map[string]string{
				"base_url": "http://localhost:5173",
			},
			"ui": map[string]string{
				"dev": "false",
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
		fmt.Println("config.yaml written")
	},
}

func init() {
	rootCmd.AddCommand(initConfigCmd)

	_ = viper.BindPFlag("database.dsn", rootCmd.PersistentFlags().Lookup("dsn"))
	_ = viper.BindPFlag("authentication.secret", rootCmd.PersistentFlags().Lookup("authentication-secret"))
}
