package config

import "github.com/spf13/viper"

func initDefaults() {
	viper.SetDefault("api.title", "Veribot API")
	viper.SetDefault("api.description", "API to authenticate/authorize discord users")
	viper.SetDefault("api.version", "0.1")
	viper.SetDefault("api.path", "/api/v1")
	viper.SetDefault("api.hostname", "whodis.netsoc.cloud")

	viper.SetDefault("discord.token", "")
}
