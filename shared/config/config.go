package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

var config = viper.New()

// Get config value.
// Env variable LV_DB_HOST will be accessible as Get("db.host")
// See viper docs for more info.
func Get(key string) string {
	return config.GetString(key)
}

func init() {
	config.SetEnvPrefix("LV")
	config.SetConfigType("json")
	config.AddConfigPath("$HOME/.logvoyage")
	config.AddConfigPath(".")
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	config.AutomaticEnv()

	err := config.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	fmt.Println("DB:", Get("db.address"))
	fmt.Println("DB:", Get("db.port"))
}
