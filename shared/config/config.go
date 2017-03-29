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
	// Main config
	mainConfig *viper.Viper
	// Config for modes.
	// E.g. if LV_MODE="dev" config.dev.json will be loaded.
	modeConfig *viper.Viper
)

func init() {
	m, err := getViperInstance("config")
	if err != nil {
		fmt.Println("Error reading main config file:", err)
		os.Exit(1)
	}
	mainConfig = m
	modeConfig, _ = getViperInstance(fmt.Sprintf("config.%s", Get("mode")))
}

func getViperInstance(configName string) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigName(configName)
	v.SetDefault("mode", MODE_DEV)
	v.SetEnvPrefix("LV")
	v.SetConfigType("json")
	v.AddConfigPath("$HOME/.logvoyage")
	v.AddConfigPath(".")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}
	return v, nil
}

// Get config value.
// Env variable LV_DB_HOST can be accessible as Get("db.host")
// See viper docs for more info.
func Get(key string) string {
	if modeConfig != nil && modeConfig.GetString(key) != "" {
		return modeConfig.GetString(key)
	}
	return mainConfig.GetString(key)
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
