package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

var config = viper.New()

// const (
// 	MODE_DEV  = "dev"
// 	MODE_PROD = "prod"
// 	MODE_TEST = "test"
// )

// Get config value.
// Env variable LV_DB_HOST can be accessible as Get("db.host")
// See viper docs for more info.
func Get(key string) string {
	return config.GetString(key)
}

func init() {
	config.SetEnvPrefix("LV")
	config.AutomaticEnv()
	config.SetConfigType("json")
	config.AddConfigPath("$HOME/.logvoyage")
	config.AddConfigPath(".")
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := config.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
}
