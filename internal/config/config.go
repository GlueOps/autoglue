package config

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
)

var File = "config.yaml"

func Load() {
	_ = godotenv.Load()

	viper.SetDefault("bind_address", "127.0.0.1")
	viper.SetDefault("bind_port", "8080")
	viper.SetDefault("database.dsn", "postgres://user:pass@localhost:5432/autoglue?sslmode=disable")

	viper.SetEnvPrefix("AUTOGLUE")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.AutomaticEnv()

	viper.SetDefault("authentication.secret", GenerateSecureSecret())

	viper.SetConfigFile(File)
	viper.SetConfigType("yaml")

	if _, err := os.Stat(File); err == nil {
		err := viper.ReadInConfig()
		if err != nil {
			log.Fatalf("Failed to read config file: %v", err)
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
	fmt.Println("Loaded configuration:")
	for k, v := range all {
		fmt.Printf("%s: %#v\n", k, v)
	}
}
