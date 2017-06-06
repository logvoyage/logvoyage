package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

const (
	MODE_DEV  = "dev" // Application will start in dev mode by default
	MODE_PROD = "prod"
	MODE_TEST = "test"
)

var (
	cfg *viper.Viper
)

func InitConfig() {
	cfg = viper.New()
	cfg.SetEnvPrefix("LV")
	cfg.SetDefault("mode", MODE_DEV)
	cfg.SetConfigType("json")
	cfg.AddConfigPath("$HOME/.logvoyage")
	cfg.AddConfigPath(".")
	cfg.AutomaticEnv()
	cfg.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Load custom config file from LV_CONFIG env variable.
	// E.g.: LV_CONFIG=test logvoyage start api
	// will try to load config.test.json
	if len(Get("config")) > 0 {
		cfg.SetConfigName(fmt.Sprintf("config.%s", Get("config")))
		err := cfg.ReadInConfig()
		if err != nil {
			fmt.Println("Error reading config:", err)
			os.Exit(-1)
		}
	}
}

// Get config value.
// Env variable LV_DB_HOST can be accessible as Get("db.host")
// See viper docs for more info.
func Get(key string) string {
	return cfg.GetString(key)
}

// IsDevMode returns true if app running in development mode
func IsDevMode() bool {
	return Get("mode") == MODE_DEV
}

// IsProdMode returns true if app running in producton mode
func IsProdMode() bool {
	return Get("mode") == MODE_PROD
}

// IsTestMode returns true if app running in testing mode
func IsTestMode() bool {
	return Get("mode") == MODE_TEST
}
