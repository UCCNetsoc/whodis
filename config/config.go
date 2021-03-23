package config

import (
	"strings"

	"github.com/spf13/viper"
)

func InitConfig() {
	initDefaults()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
}
