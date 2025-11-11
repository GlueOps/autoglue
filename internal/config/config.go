package config

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type Config struct {
	DbURL              string
	DbURLRO            string
	Port               string
	Host               string
	JWTIssuer          string
	JWTAudience        string
	JWTPrivateEncKey   string
	OAuthRedirectBase  string
	GoogleClientID     string
	GoogleClientSecret string
	GithubClientID     string
	GithubClientSecret string

	UIDev       bool
	Env         string
	Debug       bool
	Swagger     bool
	SwaggerHost string

	DBStudioEnabled bool
	DBStudioBind    string
	DBStudioPort    string
	DBStudioUser    string
	DBStudioPass    string
}

var (
	once    sync.Once
	cached  Config
	loadErr error
)

func Load() (Config, error) {
	once.Do(func() {
		_ = godotenv.Load()

		// Use a private viper to avoid global mutation/races
		v := viper.New()

		// Defaults
		v.SetDefault("bind.address", "127.0.0.1")
		v.SetDefault("bind.port", "8080")
		v.SetDefault("database.url", "postgres://user:pass@localhost:5432/db?sslmode=disable")
		v.SetDefault("database.url_ro", "")
		v.SetDefault("db_studio.enabled", false)
		v.SetDefault("db_studio.bind", "127.0.0.1")
		v.SetDefault("db_studio.port", "0") // 0 = random
		v.SetDefault("db_studio.user", "")
		v.SetDefault("db_studio.pass", "")

		v.SetDefault("ui.dev", false)
		v.SetDefault("env", "development")
		v.SetDefault("debug", false)
		v.SetDefault("swagger", false)
		v.SetDefault("swagger.host", "localhost:8080")

		// Env setup and binding
		v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		v.AutomaticEnv()

		keys := []string{
			"bind.address",
			"bind.port",
			"database.url",
			"database.url_ro",
			"jwt.issuer",
			"jwt.audience",
			"jwt.private.enc.key",
			"oauth.redirect.base",
			"google.client.id",
			"google.client.secret",
			"github.client.id",
			"github.client.secret",
			"ui.dev",
			"env",
			"debug",
			"swagger",
			"swagger.host",
			"db_studio.enabled",
			"db_studio.bind",
			"db_studio.port",
			"db_studio.user",
			"db_studio.pass",
		}
		for _, k := range keys {
			_ = v.BindEnv(k)
		}

		// Build config
		cfg := Config{
			DbURL:              v.GetString("database.url"),
			DbURLRO:            v.GetString("database.url_ro"),
			Port:               v.GetString("bind.port"),
			Host:               v.GetString("bind.address"),
			JWTIssuer:          v.GetString("jwt.issuer"),
			JWTAudience:        v.GetString("jwt.audience"),
			JWTPrivateEncKey:   v.GetString("jwt.private.enc.key"),
			OAuthRedirectBase:  v.GetString("oauth.redirect.base"),
			GoogleClientID:     v.GetString("google.client.id"),
			GoogleClientSecret: v.GetString("google.client.secret"),
			GithubClientID:     v.GetString("github.client.id"),
			GithubClientSecret: v.GetString("github.client.secret"),

			UIDev:       v.GetBool("ui.dev"),
			Env:         v.GetString("env"),
			Debug:       v.GetBool("debug"),
			Swagger:     v.GetBool("swagger"),
			SwaggerHost: v.GetString("swagger.host"),

			DBStudioEnabled: v.GetBool("db_studio.enabled"),
			DBStudioBind:    v.GetString("db_studio.bind"),
			DBStudioPort:    v.GetString("db_studio.port"),
			DBStudioUser:    v.GetString("db_studio.user"),
			DBStudioPass:    v.GetString("db_studio.pass"),
		}

		// Validate
		if err := validateConfig(cfg); err != nil {
			loadErr = err
			return
		}

		cached = cfg
	})
	return cached, loadErr
}

func validateConfig(cfg Config) error {
	var errs []string

	// Required general settings
	req := map[string]string{
		"jwt.issuer":          cfg.JWTIssuer,
		"jwt.audience":        cfg.JWTAudience,
		"jwt.private.enc.key": cfg.JWTPrivateEncKey,
		"oauth.redirect.base": cfg.OAuthRedirectBase,
	}
	for k, v := range req {
		if strings.TrimSpace(v) == "" {
			errs = append(errs, fmt.Sprintf("missing required config key %q (env %s)", k, envNameFromKey(k)))
		}
	}

	// OAuth provider requirements:
	googleOK := strings.TrimSpace(cfg.GoogleClientID) != "" && strings.TrimSpace(cfg.GoogleClientSecret) != ""
	githubOK := strings.TrimSpace(cfg.GithubClientID) != "" && strings.TrimSpace(cfg.GithubClientSecret) != ""

	// If partially configured, report what's missing for each
	if !googleOK && (cfg.GoogleClientID != "" || cfg.GoogleClientSecret != "") {
		if cfg.GoogleClientID == "" {
			errs = append(errs, fmt.Sprintf("google.client.id is missing (env %s) while google.client.secret is set", envNameFromKey("google.client.id")))
		}
		if cfg.GoogleClientSecret == "" {
			errs = append(errs, fmt.Sprintf("google.client.secret is missing (env %s) while google.client.id is set", envNameFromKey("google.client.secret")))
		}
	}
	if !githubOK && (cfg.GithubClientID != "" || cfg.GithubClientSecret != "") {
		if cfg.GithubClientID == "" {
			errs = append(errs, fmt.Sprintf("github.client.id is missing (env %s) while github.client.secret is set", envNameFromKey("github.client.id")))
		}
		if cfg.GithubClientSecret == "" {
			errs = append(errs, fmt.Sprintf("github.client.secret is missing (env %s) while github.client.id is set", envNameFromKey("github.client.secret")))
		}
	}

	// Enforce minimum: at least one full provider
	if !googleOK && !githubOK {
		errs = append(errs, "at least one OAuth provider must be fully configured: either Google (google.client.id + google.client.secret) or GitHub (github.client.id + github.client.secret)")
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}
	return nil
}

func envNameFromKey(key string) string {
	return strings.ToUpper(strings.ReplaceAll(key, ".", "_"))
}

func DebugPrintConfig() {
	cfg, _ := Load()
	b, err := yaml.Marshal(cfg)
	if err != nil {
		fmt.Println("error marshalling config:", err)
		return
	}
	fmt.Println("Loaded configuration:")
	fmt.Println(string(b))
}

func IsUIDev() bool {
	cfg, _ := Load()
	return cfg.UIDev
}

func IsDev() bool {
	cfg, _ := Load()
	return strings.EqualFold(cfg.Env, "development")
}

func IsDebug() bool {
	cfg, _ := Load()
	return cfg.Debug
}

func IsSwaggerEnabled() bool {
	cfg, _ := Load()
	return cfg.Swagger
}
