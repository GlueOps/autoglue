package config

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var File = "config.yaml"
var fileKeys = map[string]bool{}

func Load() {
	_ = godotenv.Load()

	viper.SetDefault("bind_address", "127.0.0.1")
	viper.SetDefault("bind_port", "8080")
	viper.SetDefault("database.dsn", "postgres://user:pass@localhost:5432/db?sslmode=disable")

	viper.SetDefault("ui.dev", false)

	viper.SetDefault("authentication.secret", GenerateSecureSecret())

	viper.SetDefault("smtp.enabled", false)
	viper.SetDefault("smtp.host", "smtp.example.com")
	viper.SetDefault("smtp.port", 587)
	viper.SetDefault("smtp.username", "")
	viper.SetDefault("smtp.password", "")
	viper.SetDefault("smtp.from", "no-reply@example.com")

	viper.SetDefault("frontend.base_url", "http://localhost:5173")

	viper.SetDefault("archer.instances", 2)
	viper.SetDefault("archer.timeoutSec", 60)

	viper.SetEnvPrefix("AUTOGLUE")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.SetConfigFile(File)
	viper.SetConfigType("yaml")

	if _, err := os.Stat(File); err == nil {
		err := viper.ReadInConfig()
		if err != nil {
			log.Fatalf("Failed to read config file: %v", err)
		}
		for _, k := range viper.AllKeys() {
			fileKeys[k] = true
		}
		fmt.Println("Loaded config from", File)
	}
}

func GenerateSecureSecret() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		panic("unable to generate secure secret")
	}
	return base64.URLEncoding.EncodeToString(b)
}

func GetAuthSecret() string {
	return viper.GetString("authentication.secret")
}

func DebugPrintConfig() {
	all := viper.AllSettings()

	b, err := yaml.Marshal(all)
	if err != nil {
		fmt.Println("error marshalling config:", err)
		return
	}
	fmt.Println("Loaded configuration:")
	fmt.Println(string(b))
}

func IsUIDev() bool {
	return viper.GetBool("ui.dev")
}

func SMTPEnabled() bool {
	return viper.GetBool("smtp.enabled")
}

func SMTPHost() string {
	return viper.GetString("smtp.host")
}

func SMTPPort() int {
	return viper.GetInt("smtp.port")
}

func SMTPUsername() string {
	return viper.GetString("smtp.username")
}

func SMTPPassword() string {
	return viper.GetString("smtp.password")
}

func SMTPFrom() string {
	return viper.GetString("smtp.from")
}

func FrontendBaseURL() string {
	return viper.GetString("frontend.base_url")
}
