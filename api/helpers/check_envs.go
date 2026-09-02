package helpers

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// EnvConfig holds the environment configuration values used for validation.
type EnvConfig struct {
	DBHost      string `mapstructure:"DB_HOST"`
	DBPort      string `mapstructure:"DB_PORT"`
	DBUser      string `mapstructure:"DB_USER"`
	DBPassword  string `mapstructure:"DB_PASSWORD"`
	DBName      string `mapstructure:"DB_NAME"`
	ENV         string `mapstructure:"ENV"`
	FrontendURL string `mapstructure:"FRONTEND_URL"`
}

// envOr returns the value of the environment variable named by key,
// falling back to the viper config value if the env var is empty.
func envOr(key, viperFallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return viperFallback
}

// CheckEnvs loads configuration from environment variables (primary) and
// the .env file (fallback for local development), then validates that all
// required variables are defined correctly.
func CheckEnvs() error {
	viper.SetConfigFile(".env")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("ENV", "development")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("error loading .env file: %w", err)
		}
		log.Println("No .env file found, using environment variables.")
	}

	config := EnvConfig{
		DBHost:      envOr("DB_HOST", viper.GetString("DB_HOST")),
		DBPort:      envOr("DB_PORT", viper.GetString("DB_PORT")),
		DBUser:      envOr("DB_USER", viper.GetString("DB_USER")),
		DBPassword:  envOr("DB_PASSWORD", viper.GetString("DB_PASSWORD")),
		DBName:      envOr("DB_NAME", viper.GetString("DB_NAME")),
		ENV:         envOr("ENV", viper.GetString("ENV")),
		FrontendURL: envOr("FRONTEND_URL", viper.GetString("FRONTEND_URL")),
	}

	return validate(&config)
}

func validate(cfg *EnvConfig) error {
	required := map[string]string{
		"DB_HOST":     cfg.DBHost,
		"DB_USER":     cfg.DBUser,
		"DB_PASSWORD": cfg.DBPassword,
		"DB_NAME":     cfg.DBName,
	}

	var missing []string
	for name, value := range required {
		if strings.TrimSpace(value) == "" {
			missing = append(missing, name)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing or empty environment variables: %s", strings.Join(missing, ", "))
	}

	if cfg.ENV != "development" && cfg.ENV != "production" {
		return fmt.Errorf("invalid ENV value %q: must be \"development\" or \"production\"", cfg.ENV)
	}

	if cfg.ENV == "production" && strings.TrimSpace(cfg.FrontendURL) == "" {
		return fmt.Errorf("missing or empty environment variable: FRONTEND_URL (required when ENV=production)")
	}

	return nil
}
